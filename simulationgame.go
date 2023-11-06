package simulationgame

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type SimulationConfig struct {
	Rows        int
	Cols        int
	Draw        bool
	Foods       int
	Creature1   uint
	Creature2   uint
	MaxRounds   int
	GamelogSize int
}

func RunSimulation(sc *SimulationConfig) (*Board, error) {
	// Rules should be synced with simgameserver
	if sc.Cols < 5 || sc.Cols > 150 {
		return nil, errors.New("Cols outside 5-150 interval, please adjust.")
	}

	if sc.Rows < 5 || sc.Cols > 150 {
		return nil, errors.New("Rows outside 5-150 interval, please adjust.")
	}

	maxFoods := int(sc.Cols * sc.Rows / 2)

	if sc.Foods < 1 || sc.Foods > maxFoods {
		return nil, errors.New("Foods outside interval, min = 1, max = " + strconv.Itoa(maxFoods))
	}

	maxCreatures := uint(((sc.Cols * 2) + (sc.Rows * 2) - 4) / 2)

	if sc.Creature1 < 1 && sc.Creature2 < 1 || sc.Creature1 > maxCreatures || sc.Creature2 > maxCreatures {
		return nil, errors.New("Creatures outside interval, min = 1, max = " + strconv.Itoa(int(maxCreatures)))
	}

	if sc.MaxRounds < 1 || sc.MaxRounds > 100 {
		return nil, errors.New("Max rounds outside 1-100 interval, please adjust.")
	}

	if sc.GamelogSize < 20 || sc.GamelogSize > 75 {
		return nil, errors.New("Gamelog size outside 20-75 interval, please adjust.")
	}

	board := NewBoard(sc)
	rand.Seed(time.Now().UnixNano())

	if sc.Draw {
		drawer := NewDrawer(board)
		drawer.DrawFrame(board)
		ticker := time.NewTicker(25 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			board.TickFrame()
			drawer.DrawFrame(board)

			if board.GameOn == false {
				break
			}
		}
	} else {
		for board.GameOn {
			board.TickFrame()
		}
	}

	board.Gamelog.writeGamelogToFile()

	log.Println("Saving board to DB.")

	err := writeBoardToDb(board)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Succesfully wrote board to DB.")
	}

	log.Println("Saving messages to DB.")

	err = writeMessagesToDb(board)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Succesfully wrote messages to DB.")
	}

	printResults(board)

	return board, nil
}

func GetStandardSimulationConfig() *SimulationConfig {
	return &SimulationConfig{
		Rows:        35,
		Cols:        100,
		Draw:        true,
		Foods:       70,
		Creature1:   15,
		Creature2:   15,
		MaxRounds:   50,
		GamelogSize: 40,
	}
}

func printResults(b *Board) {
	fmt.Println("A simulation was completed and these are the results:")
	fmt.Println("Total rounds: ", len(b.Rounds))
}
