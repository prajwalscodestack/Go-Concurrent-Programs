package main

import (
	"fmt"
	"sync"
)

type ChanMutex struct {
	ch chan struct{}
}

func NewMutex() *ChanMutex {
	mu := &ChanMutex{ch: make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}
func (mu *ChanMutex) Lock() {
	<-mu.ch
}
func (mu *ChanMutex) Unlock() {
	mu.ch <- struct{}{}
}

var count int

func main() {
	var wg sync.WaitGroup
	mutex := NewMutex()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mutex.Lock()
			count++
			mutex.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println(count)
}
