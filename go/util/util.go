package util

import (
	gp "crycomproj/gopy"
	"math"
)

func GeneratePowersOf2(ell int, q gp.T) gp.Arr {
	a := make(gp.Arr, ell)
	for i := 0; i < ell; i++ {
		a[i] = gp.T(math.Pow(2, float64(i))) % q
	}
	return a
}

func ApplyPowersOf2(q int, b gp.Arr, po2 gp.Arr) gp.Arr {

	b2 := make(gp.Arr, len(po2)*len(b))

	for i := 0; i < len(b); i++ {
		index := i * len(po2)
		for j := 0; j < len(po2); j++ {
			var inter int64 = (int64(b[i]) * int64(po2[j])) % int64(q)
			b2[index+j] = gp.T(inter)
		}
	}
	return b2
}

func GetEx(q float64) int {
	return int(math.Sqrt(q))
}

func GetQLogAndEll(q float64) (int, int, int) {

	qlog2 := math.Log2(q)
	qlog2_f := int(math.Floor(qlog2))
	qlog2_c := int(math.Ceil(qlog2))
	ell := qlog2_f + 1
	return qlog2_f, qlog2_c, ell
}


