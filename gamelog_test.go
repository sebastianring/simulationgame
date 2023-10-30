package simulationgame

import (
	"fmt"
	"os"
	"testing"
)

func TestGetFileString(t *testing.T) {
	str, err := getFileString()

	if err != nil {
		t.Errorf(err.Error())
	}

	expected := "logs/"

	fmt.Printf("exp: %v, str: %v \n", expected, str[0:5])

	if str[0:5] != expected {
		t.Errorf("expected: %v, str: %v \n", expected, str)
	}
}

func TestWriteGamelogToFile(t *testing.T) {
	t.Setenv("SIM_GAME_DB_PW", os.Getenv("SIM_GAME_DB_PW"))

	sc := GetStandardSimulationConfig()
	sc.Draw = false

	board := NewBoard(sc)

	addMessageToCurrentGamelog("LOL", 1)

	board.Gamelog.writeGamelogToFile()
}
