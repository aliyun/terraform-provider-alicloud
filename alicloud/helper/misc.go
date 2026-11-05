package helper

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// Misc utilities moved from the root package (alicloud/common.go);
// package alicloud re-exports them as var aliases so the flat layout
// keeps compiling unchanged.

const DefaultIntervalMini = 2

const DefaultIntervalShort = 5

const DefaultIntervalMedium = 10

const DefaultIntervalLong = 20

const (
	PageNumSmall   = 1
	PageSizeSmall  = 10
	PageSizeMedium = 20
	PageSizeLarge  = 50
	PageSizeXLarge = 100
)

const ServerSideEncryptionAes256 = "AES256"
const ServerSideEncryptionKMS = "KMS"
const ServerSideEncryptionSM4 = "SM4"

// symbol of multiIZ
const MultiIzSymbol = "MAZ"

const CommaSeparated = ","

// Protocol represents network protocol
type Protocol string

// Constants of protocol definition
const (
	Http  = Protocol("http")
	Https = Protocol("https")
	Tcp   = Protocol("tcp")
	Udp   = Protocol("udp")
	All   = Protocol("all")
	Icmp  = Protocol("icmp")
	Gre   = Protocol("gre")
)

// ValidProtocols network protocol list
var ValidProtocols = []Protocol{Http, Https, Tcp, Udp}

// simple array value check method, support string type only
func IsProtocolValid(value string) bool {
	res := false
	for _, v := range ValidProtocols {
		if string(v) == value {
			res = true
		}
	}
	return res
}

func IntMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func GetPagination(pageNumber, pageSize int) (pagination common.Pagination) {
	pagination.PageSize = pageSize
	pagination.PageNumber = pageNumber
	return
}

func UserDataHashSum(user_data string) string {
	// Check whether the user_data is not Base64 encoded.
	// Always calculate hash of base64 decoded value since we
	// check against double-encoding when setting it
	v, base64DecodeError := base64.StdEncoding.DecodeString(user_data)
	if base64DecodeError != nil {
		v = []byte(user_data)
	}
	return string(v)
}

// Remove useless blank in the string.
func Trim(v string) string {
	if len(v) < 1 {
		return v
	}
	return strings.Trim(v, " ")
}

func ConvertIntegerToInt(value requests.Integer) (v int, err error) {
	if strings.TrimSpace(string(value)) == "" {
		return
	}
	v, err = strconv.Atoi(string(value))
	if err != nil {
		return v, fmt.Errorf("Converting integer %s to int got an error: %#v.", value, err)
	}
	return
}

func GetNextpageNumber(number requests.Integer) (requests.Integer, error) {
	page, err := strconv.Atoi(string(number))
	if err != nil {
		return "", err
	}
	return requests.NewInteger(page + 1), nil
}

func TerraformToAPI(field string) string {
	var result string
	for _, v := range strings.Split(field, "_") {
		if len(v) > 0 {
			result = fmt.Sprintf("%s%s%s", result, strings.ToUpper(string(v[0])), v[1:])
		}
	}
	return result
}

// loadFileContent returns contents of a file in a given path
func LoadFileContent(v string) ([]byte, error) {
	filename, err := homedir.Expand(v)
	if err != nil {
		return nil, err
	}
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

// Return a ComplexError which including extra error message, error occurred file and path
func GetFunc(level int) string {
	pc, _, _, ok := runtime.Caller(level)
	if !ok {
		log.Printf("[ERROR] runtime.Caller error in GetFuncName.")
		return ""
	}
	return strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), ".")
}

func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// SplitSlice Divides the slice into blocks of the specified size
func SplitSlice(xs []interface{}, chunkSize int) [][]interface{} {
	if len(xs) == 0 {
		return nil
	}
	divided := make([][]interface{}, (len(xs)+chunkSize-1)/chunkSize)
	prev := 0
	i := 0
	till := len(xs) - chunkSize
	for prev < till {
		next := prev + chunkSize
		divided[i] = xs[prev:next]
		prev = next
		i++
	}
	divided[i] = xs[prev:]
	return divided
}

func SplitMultiZoneId(id string) (ids []string) {
	if !(strings.Contains(id, MultiIzSymbol) || strings.Contains(id, "(")) {
		return
	}
	firstIndex := strings.Index(id, MultiIzSymbol)
	secondIndex := strings.Index(id, "(")
	for _, p := range strings.Split(id[secondIndex+1:len(id)-1], CommaSeparated) {
		ids = append(ids, id[:firstIndex]+string(p))
	}
	return
}

func IsPagingRequest(d *schema.ResourceData) bool {
	v, ok := d.GetOk("page_number")
	return ok && v.(int) > 0
}

func SetPagingRequest(d *schema.ResourceData, request map[string]interface{}, maxPageSize int) {
	if maxPageSize == 0 {
		maxPageSize = PageSizeLarge
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = maxPageSize
	}
	return
}

func MapMerge(target, merged map[string]interface{}) map[string]interface{} {
	for key, value := range merged {
		if _, exist := target[key]; !exist {
			target[key] = value
		} else {
			// key existed in both src,target
			switch merged[key].(type) {
			case []interface{}:
				sourceSlice := value.([]interface{})
				targetSlice := make([]interface{}, len(sourceSlice))
				copy(targetSlice, target[key].([]interface{}))

				for index, val := range sourceSlice {
					switch val.(type) {
					case map[string]interface{}:
						targetMap, ok := targetSlice[index].(map[string]interface{})
						if ok {
							targetSlice[index] = MapMerge(targetMap, val.(map[string]interface{}))
						} else {
							targetSlice[index] = MapMerge(map[string]interface{}{}, val.(map[string]interface{}))
						}
					default:
						targetSlice[index] = val
					}
				}
				target[key] = targetSlice
			case map[string]interface{}:
				target[key] = MapMerge(target[key].(map[string]interface{}), merged[key].(map[string]interface{}))
			default:
				target[key] = merged[key]
			}
		}
	}
	return target
}

