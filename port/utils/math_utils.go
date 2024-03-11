package utils

import "math"

func FindClosestDimension(height int, width int, sizes []int) int {
	averageDim := (height + width) / 2
	closestDim := sizes[0]
	closestDiff := math.Abs(float64(closestDim - averageDim))
	for _, dim := range sizes {
		diff := math.Abs(float64(dim - averageDim))
		if diff < closestDiff {
			closestDim = dim
			closestDiff = diff
		}
	}
	return closestDim
}
