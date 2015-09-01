package annotate

import (
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
)

// Type is for handling a type definition
type Type interface {
	// GetSource return the source of definition of this type
	GetSource() string
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

// StructType is a struct in source code
type StructType struct {
	srcBase
	Fields []Field
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
	ArrayType
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

	Parameters []Variable
	Results    []Variable
}

//TypeName contain type and its name
type TypeName struct {
	Type Type
	Name string
	Docs Docs
}

// GetSource of the struct
func (s srcBase) GetSource() string {
	return s.src
}

func getSource(e ast.Expr, src string) string {
	start := e.Pos() - 1
	end := e.End() - 1
	// grab it in source
	return src[start:end]
}

func getType(e ast.Expr, src string) Type {
	switch t := e.(type) {
	case *ast.Ident:
		// ident is the simplest one.
		return IdentType{
			srcBase{getSource(e, src)},
			nameFromIdent(t),
		}
	case *ast.StarExpr:
		return StarType{
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
		at = ArrayType{
			srcBase{getSource(e, src)},
			t.Len == nil,
			l,
			getType(t.Elt, src),
		}
		if ellipsis {
			at = EllipsisType{at.(ArrayType)}
		}
		return at
	case *ast.MapType:
		return MapType{
			srcBase{getSource(e, src)},
			getType(t.Key, src),
			getType(t.Value, src),
		}
	case *ast.StructType:
		res := StructType{srcBase{getSource(e, src)}, nil}
		for _, s := range t.Fields.List {
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
				res.Fields = append(res.Fields, f)
			}
		}

		return res
	case *ast.InterfaceType:
		// TODO : interface may refer to itself I need more time to implement this
		iface := InterfaceType{
			srcBase{getSource(e, src)},
			nil,
		}

		for i := range t.Methods.List {
			res := Function{}
			// The method name is mandatory and always 1
			//res.Name = nameFromIdent(t.Methods.List[i].Names)
			fmt.Printf("%+v", t.Methods.List[i].Names)

			res.Docs = docsFromNodeDoc(t.Methods.List[i].Doc)
			iface.Functions = append(iface.Functions, res)
		}
		return iface
	case *ast.ChanType:
		return ChannelType{
			srcBase:   srcBase{getSource(e, src)},
			Direction: t.Dir,
			Type:      getType(t.Value, src),
		}
	case *ast.SelectorExpr:
		return SelectorType{
			srcBase: srcBase{getSource(e, src)},
			Package: nameFromIdent(t.X.(*ast.Ident)),
			Type:    getType(t.Sel, src),
		}
	case *ast.FuncType:
		return FuncType{
			srcBase:    srcBase{getSource(e, src)},
			Parameters: extractVariableList(t.Params, src),
			Results:    extractVariableList(t.Results, src),
		}
	default:
		fmt.Printf("\n%T\n%+v\n==>", t, t)
		fmt.Print(src[t.Pos()-1 : t.End()-1])
	}

	return nil
}

// NewType handle a type
func NewType(t *ast.TypeSpec, c *ast.CommentGroup, src string) TypeName {
	doc := docsFromNodeDoc(c, t.Doc)
	return TypeName{
		Docs: doc,
		Type: getType(t.Type, src),
		Name: nameFromIdent(t.Name),
	}
}
