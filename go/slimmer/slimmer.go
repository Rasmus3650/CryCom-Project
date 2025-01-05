package slimmer

import (
	gp "crycomproj/gopy"
	"crycomproj/util"
	"math/rand"
	"time"
)

type Result struct {
	Success bool
	Time    time.Duration
}

type Params struct {
	n, q, ex, m, N int
	qt             gp.T
	ell            int
	PO2            gp.Arr
}

type Slimmer struct {
	Params
}

func (b *Slimmer) Setup(q float64, n, m, ex int) {
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

func New(q float64, n, m, ex int) *Slimmer {
	b := &Slimmer{}
	b.Setup(q, n, m, ex)
	return b
}

func (b *Slimmer) Run(iter int) []Result {
	results := make([]Result, iter)

	for i := 0; i < iter; i++ {
		start := time.Now()

		var mu gp.T = gp.T(rand.Intn(2))

		sk, t, v := SecretKeyGen(b.Params)
		pk := PublicKeyGen(b.Params, sk, t)

		C := Enc(b.Params, mu, pk)

		mu_ := Dec(b.Params, v, C)

		if mu != mu_ {
			results[i].Success = false
		}
		end := time.Now()
		results[i].Time = end.Sub(start)
	}

	return results
}
