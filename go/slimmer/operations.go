package slimmer

import gp "crycomproj/gopy"

func Multiply(lhs, rhs gp.Matrix, ell int, q gp.T) gp.Matrix {
	h, w := gp.Shape(lhs)
	Cprime := make(gp.Matrix, h)

	for i := 0; i < h; i++ { // loop over k values
		arr := make(gp.Arr, w)

		for t := 0; t < w; t++ {
			var val gp.T
			for t2 := 0; t2 < w; t2++ {
				var num gp.T = lhs[i][t2]
				for j := 0; j < ell; j++ {
					if (num & 1) > 0 {
						idx := t2*ell + j
						val += rhs[idx][t]
						val %= q << 1
					} else if (num & 1) > 1 {
						panic("waaaat")
					}

					num >>= 1
				}
			}

			arr[t] += val
			arr[t] = arr[t] % (q << 1)
		}
		Cprime[i] = arr
	}

	return Cprime
}

func Add(lhs, rhs gp.Matrix, q gp.T) gp.Matrix {
	for i, row := range lhs {
		for j, ele := range row {
			lhs[i][j] = (ele + rhs[i][j]) % q
		}
	}
	return lhs
}
