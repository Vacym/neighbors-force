package game_test

import (
	"testing"

	"github.com/Vacym/neighbors-force/internal/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCasesNewGameWithBoard = []struct {
	name    string
	players int
	isValid bool
}{
	{
		name:    "valid game",
		players: 4,
		isValid: true,
	},
	{
		name:    "many players",
		players: 5,
		isValid: false,
	},
	{
		name:    "small players",
		players: 1,
		isValid: false,
	},
	{
		name:    "negative count of players",
		players: -5,
		isValid: false,
	},
}

var boardSizes = []struct {
	name       string
	rows, cols int
	isValid    bool
}{
	{
		name: "valid board size",
		rows: 5, cols: 5,
		isValid: true,
	},
	{
		name: "zero rows",
		rows: 0, cols: 5,
		isValid: false,
	},
	{
		name: "zero cols",
		rows: 5, cols: 0,
		isValid: false,
	},
	{
		name: "zero rows and cols",
		rows: 0, cols: 0,
		isValid: false,
	},
	{
		name: "negative rows",
		rows: -2, cols: 5,
		isValid: false,
	},
	{
		name: "negative cols",
		rows: 5, cols: -2,
		isValid: false,
	},
	{
		name: "negative rows and cols",
		rows: -2, cols: -2,
		isValid: false,
	},
}

var testCasesNewGame = make([]struct {
	name       string
	players    int
	rows, cols int
	isValid    bool
}, 0, len(testCasesNewGameWithBoard)*len(boardSizes))

// generate testCasesNewGame
func init() {
	for _, testCase := range testCasesNewGameWithBoard {
		for _, boardSize := range boardSizes {
			newTestCase := struct {
				name       string
				players    int
				rows, cols int
				isValid    bool
			}{
				name:    testCase.name + " & " + boardSize.name,
				players: testCase.players,
				rows:    boardSize.rows,
				cols:    boardSize.cols,
				isValid: testCase.isValid && boardSize.isValid,
			}
			testCasesNewGame = append(testCasesNewGame, newTestCase)
		}
	}
}

func TestGame_NewGame(t *testing.T) {

	for _, tc := range testCasesNewGame {
		t.Run(tc.name, func(t *testing.T) {
			game, err := game.NewGame(tc.rows, tc.cols, tc.players)

			if !tc.isValid {
				assert.Error(t, err)
				assert.Nil(t, game)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, game)
		})
	}
}

func TestGame_NewGameWithBoard(t *testing.T) {
	board := game.TestBoard()

	for _, tc := range testCasesNewGameWithBoard {
		t.Run(tc.name, func(t *testing.T) {
			players, _ := game.NewPlayersSlice(tc.players)
			game, err := game.NewGameWithBoard(board, players)

			if !tc.isValid {
				assert.Error(t, err)
				assert.Nil(t, game)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, game)
		})
	}
}

func TestGame_Attack(t *testing.T) {
	g, err := game.TestGameAttack()
	require.NoError(t, err)

	// Valid attack
	err = g.Attack(g.Players[0], g.Board.Cells[1][0], g.Board.Cells[1][1])
	assert.NoError(t, err)

	// Invalid turn
	err = g.Attack(g.Players[1], g.Board.Cells[0][1], g.Board.Cells[0][0])
	assert.Error(t, err)

	// Nil from
	err = g.Attack(g.Players[0], nil, g.Board.Cells[1][1])
	assert.Error(t, err)

	// Nil to
	err = g.Attack(g.Players[0], g.Board.Cells[1][0], nil)
	assert.Error(t, err)

	// invalid owner
	err = g.Attack(g.Players[0], g.Board.Cells[0][1], g.Board.Cells[0][0])
	assert.Error(t, err)

	// Attack finished
	g.EndAttack(g.Players[0])
	err = g.Attack(g.Players[0], g.Board.Cells[1][0], g.Board.Cells[1][1])
	assert.Error(t, err)
}

func TestGame_EndAttack(t *testing.T) {
	g, err := game.TestGameAttack()
	require.NoError(t, err)

	// Invalid end attack
	err = g.EndAttack(g.Players[1])
	assert.Error(t, err)

	// Valid end attack
	err = g.EndAttack(g.Players[0])
	assert.NoError(t, err)

	g.EndTurn(g.Players[0])

	// Invalid end attack
	err = g.EndAttack(g.Players[0])
	assert.Error(t, err)

	// Valid end attack
	err = g.EndAttack(g.Players[1])
	assert.NoError(t, err)
}
func TestGame_ToMap(t *testing.T) {
	g, err := game.TestGameAttack()
	require.NoError(t, err)

	// Invalid end attack
	gameMap := g.ToMap()
	assert.NotEmpty(t, gameMap)
}
