package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/Vacym/neighbors-force/internal/game"
)

type BotAction struct {
	Attack  [][]int `json:"attack"`
	Upgrade []int   `json:"upgrade"`
}

func DoTurn(g *game.Game, player game.Player) error {
	// return DoTurnVanilla(g, player)
	for !g.IsFinished() {
		body, err := getJSONResponse(g, "/ai_attack")
		if err != nil {
			g.EndAttack(player)
			fmt.Println("Ошибка при получении JSON-ответа:", err)
			return err
		}

		fmt.Println("Attack Response body:", string(body))

		// Чтение и размаршалирование JSON-ответа
		var action BotAction
		err = json.Unmarshal(body, &action)
		if err != nil {
			g.EndAttack(player)
			fmt.Println("Ошибка при размаршалировании JSON:", err)
			return err
		}

		// Если бот вернул None (пропустил ход)
		if action.Attack == nil {
			g.EndAttack(player)
			return fmt.Errorf("no cells found for player %v", player.Id())
		}

		cell := action.Attack[0]
		from := g.Board.Cells[cell[0]][cell[1]]

		cell = action.Attack[1]
		to := g.Board.Cells[cell[0]][cell[1]]

		g.Attack(player, from, to)
	}
	return nil
}

func DoUpgrade(g *game.Game, player game.Player) error {
	// return DoUpgradeVanilla(g, player)
	for {
		body, err := getJSONResponse(g, "/ai_upgrade")
		if err != nil {
			g.EndTurn(player)
			fmt.Println("Ошибка при получении JSON-ответа:", err)
			return err
		}

		fmt.Println("Upgrade Response body:", string(body))

		// Reading and unmarshaling JSON response
		var action BotAction
		err = json.Unmarshal(body, &action)
		if err != nil {
			g.EndTurn(player)
			fmt.Println("Ошибка при размаршалировании JSON:", err)
			return err
		}

		// If the bot returns None (skips the turn)
		if action.Upgrade == nil {
			g.EndTurn(player)
			return nil
		}

		cell := action.Upgrade
		cellToUpgrade := g.Board.Cells[cell[0]][cell[1]]

		g.Upgrade(player, cellToUpgrade, 1)
	}
}

func getJSONResponse(g *game.Game, path string) ([]byte, error) {
	// Powered by ChatGPT
	jsonData := g.ToMap()

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	url := "http://127.0.0.1:8000" + path // Replace with the actual URL of the FastAPI server
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	// Set necessary headers
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func DoTurnVanilla(g *game.Game, player game.Player) error {
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

func DoUpgradeVanilla(g *game.Game, player game.Player) error {
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

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
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
