package simulationgame

import (
	"fmt"
	"testing"
)

func TestRunSimulation(t *testing.T) {

	t.Setenv("sim_game", "valmet865")

	testSC := SimulationConfig{
		Rows:      40,
		Cols:      100,
		Foods:     100,
		Draw:      true,
		Creature1: 20,
		Creature2: 20,
	}

	result, err := RunSimulation(&testSC)

	if err != nil {
		t.Errorf("Some error %v", err.Error())
	}

	fmt.Println(result.Id)
}
