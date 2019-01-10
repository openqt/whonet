package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"sort"
)

var (
	log *logrus.Logger = nil
	LOG                = GetLogger()
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
	LOG.Debugf("Kind %v %v", v1.Kind(), v2.Kind())
	LOG.Debugf("Type %v %v", v1.Type(), v2.Type())
	LOG.Debugf("Interface %v %v", v1.Interface(), v2.Interface())

	switch v1.Kind() {
	case reflect.Slice, reflect.Array:
		switch v2.Kind() {
		case reflect.Slice, reflect.Array:
			break
		default:
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}

		for i := 0; i < v1.Len(); i++ {
			if !DeepEqual(v1.Index(i).Interface(), v2.Index(i).Interface()) {
				return false
			}
		}
		return true
	case reflect.Map:
		if v2.Kind() != reflect.Map {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}

		for _, key := range v1.MapKeys() {
			if !v2.MapIndex(key).IsValid() {
				return false
			}

			DeepEqual(v1.MapIndex(key).Interface(), v2.MapIndex(key).Interface())
		}
		return true
	default:
		if v1.Type() != v2.Type() {
			return false
		}
		return x == y
	}
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func GetLogger() *logrus.Logger {
	if log == nil {
		log = logrus.New()
		log.Level = logrus.InfoLevel
		log.Out = os.Stdout

		log.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})

		log.Info("Logger initialized.")
	}
	return log
}
