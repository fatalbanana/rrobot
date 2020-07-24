package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
)

func main() {

	// Parse command line options
	concurrency := flag.Int("concurrency", 8, "How many goroutines")
	cfgGlob := flag.String("config", "", "Config file(s) to process")
	rspamdURL := flag.String("url", "http://127.0.0.1:11333/checkv2", "Rspamd URL")
	flag.Parse()

	// Use context to stop goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Count tests
	total := 0
	totalPassed := 0
	totalFailed := 0
	defer func() {
		fmt.Println(fmt.Sprintf("\n%d tests processed, %d passed, %d failed", total, totalPassed, totalFailed))
	}()

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
		for i := 0; i < *concurrency; i++ {
			wg.Add(1)
			go processJobs(ctx, jobCh, resultCh, &wg, *rspamdURL)
		}
		wg.Wait()
		close(resultCh)
	}()

	// Catch interrupt signal
	signal.Notify(sigCh, os.Interrupt)

	// Process config files
	matches, err := filepath.Glob(*cfgGlob)
	if err != nil {
		fmt.Println("Couldn't read config: %s", err.Error())
		return
	}
	for _, match := range matches {
		realCfg, err := ReadConfig(match)
		if err != nil {
			fmt.Println("Couldn't parse config(%s): %s", match, err)
			return
		}
		for _, cfg := range realCfg.Tests {
			configCh <- cfg
		}
	}
	close(configCh)

	// Wait to die
	for {
		select {
		case <-sigCh:
			return
		case result, ok := <-resultCh:
			if !ok {
				return
			}
			if !result.Passed {
				fmt.Println(fmt.Sprintf("FAILED: %s(%s): %s", result.Name, result.File, strings.Join(result.Errors, ",")))
				totalFailed++
			} else {
				totalPassed++
			}
			total++
		}
	}
}
