package modules

//导入所有模块包，以执行所有包的Init函数
import (
	"golego/modules/bootstrap"
	"golego/modules/config"
	"golego/modules/database"
	"golego/modules/ginserver"
	"golego/modules/helloworld"
	"golego/modules/helper"
)

//取得所有包的Info信息
func GetModuleInfos() map[string]helper.Info {
	return map[string]helper.Info{
		bootstrap.GetInfo().Name:  bootstrap.GetInfo(),
		ginserver.GetInfo().Name:  ginserver.GetInfo(),
		helper.GetInfo().Name:     helper.GetInfo(),
		helloworld.GetInfo().Name: helloworld.GetInfo(),
		database.GetInfo().Name:   database.GetInfo(),
	}
}

//系统启动模块的入口
func Start() {

	defer database.Disconnect()

	helper.ModuleInfos = GetModuleInfos()
	bootstrap.Run()
	config.SaveConfig()
	helper.WaitGroup.Wait()

}
