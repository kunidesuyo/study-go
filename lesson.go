package main

import (
	"fmt"
	// "log"
	// "os"
	// "io"
)

type Vertex struct {
	X, Y int
	S string
}

func changeVertex(v Vertex) {
	v.X = 1000
}

func changeVertex2(v *Vertex) {
	(*v).X = 1000
}

func main() {
	v2 := &Vertex{X: 1, Y: 2, S: "test"}
	changeVertex2(v2)
	fmt.Println(v2)
}
