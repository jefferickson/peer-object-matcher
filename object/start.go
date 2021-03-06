package object

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// this stores a map of all objects we are going to peer
type ObjectsToPeer struct {
	Objects map[string][]*Object
	Peers 	map[string][]*Object
	N       int
}

// The main function: everything starts here
func (o *ObjectsToPeer) Run(maxPeers int, outputFile string, maxBlockSize int) {
	// open file for output
	outputCSV, err := os.Create(outputFile)
	if err != nil {
		log.Panic(err)
	}
	defer outputCSV.Close()

	// semaphore for concurrency, add for each object we have
	var wg sync.WaitGroup
	wg.Add(o.N)

	// channels to report progress
	toCounter := make(chan bool, o.N)
	toWrite := make(chan *Object, o.N)

	// for each categorical group, let's calculate the peers within since peers do not span multiple categorical groups
	for categoricalLabel, categoricalGroup := range o.Objects {
		// cache and mutex for this categorical group
		cache := cacheAndMutex{make(map[string][]string), new(sync.RWMutex)}

		// we are going to break this up into subgroups to better load the processing across go routines
		nCategoricalGroup := len(categoricalGroup)
		totalSubgroups := nCategoricalGroup/maxBlockSize + 1
		for totalProcessed := 0; totalProcessed < totalSubgroups; totalProcessed++ {
			// what are the bounds of this partition
			start := totalProcessed * maxBlockSize
			end := start + maxBlockSize
			if end > nCategoricalGroup {
				end = nCategoricalGroup
			}
			p := peerSliceAndPool{categoricalGroup[start:end], o.Peers[categoricalLabel]}

			// start go routine on this peer slice
			go func(p peerSliceAndPool, cache cacheAndMutex) {
				peerAllObjects(p, maxPeers, cache, toWrite, toCounter)
			}(p, cache)
		}
	}

	// start the counter to report progress
	go counter(toCounter, o.N)

	// start the reporter that will write out results to CSV and clean up
	go o.writeAndCleanUp(toWrite, outputCSV, &wg)

	// wait for all go routines to complete
	wg.Wait()
	close(toCounter)
	close(toWrite)
}

// Write results to disk and clean up
func (o *ObjectsToPeer) writeAndCleanUp(ch <-chan *Object, outputFile *os.File, wg *sync.WaitGroup) {
	for {
		objectToWrite, ok := <-ch
		if ok {
			outputToCSV(objectToWrite, outputFile)
			wg.Done()
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
			fmt.Printf("\rPeering: %d/%d", current, n)
		} else {
			return
		}
	}
}
