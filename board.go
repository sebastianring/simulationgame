package simulationgame

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// -------------------------------------------------- //
// -------------------------------------------------- //
// INITS AND STRUCTS -------------------------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //

type Board struct {
	Id                      string             `json:"id"`
	Rows                    int                `json:"rows"`
	Cols                    int                `json:"cols"`
	Gamelog                 *Gamelog           `json:"gamelog"`
	ObjectBoard             [][]BoardObject    `json:"object_board"`
	RoundInt                int                `json:"round_int"`
	Rounds                  []*Round           `json:"rounds"`
	CurrentRound            *Round             `json:"current_round"`
	CreatureIdCtr           map[string]int     `json:"creature_id_ctr"`
	Mutationrate            map[string]float32 `json:"mutationrate"`
	InitialFoods            int                `json:"initial_foods"`
	ConflictManager         *ConflictManager   `json:"conflict_manager"`
	AllFoodObjects          []Pos              `json:"all_food_objects"`
	AllAliveCreatureObjects []Pos              `json:"all_alive_creature_objects"`
	AllDeadCreatures        []*BoardObject     `json:"all_dead_creatures"`
	MaxRounds               int                `json:"max_rounds"`
}

type Round struct {
	Id                  int                        `json:"id"`
	Time                int                        `json:"time"`
	CreaturesSpawned    []CreatureObject           `json:"creatures_spawned"`
	CreaturesKilled     []CreatureObject           `json:"creatures_killed"`
	BoardLink           string                     `json:"board_link"`
	CreaturesSpawnedSum map[string]creatureSummary `json:"creatures_spawned_sum"`
	CreaturesKilledSum  map[string]creatureSummary `json:"creatures_killed_sum"`
}

type creatureSummary struct {
	TotalCreatures int     `json:"total_creatures"`
	TotalSpeed     int     `json:"total_speed"`
	AverageSpeed   float64 `json:"average_speed"`
}

type Pos struct {
	x int
	y int
}

type MoveType struct {
	action   Action
	conflict *ConflictInfo
}

type Action int

const (
	NoAction     Action = 0 // No action from creature
	MoveAction   Action = 1 // Action to move
	AttackAction Action = 2 // Action to attack
	WaitAction   Action = 3 // Action to wait
	AvoidAction  Action = 4 // Actively trying to avoid another creature
	FoodAction   Action = 5 // Getting food at new pos
)

type PosValidity int

const (
	InvalidPos    PosValidity = 0
	ConflictAtPos PosValidity = 1
	FoodAtPos     PosValidity = 2
	EmptyPos      PosValidity = 3
)

func InitNewBoard(sc *SimulationConfig) *Board {
	if sc.Rows < 5 || sc.Cols < 5 {
		fmt.Printf("Too few rows or cols: %v, rows: %v \n", sc.Rows, sc.Cols)
		os.Exit(1)
	}

	currentBoardId := uuid.New().String()

	newRound := Round{
		Id:               1,
		Time:             0,
		CreaturesSpawned: make([]CreatureObject, 0),
		CreaturesKilled:  make([]CreatureObject, 0),
		BoardLink:        currentBoardId,
	}

	cm, err := newConflictManager()

	if err != nil {
		fmt.Println("Error creating conflict manager, please debug.")
		os.Exit(1)
	}

	newBoard := Board{
		Id:              currentBoardId,
		Rows:            sc.Rows,
		Cols:            sc.Cols,
		Gamelog:         InitGamelog(sc.Rows, 40),
		ObjectBoard:     *createEmptyObjectsArray(sc.Rows, sc.Cols),
		RoundInt:        1,
		Rounds:          []*Round{&newRound},
		CurrentRound:    &newRound,
		CreatureIdCtr:   make(map[string]int, 0),
		Mutationrate:    make(map[string]float32, 0),
		InitialFoods:    sc.Foods,
		ConflictManager: cm,
		MaxRounds:       50,
	}

	newBoard.initBoardObjects()

	initialCreature1 := sc.Creature1
	initialCreature2 := sc.Creature2

	newBoard.spawnCreature1OnBoard(initialCreature1)
	newBoard.spawnCreature2OnBoard(initialCreature2)
	newBoard.spawnFoodOnBoard(newBoard.InitialFoods)

	db, err := openDbConnection()

	if err != nil {
		fmt.Println(err.Error())
		addMessageToCurrentGamelog(err.Error(), 1)
	}

	writeBoardToDb(db, &newBoard)

	addMessageToCurrentGamelog("Board added", 2)
	addMessageToCurrentGamelog("Welcome to the simulation game where you can simulate creatures and how they evolve!", 1)

	return &newBoard
}

