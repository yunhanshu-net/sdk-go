# SDK-GO 设计蓝图

## 1. 系统概述

SDK-GO是云函数(yunhanshu)项目的核心开发套件，是开发者创建、调试和部署云函数的基础工具。它提供了一套标准化的接口和组件，使开发者能够以Go语言快速构建基于云函数平台的应用。SDK本身只提供框架，开发者基于SDK开发的程序被称为"Runner"，用于实际执行业务逻辑。

## 2. 设计目标

1. **简单易用**：降低开发门槛，提供简洁直观的API
2. **功能完备**：支持所有云函数平台的功能特性
3. **高性能**：保证执行效率，支持高并发处理
4. **可扩展**：支持灵活的插件和中间件机制
5. **自动生成UI**：通过结构体标签自动生成前端界面
6. **标准化**：提供标准化的开发模式和项目结构

## 3. 核心架构设计

### 3.1 总体架构

```
+---------------------------+
|     Runner 应用程序       |
+---------------------------+
              |
+---------------------------+
|          SDK-GO          |
+-----------+---------------+
            |
+-----------v---------------+
|     核心组件和接口        |
|                          |
| +------------------------+
| |    路由管理            |
| +------------------------+
| |    请求处理            |
| +------------------------+
| |    响应构建            |
| +------------------------+
| |    数据存储            |
| +------------------------+
| |    渲染信息            |
| +------------------------+
| |    生命周期钩子        |
| +------------------------+
| |    API配置             |
| +------------------------+
```

### 3.2 执行流程

```
+------------------+    +-------------------+    +------------------+
|                  |    |                   |    |                  |
| 参数解析与验证   +--->+ 执行用户业务逻辑  +--->+ 构建标准化响应   |
|                  |    |                   |    |                  |
+------------------+    +-------------------+    +------------------+
```

## 4. 核心功能模块设计

### 4.1 路由系统

#### 4.1.1 基本设计
- 支持RESTful API路由注册
- 支持命令行参数映射到路由
- 支持路由分组和中间件

#### 4.1.2 路由注册函数
```go
// 注册GET请求处理函数
func Get(path string, handler interface{}, config *ApiConfig) *Route

// 注册POST请求处理函数
func Post(path string, handler interface{}, config *ApiConfig) *Route

// 注册PUT请求处理函数 
func Put(path string, handler interface{}, config *ApiConfig) *Route

// 注册DELETE请求处理函数
func Delete(path string, handler interface{}, config *ApiConfig) *Route
```

#### 4.1.3 路由组设计
```go
// 创建路由组
func Group(prefix string, middlewares ...HandlerFunc) *RouterGroup

// 路由组添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) *RouterGroup
```

### 4.2 请求处理

#### 4.2.1 请求解析
- 支持命令行参数解析
- 支持JSON请求体解析
- 支持URL查询参数和Path参数解析
- 自动参数验证

#### 4.2.2 上下文设计
```go
// 上下文结构体
type Context struct {
    Request     *Request
    Response    Response
    Params      map[string]string
    DB          *gorm.DB
    UserID      string
    Username    string
    DBCache     map[string]*gorm.DB
    // ... 其他字段
}

// 获取数据库连接
func (ctx *Context) MustGetOrInitDB(dbName string) *gorm.DB

// 获取用户信息
func (ctx *Context) GetUsername() string

// 获取用户ID
func (ctx *Context) GetUserID() string

// 获取请求参数
func (ctx *Context) Param(key string) string
```

### 4.3 响应构建

#### 4.3.1 响应类型
- JSON响应
- 表格数据响应
- 文件响应
- 图表响应

#### 4.3.2 响应构建器
```go
// 响应接口
type Response interface {
    Build() error
    JSON(data interface{}) Response
    File(path string) Response
    Table(data interface{}) TableResponse
    Chart(options *ChartOptions) ChartResponse
}

// 表格响应接口
type TableResponse interface {
    Response
    Columns(columns []TableColumn) TableResponse
    AutoPaginated(db *gorm.DB, model interface{}, pageInfo *request.PageInfo) TableResponse
}

// 图表响应接口
type ChartResponse interface {
    Response
    Options(options ChartOptions) ChartResponse
}
```

