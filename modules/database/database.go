package database

import (
	"context"
	"golego/modules/bootstrap"
	"golego/modules/config"
	"golego/modules/helper"
	"time"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	bootstrap.AddBeforeRunHook(initDataBase)
	bootstrap.AddAfterRunHook(disconnect)
}

func initDataBase() {
	cfg, ok := config.Get(GetInfo().Name)
	if !ok {
		cfg = gjson.Parse(`{"conns":{"default":"connectstring1","connName2":"connecstring2"}}`)
		config.Add(GetInfo().Name, cfg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for name, item := range cfg.Get("conns").Map() {
		dbClients[name], dbClientErrs[name] = mongo.Connect(ctx, options.Client().ApplyURI(item.String()))
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
