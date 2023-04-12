package main

import (
	"fmt"

	"github.com/Vacym/neighbor-force/internal/game"
)

func main() {
	board, _ := game.NewRectangleBoard(5, 3)

	game, _ := game.NewGame(board, 2)

	fmt.Println(&game)
}
