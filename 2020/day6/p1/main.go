package main

import (
	"errors"
	"io"
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

	Answers map[rune]bool
}

func (g *Group) Tally() {
	g.Answers = make(map[rune]bool)

	for _, f := range g.Forms {
		f.Tally()

		for a := range f.Answers {
			g.Answers[a] = true
		}
	}
}

type Form struct {
	Code string `@Code`

	Answers map[rune]bool
}

func (f *Form) Tally() {
	f.Answers = make(map[rune]bool)

	for _, c := range f.Code {
		f.Answers[c] = true
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

	for _, g := range i.Groups {
		g.Tally()

		s.Yeses += len(g.Answers)
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
