package main

import "fmt"

func main() {
	l := []string{"python", "go", "java"}

	for i, v := range l {
		fmt.Println(i, v)
	}
}
