package main

import (
	"flag"

	"gitlab.com/Oppodelldog/checkredirects/internal"
)

const defaultNumberOfConcurrentConnections = 1

func main() {
	internal.Check(ReadConcurrentConnections())
}

func ReadConcurrentConnections() int {
	concurrent := flag.Int("c", defaultNumberOfConcurrentConnections, "-c=2")

	flag.Parse()

	return *concurrent
}
