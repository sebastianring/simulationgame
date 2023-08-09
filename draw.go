package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

const edgeSymbol = byte(35) // ###### as the symbol at the edges of gui
const spaceSymbol = byte(32)

var totalWidth int
var totalHeight int
var currentOs string

func InitDrawing(b *Board) {
	currentOs = runtime.GOOS
	addMessageToCurrentGamelog("Current OS identified: "+currentOs, 1)
	totalWidth = b.cols + b.gamelog.cols + 2 + 1 + 2
	totalHeight = b.rows + 2 + 2 // rows + (edges + status bar) + (status bar line)
}

func DrawFrame(b *Board) {
	clearScreen()

	for i := 0; i < totalHeight; i++ {
		if i == 1 {
			b.printStatusLine(totalWidth)
		} else if i == 0 || i == 2 || i == totalHeight-1 {
			printSymbolLine(totalWidth)
		} else {
			printDataLine(b.objectBoard[i-3], b.gamelog, i-3)
		}
	}
}

func (b *Board) printStatusLine(totalWidth int) {
	fmt.Println("ROUND: " + strconv.Itoa(b.currentRound.id) +
		"      TIME: " + strconv.Itoa(b.currentRound.time) +
		"   CREATURES ACTIVE: " + strconv.Itoa(len(allAliveCreatureObjects)) +
		"     FOOD LEFT: " + strconv.Itoa(len(allFoodObjects)))
}

func printDataLine(boardData []BoardObject, gl *Gamelog, messageRow int) {
	line := make([]byte, totalWidth)

	line = append(line, edgeSymbol) // adding a # symbol at the start

	boardDataLine := getBoardSymbolByRow(boardData)
	for _, val := range boardDataLine {
		line = append(line, val)
	}

	line = append(line, edgeSymbol)  // adding a # symbol at the start
	line = append(line, spaceSymbol) // adding a " " symbol at the start

	gamelogDataLine := gl.getMessageByRow(messageRow)
	for _, val := range gamelogDataLine {
		line = append(line, val)
	}

	line = append(line, spaceSymbol) // adding a " " symbol at the end
	line = append(line, edgeSymbol)  // adding a # symbol at the end

	fmt.Println(string(line))
}

func printSymbolLine(length int) {
	dash := edgeSymbol
	line := make([]byte, length)

	for i := 0; i < length; i++ {
		line[i] = dash
	}

	fmt.Println(string(line))
}

func getBoardSymbolByRow(row []BoardObject) []byte {
	line := make([]byte, len(row))

	for _, object := range row {
		symbol := object.getSymbol()
		for _, symbval := range symbol {
			line = append(line, symbval)
		}
		// line = append(line, object.getSymbol())
	}

	return line
}

func clearScreen() {
	osCommand := map[string]string{
		"windows": "cls",
		"linux":   "clear",
		"darwin":  "clear",
	}

	cmd := exec.Command(osCommand[currentOs])
	cmd.Stdout = os.Stdout
	cmd.Run()
}
