package runner

import "github.com/nats-io/nats.go"

func (r *Runner) onMsg(msg *nats.Msg) {
	r.wg.Add(1)
	defer r.wg.Done()
}
