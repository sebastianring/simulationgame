package simulationgame

import (
	"fmt"
)

type MutationManager struct {
	speedMutation map[BoardObjectType]float32
	scanMutation  map[BoardObjectType]float32
}

type MutationCalc struct {
	chance int
	rate   float32
}

type MutationRate byte

const (
	flat     MutationRate = iota
	multiply MutationRate
)
