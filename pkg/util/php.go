package util

import (
	serialize "github.com/techoner/gophp"
	"reflect"
	"time"
)

func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}

func Unserialize(data string) (interface{}, error) {
	return serialize.Unserialize([]byte(data))
}

func Time() int64 {
	return time.Now().Unix()
}

func Strtotime(format, strtime string) (int64, error) {
	t, err := time.ParseInLocation(format, strtime, time.Local)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