func createEmptyObjectsArray(rows int, cols int) *[][]BoardObject {
	arr := make([][]BoardObject, rows)

	for i := 0; i < rows; i++ {
		arr[i] = make([]BoardObject, cols)
		for j := 0; j < cols; j++ {
			arr[i][j] = newEmptyObject()
		}
	}

	return &arr
}

// -------------------------------------------------- //
// -------------------------------------------------- //
// BOARD FUNCTIONS ---------------------------------- //
// -------------------------------------------------- //
// -------------------------------------------------- //

func (b *Board) spawnCreature1OnBoard(qty uint) {
	spawns := make([]Pos, 0)
	for uint(len(spawns)) < qty {
		newPos := b.randomPosAtEdgeOfMap()
		if !checkIfPosExistsInSlice(newPos, spawns) {
			spawns = append(spawns, newPos)
		}
	}

	for _, pos := range spawns {
		creaturePtr, err := b.newCreature1Object(false)

		if err != nil {
			fmt.Println("Error creating a new creature 1 object: " + err.Error())
			os.Exit(1)
		}

		b.ObjectBoard[pos.y][pos.x] = creaturePtr
		b.AllAliveCreatureObjects = append(b.AllAliveCreatureObjects, pos)
	}
}

func (b *Board) spawnCreature2OnBoard(qty uint) {
	spawns := make([]Pos, 0)
	for uint(len(spawns)) < qty {
		newPos := b.randomPosAtEdgeOfMap()
		if !checkIfPosExistsInSlice(newPos, spawns) && b.isSpotEmpty(newPos) {
			spawns = append(spawns, newPos)
		}
	}

	for _, pos := range spawns {
		creaturePtr, err := b.newCreature2Object(false)

		if err != nil {
			fmt.Println("Error creating a new creature 2 object: " + err.Error())
			os.Exit(1)
		}

		b.ObjectBoard[pos.y][pos.x] = creaturePtr
		b.AllAliveCreatureObjects = append(b.AllAliveCreatureObjects, pos)
	}
}

func (b *Board) spawnFoodOnBoard(qty int) {
	// qty := b.InitialFoods

	spawns := make([]Pos, 0)

	for len(spawns) < qty {
		newPos := b.randomPosWithinMap()
		if !checkIfPosExistsInSlice(newPos, spawns) && b.isSpotEmpty(newPos) {
			spawns = append(spawns, newPos)
		}
	}

	for _, pos := range spawns {
		b.ObjectBoard[pos.y][pos.x] = newFoodObject()
		b.AllFoodObjects = append(b.AllFoodObjects, pos)
	}
}

func (b *Board) isSpotEmpty(pos Pos) bool {
	if _, ok := b.ObjectBoard[pos.y][pos.x].(*EmptyObject); ok {
		return true
	}

	return false
}

func (b *Board) randomPosAtEdgeOfMap() Pos {
	// top = 0, right = 1, left = 2, bottom = 3
	edge := rand.Intn(4)
	var x int
	var y int

	if edge == 0 {
		y = 0
		x = rand.Intn(b.Cols - 1)
	} else if edge == 1 {
		x = b.Cols - 1
		y = rand.Intn(b.Rows - 1)
	} else if edge == 2 {
		x = 0
		y = rand.Intn(b.Rows - 1)
	} else {
		x = rand.Intn(b.Cols - 1)
		y = b.Rows - 1
	}

	return Pos{x, y}
}

func (b *Board) initBoardObjects() {
	b.CreatureIdCtr["Creature1"] = 1
	b.CreatureIdCtr["Creature2"] = 1

	b.Mutationrate = make(map[string]float32)
	b.Mutationrate["Creature1"] = 0.1
	b.Mutationrate["Creature2"] = 0.1
}

func (b *Board) randomPosWithinMap() Pos {
	minDistanceFromBorder := 3
	x := rand.Intn(b.Cols-minDistanceFromBorder*2) + minDistanceFromBorder
	y := rand.Intn(b.Rows-minDistanceFromBorder*2) + minDistanceFromBorder

	return Pos{x, y}
}

