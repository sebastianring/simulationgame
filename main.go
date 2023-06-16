package main

import (
	// "fmt"
	"time"
)

var gameOn bool

func main() {
	gameOn = true
	board := InitNewBoard(40, 100)
	InitDrawing(board)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		board.tickFrame()
		if gameOn == false {
			break
		}
	}
}
