package object

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// this stores a map of all objects we are going to peer
type ObjectsToPeer struct {
	Objects map[string][]*Object
	N       int
}

// The main function: everything starts here
func (o *ObjectsToPeer) Run(maxPeers int, outputFile string) {
	// open file for output
	outputCSV, err := os.Create(outputFile)
	if err != nil {
		log.Panic(err)
	}
	defer outputCSV.Close()

	// semaphore for concurrency
	var wg sync.WaitGroup
	wg.Add(len(o.Objects))

	// channels to report progress
	toCounter := make(chan bool)
	toWrite := make(chan *Object)

	// for each categorical group, let's calculate the peers on a separate thread
	for _, categoricalGroup := range o.Objects {
		go func(group []*Object) {
			defer wg.Done()

			// create a cache then start peering on this group
			cache := make(map[string]cachedFinalPeers)
			peerAllObjects(group, maxPeers, cache, toWrite, toCounter)
		}(categoricalGroup)
	}

	// start the counter to report progress
	go counter(toCounter, o.N)

	// start the reporter that will write out results to CSV and clean up
	go o.writeAndCleanUp(toWrite, outputCSV)

	// wait for all go routines to complete
	wg.Wait()
	close(toCounter)
	close(toWrite)
}

// Write results to disk and clean up
func (o *ObjectsToPeer) writeAndCleanUp(ch <-chan *Object, outputFile io.Writer) {
	for {
		objectToWrite, ok := <-ch
		if ok {
			outputToCSV(objectToWrite, outputFile)
			// Since everything is written out for this object, let's reset the FinalPeers
			// to save space.
			objectToWrite.FinalPeers = nil
		} else {
			return
		}
	}
}

// Report progress
func counter(ch <-chan bool, n int) {
	current := 0
	for {
		_, ok := <-ch
		if ok {
			current++
			fmt.Println(current, "/", n)
		} else {
			return
		}
	}
}
