package main

import (
	"fmt"
	"sync"
)

/*import (
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
}*/

/*USING TWO GOROUNTINES TO PRINT THE MESSAGE USING CHANNELS*/

var s string
var wg sync.WaitGroup

func updateMesaage(s string, ch chan string) {
	defer wg.Done() /* I will wait for updateMessage to finish — not for printMessage. */
	ch <- s
}

func printMessage(ch chan string) {
	msq := <-ch
	fmt.Println(msq)
}

func main() {
	message := []string{
		"hello mars",
		"hello venus",
		"hell jupiter",
		"hello pluto",
	}
	for _, x := range message {
		ch := make(chan string)

		wg.Add(1)
		go updateMesaage(x, ch)
		go printMessage(ch)
	}
	wg.Wait()
	fmt.Println("last message to be printed")
}
