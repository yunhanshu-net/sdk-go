package api

import (
	"encoding/json"
	"fmt"
	"testing"
)

type AddReq struct {
	A        int    `json:"a" form:"a" runner:"code:a;name:值a;required;type:number;example:100;placeholder:请输入值a" validate:"required,min=-1000,max=10000"`
	B        int    `json:"b" form:"b" runner:"code:b;name:值b;required;type:number;example:200;placeholder:请输入值b" validate:"required,min=-1000,max=10000"`
	Receiver string `json:"receiver" form:"receiver" runner:"code:receiver;name:接收人;widget:select;default_value:beiluo;options:admin,beiluo,user;type:string;placeholder:请输入接收人"`
	Desc     string `json:"desc" form:"desc" runner:"code:desc;name:描述;type:string;placeholder:请描述此次计算;callback:OnInputFuzzy"`
}

type AddResp struct {
	Result int `json:"result" runner:"code:result;name:计算结果;example:30000"`
}

func TestName(t *testing.T) {
	params, err := NewRequestParams(&AddReq{}, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(params)
	marshal, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}

func TestResponse(t *testing.T) {
	params, err := NewResponseParams(&AddResp{}, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(params)
	marshal, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
}
