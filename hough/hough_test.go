package hough

import (
	"flag"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"
)

var (
	in = flag.String("in", "", "input image")
)

func init() {
	flag.Parse()
}

func BenchmarkHough(b *testing.B) {
	flag.Parse()
	f, err := os.Open(*in)
	defer f.Close()
	if err != nil {
		b.Fatal(err)
	}
	input, _, err := image.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hough(input, 400, 400)
	}
}
