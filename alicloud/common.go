package alicloud

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/denverdino/aliyungo/cs"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/fc-go-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"gopkg.in/yaml.v2"

	"math"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/denverdino/aliyungo/common"
	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
)

type PayType string

const (
	PrePaid    = PayType("PrePaid")
	PostPaid   = PayType("PostPaid")
	Prepaid    = PayType("Prepaid")
	Postpaid   = PayType("Postpaid")
	Serverless = PayType("Serverless")
)

const (
	NormalMode = "normal"
	SafetyMode = "safety"
)

type DdosbgpInsatnceType string

const (
	Enterprise   = DdosbgpInsatnceType("Enterprise")
	Professional = DdosbgpInsatnceType("Professional")
)

type DdosbgpInstanceIpType string

const (
	IPv4 = DdosbgpInstanceIpType("IPv4")
	IPv6 = DdosbgpInstanceIpType("IPv6")
)

type NetType string

const (
	Internet = NetType("Internet")
	Intranet = NetType("Intranet")
)

type NetworkType string

const (
	Classic         = NetworkType("Classic")
	Vpc             = NetworkType("Vpc")
	ClassicInternet = NetworkType("classic_internet")
	ClassicIntranet = NetworkType("classic_intranet")
	PUBLIC          = NetworkType("PUBLIC")
	PRIVATE         = NetworkType("PRIVATE")
)

type NodeType string

const (
	WORKER = NodeType("WORKER")
	KIBANA = NodeType("KIBANA")
)

type ActionType string

const (
	OPEN  = ActionType("OPEN")
	CLOSE = ActionType("CLOSE")
)

type TimeType string

const (
	Hour  = TimeType("Hour")
	Day   = TimeType("Day")
	Week  = TimeType("Week")
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

type IpVersion string

const (
	IPV4 = IpVersion("ipv4")
	IPV6 = IpVersion("ipv6")
)

type Status string

const (
	Pending     = Status("Pending")
	Creating    = Status("Creating")
	Running     = Status("Running")
	Available   = Status("Available")
	Unavailable = Status("Unavailable")
	Modifying   = Status("Modifying")
	Deleting    = Status("Deleting")
	Starting    = Status("Starting")
	Stopping    = Status("Stopping")
	Stopped     = Status("Stopped")
	Normal      = Status("Normal")
	Changing    = Status("Changing")
	Online      = Status("online")
	Configuring = Status("configuring")

	Associating   = Status("Associating")
	Unassociating = Status("Unassociating")
	InUse         = Status("InUse")
	DiskInUse     = Status("In_use")

	Active   = Status("Active")
	Inactive = Status("Inactive")
	Idle     = Status("Idle")

	SoldOut = Status("SoldOut")

	InService      = Status("InService")
	Removing       = Status("Removing")
	DisabledStatus = Status("Disabled")

	Init            = Status("Init")
	Provisioning    = Status("Provisioning")
	Updating        = Status("Updating")
	FinancialLocked = Status("FinancialLocked")

	PUBLISHED   = Status("Published")
	NOPUBLISHED = Status("NonPublished")

	Deleted = Status("Deleted")
	Null    = Status("Null")

	Enable = Status("Enable")
	BINDED = Status("BINDED")
)

type IPType string

const (
	Inner   = IPType("Inner")
	Private = IPType("Private")
	Public  = IPType("Public")
)

type ResourceType string

const (
	ResourceTypeInstance      = ResourceType("Instance")
	ResourceTypeDisk          = ResourceType("Disk")
	ResourceTypeVSwitch       = ResourceType("VSwitch")
	ResourceTypeRds           = ResourceType("Rds")
	ResourceTypePolarDB       = ResourceType("PolarDB")
	IoOptimized               = ResourceType("IoOptimized")
	ResourceTypeRkv           = ResourceType("KVStore")
	ResourceTypeFC            = ResourceType("FunctionCompute")
	ResourceTypeElasticsearch = ResourceType("Elasticsearch")
	ResourceTypeSlb           = ResourceType("Slb")
	ResourceTypeMongoDB       = ResourceType("MongoDB")
	ResourceTypeGpdb          = ResourceType("Gpdb")
	ResourceTypeHBase         = ResourceType("HBase")
	ResourceTypeAdb           = ResourceType("ADB")
	ResourceTypeCassandra     = ResourceType("Cassandra")
)

type InternetChargeType string

const (
	PayByBandwidth = InternetChargeType("PayByBandwidth")
	PayByTraffic   = InternetChargeType("PayByTraffic")
	PayBy95        = InternetChargeType("PayBy95")
)

type AccountSite string

const (
	DomesticSite = AccountSite("Domestic")
	IntlSite     = AccountSite("International")
)
const (
	SnapshotCreatingInProcessing = Status("progressing")
	SnapshotCreatingAccomplished = Status("accomplished")
	SnapshotCreatingFailed       = Status("failed")

	SnapshotPolicyCreating  = Status("Creating")
	SnapshotPolicyAvailable = Status("available")
	SnapshotPolicyNormal    = Status("Normal")
)

// timeout for common product, ecs e.g.
const DefaultTimeout = 120
const Timeout5Minute = 300
const DefaultTimeoutMedium = 500

// timeout for long time progerss product, rds e.g.
const DefaultLongTimeout = 1000

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

const (
	// HeaderEnableEBTrigger header key for enabling eventbridge trigger
	// TODO: delete the header after eventbridge trigger is totally exposed to user
	HeaderEnableEBTrigger = "x-fc-enable-eventbridge-trigger"
)

// ValidProtocols network protocol list
var ValidProtocols = []Protocol{Http, Https, Tcp, Udp}

// simple array value check method, support string type only
func isProtocolValid(value string) bool {
	res := false
	for _, v := range ValidProtocols {
		if string(v) == value {
			res = true
		}
	}
	return res
}

// default region for all resource
const DEFAULT_REGION = "cn-beijing"

const INT_MAX = 2147483647

// symbol of multiIZ
const MULTI_IZ_SYMBOL = "MAZ"

const COMMA_SEPARATED = ","

const COLON_SEPARATED = ":"

const SLASH_SEPARATED = "/"

const LOCAL_HOST_IP = "127.0.0.1"

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
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
func convertListStringToListInterface(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func expandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(int))
	}
	return vs
}

