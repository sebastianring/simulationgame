package simulationgame

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

var currentGamelog *Gamelog
var emptyMessage []byte

type Gamelog struct {
	cols              int
	rows              int
	messages          []*Message
	idCtr             int
	displayedMessages []string
	createdAt         time.Time
	fileString        string
	prioThreshold     int
}

type Message struct {
	Id        int       `json:"id"`
	Prio      int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	Texts     string    `json:"text"`
}

func NewGamelog(rows int, cols int) *Gamelog {
	gl := Gamelog{
		cols:          cols,
		rows:          rows,
		messages:      []*Message{},
		idCtr:         1,
		createdAt:     time.Now(),
		fileString:    getFileString(),
		prioThreshold: 1,
	}

	for i := 0; i < cols; i++ {
		emptyMessage = append(emptyMessage, byte(32))
	}

	currentGamelog = &gl

	return &gl
}

func getFileString() string {
	logsFolder := "logs"

	_, err := os.Stat(logsFolder)

	if os.IsNotExist(err) {
		err = os.Mkdir(logsFolder, 0755)

		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	logsFolder = logsFolder + "/"

	logNamePrefix := "simulation_gamelog_"
	currentTime := time.Now()
	logNameSuffix := currentTime.Format("20060102150405") + ".txt"

	fullLogName := logsFolder + logNamePrefix + logNameSuffix

	return fullLogName
}

// Need to adapt log to be able to receive any type of values and convert them to string automatically
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

		texts = append(texts, msg[i:endSlice])
		endSlice++
	}

	newMessage := newMessage(currentGamelog.idCtr, prio, msg)
	currentGamelog.messages = append(currentGamelog.messages, newMessage)

	if prio <= currentGamelog.prioThreshold {
		for _, val := range texts {
			currentGamelog.displayedMessages = append(currentGamelog.displayedMessages, val)
		}
	}

	currentGamelog.idCtr++
}

func newMessage(id int, prio int, msg string) *Message {
	m := Message{
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
		spaces := gl.cols - len(gl.displayedMessages[messageOffset])
		var spaceString string

		for i := 0; i < spaces; i++ {
			spaceString += " "
		}

		return []byte(gl.displayedMessages[messageOffset] + spaceString)
	}
}

func (gl *Gamelog) writeGamelogToFile() error {
	file, err := os.OpenFile(gl.fileString, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return errors.New("Error opening file: " + err.Error())
	}

	defer file.Close()

	for _, message := range gl.messages {
		jsonData, err := json.Marshal(message)

		if err != nil {
			return errors.New("Error marshaling message")
		}

		_, err = file.Write([]byte(jsonData))

		if err != nil {
			return errors.New("Error writing to file: " + err.Error())
		}

		_, err = file.Write([]byte("\n"))

		if err != nil {
			return errors.New("Error writing to file: " + err.Error())
		}
	}

	return nil
}

func min(a int, b int) int {
	if a > b {
		return b
	}

	return a
}
