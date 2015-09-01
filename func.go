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
		if n != nil {
			for in := range n {
				p := variableFromExpr(nameFromIdent(n[in]), f.List[i].Type, src)
				res = append(res, p)
			}
		} else {
			// Its probably without name part (ie return variable)
			p := variableFromExpr("", f.List[i].Type, src)
			res = append(res, p)
		}
	}

	return res
}

// NewFunction return a single function annotation
func NewFunction(f *ast.FuncDecl, src string) Function {
	res := Function{}

	res.Name = nameFromIdent(f.Name)
	res.Docs = docsFromNodeDoc(f.Doc)

	if f.Recv != nil {
		// Method reciever is only one parameter
		for i := range f.Recv.List {
			n := ""
			if f.Recv.List[i].Names != nil {
				n = nameFromIdent(f.Recv.List[i].Names[0])
			}
			p := variableFromExpr(n, f.Recv.List[i].Type, src)
			res.Reciever = &p
		}
	}

	// Change the name of the function to reciver.func
	if res.Reciever != nil {
		tmp := res.Reciever.Type
		if _, ok := tmp.(StarType); ok {
			tmp = tmp.(StarType).Target
		}

		res.Name = tmp.(IdentType).Ident + "." + res.Name
	}

	res.Results = extractVariableList(f.Type.Results, src)
	res.Parameters = extractVariableList(f.Type.Params, src)

	return res
}
