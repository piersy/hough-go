package canvas

import (
	"image"
	"math"
	"testing"
)

func TestIntercepts(t *testing.T) {
	bounds := image.Rect(-10, -10, 10, 10)

	angle := 0.0
	//Vertical lines
	//	for x := -10; x <= 10; x++ {
	//		expected := []image.Point{image.Pt(x, 10), image.Pt(x, -10)}
	//		verifyExpected(t, angle, x, bounds, expected)
	//	}

	angle = math.Pi / 2
	// Horizontal lines
	for y := -10; y <= -10; y++ {
		expected := []image.Point{image.Pt(-10, y), image.Pt(10, y)}
		verifyExpected(t, angle, y, bounds, expected)
	}
}

func TestComparison(t *testing.T) {

	if !(-10.000000 >= -10.000000) {
		t.Fail()
	}
}

//Gradient: -0.000
//x distance: -0.000
//y distance: -10.000
//low x y intercept -10.000000
func verifyExpected(t *testing.T, angle float64, dist int, bounds image.Rectangle, expected []image.Point) {
	points := getIntercepts(angle, float64(dist), bounds)
	pointsCopy := points
	if len(points) != 2 {
		t.Errorf("Expecting 2 intercepts got %d", points)
	}
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
		t.Errorf("Expecting %v but got %v", expected, pointsCopy)
	}

}
