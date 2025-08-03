package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(ch chan int, i int) {
	// something
	ch <- i * 2
}

func consumer(ch chan int, wg *sync.WaitGroup) {
	for i := range ch {
		fmt.Println("process", i*1000)
		wg.Done()
	}
	fmt.Println("########################")
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go producer(ch, i)
	}
	go consumer(ch, &wg)
	wg.Wait()
	// close(ch)
	time.Sleep(2 * time.Second)
	fmt.Println("done")
}
