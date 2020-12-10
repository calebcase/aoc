package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/oriser/regroup"
)

type Puzzle struct{}

// Example Input:
//
// 1-3 a: abcde
// 1-3 b: cdefg
// 2-9 c: ccccccccc
//
type Input struct {
	Entries []*Entry
}

var entryRegex = regroup.MustCompile(`(?P<posa>[1-9][0-9]*)-(?P<posb>[1-9][0-9]*) (?P<letter>.): (?P<password>.*)`)

type Entry struct {
	PosA     int    `regroup:"posa"`
	PosB     int    `regroup:"posb"`
	Letter   string `regroup:"letter"`
	Password string `regroup:"password"`

	Reason string
}

func (e *Entry) Valid() bool {
	idxA := e.PosA - 1
	idxB := e.PosB - 1

	if idxA >= len(e.Password) {
		e.Reason = "index A is too high"

		return false
	}
	if idxB >= len(e.Password) {
		e.Reason = "index B is too high"

		return false
	}

	matchA := false
	matchB := false
	if string(e.Password[idxA]) == e.Letter {
		matchA = true
	}
	if string(e.Password[idxB]) == e.Letter {
		matchB = true
	}

	switch {
	case matchA && matchB:
		e.Reason = "A & B match"

		return false
	case matchA && !matchB:
		e.Reason = "A & !B match"

		return true
	case !matchA && matchB:
		e.Reason = "!A & B match"

		return true
	case !matchA && !matchB:
		e.Reason = "!A & !B match"

		return false
	}

	e.Reason = "invalid logic"
	return false
}

type Solution struct {
	Valid   int
	Invalid int
}

var ErrNoSolution = errors.New("no solution")

func NewInput(data io.Reader) (i *Input, err error) {
	input := &Input{}

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		entry := &Entry{}
		err := entryRegex.MatchToTarget(scanner.Text(), entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			continue
		}

		input.Entries = append(input.Entries, entry)
	}

	//spew.Dump(input)

	return input, nil
}

func (p *Puzzle) Solve(i *Input) (s *Solution, err error) {
	s = &Solution{}

	for _, e := range i.Entries {
		if e.Valid() {
			s.Valid += 1
		} else {
			s.Invalid += 1
		}
	}

	//spew.Dump(i)

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
