package user

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/runner"
)

func init() {
	getUserConfig := &runner.ApiConfig{
		OnCreated: func(ctx *runner.HttpContext) error {
			db, err := ctx.GetAndInitDB("app.db")
			if err != nil {
				return err
			}
			err = db.AutoMigrate(&User{})
			if err != nil {
				return err
			}
			return nil
		},
		AfterClose: func(ctx *runner.HttpContext) error {
			db, err := ctx.GetAndInitDB("app.db")
			if err != nil {
				return err
			}
			if err := db.Migrator().DropTable("user"); err != nil {
				return fmt.Errorf("删除表失败: %v", err)
			}
			return nil
		},
	}
	runner.Get("/user/get_user", GetUser, getUserConfig)
}

type User struct {
	ID       int    `json:"id,omitempty" gorm:"primary_key"`
	Username string `json:"username" form:"username"`
	Age      int    `json:"age" form:"age"`
}

func GetUser(ctx *runner.HttpContext) error {
	var user User
	err := ctx.Request.ShouldBindJSON(&user)
	if err != nil {
		return err
	}
	db, err := ctx.GetAndInitDB("app.db")
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	var findUser User
	db.Where("username = ?", user.Username).First(&findUser)

	if findUser.ID == 0 {
		db.Create(&user)
		findUser = user
	}
	return ctx.Response.JSON(findUser).Build()
}
