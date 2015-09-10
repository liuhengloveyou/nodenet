package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/liuhengloveyou/nodenet"
)

func main() {
	flag.Parse()

	// 从配置文件初始华
	nodenet.BuildFromConfig("nodenet.conf.example")

	// 找到名字为com1的组件, 把它的业务处理函数设置成自己实现的函数: work
	nodenet.GetComponentByName("com1").SetHandler(work1)
	nodenet.GetComponentByName("com2").SetHandler(work2)

	//启动这个组件
	go nodenet.GetComponentByName("com1").Run()
	go nodenet.GetComponentByName("com2").Run()

	// 新建一个消息
	msg, _ := nodenet.NewMessage("com2", []string{"com1"}, "demo message")

	// 把消息发送给com1组件
	nodenet.SendMsgToComponent("com1", msg)

	for {
		time.Sleep(2 * time.Second)
	}
}

// 当组件收到一条消息的时候, 会调用work函数处理
func work1(msg interface{}) (result interface{}, err error) {
	fmt.Println("work1>>>>>>", msg)

	return msg, nil
}

func work2(msg interface{}) (result interface{}, err error) {
	fmt.Println("work2>>>>>>", msg)

	return nil, nil
}
