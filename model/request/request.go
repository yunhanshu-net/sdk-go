package request

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model"
	"time"
)

// Request 表示客户端请求
type Request struct {
	Reset      bool                `json:"reset"`
	TraceID    string              `json:"trace_id"`
	Route      string              `json:"route"`
	Method     string              `json:"method"`
	Headers    map[string]string   `json:"headers"`
	Body       interface{}         `json:"body"`        // 请求JSON对象
	BodyString string              `json:"body_string"` // 请求JSON字符串
	FileMap    map[string][]string `json:"file_map"`
	Timestamp  int64               `json:"timestamp"` // 请求时间戳
}

// NewRequest 创建一个新的请求对象
func NewRequest(route, method string) *Request {
	return &Request{
		TraceID:   generateTraceID(),
		Route:     route,
		Method:    method,
		Headers:   make(map[string]string),
		Timestamp: time.Now().Unix(),
		FileMap:   make(map[string][]string),
	}
}

// generateTraceID 生成追踪ID
func generateTraceID() string {
	return fmt.Sprintf("%d-%x", time.Now().UnixNano(), time.Now().UnixNano()%1000)
}

// SetHeader 设置请求头
func (r *Request) SetHeader(key, value string) *Request {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = value
	return r
}

// GetHeader 获取请求头
func (r *Request) GetHeader(key string) string {
	if r.Headers == nil {
		return ""
	}
	return r.Headers[key]
}

// SetBody 设置请求体
func (r *Request) SetBody(body interface{}) error {
	r.Body = body

	// 如果提供了body，同时更新BodyString
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("序列化请求体失败: %w", err)
		}
		r.BodyString = string(data)
	}

	return nil
}

// DecodeJSON 解析JSON请求体到指定的结构
func (r *Request) DecodeJSON(obj interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("请求体为空")
	}

	data, err := json.Marshal(r.Body)
	if err != nil {
		logrus.Errorf("序列化请求体失败: %v", err)
		return fmt.Errorf("序列化请求体失败: %w", err)
	}

	if err := json.Unmarshal(data, obj); err != nil {
		logrus.Errorf("解析请求体失败: %v", err)
		return fmt.Errorf("解析请求体失败: %w", err)
	}

	return nil
}

// ReSetJSON 重置请求体
func (r *Request) ReSetJSON(obj interface{}) error {
	r.Reset = true
	return r.SetBody(obj)
}

// AddFile 添加文件
func (r *Request) AddFile(field, filePath string) {
	if r.FileMap == nil {
		r.FileMap = make(map[string][]string)
	}
	r.FileMap[field] = append(r.FileMap[field], filePath)
}

// Runner 表示运行器信息
type Runner struct {
	Command         string `json:"command"`           // 命令
	WorkPath        string `json:"work_path"`         // 工作路径
	Name            string `json:"name"`              // 名称
	User            string `json:"user"`              // 用户
	Version         string `json:"version"`           // 版本
	RequestJsonPath string `json:"request_json_path"` // 请求JSON路径
}

// RunnerRequest 表示运行器请求
type RunnerRequest struct {
	UUID            string                 `json:"uuid"`             // 唯一ID
	Timeout         int                    `json:"timeout"`          // 超时时间(秒)
	Runner          *model.Runner          `json:"runner"`           // 运行器信息
	TransportConfig *TransportConfig       `json:"transport_config"` // 传输配置
	Metadata        map[string]interface{} `json:"metadata"`         // 元数据
	Request         *Request               `json:"request"`          // 请求信息
	Body            string                 `json:"body"`             // 请求体字符串
	CreatedAt       int64                  `json:"created_at"`       // 创建时间
}

// NewRunnerRequest 创建一个新的RunnerRequest
func NewRunnerRequest(runner *model.Runner, req *Request) *RunnerRequest {
	return &RunnerRequest{
		Runner:    runner,
		Request:   req,
		CreatedAt: time.Now().Unix(),
		Metadata:  make(map[string]interface{}),
	}
}

// TransportConfig 传输配置
type TransportConfig struct {
	IdleTime int                    `json:"idle_time"` // 空闲时间(秒)
	Type     string                 `json:"type"`      // 传输类型
	Metadata map[string]interface{} `json:"metadata"`  // 元数据
}

// GetSubject 获取NATS主题
func (r *RunnerRequest) GetSubject() string {
	if r.Runner == nil {
		return ""
	}
	return fmt.Sprintf("runner.%s.%s.%s.run", r.Runner.User, r.Runner.Name, r.Runner.Version)
}

// Bytes 将请求序列化为JSON字节
func (r *RunnerRequest) Bytes() ([]byte, error) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("序列化RunnerRequest失败: %w", err)
	}
	return jsonBytes, nil
}

// Ping 表示心跳检测请求
type Ping struct {
	Timestamp int64 `json:"timestamp"` // 时间戳
}

// NewPing 创建新的心跳检测请求
func NewPing() *Ping {
	return &Ping{
		Timestamp: time.Now().Unix(),
	}
}
