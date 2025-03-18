package array

import "github.com/yunhanshu-net/sdk-go/runner"

func init() {
	runner.Post("/array/diff", DiffRunner)
}

type DiffReq struct {
	Base   []string `json:"base"`
	NewArr []string `json:"new_arr"`
}
type DiffResp struct {
	Add    []string `json:"add"`
	Delete []string `json:"delete"`
}

func Diff(req *DiffReq) (resp DiffResp, err error) {
	baseMap := make(map[string]struct{})
	newMap := make(map[string]struct{})

	// 记录Base元素
	for _, v := range req.Base {
		baseMap[v] = struct{}{}
	}

	// 记录NewArr元素并找出Add
	for _, v := range req.NewArr {
		newMap[v] = struct{}{}
		if _, existsInBase := baseMap[v]; !existsInBase {
			resp.Add = append(resp.Add, v)
		}
	}

	// 找出Delete：在Base中但不在NewArr中
	for _, v := range req.Base {
		if _, existsInNew := newMap[v]; !existsInNew {
			resp.Delete = append(resp.Delete, v)
		}
	}
	return resp, nil
}

func DiffRunner(ctx *runner.Context) {
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
	ctx.Response.OKWithJSON(diffResp)
}
