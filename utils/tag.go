package utils

import (
	"errors"
	"reflect"
)

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

func SetFieldByTag(i interface{}, tag, name string, value interface{}) error {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if name == f.Tag.Get(tag) {
			rf := v.Field(i)
			if !rf.CanSet() || rf.Kind() != reflect.TypeOf(value).Kind() {
				return errors.New("cannot set field")
			}

			rf.Set(reflect.ValueOf(value))
			return nil
		}
	}
	return errors.New("field not found")
}
