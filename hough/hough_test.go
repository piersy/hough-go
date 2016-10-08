package hough

import (
	"flag"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"testing"

	"github.com/piersy/hough-go/norm"
)

var (
	in = flag.String("in", "", "input image")
)

func init() {
	flag.Parse()
}

func BenchmarkHough(b *testing.B) {
	input := getImage(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hough(input, 400, 400)
	}
}

func BenchmarkNormaliseGray16(b *testing.B) {
	input := getImage(b)
	gray := Hough(input, 400, 400)
	pix := gray.Pix
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		normalise(pix, 0, math.MaxUint16/2)
	}
}

func normalise(pix []uint16, min, max uint16) {
	n := norm.NewNormaliser(float64(min), float64(max), 0, float64(math.MaxUint16))
	for i := 0; i < len(pix); i++ {
		pix[i] = uint16(n.Normalise(float64(pix[i])) + 0.5)
	}
}

func BenchmarkNormaliseGray16Method(b *testing.B) {
	input := getImage(b)
	gray := Hough(input, 400, 400)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gray.Normalise()
	}
}

func getImage(b *testing.B) image.Image {
	f, err := os.Open(*in)
	defer f.Close()
	if err != nil {
		b.Fatal(err)
	}
	input, _, err := image.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	return input
}
