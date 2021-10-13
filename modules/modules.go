package modules

//导入所有模块包，以执行所有包的Init函数
import (
	"golego/modules/bootstrap"

	_ "golego/modules/config"
	_ "golego/modules/database"
	_ "golego/modules/document"
	"golego/modules/helper"

	_ "golego/modules/auth"
	_ "golego/modules/metadata"
	_ "golego/modules/theme"
	_ "golego/modules/webserver"
)

//系统启动模块的入口
func Start() {

	bootstrap.BeforeRun()
	bootstrap.Run()
	helper.WaitGroup.Wait()
	bootstrap.AfterRun()

}
