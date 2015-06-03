/*
节点组件
*/
package nodenet

import (
	"fmt"
	"strings"

	log "github.com/golang/glog"
)

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
	Name    string           // 组件名
	Group   string           // 组件所属的组
	in      EndPoint         // 接收消息端点
	handler ComponentHandler // 业务处理函数
}

func NewComponent(name, intype string, inconf interface{}) (*Component, error) {
	sname, sintype := strings.TrimSpace(name), strings.TrimSpace(intype)
	if sname == "" {
		return nil, fmt.Errorf("Component's name empty")
	}
	if sintype == "" {
		return nil, fmt.Errorf("Component's type empty")
	}

	components[sname] = &Component{
		Name:    sname,
		in:      EndPoint{MQType: sintype, conf: inconf, mq: nil},
		handler: nil}

	return components[sname], nil
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
	p.in.mq, err = NewMQ(p.in.MQType, p.in.conf)
	if err != nil {
		return
	}

	// MQ 准备
	err = p.in.mq.Ready()
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
			log.Errorln(p.Name, "Error receiving message:", err.Error())
			continue
		}
		log.Infoln(p.Name, "Recv:", string(msg))

		go p.dealMsg(msg)
	}
}

func (p *Component) dealMsg(msg []byte) {
	comsg := &Message{}
	if e := comsg.Unmarshal(msg); e != nil {
		log.Errorln(p.Name, "msg's format error:", e.Error(), string(msg))
		return
	}
	log.Infoln(p.Name, "Recv MSG:", comsg)

	next := comsg.PopGraph()
	if next == p.Name || next == p.Group {
		// 调用工作函数
		if p.handler != nil {
			rst, e := p.handler(comsg.Payload)
			if e != nil {
				log.Errorln(p.Name, "worker error, send to entrance:", msg, e.Error())
				next = comsg.Entrance
				comsg.Payload = nil
			} else {
				log.Infoln(p.Name, "Call handler ok")
				if rst == nil {
					return // 不往下走了, 入口才会这样
				}

				comsg.Payload = rst
			}
		}

		next = comsg.PopGraph()
		if next == "" {
			log.Warningln("next is nil. send to entrance:", msg)
			next = comsg.Entrance
		}
	} else if next == "" {
		log.Warningln("next is empty. send to entrance:", msg)
		next = comsg.Entrance
	}

	err := SendMsgToNext(next, comsg)
	if err != nil {
		log.Errorln(p.Name, "send to real next ERR: ", string(msg))
	}
}
