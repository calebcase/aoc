package main

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type Puzzle struct{}

// Example:
//
// BFFFBBFRRR: row 70, column 7, seat ID 567.
// FFFBBBFRRR: row 14, column 7, seat ID 119.
// BBFFBBFRLL: row 102, column 4, seat ID 820.
//
type Input struct {
	Passes []Pass
}

type Pass struct {
	Code string
	Row  int
	Col  int
	SId  int
}

type Solution struct {
	SId int
}

var ErrNoSolution = errors.New("no solution")

func NewInput(data io.Reader) (i *Input, err error) {
	input := &Input{}

	scanner := bufio.NewScanner(data)

	for scanner.Scan() {
		line := scanner.Text()

		pass, err := NewPass(line)
		if err != nil {
			return nil, err
		}

		input.Passes = append(input.Passes, pass)
	}

	//spew.Dump(input)

	return input, nil
}

func NewPass(code string) (p Pass, err error) {
	p.Code = code

	for _, b := range code {
		switch b {
		case 'F':
			p.Row = p.Row << 1
		case 'B':
			p.Row = p.Row << 1
			p.Row++
		case 'L':
			p.Col = p.Col << 1
		case 'R':
			p.Col = p.Col << 1
			p.Col++
		}
	}

	p.SId = p.Row*8 + p.Col

	return p, nil
}

func (p *Puzzle) Solve(i *Input) (s *Solution, err error) {
	s = &Solution{}

	seats := make([]bool, 1024, 1024)

	for _, p := range i.Passes {
		seats[p.SId] = true
	}

	started := false
	for sid, taken := range seats {
		if started == false && taken != false {
			started = true
		}

		if started == true && taken != true {
			s.SId = sid
			break
		}
	}

	return s, nil
}

func main() {
	p := &Puzzle{}

	input, err := NewInput(os.Stdin)
	if err != nil {
		panic(err)
	}

	s, err := p.Solve(input)
	if err != nil {
		panic(err)
	}

	spew.Dump(s)
}
