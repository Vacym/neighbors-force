package game_test

import (
	"testing"

	"github.com/Vacym/neighbors-force/internal/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoardInputData_NewRandomBoard(t *testing.T) {
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
			name:    "cols < 1",
			rows:    5,
			cols:    -5,
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
			board, err := game.NewRandomBoard(tc.rows, tc.cols)

			if !tc.isValid {
				assert.Error(t, err)
				assert.Nil(t, board)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, board)
			require.Equal(t, board.Rows(), tc.rows)
			require.Equal(t, board.Cols(), tc.cols)
		})
	}
}

func TestBoardIntegrity_NewRandomBoard(t *testing.T) {
	repeats := 10 // How many times will each test case be checked

	testCases := []struct {
		name string
		rows int
		cols int
	}{
		{
			name: "small",
			rows: 8,
			cols: 8,
		},
		{
			name: "medium",
			rows: 28,
			cols: 28,
		},
		{
			name: "big",
			rows: 64,
			cols: 64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for attempt := 0; attempt < repeats; attempt++ {
				board, err := game.NewRandomBoard(tc.rows, tc.cols)

				require.NoError(t, err)
				require.NotNil(t, board)
				require.Equal(t, board.Rows(), tc.rows)
				require.Equal(t, board.Cols(), tc.cols)
				require.Len(t, board.Cells, tc.rows)

				for i, rowOfCells := range board.Cells {
					require.NotNil(t, rowOfCells)
					require.NotEmpty(t, rowOfCells)
					require.Len(t, rowOfCells, tc.cols)

					for j, cell := range rowOfCells {
						if cell != nil {
							row, col := cell.Coords()
							require.Equal(t, row, i)
							require.Equal(t, col, j)

							neighbors := cell.GetNeighbors(board)

							require.NotEmpty(t, neighbors)
						}
					}
				}

			}

		})
	}
}
