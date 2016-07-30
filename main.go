package main

import (
	"flag"
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
	println("yoyoyo")
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()
	midX := float64(width) / 2
	midY := float64(height) / 2

	maxDistance := math.Sqrt(float64(width*width+height*height)) / 2

	dstSize := 400
	hough := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{dstSize, dstSize}})

	for x := 0; x < i.Bounds().Dx(); x++ {
		for y := 0; y < i.Bounds().Dy(); y++ {
			r, _, _, _ := i.At(x, y).RGBA()

			px := float64(x) - midX
			py := float64(y) - midY
			//black pixel
			if r == 0 {
				//Get angle from centre
				angle := math.Atan2(py, px)
				//Get distance from centre
				distance := math.Sqrt(px*px + py*py)
				angle = normalise(-math.Pi, math.Pi, 0, float64(dstSize), angle)
				distance = normalise(0, maxDistance, 0, float64(dstSize), distance)
				g := hough.Gray16At(int(angle), int(distance))
				g.Y += uuuuu
				hough.SetGray16(int(angle), int(distance), g)
			}
		}
	}
	return hough, nil
	//r = x cos T + y sin T
}
func normalise(srcMin, srcMax, dstMin, dstMax, srcVal float64) float64 {
	return ((dstMax - dstMin) * (srcVal - srcMin) / (srcMax - srcMin)) + dstMin
}
