package nodenet

import (
	"encoding/json"
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

func NewMessage(entrance string, graphs []string, payload interface{}) (msg *Message, err error) {
	msgID := ""
	if u, e := uuid.NewV4(); e != nil {
		err = e
		return
	} else {
		msgID = u.String()
	}

	msg = &Message{
		ID:       msgID,
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
