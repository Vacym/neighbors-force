package main

import (
	"fmt"

	"github.com/Vacym/neighbor-force/internal/game"
)

func main() {
	game, _ := game.NewGame(3, 5, 2)

	fmt.Println(&game)
}
