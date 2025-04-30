package model

import (
	"fmt"
	"strconv"
	"strings"
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

func (r *Runner) GetLastVersion() (string, error) {
	all := strings.ReplaceAll(r.Version, "v", "")
	v, err := strconv.Atoi(all)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("v%d", v-1), nil
}
