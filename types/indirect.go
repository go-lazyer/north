package types

import (
	"reflect"
)

// 首先调用indirect函数将参数中可能的指针去掉。如果类型本身不是指针，那么直接返回。否则返回指针指向的值。
// 循环直到返回一个非指针的值：
func indirect(i any) (any, bool) {
	if i == nil {
		return nil, false
	}
	if t := reflect.TypeOf(i); t.Kind() != reflect.Ptr {
		return i, false
	}
	v := reflect.ValueOf(i)

	for v.Kind() == reflect.Ptr || (v.Kind() == reflect.Interface && v.Elem().Kind() == reflect.Ptr) {
		if v.IsNil() {
			return nil, true
		}

		v = v.Elem()
	}
	return v.Interface(), true
}
