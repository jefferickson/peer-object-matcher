package object

import (
	"errors"
)

// Determine distance between two peers
// Or return error if they cannot be peered
func peerObjects(obj1 *Object, obj2 *Object, distfn func([]float64, []float64) (float64, error)) (float64, error) {
	distance := 0.0

	// check to make sure the objects are not identical
	if obj1.ID == obj2.ID {
		return distance, errors.New("Objects are identical.")
	}

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
