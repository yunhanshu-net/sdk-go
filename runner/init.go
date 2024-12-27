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
