// +build ignore

// package level documents
package main

// test
import (
	"encoding/json"
	"fmt"

	"github.com/fzerorubigd/onion"
	"github.com/goraz/humanize"
)

const (
	alpha = iota
	beta
)

const (
	x1 int = iota
	y1
)

var maked = make([]int, 10)

var (
	// Booogh
	test, bogh /*doogh*/ string
)

// the hes var
var hes string

var h pq.NullTime

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

var o = onion.New()

var j, _ = AnnotatedOne("", 1, 1, 1, example{}, nil, nil)

var s func()

var alal onion.Layer

var js = 10

var im complex64

var pop = [...]int{1, 2, 6}

// Test
func main() {
	p, err := humanize.ParsePackage("github.com/goraz/annotate/test")
	if err != nil {
		panic(err)
	}
	c, _ := json.MarshalIndent(p, "", "\t")
	fmt.Print(string(c))
}

//AnnotatedOne is anotated one
// @Test is test dude "one liner"
func AnnotatedOne(p1 string, p2, p3 int, _ int, u example, i *example, j map[string]interface{}, ill ...int) (string, error) {
	return "", nil
}
