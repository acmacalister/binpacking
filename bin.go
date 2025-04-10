package binpacking

import (
	"fmt"
)

// Bin represents a container for packing boxes.
type Bin struct {
	Width      int64
	Height     int64
	Boxes      []*Box                // Boxes placed in this bin
	Placement  PlacementStrategyFunc // Strategy used for finding placement positions
	FreeSpaces []*FreeSpaceBox       // List of available free rectangles
}

// NewBin creates a new Bin instance.
func NewBin(width int64, height int64, placement PlacementStrategyFunc) *Bin {
	// Initialize FreeSpaces with one rectangle covering the entire bin
	initialFreeSpace := FreeSpaceBox{Width: width, Height: height}

	if placement == nil {
		placement = BestShortSideFit // Assuming this is desired default
	}

	return &Bin{
		Width:      width,
		Height:     height,
		Boxes:      make([]*Box, 0),
		Placement:  placement,
		FreeSpaces: []*FreeSpaceBox{&initialFreeSpace}, // Start with one large free space
	}
}

// Area returns the total area of the bin.
func (b *Bin) Area() int64 {
	return b.Width * b.Height
}

// Efficiency calculates the percentage of the bin's area occupied by packed boxes.
// Returns a float64 between 0 and 100.
func (b *Bin) Efficiency() float64 {
	boxesArea := int64(0)
	for _, box := range b.Boxes {
		boxesArea += box.Area()
	}
	binArea := b.Area()
	if binArea == 0 {
		return 0.0 // Avoid division by zero
	}
	// Use float64 for calculation to get percentage
	return (float64(boxesArea) * 100.0) / float64(binArea)
}

// Label returns a string representation of the bin including dimensions and efficiency.
func (b *Bin) Label() string {
	// %.2f formats the float with 2 decimal places
	return fmt.Sprintf("%dx%d %.2f%%", b.Width, b.Height, b.Efficiency())
}

// Insert attempts to place a box into the bin using the bin's heuristic.
// It updates the bin's state (Boxes, FreeSpaces) if successful.
// Returns true if the box was successfully packed, false otherwise.
func (b *Bin) Insert(box *Box) bool {
	if box.Packed {
		return false
	}

	placement := FindBestPlacement(box, b.FreeSpaces, b.Placement)

	if !placement.Fits {
		return false // No suitable placement found
	}

	// Apply placement
	box.X = placement.X
	box.Y = placement.Y
	box.Packed = true
	if placement.NeedsRotation {
		box.Rotate()
	}

	// Split the chosen free space
	newFreeSpaces := make([]*FreeSpaceBox, 0, len(b.FreeSpaces)+3) // Estimate capacity

	for i := 0; i < len(b.FreeSpaces); i++ {
		currentFreeSpace := b.FreeSpaces[i]
		if currentFreeSpace == placement.ChosenSpace {
			// Split this node, potentially adding 0-4 new nodes directly
			// The split function should probably return the new nodes instead of modifying b.FreeSpaces directly
			generatedSpaces := b.generateSplits(currentFreeSpace, box) // New helper needed
			newFreeSpaces = append(newFreeSpaces, generatedSpaces...)
		} else {
			// Keep non-chosen, non-split nodes
			newFreeSpaces = append(newFreeSpaces, currentFreeSpace)
		}
	}

	b.FreeSpaces = newFreeSpaces
	b.pruneFreeList()
	b.Boxes = append(b.Boxes, box)

	return true
}

// ScoreFor simulates placing the box and returns the score without modifying the bin.
// It creates a copy of the box to avoid side effects.
func (b *Bin) ScoreFor(box *Box) int64 {
	// Create a copy to pass to the placement strategy, so the original box isn't modified.
	// Assumes NewBox creates a clean copy with dimensions and rotation constraint.
	copyBox := NewBox(box.Width, box.Height, box.ConstrainRotation)
	// The placement will find the position but won't modify the original box or bin state.
	placement := FindBestPlacement(copyBox, b.FreeSpaces, b.Placement)
	return placement.Score
}

