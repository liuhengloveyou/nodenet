/*
节点组件
*/
package nodenet

import (
	"fmt"
	"strings"

	log "github.com/golang/glog"
)

var components map[string]*Component = make(map[string]*Component)

// 业务处理函数
type ComponentHandler func(interface{}) (result interface{}, err error)

// 端点
type EndPoint struct {
	MQType string
	Url    string
	conf   interface{}
	mq     MessageQueue
}

// 组件
type Component struct {
	Name    string               // 组件名
	Group   string               // 组件所属的组
	in      EndPoint             // 接收消息端点
	outs    map[string]*EndPoint // 所有下游端点
	handler ComponentHandler     // 业务处理函数
}

func NewComponent(name string, inconf interface{}) (*Component, error) {
	sname := strings.TrimSpace(name)
	if sname == "" {
		return nil, fmt.Errorf("Component's name empty")
	}

	com := &Component{
		Name:    sname,
		in:      EndPoint{MQType: "", Url: "", conf: inconf, mq: nil},
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

	_, err := p.sendToNext(next, comsg)
	if err != nil {
		log.Errorln(p.Name, "send to real next ERR: ", string(msg))
	}
}

func (p *Component) sendToNext(name string, comsg *Message) (total int, err error) {
	/*
		if url == "" {
			return 0, fmt.Errorf("sendTo nil url")
		}

		if _, ok := p.outs[url]; ok == false {
			mqtmp, err := NewMq(p.in.MQType, url)
			if err != nil {
				return 0, err
			}
			p.outs[url] = &EndPoint{Url: url, MQType: p.in.MQType, mq: mqtmp}
		}

		log.Infoln(p.Name, "sendToNext:", url, string(msg))
		total, err = p.outs[url].mq.SendToNext(msg)
	*/
	return
}
