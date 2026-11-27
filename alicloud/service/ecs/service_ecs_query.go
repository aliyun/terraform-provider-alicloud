package ecs

// Conversions, request builders and response decoding shared by the ECS query
// data sources (alicloud_ecs_auto_provisioning_groups,
// alicloud_ecs_image_pipeline_executions and alicloud_reserved_instances). The
// ECS RPC returns the same value shapes for them: a JSON number or a decimal
// string for a numeric member, an omitted member for an absent optional value,
// and the same paginated success envelope.
// These helpers keep that decoding, the request filters, the pagination
// integrity checks and the retry classification in one place. Response decoding
// tolerates the shapes the service returns without panicking on a missing key,
// and every malformed member is reported by name instead of being written to
// state.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ecsDecodeStringMember reads a string member of a record. The ECS RPC omits an
// optional string member that carries no value, so an absent member stays an
// empty string. A member that is present but is null, an object or an array is
// not a string this contract can return: it is a malformed record, and the
// member is reported by name instead of being formatted into state.
func ecsDecodeStringMember(object map[string]interface{}, key string) (types.String, error) {
	raw, ok := object[key]
	if !ok {
		return types.StringValue(""), nil
	}
	return ecsDecodeStringValue(raw, key, false)
}

// ecsDecodeNumericStringMember is ecsDecodeStringMember for the members the API
// declares as a JSON number while the schema exposes them as a string. The
// client decodes the response with encoding/json, so the value arrives as a
// float64 and is rendered with strconv.FormatFloat(..., 'f', -1, 64): an
// integral value loses its trailing ".0" and precision is limited to the float64
// representation. A json.Number is still accepted for callers that supply the
// text form. Every other shape follows ecsDecodeStringMember.
func ecsDecodeNumericStringMember(object map[string]interface{}, key string) (types.String, error) {
	raw, ok := object[key]
	if !ok {
		return types.StringValue(""), nil
	}
	return ecsDecodeStringValue(raw, key, true)
}

// ecsDecodeStringValue reports a value that is not a string for the member it
// belongs to. Only a contract-declared numeric member (numeric=true) may arrive
// as a number, because those responses carry numbers for fields the schema
// exposes as strings.
func ecsDecodeStringValue(raw interface{}, name string, numeric bool) (types.String, error) {
	switch v := raw.(type) {
	case nil:
		return types.StringNull(), fmt.Errorf("%s is null instead of a string", name)
	case string:
		return types.StringValue(v), nil
	case json.Number:
		if numeric {
			return types.StringValue(v.String()), nil
		}
	case int:
		if numeric {
			return types.StringValue(strconv.FormatInt(int64(v), 10)), nil
		}
	case int32:
		if numeric {
			return types.StringValue(strconv.FormatInt(int64(v), 10)), nil
		}
	case int64:
		if numeric {
			return types.StringValue(strconv.FormatInt(v, 10)), nil
		}
	case float32:
		if numeric {
			return types.StringValue(strconv.FormatFloat(float64(v), 'f', -1, 32)), nil
		}
	case float64:
		if numeric {
			return types.StringValue(strconv.FormatFloat(v, 'f', -1, 64)), nil
		}
	}
	return types.StringNull(), fmt.Errorf("%s is %T instead of a string", name, raw)
}

