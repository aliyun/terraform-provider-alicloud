package ecs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

// ecsFwParseInt64 converts the numeric shapes the ECS RPC returns. The boolean
// reports whether the value really was a number, so a caller that must not
// mistake a malformed value for 0 can tell the two apart.
func ecsFwParseInt64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	case json.Number:
		// The ECS RPC response is decoded with UseNumber, so integer members
		// arrive as json.Number rather than as int or float64.
		if n, err := v.Int64(); err == nil {
			return n, true
		}
		if f, err := v.Float64(); err == nil {
			return int64(f), true
		}
	case string:
		trimmed := strings.TrimSpace(v)
		if n, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
			return n, true
		}
		if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
			return int64(f), true
		}
	}
	return 0, false
}

func ecsFwInt64(value interface{}) types.Int64 {
	if n, ok := ecsFwParseInt64(value); ok {
		return types.Int64Value(n)
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

// ecsFwReservedTagPrefixes mirrors the reserved tag namespaces that the shared
// tags helper drops from data source state.
var ecsFwReservedTagPrefixes = []string{"aliyun", "acs:", "http://", "https://", "sae.do.not.delete"}

func ecsFwTagKeyReserved(key string) bool {
	for _, prefix := range ecsFwReservedTagPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// ecsFwDecodeTags decodes the `Tags` member of an ECS record. An absent member
// is a valid empty map, because the RPC omits it when the record has no tags. A
// member of the wrong shape, or an entry without a string key, is a malformed
// record: it is reported as an error instead of being dropped or panicking on a
// bare type assertion.
func ecsFwDecodeTags(tags interface{}) (types.Map, error) {
	values := map[string]string{}
	if tags == nil {
		return ecsFwStringMapValue(values), nil
	}
	tagsMap, ok := tags.(map[string]interface{})
	if !ok {
		return types.MapNull(types.StringType), fmt.Errorf("Tags is %T instead of an object", tags)
	}
	rawList, ok := tagsMap["Tag"]
	if !ok || rawList == nil {
		return ecsFwStringMapValue(values), nil
	}
	list, ok := rawList.([]interface{})
	if !ok {
		return types.MapNull(types.StringType), fmt.Errorf("Tags.Tag is %T instead of a list", rawList)
	}
	for index, entry := range list {
		item, ok := entry.(map[string]interface{})
		if !ok {
			return types.MapNull(types.StringType), fmt.Errorf("Tags.Tag[%d] is %T instead of an object", index, entry)
		}
		rawKey, ok := item["TagKey"]
		if !ok {
			rawKey = item["Key"]
		}
		key, ok := rawKey.(string)
		if !ok || key == "" {
			return types.MapNull(types.StringType), fmt.Errorf("Tags.Tag[%d] has no string key", index)
		}
		rawValue, ok := item["TagValue"]
		if !ok {
			rawValue = item["Value"]
		}
		value := ""
		if rawValue != nil {
			text, ok := rawValue.(string)
			if !ok {
				return types.MapNull(types.StringType), fmt.Errorf("Tags.Tag[%d].TagValue is %T instead of a string", index, rawValue)
			}
			value = text
		}
		if ecsFwTagKeyReserved(key) {
			continue
		}
		values[key] = value
	}
	return ecsFwStringMapValue(values), nil
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
// using the shared tags helper. A null or unknown value is left out: an empty
// string is a value the caller asked to send, while a null is not.
func ecsFwApplyTags(request map[string]interface{}, tags types.Map) {
	if tags.IsNull() || tags.IsUnknown() {
		return
	}
	list := make([]map[string]interface{}, 0, len(tags.Elements()))
	for key, value := range tags.Elements() {
		stringValue, ok := value.(types.String)
		if !ok || stringValue.IsNull() || stringValue.IsUnknown() {
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

// ValidateMap checks the key of every entry even when its value is null or
// unknown. A null value is rejected rather than serialized as an empty string:
// the two are different requests, and only an explicit empty string asks the
// API for an empty value.
func (v ecsFwTagMapValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for key, element := range req.ConfigValue.Elements() {
		if message := validateEcsFwTagKey(key); message != "" {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Tag Key", message)
			return
		}
		value, ok := element.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Tag Value", fmt.Sprintf("The value of tag key %q must be a string.", key))
			return
		}
		if value.IsUnknown() {
			// Not known yet: it cannot be validated here, and the request
			// builder leaves it out instead of sending an empty value.
			continue
		}
		if value.IsNull() {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Tag Value", fmt.Sprintf("The value of tag key %q must be set; use an empty string to send an empty value.", key))
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

// ecsFwSleepBackoff waits for the retry backoff and returns as soon as the read
// context is cancelled, so a cancelled read stops instead of sleeping through
// the remaining backoff.
func ecsFwSleepBackoff(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// ecsFwRpcPost sends one read-only ECS query with the shared retry policy:
// transient failures are retried inside the bounded window, everything else
// fails on the first response. The context carries the read budget, so the
// retry window, the backoff between attempts and the RPC itself all stop at the
// same deadline and on the same cancellation signal.
func ecsFwRpcPost(ctx context.Context, client *connectivity.AliyunClient, typeName, action string, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	backoff := 3 * time.Second
	err := sdkretry.RetryContext(ctx, ecsFwRetryTimeout, func() *sdkretry.RetryError {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return sdkretry.NonRetryableError(ctxErr)
		}
		result, rpcErr := client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if rpcErr != nil {
			if helper.NeedRetry(rpcErr) {
				if waitErr := ecsFwSleepBackoff(ctx, backoff); waitErr != nil {
					return sdkretry.NonRetryableError(waitErr)
				}
				if backoff < 30*time.Second {
					backoff += 3 * time.Second
				}
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

// ecsFwDecodeRecordList decodes a list of records nested under
// `container.member`. A missing container or member is a valid empty result:
// the ECS RPC omits them when nothing matches. A member of the wrong shape, or
// an entry that is not an object, is a malformed response and is reported as an
// error instead of being returned as an empty or partial result.
func ecsFwDecodeRecordList(response map[string]interface{}, container, member, action string) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	if response == nil {
		return nil, fmt.Errorf("%s returned an empty response body", action)
	}
	rawContainer, ok := response[container]
	if !ok || rawContainer == nil {
		return items, nil
	}
	containerMap, ok := rawContainer.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%s returned %s as %T instead of an object", action, container, rawContainer)
	}
	rawList, ok := containerMap[member]
	if !ok || rawList == nil {
		return items, nil
	}
	list, ok := rawList.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s returned %s.%s as %T instead of a list", action, container, member, rawList)
	}
	for index, entry := range list {
		item, ok := entry.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s returned %s.%s[%d] as %T instead of an object", action, container, member, index, entry)
		}
		items = append(items, item)
	}
	return items, nil
}

// ecsFwDecodeCount reads an integer response member. A missing member is
// reported as unknown so the caller can fall back; a member that is present but
// not an integer is a malformed response, not a count of 0.
func ecsFwDecodeCount(response map[string]interface{}, key, action string) (int64, bool, error) {
	raw, ok := response[key]
	if !ok || raw == nil {
		return 0, false, nil
	}
	value, ok := ecsFwParseInt64(raw)
	if !ok {
		return 0, false, fmt.Errorf("%s returned %s as %T instead of an integer", action, key, raw)
	}
	return value, true, nil
}

// ecsFwDecodeToken reads a string pagination token. A missing member means the
// server has no further page; a member of the wrong type is a malformed
// response and must not be treated as the end of the result.
func ecsFwDecodeToken(response map[string]interface{}, key, action string) (string, error) {
	raw, ok := response[key]
	if !ok || raw == nil {
		return "", nil
	}
	token, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s returned %s as %T instead of a string", action, key, raw)
	}
	return token, nil
}

// ecsFwPageSignature identifies a page by its decoded content. Two pages with
// the same signature mean the endpoint ignored the page parameter, which would
// otherwise duplicate rows or loop until the page budget is exhausted.
func ecsFwPageSignature(items []map[string]interface{}) string {
	hash := sha256.New()
	for _, item := range items {
		fmt.Fprintf(hash, "%v|", item)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// ecsFwCollectNumberedPages reads numbered pages until a short or empty page,
// the server total is reached, or the page budget is exhausted. The context
// carries the read deadline and cancellation signal, so the loop stops at the
// same point as the requests it issues. A page that repeats earlier content, or
// a full page for the whole budget, means the endpoint ignored the page
// parameters, so the result is discarded instead of returned incomplete.
func ecsFwCollectNumberedPages(ctx context.Context, maxPages int, startPage, pageSize int64, action string, fetch func(pageNumber int64) ([]map[string]interface{}, int64, error)) ([]map[string]interface{}, int64, error) {
	objects := make([]map[string]interface{}, 0)
	totalCount := int64(0)
	pageNumber := startPage
	seenPages := make(map[string]int64, maxPages)
	for page := 0; page < maxPages; page++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, 0, fmt.Errorf("%s read stopped before page %d was requested: %v; the result was not returned because it would be incomplete", action, pageNumber, ctxErr)
		}
		pageItems, pageTotal, err := fetch(pageNumber)
		if err != nil {
			return nil, 0, err
		}
		if len(pageItems) > 0 {
			signature := ecsFwPageSignature(pageItems)
			if previous, ok := seenPages[signature]; ok {
				return nil, 0, fmt.Errorf("%s returned the same content for pages %d and %d, so the page parameters are not honoured; the result was not returned because it would be incomplete", action, previous, pageNumber)
			}
			seenPages[signature] = pageNumber
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

// ecsFwCollectTokenPages follows a continuation token until the server returns
// none. A repeated token or a token for the whole page budget means the cursor
// never terminates; the accumulated rows are discarded instead of being
// returned as a partial result. The context carries the read deadline.
func ecsFwCollectTokenPages(ctx context.Context, maxPages int, initialToken string, action string, fetch func(token string) ([]map[string]interface{}, string, error)) ([]map[string]interface{}, string, error) {
	objects := make([]map[string]interface{}, 0)
	seenTokens := map[string]struct{}{}
	token := initialToken
	for page := 0; page < maxPages; page++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, "", fmt.Errorf("%s read stopped before page %d was requested: %v; the result was not returned because it would be incomplete", action, page+1, ctxErr)
		}
		pageItems, nextToken, err := fetch(token)
		if err != nil {
			return nil, "", err
		}
		objects = append(objects, pageItems...)
		if nextToken == "" {
			return objects, "", nil
		}
		if _, duplicate := seenTokens[nextToken]; duplicate {
			return nil, "", fmt.Errorf("%s returned the same continuation token %q twice, so pagination did not terminate; the result was not returned because it would be incomplete", action, nextToken)
		}
		seenTokens[nextToken] = struct{}{}
		token = nextToken
	}
	return nil, "", fmt.Errorf("%s returned a continuation token for %d consecutive requests, so pagination did not terminate; the result was not returned because it would be incomplete", action, maxPages)
}
