package simulationgame

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

type Drawer struct {
	totalWidth  int
	totalHeight int
	currentOs   string

	edgeSymbol  byte
	spaceSymbol byte
}

func NewDrawer(b *Board) *Drawer {
	newDrawer := Drawer{
		totalHeight: b.Rows + 2 + 2,                      // rows + (edges + status bar) + (status bar line)
		totalWidth:  b.Cols + b.Gamelog.cols + 2 + 1 + 2, // board cols + gamelog cols + board edges + line between board and gamelog + gamelog edges
		currentOs:   runtime.GOOS,

		edgeSymbol:  byte(35), // #### Edge of GUI
		spaceSymbol: byte(32), // "  " Space
	}

	addMessageToCurrentGamelog("Current OS identified: "+newDrawer.currentOs, 1)

	return &newDrawer
}

func (d *Drawer) DrawFrame(b *Board) {
	d.clearScreen()

	for i := 0; i < d.totalHeight; i++ {
		if i == 1 {
			b.printStatusLine(d.totalWidth)
		} else if i == 0 || i == 2 || i == d.totalHeight-1 {
			d.printSymbolLine(d.totalWidth)
		} else {
			d.printDataLine(b.ObjectBoard[i-3], b.Gamelog, i-3)
		}
	}
}

// Should be adapted so the status line always centered
func (b *Board) printStatusLine(totalWidth int) {
	fmt.Println("      ROUND: " + strconv.Itoa(b.CurrentRound.Id) +
		"      TIME: " + strconv.Itoa(b.CurrentRound.Time) +
		"   CREATURES ALIVE: " + strconv.Itoa(len(b.AliveCreatureObjects)) +
		"     FOOD LEFT: " + strconv.Itoa(len(b.AllFoodObjects)))
}

func (d *Drawer) printDataLine(boardData []BoardObject, gl *Gamelog, messageRow int) {
	line := make([]byte, d.totalWidth)

	line = append(line, d.edgeSymbol) // adding a # symbol at the start

	boardDataLine := getBoardSymbolByRow(boardData)
	for _, val := range boardDataLine {
		line = append(line, val)
	}

	line = append(line, d.edgeSymbol)  // adding a # symbol at the start
	line = append(line, d.spaceSymbol) // adding a " " symbol at the start

	gamelogDataLine := gl.getMessageByRow(messageRow)
	// addMessageToCurrentGamelog(string(gamelogDataLine), 1)

	for _, val := range gamelogDataLine {
		line = append(line, val)
	}

	line = append(line, d.spaceSymbol) // adding a " " symbol at the end
	line = append(line, d.edgeSymbol)  // adding a # symbol at the end

	fmt.Println(string(line))
}

func (d *Drawer) printSymbolLine(length int) {
	dash := d.edgeSymbol
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
	}

	return line
}

func (d *Drawer) clearScreen() {
	osCommand := map[string]string{
		"windows": "cls",
		"linux":   "clear",
		"darwin":  "clear",
	}

	cmd := exec.Command(osCommand[d.currentOs])
	cmd.Stdout = os.Stdout
	cmd.Run()
}
