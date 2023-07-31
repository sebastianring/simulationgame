package main

// import "os"
import (
	"math/rand"
	// "strconv"
)

type conflictManager struct {
	conflictMapping     [][]string
	creatureTranslation map[string]int
	actionTranslation   map[string]bool
}

type conflictInfo struct {
	attack         string
	sourceCreature CreatureObject
	targetCreature CreatureObject
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
			"creature1": 0,
			"creature2": 1,
		},

		actionTranslation: map[string]bool{
			"share":   true,
			"avoid":   false,
			"attack1": true,
			"attack2": true,
		},
	}

	return &cm, nil
}

func (cm *conflictManager) getConflict(sourceCreature CreatureObject, targetCreature CreatureObject) (bool, *conflictInfo) {
	// addMessageToCurrentGamelog("Conflict between two creatures checked", 1)
	row := cm.creatureTranslation[sourceCreature.getType()]
	col := cm.creatureTranslation[targetCreature.getType()]

	strategy := cm.conflictMapping[row][col]

	action, ok := cm.actionTranslation[strategy]

	if !ok {
		addMessageToCurrentGamelog("Strategy between creatures is not mapped correctly.", 1)
	}

	conflictInfo := conflictInfo{
		attack:         strategy,
		sourceCreature: sourceCreature,
		targetCreature: targetCreature,
	}

	return action, &conflictInfo
}

func (cm *conflictManager) share(sourceCreature CreatureObject, targetCreature CreatureObject) {
	sourceCreature.heal(sourceCreature.getOriHP() / 2)
	targetCreature.heal((targetCreature.getOriHP() / 2) * -1)

	addMessageToCurrentGamelog(sourceCreature.getIdAsString()+" shared the food of "+targetCreature.getIdAsString(), 1)
}

func (cm *conflictManager) attack1(sourceCreature CreatureObject, targetCreature CreatureObject) {
	sourceCreature.heal(sourceCreature.getOriHP())
	targetCreature.kill()

	addMessageToCurrentGamelog(sourceCreature.getIdAsString()+" killed "+targetCreature.getIdAsString()+" using attack1", 1)
}

func (cm *conflictManager) attack2(sourceCreature CreatureObject, targetCreature CreatureObject) bool {
	// function returns true if target is killed, if source is killed, it returns false
	rng := rand.Intn(2)
	if rng == 1 {
		sourceCreature.heal((sourceCreature.getOriHP() / 2) * -1)
		targetCreature.kill()

		addMessageToCurrentGamelog(sourceCreature.getIdAsString()+" killed "+targetCreature.getIdAsString()+" using attack2", 1)

		return true

	} else {
		sourceCreature.kill()
		targetCreature.heal((targetCreature.getOriHP() / 2) * -1)

		addMessageToCurrentGamelog(targetCreature.getIdAsString()+" killed "+sourceCreature.getIdAsString()+" using attack2", 1)

		return false
	}
}
