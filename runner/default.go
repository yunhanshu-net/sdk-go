package runner

import (
	"fmt"
	constants "github.com/yunhanshu-net/pkg/constants/usercall"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/api"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
	"gorm.io/gorm/schema"
	"strings"
)

func (r *Runner) registerBuiltInRouters() {
	r.get("/_env", _env)
	r.get("/_help", r.help)
	r.get("/_ping", ping)
	r.get("/_getApiInfos", r._getApiInfos)
	r.get("/_getApiInfo", r._getApiInfo)
	r.post("/_callback", r._callback)
}
func _env(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.Form(map[string]string{"version": "1.0", "lang": "go"}).Build()
}

func (r *Runner) help(ctx *Context, req *request.NoData, resp response.Response) error {
	s := ""

	for _, info := range r.routerMap {
		s += info.Method + ":" + info.Router
	}
	return resp.Form(map[string]string{"router": s}).Build()
}

func ping(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.Form(map[string]string{"ping": "pong"}).Build()
}

// buildApiInfo 从路由信息构建API信息
func (r *Runner) buildApiInfo(worker *routerInfo) (*api.Info, error) {
	config := worker.ApiInfo
	if config == nil {
		return nil, fmt.Errorf("路由配置为空")
	}

	// 构建API信息
	apiInfo := &api.Info{
		Method:      worker.Method,
		Router:      worker.Router,
		User:        r.detail.User,
		Runner:      r.detail.Name,
		ApiDesc:     config.ApiDesc,
		Tags:        config.Tags,
		Classify:    config.Classify,
		ChineseName: config.ChineseName,
		EnglishName: config.EnglishName,
	}

	if config.Request != nil {
		// 获取请求参数信息
		params, err := api.NewRequestParams(config.Request, config.RenderType)
		if err != nil {
			return nil, err
		}
		apiInfo.ParamsIn = params
	}

	if config.Response != nil {
		// 获取响应参数信息
		responseParams, err := api.NewResponseParams(config.Response, config.RenderType)
		if err != nil {
			return nil, err
		}
		apiInfo.ParamsOut = responseParams
	}

	// 获取数据表信息
	for _, table := range config.UseTables {
		if tb, ok := table.(schema.Tabler); ok {
			apiInfo.UseTables = append(apiInfo.UseTables, tb.TableName())
		}
	}
	apiInfo.UseDB = config.UseDB

	// 获取回调函数信息
	apiInfo.Callbacks = getCallbacks(config)

	return apiInfo, nil
}
func (r *routerInfo) IsDefaultRouter() bool {
	return strings.HasPrefix(strings.TrimPrefix(r.Router, "/"), "_")
}
func (r *Runner) getApiInfos() ([]*api.Info, error) {
	functions := r.routerMap
	var apis []*api.Info
	for _, worker := range functions {
		if worker.IsDefaultRouter() {
			continue
		}
		apiInfo, err := r.buildApiInfo(worker)
		if err != nil {
			fmt.Println("apiInfo err:", err)
			continue // 跳过有错误的API
		}

		apis = append(apis, apiInfo)
	}
	return apis, nil
}

func (r *Runner) getApiInfo(req *request.ApiInfoRequest) (*api.Info, error) {
	// 参数验证
	if req.Router == "" {
		return nil, fmt.Errorf("router参数不能为空")
	}

	// 如果没有指定Method，默认为GET
	if req.Method == "" {
		req.Method = "GET"
	}

	// 获取指定的路由信息
	worker, exist := r.getRouter(req.Router, req.Method)
	if !exist {
		return nil, fmt.Errorf("未找到路由: %s [%s]", req.Router, req.Method)
	}

	apiInfo, err := r.buildApiInfo(worker)
	if err != nil {
		return nil, err
	}
	return apiInfo, nil

}

func (r *Runner) _getApiInfos(ctx *Context, req *request.NoData, resp response.Response) error {
	apis, err := r.getApiInfos()
	if err != nil {
		return err
	}
	return resp.Form(apis).Build()
}

func (r *Runner) _getApiInfo(ctx *Context, req *request.ApiInfoRequest, resp response.Response) error {
	apiInfo, err := r.getApiInfo(req)
	if err != nil {
		return err
	}
	// 返回API信息
	return resp.Form(apiInfo).Build()
}

func getCallbacks(config *ApiInfo) []string {
	var callbacks []string
	if config == nil {
		return nil
	}
	if config.OnPageLoad != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnPageLoad)
	}

	// API 生命周期回调
	if config.OnApiCreated != nil {
		callbacks = append(callbacks, constants.UserCallTypeOnApiCreated)
	}
	if config.OnApiUpdated != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnApiUpdated)
	}
	if config.BeforeApiDelete != nil {
		callbacks = append(callbacks, constants.CallbackTypeBeforeApiDelete)
	}
	if config.AfterApiDeleted != nil {
		callbacks = append(callbacks, constants.CallbackTypeAfterApiDeleted)
	}

	// 运行器(Runner)生命周期回调
	if config.BeforeRunnerClose != nil {
		callbacks = append(callbacks, constants.CallbackTypeBeforeRunnerClose)
	}
	if config.AfterRunnerClose != nil {
		callbacks = append(callbacks, constants.CallbackTypeAfterRunnerClose)
	}

	// 版本控制回调
	if config.OnVersionChange != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnVersionChange)
	}

	// 输入交互回调
	if config.OnInputFuzzy != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnInputFuzzy)
	}
	if config.OnInputValidate != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnInputValidate)
	}

	// 表格操作回调
	if config.OnTableDeleteRows != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnTableDeleteRows)
	}
	if config.OnTableUpdateRow != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnTableUpdateRow)
	}
	if config.OnTableSearch != nil {
		callbacks = append(callbacks, constants.CallbackTypeOnTableSearch)
	}

	return callbacks
}
