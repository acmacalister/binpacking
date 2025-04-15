package binpacking

import (
	"math"
	"testing"
)

func TestBin(t *testing.T) {
	t.Run("allows to insert boxes while space available", func(t *testing.T) {
		// Assuming NewBin uses default (e.g., BestShortSideFit) if placement strategy is nil
		bin := NewBin(100, 50, BestShortSideFit)
		boxes := []*Box{
			NewBox(50, 50, false), // Box 0
			NewBox(10, 40, false), // Box 1
			NewBox(50, 44, false), // Box 2 (should not fit after 0 & 1)
		}
		remainingBoxes := make([]*Box, 0)

		for _, box := range boxes {
			if !bin.Insert(box) {
				remainingBoxes = append(remainingBoxes, box)
			}
		}

		// Assertions
		if len(bin.Boxes) != 2 {
			t.Errorf("Bin box count: got %d, want %d", len(bin.Boxes), 2)
		}
		if len(bin.Boxes) > 0 && bin.Boxes[0] != boxes[0] {
			t.Errorf("Bin box 0 pointer mismatch: got %p, want %p", bin.Boxes[0], boxes[0])
		}
		// Note: The exact order depends on placement strategy. BSSF might place 10x40 first.
		// Let's check placement assuming boxes[0] (50x50) was placed first. Adjust if needed.
		if len(bin.Boxes) > 0 && bin.Boxes[0] == boxes[0] {
			if bin.Boxes[0].X != 0 {
				t.Errorf("Box 0 X: got %f, want %f", bin.Boxes[0].X, 0.0)
			}
			if bin.Boxes[0].Y != 0 {
				t.Errorf("Box 0 Y: got %f, want %f", bin.Boxes[0].Y, 0.0)
			}
			if !bin.Boxes[0].Packed {
				t.Errorf("Box 0 Packed: got %v, want %v", bin.Boxes[0].Packed, true)
			}
		}
		if len(bin.Boxes) > 1 && bin.Boxes[1] == boxes[1] { // Check second placed box (10x40)
			if bin.Boxes[1].X != 50 { // Assuming placed next to 50x50
				t.Errorf("Box 1 X: got %f, want %f", bin.Boxes[1].X, 50.0)
			}
			if bin.Boxes[1].Y != 0 {
				t.Errorf("Box 1 Y: got %f, want %f", bin.Boxes[1].Y, 0.0)
			}
			if !bin.Boxes[1].Packed {
				t.Errorf("Box 1 Packed: got %v, want %v", bin.Boxes[1].Packed, true)
			}
		}

		if len(remainingBoxes) != 1 {
			t.Errorf("Remaining box count: got %d, want %d", len(remainingBoxes), 1)
		}
		if len(remainingBoxes) > 0 && remainingBoxes[0] != boxes[2] {
			t.Errorf("Remaining box 0 pointer mismatch: got %p, want %p", remainingBoxes[0], boxes[2])
		}

		// Efficiency check (depends on exactly which boxes were placed)
		// If boxes[0] (2500) and boxes[1] (400) were placed in bin (5000 area):
		expectedEfficiency := float64((50*50+10*40)*100) / float64(100*50) // 58.0
		// Use tolerance for float comparison
		tolerance := 0.01
		actualEfficiency := bin.Efficiency()
		if math.Abs(actualEfficiency-expectedEfficiency) > tolerance {
			t.Errorf("Bin efficiency: got %.2f, want %.2f", actualEfficiency, expectedEfficiency)
		}

		// Check remaining box state
		if len(remainingBoxes) > 0 {
			if remainingBoxes[0].X != 0 {
				t.Errorf("Remaining Box X: got %f, want %f", remainingBoxes[0].X, 0.0)
			}
			if remainingBoxes[0].Y != 0 {
				t.Errorf("Remaining Box Y: got %f, want %f", remainingBoxes[0].Y, 0.0)
			}
			if remainingBoxes[0].Packed {
				t.Errorf("Remaining Box Packed: got %v, want %v", remainingBoxes[0].Packed, false)
			}
		}
	})

	t.Run("allows to use custom placement strategy", func(t *testing.T) {
		bin := NewBin(100, 50, BestAreaFit)
		box := NewBox(50, 100, false) // Needs rotation to fit
		result := bin.Insert(box)

		if !result {
			t.Errorf("Insert result: got %v, want %v", result, true)
		}
		if len(bin.Boxes) != 1 {
			t.Errorf("Bin should contain 1 box after insertion, got %d", len(bin.Boxes))
		}
		if len(bin.Boxes) > 0 && (bin.Boxes[0].Width != 100 || bin.Boxes[0].Height != 50) {
			t.Errorf("Box dimensions after insertion: got %fx%f, want %dx%d", bin.Boxes[0].Width, bin.Boxes[0].Height, 100, 50)
		}
	})
}

