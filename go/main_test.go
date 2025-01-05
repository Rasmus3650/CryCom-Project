package main

import (
	gp "crycomproj/gopy"
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkFullBench1(b *testing.B) {
	run1()
}

func BenchmarkFullBench2(b *testing.B) {
	run3()
}

func BenchmarkRandDecomp1(b *testing.B) {
	M := gp.SampleMatrix(0, 100, 5000, 5000)
	BitDecompMatrix(M, 10)
}

func BenchmarkSmallDry1(b *testing.B) {
	var q float64 = 134217728
	n := 1024
	m := 1024

	var mu gp.T = gp.T(rand.Intn(2))

	params := Setup(q, n, m)
	sk, t, v := SecretKeyGen(params)
	pk := PublicKeyGen(params, sk, t)

	C := Enc(params, mu, pk)
	mu_ := Dec(params, v, C)
	C = nil

	if mu != mu_ {
		panic("fuck")
	}
}

func BenchmarkSmallDry2(b *testing.B) {
	var q float64 = 134217728
	n := 1024
	m := 1024

	var mu gp.T = gp.T(rand.Intn(2))

	params := Setup(q, n, m)
	sk, t, v := SecretKeyGen(params)
	pk := PublicKeyGen(params, sk, t)

	C := EncSlimmer(params, mu, pk)
	mu_ := DecSlimmer(params, v, C)

	if mu != mu_ {
		panic("fuck")
	}
}

func BenchmarkSmallDry3(b *testing.B) {
	var q float64 = 134217728
	n := 1024
	m := 1024

	var mu gp.T = gp.T(rand.Intn(2))

	params := Setup(q, n, m)
	sk, t, v := SecretKeyGen(params)
	pk := PublicKeyGen(params, sk, t)

	C := EncSlimmer2(params, mu, pk)
	mu_ := DecSlimmer2(params, v, C)

	if mu != mu_ {
		panic("fuck")
	}
}


func BenchmarkSmall0(b *testing.B) {
	var q float64 = 134217728
	n := 256
	m := 256

	var mu gp.T = gp.T(rand.Intn(2))

	params := Setup(q, n, m)
	sk, t, v := SecretKeyGen(params)
	pk := PublicKeyGen(params, sk, t)

	C := Enc(params, mu, pk)
	mu_ := Dec(params, v, C)

	if mu != mu_ {
		panic("fuck")
	}
}

func BenchmarkSmall1(b *testing.B) {
	var q float64 = 134217728
	n := 512
	m := 512

	var mu gp.T = gp.T(rand.Intn(2))

	params := Setup(q, n, m)
	sk, t, v := SecretKeyGen(params)
	pk := PublicKeyGen(params, sk, t)
	C := Enc(params, mu, pk)

	mu_ := Dec(params, v, C)

	if mu != mu_ {
		panic("fuck")
	}
}
func BenchmarkSmall2(b *testing.B) {
	var q float64 = 134217728
	n := 512
	m := 512

	var mu gp.T = gp.T(rand.Intn(2))

	params := Setup(q, n, m)
	sk, t, v := SecretKeyGen(params)
	pk := PublicKeyGen(params, sk, t)
	C := Enc(params, mu, pk)

	mu_ := Dec(params, v, C)

	if mu != mu_ {
		panic("fuck")
	}
}

func BenchmarkMainDry(b *testing.B) {
	var q float64 = 134217728
	n := 1024
	m := 1024

	var mu gp.T = gp.T(rand.Intn(2))

	params := Setup(q, n, m)
	sk, t, v := SecretKeyGen(params)
	pk := PublicKeyGen(params, sk, t)
	C := Enc(params, mu, pk)

	mu_ := Dec(params, v, C)

	if mu != mu_ {
		panic("fuck")
	}
}
func BenchmarkMain2(b *testing.B) {
	var q float64 = 134217728
	n := 1024
	m := 1024

	var mu gp.T = gp.T(rand.Intn(2))

	params := Setup(q, n, m)
	sk, t, v := SecretKeyGen(params)
	pk := PublicKeyGen(params, sk, t)
	C := Enc(params, mu, pk)

	mu_ := Dec(params, v, C)

	if mu != mu_ {
		panic("fuck")
	}
}

func BenchmarkMain1(b *testing.B) {
	var q float64 = 134217728
	n := 1024
	m := 1024

	var mu gp.T = gp.T(rand.Intn(2))

	params := Setup(q, n, m)
	sk, t, v := SecretKeyGen(params)
	pk := PublicKeyGen(params, sk, t)
	C := Enc(params, mu, pk)

	mu_ := Dec(params, v, C)

	if mu != mu_ {
		panic("fuck")
	}
}

func RandInt(lb, ub int) int {
	return rand.Intn(ub-lb) + lb
}

func TestAddEntryWise(t *testing.T) {
	var a gp.Arr = gp.Arr{1, 2, 3}
	var old gp.Arr = gp.Arr{1, 2, 3}
	var b gp.T = 2

	gp.AddEntrywise(a, b)

	for i, e := range a {
		if e != old[i]+b {
			t.Fatal("fuuuuck")
		}
	}
}

func TestIdentMatrix(t *testing.T) {
	n := 10
	ident := gp.Ident(n)

	if len(ident) != n {
		t.Fatal("ident matrix is incorrect row size")
	}

	for i, row := range ident {
		if len(row) != n {
			t.Fatal("ident matrix is incorrect column size")
		}
		for j, entry := range row {
			if j == i && entry != 1 {
				t.Fatal("diagonal is not exclusivly 1's")
			}
			if j != i && entry != 0 {
				t.Fatal("non-diagonal is not exclusivly 0's")
			}
		}
	}
}

func TestBitDecomp1(t *testing.T) {
	M := gp.Arr{4, 2, 1}
	Expected :=
		gp.Arr{0, 0, 1, 0, 1, 0, 1, 0, 0}

	BD := BitDecomp(M, len(M))

	for i, entry := range BD {
		if entry != Expected[i] {
			t.Fatal(fmt.Sprintf("expected %d, got %d, at %d",
				Expected[i], entry, i))
		}
	}
}

func TestBitDecomp2(t *testing.T) {
	M := gp.Arr{1, 1, 1}
	Expected :=
		gp.Arr{1, 0, 0, 1, 0, 0, 1, 0, 0}

	BD := BitDecomp(M, len(M))

	for i, entry := range BD {
		if entry != Expected[i] {
			t.Fatal(fmt.Sprintf("expected %d, got %d, at %d",
				Expected[i], entry, i))
		}
	}
}

func TestBitDecomp3(t *testing.T) {

	M := gp.Arr{8, 8, 8}
	Expected := gp.Arr{0, 0, 0, 0, 0, 0, 0, 0, 0}

	BD := BitDecomp(M, len(M))

	for i, entry := range BD {
		if entry != Expected[i] {
			t.Fatal(fmt.Sprintf("expected %d, got %d, at %d",
				Expected[i], entry, i))
		}
	}
}

func TestBitDecompMatrix1(t *testing.T) {
	M := gp.Matrix{
		{4, 2, 1},
		{1, 1, 1},
		{8, 8, 8},
		{1, 2, 3},
	}

	Expected := gp.Matrix{
		{0, 0, 1, 0, 1, 0, 1, 0, 0},
		{1, 0, 0, 1, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 1, 1, 0},
	}

	BD := BitDecompMatrix(M, 3)

	for i, row := range BD {
		for j, entry := range row {
			if entry != Expected[i][j] {
				t.Fatal(fmt.Sprintf("expected %d, got %d, at %d,%d",
					Expected[i][j], entry, i, j))
			}
		}
	}

}
func TestInvBitDecomp0(t *testing.T) {
	q := 17
	k := 1
	M := gp.Matrix{
		{1, 1, 1, 1, 1},
	}

	_, _, ell := getQLogAndEll(float64(q))

	println(ell)

	BD := InvBitDecompMatrix(M, k, ell, gp.T(q))

	gp.Out(BD)

}

func TestInvBitDecomp1(t *testing.T) {
	M :=
		gp.Arr{0, 0, 1, 0, 1, 0, 1, 0, 0}
	Expected := gp.Arr{4, 2, 1}

	BD := InvBitDecomp(M, 3, 3, 8)

	for i, entry := range BD {
		if entry != Expected[i] {
			t.Fatal(fmt.Sprintf("expected %d, got %d, at %d",
				Expected[i], entry, i))
		}
	}
}

func TestInvBitDecomp2(t *testing.T) {
	M :=
		gp.Arr{4, 2, 1}
	Expected := gp.Arr{4, 2, 1}

	BD := InvBitDecomp(M, 3, 1, 8)

	for i, entry := range BD {
		if entry != Expected[i] {
			t.Fatal(fmt.Sprintf("expected %d, got %d, at %d",
				Expected[i], entry, i))
		}
	}
}

func TestInvBitDecompMatrix1(t *testing.T) {
	M := gp.Matrix{
		gp.Arr{0, 0, 1, 0, 1, 0, 1, 0, 0},
		gp.Arr{0, 0, 1, 0, 1, 0, 1, 0, 0},
	}
	Expected := gp.Matrix{
		gp.Arr{4, 2, 1},
		gp.Arr{4, 2, 1},
	}

	BD := InvBitDecompMatrix(M, 3, 3, 8)

	for i, row := range BD {
		for j, entry := range row {
			if entry != Expected[i][j] {
				t.Fatal(fmt.Sprintf("expected %d, got %d, at %d",
					Expected[i], entry, i))
			}
		}
	}
}

func TestInvBitDecompMatrix2(t *testing.T) {
	M := gp.Matrix{
		gp.Arr{4, 2, 1},
		gp.Arr{4, 2, 1},
	}
	Expected := gp.Matrix{
		gp.Arr{4, 2, 1},
		gp.Arr{4, 2, 1},
	}

	BD := InvBitDecompMatrix(M, 3, 1, 8)

	for i, row := range BD {
		for j, entry := range row {
			if entry != Expected[i][j] {
				t.Fatal(fmt.Sprintf("expected %d, got %d, at %d",
					Expected[i], entry, i))
			}
		}
	}

}

func TestDotMatrix(t *testing.T) {

	M1 := gp.Matrix{
		{2},
	}
	M2 := gp.Matrix{
		{2, 2},
	}

	result := gp.DotMatrix(M1, M2)

	gp.Out(result)
}

func TestAssertion1(t *testing.T) {
	n, q := RandInt(5, 15), RandInt(100, 1000)
	qf := float64(q)
	k := n + 1
	_, _, ell := getQLogAndEll(qf)

	a := gp.SampleMatrix(0, q, 1, k)
	b := gp.SampleMatrix(0, q, k, 1)

	ab := gp.DotMatrixMod(a, b, gp.T(q))

	lhs := BitDecompMatrix(a, ell)
	po2 := GeneratePowersOf2(ell)
	rhs := gp.TransposeColumnMatrix(ApplyPowersOf2(q, gp.TransposeArr(b), po2))

	assertLhs := gp.DotMatrixMod(lhs, rhs, gp.T(q))
	if !gp.Equal(assertLhs, ab) {
		t.Fatal("Assertion 1 does not hold")
	}
	fmt.Println("Assertion 1 holds")
}

func TestAssertion2(t *testing.T) {
	n, q := RandInt(5, 15), RandInt(100, 1000)
	qf := float64(q)
	k := n + 1
	_, _, ell := getQLogAndEll(qf)
	po2 := GeneratePowersOf2(ell)

	a := gp.SampleMatrix(0, q, 1, k)
	b := gp.SampleMatrix(0, q, k, 1)

	lhs := BitDecompMatrix(a, ell)
	rhs := gp.TransposeColumnMatrix(ApplyPowersOf2(q, gp.TransposeArr(b), po2))
	assertLhs := gp.DotMatrixMod(lhs, rhs, gp.T(q))

	// Assert rhs
	lhs1 := InvBitDecompMatrix(lhs, k, ell, gp.T(q))

	assertRhs := gp.DotMatrixMod(lhs1, b, gp.T(q)) // Change a

	if !gp.Equal(assertLhs, assertRhs) {
		t.Fatal("Assertion 2 does not hold")
	}
	fmt.Println("Assertion 2 holds")
}

func TestAssertion3(t *testing.T) {
	n, q := RandInt(5, 15), RandInt(100, 1000)
	qf := float64(q)
	k := n + 1
	_, _, ell := getQLogAndEll(qf)
	po2 := GeneratePowersOf2(ell)

	a := gp.SampleMatrix(0, q, 1, k)
	b := gp.SampleMatrix(0, q, k, 1)

	// Assert lhs
	lhs := BitDecompMatrix(a, ell)
	rhs := gp.TransposeColumnMatrix(ApplyPowersOf2(q, gp.TransposeArr(b), po2))
	assertLhs := gp.DotMatrixMod(lhs, rhs, gp.T(q))

	// Assert rhs

	flat := FlattenMatrix(lhs, k, ell, gp.T(q))
	rhs1 := gp.TransposeColumnMatrix(ApplyPowersOf2(q, gp.TransposeArr(b), po2))

	assertRhs := gp.DotMatrixMod(flat, rhs1, gp.T(q))

	if !gp.Equal(assertLhs, assertRhs) {
		t.Fatal("Assertion 3 does not hold")
	}
	fmt.Println("Assertion 3 holds")
}

func TestAssertion4(t *testing.T) {
	n, q := RandInt(5, 15), RandInt(100, 1000)
	qf := float64(q)
	_, qlog2_c, ell := getQLogAndEll(qf)

	m := 2 * n * qlog2_c

	A := gp.SampleMatrix(0, q, m, n)

	// Assert rhs

	first := BitDecompMatrix(A, ell)
	assertRhs := InvBitDecompMatrix(first, n, ell, gp.T(q))
	assertRhs = gp.ModMatrix(assertRhs, gp.T(q))

	if !gp.Equal(A, assertRhs) {
		t.Fatal("Assertion 4 does not hold")
	}
	fmt.Println("Assertion 4 holds")
}
