package hook

import (
	"fmt"
)

func Testmain() {

	//模块 创建钩子
	hook := CreateHook(*new(Thetype))

	//子模块 挂钩子
	hook.Up(Thetype(print))

	//模块 埋钩子
	hook.In(func(f HookType) {
		f.(Thetype)("aaaa")
	})
}

type Thetype func(string)

//定义到包

//实现
func print(str string) {
	fmt.Println(str)
}
func print2(str string, index int) {
	fmt.Println(str, index)
}
