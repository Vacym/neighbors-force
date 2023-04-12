package game

type Cell struct {
	level int     // Level of cell
	power int     // Power of cell
	owner *Player // Pointer of the player who owns the cell, nil if the cell is unoccupied

	x int
	y int
}

func (c *Cell) GetNeighbors(game *Game) []*Cell {
	var neighbors []*Cell

	// Implementation

	return neighbors
}
