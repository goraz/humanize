package annotate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// Annotate represent a single annotation line
type Annotate struct {
	Identifier string
	Parameters []string
	Source     string
}

// FunctionAnnotate represent single function annotations
type FunctionAnnotate struct {
	Name        string
	Receiver    string `json:",omitempty"`
	Annotations []Annotate
}

// Package is the annotations in package and all its sub types
type Package struct {
	Name       string
	Annotation []Annotate
	Functions  []FunctionAnnotate
}

type walker struct {
	src       string
	functions []Function
	imports   []Import
	variables []Variable
	types     []TypeName
}

func nameFromIdent(i *ast.Ident) (name string) {
	if i != nil {
		name = i.String()
	}
	return
}

func (fv *walker) Visit(node ast.Node) ast.Visitor {
	switch t := node.(type) {
	case *ast.FuncDecl:
		fv.functions = append(fv.functions, NewFunction(t, fv.src))
		return nil // Do noy go deeper
	case *ast.GenDecl:
		for i := range t.Specs {
			switch decl := t.Specs[i].(type) {
			case *ast.ImportSpec:
				fv.imports = append(fv.imports, NewImport(decl))
			case *ast.ValueSpec:
				fv.variables = append(fv.variables, NewVariable(decl, fv.src)...)
			case *ast.TypeSpec:
				fv.types = append(fv.types, NewType(decl, fv.src))
			}

		}
	}

	return fv
}

// ParseFile try to parse a single file for its annotations
func ParseFile(src string) (*Package, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	fv := &walker{}
	fv.src = src

	ast.Walk(fv, f)

	fmt.Printf("%+v", fv.types)
	//	d, _ := json.MarshalIndent(fv.types, "", "\t")
	//	fmt.Print(string(d))

	return nil, nil
	/*
		cmap := ast.NewCommentMap(fset, f, f.Comments)

		p := Package{}

		for _, i := range f.Decls {
			switch i.(type) {
			case *ast.FuncDecl:
				f := i.(*ast.FuncDecl)
				fmt.Println(f.Name)
				fmt.Println(f.Type.Params)
				fmt.Printf("%+v", NewFunction(f, src))
			case *ast.GenDecl:
				s := i.(*ast.GenDecl)
				fmt.Printf("%+v", *s)
			default:
				fmt.Printf("HIII %T", i)
				continue

			}
			for j := range cmap[i] {
				fmt.Println(cmap[i][j].Text())
			}
		}
		//fmt.Println(cmap)
		os.Exit(0)
		return &p, nil
		//	w := walker{cmap: cmap, p: p}
		//	ast.Walk(&w, f)

		//	return &w.p, nil
	*/
}
