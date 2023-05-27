package main

// import "strconv"

type BoardObject interface {
	getSymbol() byte
	getType() string
	updateTick() string
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

func (f *Food) getSymbol() byte {
	return f.symbol
}

func (c *Creature1) getSymbol() byte {
	return c.symbol
}

func (eo *EmptyObject) getType() string {
	return eo.typeDesc
}

func (f *Food) getType() string {
	return f.typeDesc
}

func (f *Food) updateTick() string {
	return ""
}

func (c *Creature1) getType() string {
	return c.typeDesc
}

func (c *Creature1) getHP() int {
	return c.hp
}

func (c *Creature1) getSpeed() int {
	return c.speed
}

func (c *Creature1) updateTick() string {
	c.speed -= 1
	// addMessageToCurrentGamelog(strconv.Itoa(c.speed))
	if c.speed == 0 {
		// addMessageToCurrentGamelog("Should move now ...")
		c.speed = c.oriSpeed
		return "move"
	}

	return ""
}
