package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(workerId int, jobsChan chan int, wg *sync.WaitGroup, rateLimiter <-chan time.Time) {
	defer wg.Done()
	for job := range jobsChan {
		<-rateLimiter
		fmt.Printf("worker %d has processed job %d\n", workerId, job)
		time.Sleep(2 * time.Second)
	}
}
func main() {
	fmt.Println("Hello, World!")
	numJobs := 20
	var wg sync.WaitGroup
	jobsChan := make(chan int, numJobs)
	rate := 5
	rateLimiter := time.Tick(time.Second / time.Duration(rate))
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(i, jobsChan, &wg, rateLimiter)
	}
	for i := 1; i <= numJobs; i++ {
		jobsChan <- i
	}
	close(jobsChan)
	wg.Wait()
}