func MapSort(target map[string]string) []string {
	result := make([]string, 0)
	for key := range target {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func GetOneStringOrAllStringSlice(stringSli []interface{}, fieldName string) (interface{}, error) {
	if len(stringSli) == 1 {
		if Trim(fmt.Sprint(stringSli[0])) == "" || IsNil(stringSli[0]) {
			return nil, errs.WrapError(fmt.Errorf("[ERROR] Field \"%s\" cannot contain empty strings(\"\") ! ", fieldName))
		}

		return stringSli[0].(string), nil
	}

	sli := make([]string, len(stringSli))
	for i, v := range stringSli {
		if Trim(fmt.Sprint(v)) == "" || IsNil(v) {
			return nil, errs.WrapError(fmt.Errorf("[ERROR] Field \"%s\" cannot contain empty strings(\"\") ! ", fieldName))
		}

		sli[i] = v.(string)
	}

	return sli, nil
}

func Unique(strings []string) []string {
	dict := make(map[string]bool)
	var ss []string
	for _, s := range strings {
		if s == "" {
			continue
		}
		if _, ok := dict[s]; !ok {
			dict[s] = true
			ss = append(ss, s)
		}
	}
	return ss
}

func IsSubCollection(sub []string, full []string) bool {
	for _, s := range sub {
		var find bool
		for _, f := range full {
			if s == f {
				find = true
				break
			}
		}
		if !find {
			return false
		}
	}
	return true
}

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, m := range maps {
		for key, value := range m {
			item, existed := result[key]
			if !existed {
				result[key] = value
				continue
			}
			newValue, ok := value.([]map[string]interface{})
			if !ok || len(newValue) != 1 {
				continue
			}
			if preValue, ok := item.([]map[string]interface{}); ok && len(preValue) == 1 {
				result[key] = MergeMaps(preValue[0], newValue[0])
			}
		}
	}
	return result
}

func InArray(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

func MD5(b []byte) string {
	ctx := md5.New()
	ctx.Write(b)
	return hex.EncodeToString(ctx.Sum(nil))
}

func BytesToTB(bytes int64) float64 {
	const (
		KiB = 1024
		MiB = KiB * KiB
		GiB = MiB * KiB
		TiB = GiB * KiB
	)
	return float64(bytes) / float64(TiB)
}

func CompressIPv6OrCIDR(input string) (string, error) {
	if input == "" {
		return input, nil
	}
	if strings.Contains(input, "/") {
		ip, _, err := net.ParseCIDR(input)
		if err != nil {
			return "", err
		}
		if ip == nil {
			return input, nil
		}
		mask := strings.SplitN(input, "/", 2)[1]
		return fmt.Sprintf("%s/%s", ip.String(), mask), nil
	}
	ip := net.ParseIP(input)
	if ip == nil {
		return input, nil
	}
	return ip.String(), nil
}

func RandIntRange(min int, max int) int {
	return min + acctest.RandIntRange(min, max)
}

func IsEmpty(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.String:
		return fmt.Sprint(i) == ""
	case reflect.Int:
		return i.(int) <= 0
	case reflect.Int8:
		return i.(int8) <= 0
	case reflect.Int16:
		return i.(int16) <= 0
	case reflect.Int32:
		return i.(int32) <= 0
	case reflect.Int64:
		return i.(int64) <= 0
	case reflect.Float32:
		return i.(float32) <= 0
	case reflect.Float64:
		return i.(float64) <= 0
	case reflect.Map:
		return len(i.(map[string]interface{})) <= 0
	case reflect.Ptr:
		return reflect.ValueOf(i).IsNil()
	default:
		return false
	}
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.String:
		return fmt.Sprint(i) == ""
	case reflect.Slice:
		return len(i.([]interface{})) <= 0
	case reflect.Map:
		return len(i.(map[string]interface{})) <= 0
	case reflect.Ptr:
		return reflect.ValueOf(i).IsNil()
	default:
		return false
	}
}

func GetDaysBetween2Date(format string, date1Str string, date2Str string) (int, error) {
	var day int
	t1, err := time.ParseInLocation(format, date1Str, time.Local)
	if err != nil {
		return 0, err
	}
	t2, err := time.ParseInLocation(format, date2Str, time.Local)
	if err != nil {
		return 0, err
	}

	swap := false
	if t1.Unix() > t2.Unix() {
		t1, t2 = t2, t1
		swap = true
	}

	t1_ := t1.Add(time.Duration(t2.Sub(t1).Milliseconds()%86400000) * time.Millisecond)
	day = int(t2.Sub(t1).Hours() / 24)
	if t1_.Day() != t1.Day() {
		day += 1
	}

	if swap {
		day = -day
	}

	return day, nil
}

func GetWriteOnlyValue(d *schema.ResourceData, path cty.Path, typ cty.Type) (cty.Value, error) {
	if d.GetRawConfig().IsNull() {
		return cty.NullVal(typ), nil
	}
	value, diags := d.GetRawConfigAt(path)
	if diags.HasError() {
		return cty.NilVal, errs.WrapError(fmt.Errorf("error retrieving write-only argument: %s", diags[0].Summary))
	}
	if !value.Type().Equals(typ) {
		return cty.NilVal, errs.WrapError(fmt.Errorf("error retrieving write-only argument %s: unexpected type %s, want %s", path, value.Type().FriendlyName(), typ.FriendlyName()))
	}
	return value, nil
}
