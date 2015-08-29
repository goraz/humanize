package annotate

import (
	"go/ast"
	"strings"
)

// Import is a single import entity
type Import struct {
	Name string
	Path string
}

// NewImport extract a new import entry
func NewImport(i *ast.ImportSpec) Import {
	res := Import{"", strings.Trim(i.Path.Value, `"`)}
	if i.Name != nil {
		res.Name = i.Name.String()
	}
	return res
}
