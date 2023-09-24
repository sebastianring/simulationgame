package simulationgame

import (
	"fmt"
	"testing"
)

func TestRunSimulation(t *testing.T) {
	t.Setenv("sim_game", "valmet865")

	sc := GetStandardSimulationConfig()

	result, err := RunSimulation(sc)

	if err != nil {
		t.Errorf("Some error %v", err.Error())
	}

	fmt.Println(result.Id)
}

func TestStandardConfig(t *testing.T) {
	t.Setenv("sim_game", "valmet865")

	result := GetStandardSimulationConfig()

	if result == nil {
		t.Error("Error getting standard config")
	}
}
