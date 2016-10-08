package util

type Normaliser struct {
	ratio, srcMin, dstMin float64
}

func NewNormaliser(srcMin, srcMax, dstMin, dstMax float64) Normaliser {
	return Normaliser{
		ratio:  (dstMax - dstMin) / (srcMax - srcMin),
		srcMin: srcMin,
		dstMin: dstMin,
	}
}

func (n Normaliser) Normalise(val float64) float64 {
	return n.ratio*(val-n.srcMin) + n.dstMin
}
