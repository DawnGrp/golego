package webserver

import (
	"golego/modules/bootstrap"
	"golego/modules/config"
	"golego/modules/helper"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

//HOOK_5. 在init函数中挂载钩子
func init() {
	bootstrap.AddRunHook(startServer)
}

//实现一个开放的GetInfo方法
func GetInfo() helper.Info {
	return helper.Info{
		Name: "webserver",
	}
}

//HOOK_1. 定义钩子类型，钩子类型为一个函数类型
type set_router_hook func(r *gin.Engine)
type set_handle_hook func() (method RequestMethod, path string, handlers gin.HandlerFunc)

//HOOK_2. 钩子组，接入的钩子必须支持多个，因此需要定一个数组
var set_router_hooks = []set_router_hook{}
var set_middleWave_hooks = []gin.HandlerFunc{}
var set_handle_hooks = []set_handle_hook{}

type RequestMethod string

const (
	POST    RequestMethod = "POST"
	GET     RequestMethod = "GET"
	DELETE  RequestMethod = "DELETE"
	PATCH   RequestMethod = "PATCH"
	PUT     RequestMethod = "PUT"
	OPTIONS RequestMethod = "OPTIONS"
	HEAD    RequestMethod = "HEAD"
	ANY     RequestMethod = "ANY"
)

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

func startServer() {

	router := gin.Default()

	//在所有处理函数之前买下中间件钩子
	for _, hook := range set_middleWave_hooks {
		router.Use(hook)
	}

	//处理函数钩子
	for _, hook := range set_handle_hooks {
		m, p, f := hook()
		switch {
		case m == POST:
			router.POST(p, f)
		case m == GET:
			router.GET(p, f)
		case m == DELETE:
			router.DELETE(p, f)
		case m == PATCH:
			router.PATCH(p, f)
		case m == PUT:
			router.PUT(p, f)
		case m == OPTIONS:
			router.OPTIONS(p, f)
		case m == HEAD:
			router.HEAD(p, f)
		case m == ANY:
			router.Any(p, f)
		}
	}

	//其他钩子
	for _, hook := range set_router_hooks {
		hook(router)
	}

	//获得本模块的配置
	//如果不存在，则写入一个默认配置
	cfg, ok := config.Get(GetInfo().Name)
	if !ok {
		cfg = gjson.Parse(`{"addr":":8082"}`)
		config.Add(GetInfo().Name, cfg)
	}

	helper.WaitGroup.Add(1)
	go func() {
		router.Run(cfg.Get("addr").String())
		helper.WaitGroup.Done()
	}()
}
