// package level documents
package main

// test
import (
	"io/ioutil"
	"os"
	//fuck this shit
	"github.com/fzerorubigd/annotate"
	_ "github.com/lib/pq"
)

var (
	// Booogh
	test, bogh /*doogh*/ string
)

// the hes var
var hes string

// Commet
type example struct {
	p int `tag:"str"`
	q string
}

// On both
type (
	// On any
	any int
	// On Many
	many int64
)

// On Pointer
type pointer *struct{}

// On x chan
type x chan string

// on zz
type zz *example

// on mappex
type mapex map[string]interface{}

// on slice
type slice []pointer

// on arr
type arr [10]int

// on testing
type testing interface {
	Test(string) string
}

// The test is the test
func (m *example) Test(x string) (y string) {
	y = x
	return
}

// Test
func main() {
	r, err := os.Open("/home/f0rud/gospace/src/github.com/fzerorubigd/annotate/test/main.go")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	_, _ = annotate.ParseFile(string(data))
}

//AnnotatedOne is anotated one
// @Test is test dude "one liner"
func AnnotatedOne(p1 string, p2, p3 int, _ int, u example, i *example, j map[string]interface{}) string {
	return ""
}