### 4.4 数据存储

#### 4.4.1 数据库管理
- 自动连接SQLite数据库
- 自动创建和管理表
- 支持数据迁移和版本控制

#### 4.4.2 数据库函数
```go
// 获取数据库连接
func MustGetOrInitDB(dbName string) *gorm.DB

// 自动迁移数据表
func AutoMigrate(db *gorm.DB, models ...interface{}) error

// 数据库事务
func Transaction(db *gorm.DB, fc func(tx *gorm.DB) error) error
```

### 4.5 API配置

#### 4.5.1 API配置结构
```go
// API配置结构体
type ApiConfig struct {
    // 基本信息
    Tags        string
    EnglishName string
    ChineseName string
    ApiDesc     string
    
    // 数据相关
    UseTables   []interface{}
    
    // 生命周期钩子
    OnPageLoad      func(ctx *Context) (resetRequest interface{}, resp interface{}, err error)
    OnApiCreated    func(ctx *Context, req *request.OnApiCreated) error
    AfterApiDeleted func(ctx *Context, req *request.AfterApiDeleted) error
    OnInputValidate func(ctx *Context, req *request.OnInputValidate) (*response.OnInputValidate, error)
    OnInputFuzzy    func(ctx *Context, req *request.OnInputFuzzy) (*response.OnInputFuzzy, error)
    
    // 其他配置
    AllowAnonymous bool
    RateLimit      int
    Timeout        int
    CacheStrategy  *CacheStrategy
}
```

#### 4.5.2 生命周期钩子
1. **OnPageLoad**：页面加载时触发，用于初始化请求参数和页面数据
2. **OnApiCreated**：API创建时触发，用于初始化表结构和数据
3. **AfterApiDeleted**：API删除后触发，用于清理资源
4. **OnInputValidate**：输入验证时触发，用于校验参数
5. **OnInputFuzzy**：输入模糊匹配时触发，用于提供自动完成功能

## 5. 标签(Tag)体系设计

### 5.1 标签设计原则
- 简单易用：标签命名简洁明确
- 功能完备：覆盖各种UI渲染场景
- 易于扩展：支持自定义标签和渲染行为

### 5.2 标签体系
```go
// 通用标签
// code: 表单的英文key
// name: 表单的名称
// desc: 参数的描述信息
// required: 是否必填参数
// example: 示例值
// default: 默认值
// type: 数据类型(string/number/bool/float/file)
// widget: 渲染组件类型

// 文本输入标签
// placeholder: 占位符
// fuzzy: 是否支持模糊查询
// max: 最长字符数量
// min: 最短字符数量

// 选择器标签
// options: 选项列表
// multiple: 是否多选
// search: 是否可搜索

// 数值输入标签
// min: 最小值
// max: 最大值
// step: 步长

// 文件上传标签
// accept: 接受的文件类型
// maxSize: 最大文件大小
// multiple: 是否多文件上传

// 日期选择标签
// format: 日期格式
// range: 是否日期范围
```

### 5.3 标签使用示例
```go
type UserForm struct {
    Username string `json:"username" runner:"code:username;name:用户名;desc:登录账号;required:true;min:3;max:20"`
    Age      int    `json:"age" runner:"code:age;name:年龄;type:number;min:0;max:150"`
    Email    string `json:"email" runner:"code:email;name:邮箱;desc:用于接收通知;required:true"`
    Role     string `json:"role" runner:"code:role;name:角色;widget:select;options:admin,user,guest"`
}
```

## 6. 组件设计

### 6.1 请求组件

#### 6.1.1 请求结构体
```go
type Request struct {
    Method string
    Path   string
    Query  map[string][]string
    Body   []byte
    Headers map[string][]string
}
```

