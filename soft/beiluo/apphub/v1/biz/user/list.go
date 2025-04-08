package user

import "github.com/yunhanshu-net/sdk-go/runner"

func init() {
	runner.Get("/user/list", List, &runner.ApiConfig{
		ApiDesc:     "获取用户列表",
		IsPublicApi: true,
		Labels:      []string{"用户信息", "table"},
		ChineseName: "获取用户列表",
		EnglishName: "user_list",
		Classify:    "用户信息",
		Request:     ListReq{},
		Response:    User{},
	})
}

type ListReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type User struct {
	Name string `json:"name" runner:"name:名称;code:name"`
	Dept string `json:"dept" runner:"name:部门;code:dept"`
}
type ListResp struct {
}

func List(ctx *runner.HttpContext) {
	users := []User{
		{Name: "a", Dept: "企业IT部"},
		{Name: "b", Dept: "数据平台部"},
	}
	err := ctx.Response.Table(users).Build()
	if err != nil {
		panic(err)
	}
}