// Convert the result for an array and returns a Json string
func convertListToJsonString(configured []interface{}) string {
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

func convertJsonStringToStringList(src interface{}) (result []interface{}) {
	if err, ok := src.([]interface{}); !ok {
		panic(err)
	}
	for _, v := range src.([]interface{}) {
		result = append(result, fmt.Sprint(formatInt(v)))
	}
	return
}

func encodeToBase64String(configured []string) string {
	result := ""
	for i, v := range configured {
		result += v
		if i < len(configured)-1 {
			result += ","
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(result))
}

func decodeFromBase64String(configured string) (result []string, err error) {

	decodeString, err := base64.StdEncoding.DecodeString(configured)
	if err != nil {
		return result, err
	}

	result = strings.Split(string(decodeString), ",")
	return result, nil
}

func convertJsonStringToMap(configured string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if err := json.Unmarshal([]byte(configured), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Convert the result for an array and returns a comma separate
func convertListToCommaSeparate(configured []interface{}) string {
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

func filterEmptyStrings(arr []interface{}) []interface{} {
	var result []interface{}
	for _, str := range arr {
		if fmt.Sprint(str) != "" {
			result = append(result, str)
		}
	}
	return result
}

func convertBoolToString(configured bool) string {
	return strconv.FormatBool(configured)
}

func convertStringToBool(configured string) bool {
	v, _ := strconv.ParseBool(configured)
	return v
}

func convertIntergerToString(configured int) string {
	return strconv.Itoa(configured)
}

func convertFloat64ToString(configured float64) string {
	return strconv.FormatFloat(configured, 'E', -1, 64)
}

func convertJsonStringToList(configured string) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(configured), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func expandArrayToMap(originMap map[string]interface{}, arrayValues []interface{}, arrayKey string) map[string]interface{} {
	for i, val := range arrayValues {
		key := fmt.Sprintf("%s.%d", arrayKey, i+1)
		originMap[key] = fmt.Sprint(val)
	}
	return originMap
}

func convertJsonStringToObject(configured interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if err := json.Unmarshal([]byte(configured.(string)), &result); err != nil {
		return nil
	}

	return result
}

func convertObjectToJsonString(m interface{}) string {
	if result, err := json.Marshal(m); err != nil {
		return ""
	} else {
		return string(result)
	}
}

func convertMaptoJsonString(m map[string]interface{}) (string, error) {
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

func convertMapToJsonStringIgnoreError(m map[string]interface{}) string {
	if result, err := json.Marshal(m); err != nil {
		return ""
	} else {
		return string(result)
	}
}

func convertInterfaceToJsonString(m interface{}) (string, error) {
	if result, err := json.Marshal(m); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

func convertListMapToJsonString(configured []map[string]interface{}) (string, error) {
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

func convertMapFloat64ToJsonString(m map[string]interface{}) (string, error) {
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

func StringPointer(s string) *string {
	return &s
}

func BoolPointer(b bool) *bool {
	return &b
}

func Int32Pointer(i int32) *int32 {
	return &i
}

func Int64Pointer(i int64) *int64 {
	return &i
}

func IntMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

const ServerSideEncryptionAes256 = "AES256"
const ServerSideEncryptionKMS = "KMS"
const ServerSideEncryptionSM4 = "SM4"

type OptimizedType string

const (
	IOOptimized   = OptimizedType("optimized")
	NoneOptimized = OptimizedType("none")
)

type TagResourceType string

const (
	TagResourceImage         = TagResourceType("image")
	TagResourceInstance      = TagResourceType("instance")
	TagResourceAcl           = TagResourceType("acl")
	TagResourceCertificate   = TagResourceType("certificate")
	TagResourceDisk          = TagResourceType("disk")
	TagResourceSecurityGroup = TagResourceType("securitygroup")
	TagResourceCdn           = TagResourceType("DOMAIN")
	TagResourceApp           = TagResourceType("app")
	TagResourceTopic         = TagResourceType("topic")
	TagResourceCluster       = TagResourceType("cluster")
)

type KubernetesNodeType string

const (
	KubernetesNodeMaster = ResourceType("Master")
	KubernetesNodeWorker = ResourceType("Worker")
)

func getPagination(pageNumber, pageSize int) (pagination common.Pagination) {
	pagination.PageSize = pageSize
	pagination.PageNumber = pageNumber
	return
}

const CharityPageUrl = "http://promotion.alicdn.com/help/oss/error.html"

func userDataHashSum(user_data string) string {
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

func GetUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Get current user got an error: %#v.", err)
	}
	return usr.HomeDir, nil
}

func writeToFile(filePath string, data interface{}) error {
	var out string
	switch data.(type) {
	case string:
		out = data.(string)
		break
	case nil:
		return nil
	default:
		bs, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v got an error: %#v", data, err)
		}
		out = string(bs)
	}

	if strings.HasPrefix(filePath, "~") {
		home, err := GetUserHomeDir()
		if err != nil {
			return err
		}
		if home != "" {
			filePath = strings.Replace(filePath, "~", home, 1)
		}
	}

	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(filePath, []byte(out), 422)
}

type Invoker struct {
	catchers []*Catcher
}

type Catcher struct {
	Reason           string
	RetryCount       int
	RetryWaitSeconds int
}

var ClientErrorCatcher = Catcher{AliyunGoClientFailure, 10, 5}
var ServiceBusyCatcher = Catcher{"ServiceUnavailable", 10, 5}
var ThrottlingCatcher = Catcher{Throttling, 50, 2}

func NewInvoker() Invoker {
	i := Invoker{}
	i.AddCatcher(ClientErrorCatcher)
	i.AddCatcher(ServiceBusyCatcher)
	i.AddCatcher(ThrottlingCatcher)
	return i
}

func (a *Invoker) AddCatcher(catcher Catcher) {
	a.catchers = append(a.catchers, &catcher)
}

func (a *Invoker) Run(f func() error) error {
	err := f()

	if err == nil {
		return nil
	}

	for _, catcher := range a.catchers {
		if IsExpectedErrors(err, []string{catcher.Reason}) {
			catcher.RetryCount--

			if catcher.RetryCount <= 0 {
				return fmt.Errorf("Retry timeout and got an error: %#v.", err)
			} else {
				time.Sleep(time.Duration(catcher.RetryWaitSeconds) * time.Second)
				return a.Run(f)
			}
		}
	}
	return err
}

func buildClientToken(action string) string {
	token := strings.TrimSpace(fmt.Sprintf("TF-%s-%d-%s", action, time.Now().Unix(), strings.Trim(uuid.New().String(), "-")))
	if len(token) > 64 {
		token = token[0:64]
	}
	return token
}

func getNextpageNumber(number requests.Integer) (requests.Integer, error) {
	page, err := strconv.Atoi(string(number))
	if err != nil {
		return "", err
	}
	return requests.NewInteger(page + 1), nil
}

func terraformToAPI(field string) string {
	var result string
	for _, v := range strings.Split(field, "_") {
		if len(v) > 0 {
			result = fmt.Sprintf("%s%s%s", result, strings.ToUpper(string(v[0])), v[1:])
		}
	}
	return result
}

func compareJsonTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	obj1 := make(map[string]interface{})
	err := json.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalJson1, _ := json.Marshal(obj1)

	obj2 := make(map[string]interface{})
	err = json.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalJson2, _ := json.Marshal(obj2)

	equal := bytes.Compare(canonicalJson1, canonicalJson2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalJson1, canonicalJson2)
	}
	return equal, nil
}

func compareArrayJsonTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	obj1 := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalJson1, _ := json.Marshal(obj1)

	obj2 := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalJson2, _ := json.Marshal(obj2)

	equal := bytes.Compare(canonicalJson1, canonicalJson2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalJson1, canonicalJson2)
	}

	return equal, nil
}

func compareYamlTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	var obj1 interface{}
	err := yaml.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalYaml1, _ := yaml.Marshal(obj1)

	var obj2 interface{}
	err = yaml.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalYaml2, _ := yaml.Marshal(obj2)

	equal := bytes.Compare(canonicalYaml1, canonicalYaml2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalYaml1, canonicalYaml2)
	}
	return equal, nil
}

// loadFileContent returns contents of a file in a given path
func loadFileContent(v string) ([]byte, error) {
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

func debugOn() bool {
	for _, part := range strings.Split(os.Getenv("DEBUG"), ",") {
		if strings.TrimSpace(part) == "terraform" {
			return true
		}
	}
	return false
}

func addDebug(action, content interface{}, requestInfo ...interface{}) {
	if debugOn() {
		trace := "[DEBUG TRACE]:\n"
		for skip := 1; skip < 5; skip++ {
			_, filepath, line, _ := runtime.Caller(skip)
			trace += fmt.Sprintf("%s:%d\n", filepath, line)
		}

		if len(requestInfo) > 0 {
			var request = struct {
				Domain     string
				Version    string
				UserAgent  string
				ActionName string
				Method     string
				Product    string
				Region     string
				AK         string
			}{}
			switch requestInfo[0].(type) {
			case *requests.RpcRequest:
				tmp := requestInfo[0].(*requests.RpcRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *requests.RoaRequest:
				tmp := requestInfo[0].(*requests.RoaRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *requests.CommonRequest:
				tmp := requestInfo[0].(*requests.CommonRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *fc.Client:
				client := requestInfo[0].(*fc.Client)
				request.Version = client.Config.APIVersion
				request.Product = "FC"
				request.ActionName = fmt.Sprintf("%s", action)
			case *sls.Client:
				request.Product = "LOG"
				request.ActionName = fmt.Sprintf("%s", action)
			case *tablestore.TableStoreClient:
				request.Product = "OTS"
				request.ActionName = fmt.Sprintf("%s", action)
			case *oss.Client:
				request.Product = "OSS"
				request.ActionName = fmt.Sprintf("%s", action)
			case *datahub.DataHub:
				request.Product = "DataHub"
				request.ActionName = fmt.Sprintf("%s", action)
			case *cs.Client:
				request.Product = "CS"
				request.ActionName = fmt.Sprintf("%s", action)
			}

			requestContent := ""
			if len(requestInfo) > 1 {
				switch requestInfo[1].(type) {
				case *tea.SDKError:
					requestContent = fmt.Sprintf("%#v", requestInfo[1].(*tea.SDKError).Error())
				default:
					requestContent = fmt.Sprintf("%#v", requestInfo[1])
				}
			}

			if len(requestInfo) == 1 {
				if v, ok := requestInfo[0].(map[string]interface{}); ok {
					if res, err := json.Marshal(&v); err == nil {
						requestContent = string(res)
					}
					if res, err := json.Marshal(&content); err == nil {
						content = string(res)
					}
				}
			}

			content = fmt.Sprintf("%vDomain:%v, Version:%v, ActionName:%v, Method:%v, Product:%v, Region:%v\n\n"+
				"*************** %s Request ***************\n%#v\n",
				content, request.Domain, request.Version, request.ActionName,
				request.Method, request.Product, request.Region, request.ActionName, requestContent)
		}

		//fmt.Printf(DefaultDebugMsg, action, content, trace)
		log.Printf(DefaultDebugMsg, action, content, trace)
	}
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

func ParseResourceIds(id string) (parts []string, err error) {
	parts = strings.Split(id, ":")
	return parts, err
}

func ParseResourceId(id string, length int) (parts []string, err error) {
	parts = strings.Split(id, ":")

	if len(parts) != length {
		err = WrapError(fmt.Errorf("Invalid Resource Id %s. Expected parts' length %d, got %d", id, length, len(parts)))
	}
	return parts, err
}

func ParseResourceIdN(id string, length int) (parts []string, err error) {
	parts = strings.SplitN(id, ":", length)

	if len(parts) != length {
		err = WrapError(fmt.Errorf("Invalid Resource Id %s. Expected parts' length %d, got %d", id, length, len(parts)))
	}
	return parts, err
}

func ParseResourceIdWithEscaped(id string, length int) (parts []string, err error) {
	parts = make([]string, 0)
	var currentPart strings.Builder
	for i := 0; i < len(id); i++ {
		if id[i] == '\\' {
			i++
			if i < len(id) {
				currentPart.WriteByte(id[i])
			}
		} else if id[i] == ':' {
			parts = append(parts, currentPart.String())
			currentPart.Reset()
		} else {
			currentPart.WriteByte(id[i])
		}
	}

	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	if len(parts) != length {
		err = WrapError(fmt.Errorf("Invalid Resource Id %s. Expected parts' length %d, got %d", id, length, len(parts)))
	}
	return parts, err
}

func EscapeColons(s string) string {
	return strings.ReplaceAll(s, ":", "\\:")
}

func ParseSlbListenerId(id string) (parts []string, err error) {
	parts = strings.Split(id, ":")
	if len(parts) != 2 && len(parts) != 3 {
		err = WrapError(fmt.Errorf("Invalid alicloud_slb_listener Id %s. Expected Id format is <slb id>:<protocol>:< frontend>.", id))
	}
	return parts, err
}

func GetCenChildInstanceType(id string) (c string, e error) {
	if strings.HasPrefix(id, "vpc") {
		return ChildInstanceTypeVpc, nil
	} else if strings.HasPrefix(id, "vbr") {
		return ChildInstanceTypeVbr, nil
	} else if strings.HasPrefix(id, "ccn") {
		return ChildInstanceTypeCcn, nil
	} else {
		return c, fmt.Errorf("CEN child instance ID invalid. Now, it only supports VPC or VBR or CCN instance.")
	}
}

func BuildStateConf(pending, target []string, timeout, delay time.Duration, f resource.StateRefreshFunc) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    f,
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: 3 * time.Second,
	}
}

func incrementalWait(firstDuration time.Duration, increaseDuration time.Duration) func() {
	retryCount := 1
	return func() {
		var waitTime time.Duration
		if retryCount == 1 {
			waitTime = firstDuration
		} else if retryCount > 1 {
			waitTime += increaseDuration
		}
		time.Sleep(waitTime)
		retryCount++
	}
}

// If auto renew, the period computed from computePeriodByUnit will be changed
// This method used to compute a period accourding to current period and unit
func computePeriodByUnit(createTime, endTime interface{}, currentPeriod int, periodUnit string) (int, error) {
	if createTime == nil {
		return 0, WrapError(fmt.Errorf("createTime should not be nil"))
	}
	if endTime == nil {
		return 0, WrapError(fmt.Errorf("endTime should not be nil"))
	}
	var createTimeStr, endTimeStr string
	switch value := createTime.(type) {
	case int64:
		createTimeStr = time.Unix(createTime.(int64), 0).Format(time.RFC3339)
		endTimeStr = time.Unix(endTime.(int64), 0).Format(time.RFC3339)
	case string:
		createTimeStr = createTime.(string)
		endTimeStr = endTime.(string)
	default:
		return 0, WrapError(fmt.Errorf("Unsupported time type: %#v", value))
	}
	// currently, there is time value does not format as standard RFC3339
	UnStandardRFC3339 := "2006-01-02T15:04Z07:00"
	create, err := time.Parse(time.RFC3339, createTimeStr)
	if err != nil {
		log.Printf("Parase the CreateTime %#v failed and error is: %#v.", createTime, err)
		create, err = time.Parse(UnStandardRFC3339, createTimeStr)
		if err != nil {
			return 0, WrapError(err)
		}
	}
	end, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		log.Printf("Parase the EndTime %#v failed and error is: %#v.", endTime, err)
		end, err = time.Parse(UnStandardRFC3339, endTimeStr)
		if err != nil {
			return 0, WrapError(err)
		}
	}
	var period int
	switch periodUnit {
	case "Month":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 30))
	case "Week":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 7))
	case "Year":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 365))
	default:
		err = fmt.Errorf("Unexpected period unit %s", periodUnit)
	}
	// The period at least is 1
	if period < 1 {
		period = 1
	}
	if period > 12 {
		period = 12
	}
	// period can not be modified and if the new period is changed, using the previous one.
	if currentPeriod > 0 && currentPeriod != period {
		period = currentPeriod
	}
	return period, WrapError(err)
}

