package canvas

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

type Context interface {
	Color(color.Color)
	Move(p image.Point)
	Line(p image.Point)
	Render(im draw.Image) error
}

type op interface {
	do(i draw.Image) error
}

type line struct {
	p1, p2 image.Point
	c      color.Color
}

func (l line) do(i draw.Image) error {
	return DrawLine(l.p1, l.p2, l.c, i)
}

type context struct {
	c   color.Color
	ops []op
	p   image.Point
}

func New() Context {
	return &context{
		c: color.White,
		p: image.ZP,
	}
}

func (c *context) Color(col color.Color) {
	c.c = col
}
func (c *context) Move(pt image.Point) {
	c.p = pt
}

func (c *context) Line(pt image.Point) {
	c.ops = append(c.ops, line{c.p, pt, c.c})
}

func (c *context) Render(im draw.Image) error {
	for _, op := range c.ops {
		err := op.do(im)
		if err != nil {
			return err
		}
	}
	return nil
}

func DrawLine(p1, p2 image.Point, col color.Color, im draw.Image) error {
	b := im.Bounds()
	if !p1.In(b) {
		return fmt.Errorf("Point %v not in bounds %v", p1, b)
	}
	if !p2.In(b) {
		return fmt.Errorf("Point %v not in bounds %v", p2, b)
	}

	// determine which way to draw the line
	// top to bottom or left to right
	xdiff := p2.X - p1.X
	ydiff := p2.Y - p1.Y
	fmt.Printf("xdiff: %d, ydiff: %d\n", xdiff, ydiff)
	var yInc, xInc float64
	var iter, inc int
	//determine which direction the line is longest in
	if math.Abs(float64(xdiff)) >= math.Abs(float64(ydiff)) {
		iter = xdiff
		// increment either 1 or -1 based on the
		// value of xdiff
		inc = int(int32(uint32(xdiff)&uint32(0x80000000))>>31 | 1)
		xInc = float64(inc)
		fmt.Printf("xinc %2f\n", xInc)
		yInc = float64(ydiff) / float64(xdiff)
		//Remove possible negative sign given to yInc if xdiff was negative
		yInc *= xInc
	} else {
		iter = ydiff
		// increment either 1 or -1 based on the
		// value of Ydiff
		inc = int(int32(uint32(ydiff)&uint32(0x80000000))>>31 | 1)
		yInc = float64(inc)
		xInc = float64(xdiff) / float64(ydiff)
		//Remove possible negative sign given to xInc if ydiff was negative
		xInc *= yInc
	}

	//remove negativeness of iter
	iter *= inc
	inc = 1
	x := float64(p1.X) + 0.5
	y := float64(p1.Y) + 0.5

	//fmt.Printf("Drawing line from (%d,%d) to (%d,%d) with a longest dimension of %d and increments (%2f, %2f)\n", p1.X, p1.Y, p2.X, p2.Y, iter, xInc, yInc)
	for i := 0; i <= iter; i += inc {
		im.Set(int(x), int(y), col)
		x += xInc
		y += yInc
	}
	return nil
}
