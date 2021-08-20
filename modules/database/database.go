package database

import (
	"context"
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

var dbClient *mongo.Client
var dbClientErr error

func init() {
	bootstrap.AddBeforeRunHook(connect)
	bootstrap.AddAfterRunHook(disconnect)
	webserver.AddSetHandleHook(status)
}

func connect() {
	cfg, ok := config.Get(GetInfo().Name)
	if !ok {
		cfg = gjson.Parse(`{"conn":"mongodb://localhost:27017/db"}`)
		config.Set(GetInfo().Name, cfg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbClient, dbClientErr = mongo.Connect(ctx, options.Client().ApplyURI(cfg.Get("conn").String()))

	if dbClientErr == nil {
		dbClientErr = dbClient.Ping(ctx, readpref.Primary())
	}

}

func GetClient() (*mongo.Client, error) {
	return dbClient, dbClientErr
}

func disconnect() {
	if dbClient != nil {
		dbClient.Disconnect(context.Background())
	}
}

func status() (webserver.RequestMethod, string, gin.HandlerFunc) {

	return webserver.GET, "/status", func(c *gin.Context) {

		if dbClientErr != nil {
			c.String(200, dbClientErr.Error())
		} else {
			c.String(200, "success")
		}

	}

}
