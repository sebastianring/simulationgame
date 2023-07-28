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
	board := InitNewBoard(30, 50)
	InitDrawing(board)

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	// db, err := openDbConnection()
	//
	// if err != nil {
	// 	addMessageToCurrentGamelog(err.Error(), 1)
	// }
	//
	// defer db.Close()
	//
	// testMsg := message{
	// 	Id:    1000,
	// 	Prio:  1,
	// 	Texts: "hello",
	// }
	//
	// writeMessageToDb(db, &testMsg)

	for range ticker.C {
		board.tickFrame()

		if gameOn == false {
			break
		}
	}

	fmt.Println("ENDED")
}
