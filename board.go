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
	Id                   string                  `json:"id"`
	GameOn               bool                    `json:"game_on"`
	Rows                 int                     `json:"rows"`
	Cols                 int                     `json:"cols"`
	Gamelog              *Gamelog                `json:"gamelog"`
	ObjectBoard          [][]BoardObject         `json:"object_board"`
	Rounds               []*Round                `json:"rounds"`
	CurrentRound         *Round                  `json:"current_round"`
	CreatureIdCtr        map[BoardObjectType]int `json:"creature_id_ctr"`
	MutationManager      *MutationManager        `json:"mutation_manager"`
	InitialFoods         int                     `json:"initial_foods"`
	ConflictManager      *ConflictManager        `json:"conflict_manager"`
	AllFoodObjects       []Pos                   `json:"all_food_objects"`
	AliveCreatureObjects []CreatureObject        `json:"alive_creature_objects"`
	AllDeadCreatures     []BoardObject           `json:"all_dead_creatures"`
	MaxRounds            int                     `json:"max_rounds"`
}

type Round struct {
	Id                     int                                  `json:"id"`
	Time                   int                                  `json:"time"`
	CreaturesSpawned       []CreatureObject                     `json:"creatures_spawned"`
	CreaturesKilled        []CreatureObject                     `json:"creatures_killed"`
	CreaturesAliveAtEnd    []CreatureObject                     `json:"creatures_alive_at_end"`
	BoardLink              string                               `json:"board_link"`
	CreaturesSpawnedSum    map[BoardObjectType]*CreatureSummary `json:"creatures_spawned_sum"`
	CreaturesKilledSum     map[BoardObjectType]*CreatureSummary `json:"creatures_killed_sum"`
	CreaturesAliveAtEndSum map[BoardObjectType]*CreatureSummary `json:"creatures_alive_at_end_sum"`
}

type CreatureSummary struct {
	CreatureType      string  `json:"creature_type"`
	TotalCreatures    int     `json:"total_creatures"`
	TotalSpeed        float64 `json:"total_speed"`
	AverageSpeed      float64 `json:"average_speed"`
	TotalScanChance   float64 `json:"total_scan_chance"`
	AverageScanChance float64 `json:"average_scan_chance"`
}

type Pos struct {
	x int
	y int
}

type MoveType struct {
	action       Action
	conflictinfo *ConflictInfo
}

type Action byte

const (
	NoAction     Action = 0 // No action from creature
	MoveAction   Action = 1 // Action to move
	AttackAction Action = 2 // Action to attack
	WaitAction   Action = 3 // Action to wait
	AvoidAction  Action = 4 // Actively trying to avoid another creature
	FoodAction   Action = 5 // Getting food at new pos
)

type PosValidity byte

const (
	InvalidPos    PosValidity = 0
	ConflictAtPos PosValidity = 1
	FoodAtPos     PosValidity = 2
	EmptyPos      PosValidity = 3
)

func NewBoard(sc *SimulationConfig) *Board {
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

	mm, err := newMutationManager()

	if err != nil {
		fmt.Println("Error creating mutation manager, please debug.")
	}

	newBoard := Board{
		Id:              currentBoardId,
		GameOn:          true,
		Rows:            sc.Rows,
		Cols:            sc.Cols,
		Gamelog:         NewGamelog(sc.Rows, 40),
		ObjectBoard:     *createEmptyObjectsArray(sc.Rows, sc.Cols),
		Rounds:          []*Round{&newRound},
		CurrentRound:    &newRound,
		CreatureIdCtr:   make(map[BoardObjectType]int, 0),
		MutationManager: mm,
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
		creature, err := b.newCreature1Object(false)

		if err != nil {
			fmt.Println("Error creating a new creature 1 object: " + err.Error())
			os.Exit(1)
		}

		creature.setPos(pos)

		b.ObjectBoard[pos.y][pos.x] = creature
		b.AliveCreatureObjects = append(b.AliveCreatureObjects, creature)
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
		creature, err := b.newCreature2Object(false)

		if err != nil {
			fmt.Println("Error creating a new creature 2 object: " + err.Error())
			os.Exit(1)
		}

		creature.setPos(pos)

		b.ObjectBoard[pos.y][pos.x] = creature
		b.AliveCreatureObjects = append(b.AliveCreatureObjects, creature)
	}
}

