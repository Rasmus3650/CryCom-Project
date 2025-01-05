package big

import (
	gp "crycomproj/gopy"
	"math"
	"math/rand"
	"sync"
)

func Enc(params Params, mu gp.T, pk gp.Matrix) gp.Matrix {
	wg := sync.WaitGroup{}
	R := gp.MakeMatrix(params.N, params.m)
	wg.Add(len(R))
	for i, row := range R {
		go func() {
			for j := range row {
				R[i][j] = gp.T(rand.Intn(2))
			}
			wg.Done()
		}()

	}
	wg.Wait()

	rhs := BitDecompMatrix(gp.DotMatrix(R, pk), params.ell)

	for i := 0; i < params.N; i++ {
		rhs[i][i] += mu
	}

	C := FlattenMatrix(rhs, params.n+1, params.ell, params.qt)
	h, w := gp.Shape(C)
	if h != w || h != params.N {
		panic("bad size in enc")
	}
	return C
}

func Dec(params Params, v gp.Arr, C gp.Matrix) gp.T {
	var lb float64 = float64(params.q) / 4
	var ub float64 = float64(params.q) / 2

	for j := 0; j < params.ell; j++ {
		fe := float64(v[j])
		if fe <= ub && fe > lb {
			vi := float64(v[j])
			dot := float64(gp.InnerProdMod(C[j], v, params.qt))

			return gp.T(math.Round(dot/vi)) & 1
		}
	}
	panic("no vi within bounds found")
}

func MPDec(params Params, v gp.Arr, C gp.Matrix) gp.T {
	leng := params.ell - 1

	w := make(gp.Arr, leng)

	for i := 0; i < leng; i++ {
		var x gp.T
		row := C[i]
		for j := 0; j < len(v); j++ {
			x += (row[j] * v[j]) % params.qt
			x = x % params.qt
		}

		w[i] = x
	}

	res := make(gp.ArrSlim, leng)

	// Minus 2 because 2^ell-1 is a index ell-2
	count := 0
	for i := params.ell - 2; i >= 0; i-- {
		x := math.Round(float64(w[i]) / float64(params.PO2[i]))
		mask := gp.T(1 << count)
		res[count] = gp.S((gp.T(x) & mask) >> count)

		for j := 0; j < i; j++ {
			w[j] -= (gp.T(res[count]) << (count))

			for w[j] < 0 {
				w[j] = params.qt - w[j]
			}
			w[j] %= params.qt

		}

		count++
	}

	var x gp.T
	for i := len(res) - 1; i >= 0; i-- {
		x += gp.T(res[i])
		x <<= 1
	}
	x >>= 1

	return x % params.qt
}
