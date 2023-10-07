package simulationgame

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

func OpenDbConnection() (*sql.DB, error) {
	prefix := "postgres://"
	user := "sim_game"
	adress := "5.150.233.156"
	// adress := "192.168.0.130"
	// adress := "localhost"
	port := "5432"

	password := os.Getenv("SIM_GAME_DB_PW")

	if len(password) == 0 {
		return nil, errors.New("Missing password, will not connect to db.")
	}

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

func writeMessageToDb(b *Board, msg *Message) error {
	db, err := OpenDbConnection()

	if err != nil {
		return err
	}

	query := "INSERT INTO simulation_game.messages (id, prio, text, board_link) VALUES ($1, $2, $3, $4) RETURNiNG id"
	err = db.QueryRow(query, msg.Id, msg.Prio, msg.Texts, b.Id).Scan(&msg.Id)

	if err != nil {
		return err
	}

	return nil
}

func writeMessagesToDb(b *Board) error {
	wg := sync.WaitGroup{}
	db, err := OpenDbConnection()

	if err != nil {
		return errors.New("Error connecting to db: " + err.Error())
	}

	db.SetMaxOpenConns(50)

	for _, msg := range b.Gamelog.messages {
		wg.Add(1)

		go func(msg *Message) {
			query := "INSERT INTO simulation_game.messages (id, prio, text, board_link) VALUES ($1, $2, $3, $4) RETURNiNG id"
			_, err := db.Exec(query, msg.Id, msg.Prio, msg.Texts, b.Id)

			if err != nil {
				log.Println("Error writing to db: " + err.Error())
			}

			wg.Done()
		}(msg)
	}

	wg.Wait()
	db.Close()

	return nil
}

func writeBoardToDb(board *Board) error {
	db, err := OpenDbConnection()

	if err != nil {
		return errors.New("Error connecting to db: " + err.Error())
	}

	query := "INSERT INTO simulation_game.boards (id, rows, cols) VALUES ($1, $2, $3) RETURNiNG id"
	err = db.QueryRow(query, board.Id, board.Rows, board.Cols).Scan(&board.Id)

	if err != nil {
		return errors.New("Error writing to db: " + err.Error())
	}

	return nil
}

// need at add round to DB
func writeRoundToDb(db *sql.DB, round *Round) {

}
