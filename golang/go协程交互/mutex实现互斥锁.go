package main

import (
	"fmt"
	"sync"
)

var num int
var wg sync.WaitGroup
var mtx sync.Mutex

func add() {
	mtx.Lock()

	defer mtx.Unlock()
	defer wg.Done()

	num += 1
}

func main() {
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go add()
	}

	wg.Wait()

	fmt.Println(num)
}
