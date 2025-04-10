package binpacking

import "math"

// ScoreBoardEntry holds a potential pairing of a Box with a Bin
// and the calculated Score for that placement.
type ScoreBoardEntry struct {
	Bin   *Bin  // Pointer to the Bin being considered (allows nil)
	Box   *Box  // Pointer to the Box being placed (allows nil)
	Score int64 // Pointer to the calculated Score (allows nil initially, then set by Calculate)
}

// NewScoreBoardEntry creates a new entry linking a Bin and a Box,
// ready for score calculation.
func NewScoreBoardEntry(bin *Bin, box *Box) *ScoreBoardEntry {
	return &ScoreBoardEntry{
		Bin:   bin,
		Box:   box,
		Score: int64(math.MaxInt64),
	}
}

// Calculate determines the placement score for the entry's Box within its Bin.
// It calls the associated Bin's ScoreFor method and stores the result internally.
// It returns the calculated Score. If Bin or Box is nil, it returns nil and
// sets the internal Score to nil.
func (sbe *ScoreBoardEntry) Calculate() int64 {
	// Handle cases where Bin or Box might not be set
	if sbe.Bin == nil || sbe.Box == nil {
		return int64(math.MaxInt64)
	}

	// Call the ScoreFor method assumed to exist on the Bin type.
	// The ScoreFor method should handle returning an appropriate Score
	// (e.g., NewScore() with MaxInt values if the box doesn't fit in the bin).
	sbe.Score = sbe.Bin.ScoreFor(sbe.Box)
	return sbe.Score
}

// Fit determines if the calculated score represents a valid placement.
// Returns true if the score is NOT the maximum value (indicating a fit was found).
func (sbe *ScoreBoardEntry) Fit() bool {
	// A valid fit has a score less than the maximum possible integer value.
	return sbe.Score < int64(math.MaxInt64)
}
