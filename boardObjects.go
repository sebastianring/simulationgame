package simulationgame

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL INTERFACES AND GENERAL FUNCTIONS ------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //

type BoardObject interface {
	getSymbol() []byte
}

type CreatureObject interface {
	getSymbol() []byte
	updateTick() TickStatus
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
	getSpeed() int
	getPos() Pos
	setPos(Pos)
}

func getObjectSymbolWColor(objectname string) []byte {
	resetColor := []byte("\033[0m")

	colors := map[string][]byte{
		"green":   []byte("\033[32m"),
		"red":     []byte("\033[31m"),
		"blue":    []byte("\033[34m"),
		"yellow":  []byte("\033[33m"),
		"magenta": []byte("\033[35m"),
		"cyan":    []byte("\033[36m"),
		"white":   []byte("\033[37m"),
		"black":   []byte("\033[30m"),
	}

	drawingSymbols := map[string]byte{
		"EmptyObject": 46, // .....
		"Food":        64, // @@@@@
		"Creature1":   79, // OOOOO
		"Creature2":   87, // WWWWW
	}

	drawingColors := map[string]string{
		"EmptyObject": "black",
		"Food":        "green",
		"Creature1":   "cyan",
		"Creature2":   "red",
	}

	objectColor := drawingColors[objectname]
	returnByte := colors[objectColor]
	returnByte = append(returnByte, drawingSymbols[objectname])
	returnByte = append(returnByte, resetColor...)

	return returnByte
}

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL STRUCTS HERE AND THEIR CREATE FUNCTIONS ------ //
// -------------------------------------------------- //
// -------------------------------------------------- //

type EmptyObject struct {
	Symbol   []byte `json:"symbol"`
	TypeDesc string `json:"type_desc"`
}

func newEmptyObject() *EmptyObject {
	eo := EmptyObject{
		// symbol:   getObjectSymbol("EmptyObject"),
		Symbol:   getObjectSymbolWColor("EmptyObject"),
		TypeDesc: "EmptyObject",
	}

	// addMessageToCurrentGamelog("New empty object added", 2)

	return &eo
}

type Food struct {
	Symbol   []byte `json:"symbol"`
	TypeDesc string `json:"type_desc"`
}

func newFoodObject() *Food {
	f := Food{
		Symbol:   getObjectSymbolWColor("Food"),
		TypeDesc: "Food",
	}

	addMessageToCurrentGamelog("New food object added", 2)

	return &f
}

type TickStatus int

const (
	StatusMove  TickStatus = 0
	StatusDead  TickStatus = 1
	StatusError TickStatus = 2
)

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL THE NECESSARY INTERFACE FUNCTIONS ------------ //
// -------------------------------------------------- //
// -------------------------------------------------- //

func (eo *EmptyObject) getSymbol() []byte {
	return eo.Symbol
}

func (f *Food) getSymbol() []byte {
	return f.Symbol
}
