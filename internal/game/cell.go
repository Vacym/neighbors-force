package game

type coords struct {
	Row int // Row of the cell
	Col int // Column of the cell
}

type Cell struct {
	coords
	level int     // Level of cell
	power int     // Power of cell
	owner *Player // Pointer of the player who owns the cell, nil if the cell is unoccupied
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

func getNeighborCoords(row, col, boardRow, boardCol int) []coords {
	// Offset is necessary because the hexagonal cells
	// are not placed under each other.
	// The offset depends on the row number.
	offset := row % 2
	neighborsRelative := [6]coords{
		{-1, offset - 1},
		{-1, offset - 0},
		{+1, offset - 1},
		{+1, offset - 0},
		{0, -1},
		{0, +1},
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
