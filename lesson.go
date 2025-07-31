package main

import (
	"fmt"
	// "log"
	// "os"
	// "io"
)

func one(x *int) {
	*x = 1
}

func main() {
	var n int = 100
	one(&n)
	fmt.Println(n)
	fmt.Println(&n)
	fmt.Println(&*&n)
}
