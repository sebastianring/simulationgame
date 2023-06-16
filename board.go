package main

import (
	"fmt"
	"math/rand"
	"os"
	// "strconv"
)

var initialCreature1 int
var initialFoods int

var allFoodsObjects []Pos
var allCreatureObjects []Pos

type Board struct {
	rows int
	cols int
	// displayBoard [][]int // 0 = empty, 1 = food, 10-20 = creatures
	gamelog     *Gamelog
	objectBoard [][]BoardObject
	time        int
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

	newBoard := Board{
		rows,
		cols,
		// createBoardArray(rows, cols),
		InitTextInfo(rows),
		*createEmptyObjectsArray(rows, cols),
		0,
	}

	initialCreature1 = 20
	initialFoods = 50

	newBoard.spawnCreature1OnBoard(initialCreature1)
	newBoard.spawnFoodOnBoard(initialFoods)

	addMessageToCurrentGamelog("Board added", 2)
	addMessageToCurrentGamelog("Welcome to the simulation game where you can simulate creatures and how they evolve!", 1)

	return &newBoard
}

// No longer in use
func createBoardArray(rows int, cols int) [][]int {
	arr := make([][]int, rows)

	for i := 0; i < rows; i++ {
		arr[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			if rand.Intn(20) == 0 {
				arr[i][j] = 1
			} else {
				arr[i][j] = 0
			}
		}
	}

	return arr
}

// creates the initial array for all objects inside the board
func createObjectArray(rows int, cols int) *[][]BoardObject {
	arr := make([][]BoardObject, rows)
	edgeSpawnPoints := (rows*2 + cols*2 - 4)
	createSpawnChance := edgeSpawnPoints / 10 // 10% chance that a creature spawns at the edge

	for i := 0; i < rows; i++ {
		arr[i] = make([]BoardObject, cols)
		for j := 0; j < cols; j++ {
			// check if we are at the edge of the board, then roll the dice if a creature should be spawned
			if i == 0 || i == rows-1 || j == 0 || j == cols-1 {
				rng := rand.Intn(edgeSpawnPoints)
				if rng < createSpawnChance {
					arr[i][j] = newCreature1Object()
				} else {
					arr[i][j] = newEmptyObject()
				}
				// else, lets see if we can spawn some food, 2.5% chance to spawn
			} else {
				rng := rand.Intn(1000)
				if rng < 25 {
					arr[i][j] = newFoodObject()
				} else {
					arr[i][j] = newEmptyObject()
				}
			}
		}
	}

	return &arr
}

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

func (b *Board) spawnCreature1OnBoard(qty int) {
	spawns := make([][]int, 0)
	for len(spawns) < qty {
		newPos := b.randomPosAtEdgeOfMap()
		if !checkIfValExistsInSlice(newPos, spawns) {
			spawns = append(spawns, newPos)
		}
	}

	for _, val := range spawns {
		b.objectBoard[val[1]][val[0]] = newCreature1Object()
		allCreatureObjects = append(allCreatureObjects, Pos{x: val[0], y: val[1]})
	}
}

// Refactor with Pos struct
func (b *Board) spawnFoodOnBoard(qty int) {
	spawns := make([][]int, 0)

	for len(spawns) < qty {
		newPos := b.randomPosWithinMap()
		if !checkIfValExistsInSlice(newPos, spawns) && b.isSpotEmpty(newPos[0], newPos[1]) {
			spawns = append(spawns, newPos)
		}
	}

	for _, val := range spawns {
		b.objectBoard[val[1]][val[0]] = newFoodObject()
		allFoodsObjects = append(allFoodsObjects, Pos{x: val[0], y: val[1]})
	}
}

// Refactor with Pos struct
func (b *Board) isSpotEmpty(x int, y int) bool {
	if b.objectBoard[y][x].getType() == "empty" {
		return true
	}

	return false
}

// Refactor with Pos struct
func (b *Board) randomPosAtEdgeOfMap() []int {
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

	return []int{x, y}
}

// Refactor with Pos struct
func (b *Board) randomPosWithinMap() []int {
	minDistanceFromBorder := 3
	x := rand.Intn(b.cols-minDistanceFromBorder*2) + minDistanceFromBorder
	y := rand.Intn(b.rows-minDistanceFromBorder*2) + minDistanceFromBorder

	return []int{x, y}
}

