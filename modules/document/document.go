package document

import (
	"golego/modules/helper"
)

var me = helper.ModuleInfo{
	Name:      "document",
	HumanName: "文档模块",
}

func init() {
	helper.Register(me)
}
