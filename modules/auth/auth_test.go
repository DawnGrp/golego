package auth

import (
	"fmt"
	"runtime"
	"testing"

	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

func TestFunc(t *testing.T) {
	i := 0
	fmt.Println("i =", i)
	fmt.Println("FuncName1 =", func() string { return ff() }())

}

func ff() string {
	return runFuncName()
}

func runFuncName() string {

	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc).Name()

	return f
}

func Enforcer() *casbin.Enforcer {

	// 从字符串初始化模型
	text :=
		`
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`
	m, _ := model.NewModelFromString(text)

	// 从CSV文件adapter加载策略规则
	// 使用自己的 adapter 替换
	a := fileadapter.NewAdapter("")

	// 创建enforcer
	e, err := casbin.NewEnforcer(m, a)

	if err != nil {
		return nil
	}

	return e
}
