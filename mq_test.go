package nodenet_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

const config = `{"url":"127.0.0.1:12345"}`

func TestRpcServ(t *testing.T) {
	mq, err := nodenet.NewMQ("rpc", []byte(config))
	if err != nil {
		fmt.Println(err)
	}

	err = mq.Ready()

	msg, e := mq.RecvMessage()
	fmt.Println(string(msg), e)
}

func TestRpcClient(t *testing.T) {
	mq, err := nodenet.NewMQ("rpc", []byte(config))
	if err != nil {
		fmt.Println(err)
	}

	rst := mq.SendMessage([]byte("rpc test."))
	fmt.Println(rst)
}
