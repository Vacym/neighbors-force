package game

import (
	"math/rand"
	"time"
)

type randomMapGenerator struct {
	right, left, top, bottom int
	boardRow, boardCol       int
	currentCells             int
	cells                    [][]cell
	r                        *rand.Rand
}

// addCell adds a new cell at the given row and column.
func (g *randomMapGenerator) addCell(row, col int) {
	g.cells[row][col] = newCell(row, col)
	g.currentCells++

	if row == 0 {
		g.top++
	}
	if row == g.boardRow-1 {
		g.bottom++
	}
	if col == 0 {
		g.left++
	}
	if col == g.boardCol-1-row%2 {
		g.right++
	}
}

// timeToStop checks if the map generation should stop based on the current state.
func (g *randomMapGenerator) timeToStop() bool {
	return g.left >= 3 && g.top >= 3 && g.right >= 3 && g.bottom >= 3
}

// generateHexMap generates the hexagonal map using randomized recursive DFS.
func (g *randomMapGenerator) generateHexMap(row, col int) {
	if !g.isValid(row, col) || g.timeToStop() {
		return
	}

	g.addCell(row, col)
	neighbors := getNeighborCoords(row, col, g.boardRow, g.boardCol)
	g.r.Shuffle(len(neighbors), func(i, j int) {
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	})
	for _, neighbor := range neighbors {
		g.generateHexMap(neighbor.Row, neighbor.Col)
	}
}

// isValid checks if the current cell is within bounds and not already visited.
func (g *randomMapGenerator) isValid(row, col int) bool {
	return row >= 0 && row < g.boardRow && col >= 0 && col < g.boardCol && g.cells[row][col] == nil
}

// NewRandomBoard generates a random hexagonal game board with the given number of rows and columns.
func NewRandomBoard(rows, cols int, seed int64) (*Board, error) {
	if rows < 3 || cols < 3 || rows%2 == 0 || cols%2 == 0 {
		return nil, errIncorrectBoardSize
	}

	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	r := rand.New(rand.NewSource(seed))

	// We'll generate only a quarter of the field, then reflect it back
	halfRows := rows/2 + rows%2
	halfCols := cols/2 + cols%2

	cells := make([][]cell, halfRows, rows)

	for i := range cells {
		cells[i] = make([]cell, halfCols-i%2, cols)
	}

	startRow := r.Intn(halfRows)
	startCol := r.Intn(halfCols - startRow%2)

	// Generate the hex map
	generator := &randomMapGenerator{
		boardRow: halfRows,
		boardCol: halfCols,
		cells:    cells,
		r:        r,
	}
	generator.generateHexMap(startRow, startCol)

	// Reflect the map vertically
	for i := 0; i < halfRows; i++ {
		for j := cols/2 - 1; j >= 0; j-- {
			if cells[i][j] != nil {
				cells[i] = append(cells[i], newCell(i, len(cells[i])))
			} else {
				cells[i] = append(cells[i], nil)
			}
		}
	}

	// Reflect the map horizontally
	for i := rows/2 - 1; i >= 0; i-- {
		cells = append(cells, make([]cell, len(cells[i])))
		lastI := len(cells) - 1
		for j := range cells[lastI] {
			if cells[i][j] != nil {
				cells[lastI][j] = newCell(lastI, j)
			}
		}
	}

	board := &Board{
		rows:  rows,
		cols:  cols,
		Cells: cells,
	}

	return board, nil
}
