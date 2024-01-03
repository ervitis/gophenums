package enum

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ervitis/gotransactions"
)

const generateComment = "gophenum:generate"
const templateFile = "enum_gen.tmpl"

var templateFilePath = func() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Join(filepath.Dir(b), "..", "enum", "templates", templateFile)
}

// path -> []constant
type constant struct {
	p string // package name
	t string // type
	v string // value of type constant
}
type constants []constant
type pathsWithConsts map[string]constants
type generate struct {
	pathsWithConsts
}

func (c constant) Value() string {
	return c.t
}

func (c constant) TypeReturnValue() string {
	return c.v
}

func (c constant) Capitalize() string {
	return cases.Title(language.English).String(c.t)
}

func (cs constants) GetPackageName() string {
	if len(cs) == 0 {
		return "main"
	}
	return cs[0].p
}

func (p pathsWithConsts) First() constants {
	for _, v := range p {
		return v
	}
	return nil
}

func (p pathsWithConsts) GetOneOrNil() constants {
	for _, v := range p {
		if len(v) > 0 {
			return v
		}
	}
	return nil
}

type data struct {
	EnumData    []enumData
	PackageName string
}
type enumData struct {
	EnumIface           string
	EnumName            string
	EnumTypeValueReturn string
}

type Generator interface {
	Generate() error
}

func NewGenerator() Generator {
	return &generate{pathsWithConsts: make(pathsWithConsts)}
}

func (e *generate) search() error {
	absPath, _ := filepath.Abs(".")
	if err := filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
			if _, ok := e.pathsWithConsts[path]; ok {
				return nil
			}
			e.pathsWithConsts[path] = make(constants, 0)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("search: searching enums in files: %w", err)
	}
	return nil
}

func (e *generate) retrieveTypes() error {
	fset := token.NewFileSet()
	var commentsAssociated []string
	for fi := range e.pathsWithConsts {
		var found bool
		var data constant
		f, err := parser.ParseFile(fset, fi, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("retrieveTypes: parsing file: %w", err)
		}
		data.p = f.Name.Name
		ast.Inspect(f, func(node ast.Node) bool {
			if node == nil {
				return false
			}
			switch exp := node.(type) {
			case *ast.GenDecl:
				if len(exp.Specs) == 0 || exp.Doc == nil {
					return true
				}
				for _, commentGroup := range exp.Doc.List {
					if !strings.Contains(commentGroup.Text, generateComment) {
						return true
					}
					if tSpec, ok := exp.Specs[0].(*ast.TypeSpec); ok {
						data.t = tSpec.Name.Name
						found = true
						commentsAssociated = append(commentsAssociated, tSpec.Name.Name)
					}
				}
			case *ast.TypeSpec:
				if found &&
					exp.Name.Obj.Decl != nil {
					if tSpec, ok := exp.Name.Obj.Decl.(*ast.TypeSpec); ok {
						for _, eType := range commentsAssociated {
							if eType == tSpec.Name.Name {
								t := fmt.Sprintf("%v", tSpec.Type)
								data.v = t
								e.pathsWithConsts[fi] = append(e.pathsWithConsts[fi], data)
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil
}

func (e *generate) generateFromTemplate() error {
	tmpl, err := template.New(templateFile).ParseFiles(templateFilePath())
	if err != nil {
		return fmt.Errorf("generateFromTemplate: generating template using file %s: %w", templateFilePath(), err)
	}

	for pathConstant, consts := range e.pathsWithConsts {
		if len(consts) == 0 {
			continue
		}
		generatedFile := pathConstant[:len(pathConstant)-3] + "_gen.go"
		e.deleteGeneratedFile(generatedFile)
		tx := gotransactions.New(
			func() error {
				eData := make([]enumData, 0)
				for i := range consts {
					eData = append(eData, enumData{
						EnumIface:           consts[i].Capitalize(),
						EnumName:            consts[i].Value(),
						EnumTypeValueReturn: consts[i].TypeReturnValue(),
					})
				}
				templateData := data{
					PackageName: consts.GetPackageName(),
					EnumData:    eData,
				}

				f, err := os.Create(generatedFile)
				if err != nil {
					return fmt.Errorf("generateFromTemplate: creating gen file: %w", err)
				}
				if err := tmpl.Execute(f, templateData); err != nil {
					e.deleteGeneratedFile(generatedFile)
					_ = f.Close()
					return fmt.Errorf("generateFromTemplate: executing with template: %w", err)
				}
				return f.Close()
			},
			func() error {
				e.deleteGeneratedFile(generatedFile)
				return nil
			},
		)
		return tx.ExecuteTransaction()
	}
	return nil
}

func (e *generate) deleteGeneratedFile(generatedFile string) {
	if _, err := os.Stat(generatedFile); os.IsExist(err) {
		_ = os.Remove(generatedFile)
	}
}

func (e *generate) Generate() error {
	if err := e.search(); err != nil {
		return fmt.Errorf("[Generate] error searching: %w", err)
	}
	if err := e.retrieveTypes(); err != nil {
		return fmt.Errorf("[Generate] error retrieving types: %w", err)
	}
	if err := e.generateFromTemplate(); err != nil {
		return fmt.Errorf("[Generate] error generating template: %w", err)
	}
	return nil
}
