package slim

import gp "crycomproj/gopy"

func MultiplyConst(lhs gp.MatrixSlim, con gp.T, k, ell int, q gp.T) gp.MatrixSlim {
	h, w := gp.ShapeSlim(lhs)
	if h != w {
		panic("MultiplyConst requires symmetric matrix input")

	}
	Ma := FlattenMatrix2(gp.Diag(len(lhs), con), k, ell, q)

	result := gp.DotMatrixSlim(lhs, Ma, q)
	Ma = nil

	return FlattenMatrix2(result, k, ell, q)

}

func Multiply(lhs, rhs gp.MatrixSlim, k, ell int, q gp.T) gp.MatrixSlim {
	h1, w1 := gp.ShapeSlim(lhs)
	h2, w2 := gp.ShapeSlim(rhs)
	if h1 != h2 && w1 != w2 {
		panic("cannot mul two ciphertexts of different sizes")
	}

	// result can exceed 8-bit value
	result := gp.DotMatrixSlim(lhs, rhs, q)

	return FlattenMatrix2(result, k, ell, q)
}

func Add(lhs, rhs gp.MatrixSlim, k, ell int, q gp.T) gp.MatrixSlim {
	result := gp.AddMatrixSlim(lhs, rhs)

	return FlattenMatrix(result, k, ell, q)
}
