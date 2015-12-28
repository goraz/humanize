package humanize

import (
	"go/ast"
	"path/filepath"
	"strings"
)

// Import is a single import entity
type Import struct {
	Name string
	Path string
	Docs Docs
}

// NewImport extract a new import entry
func NewImport(i *ast.ImportSpec, c *ast.CommentGroup) Import {
	res := Import{
		Name: "",
		Path: strings.Trim(i.Path.Value, `"`),
		Docs: docsFromNodeDoc(c, i.Doc),
	}
	_, res.Name = filepath.Split(res.Path)
	if i.Name != nil {
		res.Name = i.Name.String()
	}
	return res
}
