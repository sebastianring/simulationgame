package main

import (
	// "fmt"
	"time"
)

var gameOn bool

func main() {
	board := InitNewBoard(40, 100)
	timer := 0.0

	ticker := time.Tick(time.Second / 10)
	tick := make(chan bool)

	gameOn := true

	InitDrawing(board)

	for gameOn {
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
			}
		}()

		select {}
	}
}
