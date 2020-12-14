package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	type TC struct {
		n string
		q Queue
		t int
		k int
		e error
	}

	for _, tc := range []TC{
		{"t1", Queue{1}, 9, -1, ErrNoSolution},
		{"t2", Queue{1, 2}, 9, -1, ErrNoSolution},
		{"t3", Queue{1, 2, 3}, 9, -1, ErrNoSolution},
		{"t4", Queue{1, 2, 3, 4}, 9, -1, ErrOverflow},
		{"t5", Queue{2, 3, 4}, 9, 6, nil},
		{"t6", Queue{2, 3, 4, 5}, 9, 6, nil},
	} {
		t.Run(tc.n, func(t *testing.T) {
			k, err := tc.q.SumSum(tc.t)
			require.Equal(t, tc.e, err)
			require.Equal(t, tc.k, k)
		})
	}
}

func TestSolve(t *testing.T) {
	p := &Puzzle{
		Target: 127,
	}

	f, err := os.Open("testdata/input")
	require.NoError(t, err)

	s, err := p.Solve(f)
	require.NoError(t, err)

	require.NotNil(t, s)
	require.Equal(t, &Solution{
		Key: 62,
	}, s)
}
