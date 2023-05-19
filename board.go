package main

import (
	"math/rand"
)

type Board struct {
	rows         int
	cols         int
	displayBoard [][]int // 0 = empty, 1 = food, 10-20 = creatures
	gamelog      *Gamelog
	objectBoard  [][]BoardObject
}

func InitNewBoard(rows int, cols int) *Board {
	newBoard := Board{
		rows,
		cols,
		createBoardArray(rows, cols),
		InitTextInfo(rows),
		createObjectArray(rows, cols),
	}

	newBoard.gamelog.addMessage("Board added")
	newBoard.gamelog.addMessage("Welcome to the simulation game where you can simulate creatures and how they evolve!")

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

func createObjectArray(rows int, cols int) [][]BoardObject {
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

	return arr
}
