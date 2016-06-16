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

	var maxRadius float64
	maxRadius = float64(i.Bounds().Dx()*i.Bounds().Dx() + i.Bounds().Dy()*i.Bounds().Dy())
	maxRadius = math.Sqrt(float64(maxRadius))

	hough := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{180, int(maxRadius)}})

	for x := 0; x < i.Bounds().Dx(); x++ {
		for y := 0; y < i.Bounds().Dy(); y++ {
			r, _, _, _ := i.At(x, y).RGBA()

			//black pixel
			if r == 0 {
				for t := 0; t < 180; t++ {
					radians := -(math.Pi / 2.0) + (float64(t) * math.Pi / 180.0)
					r := float64(x)*math.Cos(radians) - float64(y)*math.Sin(radians)
					g := hough.Gray16At(t, int((maxRadius/2.0)-r))
					g.Y += 300
					hough.SetGray16(t, int((maxRadius/2.0)-r), g)
				}
			}
		}
	}
	return hough, nil
	//r = x cos T + y sin T
}
