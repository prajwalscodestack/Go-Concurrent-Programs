package main

import (
	"fmt"
	"sync"
)

type Job struct {
	Data int
}

func worker(jobsChan chan Job, resultsChan chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobsChan {
		resultsChan <- Job{Data: job.Data * job.Data}
	}
}

func main() {
	jobs := []int{2, 3, 4, 5, 6, 7, 8, 9, 10}
	jobsChan := make(chan Job)
	resultsChan := make(chan Job, len(jobs))
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(jobsChan, resultsChan, &wg)
	}
	//sendjobs here
	for _, job := range jobs {
		jobsChan <- Job{Data: job}
	}
	close(jobsChan)
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	for res := range resultsChan {
		fmt.Println(res)
	}
}
