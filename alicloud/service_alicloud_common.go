package alicloud

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

const (
	// indentation symbol
	INDENTATIONSYMBOLCOMMON = " "

	// child field indend number
	CHILDINDENDCOMMON = 2
)

// deal with the parameter common method
func ValueConvert(indentation int, val reflect.Value) interface{} {
	switch val.Kind() {
	case reflect.Interface:
		return ValueConvert(indentation, reflect.ValueOf(val.Interface()))
	case reflect.String:
		valStr := val.String()
		if strings.HasPrefix(valStr, "${") && strings.HasSuffix(valStr, "}") {
			valStr = strings.TrimPrefix(valStr, "${")
			valStr = strings.TrimSuffix(valStr, "}")
			if !isValid(valStr) {
				valStr = "${" + valStr + "}"
				valStr = fmt.Sprintf("\"%s\"", valStr)
			}
		} else if !strings.HasPrefix(valStr, "var.") &&
			!strings.HasPrefix(valStr, "data.") &&
			!strings.HasPrefix(valStr, "alicloud_") {

			valStr = strings.TrimSuffix(valStr, "-update")
			valStr = strings.TrimSuffix(valStr, "_update")
			valStr = strings.TrimSuffix(valStr, "update")
			valStr = fmt.Sprintf("\"%s\"", valStr)
		}
		return valStr
	case reflect.Slice:
		return listValueCommon(indentation, val)
	case reflect.Map:
		return mapValueCommon(indentation, val)
	case reflect.Bool:
		return val.Bool()
	case reflect.Int:
		return val.Int()
	case reflect.Float64:
		return val.Float()
	default:
		log.Panicf("invalid attribute value type: %s,%#v", reflect.TypeOf(val), val)
	}
	return ""
}

func isValid(str string) bool {
	left := 0
	for _, c := range str {
		if c == '{' {
			left++
		} else if c == '}' {
			left--
		}
		if left < 0 {
			return false
		}
	}
	return left == 0
}

// deal with list parameter
func listValueCommon(indentation int, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentationCommon(indentation+CHILDINDENDCOMMON)+
			fmt.Sprint(ValueConvert(indentation+CHILDINDENDCOMMON, val.Index(i))))
	}

	return fmt.Sprintf("[\n%s\n%s]", strings.Join(valList, ",\n"), addIndentationCommon(indentation))
}

// deal with map parameter
func mapValueCommon(indentation int, val reflect.Value) string {
	var valList []string
	for _, keyV := range val.MapKeys() {
		mapVal := getRealValueTypeCommon(val.MapIndex(keyV))
		var line string
		if mapVal.Kind() == reflect.Slice && mapVal.Len() > 0 {
			eleVal := getRealValueTypeCommon(mapVal.Index(0))
			if eleVal.Kind() == reflect.Map {
				line = fmt.Sprintf(`%s%s`, addIndentationCommon(indentation),
					listValueMapChildCommon(indentation+CHILDINDENDCOMMON, keyV.String(), mapVal))
				valList = append(valList, line)
				continue
			}
		}
		value := ValueConvert(indentation+len(keyV.String())+CHILDINDENDCOMMON+3, val.MapIndex(keyV))
		switch value.(type) {
		case bool:
			line = fmt.Sprintf(`%s%s = %t`, addIndentationCommon(indentation+CHILDINDENDCOMMON), keyV.String(), value)
		case int:
			line = fmt.Sprintf(`%s%s = %d`, addIndentationCommon(indentation+CHILDINDENDCOMMON), keyV.String(), value)
		default:
			line = fmt.Sprintf(`%s%s = %s`, addIndentationCommon(indentation+CHILDINDENDCOMMON), keyV.String(), value)
		}

		valList = append(valList, line)
	}
	return fmt.Sprintf("{\n%s\n%s}", strings.Join(valList, "\n"), addIndentationCommon(indentation))
}

// deal with list parameter that child element is map
func listValueMapChildCommon(indentation int, key string, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentationCommon(indentation)+key+" "+
			mapValueCommon(indentation, getRealValueTypeCommon(val.Index(i))))
	}

	return fmt.Sprintf("%s\n%s", strings.Join(valList, "\n"), addIndentationCommon(indentation))
}

func getRealValueTypeCommon(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Interface:
		return getRealValueTypeCommon(reflect.ValueOf(value.Interface()))
	default:
		return value
	}
}

func addIndentationCommon(indentation int) string {
	return strings.Repeat(INDENTATIONSYMBOLCOMMON, indentation)
}
