package object

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

// Read in input CSV
func ProcessInputCSV(inputFile string) (map[string][]*Object, int) {
	csvfile, err := os.Open(inputFile)
	if err != nil {
		log.Panic(err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	// for each row, let's create a new object and store it in map
	// according to its categorical data (since we will want to keep them separate)
	objects := make(map[string][]*Object)
	for _, row := range rawCSVdata {
		objects[row[1]] = append(objects[row[1]], newObject(row[0], row[1], row[2], row[3:len(row)]))
	}

	return objects, len(rawCSVdata)
}

// Output results to CSV
func outputToCSV(record *Object, outputFile io.Writer) {
	w := bufio.NewWriter(outputFile)
	// create record
	fmt.Fprint(w, record.ID)
	for _, peer := range record.FinalPeers {
		fmt.Fprint(w, ",", peer)
	}
	fmt.Fprint(w, "\n")
	w.Flush()
}

// Generate a key for the cache
func genCacheKey(a string, b string, cs []string) string {
	var key bytes.Buffer
	key.WriteString(a)
	key.WriteString(",")
	key.WriteString(b)
	key.WriteString(",")
	for _, c := range cs {
		key.WriteString(c)
		key.WriteString(",")
	}

	return key.String()
}

// Generate a hash for sorting
func genMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
