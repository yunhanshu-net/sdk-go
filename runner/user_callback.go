// 用户侧的回调，用户可以在ApiConfig中配置回掉的相关逻辑，不配置就不会触发
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

// OnApiCreated 创建新的api时候的回调函数,新增一个api假如新增了一张user表， 可以在这里用gorm的db.AutoMigrate(&User)来创建表，
// 保证新版本的api可以正常使用新增的表 这个api只会在我创建这个api的时候执行一次
type OnApiCreated func(ctx *Context, req *request.OnApiCreated) error

// BeforeApiDelete  api删除前触发回调，比如该api删除的话，可以备份某些数据
type BeforeApiDelete func(ctx *Context, req *request.BeforeApiDelete) error

// AfterApiDeleted  api删除后触发回调，比如该api删除的话，可以在这里做一些操作，比如删除该api对应的表
type AfterApiDeleted func(ctx *Context, req *request.AfterApiDeleted) error

// BeforeRunnerClose 程序结束前的回调函数，可以在程序结束前做一些操作，比如上报一些数据
type BeforeRunnerClose func(ctx *Context, req *request.BeforeRunnerClose) error

// AfterRunnerClose 程序结束后的回调函数，可以在程序结束后做一些操作，比如清理某些文件
type AfterRunnerClose func(ctx *Context, req *request.AfterRunnerClose) error

// OnVersionChange 每次版本发生变更都会回调这个函数（新增/删除api）
type OnVersionChange func(ctx *Context, req *request.OnVersionChange) error

// OnInputFuzzy 模糊搜索回调函数，比如搜索用户，可以在这里做一些操作，比如根据用户名模糊搜索用户，然后返回用户列表
type OnInputFuzzy func(ctx *Context, req *request.OnInputFuzzy) (*response.OnInputFuzzy, error)

// OnInputValidate 验证输入框输入的名称是否重复或者输入是否合法
type OnInputValidate func(ctx *Context, req *request.OnInputValidate) (*response.OnInputValidate, error)

// OnTableDeleteRows 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有删除的行为，实现这个函数用来删除数据
type OnTableDeleteRows func(ctx *Context, req *request.OnTableDeleteRows) (*response.OnTableDeleteRows, error)

// OnTableUpdateRow 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有更新的行为，实现这个函数用来更新数据
type OnTableUpdateRow func(ctx *Context, req *request.OnTableUpdateRow) (*response.OnTableUpdateRow, error)

// OnTableSearch 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有搜索的行为，实现这个函数用来搜索数据
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
