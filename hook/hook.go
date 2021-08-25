package hook

import (
	"fmt"
	"reflect"
)

//定义到包
type HookType interface{}
type Hook struct {
	Group      []HookType
	defineType HookType
}

func (h *Hook) Up(f HookType) {

	if reflect.TypeOf(f).String() != reflect.TypeOf(h.defineType).String() {
		panic(fmt.Errorf("%s != %s", reflect.TypeOf(f).String(), reflect.TypeOf(h.defineType).String()))
	}

	h.Group = append(h.Group, f)

}

func (h *Hook) In(call func(HookType)) {
	for _, f := range h.Group {
		call(f)
	}
}

func CreateHook(defineType HookType) Hook {
	return Hook{defineType: defineType}
}
