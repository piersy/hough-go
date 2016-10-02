// Package hough provides an implementation of the hough line detector
package hough

import (
	"fmt"
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
	fmt.Printf("input image size (%d,%d)\n", width, height)
	midX := float64(width) / 2
	midY := float64(height) / 2
	// Precalculate angles and scale factor for sin and cos
	sinAngles := make([]float64, accAngle)
	cosAngles := make([]float64, accAngle)
	scaleFactor := make([]float64, accAngle)
	//scaleFactor := make([]float64, accAngle)
	// Converts our accumulator angle buckets into appropriate radian values between 0 and Pi
	angleN := util.NewNormaliser(0, float64(accAngle), 0, math.Pi)
	//we also need to scale distances based on the angle to account for the difference between
	//diagonal and straight edges across a pixel.
	// at 45degrees we want to scale distances by 1/sqrrt(2) at vertical and horizontal then the scale is 1
	// 45 degreees is pi/4 and 3pi/4 my problem is that this has to rise up and drop twice for the precalcs or just
	// calculate the first quarter and then reverse and append but only works for even lenghthed accumulators
	//we can use a triangle wave formula to generate the input angle for calculating the scale factor.
	po4 := math.Pi / 4.0
	for t := 0; t < accAngle; t++ {
		a := angleN.Normalise(float64(t))
		sinAngles[t] = math.Sin(a)
		cosAngles[t] = math.Cos(a)
		//cretes a triangle wave rising fom 0 to pi/4 back to 0 again
		//we use the wave as input to 1 * 1/(cos x) calculate the area accross a pixel with

		//     --------------
		//          /       |
		//         /        |
		//        /         |
		//       /          |
		//      /           |
		//     /x)          |
		//     --------------
		scaleAngle := po4 - math.Abs(math.Mod(a, 2*po4)-po4)
		scaleFactor[t] = math.Cos(scaleAngle)
		//println("scaleFac", t, scaleFactor[t])
	}

	// The max distance from centre, used for normalising the distance from the
	// source image to the size of the accumulator.
	maxDistance := math.Sqrt(float64(width*width+height*height)) / 2
	//println("maxDist", maxDistance)
	//maxDistance /= maxDistance / math.Max(midX, midY)
	//	//println("maxDist", maxDistance)
	//maxDistance = math.Max(midX, midY)
	fmt.Printf("maxDist %f\n\n", maxDistance)

	//	the max dist in fact needs to be scaled with the appropriate scale factor depending on the angle from centre image to corner
	//in fact the max distance once scaled will be the simply the greater of either half width or half height.
	distN := util.NewNormaliser(-maxDistance, maxDistance, 0, float64(accDistance-1))

	at := getRgba(input)
	acc := util.NewGray16(image.Rect(0, 0, accDistance, accAngle))
	stride := acc.Stride
	pix := acc.Pix
	var maxVal uint16

	var maxDist float64 = 0.0
	// Iterate each pixel in the source
	for x := 0; x < width; x++ {
		px := float64(x) - midX
		for y := 0; y < height; y++ {
			py := float64(y) - midY
			//Ok the second part of this should be normalised to 1 but it is not
			// see http://mathproofs.blogspot.co.uk/2005/07/mapping-square-to-circle.html
			//these are not quite working
			//sx := px * math.Sqrt(1.0-(py/midY*py/midX/2.0))
			//sy := py * math.Sqrt(1.0-(px/midX*px/midX/2.0))

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
					//square shit
					//neg := distance < 0.0
					////normalise into square coords and recalculate distance
					//nowX := distance * math.Cos(angleN.Normalise(float64(t)))
					//nowY := distance * math.Sin(angleN.Normalise(float64(t)))
					//sx := nowX * math.Sqrt(1.0-(nowY/midY*nowY/midX/2.0))
					//sy := nowY * math.Sqrt(1.0-(nowX/midX*nowX/midX/2.0))
					//distance = math.Sqrt(sx*sx + sy*sy)
					//if neg {
					//	distance = -distance
					//}

					//adjust distance by scale factor
					distance = distance * scaleFactor[t]

					//distance = distance / math.Sqrt2
					// normalize distance into accumulator range
					nDistance := int(distN.Normalise(distance))
					if distance > maxDist {
						maxDist = distance
					}
					if distance > maxDistance {
						//fmt.Printf("Raw dist %f\n", rawDist)
						fmt.Printf("scaled dist %f\n", distance)
						fmt.Printf("max dist for this spot %f\n", math.Sqrt(px*px+py*py))
						fmt.Printf("scale fac %f\n", scaleFactor[t])
						fmt.Printf("ndistatnce %d\n", nDistance)
						fmt.Printf(" acc angle %d\n", t)
						fmt.Printf(" real angle %f\n", angleN.Normalise(float64(t)))
						fmt.Printf("x:%f y:%f\n", px, py)
						fmt.Printf("\n")
						//os.Exit(1)

					}

					//g := acc.Gray16At(t, nDistance)
					//g.Y += 100
					//acc.SetGray16(t, nDistance, g)
					//update the accumulator
					//Get pixel start location
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
	fmt.Printf("Max unnormalised but scaled dist %f\n", maxDist)

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
