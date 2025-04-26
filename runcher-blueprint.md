# Runcher 设计蓝图

## 1. 系统概述

Runcher是云函数(yunhanshu)项目的执行引擎组件，负责调度和执行基于SDK-GO开发的Runner程序。它作为API-Server和Runner之间的桥梁，接收来自API-Server的请求，调用对应的Runner执行业务逻辑，并将执行结果返回给API-Server。Runcher支持混合云部署，可以部署在公有云或私有云环境中。

## 2. 设计目标

1. **高性能**：保证Runner的高效调度和执行
2. **高可用**：支持集群部署，无单点故障
3. **安全隔离**：确保不同租户之间的资源隔离
4. **资源管理**：有效管理和监控资源使用
5. **可扩展**：支持水平和垂直扩展
6. **混合云**：支持跨云环境部署

## 3. 核心架构设计

### 3.1 总体架构

```
+----------------+     +-----------------+     +----------------+
|                |     |                 |     |                |
|  API-Server    +---->+    Runcher      +---->+    Runner     |
|                |     |                 |     |                |
+----------------+     +-----------------+     +----------------+
                              |
                     +--------v---------+
                     |                  |
                     |  Runner仓库      |
                     |                  |
                     +------------------+
```

### 3.2 核心组件

1. **NATS客户端**：与API-Server通信
2. **Runner管理器**：管理Runner的生命周期
3. **调度引擎**：调度Runner执行请求
4. **资源监控**：监控资源使用情况
5. **日志收集**：收集Runner执行日志
6. **服务发现**：自动注册到API-Server
7. **安全沙箱**：隔离Runner执行环境

### 3.3 执行流程

```
+-----------------+    +------------------+    +------------------+    +------------------+
|                 |    |                  |    |                  |    |                  |
| 接收NATS消息    +--->+ 解析执行参数     +--->+ 调用Runner       +--->+ 返回执行结果     |
|                 |    |                  |    |                  |    |                  |
+-----------------+    +------------------+    +------------------+    +------------------+
```

## 4. 核心功能模块设计

### 4.1 NATS通信模块

#### 4.1.1 基本设计
- 与API-Server建立NATS连接
- 监听指定主题的消息
- 处理请求并响应结果

#### 4.1.2 通信协议
```go
// 请求消息结构
type RunnerRequest struct {
    RequestID  string          `json:"request_id"`  // 请求ID
    UserID     string          `json:"user_id"`     // 用户ID
    RunnerName string          `json:"runner_name"` // Runner名称
    Version    string          `json:"version"`     // Runner版本
    Route      string          `json:"route"`       // 路由路径
    Method     string          `json:"method"`      // HTTP方法
    Headers    map[string][]string `json:"headers"` // 请求头
    Body       []byte          `json:"body"`        // 请求体
    Timeout    int             `json:"timeout"`     // 超时时间(ms)
}

// 响应消息结构
type RunnerResponse struct {
    RequestID  string          `json:"request_id"`  // 请求ID
    StatusCode int             `json:"status_code"` // 状态码
    Headers    map[string][]string `json:"headers"` // 响应头
    Body       []byte          `json:"body"`        // 响应体
    Error      string          `json:"error"`       // 错误信息
    Metrics    *ExecutionMetrics `json:"metrics"`   // 执行指标
}

// 执行指标
type ExecutionMetrics struct {
    ExecutionTime int64 `json:"execution_time"` // 执行时间(ms)
    MemoryUsage   int64 `json:"memory_usage"`   // 内存使用(KB)
    CPUUsage      float64 `json:"cpu_usage"`    // CPU使用率(%)
}
```

#### 4.1.3 消息处理
```go
// 注册消息处理器
func RegisterHandler(subject string, handler MessageHandler)

// 处理请求
func HandleRequest(msg *nats.Msg)

// 发送响应
func SendResponse(subject string, response *RunnerResponse)
```

### 4.2 Runner管理模块

#### 4.2.1 Runner生命周期管理
- 创建：根据用户请求创建Runner项目
- 升级：更新Runner版本
- 回滚：回滚到指定版本
- 删除：删除Runner

#### 4.2.2 Runner版本管理
```go
// 创建Runner
func CreateRunner(req *CreateRunnerRequest) error

// 更新Runner
func UpdateRunner(req *UpdateRunnerRequest) error

// 回滚Runner
func RollbackRunner(runnerName, version string) error

// 删除Runner
func DeleteRunner(runnerName string) error
```

