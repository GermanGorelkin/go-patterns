package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()
	go func() {
		cond.L.Lock()
		go func() {
			cond.L.Lock()
			cond.Broadcast()
			cond.L.Unlock()
		}()
		cond.Wait()
		fmt.Print("a")
		cond.L.Unlock()
		wg.Done()
	}()
	cond.Wait()
	fmt.Print("b")
	cond.L.Unlock()
	wg.Wait()
	fmt.Println("c")
}
