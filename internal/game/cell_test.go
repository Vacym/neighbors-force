package game_test

import (
	"testing"

	"github.com/Vacym/neighbors-force/internal/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCell_attack(t *testing.T) {

	testCases := []struct {
		name          string
		row1, col1    int
		row2, col2    int
		expectedPower int  // for cell2
		isWin         bool // If a player captures a cell due to an attack
		isValid       bool
	}{
		{
			"Valid attack and win 1",
			1, 0,
			1, 1,
			4,
			true,
			true,
		},
		{
			"Valid attack and win 2",
			0, 1,
			0, 0,
			4,
			true,
			true,
		},
		{
			"Too far attack 1",
			1, 0,
			0, 2,
			-1,
			false,
			false,
		},
		{
			"Too far attack 2",
			0, 1,
			1, 2,
			-1,
			false,
			false,
		},
		{
			"Not enough power for attack 1",
			0, 0,
			0, 1,
			-1,
			false,
			false,
		},
		{
			"Not enough power for attack 2",
			1, 1,
			1, 0,
			-1,
			false,
			false,
		},
		{
			"only damage 1",
			0, 2,
			0, 1,
			2,
			false,
			true,
		},
		{
			"only damage 2",
			2, 1,
			1, 0,
			4,
			false,
			true,
		},
		{
			"attack empty ceil 1",
			1, 0,
			2, 0,
			5,
			true,
			true,
		},
		{
			"attack empty ceil 2",
			2, 1,
			3, 0,
			2,
			true,
			true,
		},
		{
			"attacking yourself 1",
			1, 0,
			0, 0,
			-1,
			false,
			false,
		},
		{
			"attacking yourself 2",
			0, 1,
			1, 1,
			-1,
			false,
			false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			board, players := game.TestBoardAttack()
			thisGame, err := game.NewGameWithBoard(board, players)
			require.NoError(t, err)

			cell1 := board.Cells[tc.row1][tc.col1]
			cell2 := board.Cells[tc.row2][tc.col2]

			if cell1.Owner().Id() != 0 {
				thisGame.EndTurn(players[0])
			}

			attacker := cell1.Owner()

			err = thisGame.Attack(attacker, cell1, cell2)

			if tc.isValid {
				require.NoError(t, err)
				// Check the result of cell1.Power() after the attack
				assert.Equal(t, 1, cell1.Power(), "Expected cell power to be %d, but got %d", 1, cell1.Power())
				// Check the result of cell2.Power() after the attack
				assert.Equal(t, tc.expectedPower, cell2.Power(), "Expected cell power to be %d, but got %d", tc.expectedPower, cell2.Power())

				// Check the result of cell2.Owner() after the attack
				assert.Equal(t, attacker, cell1.Owner(), "Expected player1 to be same as cell1.Owner()")
				if tc.isWin {
					assert.Equal(t, attacker, cell2.Owner(), "Expected player1 to be same as cell2.Owner()")
				} else {
					assert.NotEqual(t, attacker, cell2.Owner(), "Expected player2 to be same as cell2.Owner()")
				}
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestCell_upgrade(t *testing.T) {
	// TODO: add tests after releases upgrades

	testCases := []struct {
		name     string
		row, col int
		addLevel int
		isValid  bool
	}{
		{
			"valid upgrade 1",
			0, 0,
			8,
			true,
		},
		{
			"valid upgrade 2",
			0, 1,
			6,
			true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			board, players := game.TestBoardAttack()
			thisGame, err := game.NewGameWithBoard(board, players)
			require.NoError(t, err)

			cell := board.Cells[tc.row][tc.col]

			if cell.Owner().Id() != 0 {
				thisGame.EndTurn(players[0])
			}
			thisGame.EndAttack(cell.Owner())

			err = thisGame.Upgrade(cell.Owner(), cell, tc.addLevel)

			if tc.isValid {
				require.NoError(t, err)
				// Check the result of cell1.Power() after the attack
				assert.Equal(t, 1+tc.addLevel, cell.Level(), "Expected cell level to be %d, but got %d", 1+tc.addLevel, cell.Level())
			} else {
				require.Error(t, err)
			}
		})
	}
}
