package database

import (
	"context"
	"fmt"
	"golego/modules/bootstrap"
	"golego/modules/config"
	"golego/modules/helper"
	"golego/modules/webserver"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//实现一个开放的GetInfo方法

var me = helper.ModuleInfo{
	Name:      "database",
	HumanName: "数据库模块",
	Templates: []string{"status"},
}

var db *mongo.Database

var collections sync.Map

type connectd_hook func()

var connectd_hooks []connectd_hook

func init() {
	helper.Register(me)
	bootstrap.AtBeforeRun(connect)
	bootstrap.AtAfterRun(disconnect)

	webserver.AtSetHandle(func() (webserver.RequestMethod, string, gin.HandlerFunc) {
		return webserver.GET, "/collections", getCollections
	})
}

func AtConnected(h connectd_hook) {
	connectd_hooks = append(connectd_hooks, h)
}

func connect() {
	cfg, ok := config.Get(me.Name)
	if !ok {
		cfg = gjson.Parse(`{"conn":"mongodb://localhost:27017","db":"golego"}`)
		config.Set(me.Name, cfg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Get("conn").String()))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	db = client.Database(cfg.Get("db").String())

	for _, hook := range connectd_hooks {
		hook()
	}
}

//GetDataBase 获得数据库
func getDB() *mongo.Database {
	return db
}

//注册集合，避免集合冲突
func RegisterC(name string) {

	if _, ok := collections.Load(name); ok {
		err := fmt.Errorf("collection [%s] is exist", name)
		panic(err)
	}

	collections.Store(name, 1)

}

func Collections() (cs []string) {

	collections.Range(func(key interface{}, value interface{}) bool {
		cs = append(cs, key.(string))
		return true
	})

	return
}

//GetCollection 获得集合对象
func C(name string) *mongo.Collection {
	//未注册的集合返回空，不允许操作。
	if _, ok := collections.Load(name); !ok {
		return nil
	}

	return getDB().Collection(name)
}

func disconnect() {
	if db != nil {
		db.Client().Disconnect(context.Background())
	}
}

func getCollections(c *gin.Context) {
	c.String(200, strings.Join(Collections(), ","))
}
