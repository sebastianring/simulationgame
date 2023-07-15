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
	adress := "192.168.0.130"
	port := "5432"

	database_url := prefix + user + ":" +
		password + "@" + adress + ":" +
		port + "/postgres"

	addMessageToCurrentGamelog(database_url, 1)

	db, err := sql.Open("postgres", database_url)

	if err != nil {
		addMessageToCurrentGamelog(err.Error(), 1)
		return nil, err
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		addMessageToCurrentGamelog(err.Error(), 1)
		return nil, err
	}

	return db, nil
}

func writeMessageToDb(db *sql.DB, msg *message) error {

	return nil
}
