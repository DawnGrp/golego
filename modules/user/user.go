package user

import (
	"context"
	db "golego/modules/database"
	"golego/modules/helper"
	"golego/modules/metadata"
	web "golego/modules/webserver"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var me = helper.ModuleInfo{
	Name:      "user",
	HumanName: "用户模块",
	Templates: []string{
		"signup", "signin", "signout",
	},
}

//登入后设置会话的钩子，返回一个键和值，用于session.Set函数设置会话
type set_session_hook func(map[string]interface{}) (string, interface{})

var set_session_hooks []set_session_hook

func AtSetSession(h set_session_hook) {
	set_session_hooks = append(set_session_hooks, h)
}

//设置用户集合的钩子
type set_metadata_hook func() ([]metadata.Filed, []metadata.Index)

var set_metadata_hooks []set_metadata_hook

func AtUpdateMetadata(h set_metadata_hook) {
	set_metadata_hooks = append(set_metadata_hooks, h)
}

func init() {
	helper.Register(me)
	db.RegisterC(me.Name)
	db.AtConnected(updateUserMetadata)
	web.AtMiddleWave(setSession())
	web.AtMiddleWave(auth)
	web.AtSetHandle(signinGet)
	web.AtSetHandle(signinPost)
	web.AtSetHandle(signupGet)
	web.AtSetHandle(signupPost)
	web.AtSetHandle(signout)
}

//初始化会话的key，写在这里，会在每次启动的时候都用不同的会话键
var sessionKey = strings.Replace(uuid.NewV4().String(), "-", "", -1)

//创建user集合以及索引
func updateUserMetadata() {

	md := metadata.Metadata{
		Name:      me.Name, //集合名称，等于本模块名称
		HumanName: me.HumanName,
		Fileds: []metadata.Filed{
			{
				Name: "account",
				Type: metadata.FiledType_String,
			},
			{
				Name: "password",
				Type: metadata.FiledType_String,
			},
		},
		Indexs: []metadata.Index{
			{
				Type:   metadata.IndexType_Unique,
				Fileds: []string{"account"},
			},
		},
	}

	for _, h := range set_metadata_hooks {
		fs, is := h()

		md.Fileds = append(md.Fileds, fs...)
		md.Indexs = append(md.Indexs, is...)
	}

	_, err := metadata.Replace(md, true) //使用替换模式，主要是为了每次启动时，如果子模块的钩子有增减时也能立即生效
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		panic(err)
	}

}

func signinGet() (name string, paramsStructPtr interface{}, method web.Method, path string, handlers gin.HandlerFunc) {

	return "登入页面", nil, web.GET, "/signin",
		func(c *gin.Context) {

			c.HTML(http.StatusOK, "user/login", gin.H{})

		}
}

func signinPost() (name string, paramsStructPtr interface{}, method web.Method, path string, handlers gin.HandlerFunc) {
	type input struct {
		Account  string `json:"account" form:"account" name:"账户"`
		Password string `json:"password" form:"password" name:"密码"`
	}
	return "登入执行", new(input), web.POST, "/signin",
		func(c *gin.Context) {

			account := c.PostForm("account")
			password := c.PostForm("password")
			result := gin.H{
				"account": account, "password": password, "err": "",
			}

			defer c.HTML(http.StatusOK, "user/login", result)

			r := db.C(me.Name).FindOne(context.Background(), bson.D{{Key: "account", Value: account}})

			if r.Err() != nil {
				if r.Err() == mongo.ErrNoDocuments {
					result["err"] = "no this user"
				} else {
					result["err"] = r.Err().Error()
				}
				return
			}

			user := map[string]interface{}{}
			err := r.Decode(&user)
			if err != nil {
				result["err"] = err.Error()
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(password))
			if err != nil {
				result["err"] = err.Error()
				return
			}

			session := sessions.Default(c)
			defer session.Save()
			session.Set("account", account)

			for _, h := range set_session_hooks {
				session.Set(h(user))
			}
		}
}

//添加用户
func signupGet() (name string, paramsStructPtr interface{}, method web.Method, path string, handlers gin.HandlerFunc) {

	return "注册页面", nil, web.GET, "/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user/signup", gin.H{})

	}

}
func signupPost() (name string, paramsStructPtr interface{}, method web.Method, path string, handlers gin.HandlerFunc) {

	type input struct {
		Account  string `json:"account" form:"account" name:"账户"`
		Password string `json:"password" form:"password" name:"密码"`
	}

	return "注册执行", new(input), web.POST, "/signup", func(c *gin.Context) {
		account := c.PostForm("account")
		password := c.PostForm("password")

		result := gin.H{
			"err": "",
		}

		defer c.HTML(http.StatusOK, "user/login", result)

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			result["err"] = err.Error()
			return
		}

		id, err := metadata.InsertOneFromMetadata(me.Name, bson.M{"account": account, "password": string(hashPassword)})

		if err != nil {
			result["err"] = err.Error()
			return
		}

		result["id"] = id
	}

}

func signout() (string, interface{}, web.Method, string, gin.HandlerFunc) {
	return "注销", nil, web.GET, "/signout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
	}
}

func setSession() func(c *gin.Context) {

	store := memstore.NewStore([]byte(sessionKey))
	store.Options(sessions.Options{
		MaxAge: 0, // seems this works
		Path:   "/",
	})

	return sessions.Sessions("auth", store)
}

//https://github.com/open-policy-agent/opa