#### 4.2.3 Runner存储结构
```
/root/
  ├── user1/
  │    ├── runner1/
  │    │    ├── bin/
  │    │    │    ├── data/
  │    │    │    │    └── data.db
  │    │    │    ├── .request/
  │    │    │    ├── user1_runner1_v1
  │    │    │    └── user1_runner1_v2
  │    │    └── version/
  │    │         ├── v1/
  │    │         │    ├── go.mod
  │    │         │    ├── main.go
  │    │         │    └── api/
  │    │         └── v2/
  │    │              ├── go.mod
  │    │              ├── main.go
  │    │              └── api/
  │    └── runner2/
  │         ├── ...
  └── user2/
       └── ...
```

### 4.3 调度引擎

#### 4.3.1 调度策略
- 请求负载均衡
- 资源利用优化
- 超时控制
- 错误处理

#### 4.3.2 执行函数
```go
// 执行Runner
func ExecuteRunner(req *RunnerRequest) (*RunnerResponse, error)

// 构建执行命令
func BuildExecuteCommand(runnerName, version, route string, requestFile string) *exec.Cmd

// 解析Runner输出
func ParseRunnerOutput(output []byte) (*RunnerResponse, error)
```

#### 4.3.3 并发控制
```go
// Runner执行池
type RunnerPool struct {
    // 最大并发数
    MaxConcurrent int
    // 当前执行数
    CurrentExecutions int
    // 执行队列
    Queue chan *RunnerExecutionTask
    // 执行完成通知
    Done chan *RunnerExecutionResult
}

// 提交执行任务
func (p *RunnerPool) Submit(task *RunnerExecutionTask) (*RunnerExecutionResult, error)

// 处理执行任务
func (p *RunnerPool) Process()
```

### 4.4 资源监控模块

#### 4.4.1 资源监控指标
- CPU使用率
- 内存使用量
- 磁盘使用量
- 网络使用量
- Runner执行数

#### 4.4.2 资源限制
```go
// 设置Runner资源限制
func SetRunnerResourceLimit(runnerName string, limit *ResourceLimit) error

// 资源限制结构
type ResourceLimit struct {
    MaxCPU    int   `json:"max_cpu"`    // 最大CPU核心数
    MaxMemory int64 `json:"max_memory"` // 最大内存(MB)
    MaxDisk   int64 `json:"max_disk"`   // 最大磁盘(MB)
    MaxExec   int   `json:"max_exec"`   // 最大并发执行数
}
```

#### 4.4.3 资源报告
```go
// 获取Runner资源使用报告
func GetRunnerResourceReport(runnerName string) (*ResourceReport, error)

// 资源报告结构
type ResourceReport struct {
    RunnerName      string    `json:"runner_name"`       // Runner名称
    CPUUsage        float64   `json:"cpu_usage"`         // CPU使用率(%)
    MemoryUsage     int64     `json:"memory_usage"`      // 内存使用(MB)
    DiskUsage       int64     `json:"disk_usage"`        // 磁盘使用(MB)
    CurrentExec     int       `json:"current_exec"`      // 当前执行数
    TotalExec       int64     `json:"total_exec"`        // 总执行次数
    AvgExecTime     int64     `json:"avg_exec_time"`     // 平均执行时间(ms)
    LastReportTime  time.Time `json:"last_report_time"`  // 最后报告时间
}
```

### 4.5 日志管理模块

#### 4.5.1 日志收集
- 收集Runner执行日志
- 收集系统运行日志
- 收集错误和异常日志

#### 4.5.2 日志处理
```go
// 记录Runner执行日志
func LogRunnerExecution(runnerName, route string, req *RunnerRequest, resp *RunnerResponse) error

// 解析Runner日志输出
func ParseRunnerLog(output []byte) ([]LogEntry, error)

// 日志条目结构
type LogEntry struct {
    Timestamp time.Time `json:"timestamp"` // 时间戳
    Level     string    `json:"level"`     // 日志级别
    Message   string    `json:"message"`   // 日志消息
    Context   map[string]interface{} `json:"context"` // 上下文信息
}
```

### 4.6 服务发现模块

#### 4.6.1 服务注册
- 自动注册到API-Server
- 定期发送心跳
- 保持连接状态

