package connectivity

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ConvertKebabToSnake(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

func isInteger(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.String:
		if _, err := strconv.Atoi(fmt.Sprint(value)); err == nil {
			return true
		}
		return false
	default:
		return false
	}
}

func isString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}
