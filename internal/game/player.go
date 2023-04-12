package game

type Player struct {
	id        int
	points    int // Points that can be spent on upgrading cells or attacking
	attacking bool
}

// NewPlayer creates a new Player with the given ID.
func NewPlayer(id int) *Player {
	return &Player{id: id}
}

func (p *Player) Attack(cell *Cell) bool {
	if p.attacking {
		return false
	}

	p.attacking = true
	return true
}

// Method for upgrading a cell owned by a player
func (player *Player) Upgrade(cell *Cell) bool {
	// Implementation

	return true
}
