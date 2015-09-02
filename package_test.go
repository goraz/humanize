package annotate

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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
		})

		Convey("fact about the fixture", func() {
			//			st , err := p.FindType(t string)
		})
	})
}
