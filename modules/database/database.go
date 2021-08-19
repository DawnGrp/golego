package database

import (
	"context"
	"golego/modules/bootstrap"
	"golego/modules/config"
	"golego/modules/ginserver"
	"golego/modules/helper"
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
		Name: "database",
	}
}

var dbClients = map[string]*mongo.Client{}
var dbClientErrs = map[string]error{}

func init() {
	bootstrap.AddBeforeRunHook(connect)
	bootstrap.AddAfterRunHook(disconnect)
	ginserver.AddSetHandleHook(status)
}

func connect() {
	cfg, ok := config.Get(GetInfo().Name)
	if !ok {
		cfg = gjson.Parse(`{"conns":{"default":"mongodb://localhost:27017/db"}}`)
		config.Add(GetInfo().Name, cfg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for name, item := range cfg.Get("conns").Map() {
		dbClients[name], dbClientErrs[name] = mongo.Connect(ctx, options.Client().ApplyURI(item.String()))

		if dbClientErrs[name] == nil {
			dbClientErrs[name] = dbClients[name].Ping(ctx, readpref.Primary())
		}
	}
}

func GetClient(name string) (*mongo.Client, error) {
	return dbClients[name], dbClientErrs[name]
}

func disconnect() {
	for _, dbclient := range dbClients {
		if dbclient != nil {
			dbclient.Disconnect(context.Background())
		}
	}
}

func status() (ginserver.RequestMethod, string, gin.HandlerFunc) {

	return ginserver.GET, "/dbclienterrs", func(c *gin.Context) {

		for name, err := range dbClientErrs {
			if err != nil {
				c.String(200, name+":"+err.Error()+"\n")
			} else {
				c.String(200, name+":success\n")
			}

		}

	}

}
