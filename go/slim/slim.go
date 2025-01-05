package slim

import (
	gp "crycomproj/gopy"
	"crycomproj/util"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type Result struct {
	Success bool
	Time    time.Duration
	Message string
	Memory  uint64 // memory usage in bytes
}

type Params struct {
	n, q, ex, m, N int
	qt             gp.T
	ell            int
	PO2            gp.Arr
}

type Slim struct {
	Params
}

func (b *Slim) Setup(q float64, n, m, ex int) {
	_, _, ell := util.GetQLogAndEll(q)

	b.q = int(q)
	b.qt = gp.T(q)
	b.n = n
	b.m = m
	b.ex = ex

	b.N = (b.n + 1) * ell
	b.PO2 = util.GeneratePowersOf2(ell, b.qt)

	b.ell = ell

}

func New(q float64, n, m, ex int) *Slim {
	b := &Slim{}
	b.Setup(q, n, m, ex)
	return b
}

func (b *Slim) Run(iter int) chan Result {
	results := make(chan Result, iter)

	for i := 0; i < iter; i++ {
		var memStats runtime.MemStats
		result := Result{}
		start := time.Now()

		var mu gp.T = gp.T(rand.Intn(60))

		sk, t, v := SecretKeyGen(b.Params)
		pk := PublicKeyGen(b.Params, sk, t)

		var constant gp.T = 1

		C := Enc(b.Params, mu, pk)

		//C = MultiplyConst(C, constant, b.n+1, b.ell, b.qt)

		mu_ := MPDec(b.Params, v, C)

		result.Message = fmt.Sprintf("%d vs %d", mu, mu_)
		if mu*constant == mu_ {
			result.Success = true
		}

		end := time.Now()
		runtime.ReadMemStats(&memStats)
		result.Memory = memStats.Alloc
		result.Time = end.Sub(start)
		results <- result
	}
	close(results)

	return results
}

type Circuit func(gp.MatrixSlim, gp.MatrixSlim) gp.MatrixSlim
type Validator func(gp.T, gp.T) gp.T

func (s *Slim) RunCircuit(circuit Circuit, validate Validator) Result {
	var mu gp.T = gp.T(rand.Intn(60))
	var mu2 gp.T = gp.T(rand.Intn(60))

	return s.RunCircuitWithInput(circuit, validate, mu, mu2)
}

func (s *Slim) RunCircuitWithInput(circuit Circuit, validate Validator, mu, mu2 gp.T) Result {
	result := Result{}
	var memStats runtime.MemStats

	start := time.Now()

	sk, t, v := SecretKeyGen(s.Params)
	pk := PublicKeyGen(s.Params, sk, t)

	C := Enc(s.Params, mu, pk)
	C2 := Enc(s.Params, mu2, pk)

	Res := circuit(C, C2)

	mu_ := Dec(s.Params, v, Res)

	result.Message = fmt.Sprintf("(%d,%d) vs %d", mu, mu2, mu_)
	if validate(mu, mu2) == mu_ {
		result.Success = true
	}
	end := time.Now()
	runtime.ReadMemStats(&memStats)
	result.Memory = memStats.Alloc
	result.Time = end.Sub(start)

	return result
}

func (b *Slim) RunMetric() Result {
	var memStats runtime.MemStats
	result := Result{}
	start := time.Now()

	var mu gp.T = gp.T(rand.Intn(2))

	sk, t, v := SecretKeyGen(b.Params)
	pk := PublicKeyGen(b.Params, sk, t)

	C := Enc(b.Params, mu, pk)

	mu_ := Dec(b.Params, v, C)

	result.Message = fmt.Sprintf("%d vs %d", mu, mu_)
	if mu == mu_ {
		result.Success = true
	}

	end := time.Now()
	runtime.ReadMemStats(&memStats)
	result.Memory = memStats.Alloc
	result.Time = end.Sub(start)
  runtime.GC()
	return result
}
