package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/jefferickson/peer-object-matcher/object"
)

// config vars
var inputFile string
var lagFile string
var outputFile string
var maxPeers int
var maxBlockSize int
var cpuLimit int

func init() {
	// parse command line args
	flag.StringVar(&inputFile, "input", "", "The input CSV. [REQUIRED]")
	flag.StringVar(&lagFile, "lag", "", "The input CSV for lag peers.")
	flag.StringVar(&outputFile, "output", "", "The output CSV. [REQUIRED]")
	flag.IntVar(&maxPeers, "maxpeers", 50, "The maximum number of peers.")
	flag.IntVar(&maxBlockSize, "maxblocksize", 5000, "The maximum number of objects per routine.")
	flag.IntVar(&cpuLimit, "cpulimit", 0, "The max number of CPU cores to utilize.")

	flag.Parse()

	// ensure we have the minimum reqs
	if inputFile == "" || outputFile == "" {
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	// keep track of time
	start := time.Now()

	// set max number of CPUs
	runtime.GOMAXPROCS(cpuLimit)

	// process the input file and start up! if there is a lag file, process that too
	// otherwise set those in the input file to be their own peers
	objects, total_n := object.ProcessInputCSV(inputFile)
	var peers map[string][]*object.Object
	if lagFile != "" {
		peers, _ = object.ProcessInputCSV(lagFile)
	} else {
		// the objects will be the peers themselves
		peers = objects
	}
	allObjects := object.ObjectsToPeer{Objects: objects, Peers: peers, N: total_n}
	allObjects.Run(maxPeers, outputFile, maxBlockSize)

	// How long did it take?
	elapsed := time.Since(start)
	fmt.Printf("\n")
	fmt.Println("Completed in", elapsed)
}
