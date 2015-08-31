package annotate

import "go/ast"

// Variable is a string represent of a function parameter
type Variable struct {
	Name string
	Type Type
	Docs Docs
}

func variableFromExpr(name string, e ast.Expr, src string) Variable {
	return Variable{
		Name: name,
		Type: getType(e, src),
	}
}

// NewVariable return an array of variables in the scope
func NewVariable(v *ast.ValueSpec, c *ast.CommentGroup, src string) []Variable {
	var res []Variable
	for i := range v.Names {
		name := ""
		if v.Names[i] != nil {
			name = v.Names[i].String()
		}
		n := variableFromExpr(name, v.Type, src)
		n.Docs = docsFromNodeDoc(c, v.Doc)
		res = append(res, n)
	}

	return res
}
