package runner

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/nats-io/nats.go"
//	"github.com/sirupsen/logrus"
//	"github.com/yunhanshu-net/sdk-go/model/request"
//)
//
//func (r *Runner) call(msg *nats.Msg) ([]byte, error) {
//
//	data := msg.Data
//	var req request.RunnerRequest
//	err1 := json.Unmarshal(data, &req)
//	if err1 != nil {
//		logrus.Errorf("call  json.Unmarshal(data, &req) err,req:%+v err:%s", req, err1.Error())
//		return nil, fmt.Errorf("call  json.Unmarshal(data, &req) err,req:%+v err:%s", req, err1.Error())
//	}
//
//	runResponse, err1 := r.runRequest(context.Background(), req.Request)
//	if err1 != nil {
//		logrus.Errorf("call runRequest err,req:%+v err:%s", req, err1.Error())
//		return nil, fmt.Errorf("call runRequest err,req:%+v err:%s", req, err1.Error())
//	}
//	marshal, err1 := json.Marshal(runResponse)
//	if err1 != nil {
//		logrus.Errorf("call json.Marshal err,req:%+v err:%s", req, err1.Error())
//		return nil, fmt.Errorf("call json.Marshal err,req:%+v err:%s", req, err1.Error())
//	}
//
//	return marshal, nil
//}

//
//// Validatable 定义了请求参数的验证接口
//// 实现此接口的请求对象可以在处理前进行自验证
//type Validatable interface {
//	Validate() error
//}
//
////type Response interface{}
//
//// handlerMeta 包含处理函数的元数据和初始化状态
//type handlerMeta struct {
//	once      sync.Once    // 确保只初始化一次
//	meta      *runtimeMeta // 懒加载的元数据
//	initError error        // 初始化错误
//}
//
//// runtimeMeta 保存处理函数的运行时元数据
//type runtimeMeta struct {
//	fnValue     reflect.Value // 处理函数的反射值
//	reqType     reflect.Type  // 请求参数的类型
//	reqPool     *sync.Pool    // 请求对象池，减少GC压力
//	hasValidate bool          // 是否实现了Validate接口
//}
//
//var (
//	handlerCacheMap = make(map[string]*handlerMeta) // key: string -> *handlerMeta
//	handlerCacheMux = &sync.RWMutex{}               // 保护handlerCacheMap的互斥锁
//)
//
//// 运行时构建元数据
//// buildRuntimeMeta 通过反射分析处理函数，构建运行时所需的元数据
//// fn 参数必须是符合 func(*Context, *T, response.Response) error 签名的函数
//// 其中 T 必须是一个结构体类型
//func buildRuntimeMeta(fn interface{}) (*runtimeMeta, error) {
//	rawHandler := fn
//
//	// 反射分析handler
//	funcType := reflect.TypeOf(rawHandler)
//	if funcType.Kind() != reflect.Func {
//		return nil, fmt.Errorf("必须为函数类型")
//	}
//
//	// 校验参数签名
//	if funcType.NumIn() != 3 || funcType.NumOut() != 1 {
//		return nil, fmt.Errorf("函数签名必须为func(*Context, *T, response.Response) error")
//	}
//
//	reqType := funcType.In(1)
//	if reqType.Kind() != reflect.Ptr || reqType.Elem().Kind() != reflect.Struct {
//		return nil, fmt.Errorf("第二个参数必须为结构体指针")
//	}
//
//	// 构造元数据
//	meta := &runtimeMeta{
//		fnValue: reflect.ValueOf(rawHandler),
//		reqType: reqType,
//		reqPool: &sync.Pool{
//			New: func() interface{} {
//				return reflect.New(reqType.Elem()).Interface()
//			},
//		},
//	}
//
//	// 检查是否实现了Validate方法
//	tempInstance := reflect.New(reqType.Elem()).Interface()
//	if v, ok := tempInstance.(Validatable); ok && v != nil {
//		meta.hasValidate = true
//	}
//
//	return meta, nil
//}
//
//// 实际调用逻辑
//func doCall(method string, meta *runtimeMeta, ctx *Context, resp *response.Data, body interface{}) error {
//	req := meta.reqPool.Get()
//	var err error
//
//	// 确保在所有错误路径上都返回对象到池中
//	defer meta.reqPool.Put(req)
//
//	if body != nil {
//		if method == "GET" {
//			switch body.(type) {
//			case string:
//				query, err1 := url.ParseQuery(body.(string))
//				if err1 != nil {
//					return fmt.Errorf("解析查询参数失败: %w", err1)
//				}
//				err1 = form.NewDecoder().Decode(req, query)
//				if err1 != nil {
//					return fmt.Errorf("解码表单数据失败: %w", err1)
//				}
//			default:
//				return fmt.Errorf("body type faild")
//			}
//
//		} else {
//			switch body.(type) {
//			case string:
//				err = json.Unmarshal([]byte(body.(string)), req)
//				if err != nil {
//					return fmt.Errorf("JSON解析失败: %w", err)
//				}
//			case map[string]interface{}:
//				marshal, err := json.Marshal(body)
//				if err != nil {
//					return err
//				}
//				err = json.Unmarshal(marshal, req)
//				if err != nil {
//					return err
//				}
//			}
//		}
//	}
//
//	// 执行验证（如果需要的话，取消注释并实现）
//	if meta.hasValidate {
//		if v, ok := req.(Validatable); ok && v != nil {
//			if err := v.Validate(); err != nil {
//				return fmt.Errorf("验证失败: %w", err)
//			}
//		}
//	}
//
//	if resp == nil {
//		resp = new(response.Data)
//	}
//	resp.TraceID = ctx.getTraceId()
//	// 反射调用
//	results := meta.fnValue.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req), reflect.ValueOf(resp)})
//	if result := results[0].Interface(); result != nil {
//		if err, ok := result.(error); ok {
//			return err
//		}
//	}
//	return nil
//}
