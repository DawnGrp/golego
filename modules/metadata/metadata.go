package metadata

import (
	"context"
	db "golego/modules/database"
	"golego/modules/helper"
	"strings"

	"golego/modules/bootstrap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//实现一个开放的GetInfo方法
var me = helper.ModuleInfo{
	Name:      "metadata",
	HumanName: "元数据模块",
}

func init() {
	helper.Register(me)
	db.RegisterC(me.Name)
	//系统启动时创建元数据模块的集合和索引
	bootstrap.AtRun(createMetadataIndex)
}

type FiledType string //定义字段类型
type IndexType string //定义索引类型

const (
	//字段类型的项
	FiledType_String    FiledType = "string"
	FiledType_Int       FiledType = "int"
	FiledType_Float     FiledType = "float"
	FiledType_Interface FiledType = "interface"
	FiledType_Array     FiledType = "array"
	FiledType_Map       FiledType = "map"

	//索引类型的项
	IndexType_Unique IndexType = "unique"
	IndexType_Index  IndexType = "index"
)

type Metadata struct {
	Name      string  `bson:"name"`
	HumanName string  `bson:"human_name"`
	Fileds    []Filed `bson:"fileds"`
	Indexs    []Index `bson:"indexs"`
}

type Filed struct {
	Name        string        `bson:"name"`
	HumanName   string        `bson:"human_name"`
	Type        FiledType     `bson:"type"`
	Options     []interface{} `bson:"options"`
	MultiSelect bool          `bson:"multi_select"`
}

type Index struct {
	Type   IndexType `bson:"type"`
	Fileds []string  `bson:"fileds"`
}

//更新替换一个元数据
func Replace(md Metadata, no_document_to_insert bool) (id interface{}, err error) {

	err = createIndexs(md.Name, md.Indexs)
	if err != nil {
		return
	}

	opts := options.Replace().SetUpsert(no_document_to_insert)
	c := db.C(me.Name)
	ir, err := c.ReplaceOne(
		context.Background(),
		bson.D{{Key: "name", Value: md.Name}}, md, opts)

	if err != nil {
		return
	}

	id = ir.UpsertedID
	return
}

//添加元数据
func Insert(md Metadata) (id interface{}, err error) {
	err = createIndexs(md.Name, md.Indexs)
	if err != nil {
		return
	}

	c := db.C(me.Name)
	ir, err := c.InsertOne(
		context.Background(), md)

	if err != nil {
		return
	}
	id = ir.InsertedID
	return
}

//获得元数据
func Get(name string) (metadata Metadata, err error) {
	c := db.C(me.Name)
	sr := c.FindOne(context.Background(), bson.D{{Key: "name", Value: name}})
	err = sr.Err()
	if err != nil {
		return
	}
	err = sr.Decode(&metadata)

	return
}

//获得全部元数据
func GetAll() (mds []Metadata, err error) {
	c := db.C(me.Name)
	cursor, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		return
	}
	err = cursor.All(context.Background(), &mds)
	return
}

//删除元数据
func Del(name string) (err error) {
	c := db.C(me.Name)
	_, err = c.DeleteOne(context.Background(), bson.D{{Key: "name", Value: name}})
	return err
}

//通过元数据创建数据
func InsertOneFromMetadata(metadataName string, fields bson.M) (newid interface{}, err error) {
	md, err := Get(metadataName)
	if err != nil {
		return
	}

	data := bson.D{}
	for _, f := range md.Fileds {
		if field, ok := fields[f.Name]; ok {
			data = append(data, bson.E{Key: f.Name, Value: field})
		}
	}
	//todo:缺少类型检查

	c := db.C(metadataName)
	ir, err := c.InsertOne(context.Background(), data)
	if ir != nil {
		newid = ir.InsertedID
	}

	return
}

func UpdateByIDFromMetadata(metadataName string, id interface{}, fields bson.M) (err error) {
	md, err := Get(metadataName)
	if err != nil {
		return
	}

	data := bson.D{}
	for _, f := range md.Fileds {
		if field, ok := fields[f.Name]; ok {
			data = append(data, bson.E{Key: f.Name, Value: field})
		}
	}
	//todo:缺少类型检查

	c := db.C(metadataName)
	_, err = c.UpdateByID(context.Background(), id, data)

	return
}

func DeleteByIDFromMetadata(metadataName string, id interface{}) (err error) {
	_, err = Get(metadataName)
	if err != nil {
		return
	}

	c := db.C(metadataName)
	_, err = c.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}})

	return
}

func createIndexs(collectionName string, indexs []Index) (err error) {

	for _, index := range indexs {

		var isUnique = index.Type == IndexType_Unique

		var indexName = "index_"
		if isUnique {
			indexName = "unique_"
		}
		indexName += strings.Join(index.Fileds, "_")

		keys := bson.D{}
		for _, field := range index.Fileds {
			keys = append(keys, bson.E{Key: field, Value: 1})
		}

		indexModel := mongo.IndexModel{
			Keys:    keys,
			Options: &options.IndexOptions{Unique: &isUnique, Name: &indexName},
		}

		indexName, err = db.C(collectionName).Indexes().CreateOne(context.Background(), indexModel)

		if err != nil {
			break
		}
	}

	return
}

func createMetadataIndex() {

	indexs := []Index{
		{
			Type:   IndexType_Unique,
			Fileds: []string{"name"},
		},
	}

	err := createIndexs(me.Name, indexs)
	if err != nil {
		panic(err)
	}
}
