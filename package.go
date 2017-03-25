package humanize

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// pkg is list of files
type Package struct {
	Files []*File
	Path  string
	Name  string

	resolved bool
}

var (
	packageCache = make(map[string]*Package)
	lock         = sync.RWMutex{}
)

func setCache(path string, p *Package) {
	lock.Lock()
	defer lock.Unlock()

	packageCache[path] = p
}

func getCache(path string) *Package {
	lock.RLock()
	defer lock.RUnlock()

	return packageCache[path]
}

// FindType return a base type interface base on the string name of the type
func (p Package) FindType(t string) (*TypeName, error) {
	for i := range p.Files {
		for j := range p.Files[i].Types {
			if p.Files[i].Types[j].Name == t {
				return p.Files[i].Types[j], nil
			}
		}
	}

	return nil, fmt.Errorf("type with name %s not found", t)
}

// FindVariable try to find a package level variable
func (p Package) FindVariable(t string) (*Variable, error) {
	for i := range p.Files {
		for j := range p.Files[i].Variables {
			if p.Files[i].Variables[j].Name == t {
				return p.Files[i].Variables[j], nil
			}
		}
	}

	return nil, fmt.Errorf("var with name %s not found", t)
}

// FindConstant try to find a package level variable
func (p Package) FindConstant(t string) (*Constant, error) {
	for i := range p.Files {
		for j := range p.Files[i].Constants {
			if p.Files[i].Constants[j].Name == t {
				return p.Files[i].Constants[j], nil
			}
		}
	}

	return nil, fmt.Errorf("const with name %s not found", t)
}

// FindFunction try to find a package level variable
func (p Package) FindFunction(t string) (*Function, error) {
	for i := range p.Files {
		for j := range p.Files[i].Functions {
			if p.Files[i].Functions[j].Name == t {
				return p.Files[i].Functions[j], nil
			}
		}
	}
	return nil, fmt.Errorf("func with name %s not found", t)
}

// FindImport try to find an import by its full import path
func (p Package) FindImport(t string) (*Import, error) {
	if t == "" || t == "_" || t == "." {
		return nil, fmt.Errorf("import with path _/. or empty is invalid")
	}
	for i := range p.Files {
		for j := range p.Files[i].Imports {
			if p.Files[i].Imports[j].Name == t || p.Files[i].Imports[j].Path == t {
				return p.Files[i].Imports[j], nil
			}
		}
	}

	return nil, fmt.Errorf("import with name or path %s not found", t)
}

func translateToFullPath(path string) (string, error) {
	root := runtime.GOROOT()
	gopath := strings.Split(os.Getenv("GOPATH"), ":")
	gopath = append([]string{root}, gopath...)
	var (
		test string
		r    os.FileInfo
		err  error
	)
	for i := range gopath {
		test = filepath.Join(gopath[i], "src", path)
		r, err = os.Stat(test)
		if err == nil {
			break
		}
		// some hacky way to handle vendoring
		test = filepath.Join(gopath[i], "src/vendor", path)
		r, err = os.Stat(test)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", fmt.Errorf("%s is not found in GOROOT or GOPATH", path)
	}

	if !r.IsDir() {
		return "", fmt.Errorf("%s is found in %s but its not a directory", path, r.Name())
	}

	return test, nil
}

func checkTypeCast(p *Package, bi *Package, args []ast.Expr, name string) (Type, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("it can not be a typecast : %s", name)
	}

	t, err := bi.FindType(name)
	if err == nil {
		return t.Type, nil
	}

	// if the type is in this package, then simply pass an ident type
	// not the actual type, since the actual type is the definition of
	// the type
	_, err = p.FindType(name)
	if err == nil {
		return &IdentType{Ident: name, srcBase: srcBase{pkg: p}}, nil
	}

	return nil, fmt.Errorf("can not find the call for %s", name)
}

