package test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"testing"
)

func TestAst(t *testing.T) {
	fset := token.NewFileSet()
	// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
	path, _ := filepath.Abs("/Users/zeta/workspace/golego/test/astcode.go")
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Println(err)
		return
	}
	// 打印语法树
	ast.Print(fset, f)
}
