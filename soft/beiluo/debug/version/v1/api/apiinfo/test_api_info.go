package apiinfo

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
	"time"
)

var testDB = "apiinfo_test.db"

// TestModel 测试模型，只需要 gorm 标签
type TestModel struct {
	ID          int       `gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`
	Features    []string  `json:"features"`
	Priority    int       `json:"priority"`
	IsActive    bool      `json:"is_active"`
	Attachment  string    `json:"attachment"`
	CreateTime  time.Time `json:"create_time"`
}

// TestRequest 测试请求结构
type TestRequest struct {
	Name        string    `json:"name" runner:"widget:input;placeholder:请输入名称"`
	Description string    `json:"description" runner:"widget:input;mode:text_area"`
	Category    string    `json:"category" runner:"widget:select;options:类型A,类型B,类型C"`
	Status      string    `json:"status" runner:"widget:radio;options:启用,禁用;default_value:启用"`
	Features    []string  `json:"features" runner:"widget:checkbox;options:特性1,特性2,特性3"`
	Priority    int       `json:"priority" runner:"widget:slider;min:1;max:10;step:1"`
	IsActive    bool      `json:"is_active" runner:"widget:switch"`
	CreateTime  time.Time `json:"create_time"`
	*request.PageInfo
}

// TestResponse 测试响应结构
type TestResponse struct {
	ID          int       `json:"id" runner:"code:id;name:ID;widget:input"`
	Name        string    `json:"name" runner:"code:name;name:名称;widget:input;placeholder:请输入名称"`
	Description string    `json:"description" runner:"code:desc;name:描述;widget:input;mode:text_area"`
	Category    string    `json:"category" runner:"code:category;name:分类;widget:select;options:类型A,类型B,类型C"`
	Status      string    `json:"status" runner:"code:status;name:状态;widget:radio;options:启用,禁用;default_value:启用"`
	Features    []string  `json:"features" runner:"code:features;name:特性;widget:checkbox;options:特性1,特性2,特性3"`
	Priority    int       `json:"priority" runner:"code:priority;name:优先级;widget:slider;min:1;max:10;step:1"`
	IsActive    bool      `json:"is_active" runner:"code:is_active;name:是否激活;widget:switch"`
	CreateTime  time.Time `json:"create_time" runner:"code:create_time;name:创建时间"`
}

func init() {
	// 创建一个包含各种回调函数的API配置
	listConfig := &runner.ApiConfig{
		ChineseName: "API信息测试",
		EnglishName: "apiInfoTest",
		ApiDesc:     "用于测试getApiInfos功能的API",
		Tags:        "测试;API信息",
		RenderType:  "form", // 设置渲染类型
		UseTables:   []interface{}{&TestModel{}},
		Request:     &TestRequest{},
		Response:    &TestResponse{},

		// 页面加载回调
		OnPageLoad: func(ctx *runner.Context) (resetRequest interface{}, resp interface{}, err error) {
			return &TestRequest{
				Name:     "测试名称",
				Category: "类型A",
				Status:   "启用",
				IsActive: true,
				Priority: 5,
			}, nil, nil
		},

		// API 生命周期回调
		OnApiCreated: func(ctx *runner.Context, req *request.OnApiCreated) error {
			return runner.MustGetOrInitDB(testDB).AutoMigrate(&TestModel{})
		},
		BeforeApiDelete: func(ctx *runner.Context, req *request.BeforeApiDelete) error {
			logrus.Info("API即将被删除")
			return nil
		},
		AfterApiDeleted: func(ctx *runner.Context, req *request.AfterApiDeleted) error {
			return runner.MustGetOrInitDB(testDB).Migrator().DropTable(&TestModel{})
		},

		// 运行器生命周期回调
		BeforeRunnerClose: func(ctx *runner.Context, req *request.BeforeRunnerClose) error {
			logrus.Info("Runner即将关闭")
			return nil
		},
		AfterRunnerClose: func(ctx *runner.Context, req *request.AfterRunnerClose) error {
			logrus.Info("Runner已关闭")
			return nil
		},

		// 版本控制回调
		OnVersionChange: func(ctx *runner.Context, req *request.OnVersionChange) error {
			for _, change := range req.Change {
				logrus.Infof("[change]" + change.String())
			}
			return nil
		},

		// 输入交互回调
		OnInputFuzzy: func(ctx *runner.Context, req *request.OnInputFuzzy) (*response.OnInputFuzzy, error) {
			if req.Key == "category" {
				return &response.OnInputFuzzy{Values: []string{"类型A", "类型B", "类型C"}}, nil
			}
			return nil, nil
		},
		OnInputValidate: func(ctx *runner.Context, req *request.OnInputValidate) (*response.OnInputValidate, error) {
			if req.Key == "name" && len(req.Value) < 2 {
				return &response.OnInputValidate{Msg: "名称长度不能少于2个字符"}, nil
			}
			return nil, nil
		},

		// 表格操作回调
		OnTableDeleteRows: func(ctx *runner.Context, req *request.OnTableDeleteRows) (*response.OnTableDeleteRows, error) {
			logrus.Infof("删除表格行: %v", req.Ids)
			return nil, nil
		},
		OnTableUpdateRow: func(ctx *runner.Context, req *request.OnTableUpdateRow) (*response.OnTableUpdateRow, error) {
			logrus.Infof("更新表格行: %v", req.Ids)
			return nil, nil
		},
		OnTableSearch: func(ctx *runner.Context, req *request.OnTableSearch) (*response.OnTableSearch, error) {
			logrus.Infof("搜索表格: %v", req.Cond)
			return nil, nil
		},
	}
	createConfig := &runner.ApiConfig{
		ChineseName: "创建API信息",
		EnglishName: "createApiInfo",
		ApiDesc:     "创建新的API信息记录",
		Tags:        "测试;API信息",
		UseTables:   []interface{}{&TestModel{}},
		Request:     &TestRequest{},
		Response:    &TestResponse{},
	}
	updateConfig := &runner.ApiConfig{
		ChineseName: "更新API信息",
		EnglishName: "updateApiInfo",
		ApiDesc:     "更新已有的API信息记录",
		Tags:        "测试;API信息",
		UseTables:   []interface{}{&TestModel{}},
		Request:     &TestRequest{},
		Response:    &TestResponse{},
	}
	deleteConfig := &runner.ApiConfig{
		ChineseName: "删除API信息",
		EnglishName: "deleteApiInfo",
		ApiDesc:     "删除API信息记录",
		Tags:        "测试;API信息",
		UseTables:   []interface{}{&TestModel{}},
		Request:     &TestRequest{},
	}
	// 注册API路由
	runner.Get("/apiinfo/list", ListItems, listConfig)
	runner.Post("/apiinfo/create", CreateItem, createConfig)
	runner.Post("/apiinfo/update", UpdateItem, updateConfig)
	runner.Post("/apiinfo/delete", DeleteItem, deleteConfig)
}

