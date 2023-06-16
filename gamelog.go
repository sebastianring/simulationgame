package main

// import "fmt"

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
}

type message struct {
	id    int
	prio  int
	texts []string
}

func InitTextInfo(rows int) *Gamelog {
	gl := Gamelog{
		cols:     cols,
		rows:     rows,
		messages: []message{},
		idCtr:    1,
	}

	for i := 0; i < cols; i++ {
		emptyMessage = append(emptyMessage, spaceSymbol)
	}

	currentGamelog = &gl

	return &gl
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

	newMessage := newMessage(currentGamelog.idCtr, prio, texts)
	currentGamelog.messages = append(currentGamelog.messages, *newMessage)

	if prio <= prioThreshold {
		for _, val := range texts {
			currentGamelog.displayedMessages = append(currentGamelog.displayedMessages, val)
		}
	}

	currentGamelog.idCtr++
}

func newMessage(id int, prio int, texts []string) *message {
	m := message{
		id:    id,
		prio:  prio,
		texts: texts,
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

func min(a int, b int) int {
	if a > b {
		return b
	}

	return a
}
