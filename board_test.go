package simulationgame

import (
	"fmt"
	"os"
	"testing"
)

func TestRunSimulation(t *testing.T) {
	t.Setenv("SIM_GAME_DB_PW", os.Getenv("SIM_GAME_DB_PW"))

	sc := GetStandardSimulationConfig()

	result, err := RunSimulation(sc)

	if err != nil {
		t.Errorf("Some error %v", err.Error())
	}

	fmt.Println(result.Id)
}

func TestRunSimulationLowSurv(t *testing.T) {
	t.Setenv("SIM_GAME_DB_PW", os.Getenv("SIM_GAME_DB_PW"))

	sc := GetStandardSimulationConfig()
	sc.Draw = false
	sc.Foods = 10

	result, err := RunSimulation(sc)

	if err != nil {
		t.Errorf("Some error %v", err.Error())
	}

	fmt.Println(result.Id)
}

func TestRunSimulationHighSurv(t *testing.T) {
	t.Setenv("SIM_GAME_DB_PW", os.Getenv("SIM_GAME_DB_PW"))

	sc := GetStandardSimulationConfig()
	sc.Draw = false
	sc.Foods = 150

	result, err := RunSimulation(sc)

	if err != nil {
		t.Errorf("Some error %v", err.Error())
	} else {
		fmt.Println(result.Id)
	}

}

func TestStandardConfig(t *testing.T) {
	result := GetStandardSimulationConfig()

	if result == nil {
		t.Error("Error getting standard config")
	}
}

func TestRunSimulationHighSurvDraw(t *testing.T) {
	t.Setenv("SIM_GAME_DB_PW", os.Getenv("SIM_GAME_DB_PW"))

	sc := GetStandardSimulationConfig()
	sc.Foods = 150

	result, err := RunSimulation(sc)

	if err != nil {
		t.Errorf("Some error %v", err.Error())
	} else {
		fmt.Println(result.Id)
	}
}
