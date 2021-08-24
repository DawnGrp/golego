package auth

import (
	"golego/modules/helper"
)

var me = helper.ModuleInfo{
	Name:      "auth",
	HumanName: "权限模块",
}

func init() {
	helper.Register(me)
}

func Register() (err error) {

	return err
}
