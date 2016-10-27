// httputils
package utils

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

func StructToURLQuery(i interface{}, hasSubKey bool) (string, error) {
	iValue := reflect.ValueOf(i)
	if iValue.Kind() == reflect.Ptr {
		iValue = iValue.Elem()
	}
	if iValue.Kind() != reflect.Struct {
		return "", fmt.Errorf("utils:准备转换的对象不是Struct类型,实际类型为:%s", iValue.Kind())
	}
	values := url.Values{}
	setValues(iValue, values, "", hasSubKey)
	return values.Encode(), nil
}

func setValues(iValue reflect.Value, values url.Values, key string, hasSubKey bool) {
	if iValue.Kind() == reflect.Ptr {
		iValue = iValue.Elem()
	}
	switch iValue.Kind() {
	case reflect.String:
		values.Add(key, iValue.String())
		break
	case reflect.Bool:
		if iValue.Bool() {
			values.Add(key, "true")
		} else {
			values.Add(key, "false")
		}
		break
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		values.Add(key, strconv.FormatInt(iValue.Int(), 10))
		break
	case reflect.Float32, reflect.Float64:
		values.Add(key, strconv.FormatFloat(iValue.Float(), 'f', 2, 64))
		break
	case reflect.Struct:
		iType := iValue.Type()
		fieldSize := iValue.NumField()
		if key != "" {
			key = key + "."
		}
		for i := 0; i < fieldSize; i++ {
			fieldValue := iValue.Field(i)
			if !fieldValue.IsValid() {
				continue
			}
			fieldType := iType.Field(i)
			subkey := fieldType.Name
			tag := fieldType.Tag.Get("json")
			if tag != "" {
				subkey = tag
			}
			if hasSubKey {
				subkey = key + subkey
			}
			setValues(fieldValue, values, subkey, hasSubKey)
		}
		break
	case reflect.Array, reflect.Slice:
		arraySize := iValue.Len()
		for index := 0; index < arraySize; index++ {
			subValue := iValue.Index(index)
			setValues(subValue, values, key, hasSubKey)
		}
		break
	case reflect.Map:
		mapKeys := iValue.MapKeys()
		if key != "" {
			key = key + "."
		}
		for _, mapkey := range mapKeys {
			if mapkey.Type().Kind() != reflect.String {
				//日过map的KEY不是string类型则忽略
				continue
			}
			fieldValue := iValue.MapIndex(mapkey)
			subkey := mapkey.String()
			if hasSubKey {
				subkey = key + subkey
			}
			setValues(fieldValue, values, subkey, hasSubKey)
		}
		break
	default:
		break
	}
}
