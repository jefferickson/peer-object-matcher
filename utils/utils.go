package utils

import (
	"errors"
	"math"
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
