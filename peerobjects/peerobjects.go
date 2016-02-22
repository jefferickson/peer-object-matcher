package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jefferickson/peer-object-matcher/object"
)

// config vars
var inputFile string
var outputFile string
var maxPeers int
var cpuLimit int

func init() {
	// parse command line args
	flag.StringVar(&inputFile, "input", "", "The input CSV.")
	// flag.StringVar(&outputFile, "output", "", "The output CSV.")
	flag.IntVar(&maxPeers, "maxpeers", 50, "The maximum number of peers.")
	// flag.IntVar(&cpuLimit, "cpulimit", 1, "The max number of CPU cores to utilize.")

	flag.Parse()

	// ensure we have the minimum reqs
	if inputFile == "" {
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	// keep track of time
	start := time.Now()

	// process the input file
	objects := object.ProcessInputCSV(inputFile)

	// for each categorical group, let's calculate the peers
	for groupLabel, categoricalGroup := range objects {
		fmt.Println("Starting %v", groupLabel)
		object.PeerAllObjects(categoricalGroup, maxPeers)
	}

	// Diag: did it maybe work?
	fmt.Println("%v", objects["1"][0].FinalPeers)

	// How long did it take?
	elapsed := time.Since(start)
	fmt.Println("%s", elapsed)
}
