package calc

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
	"strconv"
)

var dbName = "test.db"

func init() {
	getConfig := &runner.ApiConfig{
		ChineseName: "获取计算记录",
		EnglishName: "calcGet",
		ApiDesc:     "这里可以描述的详细一点",
		Tags:        "数据管理;数据分析;记录管理",
		OnApiCreated: func(ctx *runner.Context, req *request.OnApiCreated) error {
			db := runner.MustGetOrInitDB(dbName) //这里会返回*gorm.DB
			return db.AutoMigrate(&Calc{})
		},
		AfterApiDeleted: func(ctx *runner.Context, req *request.AfterApiDeleted) error {
			//视情况决定是否要删除表
			return runner.MustGetOrInitDB(dbName).Migrator().DropTable(&Calc{})
		},

		OnPageLoad: func(ctx *runner.Context) (resetRequest interface{}, resp interface{}, err error) {
			//这里返回用户的个人信息
			return &AddReq{A: 1, B: 2, Receiver: ctx.GetUsername()}, nil, nil
		},

		OnInputFuzzy: func(ctx *runner.Context, req *request.OnInputFuzzy) (*response.OnInputFuzzy, error) {
			var values []string
			if req.Key == "a" { //
				db := ctx.MustGetOrInitDB(dbName)
				var calcs []Calc
				db.Model(&Calc{}).Where("a  like ?", "%"+req.Value+"%").Limit(10).Find(&calcs)
				for _, calc := range calcs {
					values = append(values, strconv.Itoa(calc.A))
				}

			}
			return &response.OnInputFuzzy{Values: values}, nil
		},
	}

	addConfig := &runner.ApiConfig{
		Tags:        "数据管理;数据分析;记录管理",
		EnglishName: "calcAdd",
		ChineseName: "添加计算记录",
		ApiDesc:     "这里可以描述的详细一点",
		UseTables:   []interface{}{&Calc{}}, //这里会在注册这个api的时候自动创建相关的表
		OnPageLoad: func(ctx *runner.Context) (resetRequest interface{}, resp interface{}, err error) {
			return &AddReq{Receiver: ctx.GetUsername()}, nil, nil
		},
		OnInputValidate: func(ctx *runner.Context, req *request.OnInputValidate) (*response.OnInputValidate, error) {
			msg := ""
			if req.Key == "code" {
				if len(req.Value) > 64 {
					msg = "最长不能超过64个字符"
				}
				//其他判断......
			}
			return &response.OnInputValidate{Msg: msg}, nil
		},
	}

	runner.Get("/calc/add", Add, addConfig)
	runner.Get("/calc/get", Get, getConfig)
}

type Calc struct {
	ID       int    `gorm:"primaryKey;autoIncrement" runner:"code:id;name:id"`
	A        int    `json:"a" runner:"code:a;name:a"`
	B        int    `json:"b" runner:"code:b;name:b"`
	C        int    `json:"c" runner:"code:c;name:c"`
	Receiver string `json:"receiver" runner:"code:receiver;name:receiver"`
	Code     string `json:"code" runner:"code:code;name:code"`
}

type AddReq struct {
	Receiver string `json:"receiver"`
	A        int    `json:"a" form:"a"`
	B        int    `json:"b" form:"b"`
	Code     string `json:"code" form:"code"`
}
type GetReq struct {
	ID int `json:"id" form:"id"`
	*request.PageInfo
}

type AddResp struct {
	ID     int `json:"id"`
	Result int `json:"result"`
}

// Add 拿这个处理函数举例，ctx是固定参数， req *AddReq是用户自定义的参数，根据接口请求参数自己定义，resp response.Response是固定参数，用户可以根据这个返回自己的json数据
func Add(ctx *runner.Context, req *AddReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB(dbName)
	res := Calc{A: req.A, B: req.B, C: req.A + req.B} //这里模拟处理逻辑
	err := db.Model(&Calc{}).Create(&res).Error
	if err != nil {
		logrus.Errorf("Add err:%s", err.Error())
		return err
	}
	return resp.JSON(res).Build()
}

func Get(ctx *runner.Context, req *GetReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB(dbName)
	var res []Calc
	db.Where("id > ?", req.ID)

	//这里会返回table类型的数据，前端可以直接渲染成element的表格进行展示
	//AutoPaginated 会自动把查询到的数据挂载到res上，自动添加分页等等
	return resp.Table(&res).AutoPaginated(db, &Calc{}, req.PageInfo).Build()
}
