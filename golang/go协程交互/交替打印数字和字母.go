package main

import (
	"fmt"
	"sync"
	"unicode/utf8"
)

func main() {
	number, letter := make(chan bool), make(chan bool)

	wg := sync.WaitGroup{}

	go func() {
		i := 1
		for {
			<-number //从number获取数字消息
			fmt.Printf("%d%d", i, i+1)
			i += 2
			letter <- true
		}
	}()

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		i := 0
		for {
			<-letter                              //接收管道消息
			if i >= utf8.RuneCountInString(str) { //走到字符串结尾
				wg.Done()
				return
			}
			fmt.Print(str[i : i+2])
			i += 2
			number <- true
		}
	}(&wg)

	number <- true
	wg.Wait()
}
