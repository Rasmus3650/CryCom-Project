package gopy

import "fmt"

func TransposeArr(M Matrix) Arr {
	a := make(Arr, len(M))
	for i, e := range M {
		a[i] = e[0]
	}
	return a
}

func TransposeColumnMatrix(arr Arr) Matrix {
	m := MakeMatrix(len(arr), 1)
	for i, e := range arr {
		m[i][0] = e
	}
	return m
}

func MakeMatrix(rows int, columns int) Matrix {
	a := make([][]T, rows)
	for i := range a {
		a[i] = make([]T, columns)
	}
	return a
}

func MakeMatrixSlim(rows int, columns int) MatrixSlim {
	a := make([][]S, rows)
	for i := range a {
		a[i] = make([]S, columns)
	}
	return a
}

func Shape(M Matrix) (int, int) {
	h, w := len(M), 0
	if len(M[0]) < 1 {
		return h, w
	}
	w = len(M[0])
	/*
		for _, row := range M {
			if len(row) != w {
				panic("not equal rows")
			}
		}
	*/
	return h, w
}

func Out(M Matrix) {
	msg := ""
	for _, m := range M {
		msg += fmt.Sprintf("%v\n", m)

	}
	fmt.Println(msg)
}

func OutX(M [][]string) {
	msg := ""
	for _, m := range M {
		msg += fmt.Sprintf("%v\n", m)

	}
	fmt.Println(msg)
}

func OutSlim(M MatrixSlim) {
	msg := ""
	for _, m := range M {
		msg += fmt.Sprintf("%v\n", m)

	}
	fmt.Println(msg)
}

func OutArr(A Arr) {
	fmt.Printf("%v\n", A)
}
func OutArrSlim(A ArrSlim) {
	fmt.Printf("%v\n", A)
}
