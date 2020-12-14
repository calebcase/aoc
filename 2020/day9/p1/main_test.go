package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSolve(t *testing.T) {
	p := &Puzzle{
		Preamble: 5,
	}

	f, err := os.Open("testdata/input")
	require.NoError(t, err)

	s, err := p.Solve(f)
	require.NoError(t, err)

	require.NotNil(t, s)
	require.Equal(t, &Solution{
		Invalid: 127,
	}, s)
}
