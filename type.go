package humanize

import (
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
	"strings"
)

// NOTE : src was the base on my initial work, but its not accurate and may be
// incorrect in some cases!

// Type is for handling a type definition
type Type interface {
	// GetDefinition return the definition of type
	GetDefinition() string
}

type srcBase struct {
	src string
}

// IdentType is for simple definition of a type, like int, string ...
type IdentType struct {
	srcBase
	Ident string
}

// StarType is a pointer to another type
type StarType struct {
	srcBase
	Target Type
}

// Field is a single field of a structre, a variable, with tag
type Field struct {
	Variable

	Tags reflect.StructTag
}

type Embed struct {
	Type
	Docs Docs
	Tags reflect.StructTag
}

// StructType is a struct in source code
type StructType struct {
	srcBase
	Fields []Field
	Embeds []Embed
}

// ArrayType is the base array
type ArrayType struct {
	srcBase
	Slice bool
	Len   int
	Type  Type
}

// EllipsisType is slice type but with ...type definition
type EllipsisType struct {
	*ArrayType
}

// MapType is the map type
type MapType struct {
	srcBase
	Key   Type
	Value Type
}

// InterfaceType is a single interface in system
type InterfaceType struct {
	srcBase
	Functions []Function
	Embed     []Type // IdentType or SelectorType
}

// SelectorType on my knowlege is a type from another package (I hope)
type SelectorType struct {
	srcBase
	Package string
	Type    Type
}

// ChannelType is used to handle channel type definition
type ChannelType struct {
	srcBase

	Direction ast.ChanDir
	Type      Type
}

// FuncType is for function type
type FuncType struct {
	srcBase

	Parameters []*Variable
	Results    []*Variable
}

//TypeName contain type and its name
type TypeName struct {
	Type Type
	Name string
	Docs Docs

	Methods     []*Function
	StarMethods []*Function
}

// GetDefinition return the definition of this type
func (tn TypeName) GetDefinition() string {
	return tn.Name + " " + tn.Type.GetDefinition()
}

// GetName the name of this type
func (i *IdentType) GetDefinition() string {
	return i.Ident
}

// GetName the name of this type
func (i *StarType) GetDefinition() string {
	return "*" + i.Target.GetDefinition()
}

// GetName the name of this type
func (i *ArrayType) GetDefinition() string {
	if i.Slice {
		return "[]" + i.Type.GetDefinition()
	}
	return fmt.Sprintf("[%d]%s", i.Len, i.Type.GetDefinition())
}

// GetName the name of this type
func (i *EllipsisType) GetDefinition() string {
	return fmt.Sprintf("[...]%s{}", i.Type.GetDefinition())
}

// GetName the name of this type
func (i *StructType) GetDefinition() string {
	res := "struct {\n"
	for e := range i.Embeds {
		res += "\t" + i.Embeds[e].GetDefinition() + "\n"
	}

	for f := range i.Fields {
		tags := strings.Trim(string(i.Fields[f].Tags), "`")
		if tags != "" {
			tags = "`" + tags + "`"
		}
		res += fmt.Sprintf("\t%s %s %s\n", i.Fields[f].Name, i.Fields[f].Type.GetDefinition(), tags)
	}
	return res + "}"
}

// GetName the name of this type
func (i *MapType) GetDefinition() string {
	return fmt.Sprintf("map[%s]%s", i.Key.GetDefinition(), i.Value.GetDefinition())
}

// GetName the name of this type
func (i *SelectorType) GetDefinition() string {
	return i.Package + "." + i.Type.GetDefinition()
}

// GetName the name of this type
func (i *FuncType) GetDefinition() string {
	return "func " + i.getSign()
}

// GetSign the name of this type
func (i *FuncType) getSign() string {
	var args, res []string
	for a := range i.Parameters {
		args = append(args, i.Parameters[a].Type.GetDefinition())
	}

	for a := range i.Results {
		res = append(res, i.Results[a].Type.GetDefinition())
	}

	result := "(" + strings.Join(args, ",") + ")"
	if len(res) > 1 {
		result += " (" + strings.Join(res, ",") + ")"
	} else {
		result += " " + strings.Join(res, ",")
	}

	return result
}

// GetName the name of this type
func (i *ChannelType) GetDefinition() string {
	switch i.Direction {
	case 1:
		return "chan<- " + i.Type.GetDefinition()
	case 2:
		return "<-chan " + i.Type.GetDefinition()
	default:
		return "chan " + i.Type.GetDefinition()
	}
}

