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

	user := session.Get("user")
	if user != nil {
		c.HTML(http.StatusOK, "user/signup", gin.H{})
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

	fmt.Println(" c.Request.URL.Path", c.Request.URL.Path)

	c.Next()
}

var opaQuery rego.PreparedEvalQuery

func OPAInit() {

	// Construct a Rego object that can be prepared or evaluated.
	r := rego.New(
		rego.Query("data.user"),
		rego.Load([]string{"/Users/zeta/workspace/golego/modules/user/example.rego"}, nil))
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
