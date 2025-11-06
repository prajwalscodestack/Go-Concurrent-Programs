package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func worker(workerId int, jobsChan chan int, wg *sync.WaitGroup, rateLimiter *rate.Limiter) {
	defer wg.Done()
	for job := range jobsChan {
		_ = rateLimiter.Wait(context.Background())
		fmt.Printf("worker %d has processed job %d\n", workerId, job)
		time.Sleep(2 * time.Second)
	}
}
func main() {
	numJobs := 20
	var wg sync.WaitGroup
	jobsChan := make(chan int, numJobs)
	r := 5
  burst:=5
	rateLimiter := rate.NewLimiter(rate.Limit(r), burst)
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
