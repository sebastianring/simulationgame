package main

import (
	"database/sql"
	"fmt"
	// _ "database/sql/driver"
	// "fmt"
	"os"

	_ "github.com/lib/pq"
	// _ "github.com/lib/pq"
)

func openDbConnection() (*sql.DB, error) {
	prefix := "postgres://"
	user := "sim_game"
	password := os.Getenv("SIM_GAME_DB_PW")
	adress := "192.168.0.130"
	// adress := "localhost"
	port := "5432"

	database_url := prefix + user + ":" +
		password + "@" + adress + ":" +
		port + "/postgres"

	db, err := sql.Open("postgres", database_url)

	if err != nil {
		addMessageToCurrentGamelog(err.Error(), 1)
		return nil, err
	} else {
		addMessageToCurrentGamelog("Database connection secured!", 1)
	}

	err = db.Ping()

	if err != nil {
		addMessageToCurrentGamelog(err.Error(), 1)
		return nil, err
	} else {
		addMessageToCurrentGamelog("Database ping succesful!", 1)
	}

	return db, nil
}

func writeMessageToDb(db *sql.DB, msg *message) {
	go func() {
		query := "INSERT INTO simulation_game.messages (id, prio, text, board_link) VALUES ($1, $2, $3, $4) RETURNiNG id"
		err := db.QueryRow(query, msg.Id, msg.Prio, msg.Texts, currentBoardId).Scan(&msg.Id)

		if err != nil {
			fmt.Println(err.Error())
			// addMessageToCurrentGamelog(err.Error(), 1)
		}

	}()
}

func writeBoardToDb(db *sql.DB, board *Board) {
	go func() {
		query := "INSERT INTO simulation_game.boards (id, rows, cols) VALUES ($1, $2, $3) RETURNiNG id"
		err := db.QueryRow(query, board.Id, board.rows, board.cols).Scan(&board.Id)

		if err != nil {
			fmt.Println(err.Error())
			addMessageToCurrentGamelog(err.Error(), 1)
		}

	}()
}

func writeRoundToDb(db *sql.DB, round *Round) {

}
