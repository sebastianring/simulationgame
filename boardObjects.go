package simulationgame

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL INTERFACES AND GENERAL FUNCTIONS ------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //

type BoardObject interface {
	getSymbol() []byte
	getBoardObjectType() BoardObjectType
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
	getSpeed() float64
	getPos() Pos
	setPos(Pos)
	getBoardObjectType() BoardObjectType
	getScanProcChance() float64
	scan()
}

func getObjectSymbolWColor(ObjectType BoardObjectType) []byte {
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

	drawingSymbols := map[BoardObjectType]byte{
		EmptyType:     46, // SPACE
		FoodType:      64, // @@@@@
		Creature1Type: 79, // OOOOO
		Creature2Type: 87, // WWWWW
	}

	drawingColors := map[BoardObjectType]string{
		EmptyType:     "black",
		FoodType:      "green",
		Creature1Type: "cyan",
		Creature2Type: "red",
	}

	objectColor := drawingColors[ObjectType]
	returnByte := colors[objectColor]
	returnByte = append(returnByte, drawingSymbols[ObjectType])
	returnByte = append(returnByte, resetColor...)

	return returnByte
}

// -------------------------------------------------- //
// -------------------------------------------------- //
// ALL STRUCTS HERE AND THEIR CREATE FUNCTIONS ------ //
// -------------------------------------------------- //
// -------------------------------------------------- //

type EmptyObject struct {
	Symbol          []byte          `json:"symbol"`
	TypeDesc        string          `json:"type_desc"`
	BoardObjectType BoardObjectType `json:"board_object_type"`
}

func newEmptyObject() *EmptyObject {
	eo := EmptyObject{
		Symbol:   getObjectSymbolWColor(EmptyType),
		TypeDesc: "EmptyObject",
	}

	// addMessageToCurrentGamelog("New empty object added", 2)

	return &eo
}

type Food struct {
	Symbol          []byte          `json:"symbol"`
	TypeDesc        string          `json:"type_desc"`
	BoardObjectType BoardObjectType `json:"board_object_type"`
}

func newFoodObject() *Food {
	f := Food{
		Symbol:   getObjectSymbolWColor(FoodType),
		TypeDesc: "Food",
	}

	addMessageToCurrentGamelog("New food object added", 2)

	return &f
}

type BoardObjectType byte
type TickStatus byte

const (
	StatusMove     TickStatus      = 0
	StatusDead     TickStatus      = 1
	StatusInactive TickStatus      = 2
	Creature1Type  BoardObjectType = 1
	Creature2Type  BoardObjectType = 2
	EmptyType      BoardObjectType = 100
	FoodType       BoardObjectType = 101
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

func (c *Food) getBoardObjectType() BoardObjectType {
	return c.BoardObjectType
}

func (c *EmptyObject) getBoardObjectType() BoardObjectType {
	return c.BoardObjectType
}
