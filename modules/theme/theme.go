package theme

import (
	"fmt"
	"golego/modules/bootstrap"
	"golego/modules/config"
	"golego/modules/helper"
	"golego/modules/webserver"
	"os"

	"golego/utils/path"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

var themeRoot string = "./theme"
var theme string = ""

func init() {
	webserver.AddSetRouterHook(setTemplate)
	bootstrap.AddBeforeRunHook(initTemplate)
}

//实现一个开放的GetInfo方法
func GetInfo() helper.Info {
	return helper.Info{
		Name: "theme",
	}
}

func setTemplate(r *gin.Engine) {

	//设置静态文件
	r.Static("/static", fmt.Sprintf("%s/%s/static", themeRoot, theme))

	//设置模板文件地址
	r.LoadHTMLGlob(fmt.Sprintf("%s/%s/templates/**/*", themeRoot, theme))

}

func initTemplate() {
	//获得本模块的配置
	//如果不存在，则写入一个默认配置
	cfg, ok := config.Get(GetInfo().Name)
	if !ok {
		cfg = gjson.Parse(`{"theme":"default"}`)
		config.Set(GetInfo().Name, cfg)
	}

	theme = cfg.Get("theme").String()

	//根据主题配置生成静态目录
	err := os.MkdirAll(fmt.Sprintf("%s/%s/static", themeRoot, theme), 0755)
	if err != nil {
		panic(err)
	}

	//根据模块配置，检查是否存在对应的模版文件，如果不存在，自动生成
	for _, info := range helper.ModuleInfos {
		fmt.Println(info.Name, info.Templates)

		err := os.MkdirAll(fmt.Sprintf("%s/%s/templates/%s", themeRoot, theme, info.Name), 0755)
		if err != nil {
			panic(err)
		}

		for _, template := range info.Templates {
			fmt.Println(template)

			e, _ := path.PathInfo(fmt.Sprintf("%s/%s/templates/%s/%s.tpl", themeRoot, theme, info.Name, template))
			if !e {
				f, err := os.Create(fmt.Sprintf("%s/%s/templates/%s/%s.tpl", themeRoot, theme, info.Name, template))
				if err != nil {
					panic(err)
				}

				_, err = f.Write([]byte(fmt.Sprintf("{{define \"%s/%s\"}}\n\n{{end}}", info.Name, template)))
				if err != nil {
					panic(err)
				}
				f.Close()
			}
		}
	}
}
