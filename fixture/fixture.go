package fixture

import "os"

var X, err = os.Open("the_file")

var Y *os.File

type f struct {
}

type f2 struct {
	*f
}

type T1 interface {
	Test()
}

type T2 interface {
	TestStar()
}

type T3 interface {
	Test()
	TestStar()
}

func NewF() *f {
	return &f{}
}

func NewFile() *os.File {
	return Y
}

func NoReturn() {

}

func (f) Test() {

}

func (*f) TestStar() {

}
