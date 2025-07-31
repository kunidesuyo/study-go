package main

import (
	"fmt"
	// "log"
	// "os"
	// "io"
)

func do (i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Println(v * 2)
	case string:
		fmt.Println(v + "!")
	default:
		fmt.Printf("I don't know %T\n", v)
	}
}

func main() {
	do(10)
	do("Mike")
	do(true)
}