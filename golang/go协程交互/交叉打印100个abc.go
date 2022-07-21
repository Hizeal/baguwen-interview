package main

import(
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main(){
	Achan := make(chan int,1)
	Bchan := make(chan int)
	Cchan := make(chan int)

	wg.Add(3)

	go A(Achan,Bchan)
	go B(Bchan,Cchan)
	go C(Cchan,Achan)

	wg.Wait()
}

func A(Achan chan int,Bchan chan int){
	def wg.Done()

	for i:=0;i<100;i++{
		Achan <- 1
		fmt.Println("a")
		<- Bchan 
	}
	return
}



func B(Bchan chan int,Cchan chan int){
	def wg.Done()

	for i:=0;i<100;i++{
		Bchan <- 1
		fmt.Println("b")
		<- Cchan 
	}
	return
}

func C(Cchan chan int,Achan chan int){
	def wg.Done()

	for i:=0;i<100;i++{
		Cchan <- 1
		fmt.Println("c")
		<- Achan 
	}
	return
}