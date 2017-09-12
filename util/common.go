package util

import (
	"reflect"
	"strings"
)

// ModelType get value's model type
func ModelType(value interface{}) reflect.Type {
	reflectType := reflect.Indirect(reflect.ValueOf(value)).Type()

	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}

	return reflectType
}

// ConvertNameFromCamelToSnakeLike
// e.g. "OrderItem" -> "order_item"
func ConvertNameFromCamelToSnakeLike(str string) string {
	var snake []rune
	for i, l := range str {
		if i > 0 && IsUppercase(byte(l)) {
			if (!IsUppercase(str[i-1]) && str[i-1] != '_') || (i+1 < len(str) && !IsUppercase(str[i+1]) && str[i+1] != '_' && str[i-1] != '_') {
				snake = append(snake, rune('_'))
			}
		}
		snake = append(snake, l)
	}
	return strings.ToLower(string(snake))
}

func IsUppercase(char byte) bool {
	return 'A' <= char && char <= 'Z'
}