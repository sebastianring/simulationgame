package simulationgame

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var gameOn bool

func main() {
	fmt.Println("..")
}

type SimulationConfig struct {
	Rows      int
	Cols      int
	Draw      bool
	Foods     int
	Creature1 uint
	Creature2 uint
}

func RunSimulation(sc *SimulationConfig) (*Board, error) {
	if sc.Cols < 5 {
		return nil, errors.New("Too few columns in configuration, should be at least 5.")
	}

	if sc.Rows < 5 {
		return nil, errors.New("Too few rows in configuration, should be at least 5.")
	}

	if sc.Foods < 1 {
		return nil, errors.New("Too few foods, should be at least 1 food.")
	}

	if sc.Creature1 < 1 && sc.Creature2 < 1 {
		return nil, errors.New("Too few creatures, should be at least 1 creature.")
	}

	if sc.Creature1+sc.Creature2 > uint(((sc.Cols*2)+(sc.Rows*2-4))/2) {
		return nil, errors.New("Too many creatures, need to less than half available spawn locations.")
	}

	gameOn = true

	board := InitNewBoard(sc)
	rand.Seed(time.Now().UnixNano())

	if sc.Draw {
		drawer := InitDrawing(board)

		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			board.TickFrame()
			drawer.DrawFrame(board)

			if gameOn == false {
				break
			}
		}
	} else {
		for gameOn {
			board.TickFrame()
		}
	}

	printResults(board)

	return board, nil
}

func printResults(b *Board) {
	fmt.Println("A simulation was completed and these are the results:")
	fmt.Println("Total rounds: ", len(b.Rounds))
}
