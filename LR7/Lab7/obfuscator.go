package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"math/rand"
	"os"
)

type Obfuscator struct {
	mapping map[string]string
	fset    *token.FileSet
}

func NewObfuscator() *Obfuscator {
	return &Obfuscator{
		mapping: make(map[string]string),
		fset:    token.NewFileSet(),
	}
}

func (o *Obfuscator) generateName(oldName string) string {
	chars := []rune("01IlO")
	newName := make([]rune, 12)

	if newName, ok := o.mapping[oldName]; ok {
		return newName
	}

	for i := range newName {
		if i == 0 {
			newName[i] = chars[2+rand.Intn(len(chars)-3)]
			continue
		}
		newName[i] = chars[rand.Intn(len(chars))]
	}
	o.mapping[oldName] = string(newName)
	return string(newName)
}

func (o *Obfuscator) Obfuscate(inputPath, outputPath string) error {
	node, err := parser.ParseFile(o.fset, inputPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("unable to parse file: %v", err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.StructType:
			for _, field := range x.Fields.List {
				for _, name := range field.Names {
					o.generateName(name.Name)
				}
			}
		case *ast.FuncDecl:
			if x.Name.Name != "main" {
				o.generateName(x.Name.Name)
			}
		}
		return true
	})

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if newName, ok := o.mapping[x.Name]; ok {
				x.Name = newName
				return true
			}
			if x.Name != "main" && x.Obj != nil {
				x.Name = o.generateName(x.Name)
			}

		case *ast.KeyValueExpr:
			if key, ok := x.Key.(*ast.Ident); ok {
				if newName, ok := o.mapping[key.Name]; ok {
					key.Name = newName
				}
			}
		}
		return true
	})

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return printer.Fprint(f, o.fset, node)
}

func main() {
	obs := NewObfuscator()
	err := obs.Obfuscate("lab_to_obfuscate/main.go", "lab_to_obfuscate/main_obfuscated.go")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Obfuscation complete! Check lab_to_obfuscate/main_obfuscated.go")
}
