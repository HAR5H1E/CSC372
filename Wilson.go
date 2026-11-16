/*
 * Wilson.go
 *
 * Author:     Harshvardhan kaul; Nguyen Do
 * Course:     CS 372
 * Instructor: MCANN
 * TAs:        -
 * Due date:   November 16, 2026
 *
 * PROBLEM DESCRIPTION
 *   This program generates a random maze of a size chosen by the user and
 *   displays it as text. The maze is created with Wilson's algorithm, which
 *   builds the maze as a uniform spanning tree of a rows x columns grid: every
 *   cell is reachable from every other cell by exactly one route, and there
 *   are no loops. Because the tree is uniform, every possible maze of the
 *   chosen size is equally likely to be produced.
 *
 * TECHNIQUES USED
 *   Wilson's algorithm uses a sequence of loop-erased random walks:
 *     1. One randomly chosen cell becomes the initial part of the maze.
 *     2. A random cell that is not yet in the maze is chosen, and a random
 *        walk starts there. Each step goes to a randomly chosen neighbor.
 *     3. The walk ends when it reaches a cell that is already in the maze.
 *        If the walk crosses its own path, the loop is discarded.
 *     4. Every cell of the loop-free path is added to the maze, and a passage
 *        is opened between each pair of consecutive cells.
 *     5. Steps 2-4 repeat until every cell is in the maze.
 *   Loop erasure is done with an array (nextDirection) that remembers only
 *   the most recent exit used from each cell. When the walk revisits a cell,
 *   the stored exit is overwritten, which makes the old loop unreachable.
 *   Following the stored exits from the starting cell then gives the
 *   loop-free path. The grid is stored as one flag per cell per side
 *   (north, east, south, west) saying whether there is a passage that way.
 *   Conceptually this is an undirected graph, so each passage is recorded
 *   on both of the two cells it joins.
 *
 * OPERATIONAL REQUIREMENTS
 *   Language: Go (use "go version" to check; developed with Go [VERSION]).
 *   Packages: standard library only (bufio, fmt, math/rand, os, strings, time).
 *   Compile/run: go run Wilson.go     (or: go build Wilson.go, then ./Wilson)
 *   Input: typed at the keyboard when prompted (number of rows, then number
 *   of columns). No input files are used.
 *   Output: the maze printed to the terminal; "+---+" and "|" are walls and
 *   spaces are open cells and passages. Each column is 4 characters wide, so
 *   very wide mazes may wrap in a narrow terminal window.
 *
 * FEATURES NOT INCLUDED / KNOWN BUGS
 *   - No entrance or exit is drawn (optional in the assignment).
 *   - No graphical version is provided (optional in the assignment).
 *   - Input is read with fmt.Fscan, which stops after the number. Extra text
 *     on the same line is left unread and is used by the next prompt. For
 *     example, typing "5 10" at the rows prompt supplies the columns value
 *     too, and typing "3.5" accepts 3 and rejects the ".5" at the next prompt.
 *   - No upper limit is placed on the maze size; extremely large sizes use
 *     large amounts of memory and time.
 */
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

/*
 * Constants: North, East, South, West
 *
 * Purpose: Name the four directions. iota numbers them 0 to 3 in this order,
 *   so the opposite of any direction is (direction+2)%4 (North and South are
 *   opposites, as are East and West). They are also used as indexes into
 *   rowOffset, colOffset, and the four passage flags of each cell.
 */
const (
	North = iota // direction 0: toward the previous row
	East         // direction 1: toward the next column
	South        // direction 2: toward the next row
	West         // direction 3: toward the previous column
)

// rowOffset holds the change in row number for a step in each direction,
// indexed by the direction constants.
var rowOffset = [4]int{-1, 0, 1, 0}

// colOffset holds the change in column number for a step in each direction,
// indexed by the direction constants.
var colOffset = [4]int{0, 1, 0, -1}

/*
 * Type: Maze
 *
 * Author:       [YOUR NAME]
 * Package:      main
 * Depends on:   math/rand (random numbers), strings (building the text
 *               drawing), time (seeding the random generator), fmt
 * Inheritance:  none (Go has no class inheritance; Maze does not embed
 *               any other type)
 *
 * Purpose: Represents a rectangular maze and builds it with Wilson's
 *   algorithm. Cells are numbered row by row starting at 0, so the cell in
 *   row r and column c has number r*cols + c. A fresh Maze has every wall
 *   closed and no cell in the maze; Generate then carves the passages.
 *
 * Public constants: North, East, South, West (see above): the directions.
 * Public variables: none. All fields are private to the package.
 *
 * Constructor:
 *   newMaze(rows, cols)  creates an empty maze of the given size.
 *
 * Methods implemented:
 *   Generate()                      builds the maze with Wilson's algorithm.
 *   String()                        returns the maze as a text drawing.
 *   checkBounds(cell, direction)    finds the neighbor of a cell, if any.
 *   checkLegalDirection(cell)       picks a random direction staying on grid.
 */
