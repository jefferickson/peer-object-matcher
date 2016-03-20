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
	PeerComps    []PeerComp
	PeerDists    []float64
	FinalPeers   []string
	CacheKey 	 string
}

// To store the distances to other potential peers
type PeerComp struct {
	PeerObject *Object
	Distance   float64
}

// To store cached peer comparisons
type CachedPeerComps struct {
	PeerComps []PeerComp
	PeerDists []float64
}

// Semaphore for map access
var mu sync.Mutex

// Map to store the cache itself
var peerCompCache = make(map[string]CachedPeerComps)

// Factory function to create an object
func NewObject(ID string, Categorical string, NoMatchGroup string, Coords []string) *Object {
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
		CacheKey:	  genCacheKey(Categorical, NoMatchGroup, Coords),
	}
}

// Peer all objects in a group with all other objects in a group
func PeerAllObjects(objects []*Object, n int, counter chan<- bool) {
	for _, object := range objects {
		object.findClosestPeers(objects, n)
		counter <- true
	}
}

// Function to find the closest peers
func (o *Object) findClosestPeers(peers []*Object, n int) {
	// if PeerDists doesn't exist, then we need to calculate them
	if o.PeerComps == nil {
		// Check if this calculation is already in the cache.
		if cached, ok := peerCompCache[o.CacheKey]; ok {
			o.PeerComps = cached.PeerComps
			o.PeerDists = cached.PeerDists
		} else {
			for _, peer := range peers {
				o.addPeerComp(peer)
			}
			// Store into cache for others
			mu.Lock()
			peerCompCache[o.CacheKey] = CachedPeerComps{o.PeerComps, o.PeerDists}
			mu.Unlock()
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
	// TODO: test that we have found exactly n peers
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