func checkWaitForReady(object interface{}, conditions map[string]interface{}) (bool, map[string]interface{}, error) {
	if conditions == nil {
		return false, nil, nil
	}
	objectType := reflect.TypeOf(object)
	objectValue := reflect.ValueOf(object)
	values := make(map[string]interface{})
	for key, value := range conditions {
		if _, ok := objectType.FieldByName(key); ok {
			current := objectValue.FieldByName(key)
			values[key] = current
			if fmt.Sprintf("%v", current) != fmt.Sprintf("%v", value) {
				return false, values, nil
			}
		} else {
			return false, values, WrapError(fmt.Errorf("There is missing attribute %s in the object.", key))
		}
	}
	return true, values, nil
}

// When  using teadsl, we need to convert float, int64 and int32 to int for comparison.
func formatInt(src interface{}) int {
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

func formatBool(src interface{}) bool {
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

func formatFloat64(src interface{}) float64 {
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

func convertArrayObjectToJsonString(src interface{}) (string, error) {
	res, err := json.Marshal(&src)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func convertArrayToString(src interface{}, sep string) string {
	if src == nil {
		return ""
	}
	items := make([]string, 0)
	for _, v := range src.([]interface{}) {
		items = append(items, fmt.Sprint(v))
	}
	return strings.Join(items, sep)
}

func splitMultiZoneId(id string) (ids []string) {
	if !(strings.Contains(id, MULTI_IZ_SYMBOL) || strings.Contains(id, "(")) {
		return
	}
	firstIndex := strings.Index(id, MULTI_IZ_SYMBOL)
	secondIndex := strings.Index(id, "(")
	for _, p := range strings.Split(id[secondIndex+1:len(id)-1], COMMA_SEPARATED) {
		ids = append(ids, id[:firstIndex]+string(p))
	}
	return
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

func isPagingRequest(d *schema.ResourceData) bool {
	v, ok := d.GetOk("page_number")
	return ok && v.(int) > 0
}

func setPagingRequest(d *schema.ResourceData, request map[string]interface{}, maxPageSize int) {
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

func mapMerge(target, merged map[string]interface{}) map[string]interface{} {
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
							targetSlice[index] = mapMerge(targetMap, val.(map[string]interface{}))
						} else {
							targetSlice[index] = mapMerge(map[string]interface{}{}, val.(map[string]interface{}))
						}
					default:
						targetSlice[index] = val
					}
				}
				target[key] = targetSlice
			case map[string]interface{}:
				target[key] = mapMerge(target[key].(map[string]interface{}), merged[key].(map[string]interface{}))
			default:
				target[key] = merged[key]
			}
		}
	}
	return target
}

