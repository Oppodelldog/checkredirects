package main

import (
	"flag"
	"io/ioutil"
	"bufio"
	"bytes"
	"strings"
	"net/http"
	"github.com/pkg/errors"
	"fmt"
)

func main() {
	concurrentConnections := ReadConcurrentConnections()
	workerQueue, checkResultChannel := CreateRedirectWorkers(concurrentConnections)

	redirects := ReadRedirects()
	go CheckRedirects(redirects, workerQueue)

	redirectsChecked := 0
	isRunning := true
	for isRunning {
		select {
		case checkResult, ok := <-checkResultChannel:
			if !ok {
				isRunning = false
			}

			fmt.Printf("checking '%s': ", checkResult.redirect.source)

			if checkResult.err != nil {
				fmt.Printf("%v\n", checkResult.err)
			} else {
				fmt.Println("OK")
			}
			redirectsChecked++
			if redirectsChecked == len(redirects) {
				isRunning = false
			}
		}
	}
	close(checkResultChannel)
	close(workerQueue)
}

func CreateRedirectWorkers(numberOfConcurrentChecks int) (chan Redirect, chan CheckResult) {
	workerQueue := make(chan Redirect, numberOfConcurrentChecks)
	checkResultChannel := make(chan CheckResult, numberOfConcurrentChecks)

	for i := 0; i < numberOfConcurrentChecks; i++ {
		go func() {
			for redirect := range workerQueue {
				CheckRedirect(redirect, checkResultChannel)
			}
		}()
	}

	return workerQueue, checkResultChannel
}

func CheckRedirects(redirects []Redirect, redirectQueue chan Redirect) {
	for _, redirect := range redirects {
		redirectQueue <- redirect
	}
}

func CheckRedirect(redirect Redirect, results chan CheckResult) {
	result := CheckResult{redirect: redirect}
	response, err := http.Get(redirect.source)
	if err != nil {
		result.err = err
		results <- result
		return
	}

	if response.Request.URL.String() != redirect.target {
		result.err = errors.Errorf("target url did not match: wanted '%s', got '%s'", redirect.target, response.Request.URL.String())
		results <- result
		return
	}

	results <- result

	return
}

type Redirect struct {
	source string
	target string
}

type CheckResult struct {
	redirect Redirect
	err      error
}

func ReadRedirects() (redirects []Redirect) {
	fileContent, err := ioutil.ReadFile("redirects")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(fileContent))
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, "\t")
		redirects = append(redirects, Redirect{
			cols[0],
			cols[1],
		})
	}

	return
}

func ReadConcurrentConnections() int {

	concurrent := flag.Int("c", 1, "c=2")

	flag.Parse()

	if concurrent == nil {
		return 1
	}
	return *concurrent
}
