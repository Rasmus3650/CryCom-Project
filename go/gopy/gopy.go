package gopy

import (
	"fmt"
	"sync"
)

func Ident(n int) Matrix {
	return Diag(n, 1)
}

func Diag(n int, v T) Matrix {
	if v == 0 {
		return MakeMatrix(n, n)
	}

	m := make(Matrix, n)
	for i := 0; i < n; i++ {
		row := make(Arr, n)
		row[i] = v
		m[i] = row
	}
	return m
}

func AddEntrywise(l Arr, con T) {
	for i := range l {
		l[i] += con
	}
}
func AddMatrix(l Matrix, r Matrix) Matrix {
	h, w := Shape(l)
	h0, w0 := Shape(l)
	if h != h0 && w != w0 {
		msg := fmt.Sprintf("(%d,%d)(%d,%d)", h, w, h0, w0)
		panic("cannot add matrices of size " + msg)
	}
	for i, row := range l {
		for j := range row {
			l[i][j] += r[i][j]
		}
	}
	return l
}
func AddMatrixSlim(l MatrixSlim, r MatrixSlim) MatrixSlim {
	h, w := ShapeSlim(l)
	h0, w0 := ShapeSlim(r)
	if h != h0 && w != w0 {
		msg := fmt.Sprintf("(%d,%d)(%d,%d)", h, w, h0, w0)
		panic("cannot add matrices of size " + msg)
	}
	for i, row := range l {
		for j := range row {
			l[i][j] += r[i][j]
		}
	}
	return l
}

func SubMatrix(l Matrix, r Matrix) Matrix {
	h, w := Shape(l)
	h0, w0 := Shape(l)
	if h != h0 && w != w0 {
		msg := fmt.Sprintf("(%d,%d)(%d,%d)", h, w, h0, w0)
		panic("cannot add matrices of size " + msg)
	}
	for i, row := range l {
		for j := range row {
			l[i][j] -= r[i][j]
		}
	}
	return l
}

func AddMatrixMod(l Matrix, r Matrix, m T) Matrix {
	h, w := Shape(l)
	h0, w0 := Shape(l)
	if h != h0 && w != w0 {
		msg := fmt.Sprintf("(%d,%d)(%d,%d)", h, w, h0, w0)
		panic("cannot add matrices of size " + msg)
	}
	for i, row := range l {
		for j := range row {
			l[i][j] = (l[i][j] + r[i][j]) % m
		}
	}
	return l
}

func ModMatrix(M Matrix, mod T) Matrix {
	h, w := Shape(M)
	M_ := MakeMatrix(h, w)
	return AddMatrixMod(M, M_, mod)
}

func Dot(l Arr, r Arr) Matrix {
	return Matrix{}
}
func DotMatrix(l Matrix, r Matrix) Matrix {
	if len(l) < 1 && len(l[0]) < 1 {
		panic("cant dot empty matrix")
	}
	if len(l[0]) != len(r) {
		panic("cant dot unequal matrices")

	}

	M := MakeMatrix(len(l), len(r[0]))
	wg := sync.WaitGroup{}
	wg.Add(len(M))

	for x, a := range M {
		go func(x int) {
			for y := range a {
				var res T
				for j := range r {
					res += l[x][j] * r[j][y]
				}
				M[x][y] = res
			}
			wg.Done()
		}(x)
	}
	wg.Wait()
	return M
}

func DotMatrixMod(l Matrix, r Matrix, m T) Matrix {
	if len(l) < 1 && len(l[0]) < 1 {
		panic("cant dot empty matrix")
	}
	if len(l[0]) != len(r) {
		panic("cant dot unequal matrices")

	}
	wg := sync.WaitGroup{}

	M := MakeMatrix(len(l), len(r[0]))
	wg.Add(len(M))

	for x, a := range M {
		go func() {
			defer wg.Done()
			for y := range a {
				var res T
				for j := range r {
					res += l[x][j] * r[j][y]
				}
				M[x][y] = res % m
			}
		}()
	}
	wg.Wait()
	return M
}

func DotMatrixSlim(l, r MatrixSlim, m T) Matrix {
	if len(l) < 1 && len(l[0]) < 1 {
		panic("cant dot empty matrix")
	}
	if len(l[0]) != len(r) {
		panic("cant dot unequal matrices")

	}
	wg := sync.WaitGroup{}

	M := MakeMatrix(len(l), len(r[0]))
	wg.Add(len(M))

	for x, a := range M {
		go func(x int, a Arr) {
			defer wg.Done()
			for y := range a {
				var res T
				for j := range r {
					res += T(l[x][j]) * T(r[j][y]) % m
				}
				M[x][y] = res % m
			}
		}(x, a)
	}
	wg.Wait()
	return M
}

func InnerProdMod(lhs, rhs Arr, mod T) T {
	var x T
	for i, e := range lhs {
		x += (e * rhs[i]) % mod
	}
	return x
}

func Add(l Arr, r Arr) Arr {
	return Arr{}
}
func Scale(arr Arr, s T) Arr {
	for i := range arr {
		arr[i] *= s
	}
	return arr
}

func ScaleMod(arr Arr, s, m T) Arr {
	for i, e := range arr {
		arr[i] = (e * s) % m
	}
	return arr
}

func Equal(lhs Matrix, rhs Matrix) bool {
	h1, w1 := Shape(lhs)
	h2, w2 := Shape(rhs)
	if h1 != h2 || w1 != w2 {
		return false
	}
	for i, e := range lhs {
		for j, x := range e {
			if x != rhs[i][j] {
				return false

			}
		}
	}
	return true
}
func LineEqual(lhs Matrix, rhs Matrix) {
	h1, w1 := Shape(lhs)
	h2, w2 := Shape(rhs)

	lines := make([]bool, h1)

	if h1 != h2 || w1 != w2 {
		fmt.Println(lines)
		return
	}
	for i, e := range lhs {
		v := true
		for j, x := range e {
			if x != rhs[i][j] {
				v = false
			}
		}
		lines[i] = v
	}
	fmt.Println(lines)
}
