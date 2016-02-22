package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jefferickson/peer-object-matcher/utils"
)

// config vars
var inputFile string
var outputFile string
var cpuLimit int

func init() {
	// parse command line args
	flag.StringVar(&inputFile, "input", "", "The input CSV.")
	// flag.StringVar(&outputFile, "output", "", "The output CSV.")
	// flag.IntVar(&cpuLimit, "cpulimit", 1, "The max number of CPU cores to utilize.")

	flag.Parse()

	// ensure we have the minimum reqs
	if inputFile == "" {
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	objects := utils.ProcessInputCSV(inputFile)
	fmt.Println("%v", objects)
}
