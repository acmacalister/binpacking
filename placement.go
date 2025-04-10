package binpacking

import (
	"math"
)

// PlacementInfo holds the details about the best placement found for a box
// within a set of free spaces, according to a specific placement strategy.
type PlacementInfo struct {
	// Score represents the quality of the placement, calculated by a PlacementStrategyFunc.
	// Lower scores generally indicate better fits. A score of math.MaxInt64 indicates no fit.
	Score int64
	// ChosenSpace is a pointer to the specific FreeSpaceBox where the placement should occur.
	// Will be nil if Fits is false.
	ChosenSpace *FreeSpaceBox
	// X is the calculated horizontal coordinate for the top-left corner of the box.
	X int64
	// Y is the calculated vertical coordinate for the top-left corner of the box.
	Y int64
	// NeedsRotation indicates whether the box's width and height should be swapped for this placement.
	NeedsRotation bool
	// Fits indicates whether a suitable placement satisfying the strategy was found.
	Fits bool
}

// PlacementStrategyFunc defines the signature for functions that calculate a score
// indicating how well a rectangle of given dimensions fits into a specific FreeSpaceBox.
// Lower scores are considered better fits. The returned score is a single int64 value.
type PlacementStrategyFunc func(freeSpace *FreeSpaceBox, rectWidth, rectHeight int64) int64

// FindBestPlacement iterates through available free spaces to find the best possible
// position for a given Box, according to the provided PlacementStrategyFunc.
// It considers both original and rotated orientations (if allowed by the box).
//
// Parameters:
//
//	box: The Box to find a placement for.
//	freeSpaces: A slice of available FreeSpaceBox areas to check.
//	placement: The PlacementStrategyFunc used to score potential fits.
//
// Returns:
//
//	A PlacementInfo struct containing details of the best fit found.
//	If no fit is possible, PlacementInfo.Fits will be false and Score will be math.MaxInt64.
func FindBestPlacement(box *Box, freeSpaces []*FreeSpaceBox, placement PlacementStrategyFunc) PlacementInfo {
	// Initialize with worst possible score and Fits=false
	bestInfo := PlacementInfo{Score: math.MaxInt64, Fits: false}

	for _, freeSpace := range freeSpaces {
		// Try placing the box in its original orientation
		if freeSpace.Width >= box.Width && freeSpace.Height >= box.Height {
			score := placement(freeSpace, box.Width, box.Height)
			// If this placement is better than the best found so far
			if score < bestInfo.Score {
				bestInfo = PlacementInfo{
					Score:         score,
					ChosenSpace:   freeSpace,
					X:             freeSpace.X, // Place at the top-left corner of the chosen free space
					Y:             freeSpace.Y,
					NeedsRotation: false,
					Fits:          true,
				}
			}
		}

		// Try placing the box in its rotated orientation, if allowed and different dimensions
		if !box.ConstrainRotation && box.Width != box.Height && freeSpace.Width >= box.Height && freeSpace.Height >= box.Width {
			// Calculate score using rotated dimensions
			score := placement(freeSpace, box.Height, box.Width)
			// If this placement is better than the best found so far
			if score < bestInfo.Score {
				bestInfo = PlacementInfo{
					Score:         score,
					ChosenSpace:   freeSpace,
					X:             freeSpace.X, // Place at the top-left corner of the chosen free space
					Y:             freeSpace.Y,
					NeedsRotation: true, // Mark that rotation is needed
					Fits:          true,
				}
			}
		}
	}
	// Return the best placement found (or initial state if no fit)
	return bestInfo
}

// BestAreaFit implements the PlacementStrategyFunc interface.
// It scores placements by minimizing the leftover area in the free space after placing
// the rectangle. As a tie-breaker, it adds the 'short side fit' (the smaller
// of the horizontal or vertical leftover dimensions). Lower scores are better.
func BestAreaFit(freeSpace *FreeSpaceBox, rectWidth, rectHeight int64) int64 {
	areaFit := freeSpace.Width*freeSpace.Height - rectWidth*rectHeight
	leftOverHoriz := abs(freeSpace.Width - rectWidth)
	leftOverVert := abs(freeSpace.Height - rectHeight)
	shortSideFit := min(leftOverHoriz, leftOverVert)
	// Combine area fit and short side fit into a single score
	return areaFit + shortSideFit
}

// BestShortSideFit implements the PlacementStrategyFunc interface.
// It scores placements by minimizing the sum of the leftover dimensions
// (horizontal gap + vertical gap) in the free space. Lower scores are better.
// Note: This differs from some BSSF implementations that prioritize minimizing the
// smaller gap first, then the larger gap as a tie-breaker (lexicographical score).
func BestShortSideFit(freeSpace *FreeSpaceBox, rectWidth, rectHeight int64) int64 {
	leftOverHoriz := abs(freeSpace.Width - rectWidth)
	leftOverVert := abs(freeSpace.Height - rectHeight)
	// Return the sum of the horizontal and vertical gaps
	return leftOverHoriz + leftOverVert
}

// BestLongSideFit implements the PlacementStrategyFunc interface.
// It scores placements primarily by minimizing the larger of the leftover dimensions
// (the "long side fit") in the free space. Lower scores are better.
// Note: Due to the single int64 return type limitation, the secondary tie-breaker
// (minimizing the short side fit) cannot be directly incorporated into the score
// for lexicographical comparison. This implementation returns only the long side fit value.
func BestLongSideFit(freeSpace *FreeSpaceBox, rectWidth, rectHeight int64) int64 {
	leftOverHoriz := abs(freeSpace.Width - rectWidth)
	leftOverVert := abs(freeSpace.Height - rectHeight)
	// Return the larger gap (long side fit) as the score.
	return max(leftOverHoriz, leftOverVert)
}

// BottomLeft implements the PlacementStrategyFunc interface.
// It scores placements based on a combination of the free space's top-left corner (X, Y)
// and the height of the rectangle being placed. It aims to minimize Y + X + rectHeight.
// Lower scores indicate preferred placements (lower, then left-er, considering height).
func BottomLeft(freeSpace *FreeSpaceBox, rectWidth, rectHeight int64) int64 {
	// Score prioritizes lower Y, then lower X, then lower rectangle height?
	return freeSpace.Y + freeSpace.X + rectHeight
}

// abs returns the absolute value of x.
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// min returns the smaller of x or y.
func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// max returns the larger of x or y.
func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
