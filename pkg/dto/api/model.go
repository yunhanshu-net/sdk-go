package api

type Config struct {
}

type Info struct {
	Router      string   `json:"router"`
	Method      string   `json:"method"`
	User        string   `json:"user"`
	Runner      string   `json:"runner"`
	ApiDesc     string   `json:"api_desc"`
	ChineseName string   `json:"chinese_name"`
	EnglishName string   `json:"english_name"`
	Classify    string   `json:"classify"`
	Tags        []string `json:"tags"`
	//输入参数
	ParamsIn interface{} `json:"params_in"`
	//输出参数
	ParamsOut interface{} `json:"params_out"`
	UseTables []string    `json:"use_tables"`
	UseDB     []string    `json:"use_db"`
	Callbacks []string    `json:"callbacks"`
}

type ApiLogs struct {
	Version string  `json:"version"`
	Apis    []*Info `json:"apis"`
}
