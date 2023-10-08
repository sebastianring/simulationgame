package simulationgame

import (
	"fmt"
	"os"
)

func (b *Board) spawnCreature1OnBoard(qty uint) {
	spawns := make([]Pos, 0)
	for uint(len(spawns)) < qty {
		newPos := b.randomPosAtEdgeOfMap()
		if !checkIfPosExistsInSlice(newPos, spawns) {
			spawns = append(spawns, newPos)
		}
	}

	for _, pos := range spawns {
		creature, err := b.newCreature1Object(false)

		if err != nil {
			fmt.Println("Error creating a new creature 1 object: " + err.Error())
			os.Exit(1)
		}

		creature.setPos(pos)

		b.ObjectBoard[pos.y][pos.x] = creature
		b.AliveCreatureObjects = append(b.AliveCreatureObjects, creature)
	}
}

func (b *Board) spawnCreature2OnBoard(qty uint) {
	spawns := make([]Pos, 0)
	for uint(len(spawns)) < qty {
		newPos := b.randomPosAtEdgeOfMap()
		if !checkIfPosExistsInSlice(newPos, spawns) && b.isSpotEmpty(newPos) {
			spawns = append(spawns, newPos)
		}
	}

	for _, pos := range spawns {
		creature, err := b.newCreature2Object(false)

		if err != nil {
			fmt.Println("Error creating a new creature 2 object: " + err.Error())
			os.Exit(1)
		}

		creature.setPos(pos)

		b.ObjectBoard[pos.y][pos.x] = creature
		b.AliveCreatureObjects = append(b.AliveCreatureObjects, creature)
	}
}
