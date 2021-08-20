package auth

import (
	"golego/modules/helper"
	"golego/modules/webserver"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func init() {
	webserver.AddMiddleWaveHook(setSession())
}

func GetInfo() helper.Info {
	return helper.Info{
		Name: "auth",
	}
}

func check(c *gin.Context) {
}

func register(role string, args ...string) {

}

func setSession() func(c *gin.Context) {

	uuid := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	store := memstore.NewStore([]byte(uuid))
	store.Options(sessions.Options{
		MaxAge: 0, // seems this works
		Path:   "/",
	})

	return sessions.Sessions("auth", store)
}
