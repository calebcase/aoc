package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPass(t *testing.T) {
	for _, tc := range []Pass{
		{
			"BFFFBBFRRR",
			70,
			7,
			567,
		},
		{
			"FFFBBBFRRR",
			14,
			7,
			119,
		},
		{
			"BBFFBBFRLL",
			102,
			4,
			820,
		},
		{
			"FBFBBFFRLR",
			44,
			5,
			357,
		},
	} {
		t.Run(tc.Code, func(t *testing.T) {
			p, err := NewPass(tc.Code)
			require.NoError(t, err)

			require.Equal(t, tc, p)
		})
	}
}

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
		Highest: Pass{
			Code: "BBFFBBFRLL",
			Row:  102,
			Col:  4,
			SId:  820,
		},
	}, s)
}
