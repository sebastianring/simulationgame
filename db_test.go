package simulationgame_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
	sg "github.com/sebastianring/simulationgame"
)

func TestWriteMessageToDB(t *testing.T) {
	t.Setenv("SIM_GAME_DB_PW", os.Getenv("SIM_GAME_DB_PW"))
	_ = sg.NewGamelog(10, 10)
	db, err := sg.OpenDbConnection()

	if err != nil {
		t.Fatal("Failed conneting to DB: " + err.Error())
	}

	err = db.Ping()

	if err != nil {
		t.Fatal("Can't ping db...")
	}

	randint := rand.Intn(1000)
	uuid, err := uuid.NewRandom()

	if err != nil {
		t.Fatal("Failed to generate uuid: " + err.Error())
	}

	res := 0

	query := "INSERT INTO simulation_game.messages (id, prio, text, board_link) VALUES ($1, $2, $3, $4) RETURNiNG id"
	err = db.QueryRow(query, randint, 1, "HELLO THIS IS A TEST", uuid).Scan(&res)

	if err != nil {
		t.Fatal("Failed writing to db: " + err.Error())
	}

	fmt.Print(res)
}
