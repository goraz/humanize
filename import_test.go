package humanize

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var imprt1 = `
package test

// Doc
import (
	"testing"
    onn "github.com/fzerorubigd/onion"
    _ "github.com/lib/pq"
    // Other
	. "github.com/smartystreets/goconvey/convey"
    "github.com/fzerorubigd/annotate"
)

`

func TestImport(t *testing.T) {
	Convey("Import test ", t, func() {
		var p = &Package{}
		f, err := ParseFile(imprt1, p)
		So(err, ShouldBeNil)

		p.Files = append(p.Files, f)
		Convey("Import testing", func() {
			i, err := p.FindImport("testing")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, "testing")
			So(i.Path, ShouldEqual, "testing")
			So(len(i.Docs), ShouldEqual, 1)
			So(i.Docs[0], ShouldEqual, "// Doc")
		})

		Convey("Import onion", func() {
			i, err := p.FindImport("onn")
			So(err, ShouldBeNil)
			i2, err := p.FindImport("github.com/fzerorubigd/onion")
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(i, i2), ShouldBeTrue)
			So(i.Name, ShouldEqual, "onn")
			So(i.Path, ShouldEqual, "github.com/fzerorubigd/onion")
			So(len(i.Docs), ShouldEqual, 1)
			So(i.Docs[0], ShouldEqual, "// Doc")
		})

		Convey("Import pq", func() {
			i, err := p.FindImport("github.com/lib/pq")
			So(err, ShouldBeNil)
			_, err = p.FindImport("_")
			So(err, ShouldNotBeNil)
			So(i.Name, ShouldEqual, "_")
			So(i.Path, ShouldEqual, "github.com/lib/pq")
			So(len(i.Docs), ShouldEqual, 1)
			So(i.Docs[0], ShouldEqual, "// Doc")
		})

		Convey("Import convey", func() {
			i, err := p.FindImport("github.com/smartystreets/goconvey/convey")
			So(err, ShouldBeNil)
			So(i.Name, ShouldEqual, ".")
			So(i.Path, ShouldEqual, "github.com/smartystreets/goconvey/convey")
			So(len(i.Docs), ShouldEqual, 2)
			So(i.Docs[0], ShouldEqual, "// Doc")
			So(i.Docs[1], ShouldEqual, "// Other")
		})

		Convey("Import annotate", func() {
			i, err := p.FindImport("annotate")
			So(err, ShouldBeNil)
			i2, err := p.FindImport("github.com/fzerorubigd/annotate")
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(i, i2), ShouldBeTrue)
			So(i.Name, ShouldEqual, "annotate")
			So(i.Path, ShouldEqual, "github.com/fzerorubigd/annotate")
			So(len(i.Docs), ShouldEqual, 1)
			So(i.Docs[0], ShouldEqual, "// Doc")
		})

	})
}
