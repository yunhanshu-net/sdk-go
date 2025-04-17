package excelx

import (
	"auth_delete/model"
	"fmt"
	"testing"
	"time"
)

type User struct {
	Name string  `excel:"name"`
	Age  *int    `excel:"age"`
	Desc string  `excel:"desc"`
	Sex  *string `excel:"sex"`
}

func TestName(t *testing.T) {
	var users []User
	err := Unmarshal("./测试.xlsx", &users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)
}

func TestName1(t *testing.T) {
	var access []model.Access
	now := time.Now()
	err := Unmarshal("./海外节点权限回收信息列表.xlsx", &access, "离职人员权限清理")
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Since(now))
}
