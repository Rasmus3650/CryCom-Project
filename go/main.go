package main

import (
	Scheme "crycomproj/big"
	"fmt"
	"math"
)

func main() {
	var q float64 = math.Pow(2, 27)

	msg := "q,n,m,time_in_ns,total_mem_in_byte"
	iter := 10
	upToPower := 10

	for i := 1; i <= upToPower; i++ {
		var time int64
		var mem uint64
		n := 1 << i
		m := 2*n + int(math.Log(q))
		ex := 3

		scheme := Scheme.New(q, n, m, ex)

		for j := 0; j < iter; j++ {
			result := scheme.RunMetric()
			time += result.Time.Nanoseconds()
			mem += result.Memory
		}
		avg_time := float64(time) / float64(iter)
		avg_mem := float64(mem) / float64(iter)
		msg += fmt.Sprintf("\n%d,%d,%d,%.2f,%.2f", int(q), n, m, avg_time, avg_mem)
	}
	fmt.Println(msg)
}
