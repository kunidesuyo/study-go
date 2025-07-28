package main

import "fmt"

func foo(params ...int) {
	fmt.Println(len(params), params)
	for _, param := range params {
		fmt.Println(param)
	}
}

func main() {
	s := make([]int, 3)
	fmt.Println(s)
	foo(s...)
}
