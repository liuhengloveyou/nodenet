/*
go.rpc实现的mq
*/
package nodenet

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

var msgChan chan []byte

type MqRpc struct {
	Url string `json:"url"`
}

func (p *MqRpc) config(conf []byte) error {
	return json.Unmarshal(conf, p)
}

func (p *MqRpc) Ready() error {
	err := rpc.Register(p)
	if err != nil {
		return err
	}

	rpc.HandleHTTP()

	l, e := net.Listen("tcp", p.Url)
	if e != nil {
		log.Fatal("Listen error:", e)
	}

	msgChan = make(chan []byte)
	go http.Serve(l, nil)
	return nil
}

func (p *MqRpc) RecvMessage() (msg []byte, err error) {
	return <-msgChan, nil
}

func (p *MqRpc) SendMessage(msg []byte) error {
	client, err := rpc.DialHTTP("tcp", p.Url)
	if err != nil {
		return fmt.Errorf("DialHttp ERR: %s", err.Error())
	}

	rst := 0
	err = client.Call("MqRpc.Recv", msg, &rst)
	if err != nil {
		return fmt.Errorf("Rpc ERR: %s", err.Error())
	}

	return nil
}

func (p *MqRpc) Recv(msg []byte, rst *int) error {
	msgChan <- msg
	*rst = 0
	return nil
}

func init() {
	RegisterMq("rpc", NewMqRpc)
}

func NewMqRpc(config []byte) (MessageQueue, error) {
	tmq := &MqRpc{}
	if err := tmq.config(config); err != nil {
		return nil, err
	}

	return tmq, nil
}
