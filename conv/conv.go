// Package conv provides functionality to convolve images
package conv

import (
	"image/color"
	"math"

	"github.com/piersy/hough-go/util"
)

type kernel struct {
	k      []uint16
	stride uint16
}

func AdaptiveThresh(input *util.Gray16) *util.Gray16 {
	c := float64(math.MaxUint16 / 2)
	max := color.Gray16{math.MaxUint16}
	min := color.Gray16{}
	output := util.NewGray16(input.Bounds())
	width := input.Bounds().Dx()
	height := input.Bounds().Dy()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			avg := uint16(0)
			for z := -1; z < 2; z++ {
				for a := -1; a < 2; a++ {
					avg += input.Gray16At(x+z, y+a).Y
				}
				if float64(input.Gray16At(x, y).Y) > float64(avg)/9.0+c {
					output.SetGray16(x, y, max)
				} else {
					output.SetGray16(x, y, min)
				}
			}
		}
	}
	return output
}
