package nodenet_test

import (
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

var config = map[string]interface{}{"url": "127.0.0.1:12345", "timeout": 3}

func TestTcpServ(t *testing.T) {
	mq, err := nodenet.NewMQ("tcp", config)
	if err != nil {
		t.Error(err)
	}

	mq.StartService()

	msg, e := mq.GetMessage()
	t.Log(string(msg), e)

	return
}

func TestTcpClient(t *testing.T) {
	mq, err := nodenet.NewMQ("tcp", config)
	if err != nil {
		t.Error(err)
	}

	rst := mq.SendMessage([]byte("tcp test."))
	t.Log(rst)
}

func TestByte(t *testing.T) {
	buf := make([]byte, 16)
	var tmp []byte

	tmp = buf[6:9]
	t.Log(tmp, (len(tmp) == 0))

	tmp[0] = 'a'
	tmp[1] = 'b'
	t.Log(buf)

	for i := 0; i < len(tmp); i++ {
		buf[i] = tmp[i]
	}
	t.Log(buf)
}
