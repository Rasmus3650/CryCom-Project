package slimmer

import (
	gp "crycomproj/gopy"
	"math"
	"sync"
)

func InvBitDecompMasked(A gp.Arr, k, ell int, mod, mask gp.T) gp.Arr {
	A_ := make(gp.Arr, k)
	if len(A) != k*ell {
		println(len(A))
		panic("wrong parameters for invbitdecomp")
	}
	for i := 0; i < k; i++ {
		var x gp.T
		index := i * ell
		for j := ell - 1; j >= 0; j-- {
			x += A[index+j]
			x <<= 1
		}
		x >>= 1
		A_[i] = (x % mod) & mask
	}
	return A_
}

func FlattenMatrix(M gp.Matrix, k, ell int, mod gp.T) gp.Matrix {
	mask := (gp.T(math.Pow(2, float64(ell))) << 1) - 1
	lenM := len(M)
	M_ := make(gp.Matrix, lenM)
	wg := sync.WaitGroup{}
	wg.Add(lenM)
	for i, row := range M {
		go func() {
			M_[i] = InvBitDecompMasked(row, k, ell, mod, mask)
			wg.Done()
		}()
	}
	wg.Wait()
	return M_
}
