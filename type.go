package annotate

import (
	"fmt"
	"go/ast"
	"reflect"
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

//TypeName contain type and its name
type TypeName struct {
	Type Type
	Name string
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
		return IdentType{srcBase{getSource(e, src)}, t.Name}
	case *ast.StarExpr:
		return StarType{srcBase{getSource(e, src)}, getType(t.X, src)}
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
	default:
		fmt.Printf("\n%T\n", t)
	}

	return nil
}

// NewType handle a type
func NewType(t *ast.TypeSpec, src string) TypeName {
	return TypeName{getType(t.Type, src), nameFromIdent(t.Name)}
}
