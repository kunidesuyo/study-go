package main

import (
	"fmt"
	// "log"
	// "os"
	// "io"
)

type Myint int

func (i Myint) Double() Myint {
	return i * 2
}

func main() {
	myInt := Myint(10)
	fmt.Println(myInt.Double())
}
