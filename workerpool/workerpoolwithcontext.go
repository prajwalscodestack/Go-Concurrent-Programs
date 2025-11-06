package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func worker(workerId int, jobsChan chan int, wg *sync.WaitGroup, rateLimiter *rate.Limiter, ctx context.Context) {
	defer wg.Done()
	for job := range jobsChan {
		_ = rateLimiter.Wait(context.Background())
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: Job %d cancelled\n", workerId, job)
			return
		default:
			fmt.Printf("Worker %d: Processing job %d\n", workerId, job)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
func main() {
	numJobs := 20
	var wg sync.WaitGroup
	jobsChan := make(chan int, numJobs)
	r := 5
	rateLimiter := rate.NewLimiter(rate.Limit(r), 5)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(i, jobsChan, &wg, rateLimiter, ctx)
	}
	for i := 1; i <= numJobs; i++ {
		jobsChan <- i
	}
	close(jobsChan)
	wg.Wait()
}
