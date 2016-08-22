package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"

	"github.com/piersy/hough-go/canvas"
	"github.com/piersy/hough-go/hough"
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
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	baseImage, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	acc, lines := hough.Hough(baseImage, 400, 400)
	for _, l := range lines {
		fmt.Println(l)
	}
	outFile, err := os.Create(*out)
	defer outFile.Close()
	if err != nil {
		println(err)
		os.Exit(1)
	}
	png.Encode(outFile, acc)

	testOut, err := os.Create("testout.png")
	defer testOut.Close()
	if err != nil {
		println(err)
		os.Exit(1)
	}
	//ctx := canvas.New()
	//ctx.Move(image.Pt(0, 100))
	//ctx.Line(image.Pt(300, 100))
	//ctx.Render(im)
	png.Encode(testOut, baseImage)
}

type Line interface {
	Draw(draw.Image) error
}

type line struct {
	angle, distance float64
}

func (l line) Draw(i draw.Image) error {
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()
	midX := float64(width) / 2
	midY := float64(height) / 2

	points := make([]image.Point, 0)
	//Special cases line is horizontal or vertical
	if l.angle == 0 || l.angle == math.Pi {
		//line is vertical
		points = append(points, image.Pt(int(midX+l.distance), 0))
		points = append(points, image.Pt(int(midX+l.distance), height))
	}
	if l.angle == math.Pi/2 {
		//line is vertical
		points = append(points, image.Pt(0, int(midY+l.distance)))
		points = append(points, image.Pt(width, int(midY+l.distance)))
	}

	// d = x cos t + ysin t
	// At left edge
	//   d = -midX cos t + y sin t
	// → d + (-midX cos t) = y sin t
	// → (d + (-midX cos t))/sin t = y
	// see if y is in the image bounds
	// if so we crossed an edge
	y := (l.distance + -midX*math.Cos(l.angle)) / math.Sin(l.angle)
	if y >= -midY && y <= midY {
		points = append(points, image.Pt(0, int(y+midY)))
	}
	y = (l.distance + midX*math.Cos(l.angle)) / math.Sin(l.angle)
	if y >= -midY && y <= midY {
		points = append(points, image.Pt(width, int(y+midY)))
	}
	x := (l.distance + -midY*math.Sin(l.angle)) / math.Cos(l.angle)
	if x >= -midX && x <= midX {
		points = append(points, image.Pt(int(x+midX), 0))
	}
	x = (l.distance + midY*math.Sin(l.angle)) / math.Cos(l.angle)
	if x >= -midX && x <= midX {
		points = append(points, image.Pt(int(x+midX), height))
	}
	if len(points) != 2 {

		fmt.Errorf("Something went wrong %d points found", len(points))
	}
	ctx := canvas.New()
	ctx.MoveTo(points[0])
	ctx.LineTo(points[1])
	err := ctx.Render(i)
	if err != nil {
		return err
	}
	return nil

}

//func hough(i image.Image, lineCount int) (image.Image, []Line, error) {
//	// for logging only
//	minDist := math.Pow(2, 16) - 1
//	// for logging only
//	maxDist := float64(0)
//
//	width := i.Bounds().Dx()
//	height := i.Bounds().Dy()
//	midX := float64(width) / 2
//	midY := float64(height) / 2
//
//	maxDistance := math.Sqrt(float64(width*width+height*height)) / 2
//	println(maxDistance)
//
//	dstSize := 400
//	hough := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{dstSize, dstSize}})
//
//	lines := make(map[int]lineAndScore)
//	var boundaryScore uint16
//
//	for x := 0; x < width; x++ {
//		for y := 0; y < height; y++ {
//			r, _, _, _ := i.At(x, y).RGBA()
//
//			px := float64(x) - midX
//			py := float64(y) - midY
//			//black pixel
//			if r == 0 {
//				//iterate the angles
//				for t := 0; t < dstSize; t++ {
//					angle := normalise(0, float64(dstSize), 0, math.Pi, float64(t))
//					//Get normal distance can be negative
//					distance := px*math.Cos(angle) + py*math.Sin(angle)
//					if distance < minDist {
//						minDist = distance
//					}
//					if distance > maxDist {
//						maxDist = distance
//					}
//					nAngle := int(normalise(0, math.Pi, 0, float64(dstSize), angle) + 0.5)
//					nDistance := int(normalise(-maxDistance, maxDistance, 0, float64(dstSize), distance) + 0.5)
//					g := hough.Gray16At(nAngle, nDistance)
//					g.Y += 100
//					hough.SetGray16(nAngle, nDistance, g)
//					//Update our lines result if the score is good
//					if g.Y > boundaryScore {
//						// if the map is at capacity and the new index
//						// is not part of the map then we need to delete an element
//						//Generate linear index from x and y
//						index := nAngle*nDistance + nAngle
//						_, ok := lines[index]
//						if !ok && len(lines) == lineCount {
//							var index int
//							//set minscore to be the highest possible gray16 value
//							minscore := uint16(math.Pow(2, 16) - 1)
//
//							for k, v := range lines {
//								if v.score < minscore {
//									minscore = v.score
//									index = k
//								}
//							}
//							//delete entry with lowest score
//							delete(lines, index)
//						}
//						lines[index] = lineAndScore{
//							l: line{
//								angle:    angle,
//								distance: distance,
//							},
//							score: g.Y,
//						}
//
//						if len(lines) == lineCount {
//							//set minscore to be the highest possible gray16 value
//							minscore := uint16(math.Pow(2, 16) - 1)
//
//							for k, v := range lines {
//								if v.score < minscore {
//									minscore = v.score
//									index = k
//								}
//							}
//							//Update boundary score to be lowest score in set
//							boundaryScore = lines[index].score
//						}
//
//					}
//
//				}
//			}
//		}
//	}
//	fmt.Printf("Min dst: %v, scaled:%v\n", minDist, normalise(0, maxDistance, 0, float64(dstSize), minDist))
//	fmt.Printf("Min dst: %v, scaled:%v\n", minDist, normalise(-maxDistance, maxDistance, 0, float64(dstSize), minDist))
//	fmt.Printf("Max dst: %v, scaled:%v\n", maxDist, normalise(0, maxDistance, 0, float64(dstSize), maxDist))
//	resultLines := make([]Line, lineCount)
//	c := 0
//	for k, v := range lines {
//		resultLines[c] = v.l
//		c++
//		fmt.Printf("Index %d Score: %d  Angle: %.2f Distance: %.2f\n", k, v.score, v.l.angle, v.l.distance)
//	}
//	for _, l := range resultLines {
//		fmt.Printf("Angle: %.2f Distance: %.2f\n", l.(line).angle, l.(line).distance)
//	}
//	return hough, resultLines, nil
//	//r = x cos T + y sin T
//}
func normalise(srcMin, srcMax, dstMin, dstMax, srcVal float64) float64 {
	return ((dstMax - dstMin) * (srcVal - srcMin) / (srcMax - srcMin)) + dstMin
}