func mapSort(target map[string]string) []string {
	result := make([]string, 0)
	for key := range target {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func newInstanceDiff(resourceName string, attributes, attributesDiff map[string]interface{}, state *terraform.InstanceState) (*terraform.InstanceDiff, error) {

	p := Provider().(*schema.Provider).ResourcesMap
	dOld, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	dNew, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	for key, value := range attributes {
		err := dOld.Set(key, value)
		if err != nil {
			return nil, WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}
	for key, value := range attributesDiff {
		attributes[key] = value
	}

	for key, value := range attributes {
		err := dNew.Set(key, value)
		if err != nil {
			return nil, WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}

	diff := terraform.NewInstanceDiff()
	objectKey := ""
	for _, key := range mapSort(dNew.State().Attributes) {
		newValue := dNew.State().Attributes[key]
		if objectKey != "" && !strings.HasPrefix(key, objectKey) {
			objectKey = ""
		}
		if objectKey == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					objectKey = strings.TrimSuffix(key, suffix)
					break
				}
			}
		}
		oldValue, ok := dOld.State().Attributes[key]
		if ok && oldValue == newValue {
			continue
		}
		if oldValue == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					oldValue = "0"
				}
			}
		}
		diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: oldValue, New: newValue})

		if objectKey != "" {
			for removeKey, removeValue := range dOld.State().Attributes {
				if strings.HasPrefix(removeKey, objectKey) {
					if _, ok := dNew.State().Attributes[removeKey]; !ok {
						// If the attribue has complex elements, there should remove the key, not setting it to empty
						if len(strings.Split(removeKey, ".")) > 2 {
							diff.DelAttribute(removeKey)
						} else {
							diff.SetAttribute(removeKey, &terraform.ResourceAttrDiff{Old: removeValue, New: ""})
						}
					}
				}
			}
			objectKey = ""
		}
	}
	return diff, nil
}

