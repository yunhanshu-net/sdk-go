package test

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
)

func init() {
	getConfig := &runner.ApiConfig{
		OnApiCreated: func(ctx *runner.Context, req *request.OnApiCreated) error {
			db := runner.MustGetOrInitDB("test.db")
			return db.AutoMigrate(&Calc{})
		},

		OnPageLoad: func(ctx *runner.Context) (resetRequest interface{}, resp interface{}, err error) {
			return &AddReq{A: 1, B: 2}, nil, nil
		},

		OnInputFuzzy: func(ctx *runner.Context, req *request.OnInputFuzzy) (*response.OnInputFuzzy, error) {
			if req.Value != "" {
				return &response.OnInputFuzzy{Values: []string{
					req.Value + "1",
					req.Value + "2",
					req.Value + "3",
					req.Value + "4"}}, nil
			}
			return nil, nil
		},
	}

	addConfig := &runner.ApiConfig{
		OnApiCreated: func(ctx *runner.Context, req *request.OnApiCreated) error {
			db := runner.MustGetOrInitDB("test.db")
			return db.AutoMigrate(&Calc{})
		},

		OnPageLoad: func(ctx *runner.Context) (resetRequest interface{}, resp interface{}, err error) {

			return nil, nil, err
		},
	}

	runner.Get("/test/add", Add, addConfig)
	runner.Get("/test/get", Get, getConfig)
}

type Calc struct {
	ID int `gorm:"primaryKey;autoIncrement" runner:"code:id;name:id"`
	A  int `json:"a" runner:"code:a;name:a"`
	B  int `json:"b" runner:"code:b;name:b"`
	C  int `json:"c" runner:"code:c;name:c"`
}

type AddReq struct {
	A int `json:"a" form:"a"`
	B int `json:"b" form:"b"`
}
type GetReq struct {
	ID int `json:"id" form:"id"`
	*request.PageInfo
}

type AddResp struct {
	ID     int `json:"id"`
	Result int `json:"result"`
}

// Add 拿这个处理函数举例，ctx是固定参数， req *AddReq是用户自定义的参数，根据接口请求参数自己定义，resp response.Response是固定参数，用户可以根据这个返回自己的json数据
func Add(ctx *runner.Context, req *AddReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB("test.db")
	res := Calc{A: req.A, B: req.B, C: req.A + req.B}
	err := db.Model(&Calc{}).Create(&res).Error
	if err != nil {
		logrus.Errorf("Add err:%s", err.Error())
		return err
	}
	return resp.JSON(res).Build()
}

func Get(ctx *runner.Context, req *GetReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB("test.db")
	var res []Calc
	db.Where("id > ?", req.ID)
	return resp.Table(&res).AutoPaginated(db, &Calc{}, req.PageInfo).Build()
}