// ecsParseInt64 converts the numeric shapes the ECS RPC returns. The boolean
// reports whether the value really was a whole number: a fractional value would
// otherwise be truncated silently, and a value of another type is not a number
// at all, so both are left to the caller to report.
func ecsParseInt64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case float32:
		f := float64(v)
		// math.MaxInt64 is not representable as a float64: the nearest value is
		// 2^63, which is already out of the int64 range.
		if f != math.Trunc(f) || f < math.MinInt64 || f >= math.MaxInt64 {
			return 0, false
		}
		return int64(f), true
	case float64:
		if v != math.Trunc(v) || v < math.MinInt64 || v >= math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	case json.Number:
		// The client decodes a JSON number as float64; json.Number is still
		// accepted for callers that supply the text form. A number written with
		// a zero fraction, such as 4.0, is still a whole number; a real fraction
		// is not, and must not be truncated.
		if n, err := v.Int64(); err == nil {
			return n, true
		}
		if f, err := v.Float64(); err == nil && f == math.Trunc(f) && f >= math.MinInt64 && f < math.MaxInt64 {
			return int64(f), true
		}
	case string:
		if n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			return n, true
		}
	}
	return 0, false
}

// ecsDecodeInt64Member reads a whole-number member of a record. A member the
// server did not report stays null, so state records that no value came back
// instead of a made-up zero; a member that is present but is not a whole number
// is a malformed record and is reported rather than truncated. The name is the
// member path used in the diagnostic, which for a nested object includes its
// parent.
func ecsDecodeInt64Member(object map[string]interface{}, key, name string) (types.Int64, error) {
	raw, ok := object[key]
	if !ok || raw == nil {
		return types.Int64Null(), nil
	}
	value, ok := ecsParseInt64(raw)
	if !ok {
		return types.Int64Null(), fmt.Errorf("%s is %T (%v) instead of a whole number", name, raw, raw)
	}
	return types.Int64Value(value), nil
}

// ecsDecodeBoolMember reads a boolean member of a record. A member the server did
// not report stays null instead of becoming false; a member that is present but
// is neither a boolean nor a boolean string is a malformed record.
func ecsDecodeBoolMember(object map[string]interface{}, key string) (types.Bool, error) {
	raw, ok := object[key]
	if !ok || raw == nil {
		return types.BoolNull(), nil
	}
	switch v := raw.(type) {
	case bool:
		return types.BoolValue(v), nil
	case string:
		if parsed, err := strconv.ParseBool(strings.TrimSpace(v)); err == nil {
			return types.BoolValue(parsed), nil
		}
	}
	return types.BoolNull(), fmt.Errorf("%s is %T (%v) instead of a boolean", key, raw, raw)
}

// ecsStringField declares one string member of a record and the attribute it
// populates. numeric marks the members the API returns as a JSON number while
// the schema exposes them as a string.
type ecsStringField struct {
	attribute string
	member    string
	numeric   bool
}

// ecsDecodeStringFields decodes the declared members into target. The first
// malformed member stops the record with its own path, so a partial attribute
// set is never returned.
func ecsDecodeStringFields(object map[string]interface{}, fields []ecsStringField, target map[string]attr.Value) error {
	for _, field := range fields {
		var (
			value types.String
			err   error
		)
		if field.numeric {
			value, err = ecsDecodeNumericStringMember(object, field.member)
		} else {
			value, err = ecsDecodeStringMember(object, field.member)
		}
		if err != nil {
			return err
		}
		target[field.attribute] = value
	}
	return nil
}

// ecsStringMapValue builds a string map attribute from decoded values.
func ecsStringMapValue(values map[string]string) types.Map {
	elements := make(map[string]attr.Value, len(values))
	for k, v := range values {
		elements[k] = types.StringValue(v)
	}
	return types.MapValueMust(types.StringType, elements)
}

// ecsTagMaxLength is the key and value length limit that the ECS tag filters
// document.
const ecsTagMaxLength = 128

