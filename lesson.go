package main

import (
	"fmt"
	// "log"
	// "os"
	// "io"
)

// func one(x *int) {
// 	*x = 1
// }

func main() {
	// s := make([]int, 0)
	// fmt.Printf("%T\n", s)

	// m := make(map[string]int)
	// fmt.Printf("%T\n", m)

	// var p *int = new(int)
	// fmt.Printf("%T\n", p)

	ch := make(chan int)
	fmt.Printf("%T\n", ch)

	var st = new(struct{})
	fmt.Printf("%T\n", st)
}
