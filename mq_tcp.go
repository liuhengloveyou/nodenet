package nodenet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"sync"
	"time"

	log "github.com/golang/glog"
)

type MqTcp struct {
	Url     string        `json:"url"`
	Timeout time.Duration `json:"timeout"`

	ch   chan []byte
	conn net.Conn // 连接到该节点的客户端, 是消息发送方.
	lock sync.Mutex
}

func init() {
	RegisterMq("tcp", NewMqTcp)
}

func NewMqTcp(config interface{}) (MessageQueue, error) {
	tmq := &MqTcp{}
	if e := tmq.config(config); e != nil {
		return nil, e
	}

	tmq.ch = make(chan []byte)

	return tmq, nil
}

func (p *MqTcp) config(conf interface{}) error {
	confJson, e := json.Marshal(conf)
	if e != nil {
		return e
	}

	e = json.Unmarshal(confJson, p)
	if e != nil {
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

func (p *MqTcp) GetMessage() (msg []byte, e error) {
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

		go p.heartbeat()
	}

	n, mlen := 0, len(msg)
	buff := &bytes.Buffer{}

	var head [2]byte
	binary.LittleEndian.PutUint16(head[0:], uint16(mlen))
	buff.Write(head[0:])
	buff.Write(msg)

	p.conn.SetWriteDeadline(time.Now().Add(p.Timeout * time.Second))
	n, e = p.conn.Write(buff.Bytes())
	if e != nil {
		return
	}
	if n != buff.Len() {
		log.Fatalln("tcp.Write short.")
	}

	log.Infoln("MqTcp SendMessage: ", p.Url, string(msg))

	return nil
}

func (p *MqTcp) handleConnection(conn net.Conn) {
	var (
		msgLen int
		n      int
		e      error
		tmp    []byte
	)

	buf := make([]byte, 65535)
	msg := &bytes.Buffer{}

	for {
		if len(tmp) == 0 {
			conn.SetReadDeadline(time.Now().Add(p.Timeout * time.Second * 2))
			n, e = conn.Read(buf)
			if e != nil {
				log.Infoln("tcp.Read ERR: ", e.Error())
				conn.Close()
				conn = nil
				break

				/*if e == io.EOF {
					conn.Close()
					conn = nil
					break
				} else if neterr, ok := e.(net.Error); ok && neterr.Timeout() {
					log.Println("Timeout: ", e.Error())
					continue
				}*/
			}
		} else {
			n = len(tmp)
			for i := 0; i < n; i++ {
				buf[i] = tmp[i]
			}
		}

		log.Infoln(p.Url, "buf: ", msgLen, string(buf))

		if msgLen == 0 {
			tl := binary.LittleEndian.Uint16(buf)
			msgLen = int(tl)
			if msgLen == 0 {
				log.Infoln("heartbeat")
				continue // heartbeat
			}
			log.Infoln(p.Url, "msg: ", msgLen, n, string(msg.Bytes()))

			msg.Write(buf[2:n])
			msgLen = msgLen - n + 2
		} else {
			if msgLen > n {
				msg.Write(buf[:n])
				msgLen = msgLen - n
			} else {
				msg.Write(buf[:msgLen])
				tmp = buf[n-msgLen : n]
				msgLen = 0
			}
		}

		log.Infoln(p.Url, "msg: ", msgLen, string(msg.Bytes()))

		if msgLen == 0 {
			log.Infoln(p.Url, "Pop: ", string(msg.Bytes()))
			p.ch <- msg.Bytes()
			msg.Reset()
		}
	}
}

func (p *MqTcp) heartbeat() {
	msg := []byte{0, 0}

	for {
		time.Sleep(p.Timeout * time.Second)

		p.lock.Lock()
		p.conn.SetWriteDeadline(time.Now().Add(p.Timeout * time.Second))
		n, e := p.conn.Write(msg)
		p.lock.Unlock()
		if e != nil || n != len(msg) {
			log.Infoln("heartbeat ERR:", n, e.Error())
			p.conn.Close()
			p.conn = nil
			break
		}
	}
}
