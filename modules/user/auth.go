package user

import (
	"golego/modules/webserver"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	webserver.AtMiddleWave(auth)
}

type auth_hook func(*gin.Context) bool

var auth_hooks []auth_hook

func AtAuth(h auth_hook) {
	auth_hooks = append(auth_hooks, h)
}

func auth(c *gin.Context) {

	//先检查是否登入
	session := sessions.Default(c)

	account, ok := session.Get("account").(string)
	if !ok || account == "" {
		c.HTML(http.StatusOK, "user/signup", gin.H{})
		return
	}

	for _, h := range auth_hooks {
		if !h(c) {
			return
		}
	}

	c.Next()
}
