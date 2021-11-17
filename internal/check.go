package internal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
)

const RedirectsFileName = "redirects"

type (
	Redirect struct {
		Source string
		Target string
	}

	CheckResult struct {
		Redirect Redirect
		Err      error
	}

	VerifyRedirectFuncDef func(ctx context.Context, redirect Redirect) error
)

var VerifyRedirectFunc = VerifyRedirectFuncDef(VerifyRedirect)

func Check(concurrentConnections int) {
	ctx := context.Background()
	workerQueue, checkResultChannel := CreateRedirectWorkers(ctx, concurrentConnections)

	redirects := ReadRedirects()
	go CheckRedirects(redirects, workerQueue)

	redirectsChecked := 0
	isRunning := true

	for isRunning {
		checkResult, ok := <-checkResultChannel
		if !ok {
			isRunning = false
		}

		fmt.Printf("checking '%s': ", checkResult.Redirect.Source)

		if checkResult.Err != nil {
			fmt.Printf("%v\n", checkResult.Err)
		} else {
			fmt.Println("OK")
		}

		redirectsChecked++
		if redirectsChecked == len(redirects) {
			isRunning = false
		}
	}
	close(checkResultChannel)
	close(workerQueue)
}

func CreateRedirectWorkers(ctx context.Context, numberOfConcurrentChecks int) (chan Redirect, chan CheckResult) {
	workerQueue := make(chan Redirect, numberOfConcurrentChecks)
	checkResultChannel := make(chan CheckResult, numberOfConcurrentChecks)

	for i := 0; i < numberOfConcurrentChecks; i++ {
		go func() {
			for redirect := range workerQueue {
				err := VerifyRedirectFunc(ctx, redirect)
				result := CheckResult{
					Redirect: redirect,
					Err:      err,
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
	fileContent, err := ioutil.ReadFile(RedirectsFileName)
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
