package annotate

import (
	"go/ast"
	"go/token"
)

// Variable is a string represent of a function parameter
type Variable struct {
	Name string
	Type Type
	Docs Docs

	caller *ast.CallExpr
	indx   int
}

func variableFromValue(name string, indx int, e []ast.Expr, src string) Variable {
	var t Type
	var caller *ast.CallExpr
	var ok bool
	first := e[0]
	if caller, ok = first.(*ast.CallExpr); !ok {
		switch data := e[indx].(type) {
		case *ast.CompositeLit:
			//if data.Type != nil {
			// the type is here
			t = getType(data.Type, src)
			//}
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
			//default:
			//fmt.Printf("var value => %T", data)
			//fmt.Printf("%s", src[data.Pos()-1:data.End()-1])
		}
	}
	return Variable{
		Name:   name,
		Type:   t,
		caller: caller,
		indx:   indx,
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
			if len(v.Values) != 0 {
				n = variableFromValue(name, i, v.Values, src)
			}
		}
		n.Docs = docsFromNodeDoc(c, v.Doc)
		res = append(res, n)
	}

	return res
}
