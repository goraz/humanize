package annotate

import (
	"go/ast"
	"go/token"
)

var (
	lastConst Type
)

// Constant is a string represent of a function parameter
type Constant struct {
	Name string
	Type Type
	Docs Docs

	caller *ast.CallExpr
	indx   int
}

func constantFromValue(name string, indx int, e []ast.Expr, src string) Constant {
	var t Type
	var caller *ast.CallExpr
	var ok bool
	if len(e) == 0 {
		return Constant{
			Name: name,
		}
	}
	first := e[0]
	if caller, ok = first.(*ast.CallExpr); !ok {
		switch data := e[indx].(type) {
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
		case *ast.Ident:
			t = IdentType{
				srcBase{getSource(data, src)},
				nameFromIdent(data),
			}
			//		default:
			//			fmt.Printf("\nvar value => %T", data)
			//			fmt.Printf("\n%s", src[data.Pos()-1:data.End()-1])
		}
	}
	return Constant{
		Name:   name,
		Type:   t,
		caller: caller,
		indx:   indx,
	}
}
func constantFromExpr(name string, e ast.Expr, src string) Constant {
	return Constant{
		Name: name,
		Type: getType(e, src),
	}
}

// NewConstant return an array of constant in the scope
func NewConstant(v *ast.ValueSpec, c *ast.CommentGroup, src string) []Constant {
	var res []Constant
	for i := range v.Names {
		name := nameFromIdent(v.Names[i])
		var n Constant
		n.Name = name
		if v.Type != nil {
			n = constantFromExpr(name, v.Type, src)
		} else {
			n = constantFromValue(name, i, v.Values, src)
		}
		if n.Type == nil {
			n.Type = lastConst
		} else {
			lastConst = n.Type
		}
		n.Docs = docsFromNodeDoc(c, v.Doc)
		res = append(res, n)
	}

	return res
}
