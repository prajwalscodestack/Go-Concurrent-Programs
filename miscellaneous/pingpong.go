package main

import (
	"fmt"
	"sync"
)

const N = 10

func Ping(ping, pong chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < N; i++ {
		ping <- "ping"
		fmt.Println(<-pong)
	}
}
func Pong(ping, pong chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < N; i++ {
		fmt.Println(<-ping)
		pong <- "pong"
	}
}
func main() {
	var wg sync.WaitGroup
	ping := make(chan string)
	pong := make(chan string)
	wg.Add(2)
	go Ping(ping, pong, &wg)
	go Pong(ping, pong, &wg)
	wg.Wait()
}
