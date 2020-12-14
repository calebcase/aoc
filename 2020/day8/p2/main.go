package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/davecgh/go-spew/spew"
)

type Puzzle struct{}

type Input struct {
	Instructions []*Instruction `(@@ EOL?)+`

	Idx int
	Acc int

	Perm  int
	Trial int
}

func (i *Input) Next() bool {
	for i.Perm < len(i.Instructions) {
		cmd := i.Instructions[i.Perm]

		switch cmd.Op {
		case "acc":
			i.Perm++
		case "jmp":
			switch i.Trial {
			case 0:
				i.Trial = 1
				return true
			case 1:
				cmd.Op = "nop"
				i.Trial = 2
				return true
			default:
				cmd.Op = "nop"
				i.Trial = 0
				i.Perm++
			}
		case "nop":
			switch i.Trial {
			case 0:
				i.Trial = 1
				return true
			case 1:
				cmd.Op = "jmp"
				i.Trial = 2
				return true
			default:
				cmd.Op = "jmp"
				i.Trial = 0
				i.Perm++
			}
		}
	}

	return false
}

func (i *Input) Reset() {
	i.Idx = 0
	i.Acc = 0

	for _, cmd := range i.Instructions {
		cmd.Executed = false
	}
}

func (i *Input) Execute() bool {
	i.Reset()

	for {
		if i.Idx >= len(i.Instructions) {
			return true
		}

		cmd := i.Instructions[i.Idx]

		if cmd.Executed {
			return false
		}

		fmt.Printf("%s %d (%d %d) => ", cmd.Op, cmd.Arg, i.Idx, i.Acc)

		switch cmd.Op {
		case "acc":
			i.Acc += cmd.Arg
			i.Idx++
		case "jmp":
			i.Idx += cmd.Arg
		case "nop":
			i.Idx++
		}

		fmt.Printf("%d %d\n", i.Idx, i.Acc)

		cmd.Executed = true
	}
}

type Instruction struct {
	Op  string `@Op`
	Arg int    `@Number`

	Executed bool
}

var parser = participle.MustBuild(
	&Input{},
	participle.Lexer(stateful.Must(stateful.Rules{
		"Root": {
			{"Op", `(acc|jmp|nop)`, nil},
			{"Number", `[+-]?[0-9]+`, nil},
			{"EOL", `[\n]`, nil},
			{"Whitespace", `[ ]+`, nil},
		},
	})),
	participle.Elide("Whitespace"),
)

type Solution struct {
	Acc int
}

var ErrNoSolution = errors.New("no solution")

func NewInput(data io.Reader) (i *Input, err error) {
	input := &Input{}

	err = parser.Parse("input", data, input)
	if err != nil {
		return nil, err
	}

	//spew.Dump(input)

	return input, nil
}

func (p *Puzzle) Solve(i *Input) (s *Solution, err error) {
	for i.Next() {
		//fmt.Printf("perm %d %d\n", i.Perm, i.Trial)
		if i.Execute() {
			//fmt.Printf("solution!\n")
			break
		}
	}

	return &Solution{
		Acc: i.Acc,
	}, nil
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
