// 系统回调，系统会根据这些回调来执行一些默认的操作
package runner

import (
	"fmt"
	"reflect"

	"github.com/yunhanshu-net/sdk-go/model/dto/syscallback"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"

	"github.com/yunhanshu-net/sdk-go/model/dto/api"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

const (
	SysCallbackTypeOnVersionChange = "SysOnVersionChange" // 每次版本变更时，（新增，删除api），runcher都会触发此回掉
)

func (r *Runner) _sysOnVersionChange(ctx *Context, req *syscallback.SysOnVersionChangeReq, resp response.Response) error {
	err := r.onSysVersionChange(req)
	if err != nil {
		return err
	}
	return resp.JSON(map[string]interface{}{"success": true})
}

//// sysCallback 系统回调
//func (r *Runner) _sysCallback(ctx *Context, req *syscallback.Request, resp response.Response) error {
//	var res syscallback.Response
//	switch req.CallbackType {
//	case SysCallbackTypeOnVersionChange:
//		var data syscallback.SysOnVersionChangeReq
//		err := req.DecodeData(&data)
//		if err != nil {
//			return err
//		}
//		rsp, err := r.onSysVersionChange(&data)
//		if err != nil {
//			return err
//		}
//		res.Data = rsp
//	}
//	return resp.Form(res).Build()
//}

func (r *Runner) onSysVersionChange(req *syscallback.SysOnVersionChangeReq) error {
	apis, err := r.getApiInfos()
	if err != nil {
		return err
	}
	currentApiInfos := &api.ApiLogs{Version: r.detail.Version, Apis: apis}
	//lastApiInfos := &api.ApiLogs{Version: r.detail.Version, Apis: apis}
	//会在当前目录api-logs目录下生成一个文件，文件名是version.json，然后记录当前版本的api信息
	err = jsonx.SaveFile(fmt.Sprintf("./api-logs/%s.json", r.detail.Version), currentApiInfos)
	if err != nil {
		return err
	}

	//lastVersion, err := r.detail.GetLastVersion()
	//if err != nil {
	//	return nil, err
	//}

	//err = jsonx.UnmarshalFromFile(fmt.Sprintf("./api-logs/%s.json", lastVersion), lastApiInfos)
	//if err != nil {
	//	return nil, err
	//}

	//rsp := &syscallback.SysOnVersionChangeResp{
	//	AddApi:    make([]*api.Info, 0),
	//	DelApi:    make([]*api.Info, 0),
	//	UpdateApi: make([]*api.Info, 0),
	//}
	//
	//// 创建旧API的映射，用于快速查找
	//lastApiMap := make(map[string]*api.Info)
	//for _, lastApi := range lastApiInfos.Apis {
	//	key := fmt.Sprintf("%s:%s", lastApi.Method, lastApi.Router)
	//	lastApiMap[key] = lastApi
	//}
	//
	//// 遍历当前API，查找新增和更新的API
	//for _, currentApi := range currentApiInfos.Apis {
	//	key := fmt.Sprintf("%s:%s", currentApi.Method, currentApi.Router)
	//
	//	// 检查API是否在上一版本中存在
	//	if lastApi, exists := lastApiMap[key]; exists {
	//		// 检查API是否有更新
	//		if !apiEqual(currentApi, lastApi) {
	//			rsp.UpdateApi = append(rsp.UpdateApi, currentApi)
	//		}
	//
	//		// 标记已处理过的API
	//		delete(lastApiMap, key)
	//	} else {
	//		// 新增的API
	//		rsp.AddApi = append(rsp.AddApi, currentApi)
	//	}
	//}
	//
	//// 剩余未处理的旧API即为已删除的API
	//for _, api := range lastApiMap {
	//	rsp.DelApi = append(rsp.DelApi, api)
	//}

	return nil
}

// apiEqual 比较两个API是否相等
func apiEqual(a, b *api.Info) bool {
	// 使用reflect.DeepEqual进行深度比较
	return reflect.DeepEqual(a, b)
}
