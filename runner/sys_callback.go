// 系统回调，系统会根据这些回调来执行一些默认的操作
package runner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/api"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/syscallback"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"

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

func (r *Runner) onSysVersionChange(req *syscallback.SysOnVersionChangeReq) error {
	apis, err := r.getApiInfos()
	if err != nil {
		return err
	}
	logrus.Infof("onSysVersionChange detail:%+v", r.detail)
	currentApiInfos := &api.ApiLogs{Version: r.detail.Version, Apis: apis}
	//lastApiInfos := &api.ApiLogs{Version: r.detail.Version, Apis: apis}
	//会在当前目录api-logs目录下生成一个文件，文件名是version.json，然后记录当前版本的api信息
	err = jsonx.SaveFile(fmt.Sprintf("./api-logs/%s.json", r.detail.Version), currentApiInfos)
	if err != nil {
		return err
	}

	return nil
}