type Maze struct {
	rows, cols   int        // number of rows and columns in the grid
	inMaze       []bool     // inMaze[cell] is true once the cell has joined the maze
	path         [][4]bool  // path[cell][direction] is true if a passage leaves cell that way
	randomSource *rand.Rand // generator for every random choice the maze makes
}

/*
 * Function: newMaze
 *
 * Purpose: Constructor. Creates an empty maze of the requested size with
 *   every wall closed and no cell in the maze. The random generator is
 *   seeded from the clock so every run produces a different maze.
 *
 * Pre-condition:  rows > 0 and cols > 0.
 * Post-condition: the returned maze has rows*cols cells; no cell is in the
 *                 maze and no passage is open.
 * Returns: a pointer to the new Maze.
 *
 * Parameters:
 *   rows - (in) number of rows in the grid
 *   cols - (in) number of columns in the grid
 */
func newMaze(rows int, cols int) *Maze {
	cellCount := rows * cols // total number of cells in the grid

	return &Maze{
		rows:         rows,
		cols:         cols,
		inMaze:       make([]bool, cellCount),
		path:         make([][4]bool, cellCount),
		randomSource: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

/*
 * Method: checkBounds
 *
 * Purpose: Finds the cell reached by stepping from a cell in a given
 *   direction, and reports whether that step stays inside the grid.
 *
 * Pre-condition:  0 <= cell < rows*cols and direction is North, East,
 *                 South, or West.
 * Post-condition: the maze is not changed.
 * Returns: (neighbor cell number, true) if the step stays on the grid;
 *   (0, false) if it would leave the grid. The 0 is a placeholder, so
 *   callers must check the bool before using the cell number.
 *
 * Parameters:
 *   m         - (in) the maze (method receiver)
 *   cell      - (in) number of the cell to step from
 *   direction - (in) which way to step
 */
func (m *Maze) checkBounds(cell int, direction int) (int, bool) {
	// Convert the cell number to a row and column, then apply the step.
	row := cell/m.cols + rowOffset[direction]    // row after the step
	column := cell%m.cols + colOffset[direction] // column after the step

	if row < 0 || row >= m.rows || column < 0 || column >= m.cols {
		return 0, false
	}

	// Convert the new row and column back into a cell number.
	return row*m.cols + column, true
}

/*
 * Method: checkLegalDirection
 *
 * Purpose: Chooses a direction uniformly at random from the directions
 *   that keep a walk inside the grid. A corner cell has 2 legal directions,
 *   an edge cell has 3, and an interior cell has 4.
 *
 * Pre-condition:  0 <= cell < rows*cols, and the grid has at least 2 cells
 *                 (so at least one direction is legal).
 * Post-condition: the maze is not changed.
 * Returns: a direction (North, East, South, or West) whose step is legal.
 *
 * Parameters:
 *   m    - (in) the maze (method receiver)
 *   cell - (in) number of the cell the walk is currently on
 */
func (m *Maze) checkLegalDirection(cell int) int {
	legalDirections := make([]int, 0, 4) // directions whose step stays on the grid

	for direction := 0; direction < 4; direction++ { // direction being tested
		// Only whether the step is legal matters, so the cell is ignored.
		_, isLegal := m.checkBounds(cell, direction) // true if the step stays on the grid
		if isLegal {
			legalDirections = append(legalDirections, direction)
		}
	}

	// Intn picks a random index; the value stored there is the direction.
	return legalDirections[m.randomSource.Intn(len(legalDirections))]
}

/*
 * Method: Generate
 *
 * Purpose: Builds the maze with Wilson's algorithm (see the description at
 *   the top of this file).
 *
 * Pre-condition:  the maze is empty (as created by newMaze).
 * Post-condition: every cell is in the maze, and the passages form a
 *                 spanning tree: exactly rows*cols - 1 passages, every cell
 *                 reachable from every other, no loops. Each passage is
 *                 recorded on both of the cells it joins.
 * Returns: nothing.
 *
 * Parameters:
 *   m - (in/out) the maze (method receiver); its inMaze and path fields
 *       are filled in
 */
func (m *Maze) Generate() {
	cellCount := m.rows * m.cols // total number of cells in the grid

	// nextDirection[cell] is the direction the CURRENT walk last left that
	// cell by. It is scratch space shared by all walks; only entries for
	// cells on the current walk have meaning.
	nextDirection := make([]int, cellCount)

	// Step 1: a random cell becomes the initial part of the maze.
	m.inMaze[m.randomSource.Intn(cellCount)] = true

	// Step 2: visit the cells in random order. Skipping cells that are
	// already in the maze is equivalent to repeatedly choosing a random
	// starting cell from the cells not yet in the maze.
	for _, startCell := range m.randomSource.Perm(cellCount) { // cell where a walk may begin
		if m.inMaze[startCell] {
			continue
		}

		// Step 3: random walk from startCell until it reaches the maze.
		//
		// LOOP ERASURE: each cell remembers only its most recent exit. If
		// the walk comes back to a cell, nextDirection is overwritten with
		// the new exit, so the loop that began at that cell is forgotten.
		// Following the stored exits from startCell later therefore gives a
		// loop-free path.
		currentCell := startCell // cell the walk is standing on
		for !m.inMaze[currentCell] {
			direction := m.checkLegalDirection(currentCell) // random step to take
			nextDirection[currentCell] = direction
			// The direction is already known to be legal, so the bool is ignored.
			currentCell, _ = m.checkBounds(currentCell, direction)
		}

		// Step 4: retrace the loop-free path from startCell and carve it.
		// Each passage is opened on BOTH sides: out of the current cell, and
		// back out of the neighbor in the opposite direction. Each cell is
		// then marked as part of the maze. The loop stops at the first cell
		// that was already in the maze, which is where the walk ended.
		currentCell = startCell
		for !m.inMaze[currentCell] {
			direction := nextDirection[currentCell]                  // exit the walk last used here
			neighborCell, _ := m.checkBounds(currentCell, direction) // cell that exit leads to
			m.path[currentCell][direction] = true                    // open this side of the passage
			m.path[neighborCell][(direction+2)%4] = true             // open the opposite side
			m.inMaze[currentCell] = true
			currentCell = neighborCell
		}
	}
}

/*
 * Function: readVal
 *
 * Purpose: Displays a prompt and reads a positive whole number from the
 *   keyboard. If the input is not a positive whole number, the rest of the
 *   input line is discarded, an error message is shown, and the prompt is
 *   repeated until valid input is entered.
 *
 * Pre-condition:  reader is attached to the keyboard (standard input).
 * Post-condition: a valid number has been consumed from reader.
 * Returns: the positive whole number entered.
 *
 * Parameters:
 *   reader - (in/out) buffered reader for standard input; input is consumed
 *   prompt - (in) text shown to the user before reading
 */
func readVal(reader *bufio.Reader, prompt string) int {
	for {
		fmt.Print(prompt)

		var value int                              // number typed by the user
		_, inputError := fmt.Fscan(reader, &value) // non-nil if no whole number was read
		if inputError == nil && value > 0 {
			return value
		}

		// Throw away the rest of the bad input line before asking again.
		reader.ReadString('\n')
		fmt.Println("ENTER +VE NUMBER")
	}
}

/*
 * Method: String
 *
 * Purpose: Builds the text drawing of the maze. For each row it makes two
 *   lines: a line of horizontal walls above the row, and a line of cells
 *   with vertical walls. Only the North and East flags of each cell are
 *   read; the South and West flags only duplicate the neighbors' flags.
 *   The bottom border is added after the last row.
 *
 * Pre-condition:  Generate has been called, so the passages are final.
 * Post-condition: the maze is not changed.
 * Returns: the drawing as a string with one newline after each line;
 *   it has 2*rows + 1 lines, each 4*cols + 1 characters wide.
 *
 * Parameters:
 *   m - (in) the maze (method receiver)
 */
func (m *Maze) String() string {
	var builder strings.Builder // accumulates the text of the drawing

	for row := 0; row < m.rows; row++ { // row being drawn
		// Line of horizontal walls above this row: a gap if the cell is
		// open to the north, a wall otherwise.
		builder.WriteString("+")
		for column := 0; column < m.cols; column++ { // column being drawn
			if m.path[row*m.cols+column][North] {
				builder.WriteString("   +")
			} else {
				builder.WriteString("---+")
			}
		}
		builder.WriteString("\n|") // "|" is the left border of the maze

		// Line of cells: a gap if the cell is open to the east, a wall
		// otherwise. The east flag of a last-column cell is never set, so
		// the right border is always drawn.
		for column := 0; column < m.cols; column++ { // column being drawn
			if m.path[row*m.cols+column][East] {
				builder.WriteString("    ")
			} else {
				builder.WriteString("   |")
			}
		}
		builder.WriteString("\n")
	}

	// Bottom border of the maze.
	builder.WriteString("+")
	builder.WriteString(strings.Repeat("---+", m.cols))
	builder.WriteString("\n")

	return builder.String()
}

/*
 * Function: main
 *
 * Purpose: Describes the program to the user, reads the maze dimensions
 *   from the keyboard, generates the maze, and prints it.
 *
 * Pre-condition:  none.
 * Post-condition: the maze has been printed to standard output.
 * Returns: nothing.
 *
 * Parameters: none.
 */
func main() {
	fmt.Println("Wilson's Algorithm Maze Generator")
	fmt.Println()

	inputReader := bufio.NewReader(os.Stdin)               // buffered reader for keyboard input
	rows := readVal(inputReader, "Enter Number of Rows: ") // rows requested by the user
	cols := readVal(inputReader, "Enter Number of Cols: ") // columns requested by the user

	maze := newMaze(rows, cols) // the maze to build and display
	maze.Generate()

	fmt.Println()
	fmt.Print(maze)
}
