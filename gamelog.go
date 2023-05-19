package main

// import "fmt"

const cols int = 40

var currentGamelog Gamelog

var emptyMessage []byte

type Gamelog struct {
	cols     int
	rows     int
	messages []string
}

func InitTextInfo(rows int) *Gamelog {
	gl := Gamelog{
		cols:     cols,
		rows:     rows,
		messages: []string{"Logging started"},
	}

	for i := 0; i < cols; i++ {
		emptyMessage = append(emptyMessage, spaceSymbol)
	}

	currentGamelog = gl

	return &gl
}

func (gl *Gamelog) addMessage(msg string) {
	gl.messages = append(gl.messages, msg)
}

func addMessageToCurrentGamelog(msg string) {
	currentGamelog.messages = append(currentGamelog.messages, msg)
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
