package utils

import (
	"fmt"
	"reflect"
	"testing"
)

type testStructA struct {
	A string            `json:"b" another:"dd"`
	B int               `json:"a"`
	C map[string]string `json:"map"`
	D map[string]string `json:"anotherMap"`
}

func TestGetFieldByTag(t *testing.T) {
	a := testStructA{
		A: "test",
		B: 114514,
		C: map[string]string{"a": "1"},
		D: nil,
	}

	if ret := GetFieldByTag(a, "json", "b"); ret != a.A {
		t.Error(fmt.Errorf("expects %v, got %v", a.A, ret))
	}

	if ret := GetFieldByTag(a, "json", "a"); ret != a.B {
		t.Error(fmt.Errorf("expects %v, got %v", a.B, ret))
	}

	// if ret := GetFieldByTag(a, "json", "map"); ret != a.C {
	// 	t.Error(fmt.Errorf("expects %v, got %v", a.C, ret))
	// }

	if ret := GetFieldByTag(a, "json", "anotherMap"); !reflect.ValueOf(ret).IsNil() {
		t.Error(fmt.Errorf("expects %v, got %v", a.D, ret))
	}

	if ret := GetFieldByTag(a, "another", "dd"); ret != a.A {
		t.Error(fmt.Errorf("expects %v, got %v", a.A, ret))
	}
}

func TestSetFieldByTag(t *testing.T) {
	a := testStructA{
		A: "test",
		B: 114514,
		C: map[string]string{"a": "1"},
		D: nil,
	}

	if err := SetFieldByTag(&a, "json", "a", 1919810); err != nil {
		t.Error(err)
	} else if a.B != 1919810 {
		t.Error("value not set")
	}

	if err := SetFieldByTag(&a, "json", "map", map[string]string{"a": "2"}); err != nil {
		t.Error(err)
	} else if val, ok := a.C["a"]; !ok {
		t.Error("key not found")
	} else if val != "2" {
		t.Error("value not set")
	}
}
