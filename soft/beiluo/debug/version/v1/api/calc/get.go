package calc

import (
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
	"strconv"
)

type GetReq struct {
	ID int `json:"id" form:"id"`
	*request.PageInfo
}

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
	runner.Get("/calc/get", Get, getConfig)
}

func Get(ctx *runner.Context, req *GetReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB(dbName)
	var res []Calc
	db.Where("id > ?", req.ID)

	//这里会返回table类型的数据，前端可以直接渲染成element的表格进行展示
	//AutoPaginated 会自动把查询到的数据挂载到res上，自动添加分页等等
	return resp.Table(&res).AutoPaginated(db, &Calc{}, req.PageInfo).Build()
}
