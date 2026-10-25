package convert

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	yamlV3 "gopkg.in/yaml.v3"
)

// Conversion helpers moved from the root package (alicloud/common.go).
// Package alicloud re-exports them as var aliases so the flat layout
// keeps compiling unchanged, and service packages import them here
// directly. Exported spellings are verbatim from the root originals,
// including historical typos (ConvertIntergerToString,
// ConvertMaptoJsonString), so the migration generator's rewrite map
// stays a pure capitalize-first-letter mapping.

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func ExpandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		if v == nil {
			continue
		}
		vs = append(vs, v.(string))
	}
	return vs
}

// Takes list of string to strings. Expand to an array
// of raw strings and returns a []interface{}
func ConvertListStringToListInterface(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func ExpandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(int))
	}
	return vs
}

// Convert the result for an array and returns a Json string
func ConvertListToJsonString(configured []interface{}) string {
	if len(configured) < 1 {
		return ""
	}
	result := "["
	for i, v := range configured {
		if v == nil {
			continue
		}
		result += "\"" + v.(string) + "\""
		if i < len(configured)-1 {
			result += ","
		}
	}
	result += "]"
	return result
}

func ConvertJsonStringToStringList(src interface{}) (result []interface{}, err error) {
	items, ok := src.([]interface{})
	if !ok {
		return nil, errs.Error("Not support type %s to convert to string list", reflect.TypeOf(src))
	}
	for _, v := range items {
		vInt, err := FormatInt(v)
		if err != nil {
			return nil, err
		}
		result = append(result, fmt.Sprint(vInt))
	}
	return
}

