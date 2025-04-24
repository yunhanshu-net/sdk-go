package runner

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-playground/form/v4"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"net/url"
	"reflect"
	"sync"
)

// 接口定义
type Validatable interface {
	Validate() error
}

//type Response interface{}

type handlerMeta struct {
	once      sync.Once    // 确保只初始化一次
	meta      *runtimeMeta // 懒加载的元数据
	initError error        // 初始化错误
}

type runtimeMeta struct {
	fnValue     reflect.Value
	reqType     reflect.Type
	reqPool     *sync.Pool
	hasValidate bool
}

var (
	handlerCacheMap = make(map[string]*handlerMeta) // key: string -> *handlerMeta
)

// 运行时构建元数据
func buildRuntimeMeta(fn interface{}) (*runtimeMeta, error) {
	// 实际项目中需要根据key获取原始handler
	// 这里简化处理，使用固定示例handler
	rawHandler := fn

	// 反射分析handler
	funcType := reflect.TypeOf(rawHandler)
	if funcType.Kind() != reflect.Func {
		return nil, fmt.Errorf("必须为函数类型")
	}

	// 校验参数签名
	if funcType.NumIn() != 3 || funcType.NumOut() != 1 {
		return nil, fmt.Errorf("函数签名必须为func(*Context, *T, response.Response) error")
	}

	reqType := funcType.In(1)
	if reqType.Kind() != reflect.Ptr || reqType.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("第二个参数必须为结构体指针")
	}

	// 构造元数据
	meta := &runtimeMeta{
		fnValue: reflect.ValueOf(rawHandler),
		reqType: reqType,
		reqPool: &sync.Pool{
			New: func() interface{} {
				return reflect.New(reqType.Elem()).Interface()
			},
		},
		//resPool: &sync.Pool{
		//	New: func() interface{} {
		//		return &response.Data{}
		//	},
		//},
		//ctxPool: &sync.Pool{
		//	New: func() interface{} {
		//		return &Context{
		//			Context: context.Background(),
		//		}
		//	},
		//},
	}

	// 预生成JSON解析器
	//meta.jsonParser = func(data []byte) (interface{}, error) {
	//	req := meta.reqPool.Get()
	//	if err := json.Unmarshal(data, req); err != nil {
	//		qPool.Put(req)
	//		return nil, err
	//	}
	//	return req, nil
	//}

	// 检查Validate方法（带缓存）
	//if cached, ok := validateCache.Load(reqType); ok {
	//	meta.hasValidate = cached.(bool)
	//} else {
	//	_, meta.hasValidate = reflect.New(reqType.Elem()).Interface().(Validatable)
	//	validateCache.Store(reqType, meta.hasValidate)
	//}

	return meta, nil
}

// 实际调用逻辑
func doCall(method string, meta *runtimeMeta, ctx *Context, resp *response.Data, body interface{}) error {

	req := meta.reqPool.Get()
	var err error

	if body != nil {
		if method == "GET" {

			query, err1 := url.ParseQuery(body.(string))
			if err1 != nil {
				return err1
			}
			err1 = form.NewDecoder().Decode(req, query)
			if err1 != nil {
				return err1
			}
		} else {
			err = sonic.Unmarshal([]byte(body.(string)), req)
		}
	}

	// 解析请求
	//req, err := meta.jsonParser(body)
	if err != nil {
		return fmt.Errorf("JSON解析失败: %w", err)
	}
	defer meta.reqPool.Put(req)

	// 执行验证
	//if meta.hasValidate {
	//	if err := req.(Validatable).Validate(); err != nil {
	//		return fmt.Errorf("验证失败: %w", err)
	//	}
	//}

	if resp == nil {
		resp = new(response.Data)
	}
	fmt.Println("meta:", meta)
	fmt.Println("meta.fnValue:", meta.fnValue)

	// 反射调用
	results := meta.fnValue.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req), reflect.ValueOf(resp)})
	if err := results[0].Interface(); err != nil {
		return err.(error)
	}
	return nil
}
