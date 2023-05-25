package main

type BoardObject interface {
	getSymbol() byte
	getType() string
}

type CreatureObject interface {
	getHP() int
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
	speed    int
	typeDesc string
}

func newCreature1Object() *Creature1 {
	c1 := Creature1{
		symbol:   getObjectSymbol("Creature1"),
		active:   true,
		hp:       100,
		speed:    100,
		typeDesc: "creature",
	}

	addMessageToCurrentGamelog("New creature1 object added")

	return &c1
}

func (eo EmptyObject) getSymbol() byte {
	return eo.symbol
}

func (f Food) getSymbol() byte {
	return f.symbol
}

func (c Creature1) getSymbol() byte {
	return c.symbol
}

func (eo EmptyObject) getType() string {
	return eo.typeDesc
}

func (f Food) getType() string {
	return f.typeDesc
}

func (c Creature1) getType() string {
	return c.typeDesc
}

func (c Creature1) getHP() int {
	return c.hp
}