func checkIfPosExistsInSlice(pos Pos, slice []Pos) bool {
	for _, slicePos := range slice {
		if pos.x == slicePos.x && pos.y == slicePos.y {
			return true
		}
	}

	return false
}

func (b *Board) TickFrame() {
	b.CurrentRound.Time++
	b.creatureUpdatesPerTick()

	// ----------- debugging creatures - print speed and id ---------- //

	// res := make([]string, 1)
	//
	// for _, pos := range allAliveCreatureObjects {
	// 	if obj, ok := b.ObjectBoard[pos.y][pos.x].(CreatureObject); ok {
	// 		res = append(res, strconv.Itoa(obj.getId())+":"+strconv.Itoa(obj.getId()))
	// 	}
	// }

	// addMessageToCurrentGamelog(strings.Join(res, ", "), 2)

	// ----------- end debugging ------------------------------------ //

}

func checkCreatureType(bo BoardObject) (bool, *BoardObject) {
	switch bo.(type) {
	case *Creature1:
		return true, &bo
	case *Creature2:
		return true, &bo
	default:
		return false, nil
	}
}

// debug help
func getCurrentTimeString() string {
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02 15:04:05.000")

	return timeString
}

// func (b *Board) killCreature(creature CreatureObject, pos Pos) Pos {
// 	b.CurrentRound.creaturesKilled = append(b.CurrentRound.creaturesKilled, creature)
// 	b.ObjectBoard[pos.y][pos.x] = newEmptyObject()
//
// 	return pos
// }

func (b *Board) creatureUpdatesPerTick() {
	updatedAllCreatureObjects := make([]Pos, 0)
	deadCreatures := make([]Pos, 0)

	for _, pos := range b.AllAliveCreatureObjects {
		if obj, ok := b.ObjectBoard[pos.y][pos.x].(CreatureObject); ok {
			action := obj.updateTick()

			if action == "move" {
				newPos := Pos{
					x: -1,
					y: -1}

				moveType := MoveType{}

				for {
					newPos, moveType = b.newPosAndMove(pos)

					if moveType.action != AvoidAction {
						break
					}
				}

				if moveType.action == WaitAction {
					break
				}

				// Note to self: need to update this whole section - it works but its not beautiful... at all
				if moveType.action == AttackAction {
					// addMessageToCurrentGamelog("Conflict at: "+strconv.Itoa(newPos.x)+", "+strconv.Itoa(newPos.y)+" ", 1)

					switch moveType.conflict.Conflict {
					case Share:
						b.ConflictManager.share(moveType.conflict.SourceCreature, moveType.conflict.TargetCreature)
						updatedAllCreatureObjects = append(updatedAllCreatureObjects, pos)

					case Attack1:
						b.ConflictManager.attack1(moveType.conflict.SourceCreature, moveType.conflict.TargetCreature)

						deadCreatures = append(deadCreatures, newPos)
						b.ObjectBoard[newPos.y][newPos.x] = BoardObject(obj)
						b.ObjectBoard[pos.y][pos.x] = newEmptyObject()

					case Attack2:
						killTarget := b.ConflictManager.attack2(moveType.conflict.SourceCreature, moveType.conflict.TargetCreature)

						if killTarget {
							addMessageToCurrentGamelog(moveType.conflict.SourceCreature.getIdAsString()+
								" killed "+moveType.conflict.TargetCreature.getIdAsString(), 1)

							// deadCreatures = append(deadCreatures, newPos)
							b.deleteCreature(pos, &b.ObjectBoard[newPos.y][newPos.x])
							b.ObjectBoard[newPos.y][newPos.x] = newEmptyObject()
							updatedAllCreatureObjects = append(updatedAllCreatureObjects, newPos)

						} else {
							addMessageToCurrentGamelog(moveType.conflict.TargetCreature.getIdAsString()+
								" killed "+moveType.conflict.SourceCreature.getIdAsString(), 1)

							// deadCreatures = append(deadCreatures, pos)
							b.deleteCreature(pos, &b.ObjectBoard[newPos.y][newPos.x])
							b.ObjectBoard[pos.y][pos.x] = newEmptyObject()
							updatedAllCreatureObjects = append(updatedAllCreatureObjects, newPos)
						}

					default:
						addMessageToCurrentGamelog("Conflict manager not setup properly", 1)
					}

				} else {
					b.ObjectBoard[newPos.y][newPos.x] = BoardObject(obj)

					if moveType.action == FoodAction {
						addMessageToCurrentGamelog("Food eaten by creature id: "+strconv.Itoa(obj.getId()), 2)
						obj.heal(obj.getOriHP())
						b.ObjectBoard[pos.y][pos.x] = newEmptyObject()
						b.deleteFood(newPos)
					} else {
						b.ObjectBoard[pos.y][pos.x] = newEmptyObject()
					}

					updatedAllCreatureObjects = append(updatedAllCreatureObjects, newPos)
				}

			} else if action == "dead" {
				deadCreatures = append(deadCreatures, pos)

			} else {
				updatedAllCreatureObjects = append(updatedAllCreatureObjects, pos)

			}
		}
	}

	// delete dead creatures after tick is complete
	for _, pos := range deadCreatures {
		b.deleteCreature(pos, &b.ObjectBoard[pos.y][pos.x])
		b.ObjectBoard[pos.y][pos.x] = newEmptyObject()
	}

	// update all creatures from last tick
	b.AllAliveCreatureObjects = updatedAllCreatureObjects

	if b.checkIfCreaturesAreInactive() {
		if b.checkIfCreaturesAreDead() {
			gameOn = false
			addMessageToCurrentGamelog("All creatures are dead, end the game", 1)
		}

		b.newRound()
	}
}

