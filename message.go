package nodenet

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"
)

type Message struct {
	ID       string
	Entrance string
	Graph    []string
	Payload  interface{}
	Err      error

	DispenseKey string // 均衡分发键
}

var messageTypes map[string]interface{} = make(map[string]interface{})

func NewMessage(id, entrance string, graphs []string, payload interface{}) (msg *Message) {
	if id == "" {
		var u [16]byte
		if _, err := rand.Read(u[:]); err != nil {
			panic(err)
		}
		id = fmt.Sprintf("%x%x%x%x%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	}
	if graphs == nil {
		graphs = make([]string, 0)
	}

	msg = &Message{
		ID:       id,
		Entrance: strings.TrimSpace(entrance),
		Graph:    graphs,
		Payload:  payload}

	return
}

func RegisterMessageType(message interface{}) {
	if reflect.TypeOf(message).Kind() != reflect.Struct {
		panic("Only struct.")
	}

	gob.Register(message)

	messageTypes[reflect.TypeOf(message).String()] = message
}

func GetMessageTypeByName(name string) interface{} {
	return messageTypes[name]
}

func (p *Message) SetGraph(graph []string) {
	p.Graph = graph
}

func (p *Message) TopGraph() string {
	if len(p.Graph) >= 1 {
		return p.Graph[0]
	}

	return ""
}

func (p *Message) PopGraph() string {
	if len(p.Graph) >= 1 {
		p.Graph = p.Graph[1:]
	}

	return p.TopGraph()
}

func (p *Message) Decode(data []byte) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(p)
}

func (p *Message) Encode() (data []byte) {
	var buf bytes.Buffer
	if e := gob.NewEncoder(&buf).Encode(p); e != nil {
		panic(e)
	}

	return buf.Bytes()
}
