package array

import (
	"github.com/yunhanshu-net/sdk-go/runner"
	"strings"
)

type DiffReq struct {
	Base      string `json:"base" runner:"desc:原数组;required:必填;example:1,2,3,4"`
	NewArr    string `json:"new_arr" runner:"desc:新数组;required:必填;example:3,4,5,6"`
	Separator string `json:"separator" runner:"desc:数组分隔符;required:必填;example:,;default:,"`
}
type DiffResp struct {
	Add    string `json:"add" runner:"desc:新增元素;example:5,6"`
	Delete string `json:"delete" runner:"desc:删除元素;example:1,2"`
}

func init() {

	paramsConfig := &runner.Config{
		ApiDesc:     "基于原数组和新数组进行对比，新数组存在原数组不存在的元素视为新增元素，新数组不存在原数组存在视为删除元素",
		ChineseName: "数组对比",
		Labels:      []string{"数据分析", "科研", "编程", "办公"}, //适用场景？
		EnglishName: "arrayDiff",
		Classify:    "数组",
		Tags:        "array,list,数组,集合", //表示元素本身的特性
		Request:     &DiffReq{},
		Response:    &DiffResp{},
	}
	runner.Post("/array/diff", DiffApi, paramsConfig)
}

func Diff(req *DiffReq) (resp *DiffResp, err error) {
	baseMap := make(map[string]struct{})
	newMap := make(map[string]struct{})

	// 记录Base元素
	for _, v := range strings.Split(req.Base, req.Separator) {
		baseMap[v] = struct{}{}
	}

	// 记录NewArr元素并找出Add
	add := make([]string, 0)
	del := make([]string, 0)
	for _, v := range strings.Split(req.NewArr, req.Separator) {
		newMap[v] = struct{}{}
		if _, existsInBase := baseMap[v]; !existsInBase {
			add = append(add, v)
		}
	}

	// 找出Delete：在Base中但不在NewArr中
	for _, v := range strings.Split(req.Base, req.Separator) {
		if _, existsInNew := newMap[v]; !existsInNew {
			del = append(del, v)
		}
	}
	resp.Add = strings.Join(add, req.Separator)
	resp.Delete = strings.Join(del, req.Separator)
	return resp, nil
}

func DiffApi(ctx *runner.HttpContext) {
	var req DiffReq
	err := ctx.Request.ShouldBindJSON(&req)
	if err != nil {
		ctx.Response.FailWithJSON(400, err.Error())
		return
	}
	diffResp, err := Diff(&req)
	if err != nil {
		ctx.Response.FailWithJSON(400, err.Error())
		return
	}
	err = ctx.Response.JSON(diffResp).Build()
	if err != nil {
		panic(err)
	}
}
