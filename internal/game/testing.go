package game

func TestGame() *Game {
	game, _ := NewGameWithBoard(TestBoard(), 2)
	return game
}

func TestBoard() *Board {
	board, _ := NewBoard(5, 5)
	return board
}

func TestPlayer(id int) *Player {
	player := NewPlayer(id)
	player.points = 10
	player.cellsCounter = 10
	return player
}
