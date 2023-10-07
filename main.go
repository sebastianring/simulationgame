package simulationgame

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sebastianring/simulationgame/cli"
	"log"
	"math/rand"
	"time"
)

type SimulationConfig struct {
	Rows      int
	Cols      int
	Draw      bool
	Foods     int
	Creature1 uint
	Creature2 uint
}

func main() {
	flag.Parse()
	RunSimulation(&flagConfig)
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

	board := NewBoard(sc)
	rand.Seed(time.Now().UnixNano())

	if sc.Draw {
		drawer := NewDrawer(board)
		drawer.DrawFrame(board)
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			board.TickFrame()
			drawer.DrawFrame(board)

			if board.GameOn == false {
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

				break
			}
		}
	} else {
		for board.GameOn {
			board.TickFrame()
		}
	}

	printResults(board)

	return board, nil
}

func GetStandardSimulationConfig() *SimulationConfig {
	return &SimulationConfig{
		Rows:      40,
		Cols:      100,
		Draw:      true,
		Foods:     100,
		Creature1: 15,
		Creature2: 15,
	}
}

func printResults(b *Board) {
	fmt.Println("A simulation was completed and these are the results:")
	fmt.Println("Total rounds: ", len(b.Rounds))
}
