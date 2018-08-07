package infobip

import (
	"bytes"
	"encoding/json"
)

// SMS is a message that will be sent.
// "From" field should be alphanumeric sender ID, the length of data should be between 3 and 11 characters.
// "To" is an array of message destination addresses in international format.
// "Text" is a message body.
type SMS struct {
	From string   `json:"from"` // todo: validate length
	To   []string `json:"to"`
	Text string   `json:"text"`
}

func (s *SMS) buffer() *bytes.Buffer {
	b, _ := json.Marshal(s)
	return bytes.NewBuffer(b)
}