func lateBind(p *Package) (res error) {
	builtin, err := ParsePackage("builtin")
	assertNil(err)

	for f := range p.Files {
		// Try to find variable with null type and change them to real type
	thebigLoop:
		for v := range p.Files[f].Variables {
			if p.Files[f].Variables[v].caller != nil {
				switch c := p.Files[f].Variables[v].caller.Fun.(type) {
				case *ast.Ident:
					name := nameFromIdent(c)
					bl, err := builtin.FindFunction(name)
					if err == nil {
						p.Files[f].Variables[v].Type = bl.Type
					} else {
						var t Type
						fn, err := p.FindFunction(name)
						if err == nil {
							if len(fn.Type.Results) <= p.Files[f].Variables[v].indx {
								return fmt.Errorf("%d result is available but want the %d", len(fn.Type.Results), p.Files[f].Variables[v].indx)
							}
							t = fn.Type.Results[p.Files[f].Variables[v].indx].Type
						} else {
							t, err = checkTypeCast(p, builtin, p.Files[f].Variables[v].caller.Args, name)
							if err != nil {
								return err
							}
						}

						p.Files[f].Variables[v].Type = t
					}
				case *ast.SelectorExpr:
					var pkg string
					switch c.X.(type) {
					case *ast.Ident:
						pkg = nameFromIdent(c.X.(*ast.Ident))
					case *ast.CallExpr: // TODO : Don't know why, no time for check
						continue thebigLoop
					}

					typ := nameFromIdent(c.Sel)
					imprt, err := p.FindImport(pkg)
					if err != nil {
						// TODO : package currently is not capable of parsing build tags. so ignore this :/
						continue thebigLoop
					}
					pkgDef, err := ParsePackage(imprt.Path)
					if err != nil {
						return err
					}
					var t Type
					fn, err := pkgDef.FindFunction(typ)
					if err == nil {
						if len(fn.Type.Results) <= p.Files[f].Variables[v].indx {
							return fmt.Errorf("%d result is available but want the %d", len(fn.Type.Results), p.Files[f].Variables[v].indx)
						}
						t = fn.Type.Results[p.Files[f].Variables[v].indx].Type
					} else {
						t, err = checkTypeCast(pkgDef, builtin, p.Files[f].Variables[v].caller.Args, typ)
						if err != nil {
							return err
						}
					}

					foreignTyp := t
					star := false
					if sType, ok := foreignTyp.(*StarType); ok {
						foreignTyp = sType.Target
						star = true
					}
					switch ft := foreignTyp.(type) {
					case *IdentType:
						// this is a simple hack. if the type is begin with
						// upper case, then its type on that package, else its a global type
						name := ft.Ident
						c := name[0]
						if c >= 'A' && c <= 'Z' {
							if star {
								foreignTyp = &StarType{
									ft.srcBase,
									foreignTyp,
								}
							}
							p.Files[f].Variables[v].Type = &SelectorType{
								srcBase: srcBase{p, ""}, // TODO : source?
								pkg:     getImport(imprt.Name, p.Files[f]),
								Type:    foreignTyp,
							}
						} else {
							if star {
								foreignTyp = &StarType{
									ft.srcBase,
									foreignTyp,
								}
							}
							p.Files[f].Variables[v].Type = foreignTyp
						}

					default:
						// the type is foreign to that package too
						p.Files[f].Variables[v].Type = ft
					}
				}
			}
		}
	}
	return nil
}

func findMethods(p *Package) {
	if p.resolved {
		return
	}
	p.resolved = true
	for _, f := range p.Files {
		for _, fn := range f.Functions {
			if fn.Receiver != nil {
				t := fn.Receiver.Type
				var pointer bool
				if t2, ok := t.(*StarType); ok {
					t = t2.Target
					pointer = true
				}
				nt, err := p.FindType(t.GetDefinition())
				if err != nil {
					continue
				}
				if pointer {
					nt.StarMethods = append(nt.StarMethods, fn)
				} else {
					nt.Methods = append(nt.Methods, fn)
				}
			}
		}
	}
}

func getGoFileContent(path, folder string, f os.FileInfo) (string, error) {
	if f.IsDir() {
		if path != folder {
			return "", filepath.SkipDir
		} else {
			return "", nil
		}
	}
	if filepath.Ext(path) != ".go" {
		return "", nil
	}
	// ignore test files (for now?)
	_, filename := filepath.Split(path)
	if len(filename) > 8 && filename[len(filename)-8:] == "_test.go" {
		return "", nil
	}
	r, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ParsePackage is here for loading a single package and parse all files in it
func ParsePackage(path string) (*Package, error) {
	if p := getCache(path); p != nil {
		return p, nil
	}
	var p = &Package{}
	p.Path = path
	folder, err := translateToFullPath(path)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(
		folder,
		func(path string, f os.FileInfo, err error) error {
			data, err := getGoFileContent(path, folder, f)
			if err != nil || data == "" {
				return err
			}
			fl, err := ParseFile(string(data), p)
			if err != nil {
				return err
			}
			fl.FileName = path
			p.Files = append(p.Files, fl)
			p.Name = fl.PackageName

			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	setCache(path, p)
	err = lateBind(p)
	if err != nil {
		return nil, err
	}

	findMethods(p)
	return p, nil
}

func assertNil(e interface{}) {
	if e != nil {
		panic(e)
	}
}
