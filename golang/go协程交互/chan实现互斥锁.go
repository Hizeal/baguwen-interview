package main

import (
	"fmt"
	"sync"
)

var num int

func add(h chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	h <- 1
	num += 1
	<-h

}

func main() {
	ch := make(chan int, 1)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go add(ch, &wg)
	}
	wg.Wait()

	fmt.Println(num)
}
