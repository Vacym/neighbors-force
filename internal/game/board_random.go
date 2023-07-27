package game

import (
	"math/rand"
	"time"
)

type borders struct {
	right, left, top, bottom bool
}

func NewRandomBoard(rows, cols int) (*Board, error) {
	// We'll generate only a quarter of the field, then reflect it back
	halfRows := rows/2 + rows%2
	halfCols := cols/2 + cols%2

	cells := make([][]cell, halfRows, rows)

	for i := range cells {
		cells[i] = make([]cell, halfCols, cols)
	}

	rand.Seed(time.Now().UnixNano())
	startRow := rand.Intn(halfRows)
	startCol := rand.Intn(halfCols)

	var newBorders borders

	// Generate the hex map
	generateHexMap(startRow, startCol, halfRows, halfCols, cells, &newBorders)

	// Reflect the map vertically
	for i := 0; i < halfRows; i++ {
		for j := cols/2 - 1; j >= 0; j-- {
			cells[i] = append(cells[i], cells[i][j])
		}
	}

	// Reflect the map horizontally
	for i := rows/2 - 1; i >= 0; i-- {
		cells = append(cells, cells[i])
	}

	board := &Board{
		rows:  rows,
		cols:  cols,
		Cells: cells,
	}

	return board, nil
}

func generateHexMap(row, col, boardRow, boardCol int, cells [][]cell, borders *borders) {
	// Check if the current cell is within bounds and not already visited
	if !isValid(row, col, boardRow, boardCol, cells) {
		return
	}
	if borders.left && borders.right && borders.bottom && borders.top {
		return
	}

	if row == 0 {
		borders.top = true
	}
	if row == boardRow-1 {
		borders.bottom = true
	}
	if col == 0 {
		borders.left = true
	}
	if col == boardCol-1 {
		borders.right = true
	}

	// Mark the current cell as visited or perform any other required operations
	cells[row][col] = newCell(row, col)

	// Get the neighbors of the current cell
	neighbors := getNeighborCoords(row, col, boardRow, boardCol)

	// Shuffle the order of the neighbors randomly
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(neighbors), func(i, j int) {
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	})

	// Explore the neighbors recursively
	for _, neighbor := range neighbors {
		generateHexMap(neighbor.Row, neighbor.Col, boardRow, boardCol, cells, borders)
	}
}

func isValid(row, col, boardRow, boardCol int, cells [][]cell) bool {
	// Check if the current cell is within bounds and not already visited
	return row >= 0 && row < boardRow && col >= 0 && col < boardCol && cells[row][col] == nil
}
