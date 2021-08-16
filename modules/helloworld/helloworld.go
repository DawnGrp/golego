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
	ginserver.AddSetHandleHook(sayHello)
}

func sayHello() (string, string, gin.HandlerFunc) {

	return "get", "/", func(c *gin.Context) {
		c.String(200, "Hello!!")
	}

}
