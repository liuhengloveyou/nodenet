package nodenet

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
)

type Message struct {
	ID          string                 `json:"id"`
	Entrance    string                 `json:"entrance"`
	Graph       []string               `json:"graph"`
	Context     map[string]interface{} `json:"context"`
	Payload     interface{}            `json:"payload"`
	DispenseKey string                 `json:"dispense"` // 均衡分发键
}

func NewMessage(id, entrance string, graphs []string, payload interface{}) (msg *Message, err error) {
	if id == "" {
		var u [16]byte
		if _, err = rand.Read(u[:]); err != nil {
			return
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
		Context:  make(map[string]interface{}),
		Payload:  payload}

	return
}

func (p *Message) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *Message) Marshal() ([]byte, error) {
	return json.Marshal(p)
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
