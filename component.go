/*
节点组件
*/
package cloudnet

import (
	"fmt"
	"log"
	"strings"
)

var components map[string]*Component = make(map[string]*Component)

// 业务处理函数
type ComponentHandler func(interface{}) (result interface{}, err error)

// 端点
type EndPoint struct {
	MQType string
	conf   interface{}
	mq     MessageQueue
}

// 组件
type Component struct {
	Name    string               // 组件名
	in      EndPoint             // 接收消息端点
	outs    map[string]*EndPoint // 所有下流端点
	handler ComponentHandler     // 业务处理函数
}

func NewComponent(name, mqtype string, inconf interface{}) (*Component, error) {
	sname, smqtype := strings.TrimSpace(name), strings.TrimSpace(mqtype)
	if sname == "" {
		return nil, fmt.Errorf("Component's name empty ERR")
	}
	if smqtype == "" {
		return nil, fmt.Errorf("Component's MQ type ERR")
	}

	com := &Component{
		Name:    sname,
		in:      EndPoint{MQType: smqtype, conf: inconf, mq: nil},
		outs:    make(map[string]*EndPoint),
		handler: nil}

	components[com.Name] = com

	return com, nil
}

func GetComponentByName(name string) *Component {
	if component, ok := components[name]; ok {
		return component
	}

	return nil
}

func (p *Component) SetHandler(handler ComponentHandler) {
	if handler == nil {
		panic("Set Handler nil")
	}

	p.handler = handler
}

func (p *Component) Run() (err error) {
	log.Println("Component Run...", p.Name)

	// 创建接收MQ
	p.in.mq, err = NewMQ(p.in.MQType, p.in.conf)
	if err != nil {
		return
	}

	// MQ 准备
	err = p.in.mq.Run()
	if err != nil {
		return
	}

	// 开始监听
	p.recvMonitor()

	return nil
}

func (p *Component) recvMonitor() {
	for {
		msg, err := p.in.mq.RecvMessage()
		if err != nil {
			log.Println(p.Name, "Error receiving message:", err.Error())
			continue
		}

		log.Println(string(msg))
	}
}
