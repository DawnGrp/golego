package auth

import (
	db "golego/modules/database"
	"golego/modules/helper"
)

var me = helper.ModuleInfo{
	Name:      "auth",
	HumanName: "权限模块",
}

func init() {
	helper.Register(me)
	db.RegisterC(me.Name)
}

func Register() (err error) {

	return err
}
