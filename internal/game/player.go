package game

import "errors"

var (
	errAttackTurnExpired     = errors.New("Attack turn has expired")
	errUpgradeTurnNotReached = errors.New("Upgrade time has not been reached")
	errNotEnoughPoints       = errors.New("not enough points to upgrade")
)

type player interface {
	Id() int // ID of player

	attack() error
	endAttack()
	upgrade(points int) error
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
	p.attacking = false
}

// Method for upgrading a cell owned by a player
func (p *Player) upgrade(points int) error {
	if p.attacking {
		return errUpgradeTurnNotReached
	}
	if p.points < points {
		return errNotEnoughPoints
	}
	p.points -= points

	return nil
}

func (p *Player) endUpgrade() {
	p.attacking = true
}

func (p *Player) addCell() {
	p.cellsCounter++
}

func (p *Player) deleteCell() {
	p.cellsCounter--

	// TODO: Realize loosing
}