func compareMapWithIgnoreEquivalent(tem1, tem2 map[string]interface{}, ignore []string) bool {

	if len(tem1) != len(tem2) {
		return false
	}

OuterLoop:
	for key1, val1 := range tem1 {
		for _, item := range ignore {
			if key1 == item {
				continue OuterLoop
			}
		}

		val2 := tem2[key1]
		if val2 != val1 {
			return false
		}
	}

	return true
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
		return convertStringToBool(i.(string))
	case reflect.Bool:
		return i.(bool)
	default:
		return false
	}
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

func compareCmsHybridMonitorFcTaskYamlConfigAreEquivalent(tem1, tem2 string) (bool, error) {
	type MetricList struct {
		MetricList []string `yaml:"metric_list"`
	}
	type Product struct {
		MetricInfo []MetricList `yaml:"metric_info"`
		Namespace  string       `yaml:"namespace"`
	}
	type Products struct {
		Products []Product
	}

	var P1 Products
	err := yaml.Unmarshal([]byte(tem1), &P1)
	if err != nil {
		fmt.Sprintln(false)
	}

	y1 := make([]string, 0)
	for _, product := range P1.Products {
		s1, _ := json.Marshal(product)
		y1 = append(y1, string(s1))
	}

	var P2 Products
	err = yaml.Unmarshal([]byte(tem2), &P2)
	if err != nil {
		fmt.Sprintln(false)
	}

	y2 := make([]string, 0)
	for _, product := range P2.Products {
		s2, _ := json.Marshal(product)
		y2 = append(y2, string(s2))
	}

	sort.Strings(y1)
	sort.Strings(y2)
	return reflect.DeepEqual(y1, y2), nil
}

