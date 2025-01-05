package slimmer

import (
	gp "crycomproj/gopy"
	"math"
	"sync"
)

func Enc(params Params, mu gp.T, pk gp.Matrix) gp.Matrix {
	var C gp.Matrix = gp.MakeMatrix(params.N, len(pk[0]))
	wg := sync.WaitGroup{}
	dot := func() {
		wg.Add(params.N)

		for x := 0; x < params.N; x++ {
			go func() {
				r := gp.SampleArr(0, 2, params.m)
				for y := 0; y < params.n+1; y++ {
					var res gp.T
					for j := 0; j < params.m; j++ {
						res += r[j] * pk[j][y] % params.qt
					}
					C[x][y] = res
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}

	dot()

	for i := 0; i < len(C); i++ {
		col := i / params.ell
		val := 1 << (i % params.ell)
		C[i][col] += gp.T(val) * mu
	}

	return C
}

func Dec(params Params, v gp.Arr, C gp.Matrix) gp.T {
	inner := func(i int) float64 {
		var res gp.T
		for j, vi := range v {
			col := j / params.ell
			shift := j % params.ell
			val := C[i][col]

			res += ((val >> shift) & 1) * vi
		}
		return float64(res)
	}

	var lb float64 = float64(params.q) / 4
	var ub float64 = float64(params.q) / 2

	for j := 0; j < params.ell; j++ {
		fe := float64(v[j])
		if fe <= ub && fe > lb {
			vi := float64(v[j])
			dot := inner(j)

			return gp.T(math.Round(dot/vi)) & 1
		}
	}
	panic("no vi within bounds found")
}

func MPDec2(params Params, v gp.Arr, C gp.Matrix) gp.T {
	leng := params.ell - 1

	w := make(gp.Arr, leng)

	for i := 0; i < leng; i++ {
		var x gp.T
		for j := 0; j < len(C[i]); j++ {
			elem := C[i][j]
			for j2 := 0; j2 < params.ell; j2++ {
				x += ((elem & 1) * v[j*params.ell+j2]) % params.qt

				x = x % params.qt

				elem >>= 1
			}
		}

		w[i] = x
	}

	res := make(gp.Arr, leng)

	// Minus 2 because 2^ell-1 is a index ell-2
	count := 0
	for i := params.ell - 3; i >= 0; i-- {
		x := math.Round(float64(w[i]) / float64(params.PO2[i]))
		mask := gp.T(1 << count)
		res[count] = (gp.T(x) & mask) >> count

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


