package ginserver

import (
	"fmt"
	"golego/modules/bootstrap"
	"golego/modules/helper"

	"github.com/gin-gonic/gin"
)

//HOOK_5. 在init函数中挂载钩子
func init() {
	bootstrap.AddRunHook(startServer)
}

//实现一个开放的GetInfo方法
func GetInfo() helper.Info {
	return helper.Info{
		Name: "ginserver",
	}
}

//HOOK_1. 定义钩子类型，钩子类型为一个函数类型
type set_router_hook func(r *gin.Engine)
type set_handle_hook func() (method string, path string, handlers gin.HandlerFunc)

//HOOK_2. 钩子组，接入的钩子必须支持多个，因此需要定一个数组
var set_router_hooks = []set_router_hook{}
var set_middleWave_hooks = []gin.HandlerFunc{}
var set_handle_hooks = []set_handle_hook{}

//HOOK_3. 提供挂入钩子的方法，其他模块可以将处理的函数添加到钩子组中
func AddSetRouterHook(h set_router_hook) {
	set_router_hooks = append(set_router_hooks, h)
}
func AddMiddleWaveHook(h gin.HandlerFunc) {
	set_middleWave_hooks = append(set_middleWave_hooks, h)
}
func AddSetHandleHook(h set_handle_hook) {
	set_handle_hooks = append(set_handle_hooks, h)
}

func startServer(debug bool, addr ...string) {
	fmt.Println("server run", debug)
	router := gin.Default()

	//在所有处理函数之前买下中间件钩子
	for _, hook := range set_middleWave_hooks {
		router.Use(hook)
	}

	//处理函数钩子
	for _, hook := range set_handle_hooks {
		m, p, f := hook()
		switch {
		case m == "get":
			router.GET(p, f)
		case m == "post":
			router.POST(p, f)
		}
	}

	//其他钩子
	for _, hook := range set_router_hooks {
		hook(router)
	}

	helper.WaitGroup.Add(1)
	go func() {
		router.Run(addr...)
		helper.WaitGroup.Done()
	}()
}
