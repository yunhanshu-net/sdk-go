package runner

import "runtime"

func (r *Runner) Debug(user, runner, version string, idle int) {
	r.isDebug = true
	r.idle = idle

	runtime.GOMAXPROCS(2)
	r.Get("/_env", env)
	r.Get("/_ping", ping)
	req := &TransportConfig{
		Runner:  runner,
		Version: version,
		User:    user,
	}
	//todo 长连接
	r.connect(req)
	r.listen()
}
