package simulationgame

import (
	"math/rand"
)

type ConflictManager struct {
	ConflictMapping     [][]Conflict
	CreatureTranslation map[string]int
	ActionTranslation   map[Conflict]bool
}

type ConflictInfo struct {
	Conflict       Conflict
	SourceCreature CreatureObject
	TargetCreature CreatureObject
}

type Conflict int

const (
	Share   Conflict = 0
	Avoid   Conflict = 1
	Attack1 Conflict = 2
	Attack2 Conflict = 3
)

// ------ CURRENT STRATEGY MAPPINGS ---------
//       C1         C2
//    C1 SHARE      AVOID
//    C2 ATTACK1    ATTACK2
//
//    SHARE    = BOTH CREATURES GET HALF
//    AVOID    = SOURCE CREATURE AVOIDS TARGET CREATURE
//    ATTACK 1 = SOURCE CREATURE GETS ALL FOOD
//    ATTACK 2 = CREATURES FIGHT - WINNER TAKES 50% - OTHER DIES

func newConflictManager() (*ConflictManager, error) {
	cm := ConflictManager{
		ConflictMapping: [][]Conflict{
			{Share, Avoid},
			{Attack1, Attack2},
		},

		CreatureTranslation: map[string]int{
			"Creature1": 0,
			"Creature2": 1,
		},

		ActionTranslation: map[Conflict]bool{
			Share:   true,
			Avoid:   false,
			Attack1: true,
			Attack2: true,
		},
	}

	return &cm, nil
}

func (cm *ConflictManager) getConflict(SourceCreature CreatureObject, TargetCreature CreatureObject) (bool, *ConflictInfo) {
	// addMessageToCurrentGamelog("Conflict between two creatures checked", 1)
	row := cm.CreatureTranslation[SourceCreature.getType()]
	col := cm.CreatureTranslation[TargetCreature.getType()]

	conflictType := cm.ConflictMapping[row][col]

	action, ok := cm.ActionTranslation[conflictType]

	if !ok {
		addMessageToCurrentGamelog("Strategy between creatures is not mapped correctly.", 1)
	}

	ConflictInfo := ConflictInfo{
		Conflict:       conflictType,
		SourceCreature: SourceCreature,
		TargetCreature: TargetCreature,
	}

	return action, &ConflictInfo
}

func (cm *ConflictManager) share(SourceCreature CreatureObject, TargetCreature CreatureObject) {
	SourceCreature.heal(SourceCreature.getOriHP() / 2)
	TargetCreature.heal((TargetCreature.getOriHP() / 2) * -1)

	addMessageToCurrentGamelog(SourceCreature.getIdAsString()+" shared the food of "+TargetCreature.getIdAsString(), 1)
}

func (cm *ConflictManager) attack1(SourceCreature CreatureObject, TargetCreature CreatureObject) {
	SourceCreature.heal(SourceCreature.getOriHP())
	TargetCreature.kill()

	addMessageToCurrentGamelog(SourceCreature.getIdAsString()+" killed "+TargetCreature.getIdAsString()+" using attack1", 1)
}

func (cm *ConflictManager) attack2(SourceCreature CreatureObject, TargetCreature CreatureObject) bool {
	// function returns true if target is killed, if source is killed, it returns false
	rng := rand.Intn(2)
	if rng == 1 {
		SourceCreature.heal((SourceCreature.getOriHP() / 2) * -1)
		TargetCreature.kill()

		addMessageToCurrentGamelog(SourceCreature.getIdAsString()+" killed "+TargetCreature.getIdAsString()+" using attack2", 1)

		return true

	} else {
		SourceCreature.kill()
		TargetCreature.heal((TargetCreature.getOriHP() / 2) * -1)

		addMessageToCurrentGamelog(TargetCreature.getIdAsString()+" killed "+SourceCreature.getIdAsString()+" using attack2", 1)

		return false
	}
}

func (ci *ConflictInfo) commitConflict(b *Board) {
	switch ci.Conflict {
	case Share:
		ci.share(b)
	case Attack1:
		ci.attack1(b)
	case Attack2:
		ci.attack2(b)
	default:
		addMessageToCurrentGamelog("Conflict manager does not work properly, please have a look", 1)
	}
}

func (ci *ConflictInfo) share(b *Board) {
	ci.SourceCreature.heal(ci.SourceCreature.getOriHP() / 2)
	ci.TargetCreature.heal((ci.TargetCreature.getOriHP() / 2) * -1)

	addMessageToCurrentGamelog(ci.SourceCreature.getIdAsString()+" shared the food of "+ci.TargetCreature.getIdAsString(), 1)
}

func (ci *ConflictInfo) attack1(b *Board) {
	ci.SourceCreature.heal(ci.SourceCreature.getOriHP())

	b.killCreature(ci.TargetCreature, false)
	b.moveCreature(ci.SourceCreature, ci.TargetCreature.getPos(), true)

	addMessageToCurrentGamelog(ci.SourceCreature.getIdAsString()+" killed "+ci.TargetCreature.getIdAsString()+" using attack1", 1)
}

func (ci *ConflictInfo) attack2(b *Board) {
	rng := rand.Intn(2)
	if rng == 1 {
		ci.SourceCreature.heal((ci.SourceCreature.getOriHP() / 2) * -1)

		b.killCreature(ci.TargetCreature, false)
		b.moveCreature(ci.SourceCreature, ci.TargetCreature.getPos(), true)

		addMessageToCurrentGamelog(ci.SourceCreature.getIdAsString()+" killed "+ci.TargetCreature.getIdAsString()+" using attack2", 1)

	} else {
		ci.TargetCreature.heal((ci.TargetCreature.getOriHP() / 2) * -1)

		b.killCreature(ci.SourceCreature, false)
		b.moveCreature(ci.TargetCreature, ci.SourceCreature.getPos(), true)

		addMessageToCurrentGamelog(ci.TargetCreature.getIdAsString()+" killed "+ci.SourceCreature.getIdAsString()+" using attack2", 1)
	}
}
