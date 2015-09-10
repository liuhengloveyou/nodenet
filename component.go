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

// 组件
type Component struct {
	Name  string // 组件名
	Group string // 组件所属的组

	in      MessageQueue     // 接收消息端点
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

	components[sname] = &Component{Name: sname, in: nil, handler: nil}

	var e error
	components[sname].in, e = NewMQ(intype, inconf)

	return components[sname], e
}

func (p *Component) SetHandler(handler ComponentHandler) {
	p.handler = handler
}

func (p *Component) Run() (err error) {
	log.Infoln("Component Run...", p.Name)

	// MessageQueue 启动
	p.in.StartService()

	// 开始监听
	p.recvMonitor()

	return nil
}

func (p *Component) recvMonitor() {
	for {
		msg, err := p.in.GetMessage()
		if err != nil {
			log.Errorln(p.Name, "Error receiving message:", err.Error())
			continue
		}

		go p.dealMsg(msg)
	}
}

func (p *Component) dealMsg(msg []byte) {
	comsg := &Message{}
	if e := comsg.Unmarshal(msg); e != nil {
		log.Errorln(p.Name, "msg's format error:", e.Error(), string(msg))
		return
	}
	log.Infoln(p.Name, "Recv:", comsg)

	next := comsg.PopGraph()
	if next == p.Name || next == p.Group {
		// 调用工作函数
		if p.handler != nil {
			rst, e := p.handler(comsg.Payload)
			if e != nil {
				log.Errorln(p.Name, "worker error, send to entrance:", string(msg), e.Error())
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
			log.Warningln("next is nil. send to entrance:", string(msg))
			next = comsg.Entrance
		}
	} else if next == "" {
		log.Warningln("next is empty. send to entrance:", string(msg))
		next = comsg.Entrance
	}

	if err := SendMsgToComponent(next, comsg); err != nil {
		log.Errorln(p.Name, "Send to next ERR: ", next, string(msg))
	}
}
