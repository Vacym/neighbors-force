package game

import (
	"errors"
	"math"
)

var (
	errSamePlayerCell = errors.New("target cell already belongs to the same player")
	errNotEnoughPower = errors.New("power of cell must be more then 1 for attack")
	errIsNotNeighbor  = errors.New("cell can attack only it's neighbor")
)

type Coords struct {
	Row int `json:"row"` // Row of the cell
	Col int `json:"col"` // Column of the cell
}

type Cell interface {
	Level() int
	Power() int
	Owner() Player
	Row() int
	Col() int

	attack(target Cell) error
	upgrade(points int) error
	calculatePower(*Board)

	GetNeighbors(*Board) []Cell

	toMap() map[string]interface{}
}

type cell struct {
	coords Coords
	level  int    // Level of cell
	power  int    // Power of cell
	owner  Player // Pointer of the player who owns the cell, nil if the cell is unoccupied
}

func (c cell) Level() int {
	return c.level
}

func (c cell) Power() int {
	return c.power
}

func (c cell) Owner() Player {
	return c.owner
}

func (c cell) Row() int {
	return c.coords.Row
}

func (c cell) Col() int {
	return c.coords.Col
}

func newCell(row, col int) *cell {
	return &cell{
		coords: Coords{row, col},
		power:  1,
		level:  1,
	}
}

func newCellWithParameters(row, col int, level, power int, owner Player) *cell {
	return &cell{
		coords: Coords{row, col},
		level:  level,
		power:  power,
		owner:  owner,
	}
}

// Returns neighbors of the cell
func (c cell) GetNeighbors(board *Board) []Cell {
	neighbors := make([]Cell, 0, 6)

	neighborCoords := getNeighborCoords(c.Row(), c.Col(), board.Rows(), board.Cols())

	for _, coord := range neighborCoords {
		if board.Cells[coord.Row][coord.Col] != nil {
			neighbors = append(neighbors, board.Cells[coord.Row][coord.Col])
		}
	}

	return neighbors
}

// Attacks target cell and defines new owner of target
func (c *cell) attack(targetInterface Cell) error {
	target := targetInterface.(*cell)
	if c.owner == target.owner {
		return errSamePlayerCell
	}
	if c.power <= 1 {
		return errNotEnoughPower
	}
	if !c.isNeighbor(target) {
		return errIsNotNeighbor
	}

	err := target.handleAttack(c)
	if err != nil {
		return err
	}

	c.power = 1

	return nil
}

func (c *cell) handleAttack(attacker Cell) error {
	attackPower := attacker.Power()
	c.power -= attackPower

	if c.power < 0 {
		if c.Owner() != nil {
			c.Owner().deleteCell()
		}
		attacker.Owner().addCell()

		c.power = -c.power
		c.owner = attacker.Owner()
		c.level = 1
	}

	return nil
}

func (c *cell) calculatePower(board *Board) {
	if c.owner == nil {
		return
	}

	newPower := 1

	for _, cell := range c.GetNeighbors(board) {
		if cell.Owner() == c.Owner() {
			newPower += cell.Level() - 1
		}
	}

	c.power = newPower
}

func (c *cell) upgrade(levels int) error {
	c.level += levels

	return nil
}

func getNeighborCoords(row, col, boardRow, boardCol int) []Coords {
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

func (c *cell) isNeighbor(cell2 *cell) bool {
	cell1 := c
	if cell1.Row() == cell2.Row() && math.Abs(float64(cell1.Col()-cell2.Col())) == 1 {
		return true
	} else if math.Abs(float64(cell1.Row()-cell2.Row())) == 1 {
		offset := cell1.Row() % 2
		return cell1.Col()-cell2.Col() == 0-offset || cell1.Col()-cell2.Col() == 1-offset
	} else {
		return false
	}
}

func (c *cell) toMap() map[string]interface{} {
	result := map[string]interface{}{
		"level": c.level,
		"power": c.power,
	}
	if c.owner != nil {
		result["owner_id"] = c.owner.Id()
	}
	return result
}
