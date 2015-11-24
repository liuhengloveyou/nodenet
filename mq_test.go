package nodenet_test

import (
	"encoding/gob"
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

var config = map[string]interface{}{"url": "127.0.0.1:12345", "timeout": 3}

func TestTcpServ(t *testing.T) {
	gob.Register(MsgS{})
	mq, err := nodenet.NewMQ("tcp", config)
	if err != nil {
		t.Error(err)
	}

	mq.StartService()

	msg, e := mq.GetMessage()
	t.Log(string(msg), e)

	msg1, e := mq.GetMessage()
	t.Log(string(msg1), e)

	msg2, e := mq.GetMessage()
	t.Log(string(msg2), e)

	msg3, e := mq.GetMessage()
	t.Log(string(msg3), e)

	imsg := nodenet.NewMessage("", "", []string{}, "")
	imsg.Decode([]byte(msg3))
	t.Log(imsg)

	return
}

func TestTcpClient(t *testing.T) {
	gob.Register(MsgS{})
	mq, err := nodenet.NewMQ("tcp", config)
	if err != nil {
		t.Error(err)
	}

	rst := mq.SendMessage([]byte("22222222222222222222222222222222222222222222222222222222222222222222222222222222222."))
	t.Log(rst)

	rst = mq.SendMessage([]byte("11111111111."))
	t.Log(rst)

	rst = mq.SendMessage([]byte("333."))
	t.Log(rst)

	msg := nodenet.NewMessage("", "demo", []string{"com1", "g2"}, &MsgS{Name: "aaa", Age: 18})
	rst = mq.SendMessage(msg.Encode())
	t.Log(rst)
}

func TestByte(t *testing.T) {
	buf := make([]byte, 16)
	var tmp []byte = buf[2:5]

	tmp[0] = 'a'
	tmp[1] = 'b'
	tmp[2] = 'c'
	t.Log(buf)

	copy(buf, tmp)
	t.Log(buf)

	for i := 0; i < len(tmp); i++ {
		buf[i] = tmp[i]
	}
	t.Log(buf)
}
