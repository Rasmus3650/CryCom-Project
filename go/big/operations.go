package big

import gp "crycomproj/gopy"

func MultiplyConst(lhs gp.Matrix, con gp.T, k, ell int, q gp.T) gp.Matrix {
	h, w := gp.Shape(lhs)
	if h != w {
		panic("MultiplyConst requires symmetric matrix input")

	}
	Ma := FlattenMatrix(gp.Diag(len(lhs), con), k, ell, q)
	//gp.Out(Ma)

	result := gp.DotMatrix(lhs, Ma)

	return FlattenMatrix(result, k, ell, q)

}

func Multiply(lhs, rhs gp.Matrix, k, ell int, q gp.T) gp.Matrix {
	h1, w1 := gp.Shape(lhs)
	h2, w2 := gp.Shape(rhs)
	if h1 != h2 && w1 != w2 {
		panic("cannot mul two ciphertexts of different sizes")
	}

	// result can exceed 8-bit value
	result := gp.DotMatrix(lhs, rhs)

	return FlattenMatrix(result, k, ell, q)
}

func Add(lhs, rhs gp.Matrix, k, ell int, q gp.T) gp.Matrix {
	result := gp.AddMatrix(lhs, rhs)

	return FlattenMatrix(result, k, ell, q)
}
