package simulationgame

import "errors"

// "fmt"

type MutationManager struct {
	variables map[VariableToMutate]map[BoardObjectType]*MutationCalc
	// speedMutation map[BoardObjectType]*MutationCalc
	// scanMutation  map[BoardObjectType]*MutationCalc
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

func newMutationManager() *MutationManager {
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
		rate:             1.25,
	}

	scanMutationCreature2 := MutationCalc{
		chance:           50,
		mutationRateType: multiply,
		rate:             1.25,
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

	return &mm
}

func (mm *MutationManager) getMutatedValue(variable VariableToMutate, parent CreatureObject) (float64, error) {
	values, ok := mm.variables[variable][parent.getBoardObjectType()]

	if !ok {
		return 0, errors.New("Could not find the correct config for creature or variable in the mutation manager.")
	}

	var originalValue float64

	if variable == speedVariable {
		originalValue = parent.getSpeed()
	} else if variable == scanVariable {
		originalValue = parent.getScanProcChance()
	}

	if values.mutationRateType == multiply {
		returnValue := originalValue * values.rate
		return returnValue, nil
	} else if values.mutationRateType == flat {
		returnValue := originalValue + values.rate
		return returnValue, nil
	}

	return 0.0, errors.New("Missing confing for mutation rate")
}
