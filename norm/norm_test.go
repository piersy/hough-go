package norm

import "testing"

func TestNorm(t *testing.T) {

	srcMax := 5.
	srcMin := -5.
	max := 10.
	min := -10.
	n := NewNormaliser(srcMin, srcMax, min, max)
	if n.Normalise(0) != 0 {
		t.Fatal()
	}
	if n.Normalise(5) != 10 {
		t.Fatal()
	}
	if n.Normalise(-5) != -10 {
		t.Fatal()
	}
	if n.Normalise(-2.5) != -5 {
		t.Fatal()
	}
	if n.Normalise(2.5) != 5 {
		t.Fatal()
	}
}
