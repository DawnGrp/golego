package utils

import (
	"fmt"
	"reflect"
)

func GetStruct(structPtr interface{}) (name string, r [][]string) {

	if structPtr == nil {
		return
	}

	v := reflect.ValueOf(structPtr)

	if v.Kind() != reflect.Ptr {
		panic(fmt.Errorf("kind must be ptr, current is %s", v.Kind()))
	}

	v = v.Elem()

	name = v.Type().Name()

	for i := 0; i < v.NumField(); i++ {

		r = append(r, []string{
			v.Type().Field(i).Name,
			v.Type().Field(i).Type.String(),
			v.Type().Field(i).Tag.Get("name"),
			v.Type().Field(i).Tag.Get("json"),
			v.Type().Field(i).Tag.Get("form"),
		})

	}

	return
}
