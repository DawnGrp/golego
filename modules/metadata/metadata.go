package metadata

import (
	"context"
	"golego/modules/database"
	"golego/modules/helper"

	"golego/modules/bootstrap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//实现一个开放的GetInfo方法
func GetInfo() helper.Info {
	return helper.Info{
		Name: "metadata",
	}
}

func init() {
	bootstrap.AddRunHook(createMetadataIndex)
}

type FiledType string
type IndexType string

const (
	String    FiledType = "string"
	Int       FiledType = "int"
	Float     FiledType = "float"
	Interface FiledType = "interface"
	Array     FiledType = "array"
	Map       FiledType = "map"

	Unique IndexType = "unique"
)

// var me, me2 database.Collection
var me = database.RegisterCollection("metadata")

type Metadata struct {
	Name      string
	HumanName string
	Fileds    []Filed
	Index     []Index
}

type Filed struct {
	Name        string
	HumanName   string
	Type        FiledType
	Options     []interface{}
	MultiSelect bool
}

type Index struct {
	Name  string
	Type  IndexType
	Filed []string
}

func Create(metadata Metadata) (id interface{}, err error) {
	c := database.C(me)
	ir, err := c.InsertOne(context.Background(), metadata)
	if err != nil {
		return
	}
	id = ir.InsertedID
	return
}

func GetOne(name string) (metadata Metadata, err error) {
	c := database.C(me)
	sr := c.FindOne(context.Background(), bson.D{{Key: "name", Value: name}})
	if sr.Err() != nil {
		return
	}
	err = sr.Decode(&metadata)
	return
}

func GetList() (metadatas []Metadata, err error) {
	c := database.C(me)
	cursor, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		return
	}
	err = cursor.All(context.Background(), &metadatas)
	return
}

func createMetadataIndex() {
	isUnique := true
	indexName := "unique_name"
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: &options.IndexOptions{Unique: &isUnique, Name: &indexName},
	}

	indexName, err := database.C(me).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		panic(err)
	}
}
