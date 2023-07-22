package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	// "strings"
)

// -------------------------------------------------- //
// ------------------------------------------------- //
// INITS AND STRUCTS -------------------------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //

var allFoodObjects []Pos
var allAliveCreatureObjects []Pos
var allDeadCreatures []*BoardObject

type Board struct {
	rows          int
	cols          int
	gamelog       *Gamelog
	objectBoard   [][]BoardObject
	time          int
	roundInt      int
	rounds        []*Round
	currentRound  *Round
	creatureIdCtr map[string]int
	mutationrate  map[string]float32
	initialFoods  int
}

type Round struct {
	id               int
	time             int
	creaturedSpawned []CreatureObject
	creaturedKilled  []CreatureObject
}

type Pos struct {
	x int
	y int
}

func InitNewBoard(rows int, cols int) *Board {
	if rows < 2 || cols < 2 {
		fmt.Printf("Too few rows or cols: %v, rows: %v \n", rows, cols)
		os.Exit(1)
	}

	newRound := Round{
		id:               1,
		time:             0,
		creaturedSpawned: make([]CreatureObject, 0),
		creaturedKilled:  make([]CreatureObject, 0),
	}

	newBoard := Board{
		rows:          rows,
		cols:          cols,
		gamelog:       InitTextInfo(rows),
		objectBoard:   *createEmptyObjectsArray(rows, cols),
		time:          0,
		roundInt:      1,
		rounds:        []*Round{&newRound},
		currentRound:  &newRound,
		creatureIdCtr: make(map[string]int, 0),
		mutationrate:  make(map[string]float32, 0),
		initialFoods:  100,
	}

	newBoard.initBoardObjects()

	initialCreature1 := 10
	initialCreature2 := 10

	newBoard.spawnCreature1OnBoard(initialCreature1)
	newBoard.spawnCreature2OnBoard(initialCreature2)
	newBoard.spawnFoodOnBoard()

	addMessageToCurrentGamelog("Board added", 2)
	addMessageToCurrentGamelog("Welcome to the simulation game where you can simulate creatures and how they evolve!", 1)

	return &newBoard
}

// creates the initial array for all objects inside the board
// func createObjectArray(rows int, cols int) *[][]BoardObject {
// 	arr := make([][]BoardObject, rows)
// 	edgeSpawnPoints := (rows*2 + cols*2 - 4)
// 	createSpawnChance := edgeSpawnPoints / 10 // 10% chance that a creature spawns at the edge
//
// 	for i := 0; i < rows; i++ {
// 		arr[i] = make([]BoardObject, cols)
// 		for j := 0; j < cols; j++ {
// 			// check if we are at the edge of the board, then roll the dice if a creature should be spawned
// 			if i == 0 || i == rows-1 || j == 0 || j == cols-1 {
// 				rng := rand.Intn(edgeSpawnPoints)
// 				if rng < createSpawnChance {
// 					creaturePtr, err := newCreature1Object(false)
//
// 					if err != nil {
// 						fmt.Println("Issue when creating new creature 1 object: " + err.Error())
// 						os.Exit(1)
// 					}
//
// 					arr[i][j] = creaturePtr
//
// 				} else {
// 					arr[i][j] = newEmptyObject()
// 				}
// 				// else, lets see if we can spawn some food, 2.5% chance to spawn
// 			} else {
// 				rng := rand.Intn(1000)
// 				if rng < 25 {
// 					arr[i][j] = newFoodObject()
// 				} else {
// 					arr[i][j] = newEmptyObject()
// 				}
// 			}
// 		}
// 	}
//
// 	return &arr
// }

func createEmptyObjectsArray(rows int, cols int) *[][]BoardObject {
	arr := make([][]BoardObject, rows)

	for i := 0; i < rows; i++ {
		arr[i] = make([]BoardObject, cols)
		for j := 0; j < cols; j++ {
			arr[i][j] = newEmptyObject()
		}
	}

	return &arr
}

// -------------------------------------------------- //
// -------------------------------------------------- //
// BOARD FUNCTIONS ---------------------------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //

