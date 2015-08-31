package annotate

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Variable is a string represent of a function parameter
type Variable struct {
	Name string
	Type Type
	Docs Docs
}

func variableFromValue(name string, indx int, e []ast.Expr, src string) Variable {
	var t Type
	first := e[0]
	if _, ok := first.(*ast.CallExpr); ok {
		// TODO : then it is call expr and I cannot use it, and must parse another function for that
	} else {
		switch data := e[indx].(type) {
		case *ast.CompositeLit:
			if data.Type != nil {
				// the type is here
				t = getType(data.Type, src)
			} else {
				// TODO : there is no type.
				fmt.Printf("%T", data)
			}
		case *ast.BasicLit:
			switch data.Kind {
			case token.INT:
				t = IdentType{
					srcBase{getSource(data, src)},
					"int",
				}
			case token.FLOAT:
				t = IdentType{
					srcBase{getSource(data, src)},
					"float64",
				}
			case token.IMAG:
				t = IdentType{
					srcBase{getSource(data, src)},
					"complex64",
				}
			case token.CHAR:
				t = IdentType{
					srcBase{getSource(data, src)},
					"char",
				}
			case token.STRING:
				t = IdentType{
					srcBase{getSource(data, src)},
					"string",
				}
			}
		default:
			fmt.Printf("%T", data)
		}
	}
	return Variable{
		Name: name,
		Type: t,
	}
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
		name := nameFromIdent(v.Names[i])
		var n Variable
		if v.Type != nil {
			n = variableFromExpr(name, v.Type, src)
		} else {
			n = variableFromValue(name, i, v.Values, src)
		}
		n.Docs = docsFromNodeDoc(c, v.Doc)
		res = append(res, n)
	}

	return res
}
