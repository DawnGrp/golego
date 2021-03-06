package helper

import (
	"fmt"
	"sync"
)

type ModuleInfo struct {
	Name        string
	HumanName   string
	Description string
	Templates   []string
}

//模块信息Map，用来保存所有模块的信息
//模块主动提供公共函数GetInfo，返回Info类型的对象以描述模块的信息
//在moudules包中，引用所有模块，并且调用GetInfo函数，并将模块信息写入到该变量中。
var moduleInfos = map[string]ModuleInfo{}

//公共的线程组
//在modules中，调用等待函数
//其他需要使用到异步的模块函数可以通过调用 helper.WaitGroup.Add(1) 方法添加需要等待的异步进程
//在异步进程结束时，调用helper.WaitGroup.Done() 向WaitGroup表达线程已经结束。
var WaitGroup sync.WaitGroup

var me = ModuleInfo{
	Name:      "helper",
	HumanName: "助手模块",
}

func Register(mi ModuleInfo) {
	if _, ok := moduleInfos[mi.Name]; ok {
		panic(fmt.Errorf("module [%s] is exist", mi.Name))
	}

	moduleInfos[mi.Name] = mi
}

func GetModuleInfos() map[string]ModuleInfo {
	return moduleInfos
}
