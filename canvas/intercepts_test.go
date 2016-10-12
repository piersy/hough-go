package canvas

import (
	"image"
	"testing"
)

func TestIntercepts(t *testing.T) {
	bounds := image.Rect(-10, -10, 10, 10)

	// Special cases
	// vertical line at 10x intercepts (10, -10) and (10, 10)
	//This doesn't work for special case where x or y lie on a boundary in which case the diff oneway will be 0 and 0 times anything is 0
	points := getIntercepts(0, 9, bounds)
	if len(points) != 2 {
		t.Errorf("Expecting 2 intercepts got %d", points)
	}
	expected := []image.Point{image.Pt(9, 10), image.Pt(9, -10)}
	for _, e := range expected {
		for i, p := range points {
			if e == p {
				points[i] = points[len(points)-1]
				points = points[:len(points)-1]
				break
			}
		}
	}
	if len(points) > 0 {
		t.Errorf("Expecting %v but got %v", expected, getIntercepts(0, 9, bounds))
	}
}
