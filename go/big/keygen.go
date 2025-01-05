package big

import (
	gp "crycomproj/gopy"
	"crycomproj/util"
	"math/rand"
	"sync"
)

func SecretKeyGen(params Params) (gp.Arr, gp.Matrix, gp.Arr) {
	sk := make(gp.Arr, params.n+1)
	sk[0] = 1
	t := gp.MakeMatrix(params.n, 1)
	for i := 1; i < params.n+1; i++ {
		ti := gp.T(rand.Intn(params.q))
		if ti == 0 {
			continue
		}

		sk[i] = params.qt - ti
		t[i-1][0] = ti
	}

	v := util.ApplyPowersOf2(params.q, sk, params.PO2)
	return sk, t, v
}

func PublicKeyGen(params Params, s gp.Arr, t gp.Matrix) gp.Matrix {
	var B gp.Matrix
	var e gp.Matrix
	var b gp.Matrix
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		B = gp.SampleMatrix(0, params.q, params.m, params.n)
		wg.Done()
	}()
	go func() {
		e = gp.SampleMatrix(-params.ex, params.ex+1, params.m, 1)
		wg.Done()
	}()
	wg.Wait()

	b = gp.DotMatrix(B, t)
	b = gp.AddMatrixMod(b, e, params.qt)
	_, w := gp.Shape(b)
	if w != 1 {
		panic("prod vector containing more columns")
	}

	A := gp.MakeMatrix(params.m, params.n+1)
	wg.Add(len(B))
	for i, row := range B {
		go func() {
			A[i][0] = b[i][0]
			for j, entry := range row {
				A[i][j+1] = entry
			}
			wg.Done()
		}()
	}
	wg.Wait()

	/*
	e_ := gp.DotMatrixMod(A, gp.TransposeColumnMatrix(s), params.qt)
  if !Equal(e_, e) {
    panic("A dot s is not gp.equal e")
  }
	*/

	return A
}
