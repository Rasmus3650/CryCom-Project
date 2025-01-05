package big

import (
	gp "crycomproj/gopy"
	"sync"
)

func BitDecomp(A gp.Arr, ell int) gp.Arr {

	A_ := make(gp.Arr, len(A)*ell)
	for i, a := range A {
		index := i * ell
		for j := 0; j < ell; j++ {
			A_[index+j] = (a >> j) & 1
		}
	}

	return A_
}

func BitDecompMatrix(M gp.Matrix, ell int) gp.Matrix {
	M_ := make(gp.Matrix, len(M))
	wg := sync.WaitGroup{}
	wg.Add(len(M))
	for i, row := range M {
		go func(j int) {
			M_[j] = BitDecomp(row, ell)
			wg.Done()
		}(i)
	}
	wg.Wait()
	return M_
}

func InvBitDecomp(A gp.Arr, k int, ell int, mod gp.T) gp.Arr {
	A_ := make(gp.Arr, k)
	if len(A) != k*ell {
		println(len(A), k*ell)
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
		A_[i] = x % mod
	}
	return A_
}

func InvBitDecompMatrix(M gp.Matrix, k int, ell int, mod gp.T) gp.Matrix {
	lenM := len(M)
	M_ := make(gp.Matrix, lenM)
	wg := sync.WaitGroup{}
	wg.Add(lenM)
	for i, row := range M {
		go func(i int) {
			M_[i] = InvBitDecomp(row, k, ell, mod)
			wg.Done()
		}(i)
	}
	wg.Wait()
	return M_
}

func Flatten(a gp.Arr, k, ell int, mod gp.T) gp.Arr {
	return BitDecomp(InvBitDecomp(a, k, ell, mod), ell)
}

func FlattenMatrix(M gp.Matrix, k, ell int, mod gp.T) gp.Matrix {
	return BitDecompMatrix(
		InvBitDecompMatrix(M, k, ell, mod), ell,
	)
}
