package bootstrap

import "golego/modules/helper"

//HOOK_0. 为函数提供一个描述函数
func GetInfo() helper.Info {
	return helper.Info{
		Name: "bootstrap",
	}
}

//HOOK_1. 定义钩子类型，钩子类型为一个函数类型
type run_hook func(debug bool, addr ...string)

//HOOK_2. 钩子组，接入的钩子必须支持多个，因此需要定一个数组
var run_hooks = []run_hook{}

//HOOK_3. 提供挂入钩子的方法，其他模块可以将处理的函数添加到钩子组中
func AddRunHook(h run_hook) {
	run_hooks = append(run_hooks, h)
}

func Run() {
	//HOOK_4. 埋钩子，实现所有挂入钩子组的钩子函数
	for _, hook := range run_hooks {
		hook(helper.Config.Debug, helper.Config.Addr...)
	}
}
