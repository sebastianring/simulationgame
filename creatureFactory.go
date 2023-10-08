package simulationgame

import (
	"errors"
	"strconv"
)

type CreatureObjectType byte

const (
	creature1 CreatureObjectType = 0
	creature2 CreatureObjectType = 1
)

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
	getCreatureObjectType() CreatureObjectType
	getScanProcChance() float64
	scan()
}

func (b *Board) newCreatureObject(objectType CreatureObjectType, parent ...CreatureObject) (CreatureObject, error) {
	var creature CreatureObject

	var speed float64
	var procScanChance float64

	if len(parent) == 0 {
		speed = 5
		procScanChance = 50
	} else if len(parent) == 1 {
		tempspeed, err := b.MutationManager.getVariableValue(speedVariable, parent[0])

		if err != nil {
			addMessageToCurrentGamelog(err.Error(), 1)
			return nil, errors.New("Error creating new creature.")
		}

		speed = tempspeed

		tempProcChance, err := b.MutationManager.getVariableValue(scanVariable, parent[0])

		if err != nil {
			addMessageToCurrentGamelog(err.Error(), 1)
			return nil, errors.New("Error creating new creature.")
		}

		procScanChance = tempProcChance

	} else {
		return nil, errors.New("Too many parents")
	}

	switch objectType {
	case creature1:
		creature = &Creature1{
			Id:                 b.CreatureIdCtr[Creature1Type],
			Symbol:             getObjectSymbolWColor(Creature1Type),
			OriHP:              500,
			Hp:                 500,
			Speed:              speed,
			OriSpeed:           speed,
			ProcScanChance:     procScanChance,
			TypeDesc:           "Creature1",
			BoardObjectType:    Creature1Type,
			CreatureObjectType: creature1,
			Moving:             true,
		}

		b.CreatureIdCtr[Creature1Type] += 1

	case creature2:
		creature = &Creature2{
			Id:                 b.CreatureIdCtr[Creature2Type],
			Symbol:             getObjectSymbolWColor(Creature2Type),
			OriHP:              500,
			Hp:                 500,
			Speed:              speed,
			OriSpeed:           speed,
			ProcScanChance:     procScanChance,
			TypeDesc:           "Creature2",
			BoardObjectType:    Creature2Type,
			CreatureObjectType: creature2,
			Moving:             true,
		}

		b.CreatureIdCtr[Creature2Type] += 1

	default:
		return nil, errors.New("Trying to create an invalid creature")
	}

	addMessageToCurrentGamelog("Creature1 object with ID: "+
		strconv.Itoa(creature.getId())+" added to the board", 2)

	return creature, nil
}
