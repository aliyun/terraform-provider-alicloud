package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// rocketmqDataSourcePageSize is inside the 10-200 range that the RocketMQ
	// list APIs accept and keeps the number of round trips low.
	rocketmqDataSourcePageSize = 100
	// rocketmqDataSourceMaxPages bounds a pagination loop that never shrinks.
	// RocketMQ quotas are far below this, so hitting it means the endpoint
	// ignored pageNumber and returning the accumulated rows would silently be
	// partial.
	rocketmqDataSourceMaxPages = 1000
)

type rocketmqTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

var rocketmqTagObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}}

type rocketmqTagRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// rocketmqNameRegexValidator rejects a `name_regex` that cannot compile, so the
// mismatch surfaces at plan time and no request is sent.
type rocketmqNameRegexValidator struct{}

func (v rocketmqNameRegexValidator) Description(context.Context) string {
	return "The value must be a valid regular expression."
}

func (v rocketmqNameRegexValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v rocketmqNameRegexValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if _, err := regexp.Compile(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Name Regex", err.Error())
	}
}

// rocketmqCollectPages reads every page of a list API and drops rows whose
// primary key was already returned by an earlier page, so an unstable page
// cursor cannot make the same object appear twice. The RocketMQ list responses
// carry totalCount, but the loop stops on the first page shorter than the
// requested size instead so an outdated count cannot truncate the result.
func rocketmqCollectPages(action string, key func(item map[string]interface{}) string, fetch func(pageNumber int) ([]map[string]interface{}, error)) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	seen := make(map[string]struct{})
	for pageNumber := 1; pageNumber <= rocketmqDataSourceMaxPages; pageNumber++ {
		pageItems, err := fetch(pageNumber)
		if err != nil {
			return nil, err
		}
		for _, item := range pageItems {
			id := key(item)
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			items = append(items, item)
		}
		// The termination decision uses the raw page length, not the number of
		// rows kept, so duplicates cannot end a walk that still has pages.
		if len(pageItems) < rocketmqDataSourcePageSize {
			return items, nil
		}
	}
	return nil, fmt.Errorf("%s returned a full page for %d consecutive requests (pageSize=%d), so pagination did not terminate; the result was not returned because it would be incomplete", action, rocketmqDataSourceMaxPages, rocketmqDataSourcePageSize)
}

// rocketmqListItems extracts data.list from a list response. The list APIs omit
// data.list when nothing matched, so totalCount==0 is the legitimate empty
// result; an empty array is empty as well. Everything else - a body-less
// response, a missing list while rows are reported, a null or non-array list,
// or a non-object element - is a parse failure that must reach the caller
// instead of looking like a shorter page.
func rocketmqListItems(response map[string]interface{}) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	if response == nil {
		return nil, fmt.Errorf("RocketMQ returned an empty response body, so the result list cannot be read")
	}
	if success, ok := response["success"].(bool); ok && !success {
		return nil, fmt.Errorf("RocketMQ reported an unsuccessful response (code %q): %v", rocketmqString(response, "code"), response["message"])
	}
	raw, err := jsonpath.Get("$.data.list", response)
	if err != nil {
		// ListInstances returns {"data":{"pageNumber":1,"pageSize":10,"totalCount":0}}
		// with no list key at all when no instance matches.
		totalCount, countErr := jsonpath.Get("$.data.totalCount", response)
		if countErr == nil && rocketmqCountIsZero(totalCount) {
			return items, nil
		}
		return nil, fmt.Errorf("RocketMQ response has no readable $.data.list and does not report an empty result: %w", err)
	}
	if raw == nil {
		return nil, fmt.Errorf("RocketMQ response carries $.data.list as null instead of an array")
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("RocketMQ response carries $.data.list as %T instead of an array", raw)
	}
	for index, v := range list {
		item, ok := v.(map[string]interface{})
		if !ok {
			// Dropping the element would shrink the page and let the walk stop
			// early with a partial result.
			return nil, fmt.Errorf("RocketMQ response carries a non-object element at $.data.list[%d] (type %T)", index, v)
		}
		items = append(items, item)
	}
	return items, nil
}

// rocketmqCountIsZero reports whether a decoded JSON number is zero. A missing
// or unreadable count is not zero, so the caller keeps treating the response as
// malformed.
func rocketmqCountIsZero(value interface{}) bool {
	switch v := value.(type) {
	case float64:
		return v == 0
	case json.Number:
		n, err := v.Float64()
		return err == nil && n == 0
	case int:
		return v == 0
	case int64:
		return v == 0
	}
	return false
}

func rocketmqString(item map[string]interface{}, key string) string {
	if v, ok := item[key].(string); ok {
		return v
	}
	return ""
}

func rocketmqInt64(item map[string]interface{}, key string) types.Int64 {
	switch v := item[key].(type) {
	case float64:
		return types.Int64Value(int64(v))
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return types.Int64Value(n)
		}
	case int:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	}
	return types.Int64Value(0)
}

func rocketmqBool(item map[string]interface{}, key string) types.Bool {
	if v, ok := item[key].(bool); ok {
		return types.BoolValue(v)
	}
	return types.BoolValue(false)
}

func rocketmqTagsFromAPI(raw interface{}) []rocketmqTagModel {
	tags := make([]rocketmqTagModel, 0)
	list, ok := raw.([]interface{})
	if !ok {
		return tags
	}
	for _, v := range list {
		tag, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		tags = append(tags, rocketmqTagModel{
			Key:   types.StringValue(rocketmqString(tag, "key")),
			Value: types.StringValue(rocketmqString(tag, "value")),
		})
	}
	sort.Slice(tags, func(i, j int) bool { return tags[i].Key.ValueString() < tags[j].Key.ValueString() })
	return tags
}

// rocketmqTagsRequestValue renders the configured tags in the JSON form that the
// RocketMQ list APIs accept for their `tags` query parameter.
func rocketmqTagsRequestValue(tags map[string]string) string {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	payload := make([]rocketmqTagRequest, 0, len(keys))
	for _, k := range keys {
		payload = append(payload, rocketmqTagRequest{Key: k, Value: tags[k]})
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(encoded)
}

// rocketmqTagsMatch reports whether every configured tag is present on the
// resource, so a server that ignores the `tags` parameter cannot widen the
// result.
func rocketmqTagsMatch(item map[string]interface{}, want map[string]string) bool {
	if len(want) == 0 {
		return true
	}
	raw := make(map[string]string)
	list, ok := item["tags"].([]interface{})
	if !ok {
		return false
	}
	for _, v := range list {
		tag, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		raw[rocketmqString(tag, "key")] = rocketmqString(tag, "value")
	}
	for k, v := range want {
		if got, ok := raw[k]; !ok || got != v {
			return false
		}
	}
	return true
}

func rocketmqDataSourceId(values ...string) string {
	return helper.DataResourceIdHash(values)
}
