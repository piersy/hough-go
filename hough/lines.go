package hough

import (
	"fmt"
	"math"
)

type Line struct {
	angle, distance float64
}

func (l Line) String() string {
	return fmt.Sprintf("Angle: %.2f Distance: %.2f", l.angle*180/math.Pi, l.distance)
}

type topLines struct {
	l         []lineAndScore
	boundary  uint16
	lineCount int
}

type lineAndScore struct {
	line  Line
	score uint16
}

func newLineAndScore(angle, distance float64, score uint16) lineAndScore {
	return lineAndScore{
		line: Line{
			angle:    angle,
			distance: distance,
		},
		score: score,
	}
}

// addLine adds a line to the list. The caller must make sure that any line
// added has a score greater than the current boundary score available via the
// boundary field of the struct.
func (t *topLines) addLine(angle, distance float64, score uint16) {
	// if the map is at capacity and the new index
	// is not part of the map then we need to delete an element
	newLine := newLineAndScore(angle, distance, score)
	if len(t.l) == t.lineCount {
		t.l[t.minElement()] = newLine
	} else {
		t.l = append(t.l, newLine)
	}
	//If we are at max capacity set the boundary
	if len(t.l) == t.lineCount {
		t.boundary = t.l[t.minElement()].score
	}

}

// minElement returns the index of the line having the minimum score
func (t *topLines) minElement() int {
	//set minscore to be the highest possible gray16 value
	minscore := uint16(math.MaxUint16)
	var minIndex int
	for i, v := range t.l {
		if v.score < minscore {
			minscore = v.score
			minIndex = i
		}
	}
	return minIndex
}

func (t *topLines) lines() []Line {
	res := make([]Line, len(t.l))
	for i, _ := range t.l {
		res[i] = t.l[i].line
	}
	return res
}

//func TopLines(lineCount int, hough image.Gray16) []Line {
//	b := hough.Bounds()
//	tl := &topLines{}
//	for y := b.Min.Y; y < b.Max.Y; y++ {
//		for x := b.Min.X; x < b.Max.X; x++ {
//			score := hough.Gray16At(x, y).Y
//			if score > tl.boundary {
//				tl.addLine(x, y, score)
//			}
//		}
//	}
//	return tl.lines()
//}

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
