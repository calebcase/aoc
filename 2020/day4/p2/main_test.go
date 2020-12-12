package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSolveValid(t *testing.T) {
	p := &Puzzle{}

	f, err := os.Open("testdata/valid")
	require.NoError(t, err)

	i, err := NewInput(f)
	require.NoError(t, err)

	s, err := p.Solve(i)
	require.NoError(t, err)

	require.NotNil(t, s)
	require.Equal(t, &Solution{
		Valid:   4,
		Invalid: 0,
	}, s)
}

func TestSolveInvalid(t *testing.T) {
	p := &Puzzle{}

	f, err := os.Open("testdata/invalid")
	require.NoError(t, err)

	i, err := NewInput(f)
	require.NoError(t, err)

	s, err := p.Solve(i)
	require.NoError(t, err)

	require.NotNil(t, s)
	require.Equal(t, &Solution{
		Valid:   0,
		Invalid: 4,
	}, s)
}
