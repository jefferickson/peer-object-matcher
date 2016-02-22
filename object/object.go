package object

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/jefferickson/peer-object-matcher/utils"
)

// The base objects we are going to be matching.
type Object struct {
	ID           string
	Categorical  string
	NoMatchGroup string
	Coords       []float64
	PeerComps    []PeerComp
	PeerDists    []float64
	FinalPeers   []string
}

// To store the distances to other potential peers
type PeerComp struct {
	PeerObject *Object
	Distance   float64
}

// Factory function to create an object
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

// Function to add a distance to another peer
func (o *Object) addPeerComp(peer *Object) {
	distanceToPeer, err := PeerObjects(o, peer, utils.EuclideanDistance)
	if err != nil {
		return
	}

	peerCompTemp := PeerComp{
		PeerObject: peer,
		Distance:   distanceToPeer,
	}

	o.PeerComps = append(o.PeerComps, peerCompTemp)
	o.PeerDists = append(o.PeerDists, distanceToPeer)
}

// Function to find the closest peers
func (o *Object) findClosestPeers(peers []*Object, n int) {
	// if PeerDists doesn't exist, then we need to calculate them
	if o.PeerComps == nil {
		for _, peer := range peers {
			o.addPeerComp(peer)
		}
	}

	// make sure we have the same length of comps and dists, otherwise something went wrong
	if len(o.PeerComps) != len(o.PeerDists) {
		log.Panic(errors.New("PeerComps and PeerDists lengths differ."))
	}

	// find the closest n peers
	var maxDistance float64
	if len(o.PeerDists) > n {
		sort.Float64s(o.PeerDists)
		maxDistance = o.PeerDists[n-1]
	} else {
		maxDistance = o.PeerDists[len(o.PeerDists)-1]
	}

	for _, finalPeer := range o.PeerComps {
		if finalPeer.Distance <= maxDistance {
			o.FinalPeers = append(o.FinalPeers, finalPeer.PeerObject.ID)
		}
	}
	sort.Strings(o.FinalPeers)
	// TODO: 'smartly' delete those over n so that we have exactly n
}

// Peer all objects in a group with all other objects in a group
func PeerAllObjects(objects []*Object, n int) {
	for _, object := range objects {
		object.findClosestPeers(objects, n)
	}
}

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
	objects := make(map[string][]*Object)
	for _, row := range rawCSVdata {
		objects[row[1]] = append(objects[row[1]], NewObject(row[0], row[1], row[2], row[3:len(row)]))
	}

	return objects
}