// IsLargerThan checks if the bin is large enough to potentially hold the box
// (considering rotation if allowed by the box).
func (b *Bin) IsLargerThan(box *Box) bool {
	canFitOriginal := b.Width >= box.Width && b.Height >= box.Height
	canFitRotated := !box.ConstrainRotation && b.Height >= box.Width && b.Width >= box.Height
	return canFitOriginal || canFitRotated
}

// Helper to generate splits without modifying the list directly during split logic.
func (b *Bin) generateSplits(freeNode *FreeSpaceBox, usedNode *Box) []*FreeSpaceBox {
	// Based on your original split logic, but appends to a local slice instead of b.FreeSpaces
	splits := make([]*FreeSpaceBox, 0, 4)

	// Separating Axis Theorem (SAT) intersection test.
	if usedNode.X >= freeNode.X+freeNode.Width ||
		usedNode.X+usedNode.Width <= freeNode.X ||
		usedNode.Y >= freeNode.Y+freeNode.Height ||
		usedNode.Y+usedNode.Height <= freeNode.Y {
		// Should not happen if called on the chosen node, but check anyway
		return splits // Return empty slice
	}

	// Try vertical splits (Top/Bottom)
	if usedNode.X < freeNode.X+freeNode.Width && usedNode.X+usedNode.Width > freeNode.X {
		// Top
		if usedNode.Y > freeNode.Y {
			newNode := *freeNode
			newNode.Height = usedNode.Y - newNode.Y
			splits = append(splits, &newNode)
		}
		// Bottom
		usedBottomY := usedNode.Y + usedNode.Height
		freeBottomY := freeNode.Y + freeNode.Height
		if usedBottomY < freeBottomY {
			newNode := *freeNode
			newNode.Y = usedBottomY
			newNode.Height = freeBottomY - usedBottomY
			splits = append(splits, &newNode)
		}
	}

	// Try horizontal splits (Left/Right)
	if usedNode.Y < freeNode.Y+freeNode.Height && usedNode.Y+usedNode.Height > freeNode.Y {
		// Left
		if usedNode.X > freeNode.X {
			newNode := *freeNode
			newNode.Width = usedNode.X - newNode.X
			splits = append(splits, &newNode)
		}
		// Right
		usedRightX := usedNode.X + usedNode.Width
		freeRightX := freeNode.X + freeNode.Width
		if usedRightX < freeRightX {
			newNode := *freeNode
			newNode.X = usedRightX
			newNode.Width = freeRightX - usedRightX
			splits = append(splits, &newNode)
		}
	}

	return splits
}

// pruneFreeList removes redundant free spaces (those fully contained within another).
func (b *Bin) pruneFreeList() {
	// Create a new list to store non-contained free spaces.
	// Pre-allocate capacity close to original for efficiency.
	prunedList := make([]*FreeSpaceBox, 0, len(b.FreeSpaces))

	for i := 0; i < len(b.FreeSpaces); i++ {
		rectA := b.FreeSpaces[i]
		isContained := false

		// Check if rectA is contained within any *other* rectangle
		for j := 0; j < len(b.FreeSpaces); j++ {
			if i == j {
				continue // Don't compare with self
			}
			rectB := b.FreeSpaces[j]
			if b.isContainedIn(rectA, rectB) {
				isContained = true
				break // Found a container, no need to check further
			}
		}

		// If rectA was not contained in any other rectangle, keep it.
		if !isContained {
			prunedList = append(prunedList, rectA)
		}
	}

	b.FreeSpaces = prunedList
}

// isContainedIn checks if rectA is fully contained within rectB.
func (b *Bin) isContainedIn(rectA, rectB *FreeSpaceBox) bool {
	// Basic nil check for safety, although unlikely if called from pruneFreeList
	if rectA == nil || rectB == nil {
		return false
	}
	return rectA.X >= rectB.X &&
		rectA.Y >= rectB.Y &&
		rectA.X+rectA.Width <= rectB.X+rectB.Width &&
		rectA.Y+rectA.Height <= rectB.Y+rectB.Height
}
