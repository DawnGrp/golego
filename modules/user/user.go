package user

import (
	"context"
	db "golego/modules/database"
	"golego/modules/helper"
	"golego/modules/metadata"
	"golego/modules/webserver"
	"net/http"
	"strings"

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
		"login",
	},
}

func init() {
	helper.Register(me)
	db.RegisterC(me.Name)
	db.AtConnected(createUserMetadata)
	webserver.AtSetHandle(func() (webserver.RequestMethod, string, gin.HandlerFunc) {
		return webserver.GET, "/login", login
	})
	webserver.AtSetHandle(func() (webserver.RequestMethod, string, gin.HandlerFunc) {
		return webserver.POST, "/login", loginPost
	})
	webserver.AtSetHandle(func() (webserver.RequestMethod, string, gin.HandlerFunc) {
		return webserver.POST, "/adduser", addUserPost
	})
}

var sessionKey = strings.Replace(uuid.NewV4().String(), "-", "", -1)

//创建user集合以及索引
func createUserMetadata() {

	md := metadata.Metadata{
		Name:      me.Name,
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

	_, err := metadata.Insert(md)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		panic(err)
	}

}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "user/login", gin.H{})
}

func loginPost(c *gin.Context) {

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

	user := map[string]string{}
	err := r.Decode(&user)
	if err != nil {
		result["err"] = err.Error()
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user["password"]), []byte(password))
	if err != nil {
		result["err"] = err.Error()
		return
	}

}

func addUserPost(c *gin.Context) {
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
