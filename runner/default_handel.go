package runner

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/dto"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
	"gorm.io/gorm/schema"
	"reflect"
)

func env(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.JSON(map[string]string{"version": "1.0", "lang": "go"}).Build()
}

func ping(ctx *Context, req *request.NoData, resp response.Response) error {
	return resp.JSON(map[string]string{"ping": "pong"}).Build()
}

func (r *Runner) routerListInfo(ctx *Context, req *request.NoData, resp response.Response) error {
	functions := r.routerMap
	var configs []*ApiConfig
	for _, worker := range functions {
		if worker.IsDefaultRouter() {
			continue
		}
		worker.Config.Method = worker.Method
		worker.Config.Router = worker.Router
		if worker.Config != nil {
			if worker.Config.Request != nil {
				params, err := worker.Config.getParams(worker.Config.Request, "in")
				if err != nil {
					continue
				}
				worker.Config.ParamsIn = params
			}

			if worker.Config.Response != nil {
				params, err := worker.Config.getParams(worker.Config.Response, "out")
				if err != nil {
					continue
				}
				worker.Config.ParamsOut = params
			}
			configs = append(configs, worker.Config)
		}
	}

	return resp.JSON(configs).Build()
}

func (r *Runner) apiInfoGetAll(ctx *Context, req *request.NoData, resp response.Response) error {
	functions := r.routerMap
	var configs []*dto.ApiInfo
	for _, worker := range functions {
		if worker.IsDefaultRouter() {
			continue
		}
		config := worker.Config
		if config == nil {
			continue
		}
		api := &dto.ApiInfo{
			Method:      worker.Method,
			Router:      worker.Router,
			User:        r.detail.User,
			Runner:      r.detail.Name,
			ApiDesc:     config.ApiDesc,
			Labels:      config.Labels,
			ChineseName: config.ChineseName,
			EnglishName: config.EnglishName,
		}

		typeOf := reflect.TypeOf(worker.Config.Request)

		if typeOf.Kind() != reflect.Struct {
			return fmt.Errorf("输入参数仅支持Struct类型")
		}

		reqFields, err := tagx.ParseStructFieldsTypeOf(typeOf, "runner")
		if err != nil {
			return err
		}
		params, err := dto.NewParams(reqFields, config.Render)
		if err != nil {
			return err
		}
		api.ParamsIn = params

		rspType := reflect.TypeOf(worker.Config.Response)

		if rspType.Kind() != reflect.Struct || rspType.Kind() != reflect.Slice {
			return fmt.Errorf("输出参数仅支持Struct和Slice类型")
		}

		resFields, err := tagx.ParseStructFieldsTypeOf(rspType, "runner")
		if err != nil {
			return err
		}
		paramsOut, err := dto.NewParams(resFields, config.Render)
		if err != nil {
			return err
		}
		api.ParamsOut = paramsOut

		for _, table := range config.UseTables {
			if tb, ok := table.(schema.Tabler); ok {
				api.UseTables = append(api.UseTables, tb.TableName())
			}
		}

		//worker.Config.Method = worker.Method
		//worker.Config.Router = worker.Router
		//if worker.Config != nil {
		//	typeOf := reflect.TypeOf(worker.Config.Request)
		//
		//	if typeOf.Kind() != reflect.Struct {
		//		return fmt.Errorf("输入参数仅支持Struct类型")
		//	}
		//
		//	reqFields, err := tagx.ParseStructFieldsTypeOf(typeOf, "runner")
		//	if err != nil {
		//		return err
		//	}
		//
		//}
	}

	return resp.JSON(configs).Build()
}
