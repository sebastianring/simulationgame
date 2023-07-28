package main

import "math/rand"

type conflictManager struct {
	conflictMapping     [][]string
	creatureTranslation map[string]int
}

// ------ CURRENT STRATEGY MAPPINGS ---------
//       C1         C2
//    C1 SHARE      AVOID
//    C2 ATTACK1    ATTACK2
//
//    SHARE    = BOTH CREATURES GET HALF
//    AVOID    = SOURCE CREATURE AVOIDS TARGET CREATURE
//    ATTACK 1 = SOURCE CREATURE GETS ALL FOOD
//    ATTACK 2 = CREATURES FIGHT - WINNER TAKES 50% - OTHER DIES

func newConflictManager() (*conflictManager, error) {
	cm := conflictManager{
		conflictMapping: [][]string{
			{"share", "avoid"},
			{"attack1", "attack2"},
		},

		creatureTranslation: map[string]int{
			"creature1": 1,
			"creature2": 2,
		},
	}

	return &cm, nil
}

func (cm *conflictManager) getConflict(sourceCreature CreatureObject, targetCreature CreatureObject) bool {
	row := cm.creatureTranslation[sourceCreature.getType()]
	col := cm.creatureTranslation[targetCreature.getType()]

	strategy := cm.conflictMapping[row][col]

	switch strategy {
	case "avoid":
		return false

	case "share":
		sourceCreature.heal(sourceCreature.getOriHP() / 2)
		targetCreature.heal((targetCreature.getOriHP() / 2) * -1)
		return true

	case "attack1":
		sourceCreature.heal(sourceCreature.getOriHP())
		targetCreature.kill()
		return true

	case "attack2":
		rng := rand.Intn(2)
		if rng == 1 {
			sourceCreature.heal((sourceCreature.getOriHP() / 2) * -1)
			targetCreature.kill()
		} else {
			sourceCreature.kill()
			targetCreature.heal((targetCreature.getOriHP() / 2) * -1)
		}

	default:
		return false
	}

	return false
}
