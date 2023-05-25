package main

import (
	"fmt"
	"math/rand"
	"os"
)

var initialCreature1 int
var initialFoods int

var allFoodsObjects []*Food
var allCreatureObjects []*CreatureObject

type Board struct {
	rows         int
	cols         int
	displayBoard [][]int // 0 = empty, 1 = food, 10-20 = creatures
	gamelog      *Gamelog
	objectBoard  [][]BoardObject
}

func InitNewBoard(rows int, cols int) *Board {
	if rows < 2 || cols < 2 {
		fmt.Printf("Too few rows or cols: %v, rows: %v \n", rows, cols)
		os.Exit(1)
	}

	newBoard := Board{
		rows,
		cols,
		createBoardArray(rows, cols),
		InitTextInfo(rows),
		*createEmptyObjectsArray(rows, cols),
		// createObjectArray(rows, cols),
	}

	initialCreature1 = 20
	initialFoods = 50

	newBoard.spawnCreature1OnBoard(initialCreature1)
	newBoard.spawnFoodOnBoard(initialFoods)

	addMessageToCurrentGamelog("Board added")
	addMessageToCurrentGamelog("Welcome to the simulation game where you can simulate creatures and how they evolve!")

	return &newBoard
}

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

	fmt.Println(spawns)

	for _, val := range spawns {
		b.objectBoard[val[1]][val[0]] = newCreature1Object()
	}
}

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
	}
}

func (b *Board) isSpotEmpty(x int, y int) bool {
	if b.objectBoard[y][x].getType() == "empty" {
		return true
	}

	return false
}

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
