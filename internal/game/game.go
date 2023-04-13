package game

import "errors"

var (
	errTooSmallPlayers = errors.New("game cannot be played with less than 2 players")
	errTooManyPlayers  = errors.New("game cannot be played with more than 4 players")
)

type Game struct {
	board   *Board    // Instance of the board
	players []*Player // List of players
	turn    int       // ID of the player whose turn it is
}

// Creates a new Game with the random Board with given count of rows and cols and number of players
func NewGame(rows, cols int, numPlayers int) (*Game, error) {
	board, err := NewRandomBoard(rows, cols)
	if err != nil {
		return nil, err
	}

	game, err := NewGameWithBoard(board, numPlayers)
	if err != nil {
		return nil, err
	}

	game.placePlayers()

	return game, nil
}

// Creates a new Game with the given Board and number of players
func NewGameWithBoard(board *Board, numPlayers int) (*Game, error) {
	if numPlayers < 2 {
		return nil, errTooSmallPlayers
	} else if numPlayers > 4 {
		return nil, errTooManyPlayers
	}

	players := make([]*Player, numPlayers)

	for i := range players {
		players[i] = NewPlayer(i)
	}

	game := &Game{
		board:   board,
		players: players,
		turn:    0,
	}

	game.placePlayers()

	return game, nil
}

// A method that automatically places players as far apart as possible depending on the board
// Necessary and can be called only if there are no players on the field
func (g *Game) placePlayers() {
	// PLUG
	// TODO: make normal placing

	for idx, player := range g.players {
		var cell *Cell

		switch idx {
		case 0:
			cell = g.board.Cells[0][0]
		case 1:
			cell = g.board.Cells[0][g.board.cols-1]
		case 2:
			cell = g.board.Cells[g.board.rows-1][g.board.cols-1]
		case 3:
			cell = g.board.Cells[g.board.rows-1][0]
		}

		cell.level = 1
		cell.power = 2
		cell.owner = player
	}
}

// Method for executing a move in the game
func (g *Game) Move(from, to int, numUnits int) bool {
	// Implementation
	return false
}

// Method for ending a player's turn
func (g *Game) EndTurn() bool {
	// Implementation
	return false
}

func (g *Game) NextTurn() {
	g.turn = (g.turn + 1) % len(g.players)
}
