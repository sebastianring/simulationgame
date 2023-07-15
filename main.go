package main

import (
	// "fmt"
	"math/rand"
	"time"
)

var gameOn bool

func main() {
	rand.Seed(time.Now().UnixNano())
	gameOn = true
	board := InitNewBoard(40, 100)
	InitDrawing(board)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	_, err := openDbConnection()

	if err != nil {
		addMessageToCurrentGamelog(err.Error(), 1)
	}

	for range ticker.C {
		board.tickFrame()
		if gameOn == false {
			break
		}
	}
}
