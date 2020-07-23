package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
)

func main() {

	// Use context to stop goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configs discovered by worker
	configCh := make(chan Config, 1)
	// Jobs compiled from configs
	jobCh := make(chan Job, 1)
	// Results of jobs
	resultCh := make(chan Result, 1)
	// OS signals
	sigCh := make(chan os.Signal, 1)

	// createJobs creates a task for each file to be scanned
	go createJobs(ctx, configCh, jobCh)

	// processJobs runs the individual jobs
	go func() {
		var wg sync.WaitGroup
		// FIXME: configurability
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go processJobs(ctx, jobCh, resultCh, &wg)
		}
		wg.Wait()
		close(resultCh)
	}()

	// Catch interrupt signal
	signal.Notify(sigCh, os.Interrupt)

	// FIXME - walk some directory for HCL files
	err, cfg := ReadConfig("test.hcl")
	if err != nil {
		panic(err)
	}
	configCh <- cfg
	close(configCh)

	// Wait to die
	total := 0
	for {
		select {
		case <-sigCh:
			return
		case result, ok := <-resultCh:
			if !ok {
				fmt.Println(fmt.Sprintf("%d tests processed", total))
				return
			}
			if !result.Passed {
				fmt.Println(fmt.Sprintf("FAILED: %s(%s): %s", result.Name, result.File, strings.Join(result.Errors, ",")))
			}
			total++
		}
	}
}
