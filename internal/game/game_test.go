package game_test

import (
	"testing"

	"github.com/Vacym/neighbors-force/internal/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGame_NewGameWithBoard(t *testing.T) {
	board := game.TestBoard()

	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game, err := game.NewGameWithBoard(board, tc.players)

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
