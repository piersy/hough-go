// Package hough provides an implementation of the hough line detector
package hough

import (
	"image"
	"math"

	"github.com/piersy/hough-go/util"
)

type houghParams struct {
	sinAngles, cosAngles []float64
	width, height        int
	midX, midY           float64
	distanceDivisions    int
	distanceNormaliser   util.Normaliser
}

// HoughParams generates a set of parameters with which a hough transform can
// be applied the same parameters can be reused on many transforms having the
// same input and output image sizes. The distance divisions and angle
// divisions define the hough transform output size, with distanceDivisions
// being pixels in the y dimension and angleDivisions being pixels in the x
// dimension.
func HoughParams(input image.Image, distanceDivisions, angleDivisions int) houghParams {
	params := houghParams{
		sinAngles:         make([]float64, angleDivisions),
		cosAngles:         make([]float64, angleDivisions),
		width:             input.Bounds().Dx(),
		height:            input.Bounds().Dy(),
		distanceDivisions: distanceDivisions,
	}
	params.midX = float64(params.width) / 2
	params.midY = float64(params.height) / 2
	// Converts our accumulator angle buckets into appropriate radian values between 0 and Pi
	angleN := util.NewNormaliser(0, float64(angleDivisions), 0, math.Pi)
	for t := 0; t < angleDivisions; t++ {
		a := angleN.Normalise(float64(t))
		params.sinAngles[t] = math.Sin(a)
		params.cosAngles[t] = math.Cos(a)
	}

	// The max distance from centre, used for normalising the distance from the
	// source image to the size of the accumulator.
	maxDistance := math.Sqrt(float64(params.width*params.width+params.height*params.height)) / 2
	params.distanceNormaliser = util.NewNormaliser(-maxDistance, maxDistance, 0, float64(distanceDivisions))
	return params
}

// Hough takes an input image and and a houghParams instance and returns the
// hough transform of that image. The y axis  of the output image represents the
// line distance from centre of the input and the x axis represents the angle
// of the line. Only black pixels are considered as contributing to the hough
// transform.
func Hough(input image.Image, hp houghParams) *util.Gray16 {
	at := getRgba(input)
	acc := util.NewGray16(image.Rect(0, 0, hp.distanceDivisions, len(hp.cosAngles)))
	stride := acc.Stride
	pix := acc.Pix
	var maxVal uint16

	// Iterate each pixel in the source
	for x := 0; x < hp.width; x++ {
		px := float64(x) - hp.midX
		for y := 0; y < hp.height; y++ {
			py := float64(y) - hp.midY

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
				for t := 0; t < len(hp.sinAngles); t++ {
					//Get normal distance - can be negative
					distance := px*hp.cosAngles[t] + py*hp.sinAngles[t]
					// normalize distance into accumulator range.
					// Accumulator range cannot benegative
					dist := hp.distanceNormaliser.Normalise(distance)
					// The distance is likely to fall between two of our
					// accumulator buckets so we divide the score
					// appropriately between the buckets.
					intDist := int(dist)
					floatingPointPart := dist - float64(intDist)
					//find different components of the score
					further := 10.0 * floatingPointPart
					nearer := 10.0 * (1.0 - floatingPointPart)
					// Update the further pixel
					pixel := (intDist+1)*stride + t
					if pixel < len(pix) {
						increment(uint16(further), &pix[pixel], &maxVal)
					}
					// Update the nearer pixel
					pixel = intDist*stride + t
					if pixel < len(pix) {
						increment(uint16(nearer), &pix[pixel], &maxVal)
					}
				}
			}
		}
	}
	// Set the max val on the acc so that it can be normalised correctly
	acc.MaxVal = maxVal
	return acc
}

func increment(inc uint16, initial, max *uint16) {
	result := *initial + inc
	if result > *max {
		*max = result
	}
	//Overflow simply hard limit at max
	if result < *initial {
		*max = math.MaxUint16
		result = math.MaxUint16
	}
	*initial = result
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
