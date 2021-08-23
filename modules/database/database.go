package database

import (
	"context"
	"fmt"
	"golego/modules/bootstrap"
	"golego/modules/config"
	"golego/modules/helper"
	"golego/modules/webserver"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//实现一个开放的GetInfo方法
func GetInfo() helper.Info {
	return helper.Info{
		Name:      "database",
		Templates: []string{"status"},
	}
}

var db *mongo.Database
var err error

type Collection string

var collections = map[string]int{}

func init() {
	bootstrap.AddBeforeRunHook(connect)
	bootstrap.AddAfterRunHook(disconnect)
	webserver.AddSetHandleHook(func() (webserver.RequestMethod, string, gin.HandlerFunc) {
		return webserver.GET, "/status", status
	})
}

func RegisterCollection(C string) Collection {
	if _, ok := collections[C]; ok {
		panic(fmt.Errorf("%s exist", C))
	}

	collections[C] = 1

	return Collection(C)
}

func connect() {
	cfg, ok := config.Get(GetInfo().Name)
	if !ok {
		cfg = gjson.Parse(`{"conn":"mongodb://localhost:27017","db":"golego"}`)
		config.Set(GetInfo().Name, cfg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client *mongo.Client
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.Get("conn").String()))

	if err == nil {
		err = client.Ping(ctx, readpref.Primary())
	}

	db = client.Database(cfg.Get("db").String())

}

//GetDataBase 获得数据库
func DB() *mongo.Database {
	return db
}

//GetCollection 获得集合对象
func C(name Collection) *mongo.Collection {
	return DB().Collection(string(name))
}

func disconnect() {
	if db != nil {
		db.Client().Disconnect(context.Background())
	}
}

func status(c *gin.Context) {
	if err != nil {
		c.String(200, err.Error())
	} else {
		c.String(200, "success")
	}

}
