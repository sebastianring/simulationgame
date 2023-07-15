package main

import (
	"database/sql"
	// _ "database/sql/driver"
	// "fmt"
	_ "github.com/lib/pq"
	"os"
	// _ "github.com/lib/pq"
)

func openDb() {
	db, err := sql.Open("postgres", "postgres://sim_game:valmet865@localhost:5432/postgres")

	if err != nil {
		addMessageToCurrentGamelog(err.Error(), 1)
		os.Exit(1)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		addMessageToCurrentGamelog(err.Error(), 1)
	}
}
