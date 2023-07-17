package main

import (
	"database/sql"
	// _ "database/sql/driver"
	// "fmt"
	_ "github.com/lib/pq"
	"os"
	// _ "github.com/lib/pq"
)

func openDbConnection() (*sql.DB, error) {
	prefix := "postgres://"
	user := "sim_game"
	password := os.Getenv("SIM_GAME_DB_PW")
	// adress := "192.168.0.130"
	adress := "localhost"
	port := "5432"

	database_url := prefix + user + ":" +
		password + "@" + adress + ":" +
		port + "/postgres"

	// addMessageToCurrentGamelog(database_url, 1)

	db, err := sql.Open("postgres", database_url)

	if err != nil {
		addMessageToCurrentGamelog(err.Error(), 1)
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		addMessageToCurrentGamelog(err.Error(), 1)
		return nil, err
	}

	return db, nil
}

func writeMessageToDb(db *sql.DB, msg *message) {
	go func() {
		const SCHEMA = "simulation_game."
		query := "INSERT INTO simulation_game.message (id, prio, text) VALUES ($1, $2, $3) RETURNiNG id"
		err := db.QueryRow(query, msg.Id, msg.Prio, msg.Texts).Scan(&msg.Id)

		if err != nil {
			addMessageToCurrentGamelog(err.Error(), 1)
		}

	}()
}
