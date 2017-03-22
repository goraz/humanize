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
	// pkg return the package where this type is inside it
	Package() *Package
	// TypeName return the type with name of this type
	//TypeName() *TypeName
}

type srcBase struct {
	pkg *Package
	src string
}

func (s srcBase) Package() *Package {
	return s.pkg
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

// Field is a single field of a structure, a variable, with tag
type Field struct {
	Variable
	Tags reflect.StructTag
}

// Embed is the embeded type in the struct or interface
type Embed struct {
	Type
	Docs Docs
	Tags reflect.StructTag
}

// StructType is a struct in source code
type StructType struct {
	srcBase
	Fields []*Field
	Embeds []*Embed
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
	Functions []*Function
	Embed     []Type // IdentType or SelectorType
}

// SelectorType a type from another package
type SelectorType struct {
	srcBase
	pkg   *Import
	ident *IdentType
	Type  Type
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

//TypeName contain type and its name, means the type is in this package
type TypeName struct {
	Type Type
	Name string
	Docs Docs

	Methods     []*Function
	StarMethods []*Function
}

// Package in selector type is not this package
func (st *SelectorType) Package() *Package {
	p := st.pkg.LoadPackage()
	return p
}

func (st *SelectorType) IdentType() *IdentType {
	if st.ident == nil {
		st.ident = &IdentType{
			Ident: st.Type.GetDefinition(),
			srcBase: srcBase{
				pkg: st.pkg.LoadPackage(),
				src: st.src,
			},
		}
	}

	return st.ident
}

// GetDefinition return the definition of this type
func (tn TypeName) GetDefinition() string {
	return tn.Name + " " + tn.Type.GetDefinition()
}

func getTypeName(t Type) (*TypeName, bool) {
	var pointer bool
	if t2, ok := t.(*StarType); ok {
		t = t2.Target
		pointer = true
	}

	if t2, ok := t.(*SelectorType); ok {
		// its in another package, load it from there
		p := t2.pkg.LoadPackage()
		t3, err := p.FindType(t2.Type.GetDefinition())
		assertNil(err)
		return t3, pointer
	}

	tn, err := t.Package().FindType(t.GetDefinition())
	assertNil(err)
	return tn, pointer
}

func (tn TypeName) GetAllMethods(pointer bool) []*Function {
	met := tn.Methods
	if pointer {
		met = append(met, tn.StarMethods...)
	}
	// there is a small problem with struct types. if there is an struct type,
	// then the embeded functions are important too.
	if st, ok := tn.Type.(*StructType); ok && len(st.Embeds) > 0 {
		// there is two case. one is when the struct is inside this package, other when its a selector
		// both can be pointer, if pointer then the StartMethods are available
		for i := range st.Embeds {
			tn, pn := getTypeName(st.Embeds[i].Type)
			met = append(met, tn.GetAllMethods(pn)...)
		}
	}

	return met
}

func getInterfaceFunc(in *InterfaceType) []*Function {
	fn := in.Functions
	for i := range in.Embed {
		p := in.Embed[i].Package()
		tn, err := p.FindType(removeReceiver(in.Embed[i].GetDefinition()))
		assertNil(err)
		ni := tn.Type.(*InterfaceType)
		fn = append(fn, getInterfaceFunc(ni)...)
	}
	return fn
}

// Support return true if the type support the interface, if pointer is true then it checked with
// pointer receiver
func (tn TypeName) Support(in *InterfaceType, pointer bool) bool {
	two := tn.GetAllMethods(pointer)

	one := getInterfaceFunc(in)

	return compare(one, two)
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
	if len(i.Embeds) == 0 && len(i.Fields) == 0 {
		return "struct{}"
	}
	res := "struct{\n"
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
	return i.pkg.Name + "." + i.Type.GetDefinition()
}

// GetName the name of this type
func (i *FuncType) GetDefinition() string {
	return "func " + i.getSign()
}

func (i *FuncType) getDefinitionWithName(name string) string {
	return "func " + name + i.getSign()
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
	if len(i.Embed) == 0 && len(i.Functions) == 0 {
		return "interface{}"
	}

	res := "interface{\n"
	for e := range i.Embed {
		res += "\t" + i.Embed[e].GetDefinition() + "\n"
	}
	for f := range i.Functions {
		res += "\t" + i.Functions[f].Type.getDefinitionWithName(i.Functions[f].Name) + "\n"
	}
	return res + "}"
}

func getSource(e ast.Expr, src string) string {
	res := ""
	start := e.Pos() - 1
	end := e.End() - 1
	// grab it in source
	if len(src) >= int(end) {
		res = src[start:end]
	}
	return res
}

func getImport(name string, f *File) (res *Import) {
	for _, i := range f.Imports {
		if i.Name == name {
			res = i
			break
		}
	}
	return
}

func getType(e ast.Expr, src string, f *File, p *Package) Type {
	switch t := e.(type) {
	case *ast.Ident:
		// ident is the simplest one.
		return &IdentType{
			srcBase{p, getSource(e, src)},
			nameFromIdent(t),
		}
	case *ast.StarExpr:
		return &StarType{
			srcBase{p, getSource(e, src)},
			getType(t.X, src, f, p),
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
			srcBase{p, getSource(e, src)},
			t.Len == nil,
			l,
			getType(t.Elt, src, f, p),
		}
		if ellipsis {
			at = &EllipsisType{at.(*ArrayType)}
		}
		return at
	case *ast.MapType:
		return &MapType{
			srcBase{p, getSource(e, src)},
			getType(t.Key, src, f, p),
			getType(t.Value, src, f, p),
		}

	case *ast.StructType:
		res := &StructType{srcBase{p, getSource(e, src)}, nil, nil}
		for _, s := range t.Fields.List {
			if s.Names != nil {
				for i := range s.Names {
					v := Variable{
						Name: nameFromIdent(s.Names[i]),
						Type: getType(s.Type, src, f, p),
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
					res.Fields = append(res.Fields, &f)
				}
			} else {
				e := Embed{
					Type: getType(s.Type, src, f, p),
				}
				if s.Tag != nil {
					e.Tags = reflect.StructTag(s.Tag.Value)
					e.Tags = e.Tags[1 : len(e.Tags)-1]
				}
				e.Docs = docsFromNodeDoc(s.Doc)
				res.Embeds = append(res.Embeds, &e)
			}
		}

		return res
	case *ast.InterfaceType:
		// TODO : interface may refer to itself I need more time to implement this
		iface := &InterfaceType{
			srcBase: srcBase{p, getSource(e, src)},
		}
		for i := range t.Methods.List {
			res := Function{}
			// The method name is mandatory and always 1
			if len(t.Methods.List[i].Names) > 0 {
				res.Name = nameFromIdent(t.Methods.List[i].Names[0])

				res.Docs = docsFromNodeDoc(t.Methods.List[i].Doc)
				typ := getType(t.Methods.List[i].Type, src, f, p)
				res.Type = typ.(*FuncType)
				iface.Functions = append(iface.Functions, &res)
			} else {
				// This is the embeded interface
				embed := getType(t.Methods.List[i].Type, src, f, p)
				iface.Embed = append(iface.Embed, embed)
			}

		}
		return iface
	case *ast.ChanType:
		return &ChannelType{
			srcBase:   srcBase{p, getSource(e, src)},
			Direction: t.Dir,
			Type:      getType(t.Value, src, f, p),
		}
	case *ast.SelectorExpr:
		return &SelectorType{
			srcBase: srcBase{p, getSource(e, src)},
			pkg:     getImport(nameFromIdent(t.X.(*ast.Ident)), f),
			Type:    getType(t.Sel, src, f, p),
		}
	case *ast.FuncType:
		return &FuncType{
			srcBase:    srcBase{p, getSource(e, src)},
			Parameters: extractVariableList(t.Params, src, f, p),
			Results:    extractVariableList(t.Results, src, f, p),
		}
	}

	return nil
}

// NewType handle a type
func NewType(t *ast.TypeSpec, c *ast.CommentGroup, src string, f *File, p *Package) *TypeName {
	doc := docsFromNodeDoc(c, t.Doc)
	return &TypeName{
		Docs: doc,
		Type: getType(t.Type, src, f, p),
		Name: nameFromIdent(t.Name),
	}
}
