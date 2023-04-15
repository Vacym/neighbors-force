package main

import (
	"fmt"

	"github.com/Vacym/neighbor-force/internal/game"
)

func main() {
	g, _ := game.NewGame(3, 5, 2)

	fmt.Println(&g)

}
