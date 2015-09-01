package annotate

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var typ = `
package test

type INT int

type POINTER *float64

type ARRAY [10]int

type SLICE []string

var ELLIPSIS = [...]int{1,2,3}

type MAP map[INT]string

type CHAN chan int

type FUNC func(int)string

`

func TestType(t *testing.T) {
	Convey("Variable test", t, func() {
		f, err := ParseFile(typ)
		So(err, ShouldBeNil)
		var p Package
		p = append(p, f)
		Convey("ident type", func() {
			t, err := p.FindType("INT")
			So(err, ShouldBeNil)
			So(t.Type.(IdentType).Ident, ShouldEqual, "int")
			So(t.Name, ShouldEqual, "INT")
		})

		Convey("pointer type", func() {
			t, err := p.FindType("POINTER")
			So(err, ShouldBeNil)
			So(t.Type.(StarType).Target.(IdentType).Ident, ShouldEqual, "float64")
			So(t.Name, ShouldEqual, "POINTER")
		})

		Convey("array type", func() {
			t, err := p.FindType("ARRAY")
			So(err, ShouldBeNil)
			So(t.Type.(ArrayType).Type.(IdentType).Ident, ShouldEqual, "int")
			So(t.Type.(ArrayType).Len, ShouldEqual, 10)
			So(t.Type.(ArrayType).Slice, ShouldBeFalse)
			So(t.Name, ShouldEqual, "ARRAY")
		})

		Convey("slice type", func() {
			t, err := p.FindType("SLICE")
			So(err, ShouldBeNil)
			So(t.Type.(ArrayType).Type.(IdentType).Ident, ShouldEqual, "string")
			So(t.Type.(ArrayType).Len, ShouldEqual, 0)
			So(t.Type.(ArrayType).Slice, ShouldBeTrue)
			So(t.Name, ShouldEqual, "SLICE")
		})

		Convey("Ellipsis type", func() {
			t, err := p.FindVariable("ELLIPSIS")
			So(err, ShouldBeNil)
			So(t.Type.(EllipsisType).Type.(IdentType).Ident, ShouldEqual, "int")
			So(t.Type.(EllipsisType).Len, ShouldEqual, 0)
			So(t.Type.(EllipsisType).Slice, ShouldBeFalse)
		})

		Convey("map type", func() {
			t, err := p.FindType("MAP")
			So(err, ShouldBeNil)
			So(t.Type.(MapType).Key.(IdentType).Ident, ShouldEqual, "INT")
			So(t.Type.(MapType).Value.(IdentType).Ident, ShouldEqual, "string")
		})

		Convey("chan type", func() {
			t, err := p.FindType("CHAN")
			So(err, ShouldBeNil)
			So(t.Type.(ChannelType).Type.(IdentType).Ident, ShouldEqual, "int")
			So(t.Type.(ChannelType).Direction, ShouldEqual, 3)
		})

		Convey("func type", func() {
			t, err := p.FindType("FUNC")
			So(err, ShouldBeNil)
			So(len(t.Type.(FuncType).Parameters), ShouldEqual, 1)
			So(t.Type.(FuncType).Parameters[0].Type.(IdentType).Ident, ShouldEqual, "int")
			So(len(t.Type.(FuncType).Results), ShouldEqual, 1)
			So(t.Type.(FuncType).Results[0].Type.(IdentType).Ident, ShouldEqual, "string")
		})

	})
}