// ListItems 列出所有项目
func ListItems(ctx *runner.Context, req *TestRequest, resp response.Response) error {
	db := runner.MustGetOrInitDB(testDB)
	var items []TestModel

	query := db.Model(&TestModel{})
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	return resp.Table(&items).
		AutoPaginated(query, &TestModel{}, req.PageInfo).
		Build()
}

// CreateItem 创建新项目
func CreateItem(ctx *runner.Context, req *TestRequest, resp response.Response) error {
	item := TestModel{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Status:      req.Status,
		Features:    req.Features,
		Priority:    req.Priority,
		IsActive:    req.IsActive,
		CreateTime:  time.Now(),
	}

	if err := runner.MustGetOrInitDB(testDB).Create(&item).Error; err != nil {
		logrus.Errorf("创建项目失败: %v", err)
		return err
	}

	return resp.Form(TestResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Category:    item.Category,
		Status:      item.Status,
		Features:    item.Features,
		Priority:    item.Priority,
		IsActive:    item.IsActive,
		CreateTime:  item.CreateTime,
	}).Build()
}

// UpdateItem 更新项目
func UpdateItem(ctx *runner.Context, req *TestRequest, resp response.Response) error {
	var item TestModel
	db := runner.MustGetOrInitDB(testDB)

	if err := db.Where("name = ?", req.Name).First(&item).Error; err != nil {
		logrus.Errorf("查找项目失败: %v", err)
		return err
	}

	// 更新字段
	item.Description = req.Description
	item.Category = req.Category
	item.Status = req.Status
	item.Features = req.Features
	item.Priority = req.Priority
	item.IsActive = req.IsActive

	if err := db.Save(&item).Error; err != nil {
		logrus.Errorf("更新项目失败: %v", err)
		return err
	}

	return resp.Form(TestResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Category:    item.Category,
		Status:      item.Status,
		Features:    item.Features,
		Priority:    item.Priority,
		IsActive:    item.IsActive,
		CreateTime:  item.CreateTime,
	}).Build()
}

// DeleteItem 删除项目
func DeleteItem(ctx *runner.Context, req *TestRequest, resp response.Response) error {
	result := runner.MustGetOrInitDB(testDB).
		Where("name = ?", req.Name).
		Delete(&TestModel{})

	if result.Error != nil {
		logrus.Errorf("删除项目失败: %v", result.Error)
		return result.Error
	}

	return resp.Form(map[string]interface{}{
		"deleted_count": result.RowsAffected,
	}).Build()
}
