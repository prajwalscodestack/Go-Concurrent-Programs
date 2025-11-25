package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(id int, jobs chan int, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	for {
		select {
		// Option A: Check context for cancellation (urgent stop)
		case <-ctx.Done():
			// If the work item was long-running, we'd check ctx here.
			fmt.Printf("worker %d closed by context cancellation\n", id)
			return

		// Option B: Check job channel for work (standard pool exit)
		case j, ok := <-jobs:
			if !ok {
				// The job channel was closed, meaning no more jobs are coming.
				fmt.Printf("worker %d: job channel closed, exiting gracefully\n", id)
				return
			}
			// Process the job
			time.Sleep(100 * time.Millisecond) // Use a smaller duration for faster test
			fmt.Printf("worker %d processed %d\n", id, j)
		}
	}
}
func main() {
	//graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan int)
	go func() {
		<-stop
		fmt.Println("stopped..")
		cancel()
	}()
	var wg sync.WaitGroup
	workerCount := 5
	jobCount := 100
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg, ctx)
	}
	for i := 0; i < jobCount; i++ {
		select {
		case <-ctx.Done():
			// Shutdown requested, stop submitting jobs
			fmt.Printf("Shutdown requested. Submitted %d jobs.\n", i)
			goto submissionComplete // Use goto for clean exit from loop
		case jobs <- i:
			// Job submitted successfully
		}
	}
submissionComplete:
	close(jobs)

	wg.Wait()
}
