package main

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL INTERFACES AND GENERAL FUNCTIONS ------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //
// I really need to change architecture of the board .. this is abuse of interfaces. //

type BoardObject interface {
	getSymbol() byte
}

type CreatureObject interface {
	getSymbol() byte
	updateTick() string
	ifOffspring() bool
	getHP() int
	getId() int
	resetValues()
	heal(int)
	isMoving() bool
	isDead() bool
	getType() string
	kill()
	getOriHP() int
	getIdAsString() string
}

func getObjectSymbol(objectname string) byte {
	drawingSymbols := map[string]byte{
		"EmptyObject": 46, // .....
		"Food":        64, // @@@@@
		"Creature1":   79, // OOOOO
		"Creature2":   87, // WWWWW
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
