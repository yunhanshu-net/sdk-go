package office

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
	"time"
)

var taskDB = "tasks.db"

type Task struct {
	ID          int       `gorm:"primaryKey" runner:"code:id;name:任务ID"`
	Title       string    `runner:"code:title;name:任务标题"`
	Description string    `runner:"code:desc;name:任务描述"`
	Assignee    string    `runner:"code:assignee;name:负责人"`
	Deadline    time.Time `runner:"code:deadline;name:截止时间"`
	Status      string    `runner:"code:status;name:任务状态"`
}

func init() {
	taskConfig := &runner.ApiConfig{
		ChineseName: "任务管理",
		EnglishName: "taskManager",
		ApiDesc:     "办公任务管理系统",
		Tags:        "办公管理;任务管理",
		OnApiCreated: func(ctx *runner.Context, req *request.OnApiCreated) error {
			return runner.MustGetOrInitDB(taskDB).AutoMigrate(&Task{})
		},
		OnInputValidate: func(ctx *runner.Context, req *request.OnInputValidate) (*response.OnInputValidate, error) {
			if req.Key == "deadline" {
				if _, err := time.Parse("2006-01-02", req.Value); err != nil {
					return &response.OnInputValidate{Msg: "日期格式错误，请使用YYYY-MM-DD格式"}, nil
				}
			}
			return nil, nil
		},
	}

	runner.Get("/tasks", ListTasks, taskConfig)
	runner.Post("/tasks/create", CreateTask, taskConfig)
	runner.Post("/tasks/update", UpdateTaskStatus, taskConfig)
}

type TaskRequest struct {
	Title    string    `json:"title"`
	Desc     string    `json:"desc"`
	Assignee string    `json:"assignee"`
	Deadline time.Time `json:"deadline"`
	Status   string    `json:"status"` // 新增状态字段
	*request.PageInfo
}

func ListTasks(ctx *runner.Context, req *TaskRequest, resp response.Response) error {
	db := runner.MustGetOrInitDB(taskDB)
	var tasks []Task

	query := db.Model(&Task{})
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	return resp.Table(&tasks).
		AutoPaginated(query, &Task{}, req.PageInfo).
		Build()
}

func CreateTask(ctx *runner.Context, req *TaskRequest, resp response.Response) error {
	task := Task{
		Title:       req.Title,
		Description: req.Desc,
		Assignee:    req.Assignee,
		Deadline:    req.Deadline,
		Status:      "未开始",
	}

	if err := runner.MustGetOrInitDB(taskDB).Create(&task).Error; err != nil {
		logrus.Errorf("创建任务失败: %v", err)
		return err
	}
	return resp.Form(task).Build()
}

func UpdateTaskStatus(ctx *runner.Context, req *TaskRequest, resp response.Response) error {
	result := runner.MustGetOrInitDB(taskDB).
		Model(&Task{}).
		Where("title = ?", req.Title).
		Update("status", req.Status)

	if result.Error != nil {
		logrus.Errorf("更新任务状态失败: %v", result.Error)
		return result.Error
	}
	return resp.Form(map[string]interface{}{"updated_count": result.RowsAffected}).Build()
}
