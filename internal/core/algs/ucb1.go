package algs

import "math"

// Ucb1 x - доход от ручки, n – то, сколько раз мы дёргали за все ручки, nj – то, сколько раз мы дёрнули за ручку j.
func Ucb1(x, n, nj int) float64 {
	fnj := float64(nj)
	fn := float64(n)
	fx := float64(x)

	baseWeight := 1.0

	if nj == 0 {
		return baseWeight
	}

	return (fx + math.Sqrt(math.Log(fn)/fnj)) + baseWeight
}
