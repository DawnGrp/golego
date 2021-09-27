package user

import (
	"context"
	"encoding/json"
	"fmt"
	"golego/modules/bootstrap"
	web "golego/modules/webserver"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/open-policy-agent/opa/rego"
)

func init() {
	bootstrap.AtRun(OPAInit)
	web.AtMiddleWave(setSession())
	web.AtMiddleWave(auth)
}

func setSession() func(c *gin.Context) {

	store := memstore.NewStore([]byte(sessionKey))
	store.Options(sessions.Options{
		MaxAge: 0, // seems this works
		Path:   "/",
	})

	return sessions.Sessions("auth", store)
}

func auth(c *gin.Context) {

	//先检查是否登入
	session := sessions.Default(c)

	user, ok := session.Get("user").(map[string]interface{})
	if user == nil || !ok {

		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusForbidden, gin.H{"error": "unsignin"})
		} else {
			c.Redirect(302, "/signin")
		}

		c.Abort()
		return
	}

	c.Request.ParseForm()
	params := map[string]interface{}{}
	for k, v := range c.Request.Form {
		params[k] = v[0]
	}

	OPAQuery(map[string]interface{}{
		"user":   user,
		"params": params,
		"method": c.Request.Method,
		"path":   strings.Trim(c.Request.URL.Path, "/"),
	}) //OPA 权限验证

	fmt.Println("c.Request.URL.Path", c.Request.URL.Path)

	c.Next()
}

var opaQuery rego.PreparedEvalQuery

const (
	auth_rego_path = "./asset/auth.rego"
	auth_package   = "data.auth"
)

func OPAInit() {

	// Construct a Rego object that can be prepared or evaluated.
	r := rego.New(
		rego.Query(auth_package),
		rego.Load([]string{auth_rego_path}, nil))
	// Create a prepared query that can be evaluated.
	var err error
	opaQuery, err = r.PrepareForEval(context.Background())
	if err != nil {
		panic(err)
	}

}

func OPAQuery(input interface{}) {
	rs, err := opaQuery.Eval(context.Background(), rego.EvalInput(input))
	if err != nil {
		log.Fatal("err", err)
	}

	// Do something with the result.
	rsJson, _ := json.Marshal(rs)

	fmt.Println("rs", string(rsJson))

}
