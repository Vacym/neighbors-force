package game

import (
	"errors"
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

type Board struct {
	rows  int      // Rows of the hexagonal board
	cols  int      // Columns of the hexagonal board
	Cells [][]cell // 2D array of cells representing the board
}

// NewBoard creates a new Board with the given number of rows and columns
func NewBoard(rows, cols int) (*Board, error) {
	if rows < 2 || cols < 2 {
		return nil, errIncorrectBoardSize
	}

	cells := make([][]cell, rows)

	for i := range cells {
		cells[i] = make([]cell, cols-i%2)

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

// Returns num of rows of the hexagonal board
func (b *Board) Rows() int {
	return b.rows
}

// Returns num of columns of the hexagonal board
func (b *Board) Cols() int {
	return b.cols
}

func (b *Board) calculatePower(player player) {
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.Owner() == player {
				cell.calculatePower(b)
			}
		}
	}
}

func (b *Board) IsInsideBoard(coords Coords) bool {
	if coords.Row < 0 || coords.Col < 0 || coords.Row >= b.rows || coords.Col >= b.cols {
		return false
	}
	return true
}

func (b *Board) GetCell(coords Coords) (cell, error) {
	if !b.IsInsideBoard(coords) {
		return nil, errIndexOutOfRange
	}

	return b.Cells[coords.Row][coords.Col], nil
}

func (b *Board) toMap() map[string]interface{} {
	return map[string]interface{}{
		"rows":  b.rows,
		"cols":  b.cols,
		"cells": toCellInterfaceSlice(b.Cells),
	}
}

func toCellInterfaceSlice(cells [][]cell) [][]interface{} {
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
