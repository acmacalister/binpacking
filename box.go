package binpacking

import "fmt"

// FreeSpaceBox represents a rectangular area typically used to track
// available space in packing algorithms.
type FreeSpaceBox struct {
	X      float64 // X-coordinate of the top-left corner of the free space
	Y      float64 // Y-coordinate of the top-left corner of the free space
	Width  float64 // Width of the free space area
	Height float64 // Height of the free space area
}

// Box represents a rectangle with dimensions, position, and packing status.
type Box struct {
	Width             float64 // Width of the box
	Height            float64 // Height of the box
	ConstrainRotation bool    // If true, rotation during packing should be avoided
	X                 float64 // X-coordinate of the top-left corner
	Y                 float64 // Y-coordinate of the top-left corner
	Packed            bool    // Flag indicating if the box has been packed
}

// NewBox creates a new Box instance with specified dimensions and rotation constraint.
// X, Y, and Packed are initialized to their zero values (0, 0, false).
func NewBox(width float64, height float64, constrainRotation bool) *Box {
	return &Box{
		Width:             width,
		Height:            height,
		ConstrainRotation: constrainRotation,
		// X, Y, and Packed default to 0, 0, and false respectively (Go's zero values)
	}
}

// Rotate swaps the Width and Height of the Box.
// This method modifies the receiver Box (b).
func (b *Box) Rotate() {
	// Simple swap in Go
	b.Width, b.Height = b.Height, b.Width
}

// Label returns a formatted string describing the box's dimensions and position.
func (b *Box) Label() string {
	// Use %g which trims trailing zeros for cleaner output
	return fmt.Sprintf("%gx%g at [%g,%g]", b.Width, b.Height, b.X, b.Y)
}

// Area calculates and returns the area of the box (Width * Height).
func (b *Box) Area() float64 {
	return b.Width * b.Height
}
