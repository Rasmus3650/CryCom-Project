package gopy

func ScaleSlim(arr ArrSlim, s S) ArrSlim {
	for i := range arr {
		arr[i] *= s

	}
	return arr
}

func ShapeSlim(M MatrixSlim) (int, int) {
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

func InnerProdModSlim(lhs ArrSlim, rhs Arr, mod T) T {
	var x T
	for i, e := range lhs {
		x += (T(e) * rhs[i]) % mod
	}
	return x
}