// GetName the name of this type
func (i *InterfaceType) GetDefinition() string {
	res := "interface {\n"
	for e := range i.Embed {
		res += "\t" + i.Embed[e].GetDefinition() + "\n"
	}
	for f := range i.Functions {
		res += "\t" + i.Functions[f].Type.GetDefinition() + "\n"
	}
	return res + "}"
}

func getSource(e ast.Expr, src string) string {
	start := e.Pos() - 1
	end := e.End() - 1
	// grab it in source
	if len(src) >= int(end) {
		return src[start:end]
	}
	return ""
}

func getType(e ast.Expr, src string) Type {
	switch t := e.(type) {
	case *ast.Ident:
		// ident is the simplest one.
		return &IdentType{
			srcBase{getSource(e, src)},
			nameFromIdent(t),
		}
	case *ast.StarExpr:
		return &StarType{
			srcBase{getSource(e, src)},
			getType(t.X, src),
		}
	case *ast.ArrayType:
		slice := t.Len == nil
		ellipsis := false
		l := 0
		if !slice {
			var (
				ls string
			)
			switch t.Len.(type) {
			case *ast.BasicLit:
				ls = t.Len.(*ast.BasicLit).Value
			case *ast.Ellipsis:
				ls = "0"
				ellipsis = true
			}
			l, _ = strconv.Atoi(ls)
		}
		var at Type
		at = &ArrayType{
			srcBase{getSource(e, src)},
			t.Len == nil,
			l,
			getType(t.Elt, src),
		}
		if ellipsis {
			at = &EllipsisType{at.(*ArrayType)}
		}
		return at
	case *ast.MapType:
		return &MapType{
			srcBase{getSource(e, src)},
			getType(t.Key, src),
			getType(t.Value, src),
		}

	case *ast.StructType:
		res := &StructType{srcBase{getSource(e, src)}, nil, nil}
		for _, s := range t.Fields.List {
			if s.Names != nil {
				for i := range s.Names {
					v := Variable{
						Name: nameFromIdent(s.Names[i]),
						Type: getType(s.Type, src),
					}

					f := Field{
						v,
						"",
					}
					if s.Tag != nil {
						f.Tags = reflect.StructTag(s.Tag.Value)
						f.Tags = f.Tags[1 : len(f.Tags)-1]
					}
					f.Docs = docsFromNodeDoc(s.Doc)
					res.Fields = append(res.Fields, f)
				}
			} else {
				e := Embed{
					Type: getType(s.Type, src),
				}
				if s.Tag != nil {
					e.Tags = reflect.StructTag(s.Tag.Value)
					e.Tags = e.Tags[1 : len(e.Tags)-1]
				}
				e.Docs = docsFromNodeDoc(s.Doc)
				res.Embeds = append(res.Embeds, e)
			}
		}

		return res
	case *ast.InterfaceType:
		// TODO : interface may refer to itself I need more time to implement this
		iface := &InterfaceType{
			srcBase: srcBase{getSource(e, src)},
		}
		for i := range t.Methods.List {
			res := Function{}
			// The method name is mandatory and always 1
			if len(t.Methods.List[i].Names) > 0 {
				res.Name = nameFromIdent(t.Methods.List[i].Names[0])

				res.Docs = docsFromNodeDoc(t.Methods.List[i].Doc)
				typ := getType(t.Methods.List[i].Type, src)
				res.Type = typ.(*FuncType)
				iface.Functions = append(iface.Functions, res)
			} else {
				// This is the embeded interface
				embed := getType(t.Methods.List[i].Type, src)
				iface.Embed = append(iface.Embed, embed)
			}

		}
		return iface
	case *ast.ChanType:
		return &ChannelType{
			srcBase:   srcBase{getSource(e, src)},
			Direction: t.Dir,
			Type:      getType(t.Value, src),
		}
	case *ast.SelectorExpr:
		return &SelectorType{
			srcBase: srcBase{getSource(e, src)},
			Package: nameFromIdent(t.X.(*ast.Ident)),
			Type:    getType(t.Sel, src),
		}
	case *ast.FuncType:
		return &FuncType{
			srcBase:    srcBase{getSource(e, src)},
			Parameters: extractVariableList(t.Params, src),
			Results:    extractVariableList(t.Results, src),
		}
	}

	return nil
}

// NewType handle a type
func NewType(t *ast.TypeSpec, c *ast.CommentGroup, src string) *TypeName {
	doc := docsFromNodeDoc(c, t.Doc)
	return &TypeName{
		Docs: doc,
		Type: getType(t.Type, src),
		Name: nameFromIdent(t.Name),
	}
}
