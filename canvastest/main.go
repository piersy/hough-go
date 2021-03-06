package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"reflect"

	"github.com/piersy/hough-go/canvas"
)

var (
	in  = flag.String("in", "", "input image")
	out = flag.String("out", "", "output image")
)

// This program exercises the LineTo functionality of the drawing context.
func main() {
	err := dotest()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error drawing line :: %v\n", err)
		os.Exit(1)
	}
}
func dotest() error {
	flag.Parse()
	im := image.NewGray(image.Rect(0, 0, 100, 100))
	fmt.Printf("image type: %s\n", reflect.TypeOf(im))
	fmt.Printf("image bounds: %v\n", im.Bounds())
	c := canvas.New()
	c.Color(color.White)
	centre := image.Pt(49, 49)
	incr := math.Pi / (2 * 10)
	t := float64(0)
	for ; t <= math.Pi*2; t += incr {
		c.MoveTo(centre)
		y := math.Sin(t) * 50
		x := math.Cos(t) * 50
		c.LineTo(image.Pt(int(x+0.5)+centre.X, int(y+0.5)+centre.Y))
	}
	c.Render(im)
	outFile, err := os.Create(*out)
	defer outFile.Close()
	if err != nil {
		return err
	}
	png.Encode(outFile, im)
	return nil

}
