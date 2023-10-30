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
		prioThreshold: 1,
	}

	for i := 0; i < cols; i++ {
		emptyMessage = append(emptyMessage, byte(32))
	}

	currentGamelog = &gl

	return &gl
}

func getFileString() (string, error) {
	logsFolder := "logs"

	_, err := os.Stat(logsFolder)

	if os.IsNotExist(err) {
		err = os.Mkdir(logsFolder, 0755)

		if err != nil {
			return "", errors.New("Issue creating folder: " + err.Error())
		}
	} else if err != nil {
		return "", errors.New("Issue accessing folder: " + err.Error())
	}

	logsFolder = logsFolder + "/"

	logNamePrefix := "simulation_gamelog_"
	currentTime := time.Now()
	logNameSuffix := currentTime.Format("20060102150405") + ".txt"

	fullLogName := logsFolder + logNamePrefix + logNameSuffix

	return fullLogName, nil
}

// Need to adapt log to be able to receive any type of values and convert them to string automatically
func addMessageToCurrentGamelog(msg string, prio int) {
<<<<<<< HEAD
	endSlice := 0
	msgLen := len(msg)

	var texts []string

	for i := 0; i <= msgLen; i = endSlice {
		endSlice = min(currentGamelog.cols+i, msgLen)
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

=======
>>>>>>> ab92c6c64c640c6d0ea61914f6f53e9c07be038a
	newMessage := newMessage(currentGamelog.idCtr, prio, msg)
	currentGamelog.messages = append(currentGamelog.messages, newMessage)
	currentGamelog.idCtr++

	if prio <= currentGamelog.prioThreshold {
		endSlice := 0
		msgLen := len(msg)

		for i := 0; i <= msgLen; i = endSlice {
			endSlice = min(currentGamelog.cols+i, msgLen)
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

			currentGamelog.displayedMessages = append(currentGamelog.displayedMessages, msg[i:endSlice])
			endSlice++
		}
	}
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
	folder, err := getFileString()

	if err != nil {
		return err
	}

	file, err := os.OpenFile(folder, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

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
