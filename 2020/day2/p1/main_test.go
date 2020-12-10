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
		Valid:   2,
		Invalid: 1,
	}, s)
}
