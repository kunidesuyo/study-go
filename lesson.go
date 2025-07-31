package main

import (
	"fmt"
	// "log"
	// "os"
	// "io"
)

type Human interface {
	Say() string
}

type Person struct {
	Name string
}

func (p *Person) Say() string {
	p.Name = "Mr." + p.Name
	fmt.Println(p.Name)
	return p.Name
}

func DriveCar(human Human) {
	if human.Say() == "Mr. Mike" {
		fmt.Println("Run")
	} else {
		fmt.Println("Get out")
	}
}

type Dog struct {
	Name string
}

func main() {
	var mike Human = &Person{"Mike"}
	DriveCar(mike)

	var x Human = &Person{"X"}
	DriveCar(x)

	var dog Dog = Dog{"dog"}
	DriveCar(dog)

}
