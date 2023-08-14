package main

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var gameOn bool

func main() {
	resultBoard := runSimulation(false)
	printResults(resultBoard)

	fmt.Println("ENDED")
}

func runServer() {
	http.HandleFunc("/api/new_sim", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Something is happening")
		resultBoard := runSimulation(false)
		jsonBytes, err := json.Marshal(resultBoard)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			os.Exit(1)
		}

		w.Header().Set("Content-type", "application/json")
		w.Write(jsonBytes)
	})

	fmt.Println("Server running at port 8080")
	http.ListenAndServe(":8080", nil)
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
	fmt.Println("Total rounds: ", len(b.rounds))
}
