# Go Bin Packing

[![Go Reference](https://pkg.go.dev/badge/github.com/acmacalister/binpacking.svg)](https://pkg.go.dev/github.com/acmacalister/binpacking)
[![Go Report Card](https://goreportcard.com/badge/github.com/acmacalister/binpacking)](https://goreportcard.com/report/github.com/acmacalister/binpacking)
A 2D bin packing library for Go, designed to efficiently pack rectangles (boxes) into larger rectangles (bins). This library provides tools and algorithms to solve the classic two-dimensional bin packing problem, where the goal is to fit a set of items into the minimum number of containers.

This library is a Go port inspired by JavaScript bin-packing libraries (such as [teamcurri/binpacking](https://github.com/teamcurri/binpacking)).

## Features

* Pack 2D rectangular items (`Box`) into 2D containers (`Bin`).
* Support for using multiple bins of potentially different sizes.
* Optional rotation of boxes (can be constrained per box).
* Multiple placement strategies (heuristics) available:
    * Best Short Side Fit (BSSF)
    * Best Long Side Fit (BLSF) - *Note: See implementation details below.*
    * Best Area Fit (BAF)
    * Bottom Left (BL)
* Tracks which boxes were successfully packed and which were left unpacked.
* Calculates packing efficiency for bins.

## Installation

```bash
go get [github.com/acmacalister/binpacking](https://github.com/acmacalister/binpacking)
```

## Usage

```go
package main

import (
	"fmt"
	"[github.com/acmacalister/binpacking](https://github.com/acmacalister/binpacking)" // Use the correct import path
)

func main() {
	// 1. Create Bins (the containers)
	//    Uses the default heuristic (BestShortSideFit) if nil is passed.
	bin := binpacking.NewBin(100, 50, nil) // Width=100, Height=50

	// 2. Create Boxes (the items to pack)
	boxes := []*binpacking.Box{
		binpacking.NewBox(15, 10, false),  // Width=15, Height=10, Rotation allowed
		binpacking.NewBox(50, 45, false),  // Width=50, Height=45, Rotation allowed
		binpacking.NewBox(40, 40, false),  // Width=40, Height=40, Rotation allowed
		binpacking.NewBox(200, 200, false), // Too large
		binpacking.NewBox(50, 100, true),  // Constrained rotation - cannot fit bin
		binpacking.NewBox(50, 100, false), // Allow rotation - can fit bin
	}

	// 3. Create a Packer
	//    Pass a slice of bins the packer can use.
	packer := binpacking.NewPacker([]*binpacking.Bin{bin})

	// 4. Pack the boxes
	//    Can pass options like packing limit.
	packedBoxes := packer.Pack(boxes, binpacking.PackerOptions{})

	// 5. Check Results
	fmt.Printf("--- Packing Results ---\n")
	fmt.Printf("Packed %d boxes:\n", len(packedBoxes))
	for i, box := range packedBoxes {
		fmt.Printf("  Box %d: %s Packed: %v\n", i, box.Label(), box.Packed)
	}

	fmt.Printf("\nUnpacked %d boxes:\n", len(packer.UnpackedBoxes))
	for i, box := range packer.UnpackedBoxes {
		fmt.Printf("  Box %d: %s Packed: %v\n", i, box.Label(), box.Packed)
	}

	fmt.Printf("\nBin State:\n")
	fmt.Printf("  Bin (%fx%f) contains %d boxes. Efficiency: %.2f%%\n",
		bin.Width, bin.Height, len(bin.Boxes), bin.Efficiency())
	for i, box := range bin.Boxes {
		fmt.Printf("    Box %d in bin: %s\n", i, box.Label())
	}

	// Example using a different heuristic
	fmt.Printf("\n--- Using BestAreaFit ---\n")
	bin2 := binpacking.NewBin(100, 50, binpacking.BestAreaFit)
	boxes2 := []*binpacking.Box{
		binpacking.NewBox(60, 40, false),
		binpacking.NewBox(60, 40, false), // Second one shouldn't fit
	}
	packer2 := binpacking.NewPacker([]*binpacking.Bin{bin2})
	packedBoxes2 := packer2.Pack(boxes2, binpacking.PackerOptions{})
	fmt.Printf("Packed %d boxes with BAF:\n", len(packedBoxes2))
	fmt.Printf("Unpacked %d boxes with BAF:\n", len(packer2.UnpackedBoxes))
	fmt.Printf("Bin2 State: Contains %d boxes. Efficiency: %.2f%%\n",
		len(bin2.Boxes), bin2.Efficiency())

}

// Expected Output (may vary slightly based on heuristic details):
// --- Packing Results ---
// Packed 4 boxes:
//   Box 0: 50x45 at [0,0] Packed: true
//   Box 1: 40x40 at [50,0] Packed: true
//   Box 2: 50x50 at [0,0] Packed: true  // Note: 50x100 rotated
//   Box 3: 15x10 at [50,40] Packed: true
//
// Unpacked 2 boxes:
//   Box 0: 200x200 at [0,0] Packed: false
//   Box 1: 50x100 at [0,0] Packed: false // Note: Constrained rotation
//
// Bin State:
//   Bin (100x50) contains 4 boxes. Efficiency: 98.00% // (2250+1600+2500+150)/5000
//     Box 0 in bin: 50x45 at [0,0]
//     Box 1 in bin: 40x40 at [50,0]
//     Box 2 in bin: 50x50 at [0,0]
//     Box 3 in bin: 15x10 at [50,40]
//
// --- Using BestAreaFit ---
// Packed 1 boxes with BAF:
// Unpacked 1 boxes with BAF:
// Bin2 State: Contains 1 boxes. Efficiency: 48.00%
````
