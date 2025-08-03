package main

import (
	"fmt"
	"time"
)

func producer(first chan int) {
	defer close(first)
	for i := 0; i < 10; i++ {
		first <- i
	}
}

func multi2(first <-chan int, second chan<- int) {
	defer close(second)
	for i := range first {
		second <- i * 2
	}
}

func multi4(second <-chan int, third chan<- int) {
	defer close(third)
	for i := range second {
		third <- i * 4
	}
}

func goroutine1(ch chan string) {
	for {
		ch <- "packet from 1"
		time.Sleep(3 * time.Second)
	}
}

func goroutine2(ch chan int) {
	for {
		ch <- 100
		time.Sleep(1 * time.Second)
	}
}

// func consumer(ch chan int, wg *sync.WaitGroup) {
// 	for i := range ch {
// 		fmt.Println("process", i*1000)
// 		wg.Done()
// 	}
// 	fmt.Println("########################")
// }

func main() {
	c1 := make(chan string)
	c2 := make(chan int)
	go goroutine1(c1)
	go goroutine2(c2)
	for {
		select {
		case msg1 := <-c1:
			fmt.Println(msg1)
		case msg2 := <-c2:
			fmt.Println(msg2)
		}
	}

}
