/*
节点组件
*/
package nodenet

import (
	"fmt"
	"reflect"
	"strings"

	log "github.com/golang/glog"
	gocommon "github.com/liuhengloveyou/go-common"
)

// 业务消息处理函数
type MessageHandler func(interface{}) (result interface{}, err error)

// 组件
type Component struct {
	Name  string // 组件名
	Group string // 组件所属的组

	in MessageQueue // 接收消息端点

	//每个组件可以处理多种消息
	handlers map[string]MessageHandler
}

func NewComponent(name, intype string, inconf interface{}) (*Component, error) {
	sname, sintype := strings.TrimSpace(name), strings.TrimSpace(intype)
	if sname == "" {
		return nil, fmt.Errorf("Component's name empty")
	}
	if sintype == "" {
		return nil, fmt.Errorf("Component's type empty")
	}

	components[sname] = &Component{Name: sname, in: nil, handlers: make(map[string]MessageHandler)}

	var e error
	components[sname].in, e = NewMQ(intype, inconf)

	return components[sname], e
}

func (p *Component) RegisterHandler(message interface{}, handler MessageHandler) {
	p.handlers[reflect.TypeOf(message).String()] = handler
}

func (p *Component) Run() error {
	gocommon.SingleInstane("/tmp/nodenet." + p.Name + ".pid")
	log.Infof("Component Run: [%v]", p.Name)

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

func (p *Component) dealMsg(msg string) {
	comsg := &Message{}
	if e := comsg.Decode([]byte(msg)); e != nil {
		log.Errorln(p.Name, "msg's format error:", e.Error(), msg)
		return
	}

	next := comsg.TopGraph()
	log.Infoln(p.Name, "Recv:", comsg, reflect.TypeOf(comsg.Payload), next, p.Name, p.Group)

	if next == p.Name || next == p.Group {
		// 调用工作函数
		if handler, ok := p.handlers[reflect.TypeOf(comsg.Payload).String()]; ok {
			rst, e := handler(comsg.Payload)
			if e != nil {
				log.Errorln(p.Name, "worker error, send to entrance:", msg, e.Error())
				next = comsg.Entrance
				comsg.Err = e
				comsg.Payload = nil
			} else {
				log.Infoln(p.Name, "Call handler ok")
				if rst == nil {
					return // 不往下走了
				}

				comsg.Payload = rst
			}
		} else {
			panic(fmt.Errorf("No handler for message: %s. %v", reflect.TypeOf(comsg.Payload).String(), comsg))
		}

		next = comsg.PopGraph()
		if next == "" {
			log.Infoln("next is nil:", comsg)
			// next = comsg.Entrance
		}
	} else if next == "" {
		log.Infoln("next is nil:", comsg)
		// next = comsg.Entrance
	}

	if next != "" {
		log.Infoln("Send msg to next:", next, comsg)
		if err := SendMsgToComponent(next, comsg); err != nil {
			log.Errorln(p.Name, "Send to next ERR: ", next, comsg)
		}
	}
}
