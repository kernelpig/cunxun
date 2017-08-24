package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func IsTypeDefaultValue(v interface{}) bool {
	vType := reflect.TypeOf(v)

	if vType.Name() == "string" {
		if v == "" {
			return true
		}
	} else if vType.Name() == "int" || vType.Name() == "int32" || vType.Name() == "int64" {
		if v == 0 {
			return true
		}
	} else {
		panic(fmt.Sprintf("not support type: %s", vType.Name()))
	}

	return false
}

func Struct2MapWithAddr(st interface{}, tagKey string) (string, map[string]interface{}) {
	result := make(map[string]interface{})

	valueElements := reflect.ValueOf(st).Elem()
	keyElements := reflect.TypeOf(st).Elem()

	// 只处理可导出的成员
	for i := 0; i < valueElements.NumField(); i++ {
		if !valueElements.Field(i).CanSet() {
			continue
		}
		tagValue := keyElements.Field(i).Name
		if len(tagKey) != 0 {
			tagValue = keyElements.Field(i).Tag.Get(tagKey)
		}
		result[tagValue] = valueElements.Field(i).Addr().Interface()
	}

	return strings.ToLower(keyElements.Name()), result
}

func Struct2MapWithValue(st interface{}, tagKey string, isParseDefault bool) (string, map[string]interface{}) {
	result := make(map[string]interface{})

	valueElements := reflect.ValueOf(st).Elem()
	keyElements := reflect.TypeOf(st).Elem()

	// 只处理可导出的成员, 开启过滤时默认类型值不处理
	for i := 0; i < valueElements.NumField(); i++ {
		if !valueElements.Field(i).CanSet() {
			continue
		}
		valueElement := valueElements.Field(i).Interface()
		if isParseDefault && IsTypeDefaultValue(valueElement) {
			continue
		}
		tagValue := keyElements.Field(i).Name
		if len(tagKey) != 0 {
			tagValue = keyElements.Field(i).Tag.Get(tagKey)
		}
		result[tagValue] = valueElement
	}

	return strings.ToLower(keyElements.Name()), result
}

func StructGetFieldName(st interface{}, tagKey string) (string, []string) {
	result := make([]string, 0)

	valueElements := reflect.ValueOf(st).Elem()
	keyElements := reflect.TypeOf(st).Elem()

	// 只处理可导出的成员
	for i := 0; i < valueElements.NumField(); i++ {
		if !valueElements.Field(i).CanSet() {
			continue
		}
		tagValue := keyElements.Field(i).Name
		if len(tagKey) != 0 {
			tagValue = keyElements.Field(i).Tag.Get(tagKey)
		}
		result = append(result, tagValue)
	}

	return strings.ToLower(keyElements.Name()), result
}

func StructGetFieldAddr(st interface{}) (string, []interface{}) {
	result := make([]interface{}, 0)

	valueElements := reflect.ValueOf(st).Elem()
	keyElements := reflect.TypeOf(st).Elem()

	// 只处理可导出的成员
	for i := 0; i < valueElements.NumField(); i++ {
		if !valueElements.Field(i).CanSet() {
			continue
		}
		result = append(result, valueElements.Field(i).Addr())
	}

	return strings.ToLower(keyElements.Name()), result
}

func StructGetFieldValue(st interface{}) (string, []interface{}) {
	result := make([]interface{}, 0)

	valueElements := reflect.ValueOf(st).Elem()
	keyElements := reflect.TypeOf(st).Elem()

	// 只处理可导出的成员
	for i := 0; i < valueElements.NumField(); i++ {
		if !valueElements.Field(i).CanSet() {
			continue
		}
		result = append(result, valueElements.Field(i).Interface())
	}

	return strings.ToLower(keyElements.Name()), result
}