func (b *Board) spawnCreature1OnBoard(qty int) {
	spawns := make([]Pos, 0)
	for len(spawns) < qty {
		newPos := b.randomPosAtEdgeOfMap()
		if !checkIfPosExistsInSlice(newPos, spawns) {
			spawns = append(spawns, newPos)
		}
	}

	for _, pos := range spawns {
		creaturePtr, err := b.newCreature1Object(false)

		if err != nil {
			fmt.Println("Error creating a new creature 1 object: " + err.Error())
			os.Exit(1)
		}

		b.objectBoard[pos.y][pos.x] = creaturePtr
		allAliveCreatureObjects = append(allAliveCreatureObjects, pos)
	}
}

func (b *Board) spawnCreature2OnBoard(qty int) {
	spawns := make([]Pos, 0)
	for len(spawns) < qty {
		newPos := b.randomPosAtEdgeOfMap()
		if !checkIfPosExistsInSlice(newPos, spawns) && b.isSpotEmpty(newPos) {
			spawns = append(spawns, newPos)
		}
	}

	for _, pos := range spawns {
		creaturePtr, err := b.newCreature2Object(false)

		if err != nil {
			fmt.Println("Error creating a new creature 2 object: " + err.Error())
			os.Exit(1)
		}

		b.objectBoard[pos.y][pos.x] = creaturePtr
		allAliveCreatureObjects = append(allAliveCreatureObjects, pos)
	}
}

func (b *Board) spawnFoodOnBoard() {
	qty := b.initialFoods

	spawns := make([]Pos, 0)

	for len(spawns) < qty {
		newPos := b.randomPosWithinMap()
		if !checkIfPosExistsInSlice(newPos, spawns) && b.isSpotEmpty(newPos) {
			spawns = append(spawns, newPos)
		}
	}

	for _, pos := range spawns {
		b.objectBoard[pos.y][pos.x] = newFoodObject()
		allFoodObjects = append(allFoodObjects, pos)
	}
}

func (b *Board) isSpotEmpty(pos Pos) bool {
	if _, ok := b.objectBoard[pos.y][pos.x].(*EmptyObject); ok {
		return true
	}

	return false
}

func (b *Board) randomPosAtEdgeOfMap() Pos {
	// top = 0, right = 1, left = 2, bottom = 3
	edge := rand.Intn(4)
	var x int
	var y int

	if edge == 0 {
		y = 0
		x = rand.Intn(b.cols - 1)
	} else if edge == 1 {
		x = b.cols - 1
		y = rand.Intn(b.rows - 1)
	} else if edge == 2 {
		x = 0
		y = rand.Intn(b.rows - 1)
	} else {
		x = rand.Intn(b.cols - 1)
		y = b.rows - 1
	}

	return Pos{x, y}
}

func (b *Board) initBoardObjects() {
	fmt.Println("does this work?")
	b.creatureIdCtr["creature1"] = 1
	b.creatureIdCtr["creature2"] = 1

	fmt.Println("how about this?")
	b.mutationrate = make(map[string]float32)
	b.mutationrate["creature1"] = 0.1
	b.mutationrate["creature2"] = 0.1
}

func (b *Board) randomPosWithinMap() Pos {
	minDistanceFromBorder := 3
	x := rand.Intn(b.cols-minDistanceFromBorder*2) + minDistanceFromBorder
	y := rand.Intn(b.rows-minDistanceFromBorder*2) + minDistanceFromBorder

	return Pos{x, y}
}

func checkIfPosExistsInSlice(pos Pos, slice []Pos) bool {
	for _, slicePos := range slice {
		if pos.x == slicePos.x && pos.y == slicePos.y {
			return true
		}
	}

	return false
}

func (b *Board) tickFrame() {
	b.currentRound.time++
	b.creatureUpdatesPerTick()

	// ----------- debugging creatures - print speed and id ---------- //

	res := make([]string, 1)

	for _, pos := range allAliveCreatureObjects {
		if obj, ok := b.objectBoard[pos.y][pos.x].(CreatureObject); ok {
			res = append(res, strconv.Itoa(obj.getId())+":"+strconv.Itoa(obj.getId()))
		}
	}

	// addMessageToCurrentGamelog(strings.Join(res, ", "), 2)

	// ----------- end debugging ------------------------------------ //

	DrawFrame(b)
}

func checkCreatureType(bo BoardObject) (bool, *BoardObject) {
	switch bo.(type) {
	case *Creature1:
		return true, &bo
	case *Creature2:
		return true, &bo
	default:
		return false, nil
	}
}

