package nodenet_test

import (
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

func TestNewMessage(t *testing.T) {
	msg, err := nodenet.NewMessage("demo", []string{"com1"}, "demo message")
	t.Log(msg, err)
}

func TestMsgGraph(t *testing.T) {
	msg, err := nodenet.NewMessage("demo", []string{"com1"}, "demo message")
	t.Log(msg, err)
	msg.PopGraph()
	t.Log(msg)
}