func getOneStringOrAllStringSlice(stringSli []interface{}) interface{} {
	if len(stringSli) == 1 {
		return stringSli[0].(string)
	}
	sli := make([]string, len(stringSli))
	for i, v := range stringSli {
		sli[i] = v.(string)
	}
	return sli
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

func rpcParam(method, version, action string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Pathname:    tea.String("/"),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
}
func roaParam(method, version, action, path string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Pathname:    tea.String(path),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
}

func xmlParam(method, version, action, path string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String(path),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		ReqBodyType: tea.String("xml"),
		BodyType:    tea.String("xml"),
	}
}

func jsonXmlParam(method, version, action, path string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String(path),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("xml"),
	}
}

func xmlJsonParam(method, version, action, path string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String(path),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		ReqBodyType: tea.String("xml"),
		BodyType:    tea.String("json"),
	}
}

type MyMap map[string]interface{}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   interface{} `xml:",chardata"`
}

func (m MyMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

func expandSingletonToList(singleton interface{}) []interface{} {
	vs := make([]interface{}, 0)
	vs = append(vs, singleton)
	return vs
}

func MD5(b []byte) string {
	ctx := md5.New()
	ctx.Write(b)
	return hex.EncodeToString(ctx.Sum(nil))
}

func ConvertTags(tagsMap map[string]interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, 0)
	for key, value := range tagsMap {
		if value != nil {
			if v, ok := value.(string); ok {
				tags = append(tags, map[string]interface{}{
					"Key":   key,
					"Value": v,
				})
			}
		}
	}

	return tags
}

