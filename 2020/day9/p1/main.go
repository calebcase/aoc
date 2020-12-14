package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

type Puzzle struct {
	Preamble int
}

type Solution struct {
	Invalid int
}

var ErrNoSolution = errors.New("no solution")

type Lookback struct {
	Max   int
	Queue []int

	Idx    int
	Filled bool
}

func NewLookback(max int) *Lookback {
	return &Lookback{
		Max:   max,
		Queue: make([]int, max, max),
	}
}

func (l *Lookback) Push(i int) {
	l.Idx++

	if !l.Filled && l.Idx == l.Max {
		l.Filled = true
	}

	l.Idx = l.Idx % l.Max

	l.Queue[l.Idx] = i
}

func (l *Lookback) Contains(i int) bool {
	for _, v := range l.Queue {
		if v == i {
			return true
		}
	}

	return false
}

func (l *Lookback) Summable(i int) bool {
	/*
		sorted := make([]int, len(l.Queue))
		copy(sorted, l.Queue)
		sort.Ints(sorted)

		if i < sorted[0] {
			return false
		}

		if i > sorted[len(sorted)-1] {
			return false
		}
	*/

	for _, v := range l.Queue {
		if l.Contains(i - v) {
			return true
		}
	}

	return false
}

func (p *Puzzle) Solve(data io.Reader) (s *Solution, err error) {
	s = &Solution{}

	lb := NewLookback(p.Preamble)

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()

		i, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}

		if !lb.Filled {
			lb.Push(i)
			continue
		}

		if !lb.Summable(i) {
			s.Invalid = i
			break
		}

		lb.Push(i)
	}

	return s, nil
}

func main() {
	p := &Puzzle{
		Preamble: 25,
	}

	s, err := p.Solve(os.Stdin)
	if err != nil {
		panic(err)
	}

	spew.Dump(s)
}
