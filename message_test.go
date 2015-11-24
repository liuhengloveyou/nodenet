package nodenet_test

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

type MsgS struct {
	Name string
	Age  int32
}

func TestNewMessage(t *testing.T) {
	gob.Register(MsgS{})

	msg := nodenet.NewMessage("", "demo", []string{"com1", "g2"}, &MsgS{Name: "aaa", Age: 18})
	t.Log(msg)

	b := msg.Encode()
	t.Log(string(b))

	msg1 := nodenet.NewMessage("", "", nil, "")
	msg1.Decode(b)
	t.Log(msg1, reflect.TypeOf(msg1), reflect.TypeOf(msg1.Payload))

}

func TestMsgGraph(t *testing.T) {
	msg := nodenet.NewMessage("", "demo", []string{"com1"}, "demo message")
	t.Log(msg)
	msg.PopGraph()
	t.Log(msg)
}

func TestReflect(t *testing.T) {

	mt := reflect.TypeOf(&nodenet.Message{})
	et := mt.Elem()

	t.Log(mt, et)

}

func TestGob(t *testing.T) {
	type P struct {
		X, Y, Z int
		Name    string
	}

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.

	// Encode (send) some values.
	err := enc.Encode(P{3, 4, 5, "Pythagoras"})
	if err != nil {
		t.Log("encode error:", err)
	}
	err = enc.Encode(P{1782, 1841, 1922, "Treehouse"})
	if err != nil {
		t.Log("encode error:", err)
	}

	t.Log(string(network.Bytes()))
}
