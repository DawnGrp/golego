package bootstrap

import "golego/modules/helper"

var me = helper.ModuleInfo{
	Name:      "bootstrap",
	HumanName: "启动模块",
}

func init() {
	helper.Register(me)
}

//HOOK_1. 定义钩子类型，钩子类型为一个函数类型
type before_run_hook func()
type run_hook func()
type after_run_hook func()

//HOOK_2. 钩子组，接入的钩子必须支持多个，因此需要定一个数组
var before_run_hooks = []before_run_hook{}
var run_hooks = []run_hook{}
var after_run_hooks = []after_run_hook{}

//HOOK_3. 提供挂入钩子的方法，其他模块可以将处理的函数添加到钩子组中
func AtBeforeRun(h before_run_hook) {
	before_run_hooks = append(before_run_hooks, h)
}

func AtRun(h run_hook) {
	run_hooks = append(run_hooks, h)
}

func AtAfterRun(h after_run_hook) {
	after_run_hooks = append(after_run_hooks, h)
}

func BeforeRun() {
	//HOOK_4. 埋钩子，实现所有挂入钩子组的钩子函数
	for _, hook := range before_run_hooks {
		hook()
	}
}

func Run() {
	//HOOK_4. 埋钩子，实现所有挂入钩子组的钩子函数
	for _, hook := range run_hooks {
		hook()
	}
}

func AfterRun() {
	for _, hook := range after_run_hooks {
		hook()
	}
}