#### 4.6.2 服务注册函数
```go
// 注册Runcher服务
func RegisterRuncher(config *RuncherConfig) error

// 发送心跳
func SendHeartbeat() error

// Runcher配置
type RuncherConfig struct {
    ID          string   `json:"id"`           // Runcher ID
    Name        string   `json:"name"`         // Runcher名称
    Host        string   `json:"host"`         // 主机地址
    Port        int      `json:"port"`         // 端口
    ApiServer   string   `json:"api_server"`   // API-Server地址
    MaxRunners  int      `json:"max_runners"`  // 最大Runner数
    Tags        []string `json:"tags"`         // 标签
    Region      string   `json:"region"`       // 区域
}
```

### 4.7 安全模块

#### 4.7.1 Runner隔离
- 进程级隔离
- 资源限制
- 文件系统隔离

#### 4.7.2 权限控制
```go
// 检查执行权限
func CheckExecutePermission(userID, runnerName string) bool

// 验证API-Server请求
func VerifyApiServerRequest(req *http.Request) bool
```

## 5. 部署架构

### 5.1 单机部署
```
+-----------------+        +-----------------+
|                 |        |                 |
|   API-Server    |<------>|    Runcher      |
|                 |        |                 |
+-----------------+        +-----------------+
                                   |
                           +-------v--------+
                           |                |
                           |    Runners     |
                           |                |
                           +----------------+
```

### 5.2 集群部署
```
                   +------------------+
                   |                  |
                   |   API-Server     |
                   |                  |
                   +--------+---------+
                            |
            +---------------v----------------+
            |                                |
    +-------v-------+             +---------v-------+
    |               |             |                 |
    |   Runcher-1   |     ...     |   Runcher-N     |
    |               |             |                 |
    +-------+-------+             +---------+-------+
            |                               |
    +-------v-------+             +---------v-------+
    |               |             |                 |
    |   Runners-1   |     ...     |   Runners-N     |
    |               |             |                 |
    +---------------+             +-----------------+
```

### 5.3 混合云部署
```
                     +------------------+
                     |                  |
                     |   API-Server     |
                     |   (公有云)       |
                     |                  |
                     +--------+---------+
                              |
                     +--------v---------+
                     |                  |
                     |    NATS服务      |
                     |   (公有云)       |
                     |                  |
                     +--------+---------+
                              |
            +-----------------+------------------+
            |                                    |
   +--------v---------+               +----------v---------+
   |                  |               |                    |
   |    Runcher-1     |               |    Runcher-2       |
   |    (公有云)      |               |    (私有云)        |
   |                  |               |                    |
   +--------+---------+               +----------+---------+
            |                                    |
   +--------v---------+               +----------v---------+
   |                  |               |                    |
   |    Runners-1     |               |    Runners-2       |
   |    (公有云)      |               |    (私有云)        |
   |                  |               |                    |
   +------------------+               +--------------------+
```

## 6. 通信协议详细设计

### 6.1 NATS消息主题设计
```
yunhanshu.runcher.{region}.{runcher_id}.request     // 请求主题
yunhanshu.runcher.{region}.{runcher_id}.response    // 响应主题
yunhanshu.runcher.{region}.{runcher_id}.heartbeat   // 心跳主题
yunhanshu.runcher.{region}.{runcher_id}.admin       // 管理主题
yunhanshu.runcher.{region}.{runcher_id}.log         // 日志主题
```

### 6.2 Runner执行请求消息
```json
{
  "request_id": "req-123456789",
  "user_id": "user1",
  "runner_name": "runner1",
  "version": "v1",
  "route": "/calc/add",
  "method": "POST",
  "headers": {
    "Content-Type": ["application/json"],
    "Authorization": ["Bearer token"]
  },
  "body": "{\"a\":10,\"b\":20}",
  "timeout": 5000
}
```

### 6.3 Runner执行响应消息
```json
{
  "request_id": "req-123456789",
  "status_code": 200,
  "headers": {
    "Content-Type": ["application/json"]
  },
  "body": "{\"id\":1,\"result\":30}",
  "error": "",
  "metrics": {
    "execution_time": 120,
    "memory_usage": 1024,
    "cpu_usage": 0.5
  }
}
```

## 7. Runner生命周期管理详细设计

### 7.1 Runner创建流程
1. 接收API-Server创建请求
2. 在指定用户目录创建Runner目录结构
3. 初始化Go项目
4. 生成main.go和go.mod文件
5. 编译Runner可执行程序
6. 注册Runner到系统

### 7.2 Runner更新流程
1. 接收API-Server更新请求
2. 创建新版本目录
3. 复制已有代码到新版本目录
4. 更新指定文件
5. 编译新版本Runner
6. 更新系统中的Runner版本信息

### 7.3 Runner回滚流程
1. 接收API-Server回滚请求
2. 验证目标版本存在性
3. 更新系统中的Runner版本指向
4. 无需重新编译，直接使用历史版本

