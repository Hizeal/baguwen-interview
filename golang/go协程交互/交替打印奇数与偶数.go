package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i++ {
			ch <- true
			if i%2 == 1 {
				fmt.Printf("routine 1: %d\n", i)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i < 101; i++ {
			<-ch
			if i%2 == 0 {
				fmt.Printf("routine 2: %d\n", i)
			}

		}
	}()
	wg.Wait()
}
