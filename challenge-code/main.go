package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func printMessage(s string, wg *sync.WaitGroup) {

	defer wg.Done()
	fmt.Println(s)

}

func main() {
	message := []string{
		"hello mars",
		"hello venus",
		"hello Jupiter ",
		"hello saturn",
	}
	wg.Add(len(message))
	for _, x := range message {
		go printMessage(x, &wg)

	}
	wg.Wait()

	fmt.Println("last messsage to be printed")
}
