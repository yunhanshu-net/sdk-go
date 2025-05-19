// 用户侧的回调，用户可以在ApiConfig中配置回掉的相关逻辑，不配置就不会触发
package runner

import (
	"encoding/json"
	"fmt"
	consts "github.com/yunhanshu-net/sdk-go/pkg/constants"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/response"
	"github.com/yunhanshu-net/sdk-go/pkg/dto/usercall"
	"github.com/yunhanshu-net/sdk-go/pkg/logger"
)

//
//const (
//	// 页面事件
//	consts.CallbackTypeOnPageLoad = "OnPageLoad" // 页面加载时
//
//	// UserCallTypeOnApiCreated API 生命周期
//	UserCallTypeOnApiCreated    = "OnApiCreatedReq"    // API创建完成时
//	consts.CallbackTypeOnApiUpdated    = "OnApiUpdatedReq"    // API更新时
//	consts.CallbackTypeBeforeApiDelete = "BeforeApiDeleteReq" // API删除前
//	consts.CallbackTypeAfterApiDeleted = "AfterApiDeletedReq" // API删除后
//
//	// 运行器(Runner)生命周期
//	consts.CallbackTypeBeforeRunnerClose = "BeforeRunnerCloseReq" // 运行器关闭前
//	consts.CallbackTypeAfterRunnerClose  = "AfterRunnerCloseReq"  // 运行器关闭后
//
//	// 版本控制
//	consts.CallbackTypeOnVersionChange = "OnVersionChangeReq" // 版本变更时
//
//	// 输入交互
//	consts.CallbackTypeOnInputFuzzy    = "OnInputFuzzyReq"    // 输入模糊匹配
//	consts.CallbackTypeOnInputValidate = "OnInputValidateReq" // 输入校验
//
//	// 表格操作
//	consts.CallbackTypeOnTableDeleteRows = "OnTableDeleteRowsReq" // 删除表格行
//	consts.CallbackTypeOnTableUpdateRow  = "OnTableUpdateRowReq"  // 更新表格行
//	consts.CallbackTypeOnTableSearch     = "OnTableSearchReq"     // 表格搜索
//)

// OnPageLoad 当用户进入某个函数的页面后，函数默认调用的行为，用户可以通过这个来初始化表单数据，resetRequest可以返回初始化后的表单数据
type OnPageLoad func(ctx *Context) (resetRequest interface{}, resp interface{}, err error)

// OnApiCreated 创建新的api时候的回调函数,新增一个api假如新增了一张user表， 可以在这里用gorm的db.AutoMigrate(&User)来创建表，
// 保证新版本的api可以正常使用新增的表 这个api只会在我创建这个api的时候执行一次
type OnApiCreated func(ctx *Context, req *usercall.OnApiCreatedReq) error

// OnApiUpdated 当api发生变更时候的回调函数
type OnApiUpdated func(ctx *Context, req *usercall.OnApiUpdatedReq) error

// BeforeApiDelete  api删除前触发回调，比如该api删除的话，可以备份某些数据
type BeforeApiDelete func(ctx *Context, req *usercall.BeforeApiDeleteReq) error

// AfterApiDeleted  api删除后触发回调，比如该api删除的话，可以在这里做一些操作，比如删除该api对应的表
type AfterApiDeleted func(ctx *Context, req *usercall.AfterApiDeletedReq) error

// BeforeRunnerClose 程序结束前的回调函数，可以在程序结束前做一些操作，比如上报一些数据
type BeforeRunnerClose func(ctx *Context, req *usercall.BeforeRunnerCloseReq) error

// AfterRunnerClose 程序结束后的回调函数，可以在程序结束后做一些操作，比如清理某些文件
type AfterRunnerClose func(ctx *Context, req *usercall.AfterRunnerCloseReq) error

// OnVersionChange 每次版本发生变更都会回调这个函数（新增/删除api）
type OnVersionChange func(ctx *Context, req *usercall.OnVersionChangeReq) error

// OnInputFuzzy 模糊搜索回调函数，比如搜索用户，可以在这里做一些操作，比如根据用户名模糊搜索用户，然后返回用户列表
type OnInputFuzzy func(ctx *Context, req *usercall.OnInputFuzzyReq) (*usercall.OnInputFuzzyResp, error)

// OnInputValidate 验证输入框输入的名称是否重复或者输入是否合法
type OnInputValidate func(ctx *Context, req *usercall.OnInputValidateReq) (*usercall.OnInputValidateResp, error)

// OnTableDeleteRows 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有删除的行为，实现这个函数用来删除数据
type OnTableDeleteRows func(ctx *Context, req *usercall.OnTableDeleteRowsReq) (*usercall.OnTableDeleteRowsResp, error)

// OnTableUpdateRow 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有更新的行为，实现这个函数用来更新数据
type OnTableUpdateRow func(ctx *Context, req *usercall.OnTableUpdateRowReq) (*usercall.OnTableUpdateRowResp, error)

// OnTableSearch 当返回前端的数据是table类型时候，前端会把数据渲染成表格，这时候表格数据会有搜索的行为，实现这个函数用来搜索数据
type OnTableSearch func(ctx *Context, req *usercall.OnTableSearchReq) (*usercall.OnTableSearchResp, error)

