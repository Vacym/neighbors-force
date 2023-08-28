package game

import "errors"

var (
	errAttackTurnExpired     = errors.New("attack turn has expired")
	errUpgradeTurnNotReached = errors.New("upgrade time has not been reached")
	errNotEnoughPoints       = errors.New("not enough points to upgrade")
	errAttackAlreadyFinished = errors.New("attack already finished")
)

// Player represents a player in the game.
type Player interface {
	// Id returns the ID of the player.
	Id() int

	// Points returns the points of the player.
	Points() int

	// CellsCount returns number of cells owned by the player
	CellsCount() int

	// attack initiates an attack for the player's turn.
	attack() error

	// endAttack concludes the attack phase for the player.
	endAttack() error

	// upgrade upgrades a cell owned by the player.
	upgrade(cell Cell, levels int) error

	// endUpgrade concludes the upgrade phase for the player's turn.
	endUpgrade()

	// addCell increments the player's cell count.
	addCell()

	// deleteCell decrements the player's cell count.
	// It returns true if there are still cells and the player remaining, otherwise false.
	deleteCell() bool

	// toMap converts the player's information into a map for serialization.
	toMap() map[string]interface{}
}

// player implements the Player interface.
type player struct {
	id         int  // ID of player
	points     int  // Points that can be spent on upgrading cells or attacking
	cellsCount int  // Count of cells, owned user
	attacking  bool // phase of player turn
}

// newPlayer creates a new Player with the given ID.
func newPlayer(id int) *player {
	return &player{
		id:        id,
		attacking: true,
	}
}

// Id returns the ID of the player.
func (p *player) Id() int {
	return p.id
}

// Points returns the points of the player.
func (p *player) Points() int {
	return p.points
}

// CellsCount returns number of cells owned by the player
func (p *player) CellsCount() int {
	return p.cellsCount
}

// attack performs an attack for the player.
func (p *player) attack() error {
	if !p.attacking {
		return errAttackTurnExpired
	}

	return nil
}

// endAttack ends the attack phase for the player.
func (p *player) endAttack() error {
	if !p.attacking {
		return errAttackAlreadyFinished
	}

	p.attacking = false
	p.countPoints()
	return nil
}

// countPoints calculates the points based on the player's cell count.
func (p *player) countPoints() {
	p.points += p.cellsCount
}

// upgrade upgrades a cell owned by the player.
func (p *player) upgrade(cell Cell, levels int) error {
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

// endUpgrade ends the upgrade phase for the player.
func (p *player) endUpgrade() {
	if p.attacking {
		p.endAttack()
	}

	p.attacking = true
}

// addCell increments the count of cells owned by the player.
func (p *player) addCell() {
	p.cellsCount++
}

// deleteCell decrements the count of cells owned by the player.
//It returns true if no cells remain and the player has no cells, otherwise false.
func (p *player) deleteCell() bool {
	p.cellsCount--
	return p.cellsCount == 0
}

// toMap converts the player's information into a map for serialization.
func (p *player) toMap() map[string]interface{} {
	return map[string]interface{}{
		"id":          p.id,
		"points":      p.points,
		"cells_count": p.cellsCount,
		"attacking":   p.attacking,
	}
}
