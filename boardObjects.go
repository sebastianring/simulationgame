package main

import (
	"errors"
	"math/rand"
	"strconv"
)

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL INTERFACES AND GENERAL FUNCTIONS ------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //
// I really need to change architecture of the board .. this is abuse of interfaces. //

var mutationrate map[string]float32
var Creature1IdCtr int
var creatureIdCtr map[string]int

type BoardObject interface {
	getSymbol() byte
}

type CreatureObject interface {
	updateTick() string
	ifOffspring() bool
	getHP() int
	getId() int
	resetValues()
	updateVal(string)
}

func getObjectSymbol(objectname string) byte {
	drawingSymbols := map[string]byte{
		"EmptyObject": 46,  // .....
		"Food":        64,  // @@@@@
		"Creature1":   65,  // AAAAA
		"Creature2":   126, // ~~~~~
	}

	return drawingSymbols[objectname]
}

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL STRUCTS HERE AND THEIR CREATE FUNCTIONS ------ //
// -------------------------------------------------- //
// -------------------------------------------------- //

func initBoardObjects() {
	Creature1IdCtr = 1
	mutationrate = make(map[string]float32)
	mutationrate["creature1"] = 0.1
}

type EmptyObject struct {
	symbol   byte
	typeDesc string
}

func newEmptyObject() *EmptyObject {
	eo := EmptyObject{
		symbol:   getObjectSymbol("EmptyObject"),
		typeDesc: "empty",
	}

	// addMessageToCurrentGamelog("New empty object added", 2)

	return &eo
}

type Food struct {
	symbol   byte
	typeDesc string
}

func newFoodObject() *Food {
	f := Food{
		symbol:   getObjectSymbol("Food"),
		typeDesc: "food",
	}

	addMessageToCurrentGamelog("New food object added", 2)

	return &f
}

type Creature1 struct {
	// CreatureObject
	id       int
	symbol   byte
	hp       int
	oriHP    int
	speed    int
	oriSpeed int
	typeDesc string
	moving   bool
}

func newCreature1Object(mutate bool, parent ...*Creature1) (*Creature1, error) {
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

	if Creature1IdCtr < 1 {
		Creature1IdCtr = 1
	}

	c1 := Creature1{
		id:       Creature1IdCtr,
		symbol:   getObjectSymbol("Creature1"),
		oriHP:    175,
		hp:       175,
		speed:    speed,
		oriSpeed: speed,
		typeDesc: "creature",
		moving:   true,
	}

	Creature1IdCtr++
	addMessageToCurrentGamelog("Creature1 object with ID: "+
		strconv.Itoa(c1.id)+" added to the board", 2)

	return &c1, nil
}

func getMutationChanges(creaturename string, oriqty float32) {
	mutationinterval := int(mutationrate[creaturename] * oriqty)

	if mutationinterval < 10 {

	}

}

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL THE SPECIFIC CREATURE FUNCTIONS -------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //

func (c *Creature1) updateTick() string {
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

func (c *Creature1) updateVal(val string) {
	if val == "heal" {
		addMessageToCurrentGamelog("Creature 1 with id "+
			strconv.Itoa(c.id)+" healed for: "+
			strconv.Itoa(c.oriHP), 2)
		c.hp += c.oriHP
		c.moving = false
	}
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

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL THE NECESSARY INTERFACE FUNCTIONS ------------ //
// -------------------------------------------------- //
// -------------------------------------------------- //

func (eo *EmptyObject) getSymbol() byte {
	return eo.symbol
}

func (f *Food) getSymbol() byte {
	return f.symbol
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
