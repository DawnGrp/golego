package modules

//导入所有模块包，以执行所有包的Init函数
import (
	"golego/modules/bootstrap"
	"golego/modules/config"
	"golego/modules/database"
	"golego/modules/helper"
	"golego/modules/theme"
	"golego/modules/webserver"
)

//取得所有包的Info信息
func GetModuleInfos() map[string]helper.Info {
	return map[string]helper.Info{
		bootstrap.GetInfo().Name: bootstrap.GetInfo(),
		webserver.GetInfo().Name: webserver.GetInfo(),
		helper.GetInfo().Name:    helper.GetInfo(),
		theme.GetInfo().Name:     theme.GetInfo(),
		database.GetInfo().Name:  database.GetInfo(),
		config.GetInfo().Name:    config.GetInfo(),
	}
}

//系统启动模块的入口
func Start() {

	helper.ModuleInfos = GetModuleInfos()
	bootstrap.BeforeRun()
	bootstrap.Run()
	bootstrap.AfterRun()
	helper.WaitGroup.Wait()

}
