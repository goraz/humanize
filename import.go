package annotate

import (
	"go/ast"
	"strings"
)

// Import is a single import entity
type Import struct {
	Name string
	Path string
	Docs Docs
}

// NewImport extract a new import entry
func NewImport(i *ast.ImportSpec) Import {
	res := Import{
		Name: "",
		Path: strings.Trim(i.Path.Value, `"`),
		Docs: docsFromNodeDoc(i.Doc),
	}
	if i.Name != nil {
		res.Name = i.Name.String()
	}
	return res
}
