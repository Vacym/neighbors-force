package game_test

import (
	"testing"

	"github.com/Vacym/neighbor-force/internal/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoard_NewRectangleBoard(t *testing.T) {
	testCases := []struct {
		name    string
		rows    int
		cols    int
		isValid bool
	}{
		{
			name:    "valid board",
			rows:    5,
			cols:    5,
			isValid: true,
		},
		{
			name:    "rows < 1",
			rows:    -5,
			cols:    5,
			isValid: false,
		},
		{
			name:    "rows < 1",
			rows:    -5,
			cols:    5,
			isValid: false,
		},
		{
			name:    "rows and cols < 1",
			rows:    -5,
			cols:    -5,
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			board, err := game.NewRectangleBoard(tc.rows, tc.cols)

			if !tc.isValid {
				assert.Error(t, err)
				assert.Nil(t, board)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, board)
			require.Equal(t, board.Rows(), tc.rows)
			require.Equal(t, board.Cols(), tc.cols)
		})
	}
}
