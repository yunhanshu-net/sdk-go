package calc

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
)

var dbName = "test.db"

type Calc struct {
	ID       int    `gorm:"primaryKey;autoIncrement" runner:"code:id;name:id"`
	A        int    `json:"a" runner:"code:a;name:a"`
	B        int    `json:"b" runner:"code:b;name:b"`
	C        int    `json:"c" runner:"code:c;name:c"`
	Receiver string `json:"receiver" runner:"code:receiver;name:receiver"`
	Code     string `json:"code" runner:"code:code;name:code"`
}

func init() {
	addConfig := &runner.ApiConfig{
		Tags:        "数据管理;数据分析;记录管理",
		EnglishName: "calcAdd",
		ChineseName: "添加计算记录",
		ApiDesc:     "这里可以描述的详细一点",
		UseTables:   []interface{}{&Calc{}}, //这里会在注册这个api的时候自动创建相关的表
		OnPageLoad: func(ctx *runner.Context) (resetRequest interface{}, resp interface{}, err error) {
			return &AddReq{Receiver: ctx.GetUsername()}, nil, nil
		},
		OnInputValidate: func(ctx *runner.Context, req *request.OnInputValidate) (*response.OnInputValidate, error) {
			msg := ""
			if req.Key == "code" {
				if len(req.Value) > 64 {
					msg = "最长不能超过64个字符"
				}
				//其他判断......
			}
			return &response.OnInputValidate{Msg: msg}, nil
		},
	}

	runner.Get("/calc/add", Add, addConfig)
}

type AddReq struct {
	Receiver string `json:"receiver"`
	A        int    `json:"a" form:"a"`
	B        int    `json:"b" form:"b"`
	Code     string `json:"code" form:"code"`
}

type AddResp struct {
	ID     int `json:"id"`
	Result int `json:"result"`
}

// Add 拿这个处理函数举例，ctx是固定参数， req *AddReq是用户自定义的参数，根据接口请求参数自己定义，resp response.Response是固定参数，用户可以根据这个返回自己的json数据
func Add(ctx *runner.Context, req *AddReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB(dbName)
	res := Calc{A: req.A, B: req.B, C: req.A + req.B} //这里模拟处理逻辑
	err := db.Model(&Calc{}).Create(&res).Error
	if err != nil {
		logrus.Errorf("Add err:%s", err.Error())
		return err
	}
	return resp.Form(res).Build()
}