func checkIfValExistsInSlice(val []int, slice [][]int) bool {
	for _, val2 := range slice {
		if len(val) == len(val2) {
			for i := 0; i < len(val); i++ {
				if val[i] == val2[i] {
					return false
				}
			}
		}
	}

	return false
}

func (b *Board) tickFrame() {
	b.time++
	// Update all the creatures on board
	for i, pos := range allCreatureObjects {
		action := b.objectBoard[pos.y][pos.x].updateTick()
		if action == "move" {
			// addMessageToCurrentGamelog("OLD POS: " + strconv.Itoa(pos.x) + " " + strconv.Itoa(pos.y))
			newPos, moveType := b.newPosAndMove(pos)
			tempObject := b.objectBoard[newPos.y][newPos.x]

			b.objectBoard[newPos.y][newPos.x] = b.objectBoard[pos.y][pos.x]
			if moveType == "food" {
				b.objectBoard[newPos.y][newPos.x].updateVal("heal")
				b.objectBoard[pos.y][pos.x] = newEmptyObject()
				deleteFood(newPos)
			} else {
				b.objectBoard[pos.y][pos.x] = tempObject
			}
			// addMessageToCurrentGamelog("Object moved to " + strconv.Itoa(newPos.x) + " " + strconv.Itoa(newPos.y))
			allCreatureObjects[i] = newPos

			// addMessageToCurrentGamelog("New POS: " + strconv.Itoa(pos.x) + " " + strconv.Itoa(pos.y))
		} else if action == "dead" {
			b.objectBoard[pos.y][pos.x] = newEmptyObject()
			deleteCreature(pos)
		}
	}

	if b.checkIfCreaturesAreInactive() == true {
		gameOn = false
		addMessageToCurrentGamelog("Game should end now", 2)
	}

	DrawFrame(b)
}

func (b *Board) checkIfCreaturesAreDead() bool {
	for _, pos := range allCreatureObjects {
		dead := b.objectBoard[pos.y][pos.x].isDead()
		// moving := b.objectBoard[pos.y][pos.x].isMoving()
		// addMessageToCurrentGamelog("DEAD:" + strconv.FormatBool(dead) + " MOVING: " + strconv.FormatBool(moving))

		if !dead {
			return false
		}
	}

	return true
}

func (b *Board) checkIfCreaturesAreInactive() bool {
	for _, pos := range allCreatureObjects {
		dead := b.objectBoard[pos.y][pos.x].isDead()
		moving := b.objectBoard[pos.y][pos.x].isMoving()

		// addMessageToCurrentGamelog("Current counter: " + strconv.Itoa(i) + "total length: " + strconv.Itoa(len(allCreatureObjects)))
		// addMessageToCurrentGamelog("DEAD:" + strconv.FormatBool(dead) + " MOVING: " + strconv.FormatBool(moving))

		if !dead && moving || dead {
			return false
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
		}
		return newPos, moveType
	}

	return newPos, "empty"
}

func (b *Board) checkIfNewPosIsValid(x int, y int) (bool, string) {
	if x < 0 || x >= b.cols || y < 0 || y >= b.rows {
		return false, ""
	}
	objectType := b.objectBoard[y][x].getType()
	if objectType == "food" {
		return true, "food"
	}

	return true, "empty"
}

func deleteFood(pos Pos) {
	var element int
	for i, val := range allFoodsObjects {
		if val.x == pos.x && val.y == pos.y {
			element = i
			break
		}
	}

	allFoodsObjects = deleteIndexInPosSlice(allFoodsObjects, element)
}

func deleteCreature(pos Pos) {
	var element int
	for i, val := range allCreatureObjects {
		if val.x == pos.x && val.y == pos.y {
			element = i
			break
		}
	}

	allCreatureObjects = deleteIndexInPosSlice(allCreatureObjects, element)
}

func deleteIndexInPosSlice(posSlice []Pos, index int) []Pos {
	posSlice[index] = posSlice[len(posSlice)-1]
	return posSlice[:len(posSlice)-1]
}
