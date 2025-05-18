package runner

func (r *Runner) registerBuiltInRouters() {
	r.get("/_env", env)
	r.get("/_help", r.help)
	r.get("/_ping", ping)
	r.get("/_getApiInfos", r._getApiInfos)
	r.get("/_getApiInfo", r._getApiInfo)
	r.post("/_callback", r._callback)
}
