package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var typ = `
package test

import (
	"net/http"
	"github.com/fzerorubigd/onion"
)

type INT int

type POINTER *float64

type ARRAY [10]int

type SLICE []string

var ELLIPSIS = [...]int{1,2,3}

type MAP map[INT]string

type CHAN chan int

type CHAN2 chan<- int

type CHAN3 <-chan int

type FUNC func(int)string

type SEL onion.Layer

type STRUCT struct {
   N SEL   ` + "`json:\"tag\"`" + `
   M MAP
   X int
}

type EMBEDSTRUCT struct {
   STRUCT  ` + "`json:\"tag\"`" + `
}

type INTERFACE interface {
   Test(int, INT, FUNC) (FUNC, error)
}

type EMBEDINTERFACE interface {
   INTERFACE
}

type EMPTYSTRUCT struct{}

type EMPTYINTERFACE interface{}

var on = onion.New()

var xx = http.ConnState(10)

`

const wrongTypeCast = `
package example

var test = invalidFuncAndCast(10)
`

const wrongTypeCast2 = `
package example

import "net/http"

var test = http.invalidFuncAndCast(10)
`

const validTypeCast = `
package example

type XX int

var test = int64(1)

var test2 = XX(2)
`

func TestType(t *testing.T) {
	Convey("Variable test", t, func() {
		var p = &Package{}

		f, err := ParseFile(typ, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)
		Convey("ident type", func() {
			t, err := p.FindType("INT")
			So(err, ShouldBeNil)
			So(t.Type.(*IdentType).Ident, ShouldEqual, "int")
			So(t.Name, ShouldEqual, "INT")
			So(t.Type.GetDefinition(), ShouldEqual, "int")
			So(t.GetDefinition(), ShouldEqual, "INT int")
		})

		Convey("pointer type", func() {
			t, err := p.FindType("POINTER")
			So(err, ShouldBeNil)
			So(t.Type.(*StarType).Target.(*IdentType).Ident, ShouldEqual, "float64")
			So(t.Name, ShouldEqual, "POINTER")
			So(t.Type.GetDefinition(), ShouldEqual, "*float64")
			So(t.GetDefinition(), ShouldEqual, "POINTER *float64")
		})

		Convey("array type", func() {
			t, err := p.FindType("ARRAY")
			So(err, ShouldBeNil)
			So(t.Type.(*ArrayType).Type.(*IdentType).Ident, ShouldEqual, "int")
			So(t.Type.(*ArrayType).Len, ShouldEqual, 10)
			So(t.Type.(*ArrayType).Slice, ShouldBeFalse)
			So(t.Name, ShouldEqual, "ARRAY")
			So(t.Type.GetDefinition(), ShouldEqual, "[10]int")
		})

		Convey("slice type", func() {
			t, err := p.FindType("SLICE")
			So(err, ShouldBeNil)
			So(t.Type.(*ArrayType).Type.(*IdentType).Ident, ShouldEqual, "string")
			So(t.Type.(*ArrayType).Len, ShouldEqual, 0)
			So(t.Type.(*ArrayType).Slice, ShouldBeTrue)
			So(t.Name, ShouldEqual, "SLICE")
			So(t.Type.GetDefinition(), ShouldEqual, "[]string")
		})

		Convey("Ellipsis type", func() {
			t, err := p.FindVariable("ELLIPSIS")
			So(err, ShouldBeNil)
			So(t.Type.(*EllipsisType).Type.(*IdentType).Ident, ShouldEqual, "int")
			So(t.Type.(*EllipsisType).Len, ShouldEqual, 0)
			So(t.Type.(*EllipsisType).Slice, ShouldBeFalse)
			So(t.Type.GetDefinition(), ShouldEqual, "[...]int{}")
		})

		Convey("map type", func() {
			t, err := p.FindType("MAP")
			So(err, ShouldBeNil)
			So(t.Type.(*MapType).Key.(*IdentType).Ident, ShouldEqual, "INT")
			So(t.Type.(*MapType).Value.(*IdentType).Ident, ShouldEqual, "string")
			So(t.Type.GetDefinition(), ShouldEqual, "map[INT]string")
		})

		Convey("chan type", func() {
			t, err := p.FindType("CHAN")
			So(err, ShouldBeNil)
			So(t.Type.(*ChannelType).Type.(*IdentType).Ident, ShouldEqual, "int")
			So(t.Type.(*ChannelType).Direction, ShouldEqual, 3)
			So(t.Type.GetDefinition(), ShouldEqual, "chan int")

			t, err = p.FindType("CHAN2")
			So(err, ShouldBeNil)
			So(t.Type.(*ChannelType).Type.(*IdentType).Ident, ShouldEqual, "int")
			So(t.Type.(*ChannelType).Direction, ShouldEqual, 1)
			So(t.Type.GetDefinition(), ShouldEqual, "chan<- int")

			t, err = p.FindType("CHAN3")
			So(err, ShouldBeNil)
			So(t.Type.(*ChannelType).Type.(*IdentType).Ident, ShouldEqual, "int")
			So(t.Type.(*ChannelType).Direction, ShouldEqual, 2)
			So(t.Type.GetDefinition(), ShouldEqual, "<-chan int")

		})

		Convey("func type", func() {
			t, err := p.FindType("FUNC")
			So(err, ShouldBeNil)
			So(len(t.Type.(*FuncType).Parameters), ShouldEqual, 1)
			So(t.Type.(*FuncType).Parameters[0].Type.(*IdentType).Ident, ShouldEqual, "int")
			So(len(t.Type.(*FuncType).Results), ShouldEqual, 1)
			So(t.Type.(*FuncType).Results[0].Type.(*IdentType).Ident, ShouldEqual, "string")
			So(t.Type.GetDefinition(), ShouldEqual, "func (int) string")
		})

		Convey("select type", func() {
			t, err := p.FindType("SEL")
			So(err, ShouldBeNil)
			So(t.Type.(*SelectorType).pkg.Name, ShouldEqual, "onion")
			So(t.Type.(*SelectorType).Type.(*IdentType).Ident, ShouldEqual, "Layer")
			So(t.Type.GetDefinition(), ShouldEqual, "onion.Layer")
			So(t.Type.(*SelectorType).IdentType().GetDefinition(), ShouldEqual, "Layer")
		})

		Convey("struct type", func() {
			t, err := p.FindType("STRUCT")
			So(err, ShouldBeNil)
			So(len(t.Type.(*StructType).Fields), ShouldEqual, 3)
			So(t.Type.(*StructType).Fields[0].Name, ShouldEqual, "N")
			So(t.Type.(*StructType).Fields[0].Tags.Get("json"), ShouldEqual, "tag")
			So(t.Type.(*StructType).Fields[0].Type.(*IdentType).Ident, ShouldEqual, "SEL")

			So(t.Type.(*StructType).Fields[1].Name, ShouldEqual, "M")
			So(t.Type.(*StructType).Fields[1].Tags, ShouldEqual, "")
			So(t.Type.(*StructType).Fields[1].Type.(*IdentType).Ident, ShouldEqual, "MAP")

			So(t.Type.(*StructType).Fields[2].Name, ShouldEqual, "X")
			So(t.Type.(*StructType).Fields[2].Tags, ShouldEqual, "")
			So(t.Type.(*StructType).Fields[2].Type.(*IdentType).Ident, ShouldEqual, "int")
			So(t.Type.GetDefinition(), ShouldEqual, "struct{\n\tN SEL `json:\"tag\"`\n\tM MAP \n\tX int \n}")
		})

		Convey("embed struct type", func() {
			t, err := p.FindType("EMBEDSTRUCT")
			So(err, ShouldBeNil)
			So(len(t.Type.(*StructType).Fields), ShouldEqual, 0)
			So(len(t.Type.(*StructType).Embeds), ShouldEqual, 1)
			So(t.Type.(*StructType).Embeds[0].Type.(*IdentType).Ident, ShouldEqual, "STRUCT")
			So(t.Type.(*StructType).Embeds[0].Tags.Get("json"), ShouldEqual, "tag")
			So(t.Type.GetDefinition(), ShouldEqual, "struct{\n\tSTRUCT\n}")
		})

		Convey("empty struct type", func() {
			t, err := p.FindType("EMPTYSTRUCT")
			So(err, ShouldBeNil)
			So(len(t.Type.(*StructType).Fields), ShouldEqual, 0)
			So(len(t.Type.(*StructType).Embeds), ShouldEqual, 0)
			So(t.Type.GetDefinition(), ShouldEqual, "struct{}")
		})

		Convey("interface type", func() {
			t, err := p.FindType("INTERFACE")
			So(err, ShouldBeNil)
			So(len(t.Type.(*InterfaceType).Functions), ShouldEqual, 1)
			So(t.Type.(*InterfaceType).Functions[0].Name, ShouldEqual, "Test")

			So(len(t.Type.(*InterfaceType).Functions[0].Type.Parameters), ShouldEqual, 3)
			So(len(t.Type.(*InterfaceType).Functions[0].Type.Results), ShouldEqual, 2)
			So(t.Type.GetDefinition(), ShouldEqual, `interface{
	func Test(int,INT,FUNC) (FUNC,error)
}`)
		})

		Convey("embed interface type", func() {
			t, err := p.FindType("EMBEDINTERFACE")
			So(err, ShouldBeNil)
			So(len(t.Type.(*InterfaceType).Functions), ShouldEqual, 0)
			So(len(t.Type.(*InterfaceType).Embed), ShouldEqual, 1)
			So(t.Type.(*InterfaceType).Embed[0].(*IdentType).Ident, ShouldEqual, "INTERFACE")
			So(t.Type.GetDefinition(), ShouldEqual, "interface{\n\tINTERFACE\n}")
		})

		Convey("empty interface type", func() {
			t, err := p.FindType("EMPTYINTERFACE")
			So(err, ShouldBeNil)
			So(len(t.Type.(*InterfaceType).Functions), ShouldEqual, 0)
			So(len(t.Type.(*InterfaceType).Embed), ShouldEqual, 0)
			So(t.Type.GetDefinition(), ShouldEqual, "interface{}")
		})

		Convey("selector type", func() {
			lateBind(p)
			t, err := p.FindVariable("on")
			So(err, ShouldBeNil)
			sel, ok := t.Type.(*SelectorType)
			So(ok, ShouldBeTrue)
			p2 := sel.Package()
			So(p2.Path, ShouldEqual, "github.com/fzerorubigd/onion")

			v, err := p.FindVariable("xx")
			So(err, ShouldBeNil)
			So(v.Type.GetDefinition(), ShouldEqual, "http.ConnState")
		})

	})

	Convey("invalid type cast", t, func() {
		p := &Package{}
		f, err := ParseFile(wrongTypeCast, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)
		So(lateBind(p), ShouldNotBeNil)
	})

	Convey("invalid type cast 2", t, func() {
		p := &Package{}
		f, err := ParseFile(wrongTypeCast2, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)
		So(lateBind(p), ShouldNotBeNil)
	})

	Convey("valid type cast", t, func() {
		p := &Package{}
		f, err := ParseFile(validTypeCast, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)
		So(lateBind(p), ShouldBeNil)
		v, err := p.FindVariable("test")
		So(err, ShouldBeNil)
		So(v.Name, ShouldEqual, "test")
		So(v.Type.GetDefinition(), ShouldEqual, "int64")

		v, err = p.FindVariable("test2")
		So(err, ShouldBeNil)
		So(v.Name, ShouldEqual, "test2")
		So(v.Type.GetDefinition(), ShouldEqual, "XX")

	})
}
