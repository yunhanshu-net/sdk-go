package unix

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"sync"
	"time"
)

type Message struct {
	Data     []byte
	id       string         // 内部使用的UUID
	respChan chan *Response // 响应通道
}

type Response struct {
	Data []byte
	Err  error
}
type MessageConn struct {
	conn     net.Conn
	requests sync.Map // 存储 UUID -> *Message
}

func NewMessageConn(conn net.Conn) *MessageConn {
	mc := &MessageConn{conn: conn}
	go mc.processMessages()
	return mc
}

// RequestMsg 发送请求并等待响应（NATS风格）
func (mc *MessageConn) RequestMsg(data []byte, timeout time.Duration) (*Response, error) {
	msg := &Message{
		Data:     data,
		id:       uuid.New().String(),
		respChan: make(chan *Response, 1),
	}

	mc.requests.Store(msg.id, msg)
	defer mc.requests.Delete(msg.id)

	err := mc.writeWithUUID(msg.id, data)
	if err != nil {
		return nil, err
	}

	select {
	case resp := <-msg.respChan:
		return resp, resp.Err
	case <-time.After(timeout):
		return nil, fmt.Errorf("request timeout")
	}
}

// Respond 发送响应（服务端使用）
func (mc *MessageConn) Respond(request *Message, data []byte) error {
	return mc.writeWithUUID(request.id, data)
}

// 内部方法：处理消息循环
func (mc *MessageConn) processMessages() {
	for {
		id, data, err := mc.readWithUUID()
		if err != nil {
			mc.handleError(err)
			return
		}

		if msg, ok := mc.requests.Load(id); ok {
			msg.(*Message).respChan <- &Response{Data: data}
		} else {
			go func() {
				resp := &Response{Data: data}
				err := mc.writeWithUUID(id, resp.Data)
				if err != nil {
					panic(err)
				}
			}()
		}
	}
}

// 底层读写方法
func (mc *MessageConn) writeWithUUID(id string, data []byte) error {
	header := make([]byte, 4+16)
	uuidBytes, _ := uuid.Parse(id)
	binary.BigEndian.PutUint32(header[:4], uint32(16+len(data)))
	copy(header[4:], uuidBytes[:])
	_, err := mc.conn.Write(header)
	if err != nil {
		return err
	}
	_, err = mc.conn.Write(data)
	return err
}

func (mc *MessageConn) readWithUUID() (string, []byte, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(mc.conn, lengthBuf)
	if err != nil {
		return "", nil, err
	}
	totalLength := binary.BigEndian.Uint32(lengthBuf)

	uuidData := make([]byte, totalLength)
	_, err = io.ReadFull(mc.conn, uuidData)
	if err != nil {
		return "", nil, err
	}

	uuidBytes := uuidData[:16]
	data := uuidData[16:]
	id := uuid.UUID(uuidBytes).String()
	return id, data, nil
}

func (mc *MessageConn) handleError(err error) {
	mc.requests.Range(func(key, value interface{}) bool {
		msg := value.(*Message)
		msg.respChan <- &Response{Err: err}
		return true
	})
}

func (mc *MessageConn) Close() error {
	return mc.conn.Close()
}
