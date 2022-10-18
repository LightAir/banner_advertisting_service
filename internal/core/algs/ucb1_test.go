package algs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type UCBCases struct {
	x      int
	n      int
	nj     int
	result float64
}

func TestStorage(t *testing.T) {
	t.Run("ucb1 tests", func(t *testing.T) {
		cases := []*UCBCases{
			{
				x:      0,
				n:      0,
				nj:     0,
				result: 1,
			},
			{
				x:      0,
				n:      100,
				nj:     0,
				result: 1,
			},
			{
				x:      0,
				n:      100,
				nj:     10,
				result: 1.6786140424415112,
			},
			{
				x:      0,
				n:      100,
				nj:     100,
				result: 1.2145966026289348,
			},
			{
				x:      1,
				n:      100,
				nj:     10,
				result: 2.6786140424415112,
			},
			{
				x:      5,
				n:      100,
				nj:     10,
				result: 6.678614042441511,
			},
		}

		for _, oneCase := range cases {
			result := Ucb1(oneCase.x, oneCase.n, oneCase.nj)

			require.Equal(t, oneCase.result, result)
		}
	})
}
