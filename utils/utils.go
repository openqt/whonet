package utils

import (
	"reflect"
	"sort"
	)

// 对字典的key排序
func SortKeys(l []reflect.Value) []string {
	var keys []string
	for _, v := range l {
		keys = append(keys, v.String())
	}
	sort.Strings(keys)
	return keys
}

// 比较两个数组的值是否一致，忽略类型信息，数组值的顺序也必须一致
func DeepEqual(x, y interface{}) bool {
	if x == nil || y == nil {
		return x == y
	}

	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	//if v1.Type() != v2.Type() {
	//	return false
	//}
	if v1.Len() != v2.Len() {
		return false
	}
	//for i:=0;i<v1.Len();i++{
	//
	//	switch e1 := v1.Index(i); e1.Kind() {
	//	case reflect.Interface:
	//	}
	//	if v1.Index(i).Elem() != v2.Index(i).Elem() {
	//		return false
	//	}
	//}
	return true
	//return deepValueEqual(v1, v2, make(map[visit]bool), 0)}
}
