package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
)

var (
	in  = flag.String("in", "", "input image")
	out = flag.String("out", "", "output image")
)

type cn struct {
	i float64
	r float64
}

func main() {
	flag.Parse()
	f, err := os.Open(*in)
	defer f.Close()
	if err != nil {
		println(err)
		os.Exit(1)
	}
	im, _, err := image.Decode(f)
	//im := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{1000, 1000}})
	hough, _ := hough(im)

	outFile, err := os.Create(*out)
	defer outFile.Close()
	if err != nil {
		println(err)
		os.Exit(1)
	}
	png.Encode(outFile, hough)
}

func hough(i image.Image) (image.Image, error) {
	minDist := math.Pow(2, 16) - 1
	maxDist := float64(0)
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()
	midX := float64(width) / 2
	midY := float64(height) / 2

	maxDistance := math.Sqrt(float64(width*width+height*height)) / 2
	println(maxDistance)

	dstSize := 400
	hough := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{dstSize, dstSize}})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, _, _, _ := i.At(x, y).RGBA()

			px := float64(x) - midX
			py := float64(y) - midY
			//black pixel
			if r == 0 {
				//iterate the angles
				for t := 0; t < dstSize; t++ {
					angle := normalise(0, float64(dstSize), 0, math.Pi, float64(t))
					//Get normal distance can be negative
					distance := px*math.Cos(angle) + py*math.Sin(angle)
					if distance < minDist {
						minDist = distance
					}
					if distance > maxDist {
						maxDist = distance
					}
					angle = normalise(0, math.Pi, 0, float64(dstSize), angle)
					distance = normalise(-maxDistance, maxDistance, 0, float64(dstSize), distance)
					g := hough.Gray16At(int(angle+0.5), int(distance+0.5))
					g.Y += 100
					hough.SetGray16(int(angle+0.5), int(distance+0.5), g)
				}
			}
		}
	}
	fmt.Printf("Min dst: %v, scaled:%v\n", minDist, normalise(0, maxDistance, 0, float64(dstSize), minDist))
	fmt.Printf("Max dst: %v, scaled:%v\n", maxDist, normalise(0, maxDistance, 0, float64(dstSize), maxDist))
	return hough, nil
	//r = x cos T + y sin T
}
func normalise(srcMin, srcMax, dstMin, dstMax, srcVal float64) float64 {
	return ((dstMax - dstMin) * (srcVal - srcMin) / (srcMax - srcMin)) + dstMin
}
