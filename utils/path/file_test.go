package path

import (
	"os"
	"testing"
)

func TestCreateDir(t *testing.T) {
	err := os.MkdirAll("/Users/zeta/workspace/golego/utils/path/a/b/c", 0755)
	t.Log(err)
}
