package main

import (
	"fmt"
	"sync"
)

func secondPrintSomething(s string, wg2 *sync.WaitGroup) {
	defer wg2.Done()
	fmt.Println(s)

}
func main1() {
	var wg2 sync.WaitGroup
	words := []string{
		"swayam",
		"shailesh",
		"siddharth",
		"subham",
		"preet",
	}
	wg2.Add(len(words))
	for i, x := range words {

		go secondPrintSomething(fmt.Sprintf("%d : %s", i, x), &wg2)

	}
	wg2.Wait()
	wg2.Add(1)
	secondPrintSomething("this is the last print", &wg2)
}
