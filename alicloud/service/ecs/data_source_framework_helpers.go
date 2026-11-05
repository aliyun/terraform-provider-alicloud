package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const ecsFwRetryTimeout = 5 * time.Minute

// Shared conversions, request helpers and remote-call plumbing for the
// Framework data sources in this package. They keep response decoding tolerant
// of the two shapes the ECS RPC returns for the same field (JSON number or
// decimal string) without panicking on a missing key, and keep the query,
// pagination and retry behavior identical across the data sources.

func ecsFwString(value interface{}) types.String {
	if value == nil {
		return types.StringValue("")
	}
	if s, ok := value.(string); ok {
		return types.StringValue(s)
	}
	return types.StringValue(fmt.Sprint(value))
}

func ecsFwInt64(value interface{}) types.Int64 {
	switch v := value.(type) {
	case nil:
		return types.Int64Value(0)
	case int:
		return types.Int64Value(int64(v))
	case int32:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case float32:
		return types.Int64Value(int64(v))
	case float64:
		return types.Int64Value(int64(v))
	case json.Number:
		// The ECS RPC response is decoded with UseNumber, so integer members
		// arrive as json.Number rather than as int or float64.
		if n, err := v.Int64(); err == nil {
			return types.Int64Value(n)
		}
		if f, err := v.Float64(); err == nil {
			return types.Int64Value(int64(f))
		}
	case string:
		trimmed := strings.TrimSpace(v)
		if n, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
			return types.Int64Value(n)
		}
		if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
			return types.Int64Value(int64(f))
		}
	}
	return types.Int64Value(0)
}

func ecsFwBool(value interface{}) types.Bool {
	switch v := value.(type) {
	case bool:
		return types.BoolValue(v)
	case string:
		if b, err := strconv.ParseBool(strings.TrimSpace(v)); err == nil {
			return types.BoolValue(b)
		}
	}
	return types.BoolValue(false)
}

func ecsFwStringMapValue(values map[string]string) types.Map {
	elements := make(map[string]attr.Value, len(values))
	for k, v := range values {
		elements[k] = types.StringValue(v)
	}
	return types.MapValueMust(types.StringType, elements)
}

// ecsFwTagsValue decodes the `Tags` member of an ECS object. The RPC wraps the
// entries in a `Tag` list; an absent or differently shaped member is an empty
// map rather than a failure.
func ecsFwTagsValue(tags interface{}) types.Map {
	values := map[string]string{}
	if tagsMap, ok := tags.(map[string]interface{}); ok {
		if list, ok := tagsMap["Tag"]; ok {
			for key, value := range tagsToMap(list) {
				values[key] = fmt.Sprint(value)
			}
		}
	}
	return ecsFwStringMapValue(values)
}

func ecsFwStringListValue(values []string) types.List {
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	return types.ListValueMust(types.StringType, elements)
}

func ecsFwElementsAsStrings(value types.List) []string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	out := make([]string, 0, len(value.Elements()))
	for _, element := range value.Elements() {
		s, ok := element.(types.String)
		if !ok {
			continue
		}
		out = append(out, s.ValueString())
	}
	return out
}

func ecsFwSetRequestString(request map[string]interface{}, key string, value types.String) {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return
	}
	request[key] = value.ValueString()
}

// ecsFwApplyRepeatList writes a repeat-list filter in the 1-based
// `Parameter.N` form using the shared convert helper.
func ecsFwApplyRepeatList(request map[string]interface{}, parameter string, values []string) {
	raw := make([]interface{}, 0, len(values))
	for _, value := range values {
		raw = append(raw, value)
	}
	expandArrayToMap(request, raw, parameter)
}

// ecsFwApplyTags maps the Terraform tag map onto the `Tag.N.Key/Value` filter
// using the shared tags helper.
func ecsFwApplyTags(request map[string]interface{}, tags types.Map) {
	if tags.IsNull() || tags.IsUnknown() {
		return
	}
	list := make([]map[string]interface{}, 0, len(tags.Elements()))
	for key, value := range tags.Elements() {
		stringValue, ok := value.(types.String)
		if !ok {
			continue
		}
		list = append(list, map[string]interface{}{
			"Key":   key,
			"Value": stringValue.ValueString(),
		})
	}
	expandTagsToMap(request, list)
}

const ecsFwTagMaxLength = 128

