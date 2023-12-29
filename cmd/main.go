package main

import (
	"github.com/ervitis/gophenums/enum"
	"log"
)

func main() {
	enumGenerator := enum.NewGenerator()
	log.Println(enumGenerator.Generate())
}