func (b *Board) newRound() {
	addMessageToCurrentGamelog("All creatures are dead or have eaten, starting new round", 1)
	b.writeSummaryOfRound()
	b.Gamelog.writeGamelogToFile()

	if len(b.Rounds) >= b.MaxRounds {
		gameOn = false

		addMessageToCurrentGamelog("Max number of rounds reached, ending the game.", 1)
		fmt.Println("Max number of rounds reached, ending the game.")

	} else {
		b.spawnOffsprings()
		b.findPosForAllCreatures()
		b.deleteAndSpawnFood()

		newRound := Round{
			Id:               b.CurrentRound.Id + 1,
			Time:             0,
			CreaturesSpawned: make([]CreatureObject, 0),
			CreaturesKilled:  make([]CreatureObject, 0),
		}

		b.CurrentRound = &newRound
		b.Rounds = append(b.Rounds, &newRound)

		addMessageToCurrentGamelog("---- NEW ROUND ----", 1)
	}

}

func (b *Board) deleteAndSpawnFood() {
	for _, pos := range b.AllFoodObjects {
		b.ObjectBoard[pos.y][pos.x] = newEmptyObject()
	}

	b.AllFoodObjects = make([]Pos, 0)
	b.spawnFoodOnBoard(b.InitialFoods)
}

func (b *Board) writeSummaryOfRound() {
	creaturesSpawned := make(map[string]creatureSummary)

	for _, creature := range b.CurrentRound.CreaturesSpawned {
		creatureType := creature.getType()

		if obj, ok := creaturesSpawned[creatureType]; ok {
			mapPtr := &obj
			mapPtr.TotalCreatures += 1
			mapPtr.TotalSpeed += creature.getSpeed()

		} else {
			newCreatureSummary := creatureSummary{
				TotalCreatures: 1,
				TotalSpeed:     creature.getSpeed(),
			}

			creaturesSpawned[creatureType] = newCreatureSummary
		}
	}

	b.CurrentRound.CreaturesSpawnedSum = creaturesSpawned

	for c, cs := range b.CurrentRound.CreaturesSpawnedSum {
		cs.AverageSpeed = float64(cs.TotalSpeed) / float64(cs.TotalCreatures)

		addMessageToCurrentGamelog("In last round, "+strconv.Itoa(cs.TotalCreatures)+
			" x "+c+" was spawned with the average speed of: "+strconv.FormatFloat(cs.AverageSpeed, 'f', 2, 64), 1)
	}

	creaturesKilled := make(map[string]creatureSummary)

	for _, creature := range b.CurrentRound.CreaturesKilled {
		creatureType := creature.getType()

		if obj, ok := creaturesKilled[creatureType]; ok {
			mapPtr := &obj
			mapPtr.TotalCreatures += 1
			mapPtr.TotalSpeed += creature.getSpeed()

		} else {
			newCreatureSummary := creatureSummary{
				TotalCreatures: 1,
				TotalSpeed:     creature.getSpeed(),
			}

			creaturesKilled[creatureType] = newCreatureSummary
		}
	}

	b.CurrentRound.CreaturesKilledSum = creaturesKilled

	for c, cs := range b.CurrentRound.CreaturesKilledSum {
		cs.AverageSpeed = float64(cs.TotalSpeed) / float64(cs.TotalCreatures)

		addMessageToCurrentGamelog("In last round, "+strconv.Itoa(cs.TotalCreatures)+
			" x "+c+" was killed with the average speed of: "+strconv.FormatFloat(cs.AverageSpeed, 'f', 2, 64), 1)
	}
}

