package main

import (
	"fmt"
	"sync"
)

//simple ping-pong game between two go-routines

func Ping(pingCh, pongCh chan string, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		// Ping RECEIVES "pong"
		message := <-pongCh
		fmt.Println(message)

		// On the last iteration, don’t send "ping" again (avoid deadlock)
		if i < n-1 {
			pingCh <- "ping"
		}
	}
}

func Pong(pingCh, pongCh chan string, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		// Pong RECEIVES "ping"
		message := <-pingCh
		fmt.Println(message)

		// Always reply with "pong"
		pongCh <- "pong"
	}
}

func main() {
	const N = 10
	//we can also do this using buffered channels
	pingCh := make(chan string)
	pongCh := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	go Ping(pingCh, pongCh, N, &wg)
	go Pong(pingCh, pongCh, N, &wg)

	// Start with "ping"
	pingCh <- "ping"

	wg.Wait()
}
