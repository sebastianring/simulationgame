package main

import (
	"encoding/json"
	"os"
	// "strconv"
	"time"
)

const cols int = 40
const prioThreshold int = 1

var currentGamelog *Gamelog
var emptyMessage []byte

type Gamelog struct {
	cols              int
	rows              int
	messages          []message
	idCtr             int
	displayedMessages []string
	createdAt         time.Time
	fileString        string
}

type message struct {
	Id        int       `json:"id"`
	Prio      int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	Texts     string    `json:"text"`
}

func InitTextInfo(rows int) *Gamelog {
	gl := Gamelog{
		cols:       cols,
		rows:       rows,
		messages:   []message{},
		idCtr:      1,
		createdAt:  time.Now(),
		fileString: getFileString(),
	}

	for i := 0; i < cols; i++ {
		emptyMessage = append(emptyMessage, spaceSymbol)
	}

	currentGamelog = &gl

	return &gl
}

func getFileString() string {
	logsFolder := "logs"

	dir, err := os.Open(logsFolder)
	defer dir.Close()

	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(logsFolder, 0755)

			if err != nil {
				panic(err)
			}
			addMessageToCurrentGamelog("New folder created", 2)
		}
	}
	logsFolder = logsFolder + "/"

	logNamePrefix := "simulation_gamelog_"
	currentTime := time.Now()
	logNameSuffix := currentTime.Format("20060102150405") + ".txt"

	fullLogName := logsFolder + logNamePrefix + logNameSuffix

	return fullLogName
}

func addMessageToCurrentGamelog(msg string, prio int) {
	endSlice := 0
	msgLen := len(msg)

	var texts []string

	for i := 0; i <= msgLen; i = endSlice {
		endSlice = min(currentGamelog.rows+i, msgLen)
		// do we need to split the message?
		if endSlice < msgLen {
			for j := endSlice; j > i; j-- {
				// find the first space, backwards, and break there!
				if msg[j] == byte(32) {
					endSlice = j
					break
				}
			}
		}

		// log := fmt.Sprintf("%v %v %v", i, endSlice, msgLen)
		texts = append(texts, msg[i:endSlice])
		// currentGamelog.messages = append(currentGamelog.messages, log)
		endSlice++
	}

	newMessage := newMessage(currentGamelog.idCtr, prio, msg)
	currentGamelog.messages = append(currentGamelog.messages, *newMessage)

	if prio <= prioThreshold {
		for _, val := range texts {
			currentGamelog.displayedMessages = append(currentGamelog.displayedMessages, val)
		}
	}

	// if prio == 1 {
	// 	db, err := openDbConnection()
	//
	// 	if err != nil {
	// 		os.Exit(1)
	// 	}
	//
	// 	writeMessageToDb(db, newMessage)
	// }
	//
	currentGamelog.idCtr++
}

func newMessage(id int, prio int, msg string) *message {
	m := message{
		Id:        id,
		Prio:      prio,
		CreatedAt: time.Now(),
		Texts:     msg,
	}

	return &m
}

func (gl *Gamelog) getNumberOfMessageRows() int {
	rows := 0
	for _, msg := range gl.messages {
		rows += len(msg.Texts)
	}

	return rows
}

func (gl *Gamelog) getMessageByRow(row int) []byte {
	numberOfMessageRows := len(gl.displayedMessages)

	messageOffset := numberOfMessageRows - (gl.rows - row)

	if messageOffset > numberOfMessageRows-1 || messageOffset < 0 {
		return emptyMessage
	} else {
		spaces := cols - len(gl.displayedMessages[messageOffset])
		var spaceString string
		for i := 0; i < spaces; i++ {
			spaceString += " "
		}
		return []byte(gl.displayedMessages[messageOffset] + spaceString)
	}
}

func (gl *Gamelog) writeGamelogToFile() {
	file, err := os.OpenFile(gl.fileString, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	for _, message := range gl.messages {
		jsonData, err := json.Marshal(message)

		if err != nil {
			panic(err)
		}

		_, err = file.Write([]byte(jsonData))
		_, err = file.Write([]byte("\n"))
	}

	if err != nil {
		panic(err)
	}
}

func min(a int, b int) int {
	if a > b {
		return b
	}

	return a
}
