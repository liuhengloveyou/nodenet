/*
消息接口
*/
package cloudnet

import (
	"fmt"
)

var mqs map[string]mqType = make(map[string]mqType)

// 消息接口
type MessageQueue interface {
	Run() error                   // 启动服务监听
	RecvMessage() ([]byte, error) // 读一条消息
	SendMessage([]byte) error     // 发送一条消息到给定节点
}

type mqType func(interface{}) (MessageQueue, error)

func RegisterMq(name string, one mqType) {
	if one == nil {
		panic("Register MQ nil")
	}

	if _, dup := mqs[name]; dup {
		panic("Register MQ duplicate for " + name)
	}

	mqs[name] = one
}

func NewMQ(typeName string, config interface{}) (mq MessageQueue, err error) {
	if newFunc, ok := mqs[typeName]; ok {
		return newFunc(config)
	}

	return nil, fmt.Errorf("No MQ types " + typeName)
}
