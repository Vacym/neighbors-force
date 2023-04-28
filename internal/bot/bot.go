package bot

import (
	"fmt"
	"math/rand"

	"github.com/Vacym/neighbors-force/internal/game"
)

func DoTurn(g *game.Game, player *game.Player) error {
	// temporary implementation
	for _, row := range g.Board.Cells {
		for _, cell := range row {
			if cell == nil {
				continue
			}

			if cell.Owner() == player {
				neighbors := cell.GetNeighbors(g.Board)
				if len(neighbors) > 0 {
					to := neighbors[rand.Intn(len(neighbors))]
					err := g.Attack(player, cell, to)
					g.EndTurn(player)
					if err != nil {
						return err
					}
					return nil
				}
			}
		}
	}
	return fmt.Errorf("no cells found for player %v", player.Id())
}
