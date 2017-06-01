// Copyright 2017 <Company Name>, Inc. All Rights Reserved.

package main

import (
	"flag"
	"fmt"
	"github.com/harapr-jc/hashgen"
	"os"
)

// Filled by command line flags
var (
	host string
	port string
)

func init() {
	const (
		defaultHost = "localhost"
		defaultPort = "8080"
	)
	flag.StringVar(&host, "host", defaultHost, "Server Host")
	flag.StringVar(&port, "port", defaultPort, "Listen Port")
	// Custom usage message
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Starts crypto hash server\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {

	flag.Parse()

	// go install github.com/harapr-jc/jchash

	hashgen.StartServer(host, port)

	// Block (waiting for exit signal)
	<-hashgen.Exit
}
