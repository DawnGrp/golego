package helloworld

import (
	"golego/modules/ginserver"
	"golego/modules/helper"

	"github.com/gin-gonic/gin"
)

func GetInfo() helper.Info {
	return helper.Info{
		Name:        "HelloWorld",
		Description: "测试",
	}
}

func init() {
	//挂钩子
	ginserver.AddSetHandleHook(sayHello)
}

func sayHello() (ginserver.RequestMethod, string, gin.HandlerFunc) {

	return ginserver.GET, "/", func(c *gin.Context) {
		c.String(200, "Hello!!")
	}

}
