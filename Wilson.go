package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	North = iota
	East
	South
	West
)

var dRow = [4]int{-1, 0, 1, 0}
var dCol = [4]int{0, 1, 0, -1}

type Maze struct {
	rows, col int
	inMaze    []bool
	path      [][4]bool
	rng       *rand.Rand
}

func newMaze(rows int, col int) *Maze {
	n := rows * col
	return &Maze{
		rows:   rows,
		col:    col,
		inMaze: make([]bool, n),
		path:   make([][4]bool, n),
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *Maze) checkBounds(cell int, d int) (int, bool) {
	r := cell/m.col + dRow[d]
	c := cell%m.col + dCol[d]
	if r < 0 || r >= m.rows || c < 0 || c >= m.col {
		return 0, false
	}
	return r*m.col + c, true

}

func (m *Maze) checkLegalDirection(cell int) int {
	legal := make([]int, 0, 4)
	for d := 0; d < 4; d++ {
		_, ok := m.checkBounds(cell, d)
		if ok {
			legal = append(legal, d)
		}

	}

	return legal[m.rng.Intn(len(legal))]

}

func (m *Maze) Generate() {
	n := m.rows * m.col
	next := make([]int, n)

	m.inMaze[m.rng.Intn(n)] = true

	for _, start := range m.rng.Perm(n) {
		if m.inMaze[start] {
			continue
		}

		cur := start
		for !m.inMaze[cur] {
			d := m.checkLegalDirection(cur)
			next[cur] = d
			cur, _ = m.checkBounds(cur, d)
		}

		cur = start
		for !m.inMaze[cur] {
			d := next[cur]
			nb, _ := m.checkBounds(cur, d)
			m.path[cur][d] = true
			m.path[nb][(d+2)%4] = true
			m.inMaze[cur] = true
			cur = nb
		}
	}

}

func readVal(in *bufio.Reader, prompt string) int {

	fmt.Print(prompt)
	for {
		var v int
		_, err := fmt.Fscan(in, &v)
		if err == nil && v > 0 {
			return v
		}

		in.ReadString('\n')
		fmt.Println("ENTER +VE NUMBER")

	}

}

func (m *Maze) String() string {
	var sb strings.Builder
	for r := 0; r < m.rows; r++ {
		// Line of horizontal walls above row r.
		sb.WriteString("+")
		for c := 0; c < m.col; c++ {
			if m.path[r*m.col+c][North] {
				sb.WriteString("   +")
			} else {
				sb.WriteString("---+")
			}
		}
		sb.WriteString("\n|")

		for c := 0; c < m.col; c++ {
			if m.path[r*m.col+c][East] {
				sb.WriteString("    ")
			} else {
				sb.WriteString("   |")
			}
		}
		sb.WriteString("\n")
	}
	sb.WriteString("+")
	sb.WriteString(strings.Repeat("---+", m.col))
	sb.WriteString("\n")
	return sb.String()
}

func main() {

	fmt.Println("Wilson's Algorithm")

	input := bufio.NewReader(os.Stdin)
	rows := readVal(input, "Enter Number of Rows: ")
	cols := readVal(input, "Enter Number of Cols: ")
	m := newMaze(rows, cols)
	m.Generate()
	fmt.Println()
	fmt.Print(m)

}
