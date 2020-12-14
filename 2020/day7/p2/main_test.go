package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSolve(t *testing.T) {
	type TC struct {
		path     string
		contains int
	}

	for _, tc := range []TC{
		{"testdata/input1", 32},
		{"testdata/input2", 126},
	} {
		t.Run(tc.path, func(t *testing.T) {
			p := &Puzzle{}

			f, err := os.Open(tc.path)
			require.NoError(t, err)

			i, err := NewInput(f)
			require.NoError(t, err)

			s, err := p.Solve(i)
			require.NoError(t, err)

			require.NotNil(t, s)
			require.Equal(t, &Solution{
				Contains: tc.contains,
			}, s)
		})
	}
}
