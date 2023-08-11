package bot

import (
	"fmt"
	"math/rand"

	"github.com/Vacym/neighbors-force/internal/game"
)

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func DoTurn(g *game.Game, player game.Player) error {
	// temporary implementation
	for _, row := range g.Board.Cells {
		for _, cell := range row {
			if cell == nil {
				continue
			}

			if cell.Owner() == player {
				err := attack(g, cell)
				if err != nil {
					return err
				}
			}
		}
	}

	g.EndAttack(player)
	return fmt.Errorf("no cells found for player %v", player.Id())
}

func attack(g *game.Game, cell game.Cell) error {
	player := cell.Owner()
	neighbors := cell.GetNeighbors(g.Board)
	alienNeighbors := filter(neighbors, func(neighbor game.Cell) bool { return neighbor.Owner() != player })
	if len(alienNeighbors) > 0 {
		to := alienNeighbors[rand.Intn(len(alienNeighbors))]
		g.Attack(player, cell, to)
		if to.Owner() == cell.Owner() {
			err := attack(g, to)
			return err
		}
	}
	return nil
}

func DoUpgrade(g *game.Game, player game.Player) error {
	// ChatGPT powered
	// Find all cells owned by the player and store them in an array
	var ownedCells []game.Cell
	for _, row := range g.Board.Cells {
		for _, cell := range row {
			if cell == nil {
				continue
			}
			if cell.Owner() == player {
				ownedCells = append(ownedCells, cell)
			}
		}
	}

	// Upgrade a random cell from the ownedCells array if the player has enough points
	for len(ownedCells) > 0 && player.Points() > 0 {
		// Choose a random index from the array of ownedCells
		randomIndex := rand.Intn(len(ownedCells))
		cellToUpgrade := ownedCells[randomIndex]

		// Calculate the cost to upgrade the cell
		upgradeCost := (cellToUpgrade.Level() * (cellToUpgrade.Level() + 1)) / 2

		// Check if the player has enough points to perform the upgrade
		if player.Points() >= upgradeCost {
			// Perform the upgrade and deduct the points from the player
			err := g.Upgrade(player, cellToUpgrade, 1)
			if err != nil {
				return err
			}
		}

		// Remove the cell from the ownedCells array regardless of the upgrade result
		ownedCells = append(ownedCells[:randomIndex], ownedCells[randomIndex+1:]...)
	}

	// End the turn for the player
	g.EndTurn(player)
	return nil
}
