package main

import (
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/davecgh/go-spew/spew"
	"github.com/oriser/regroup"
)

type Puzzle struct{}

// Example:
//
// ecl:gry pid:860033327 eyr:2020 hcl:#fffffd
// byr:1937 iyr:2017 cid:147 hgt:183cm
//
// iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884
// hcl:#cfa07d byr:1929
//
// hcl:#ae17e1 iyr:2013
// eyr:2024
// ecl:brn pid:760753108 byr:1931
// hgt:179cm
//
// hcl:#cfa07d eyr:2025 pid:166559648
// iyr:2011 ecl:brn hgt:59in
//
type Input struct {
	Passports []*Passport `(@@ Break?)+`
}

type Passport struct {
	Pos    lexer.Position
	EndPos lexer.Position
	Tokens []lexer.Token

	Fields []*Field `@@+`
}

type Field struct {
	Id    string `@Id Colon`
	Value string `@Value`
}

var parser = participle.MustBuild(
	&Input{},
	participle.Lexer(stateful.Must(stateful.Rules{
		"Root": {
			{"Id", `[a-z]{3}`, nil},
			{"Colon", `:`, stateful.Push("Field")},
			{"Break", `[\n][\n]`, nil},
			{"Whitespace", `[ \t\n]+`, nil},
		},
		"Field": {
			{"Value", `[^ \t\n]+`, stateful.Pop()},
		},
	})),
	participle.Elide("Whitespace"),
)

type Solution struct {
	Valid   int
	Invalid int
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

func (p *Passport) Valid() bool {
	// byr (Birth Year)
	// iyr (Issue Year)
	// eyr (Expiration Year)
	// hgt (Height)
	// hcl (Hair Color)
	// ecl (Eye Color)
	// pid (Passport ID)
	// cid (Country ID)
	counters := map[string]int{
		"byr": 0,
		"iyr": 0,
		"eyr": 0,
		"hgt": 0,
		"hcl": 0,
		"ecl": 0,
		"pid": 0,
		"cid": 0,
	}

	for _, field := range p.Fields {
		if !field.Valid() {
			log.Printf("line %d: invalid field: %+v\n%s\n", p.Pos.Line, field, p.Tokens)
			return false
		}

		counters[field.Id]++
	}

	for key, value := range counters {
		if value > 1 {
			log.Printf("line %d: duplicate values\n%s\n", p.Pos.Line, p.Tokens)

			return false
		}
		if value == 0 && key != "cid" {
			log.Printf("line %d: missing required values\n%s\n", p.Pos.Line, p.Tokens)

			return false
		}
	}

	return true
}

type Height struct {
	Amount int    `regroup:"Amount,required"`
	Unit   string `regroup:"Unit,required"`
}

func (h Height) Valid() bool {
	switch h.Unit {
	case "cm":
		if h.Amount < 150 || h.Amount > 193 {
			return false
		}
	case "in":
		if h.Amount < 59 || h.Amount > 76 {
			return false
		}
	}

	return true
}

func (f *Field) Valid() bool {
	switch f.Id {
	case "byr":
		i, err := strconv.Atoi(f.Value)
		if err != nil {
			return false
		}
		if i < 1920 || i > 2002 {
			return false
		}
	case "iyr":
		i, err := strconv.Atoi(f.Value)
		if err != nil {
			return false
		}
		if i < 2010 || i > 2020 {
			return false
		}
	case "eyr":
		i, err := strconv.Atoi(f.Value)
		if err != nil {
			return false
		}
		if i < 2020 || i > 2030 {
			return false
		}
	case "hgt":
		re := regroup.MustCompile(`^(?P<Amount>\d+)(?P<Unit>(cm|in))$`)

		h := &Height{}
		err := re.MatchToTarget(f.Value, h)
		if err != nil {
			spew.Dump(err)

			return false
		}

		return h.Valid()
	case "hcl":
		re := regexp.MustCompile(`^#[0-9a-f]{6}$`)

		return re.MatchString(f.Value)
	case "ecl":
		re := regexp.MustCompile(`^(amb|blu|brn|gry|grn|hzl|oth)$`)

		return re.MatchString(f.Value)
	case "pid":
		re := regexp.MustCompile(`^[0-9]{9}$`)

		return re.MatchString(f.Value)
	case "cid":
		return true
	}

	return true
}

func (p *Puzzle) Solve(i *Input) (s *Solution, err error) {
	s = &Solution{}

	for _, p := range i.Passports {
		if p.Valid() {
			s.Valid++
		} else {
			s.Invalid++
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
