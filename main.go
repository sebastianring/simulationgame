package main

import (
	// "fmt"
	"fmt"
	"math/rand"
	"time"
)

var gameOn bool

func main() {
	rand.Seed(time.Now().UnixNano())
	gameOn = true
	board := InitNewBoard(40, 100)
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

	fmt.Println("ENDED")
}
