package simulationgame

import (
	"errors"
	"math/rand"
	"strconv"
)

type Creature2 struct {
	Id       int    `json:"id"`
	Symbol   []byte `json:"symbol"`
	Pos      Pos    `json:"pos"`
	Hp       int    `json:"hp"`
	OriHP    int    `json:"ori_hp"`
	Speed    int    `json:"speed"`
	OriSpeed int    `json:"ori_speed"`
	TypeDesc string `json:"type_desc"`
	Moving   bool   `json:"moving"`
}

func (b *Board) newCreature2Object(mutate bool, parent ...*Creature2) (*Creature2, error) {
	var speed int

	if len(parent) == 0 {
		speed = 5
	} else if len(parent) == 1 {
		for _, creature := range parent {
			speed = creature.Speed
		}
	} else {
		return nil, errors.New("Too many parents")
	}

	if mutate {
		chance := rand.Intn(100)

		if chance < 33 {
			speed++
		} else if chance < 67 {
			speed--
		}
	}

	c2 := Creature2{
		Id: b.CreatureIdCtr["Creature2"],
		// symbol:   getObjectSymbol("Creature2"),
		Symbol:   getObjectSymbolWColor("Creature2"),
		OriHP:    250,
		Hp:       250,
		Speed:    speed,
		OriSpeed: speed,
		TypeDesc: "Creature2",
		Moving:   true,
	}

	b.CreatureIdCtr["Creature2"] += 1

	addMessageToCurrentGamelog("Creature2 object with ID: "+
		strconv.Itoa(c2.Id)+" added to the board", 2)

	return &c2, nil
}

//
// func getMutationChanges(creaturename string, oriqty float32) {
// 	mutationinterval := int(mutationrate[creaturename] * oriqty)
//
// 	if mutationinterval < 10 {
//
// 	}
//
// }

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL THE SPECIFIC CREATURE FUNCTIONS -------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //

func (c *Creature2) updateTick() TickStatus {
	if c.Moving && c.Hp > 0 {
		c.Speed -= 1
		if c.Speed == 0 {
			c.Speed = c.OriSpeed
			c.Hp -= 5 + (10 / c.Speed)
			return StatusMove
		}
	} else if c.Hp <= 0 {
		return StatusDead
	}

	return StatusInactive
}

func (c *Creature2) heal(val int) {
	addMessageToCurrentGamelog("Creature 2 with id "+
		strconv.Itoa(c.Id)+" healed for: "+
		strconv.Itoa(c.OriHP), 2)
	c.Hp += val
	c.Moving = false
}

func (c *Creature2) isDead() bool {
	if c.Hp <= 0 {
		return true
	}

	return false
}

func (c *Creature2) resetValues() {
	c.Hp = c.OriHP
	c.Speed = c.OriSpeed
	c.Moving = true
}

func (c *Creature2) ifOffspring() bool {
	if c.Hp > int(float32(c.OriHP)*1.25) {
		return true
	}

	return false
}

func (c *Creature2) getHP() int {
	return c.Hp
}

func (c *Creature2) getId() int {
	return c.Id
}

func (c *Creature2) getSymbol() []byte {
	return c.Symbol
}

func (c *Creature2) getSpeed() int {
	return c.Speed
}

func (c *Creature2) isMoving() bool {
	return c.Moving
}

func (c *Creature2) getType() string {
	return c.TypeDesc
}

func (c *Creature2) kill() {
	c.Hp = 0
}

func (c *Creature2) getOriHP() int {
	return c.OriHP
}

func (c *Creature2) getIdAsString() string {
	return c.TypeDesc + " (" + strconv.Itoa(c.Id) + ")"
}

func (c *Creature2) getPos() Pos {
	return c.Pos
}

func (c *Creature2) setPos(pos Pos) {
	c.Pos = pos
}
