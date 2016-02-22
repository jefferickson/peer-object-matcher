package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/jefferickson/peer-object-matcher/object"
)

// Euclidean distance
func EuclideanDistance(coords1 []float64, coords2 []float64) (float64, error) {
	sum := 0.0

	// must be the same number of dims
	if len(coords1) != len(coords2) {
		return sum, errors.New("Different number of dimensions.")
	}

	// calculate sum of square of diffs
	for i := 0; i < len(coords1); i++ {
		sum += math.Pow(coords1[i]-coords2[i], 2)
	}

	return math.Sqrt(sum), nil
}

// Read in input CSV
func ProcessInputCSV(inputFile string) *map[string]*object.Object {
	csvfile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	// for each row, let's create a new object and store it in map
	// according to its categorical data
	objects := make(map[string]*object.Object)
	for _, row := range rawCSVdata {
		objects[row[1]] = object.NewObject(row[0], row[1], row[2], row[3:len(row)])
	}

	return &objects
}
