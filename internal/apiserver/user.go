package apiserver

import (
	"github.com/Vacym/neighbors-force/internal/game"
)

type User struct {
	Game *game.Game
}

func NewUser() *User {
	return &User{}
}