### 7.4 Runner执行流程
1. 接收API-Server执行请求
2. 解析请求参数
3. 生成请求JSON文件
4. 构建命令行调用命令
5. 执行Runner并捕获输出
6. 解析输出结果
7. 返回执行结果

## 8. 安全设计

### 8.1 Runner执行安全
1. **进程隔离**：每个Runner在独立进程中执行
2. **资源限制**：限制CPU、内存、磁盘使用
3. **超时控制**：防止无限执行
4. **文件系统隔离**：限制访问范围

### 8.2 通信安全
1. **TLS加密**：NATS通信加密
2. **Token认证**：API-Server与Runcher之间的认证
3. **消息签名**：验证消息完整性

### 8.3 数据安全
1. **用户数据隔离**：不同用户的数据存储在独立目录
2. **敏感信息保护**：敏感配置加密存储
3. **日志脱敏**：敏感数据日志脱敏

## 9. 可扩展性设计

### 9.1 水平扩展
- 支持多Runcher节点部署
- 基于标签的请求路由
- 节点动态扩缩容

### 9.2 垂直扩展
- 单节点资源动态调整
- Runner资源上限配置
- 按需分配资源

### 9.3 多区域支持
- 支持多区域部署
- 区域间通信和同步
- 就近接入和执行

## 10. 监控与运维

### 10.1 健康检查
- Runcher自身健康监控
- Runner健康状态检查
- 自动恢复机制

### 10.2 指标监控
- Runcher资源使用指标
- Runner执行指标
- API请求统计

### 10.3 日志管理
- 系统运行日志
- Runner执行日志
- 错误和异常日志

### 10.4 告警机制
- 资源使用告警
- 执行错误告警
- 系统异常告警

## 11. 高可用设计

### 11.1 无单点故障
- 多节点部署
- 节点故障自动转移
- 请求重试机制

### 11.2 容灾设计
- 多区域部署
- 数据备份与恢复
- 灾难恢复策略

### 11.3 降级策略
- 负载过高时自动降级
- 关键服务优先保障
- 限流与排队机制

## 12. 完整示例

### 12.1 Runcher启动配置示例
```json
{
  "id": "runcher-001",
  "name": "Runcher-Production",
  "host": "192.168.1.100",
  "port": 4222,
  "api_server": "https://api.yunhanshu.com",
  "max_runners": 100,
  "tags": ["production", "general"],
  "region": "cn-north",
  "root_dir": "/data/runners",
  "log_dir": "/var/log/runcher",
  "max_concurrent": 50,
  "timeout": 30000
}
```

### 12.2 Runner执行命令示例
```bash
# 执行Runner命令
/data/runners/user1/runner1/bin/user1_runner1_v1 /calc/add /data/runners/user1/runner1/bin/.request/req-123456789.json

# 请求参数文件内容 (req-123456789.json)
{
  "method": "POST",
  "headers": {
    "Content-Type": ["application/json"],
    "Authorization": ["Bearer token"]
  },
  "body": "{\"a\":10,\"b\":20}"
}
```

### 12.3 执行响应输出示例
```
<Response>{"id":1,"result":30}</Response>
<Log>2023-05-10 14:30:45 [INFO] Calculating 10 + 20</Log>
<Metrics>{"execution_time":123,"memory_usage":1024,"cpu_usage":0.5}</Metrics>
```

## 13. 实现注意事项

### 13.1 性能优化
1. 使用资源池管理Runner执行
2. 缓存热点Runner和请求
3. 异步处理非关键操作
4. 批量处理日志和监控数据

### 13.2 稳定性保障
1. 完善的错误处理和恢复机制
2. 限制单个Runner的资源使用
3. 防止恶意或错误代码影响系统
4. 监控系统关键指标

### 13.3 可维护性
1. 清晰的代码结构和模块划分
2. 完善的日志记录
3. 可配置的系统参数
4. 友好的管理接口

## 14. 未来扩展计划

### 14.1 功能扩展
1. 支持更多类型的Runtime
2. WebSocket和长连接支持
3. 文件和数据处理增强
4. 与CI/CD系统集成

### 14.2 性能提升
1. 更精细的资源调度算法
2. 预热和冷启动优化
3. 请求预测和资源预分配
4. 分布式执行引擎

### 14.3 生态集成
1. 支持容器化执行环境
2. 与云服务商资源集成
3. 第三方监控和日志系统集成
4. 多语言SDK支持 