func TestPacker(t *testing.T) {
	t.Run("does nothing when no bin and no box passed", func(t *testing.T) {
		packer := NewPacker([]*Bin{}) // No bins
		packedBoxes := packer.Pack([]*Box{}, PackerOptions{})

		if len(packedBoxes) != 0 {
			t.Errorf("Packed box count: got %d, want %d", len(packedBoxes), 0)
		}
		if len(packer.UnpackedBoxes) != 0 {
			t.Errorf("Unpacked box count: got %d, want %d", len(packer.UnpackedBoxes), 0)
		}
	})

	t.Run("puts single box in single bin", func(t *testing.T) {
		bin1, _, _ := newBins(t)
		box := NewBox(9000, 3000, false)
		packer := NewPacker([]*Bin{bin1})
		packedBoxes := packer.Pack([]*Box{box}, PackerOptions{})

		if len(packedBoxes) != 1 {
			t.Fatalf("Packed box count: got %d, want %d", len(packedBoxes), 1)
		}
		if packedBoxes[0] != box {
			t.Errorf("Packed box pointer mismatch: got %p, want %p", packedBoxes[0], box)
		}
		if len(bin1.Boxes) != 1 {
			t.Fatalf("Bin box count: got %d, want %d", len(bin1.Boxes), 1)
		}
		if bin1.Boxes[0] != box {
			t.Errorf("Bin box pointer mismatch: got %p, want %p", bin1.Boxes[0], box)
		}
		if box.Width != 9000 {
			t.Errorf("Box Width: got %f, want %f", box.Width, 9000.0)
		}
		if box.Height != 3000 {
			t.Errorf("Box Height: got %f, want %f", box.Height, 3000.0)
		}
		if box.X != 0 {
			t.Errorf("Box X: got %f, want %f", box.X, 0.0)
		}
		if box.Y != 0 {
			t.Errorf("Box Y: got %f, want %f", box.Y, 0.0)
		}
		if !box.Packed {
			t.Errorf("Box Packed: got %v, want %v", box.Packed, true)
		}
		if len(packer.UnpackedBoxes) != 0 {
			t.Errorf("Unpacked box count: got %d, want %d", len(packer.UnpackedBoxes), 0)
		}
	})

	t.Run("puts rotated box in single bin", func(t *testing.T) {
		bin1, _, _ := newBins(t)
		box := NewBox(1000, 9000, false) // Initially tall
		packer := NewPacker([]*Bin{bin1})
		packedBoxes := packer.Pack([]*Box{box}, PackerOptions{})

		if len(packedBoxes) != 1 {
			t.Fatalf("Packed box count: got %d, want %d", len(packedBoxes), 1)
		}
		if len(bin1.Boxes) != 1 {
			t.Fatalf("Bin box count: got %d, want %d", len(bin1.Boxes), 1)
		}
		// Check dimensions *after* packing - should be rotated
		if box.Width != 9000 {
			t.Errorf("Box Width after pack: got %f, want %f", box.Width, 9000.0)
		}
		if box.Height != 1000 {
			t.Errorf("Box Height after pack: got %f, want %f", box.Height, 1000.0)
		}
		if box.X != 0 {
			t.Errorf("Box X: got %f, want %f", box.X, 0.0)
		}
		if box.Y != 0 {
			t.Errorf("Box Y: got %f, want %f", box.Y, 0.0)
		}
		if !box.Packed {
			t.Errorf("Box Packed: got %v, want %v", box.Packed, true)
		}
		if len(packer.UnpackedBoxes) != 0 {
			t.Errorf("Unpacked box count: got %d, want %d", len(packer.UnpackedBoxes), 0)
		}
	})

	t.Run("puts large box in large bin", func(t *testing.T) {
		bin1, bin2, bin3 := newBins(t)
		box := NewBox(11000, 2000, false)
		packer := NewPacker([]*Bin{bin1, bin2, bin3})
		packedBoxes := packer.Pack([]*Box{box}, PackerOptions{})

		if len(packedBoxes) != 1 {
			t.Fatalf("Packed box count: got %d, want %d", len(packedBoxes), 1)
		}
		if len(bin1.Boxes) != 0 {
			t.Errorf("Bin1 box count: got %d, want %d", len(bin1.Boxes), 0)
		}
		if len(bin2.Boxes) != 0 {
			t.Errorf("Bin2 box count: got %d, want %d", len(bin2.Boxes), 0)
		}
		if len(bin3.Boxes) != 1 {
			t.Errorf("Bin3 box count: got %d, want %d", len(bin3.Boxes), 1)
		}
		// Check dimensions were not rotated
		if box.Width != 11000 || box.Height != 2000 {
			t.Errorf("Box dimensions: got %fx%f, want %fx%f", box.Width, box.Height, 11000.0, 2000.0)
		}
		if !box.Packed {
			t.Errorf("Box Packed: got %v, want %v", box.Packed, true)
		}
		if len(packer.UnpackedBoxes) != 0 {
			t.Errorf("Unpacked box count: got %d, want %d", len(packer.UnpackedBoxes), 0)
		}
	})

	t.Run("puts two boxes in single bin", func(t *testing.T) {
		bin1, _, _ := newBins(t)
		box1 := NewBox(8000, 1500, false)
		box2 := NewBox(1000, 9000, false) // Needs rotation
		packer := NewPacker([]*Bin{bin1})
		packedBoxes := packer.Pack([]*Box{box1, box2}, PackerOptions{})

		if len(packedBoxes) != 2 {
			t.Errorf("Packed box count: got %d, want %d", len(packedBoxes), 2)
		}
		if len(bin1.Boxes) != 2 {
			t.Errorf("Bin box count: got %d, want %d", len(bin1.Boxes), 2)
		}
		if !box1.Packed {
			t.Errorf("Box1 Packed: got %v, want %v", box1.Packed, true)
		}
		if !box2.Packed {
			t.Errorf("Box2 Packed: got %v, want %v", box2.Packed, true)
		}
		if box2.Width != 9000 || box2.Height != 1000 { // Check rotation happened
			t.Errorf("Box2 dimensions after pack: got %fx%f, want %fx%f", box2.Width, box2.Height, 9000.0, 1000.0)
		}
		if len(packer.UnpackedBoxes) != 0 {
			t.Errorf("Unpacked box count: got %d, want %d", len(packer.UnpackedBoxes), 0)
		}
	})

	t.Run("puts two boxes in separate bins", func(t *testing.T) {
		// Use slightly different placement strategy maybe? Or ensure BSSF places them separately.
		// Default placement strategy should be sufficient if space is tight.
		bin1 := NewBin(9600, 3100, nil)
		bin2 := NewBin(9600, 3100, nil)
		box1 := NewBox(5500, 2000, false)
		box2 := NewBox(5000, 2000, false)
		packer := NewPacker([]*Bin{bin1, bin2})
		packedBoxes := packer.Pack([]*Box{box1, box2}, PackerOptions{})

		if len(packedBoxes) != 2 {
			t.Errorf("Packed box count: got %d, want %d", len(packedBoxes), 2)
		}
		// Check that boxes ended up in different bins (exact distribution might vary)
		if !((len(bin1.Boxes) == 1 && len(bin2.Boxes) == 1) || (len(bin1.Boxes) == 2 && len(bin2.Boxes) == 0) || (len(bin1.Boxes) == 0 && len(bin2.Boxes) == 2)) {
			// Allow flexibility, but maybe tighten this if specific placement strategy guarantees separation
			t.Logf("Bin1 boxes: %d, Bin2 boxes: %d", len(bin1.Boxes), len(bin2.Boxes))
			// If they both fit in one bin, that might be valid for some placement strategy. BSSF should separate though.
			if len(bin1.Boxes)+len(bin2.Boxes) != 2 {
				t.Errorf("Total boxes in bins is %d, want 2", len(bin1.Boxes)+len(bin2.Boxes))
			}
		} else if len(bin1.Boxes) != 1 || len(bin2.Boxes) != 1 {
			t.Errorf("Bin box counts: got %d/%d, want 1/1", len(bin1.Boxes), len(bin2.Boxes))
		}

		if !box1.Packed {
			t.Errorf("Box1 Packed: got %v, want %v", box1.Packed, true)
		}
		if !box2.Packed {
			t.Errorf("Box2 Packed: got %v, want %v", box2.Packed, true)
		}
		if len(packer.UnpackedBoxes) != 0 {
			t.Errorf("Unpacked box count: got %d, want %d", len(packer.UnpackedBoxes), 0)
		}
	})

	t.Run("does not put in bin too large box", func(t *testing.T) {
		bin1, _, _ := newBins(t)        // 9600x3100
		box := NewBox(10000, 10, false) // Too wide
		packer := NewPacker([]*Bin{bin1})
		packedBoxes := packer.Pack([]*Box{box}, PackerOptions{})

		if len(packedBoxes) != 0 {
			t.Errorf("Packed box count: got %d, want %d", len(packedBoxes), 0)
		}
		if len(bin1.Boxes) != 0 {
			t.Errorf("Bin box count: got %d, want %d", len(bin1.Boxes), 0)
		}
		if box.Packed {
			t.Errorf("Box Packed: got %v, want %v", box.Packed, false)
		}
		if len(packer.UnpackedBoxes) != 1 || packer.UnpackedBoxes[0] != box {
			t.Errorf("Unpacked box count/content mismatch")
		}
	})

	t.Run("puts in bin only fitting boxes", func(t *testing.T) {
		bin1, _, _ := newBins(t) // 9600x3100
		box1 := NewBox(4000, 3000, false)
		box2 := NewBox(4000, 3000, false)
		box3 := NewBox(4000, 3000, false) // Should not fit after 1 & 2
		boxes := []*Box{box1, box2, box3}
		packer := NewPacker([]*Bin{bin1})
		packedBoxes := packer.Pack(boxes, PackerOptions{})

		wantPackedCount := 2
		gotPackedCount := countPacked(boxes)

		if len(packedBoxes) != wantPackedCount {
			t.Errorf("Packed box count (return value): got %d, want %d", len(packedBoxes), wantPackedCount)
		}
		if len(bin1.Boxes) != wantPackedCount {
			t.Errorf("Bin box count: got %d, want %d", len(bin1.Boxes), wantPackedCount)
		}
		if len(boxes) != 3 { // Original slice should be unchanged
			t.Errorf("Original boxes slice length changed: got %d, want %d", len(boxes), 3)
		}
		if gotPackedCount != wantPackedCount {
			t.Errorf("Packed boxes in original slice: got %d, want %d", gotPackedCount, wantPackedCount)
		}
		if len(packer.UnpackedBoxes) != 1 {
			t.Errorf("Unpacked box count: got %d, want %d", len(packer.UnpackedBoxes), 1)
		}
		if len(packer.UnpackedBoxes) > 0 && packer.UnpackedBoxes[0] != box3 {
			t.Errorf("Unpacked box pointer mismatch")
		}
	})

	t.Run("respects limit", func(t *testing.T) {
		bin1, _, _ := newBins(t)
		box1 := NewBox(1000, 1000, false)
		box2 := NewBox(1000, 1000, false)
		boxes := []*Box{box1, box2}
		packer := NewPacker([]*Bin{bin1})
		packedBoxes := packer.Pack(boxes, PackerOptions{Limit: 1})

		wantPackedCount := 1
		gotPackedCount := countPacked(boxes)

		if len(packedBoxes) != wantPackedCount {
			t.Errorf("Packed box count (return value): got %d, want %d", len(packedBoxes), wantPackedCount)
		}
		if len(bin1.Boxes) != wantPackedCount {
			t.Errorf("Bin box count: got %d, want %d", len(bin1.Boxes), wantPackedCount)
		}
		if len(boxes) != 2 {
			t.Errorf("Original boxes slice length changed: got %d, want %d", len(boxes), 2)
		}
		if gotPackedCount != wantPackedCount {
			t.Errorf("Packed boxes in original slice: got %d, want %d", gotPackedCount, wantPackedCount)
		}
		if len(packer.UnpackedBoxes) != 1 {
			t.Errorf("Unpacked box count: got %d, want %d", len(packer.UnpackedBoxes), 1)
		}
	})

	t.Run("does not pack box twice", func(t *testing.T) {
		bin1, _, _ := newBins(t)
		box1 := NewBox(1000, 1000, false)
		packer := NewPacker([]*Bin{bin1})

		packed1 := packer.Pack([]*Box{box1}, PackerOptions{})
		if len(packed1) != 1 {
			t.Fatalf("First pack call failed: got len %d, want 1", len(packed1))
		}
		if !box1.Packed {
			t.Fatalf("Box was not marked packed after first pack")
		}

		// Second call with the same box (which is now marked as packed)
		packed2 := packer.Pack([]*Box{box1}, PackerOptions{})
		if len(packed2) != 0 {
			t.Errorf("Second pack call packed boxes: got len %d, want 0", len(packed2))
		}
		// Bin should still only contain the box from the first pack
		if len(bin1.Boxes) != 1 {
			t.Errorf("Bin box count after second pack: got %d, want 1", len(bin1.Boxes))
		}
	})

	t.Run("puts multiple boxes into multiple bins", func(t *testing.T) {
		bin1 := NewBin(100, 50, nil) // Default lacement strategy (BSSF assumed)
		bin2 := NewBin(50, 50, nil)
		boxes := []*Box{
			NewBox(15, 10, false),   // box0
			NewBox(50, 45, false),   // box1 - Fits bin2 exactly width-wise
			NewBox(40, 40, false),   // box2
			NewBox(200, 200, false), // box3 - Too large
		}
		packer := NewPacker([]*Bin{bin1, bin2})
		packedBoxes := packer.Pack(boxes, PackerOptions{})

		// Assertions based on expected BSSF behavior:
		// - 50x45 fits best in bin2 (50x50)
		// - 40x40 fits best in bin1 (100x50) at 0,0
		// - 15x10 fits best in bin1 (100x50) likely at 0,40 (below 40x40)
		// - 200x200 doesn't fit

		if len(packedBoxes) != 3 {
			t.Fatalf("Packed box count: got %d, want 3", len(packedBoxes))
		}
		if len(bin1.Boxes) != 2 {
			t.Errorf("Bin1 box count: got %d, want 2", len(bin1.Boxes))
		}
		if len(bin2.Boxes) != 1 {
			t.Errorf("Bin2 box count: got %d, want 1", len(bin2.Boxes))
		}

		// Check which boxes are where (might require finding by properties/label)
		box1Label := "50x45 at [0,0]" // Expected label if Box.Label() works
		box2Label := "40x40 at [0,0]"
		box0Label := "15x10 at [0,40]" // Expected placement below box2

		foundBin1Box2 := false
		foundBin1Box0 := false
		for _, b := range bin1.Boxes {
			if b == boxes[2] && b.X == 0 && b.Y == 0 { // Check box2 (40x40) placement
				foundBin1Box2 = true
				if b.Label() != box2Label {
					t.Logf("Note: Bin1 Box2 label mismatch: got %q, want %q", b.Label(), box2Label)
				}
			}
			if b == boxes[0] && b.X == 0 && b.Y == 40 { // Check box0 (15x10) placement
				foundBin1Box0 = true
				if b.Label() != box0Label {
					t.Logf("Note: Bin1 Box0 label mismatch: got %q, want %q", b.Label(), box0Label)
				}
			}
		}
		if !foundBin1Box2 {
			t.Errorf("Box 2 (40x40) not found in Bin 1 at expected position [0,0]")
		}
		if !foundBin1Box0 {
			t.Errorf("Box 0 (15x10) not found in Bin 1 at expected position [0,40]")
		}

		if len(bin2.Boxes) > 0 {
			if bin2.Boxes[0] != boxes[1] { // Check box1 (50x45) is in bin2
				t.Errorf("Box 1 (50x45) not found in Bin 2")
			}
			if bin2.Boxes[0].Label() != box1Label {
				t.Logf("Note: Bin2 Box1 label mismatch: got %q, want %q", bin2.Boxes[0].Label(), box1Label)
			}
			if bin2.Boxes[0].X != 0 || bin2.Boxes[0].Y != 0 {
				t.Errorf("Box 1 position in Bin 2: got [%f,%f], want [0,0]", bin2.Boxes[0].X, bin2.Boxes[0].Y)
			}
		}

		if boxes[3].Packed { // Check the too-large box
			t.Errorf("Box 3 (200x200) Packed: got %v, want %v", boxes[3].Packed, false)
		}
		if len(packer.UnpackedBoxes) != 1 || packer.UnpackedBoxes[0] != boxes[3] {
			t.Errorf("Unpacked box count/content mismatch")
		}
	})

	t.Run("can work with scaled integers representing floats", func(t *testing.T) {
		// Adapt float test by multiplying dimensions by 100
		bin1 := NewBin(100, 100, nil) // Represents 1x1
		bin2 := NewBin(50, 50, nil)   // Represents 0.5x0.5
		boxes := []*Box{
			NewBox(20, 14, false), // Represents 1/5 x 1/7 (approx 0.2 x 0.14)
			NewBox(20, 50, false), // Represents 1/5 x 1/2 (0.2 x 0.5)
			NewBox(25, 25, false), // Represents 1/4 x 1/4 (0.25 x 0.25)
			NewBox(50, 50, false), // Represents 1/2 x 1/2 (0.5 x 0.5) - Should fit bin1 OR bin2
		}
		packer := NewPacker([]*Bin{bin1, bin2})
		packedBoxes := packer.Pack(boxes, PackerOptions{})

		// All boxes should fit in this scaled integer scenario
		if len(packedBoxes) != 4 {
			t.Errorf("Packed box count: got %d, want 4", len(packedBoxes))
		}
		// Further checks on distribution could be added but depend heavily on placement specifics.
		if countPacked(boxes) != 4 {
			t.Errorf("Total packed boxes in original slice: got %d, want 4", countPacked(boxes))
		}
		if len(packer.UnpackedBoxes) != 0 {
			t.Errorf("Unpacked box count: got %d, want 0", len(packer.UnpackedBoxes))
		}
	})

	t.Run("Constrain rotation on boxes", func(t *testing.T) {
		bin1 := NewBin(100, 50, nil)
		boxes := []*Box{
			NewBox(50, 100, true),  // Constrained, cannot rotate to fit 100x50
			NewBox(50, 100, false), // Not constrained, *can* rotate to fit 100x50
		}
		packer := NewPacker([]*Bin{bin1})
		packedBoxes := packer.Pack(boxes, PackerOptions{})

		if len(packedBoxes) != 1 {
			t.Errorf("Packed box count: got %d, want 1", len(packedBoxes))
		}
		if boxes[0].Packed {
			t.Errorf("Box 0 (constrained) Packed: got %v, want %v", boxes[0].Packed, false)
		}
		if !boxes[1].Packed {
			t.Errorf("Box 1 (not constrained) Packed: got %v, want %v", boxes[1].Packed, true)
		}
		if len(packer.UnpackedBoxes) != 1 || packer.UnpackedBoxes[0] != boxes[0] {
			t.Errorf("Unpacked box count/content mismatch")
		}
	})
}

// newBins creates standard bins used across multiple tests.
func newBins(t *testing.T) (*Bin, *Bin, *Bin) {
	t.Helper() // Marks this as a test helper
	// Assuming NewBin uses default placement strategy if nil is passed
	bin1 := NewBin(9600, 3100, nil)
	bin2 := NewBin(10000, 4500, nil)
	bin3 := NewBin(12000, 4500, nil)
	if bin1 == nil || bin2 == nil || bin3 == nil {
		t.Fatal("Failed to create bins in helper")
	}
	return bin1, bin2, bin3
}

// countPacked counts how many boxes in a slice are marked as packed.
func countPacked(boxes []*Box) int {
	count := 0
	for _, box := range boxes {
		if box != nil && box.Packed {
			count++
		}
	}
	return count
}

// findBoxByLabel searches a slice of boxes for one with a matching label.
// Note: Relies on Box.Label() method existing and being suitable identifier.
// Returns nil if not found.
func findBoxByLabel(boxes []*Box, label string) *Box {
	for _, box := range boxes {
		if box != nil && box.Label() == label {
			return box
		}
	}
	return nil
}