func (b *Board) spawnFoodOnBoard(qty int) {
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

	for {
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

		if b.cornerAvoided(Pos{x, y}) {
			break
		}
	}

	return Pos{x, y}
}

func (b *Board) cornerAvoided(pos Pos) bool {
	if pos.x == 0 && pos.y == 0 ||
		pos.x == b.Cols-1 && pos.y == b.Rows-1 ||
		pos.x == 0 && pos.y == b.Rows-1 ||
		pos.x == b.Cols-1 && pos.y == 0 {

		return false
	}

	return true
}

func (b *Board) initBoardObjects() {
	b.CreatureIdCtr[Creature1Type] = 1
	b.CreatureIdCtr[Creature2Type] = 1
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

func checkBoardObjectType(bo BoardObject) (bool, *BoardObject) {
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

func (b *Board) creatureUpdatesPerTick() {
	tempAliveCreatureObjects := b.AliveCreatureObjects

	for _, sourceCreature := range tempAliveCreatureObjects {
		status := sourceCreature.updateTick()

		if status == StatusMove {
			newPos := Pos{
				x: -1,
				y: -1,
			}

			newMove := MoveType{}

			for {
				newPos, newMove = b.newPosAndMove(sourceCreature)

				if newMove.action != AvoidAction {
					break
				}
			}

			switch newMove.action {
			case WaitAction:
				break

			case AttackAction:
				newMove.conflictinfo.commitConflict(b)

			case FoodAction:
				addMessageToCurrentGamelog("Food eaten by creature id: "+strconv.Itoa(sourceCreature.getId()), 2)
				b.deleteFood(newPos)
				b.moveCreature(sourceCreature, newPos, true)
				sourceCreature.heal(sourceCreature.getOriHP())

			case MoveAction:
				b.moveCreature(sourceCreature, newPos, true)

			default:
				addMessageToCurrentGamelog("Issue with MoveType newMove, please have a look.", 1)
			}

		} else if status == StatusDead {
			b.killCreature(sourceCreature, true)
		}
	}

	if b.checkIfCreaturesAreInactive() {
		if b.checkIfCreaturesAreDead() {
			b.GameOn = false
			addMessageToCurrentGamelog("All creatures are dead, end the game", 1)
		}

		b.newRound()
	}
}

func (b *Board) moveCreature(creature CreatureObject, newPos Pos, placeEmptyObject bool) {

	if placeEmptyObject {
		oldPos := creature.getPos()
		b.ObjectBoard[oldPos.y][oldPos.x] = newEmptyObject()
	}

	creature.setPos(newPos)
	b.ObjectBoard[newPos.y][newPos.x] = creature
}

func (b *Board) killCreature(creature CreatureObject, placeEmptyObject bool) {
	addMessageToCurrentGamelog("Killing creature: "+creature.getIdAsString()+
		" at pos: "+strconv.Itoa(creature.getPos().x)+" "+strconv.Itoa(creature.getPos().y), 2)
	pos := creature.getPos()

	if creature.getHP() > 0 {
		creature.kill()
	}

	b.deleteCreatureFromAliveSlice(creature)
	b.CurrentRound.CreaturesKilled = append(b.CurrentRound.CreaturesKilled, creature)

	if placeEmptyObject {
		b.ObjectBoard[pos.y][pos.x] = newEmptyObject()
	}
}

func (b *Board) newRound() {
	addMessageToCurrentGamelog("All creatures are dead or have eaten, starting new round", 1)
	b.spawnOffsprings()
	b.CurrentRound.CreaturesAliveAtEnd = b.AliveCreatureObjects
	b.writeSummaryOfRound()
	b.Gamelog.writeGamelogToFile()

	if len(b.Rounds) >= b.MaxRounds {
		b.GameOn = false

		addMessageToCurrentGamelog("Max number of rounds reached, ending the game.", 1)
		fmt.Println("Max number of rounds reached, ending the game.")

	} else {
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
	b.CurrentRound.CreaturesSpawnedSum = getSummary(b.CurrentRound.CreaturesSpawned)
	summaryStrings := getSummariesAsString(b.CurrentRound.CreaturesSpawnedSum, "spawned")

	for _, msg := range summaryStrings {
		addMessageToCurrentGamelog(msg, 1)
	}

	b.CurrentRound.CreaturesKilledSum = getSummary(b.CurrentRound.CreaturesKilled)
	summaryStrings = getSummariesAsString(b.CurrentRound.CreaturesKilledSum, "killed")

	for _, msg := range summaryStrings {
		addMessageToCurrentGamelog(msg, 1)
	}

	b.CurrentRound.CreaturesAliveAtEndSum = getSummary(b.CurrentRound.CreaturesAliveAtEnd)
	summaryStrings = getSummariesAsString(b.CurrentRound.CreaturesAliveAtEndSum, "alive")

	for _, msg := range summaryStrings {
		addMessageToCurrentGamelog(msg, 1)
	}

}

func getSummary(creatureList []CreatureObject) map[BoardObjectType]*CreatureSummary {
	returnCreatureSummary := make(map[BoardObjectType]*CreatureSummary)

	for _, creature := range creatureList {
		creatureType := creature.getBoardObjectType()

		if obj, ok := returnCreatureSummary[creatureType]; ok {
			obj.TotalCreatures += 1
			obj.TotalSpeed += creature.getSpeed()
			obj.TotalScanChance += creature.getScanProcChance()

		} else {
			newCreatureSummary := CreatureSummary{
				CreatureType:    creature.getType(),
				TotalCreatures:  1,
				TotalSpeed:      creature.getSpeed(),
				TotalScanChance: creature.getScanProcChance(),
			}

			returnCreatureSummary[creatureType] = &newCreatureSummary
		}
	}

	return returnCreatureSummary
}

func getSummariesAsString(summaries map[BoardObjectType]*CreatureSummary, action string) []string {
	returnString := []string{}

	for _, cs := range summaries {
		cs.AverageSpeed = float64(cs.TotalSpeed) / float64(cs.TotalCreatures)
		cs.AverageScanChance = cs.TotalScanChance / float64(cs.TotalCreatures)
		summary := "In last round, " + strconv.Itoa(cs.TotalCreatures) +
			" x " + cs.CreatureType + " was " + action + " with the average speed of: " +
			strconv.FormatFloat(cs.AverageSpeed, 'f', 2, 64) + " and average scan chance: " +
			strconv.FormatFloat(cs.AverageScanChance, 'f', 2, 64)

		returnString = append(returnString, summary)
	}

	return returnString
}

func (b *Board) findPosForAllCreatures() {
	addMessageToCurrentGamelog("Finding new pos for creatures", 2)
	for _, creature := range b.AliveCreatureObjects {
		creature.resetValues()

		for {
			newPos := b.randomPosAtEdgeOfMap()
			if b.isSpotEmpty(newPos) {
				b.moveCreature(creature, newPos, true)
				break
			}
		}
	}
}

func (b *Board) spawnOffsprings() {
	addMessageToCurrentGamelog("Spawning offsprings from last round", 1)

	// Refactor using enums
	creatureQty := map[string]uint{
		"creature1": 0,
		"creature2": 0,
	}

	tempAliveCreatureObjects := b.AliveCreatureObjects

	for _, obj := range tempAliveCreatureObjects {

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
			b.AliveCreatureObjects = append(b.AliveCreatureObjects, offspring)

			newPos := b.randomPosAtEdgeOfMap()

			for !b.isSpotEmpty(newPos) {
				newPos = b.randomPosAtEdgeOfMap()
			}

			b.moveCreature(offspring, newPos, false)
		}
	}

	for key, val := range creatureQty {
		if val > 0 {
			addMessageToCurrentGamelog(strconv.Itoa(int(val))+" Creatures of type "+key+" spawned", 1)
		}
	}

}

func (b *Board) checkIfCreaturesAreDead() bool {
	if len(b.AliveCreatureObjects) == 0 {
		return true
	}

	return false
}

func (b *Board) checkIfCreaturesAreInactive() bool {

	for _, obj := range b.AliveCreatureObjects {
		if !obj.isDead() && obj.isMoving() || obj.isDead() {
			return false
		}
	}

	return true
}

func (b *Board) newPosAndMove(creature CreatureObject) (Pos, MoveType) {
	// Check if food is nearby
	moveType := MoveType{
		action: NoAction,
	}

	chance := rand.Intn(100)

	if creature.getSpeed()+6 < float64(creature.getHP()) && chance < int(creature.getScanProcChance()) {
		creature.scan()
		foodFound, newPos := b.scanForFood(creature)
		if foodFound {
			moveType.action = FoodAction
			// addMessageToCurrentGamelog("FOOD FOUND AT: "+strconv.Itoa(newPos.x)+" "+strconv.Itoa(newPos.y)+"with moveaction: "+strconv.Itoa(int(moveType.action)), 1)
			return newPos, moveType
		}
	}

	currentPos := creature.getPos()
	newPos := Pos{-1, -1}

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

			if newPosType == FoodAtPos {
				moveType.action = FoodAction
				valid = true
				break

			} else if newPosType == EmptyPos {
				moveType.action = MoveAction
				valid = true
				break

			}

			if newPosType == ConflictAtPos {
				sourceCreature, err := b.getCreatureObjectFromBoard(currentPos)

				if err != nil {
					addMessageToCurrentGamelog(err.Error(), 1)
				}

				targetCreature, err := b.getCreatureObjectFromBoard(Pos{x: x, y: y})

				if err != nil {
					addMessageToCurrentGamelog(err.Error(), 1)
				}

				action, conflictinfo := b.ConflictManager.getConflict(sourceCreature, targetCreature)

				if action {
					valid = true
					moveType.action = AttackAction
					moveType.conflictinfo = conflictinfo
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

func (b *Board) deleteCreatureFromAliveSlice(creature CreatureObject) {
	var element int

	for i, aliveCreature := range b.AliveCreatureObjects {
		if aliveCreature == creature {
			element = i
			break
		}
	}

	b.AliveCreatureObjects = deleteIndexInCreatureSlice(b.AliveCreatureObjects, element)
}

func (b *Board) scanForFood(creature CreatureObject) (bool, Pos) {
	maxY := min(creature.getPos().y+1, b.Rows-1)
	maxX := min(creature.getPos().x+1, b.Cols-1)

	minY := max(creature.getPos().y-1, 0)
	minX := max(creature.getPos().x-1, 0)

	for minY <= maxY {
		for minX <= maxX {
			if _, ok := b.ObjectBoard[minY][minX].(*Food); ok {
				addMessageToCurrentGamelog(creature.getIdAsString()+" found food, by scanning for it, at: "+strconv.Itoa(minX)+" "+strconv.Itoa(minY), 1)
				return true, Pos{y: minY, x: minX}
			}
			minX++
		}
		minY++
	}

	return false, Pos{y: -1, x: -1}
}

func deleteIndexInPosSlice(posSlice []Pos, index int) []Pos {
	posSlice[index] = posSlice[len(posSlice)-1]
	return posSlice[:len(posSlice)-1]
}

func deleteIndexInCreatureSlice(creatureSlice []CreatureObject, index int) []CreatureObject {
	creatureSlice[index] = creatureSlice[len(creatureSlice)-1]
	return creatureSlice[:len(creatureSlice)-1]
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

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
