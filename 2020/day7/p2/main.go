package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/davecgh/go-spew/spew"
)

type Puzzle struct{}

type Input struct {
	Bags []*Bag `(@@ EOL?)+`
}

type Command int

const (
	Continue Command = iota
	Stop
	Trim
	Return
)

func (i *Input) GetBag(color string) *Bag {
	for _, b := range i.Bags {
		if b.Color == color {
			return b
		}
	}

	return nil
}

func (i *Input) Walk(visit func(b *Bag, p Path) Command) {
	for _, b := range i.Bags {
		var cmd Command

		p := Path{&NBags{1, b.Color}}

		cmd = b.Walk(i, p, visit)
		switch cmd {
		case Stop:
			break
		case Trim:
			continue
		case Return:
			continue
		}

		cmd = visit(b, p)
		switch cmd {
		case Stop:
			break
		case Trim:
			continue
		}
	}
}

type Bag struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Color       string        `@Color Bag Contain`
	Constraints []*Constraint `(@@ Comma?)+ Period`
}

func (b *Bag) IsLeaf() bool {
	if len(b.Constraints) > 1 {
		return false
	}

	if b.Constraints[0].NBags != nil {
		return false
	}

	return true
}

func (b *Bag) Walk(i *Input, p Path, visit func(b *Bag, p Path) Command) Command {
	for _, c := range b.Constraints {
		if c.NBags != nil {
			var cmd Command

			if p.Contains(c.NBags) {
				continue
			}

			cmd = func() Command {
				p.Push(c.NBags)
				defer p.Pop()

				next := i.GetBag(c.NBags.Color)

				cmd = next.Walk(i, p, visit)
				if cmd != Continue {
					return cmd
				}

				cmd = visit(next, p)
				if cmd != Continue {
					return cmd
				}

				return Continue
			}()
			switch cmd {
			case Stop:
				return cmd
			case Trim:
				continue
			case Return:
				return cmd
			}
		}
	}

	return Continue
}

type Constraint struct {
	NBags *NBags  `(@@`
	Empty *string `| @Empty)`
}

type NBags struct {
	Count int    `@Number`
	Color string `@Color Bag`
}

func (nb *NBags) String() string {
	return fmt.Sprintf("%+v", *nb)
}

type Path []*NBags

func (p *Path) Push(nb *NBags) {
	*p = append(*p, nb)
}

func (p *Path) Pop() {
	*p = (*p)[:len(*p)-1]
}

func (p Path) Contains(target *NBags) bool {
	for _, nb := range p {
		if nb.Color == target.Color {
			return true
		}
	}

	return false
}

var parser = participle.MustBuild(
	&Input{},
	participle.Lexer(stateful.Must(stateful.Rules{
		"Root": {
			{"Empty", `no other bags`, nil},
			{"Bag", `bags?`, nil},
			{"Contain", `contain`, nil},
			{"Comma", `[,]`, nil},
			{"Period", `[.]`, nil},
			{"Color", `[a-z]+ [a-z]+`, nil},
			{"Number", `[1-9][0-9]*`, nil},
			{"EOL", `[\n]`, nil},
			{"Whitespace", `[ ]+`, nil},
		},
	})),
	participle.Elide("Whitespace"),
)

type Solution struct {
	Contains int
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
	s = &Solution{}

	sgb := i.GetBag("shiny gold")
	sgb.Walk(i, nil, func(b *Bag, p Path) Command {
		multiplier := 1
		for _, nb := range p {
			multiplier *= nb.Count
		}

		//fmt.Printf("%+v = %d\n", p, multiplier)

		s.Contains += multiplier

		return Continue
	})

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
