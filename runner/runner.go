package runner

type Conn struct {
}

type Runner struct {
	conn *Conn

	WorkPath   string `json:"work_path"`
	RunnerType string `json:"runner_type"`
	Version    string `json:"version"`
	Command    string `json:"command"` //命令
	User       string `json:"user"`    //软件所属的用户
	Soft       string `json:"soft"`    //软件名
	OssPath    string `json:"oss_path"`
}
