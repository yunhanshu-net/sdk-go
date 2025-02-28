package runner

import "github.com/nats-io/nats.go"

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
	if t.msgKind == "nats" {
		tsNats := t.transport.(*transportNats)
		defer func() {
			tsNats.wg.Done()
			tsNats.responseMsgCount++
		}()
		msg := nats.NewMsg(t.Subject)
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
