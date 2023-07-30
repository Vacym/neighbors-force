package apiserver

import (
	"errors"

	"github.com/Vacym/neighbors-force/internal/game"
)

var (
	errGameIsNotExist  = errors.New("game has not been created yet")
	errIndexOutOfRange = errors.New("index out of range")
)

type gameBox struct {
	Game   *game.Game
	UserId int
}

type User struct {
	GameBox gameBox
}

func NewUser() *User {
	return &User{}
}

func (u *User) me() game.PlayerInterface {
	return u.GameBox.Game.Players[u.GameBox.UserId]
}

func (u *User) createGame(g *game.Game, id int) {
	u.GameBox.Game = g
	u.GameBox.UserId = id
}

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

func (u *User) endAttack() error {
	g := u.GameBox.Game

	if g == nil {
		return errGameIsNotExist
	}

	return g.EndAttack(u.me())
}

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

func (u *User) endTurn() error {
	g := u.GameBox.Game

	if g == nil {
		return errGameIsNotExist
	}

	return g.EndTurn(g.Players[u.GameBox.UserId])
}
