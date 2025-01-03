package runner

import (
	"os"
	"runtime"
)

func (r *Runner) init() {
	r.args = os.Args
	runtime.GOMAXPROCS(2)
	r.Get("/_env", env)
	r.Get("/_ping", ping)

}

func (r *Runner) GetCommand() string {
	return r.args[1]
}
