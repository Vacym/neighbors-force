package game

func TestGame() *Game {
	game, _ := NewGame(5, 5, 2)
	return game
}

func TestBoard() *Board {
	board, _ := NewBoard(5, 5)
	return board
}

func TestPlayer(id int) *Player {
	player := NewPlayer(id)
	player.points = 20
	player.cellsCount = 10
	return player
}

func TestBoardAttack() (*Board, []player) {
	// This method of map creation will only be used until
	// the implementation of saving custom boards
	/*
		(1) {6} (5) (1) ( )
		  (6) {1} ( ) ( ) ( )
		( ) {6} ( ) ( ) ( )
		  ( ) ( ) ( ) ( ) ( )
		( ) ( ) ( ) ( ) ( )

		(power) - p1
		{power} - p2
	*/

	p1 := TestPlayer(0)
	p2 := TestPlayer(1)

	board, _ := NewBoard(5, 5)
	board.Cells[0][0] = newCellWithParameters(0, 0, 6, 1, p1)
	board.Cells[1][0] = newCellWithParameters(1, 0, 1, 6, p1)
	board.Cells[0][1] = newCellWithParameters(0, 1, 1, 6, p2)
	board.Cells[1][1] = newCellWithParameters(1, 1, 6, 1, p2)
	board.Cells[0][2] = newCellWithParameters(0, 2, 1, 5, p1)
	board.Cells[0][3] = newCellWithParameters(0, 3, 5, 1, p1)
	board.Cells[2][1] = newCellWithParameters(2, 1, 1, 3, p2)

	return board, []player{p1, p2}
}

func TestGameAttack() (*Game, error) {
	return NewGameWithBoard(TestBoardAttack())
}
