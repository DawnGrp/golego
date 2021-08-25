package helper

import (
	"fmt"
	"testing"
)

func CreateHook(f interface{}) (at func(interface{}), group []interface{}, exec func()) {

	return
}

type myfunc func()

func hello() {
	fmt.Println("Hello")
}

func TestCreateHook(t *testing.T) {

}
