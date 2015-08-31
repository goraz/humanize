package annotate

import "go/ast"

// Function is annotation data for a single function
type Function struct {
	Name       string
	Reciever   *Variable // Nil means function
	Docs       Docs
	Parameters []Variable
	Results    []Variable
	Annotates  []Annotate
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

	if f.Type.Results != nil {
		for i := range f.Type.Results.List {
			n := f.Type.Results.List[i].Names
			for in := range n {
				p := variableFromExpr(nameFromIdent(n[in]), f.Type.Results.List[i].Type, src)
				res.Results = append(res.Results, p)
			}
		}
	}

	for i := range f.Type.Params.List {
		n := f.Type.Params.List[i].Names
		for in := range n {
			p := variableFromExpr(nameFromIdent(n[in]), f.Type.Params.List[i].Type, src)
			res.Parameters = append(res.Parameters, p)
		}
	}

	return res
}