#### 6.1.2 请求绑定
```go
// 绑定请求参数到结构体
func (r *Request) Bind(obj interface{}) error

// 绑定查询参数到结构体
func (r *Request) BindQuery(obj interface{}) error

// 绑定JSON请求体到结构体
func (r *Request) BindJSON(obj interface{}) error
```

### 6.2 日志组件

#### 6.2.1 日志级别
- DEBUG：调试信息
- INFO：一般信息
- WARN：警告信息
- ERROR：错误信息
- FATAL：致命错误

#### 6.2.2 日志函数
```go
// 设置日志级别
func SetLogLevel(level LogLevel)

// 记录调试日志
func Debug(format string, args ...interface{})

// 记录一般信息
func Info(format string, args ...interface{})

// 记录警告信息
func Warn(format string, args ...interface{})

// 记录错误信息
func Error(format string, args ...interface{})

// 记录致命错误并退出
func Fatal(format string, args ...interface{})
```

### 6.3 数据校验组件

#### 6.3.1 校验规则
- 必填校验
- 长度校验
- 范围校验
- 格式校验
- 自定义校验

#### 6.3.2 校验函数
```go
// 校验结构体
func Validate(obj interface{}) error

// 校验单个字段
func ValidateField(field interface{}, rules ...string) error

// 添加自定义校验器
func AddValidator(name string, fn ValidatorFunc)
```

### 6.4 安全组件

#### 6.4.1 用户认证
```go
// 获取当前用户
func GetCurrentUser(ctx *Context) (*User, error)

// 检查用户权限
func CheckPermission(ctx *Context, resource string, action string) bool
```

#### 6.4.2 数据加密
```go
// 加密数据
func Encrypt(data []byte, key []byte) ([]byte, error)

// 解密数据
func Decrypt(ciphertext []byte, key []byte) ([]byte, error)

// 生成哈希
func Hash(data []byte) string
```

## 7. 接口设计

### 7.1 处理函数接口
```go
// 标准处理函数定义
type HandlerFunc func(ctx *Context, req interface{}, resp Response) error

// 中间件定义
type MiddlewareFunc func(ctx *Context) error
```

### 7.2 生命周期钩子接口
```go
// 页面加载钩子
type OnPageLoadFunc func(ctx *Context) (resetRequest interface{}, resp interface{}, err error)

// API创建钩子
type OnApiCreatedFunc func(ctx *Context, req *request.OnApiCreated) error

// API删除后钩子
type AfterApiDeletedFunc func(ctx *Context, req *request.AfterApiDeleted) error

// 输入验证钩子
type OnInputValidateFunc func(ctx *Context, req *request.OnInputValidate) (*response.OnInputValidate, error)

// 输入模糊匹配钩子
type OnInputFuzzyFunc func(ctx *Context, req *request.OnInputFuzzy) (*response.OnInputFuzzy, error)
```

### 7.3 请求/响应模型

#### 7.3.1 通用请求模型
```go
// 分页信息
type PageInfo struct {
    Page     int `json:"page" form:"page"`
    PageSize int `json:"pageSize" form:"pageSize"`
}

// 排序信息
type SortInfo struct {
    Field     string `json:"field" form:"field"`
    Order     string `json:"order" form:"order"`
}

// 输入验证请求
type OnInputValidate struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

// 输入模糊匹配请求
type OnInputFuzzy struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}
```

#### 7.3.2 通用响应模型
```go
// 分页响应
type PageResult struct {
    Total    int64       `json:"total"`
    Page     int         `json:"page"`
    PageSize int         `json:"pageSize"`
    Data     interface{} `json:"data"`
}

// 输入验证响应
type OnInputValidate struct {
    Msg string `json:"msg"`
}

// 输入模糊匹配响应
type OnInputFuzzy struct {
    Values []string `json:"values"`
}
```

## 8. 自动UI生成机制

### 8.1 UI渲染规则
1. **请求参数**：根据结构体标签自动生成输入表单
2. **响应结果**：根据响应类型和标签自动生成展示界面
3. **表格数据**：自动生成带分页、排序的数据表格
4. **图表数据**：自动生成各类统计图表

