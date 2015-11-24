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

var messageTypes map[string]reflect.Type = make(map[string]reflect.Type)

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

func RegisterMessageType(value interface{}) {
	gob.Register(value)
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
