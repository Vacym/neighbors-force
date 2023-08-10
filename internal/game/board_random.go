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
	maxNeighborCount                     int
	swapProbability                      float64
	cells                                [][]cell
	neighbors                            [][]int
	r                                    *rand.Rand
}

// NewRandomMapGenerator creates a new instance of randomMapGenerator.
// This function allows customization of the generator's behavior using functional options.
// If values for minRight, minBottom, minLeft, and minTop are not provided,
// the default value of 1 will be used for all of them.
func NewRandomMapGenerator(boardRow, boardCol int, cells [][]cell, r *rand.Rand, options ...func(*randomMapGenerator)) *randomMapGenerator {
	gen := &randomMapGenerator{
		minRight:         1,
		minLeft:          1,
		minTop:           1,
		minBottom:        1,
		boardRow:         boardRow,
		boardCol:         boardCol,
		maxNeighborCount: 2,
		swapProbability:  0.3,
		cells:            cells,
		neighbors:        make([][]int, boardRow),
		r:                r,
	}

	for _, option := range options {
		option(gen)
	}

	for i := range gen.neighbors {
		gen.neighbors[i] = make([]int, gen.boardCol)
	}

	return gen
}

// WithMinRight sets the minimum right value.
func WithMinRight(minRight int) func(*randomMapGenerator) {
	return func(r *randomMapGenerator) {
		r.minRight = minRight
	}
}

// WithMinBottom sets the minimum bottom value.
func WithMinBottom(minBottom int) func(*randomMapGenerator) {
	return func(r *randomMapGenerator) {
		r.minBottom = minBottom
	}
}

// WithMinLeft sets the minimum left value.
func WithMinLeft(minLeft int) func(*randomMapGenerator) {
	return func(r *randomMapGenerator) {
		r.minLeft = minLeft
	}
}

// WithMinTop sets the minimum top value.
func WithMinTop(minTop int) func(*randomMapGenerator) {
	return func(r *randomMapGenerator) {
		r.minTop = minTop
	}
}

// WithMaxNeighborCount sets the maximum neighbor count.
func WithMaxNeighborCount(maxNeighborCount int) func(*randomMapGenerator) {
	return func(r *randomMapGenerator) {
		r.maxNeighborCount = maxNeighborCount
	}
}

// WithSwapProbability sets the swap probability.
func WithSwapProbability(swapProbability float64) func(*randomMapGenerator) {
	return func(r *randomMapGenerator) {
		r.swapProbability = swapProbability
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
		if g.neighbors[coord.Row][coord.Col] < g.maxNeighborCount {
			g.neighbors[coord.Row][coord.Col]++
		}
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

	sort.Slice(coords, func(i, j int) bool {
		return g.neighbors[coords[i].Row][coords[i].Col] >
			g.neighbors[coords[j].Row][coords[j].Col]
	})

	g.r.Shuffle(len(coords), func(i, j int) {
		neighborsI := g.neighbors[coords[i].Row][coords[i].Col]
		neighborsJ := g.neighbors[coords[j].Row][coords[j].Col]

		if neighborsI == neighborsJ || g.r.Float64() <= g.swapProbability {
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
		WithMinRight(max(2, rows/5)),
		WithMinBottom(max(2, rows/5)),
		WithMinLeft(2),
		WithMinTop(2),
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
