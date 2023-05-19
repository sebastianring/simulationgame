package main

import "fmt"

func main() {
	fmt.Println()

	board := InitNewBoard(40, 100)

	InitDrawing(board)
	DrawFrame(board)
}
