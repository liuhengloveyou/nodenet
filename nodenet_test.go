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
	// 从配置文件初始华
	nodenet.BuildFromConfig("example/nodenet.conf.example")

	// 找到名字为com1的组件, 把它的业务处理函数设置成自己实现的函数: work
	nodenet.GetComponentByName("com1").SetHandler(work)

	//启动这个组件
	go nodenet.GetComponentByName("com1").Run()

	// 新建一个消息
	msg, _ := nodenet.NewMessage("demo", []string{"com1"}, "demo message")
	t.Log(msg)

	// 把消息发送给com1组件
	err := nodenet.SendMsgToNext("com1", msg)
	t.Log(err)
}

// 当组件收到一条消息的时候, 会调用work函数处理
func work(msg interface{}) (result interface{}, err error) {
	fmt.Println(">>>>>>", msg)

	return nil, nil
}
