package annotate

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
import "os"
var s = os.invalid_func()

`
const invalidImport4 = `
package invalid
import "os"
var s = os.Clearenv()

`
const validImport5 = `
package invalid
import "github.com/goraz/annotate/fixture"
var s = fixture.NewF()
var f = fixture.NewFile()
`

func TestCurrentPackage(t *testing.T) {
	Convey("Current package", t, func() {
		p, err := ParsePackage("github.com/goraz/annotate")
		So(err, ShouldBeNil)
		Convey("not found items", func() {
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
			_, err = translateToFullPath("github.com/goraz/annotate/type.go")
			So(err, ShouldNotBeNil)

		})

		Convey("fact about the fixture", func() {
			Convey("return unexported", func() {
				f, err := ParseFile(validImport5)
				So(err, ShouldBeNil)
				var p Package
				p = append(p, f)

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
		f, err := ParseFile(invalidFunc)
		So(err, ShouldBeNil)
		var p Package
		p = append(p, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})

	Convey("invalid file 2", t, func() {
		f, err := ParseFile(invalidFunc2)
		So(err, ShouldBeNil)
		var p Package
		p = append(p, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})

	Convey("invalid import", t, func() {
		f, err := ParseFile(invalidImport)
		So(err, ShouldBeNil)
		var p Package
		p = append(p, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})

	Convey("invalid import 2", t, func() {
		f, err := ParseFile(invalidImport2)
		So(err, ShouldBeNil)
		var p Package
		p = append(p, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})

	Convey("invalid import 3", t, func() {
		f, err := ParseFile(invalidImport3)
		So(err, ShouldBeNil)
		var p Package
		p = append(p, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})

	Convey("invalid import 4", t, func() {
		f, err := ParseFile(invalidImport4)
		So(err, ShouldBeNil)
		var p Package
		p = append(p, f)

		err = lateBind(p)
		So(err, ShouldNotBeNil)
	})
}
