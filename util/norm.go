package util

type noramliser struct {
	ratio, srcMin, dstMin float64
}

func NewNormaliser(srcMin, srcMax, dstMin, dstMax float64) noramliser {
	return noramliser{
		ratio:  (dstMax - dstMin) / (srcMax - srcMin),
		srcMin: srcMin,
		dstMin: dstMin,
	}
}

func (n noramliser) Normalise(val float64) float64 {
	return n.ratio*(val-n.srcMin) + n.dstMin
}
