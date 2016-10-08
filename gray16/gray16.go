package gray16

import (
	"image"
	"image/color"
	"math"
)

// Gray16 is an in-memory image whose At method returns color.Gray16 values.
type Gray16 struct {
	// Pix holds the image's pixels, as gray values in big-endian format. The pixel at
	// (x, y) is at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)].
	Pix []uint16
	// MaxVal is the highest value occurring in this image
	MaxVal uint16
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *Gray16) ColorModel() color.Model { return color.Gray16Model }

func (p *Gray16) Bounds() image.Rectangle { return p.Rect }

func (p *Gray16) At(x, y int) color.Color {
	return p.Gray16At(x, y)
}

func (p *Gray16) Gray16At(x, y int) color.Gray16 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.Gray16{}
	}
	i := p.PixOffset(x, y)
	return color.Gray16{p.Pix[i]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *Gray16) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x - p.Rect.Min.X)
}

func (p *Gray16) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.Gray16Model.Convert(c).(color.Gray16)
	p.Pix[i] = c1.Y
}

func (p *Gray16) SetGray16(x, y int, c color.Gray16) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = c.Y
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *Gray16) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are image.Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &Gray16{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &Gray16{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *Gray16) Opaque() bool {
	return true
}

// NewGray16 returns a new Gray16 with the given bounds.
func NewGray16(r image.Rectangle) *Gray16 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint16, w*h)
	return &Gray16{
		Pix:    pix,
		Stride: w,
		Rect:   r,
	}
}

func (p *Gray16) Normalise() {
	ratio := float64(math.MaxUint16) / float64(p.MaxVal)
	for i, v := range p.Pix {
		p.Pix[i] = uint16(float64(v) * ratio)
	}
	p.MaxVal = math.MaxUint16
}
