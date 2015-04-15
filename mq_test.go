package cloudnet_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/cloudnet"
)

func TestServ(t *testing.T) {
	mq, err := cloudnet.NewMQ("rpc", ":1234")
	if err != nil {
		fmt.Println(err)
	}

	err = mq.ListenAndServe()

	msg, e := mq.RecvMessage()
	fmt.Println(msg, e)
}

func TestClient(t *testing.T) {
	mq, err := cloudnet.NewMQ("rpc", ":1234")
	if err != nil {
		fmt.Println(err)
	}

	rst := mq.SendMessage([]byte("rpc test."))
	fmt.Println(rst)
}
