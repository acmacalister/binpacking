package binpacking

// PackerOptions defines optional parameters for the packing process.
type PackerOptions struct {
	// Limit specifies the maximum number of boxes to pack.
	// If zero or negative, packing continues until no more boxes fit
	// or all boxes are packed.
	Limit int64
}

// Packer orchestrates the bin packing process by coordinating
// bins, boxes, and the scoreboard evaluating potential fits.
type Packer struct {
	Bins          []*Bin // Bins available for packing. Owned/managed by the Packer instance.
	UnpackedBoxes []*Box // Boxes that could not be packed in the last call to Pack.
}

// NewPacker creates a new Packer instance with a given set of initial bins.
// It takes ownership of the provided bin slice.
func NewPacker(bins []*Bin) *Packer {
	// Create a copy of the bins slice if you want the packer to manage its own list
	// independent of the caller after creation. Otherwise, just assign.
	packerBins := make([]*Bin, len(bins))
	copy(packerBins, bins) // Creates a new slice with copies of pointers

	return &Packer{
		Bins:          packerBins,
		UnpackedBoxes: make([]*Box, 0), // Initialize as empty slice
	}
}

// Pack attempts to pack the given boxes into the packer's bins using a best-fit strategy.
//
// Args:
//
//	boxes: A slice of Box pointers to attempt packing. Boxes marked as Packed=true are skipped.
//	options: PackerOptions allowing specification of limits, etc.
//
// Returns:
//
//	A slice containing pointers to the boxes that were successfully packed in this run.
//
// Note: This method updates the Packer's UnpackedBoxes field with boxes that could not be placed.
func (p *Packer) Pack(boxes []*Box, options PackerOptions) []*Box {
	packedBoxes := make([]*Box, 0)
	p.UnpackedBoxes = make([]*Box, 0) // Clear previous unpacked boxes from this packer instance

	// 1. Filter out nil boxes and those already marked as packed.
	boxesToPack := make([]*Box, 0, len(boxes))
	for _, box := range boxes {
		if box != nil && !box.Packed {
			boxesToPack = append(boxesToPack, box)
		}
	}

	// Return early if no boxes need packing.
	if len(boxesToPack) == 0 {
		return packedBoxes
	}

	// 2. Determine packing limit
	limit := options.Limit
	useLimit := limit > 0 // Only use the limit if it's positive

	// 3. Set up the ScoreBoard.
	// Use the packer's current set of bins and the filtered list of boxes.
	board := NewScoreBoard(p.Bins, boxesToPack)

	// 4. Main packing loop: Continues as long as a best fit can be found.
	for {
		bestEntry := board.BestFit()

		// If BestFit returns nil, no more boxes can be placed in any bin.
		if bestEntry == nil {
			break // Exit the packing loop
		}

		// Safeguard: Ensure the best entry has valid Bin and Box pointers.
		// This condition should ideally not be met if ScoreBoard is working correctly.
		if bestEntry.Bin == nil || bestEntry.Box == nil {
			// Log this anomaly if possible.
			// fmt.Printf("Warning: BestFit returned entry with nil Bin (%v) or Box (%v)\n", bestEntry.Bin == nil, bestEntry.Box == nil)
			// Attempt to remove the problematic box (if identifiable) from the board to prevent infinite loops.
			if bestEntry.Box != nil {
				board.RemoveBox(bestEntry.Box)
			} else {
				// If the box itself is nil, we can't easily remove it; break to avoid potential issues.
				break
			}
			continue // Try finding the next best fit
		}

		// Attempt to insert the chosen box into the chosen bin.
		// The Bin.Insert method should handle the actual placement logic,
		// update the bin's free spaces, and mark the box as Packed = true.
		inserted := bestEntry.Bin.Insert(bestEntry.Box)

		// If insertion failed (e.g., bin state changed concurrently, strategy strategy inconsistency),
		// remove the box from consideration to avoid potential infinite loops.
		if !inserted {
			board.RemoveBox(bestEntry.Box)
			continue // Try the next best fit
		}

		// Add the successfully placed box to the list of packed boxes for this run.
		packedBoxes = append(packedBoxes, bestEntry.Box)

		// Remove the now-packed box from the ScoreBoard so it's not considered again.
		board.RemoveBox(bestEntry.Box)

		// Recalculate scores for the bin that was just modified, as its free spaces changed.
		board.RecalculateBin(bestEntry.Bin)

		// Check if the packing limit has been reached.
		if useLimit && int64(len(packedBoxes)) >= limit {
			break // Exit loop if limit reached
		}
	} // End packing loop

	// 5. Determine which boxes remain unpacked.
	// These are the boxes still present on the scoreboard after the loop.
	// Note: CurrentBoxes() returns unique boxes remaining in scoreboard entries.
	p.UnpackedBoxes = board.CurrentBoxes()

	return packedBoxes
}
