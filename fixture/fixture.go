package fixture

import "os"

var X, err = os.Open("the_file")

var Y *os.File

type f struct {
}

func NewF() *f {
	return &f{}
}

func NewFile() *os.File {
	return Y
}
