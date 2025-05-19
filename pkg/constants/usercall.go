package constants

const (
	// 页面事件
	CallbackTypeOnPageLoad = "OnPageLoad" // 页面加载时

	// UserCallTypeOnApiCreated API 生命周期
	UserCallTypeOnApiCreated    = "OnApiCreatedReq"    // API创建完成时
	CallbackTypeOnApiUpdated    = "OnApiUpdatedReq"    // API更新时
	CallbackTypeBeforeApiDelete = "BeforeApiDeleteReq" // API删除前
	CallbackTypeAfterApiDeleted = "AfterApiDeletedReq" // API删除后

	// 运行器(Runner)生命周期
	CallbackTypeBeforeRunnerClose = "BeforeRunnerCloseReq" // 运行器关闭前
	CallbackTypeAfterRunnerClose  = "AfterRunnerCloseReq"  // 运行器关闭后

	// 版本控制
	CallbackTypeOnVersionChange = "OnVersionChangeReq" // 版本变更时

	// 输入交互
	CallbackTypeOnInputFuzzy    = "OnInputFuzzyReq"    // 输入模糊匹配
	CallbackTypeOnInputValidate = "OnInputValidateReq" // 输入校验

	// 表格操作
	CallbackTypeOnTableDeleteRows = "OnTableDeleteRowsReq" // 删除表格行
	CallbackTypeOnTableUpdateRow  = "OnTableUpdateRowReq"  // 更新表格行
	CallbackTypeOnTableSearch     = "OnTableSearchReq"     // 表格搜索
)