func (b *Board) creatureUpdatesPerTick() {
	updatedAllCreatureObjects := make([]Pos, 0)
	deadCreatures := make([]Pos, 0)

	for _, pos := range allAliveCreatureObjects {
		if obj, ok := b.objectBoard[pos.y][pos.x].(CreatureObject); ok {
			// ifCreature, obj := checkCreatureType(b.objectBoard[pos.y][pos.x])
			// addMessageToCurrentGamelog(strconv.Itoa(obj.id)+" "+strconv.Itoa(i)+" "+strconv.Itoa(pos.x)+" "+strconv.Itoa(pos.y), 2)
			action := obj.updateTick()

			if action == "move" {
				newPos, moveType := b.newPosAndMove(pos)
				b.objectBoard[newPos.y][newPos.x] = BoardObject(obj)

				if moveType == "food" {
					addMessageToCurrentGamelog("Food eaten by creature id: "+strconv.Itoa(obj.getId()), 2)
					obj.updateVal("heal")
					b.objectBoard[pos.y][pos.x] = newEmptyObject()
					deleteFood(newPos)
				} else {
					b.objectBoard[pos.y][pos.x] = newEmptyObject()
				}

				updatedAllCreatureObjects = append(updatedAllCreatureObjects, newPos)

			} else if action == "dead" {
				deadCreatures = append(deadCreatures, pos)

			} else {
				updatedAllCreatureObjects = append(updatedAllCreatureObjects, pos)
			}
		}
	}

	// delete dead creatures after tick is complete
	for _, pos := range deadCreatures {
		deleteCreature(pos, &b.objectBoard[pos.y][pos.x])
		b.objectBoard[pos.y][pos.x] = newEmptyObject()
	}

	// update all creatures from last tick
	allAliveCreatureObjects = updatedAllCreatureObjects

	if b.checkIfCreaturesAreInactive() {
		if b.checkIfCreaturesAreDead() {
			gameOn = false
			addMessageToCurrentGamelog("All creatures are dead, end the game", 1)
		}

		b.newRound()
	}
}

func (b *Board) newRound() {
	addMessageToCurrentGamelog("All creatures inactive, starting new round", 1)
	b.spawnOffsprings()
	b.findPosForAllCreatures()
	b.deleteAndSpawnFood()

	newRound := Round{
		id:               b.currentRound.id + 1,
		time:             0,
		creaturedSpawned: make([]CreatureObject, 0),
		creaturedKilled:  make([]CreatureObject, 0),
	}

	b.currentRound = &newRound
	b.rounds = append(b.rounds, &newRound)

	b.gamelog.writeGamelogToFile()
}

func (b *Board) deleteAndSpawnFood() {
	for _, pos := range allFoodObjects {
		b.objectBoard[pos.y][pos.x] = newEmptyObject()
	}

	allFoodObjects = make([]Pos, 0)
	b.spawnFoodOnBoard()
}

func (b *Board) findPosForAllCreatures() {
	for i, creaturePos := range allAliveCreatureObjects {
		if obj, ok := b.objectBoard[creaturePos.y][creaturePos.x].(CreatureObject); ok {
			findNewPos := false
			for !findNewPos {
				newPos := b.randomPosAtEdgeOfMap()
				if b.isSpotEmpty(newPos) {
					b.objectBoard[newPos.y][newPos.x] = obj
					obj.resetValues()
					b.objectBoard[creaturePos.y][creaturePos.x] = newEmptyObject()
					allAliveCreatureObjects[i] = newPos
					findNewPos = true
				}
			}

		}
	}
}

func (b *Board) spawnOffsprings() {
	creatureQty := map[string]int{
		"creature1": 0,
		"creature2": 0,
	}

	for _, pos := range allAliveCreatureObjects {
		if obj, ok := b.objectBoard[pos.y][pos.x].(CreatureObject); ok {
			if obj.ifOffspring() {
				var offspring CreatureObject
				var err error

				if obj2, ok := obj.(*Creature1); ok {
					offspring, err = b.newCreature1Object(true, obj2)

					if err != nil {
						fmt.Println("Error creating offspring: " + err.Error())
					}

					creatureQty["creature1"]++

				} else if obj2, ok := obj.(*Creature2); ok {
					offspring, err = b.newCreature2Object(true, obj2)
					if err != nil {
						fmt.Println("Error creating offspring: " + err.Error())
					}

					creatureQty["creature2"]++
				}

				b.currentRound.creaturedSpawned = append(b.currentRound.creaturedSpawned, offspring)

				newPos := b.randomPosAtEdgeOfMap()
				for !b.isSpotEmpty(newPos) {
					newPos = b.randomPosAtEdgeOfMap()
				}

				b.objectBoard[newPos.y][newPos.x] = offspring
				allAliveCreatureObjects = append(allAliveCreatureObjects, newPos)

				// qty++
			}
		}
	}

	sum := 0
	for _, val := range creatureQty {
		sum = +val
	}

	addMessageToCurrentGamelog(strconv.Itoa(sum)+" new creatures spawned", 1)

	b.spawnCreature1OnBoard(creatureQty["creature1"])
	b.spawnCreature2OnBoard(creatureQty["creature2"])
}

