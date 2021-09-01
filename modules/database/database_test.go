package database

import (
	"fmt"
	"golego/utils"
	"reflect"
	"testing"
)

type mystruct struct {
	Name string `json:"name" form:"name" name:"姓名" option:"a,b,c,d,e"`
	Age  int    `json:"age" form:"age" name:"年龄" option:"1,2,3,4,5"`
}

func (m1 *mystruct) Say() {}

func TestDB(t *testing.T) {

	n, a := utils.GetStruct(new(mystruct))

	fmt.Println(n)

	for _, v := range a {
		fmt.Println(a[0], a[1])
		if len(v) > 2 {
			option := reflect.StructTag(v[2]).Get("option")
			fmt.Println(option)
		}
	}

}
