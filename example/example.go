package main

import (
	"flag"
	"time"

	log "github.com/golang/glog"
	"github.com/liuhengloveyou/nodenet"
)

// 第1步: 定义业务数据结构
type P struct {
	Name string
	Age  int
	Sex  bool
}

type W struct {
	Work  string
	Addr  int
	Start time.Time
}

// 第2步: 定义业处理函数
func work1(msg interface{}) (result interface{}, err error) {
	p := msg.(P)

	log.Infoln("work1>>>", p.Name, p.Age, p.Sex)

	return msg, nil
}

func work2(msg interface{}) (result interface{}, err error) {
	w := msg.(W)

	log.Infoln("work2>>>", w.Work, w.Addr, w.Start)

	return nil, nil // 不会再往下发, 只能入口才应该返回nil.
}

// 第3步: 注册本组件可以处理的消息和相应的处理函数
func main() {
	flag.Parse()

	// 从配置文件初始华
	nodenet.BuildFromConfig("nodenet.conf.example")

	// 找到名字为com1的组件, 注册组件可以处理的消息和相应的处理函数
	nodenet.GetComponentByName("com1").RegisterHandler(P{}, work1)
	nodenet.GetComponentByName("com2").RegisterHandler(W{}, work2)
	nodenet.GetComponentByName("com2").RegisterHandler(P{}, work1)

	//启动这个组件
	go nodenet.GetComponentByName("com1").Run()
	go nodenet.GetComponentByName("com2").Run()

	// 新建一个消息

	msg := nodenet.NewMessage("", "com2", nodenet.GetGraphByName("demo"), &P{Name: "aaa", Age: 18, Sex: false})

	// 把消息发送出去
	nodenet.SendMsgToNext(msg)

	for {
		time.Sleep(2 * time.Second)
	}
}
