package binpacking

import "math"

// ScoreBoardEntry holds a potential pairing of a Box with a Bin
// and the calculated Score for that placement.
type ScoreBoardEntry struct {
	Bin   *Bin    // Pointer to the Bin being considered (allows nil)
	Box   *Box    // Pointer to the Box being placed (allows nil)
	Score float64 // Pointer to the calculated Score (allows nil initially, then set by Calculate)
}

// NewScoreBoardEntry creates a new entry linking a Bin and a Box,
// ready for score calculation.
func NewScoreBoardEntry(bin *Bin, box *Box) *ScoreBoardEntry {
	return &ScoreBoardEntry{
		Bin:   bin,
		Box:   box,
		Score: math.MaxFloat64, // Initialize score to indicate no calculation/fit yet
	}
}

// Calculate determines the placement score for the entry's Box within its Bin.
// It calls the associated Bin's ScoreFor method and stores the result internally.
// It returns the calculated Score. If Bin or Box is nil, it returns math.MaxFloat64 and
// sets the internal Score appropriately.
func (sbe *ScoreBoardEntry) Calculate() float64 {
	// Handle cases where Bin or Box might not be set
	if sbe.Bin == nil || sbe.Box == nil {
		sbe.Score = math.MaxFloat64 // Ensure score is set to max value
		return math.MaxFloat64
	}

	// Call the ScoreFor method assumed to exist on the Bin type.
	// This will return math.MaxFloat64 if the box doesn't fit in the bin.
	sbe.Score = sbe.Bin.ScoreFor(sbe.Box)
	return sbe.Score
}

// Fit determines if the calculated score represents a valid placement.
// Returns true if the score is NOT the maximum float value (indicating a fit was found).
func (sbe *ScoreBoardEntry) Fit() bool {
	// A valid fit has a score less than the maximum possible float value.
	return sbe.Score < math.MaxFloat64
}
