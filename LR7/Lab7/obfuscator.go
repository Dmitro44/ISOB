package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"math/rand"
	"os"
	"strconv"
)

type Obfuscator struct {
	mapping map[string]string
	fset    *token.FileSet
	key     byte
}

func NewObfuscator() *Obfuscator {
	return &Obfuscator{
		mapping: make(map[string]string),
		fset:    token.NewFileSet(),
		key:     0xAA,
	}
}

func (o *Obfuscator) generateNewName(oldName string) string {
	if newName, ok := o.mapping[oldName]; ok {
		return newName
	}

	noDigits := []rune("IlO")
	chars := []rune("10IlO")
	res := make([]rune, 12)
	res[0] = noDigits[rand.Intn(len(noDigits))]

	for i := 1; i < len(res); i++ {
		res[i] = chars[rand.Intn(len(chars))]
	}
	o.mapping[oldName] = string(res)
	return string(res)
}

func (o *Obfuscator) collectIdentifiers(node *ast.File) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.StructType:
			for _, field := range x.Fields.List {
				for _, name := range field.Names {
					o.generateNewName(name.Name)
				}
			}
		case *ast.FuncDecl:
			if x.Name.Name != "main" {
				o.generateNewName(x.Name.Name)
			}
		}
		return true
	})
}

func (o *Obfuscator) createDecryptCall(val string) *ast.CallExpr {
	elements := []ast.Expr{}
	for j := 0; j < len(val); j++ {
		elements = append(elements, &ast.BasicLit{
			Kind:  token.INT,
			Value: fmt.Sprintf("0x%02x", val[j]^o.key),
		})
	}
	return &ast.CallExpr{
		Fun: &ast.Ident{Name: o.generateNewName("decode")},
		Args: []ast.Expr{
			&ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: &ast.Ident{Name: "byte"},
				},
				Elts: elements,
			},
		},
	}
}

func (o *Obfuscator) obfuscateIdentifiersAndStrings(node *ast.File) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if newName, ok := o.mapping[x.Name]; ok {
				x.Name = newName
				return true
			}
			if x.Name != "main" && x.Obj != nil && x.Name != "d" {
				x.Name = o.generateNewName(x.Name)
			}

		case *ast.KeyValueExpr:
			if key, ok := x.Key.(*ast.Ident); ok {
				if newName, ok := o.mapping[key.Name]; ok {
					key.Name = newName
				}
			}

		case *ast.CallExpr:
			for i, arg := range x.Args {
				if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					val, err := strconv.Unquote(lit.Value)
					if err != nil || len(val) == 0 {
						continue
					}
					x.Args[i] = o.createDecryptCall(val)
				}
			}
		}
		return true
	})
}

func (o *Obfuscator) injectDecryptFunc(node *ast.File) {
	decryptFunc := &ast.FuncDecl{
		Name: &ast.Ident{Name: o.generateNewName("decode")},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{{
					Names: []*ast.Ident{{Name: o.generateNewName("b")}},
					Type:  &ast.ArrayType{Elt: &ast.Ident{Name: "byte"}},
				}},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{{Type: &ast.Ident{Name: "string"}}},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.RangeStmt{
					Key: &ast.Ident{Name: o.generateNewName("i")},
					Tok: token.DEFINE,
					X:   &ast.Ident{Name: o.generateNewName("b")},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: o.generateNewName("b")}, Index: &ast.Ident{Name: o.generateNewName("i")}}},
								Tok: token.XOR_ASSIGN,
								Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("0x%02x", o.key)}},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{Fun: &ast.Ident{Name: "string"}, Args: []ast.Expr{&ast.Ident{Name: o.generateNewName("b")}}},
					},
				},
			},
		},
	}
	node.Decls = append(node.Decls, decryptFunc)
}

func (o *Obfuscator) Obfuscate(inputPath, outputPath string) error {
	node, err := parser.ParseFile(o.fset, inputPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	o.collectIdentifiers(node)
	o.obfuscateIdentifiersAndStrings(node)
	o.injectDecryptFunc(node)

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return printer.Fprint(f, o.fset, node)
}

func main() {
	obfs := NewObfuscator()
	err := obfs.Obfuscate("lab_to_obfuscate/main.go", "lab_to_obfuscate/main_obfuscated.go")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Obfuscation complete!")
}
