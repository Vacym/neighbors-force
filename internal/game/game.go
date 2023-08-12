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

// Game represents the core structure that encapsulates the state and logic of the game.
// It holds references to the game board, a list of players, and the current turn's player ID.
type Game struct {
	Board   *Board   // Instance of the board
	Players []Player // List of players
	turn    int      // ID of the player whose turn it is
}

// NewGame creates a new Game with a random Board and specified number of players.
func NewGame(rows, cols int, numPlayers int, seed int64) (*Game, error) {
	players, err := NewPlayersSlice(numPlayers)
	if err != nil {
		return nil, err
	}

	board, err := NewRandomBoard(rows, cols, seed)
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

// NewFullGame creates a new Game with a fully filled board (for testing purposes).
func NewFullGame(rows, cols int, numPlayers int) (*Game, error) {
	players, err := NewPlayersSlice(numPlayers)
	if err != nil {
		return nil, err
	}

	board, err := NewBoard(rows, cols)
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

// NewGameWithBoard creates a new Game with a given Board and player list.
func NewGameWithBoard(board *Board, players []Player) (*Game, error) {
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

// NewPlayersSlice creates a slice of Player instances based on the given number.
func NewPlayersSlice(numPlayers int) ([]Player, error) {
	if numPlayers < 0 {
		return nil, errNegativePlayers
	}

	players := make([]Player, numPlayers)

	for i := range players {
		players[i] = newPlayer(i)
	}

	return players, nil
}

// Turn returns the ID of the current player's turn.
func (g *Game) Turn() int {
	return g.turn
}

// placePlayers places players on the board at specific locations.
func (g *Game) placePlayers() {
	const startPower = 2
	const startLevel = 1

	for idx, player := range g.Players {
		var c Cell
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

// findNearestCell finds the nearest cell based on row and column coordinates.
func findNearestCell(board *Board, row, col int) Cell {
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

// countPlayersCell updates player cell counts on the board.
func (g *Game) countPlayersCell() {
	for _, row := range g.Board.Cells {
		for _, cell := range row {
			if cell != nil && cell.Owner() != nil {
				cell.Owner().addCell()
			}
		}
	}
}

// Attack performs an attack from one cell to another.
func (g *Game) Attack(player Player, from, to Cell) error {
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

// EndAttack ends the attack phase for the current player.
func (g *Game) EndAttack(player Player) error {
	if player.Id() != g.turn {
		return errNotPlayerTurn
	}

	return player.endAttack()
}

// Upgrade upgrades a target cell's level by a specified number of levels.
func (g *Game) Upgrade(player Player, target Cell, levels int) error {
	if player.Id() != g.turn {
		return errNotPlayerTurn
	}
	if target == nil {
		return errNilPointer
	}
	if target.Owner() != player {
		return errInvalidUpgradingCell
	}
	if err := player.upgrade(target, levels); err != nil {
		return err
	}

	err := target.upgrade(levels)

	return err
}

// EndTurn ends the current player's turn, updating the board and switching turns.
func (g *Game) EndTurn(player Player) error {
	if player.Id() != g.turn {
		return errNotPlayerTurn
	}
	player.endUpgrade()

	g.Board.calculatePower(player)

	g.nextTurn()

	return nil
}

// nextTurn advances the turn to the next player.
func (g *Game) nextTurn() {
	g.turn = (g.turn + 1) % len(g.Players)
}

func (g *Game) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"board":   g.Board.toMap(),
		"players": toPlayerInterfaceSlice(g.Players),
		"turn":    g.turn,
	}
}

// ToMap converts the game state into a map for serialization.
func toPlayerInterfaceSlice(players []Player) []interface{} {
	result := make([]interface{}, len(players))
	for i, p := range players {
		result[i] = p.toMap()
	}
	return result
}

// min returns the minimum of two integers.
func min(x, y int) int {
	// Delete after Go 1.21 (Q3 2023)
	if x < y {
		return x
	}
	return y
}

// max returns the maximum of two integers.
func max(x, y int) int {
	// Delete after Go 1.21 (Q3 2023)
	if x > y {
		return x
	}
	return y
}
