package humanize

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// File is the annotations in package and all its sub types
type File struct {
	FileName    string
	PackageName string
	Docs        Docs
	Functions   []*Function
	Imports     []*Import
	Variables   []*Variable
	Constants   []*Constant
	Types       []*TypeName
}

type walker struct {
	src     string
	File    *File
	Package *Package
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
		switch t := node.(type) {
		case *ast.File:
			fv.File.PackageName = nameFromIdent(t.Name)
			fv.File.Docs = docsFromNodeDoc(t.Doc)
		case *ast.FuncDecl:
			fv.File.Functions = append(fv.File.Functions, NewFunction(t, fv.src, fv.File, fv.Package))
			return nil // Do not go deeper
		case *ast.GenDecl:
			// Constants :/
			for i := range t.Specs {
				switch decl := t.Specs[i].(type) {
				case *ast.ImportSpec:
					fv.File.Imports = append(fv.File.Imports, NewImport(decl, t.Doc))
				case *ast.ValueSpec:
					if t.Tok.String() == "var" {
						fv.File.Variables = append(fv.File.Variables, NewVariable(decl, t.Doc, fv.src, fv.File, fv.Package)...)
					} else if t.Tok.String() == "const" {
						fv.File.Constants = append(fv.File.Constants, NewConstant(decl, t.Doc, fv.src, fv.File, fv.Package)...)
					}
				case *ast.TypeSpec:
					fv.File.Types = append(fv.File.Types, NewType(decl, t.Doc, fv.src, fv.File, fv.Package))
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
func ParseFile(src string, p *Package) (*File, error) {
	lastConst = nil
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return &File{}, err
	}

	fv := &walker{}
	fv.src = src
	fv.File = &File{}
	fv.Package = p

	ast.Walk(fv, f)

	return fv.File, nil
}
