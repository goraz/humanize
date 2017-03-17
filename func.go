package humanize

import (
	"go/ast"
	"strings"
)

// Function is annotation data for a single function
type Function struct {
	Name     string
	Receiver *Variable // Nil means normal function
	Docs     Docs
	Type     *FuncType
}

func compareVariable(one, two []*Variable) bool {
	if len(one) != len(two) {
		return false
	}

	for i := range one {
		if one[i].Type.GetDefinition() != two[i].Type.GetDefinition() {
			return false
		}
	}

	return true
}

func removeReceiver(fn string) string {
	split := strings.Split(fn, ".")
	if len(split) == 2 {
		return split[1]
	}

	return fn
}

func compareFunc(one, two *Function) bool {
	// the receiver is not important, since we want to check interface match

	if removeReceiver(one.Name) != removeReceiver(two.Name) {
		return false
	}

	if !compareVariable(one.Type.Parameters, two.Type.Parameters) {
		return false
	}

	if !compareVariable(one.Type.Results, two.Type.Results) {
		return false
	}

	return true
}

func compare(one, two []*Function) bool {
bigLoop:
	for i := range one {
		for j := range two {
			if compareFunc(one[i], two[j]) {
				continue bigLoop
			}
		}
		return false
	}

	return true
}

func extractVariableList(f *ast.FieldList, src string, fl *File, p *Package) []*Variable {
	var res []*Variable
	if f == nil {
		return res
	}
	for i := range f.List {
		n := f.List[i].Names
		if n != nil {
			for in := range n {
				p := variableFromExpr(nameFromIdent(n[in]), f.List[i].Type, src, fl, p)
				res = append(res, p)
			}
		} else {
			// Its probably without name part (ie return variable)
			p := variableFromExpr("", f.List[i].Type, src, fl, p)
			res = append(res, p)
		}
	}

	return res
}

// NewFunction return a single function annotation
func NewFunction(f *ast.FuncDecl, src string, fl *File, p *Package) *Function {
	res := &Function{}

	res.Name = nameFromIdent(f.Name)
	res.Docs = docsFromNodeDoc(f.Doc)

	if f.Recv != nil {
		// Method receiver is only one parameter
		for i := range f.Recv.List {
			n := ""
			if f.Recv.List[i].Names != nil {
				n = nameFromIdent(f.Recv.List[i].Names[0])
			}
			p := variableFromExpr(n, f.Recv.List[i].Type, src, fl, p)
			res.Receiver = p
		}
	}

	// Change the name of the function to receiver.func
	if res.Receiver != nil {
		tmp := res.Receiver.Type
		if _, ok := tmp.(*StarType); ok {
			tmp = tmp.(*StarType).Target
		}

		res.Name = tmp.(*IdentType).Ident + "." + res.Name
	}

	res.Type = &FuncType{
		srcBase:    srcBase{p, ""},
		Parameters: extractVariableList(f.Type.Params, src, fl, p),
		Results:    extractVariableList(f.Type.Results, src, fl, p),
	}

	return res
}
