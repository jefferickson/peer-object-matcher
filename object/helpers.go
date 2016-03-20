package object

import (
    "bytes"
    "encoding/csv"
    "errors"
    "log"
    "os"
)

// Determine distance between two peers
// Or return error if they cannot be peered
func PeerObjects(obj1 *Object, obj2 *Object, distfn func([]float64, []float64) (float64, error)) (float64, error) {
    distance := 0.0

    // check to make sure they have the same categorical values
    if obj1.Categorical != obj2.Categorical {
        return distance, errors.New("Categorical data does not match.")
    }

    // check to make sure they aren't in the same no match group
    if obj1.NoMatchGroup == obj2.NoMatchGroup {
        return distance, errors.New("noMatchGroups match.")
    }

    // these objects can be peered
    // calc the distance between them
    distance, err := distfn(obj1.Coords, obj2.Coords)
    return distance, err
}

// Read in input CSV
func ProcessInputCSV(inputFile string) map[string][]*Object {
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
    // according to its categorical data
    objects := make(map[string][]*Object)
    for _, row := range rawCSVdata {
        objects[row[1]] = append(objects[row[1]], NewObject(row[0], row[1], row[2], row[3:len(row)]))
    }

    return objects
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
