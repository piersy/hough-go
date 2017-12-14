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
	MoveTo(p image.Point)
	LineTo(p image.Point)
	Line(dist, angle float64)
	Render(im draw.Image)
}

type op interface {
	do(i draw.Image)
}

type lineTo struct {
	p1, p2 image.Point
	c      color.Color
}

func (l lineTo) do(i draw.Image) {
	DrawLine(l.p1, l.p2, l.c, i)
}

type line struct {
	d, a float64
	c    color.Color
}

func (l line) do(i draw.Image) {
	points := getIntercepts(l.a, l.d, i.Bounds())
	if len(points) == 2 {
		DrawLine(points[0], points[1], l.c, i)
	}
}

// getIntercepts returns the points at which the line defined by angle in
// radians and distance in pixels from the centre of the given rectangle
// intercept with the bounds of the rectangle. If the returned slice is empty
// it indicates that the given line does not intercept the bounds of the
// rectangle.
func getIntercepts(angle, distance float64, b image.Rectangle) []image.Point {
	//Special cases for lines without gradient, this avoids division by zero errors.
	// vertical
	if angle == 0 || math.Mod(angle, math.Pi) == 0 {
		return []image.Point{image.Point{int(distance), b.Min.Y}, image.Point{int(distance), b.Max.Y}}

	}
	//horizontal
	if math.Mod(angle, math.Pi/2.0) == 0 {
		return []image.Point{image.Point{b.Min.X, int(distance)}, image.Point{b.Max.X, int(distance)}}

	}
	// Find intercepts with bounds
	var intercept float64
	xmin := float64(b.Min.X)
	xmax := float64(b.Max.X)
	ymin := float64(b.Min.Y)
	ymax := float64(b.Max.Y)
	fmt.Printf("xmin: %.3f\n", xmin)
	fmt.Printf("xmax: %.3f\n", xmax)
	fmt.Printf("ymin: %.3f\n", ymin)
	fmt.Printf("ymax: %.3f\n", ymax)

	gradient := math.Tan(angle + math.Pi/2.0)
	fmt.Printf("Gradient: %.3f\n", gradient)
	// calculate the coords of the intersection of the line and the
	// perpendicular relative to the centre of the given rectangle.  We will
	// then project both ways from this point, using the gradient to direct the
	// projection and see where that intersects with the bounds.
	x := float64(b.Bounds().Dx())/2.0 + xmin + distance*math.Cos(angle)
	y := float64(b.Bounds().Dy())/2.0 + ymin + distance*math.Sin(angle)

	fmt.Printf("x distance: %.3f\n", x)
	fmt.Printf("y distance: %.3f\n", y)
	var points []image.Point

	// y intercept at lower x bound
	intercept = y - ((x - xmin) * gradient)
	if intercept >= ymin && intercept <= ymax {
		fmt.Printf("low x, y intercept %3f\n", intercept)
		points = append(points, image.Point{int(xmin), int(intercept)})
	}
	// y intercept at upper x bound
	intercept = y + ((xmax - x) * gradient)
	fmt.Printf("high x, y intercept %3f\n", intercept)
	fmt.Printf("%g >= %g && %g <= %g, %t\n", intercept, ymin, intercept, ymax, intercept >= ymin)
	if intercept >= ymin && intercept <= ymax {
		fmt.Println("doing the high x intercept")
		fmt.Printf("high x, y intercept %3f\n", intercept)
		points = append(points, image.Point{int(xmax), int(intercept)})
	}

	if len(points) == 2 {
		return points
	}
	// x intercept at lower y bound
	invGradient := 1.0 / gradient
	fmt.Printf("invGradient: %.3f\n", invGradient)
	intercept = x - ((y - ymin) * invGradient)
	if intercept >= xmin && intercept <= xmax {
		fmt.Printf("low y, x intercept %3f\n", intercept)
		points = append(points, image.Point{int(intercept), int(ymin)})
	}
	if len(points) == 2 {
		return points
	}
	// x intercept at upper y bound
	intercept = x + ((ymax - y) * invGradient)
	if intercept >= ymin && intercept <= ymax {
		fmt.Printf("high y, x intercept %3f\n", intercept)
		points = append(points, image.Point{int(intercept), int(ymax)})
	}
	return points
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
func (c *context) MoveTo(pt image.Point) {
	c.p = pt
}

func (c *context) LineTo(pt image.Point) {
	c.ops = append(c.ops, lineTo{c.p, pt, c.c})
}

func (c *context) Line(d, a float64) {
	c.ops = append(c.ops, line{d, a, c.c})
}

func (c *context) Render(im draw.Image) {
	for _, op := range c.ops {
		op.do(im)
	}
}

func DrawLine(p1, p2 image.Point, col color.Color, im draw.Image) {
	fmt.Printf("Drawing line %v %v \n", p1, p2)

	// determine which way to draw the line
	// top to bottom or left to right
	xdiff := p2.X - p1.X
	ydiff := p2.Y - p1.Y
	//fmt.Printf("xdiff: %d, ydiff: %d\n", xdiff, ydiff)
	var yInc, xInc float64
	var iter, inc int
	//determine which direction the line is longest in
	if math.Abs(float64(xdiff)) >= math.Abs(float64(ydiff)) {
		iter = xdiff
		// increment either 1 or -1 based on the
		// value of xdiff
		inc = int(int32(uint32(xdiff)&uint32(0x80000000))>>31 | 1)
		xInc = float64(inc)
		//fmt.Printf("xinc %2f\n", xInc)
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
}
