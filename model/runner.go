package model

import (
	"fmt"
)

type Runner struct {
	Kind string `json:"kind"` //类型，可执行程序，so文件等等

	//WorkPath        string `json:"work_path"`
	//Command         string `json:"command"`
	//RequestJsonPath string `json:"request_json_path"`
	//Language        string `json:"language"`   //编程语言
	//StoreRoot       string `json:"store_root"` //oss 存储的跟路径
	Name string `json:"name"` //应用名称（英文标识）
	//ToolType        string `json:"tool_type"`  //工具类型
	Version string `json:"version"` //应用版本
	//OssPath         string `json:"oss_path"`   //文件地址
	User string `json:"user"` //所属租户
}

func (r *Runner) GetRequestSubject() string {
	return fmt.Sprintf("runner.%s.%s.%s.run", r.User, r.Name, r.Version)
}

func (r *Runner) GetCloseSubject() string {
	return fmt.Sprintf("runner.close.%s.%s.%s", r.User, r.Name, r.Version)
}

func (r *Runner) GetAddr() string {
	return fmt.Sprintf("unix.%s.%s.%s", r.User, r.Name, r.Version)
}
func (r *Runner) GetUnixPath() string {
	return fmt.Sprintf("%s_%s_%s.sock", r.User, r.Name, r.Version)
}

func (r *Runner) Check() error {
	if r.Name == "" {
		return fmt.Errorf("name 不能为空")
	}
	if r.Version == "" {
		return fmt.Errorf("version 不能为空")
	}

	if r.User == "" {
		return fmt.Errorf("user 不能为空")
	}
	return nil
}

type UpdateVersion struct {
	RunnerConf *Runner `json:"runner_conf"`
	OldVersion string  `json:"old_version"`
	//NewVersion        string  `json:"new_version"`
	//NewVersionOssPath string  `json:"new_version_oss_path"`
}