func (b *Board) checkIfCreaturesAreDead() bool {
	for _, pos := range allAliveCreatureObjects {
		if obj, ok := b.objectBoard[pos.y][pos.x].(CreatureObject); ok {
			dead := obj.isDead()

			if !dead {
				return false
			}
		}
	}

	return true
}

func (b *Board) checkIfCreaturesAreInactive() bool {
	for _, pos := range allAliveCreatureObjects {
		if obj, ok := b.objectBoard[pos.y][pos.x].(CreatureObject); ok {
			dead := obj.isDead()
			moving := obj.isMoving()

			// addMessageToCurrentGamelog("Current counter: " + strconv.Itoa(i) + "total length: " + strconv.Itoa(len(allAliveCreatureObjects)))
			// addMessageToCurrentGamelog("DEAD:" + strconv.FormatBool(dead) + " MOVING: " + strconv.FormatBool(moving))

			if !dead && moving || dead {
				return false
			}
		}

	}

	return true
}

func (b *Board) newPosAndMove(currentPos Pos) (Pos, string) {
	newPos := Pos{-1, -1}

	// HOW TO MAKE THE CREATURES MOVE INWARDS TO LOOK FOR FOOD?
	// The closer they are to one edge, the more probable they are to move towards the other edge?
	// Example: x = 99, y = 40
	// Width-x = the probability to move the left
	// Height-y = the probability to move upwards?

	for newPos.x == -1 || newPos.y == -1 {
		direction := rand.Intn(2) // 0 = x movement, 1 = y-movement
		var x int
		var y int
		if direction == 0 {
			xdirection := rand.Intn(b.cols)
			xprobability := b.cols - 1 - currentPos.x
			if xdirection < xprobability {
				x = currentPos.x + 1
			} else {
				x = currentPos.x - 1
			}
			y = currentPos.y
		} else {
			ydirection := rand.Intn(b.rows)
			yprobability := b.rows - 1 - currentPos.y
			if ydirection < yprobability {
				y = currentPos.y + 1
			} else {
				y = currentPos.y - 1
			}
			x = currentPos.x
		}

		valid, moveType := b.checkIfNewPosIsValid(x, y)

		if valid {
			newPos.x = x
			newPos.y = y
			return newPos, moveType
		}
	}

	return newPos, "empty"
}

func (b *Board) checkIfNewPosIsValid(x int, y int) (bool, string) {
	if x < 0 || x >= b.cols || y < 0 || y >= b.rows {
		return false, ""
	}

	if _, ok := b.objectBoard[y][x].(*EmptyObject); ok {
		return true, "empty"
	} else if _, ok := b.objectBoard[y][x].(*Food); ok {
		return true, "food"
	}

	return false, ""
}

func deleteFood(pos Pos) {
	var element int
	for i, val := range allFoodObjects {
		if val.x == pos.x && val.y == pos.y {
			element = i
			break
		}
	}

	allFoodObjects = deleteIndexInPosSlice(allFoodObjects, element)
}

func deleteCreature(pos Pos, creature *BoardObject) {
	allDeadCreatures = append(allDeadCreatures, creature)
	var element int
	for i, val := range allAliveCreatureObjects {
		if val.x == pos.x && val.y == pos.y {
			element = i
			break
		}
	}

	allAliveCreatureObjects = deleteIndexInPosSlice(allAliveCreatureObjects, element)
}

func deleteIndexInPosSlice(posSlice []Pos, index int) []Pos {
	posSlice[index] = posSlice[len(posSlice)-1]
	return posSlice[:len(posSlice)-1]
}
