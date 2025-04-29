package crm

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
)

var crmDB = "crm.db"

type Customer struct {
	ID       int    `gorm:"primaryKey" runner:"code:id;name:客户ID"`
	Name     string `runner:"code:name;name:客户名称"`
	Industry string `runner:"code:industry;name:所属行业"`
	Contact  string `runner:"code:contact;name:联系人"`
	Phone    string `runner:"code:phone;name:联系电话"`
	SalesRep string `runner:"code:sales_rep;name:销售代表"`
}

func init() {
	crmConfig := &runner.ApiConfig{
		ChineseName: "客户管理",
		EnglishName: "customerManager",
		ApiDesc:     "企业客户关系管理系统",
		Tags:        "CRM;客户管理",
		OnApiCreated: func(ctx *runner.Context, req *request.OnApiCreated) error {
			return runner.MustGetOrInitDB(crmDB).AutoMigrate(&Customer{})
		},
		OnInputFuzzy: func(ctx *runner.Context, req *request.OnInputFuzzy) (*response.OnInputFuzzy, error) {
			if req.Key == "industry" {
				var industries []string
				runner.MustGetOrInitDB(crmDB).
					Model(&Customer{}).
					Where("industry LIKE ?", "%"+req.Value+"%").
					Pluck("DISTINCT industry", &industries)
				return &response.OnInputFuzzy{Values: industries}, nil
			}
			return nil, nil
		},
	}

	runner.Get("/crm/customers", ListCustomers, crmConfig)
	runner.Post("/crm/customers/add", AddCustomer, crmConfig)
}

type CustomerRequest struct {
	Name     string `json:"name"`
	Industry string `json:"industry"`
	Contact  string `json:"contact"`
	Phone    string `json:"phone"`
	*request.PageInfo
}

func ListCustomers(ctx *runner.Context, req *CustomerRequest, resp response.Response) error {
	db := runner.MustGetOrInitDB(crmDB)
	var customers []Customer

	query := db.Model(&Customer{})
	if req.Industry != "" {
		query = query.Where("industry = ?", req.Industry)
	}

	return resp.Table(&customers).
		AutoPaginated(query, &Customer{}, req.PageInfo).
		Build()
}

func AddCustomer(ctx *runner.Context, req *CustomerRequest, resp response.Response) error {
	customer := Customer{
		Name:     req.Name,
		Industry: req.Industry,
		Contact:  req.Contact,
		Phone:    req.Phone,
		SalesRep: ctx.GetUsername(),
	}

	if err := runner.MustGetOrInitDB(crmDB).Create(&customer).Error; err != nil {
		logrus.Errorf("添加客户失败: %v", err)
		return err
	}
	return resp.Form(customer).Build()
}
