package test

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
	"sync"
)

func init() {
	db := runner.MustGetOrInitDB("test.db")
	db.AutoMigrate(&Calc{})
	runner.Get("/test/add", Add)
	runner.Get("/test/get", Get)
}

type Calc struct {
	ID int `gorm:"primaryKey;autoIncrement"`
	A  int `json:"a"`
	B  int `json:"b"`
	C  int `json:"c"`
}

type AddReq struct {
	A int `json:"a" form:"a"`
	B int `json:"b" form:"b"`
}
type GetReq struct {
	ID int `json:"id" form:"id"`
}

type AddResp struct {
	ID     int `json:"id"`
	Result int `json:"result"`
}

var lk = new(sync.Mutex)

func Add(ctx *runner.Context, req *AddReq, resp response.Response) error {
	//lk.Lock()
	//defer lk.Unlock()
	db := ctx.MustGetOrInitDB("test.db")
	//db.AutoMigrate(&Calc{})
	res := Calc{A: req.A, B: req.B, C: req.A + req.B}
	err := db.Model(&Calc{}).Create(&res).Error
	if err != nil {
		logrus.Errorf("Add err:%s", err.Error())
		return err
	}
	return resp.JSON(res).Build()
}

func Get(ctx *runner.Context, req *GetReq, resp response.Response) error {
	//lk.Lock()
	//defer lk.Unlock()
	db := ctx.MustGetOrInitDB("test.db")
	//db.AutoMigrate(&Calc{})
	res := Calc{}
	err := db.Model(&Calc{}).Where("id = ?", req.ID).First(&res).Error
	if err != nil {
		logrus.Errorf("get err:%s", err.Error())
		return err
	}
	return resp.JSON(res).Build()
}
