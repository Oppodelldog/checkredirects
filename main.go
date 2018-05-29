package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	defaultNumberOfConcurrentConnections = 1
	redirectsFileName                    = "redirects"
)

type (
	Redirect struct {
		source string
		target string
	}

	CheckResult struct {
		redirect Redirect
		err      error
	}

	verifyRedirectFuncDef func(redirect Redirect) error
)

var verifyRedirectFunc = verifyRedirectFuncDef(VerifyRedirect)

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
				err := verifyRedirectFunc(redirect)
				result := CheckResult{
					redirect: redirect,
					err:      err,
				}

				checkResultChannel <- result
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

func ReadRedirects() (redirects []Redirect) {
	fileContent, err := ioutil.ReadFile(redirectsFileName)
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

	concurrent := flag.Int("c", defaultNumberOfConcurrentConnections, "c=2")

	flag.Parse()

	return *concurrent
}