### 8.2 输入组件映射
- string -> 文本输入框
- number/int/float -> 数字输入框
- bool -> 开关/复选框
- []string -> 多选框
- time.Time -> 日期选择器
- file -> 文件上传组件

### 8.3 输出组件映射
- 基本类型 -> 文本显示
- []interface{} -> 列表显示
- map[string]interface{} -> 键值对显示
- 表格数据 -> 分页表格
- 图表数据 -> 相应图表

## 9. 扩展机制设计

### 9.1 中间件系统
```go
// 添加全局中间件
func Use(middlewares ...MiddlewareFunc)

// 添加路由中间件
func (route *Route) Use(middlewares ...MiddlewareFunc) *Route
```

### 9.2 插件系统
```go
// 插件接口
type Plugin interface {
    Name() string
    Init(app *Application) error
    Destroy() error
}

// 注册插件
func RegisterPlugin(plugin Plugin) error

// 获取插件
func GetPlugin(name string) (Plugin, error)
```

### 9.3 自定义标签处理器
```go
// 标签处理函数
type TagProcessor func(field reflect.StructField, value reflect.Value, tag string) interface{}

// 注册标签处理器
func RegisterTagProcessor(tagName string, processor TagProcessor)
```

## 10. 命令行集成

### 10.1 命令行参数映射
```
runner <路由路径> <请求参数文件>
示例: tencent_tencentOa_v1 /calc/get ./request.json
```

### 10.2 输出格式
```
<Response>{JSON响应数据}</Response>
<Log>{日志信息}</Log>
```

## 11. 完整示例

### 11.1 一个计算服务的完整示例
```go
package calc

import (
    "github.com/sirupsen/logrus"
    "github.com/yunhanshu-net/sdk-go/model/request"
    "github.com/yunhanshu-net/sdk-go/model/response"
    "github.com/yunhanshu-net/sdk-go/runner"
)

var dbName = "calc.db"

// 数据模型
type Calc struct {
    ID       int    `gorm:"primaryKey;autoIncrement" runner:"code:id;name:ID"`
    A        int    `json:"a" runner:"code:a;name:A值;required:true"`
    B        int    `json:"b" runner:"code:b;name:B值;required:true"`
    C        int    `json:"c" runner:"code:c;name:结果"`
    Receiver string `json:"receiver" runner:"code:receiver;name:接收人"`
    Code     string `json:"code" runner:"code:code;name:编码;max:64"`
}

// 初始化路由
func init() {
    addConfig := &runner.ApiConfig{
        Tags:        "数据管理;数据分析;记录管理",
        EnglishName: "calcAdd",
        ChineseName: "添加计算记录",
        ApiDesc:     "添加两个数值的计算记录，并保存结果",
        UseTables:   []interface{}{&Calc{}},
        OnPageLoad: func(ctx *runner.Context) (resetRequest interface{}, resp interface{}, err error) {
            return &AddReq{Receiver: ctx.GetUsername()}, nil, nil
        },
        OnInputValidate: func(ctx *runner.Context, req *request.OnInputValidate) (*response.OnInputValidate, error) {
            msg := ""
            if req.Key == "code" && len(req.Value) > 64 {
                msg = "最长不能超过64个字符"
            }
            return &response.OnInputValidate{Msg: msg}, nil
        },
    }

    // 注册API路由
    runner.Post("/calc/add", Add, addConfig)
}

// 请求参数
type AddReq struct {
    Receiver string `json:"receiver" runner:"code:receiver;name:接收人"`
    A        int    `json:"a" form:"a" runner:"code:a;name:A值;required:true"`
    B        int    `json:"b" form:"b" runner:"code:b;name:B值;required:true"`
    Code     string `json:"code" form:"code" runner:"code:code;name:编码;max:64"`
}

// 响应参数
type AddResp struct {
    ID     int `json:"id" runner:"code:id;name:记录ID"`
    Result int `json:"result" runner:"code:result;name:计算结果"`
}

// 处理函数
func Add(ctx *runner.Context, req *AddReq, resp response.Response) error {
    // 获取数据库连接
    db := ctx.MustGetOrInitDB(dbName)
    
    // 业务处理
    res := Calc{
        A: req.A, 
        B: req.B, 
        C: req.A + req.B,
        Receiver: req.Receiver,
        Code: req.Code,
    }
    
    // 保存数据
    err := db.Model(&Calc{}).Create(&res).Error
    if err != nil {
        logrus.Errorf("Add err:%s", err.Error())
        return err
    }
    
    // 返回结果
    return resp.JSON(&AddResp{
        ID: res.ID,
        Result: res.C,
    }).Build()
}
```