func ConvertTagsForKms(tagsMap map[string]interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, 0)
	for key, value := range tagsMap {
		if value != nil {
			if v, ok := value.(string); ok {
				tags = append(tags, map[string]interface{}{
					"TagKey":   key,
					"TagValue": v,
				})
			}
		}
	}

	return tags
}

func expandTagsToMap(originMap map[string]interface{}, tags []map[string]interface{}) map[string]interface{} {
	for i, tag := range tags {
		for key, value := range tag {
			if key == "Key" || key == "Value" {
				newKey := "Tag" + "." + strconv.Itoa(i+1) + "." + key
				originMap[newKey] = fmt.Sprintf("%v", value)
			}
		}
	}
	return originMap
}

func convertChargeTypeToPaymentType(source interface{}) interface{} {
	switch source {
	case "PostPaid", "Postpaid":
		return "PayAsYouGo"
	case "PrePaid", "Prepaid":
		return "Subscription"
	}
	return source
}

func convertPaymentTypeToChargeType(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}

func bytesToTB(bytes int64) float64 {
	const (
		KiB = 1024
		MiB = KiB * KiB
		GiB = MiB * KiB
		TiB = GiB * KiB
	)
	return float64(bytes) / float64(TiB)
}

func compressIPv6OrCIDR(input string) (string, error) {
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

func randIntRange(min int, max int) int {
	return min + acctest.RandIntRange(min, max)
}
