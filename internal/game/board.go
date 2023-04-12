package game

import "errors"

var (
	errIncorrectBoardSize = errors.New("board size cannot be less then 1")
)

type Board struct {
	rows  int      // Rows of the hexagonal board
	cols  int      // Columns of the hexagonal board
	cells [][]Cell // 2D array of cells representing the board
}

// NewBoard creates a new Board with the given number of rows and columns
func NewRectangleBoard(rows, cols int) (*Board, error) {
	if rows < 1 || cols < 1 {
		return nil, errIncorrectBoardSize
	}

	cells := make([][]Cell, rows)

	for i := range cells {
		cells[i] = make([]Cell, cols)

		for j := range cells[i] {
			cells[i][j].x = j
			cells[i][j].y = i
		}
	}

	board := &Board{
		rows:  rows,
		cols:  cols,
		cells: cells,
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
