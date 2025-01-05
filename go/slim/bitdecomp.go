package slim

import (
	"crycomproj/big"
	gp "crycomproj/gopy"
	"fmt"
	"sync"
)

func BitDecompAsync(row gp.Arr, ch chan gp.ArrSlim, ell int) gp.ArrSlim {
	A := <-ch
	for i, a := range row {
		index := i * ell
		for y := 0; y < ell; y++ {
			A[index+y] = gp.S((a >> y) & 1)
		}
	}

	ch <- make(gp.ArrSlim, len(A))
	return A
}

func BitDecompMatrix(M gp.Matrix, ell int) gp.MatrixSlim {
	limit := 20
	mlen := len(M)
	M_ := make(gp.MatrixSlim, mlen)
	wg := sync.WaitGroup{}
	wg.Add(len(M))
	arrays := make(chan gp.ArrSlim, limit)
	for i, row := range M {
		go func(i int, row gp.Arr) {
			defer wg.Done()
			M_[i] = BitDecompAsync(row, arrays, ell)
		}(i, row)
	}

	for i := 0; i < limit; i++ {
		arrays <- make(gp.ArrSlim, len(M[0])*ell)
	}
	wg.Wait()

	return M_
}

func FlattenMatrix(M gp.MatrixSlim, k, ell int, mod gp.T) gp.MatrixSlim {
	return BitDecompMatrix(
		InvBitDecompMatrix(M, k, ell, mod), ell,
	)
}

// Flatten when input can exceed limit of gopy.S type
func FlattenMatrix2(M gp.Matrix, k, ell int, mod gp.T) gp.MatrixSlim {
	return BitDecompMatrix(
		big.InvBitDecompMatrix(M, k, ell, mod), ell,
	)
}

func InvBitDecomp(A gp.ArrSlim, k int, ell int, mod gp.T) gp.Arr {
	A_ := make(gp.Arr, k)
	if len(A) != k*ell {
		fmt.Println(len(A), k*ell)
		panic("wrong parameters for invbitdecomp")
	}
	for i := 0; i < k; i++ {
		var x gp.T
		index := i * ell
		for j := ell - 1; j >= 0; j-- {
			x += gp.T(A[index+j])
			x <<= 1
		}
		x >>= 1
		A_[i] = x % (mod << 1)
	}
	//gp.OutArr(A_)
	//panic("fuck")
	return A_
}

func InvBitDecompMatrix(M gp.MatrixSlim, k int, ell int, mod gp.T) gp.Matrix {
	lenM := len(M)
	M_ := make(gp.Matrix, lenM)
	wg := sync.WaitGroup{}
	wg.Add(len(M))
	for i, row := range M {
		go func(t int) {
			M_[t] = InvBitDecomp(row, k, ell, mod)
			wg.Done()
		}(i)
	}
	wg.Wait()
	return M_
}
