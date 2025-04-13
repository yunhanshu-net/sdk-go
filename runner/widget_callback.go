package runner

//// ApiCallback 此时无论程序是否运行都会调用
//type ApiCallback struct {
//
//	//创建新的api时候的回调函数,新增一个api假如新增了一张user表，
//	//可以在这里用gorm的db.AutoMigrate(&User)来创建表，保证新版本的api可以正常使用新增的表
//	//这个api只会在我创建这个api的时候执行一次
//	OnCreated func(ctx *HttpContext) error `json:"-"`
//
//	//api删除后触发回调，比如该api删除的话，可以在这里做一些操作，比如删除该api对应的表
//	AfterDelete func(ctx *HttpContext) error `json:"-"`
//
//	//每次版本发生变更都会回调这个函数（新增/删除api）
//	OnVersionChange func(ctx *HttpContext) error `json:"-"`
//}
//
//type WebCallback struct {
//	//假如该接口有对应的前端界面，渲染该界面后会调用该函数来加载默认请求数据，
//	//比如一个用户订单列表的页面，在点进去页面后会调用该回调
//	//此时已经知道是哪个用户的了，然后可以根据用户信息，展示该用户的默认数据。
//	//这样就省的用户自己输入用户名然后再点击运行按钮展示出来了
//	OnPageLoad func(ctx *HttpContext) error `json:"-"`
//}
//
//type ServerCallback struct {
//	//程序结束前的回调函数，可以在程序结束前做一些操作，比如上报一些数据
//	BeforeClose func(ctx *HttpContext) error `json:"-"`
//	//程序结束后的回调函数，可以在程序结束后做一些操作，比如清理某些文件
//	AfterClose func(ctx *HttpContext) error `json:"-"`
//}
//
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
