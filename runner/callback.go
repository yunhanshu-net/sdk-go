package runner

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
)

const (
	// 页面事件
	CallbackTypeOnPageLoad = "OnPageLoad" // 页面加载时

	// API 生命周期
	CallbackTypeOnApiCreated    = "OnApiCreated"    // API创建完成时
	CallbackTypeBeforeApiDelete = "BeforeApiDelete" // API删除前
	CallbackTypeAfterApiDeleted = "AfterApiDeleted" // API删除后

	// 运行器(Runner)生命周期
	CallbackTypeBeforeRunnerClose = "BeforeRunnerClose" // 运行器关闭前
	CallbackTypeAfterRunnerClose  = "AfterRunnerClose"  // 运行器关闭后

	// 版本控制
	CallbackTypeOnVersionChange = "OnVersionChange" // 版本变更时

	// 输入交互
	CallbackTypeOnInputFuzzy    = "OnInputFuzzy"    // 输入模糊匹配
	CallbackTypeOnInputValidate = "OnInputValidate" // 输入校验

	// 表格操作
	CallbackTypeOnTableDeleteRows = "OnTableDeleteRows" // 删除表格行
	CallbackTypeOnTableUpdateRow  = "OnTableUpdateRow"  // 更新表格行
	CallbackTypeOnTableSearch     = "OnTableSearch"     // 表格搜索
)

// OnPageLoad 当用户进入某个函数的页面后，函数默认调用的行为，用户可以通过这个来初始化表单数据，resetRequest可以返回初始化后的表单数据
type OnPageLoad func(ctx *Context) (resetRequest interface{}, resp interface{}, err error)

type OnApiCreated func(ctx *Context, req *request.OnApiCreated) error
type BeforeApiDelete func(ctx *Context, req *request.BeforeApiDelete) error
type AfterApiDeleted func(ctx *Context, req *request.AfterApiDeleted) error

type BeforeRunnerClose func(ctx *Context, req *request.BeforeRunnerClose) error
type AfterRunnerClose func(ctx *Context, req *request.AfterRunnerClose) error
type OnVersionChange func(ctx *Context, req *request.OnVersionChange) error

type OnInputFuzzy func(ctx *Context, req *request.OnInputFuzzy) (*response.OnInputFuzzy, error)
type OnInputValidate func(ctx *Context, req *request.OnInputValidate) (*response.OnInputValidate, error)

type OnTableDeleteRows func(ctx *Context, req *request.OnTableDeleteRows) (*response.OnTableDeleteRows, error)
type OnTableUpdateRow func(ctx *Context, req *request.OnTableUpdateRow) (*response.OnTableUpdateRow, error)
type OnTableSearch func(ctx *Context, req *request.OnTableSearch) (*response.OnTableSearch, error)

