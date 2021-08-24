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
func GetInfo() helper.Info {
	return helper.Info{
		Name:      "metadata",
		HumanName: "元数据模块",
	}
}

func init() {
	bootstrap.AddRunHook(createMetadataIndex)
}

type FiledType string
type IndexType string

const (
	FiledType_String    FiledType = "string"
	FiledType_Int       FiledType = "int"
	FiledType_Float     FiledType = "float"
	FiledType_Interface FiledType = "interface"
	FiledType_Array     FiledType = "array"
	FiledType_Map       FiledType = "map"

	IndexType_Unique IndexType = "unique"
	IndexType_Index  IndexType = "index"
)

var me = db.RegisterCollection("metadata")

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

//todo: 在什么时候添加index呢 ？，应该在Replace和Insert的时候添加

func Replace(md Metadata, no_document_to_insert bool) (id interface{}, err error) {

	err = createIndexs(db.Collection(md.Name), md.Indexs)
	if err != nil {
		return
	}

	opts := options.Replace().SetUpsert(no_document_to_insert)
	c := db.C(me)
	ir, err := c.ReplaceOne(
		context.Background(),
		bson.D{{Key: "name", Value: md.Name}}, md, opts)

	if err != nil {
		return
	}

	id = ir.UpsertedID
	return
}

func Insert(md Metadata) (id interface{}, err error) {
	err = createIndexs(db.Collection(md.Name), md.Indexs)
	if err != nil {
		return
	}

	c := db.C(me)
	ir, err := c.InsertOne(
		context.Background(), md)

	if err != nil {
		return
	}
	id = ir.InsertedID
	return
}

func Get(name string) (metadata Metadata, err error) {
	c := db.C(me)
	sr := c.FindOne(context.Background(), bson.D{{Key: "name", Value: name}})
	err = sr.Err()
	if err != nil {
		return
	}
	err = sr.Decode(&metadata)

	return
}

func GetAll() (mds []Metadata, err error) {
	c := db.C(me)
	cursor, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		return
	}
	err = cursor.All(context.Background(), &mds)
	return
}

func Del(name string) (err error) {
	c := db.C(me)
	_, err = c.DeleteOne(context.Background(), bson.D{{Key: "name", Value: name}})
	return err
}

//通过元数据创建数据
func InsertOneFromMetadata(metadataName string, fields map[string]interface{}) (newid interface{}, err error) {
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

	c := db.C(db.Collection(metadataName))
	ir, err := c.InsertOne(context.Background(), data)
	newid = ir.InsertedID
	return
}

func UpdateByIDFromMetadata(metadataName string, id interface{}, fields map[string]interface{}) (err error) {
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

	c := db.C(db.Collection(metadataName))
	_, err = c.UpdateByID(context.Background(), id, data)

	return
}

func DeleteByIDFromMetadata(metadataName string, id interface{}) (err error) {
	_, err = Get(metadataName)
	if err != nil {
		return
	}

	c := db.C(db.Collection(metadataName))
	_, err = c.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}})

	return
}

func createIndexs(collectionName db.Collection, indexs []Index) (err error) {

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

	err := createIndexs(me, indexs)
	if err != nil {
		panic(err)
	}
}
