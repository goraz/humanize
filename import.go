package humanize

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Import is a single import entity
type Import struct {
	Name string
	Path string
	Docs Docs
}

type importWalker struct {
	pkgName string
}

func (fv *importWalker) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch t := node.(type) {
		case *ast.File:
			fv.pkgName = nameFromIdent(t.Name)
		default:
		}
	}

	return fv
}

// LoadPackage is the function to load import package
func (i Import) LoadPackage() *Package {
	// XXX : Watch this.
	pkg, _ := ParsePackage(i.Path)
	return pkg
}

func peekPackageName(pkg string) (xx string) {
	_, name := filepath.Split(pkg)
	folder, err := translateToFullPath(pkg)
	if err != nil {
		return name
	}
	fv := &importWalker{}
	err = filepath.Walk(
		folder,
		func(path string, f os.FileInfo, err error) error {
			data, err := getGoFileContent(path, folder, f)
			if err != nil || data == "" {
				return err
			}
			fset := token.NewFileSet()
			fle, err := parser.ParseFile(fset, "", data, parser.PackageClauseOnly)
			if err != nil {
				return nil // try another file?
			}

			ast.Walk(fv, fle)
			// no need to continue
			return filepath.SkipDir
		},
	)
	if fv.pkgName != "" {
		name = fv.pkgName
	}
	// can not parse it, use the folder name
	return name
}

// NewImport extract a new import entry
func NewImport(i *ast.ImportSpec, c *ast.CommentGroup) *Import {
	res := &Import{
		Name: "",
		Path: strings.Trim(i.Path.Value, `"`),
		Docs: docsFromNodeDoc(c, i.Doc),
	}
	if i.Name != nil {
		res.Name = i.Name.String()
	} else {
		res.Name = peekPackageName(res.Path)
	}
	return res
}
