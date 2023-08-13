package game

import (
	"errors"
	"math"
)

var (
	errIncorrectBoardSize = errors.New("board size cannot be even or less than 3")
	errIndexOutOfRange    = errors.New("cell index out of range")
)

/*
In every board even row shorter of the odd

for example NewBoard(5, 6)

⬢ ⬢ ⬢ ⬢ ⬢ ⬢
 ⬢ ⬢ ⬢ ⬢ ⬢
⬢ ⬢ ⬢ ⬢ ⬢ ⬢
 ⬢ ⬢ ⬢ ⬢ ⬢
⬢ ⬢ ⬢ ⬢ ⬢ ⬢
*/

// Board represents a hexagonal game board.
type Board struct {
	rows  int      // Rows of the hexagonal board
	cols  int      // Columns of the hexagonal board
	Cells [][]Cell // 2D array of cells representing the board
}

// NewBoard creates a new hexagonal game board with the specified number of rows and columns.
func NewBoard(rows, cols int) (*Board, error) {
	if rows < 2 || cols < 2 {
		return nil, errIncorrectBoardSize
	}

	cells := make([][]Cell, rows)

	for i := range cells {
		cells[i] = make([]Cell, cols-i%2)

		for j := range cells[i] {
			cells[i][j] = newCell(i, j)
		}
	}

	board := &Board{
		rows:  rows,
		cols:  cols,
		Cells: cells,
	}

	return board, nil
}

// Rows returns the number of rows on the hexagonal board.
func (b *Board) Rows() int {
	return b.rows
}

// Cols returns the number of columns on the hexagonal board.
func (b *Board) Cols() int {
	return b.cols
}

// calculatePower updates the power levels of cells owned by a player on the board.
func (b *Board) calculatePower(player Player) {
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell != nil && cell.Owner() == player {
				cell.calculatePower(b)
			}
		}
	}
}

// IsInsideBoard checks if the given coordinates are inside the boundaries of the board.
func (b *Board) IsInsideBoard(coords Coords) bool {
	if coords.Row < 0 || coords.Col < 0 || coords.Row >= b.rows || coords.Col >= b.cols {
		return false
	}
	return true
}

// HasCellAt checks if there is a cell at the given coordinates on the board.
func (b *Board) HasCellAt(coords Coords) bool {
	cell, err := b.GetCell(coords)
	if err != nil || cell == nil {
		return false
	}
	return true
}

// GetCell returns the cell at the specified coordinates on the board.
func (b *Board) GetCell(coords Coords) (Cell, error) {
	if !b.IsInsideBoard(coords) {
		return nil, errIndexOutOfRange
	}

	return b.Cells[coords.Row][coords.Col], nil
}

// toMap converts the board's information into a map for serialization.
func (b *Board) toMap() map[string]interface{} {
	return map[string]interface{}{
		"rows":  b.rows,
		"cols":  b.cols,
		"cells": toCellSlice(b.Cells),
	}
}

// toCellSlice converts a 2D array of cells into a slice of interface slices for serialization.
func toCellSlice(cells [][]Cell) [][]interface{} {
	result := make([][]interface{}, len(cells))
	for i, row := range cells {
		result[i] = make([]interface{}, len(row))
		for j, c := range row {
			if c == nil {
				result[i][j] = nil
			} else {
				result[i][j] = c.toMap()
			}
		}
	}
	return result
}

// IsNeighborCoords checks if the given coordinates are neighbors.
func IsNeighborCoords(coord1, coord2 Coords) bool {
	if coord1.Row == coord2.Row && math.Abs(float64(coord1.Col-coord2.Col)) == 1 {
		return true
	} else if math.Abs(float64(coord1.Row-coord2.Row)) == 1 {
		offset := coord1.Row % 2
		return coord1.Col-coord2.Col == 0-offset || coord1.Col-coord2.Col == 1-offset
	}
	return false
}

// GetNeighborCoords returns the coordinates of neighboring cells.
func GetNeighborCoords(coords Coords, boardRow, boardCol int) []Coords {
	row, col := coords.Row, coords.Col
	// Offset is necessary because the hexagonal cells
	// are not placed under each other.
	// The offset depends on the row number.

	offset := row % 2
	neighborsRelative := [6]Coords{
		{-1, offset - 1}, // up-left
		{-1, offset - 0}, // up-right
		{+1, offset - 1}, // down-left
		{+1, offset - 0}, // down-right
		{0, -1},          // left
		{0, +1},          // right
	}

	neighborCoords := make([]Coords, 0, 6)

	for _, relative := range neighborsRelative {
		neighborRow := row + relative.Row
		neighborCol := col + relative.Col
		if neighborRow >= 0 && neighborCol >= 0 && neighborRow < boardRow && neighborCol < boardCol-neighborRow%2 {
			neighborCoords = append(neighborCoords, Coords{neighborRow, neighborCol})
		}
	}

	return neighborCoords
}
