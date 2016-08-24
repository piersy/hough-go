// Package hough provides an implementation of the hough line detector
package hough

import (
	"image"
	"image/draw"
	"math"
)

// Hough takes an input image and returns the hough transform of that image
// with size accDistance * accAngles. The y axis  of this image represents the
// line distance from centre of the input and the x axis represents the angle
// of the line. Only black pixels are considered as contributing to the hough
// transform.
func Hough(input image.Image, accDistance, accAngle int) (draw.Image, []Line) {
	width := input.Bounds().Dx()
	height := input.Bounds().Dy()
	midX := float64(width) / 2
	midY := float64(height) / 2
	// Precalculate angles for sin and cos
	sinAngles := make([]float64, accDistance)
	cosAngles := make([]float64, accDistance)
	// Converts our accumulator angle buckets into appropriate radian values between 0 and Pi
	angleN := NewNormaliser(0, float64(accAngle), 0, math.Pi)
	for t := 0; t < accAngle; t++ {
		a := angleN.normalise(float64(t))
		sinAngles[t] = math.Sin(a)
		cosAngles[t] = math.Cos(a)
	}

	// The max distance from centre, used for normalising the distance from the
	// source image to the size of the accumulator.
	maxDistance := math.Sqrt(float64(width*width+height*height)) / 2
	distN := NewNormaliser(-maxDistance, maxDistance, 0, float64(accDistance))

	nrgba := *input.(*image.RGBA)
	acc := NewGray16(image.Rect(0, 0, accDistance, accAngle))
	stride := acc.Stride
	pix := acc.Pix

	// keeps track of the highest scoring lines
	tl := topLines{
		lineCount: 100,
	}
	// Iterate each pixel in the source
	for x := 0; x < width; x++ {
		px := float64(x) - midX
		for y := 0; y < height; y++ {
			py := float64(y) - midY

			// check black pixel
			r, g, b, _ := nrgba.RGBAAt(x, y).RGBA()
			if r&g&b == 0 {
				// For all angles represented in the accumulator, calculate
				// perpendicular distance to the center of the input for a line
				// through (x, y) at each angle and plot (dist, angle) in the
				// accumulator.  Subsequent pixels that form a line of angle t
				// with this pixel will share the same perpendicular distance
				// at angle t and hence the point (d(t), t) will conicide for
				// all pixels along the line.
				for t := 0; t < accAngle; t++ {
					//Get normal distance can be negative
					distance := px*cosAngles[t] + py*sinAngles[t]
					// normalize distance into accumulator range
					nDistance := int(distN.normalise(distance) + 0.5)
					//		g := acc.Gray16At(t, nDistance)
					//		g.Y += 100
					//acc.SetGray16(t, nDistance, g)
					// update the accumulator
					// Get pixel start location
					pixStart := nDistance*stride + t
					val := pix[pixStart]
					val += 100
					pix[pixStart] = val
					// Add the line to the top lines if it scores above the boundary
					if val > tl.boundary {
						tl.addLine(angleN.normalise(float64(t)), distance, val)
					}
				}
			}
		}
	}
	return acc, tl.lines()
}

func NewNormaliser(srcMin, srcMax, dstMin, dstMax float64) noramliser {
	return noramliser{
		ratio:  (dstMax - dstMin) / (srcMax - srcMin),
		srcMin: srcMin,
		dstMin: dstMin,
	}
}

type noramliser struct {
	ratio, srcMin, dstMin float64
}

func (n noramliser) normalise(val float64) float64 {
	return n.ratio*(val-n.srcMin) + n.dstMin
}
