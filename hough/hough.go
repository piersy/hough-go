// Package hough provides an implementation of the hough line detector
package hough

import (
	"image"
	"math"

	"github.com/piersy/hough-go/util"
)

// Hough takes an input image and returns the hough transform of that image
// with size accDistance * accAngles. The y axis  of this image represents the
// line distance from centre of the input and the x axis represents the angle
// of the line. Only black pixels are considered as contributing to the hough
// transform.
func Hough(input image.Image, accDistance, accAngle int) *util.Gray16 {
	width := input.Bounds().Dx()
	height := input.Bounds().Dy()
	midX := float64(width) / 2
	midY := float64(height) / 2
	// Precalculate angles for sin and cos
	sinAngles := make([]float64, accDistance)
	cosAngles := make([]float64, accDistance)
	// Converts our accumulator angle buckets into appropriate radian values between 0 and Pi
	angleN := util.NewNormaliser(0, float64(accAngle), 0, math.Pi)
	for t := 0; t < accAngle; t++ {
		a := angleN.Normalise(float64(t))
		sinAngles[t] = math.Sin(a)
		cosAngles[t] = math.Cos(a)
	}

	// The max distance from centre, used for normalising the distance from the
	// source image to the size of the accumulator.
	maxDistance := math.Sqrt(float64(width*width+height*height)) / 2
	distN := util.NewNormaliser(-maxDistance, maxDistance, 0, float64(accDistance))

	at := getRgba(input)
	acc := util.NewGray16(image.Rect(0, 0, accDistance, accAngle))
	stride := acc.Stride
	pix := acc.Pix
	var maxVal uint16

	// Iterate each pixel in the source
	for x := 0; x < width; x++ {
		px := float64(x) - midX
		for y := 0; y < height; y++ {
			py := float64(y) - midY

			// check black pixel
			r, g, b, _ := at(x, y)
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
					nDistance := int(distN.Normalise(distance) + 0.5)
					//		g := acc.Gray16At(t, nDistance)
					//		g.Y += 100
					//acc.SetGray16(t, nDistance, g)
					// update the accumulator
					// Get pixel start location
					pixStart := nDistance*stride + t
					val := pix[pixStart]
					val++
					if val > maxVal {
						maxVal = val
					}
					pix[pixStart] = val
				}
			}
		}
	}
	// Set the max val on the acc so that it can be normalised correctly
	acc.MaxVal = maxVal
	return acc
}

// getAt reruns the At method defined on the underlying struct implementing
// image.Image, if the image type is unknown then this function panics.
// This performs about 25% better than calling at on an interface, I will accept it for now.
// best performance is achieved by calling the types at method in the main loop of hough
// but that would mean writing the algorithm once for each image type.
func getRgba(i image.Image) func(int, int) (uint32, uint32, uint32, uint32) {
	switch t := i.(type) {
	case *image.Alpha:
		return nil
	case *image.Alpha16:
		return nil
	case *image.Gray:
		return nil
	case *image.Gray16:
		return nil
	case *image.NRGBA:
		return func(x, y int) (uint32, uint32, uint32, uint32) {
			return t.NRGBAAt(x, y).RGBA()
		}
	case *image.NRGBA64:
		return nil
	case *image.Paletted:
		return nil
	case *image.RGBA:
		return func(x, y int) (uint32, uint32, uint32, uint32) {
			return t.RGBAAt(x, y).RGBA()
		}
	case *image.RGBA64:
		return nil
	default:
		panic("unrecognised image type")
	}
}
