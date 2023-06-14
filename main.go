package main

import (
	// "fmt"
	"fmt"
	"time"
)

var gameOn bool

func main() {
	gameOn = true
	board := InitNewBoard(40, 100)
	timer := 0.0

	ticker := time.Tick(time.Second / 10)
	tick := make(chan bool)

	gameOff := make(chan bool)

	InitDrawing(board)

	go func() {
		for {
			<-ticker // blocker statement which waits for input from the channel ticker - which is waiting for the time to run up (10th of a second)
			tick <- true
			timer += 0.1
		}
	}()

	go func() {
		for {
			<-tick // blocker statement which waits for tick to get a value...
			board.tickFrame()
			if gameOn == false {
				addMessageToCurrentGamelog("GAME SHOULD END")
				gameOff <- true
			}
		}
	}()

	select {
	case <-gameOff:
		fmt.Println("All creatures deaded")
	}
}
