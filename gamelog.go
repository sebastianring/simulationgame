package main

import (
	// "io/ioutil"
	"os"
	"strconv"
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
	id        int
	prio      int
	createdAt time.Time
	texts     string
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
	logsFolder := "logs/"
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

	currentGamelog.idCtr++
}

func newMessage(id int, prio int, msg string) *message {
	m := message{
		id:        id,
		prio:      prio,
		createdAt: time.Now(),
		texts:     msg,
	}

	return &m
}

func (gl *Gamelog) getNumberOfMessageRows() int {
	rows := 0
	for _, msg := range gl.messages {
		rows += len(msg.texts)
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

	// testcontent := []byte("Testing" + gl.createdAt.Format("20060102150405") + "\n")
	// _, err = file.Write(testcontent)

	for _, message := range gl.messages {
		id := ("MESSAGE ID: " + strconv.Itoa(message.id) + "\n")
		prio := ("PRIORITY: " + strconv.Itoa(message.prio) + "\n")
		createdAt := ("CREATED AT: " + message.createdAt.Format("20060102150405") + "\n")
		text := ("MESSAGE: " + message.texts)

		log := []byte(id + prio + createdAt + text + "\n" + "\n")

		_, err = file.Write(log)
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
