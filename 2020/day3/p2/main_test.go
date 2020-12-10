package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSolve(t *testing.T) {
	p := &Puzzle{}

	f, err := os.Open("testdata/input")
	require.NoError(t, err)

	i, err := NewInput(f)
	require.NoError(t, err)

	s, err := p.Solve(i)
	require.NoError(t, err)

	require.NotNil(t, s)
	require.Equal(t, &Solution{
		HitTrees: map[Slope]int{
			Slope{1, 1}: 2,
			Slope{3, 1}: 7,
			Slope{5, 1}: 3,
			Slope{7, 1}: 4,
			Slope{1, 2}: 2,
		},
		Multiple: 336,
	}, s)
}
