package user

import (
	"golego/modules/webserver"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func init() {
	webserver.AtMiddleWave(setSession())
}
func setSession() func(c *gin.Context) {

	store := memstore.NewStore([]byte(sessionKey))
	store.Options(sessions.Options{
		MaxAge: 0, // seems this works
		Path:   "/",
	})

	return sessions.Sessions("auth", store)
}
