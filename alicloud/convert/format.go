package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// Format helpers moved from the root package (alicloud/common.go);
// package alicloud re-exports them as var aliases so the flat layout
// keeps compiling unchanged. Abnormal API data is reported through the
// returned error instead of panicking the provider process.

// When  using teadsl, we need to convert float, int64 and int32 to int for comparison.
func FormatInt(src interface{}) (int, error) {
	if src == nil {
		return 0, nil
	}
	attrType := reflect.TypeOf(src)
	switch attrType.String() {
	case "float64":
		return int(src.(float64)), nil
	case "float32":
		return int(src.(float32)), nil
	case "int64":
		return int(src.(int64)), nil
	case "int32":
		return int(src.(int32)), nil
	case "int":
		return src.(int), nil
	case "string":
		vv := fmt.Sprint(src)
		if vv == "" {
			return 0, nil
		}
		v, err := strconv.Atoi(vv)
		if err != nil {
			return 0, errs.Error("Converting %s value %q to int got an error: %#v", attrType.String(), vv, err)
		}
		return v, nil
	case "json.Number":
		v, err := strconv.Atoi(src.(json.Number).String())
		if err != nil {
			return 0, errs.Error("Converting %s value %s to int got an error: %#v", attrType.String(), src.(json.Number).String(), err)
		}
		return v, nil
	default:
		return 0, errs.Error("Not support type %s", attrType.String())
	}
}

func FormatBool(src interface{}) (bool, error) {
	if src == nil {
		return false, nil
	}
	attrType := reflect.TypeOf(src)
	switch attrType.String() {
	case "bool":
		return src.(bool), nil
	case "string":
		vv := fmt.Sprint(src)
		if vv == "" {
			return false, nil
		}
		v, err := strconv.ParseBool(vv)
		if err != nil {
			return false, errs.Error("Converting %s value %q to bool got an error: %#v", attrType.String(), vv, err)
		}
		return v, nil
	default:
		return false, errs.Error("Not support type %s", attrType.String())
	}
}

func FormatFloat64(src interface{}) (float64, error) {
	if src == nil {
		return 0, nil
	}
	attrType := reflect.TypeOf(src)
	switch attrType.String() {
	case "float64":
		return src.(float64), nil
	case "float32":
		return float64(src.(float32)), nil
	case "int64":
		return float64(src.(int64)), nil
	case "int32":
		return float64(src.(int32)), nil
	case "int":
		return float64(src.(int)), nil
	case "string":
		vv := fmt.Sprint(src)
		if vv == "" {
			return 0, nil
		}
		v, err := strconv.ParseFloat(vv, 64)
		if err != nil {
			return 0, errs.Error("Converting %s value %q to float64 got an error: %#v", attrType.String(), vv, err)
		}
		return v, nil
	case "json.Number":
		v, err := src.(json.Number).Float64()
		if err != nil {
			return 0, errs.Error("Converting %s value %s to float64 got an error: %#v", attrType.String(), src.(json.Number).String(), err)
		}
		return v, nil
	default:
		return 0, errs.Error("Not support type %s", attrType.String())
	}
}
