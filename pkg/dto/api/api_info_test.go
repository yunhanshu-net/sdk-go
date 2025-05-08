package api

import (
	"fmt"
	"testing"
)

type AddReq struct {
	A        int    `json:"a" form:"a" runner:"code:a;name:值a;required;type:number;number_limit:[-20000,10000];example:100;placeholder:请输入值a"`
	B        int    `json:"b" form:"b" runner:"code:b;name:值b;required;type:number;number_limit:[-20000,10000];example:200;placeholder:请输入值b"`
	Receiver string `json:"receiver" form:"receiver" runner:"code:receiver;name:接收人;default_value:beiluo;type:string;text_limit:1-20;placeholder:请输入接收人"`
	Desc     string `json:"desc" form:"desc" runner:"code:desc;name:描述;type:string;text_limit:1-50;placeholder:请描述此次计算"`
}

func TestName(t *testing.T) {
	params, err := NewRequestParams(&AddReq{}, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(params)
}
