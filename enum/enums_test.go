package enum

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_generate_search(t *testing.T) {
	tests := []struct {
		name        string
		wantErr     bool
		checkResult func(gn *generate) error
	}{
		{
			name:    "get all go files in this project",
			wantErr: false,
			checkResult: func(gn *generate) error {
				if len(gn.pathsWithConsts) == 0 {
					return fmt.Errorf("pathsWithConsts should not be len 0")
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &generate{pathsWithConsts: make(pathsWithConsts)}
			if err := e.search(); (err != nil) != tt.wantErr {
				t.Errorf("search() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResult(e); err != nil {
				t.Errorf("checkResult function error: %v", err)
			}
		})
	}
}

func Test_generate_retrieveTypes(t *testing.T) {
	tests := []struct {
		name                    string
		wantErr                 bool
		generatePathsWithConsts func() pathsWithConsts
		checkResult             func(gn *generate) error
	}{
		{
			name:    "success retrieving types from _tests folder",
			wantErr: false,
			generatePathsWithConsts: func() pathsWithConsts {
				absPath, _ := filepath.Abs(".")
				absPath = filepath.Join(
					absPath,
					"..",
					"_tests",
					"simple",
					"test1.go",
				)
				return pathsWithConsts{
					absPath: make(constants, 0),
				}
			},
			checkResult: func(gn *generate) error {
				first := gn.pathsWithConsts.GetOneOrNil()
				if first == nil {
					return fmt.Errorf("checkResult, at least i should have one element")
				}

				if first.GetPackageName() == "main" {
					return fmt.Errorf("checkResult, package should not be main")
				}

				type testResults struct {
					testType      string
					testValueType string
					checked       bool
				}

				results := map[string]*testResults{
					"car": {
						testType:      "car",
						testValueType: "int",
					},
					"color": {
						testType:      "color",
						testValueType: "string",
					},
				}
				for _, consts := range gn.pathsWithConsts {
					for _, c := range consts {
						v, ok := results[c.Value()]
						if !ok {
							continue
						}
						if v.testType == c.Value() && v.testValueType == c.TypeReturnValue() {
							v.checked = true
						}
					}
				}
				for i := range results {
					if !results[i].checked {
						return fmt.Errorf("the %+v is not checked as true", results[i])
					}
				}

				return nil
			},
		},
		{
			name: "fail, no file is created",
			generatePathsWithConsts: func() pathsWithConsts {
				return pathsWithConsts{}
			},
			checkResult: func(gn *generate) error {
				first := gn.pathsWithConsts.GetOneOrNil()
				if first != nil {
					return fmt.Errorf("checkResult, there should be no element")
				}

				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &generate{
				pathsWithConsts: tt.generatePathsWithConsts(),
			}
			if err := e.retrieveTypes(); (err != nil) != tt.wantErr {
				t.Errorf("retrieveTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResult(e); err != nil {
				t.Errorf("checkResult function error: %v", err)
			}
		})
	}
}

func Test_generate_generateFromTemplate(t *testing.T) {
	tests := []struct {
		name                    string
		generatePathsWithConsts func() pathsWithConsts
		checkResult             func(gn *generate) error
		wantErr                 bool
	}{
		{
			name: "success",
			generatePathsWithConsts: func() pathsWithConsts {
				absPath, _ := filepath.Abs(".")
				absPath = filepath.Join(
					absPath,
					"..",
					"_tests",
					"simple",
					"test1.go",
				)
				return pathsWithConsts{
					absPath: constants{
						{
							p: "simple",
							t: "color",
							v: "string",
						},
						{
							p: "simple",
							t: "car",
							v: "int",
						},
					},
				}
			},
			checkResult: func(gn *generate) error {
				absPathRoot, _ := filepath.Abs(".")
				absPath := filepath.Join(
					absPathRoot,
					"..",
					"_tests",
					"simple",
					"result_gen.go.txt",
				)
				f, err := os.Open(absPath)
				if err != nil {
					return err
				}
				b1, err := io.ReadAll(f)
				if err != nil {
					return err
				}
				_ = f.Close()
				absPath = filepath.Join(
					absPathRoot,
					"..",
					"_tests",
					"simple",
					"test1_gen.go",
				)
				f, err = os.Open(absPath)
				if err != nil {
					return err
				}
				b2, err := io.ReadAll(f)
				if err != nil {
					return err
				}
				_ = f.Close()
				if !reflect.DeepEqual(b1, b2) {
					return fmt.Errorf("result != generated:\n--------\n%v\n------\n%v\n--------", string(b1), string(b2))
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &generate{
				pathsWithConsts: tt.generatePathsWithConsts(),
			}
			if err := e.generateFromTemplate(); (err != nil) != tt.wantErr {
				t.Errorf("generateFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResult(e); err != nil {
				t.Errorf("checkResult function error: %v", err)
			}
		})
	}
}

func Test_generate_Generate(t *testing.T) {
	tests := []struct {
		name        string
		checkResult func(gn Generator) error
		wantErr     bool
	}{
		{
			name: "success",
			checkResult: func(gn Generator) error {
				absPathRoot, _ := filepath.Abs(".")
				absPath := filepath.Join(
					absPathRoot,
					"..",
					"_tests",
					"simple",
					"result_gen.go.txt",
				)
				f, err := os.Open(absPath)
				if err != nil {
					return err
				}
				b1, err := io.ReadAll(f)
				if err != nil {
					return err
				}
				_ = f.Close()
				absPath = filepath.Join(
					absPathRoot,
					"..",
					"_tests",
					"simple",
					"test1_gen.go",
				)
				f, err = os.Open(absPath)
				if err != nil {
					return err
				}
				b2, err := io.ReadAll(f)
				if err != nil {
					return err
				}
				_ = f.Close()
				if !reflect.DeepEqual(b1, b2) {
					return fmt.Errorf("result != generated:\n--------\n%v\n------\n%v\n--------", string(b1), string(b2))
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewGenerator()
			if err := e.Generate(); (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResult(e); err != nil {
				t.Errorf("checkResult with err: %v", err)
			}
		})
	}
}
