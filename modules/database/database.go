package database

import (
	"context"
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

var dbClients map[string]*mongo.Client

func init() {
	initDataBase()
}

func initDataBase() {
	cfg, ok := config.Get(GetInfo().Name)
	if !ok {
		cfg = gjson.Parse(`{"conns":{"default":"connectstring1","connName2":"connecstring2"}}`)
		config.Add(GetInfo().Name, cfg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	for name, item := range cfg.Get("conns").Map() {
		dbClients[name], err = mongo.Connect(ctx, options.Client().ApplyURI(item.String()))
		if err != nil {
			panic(err)
		}
	}
}

func GetClient(name string) *mongo.Client {
	return dbClients[name]
}

func Disconnect() {
	for _, dbclient := range dbClients {
		dbclient.Disconnect(context.Background())
	}
}
