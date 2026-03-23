package main

import (
	"io"
	"os"
	"sync"
)

func NewprintSomething(s string, wg *sync.WaitGroup) {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	var wg1 *sync.WaitGroup
	wg.Add(1)
	go printSomething("lamba", wg1)
	wg1.Wait()

	w.Close()
	result, _ := io.ReadAll(r)
	output := string(result)
	os.Stdout = stdOut
	println(output)

}
