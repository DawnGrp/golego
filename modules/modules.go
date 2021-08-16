package modules

//导入所有模块包，以执行所有包的Init函数
import (
	"golego/modules/bootstrap"
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
	}
}

//系统启动模块的入口
func Start() {

	helper.ModuleInfos = GetModuleInfos()
	bootstrap.Run()
	helper.WaitGroup.Wait()
}
