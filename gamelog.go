package main

// import "fmt"

const cols int = 40

const prioThreshold int = 1

var currentGamelog *Gamelog

var emptyMessage []byte

type Gamelog struct {
	cols           int
	rows           int
	messages       []string
	lowPrioMessage []string
}

func InitTextInfo(rows int) *Gamelog {
	gl := Gamelog{
		cols:           cols,
		rows:           rows,
		messages:       []string{"Logging started"},
		lowPrioMessage: []string{},
	}

	for i := 0; i < cols; i++ {
		emptyMessage = append(emptyMessage, spaceSymbol)
	}

	currentGamelog = &gl

	return &gl
}

func addMessageToCurrentGamelog(msg string) {
	endSlice := 0
	msgLen := len(msg)

	for i := 0; i <= msgLen; i = endSlice {
		endSlice = min(currentGamelog.rows+i, msgLen)
		// do we need to split the message?
		if endSlice < msgLen {
			for j := endSlice; j > i; j-- {
				// find the first space and break there!
				if msg[j] == byte(32) {
					endSlice = j
					break
				}
			}
		}

		// log := fmt.Sprintf("%v %v %v", i, endSlice, msgLen)
		currentGamelog.messages = append(currentGamelog.messages, msg[i:endSlice])
		// currentGamelog.messages = append(currentGamelog.messages, log)
		endSlice++
	}
}

func (gl *Gamelog) getMessageByRow(row int) []byte {
	messageOffset := len(gl.messages) - (gl.rows - row)

	if messageOffset > len(gl.messages)-1 || messageOffset < 0 {
		return emptyMessage
	} else {
		spaces := cols - len(gl.messages[messageOffset])
		var spaceString string
		for i := 0; i < spaces; i++ {
			spaceString += " "
		}
		return []byte(gl.messages[messageOffset] + spaceString)
	}
}

func min(a int, b int) int {
	if a > b {
		return b
	}

	return a
}
