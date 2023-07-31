package main

import (
	"errors"
	"math/rand"
	"strconv"
)

type Creature1 struct {
	id       int
	symbol   byte
	hp       int
	oriHP    int
	speed    int
	oriSpeed int
	typeDesc string
	moving   bool
}

func (b *Board) newCreature1Object(mutate bool, parent ...*Creature1) (*Creature1, error) {
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

	c1 := Creature1{
		id:       b.creatureIdCtr["creature1"],
		symbol:   getObjectSymbol("Creature1"),
		oriHP:    250,
		hp:       250,
		speed:    speed,
		oriSpeed: speed,
		typeDesc: "creature1",
		moving:   true,
	}

	b.creatureIdCtr["creature1"] += 1

	addMessageToCurrentGamelog("Creature1 object with ID: "+
		strconv.Itoa(c1.id)+" added to the board", 2)

	return &c1, nil
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

func (c *Creature1) updateTick() string {
	if c.moving && c.hp > 0 {
		c.speed -= 1
		if c.speed <= 0 {
			c.speed = c.oriSpeed
			c.hp -= 5 + (10 / c.speed)
			return "move"
		}
	} else if c.hp <= 0 {
		return "dead"
	}

	return "error"
}

func (c *Creature1) heal(val int) {
	prio := 2

	if val < 0 {
		prio = 1
	}

	addMessageToCurrentGamelog("Creature 1 with id "+
		strconv.Itoa(c.id)+" healed for: "+
		strconv.Itoa(val), prio)
	c.hp += val
	c.moving = false
}

func (c *Creature1) isDead() bool {
	if c.hp <= 0 {
		return true
	}

	return false
}

func (c *Creature1) resetValues() {
	c.hp = c.oriHP
	c.speed = c.oriSpeed
	c.moving = true
}

func (c *Creature1) ifOffspring() bool {
	if c.hp > int(float32(c.oriHP)*1.1) {
		return true
	}

	return false
}

func (c *Creature1) getHP() int {
	return c.hp
}

func (c *Creature1) getId() int {
	return c.id
}

func (c *Creature1) getSymbol() byte {
	return c.symbol
}

func (c *Creature1) getSpeed() int {
	return c.speed
}

func (c *Creature1) isMoving() bool {
	return c.moving
}

func (c *Creature1) getType() string {
	return c.typeDesc
}

func (c *Creature1) kill() {
	c.hp = 0
}

func (c *Creature1) getOriHP() int {
	return c.oriHP
}

func (c *Creature1) getIdAsString() string {
	return c.typeDesc + " (" + strconv.Itoa(c.id) + ")"
}
