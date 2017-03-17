package humanize

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const invalidFunc = `
package invalid
var s = invalid_func()

`
const invalidFunc2 = `
package invalid

func invalid_func() {

}

var s = invalid_func()

`
const invalidImport = `
package invalid

var s = im.invalid_func()

`
const invalidImport2 = `
package invalid

import "im"
var s = im.invalid_func()

`

const invalidImport3 = `
package invalid
import (
    os "github.com/goraz/annotate/fixture"
)
var s = os.invalid_func()

`
const invalidImport4 = `
package invalid
import (
     os "github.com/goraz/humanize/fixture"
)
var s = os.NoReturn()

`
const validImport5 = `
package invalid
import "github.com/goraz/humanize/fixture"
var s = fixture.NewF()
var f = fixture.NewFile()
`

func TestCurrentPackage(t *testing.T) {
	Convey("Current package", t, func() {
		p, err := ParsePackage("github.com/goraz/humanize")
		So(err, ShouldBeNil)
		Convey("not found items", func() {
			So(p.Path, ShouldEqual, "github.com/goraz/humanize")
			So(p.Name, ShouldEqual, "humanize")
			_, err := p.FindFunction("invalid_function")
			So(err, ShouldNotBeNil)
			_, err = p.FindImport("invalid_import")
			So(err, ShouldNotBeNil)
			_, err = p.FindType("invalid_type")
			So(err, ShouldNotBeNil)
			_, err = p.FindVariable("invalid_variable")
			So(err, ShouldNotBeNil)
			_, err = p.FindConstant("invalid_constant")
			So(err, ShouldNotBeNil)
		})

		Convey("translate path", func() {
			_, err := translateToFullPath("invalid_path")
			So(err, ShouldNotBeNil)
			_, err = translateToFullPath("github.com/goraz/humanize/type.go")
			So(err, ShouldNotBeNil)

		})

		Convey("fact about the fixture", func() {
			p, err := ParsePackage("github.com/goraz/humanize/fixture")
			So(err, ShouldBeNil)
			tt, err := p.FindType("f")
			So(err, ShouldBeNil)
			So(len(tt.Methods), ShouldEqual, 1)
			So(tt.Methods[0].Name, ShouldEqual, "f.Test")
			So(len(tt.StarMethods), ShouldEqual, 1)
			So(tt.StarMethods[0].Name, ShouldEqual, "f.TestStar")
			So(len(tt.GetAllMethods(false)), ShouldEqual, 1)
			So(len(tt.GetAllMethods(true)), ShouldEqual, 2)

			tt2, err := p.FindType("f2")
			So(err, ShouldBeNil)
			So(len(tt2.Methods), ShouldEqual, 0)
			So(len(tt2.StarMethods), ShouldEqual, 0)
			So(len(tt2.GetAllMethods(false)), ShouldEqual, 2)
			So(len(tt2.GetAllMethods(true)), ShouldEqual, 2)

			t1, err := p.FindType("T1")
			So(err, ShouldBeNil)
			it1 := t1.Type.(*InterfaceType)
			So(tt.Support(it1, false), ShouldBeTrue)
			So(tt.Support(it1, true), ShouldBeTrue)
			So(tt2.Support(it1, true), ShouldBeTrue)

			t2, err := p.FindType("T2")
			So(err, ShouldBeNil)
			it2 := t2.Type.(*InterfaceType)
			So(tt.Support(it2, false), ShouldBeFalse)
			So(tt.Support(it2, true), ShouldBeTrue)
			So(tt2.Support(it2, true), ShouldBeTrue)

			t3, err := p.FindType("T3")
			So(err, ShouldBeNil)
			it3 := t3.Type.(*InterfaceType)
			So(tt.Support(it3, false), ShouldBeFalse)
			So(tt.Support(it3, true), ShouldBeTrue)
			So(tt2.Support(it3, true), ShouldBeTrue)

			Convey("return unexported", func() {
				var p = &Package{}
				f, err := ParseFile(validImport5, p)
				So(err, ShouldBeNil)
				p.Files = append(p.Files, f)

				err = lateBind(p)
				So(err, ShouldBeNil)

				s, err := p.FindVariable("s")
				So(err, ShouldBeNil)
				So(s.Name, ShouldEqual, "s")
			})

		})
	})
}

func TestErrors(t *testing.T) {
	Convey("invalid file", t, func() {
		var p = &Package{}
		f, err := ParseFile(invalidFunc, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})

	Convey("invalid file 2", t, func() {
		var p = &Package{}
		f, err := ParseFile(invalidFunc2, p)
		So(err, ShouldBeNil)

		p.Files = append(p.Files, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})

	Convey("invalid import", t, func() {
		var p = &Package{}
		f, err := ParseFile(invalidImport, p)
		So(err, ShouldBeNil)

		p.Files = append(p.Files, f)

		// TODO : support build tags and enable this two line again
		//err = lateBind(p)
		//So(err, ShouldNotBeNil)
	})

	Convey("invalid import 2", t, func() {
		var p = &Package{}
		f, err := ParseFile(invalidImport2, p)
		So(err, ShouldBeNil)

		p.Files = append(p.Files, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})

	Convey("invalid import 3", t, func() {
		var p = &Package{}
		f, err := ParseFile(invalidImport3, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})

	Convey("invalid import 4", t, func() {
		var p = &Package{}

		f, err := ParseFile(invalidImport4, p)
		So(err, ShouldBeNil)
		p.Files = append(p.Files, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})
}
