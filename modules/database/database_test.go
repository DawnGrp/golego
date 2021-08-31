package database

import (
	"fmt"
	"golego/utils"
	"testing"
)

type mystruct struct {
	Name string `json:"name" form:"name" name:"姓名"`
	Age  int    `json:"age" form:"age" name:"年龄"`
}

type mystruct2 struct {
	mystruct
	Role string `json:"role" form:"role" name:"角色"`
}

type istruct interface {
	Say()
}

func (m1 *mystruct) Say() {}

func TestDB(t *testing.T) {
	m := mystruct2{}
	m.Name = "1"
	m.Age = 1

	func() {
		a := utils.GetStruct(&m)
		fmt.Println(a)
	}()
}
