package object

import (
	"errors"
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/jefferickson/peer-object-matcher/utils"
)

// The base objects we are going to be matching.
type Object struct {
	ID           string
	Categorical  string
	NoMatchGroup string
	Coords       []float64
	PeerComps    []peerComp
	PeerDists    []float64
	FinalPeers   []string
	CacheKey     string
}

// To store the distances to other potential peers
type peerComp struct {
	PeerObject *Object
	Distance   float64
}

// To store cached peer comparisons
type cachedFinalPeers []string

// Semaphore for map access
var mu sync.Mutex

// Peer all objects in a group with all other objects in a group
func peerAllObjects(objects []*Object, n int, cache map[string]cachedFinalPeers, writer chan<- *Object, counter chan<- bool) {
	for _, object := range objects {
		object.findClosestPeers(objects, n, cache)
		writer <- object
		counter <- true
	}
}

// Factory function to create an object
func newObject(ID string, Categorical string, NoMatchGroup string, Coords []string) *Object {
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
	}
}

// Function to find the closest peers
func (o *Object) findClosestPeers(peers []*Object, n int, cache map[string]cachedFinalPeers) {
	// if PeerDists doesn't exist, then we need to calculate them
	if o.PeerComps == nil {
		// first check to see if they are already calculated and cached
		if cached, ok := cache[o.CacheKey]; ok {
			o.FinalPeers = cached
			return
		} else {
			for _, peer := range peers {
				o.addPeerComp(peer)
			}
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
	// TODO: 'smartly' delete those over n so that we have exactly n
	// TODO: test that we have found exactly n peers

	// Sort the final peers by ID
	sort.Strings(o.FinalPeers)

	// We no longer need the actual comps so clear up that space
	o.PeerComps = nil
	o.PeerDists = nil

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
	o.PeerDists = append(o.PeerDists, distanceToPeer)
}