func (r *Runner) callback(ctx *Context, req *request.Callback, resp response.Response) (err error) {
	var res response.Callback

	worker, exist := r.getRouter(req.Router, req.Method)
	if !exist {
		return fmt.Errorf("router not found")
	}
	if worker.Config == nil {
		return fmt.Errorf("router config nil")
	}
	apiConf := worker.Config

	switch req.Type {
	// 页面加载回调
	case CallbackTypeOnPageLoad:
		if apiConf.OnPageLoad == nil {
			return fmt.Errorf("OnPageLoad handler not configured")
		}
		resetReq, rsp, err := apiConf.OnPageLoad(ctx)
		if err != nil {
			return fmt.Errorf("OnPageLoad failed: %w", err)
		}
		res.Request = resetReq
		res.Response = rsp

	// API 生命周期回调
	case CallbackTypeOnApiCreated:
		var reqData request.OnApiCreated
		if apiConf.OnApiCreated == nil {
			return fmt.Errorf("OnApiCreated handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("OnApiCreated decode failed: %w", err)
		}

		if err := apiConf.OnApiCreated(ctx, &reqData); err != nil {
			return fmt.Errorf("OnApiCreated failed: %w", err)
		}

	case CallbackTypeBeforeApiDelete:
		var reqData request.BeforeApiDelete
		if apiConf.BeforeApiDelete == nil {
			return fmt.Errorf("BeforeApiDelete handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("BeforeApiDelete decode failed: %w", err)
		}
		if err := apiConf.BeforeApiDelete(ctx, &reqData); err != nil {
			return fmt.Errorf("BeforeApiDelete failed: %w", err)
		}

	case CallbackTypeAfterApiDeleted:
		var reqData request.AfterApiDeleted
		if apiConf.AfterApiDeleted == nil {
			return fmt.Errorf("AfterApiDeleted handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("AfterApiDeleted decode failed: %w", err)
		}
		if err := apiConf.AfterApiDeleted(ctx, &reqData); err != nil {
			return fmt.Errorf("AfterApiDeleted failed: %w", err)
		}

	// Runner 生命周期回调
	case CallbackTypeBeforeRunnerClose:
		var reqData request.BeforeRunnerClose
		if apiConf.BeforeRunnerClose == nil {
			return fmt.Errorf("BeforeRunnerClose handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("BeforeRunnerClose decode failed: %w", err)
		}
		if err := apiConf.BeforeRunnerClose(ctx, &reqData); err != nil {
			return fmt.Errorf("BeforeRunnerClose failed: %w", err)
		}

	case CallbackTypeAfterRunnerClose:
		var reqData request.AfterRunnerClose
		if apiConf.AfterRunnerClose == nil {
			return fmt.Errorf("AfterRunnerClose handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("AfterRunnerClose decode failed: %w", err)
		}
		if err := apiConf.AfterRunnerClose(ctx, &reqData); err != nil {
			return fmt.Errorf("AfterRunnerClose failed: %w", err)
		}

	// 版本控制回调
	case CallbackTypeOnVersionChange:
		var reqData request.OnVersionChange
		if apiConf.OnVersionChange == nil {
			return fmt.Errorf("OnVersionChange handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("OnVersionChange decode failed: %w", err)
		}
		if err := apiConf.OnVersionChange(ctx, &reqData); err != nil {
			return fmt.Errorf("OnVersionChange failed: %w", err)
		}

	// 输入交互回调
	case CallbackTypeOnInputFuzzy:
		var reqData request.OnInputFuzzy
		if apiConf.OnInputFuzzy == nil {
			return fmt.Errorf("OnInputFuzzy handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("OnInputFuzzy decode failed: %w", err)
		}
		respData, err := apiConf.OnInputFuzzy(ctx, &reqData)
		if err != nil {
			return fmt.Errorf("OnInputFuzzy failed: %w", err)
		}
		res.Response = respData

	case CallbackTypeOnInputValidate:
		var reqData request.OnInputValidate
		if apiConf.OnInputValidate == nil {
			return fmt.Errorf("OnInputValidate handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("OnInputValidate decode failed: %w", err)
		}
		respData, err := apiConf.OnInputValidate(ctx, &reqData)
		if err != nil {
			return fmt.Errorf("OnInputValidate failed: %w", err)
		}
		res.Response = respData

	// 表格操作回调
	case CallbackTypeOnTableDeleteRows:
		var reqData request.OnTableDeleteRows
		if apiConf.OnTableDeleteRows == nil {
			return fmt.Errorf("OnTableDeleteRows handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("OnTableDeleteRows decode failed: %w", err)
		}
		respData, err := apiConf.OnTableDeleteRows(ctx, &reqData)
		if err != nil {
			return fmt.Errorf("OnTableDeleteRows failed: %w", err)
		}
		res.Response = respData

	case CallbackTypeOnTableUpdateRow:
		var reqData request.OnTableUpdateRow
		if apiConf.OnTableUpdateRow == nil {
			return fmt.Errorf("OnTableUpdateRow handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("OnTableUpdateRow decode failed: %w", err)
		}
		respData, err := apiConf.OnTableUpdateRow(ctx, &reqData)
		if err != nil {
			return fmt.Errorf("OnTableUpdateRow failed: %w", err)
		}
		res.Response = respData

	case CallbackTypeOnTableSearch:
		var reqData request.OnTableSearch
		if apiConf.OnTableSearch == nil {
			return fmt.Errorf("OnTableSearch handler not configured")
		}
		if err := req.DecodeData(&reqData); err != nil {
			return fmt.Errorf("OnTableSearch decode failed: %w", err)
		}
		respData, err := apiConf.OnTableSearch(ctx, &reqData)
		if err != nil {
			return fmt.Errorf("OnTableSearch failed: %w", err)
		}
		res.Response = respData

	default:
		return fmt.Errorf("unsupported callback type: %s", req.Type)
	}

	return resp.JSON(res).Build()
}

//type InputCallback struct {
//	//模糊搜索回调函数，比如搜索用户，可以在这里做一些操作，比如根据用户名模糊搜索用户，然后返回用户列表
//	OnFuzzy func(ctx *HttpContext) error `json:"-"`
//	//验证输入框输入的名称是否重复或者输入是否合法
//	OnValidate func(ctx *HttpContext) error `json:"-"`
//}
//
//type OnTableCallback struct {
//	OnDeleteRows func(ctx *HttpContext) error `json:"-"`
//	OnUpdateRow  func(ctx *HttpContext) error `json:"-"`
//	OnSearch     func(ctx *HttpContext) error `json:"-"`
//}

//func (r *Runner) callback(ctx *HttpContext) error {
//	//var call callback
//	//err := ctx.Request.DecodeJSON(&call)
//	//if err != nil {
//	//	return err
//	//}
//	//worker, exist := r.getRouterWorker(call.Router, call.Method)
//	//if exist {
//	//	return nil
//	//}
//	//if worker.Config == nil {
//	//	return nil
//	//}
//	//
//	//var callbackFunc func(ctx *HttpContext) error
//	//switch call.Type {
//	//case callbackTypeOnCreated:
//	//	if worker.Config.OnCreated != nil {
//	//		callbackFunc = worker.Config.OnCreated
//	//	}
//	//case callbackTypeOnVersionChange:
//	//	//遍历所有路由，只要有这个回调的，就执行
//	//	if worker.Config.OnVersionChange != nil {
//	//		callbackFunc = worker.Config.OnVersionChange
//	//	}
//	//case callbackTypeAfterDelete:
//	//	if worker.Config.AfterDelete != nil {
//	//		callbackFunc = worker.Config.AfterDelete
//	//	}
//	//case callbackTypeAfterClose:
//	//	//遍历所有路由，只要有这个回调的，就执行
//	//	if worker.Config.AfterClose != nil {
//	//		callbackFunc = worker.Config.AfterClose
//	//	}
//	//case callbackTypeBeforeClose:
//	//	//遍历所有路由，只要有这个回调的，就执行
//	//	if worker.Config.BeforeClose != nil {
//	//		callbackFunc = worker.Config.BeforeClose
//	//	}
//	//}
//	//if callbackFunc == nil {
//	//	return nil
//	//}
//	//
//	//err = callbackFunc(ctx)
//	//if err != nil {
//	//	return err
//	//}
//	//return nil
//
//}
