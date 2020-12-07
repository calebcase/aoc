package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Puzzle struct {
	Target int64
}

type Input map[int64]bool

type Solution struct {
	A int64
	B int64
	M int64
}

var ErrNoSolution = errors.New("no solution")

func NewInput(data io.Reader) (i Input, err error) {
	input := Input{}

	scanner := bufio.NewScanner(data)

	for scanner.Scan() {
		text := scanner.Text()

		number, err := strconv.ParseInt(text, 10, 0)
		if err != nil {
			log.Print(err)

			continue
		}

		input[number] = true
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return input, nil
}

func (p *Puzzle) Solve(input Input) (s *Solution, err error) {
	for key := range input {
		if _, ok := input[p.Target-key]; ok {
			return &Solution{
				A: p.Target,
				B: p.Target - key,
				M: key * (p.Target - key),
			}, nil
		}
	}

	return nil, ErrNoSolution
}

func main() {
	p := &Puzzle{
		Target: 2020,
	}

	input, err := NewInput(os.Stdin)
	if err != nil {
		panic(err)
	}

	s, err := p.Solve(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", s)
}
