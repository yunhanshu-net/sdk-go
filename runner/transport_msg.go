package runner

import (
	"github.com/nats-io/nats.go"
	"strings"
)

func newTransportMsg(t *TransportMsg) *TransportMsg {
	r := &TransportMsg{
		Subject:   t.Subject,
		transport: t.transport,
	}
	return r
}

type TransportMsg struct {
	Data     []byte
	Headers  MsgHeader
	MetaData map[string]interface{}
	Subject  string

	msg       interface{}
	transport Transport
	msgKind   string
}

func (t *TransportMsg) Reply(req *TransportMsg) error {
	if t.msgKind == "nats" || t.msgKind == "" {
		tsNats := t.transport.(*transportNats)
		defer func() {
			tsNats.wg.Done()
			tsNats.responseMsgCount++
		}()
		subject := strings.Replace(t.Subject, "runner.", "runcher.", 1)
		msg := nats.NewMsg(subject)
		for k, v := range req.Headers {
			msg.Header[k] = v
		}
		msg.Data = req.Data
		err := t.msg.(*nats.Msg).RespondMsg(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
