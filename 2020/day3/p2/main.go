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
	Board *Board
}

type Player struct {
	X int
	Y int
}

type Board struct {
	Grid   []Row
	Player *Player
}

type Slope struct {
	Right int
	Down  int
}

func (b *Board) Reset() {
	b.Player = &Player{}
}

func (b *Board) Move(s Slope) (ok bool) {
	if b.Player.Y+s.Down >= len(b.Grid) {
		return false
	}

	b.Player.X += s.Right
	b.Player.Y += s.Down

	return true
}

func (b *Board) HitTree() bool {
	r := b.Grid[b.Player.Y]

	return r.Get(b.Player.X)
}

type Row []bool

func (r Row) Get(i int) bool {
	return r[i%len(r)]
}

type Solution struct {
	HitTrees map[Slope]int
	Multiple int
}

func NewSolution() *Solution {
	return &Solution{
		HitTrees: map[Slope]int{},
		Multiple: 1,
	}
}

var ErrNoSolution = errors.New("no solution")

func NewInput(data io.Reader) (i *Input, err error) {
	input := &Input{}
	input.Board = &Board{}
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
	s = NewSolution()

	slopes := []Slope{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}

	for _, slope := range slopes {
		i.Board.Reset()

		for i.Board.Move(slope) {
			//spew.Dump(i.Board.Player)

			if i.Board.HitTree() {
				s.HitTrees[slope]++
			}
		}
	}

	for _, hits := range s.HitTrees {
		s.Multiple *= hits
	}

	//spew.Dump(s)

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
