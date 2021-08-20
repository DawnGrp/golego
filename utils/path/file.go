package path

import "os"

//判断路径是否存在
func PathInfo(path string) (exist bool, isdir bool) {
	s, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err), false
	}
	return true, s.IsDir()
}
