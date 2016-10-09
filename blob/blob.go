// Package blob provides functions for finding blobs (connected pixel regions) within an image.
package blob

import (
	"image"
	"image/color"
	"math"

	"github.com/piersy/hough-go/gray16"
	"github.com/piersy/hough-go/point"
)

var (
	// BlobColor is the color a pixel must have to be considered part of a
	// blob. By default it is white.
	BlobColor = color.Gray16{math.MaxUint16}
)

type Blob struct {
	points []image.Point
}

func (b *Blob) Centre() point.Point {
	var totX, totY float64
	numPoints := len(b.points)
	for _, p := range b.points {
		totX += float64(p.X)
		totY += float64(p.Y)
	}
	return point.Point{totX / float64(numPoints), totY / float64(numPoints)}
}

// Find finds the blobs in an image. The input image is searched for connected
// regions of color equal to BlobColor. Pixels belonging to regions are put
// into Blobs and a slice of blobs is returned.
func Find(i *gray16.Gray16) []*Blob {
	found := make(map[image.Point]struct{})
	var blobs []*Blob

	b := i.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			p := image.Point{x, y}
			// If the point has not already been found and is of the right color
			if _, ok := found[p]; !ok && i.Gray16At(p.X, p.Y) == BlobColor {
				bl := &Blob{}
				findConnected(i, p, bl, found)
				blobs = append(blobs, bl)
			}
		}
	}
	return blobs
}

// findConnected checks to see that the given point is of BlobColor and is
// within bounds if not it returns. Otherwise it adds it to the given blob and
// the given map and then calls itself for the four adjacent pixels.
func findConnected(i *gray16.Gray16, p image.Point, b *Blob, found map[image.Point]struct{}) {

	if _, ok := found[p]; ok || i.Gray16At(p.X, p.Y) != BlobColor || !(p.In(i.Bounds())) {
		return
	}
	// Add the point to found
	found[p] = struct{}{}
	// Add the point to the blob
	b.points = append(b.points, p)

	findConnected(i, image.Point{p.X, p.Y - 1}, b, found)
	findConnected(i, image.Point{p.X, p.Y + 1}, b, found)
	findConnected(i, image.Point{p.X - 1, p.Y}, b, found)
	findConnected(i, image.Point{p.X + 1, p.Y}, b, found)
}
