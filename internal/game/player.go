package game

import "errors"

var (
	errAttackTurnExpired     = errors.New("Attack turn has expired")
	errUpgradeTurnNotReached = errors.New("Upgrade time has not been reached")
	errNotEnoughPoints       = errors.New("Not enough points to upgrade")
	errAttackAlreadyFinished = errors.New("Attack already finished")
)

type player interface {
	Id() int     // ID of player
	Points() int // points of player

	attack() error
	endAttack() error
	upgrade(cell cell, levels int) error
	endUpgrade()

	addCell()
	deleteCell()

	toMap() map[string]interface{}
}

type PlayerInterface = player

type Player struct {
	id         int  // ID of player
	points     int  // Points that can be spent on upgrading cells or attacking
	cellsCount int  // Count of cells, owned user
	attacking  bool // phase of player turn
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

func (p *Player) Points() int {
	return p.points
}

func (p *Player) attack() error {
	if !p.attacking {
		return errAttackTurnExpired
	}

	return nil
}

func (p *Player) endAttack() error {
	if p.attacking == false {
		return errAttackAlreadyFinished
	}

	p.attacking = false
	p.countPoints()
	return nil
}

func (p *Player) countPoints() {
	p.points += p.cellsCount
}

// Method for upgrading a cell owned by a player
func (p *Player) upgrade(cell cell, levels int) error {
	if p.attacking {
		return errUpgradeTurnNotReached
	}

	// Sum of triangle numbers
	targetLevel := cell.Level() + levels
	cost := ((targetLevel-1)*(targetLevel)*(targetLevel+1) - (cell.Level()-1)*(cell.Level())*(cell.Level()+1)) / 6

	if p.points < cost {
		return errNotEnoughPoints
	}
	p.points -= cost

	return nil
}

func (p *Player) endUpgrade() {
	if p.attacking == true {
		p.endAttack()
	}

	p.attacking = true
}

func (p *Player) addCell() {
	p.cellsCount++
}

func (p *Player) deleteCell() {
	p.cellsCount--

	// TODO: Realize loosing
}

func (p *Player) toMap() map[string]interface{} {
	return map[string]interface{}{
		"id":          p.id,
		"points":      p.points,
		"cells_count": p.cellsCount,
		"attacking":   p.attacking,
	}
}
