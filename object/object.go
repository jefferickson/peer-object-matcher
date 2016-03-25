package object

import (
	"errors"
	"log"
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
	PeerComps    peerComps
	FinalPeers   []string
	CacheKey     string
	HashedID     string
	MaxPeers     int
}

// The slice and pool you want to peer
type peerSliceAndPool struct {
	slice []*Object
	pool  []*Object
}

// To store the distances to other potential peers
type peerComp struct {
	PeerObject *Object
	Distance   float64
}
type peerComps []peerComp

// To store cached peer comparisons
type cachedFinalPeers []string
type peeredCache map[string]cachedFinalPeers

// Peer all objects in a group with all other objects in a group
func peerAllObjects(p peerSliceAndPool, n int, cache peeredCache, writer chan<- *Object, counter chan<- bool) {
	for _, object := range p.slice {
		object.findClosestPeers(p.pool, n, cache)
		writer <- object
		counter <- true
	}
}

// Factory function to create an object
func newObject(ID string, Categorical string, NoMatchGroup string, Coords []string, MaxPeers int) *Object {
	// convert coords to floats
	var coordsAsFloat []float64
	for _, s := range Coords {
		coord, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Panic(err)
		}
		coordsAsFloat = append(coordsAsFloat, coord)
	}

	return &Object{
		ID:           ID,
		Categorical:  Categorical,
		NoMatchGroup: NoMatchGroup,
		Coords:       coordsAsFloat,
		CacheKey:     genCacheKey(Categorical, NoMatchGroup, Coords),
		HashedID:     genMD5Hash(ID),
		MaxPeers:     MaxPeers,
	}
}

// Function to find the closest peers
func (o *Object) findClosestPeers(peers []*Object, n int, cache map[string]cachedFinalPeers) {
	// first check to see if they are already calculated and cached
	if cached, ok := cache[o.CacheKey]; ok {
		o.FinalPeers = cached
		return
	} else {
		// allocate memory to store all of the peer comparisons
		o.PeerComps = make([]peerComp, 0, o.MaxPeers)
		for _, peer := range peers {
			o.addPeerComp(peer)
		}
	}

	// find the closest n peers
	numOfPeers := len(o.PeerComps)
	if numOfPeers > n {
		sort.Sort(o.PeerComps)
		numOfPeers = n
	}
	o.FinalPeers = make([]string, 0, numOfPeers)
	for i := 0; i < numOfPeers; i++ {
		o.FinalPeers = append(o.FinalPeers, o.PeerComps[i].PeerObject.ID)
	}

	// make sure we didn't exceed the allowed number
	if len(o.FinalPeers) > n {
		log.Panic(errors.New("Too many peers!"))
	}

	// Sort the final peers by ID
	sort.Strings(o.FinalPeers)

	// We no longer need the actual comps so clear up that space
	o.PeerComps = nil

	// Store the results into the cache
	cache[o.CacheKey] = o.FinalPeers
}

// Function to add a distance to another peer
func (o *Object) addPeerComp(peer *Object) {
	distanceToPeer, err := peerObjects(o, peer, utils.EuclideanDistance)
	if err != nil {
		// we don't want to peer these objects
		return
	}

	peerCompTemp := peerComp{
		PeerObject: peer,
		Distance:   distanceToPeer,
	}

	o.PeerComps = append(o.PeerComps, peerCompTemp)
}
