package webserver

import (
	"fmt"
	"golego/modules/bootstrap"
	"golego/modules/config"
	"golego/modules/helper"
	"golego/utils"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

//HOOK_5. 在init函数中挂载钩子
func init() {
	helper.Register(me)
	bootstrap.AtRun(startServer)
	AtSetHandle(getActions)
}

//实现一个开放的GetInfo方法
var me = helper.ModuleInfo{
	Name:      "webserver",
	HumanName: "网络服务模块",
}

//HOOK_1. 定义钩子类型，钩子类型为一个函数类型
type set_router_hook func(r *gin.Engine)
type set_handle_hook func() (
	human_name string,
	paramsStructPtr interface{},
	method Method,
	path string,
	handlers gin.HandlerFunc)

//HOOK_2. 钩子组，接入的钩子必须支持多个，因此需要定一个数组
var set_router_hooks = []set_router_hook{}
var set_middleWave_hooks = []gin.HandlerFunc{}
var set_handle_hooks = []set_handle_hook{}

type action struct {
	Name   string
	Action string
	Params [][]string
}

var actions = map[string]action{} //不能用普通的map，用数组吧

type Method string

const (
	POST    Method = "POST"
	GET     Method = "GET"
	DELETE  Method = "DELETE"
	PATCH   Method = "PATCH"
	PUT     Method = "PUT"
	OPTIONS Method = "OPTIONS"
	HEAD    Method = "HEAD"
	ANY     Method = "ANY"
)

//HOOK_3. 提供挂入钩子的方法，其他模块可以将处理的函数添加到钩子组中
func AtSetRouter(h set_router_hook) {
	set_router_hooks = append(set_router_hooks, h)
}
func AtMiddleWave(h gin.HandlerFunc) {
	set_middleWave_hooks = append(set_middleWave_hooks, h)
}
func AtSetHandle(h set_handle_hook) {
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

		n, s, m, p, f := hook()

		_, ok := actions[n]
		if ok {
			panic(fmt.Errorf("%s exist", n))
		}

		_, fields := utils.GetStruct(s)
		actions[n] = action{
			Name:   n,
			Action: fmt.Sprintf("%s:%s", m, p),
			Params: fields,
		}

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
	cfg, ok := config.Get(me.Name)
	if !ok {
		cfg = gjson.Parse(`{"addr":":8082"}`)
		config.Set(me.Name, cfg)
	}

	helper.WaitGroup.Add(1)
	go func() {
		router.Run(cfg.Get("addr").String())
		helper.WaitGroup.Done()
	}()
}

func Actions() map[string]action {
	return actions
}

func getActions() (name string, paramsStructPtr interface{}, method Method, path string, handlers gin.HandlerFunc) {
	return "获取参数列表", nil, GET, "/getactions", func(c *gin.Context) {
		c.JSON(200, gin.H{"actions": Actions()})
	}
}