// ecsFwTagMapValidator enforces the key and value limits that the ECS tag
// filters document: the key must be present and must not use the reserved
// aliyun/acs: prefixes, the value length is capped, and neither side may carry
// a URL. The value emptiness rule differs per API, so it is configurable.
type ecsFwTagMapValidator struct {
	allowEmptyValue        bool
	valueForbiddenPrefixes []string
}

func (v ecsFwTagMapValidator) Description(context.Context) string {
	return "Each tag key and value must satisfy the tag limits documented for the API."
}

func (v ecsFwTagMapValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ecsFwTagMapValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for key, element := range req.ConfigValue.Elements() {
		value, ok := element.(types.String)
		if !ok || value.IsNull() || value.IsUnknown() {
			continue
		}
		if message := validateEcsFwTagKey(key); message != "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Tag Key", message)
			return
		}
		if message := validateEcsFwTagValue(v, value.ValueString()); message != "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Tag Value", message)
			return
		}
	}
}

func validateEcsFwTagKey(key string) string {
	if key == "" {
		return "A tag key must not be empty."
	}
	if utf8.RuneCountInString(key) > ecsFwTagMaxLength {
		return fmt.Sprintf("A tag key must be at most %d characters.", ecsFwTagMaxLength)
	}
	for _, prefix := range []string{"aliyun", "acs:"} {
		if strings.HasPrefix(key, prefix) {
			return fmt.Sprintf("A tag key must not start with %q.", prefix)
		}
	}
	if strings.Contains(key, "http://") || strings.Contains(key, "https://") {
		return "A tag key must not contain a URL."
	}
	return ""
}

func validateEcsFwTagValue(v ecsFwTagMapValidator, value string) string {
	if value == "" && !v.allowEmptyValue {
		return "A tag value must not be empty."
	}
	if utf8.RuneCountInString(value) > ecsFwTagMaxLength {
		return fmt.Sprintf("A tag value must be at most %d characters.", ecsFwTagMaxLength)
	}
	for _, prefix := range v.valueForbiddenPrefixes {
		if strings.HasPrefix(value, prefix) {
			return fmt.Sprintf("A tag value must not start with %q.", prefix)
		}
	}
	if strings.Contains(value, "http://") || strings.Contains(value, "https://") {
		return "A tag value must not contain a URL."
	}
	return ""
}

// ecsFwRpcPost sends one read-only ECS query with the shared retry policy:
// transient failures are retried inside the bounded window, everything else
// fails on the first response.
func ecsFwRpcPost(ctx context.Context, client *connectivity.AliyunClient, typeName, action string, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := sdkretry.RetryContext(ctx, ecsFwRetryTimeout, func() *sdkretry.RetryError {
		result, rpcErr := client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if rpcErr != nil {
			if helper.NeedRetry(rpcErr) {
				wait()
				return sdkretry.RetryableError(rpcErr)
			}
			return sdkretry.NonRetryableError(rpcErr)
		}
		response = result
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, typeName, action, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// ecsFwCollectNumberedPages reads numbered pages until a short or empty page,
// the server total is reached, or the page budget is exhausted. A full page for
// the whole budget means the endpoint ignored the page parameters, so the
// result is discarded instead of returned incomplete.
func ecsFwCollectNumberedPages(maxPages int, startPage, pageSize int64, action string, fetch func(pageNumber int64) ([]map[string]interface{}, int64, error)) ([]map[string]interface{}, int64, error) {
	objects := make([]map[string]interface{}, 0)
	totalCount := int64(0)
	pageNumber := startPage
	for page := 0; page < maxPages; page++ {
		pageItems, pageTotal, err := fetch(pageNumber)
		if err != nil {
			return nil, 0, err
		}
		objects = append(objects, pageItems...)
		if pageTotal > 0 {
			totalCount = pageTotal
		}
		if len(pageItems) == 0 || int64(len(pageItems)) < pageSize {
			break
		}
		if totalCount > 0 && int64(len(objects)) >= totalCount {
			break
		}
		pageNumber++
		if page == maxPages-1 {
			return nil, 0, fmt.Errorf("%s returned a full page for %d consecutive requests (PageSize=%d), so pagination did not terminate; the result was not returned because it would be incomplete", action, maxPages, pageSize)
		}
	}
	if totalCount == 0 && len(objects) > 0 {
		totalCount = int64(len(objects))
	}
	return objects, totalCount, nil
}