func (b *Board) findPosForAllCreatures() {
	for i, creaturePos := range b.AllAliveCreatureObjects {
		if obj, ok := b.ObjectBoard[creaturePos.y][creaturePos.x].(CreatureObject); ok {
			findNewPos := false

			for !findNewPos {
				newPos := b.randomPosAtEdgeOfMap()
				if b.isSpotEmpty(newPos) {
					b.ObjectBoard[newPos.y][newPos.x] = obj
					obj.resetValues()
					b.ObjectBoard[creaturePos.y][creaturePos.x] = newEmptyObject()
					b.AllAliveCreatureObjects[i] = newPos
					findNewPos = true
				}
			}
		}
	}
}

func (b *Board) spawnOffsprings() {
	creatureQty := map[string]uint{
		"creature1": 0,
		"creature2": 0,
	}

	for _, pos := range b.AllAliveCreatureObjects {
		if obj, ok := b.ObjectBoard[pos.y][pos.x].(CreatureObject); ok {
			if obj.ifOffspring() {
				var offspring CreatureObject
				var err error

				if obj2, ok := obj.(*Creature1); ok {
					offspring, err = b.newCreature1Object(true, obj2)

					if err != nil {
						fmt.Println("Error creating offspring: " + err.Error())
					}

					creatureQty["creature1"]++

				} else if obj2, ok := obj.(*Creature2); ok {
					offspring, err = b.newCreature2Object(true, obj2)
					if err != nil {
						fmt.Println("Error creating offspring: " + err.Error())
					}

					creatureQty["creature2"]++
				}

				b.CurrentRound.CreaturesSpawned = append(b.CurrentRound.CreaturesSpawned, offspring)

				newPos := b.randomPosAtEdgeOfMap()
				for !b.isSpotEmpty(newPos) {
					newPos = b.randomPosAtEdgeOfMap()
				}

				b.ObjectBoard[newPos.y][newPos.x] = offspring
				b.AllAliveCreatureObjects = append(b.AllAliveCreatureObjects, newPos)
			}
		}
	}

	for key, val := range creatureQty {
		if val > 0 {
			addMessageToCurrentGamelog(strconv.Itoa(int(val))+" Creatures of type "+key+" spawned", 1)
		}
	}

	b.spawnCreature1OnBoard(creatureQty["creature1"])
	b.spawnCreature2OnBoard(creatureQty["creature2"])
}

func (b *Board) checkIfCreaturesAreDead() bool {
	for _, pos := range b.AllAliveCreatureObjects {
		if obj, ok := b.ObjectBoard[pos.y][pos.x].(CreatureObject); ok {
			dead := obj.isDead()

			if !dead {
				return false
			}
		}
	}

	return true
}

func (b *Board) checkIfCreaturesAreInactive() bool {
	for _, pos := range b.AllAliveCreatureObjects {
		if obj, ok := b.ObjectBoard[pos.y][pos.x].(CreatureObject); ok {
			dead := obj.isDead()
			moving := obj.isMoving()

			if !dead && moving || dead {
				return false
			}
		}

	}

	return true
}

