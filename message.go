package cloudnet

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
)

type Message struct {
	ID       string                 `json:"id"`
	Entrance string                 `json:"entrance"`
	Graph    []string               `json:"graph"`
	Context  map[string]interface{} `json:"context"`
	Payload  interface{}            `json:"payload"`
}

func NewComponentMessage(entrance string, payload interface{}) (msg *Message, err error) {
	msgID := ""
	if u, e := uuid.NewV4(); e != nil {
		err = e
		return
	} else {
		msgID = u.String()
	}

	msg = &ComponentMessage{
		ID:       msgID,
		Entrance: strings.TrimSpace(entrance),
		Graph:    nil,
		Payload:  payload}

	return
}
