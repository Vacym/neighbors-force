package game

import (
	"errors"
)

var (
	errTooSmallPlayers      = errors.New("game cannot be played with less than 2 players")
	errTooManyPlayers       = errors.New("game cannot be played with more than 4 players")
	errNegativePlayers      = errors.New("players cannot be less than 0")
	errNotPlayerTurn        = errors.New("not player's turn to move")
	errNilPointer           = errors.New("nil pointer error")
	errInvalidAttackingCell = errors.New("attacking cell is not owned by attacking player")
	errInvalidUpgradingCell = errors.New("upgrading cell is not owned by player")
)

type Game struct {
	Board   *Board   // Instance of the board
	Players []player // List of players
	turn    int      // ID of the player whose turn it is
}

// Creates a new Game with the random Board with given count of rows and cols and number of players
func NewGame(rows, cols int, numPlayers int) (*Game, error) {
	players, err := NewPlayersSlice(numPlayers)
	if err != nil {
		return nil, err
	}

	board, err := NewRandomBoard(rows, cols)
	if err != nil {
		return nil, err
	}

	game, err := NewGameWithBoard(board, players)
	if err != nil {
		return nil, err
	}

	game.placePlayers()
	game.countPlayersCell()

	return game, nil
}

// Creates a new Game with the given Board and number of players
func NewGameWithBoard(board *Board, players []player) (*Game, error) {
	if len(players) < 2 {
		return nil, errTooSmallPlayers
	} else if len(players) > 4 {
		return nil, errTooManyPlayers
	}

	game := &Game{
		Board:   board,
		Players: players,
		turn:    0,
	}

	return game, nil
}

func NewPlayersSlice(numPlayers int) ([]player, error) {
	if numPlayers < 0 {
		return nil, errNegativePlayers
	}

	players := make([]player, numPlayers)

	for i := range players {
		players[i] = NewPlayer(i)
	}

	return players, nil
}

func (g *Game) Turn() int {
	return g.turn
}

// A method that automatically places players as far apart as possible depending on the board
// Necessary and can be called only if there are no players on the field
func (g *Game) placePlayers() {
	const startPower = 2
	const startLevel = 1

	for idx, player := range g.Players {
		var c cell
		switch idx {
		case 0:
			c = findNearestCell(g.Board, 0, 0)
		case 1:
			c = findNearestCell(g.Board, g.Board.rows-1, g.Board.cols-1)
		case 2:
			c = findNearestCell(g.Board, g.Board.rows-1, 0)
		case 3:
			c = findNearestCell(g.Board, 0, g.Board.cols-1)
		}
		g.Board.Cells[c.Row()][c.Col()] = newCellWithParameters(c.Row(), c.Col(), startLevel, startPower, player)
	}
}

// findNearestCell finds the nearest cell to a given row and column coordinates
// considering the constraints of the board's dimensions.
func findNearestCell(board *Board, row, col int) cell {
	// Offset for mirror columns
	var offset int
	if col > board.cols/2 {
		offset = 1
		if row%2 == 0 {
			offset = -offset
		}
	}

	for dist := 0; dist < max(board.cols, board.rows); dist++ {
		// Search left
		if col-dist >= 0 && board.Cells[row][col-dist] != nil {
			return board.Cells[row][col-dist]
		}
		// Search right
		if col+dist < board.cols-row%2 && board.Cells[row][col+dist] != nil {
			return board.Cells[row][col+dist]
		}

		var locOffset int
		if dist%2 == 1 {
			locOffset = offset
		}

		// Search up
		if row-dist >= 0 && board.Cells[row-dist][col+locOffset] != nil {
			return board.Cells[row-dist][col+locOffset]
		}
		// Search down
		if row+dist < board.rows && board.Cells[row+dist][col+locOffset] != nil {
			return board.Cells[row+dist][col+locOffset]
		}
	}
	return nil
}

func (g *Game) countPlayersCell() {
	for _, row := range g.Board.Cells {
		for _, cell := range row {
			if cell != nil && cell.Owner() != nil {
				cell.Owner().addCell()
			}
		}
	}
}

// Method for executing a attack in the game
func (g *Game) Attack(player player, from, to cell) error {
	if player.Id() != g.turn {
		return errNotPlayerTurn
	}
	if from == nil || to == nil {
		return errNilPointer
	}
	if from.Owner() != player {
		return errInvalidAttackingCell
	}
	if err := player.attack(); err != nil {
		return err
	}

	err := from.attack(to)

	return err
}

func (g *Game) EndAttack(player player) error {
	if player.Id() != g.turn {
		return errNotPlayerTurn
	}
	player.endAttack()

	return nil
}

func (g *Game) Upgrade(player player, target cell, points int) error {
	if player.Id() != g.turn {
		return errNotPlayerTurn
	}
	if target == nil {
		return errNilPointer
	}
	if target.Owner() != player {
		return errInvalidUpgradingCell
	}
	if err := player.upgrade(points); err != nil {
		return err
	}

	err := target.upgrade(points)

	return err
}

// Method for ending a player's turn
func (g *Game) EndTurn(player player) error {
	if player.Id() != g.turn {
		return errNotPlayerTurn
	}
	player.endUpgrade()

	// Implementation

	g.NextTurn()

	return nil
}

func (g *Game) NextTurn() {
	g.turn = (g.turn + 1) % len(g.Players)
}

func (g *Game) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"board":   g.Board.toMap(),
		"players": toPlayerInterfaceSlice(g.Players),
		"turn":    g.turn,
	}
}

func toPlayerInterfaceSlice(players []player) []interface{} {
	result := make([]interface{}, len(players))
	for i, p := range players {
		result[i] = p.toMap()
	}
	return result
}

func min(x, y int) int {
	// Delete after Go 1.21 (Q3 2023)
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	// Delete after Go 1.21 (Q3 2023)
	if x > y {
		return x
	}
	return y
}