func (b *Board) newPosAndMove(currentPos Pos) (Pos, MoveType) {
	newPos := Pos{-1, -1}
	moveType := MoveType{
		action: NoAction,
	}

	// validPosTypes := []PosValidity{EmptyPos, FoodAtPos}
	counter := 0

	// HOW TO MAKE THE CREATURES MOVE INWARDS TO LOOK FOR FOOD?
	// The closer they are to one edge, the more probable they are to move towards the other edge?
	// Example: x = 99, y = 40
	// Width-x = the probability to move the left
	// Height-y = the probability to move upwards?

	for newPos.x == -1 || newPos.y == -1 {
		direction := rand.Intn(2) // 0 = x movement, 1 = y-movement
		var x int
		var y int
		if direction == 0 {
			xdirection := rand.Intn(b.Cols)
			xprobability := b.Cols - 1 - currentPos.x
			if xdirection < xprobability {
				x = currentPos.x + 1
			} else {
				x = currentPos.x - 1
			}
			y = currentPos.y
		} else {
			ydirection := rand.Intn(b.Rows)
			yprobability := b.Rows - 1 - currentPos.y
			if ydirection < yprobability {
				y = currentPos.y + 1
			} else {
				y = currentPos.y - 1
			}
			x = currentPos.x
		}

		valid := false

		for {
			newPosType := b.checkIfNewPosIsValid(x, y)

			// if containsString(validMoveTypes, moveType.action) {
			// 	valid = true
			// 	break
			// }

			if newPosType == FoodAtPos {
				moveType.action = FoodAction
				valid = true
				break
			} else if newPosType == EmptyPos {
				moveType.action = MoveAction
				valid = true
				break
			}

			// if moveType.action == "conflict" {
			if newPosType == ConflictAtPos {
				sourceCreature, err := b.getCreatureObjectFromBoard(currentPos)

				if err != nil {
					addMessageToCurrentGamelog(err.Error(), 1)
				}

				targetCreature, err := b.getCreatureObjectFromBoard(Pos{x: x, y: y})

				if err != nil {
					addMessageToCurrentGamelog(err.Error(), 1)
				}

				action, conflict := b.ConflictManager.getConflict(sourceCreature, targetCreature)

				if action {
					valid = true
					moveType.action = AttackAction
					moveType.conflict = conflict
					break
				} else {
					break
				}
			} else {
				break
			}
		}

		if counter > 10 {
			moveType.action = NoAction
			return newPos, moveType
		}

		if valid {
			newPos.x = x
			newPos.y = y
			return newPos, moveType
		}
	}

	return newPos, moveType
}

func (b *Board) getCreatureObjectFromBoard(pos Pos) (CreatureObject, error) {
	if creature, ok := b.ObjectBoard[pos.y][pos.x].(CreatureObject); ok {
		return creature, nil
	}

	return nil, errors.New("Was not a creature")
}

func (b *Board) checkIfNewPosIsValid(x int, y int) PosValidity {
	if x < 0 || x >= b.Cols || y < 0 || y >= b.Rows {
		return InvalidPos
	}

	if _, ok := b.ObjectBoard[y][x].(*EmptyObject); ok {
		return EmptyPos
	} else if _, ok := b.ObjectBoard[y][x].(*Food); ok {
		return FoodAtPos
	} else if obj, ok := b.ObjectBoard[y][x].(CreatureObject); ok {

		if obj.isMoving() {
			return InvalidPos
		}

		return ConflictAtPos
	}

	return InvalidPos
}

func (b *Board) deleteFood(pos Pos) {
	var element int
	for i, val := range b.AllFoodObjects {
		if val.x == pos.x && val.y == pos.y {
			element = i
			break
		}
	}

	b.AllFoodObjects = deleteIndexInPosSlice(b.AllFoodObjects, element)
}

func (b *Board) deleteCreature(pos Pos, creature *BoardObject) {
	b.AllDeadCreatures = append(b.AllDeadCreatures, creature)
	var element int
	for i, val := range b.AllAliveCreatureObjects {
		if val.x == pos.x && val.y == pos.y {
			element = i
			break
		}
	}

	b.AllAliveCreatureObjects = deleteIndexInPosSlice(b.AllAliveCreatureObjects, element)
}

func deleteIndexInPosSlice(posSlice []Pos, index int) []Pos {
	posSlice[index] = posSlice[len(posSlice)-1]
	return posSlice[:len(posSlice)-1]
}

func containsString(slice []string, target string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}

	return false
}

func containsValidPosType(slice []PosValidity, posType PosValidity) bool {
	for _, element := range slice {
		if element == posType {
			return true
		}
	}

	return false
}
