package infobip

import (
	"bytes"
	"encoding/json"
)

type SMS struct {
	From string   `json:"from"`
	To   []string `json:"to"`
	Text string   `json:"text"`
}

func (s *SMS) Buffer() *bytes.Buffer {
	b, _ := json.Marshal(s)
	return bytes.NewBuffer(b)
}
