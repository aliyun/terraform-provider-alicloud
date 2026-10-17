package paiplugin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PAIPlugin is a REST/ROA-style POP product: every operation is an HTTP verb
// against an /api/v2/* path, with body parameters (POST/PUT) or query
// parameters (GET). The connectivity client's Roa* methods speak that contract.
const (
	paiPluginProduct = "PaiPlugin"
	paiPluginVersion = "2022-01-12"

	// Base paths. A path parameter is appended with a leading slash.
	pathJobConfigs = "/api/v2/contentgeneration/jobconfigs"
	pathJobs       = "/api/v2/contentgeneration/jobs"
	pathMessages   = "/api/v2/messages"
	pathMsgMetrics = "/api/v2/messages/metrics"
	pathAlgorithms = "/api/v2/algorithms"
	pathUsers      = "/api/v2/users"
	pathMsgConfig  = "/api/v2/users/messageConfig"
)

// ptr returns the address of s, the shape Roa* query/header maps require.
func ptr(s string) *string { return &s }

// dataObject returns the Data object of a PAIPlugin response. The gateway
// wraps every payload in {"Data": …, "ErrorCode": 0, "ErrorMessage": "OK",
// "RequestId": "…"}; roaRequest already unwraps the tea client's "body" nesting,
// so the map passed in is the raw response shape and the fields live under Data.
func dataObject(response map[string]interface{}) map[string]interface{} {
	if response == nil {
		return nil
	}
	if d, ok := response["Data"].(map[string]interface{}); ok {
		return d
	}
	return nil
}

// errorCode reports the application ErrorCode carried in a response body. Zero
// means success; the cloudspec marks 400 as the resource-not-found signal.
// The tea decoder may surface the value as a json.Number or a float64, so both
// are accepted rather than assuming one concrete type.
func errorCode(response map[string]interface{}) int {
	if response == nil {
		return 0
	}
	raw, ok := response["ErrorCode"]
	if !ok || raw == nil {
		return 0
	}
	switch v := raw.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case jsonNumber:
		n, err := v.Int64()
		if err != nil {
			return 0
		}
		return int(n)
	case string:
		n, _ := strconv.Atoi(v)
		return n
	}
	return 0
}

// jsonNumber matches the methods of encoding/json.Number, so a decoded
// json.Number satisfies it without the service package importing the tea or
// encoding/json decoder types directly.
type jsonNumber interface {
	Int64() (int64, error)
	Float64() (float64, error)
	String() string
}

// isNotFound reports whether a Get response signals that the resource is gone.
// The cloudspec binds the not-found condition to ErrorCode 400, which the
// gateway may deliver as an HTTP 400 (surfaced as a Go error) or as a 200 body
// carrying ErrorCode 400. Both paths are checked so the Read path removes the
// resource from state instead of erroring on a drift.
func isNotFound(err error, response map[string]interface{}) bool {
	if err != nil {
		if strings.Contains(err.Error(), "400") || strings.Contains(err.Error(), "NotFound") {
			return true
		}
	}
	return errorCode(response) == 400
}

// str pulls a string field from a response object, mapping the JSON decoder's
// possible concrete types.
func str(obj map[string]interface{}, key string) string {
	if obj == nil {
		return ""
	}
	switch v := obj[key].(type) {
	case string:
		return v
	case jsonNumber:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return ""
}

// strValue maps a response string to a framework String, with empty → null so
// optional fields the API omits read back as absent rather than "".
func strValue(obj map[string]interface{}, key string) types.String {
	if s := str(obj, key); s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

// intValue pulls an integer field from a response object.
func intValue(obj map[string]interface{}, key string) int64 {
	if obj == nil {
		return 0
	}
	switch v := obj[key].(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	case jsonNumber:
		n, _ := v.Int64()
		return n
	case string:
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	}
	return 0
}

// floatValue pulls a float field from a response object.
func floatValue(obj map[string]interface{}, key string) float64 {
	if obj == nil {
		return 0
	}
	switch v := obj[key].(type) {
	case float64:
		return v
	case jsonNumber:
		f, _ := v.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	}
	return 0
}

// strSlice pulls a []string field from a response object. The JSON decoder
// delivers arrays as []interface{}; non-string elements are stringified to
// avoid silently dropping a value the API actually returned.
func strSlice(obj map[string]interface{}, key string) []string {
	if obj == nil {
		return nil
	}
	arr, ok := obj[key].([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, e := range arr {
		switch v := e.(type) {
		case string:
			out = append(out, v)
		case jsonNumber:
			out = append(out, v.String())
		case float64:
			out = append(out, strconv.FormatFloat(v, 'f', -1, 64))
		default:
			out = append(out, fmt.Sprintf("%v", v))
		}
	}
	return out
}

// strSliceValue maps a response string array to a framework List, returning
// null when the API omits the field entirely.
func strSliceValue(obj map[string]interface{}, key string) types.List {
	values := strSlice(obj, key)
	if values == nil {
		return types.ListNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	list, diags := types.ListValue(types.StringType, elements)
	_ = diags
	return list
}

// strSliceFromList extracts a []string from a framework String list, returning
// nil for a null/unknown list. It is the inverse of strSliceValue, used when
// building a request body from a planned state value.
func strSliceFromList(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	elements := make([]types.String, 0, len(list.Elements()))
	// List.Elements returns []attr.Value; each element of a String list is a
	// types.String, but asserting through attr.Value keeps the cast explicit.
	for _, el := range list.Elements() {
		if s, ok := el.(types.String); ok {
			elements = append(elements, s)
		}
	}
	out := make([]string, 0, len(elements))
	for _, e := range elements {
		out = append(out, e.ValueString())
	}
	return out
}

// intSliceValue maps a response int array (under a list of objects keyed by
// field) to a framework Int64 list. The caller has already projected each
// element to the integer field; this helper turns the slice into a List value.
func intSliceValue(values []int64) types.List {
	if values == nil {
		return types.ListNull(types.Int64Type)
	}
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.Int64Value(v))
	}
	list, diags := types.ListValue(types.Int64Type, elements)
	_ = diags
	return list
}

// floatSliceValue maps a float slice to a framework Float64 list.
func floatSliceValue(values []float64) types.List {
	if values == nil {
		return types.ListNull(types.Float64Type)
	}
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.Float64Value(v))
	}
	list, diags := types.ListValue(types.Float64Type, elements)
	_ = diags
	return list
}

// queryFrom builds a Roa* query map from non-empty string inputs, skipping
// empty values so they do not override server-side defaults.
func queryFrom(pairs ...string) map[string]*string {
	q := map[string]*string{}
	for i := 0; i+1 < len(pairs); i += 2 {
		if k, v := pairs[i], pairs[i+1]; v != "" {
			q[k] = ptr(v)
		}
	}
	return q
}

// roaGet is a thin wrapper around the connectivity ROA GET so call sites stay
// uniform and the product/version constants live in one place.
func roaGet(client *connectivity.AliyunClient, path string, query map[string]*string) (map[string]interface{}, error) {
	return client.RoaGet(paiPluginProduct, paiPluginVersion, path, query, nil, nil)
}
