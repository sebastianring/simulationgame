package main

import (
	"errors"
	"math/rand"
	"strconv"
)

type Creature2 struct {
	id       int
	symbol   byte
	hp       int
	oriHP    int
	speed    int
	oriSpeed int
	typeDesc string
	moving   bool
}

func (b *Board) newCreature2Object(mutate bool, parent ...*Creature2) (*Creature2, error) {
	var speed int

	if len(parent) == 0 {
		speed = 5
	} else if len(parent) == 1 {
		for _, creature := range parent {
			speed = creature.speed
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
		id:       b.creatureIdCtr["creature2"],
		symbol:   getObjectSymbol("Creature2"),
		oriHP:    175,
		hp:       175,
		speed:    speed,
		oriSpeed: speed,
		typeDesc: "creature2",
		moving:   true,
	}

	b.creatureIdCtr["creature2"] += 1

	addMessageToCurrentGamelog("Creature2 object with ID: "+
		strconv.Itoa(c2.id)+" added to the board", 2)

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

func (c *Creature2) updateTick() string {
	if c.moving && c.hp > 0 {
		c.speed -= 1
		if c.speed == 0 {
			c.speed = c.oriSpeed
			c.hp -= 5 + (10 / c.speed)
			return "move"
		}
	} else if c.hp <= 0 {
		return "dead"
	}

	return "error"
}

func (c *Creature2) heal(val int) {
	addMessageToCurrentGamelog("Creature 2 with id "+
		strconv.Itoa(c.id)+" healed for: "+
		strconv.Itoa(c.oriHP), 2)
	c.hp += val
	c.moving = false
}

func (c *Creature2) isDead() bool {
	if c.hp <= 0 {
		return true
	}

	return false
}

func (c *Creature2) resetValues() {
	c.hp = c.oriHP
	c.speed = c.oriSpeed
	c.moving = true
}

func (c *Creature2) ifOffspring() bool {
	if c.hp > int(float32(c.oriHP)*1.1) {
		return true
	}

	return false
}

func (c *Creature2) getHP() int {
	return c.hp
}

func (c *Creature2) getId() int {
	return c.id
}

func (c *Creature2) getSymbol() byte {
	return c.symbol
}

func (c *Creature2) getSpeed() int {
	return c.speed
}

func (c *Creature2) isMoving() bool {
	return c.moving
}

func (c *Creature2) getType() string {
	return c.typeDesc
}

func (c *Creature2) kill() {
	c.hp = 0
}

func (c *Creature2) getOriHP() int {
	return c.oriHP
}
