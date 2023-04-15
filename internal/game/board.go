package game

import "errors"

var (
	errIncorrectBoardSize = errors.New("board size cannot be less then 1")
)

/*
In every board even row lefter of the odd

for example NewBoard(5, 6)

⬢ ⬢ ⬢ ⬢ ⬢ ⬢
 ⬢ ⬢ ⬢ ⬢ ⬢ ⬢
⬢ ⬢ ⬢ ⬢ ⬢ ⬢
 ⬢ ⬢ ⬢ ⬢ ⬢ ⬢
⬢ ⬢ ⬢ ⬢ ⬢ ⬢
*/

type Board struct {
	rows  int      // Rows of the hexagonal board
	cols  int      // Columns of the hexagonal board
	Cells [][]cell // 2D array of cells representing the board
}

// NewBoard creates a new Board with the given number of rows and columns
func NewBoard(rows, cols int) (*Board, error) {
	if rows < 1 || cols < 1 {
		return nil, errIncorrectBoardSize
	}

	cells := make([][]cell, rows)

	for i := range cells {
		cells[i] = make([]cell, cols)

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

func NewRandomBoard(rows, cols int) (*Board, error) {
	// PLUG
	// TODO: make generation random board

	return NewBoard(rows, cols)
}

// Returns num of rows of the hexagonal board
func (b *Board) Rows() int {
	return b.rows
}

// Returns num of columns of the hexagonal board
func (b *Board) Cols() int {
	return b.cols
}
