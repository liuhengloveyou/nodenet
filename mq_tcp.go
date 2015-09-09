package nodenet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"sync/atomic"
	"time"
)

type MqTcp struct {
	Url     string        `json:"url"`
	Timeout time.Duration `json:"timeout"`

	ch   chan []byte
	conn net.Conn
	ws   int32
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
		log.Panicln(e)
	}

	go func() {
		for {
			p.conn, e = listener.Accept()
			if e != nil {
				log.Panicln(e)
				//continue
			}

			go p.handleConnection()
		}
	}()
}

func (p *MqTcp) GetMessage() (msg []byte, e error) {
	return <-p.ch, nil
}

func (p *MqTcp) SendMessage(msg []byte) (e error) {
	if p.conn == nil {
		p.conn, e = net.Dial("tcp", p.Url)
		if e != nil {
			return
		}

		go p.heartbeat()
	}

	n, mlen := 0, len(msg)
	buff := &bytes.Buffer{}
	e = binary.Write(buff, binary.LittleEndian, uint16(mlen))
	if e != nil {
		return e
	}
	e = binary.Write(buff, binary.LittleEndian, msg)
	if e != nil {
		return e
	}

	atomic.AddInt32(&p.ws, 1)
	p.conn.SetWriteDeadline(time.Now().Add(p.Timeout * time.Second))
	n, e = p.conn.Write(buff.Bytes())
	atomic.AddInt32(&p.ws, -1)
	if e != nil {
		return
	}
	if n != buff.Len() {
		log.Panicln("tcp.Write short.")
	}

	log.Println("MqTcp SendMessage: ", p.Url, string(msg))

	return nil
}

func (p *MqTcp) handleConnection() {
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
			p.conn.SetReadDeadline(time.Now().Add(p.Timeout * time.Second * 2))
			n, e = p.conn.Read(buf)
			if e != nil {
				log.Println("tcp.Read ERR: ", e.Error())
				p.conn.Close()
				p.conn = nil
				break

				/*if e == io.EOF {
					p.conn.Close()
					p.conn = nil
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

		if msgLen == 0 {
			head := []byte{buf[0], buf[1]}
			msgLen64, _ := binary.Uvarint(head)
			msgLen = int(msgLen64)
			if msgLen == 0 {
				log.Println("heartbeat")
				continue // heartbeat
			}
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

		if msgLen == 0 {
			p.ch <- msg.Bytes()
			msg.Reset()
		}
	}
}

func (p *MqTcp) heartbeat() {
	msg := []byte{0, 0}

	for {
		time.Sleep(p.Timeout * time.Second)

		if p.ws > 0 {
			continue
		}

		p.conn.SetWriteDeadline(time.Now().Add(p.Timeout * time.Second))
		n, e := p.conn.Write(msg)
		if e != nil || n != len(msg) {
			log.Println("heartbeat ERR:", n, e.Error())
			p.conn.Close()
			p.conn = nil
			break
		}
	}
}
