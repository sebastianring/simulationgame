package cli

import (
	"flag"
)

var flagConfig SimulationConfig

func init() {
	flag.IntVar(&flagConfig.Rows, "rows", 0, "Input number of rows for simulation")
	flag.IntVar(&flagConfig.Cols, "cols", 0, "Input number of cols for simulation")
	flag.IntVar(&flagConfig.Foods, "foods", 0, "Input number of foods for simulation")
	flag.BoolVar(&flagConfig.Draw, "draw", false, "Input if it should draw the simulation")
	flag.UintVar(&flagConfig.Creature1, "creature1", 0, "Input number of creature 1.")
	flag.UintVar(&flagConfig.Creature2, "creature2", 0, "Input number of creature 2.")

	flag.PrintDefaults()
}