### 11.2 查询服务示例
```go
package calc

import (
    "github.com/yunhanshu-net/sdk-go/model/request"
    "github.com/yunhanshu-net/sdk-go/model/response"
    "github.com/yunhanshu-net/sdk-go/runner"
    "strconv"
)

// 请求参数
type GetReq struct {
    ID int `json:"id" form:"id" runner:"code:id;name:记录ID"`
    *request.PageInfo
}

// 初始化路由
func init() {
    getConfig := &runner.ApiConfig{
        ChineseName: "获取计算记录",
        EnglishName: "calcGet",
        ApiDesc:     "查询计算记录，支持分页和条件过滤",
        Tags:        "数据管理;数据分析;记录管理",
        OnApiCreated: func(ctx *runner.Context, req *request.OnApiCreated) error {
            db := runner.MustGetOrInitDB(dbName)
            return db.AutoMigrate(&Calc{})
        },
        AfterApiDeleted: func(ctx *runner.Context, req *request.AfterApiDeleted) error {
            return runner.MustGetOrInitDB(dbName).Migrator().DropTable(&Calc{})
        },
        OnInputFuzzy: func(ctx *runner.Context, req *request.OnInputFuzzy) (*response.OnInputFuzzy, error) {
            var values []string
            if req.Key == "a" {
                db := ctx.MustGetOrInitDB(dbName)
                var calcs []Calc
                db.Model(&Calc{}).Where("a LIKE ?", "%"+req.Value+"%").Limit(10).Find(&calcs)
                for _, calc := range calcs {
                    values = append(values, strconv.Itoa(calc.A))
                }
            }
            return &response.OnInputFuzzy{Values: values}, nil
        },
    }
    
    // 注册API路由
    runner.Get("/calc/get", Get, getConfig)
}

// 处理函数
func Get(ctx *runner.Context, req *GetReq, resp response.Response) error {
    // 获取数据库连接
    db := ctx.MustGetOrInitDB(dbName)
    
    // 构建查询条件
    query := db.Model(&Calc{})
    if req.ID > 0 {
        query = query.Where("id > ?", req.ID)
    }
    
    // 返回表格数据，自动分页
    var results []Calc
    return resp.Table(&results).AutoPaginated(query, &Calc{}, req.PageInfo).Build()
}
```

## 12. 实现注意事项

### 12.1 性能优化
1. 使用高效的路由匹配算法
2. 减少反射使用，提高执行效率
3. 实现请求和响应对象池，减少内存分配
4. 优化数据库连接池和查询缓存

### 12.2 安全措施
1. 输入参数严格验证和过滤
2. 防止SQL注入和XSS攻击
3. 敏感数据加密存储
4. 权限严格控制

### 12.3 可维护性
1. 清晰的项目结构和代码组织
2. 完善的文档和注释
3. 单元测试和集成测试
4. 版本控制和兼容性保证

## 13. 未来扩展计划

### 13.1 多语言SDK
1. Python SDK
2. Java SDK
3. Node.js SDK

### 13.2 高级功能
1. 分布式事务支持
2. WebSocket长连接支持
3. 事件驱动和消息队列集成
4. 微服务架构支持

### 13.3 AI集成
1. 代码生成助手
2. 自然语言处理接口
3. 智能错误诊断
4. 性能优化建议 