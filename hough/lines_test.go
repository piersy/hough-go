package hough

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"
)

var (
	l1 = newLineAndScore(1, 1, 1)
	l2 = newLineAndScore(2, 2, 2)
	l3 = newLineAndScore(3, 3, 3)
)

func TestLines(t *testing.T) {
	tl := topLines{
		lineCount: 3,
	}
	tl.addLine(1, 1, 1)
	if tl.boundary != 0 {
		t.Errorf("expecting boundary to be 0 but was %d", tl.boundary)
	}
	tl.addLine(2, 2, 2)
	if tl.boundary != 0 {
		t.Errorf("expecting boundary to be 0 but was %d", tl.boundary)
	}
	tl.addLine(3, 3, 3)
	// At this point we expect the boundary to have updated to the lowest
	// scoring line in the set.
	if tl.boundary != 1 {
		t.Errorf("expecting boundary to be 1 but was %d", tl.boundary)
	}

}
