package simulationgame

import (
	// "fmt"
	"fmt"
	"math/rand"
	"time"
)

var gameOn bool

func main() {
	fmt.Println("..")
}

func runSimulation(draw bool) *Board {
	board := InitNewBoard(40, 100)
	gameOn = true
	rand.Seed(time.Now().UnixNano())

	if draw {
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
	return board
}

func printResults(b *Board) {
	fmt.Println("A simulation was completed and these are the results:")
	fmt.Println("Total rounds: ", len(b.Rounds))
}
