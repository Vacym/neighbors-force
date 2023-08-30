package apiserver

import (
	"errors"

	"github.com/Vacym/neighbors-force/internal/game"
)

var (
	errGameIsNotExist           = errors.New("game has not been created yet")
	errIndexOutOfRange          = errors.New("index out of range")
	errIncorrectDifficultiesLen = errors.New("incorrect slice of difficulties")
)

// gameBox holds a reference to the current game and the user's ID.
type gameBox struct {
	Game         *game.Game
	UserId       int
	difficulties []int
}

// User represents a user and their actions in the game.
type User struct {
	GameBox gameBox
}

// NewUser creates a new User instance.
func NewUser() *User {
	return &User{}
}

// me returns the current user's player instance.
func (u *User) me() game.Player {
	return u.GameBox.Game.Players[u.GameBox.UserId]
}

// createGame sets the current game and user's ID in the user's GameBox.
func (u *User) createGame(g *game.Game, id int, difficulties []int) {
	u.GameBox.Game = g
	u.GameBox.UserId = id

	// Ensure that the difficulties slice has the same length as the number of players in the game.
	difficultyCount := len(g.Players)
	if len(difficulties) < difficultyCount {
		difficulties = append(difficulties, make([]int, difficultyCount-len(difficulties))...)
	}
	u.GameBox.difficulties = difficulties[:difficultyCount]
}

// attack performs an attack from a source cell to a target cell.
func (u *User) attack(from, to game.Coords) error {
	g := u.GameBox.Game

	if g == nil {
		return errGameIsNotExist
	}

	if !g.Board.IsInsideBoard(from) || !g.Board.IsInsideBoard(to) {
		return errIndexOutOfRange
	}

	fromCell := g.Board.Cells[from.Row][from.Col]
	toCell := g.Board.Cells[to.Row][to.Col]
	return g.Attack(g.Players[u.GameBox.UserId], fromCell, toCell)
}

// endAttack ends the current attack phase for the user.
func (u *User) endAttack() error {
	g := u.GameBox.Game

	if g == nil {
		return errGameIsNotExist
	}

	return g.EndAttack(u.me())
}

// makeUpgrade upgrades a cell owned by the user.
func (u *User) makeUpgrade(cellCoords game.Coords, levels int) error {
	g := u.GameBox.Game

	if g == nil {
		return errGameIsNotExist
	}

	cell, err := g.Board.GetCell(cellCoords)
	if err != nil {
		return err
	}

	return g.Upgrade(u.me(), cell, levels)
}

// endTurn ends the current turn for the user.
func (u *User) endTurn() error {
	g := u.GameBox.Game

	if g == nil {
		return errGameIsNotExist
	}

	return g.EndTurn(g.Players[u.GameBox.UserId])
}
