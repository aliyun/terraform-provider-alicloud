package agentloop

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Package-local counterparts of the flat alicloud package helpers that the
// shared errs/helper/logs packages do not export.

const failedToReachTargetStatus = "Failed to reach target status. Last status: %s."

func isExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*errs.ComplexError); ok {
		return isExpectedErrors(e.Cause, expectCodes)
	}

	if e, ok := err.(*tea.SDKError); ok {
		for _, code := range expectCodes {
			// The second statement aims to match the tea sdk history bug
			if tea.StringValue(e.Code) == code || strings.HasPrefix(code, tea.StringValue(e.Code)) || strings.Contains(tea.StringValue(e.Data), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*errs.ProviderError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	for _, code := range expectCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

func buildStateConf(pending, target []string, timeout, delay time.Duration, f retry.StateRefreshFunc) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    f,
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: 3 * time.Second,
	}
}

func tagsSchemaWithElements() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func stringPointer(s string) *string {
	return &s
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.String:
		return fmt.Sprint(i) == ""
	case reflect.Slice:
		return reflect.ValueOf(i).Len() <= 0
	case reflect.Map:
		return reflect.ValueOf(i).Len() <= 0
	case reflect.Ptr:
		return reflect.ValueOf(i).IsNil()
	default:
		return false
	}
}

func convertToInterfaceArray(v interface{}) []interface{} {
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

// formatInt converts the numeric types returned by teadsl to int.
func formatInt(src interface{}) int {
	if src == nil {
		return 0
	}
	switch v := src.(type) {
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int64:
		return int(v)
	case int32:
		return int(v)
	case int:
		return v
	case string:
		if v == "" {
			return 0
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		return i
	case json.Number:
		i, err := strconv.Atoi(v.String())
		if err != nil {
			panic(err)
		}
		return i
	default:
		panic(fmt.Sprintf("Not support type %T", src))
	}
}
