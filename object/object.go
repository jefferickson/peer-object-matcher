package object

import (
	"fmt"
	"log"
	"strconv"
)

// The base objects we are going to be matching.
type Object struct {
	ID           string
	Categorical  string
	NoMatchGroup string
	Coords       []float64
}

func NewObject(ID string, Categorical string, NoMatchGroup string, Coords []string) *Object {
	// convert coords to floats
	var coordsAsFloat []float64
	for _, s := range Coords {
		coord, err := strconv.ParseFloat(s, 64)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		coordsAsFloat = append(coordsAsFloat, coord)
	}

	return &Object{
		ID:           ID,
		Categorical:  Categorical,
		NoMatchGroup: NoMatchGroup,
		Coords:       coordsAsFloat,
	}
}
