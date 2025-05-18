// 用户侧的回调，用户可以在ApiConfig中配置回掉的相关逻辑，不配置就不会触发
package runner

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/sdk-go/model/response"
	callback2 "github.com/yunhanshu-net/sdk-go/pkg/dto/callback"
	"github.com/yunhanshu-net/sdk-go/pkg/logger"
)

const (
	// 页面事件
	CallbackTypeOnPageLoad = "OnPageLoad" // 页面加载时

	// API 生命周期
	CallbackTypeOnApiCreated    = "OnApiCreated"    // API创建完成时
	CallbackTypeOnApiUpdated    = "OnApiUpdated"    // API更新时
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
type OnApiCreated func(ctx *Context, req *callback2.OnApiCreated) error

// OnApiUpdated 当api发生变更时候的回调函数
type OnApiUpdated func(ctx *Context, req *callback2.OnApiUpdated) error

// BeforeApiDelete  api删除前触发回调，比如该api删除的话，可以备份某些数据
type BeforeApiDelete func(ctx *Context, req *callback2.BeforeApiDelete) error

// AfterApiDeleted  api删除后触发回调，比如该api删除的话，可以在这里做一些操作，比如删除该api对应的表
type AfterApiDeleted func(ctx *Context, req *callback2.AfterApiDeleted) error

// BeforeRunnerClose 程序结束前的回调函数，可以在程序结束前做一些操作，比如上报一些数据
type BeforeRunnerClose func(ctx *Context, req *callback2.BeforeRunnerClose) error

// AfterRunnerClose 程序结束后的回调函数，可以在程序结束后做一些操作，比如清理某些文件
type AfterRunnerClose func(ctx *Context, req *callback2.AfterRunnerClose) error

// OnVersionChange 每次版本发生变更都会回调这个函数（新增/删除api）
type OnVersionChange func(ctx *Context, req *callback2.OnVersionChange) error

// OnInputFuzzy 模糊搜索回调函数，比如搜索用户，可以在这里做一些操作，比如根据用户名模糊搜索用户，然后返回用户列表
type OnInputFuzzy func(ctx *Context, req *callback2.OnInputFuzzy) (*response.OnInputFuzzy, error)

// OnInputValidate 验证输入框输入的名称是否重复或者输入是否合法
type OnInputValidate func(ctx *Context, req *callback2.OnInputValidate) (*response.OnInputValidate, error)

// OnTableDeleteRows 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有删除的行为，实现这个函数用来删除数据
type OnTableDeleteRows func(ctx *Context, req *callback2.OnTableDeleteRows) (*response.OnTableDeleteRows, error)

// OnTableUpdateRow 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有更新的行为，实现这个函数用来更新数据
type OnTableUpdateRow func(ctx *Context, req *callback2.OnTableUpdateRow) (*response.OnTableUpdateRow, error)

// OnTableSearch 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有搜索的行为，实现这个函数用来搜索数据
type OnTableSearch func(ctx *Context, req *callback2.OnTableSearch) (*response.OnTableSearch, error)

func (r *Runner) _callback(ctx *Context, req *callback2.Request, resp response.Response) (err error) {
	var res callback2.Response

	// 记录请求参数
	reqJSON, _ := json.Marshal(req)
	logger.InfoContextf(ctx, "处理回调 [类型:%s] [路由:%s] [方法:%s] 请求参数: %s", req.Type, req.Router, req.Method, string(reqJSON))

	worker, exist := r.getRouter(req.Router, req.Method)
	if !exist {
		err = fmt.Errorf("router not found")
		logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
		return err
	}
	if worker.ApiInfo == nil {
		err = fmt.Errorf("router config nil")
		logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
		return err
	}
	apiConf := worker.ApiInfo

	switch req.Type {
	// 页面加载回调
	case CallbackTypeOnPageLoad:
		if apiConf.OnPageLoad == nil {
			err = fmt.Errorf("OnPageLoad handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}

		resetReq, rsp, err := apiConf.OnPageLoad(ctx)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnPageLoad failed: %w", err)
		}
		res.Request = resetReq
		res.Response = rsp

		// 记录响应参数
		resetReqJSON, _ := json.Marshal(resetReq)
		rspJSON, _ := json.Marshal(rsp)
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s] 重置请求: %s, 响应: %s", req.Type, resetReqJSON, rspJSON)

	// API 生命周期回调
	case CallbackTypeOnApiCreated:
		var reqData callback2.OnApiCreated
		if apiConf.OnApiCreated == nil {
			err = fmt.Errorf("OnApiCreated handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnApiCreated decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.OnApiCreated(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnApiCreated failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	case CallbackTypeOnApiUpdated:
		var reqData callback2.OnApiUpdated
		if apiConf.OnApiUpdated == nil {
			err = fmt.Errorf("OnApiUpdated handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnApiUpdated decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.OnApiUpdated(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnApiUpdated failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	case CallbackTypeBeforeApiDelete:
		var reqData callback2.BeforeApiDelete
		if apiConf.BeforeApiDelete == nil {
			err = fmt.Errorf("BeforeApiDelete handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("BeforeApiDelete decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.BeforeApiDelete(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("BeforeApiDelete failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	case CallbackTypeAfterApiDeleted:
		var reqData callback2.AfterApiDeleted
		if apiConf.AfterApiDeleted == nil {
			err = fmt.Errorf("AfterApiDeleted handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("AfterApiDeleted decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.AfterApiDeleted(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("AfterApiDeleted failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	// Runner 生命周期回调
	case CallbackTypeBeforeRunnerClose:
		var reqData callback2.BeforeRunnerClose
		if apiConf.BeforeRunnerClose == nil {
			err = fmt.Errorf("BeforeRunnerClose handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("BeforeRunnerClose decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.BeforeRunnerClose(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("BeforeRunnerClose failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	case CallbackTypeAfterRunnerClose:
		var reqData callback2.AfterRunnerClose
		if apiConf.AfterRunnerClose == nil {
			err = fmt.Errorf("AfterRunnerClose handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("AfterRunnerClose decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.AfterRunnerClose(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("AfterRunnerClose failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	// 版本控制回调
	case CallbackTypeOnVersionChange:
		var reqData callback2.OnVersionChange
		if apiConf.OnVersionChange == nil {
			err = fmt.Errorf("OnVersionChange handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnVersionChange decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.OnVersionChange(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnVersionChange failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	// 输入交互回调
	case CallbackTypeOnInputFuzzy:
		var reqData callback2.OnInputFuzzy
		if apiConf.OnInputFuzzy == nil {
			err = fmt.Errorf("OnInputFuzzy handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnInputFuzzy decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := apiConf.OnInputFuzzy(ctx, &reqData)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnInputFuzzy failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)

	case CallbackTypeOnInputValidate:
		var reqData callback2.OnInputValidate
		if apiConf.OnInputValidate == nil {
			err = fmt.Errorf("OnInputValidate handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnInputValidate decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := apiConf.OnInputValidate(ctx, &reqData)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnInputValidate failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)

	// 表格操作回调
	case CallbackTypeOnTableDeleteRows:
		var reqData callback2.OnTableDeleteRows
		if apiConf.OnTableDeleteRows == nil {
			err = fmt.Errorf("OnTableDeleteRows handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnTableDeleteRows decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := apiConf.OnTableDeleteRows(ctx, &reqData)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnTableDeleteRows failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)

	case CallbackTypeOnTableUpdateRow:
		var reqData callback2.OnTableUpdateRow
		if apiConf.OnTableUpdateRow == nil {
			err = fmt.Errorf("OnTableUpdateRow handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnTableUpdateRow decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := apiConf.OnTableUpdateRow(ctx, &reqData)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnTableUpdateRow failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)

	case CallbackTypeOnTableSearch:
		var reqData callback2.OnTableSearch
		if apiConf.OnTableSearch == nil {
			err = fmt.Errorf("OnTableSearch handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnTableSearch decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := apiConf.OnTableSearch(ctx, &reqData)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnTableSearch failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)

	default:
		err = fmt.Errorf("unsupported callback type: %s", req.Type)
		logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 不支持的回调类型", req.Type)
		return err
	}

	err = resp.Form(res).Build()
	if err != nil {
		logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 构建响应失败 %v", req.Type, err)
		return err
	}

	// 在没有响应参数被记录的情况下，记录最终响应
	if res.Response != nil && (req.Type == CallbackTypeOnApiCreated ||
		req.Type == CallbackTypeOnApiUpdated ||
		req.Type == CallbackTypeBeforeApiDelete ||
		req.Type == CallbackTypeAfterApiDeleted ||
		req.Type == CallbackTypeBeforeRunnerClose ||
		req.Type == CallbackTypeAfterRunnerClose ||
		req.Type == CallbackTypeOnVersionChange) {
		resJSON, _ := json.Marshal(res.Response)
		logger.InfoContextf(ctx, "回调处理完成 [类型:%s] 响应: %s", req.Type, resJSON)
	}

	return nil
}
