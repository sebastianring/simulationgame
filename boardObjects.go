package main

// import "strconv"

type BoardObject interface {
	getSymbol() byte
	getType() string
	updateTick() string
	updateVal(string)
	getIntData(string) int
	isDead() bool
	isMoving() bool
}

type CreatureObject interface {
	getHP() int
	updateTick() string
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

type EmptyObject struct {
	symbol   byte
	typeDesc string
}

func newEmptyObject() *EmptyObject {
	eo := EmptyObject{
		symbol:   getObjectSymbol("EmptyObject"),
		typeDesc: "empty",
	}

	addMessageToCurrentGamelog("New empty object added")

	return &eo
}

type Food struct {
	symbol   byte
	active   bool
	typeDesc string
}

func newFoodObject() *Food {
	f := Food{
		symbol:   getObjectSymbol("Food"),
		typeDesc: "food",
	}

	addMessageToCurrentGamelog("New food object added")

	return &f
}

type Creature1 struct {
	symbol   byte
	active   bool
	hp       int
	oriHP    int
	speed    int
	oriSpeed int
	typeDesc string
	moving   bool
}

func newCreature1Object() *Creature1 {
	c1 := Creature1{
		symbol:   getObjectSymbol("Creature1"),
		active:   true,
		oriHP:    100,
		hp:       100,
		speed:    15,
		oriSpeed: 15,
		typeDesc: "creature",
		moving:   true,
	}

	addMessageToCurrentGamelog("New creature1 object added")

	return &c1
}

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

func (c *Creature1) getType() string {
	return c.typeDesc
}

func (c *Creature1) getHP() (int, bool) {
	return c.hp, c.moving
}

func (c *Creature1) getSpeed() int {
	return c.speed
}

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

	return ""
}

func (c *Creature1) getSymbol() byte {
	return c.symbol
}

func (c *Creature1) updateVal(val string) {
	if val == "heal" {
		addMessageToCurrentGamelog("FOOD EATEN")
		c.hp += c.oriHP
		c.moving = false
	}
}

func (c *Creature1) getIntData(data string) int {
	if data == "hp" {
		return c.hp
	}

	return 0
}

func (c *Creature1) isDead() bool {
	if c.hp <= 0 {
		return true
	}

	return false
}

func (c *Creature1) isMoving() bool {
	return c.moving
}
