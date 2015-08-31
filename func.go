package annotate

import "go/ast"

// Function is annotation data for a single function
type Function struct {
	Name       string
	Reciever   *Variable // Nil means function
	Docs       Docs
	Parameters []Variable
	Results    []Variable
}

func extractVariableList(f *ast.FieldList, src string) []Variable {
	var res []Variable
	if f == nil {
		return res
	}
	for i := range f.List {
		n := f.List[i].Names
		for in := range n {
			p := variableFromExpr(nameFromIdent(n[in]), f.List[i].Type, src)
			res = append(res, p)
		}
	}

	return res
}

// NewFunction return a single function annotation
func NewFunction(f *ast.FuncDecl, src string) Function {
	res := Function{}

	res.Name = f.Name.String()
	res.Docs = docsFromNodeDoc(f.Doc)

	if f.Recv != nil {
		// Method reciever is only one parameter
		for i := range f.Recv.List {
			n := f.Recv.List[i].Names
			for in := range n {
				p := variableFromExpr(nameFromIdent(n[in]), f.Recv.List[i].Type, src)
				if res.Reciever != nil {
					panic("method with two receiever")
				}
				res.Reciever = &p
			}
		}
	}

	res.Results = extractVariableList(f.Type.Results, src)
	res.Parameters = extractVariableList(f.Type.Params, src)

	return res
}
