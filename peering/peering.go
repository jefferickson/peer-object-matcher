package peering

import (
	"errors"

	"github.com/jefferickson/peer-object-matcher/object"
)

// Determine distance between two peers
// Or return error if they cannot be peered
func PeerObjects(obj1 *object.Object, obj2 *object.Object, distfn func([]float64, []float64) (float64, error)) (float64, error) {
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