// ecsValidateTagKey enforces the key limits the ECS tag filters document: the
// key must be present and must not use the reserved aliyun/acs: prefixes, and it
// must not carry a URL. The data sources differ only in what they accept as a
// tag value, so the key rule is shared.
func ecsValidateTagKey(key string) string {
	if key == "" {
		return "A tag key must not be empty."
	}
	if utf8.RuneCountInString(key) > ecsTagMaxLength {
		return fmt.Sprintf("A tag key must be at most %d characters.", ecsTagMaxLength)
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

// ecsSaeReservedTagKeyPattern keeps the historical SAE rule of the shared tag
// filter: a regular expression in which each dot matches any single character,
// so keys such as "sae-do-not-delete-x" are reserved as well.
var ecsSaeReservedTagKeyPattern = regexp.MustCompile(`^sae.do.not.delete`)

// ecsProviderReservedTagKey reports whether a tag key belongs to a namespace
// that the provider does not surface in state. It applies the provider-wide
// filter (four literal prefixes plus the historical SAE match) without calling
// tags.TagIgnored, which logs every key and value it inspects: this decoder runs
// on every tag of every record, and tag values are user data.
func ecsProviderReservedTagKey(key string) bool {
	for _, prefix := range []string{"aliyun", "acs:", "http://", "https://", "sae.do.not.delete"} {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return ecsSaeReservedTagKeyPattern.MatchString(key)
}

// ecsReservedInstancesTagKey keeps the historical tag rule of
// alicloud_reserved_instances: a key that starts with "sae" but not with the
// literal "sae.do.not.delete" prefix survives the shared filter, while every
// other key follows the provider-wide rule. Unifying this exception with
// ecsProviderReservedTagKey changes the tags that data source reports in state,
// so it needs its own decision rather than riding along with this change.
func ecsReservedInstancesTagKey(key string) bool {
	if strings.HasPrefix(key, "sae") && !strings.HasPrefix(key, "sae.do.not.delete") {
		return false
	}
	return ecsProviderReservedTagKey(key)
}

// ecsTagDecodeOptions selects the two tag rules the ECS data sources still
// differ on. reservedKey reports whether a key is hidden from state, and
// alternateMembers accepts `Key`/`Value` as alternatives to `TagKey`/`TagValue`.
type ecsTagDecodeOptions struct {
	reservedKey      func(string) bool
	alternateMembers bool
}

// ecsDecodeTags decodes the `Tags` member of an ECS record. A missing member is
// a valid empty map, because the RPC omits it when the record has no tags; an
// explicit null member or list is a malformed record, because the service never
// sends one. A member of the wrong shape, or an entry without a string key, is
// reported for the same reason instead of being dropped or panicking on a bare
// type assertion. The record itself is the argument so that a missing member
// stays distinguishable from an explicit null one.
func ecsDecodeTags(object map[string]interface{}, options ecsTagDecodeOptions) (types.Map, error) {
	values := map[string]string{}
	rawTags, ok := object["Tags"]
	if !ok {
		return ecsStringMapValue(values), nil
	}
	if rawTags == nil {
		return types.MapNull(types.StringType), fmt.Errorf("Tags is null instead of an object")
	}
	tagsMap, ok := rawTags.(map[string]interface{})
	if !ok {
		return types.MapNull(types.StringType), fmt.Errorf("Tags is %T instead of an object", rawTags)
	}
	rawList, ok := tagsMap["Tag"]
	if !ok {
		return ecsStringMapValue(values), nil
	}
	if rawList == nil {
		return types.MapNull(types.StringType), fmt.Errorf("Tags.Tag is null instead of a list")
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
		if !ok && options.alternateMembers {
			rawKey = item["Key"]
		}
		key, ok := rawKey.(string)
		if !ok || key == "" {
			return types.MapNull(types.StringType), fmt.Errorf("Tags.Tag[%d] has no string key", index)
		}
		rawValue, hasValue := item["TagValue"]
		if !hasValue && options.alternateMembers {
			rawValue, hasValue = item["Value"]
		}
		value := ""
		if hasValue {
			if rawValue == nil {
				return types.MapNull(types.StringType), fmt.Errorf("Tags.Tag[%d].TagValue is null instead of a string", index)
			}
			text, ok := rawValue.(string)
			if !ok {
				return types.MapNull(types.StringType), fmt.Errorf("Tags.Tag[%d].TagValue is %T instead of a string", index, rawValue)
			}
			value = text
		}
		if options.reservedKey != nil && options.reservedKey(key) {
			continue
		}
		values[key] = value
	}
	return ecsStringMapValue(values), nil
}

// ecsStringListValue builds a list attribute from decoded strings.
func ecsStringListValue(values []string) types.List {
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	return types.ListValueMust(types.StringType, elements)
}

// ecsElementsAsStrings reads the string elements of a configured list. A null
// or unknown list contributes no filter values.
func ecsElementsAsStrings(value types.List) []string {
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

// ecsSetRequestString writes a string filter only when the configuration set a
// value: a null, an unknown or an empty string asks for no filter.
func ecsSetRequestString(request map[string]interface{}, key string, value types.String) {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return
	}
	request[key] = value.ValueString()
}

// ecsApplyRepeatList writes a repeat-list filter in the 1-based `Parameter.N`
// form using the shared convert helper. Blank and repeated values are sent once,
// so a duplicated input does not change the request or make the result count
// depend on the server's duplicate handling.
func ecsApplyRepeatList(request map[string]interface{}, parameter string, values []string) {
	unique := ecsUniqueStrings(values)
	raw := make([]interface{}, 0, len(unique))
	for _, value := range unique {
		raw = append(raw, value)
	}
	expandArrayToMap(request, raw, parameter)
}

// ecsApplyTags maps the Terraform tag map onto the `Tag.N.Key/Value` filter
// using the shared tags helper. A null or unknown value is left out: an empty
// string is a value the caller asked to send, while a null is not. The keys are
// sorted so the same configuration always produces the same request.
func ecsApplyTags(request map[string]interface{}, tags types.Map) {
	if tags.IsNull() || tags.IsUnknown() {
		return
	}
	elements := tags.Elements()
	keys := make([]string, 0, len(elements))
	for key := range elements {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	list := make([]map[string]interface{}, 0, len(keys))
	for _, key := range keys {
		stringValue, ok := elements[key].(types.String)
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

// ecsNeedRetry reports whether a failed read is worth retrying inside the read
// budget. The context bound client reports throttling and decodable 5xx as tea
// SDK errors, a transport failure as a URL error, and a 5xx whose body cannot be
// decoded as UnreadableResponseError; the shared helper classifies the first
// two, and the status of the third is checked here because the helper does not
// know that type. A cancelled or expired read is never retried: its budget is
// gone.
func ecsNeedRetry(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	var unreadable *connectivity.UnreadableResponseError
	if errors.As(err, &unreadable) && unreadable.StatusCode >= 500 {
		return true
	}
	return helper.NeedRetry(err)
}

// ecsValidateQueryEnvelope rejects a 2xx response that is not the success
// envelope of the query that was sent. Decoding such a document would produce no
// records and overwrite the state with a zero result even though the service
// never answered the question.
//
// Every ECS success response carries the RequestId of the request it answered,
// so it is validated first, including for a response that carries the
// container: a document that merely happens to contain that key is not an
// answer to this query. The ECS RPC omits the container when nothing matched,
// and only that documented empty envelope is accepted without it: it still
// carries the total count the action documents. An explicit null container is
// not an empty result either. A present container is decoded by
// ecsDecodeRecordList, which reports a member or entry of the wrong shape.
//
// Whether a missing container may carry a positive count depends on the page
// the read asked for. At the first page a positive count without the container
// means the records of a non-empty result were dropped, so it is reported.
// A later page of a data source that accepts an arbitrary page_number may
// legitimately come back without the container once the page lies past the end
// of the result: the service total still describes the whole result, so the
// paginator, not this envelope check, decides whether every record was read.
// The caller passes fromFirstPage according to the page it starts reading
// from.
func ecsValidateQueryEnvelope(response map[string]interface{}, action, container, countKey string, fromFirstPage bool) error {
	if response == nil {
		return fmt.Errorf("%s returned an empty response body", action)
	}
	if _, hasCode := response["Code"]; hasCode {
		// An error document is not an empty result of the query, even when it
		// carries a RequestId.
		return fmt.Errorf("%s returned an error document instead of the %s member of its success envelope", action, container)
	}
	raw, ok := response["RequestId"]
	if !ok {
		return fmt.Errorf("%s returned a document without a RequestId; it is not this action's success envelope and is not treated as an empty result", action)
	}
	requestId, isString := raw.(string)
	if !isString || strings.TrimSpace(requestId) == "" {
		return fmt.Errorf("%s returned a document with %T as its RequestId instead of a non-empty string; it is not this action's success envelope and is not treated as an empty result", action, raw)
	}
	rawContainer, hasContainer := response[container]
	if hasContainer {
		if rawContainer == nil {
			return fmt.Errorf("%s returned %s as null instead of the documented object of its success envelope; the response is not treated as an empty result", action, container)
		}
		return nil
	}
	// The container is omitted when nothing matched, and the count member is
	// part of the documented success envelope of this action. Requiring a zero
	// count keeps a response whose records were dropped from being read as an
	// empty result at the first page.
	rawCount, ok := response[countKey]
	if !ok || rawCount == nil {
		return fmt.Errorf("%s returned neither the %s member nor %s, which its empty success envelope still carries; it is not this action's success envelope and is not treated as an empty result", action, container, countKey)
	}
	count, valid := ecsParseCount(rawCount)
	if !valid {
		return fmt.Errorf("%s returned %s as %T instead of a non-negative integer; it is not this action's success envelope and is not treated as an empty result", action, countKey, rawCount)
	}
	if count != 0 && fromFirstPage {
		return fmt.Errorf("%s returned %s %d without the %s member, so the records of a non-empty result were dropped; it is not treated as an empty result", action, countKey, count, container)
	}
	return nil
}

// ecsDecodeRecordList decodes a list of records nested under
// `container.member`. The service omits the container or the member when nothing
// matched, so an absent container or member is a valid empty result. An explicit
// null container or member is a malformed response and is reported instead of
// being read as an empty result. A member of the wrong shape, or an entry that
// is not an object, is reported as well, so a partial or empty result is never
// returned for a response that carried records.
func ecsDecodeRecordList(response map[string]interface{}, action, container, member string) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	if response == nil {
		return nil, fmt.Errorf("%s returned an empty response body", action)
	}
	rawContainer, ok := response[container]
	if !ok {
		return items, nil
	}
	if rawContainer == nil {
		return nil, fmt.Errorf("%s returned %s as null instead of an object; the response is not treated as an empty result", action, container)
	}
	containerMap, ok := rawContainer.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%s returned %s as %T instead of an object", action, container, rawContainer)
	}
	rawList, ok := containerMap[member]
	if !ok {
		return items, nil
	}
	if rawList == nil {
		return nil, fmt.Errorf("%s returned %s.%s as null instead of a list; the response is not treated as an empty result", action, container, member)
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

// ecsParseCount parses a response count. Only a non-negative whole number is
// a valid count: a fractional value would be silently truncated, and a negative
// value is never a count.
func ecsParseCount(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case json.Number:
		n, err := v.Int64()
		if err != nil || n < 0 {
			return 0, false
		}
		return n, true
	case int:
		if v < 0 {
			return 0, false
		}
		return int64(v), true
	case int32:
		if v < 0 {
			return 0, false
		}
		return int64(v), true
	case int64:
		if v < 0 {
			return 0, false
		}
		return v, true
	case float64:
		// math.MaxInt64 is not representable as a float64: the nearest value is
		// 2^63, which is already out of the int64 range.
		if v < 0 || v != math.Trunc(v) || v >= math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	case float32:
		f := float64(v)
		if f < 0 || f != math.Trunc(f) || f >= math.MaxInt64 {
			return 0, false
		}
		return int64(f), true
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil || n < 0 {
			return 0, false
		}
		return n, true
	}
	return 0, false
}

// ecsDecodeCount reads an integer response member. A missing member is
// reported as unknown so the caller can fall back; a member that is present but
// is not a non-negative whole number is a malformed response, not a count of 0.
func ecsDecodeCount(response map[string]interface{}, action, key string) (int64, bool, error) {
	raw, ok := response[key]
	if !ok || raw == nil {
		return 0, false, nil
	}
	value, ok := ecsParseCount(raw)
	if !ok {
		return 0, false, fmt.Errorf("%s returned %s as %T instead of a non-negative integer", action, key, raw)
	}
	return value, true, nil
}

// ecsOptionalObjectList decodes a single nested object of the record into a one
// element list. Only an absent field is an empty list; an explicit null is a
// malformed record, not an absent one, and is reported. The record itself is
// the argument so that a missing field stays distinguishable from an explicit
// null one, which a nil value alone cannot express. A field of the wrong type
// is reported as well. The build callback returns the decoded attributes with
// an error so a malformed nested member fails the record instead of being
// dropped.
func ecsOptionalObjectList(object map[string]interface{}, field string, objectType types.ObjectType, build func(object map[string]interface{}) (map[string]attr.Value, error)) (types.List, error) {
	value, ok := object[field]
	if !ok {
		return types.ListValueMust(objectType, []attr.Value{}), nil
	}
	if value == nil {
		return types.ListNull(objectType), fmt.Errorf("%s is null instead of an object", field)
	}
	nested, ok := value.(map[string]interface{})
	if !ok {
		return types.ListNull(objectType), fmt.Errorf("%s is %T instead of an object", field, value)
	}
	attributes, err := build(nested)
	if err != nil {
		return types.ListNull(objectType), fmt.Errorf("%s: %w", field, err)
	}
	return types.ListValueMust(objectType, []attr.Value{types.ObjectValueMust(objectType.AttrTypes, attributes)}), nil
}

// ecsOptionalObjectListMember decodes a nested list member of the record into a
// list of objects. Only an absent container or member is an empty list; an
// explicit null container or member is a malformed record, not an absent one,
// and is reported: reading it as an empty list would hide that the service
// answered with no usable content. The record itself is the argument so that a
// missing field stays distinguishable from an explicit null one. A container,
// member or entry of the wrong type is reported as well, so a partial result
// cannot be written to state.
func ecsOptionalObjectListMember(object map[string]interface{}, field, member string, objectType types.ObjectType, build func(item map[string]interface{}) (attr.Value, error)) (types.List, error) {
	value, ok := object[field]
	if !ok {
		return types.ListValueMust(objectType, []attr.Value{}), nil
	}
	if value == nil {
		return types.ListNull(objectType), fmt.Errorf("%s is null instead of an object", field)
	}
	container, ok := value.(map[string]interface{})
	if !ok {
		return types.ListNull(objectType), fmt.Errorf("%s is %T instead of an object", field, value)
	}
	raw, ok := container[member]
	if !ok {
		return types.ListValueMust(objectType, []attr.Value{}), nil
	}
	if raw == nil {
		return types.ListNull(objectType), fmt.Errorf("%s.%s is null instead of a list", field, member)
	}
	list, ok := raw.([]interface{})
	if !ok {
		return types.ListNull(objectType), fmt.Errorf("%s.%s is %T instead of a list", field, member, raw)
	}
	elements := make([]attr.Value, 0, len(list))
	for index, entry := range list {
		item, ok := entry.(map[string]interface{})
		if !ok {
			return types.ListNull(objectType), fmt.Errorf("%s.%s[%d] is %T instead of an object", field, member, index, entry)
		}
		element, err := build(item)
		if err != nil {
			return types.ListNull(objectType), fmt.Errorf("%s.%s[%d]: %w", field, member, index, err)
		}
		elements = append(elements, element)
	}
	return types.ListValueMust(objectType, elements), nil
}

// ecsOptionalStringList decodes a nested list of strings. Only an absent
// container or member is an empty list; an explicit null container or member
// is a malformed record and is reported instead of being read as an empty
// list. The record itself is the argument so that a missing field stays
// distinguishable from an explicit null one. A container, member or entry of
// the wrong type is reported as well.
func ecsOptionalStringList(object map[string]interface{}, field, member string) (types.List, error) {
	value, ok := object[field]
	if !ok {
		return types.ListValueMust(types.StringType, []attr.Value{}), nil
	}
	if value == nil {
		return types.ListNull(types.StringType), fmt.Errorf("%s is null instead of an object", field)
	}
	container, ok := value.(map[string]interface{})
	if !ok {
		return types.ListNull(types.StringType), fmt.Errorf("%s is %T instead of an object", field, value)
	}
	raw, ok := container[member]
	if !ok {
		return types.ListValueMust(types.StringType, []attr.Value{}), nil
	}
	if raw == nil {
		return types.ListNull(types.StringType), fmt.Errorf("%s.%s is null instead of a list", field, member)
	}
	list, ok := raw.([]interface{})
	if !ok {
		return types.ListNull(types.StringType), fmt.Errorf("%s.%s is %T instead of a list", field, member, raw)
	}
	elements := make([]attr.Value, 0, len(list))
	for index, entry := range list {
		text, ok := entry.(string)
		if !ok {
			return types.ListNull(types.StringType), fmt.Errorf("%s.%s[%d] is %T instead of a string", field, member, index, entry)
		}
		elements = append(elements, types.StringValue(text))
	}
	return types.ListValueMust(types.StringType, elements), nil
}

// ecsRecordIdentity returns the stable identity of a record. The identity must be
// a non-empty string: a value of another type cannot be deduplicated or
// addressed, so it is reported as a malformed result instead of being turned
// into an identity by formatting whatever the server sent.
func ecsRecordIdentity(object map[string]interface{}, action, key string) (string, error) {
	raw, ok := object[key]
	if !ok || raw == nil {
		return "", fmt.Errorf("%s returned a record without %s; the result was not returned because the record cannot be identified", action, key)
	}
	identity, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s returned %s as %T instead of a string; the result was not returned because the record cannot be identified", action, key, raw)
	}
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return "", fmt.Errorf("%s returned an empty %s; the result was not returned because the record cannot be identified", action, key)
	}
	return identity, nil
}

// ecsUniqueStrings drops blank and repeated values while preserving order, so
// a repeated filter value is sent once.
func ecsUniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	unique := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}
	return unique
}

// ecsNumberedPage is one page of a numbered query. totalKnown distinguishes an
// explicit TotalCount of 0 from a response that does not report the total.
type ecsNumberedPage struct {
	items      []map[string]interface{}
	total      int64
	totalKnown bool
}

// ecsCollectNumberedPages reads numbered pages until a short or empty page, the
// server total is reached, or the page budget is exhausted. The context carries
// the read deadline and cancellation signal, so the loop stops at the same point
// as the requests it issues. Integrity is enforced before anything is returned:
// every record must carry a unique, non-empty string identity, the server total
// must stay consistent, a short or empty page is only accepted when it accounts
// for every row that is readable from the start page, and a full page for the
// whole budget means the endpoint ignored the page parameters.
//
// A read that starts after the first page cannot see the rows of the earlier
// pages, so completeness is measured against the readable suffix of the total
// (TotalCount minus the rows the start page skips), not against the whole total.
func ecsCollectNumberedPages(ctx context.Context, action, identityKey string, maxPages int, startPage, pageSize int64, fetch func(pageNumber int64) (ecsNumberedPage, error)) ([]map[string]interface{}, int64, error) {
	// PageNumber and PageSize are int32 on the wire. Bounding the two inputs to
	// that range keeps `skipped` inside int64 and keeps every page number this
	// loop sends representable, instead of wrapping around after it is cast:
	// (MaxInt32-1)*MaxInt32 is still far below MaxInt64.
	if startPage < 1 || startPage > math.MaxInt32 {
		return nil, 0, fmt.Errorf("%s cannot start at page %d: the start page must be between 1 and %d", action, startPage, math.MaxInt32)
	}
	if pageSize < 1 || pageSize > math.MaxInt32 {
		return nil, 0, fmt.Errorf("%s cannot use PageSize %d: the page size must be between 1 and %d", action, pageSize, math.MaxInt32)
	}
	objects := make([]map[string]interface{}, 0)
	totalCount := int64(0)
	totalKnown := false
	skipped := (startPage - 1) * pageSize
	pageNumber := startPage
	seen := make(map[string]struct{})
	for page := 0; page < maxPages; page++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, 0, fmt.Errorf("%s read stopped before page %d was requested: %w; the result was not returned because it would be incomplete", action, pageNumber, ctxErr)
		}
		result, err := fetch(pageNumber)
		if err != nil {
			return nil, 0, err
		}
		if result.totalKnown {
			if totalKnown && result.total != totalCount {
				return nil, 0, fmt.Errorf("%s reported TotalCount %d on page %d after %d on an earlier page; the result was not returned because the server total is inconsistent", action, result.total, pageNumber, totalCount)
			}
			totalCount = result.total
			totalKnown = true
		}
		readableTotal := int64(0)
		if totalKnown && totalCount > skipped {
			readableTotal = totalCount - skipped
		}
		for _, item := range result.items {
			identity, err := ecsRecordIdentity(item, action, identityKey)
			if err != nil {
				return nil, 0, err
			}
			if _, duplicate := seen[identity]; duplicate {
				return nil, 0, fmt.Errorf("%s returned %s %q more than once; the result was not returned because it would contain duplicate rows", action, identityKey, identity)
			}
			seen[identity] = struct{}{}
		}
		objects = append(objects, result.items...)
		if totalKnown && int64(len(objects)) > readableTotal {
			return nil, 0, fmt.Errorf("%s returned %d rows from page %d while TotalCount %d leaves %d readable rows; the result was not returned because the response is inconsistent", action, len(objects), startPage, totalCount, readableTotal)
		}
		if len(result.items) == 0 {
			if totalKnown && int64(len(objects)) < readableTotal {
				return nil, 0, fmt.Errorf("%s returned an empty page at page %d while TotalCount %d still expects %d more readable rows; the result was not returned because it would be incomplete", action, pageNumber, totalCount, readableTotal-int64(len(objects)))
			}
			break
		}
		if int64(len(result.items)) < pageSize {
			if totalKnown && int64(len(objects)) < readableTotal {
				return nil, 0, fmt.Errorf("%s returned a short page at page %d while TotalCount %d still expects %d more readable rows; the result was not returned because it would be incomplete", action, pageNumber, totalCount, readableTotal-int64(len(objects)))
			}
			break
		}
		if totalKnown && int64(len(objects)) >= readableTotal {
			break
		}
		if pageNumber >= math.MaxInt32 {
			return nil, 0, fmt.Errorf("%s cannot continue past page %d, which is the largest page number the API accepts; the result was not returned because it would be incomplete", action, pageNumber)
		}
		pageNumber++
		if page == maxPages-1 {
			return nil, 0, fmt.Errorf("%s returned a full page for %d consecutive requests (PageSize=%d), so pagination did not terminate; the result was not returned because it would be incomplete", action, maxPages, pageSize)
		}
	}
	if !totalKnown && len(objects) > 0 {
		totalCount = int64(len(objects))
	}
	return objects, totalCount, nil
}
