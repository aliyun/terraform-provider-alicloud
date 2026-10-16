package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Format helpers moved from the root package (alicloud/common.go);
// package alicloud re-exports them as var aliases so the flat layout
// keeps compiling unchanged.

// When  using teadsl, we need to convert float, int64 and int32 to int for comparison.
func FormatInt(src interface{}) int {
	if src == nil {
		return 0
	}
	attrType := reflect.TypeOf(src)
	switch attrType.String() {
	case "float64":
		return int(src.(float64))
	case "float32":
		return int(src.(float32))
	case "int64":
		return int(src.(int64))
	case "int32":
		return int(src.(int32))
	case "int":
		return src.(int)
	case "string":
		vv := fmt.Sprint(src)
		if vv == "" {
			return 0
		}
		v, err := strconv.Atoi(vv)
		if err != nil {
			panic(err)
		}
		return v
	case "json.Number":
		v, err := strconv.Atoi(src.(json.Number).String())
		if err != nil {
			panic(err)
		}
		return v
	default:
		panic(fmt.Sprintf("Not support type %s", attrType.String()))
	}
}

func FormatBool(src interface{}) bool {
	if src == nil {
		return false
	}
	attrType := reflect.TypeOf(src)
	switch attrType.String() {
	case "bool":
		return src.(bool)
	case "string":
		vv := fmt.Sprint(src)
		if vv == "" {
			return false
		}
		v, err := strconv.ParseBool(vv)
		if err != nil {
			panic(err)
		}
		return v
	default:
		panic(fmt.Sprintf("Not support type %s", attrType.String()))
	}
}

func FormatFloat64(src interface{}) float64 {
	if src == nil {
		return 0
	}
	attrType := reflect.TypeOf(src)
	switch attrType.String() {
	case "float64":
		return src.(float64)
	case "float32":
		return float64(src.(float32))
	case "int64":
		return float64(src.(int64))
	case "int32":
		return float64(src.(int32))
	case "int":
		return float64(src.(int))
	case "string":
		vv := fmt.Sprint(src)
		if vv == "" {
			return 0
		}
		v, err := strconv.ParseFloat(vv, 64)
		if err != nil {
			panic(err)
		}
		return v
	case "json.Number":
		v, err := src.(json.Number).Float64()
		if err != nil {
			panic(err)
		}
		return v
	default:
		panic(fmt.Sprintf("Not support type %s", attrType.String()))
	}
}
