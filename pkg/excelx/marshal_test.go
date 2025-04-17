package excelx

import (
	"testing"
	"time"
)

type Employee struct {
	Name     string    `excel:"姓名"`
	Age      int       `excel:"年龄"`
	Salary   float64   `excel:"薪资"`
	JoinedAt time.Time `excel:"入职时间"`
	IsAdmin  *bool     `excel:"管理员"`
}

func TestM(t *testing.T) {
	admin := true
	data := []Employee{
		{
			Name:     "张三",
			Age:      28,
			Salary:   15000.50,
			JoinedAt: time.Now(),
			IsAdmin:  &admin,
		},
		{
			Name:     "李四",
			Age:      35,
			Salary:   20000.00,
			JoinedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			IsAdmin:  nil,
		},
	}

	if err := Marshal("./employees.xlsx", data); err != nil {
		panic(err)
	}

}