func EncodeToBase64String(configured []string) string {
	result := ""
	for i, v := range configured {
		result += v
		if i < len(configured)-1 {
			result += ","
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(result))
}

func DecodeFromBase64String(configured string) (result []string, err error) {

	decodeString, err := base64.StdEncoding.DecodeString(configured)
	if err != nil {
		return result, err
	}

	result = strings.Split(string(decodeString), ",")
	return result, nil
}

func ConvertJsonStringToMap(configured string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if err := json.Unmarshal([]byte(configured), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Convert the result for an array and returns a comma separate
func ConvertListToCommaSeparate(configured []interface{}) string {
	if len(configured) < 1 {
		return ""
	}
	result := ""
	for i, v := range configured {
		if v == nil {
			continue
		}
		rail := ","
		if i == len(configured)-1 {
			rail = ""
		}
		result += fmt.Sprint(v) + rail
	}
	return result
}

func FilterEmptyStrings(arr []interface{}) []interface{} {
	var result []interface{}
	for _, str := range arr {
		if fmt.Sprint(str) != "" {
			result = append(result, str)
		}
	}
	return result
}

func ConvertBoolToString(configured bool) string {
	return strconv.FormatBool(configured)
}

func ConvertStringToBool(configured string) bool {
	v, _ := strconv.ParseBool(configured)
	return v
}

func ConvertIntergerToString(configured int) string {
	return strconv.Itoa(configured)
}

func ConvertFloat64ToString(configured float64) string {
	return strconv.FormatFloat(configured, 'E', -1, 64)
}

func ConvertJsonStringToList(configured string) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(configured), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func ExpandArrayToMap(originMap map[string]interface{}, arrayValues []interface{}, arrayKey string) map[string]interface{} {
	for i, val := range arrayValues {
		key := fmt.Sprintf("%s.%d", arrayKey, i+1)
		originMap[key] = fmt.Sprint(val)
	}
	return originMap
}

func ConvertJsonStringToObject(configured interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if err := json.Unmarshal([]byte(configured.(string)), &result); err != nil {
		return nil
	}

	return result
}

func ConvertObjectToJsonString(m interface{}) string {
	if result, err := json.Marshal(m); err != nil {
		return ""
	} else {
		return string(result)
	}
}

func ConvertMaptoJsonString(m map[string]interface{}) (string, error) {
	//sm := make(map[string]string, len(m))
	//for k, v := range m {
	//	sm[k] = v.(string)
	//}

	if result, err := json.Marshal(m); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

func ConvertMapToJsonStringIgnoreError(m map[string]interface{}) string {
	if result, err := json.Marshal(m); err != nil {
		return ""
	} else {
		return string(result)
	}
}

func ConvertInterfaceToJsonString(m interface{}) (string, error) {
	if result, err := json.Marshal(m); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

func ConvertListMapToJsonString(configured []map[string]interface{}) (string, error) {
	if len(configured) < 1 {
		return "[]", nil
	}

	result := "["
	for i, m := range configured {
		if m == nil {
			continue
		}

		sm := make(map[string]interface{}, len(m))
		for k, v := range m {
			sm[k] = v
		}

		item, err := json.Marshal(sm)
		if err == nil {
			result += string(item)
			if i < len(configured)-1 {
				result += ","
			}
		}
	}
	result += "]"
	return result, nil
}

func ConvertMapFloat64ToJsonString(m map[string]interface{}) (string, error) {
	sm := make(map[string]json.Number, len(m))

	for k, v := range m {
		sm[k] = v.(json.Number)
	}

	if result, err := json.Marshal(sm); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

func ConvertArrayObjectToJsonString(src interface{}) (string, error) {
	res, err := json.Marshal(&src)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func ConvertArrayToString(src interface{}, sep string) string {
	if src == nil {
		return ""
	}
	items := make([]string, 0)
	for _, v := range src.([]interface{}) {
		items = append(items, fmt.Sprint(v))
	}
	return strings.Join(items, sep)
}

func ConvertToJsonWithoutEscapeHTML(m map[string]interface{}) (string, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(m); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func ConvertToInterfaceArray(v interface{}) []interface{} {
	if v == nil {
		return []interface{}{}
	}

	if set, ok := v.(*schema.Set); ok {
		return set.List()
	}

	if array, ok := v.([]interface{}); ok {
		return array
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Slice {
		result := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = val.Index(i).Interface()
		}
		return result
	}

	return []interface{}{v}
}

func ConvertYamlToObject(configured interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if err := yamlV3.Unmarshal([]byte(configured.(string)), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func Interface2String(val interface{}) string {
	if v, ok := val.(string); ok {
		return v
	}
	return fmt.Sprint(val)
}

func Interface2StrSlice(ii []interface{}) []string {
	ss := make([]string, 0, len(ii))
	for _, i := range ii {
		s := Interface2String(i)
		ss = append(ss, s)
	}
	return ss
}

func Str2InterfaceSlice(ss []string) []interface{} {
	ii := make([]interface{}, 0, len(ss))
	for _, s := range ss {
		ii = append(ii, s)
	}
	return ii
}

func Interface2Bool(i interface{}) bool {
	if i == nil {
		return false
	}
	t := reflect.TypeOf(i).Kind()
	switch t {
	case reflect.String:
		return ConvertStringToBool(i.(string))
	case reflect.Bool:
		return i.(bool)
	default:
		return false
	}
}

func ExpandSingletonToList(singleton interface{}) []interface{} {
	vs := make([]interface{}, 0)
	vs = append(vs, singleton)
	return vs
}

func ConvertChargeTypeToPaymentType(source interface{}) interface{} {
	switch source {
	case "PostPaid", "Postpaid":
		return "PayAsYouGo"
	case "PrePaid", "Prepaid":
		return "Subscription"
	}
	return source
}

func ConvertPaymentTypeToChargeType(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}

func NormalizeMap(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range data {
		normalizedValue, err := NormalizeValue(v)
		if err != nil {
			result[k] = v
		} else {
			result[k] = normalizedValue
		}
	}
	return result
}

func NormalizeValue(value interface{}) (interface{}, error) {
	switch val := value.(type) {
	case string:
		trimmed := strings.TrimSpace(val)
		// handle JSON array
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			var arr []interface{}
			if err := json.Unmarshal([]byte(val), &arr); err == nil {
				result := make([]interface{}, len(arr))
				for i, item := range arr {
					normalizedItem, err := NormalizeValue(item)
					if err != nil {
						return nil, err
					}
					result[i] = normalizedItem
				}
				return result, nil
			}
		}

		// handle JSON object
		if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(val), &obj); err == nil {
				normalizedObj := NormalizeMap(obj)
				return normalizedObj, nil
			}
		}

		switch val {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}

		if num, err := strconv.Atoi(val); err == nil {
			return num, nil
		}

		return val, nil

	case map[string]interface{}:
		normalizedMap := NormalizeMap(val)
		return normalizedMap, nil

	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			normalizedItem, err := NormalizeValue(item)
			if err != nil {
				return nil, err
			}
			result[i] = normalizedItem
		}
		return result, nil

	default:
		return val, nil
	}
}

// NormalizeJsonObjectSliceByName normalizes a JSON object slice so that
// HasChange comparisons are not affected by element ordering. It accepts either
// a JSON-encoded string or a []interface{} of map[string]interface{}, sorts the
// objects by their "name" field with a stable sort, and returns the
// re-serialized JSON string.
// A JSON string that cannot be unmarshalled is returned unchanged. Any other
// type is rendered via fmt.Sprintf("%v") so the caller still gets a stable,
// non-empty string.
func NormalizeJsonObjectSliceByName(raw interface{}) string {
	switch t := raw.(type) {
	case string:
		var arr []map[string]interface{}
		if err := json.Unmarshal([]byte(t), &arr); err != nil {
			return t
		}
		sort.SliceStable(arr, func(i, j int) bool {
			n1, _ := arr[i]["name"].(string)
			n2, _ := arr[j]["name"].(string)
			return n1 < n2
		})
		b, _ := json.Marshal(arr)
		return string(b)
	case []interface{}:
		// directly convert to []map
		var arr []map[string]interface{}
		for _, v := range t {
			m, ok := v.(map[string]interface{})
			if ok {
				arr = append(arr, m)
			}
		}
		sort.SliceStable(arr, func(i, j int) bool {
			n1, _ := arr[i]["name"].(string)
			n2, _ := arr[j]["name"].(string)
			return n1 < n2
		})
		b, _ := json.Marshal(arr)
		return string(b)
	default:
		return fmt.Sprintf("%v", raw)
	}
}
