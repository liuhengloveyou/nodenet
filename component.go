/*
节点组件
*/
package cloudnet

import (
	"log"
)

var components map[string]*Component = make(map[string]*Component)

// 业务处理函数
type ComponentHandler func(*Payload) (result interface{}, err error)

// 端点
type EndPoint struct {
	MQType string
	Conf   interface{}
	mq     *MessageQueue
}

// 组件
type Component struct {
	Name    string               // 组件名
	in      EndPoint             // 接收消息端点
	outs    map[string]*EndPoint // 所有下流端点
	handler ComponentHandler     // 业务处理函数
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
	log.Infoln("Component Run...", p.Name)

	// 创建接收MQ
	p.in.mq, err = NewMq(p.in.MQType, p.in.Conf)
	if err != nil {
		return
	}

	// MQ 准备
	err = p.in.mq.ListenAndServe()
	if err != nil {
		return
	}

	// 开始监听
	go p.recvMonitor()

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
