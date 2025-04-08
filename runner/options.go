package runner

import (
	"fmt"
	"reflect"
	"strings"
)

type Option func(*ApiConfig)
type FuncParam struct {
	Code          string `json:"code,omitempty"`
	Desc          string `json:"desc,omitempty"`
	Mode          string `json:"mode,omitempty"`
	Type          string `json:"type,omitempty"`
	Value         string `json:"value,omitempty"`
	Options       string `json:"options,omitempty"`
	Required      string `json:"required,omitempty"`
	MockData      string `json:"mock_data,omitempty"`
	InputMode     string `json:"input_mode,omitempty"`
	TextLimit     string `json:"text_limit,omitempty"`
	NumberLimit   string `json:"number_limit,omitempty"`
	SelectOptions string `json:"select_options,omitempty"`
	FileSizeLimit string `json:"file_size_limit,omitempty"`
	FileTypeLimit string `json:"file_type_limit,omitempty"`
	IsTableField  bool   `json:"is_table_field"`
}

func NewConfig(opts ...Option) *ApiConfig {
	config := &ApiConfig{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

type ApiConfig struct {
	Router      string      `json:"router"`
	Method      string      `json:"method"`
	ApiDesc     string      `json:"api_desc"`
	IsPublicApi bool        `json:"is_public_api"`
	Labels      []string    `json:"labels"`
	ChineseName string      `json:"chinese_name"`
	EnglishName string      `json:"english_name"`
	Classify    string      `json:"classify"`
	Tags        string      `json:"tags"`
	ParamsIn    []FuncParam `json:"params_in"`
	ParamsOut   []FuncParam `json:"params_out"`

	Request  interface{} `json:"-"`
	Response interface{} `json:"-"`

	//创建新的api时候的回调函数,新增一个api假如新增了一张user表，
	//可以在这里用gorm的db.AutoMigrate(&User)来创建表，保证新版本的api可以正常使用新增的表
	//这个api只会在我创建这个api的时候执行一次
	OnCreated func(ctx *HttpContext) error `json:"-"`

	//api删除后触发回调，比如新增了一个字段，可以在这里做一些操作
	AfterDelete func(ctx *HttpContext) error `json:"-"`

	//每次版本发生变更都会回调这个函数（新增/删除api）
	OnVersionChange func(ctx *HttpContext) error `json:"-"`

	//程序结束前的回调函数，可以在程序结束前做一些操作，比如上报一些数据
	BeforeClose func(ctx *HttpContext) error `json:"-"`
	//程序结束后的回调函数，可以在程序结束后做一些操作，比如清理某些文件
	AfterClose func(ctx *HttpContext) error `json:"-"`
}

func getRunnerTag(runnerTag string) FuncParam {
	funcP := FuncParam{}
	split := strings.Split(runnerTag, ";")
	mp := make(map[string]string)
	for _, s := range split {
		vals := strings.Split(s, ":")
		if len(vals) == 2 {
			mp[vals[0]] = vals[1]
		}
	}
	valueOf := reflect.ValueOf(&funcP).Elem()
	typeOf := reflect.TypeOf(funcP)
	for i := 0; i < valueOf.NumField(); i++ {
		field := typeOf.Field(i)
		if !field.IsExported() {
			continue
		}

		value := valueOf.Field(i)
		tag := field.Tag.Get("json")
		if v := mp[strings.Split(tag, ",")[0]]; v != "" {
			if value.CanSet() {
				value.SetString(v)
			}
		}
	}
	return funcP
}

func (c *ApiConfig) getParams(p interface{}, mode string) (params []FuncParam, err error) {
	if p == nil {
		return nil, nil
	}
	val := reflect.ValueOf(p)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		fmt.Println("input is not a struct", val.Kind())
		return nil, fmt.Errorf("input is not a struct")
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		par := FuncParam{Mode: mode, Code: typ.Field(i).Name}
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			par.Code = strings.Split(jsonTag, ",")[0]
		}
		//runnerTag := field.Tag.Get("runner")
		runnerTag := field.Tag.Get("runner")
		if runnerTag == "" {

		} else {
			p1 := par
			par = getRunnerTag(runnerTag)
			if par.Mode == "" {
				par.Mode = p1.Mode
			}
			if par.Code == "" {
				par.Code = p1.Code
			}
		}

		if par.Type == "" {
			switch field.Type.Kind() {
			case reflect.Float32, reflect.Float64:
				par.Type = "float"
			case reflect.String:
				par.Type = "string"
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				par.Type = "number"
			default:
				par.Type = "string"
			}
		}

		//pp := getRunnerTag(runnerTag)

		params = append(params, par)
	}
	return params, nil
}

func (c *ApiConfig) GetParams() ([]FuncParam, error) {
	var list []FuncParam
	if c.Request != nil {
		params, err := c.getParams(c.Request, "in")
		if err != nil {
			return nil, err
		}
		list = append(list, params...)
	}
	if c.Response != nil {
		params, err := c.getParams(c.Response, "out")
		if err != nil {
			return nil, err
		}
		list = append(list, params...)
	}
	return list, nil
}

func WithPublicApi() Option {
	return func(c *ApiConfig) {
		c.IsPublicApi = true
	}
}
func WithApiDesc(apiDesc string) Option {
	return func(c *ApiConfig) {
		c.ApiDesc = apiDesc
	}
}
