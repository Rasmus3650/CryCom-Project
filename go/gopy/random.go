package gopy

import (
	"math/rand"
)

func SampleArr(lb int, ub int, n int) Arr {
	a := make(Arr, n)
	for i := range a {
		a[i] = T(rand.Intn(ub-lb) - lb)
	}
	return a
}

func SampleMatrix(lb int, ub int, n int, m int) Matrix {
	a := MakeMatrix(n, m)
	for i, row := range a {
		for j := range row {
			a[i][j] = T(rand.Intn(ub-lb) + lb)
		}
	}
	return a
}
