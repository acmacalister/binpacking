package binpacking

// ScoreBoard manages the evaluation of potential placements (ScoreBoardEntry)
// for a set of boxes into a set of bins.
type ScoreBoard struct {
	Entries []*ScoreBoardEntry // All calculated bin/box placement evaluations
	Bins    []*Bin             // The list of available bins
	// Note: Storing the original boxes list might be redundant if CurrentBoxes() is sufficient.
	// Consider if this field is truly needed or if it should be InitialBoxes.
	Boxes []*Box // The initial list of boxes provided
}

// NewScoreBoard creates a new ScoreBoard, initializing entries by calculating
// the score for each initial box against each initial bin.
func NewScoreBoard(bins []*Bin, boxes []*Box) *ScoreBoard {
	sb := &ScoreBoard{
		Entries: make([]*ScoreBoardEntry, 0, len(bins)*len(boxes)), // Pre-allocate slice capacity
		Bins:    bins,
		Boxes:   boxes,
	}

	// Populate initial entries
	for _, bin := range bins {
		sb.addBinEntries(bin, boxes)
	}

	return sb
}

// CurrentBoxes returns a slice containing unique pointers to all boxes
// currently represented in the ScoreBoard entries.
func (sb *ScoreBoard) CurrentBoxes() []*Box {
	// Use a map to efficiently track unique box pointers
	boxSet := make(map[*Box]struct{})
	for _, entry := range sb.Entries {
		if entry != nil && entry.Box != nil {
			boxSet[entry.Box] = struct{}{}
		}
	}

	// Convert the map keys (unique box pointers) back into a slice
	uniqueBoxes := make([]*Box, 0, len(boxSet))
	for boxPtr := range boxSet {
		uniqueBoxes = append(uniqueBoxes, boxPtr)
	}
	return uniqueBoxes
}

// AnyBoxesLeft checks if there are any boxes currently tracked by the scoreboard entries.
func (sb *ScoreBoard) AnyBoxesLeft() bool {
	return len(sb.CurrentBoxes()) > 0
}

// BestFit finds the ScoreBoardEntry representing the best possible placement
// (lowest score) among all entries that indicate a valid fit.
// Returns nil if no fitting placement exists in the current entries.
func (sb *ScoreBoard) BestFit() *ScoreBoardEntry {
	var bestEntry *ScoreBoardEntry = nil // Initialize best to nil

	for _, entry := range sb.Entries {
		// Check if the entry represents a valid fit.
		// Fit() internally checks entry.Score != nil && !entry.Score.IsBlank()
		if entry == nil || !entry.Fit() {
			continue // Skip invalid entries or those that don't fit
		}

		// If this is the first valid entry found, it's the best so far.
		if bestEntry == nil {
			bestEntry = entry
			continue
		}

		// Compare current entry's score value with the best score value found so far.
		if entry.Score < bestEntry.Score {
			bestEntry = entry
		}
	}
	return bestEntry
}

// RemoveBox removes all ScoreBoardEntry instances associated with the specified box
// from the scoreboard.
func (sb *ScoreBoard) RemoveBox(boxToRemove *Box) {
	if boxToRemove == nil {
		return // Nothing to remove
	}
	// Create a new slice to hold entries that don't match the box to remove.
	// Pre-allocate capacity for efficiency.
	filteredEntries := make([]*ScoreBoardEntry, 0, len(sb.Entries))
	for _, entry := range sb.Entries {
		// Keep the entry if its Box pointer is not the one to remove.
		if entry != nil && entry.Box != boxToRemove {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	sb.Entries = filteredEntries // Replace the old slice with the filtered one
}

// AddBin incorporates a new bin into the scoreboard.
// It calculates and adds entries for this new bin against all currently tracked boxes.
func (sb *ScoreBoard) AddBin(bin *Bin) {
	if bin == nil {
		return // Cannot add a nil bin
	}
	sb.Bins = append(sb.Bins, bin) // Add bin to the list of bins
	// Add entries for the new bin against boxes currently in the scoreboard
	sb.addBinEntries(bin, sb.CurrentBoxes())
}

// RecalculateBin updates the scores for all entries associated with a specific bin.
// Useful if the bin's state (e.g., free spaces) has changed.
func (sb *ScoreBoard) RecalculateBin(bin *Bin) {
	if bin == nil {
		return
	}
	for _, entry := range sb.Entries {
		// If the entry belongs to the specified bin, recalculate its score.
		if entry != nil && entry.Bin == bin {
			entry.Calculate()
		}
	}
}

// addBinEntries creates ScoreBoardEntry objects for a given bin and list of boxes,
// calculates their scores, and adds them to the scoreboard's entries.
func (sb *ScoreBoard) addBinEntries(bin *Bin, boxes []*Box) {
	for _, box := range boxes {
		if bin == nil || box == nil {
			continue // Skip nil inputs
		}
		entry := NewScoreBoardEntry(bin, box)
		entry.Calculate() // Calculate the score for this bin/box pair
		sb.Entries = append(sb.Entries, entry)
	}
}
