package main

import (
	"fmt"

	"github.com/Vacym/neighbors-force/internal/game"
)

func main() {
	g, _ := game.NewGame(3, 5, 2)

	fmt.Println(&g)

}
