package main

import (
	"errors"
	"strconv"
)

// import "strconv"

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL INTERFACES AND GENERAL FUNCTIONS ------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //
// I really need to change architecture of the board .. this is abuse of interfaces. //

type BoardObject interface {
	getSymbol() byte
	getType() string
	updateTick() string
	updateVal(string)
	getIntData(string) int
	isDead() bool
	isMoving() bool
	resetValues()
	ifOffspring() bool
}

type CreatureObject interface {
	getHP() int
	updateTick() string
	ifOffspring() bool
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

type EmptyObject struct {
	symbol   byte
	typeDesc string
}

func newEmptyObject() *EmptyObject {
	eo := EmptyObject{
		symbol:   getObjectSymbol("EmptyObject"),
		typeDesc: "empty",
	}

	addMessageToCurrentGamelog("New empty object added", 2)

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

var Creature1IdCtr int

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

func newCreature1Object(parent ...*Creature1) (*Creature1, error) {
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

	if Creature1IdCtr < 1 {
		Creature1IdCtr = 1
	}

	c1 := Creature1{
		id:       Creature1IdCtr,
		symbol:   getObjectSymbol("Creature1"),
		oriHP:    250,
		hp:       250,
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
			c.hp -= 10
			return "move"
		}
	} else if c.hp <= 0 {
		return "dead"
	}

	return "error"
}

func (c *Creature1) getSymbol() byte {
	return c.symbol
}

func (c *Creature1) updateVal(val string) {
	if val == "heal" {
		addMessageToCurrentGamelog("", 2)
		c.hp += c.oriHP
		c.moving = false
	}
}

func (c *Creature1) getIntData(data string) int {
	if data == "hp" {
		return c.hp
	} else if data == "speed" {
		return c.speed
	} else if data == "id" {
		return c.id
	}

	return 0
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

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL THE NECESSARY INTERFACE FUNCTIONS ------------ //
// -------------------------------------------------- //
// -------------------------------------------------- //

func (eo *EmptyObject) getSymbol() byte {
	return eo.symbol
}

func (eo *EmptyObject) updateTick() string {
	return ""
}

func (eo *EmptyObject) getData() map[string]int {
	returnMap := make(map[string]int, 0)

	return returnMap
}

func (eo *EmptyObject) updateVal(val string) {

}

func (eo *EmptyObject) getType() string {
	return eo.typeDesc
}

func (eo *EmptyObject) getIntData(data string) int {
	return 0
}

func (eo *EmptyObject) isDead() bool {
	return false
}

func (eo *EmptyObject) isMoving() bool {
	return false
}

func (eo *EmptyObject) resetValues() {

}

func (eo *EmptyObject) ifOffspring() bool {
	return false
}

func (f *Food) getSymbol() byte {
	return f.symbol
}

func (f *Food) getType() string {
	return f.typeDesc
}

func (f *Food) updateTick() string {
	return ""
}

func (f *Food) updateVal(val string) {

}

func (f *Food) getIntData(data string) int {
	return 0
}

func (f *Food) isDead() bool {
	return false
}

func (f *Food) isMoving() bool {
	return false
}

func (f *Food) resetValues() {

}

func (f *Food) ifOffspring() bool {
	return false
}

func (c *Creature1) getType() string {
	return c.typeDesc
}

func (c *Creature1) getHP() (int, bool) {
	return c.hp, c.moving
}

func (c *Creature1) ifOffspring() bool {
	if c.hp > int(float32(c.oriHP)*1.1) {
		return true
	}

	return false
}

func (c *Creature1) getSpeed() int {
	return c.speed
}

func (c *Creature1) isMoving() bool {
	return c.moving
}
