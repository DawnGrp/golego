package user

import (
	"golego/modules/bootstrap"
	"golego/modules/helper"
	"golego/modules/metadata"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetInfo() helper.Info {
	return helper.Info{
		Name:      "user",
		HumanName: "用户模块",
	}
}

func init() {
	bootstrap.AddRunHook(createUserMetadata)
}

//创建user集合以及索引
func createUserMetadata() {

	md := metadata.Metadata{
		Name:      GetInfo().Name,
		HumanName: GetInfo().HumanName,
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
