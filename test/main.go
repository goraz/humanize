package main

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

// Commet
type example struct {
	p int `tag:"str"`
	q string
}

type (
	any  int
	many int64
)

type pointer *struct{}

type x chan string

type zz *example

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
