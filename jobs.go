package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

type Job struct {
	Assertions []*vm.Program
	Name       string
	File       string
	Headers    map[string]string
}

type JobEnv struct {
	Result RspamdResult
}

func createJobs(ctx context.Context, configCh chan Config, jobCh chan Job) {

	defer close(jobCh)
	exprEnv := expr.Env(JobEnv{})

	for {
		select {
		case <-ctx.Done():
			return
		case config, ok := <-configCh:
			if !ok {
				return
			}
			var err error
			job := Job{
				Assertions: make([]*vm.Program, len(config.Assertions)),
				Name:       config.Name,
				Headers:    config.Headers,
			}
			for i, v := range config.Assertions {
				job.Assertions[i], err = expr.Compile(v, exprEnv)
				if err != nil {
					fmt.Println("Compile failed (%s): %s", v, err)
					return
				}
			}
			for _, input := range config.Inputs {
				matches, err := filepath.Glob(input)
				if err != nil {
					fmt.Println("Glob failed (%s): %s", input, err)
					return
				}
				if len(matches) == 0 {
					fmt.Println("Glob matched nothing: %s", input)
					return
				}
				for _, match := range matches {
					job.File = match
					jobCh <- job
				}
			}
		}
	}
}

func processJobs(ctx context.Context, jobCh chan Job, resultCh chan Result, wg *sync.WaitGroup, url string) {

	defer wg.Done()
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	myVM := vm.VM{}

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobCh:
			if !ok {
				return
			}
			jobResult := Result{
				File:   job.File,
				Name:   job.Name,
				Passed: true,
			}
			file, err := os.Open(job.File)
			if err != nil {
				jobResult.Passed = false
				jobResult.Errors = append(jobResult.Errors,
					fmt.Sprintf("Open failed (%s): %s", job.File, err))
				resultCh <- jobResult
				continue
			}
			req, err := http.NewRequestWithContext(ctx, "POST", url, file)
			if err != nil {
				jobResult.Passed = false
				jobResult.Errors = append(jobResult.Errors,
					fmt.Sprintf("Creating request failed: %s", err))
				resultCh <- jobResult
				continue
			}
			for k, v := range job.Headers {
				req.Header.Add(k, v)
			}
			resp, err := client.Do(req)
			if err != nil {
				jobResult.Passed = false
				jobResult.Errors = append(jobResult.Errors,
					fmt.Sprintf("Sending request failed: %s", err))
				resultCh <- jobResult
				continue
			}
			jobEnv := JobEnv{}
			rspamdResult := RspamdResult{}
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(&rspamdResult)
			resp.Body.Close()
			if err != nil {
				jobResult.Passed = false
				jobResult.Errors = append(jobResult.Errors,
					fmt.Sprintf("Decoding response failed: %s", err))
				resultCh <- jobResult
				continue
			}
			jobEnv.Result = rspamdResult
			for _, assertion := range job.Assertions {
				output, err := myVM.Run(assertion, jobEnv)
				if err != nil || output != true {
					jobResult.Passed = false
					var strErr string
					if err != nil {
						strErr = fmt.Sprintf("Run failed: %s", err.Error())
					} else {
						strErr = fmt.Sprintf("Assertion failed: %s", assertion.Source.Content())
					}
					jobResult.Errors = append(jobResult.Errors, strErr)
				}
			}
			resultCh <- jobResult
		}
	}
}
