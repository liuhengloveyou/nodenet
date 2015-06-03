package nodenet_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

func TestNewMessage(t *testing.T) {
	msg, err := nodenet.NewMessage("demo", []string{"com1"}, "demo message")
	fmt.Println(msg, err)
}

func TestMsgGraph(t *testing.T) {
	msg, err := nodenet.NewMessage("demo", []string{"com1"}, "demo message")
	fmt.Println(msg, err)
	msg.PopGraph()
	fmt.Println(msg)
}
