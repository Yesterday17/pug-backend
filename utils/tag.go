package utils

import "reflect"

func GetFieldByTag(i interface{}, tag, name string) interface{} {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if name == f.Tag.Get(tag) {
			return v.Field(i).Interface()
		}
	}
	return nil
}
