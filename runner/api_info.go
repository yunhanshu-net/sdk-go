package runner

type ApiInfo struct {
	Router      string `json:"router"`
	Method      string `json:"method"`
	ApiDesc     string `json:"api_desc"`
	IsPublicApi bool   `json:"is_public_api"`
	//Labels      []string `json:"labels"`
	ChineseName string   `json:"chinese_name"`
	EnglishName string   `json:"english_name"`
	Classify    string   `json:"classify"`
	Tags        []string `json:"tags"`

	//form，table，
	RenderType string `json:"widget"`

	UseTables []interface{} `json:"use_tables"` //这里注册使用到的数据表
	UseDB     []string      `json:"use_db"`     //用到的db文件
	Request   interface{}   `json:"-"`
	Response  interface{}   `json:"-"`

	//用map的都是字段级别的回调，其他的都是接口级别回调

	OnPageLoad OnPageLoad `json:"-"`

	OnApiCreated    OnApiCreated    `json:"-"`
	OnApiUpdated    OnApiUpdated    `json:"-"`
	BeforeApiDelete BeforeApiDelete `json:"-"`
	AfterApiDeleted AfterApiDeleted `json:"-"`

	BeforeRunnerClose BeforeRunnerClose `json:"-"`
	AfterRunnerClose  AfterRunnerClose  `json:"-"`
	OnVersionChange   OnVersionChange   `json:"-"`

	OnInputFuzzy    OnInputFuzzy    `json:"-"` //这两个可以删掉
	OnInputValidate OnInputValidate `json:"-"` //这两个可以删掉

	OnTableDeleteRows OnTableDeleteRows `json:"-"`
	OnTableUpdateRow  OnTableUpdateRow  `json:"-"`
	OnTableSearch     OnTableSearch     `json:"-"`

	OnInputFuzzyMap    map[string]OnInputFuzzy    `json:"-"` //key是字段的code，字段级回调
	OnInputValidateMap map[string]OnInputValidate `json:"-"` //key是字段的code，字段级回调

}
