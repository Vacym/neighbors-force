package game_test

import (
	"testing"

	"github.com/Vacym/neighbors-force/internal/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golang-collections/collections/queue"
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
			name:    "rows < 3",
			rows:    1,
			cols:    5,
			isValid: false,
		},
		{
			name:    "cols < 3",
			rows:    5,
			cols:    1,
			isValid: false,
		},
		{
			name:    "rows and cols < 3",
			rows:    1,
			cols:    1,
			isValid: false,
		},
		{
			name:    "rows is even",
			rows:    4,
			cols:    5,
			isValid: false,
		},
		{
			name:    "cols is even",
			rows:    5,
			cols:    4,
			isValid: false,
		},
		{
			name:    "rows and cols are even",
			rows:    4,
			cols:    4,
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			board, err := game.NewRandomBoard(tc.rows, tc.cols, 0)

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
			rows: 9,
			cols: 9,
		},
		{
			name: "medium",
			rows: 31,
			cols: 31,
		},
		{
			name: "big",
			rows: 63,
			cols: 63,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for attempt := 0; attempt < repeats; attempt++ {
				board, err := game.NewRandomBoard(tc.rows, tc.cols, 0)

				require.NoError(t, err)
				require.NotNil(t, board)
				require.Equal(t, board.Rows(), tc.rows)
				require.Equal(t, board.Cols(), tc.cols)

				require.NotNil(t, board.Cells)
				require.Len(t, board.Cells, tc.rows)

				for i, rowOfCells := range board.Cells {
					require.NotNil(t, rowOfCells)
					require.NotEmpty(t, rowOfCells)
					require.Len(t, rowOfCells, tc.cols-i%2)

					for j, cell := range rowOfCells {
						if cell != nil {
							require.Equal(t, cell.Row(), i)
							require.Equal(t, cell.Col(), j)

							neighbors := cell.GetNeighbors(board)

							require.NotEmpty(t, neighbors)
						}
					}
				}

				checkConnectivity(t, *board)
			}
		})
	}
}

func checkConnectivity(t *testing.T, board game.Board) {
	visited := make(map[game.Cell]bool)
	q := queue.New()

	// Find the first non-empty cell on the board and start Breadth-first search from it
	var startCell game.Cell
	for _, rowOfCells := range board.Cells {
		for _, cell := range rowOfCells {
			if cell != nil {
				startCell = cell
				break
			}
		}
		if startCell != nil {
			break
		}
	}

	q.Enqueue(startCell)
	visited[startCell] = true

	for q.Len() > 0 {
		currCell := q.Dequeue().(game.Cell)
		neighbors := currCell.GetNeighbors(&board)

		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				q.Enqueue(neighbor)
				visited[neighbor] = true
			}
		}
	}

	// Check that all cells on the board have been visited
	for _, rowOfCells := range board.Cells {
		for _, cell := range rowOfCells {
			if cell != nil {
				_, ok := visited[cell]
				require.True(t, ok, "All cells must be connected")
			}
		}
	}
}
