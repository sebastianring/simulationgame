package simulationgame

import (
	"errors"
	"math/rand"
	"strconv"
)

type MutationManager struct {
	variables map[VariableToMutate]map[BoardObjectType]*MutationCalc
}

type MutationCalc struct {
	chance           int // chance is X out of a 100 - so if this variable is 40, then it has a 40% chance of procing
	mutationRateType MutationRateType
	rate             float64
}

type VariableToMutate byte

const (
	speedVariable VariableToMutate = iota
	scanVariable
)

type MutationRateType byte

const (
	flat MutationRateType = iota
	multiply
)

func newMutationManager() (*MutationManager, error) {
	speedMutationCreature1 := MutationCalc{
		chance:           25,
		mutationRateType: flat,
		rate:             1,
	}

	speedMutationCreature2 := MutationCalc{
		chance:           25,
		mutationRateType: flat,
		rate:             1,
	}

	scanMutationCreature1 := MutationCalc{
		chance:           50,
		mutationRateType: multiply,
		rate:             0.10,
	}

	scanMutationCreature2 := MutationCalc{
		chance:           50,
		mutationRateType: multiply,
		rate:             0.10,
	}

	mm := MutationManager{
		variables: map[VariableToMutate]map[BoardObjectType]*MutationCalc{
			speedVariable: {
				Creature1Type: &speedMutationCreature1,
				Creature2Type: &speedMutationCreature2,
			},
			scanVariable: {
				Creature1Type: &scanMutationCreature1,
				Creature2Type: &scanMutationCreature2,
			},
		},
	}

	return &mm, nil
}

func (mm *MutationManager) getVariableValue(variable VariableToMutate, parent CreatureObject) (float64, error) {
	trigger, err := mm.rollMutationDice(variable, parent)
	var returnValue float64

	if err != nil {
		// addMessageToCurrentGamelog("Error when mutating creature speed value: "+err.Error(), 1)
		return returnValue, errors.New("Error when mutating creature " + strconv.Itoa(int(variable)) + " value: " + err.Error())
	}

	if trigger {
		returnValue, err = mm.getMutatedValue(speedVariable, parent)

		if err != nil {
			// addMessageToCurrentGamelog("Error when getting speed value "+err.Error(), 1)
			return returnValue, errors.New("Error when getting " + strconv.Itoa(int(variable)) + " value: " + err.Error())
		}
	} else {
		returnValue = parent.getSpeed()
	}

	return returnValue, nil
}

func (mm *MutationManager) rollMutationDice(variable VariableToMutate, parent CreatureObject) (bool, error) {
	values, ok := mm.variables[variable][parent.getBoardObjectType()]

	if !ok {
		return false, errors.New("Could not find the correct config for creature or variable in the mutation manager. Parent board type: " + strconv.Itoa(int(parent.getBoardObjectType())))
	}

	chance := rand.Intn(100)

	if chance < values.chance {
		return true, nil
	}

	return false, nil
}

func (mm *MutationManager) getMutatedValue(variable VariableToMutate, parent CreatureObject) (float64, error) {
	values, ok := mm.variables[variable][parent.getBoardObjectType()]

	// addMessageToCurrentGamelog("Parents speed: "+strconv.FormatFloat(parent.getSpeed(), 'f', 2, 64), 1)

	if !ok {
		return 0, errors.New("Could not find the correct config for creature or variable in the mutation manager.")
	}

	var originalValue float64

	if variable == speedVariable {
		originalValue = parent.getSpeed()
	} else if variable == scanVariable {
		originalValue = parent.getScanProcChance()
	}

	negativeOrPositive := rand.Intn(2)
	var adjustor float64
	if negativeOrPositive == 0 {
		adjustor = -1.00
	} else {
		adjustor = 1.00
	}

	if values.mutationRateType == multiply {
		returnValue := originalValue + ((originalValue * values.rate) * adjustor)
		return returnValue, nil
	} else if values.mutationRateType == flat {
		returnValue := originalValue + (values.rate * adjustor)
		return returnValue, nil
	}

	return 0.0, errors.New("Missing confing for mutation rate")
}
