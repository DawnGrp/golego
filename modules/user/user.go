package user

import (
	"golego/modules/bootstrap"
	"golego/modules/helper"
	"golego/modules/metadata"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usermdname      = "user"
	usermdhumanname = "用户"
)

func GetInfo() helper.Info {
	return helper.Info{
		Name: "user",
	}
}

func init() {
	bootstrap.AddRunHook(createUserMetadata)
}

func createUserMetadata() {

	md := metadata.Metadata{
		Name:      usermdname,
		HumanName: usermdhumanname,
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
		Indexs: []metadata.Index{},
	}

	_, err := metadata.Insert(md)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		panic(err)
	}

}
