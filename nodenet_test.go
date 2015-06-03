package nodenet_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/nodenet"
)

func TestBuildFromConfig(t *testing.T) {
	nodenet.BuildFromConfig("example/nodenet.conf.example")
}

func TestNodenet(t *testing.T) {
	nodenet.BuildFromConfig("example/nodenet.conf.example")

	nodenet.GetComponentByName("com1").SetHandler(work)
	go nodenet.GetComponentByName("com1").Run()

	msg, _ := nodenet.NewMessage("demo", []string{"com1"}, "demo message")
	fmt.Println(msg)

	err := nodenet.SendMsgToNext("com1", msg)
	fmt.Println(err)
}

func work(msg interface{}) (result interface{}, err error) {
	fmt.Println(">>>>>>", msg)

	return nil, nil
}
