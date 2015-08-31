package annotate

import (
	"encoding/json"
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

// Docs is use to store documents
type Docs []string

// File is the annotations in package and all its sub types
type File struct {
	Name       string
	Annotation []Annotate
	Docs       Docs
}

type walker struct {
	src       string
	File      File
	Functions []Function
	Imports   []Import
	Variables []Variable
	Types     []TypeName
}

func nameFromIdent(i *ast.Ident) (name string) {
	if i != nil {
		name = i.String()
	}
	return
}

func docsFromNodeDoc(cgs ...*ast.CommentGroup) Docs {
	var res Docs
	for _, cg := range cgs {
		if cg != nil {
			for i := range cg.List {
				res = append(res, cg.List[i].Text)
			}
		}
	}
	return res
}

func (fv *walker) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		//fmt.Printf("\n%T\n", node)
		switch t := node.(type) {
		case *ast.File:
			fv.File.Name = nameFromIdent(t.Name)
			fv.File.Docs = docsFromNodeDoc(t.Doc)
		case *ast.FuncDecl:
			fv.Functions = append(fv.Functions, NewFunction(t, fv.src))
			return nil // Do not go deeper
		case *ast.GenDecl:
			for i := range t.Specs {
				switch decl := t.Specs[i].(type) {
				case *ast.ImportSpec:
					fv.Imports = append(fv.Imports, NewImport(decl))
				case *ast.ValueSpec:
					fv.Variables = append(fv.Variables, NewVariable(decl, t.Doc, fv.src)...)
				case *ast.TypeSpec:
					fv.Types = append(fv.Types, NewType(decl, t.Doc, fv.src))
				}
			}
			return nil
		default:
			//fmt.Printf("\n%T\n", t)
		}
	}
	return fv
}

// ParseFile try to parse a single file for its annotations
func ParseFile(src string) (*File, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	fv := &walker{}
	fv.src = src

	ast.Walk(fv, f)

	//fmt.Printf("%+v", fv.types)
	d, _ := json.MarshalIndent(fv, "", "\t")
	fmt.Print(string(d))

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
