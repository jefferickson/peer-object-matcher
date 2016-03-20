package main

import (
	"flag"
	"fmt"
	"io"
	"log"
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
	flag.StringVar(&outputFile, "output", "", "The output CSV.")
	flag.IntVar(&maxPeers, "maxpeers", 50, "The maximum number of peers.")
	// flag.IntVar(&cpuLimit, "cpulimit", 1, "The max number of CPU cores to utilize.")

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

	// process the input file, open file for output
	objects, total_n := object.ProcessInputCSV(inputFile)
	outputCSV, err := os.Create(outputFile)
	if err != nil {
		log.Panic(err)
	}
	defer outputCSV.Close()

	// semaphore for concurreny
	var wg sync.WaitGroup
	wg.Add(len(objects))

	// channels to report progress
	counter := make(chan bool)
	toWrite := make(chan string)

	// TODO: move the creation of go routines to helper functions.
	// for each categorical group, let's calculate the peers
	for groupLabel, categoricalGroup := range objects {
		go func(group []*object.Object, groupLabel string, counter chan<- bool, writer chan<- string) {
			defer wg.Done()
			object.PeerAllObjects(group, maxPeers, counter)
			writer <- groupLabel
		}(categoricalGroup, groupLabel, counter, toWrite)
	}

	// start the counter to report progress
	go func(counter <-chan bool, n int) {
		current := 0
		for {
			_, ok := <-counter
			if ok {
				current++
				fmt.Println(current, "/", n)
			} else {
				return
			}
		}
	}(counter, total_n)

	// start the reporter that will write out results to CSV and clean up
	go func(toWrite <-chan string, outputFile io.Writer) {
		for {
			categoricalGroup, ok := <-toWrite
			if ok {
				object.OutputToCSV(objects[categoricalGroup], outputFile)
				delete(objects, categoricalGroup)
			} else {
				return
			}
		}
	}(toWrite, outputCSV)

	// wait for all go routines to complete
	wg.Wait()
	close(counter)
	close(toWrite)

	// How long did it take?
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
