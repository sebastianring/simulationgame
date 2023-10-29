package simulationgame

import (
	"errors"
	"strconv"
)

// Potentially create a CREATURE FACTORY?

type Creature1 struct {
	Id              int             `json:"id"`
	Symbol          []byte          `json:"symbol"`
	Pos             Pos             `json:"pos"`
	Hp              int             `json:"hp"`
	OriHP           int             `json:"ori_hp"`
	Speed           float64         `json:"speed"`
	OriSpeed        float64         `json:"ori_speed"`
	ProcScanChance  float64         `json:"proc_scan_chance"`
	TypeDesc        string          `json:"type_desc"`
	BoardObjectType BoardObjectType `json:"board_object_type"`
	// CreatureObjectType CreatureObjectType `json:"creature_object_type"`
	Moving bool `json:"moving"`
}

func (b *Board) newCreature1Object(mutate bool, parent ...*Creature1) (*Creature1, error) {
	var speed float64
	var procScanChance float64

	if len(parent) == 0 {
		speed = 5
		procScanChance = 50
	} else if len(parent) == 1 {
		if mutate {
			tempspeed, err := b.MutationManager.getVariableValue(speedVariable, parent[0])

			if err != nil {
				addMessageToCurrentGamelog(err.Error(), 1)
				return nil, errors.New("Error creating new creature.")
			}

			speed = tempspeed

			tempProcChance, err := b.MutationManager.getVariableValue(scanVariable, parent[0])

			if err != nil {
				addMessageToCurrentGamelog(err.Error(), 1)
				return nil, errors.New("Error creating new creature.")
			}

			procScanChance = tempProcChance
		}

	} else {
		return nil, errors.New("Too many parents")
	}

	c1 := Creature1{
		Id:              b.CreatureIdCtr[Creature1Type],
		Symbol:          getObjectSymbolWColor(Creature1Type),
		OriHP:           500,
		Hp:              500,
		Speed:           speed,
		OriSpeed:        speed,
		ProcScanChance:  procScanChance,
		TypeDesc:        "Creature1",
		BoardObjectType: Creature1Type,
		Moving:          true,
	}

	b.CreatureIdCtr[Creature1Type] += 1

	addMessageToCurrentGamelog("Creature1 object with ID: "+
		strconv.Itoa(c1.Id)+" added to the board", 2)

	return &c1, nil
}

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL THE SPECIFIC CREATURE FUNCTIONS -------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //

func (c *Creature1) updateTick() TickStatus {
	if c.Moving && c.Hp > 0 {
		c.Speed -= 1
		if c.Speed <= 0 {
			c.Speed = c.OriSpeed
			c.Hp -= 5 + int(c.Speed)
			return StatusMove
		}
	} else if c.Hp <= 0 {
		return StatusDead
	}

	return StatusInactive
}

func (c *Creature1) heal(val int) {
	prio := 2
	if val < 0 {
		prio = 1
	}

	addMessageToCurrentGamelog("Creature 1 with id "+
		strconv.Itoa(c.Id)+" healed for: "+
		strconv.Itoa(val), prio)

	c.Hp += val
	c.isFull()
}

func (c *Creature1) isDead() bool {
	if c.Hp <= 0 {
		return true
	}

	return false
}

func (c *Creature1) resetValues() {
	c.Hp = c.OriHP
	c.Speed = c.OriSpeed
	c.Moving = true
}

func (c *Creature1) ifOffspring() bool {
	if c.Hp > int(float32(c.OriHP)*1.25) {
		return true
	}

	return false
}

func (c *Creature1) getHP() int {
	return c.Hp
}

func (c *Creature1) getId() int {
	return c.Id
}

func (c *Creature1) getSymbol() []byte {
	return c.Symbol
}

func (c *Creature1) getSpeed() float64 {
	return c.Speed
}

func (c *Creature1) isMoving() bool {
	return c.Moving
}

func (c *Creature1) getType() string {
	return c.TypeDesc
}

func (c *Creature1) kill() {
	c.Hp = 0
}

func (c *Creature1) getOriHP() int {
	return c.OriHP
}

func (c *Creature1) getIdAsString() string {
	return c.TypeDesc + " (" + strconv.Itoa(c.Id) + ")"
}

func (c *Creature1) getPos() Pos {
	return c.Pos
}

func (c *Creature1) setPos(pos Pos) {
	c.Pos = pos
}

func (c *Creature1) getBoardObjectType() BoardObjectType {
	return c.BoardObjectType
}

func (c *Creature1) getScanProcChance() float64 {
	return c.ProcScanChance
}

func (c *Creature1) scan() {
	c.Hp -= 5 + int(c.Speed)
}

func (c *Creature1) isFull() {
	if c.Hp > c.OriHP+c.OriHP/2 {
		c.Moving = false
	}
}
