package main

import (
	"errors"
	"io"
	"math/bits"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/davecgh/go-spew/spew"
)

type Puzzle struct{}

type Input struct {
	Groups []*Group `(@@ Break?)+`
}

type Group struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Forms []*Form `(@@ EOL?)+`

	Answers uint32
}

func (g *Group) Tally() {
	for i, f := range g.Forms {
		f.Tally()

		if i == 0 {
			g.Answers = f.Answers
		} else {
			g.Answers &= f.Answers
		}
	}
}

type Form struct {
	Code string `@Code`

	Answers uint32
}

func (f *Form) Tally() {
	for _, c := range f.Code {
		i := int(c)
		i = i - 97
		f.Answers = f.Answers | 1<<i
	}
}

var parser = participle.MustBuild(
	&Input{},
	participle.Lexer(stateful.Must(stateful.Rules{
		"Root": {
			{"Code", `[a-z]{1,26}`, nil},
			{"Break", `[\n][\n]`, nil},
			{"EOL", `[\n]`, nil},
		},
	})),
)

type Solution struct {
	Yeses int
}

var ErrNoSolution = errors.New("no solution")

func NewInput(data io.Reader) (input *Input, err error) {
	input = &Input{}

	err = parser.Parse("input", data, input)
	if err != nil {
		return nil, err
	}

	for _, g := range input.Groups {
		g.Tally()
	}

	//spew.Dump(input)

	return input, nil
}

func (p *Puzzle) Solve(i *Input) (s *Solution, err error) {
	s = &Solution{}

	for _, g := range i.Groups {
		s.Yeses += bits.OnesCount32(g.Answers)
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
