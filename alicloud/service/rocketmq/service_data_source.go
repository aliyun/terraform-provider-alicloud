package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// Hitting the page bound means the endpoint ignored pageNumber; returning the rows would silently be partial.
	rocketmqDataSourceMaxPages = 1000
)

// rocketmqNameRegexValidator rejects a `name_regex` that cannot compile, so the mismatch surfaces at plan time.
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

type rocketmqPage struct {
	Items      []map[string]interface{}
	PageNumber int64
	PageSize   int64
	TotalCount int64
}

// A short page must account for totalCount with unique items; a full page always reads another page so a stale low count cannot truncate results.
func rocketmqCollectPages(action string, key func(item map[string]interface{}) string, fetch func(pageNumber int) (rocketmqPage, error)) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	seen := make(map[string]struct{})
	for pageNumber := 1; pageNumber <= rocketmqDataSourceMaxPages; pageNumber++ {
		page, err := fetch(pageNumber)
		if err != nil {
			return nil, err
		}
		if page.PageNumber != int64(pageNumber) {
			return nil, fmt.Errorf("%s returned pageNumber %d for requested page %d; the result was not returned because pagination is inconsistent", action, page.PageNumber, pageNumber)
		}
		if page.PageSize <= 0 {
			return nil, fmt.Errorf("%s returned invalid pageSize %d on page %d", action, page.PageSize, pageNumber)
		}
		if page.TotalCount < 0 {
			return nil, fmt.Errorf("%s returned invalid totalCount %d on page %d", action, page.TotalCount, pageNumber)
		}
		if int64(len(page.Items)) > page.PageSize {
			return nil, fmt.Errorf("%s returned %d items with pageSize %d on page %d", action, len(page.Items), page.PageSize, pageNumber)
		}
		for _, item := range page.Items {
			id := key(item)
			if id == "" {
				return nil, fmt.Errorf("%s returned an item with an empty primary key on page %d; the result was not returned because it would be incomplete", action, pageNumber)
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			items = append(items, item)
		}
		if int64(len(page.Items)) < page.PageSize {
			if int64(len(items)) < page.TotalCount {
				return nil, fmt.Errorf("%s ended on page %d after %d unique of %d reported rows; the result was not returned because it would be incomplete", action, pageNumber, len(items), page.TotalCount)
			}
			return items, nil
		}
	}
	return nil, fmt.Errorf("%s returned a full page for %d consecutive requests (pageSize=%d), so pagination did not terminate; the result was not returned because it would be incomplete", action, rocketmqDataSourceMaxPages, helper.PageSizeXLarge)
}

func rocketmqListPage(response map[string]interface{}) (rocketmqPage, error) {
	page := rocketmqPage{Items: make([]map[string]interface{}, 0)}
	if response == nil {
		return page, fmt.Errorf("RocketMQ returned an empty response body, so the result list cannot be read")
	}
	if success, ok := response["success"].(bool); ok && !success {
		return page, fmt.Errorf("RocketMQ reported an unsuccessful response (code %q): %v", rocketmqString(response, "code"), response["message"])
	}
	var err error
	if page.PageNumber, err = rocketmqMetadataInt64(response, "$.data.pageNumber"); err != nil {
		return page, err
	}
	if page.PageSize, err = rocketmqMetadataInt64(response, "$.data.pageSize"); err != nil {
		return page, err
	}
	if page.TotalCount, err = rocketmqMetadataInt64(response, "$.data.totalCount"); err != nil {
		return page, err
	}
	if page.PageNumber < 1 || page.PageSize < 1 || page.TotalCount < 0 {
		return page, fmt.Errorf("RocketMQ returned invalid pagination metadata: pageNumber=%d pageSize=%d totalCount=%d", page.PageNumber, page.PageSize, page.TotalCount)
	}
	raw, err := jsonpath.Get("$.data.list", response)
	if err != nil {
		if page.TotalCount == 0 {
			return page, nil
		}
		return page, fmt.Errorf("RocketMQ response has no readable $.data.list and reports %d rows: %w", page.TotalCount, err)
	}
	if raw == nil {
		return page, fmt.Errorf("RocketMQ response carries $.data.list as null instead of an array")
	}
	list, ok := raw.([]interface{})
	if !ok {
		return page, fmt.Errorf("RocketMQ response carries $.data.list as %T instead of an array", raw)
	}
	if int64(len(list)) > page.PageSize {
		return page, fmt.Errorf("RocketMQ response carries %d rows with pageSize %d", len(list), page.PageSize)
	}
	for index, v := range list {
		item, ok := v.(map[string]interface{})
		if !ok {
			return page, fmt.Errorf("RocketMQ response carries a non-object element at $.data.list[%d] (type %T)", index, v)
		}
		page.Items = append(page.Items, item)
	}
	return page, nil
}

func rocketmqMetadataInt64(response map[string]interface{}, path string) (int64, error) {
	raw, err := jsonpath.Get(path, response)
	if err != nil {
		return 0, fmt.Errorf("RocketMQ response has no readable %s: %w", path, err)
	}
	value, err := rocketmqInt64Value(raw)
	if err != nil {
		return 0, fmt.Errorf("RocketMQ response field %s is not an integer: %w", path, err)
	}
	return value, nil
}

func rocketmqString(item map[string]interface{}, key string) string {
	if v, ok := item[key].(string); ok {
		return v
	}
	return ""
}

func rocketmqInt64Value(value interface{}) (int64, error) {
	switch v := value.(type) {
	case float64:
		if math.Trunc(v) == v && v >= math.MinInt64 && v < math.MaxInt64 {
			return int64(v), nil
		}
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return n, nil
		}
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	}
	return 0, fmt.Errorf("value %v (type %T) is not an integer", value, value)
}

func rocketmqInt64(item map[string]interface{}, key string) (types.Int64, error) {
	raw, ok := item[key]
	if !ok || raw == nil {
		return types.Int64Null(), fmt.Errorf("RocketMQ response is missing expected integer field %q", key)
	}
	value, err := rocketmqInt64Value(raw)
	if err != nil {
		return types.Int64Null(), fmt.Errorf("RocketMQ response field %q is invalid: %w", key, err)
	}
	return types.Int64Value(value), nil
}

func rocketmqOptionalInt64(item map[string]interface{}, key string) (types.Int64, error) {
	raw, ok := item[key]
	if !ok || raw == nil {
		return types.Int64Null(), nil
	}
	return rocketmqInt64(item, key)
}

func rocketmqDataSourceId(values ...string) string {
	return helper.DataResourceIdHash(values)
}
