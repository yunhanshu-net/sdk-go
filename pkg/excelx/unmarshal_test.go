package excelx

import (
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

type Access struct {
	Id             int    `excel:"id"`
	ProductPid     string `excel:"product_pid"`
	SubjectType    int    `excel:"subject_type"`
	Subject        string `excel:"subject"`
	Username       string `excel:"username"`
	Dept1          string `excel:"dept1"`
	Dept1Code      string `excel:"dept1_code"`
	Department     string `excel:"department"`
	DepartmentCode string `excel:"department_code"`
	Status         string `excel:"离职状态"`
	ObjectType     int    `excel:"object_type"`
	Object         string `excel:"object"`
	RoleRid        string `excel:"role_rid"`
	RoleProductPid string `excel:"role_product_pid"`
	CreatedAt      string `excel:"created_at"`
	UpdatedAt      string `excel:"updated_at"`
	ValidityType   int    `excel:"validity_type"`
	BeginTime      string `excel:"begin_time"`
	EndTime        string `excel:"end_time"`
	IsDeleted      int    `excel:"is_deleted"`
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
	var access []Access
	now := time.Now()
	err := Unmarshal("./海外节点权限回收信息列表.xlsx", &access, "离职人员权限清理")
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Since(now))
}