func (r *Runner) _callback(ctx *Context, req *usercall.Request, resp response.Response) (err error) {
	var res usercall.Response

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
	case consts.CallbackTypeOnPageLoad:
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
	case consts.UserCallTypeOnApiCreated:
		var reqData usercall.OnApiCreatedReq
		if apiConf.OnApiCreated == nil {
			err = fmt.Errorf("OnApiCreatedReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		//if err := req.DecodeData(&reqData); err != nil {
		//	logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
		//	return fmt.Errorf("OnApiCreatedReq decode failed: %w", err)
		//}
		//
		//// 记录请求详情
		//reqDataJSON, _ := json.Marshal(reqData)
		//logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.OnApiCreated(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnApiCreatedReq failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	case consts.CallbackTypeOnApiUpdated:
		var reqData usercall.OnApiUpdatedReq
		if apiConf.OnApiUpdated == nil {
			err = fmt.Errorf("OnApiUpdatedReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnApiUpdatedReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.OnApiUpdated(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnApiUpdatedReq failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	case consts.CallbackTypeBeforeApiDelete:
		var reqData usercall.BeforeApiDeleteReq
		if apiConf.BeforeApiDelete == nil {
			err = fmt.Errorf("BeforeApiDeleteReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("BeforeApiDeleteReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.BeforeApiDelete(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("BeforeApiDeleteReq failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	case consts.CallbackTypeAfterApiDeleted:
		var reqData usercall.AfterApiDeletedReq
		if apiConf.AfterApiDeleted == nil {
			err = fmt.Errorf("AfterApiDeletedReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("AfterApiDeletedReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.AfterApiDeleted(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("AfterApiDeletedReq failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	// Runner 生命周期回调
	case consts.CallbackTypeBeforeRunnerClose:
		var reqData usercall.BeforeRunnerCloseReq
		if apiConf.BeforeRunnerClose == nil {
			err = fmt.Errorf("BeforeRunnerCloseReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("BeforeRunnerCloseReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.BeforeRunnerClose(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("BeforeRunnerCloseReq failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	case consts.CallbackTypeAfterRunnerClose:
		var reqData usercall.AfterRunnerCloseReq
		if apiConf.AfterRunnerClose == nil {
			err = fmt.Errorf("AfterRunnerCloseReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("AfterRunnerCloseReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.AfterRunnerClose(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("AfterRunnerCloseReq failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	// 版本控制回调
	case consts.CallbackTypeOnVersionChange:
		var reqData usercall.OnVersionChangeReq
		if apiConf.OnVersionChange == nil {
			err = fmt.Errorf("OnVersionChangeReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnVersionChangeReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		if err := apiConf.OnVersionChange(ctx, &reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnVersionChangeReq failed: %w", err)
		}
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s]", req.Type)

	// 输入交互回调
	case consts.CallbackTypeOnInputFuzzy:
		var reqData usercall.OnInputFuzzyReq
		if apiConf.OnInputFuzzy == nil {
			err = fmt.Errorf("OnInputFuzzyReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnInputFuzzyReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := apiConf.OnInputFuzzy(ctx, &reqData)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnInputFuzzyReq failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)

	case consts.CallbackTypeOnInputValidate:
		var reqData usercall.OnInputValidateReq
		if apiConf.OnInputValidate == nil {
			err = fmt.Errorf("OnInputValidateReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnInputValidateReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := apiConf.OnInputValidate(ctx, &reqData)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnInputValidateReq failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)

	// 表格操作回调
	case consts.CallbackTypeOnTableDeleteRows:
		var reqData usercall.OnTableDeleteRowsReq
		if apiConf.OnTableDeleteRows == nil {
			err = fmt.Errorf("OnTableDeleteRowsReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnTableDeleteRowsReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := apiConf.OnTableDeleteRows(ctx, &reqData)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnTableDeleteRowsReq failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)

	case consts.CallbackTypeOnTableUpdateRow:
		var reqData usercall.OnTableUpdateRowReq
		if apiConf.OnTableUpdateRow == nil {
			err = fmt.Errorf("OnTableUpdateRowReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnTableUpdateRowReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := apiConf.OnTableUpdateRow(ctx, &reqData)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnTableUpdateRowReq failed: %w", err)
		}
		res.Response = respData

		// 记录响应参数
		respDataJSON, _ := json.Marshal(respData)
		logger.InfoContextf(ctx, "回调处理成功 [类型:%s] 响应: %s", req.Type, respDataJSON)

	case consts.CallbackTypeOnTableSearch:
		var reqData usercall.OnTableSearchReq
		if apiConf.OnTableSearch == nil {
			err = fmt.Errorf("OnTableSearchReq handler not configured")
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return err
		}
		if err := req.DecodeData(&reqData); err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: 解码失败 %v", req.Type, err)
			return fmt.Errorf("OnTableSearchReq decode failed: %w", err)
		}

		// 记录请求详情
		reqDataJSON, _ := json.Marshal(reqData)
		logger.InfoContextf(ctx, "回调处理中 [类型:%s] 请求详情: %s", req.Type, reqDataJSON)

		respData, err := apiConf.OnTableSearch(ctx, &reqData)
		if err != nil {
			logger.InfoContextf(ctx, "回调处理失败 [类型:%s]: %v", req.Type, err)
			return fmt.Errorf("OnTableSearchReq failed: %w", err)
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
	if res.Response != nil && (req.Type == consts.UserCallTypeOnApiCreated ||
		req.Type == consts.CallbackTypeOnApiUpdated ||
		req.Type == consts.CallbackTypeBeforeApiDelete ||
		req.Type == consts.CallbackTypeAfterApiDeleted ||
		req.Type == consts.CallbackTypeBeforeRunnerClose ||
		req.Type == consts.CallbackTypeAfterRunnerClose ||
		req.Type == consts.CallbackTypeOnVersionChange) {
		resJSON, _ := json.Marshal(res.Response)
		logger.InfoContextf(ctx, "回调处理完成 [类型:%s] 响应: %s", req.Type, resJSON)
	}

	return nil
}
