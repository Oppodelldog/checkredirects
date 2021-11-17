package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/Oppodelldog/checkredirects/internal"
)

const defaultNumberOfConcurrentConnections = 1
const defaultDelimiter = "\t"
const defaultFilename = "redirects"

func main() {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	var filename = flagSet.String("f", defaultFilename, "-f=redirects.txt")
	if len(*filename) == 0 {
		*filename = defaultFilename
	}

	var concurrent = flagSet.Int("c", defaultNumberOfConcurrentConnections, "-c=2")
	if *concurrent == 0 {
		*concurrent = defaultNumberOfConcurrentConnections
	}

	var delimiter = flagSet.String("d", defaultDelimiter, "-d=;")
	if len(*delimiter) == 0 {
		*delimiter = defaultDelimiter
	}

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		return
	}

	internal.Check(*filename, *concurrent, *delimiter)
}
