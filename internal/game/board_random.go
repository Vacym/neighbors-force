package game

import (
	"math/rand"
	"sort"
	"time"
)

type randomMapGenerator struct {
	right, left, top, bottom             int
	minRight, minLeft, minTop, minBottom int
	boardRow, boardCol                   int
	currentCells                         int
	cells                                [][]cell
	neighbors                            [][]int
	r                                    *rand.Rand
}

// NewRandomMapGenerator creates a new instance of randomMapGenerator.
// Values can be provided for minRight, minBottom, minLeft, and minTop.
// If the values are not provided, the default value of 1 will be used for all of them.
func NewRandomMapGenerator(boardRow, boardCol int, cells [][]cell, r *rand.Rand, mins ...int) *randomMapGenerator {
	minRight, minBottom, minLeft, minTop := 1, 1, 1, 1
	if len(mins) > 0 {
		minRight = mins[0]
	}
	if len(mins) > 1 {
		minBottom = mins[1]
	}
	if len(mins) > 2 {
		minLeft = mins[2]
	}
	if len(mins) > 3 {
		minTop = mins[3]
	}

	neighbors := make([][]int, boardRow)
	for i := range neighbors {
		neighbors[i] = make([]int, boardCol)
	}

	return &randomMapGenerator{
		minRight:  minRight,
		minLeft:   minLeft,
		minTop:    minTop,
		minBottom: minBottom,
		boardRow:  boardRow,
		boardCol:  boardCol,
		cells:     cells,
		neighbors: neighbors,
		r:         r,
	}
}

// addCell adds a new cell at the given row and column.
func (g *randomMapGenerator) addCell(row, col int) cell {
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

	for _, coord := range getNeighborCoords(row, col, g.boardRow, g.boardCol) {
		g.neighbors[coord.Row][coord.Col]++
	}

	return g.cells[row][col]
}

// timeToStop checks if the map generation should stop based on the current state.
func (g *randomMapGenerator) timeToStop() bool {
	return g.left >= g.minLeft &&
		g.top >= g.minTop &&
		g.right >= g.minRight &&
		g.bottom >= g.minBottom
}

// shuffleCoords shuffle slice of coords
func (g *randomMapGenerator) shuffleCoords(coords []Coords) {
	limit := 2
	swapDifferentNeighborsProbability := 0.3

	sort.Slice(coords, func(i, j int) bool {
		return min(limit, g.neighbors[coords[i].Row][coords[i].Col]) >
			min(limit, g.neighbors[coords[j].Row][coords[j].Col])
	})

	g.r.Shuffle(len(coords), func(i, j int) {
		neighborsI := min(limit, g.neighbors[coords[i].Row][coords[i].Col])
		neighborsJ := min(limit, g.neighbors[coords[j].Row][coords[j].Col])

		if neighborsI == neighborsJ || g.r.Float64() <= swapDifferentNeighborsProbability {
			coords[i], coords[j] = coords[j], coords[i]
		}
	})
}

// generateHexMap generates the hexagonal map using randomized recursive DFS.
func (g *randomMapGenerator) generateHexMap(row, col int) {
	if !g.isValid(row, col) || g.timeToStop() {
		return
	}

	g.addCell(row, col)
	neighbors := getNeighborCoords(row, col, g.boardRow, g.boardCol)
	g.shuffleCoords(neighbors)

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
	halfRows := rows/2 + 1
	halfCols := cols/2 + 1

	cells := make([][]cell, halfRows, rows)

	for i := range cells {
		cells[i] = make([]cell, halfCols-i%2, cols)
	}

	startRow := r.Intn(halfRows)
	startCol := r.Intn(halfCols - startRow%2)

	// Generate the hex map
	generator := NewRandomMapGenerator(
		halfRows, halfCols,
		cells, r,
		max(2, rows/5),
		max(2, rows/5),
		2, 2,
	)
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
