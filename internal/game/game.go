package game

import "errors"

var (
	errTooSmallPlayers      = errors.New("game cannot be played with less than 2 players")
	errTooManyPlayers       = errors.New("game cannot be played with more than 4 players")
	errNegativePlayers      = errors.New("players cannot be less than 0")
	errNotPlayerTurn        = errors.New("not player's turn to move")
	errNilPointer           = errors.New("nil pointer error")
	errInvalidAttackingCell = errors.New("attacking cell is not owned by attacking player")
)

type Game struct {
	Board   *Board   // Instance of the board
	Players []player // List of players
	turn    int      // ID of the player whose turn it is
}

// Creates a new Game with the random Board with given count of rows and cols and number of players
func NewGame(rows, cols int, numPlayers int) (*Game, error) {
	board, err := NewRandomBoard(rows, cols)
	if err != nil {
		return nil, err
	}

	players, err := NewPlayersSlice(numPlayers)
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

	game.countPlayersCell()

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

// A method that automatically places players as far apart as possible depending on the board
// Necessary and can be called only if there are no players on the field
func (g *Game) placePlayers() {
	// PLUG
	// TODO: make normal placing

	// It's temporary!!!
	for idx, player := range g.Players {
		switch idx {
		case 0:
			g.Board.Cells[0][0] = newCellWithParameters(0, 0, 1, 2, player)
		case 1:
			g.Board.Cells[0][g.Board.cols-1] = newCellWithParameters(0, g.Board.cols-1, 1, 2, player)
		case 2:
			g.Board.Cells[g.Board.rows-1][g.Board.cols-1] = newCellWithParameters(g.Board.rows-1, g.Board.cols-1, 1, 2, player)
		case 3:
			g.Board.Cells[g.Board.rows-1][0] = newCellWithParameters(g.Board.rows-1, 0, 1, 2, player)
		}
	}
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

func (g *Game) Upgrade(player player) error {
	// PLUG
	return nil
}

// Method for ending a player's turn
func (g *Game) EndTurn(player player) error {
	player.endUpgrade()

	// Implementation

	g.NextTurn()

	return nil
}

func (g *Game) NextTurn() {
	g.turn = (g.turn + 1) % len(g.Players)
}
