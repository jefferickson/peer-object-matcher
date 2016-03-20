package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
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
	objects, total_n := object.ProcessInputCSV(inputFile)

	// semaphore for concurreny
	var wg sync.WaitGroup
	wg.Add(len(objects))

	// channel to report progress
	counter := make(chan bool)

	// for each categorical group, let's calculate the peers
	for _, categoricalGroup := range objects {
		go func(group []*object.Object, counter chan<- bool) {
			defer wg.Done()
			object.PeerAllObjects(group, maxPeers, counter)
		}(categoricalGroup, counter)
	}

	// start the counter to report progress
	go func(counter <-chan bool, n int) {
		current := 0
		for {
			_, ok := <- counter
			if ok {
				current++
				fmt.Println(current, "/", n)
			} else {
				return
			}
		}
	}(counter, total_n)

	// wait for all go routines to complete
	wg.Wait()
	close(counter)

	fmt.Println(objects["1"][0].FinalPeers)

	// How long did it take?
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
