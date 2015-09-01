package annotate

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Package is list of files
type Package []File

var (
	packageCache map[string]Package
)

// FindType return a base type interface base on the string name of the type
func (p Package) FindType(t string) (*TypeName, error) {
	for i := range p {
		for j := range p[i].Types {
			if p[i].Types[j].Name == t {
				return &p[i].Types[j], nil
			}
		}
	}

	return nil, fmt.Errorf("type with name %s not found", t)
}

// FindVariable try to find a package level variable
func (p Package) FindVariable(t string) (*Variable, error) {
	for i := range p {
		for j := range p[i].Variables {
			if p[i].Variables[j].Name == t {
				return &p[i].Variables[j], nil
			}
		}
	}

	return nil, fmt.Errorf("var with name %s not found", t)
}

// FindFunction try to find a package level variable
func (p Package) FindFunction(t string) (*Function, error) {
	for i := range p {
		for j := range p[i].Functions {
			if p[i].Functions[j].Name == t {
				return &p[i].Functions[j], nil
			}
		}
	}

	return nil, fmt.Errorf("func with name %s not found", t)
}

// FindImport try to find an import by its full import path
func (p Package) FindImport(t string) (*Import, error) {
	if t == "" || t == "_" {
		return nil, fmt.Errorf("import with path _ or empty is invalid")
	}
	for i := range p {
		for j := range p[i].Imports {
			if p[i].Imports[j].Name == t || p[i].Imports[j].Path == t {
				return &p[i].Imports[j], nil
			}
		}
	}

	return nil, fmt.Errorf("import with name or path %s not found", t)
}

func translateToFullPath(path string) (string, error) {
	root := runtime.GOROOT()
	gopath := strings.Split(os.Getenv("GOPATH"), ":")

	test := filepath.Join(root, "src", path)
	r, err := os.Stat(test)
	if err != nil {
		for i := range gopath {
			test = filepath.Join(gopath[i], "src", path)
			r, err = os.Stat(test)
			if err == nil {
				break
			}
		}
		if err != nil {
			return "", fmt.Errorf("%s is not found in GOROOT or GOPATH", path)
		}
	}

	if !r.IsDir() {
		return "", fmt.Errorf("%s is found in %s but its not a directory", path, r.Name())
	}

	return test, nil
}

func lateBind(p Package) error {
	for f := range p {
		// Try to find variable with null type and change them to real type
		for v := range p[f].Variables {
			if p[f].Variables[v].caller != nil {
				switch c := p[f].Variables[v].caller.Fun.(type) {
				case *ast.Ident:
					name := nameFromIdent(c)
					fn, err := p.FindFunction(name)
					if err != nil {
						// Todo is this possible?
						return err
					}
					if len(fn.Results) <= p[f].Variables[v].indx {
						return fmt.Errorf("%d result is available but want the %d", len(fn.Results), p[f].Variables[v].indx)
					}
					p[f].Variables[v].Type = fn.Results[p[f].Variables[v].indx].Type
				case *ast.SelectorExpr:
					pkg := nameFromIdent(c.X.(*ast.Ident))
					typ := nameFromIdent(c.Sel)
					imprt, err := p.FindImport(pkg)
					if err != nil {
						return err
					}
					pkgDef, err := ParsePackage(imprt.Path)
					if err != nil {
						return err
					}
					fn, err := pkgDef.FindFunction(typ)
					if err != nil {
						return err
					}

					if len(fn.Results) <= p[f].Variables[v].indx {
						return fmt.Errorf("%d result is available but want the %d", len(fn.Results), p[f].Variables[v].indx)
					}
					foreignTyp := fn.Results[p[f].Variables[v].indx].Type
					if _, ok := foreignTyp.(SelectorType); ok {
						// the type is foreign to that package too
					}
					p[f].Variables[v].Type = SelectorType{
						srcBase: srcBase{""}, // TODO : source?
						Package: imprt.Name,
						Type:    foreignTyp,
					}
				}
			}
		}
	}
	return nil
}

// ParsePackage is here for loading a single package and parse all files in it
func ParsePackage(path string) (Package, error) {
	var p Package
	var ok bool
	if p, ok = packageCache[path]; ok {
		return p, nil
	}
	folder, err := translateToFullPath(path)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(
		folder,
		func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				if filepath.Ext(path) == ".go" {
					r, err := os.Open(path)
					if err != nil {
						return err
					}
					data, err := ioutil.ReadAll(r)
					if err != nil {
						return err
					}

					f, err := ParseFile(string(data))
					if err != nil {
						return err
					}
					f.FileName = path
					p = append(p, f)
				}
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	packageCache[path] = p

	err = lateBind(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func init() {
	packageCache = make(map[string]Package)
}
