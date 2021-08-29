package auth

import (
	"fmt"
	"testing"

	"github.com/casbin/casbin/v2"
)

func TestFunc(t *testing.T) {

	// sub := map[string]interface{}{
	// 	"name": "zeta",
	// 	"age":  18,
	// }

	sub := "zeta"
	obj := map[string]interface{}{
		"role":      "admin",
		"limit_age": 20,
	}
	act := "read" // 用户对资源执行的操作。

	e := Enforcer()

	ok, err := e.Enforce(sub, obj, act)

	if err != nil {
		panic(err)
	}
	fmt.Println(ok)
}

func Enforcer() *casbin.Enforcer {

	// 创建enforcer
	e, err := casbin.NewEnforcer("/Users/zeta/workspace/golego/modules/auth/model.conf",
		"/Users/zeta/workspace/golego/modules/auth/policy.csv")

	if err != nil {
		fmt.Println(err)
	}
	return e
}
