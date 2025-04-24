package runner

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/dto/api"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/pkg/jsonx"
	"gorm.io/gorm/schema"
)

func env(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.JSON(map[string]string{"version": "1.0", "lang": "go"}).Build()
}

func ping(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.JSON(map[string]string{"ping": "pong"}).Build()
}

// buildApiInfo 从路由信息构建API信息
func (r *Runner) buildApiInfo(worker *routerInfo) (*api.Info, error) {
	config := worker.Config
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
		Labels:      config.Labels,
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

	// 获取回调函数信息
	apiInfo.Callbacks = getCallbacks(config)

	return apiInfo, nil
}

func (r *Runner) getApiInfos(ctx *Context, req *request.NoData, resp response.Response) error {
	functions := r.routerMap
	var apis []*api.Info
	fmt.Println("routerMap:", jsonx.String(r.routerMap))
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
	return resp.JSON(apis).Build()
}

func (r *Runner) getApiInfo(ctx *Context, req *request.ApiInfoRequest, resp response.Response) error {
	// 参数验证
	if req.Router == "" {
		return resp.FailWithJSON(nil, "router参数不能为空")
	}

	// 如果没有指定Method，默认为GET
	if req.Method == "" {
		req.Method = "GET"
	}

	// 获取指定的路由信息
	worker, exist := r.getRouter(req.Router, req.Method)
	if !exist {
		return resp.FailWithJSON(nil, fmt.Sprintf("未找到路由: %s [%s]", req.Router, req.Method))
	}

	apiInfo, err := r.buildApiInfo(worker)
	if err != nil {
		return resp.FailWithJSON(nil, err.Error())
	}

	// 返回API信息
	return resp.JSON(apiInfo).Build()
}

func getCallbacks(config *ApiConfig) []string {
	var callbacks []string
	if config == nil {
		return nil
	}
	if config.OnPageLoad != nil {
		callbacks = append(callbacks, CallbackTypeOnPageLoad)
	}

	// API 生命周期回调
	if config.OnApiCreated != nil {
		callbacks = append(callbacks, CallbackTypeOnApiCreated)
	}
	if config.BeforeApiDelete != nil {
		callbacks = append(callbacks, CallbackTypeBeforeApiDelete)
	}
	if config.AfterApiDeleted != nil {
		callbacks = append(callbacks, CallbackTypeAfterApiDeleted)
	}

	// 运行器(Runner)生命周期回调
	if config.BeforeRunnerClose != nil {
		callbacks = append(callbacks, CallbackTypeBeforeRunnerClose)
	}
	if config.AfterRunnerClose != nil {
		callbacks = append(callbacks, CallbackTypeAfterRunnerClose)
	}

	// 版本控制回调
	if config.OnVersionChange != nil {
		callbacks = append(callbacks, CallbackTypeOnVersionChange)
	}

	// 输入交互回调
	if config.OnInputFuzzy != nil {
		callbacks = append(callbacks, CallbackTypeOnInputFuzzy)
	}
	if config.OnInputValidate != nil {
		callbacks = append(callbacks, CallbackTypeOnInputValidate)
	}

	// 表格操作回调
	if config.OnTableDeleteRows != nil {
		callbacks = append(callbacks, CallbackTypeOnTableDeleteRows)
	}
	if config.OnTableUpdateRow != nil {
		callbacks = append(callbacks, CallbackTypeOnTableUpdateRow)
	}
	if config.OnTableSearch != nil {
		callbacks = append(callbacks, CallbackTypeOnTableSearch)
	}

	return callbacks
}
