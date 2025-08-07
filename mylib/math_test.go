package mylib

import (
	"fmt"
	"testing"
)

var Debug bool = true

func TestAverage(t *testing.T) {
	if Debug {
		t.Skip("Skiping testing")
	}

	v := Average([]int{1, 2, 3, 4, 5})
	if v != 3 {
		t.Error("Expected 3, got", v)
	}
}

func ExampleAverage() {
	v := Average([]int{1, 2, 3, 4, 5})
	fmt.Println(v)
}

func Example() {
	v := Average([]int{1, 2, 3, 4, 5, 6})
	fmt.Println(v)
}

func ExamplePerson2_Say() {
	p := Person2{Name: "Mike", Age: 20}
	p.Say()
}
