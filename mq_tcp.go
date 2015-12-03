package nodenet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"sync"
	"time"

	log "github.com/golang/glog"
)

type MqTcp struct {
	Url     string `json:"url"`
	Timeout int64  `json:"timeout"`

	ch   chan string
	conn net.Conn // 连接到该节点的客户端, 是消息发送方.
	lock sync.Mutex

	hb int32 // 是否已经启动心跳
}

func init() {
	RegisterMq("tcp", NewMqTcp)
}

func NewMqTcp(config interface{}) (MessageQueue, error) {
	tmq := &MqTcp{}
	if e := tmq.config(config); e != nil {
		return nil, e
	}

	tmq.ch = make(chan string)

	// 心跳间隔至少30秒
	if tmq.Timeout < 30 {
		tmq.Timeout = 30
	}

	return tmq, nil
}

func (p *MqTcp) config(conf interface{}) (e error) {
	var confJson []byte
	if confJson, e = json.Marshal(conf); e != nil {
		return e
	}

	if e = json.Unmarshal(confJson, p); e != nil {
		return e
	}

	return nil
}

func (p *MqTcp) StartService() {
	listener, e := net.Listen("tcp", p.Url)
	if e != nil {
		log.Fatalln(e)
	}
	log.Infoln("Listen:", p.Url)

	go func() {
		for {
			conn, e := listener.Accept()
			if e != nil {
				log.Fatalln(e)
				//continue
			}
			log.Infoln("Accept:", conn.LocalAddr(), conn.RemoteAddr())
			go p.handleConnection(conn)
		}
	}()
}

func (p *MqTcp) GetMessage() (msg string, e error) {
	return <-p.ch, nil
}

func (p *MqTcp) SendMessage(msg []byte) (e error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.conn == nil {
		p.conn, e = net.Dial("tcp", p.Url)
		if e != nil {
			return
		}

		if p.hb == 0 {
			p.hb++
			go p.heartbeat()
		}
	}

	n, mlen := 0, len(msg)
	buff := &bytes.Buffer{}

	var head [2]byte
	binary.LittleEndian.PutUint16(head[0:], uint16(mlen))
	buff.Write(head[0:])
	buff.Write(msg)

	p.conn.SetWriteDeadline(time.Now().Add(time.Duration(p.Timeout) * time.Second))
	n, e = p.conn.Write(buff.Bytes())
	if e != nil {
		log.Errorln("tcp.Write ERR:", e)
		p.conn.Close()
		p.conn = nil
		return
	}
	if n != buff.Len() {
		log.Fatalln("tcp.Write short.")
	}

	// log.Infoln("MqTcp SendMessage: ", p.Url, mlen, string(msg))

	return nil
}

func (p *MqTcp) handleConnection(conn net.Conn) {
	var (
		msgLen int // 当前一条消息的长度
		rn     int // 当前消息已读长度
	)

	buf := make([]byte, 65535)
	msg := &bytes.Buffer{}

	for {
	READ:
		if rn == 0 {
			conn.SetReadDeadline(time.Now().Add(time.Duration(p.Timeout) * time.Second * 2))
			n, e := conn.Read(buf[rn:])
			if e != nil {
				if e == io.EOF {
					conn.Close()
					conn = nil
					break
				} else if neterr, ok := e.(net.Error); ok && neterr.Timeout() {
					log.Infoln(p.Timeout, e)
					continue
				}
			}
			rn += n
		}

		if rn < 2 {
			log.Warningln("tcp.Read too short:", rn)
			continue
		}

		if msgLen == 0 {
			for p := 0; ; {
				msgLen = int(binary.LittleEndian.Uint16(buf[p:]))
				p += 2

				if msgLen > 0 {
					if msgLen >= (rn - p) {
						msg.Write(buf[p:rn])
						msgLen -= (rn - p)
						rn = 0
					} else {
						msg.Write(buf[p : p+msgLen])
						for i, j := 0, p+msgLen; j < rn; i, j = i+1, j+1 {
							buf[i] = buf[j]
						}
						rn -= (p + msgLen)
						msgLen = 0
					}
					break
				} else if msgLen < 0 {
					log.Exitln("tcp.Read ERR:", msgLen, string(buf)) // 错误退出
				} else if (p + 2) > rn /* msgLen == 0 heartbeat */ {
					rn = 0 /* 极小的机率会出问题 */
					goto READ
				}

			}
		} else {
			if msgLen >= rn {
				msg.Write(buf[:rn])
				msgLen -= rn
				rn = 0
			} else {
				msg.Write(buf[:msgLen])
				for i, j := 0, msgLen; j < rn; i, j = i+1, j+1 {
					buf[i] = buf[j]
				}
				rn -= msgLen
				msgLen = 0
			}
		}

		// log.Infoln(p.Url, "msg: ", msgLen, msg.String())

		if msgLen == 0 {
			p.ch <- msg.String()
			msg.Reset()
		}
	}
}

func (p *MqTcp) heartbeat() {
	msg := []byte{0, 0}

	for {
		time.Sleep(time.Duration(p.Timeout) * time.Second)

		if p.conn != nil {
			p.lock.Lock()
			p.conn.SetWriteDeadline(time.Now().Add(time.Duration(p.Timeout) * time.Second))
			n, e := p.conn.Write(msg)
			if e != nil || n != len(msg) {
				log.Errorln("heartbeat ERR:", n, e.Error())
				p.conn.Close()
				p.conn = nil
			}
			p.lock.Unlock()
		}
	}

	p.hb--
}
