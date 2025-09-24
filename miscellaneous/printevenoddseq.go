package main

import (
	"fmt"
	"sync"
)

func printOdd(oddCh, evenCh chan struct{}, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= n; i += 2 {
		<-oddCh // wait for turn
		fmt.Println(i)
		if i+1 <= n { // only signal if there's a next even number
			evenCh <- struct{}{}
		}
	}
}

func printEven(oddCh, evenCh chan struct{}, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 2; i <= n; i += 2 {
		<-evenCh // wait for turn
		fmt.Println(i)
		if i+1 <= n { // only signal if there's a next odd number
			oddCh <- struct{}{}
		}
	}
}

func main() {
	const N = 10
	oddCh := make(chan struct{})
	evenCh := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(2)
	go printOdd(oddCh, evenCh, N, &wg)
	go printEven(oddCh, evenCh, N, &wg)

	// Start with odd
	oddCh <- struct{}{}

	wg.Wait()
}
