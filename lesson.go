package main

import (
	"fmt"
	"strconv"
)

func main() {
	var s string = "14"
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("ERROR")
	}
	fmt.Printf("%T %v\n", i, i)
}
