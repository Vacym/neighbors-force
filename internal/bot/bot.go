package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"sort"

	"github.com/Vacym/neighbors-force/internal/game"
)

var (
	errIncorrectDifficulty = errors.New("incorrect difficulty level")
	errNoPlayerCells       = func(playerId int) error {
		return fmt.Errorf("no cells found for player %v", playerId)
	}
)

type BotAction struct {
	Attack  [][]int `json:"attack"`
	Upgrade []int   `json:"upgrade"`
}

var attackHandlers = [...]func(g *game.Game, player game.Player) error{
	DoAttackEasy, DoAttackMedium, DoAttackAPI, DoAttackAPI,
}

var upgradeHandlers = [...]func(g *game.Game, player game.Player) error{
	DoUpgradeEasy, DoUpgradeMedium, DoUpgradeAPI, DoUpgradeAPI,
}

func DoAttack(g *game.Game, player game.Player, difficulty int) error {
	if difficulty > len(attackHandlers) {
		return errIncorrectDifficulty
	}

	for !g.IsFinished() {
		err := attackHandlers[difficulty](g, player)
		if err != nil {
			g.EndAttack(player)
			return err
		}
	}
	return nil
}

func DoUpgrade(g *game.Game, player game.Player, difficulty int) error {
	if difficulty > len(upgradeHandlers) {
		return errIncorrectDifficulty
	}

	for !g.IsFinished() {
		err := upgradeHandlers[difficulty](g, player)
		if err != nil {
			g.EndTurn(player)
			return err
		}
	}
	return nil
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

func DoAttackAPI(g *game.Game, player game.Player) error {
	body, err := getJSONResponse(g, "/ai_attack")
	if err != nil {
		fmt.Println("Error for receiving JSON:", err)
		return err
	}

	fmt.Println("Attack Response body:", string(body))

	// Reading and unmarshaling JSON response
	var action BotAction
	err = json.Unmarshal(body, &action)
	if err != nil {
		fmt.Println("Error for unmarshaling JSON:", err)
		return err
	}

	// If the bot returns None (skips the turn)
	if action.Attack == nil {
		return fmt.Errorf("no cells found for player %v", player.Id())
	}

	cell := action.Attack[0]
	from := g.Board.Cells[cell[0]][cell[1]]

	cell = action.Attack[1]
	to := g.Board.Cells[cell[0]][cell[1]]

	g.Attack(player, from, to)
	return nil
}

func DoUpgradeAPI(g *game.Game, player game.Player) error {
	body, err := getJSONResponse(g, "/ai_upgrade")
	if err != nil {
		fmt.Println("Error for receiving JSON:", err)
		return err
	}

	fmt.Println("Upgrade Response body:", string(body))

	// Reading and unmarshaling JSON response
	var action BotAction
	err = json.Unmarshal(body, &action)
	if err != nil {
		fmt.Println("Error for unmarshaling JSON:", err)
		return err
	}

	// If the bot returns None (skips the turn)
	if action.Upgrade == nil {
		return fmt.Errorf("no cells found for player %v", player.Id())
	}

	cell := action.Upgrade
	cellToUpgrade := g.Board.Cells[cell[0]][cell[1]]

	g.Upgrade(player, cellToUpgrade, 1)
	return nil
}

func DoAttackMedium(g *game.Game, player game.Player) error {
	var bestFrom, bestTo game.Cell
	bestScore := 0
	for _, row := range g.Board.Cells {
		for _, cell := range row {
			if cell == nil {
				continue
			}

			if cell.Owner() == player && cell.Power() > 1 {
				neighbors := cell.GetNeighbors(g.Board)
				alienNeighbors := filter(neighbors, func(neighbor game.Cell) bool { return neighbor.Owner() != player })
				for _, to := range alienNeighbors {
					score := calculateScore(cell, to)
					if bestScore == 0 || score > bestScore {
						bestScore = score
						bestFrom = cell
						bestTo = to
					}
				}
			}
		}
	}

	if bestScore > 0 {
		g.Attack(player, bestFrom, bestTo)
		return nil
	}
	return fmt.Errorf("no cells found for player %v", player.Id())
}

func DoUpgradeMedium(g *game.Game, player game.Player) error {
	if player.Points() == 0 {
		return errNoPlayerCells(player.Id())
	}

	// ChatGPT powered
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

	maxRow, maxCol := g.Board.Rows(), g.Board.Cols()
	mainCell := [4][2]int{
		{0, 0},
		{maxRow, maxCol},
		{maxRow, 0},
		{0, maxCol},
	}[player.Id()]

	sort.Slice(ownedCells, func(i, j int) bool {
		cellI, cellJ := ownedCells[i], ownedCells[j]
		distanceI := math.Abs(float64(mainCell[0]-cellI.Row())) + math.Abs(float64(mainCell[1]-cellI.Col()))
		distanceJ := math.Abs(float64(mainCell[0]-cellJ.Row())) + math.Abs(float64(mainCell[1]-cellJ.Col()))
		if distanceI != distanceJ {
			return distanceI > distanceJ
		}
		return cellI.Level() > cellJ.Level()
	})

	for _, cell := range ownedCells {
		upgradeCost := cell.Level() * (cell.Level() + 1) / 2
		if player.Points() >= upgradeCost {
			return g.Upgrade(player, cell, 1)
		}
	}
	return errNoPlayerCells(player.Id())
}

func calculateScore(from game.Cell, to game.Cell) int {
	// Simple Evaluation Function implementation (temporary)

	// If cell is free
	if to.Owner() == nil {
		return from.Power()
	}

	// If our cell power lower than other
	if from.Power() <= to.Power() {
		return 100 + (to.Power() - from.Power())
	}

	// Otherwise
	return 10000 + (to.Power() - from.Power())
}

func DoAttackEasy(g *game.Game, player game.Player) error {
	for _, row := range g.Board.Cells {
		for _, cell := range row {
			if cell == nil {
				continue
			}

			if cell.Owner() == player && cell.Power() > 1 {
				attackEasy(g, cell)
			}
		}
	}
	return errNoPlayerCells(player.Id())
}

func DoUpgradeEasy(g *game.Game, player game.Player) error {
	if player.Points() == 0 {
		return errNoPlayerCells(player.Id())
	}

	// ChatGPT powered
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

	for len(ownedCells) > 0 && player.Points() > 0 {
		randomIndex := rand.Intn(len(ownedCells))
		cellToUpgrade := ownedCells[randomIndex]

		upgradeCost := (cellToUpgrade.Level() * (cellToUpgrade.Level() + 1)) / 2
		if player.Points() >= upgradeCost {
			err := g.Upgrade(player, cellToUpgrade, 1)
			if err != nil {
				return err
			}
		}
		ownedCells = append(ownedCells[:randomIndex], ownedCells[randomIndex+1:]...)
	}
	return errNoPlayerCells(player.Id())
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func attackEasy(g *game.Game, cell game.Cell) error {
	player := cell.Owner()
	neighbors := cell.GetNeighbors(g.Board)
	alienNeighbors := filter(neighbors, func(neighbor game.Cell) bool { return neighbor.Owner() != player })
	if len(alienNeighbors) > 0 {
		to := alienNeighbors[rand.Intn(len(alienNeighbors))]
		g.Attack(player, cell, to)
		if to.Owner() == cell.Owner() {
			return attackEasy(g, to)
		}
	}
	return nil
}
