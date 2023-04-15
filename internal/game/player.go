package game

import "errors"

var (
	errAttackTurnExpired     = errors.New("Attack turn has expired")
	errUpgradeTurnNotReached = errors.New("Upgrade time has not been reached")
)

type player interface {
	Id() int // ID of player

	attack() error
	endAttack()
	upgrade(cell *Cell) bool
	endUpgrade()

	addCell()
	deleteCell()
}

type Player struct {
	id           int  // ID of player
	points       int  // Points that can be spent on upgrading cells or attacking
	cellsCounter int  // Count of cells, owned user
	attacking    bool // phase of player turn
}

// NewPlayer creates a new Player with the given ID.
func NewPlayer(id int) *Player {
	return &Player{
		id:        id,
		attacking: true,
	}
}

func (p *Player) Id() int {
	return p.id
}

func (p *Player) attack() error {
	if !p.attacking {
		return errAttackTurnExpired
	}

	return nil
}

func (p *Player) endAttack() {
	p.attacking = true
}

// Method for upgrading a cell owned by a player
func (player *Player) upgrade(cell *Cell) bool {
	// Implementation

	return true
}

func (p *Player) endUpgrade() {
	p.attacking = false
}

func (p *Player) addCell() {
	p.cellsCounter++
}

func (p *Player) deleteCell() {
	p.cellsCounter--

	// TODO: Realize loosing
}
