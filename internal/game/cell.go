package game

import "errors"

var (
	errSamePlayerCell = errors.New("target cell already belongs to the same player")
	errNotEnoughPower = errors.New("power of cell must be more then 1 for attack")
)

type coords struct {
	Row int // Row of the cell
	Col int // Column of the cell
}

type Cell struct {
	coords
	level int    // Level of cell
	power int    // Power of cell
	owner player // Pointer of the player who owns the cell, nil if the cell is unoccupied
}

func (c Cell) Owner() player {
	return c.owner
}

func newCell(row, col int) *Cell {
	return &Cell{
		coords: coords{row, col},
	}
}

// Returns neighbors of the cell
func (c Cell) GetNeighbors(board *Board) []*Cell {
	var neighbors []*Cell

	neighborCoords := getNeighborCoords(c.Row, c.Col, board.Rows(), board.Cols())

	for _, coord := range neighborCoords {
		if board.Cells[coord.Row][coord.Col] != nil {
			neighbors = append(neighbors, board.Cells[coord.Row][coord.Col])
		}
	}

	return neighbors
}

// Attacks target cell and defines new owner of target
func (c *Cell) Attack(target *Cell) error {
	if c.owner == target.owner {
		return errSamePlayerCell
	}
	if c.power <= 1 {
		return errNotEnoughPower
	}

	attackPower := c.power - 1
	c.power = 1
	target.power -= attackPower

	if target.power < 0 {
		target.owner.deleteCell()
		c.owner.addCell()

		target.power = -target.power
		target.owner = c.owner
	}

	return nil
}

func getNeighborCoords(row, col, boardRow, boardCol int) []coords {
	// Offset is necessary because the hexagonal cells
	// are not placed under each other.
	// The offset depends on the row number.
	offset := row % 2
	neighborsRelative := [6]coords{
		{-1, offset - 1}, // up-left
		{-1, offset - 0}, // up-right
		{+1, offset - 1}, // down-left
		{+1, offset - 0}, // down-right
		{0, -1},          // left
		{0, +1},          // right
	}

	neighborCoords := make([]coords, 0, 6)

	for _, relative := range neighborsRelative {
		neighborRow := row + relative.Row
		neighborCol := col + relative.Col
		if neighborRow >= 0 && neighborCol >= 0 && neighborRow < boardRow && neighborCol < boardCol {
			neighborCoords = append(neighborCoords, coords{neighborRow, neighborCol})
		}
	}

	return neighborCoords
}
