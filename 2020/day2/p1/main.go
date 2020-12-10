package main

import (
	"bufio"
	"errors"
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

var entryRegex = regroup.MustCompile(`(?P<min>[1-9][0-9]*)-(?P<max>[1-9][0-9]*) (?P<letter>.):(?P<password>.*)`)

type Entry struct {
	Min      int    `regroup:"min"`
	Max      int    `regroup:"max"`
	Letter   string `regroup:"letter"`
	Password string `regroup:"password"`
}

func (e *Entry) Valid() bool {
	count := 0
	for _, r := range e.Password {
		if r == []rune(e.Letter)[0] {
			count += 1
		}
	}

	if count < e.Min {
		return false
	}

	if count > e.Max {
		return false
	}

	return true
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
			continue
		}

		input.Entries = append(input.Entries, entry)
	}

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
