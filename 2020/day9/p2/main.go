package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

type Puzzle struct {
	Target int
}

type Solution struct {
	Key int
}

var ErrOverflow = errors.New("overflow")
var ErrNoSolution = errors.New("no solution")

type Queue []int

func (q Queue) SumSum(target int) (value int, err error) {
	min := 0
	max := 0
	sum := 0

	for _, i := range q {
		if min == 0 {
			min = i
		}
		if i < min {
			min = i
		}

		if i > max {
			max = i
		}

		sum += i

		if sum > target {
			return -1, ErrOverflow
		}

		if sum == target {
			//fmt.Printf("Solution! %d %d %d %d\n", q[0], i, min, max)
			return min + max, nil
		}
	}

	return -1, ErrNoSolution
}

func (q *Queue) Push(i int) {
	*q = append(*q, i)
}

func (q *Queue) Pop() {
	*q = (*q)[1:]
}

func (p *Puzzle) Solve(data io.Reader) (s *Solution, err error) {
	s = &Solution{}

	q := Queue{}

	scanner := bufio.NewScanner(data)

outer:
	for scanner.Scan() {
		line := scanner.Text()

		i, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}

		q.Push(i)

	inner:
		for {
			v, err := q.SumSum(p.Target)
			switch err {
			case ErrOverflow:
				q.Pop()
			case ErrNoSolution:
				break inner
			case nil:
				s.Key = v
				break outer
			}

		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return s, nil
}

func main() {
	p := &Puzzle{
		Target: 18272118,
	}

	s, err := p.Solve(os.Stdin)
	if err != nil {
		panic(err)
	}

	spew.Dump(s)
}
