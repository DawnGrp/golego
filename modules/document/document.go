package document

import (
	"golego/modules/helper"
)

func GetInfo() helper.Info {
	return helper.Info{
		Name:      "document",
		HumanName: "文档模块",
	}
}
