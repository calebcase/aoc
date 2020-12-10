package main

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type Puzzle struct{}

// Example Input:
//
// ..##.......
// #...#...#..
// .#....#..#.
// ..#.#...#.#
// .#...##..#.
// ..#.##.....
// .#.#.#....#
// .#........#
// #.##...#...
// #...##....#
// .#..#...#.#
//
type Input struct {
	Board Board
}

type Player struct {
	X int
	Y int
}

type Board struct {
	Grid   []Row
	Player *Player
}

func (b Board) Move() (ok bool) {
	if b.Player.Y+1 >= len(b.Grid) {
		return false
	}

	b.Player.X += 3
	b.Player.Y += 1

	return true
}

func (b Board) HitTree() bool {
	r := b.Grid[b.Player.Y]

	return r.Get(b.Player.X)
}

type Row []bool

func (r Row) Get(i int) bool {
	return r[i%len(r)]
}

type Solution struct {
	HitTrees int
}

var ErrNoSolution = errors.New("no solution")

func NewInput(data io.Reader) (i *Input, err error) {
	input := &Input{}
	input.Board.Player = &Player{}

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()
		row := make(Row, len(line), len(line))

		for i, b := range line {
			if b == '#' {
				row[i] = true
			}
		}

		input.Board.Grid = append(input.Board.Grid, row)
	}

	//spew.Dump(input)

	return input, nil
}

func (p *Puzzle) Solve(i *Input) (s *Solution, err error) {
	s = &Solution{}

	for i.Board.Move() {
		//spew.Dump(i.Board.Player)

		if i.Board.HitTree() {
			s.HitTrees += 1
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
