package ecs

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	reservedInstancesTypeName   = "alicloud_reserved_instances"
	reservedInstancesAction     = "DescribeReservedInstances"
	reservedInstancesPageSize   = 100
	reservedInstancesMaxPages   = 1000
	reservedInstancesReadBudget = 30 * time.Minute
)

// reservedInstanceStatuses keeps `Expired` accepted because existing
// configurations may still send it, even though DescribeReservedInstances marks
// that filter value deprecated. The documentation states the deprecation.
var reservedInstanceStatuses = []string{
	"Creating",
	"Active",
	"Retired",
	"Updating",
	"Expired",
}

var (
	_ datasource.DataSource              = (*reservedInstancesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*reservedInstancesDataSource)(nil)
)

func NewReservedInstancesDataSource() datasource.DataSource {
	return &reservedInstancesDataSource{}
}

type reservedInstancesDataSource struct {
	fwadapt.DataSourceBase
}

type reservedInstancesModel struct {
	Id                   types.String `tfsdk:"id"`
	Ids                  types.List   `tfsdk:"ids"`
	ReservedInstanceName types.String `tfsdk:"reserved_instance_name"`
	ZoneId               types.String `tfsdk:"zone_id"`
	LockReason           types.String `tfsdk:"lock_reason"`
	InstanceType         types.String `tfsdk:"instance_type"`
	InstanceTypeFamily   types.String `tfsdk:"instance_type_family"`
	Scope                types.String `tfsdk:"scope"`
	OfferingType         types.String `tfsdk:"offering_type"`
	AllocationType       types.String `tfsdk:"allocation_type"`
	Status               types.List   `tfsdk:"status"`
	Tags                 types.Map    `tfsdk:"tags"`
	PageNumber           types.Int64  `tfsdk:"page_number"`
	PageSize             types.Int64  `tfsdk:"page_size"`
	TotalCount           types.Int64  `tfsdk:"total_count"`
	Names                types.List   `tfsdk:"names"`
	Instances            types.List   `tfsdk:"instances"`
}

var reservedInstanceOperationLockAttrTypes = map[string]attr.Type{
	"lock_reason": types.StringType,
}

var reservedInstanceOperationLockObjectType = types.ObjectType{AttrTypes: reservedInstanceOperationLockAttrTypes}

var reservedInstanceAttrTypes = map[string]attr.Type{
	"reserved_instance_id":   types.StringType,
	"reserved_instance_name": types.StringType,
	"status":                 types.StringType,
	"creation_time":          types.StringType,
	"start_time":             types.StringType,
	"expired_time":           types.StringType,
	"instance_type":          types.StringType,
	"instance_amount":        types.Int64Type,
	"region_id":              types.StringType,
	"zone_id":                types.StringType,
	"scope":                  types.StringType,
	"platform":               types.StringType,
	"offering_type":          types.StringType,
	"allocation_status":      types.StringType,
	"description":            types.StringType,
	"resource_group_id":      types.StringType,
	"operation_locks":        types.ListType{ElemType: reservedInstanceOperationLockObjectType},
	"tags":                   types.MapType{ElemType: types.StringType},
}

var reservedInstanceObjectType = types.ObjectType{AttrTypes: reservedInstanceAttrTypes}

func (d *reservedInstancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = reservedInstancesTypeName
}

func (d *reservedInstancesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of reserved instances (reserved instance coupons) in a region.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, derived from the matched reserved instance IDs.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of reserved instance IDs to filter by. The filter is sent to the server. A maximum of 100 entries is allowed. Blank IDs are dropped and repeated IDs are sent once. The attribute also reports the IDs that matched.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(100),
				},
			},
			"reserved_instance_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the reserved instance. The match is exact and is sent to the server.",
			},
			"zone_id": schema.StringAttribute{
				Optional:    true,
				Description: "The zone ID. It is required when `scope` is `Zone` and rejected when `scope` is `Region`, because a region-scoped reserved instance has no zone. When `zone_id` is set without `scope`, the data source requests `Scope = Zone`; when neither is set, the query is not restricted to a scope.",
			},
			"lock_reason": schema.StringAttribute{
				Optional:    true,
				Description: "The reason why the reserved instance is locked. Valid values: `financial` and `security`.",
				Validators: []validator.String{
					stringvalidator.OneOf("financial", "security"),
				},
			},
			"instance_type": schema.StringAttribute{
				Optional:    true,
				Description: "The instance type.",
			},
			"instance_type_family": schema.StringAttribute{
				Optional:    true,
				Description: "The instance type family.",
			},
			"scope": schema.StringAttribute{
				Optional:    true,
				Description: "The scope of the reserved instance. Valid values: `Region` and `Zone`. When set to `Zone`, `zone_id` is required; when set to `Region`, `zone_id` must not be set.",
				Validators: []validator.String{
					stringvalidator.OneOf("Region", "Zone"),
				},
			},
			"offering_type": schema.StringAttribute{
				Optional:    true,
				Description: "The payment type of the reserved instance. Valid values: `No Upfront`, `Partial Upfront` and `All Upfront`.",
				Validators: []validator.String{
					stringvalidator.OneOf("No Upfront", "Partial Upfront", "All Upfront"),
				},
			},
			"allocation_type": schema.StringAttribute{
				Optional:    true,
				Description: "The allocation type of the reserved instance. Valid values: `Normal` and `Shared`. Defaults to `Normal` on the server side.",
				Validators: []validator.String{
					stringvalidator.OneOf("Normal", "Shared"),
				},
			},
			"status": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A list of statuses to filter by. Valid values: `Creating`, `Active`, `Retired`, `Updating` and `Expired`. `Expired` is deprecated by `DescribeReservedInstances` and may be removed from the API; it is still accepted for compatibility, and new configurations should filter with the other values.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(100),
					listvalidator.ValueStringsAre(stringvalidator.OneOf(reservedInstanceStatuses...)),
				},
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A mapping of tags to filter by. Each entry is sent as a `Tag` filter. A maximum of 20 entries is allowed. A key and a value must be 1 to 128 characters, so an empty value is rejected; a key must not start with `aliyun` or `acs:`, a value must not start with `acs:`, and neither may contain a URL. A null value is rejected as well: omit the entry to leave the tag out of the request. A tag filter returns at most 1,000 reserved instances; when more may match, query `ListTagResources` instead.",
				Validators: []validator.Map{
					mapvalidator.SizeAtMost(20),
					reservedInstancesTagMapValidator{},
				},
			},
			"page_number": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of the first page to read. Valid values: 1 to 2147483647 (`PageNumber` is an int32 on the wire). Defaults to 1. The rows of the pages before `page_number` are not read and are not part of the result. When the API reports `TotalCount`, `total_count` still includes those skipped rows; an unknown-total fallback contains only the rows read. The data source continues with the following pages until the rows that are readable from `page_number` have been read, or the read budget is exhausted.",
				Validators: []validator.Int64{
					int64validator.Between(1, math.MaxInt32),
				},
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page. Valid values: 1 to 100. Defaults to 100. A page shorter than this value ends the read only when the rows read account for every row that is readable from `page_number`; otherwise the read is reported as incomplete.",
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The `TotalCount` reported by `DescribeReservedInstances` for the configured filters. When present, it counts every matching reserved instance in the region, including rows skipped by `page_number`. A total that changes between pages, or that disagrees with the rows read, is reported as an error instead of being written. When a valid `ReservedInstances` container is present but `TotalCount` is omitted or null, the total falls back to the number of rows read, which excludes skipped rows. A response with neither the container nor a usable `TotalCount` is rejected.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of reserved instance names that match the configured filters.",
			},
			"instances": schema.ListAttribute{
				Computed:    true,
				ElementType: reservedInstanceObjectType,
				Description: "A list of reserved instances that match the configured filters.",
			},
		},
	}
}

func (d *reservedInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// A configuration that still carries unknown values cannot be turned into a
	// complete request: every filter builder treats an unknown value as absent,
	// so the read would answer a wider question than the one asked and could
	// write that answer to state. Defer until the values are known when the
	// client supports it; otherwise stop before any request is sent.
	if !req.Config.Raw.IsFullyKnown() {
		if req.ClientCapabilities.DeferralAllowed {
			resp.Deferred = &datasource.Deferred{Reason: datasource.DeferredReasonDataSourceConfigUnknown}
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: the configuration is not known yet", reservedInstancesTypeName),
			"The configuration still contains unknown values, so the filters they define cannot be sent. Apply the configuration that defines them and read the data source again.",
		)
		return
	}

	var model reservedInstancesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope := ""
	if !model.Scope.IsNull() && !model.Scope.IsUnknown() {
		scope = model.Scope.ValueString()
	}
	zoneId := ""
	if !model.ZoneId.IsNull() && !model.ZoneId.IsUnknown() {
		zoneId = model.ZoneId.ValueString()
	}
	scope, err := reservedInstancesResolveScope(scope, zoneId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: invalid configuration", reservedInstancesTypeName),
			err.Error(),
		)
		return
	}

	pageNumber := int64(1)
	if !model.PageNumber.IsNull() && !model.PageNumber.IsUnknown() {
		pageNumber = model.PageNumber.ValueInt64()
	}
	pageSize := int64(reservedInstancesPageSize)
	if !model.PageSize.IsNull() && !model.PageSize.IsUnknown() {
		pageSize = model.PageSize.ValueInt64()
	}

	// One deadline covers every page, every retry and the final state write.
	ctx, cancel := context.WithTimeout(ctx, reservedInstancesReadBudget)
	defer cancel()

	model, err = d.load(ctx, d.Meta(), d.Meta().RegionId, scope, pageNumber, pageSize, model)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", reservedInstancesTypeName, reservedInstancesAction),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// load performs the whole read: it builds the request, follows every page,
// decodes the records and returns the state model. Read only adds the deadline,
// diagnostics and the state write, so the same request-to-state path can be
// exercised directly in unit tests.
func (d *reservedInstancesDataSource) load(ctx context.Context, client *connectivity.AliyunClient, regionId, scope string, pageNumber, pageSize int64, model reservedInstancesModel) (reservedInstancesModel, error) {
	request := reservedInstancesRequest(model, regionId, scope, pageNumber, pageSize)

	objects, totalCount, err := d.reservedInstancesDescribeAll(ctx, client, request, pageNumber, pageSize)
	if err != nil {
		return model, err
	}

	// A cancelled or expired read must not return a partial result to state.
	if ctxErr := ctx.Err(); ctxErr != nil {
		return model, fmt.Errorf("the read was cancelled or ran out of time before the result was written: %w", ctxErr)
	}

	return reservedInstancesResultModel(objects, totalCount, model)
}

// reservedInstancesResultModel turns the collected records into the state model.
// Records are ordered by identity, so the same result always produces the same
// lists and the same data source ID. The mapping is a pure function of the
// records and the server total, which lets the non-empty request-to-state path be
// asserted in unit tests without a live account.
func reservedInstancesResultModel(objects []map[string]interface{}, totalCount int64, model reservedInstancesModel) (reservedInstancesModel, error) {
	type matchedInstance struct {
		id     string
		name   string
		object map[string]interface{}
	}
	sorted := make([]matchedInstance, 0, len(objects))
	for _, object := range objects {
		id, err := reservedInstancesRecordIdentity(object)
		if err != nil {
			return model, err
		}
		name, err := reservedInstancesDecodeStringMember(object, "ReservedInstanceName")
		if err != nil {
			return model, fmt.Errorf("%s: %w", id, err)
		}
		sorted = append(sorted, matchedInstance{id: id, name: name.ValueString(), object: object})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].id < sorted[j].id
	})

	ids := make([]string, 0, len(sorted))
	names := make([]string, 0, len(sorted))
	elements := make([]attr.Value, 0, len(sorted))
	for _, instance := range sorted {
		ids = append(ids, instance.id)
		names = append(names, instance.name)
		element, err := reservedInstanceObject(instance.object)
		if err != nil {
			return model, err
		}
		elements = append(elements, element)
	}

	model.Id = types.StringValue(helper.DataResourceIdHash(ids))
	model.Ids = reservedInstancesStringListValue(ids)
	model.Names = reservedInstancesStringListValue(names)
	model.TotalCount = types.Int64Value(totalCount)
	model.Instances = types.ListValueMust(reservedInstanceObjectType, elements)
	return model, nil
}

// reservedInstancesRequest builds the full DescribeReservedInstances request
// from the configuration. It is kept separate from Read so the request-to-state
// mapping can be asserted directly.
func reservedInstancesRequest(model reservedInstancesModel, regionId, scope string, pageNumber, pageSize int64) map[string]interface{} {
	request := map[string]interface{}{
		"RegionId":   regionId,
		"PageNumber": int(pageNumber),
		"PageSize":   int(pageSize),
	}
	reservedInstancesSetRequestString(request, "ReservedInstanceName", model.ReservedInstanceName)
	reservedInstancesSetRequestString(request, "ZoneId", model.ZoneId)
	reservedInstancesSetRequestString(request, "LockReason", model.LockReason)
	reservedInstancesSetRequestString(request, "InstanceType", model.InstanceType)
	reservedInstancesSetRequestString(request, "InstanceTypeFamily", model.InstanceTypeFamily)
	reservedInstancesSetRequestString(request, "Scope", types.StringValue(scope))
	reservedInstancesSetRequestString(request, "OfferingType", model.OfferingType)
	reservedInstancesSetRequestString(request, "AllocationType", model.AllocationType)
	reservedInstancesApplyRepeatList(request, "ReservedInstanceId", reservedInstancesElementsAsStrings(model.Ids))
	reservedInstancesApplyRepeatList(request, "Status", reservedInstancesElementsAsStrings(model.Status))
	reservedInstancesApplyTags(request, model.Tags)
	return request
}

func (d *reservedInstancesDataSource) reservedInstancesDescribeAll(ctx context.Context, client *connectivity.AliyunClient, request map[string]interface{}, startPage, pageSize int64) ([]map[string]interface{}, int64, error) {
	return reservedInstancesCollectNumberedPages(ctx, reservedInstancesMaxPages, startPage, pageSize, func(pageNumber int64) (reservedInstancesNumberedPage, error) {
		request["PageNumber"] = int(pageNumber)
		request["PageSize"] = int(pageSize)
		response, err := reservedInstancesRpcPost(ctx, client, request)
		if err != nil {
			return reservedInstancesNumberedPage{}, err
		}
		if err := reservedInstancesValidateQueryEnvelope(response, pageNumber == 1); err != nil {
			return reservedInstancesNumberedPage{}, err
		}
		items, err := reservedInstancesDecodeRecordList(response)
		if err != nil {
			return reservedInstancesNumberedPage{}, err
		}
		totalCount, totalKnown, err := reservedInstancesDecodeCount(response)
		if err != nil {
			return reservedInstancesNumberedPage{}, err
		}
		return reservedInstancesNumberedPage{items: items, total: totalCount, totalKnown: totalKnown}, nil
	})
}

func reservedInstanceObject(object map[string]interface{}) (types.Object, error) {
	identity, err := reservedInstancesRecordIdentity(object)
	if err != nil {
		return types.ObjectNull(reservedInstanceAttrTypes), err
	}
	attributes, err := reservedInstanceAttributes(object)
	if err != nil {
		return types.ObjectNull(reservedInstanceAttrTypes), fmt.Errorf("%s: %w", identity, err)
	}
	return types.ObjectValueMust(reservedInstanceAttrTypes, attributes), nil
}

// reservedInstanceAttributes decodes every member of one reserved instance.
// Every contract string is decoded strictly by name, so a null or a container
// fails the record with its path instead of being written to state.
func reservedInstanceAttributes(object map[string]interface{}) (map[string]attr.Value, error) {
	attributes := make(map[string]attr.Value, 19)
	if err := reservedInstancesDecodeStringFields(object, []reservedInstancesStringField{
		{"reserved_instance_id", "ReservedInstanceId"},
		{"reserved_instance_name", "ReservedInstanceName"},
		{"status", "Status"},
		{"creation_time", "CreationTime"},
		{"start_time", "StartTime"},
		{"expired_time", "ExpiredTime"},
		{"instance_type", "InstanceType"},
		{"region_id", "RegionId"},
		{"zone_id", "ZoneId"},
		{"scope", "Scope"},
		{"platform", "Platform"},
		{"offering_type", "OfferingType"},
		{"allocation_status", "AllocationStatus"},
		{"description", "Description"},
		{"resource_group_id", "ResourceGroupId"},
	}, attributes); err != nil {
		return nil, err
	}

	tags, err := reservedInstancesDecodeTags(object)
	if err != nil {
		return nil, err
	}
	attributes["tags"] = tags

	operationLocks, err := reservedInstancesDecodeOperationLocks(object)
	if err != nil {
		return nil, err
	}
	attributes["operation_locks"] = operationLocks

	instanceAmount, err := reservedInstancesDecodeInt64Member(object, "InstanceAmount")
	if err != nil {
		return nil, err
	}
	attributes["instance_amount"] = instanceAmount

	return attributes, nil
}

const reservedInstancesRetryTimeout = 5 * time.Minute

// reservedInstancesRetryBackoff is the first wait between two attempts of one
// read; it grows by the same step up to 30 seconds.
var reservedInstancesRetryBackoff = 3 * time.Second

const reservedInstancesTagMaxLength = 128

// Conversions, request helpers and remote-call plumbing used by
// alicloud_reserved_instances. They keep response decoding tolerant of the two
// shapes the ECS RPC returns for the same field (JSON number or decimal string)
// without panicking on a missing key, and keep the query, pagination and retry
// behavior of the data source in one place.
//
// The names are private and carry this data source's prefix so they cannot
// collide with helpers of other data sources in this package.

// reservedInstancesDecodeStringMember reads a string member of a record. The ECS
// RPC omits an optional string member that carries no value, so an absent member
// stays an empty string. A member that is present but is null, an object or an
// array is not a string this contract can return: it is a malformed record, and
// the member is reported by name instead of being formatted into state.
func reservedInstancesDecodeStringMember(object map[string]interface{}, key string) (types.String, error) {
	raw, ok := object[key]
	if !ok {
		return types.StringValue(""), nil
	}
	if text, ok := raw.(string); ok {
		return types.StringValue(text), nil
	}
	if raw == nil {
		return types.StringNull(), fmt.Errorf("%s is null instead of a string", key)
	}
	return types.StringNull(), fmt.Errorf("%s is %T instead of a string", key, raw)
}

// reservedInstancesParseInt64 converts the numeric shapes the RPC decoder can
// produce: a JSON number arrives as float64, while a decimal string remains a
// string. The boolean reports whether the value really was a whole number: a
// fractional value would otherwise be truncated silently, and a value of another
// type is not a number at all, so both are left to the caller to report.
func reservedInstancesParseInt64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case float64:
		if v != math.Trunc(v) || v < math.MinInt64 || v >= math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	case string:
		if n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			return n, true
		}
	}
	return 0, false
}

// reservedInstancesDecodeInt64Member reads a whole-number member of a record. A
// member the server did not report stays null, so state records that no value
// came back instead of a made-up zero; a member that is present but is not a
// whole number is a malformed record and is reported rather than truncated.
func reservedInstancesDecodeInt64Member(object map[string]interface{}, key string) (types.Int64, error) {
	raw, ok := object[key]
	if !ok || raw == nil {
		return types.Int64Null(), nil
	}
	value, ok := reservedInstancesParseInt64(raw)
	if !ok {
		return types.Int64Null(), fmt.Errorf("%s is %T (%v) instead of a whole number", key, raw, raw)
	}
	return types.Int64Value(value), nil
}

// reservedInstancesStringField declares one string member of a record and the
// attribute it populates.
type reservedInstancesStringField struct {
	attribute string
	member    string
}

// reservedInstancesDecodeStringFields decodes the declared members into target.
// The first malformed member stops the record with its own path, so a partial
// attribute set is never returned.
func reservedInstancesDecodeStringFields(object map[string]interface{}, fields []reservedInstancesStringField, target map[string]attr.Value) error {
	for _, field := range fields {
		value, err := reservedInstancesDecodeStringMember(object, field.member)
		if err != nil {
			return err
		}
		target[field.attribute] = value
	}
	return nil
}

// reservedInstancesStringMapValue builds a string map attribute from decoded
// values.
func reservedInstancesStringMapValue(values map[string]string) types.Map {
	elements := make(map[string]attr.Value, len(values))
	for k, v := range values {
		elements[k] = types.StringValue(v)
	}
	return types.MapValueMust(types.StringType, elements)
}

// reservedInstancesDecodeTags decodes the `Tags` member of an ECS record. A
// missing member is a valid empty map, because the RPC omits it when the record
// has no tags; an explicit null member or list is a malformed record, because
// the service never sends one. A member of the wrong shape, or an entry without
// a string key, is reported for the same reason instead of being dropped or
// panicking on a bare type assertion. The record itself is the argument so that
// a missing member stays distinguishable from an explicit null one.
func reservedInstancesDecodeTags(object map[string]interface{}) (types.Map, error) {
	values := map[string]string{}
	rawTags, ok := object["Tags"]
	if !ok {
		return reservedInstancesStringMapValue(values), nil
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
		return reservedInstancesStringMapValue(values), nil
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
		key, ok := rawKey.(string)
		if !ok || key == "" {
			return types.MapNull(types.StringType), fmt.Errorf("Tags.Tag[%d] has no string key", index)
		}
		rawValue, hasValue := item["TagValue"]
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
		if strings.HasPrefix(key, "sae.do.not.delete") || (!strings.HasPrefix(key, "sae") && tagIgnored(key, value)) {
			continue
		}
		values[key] = value
	}
	return reservedInstancesStringMapValue(values), nil
}

func reservedInstancesStringListValue(values []string) types.List {
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	return types.ListValueMust(types.StringType, elements)
}

func reservedInstancesElementsAsStrings(value types.List) []string {
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

func reservedInstancesSetRequestString(request map[string]interface{}, key string, value types.String) {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return
	}
	request[key] = value.ValueString()
}

// reservedInstancesResolveScope applies the documented rules of the `scope` and
// `zone_id` pair and returns the scope that is sent to the API. `Zone` and
// `Region` cannot both be expressed: a conflicting pair is rejected instead of
// being sent, because the API would answer it with a silently wrong filter (a
// region-scoped reserved instance has no zone, so the conjunction matches
// nothing, while a server that ignores ZoneId returns every region-scoped
// instance). `zone_id` alone keeps the documented compatibility behaviour and
// means `Zone`.
func reservedInstancesResolveScope(scope, zoneId string) (string, error) {
	switch {
	case scope == "Zone" && zoneId == "":
		return "", fmt.Errorf("`zone_id` must be set when `scope` is `Zone`.")
	case scope == "Region" && zoneId != "":
		return "", fmt.Errorf("`zone_id` cannot be set when `scope` is `Region`: a region-scoped reserved instance has no zone. Remove `zone_id` or set `scope` to `Zone`.")
	case scope == "" && zoneId != "":
		return "Zone", nil
	}
	return scope, nil
}

// reservedInstancesApplyRepeatList writes a repeat-list filter in the 1-based
// `Parameter.N` form using the shared convert helper. Blank and repeated values
// are sent once, so a duplicated input does not change the request or make the
// result count depend on the server's duplicate handling.
func reservedInstancesApplyRepeatList(request map[string]interface{}, parameter string, values []string) {
	unique := reservedInstancesUniqueStrings(values)
	raw := make([]interface{}, 0, len(unique))
	for _, value := range unique {
		raw = append(raw, value)
	}
	expandArrayToMap(request, raw, parameter)
}

// reservedInstancesApplyTags maps the Terraform tag map onto the
// `Tag.N.Key/Value` filter using the shared tags helper. A null or unknown value
// is left out: an empty string is a value the caller asked to send, while a null
// is not.
func reservedInstancesApplyTags(request map[string]interface{}, tags types.Map) {
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

// reservedInstancesTagMapValidator enforces the key and value limits that the
// ECS tag filters document: the key must be present and must not use the
// reserved aliyun/acs: prefixes, the value length is capped, and neither side
// may carry a URL. This API rejects an empty value, so every empty value is
// reported.
type reservedInstancesTagMapValidator struct{}

func (v reservedInstancesTagMapValidator) Description(context.Context) string {
	return "Each tag key and value must satisfy the tag limits documented for the API."
}

func (v reservedInstancesTagMapValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateMap checks the key of every entry even when its value is null or
// unknown. A null value is rejected rather than serialized as an empty string:
// the two are different requests, and only an explicit empty string asks the
// API for an empty value. Every invalid entry is reported, each on its own map
// key, so one plan names all the entries to fix instead of the first one only.
func (v reservedInstancesTagMapValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for key, element := range req.ConfigValue.Elements() {
		path := req.Path.AtMapKey(key)
		if message := validateReservedInstancesTagKey(key); message != "" {
			resp.Diagnostics.AddAttributeError(path, "Invalid Tag Key", message)
			continue
		}
		value, ok := element.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(path, "Invalid Tag Value", fmt.Sprintf("The value of tag key %q must be a string.", key))
			continue
		}
		if value.IsUnknown() {
			// Not known yet: it cannot be validated here, and the request
			// builder leaves it out instead of sending an empty value.
			continue
		}
		if value.IsNull() {
			resp.Diagnostics.AddAttributeError(path, "Invalid Tag Value", fmt.Sprintf("The value of tag key %q must be set; omit the entry to leave the tag out of the request, because this API rejects an empty value.", key))
			continue
		}
		if message := validateReservedInstancesTagValue(value.ValueString()); message != "" {
			resp.Diagnostics.AddAttributeError(path, "Invalid Tag Value", message)
		}
	}
}

func validateReservedInstancesTagKey(key string) string {
	if key == "" {
		return "A tag key must not be empty."
	}
	if utf8.RuneCountInString(key) > reservedInstancesTagMaxLength {
		return fmt.Sprintf("A tag key must be at most %d characters.", reservedInstancesTagMaxLength)
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

func validateReservedInstancesTagValue(value string) string {
	if value == "" {
		return "A tag value must not be empty."
	}
	if utf8.RuneCountInString(value) > reservedInstancesTagMaxLength {
		return fmt.Sprintf("A tag value must be at most %d characters.", reservedInstancesTagMaxLength)
	}
	if strings.HasPrefix(value, "acs:") {
		return "A tag value must not start with \"acs:\"."
	}
	if strings.Contains(value, "http://") || strings.Contains(value, "https://") {
		return "A tag value must not contain a URL."
	}
	return ""
}

// reservedInstancesSleepBackoff waits for the retry backoff and returns as soon
// as the read context is cancelled, so a cancelled read stops instead of
// sleeping through the remaining backoff.
func reservedInstancesSleepBackoff(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// reservedInstancesNeedRetry reports whether a failed read is worth retrying
// inside the read budget. The context bound SDK path reports throttling and 5xx
// as tea SDK errors, a transport failure as a URL error, and a 5xx whose body
// cannot be decoded as UnreadableResponseError; the shared helper classifies the
// first two, and the status of the third is checked here because the helper does
// not know that type. A cancelled or expired read is never retried: its budget
// is gone.
func reservedInstancesNeedRetry(err error) bool {
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

// reservedInstancesRpcClient is the context-bound RPC capability used by the
// read loop. Keeping the boundary explicit lets its error classification be
// regression-tested without real credentials or timing-dependent transport.
type reservedInstancesRpcClient interface {
	RpcPostWithContext(ctx context.Context, apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}) (map[string]interface{}, error)
}

// reservedInstancesRpcPost sends one read-only ECS query with the shared retry
// policy: transient failures are retried inside the bounded window, everything
// else fails on the first response. The context carries the read budget: it
// gates every attempt and the backoff, and the request itself is built from it,
// so an in-flight call is aborted instead of running on after the read returns.
// RpcPostWithContext issues exactly one transport attempt per call, so the loop
// below owns every retry and runs inline: no retry callback can write a
// response after the caller has stopped waiting.
func reservedInstancesRpcPost(ctx context.Context, client reservedInstancesRpcClient, request map[string]interface{}) (map[string]interface{}, error) {
	backoff := reservedInstancesRetryBackoff
	retryDeadline := time.Now().Add(reservedInstancesRetryTimeout)
	for {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("%s was not sent: %w", reservedInstancesAction, ctxErr)
		}
		response, rpcErr := client.RpcPostWithContext(ctx, "Ecs", "2014-05-26", reservedInstancesAction, nil, request)
		if rpcErr == nil {
			logs.AddDebug(reservedInstancesAction, response, request)
			return response, nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("%s was cancelled while the request was in flight: %w", reservedInstancesAction, ctxErr)
		}
		if errors.Is(rpcErr, context.Canceled) || errors.Is(rpcErr, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%s reached its own cancellation or request deadline before the read context ended: %w", reservedInstancesAction, rpcErr)
		}
		if !reservedInstancesNeedRetry(rpcErr) || !time.Now().Add(backoff).Before(retryDeadline) {
			return nil, errs.WrapErrorf(rpcErr, errs.DataDefaultErrorMsg, reservedInstancesTypeName, reservedInstancesAction, errs.AlibabaCloudSdkGoERROR)
		}
		if waitErr := reservedInstancesSleepBackoff(ctx, backoff); waitErr != nil {
			return nil, fmt.Errorf("%s was interrupted while waiting to retry: %w", reservedInstancesAction, waitErr)
		}
		if backoff < 30*time.Second {
			backoff += 3 * time.Second
		}
	}
}

// reservedInstancesValidateQueryEnvelope rejects a 2xx response that is not the
// success envelope of the query that was sent. Decoding such a document would
// produce no records and overwrite the state with a zero result even though the
// service never answered the question.
//
// Every ECS success response carries the RequestId of the request it answered,
// so it is validated first, including for a response that carries the container:
// a document that merely happens to contain that key is not an answer to this
// query. The ECS RPC omits the container when nothing matched, and only that
// documented empty envelope is accepted without it: it still carries the total
// count the action documents. An explicit null container is not an empty result
// either. A present container is decoded by reservedInstancesDecodeRecordList,
// which reports a member or entry of the wrong shape.
//
// Whether a missing container may carry a positive count depends on the page the
// read asked for. At the first page a positive count without the container means
// the records of a non-empty result were dropped, so it is reported. A later page
// of a data source that accepts an arbitrary page_number may legitimately come
// back without the container once the page lies past the end of the result: the
// service total still describes the whole result, so the paginator, not this
// envelope check, decides whether every record was read.
func reservedInstancesValidateQueryEnvelope(response map[string]interface{}, fromFirstPage bool) error {
	if response == nil {
		return fmt.Errorf("%s returned an empty response body", reservedInstancesAction)
	}
	if _, hasCode := response["Code"]; hasCode {
		// An error document is not an empty result of the query, even when it
		// carries a RequestId.
		return fmt.Errorf("%s returned an error document instead of the ReservedInstances member of its success envelope", reservedInstancesAction)
	}
	raw, ok := response["RequestId"]
	if !ok {
		return fmt.Errorf("%s returned a document without a RequestId; it is not this action's success envelope and is not treated as an empty result", reservedInstancesAction)
	}
	requestId, isString := raw.(string)
	if !isString || strings.TrimSpace(requestId) == "" {
		return fmt.Errorf("%s returned a document with %T as its RequestId instead of a non-empty string; it is not this action's success envelope and is not treated as an empty result", reservedInstancesAction, raw)
	}
	rawContainer, hasContainer := response["ReservedInstances"]
	if hasContainer {
		if rawContainer == nil {
			return fmt.Errorf("%s returned ReservedInstances as null instead of the documented object of its success envelope; the response is not treated as an empty result", reservedInstancesAction)
		}
		return nil
	}
	// The container is omitted when nothing matched, and the count member is
	// part of the documented success envelope of this action. Requiring a zero
	// count keeps a response whose records were dropped from being read as an
	// empty result at the first page.
	rawCount, ok := response["TotalCount"]
	if !ok || rawCount == nil {
		return fmt.Errorf("%s returned neither the ReservedInstances member nor TotalCount, which its empty success envelope still carries; it is not this action's success envelope and is not treated as an empty result", reservedInstancesAction)
	}
	count, valid := reservedInstancesParseCount(rawCount)
	if !valid {
		return fmt.Errorf("%s returned TotalCount as %T instead of a non-negative integer; it is not this action's success envelope and is not treated as an empty result", reservedInstancesAction, rawCount)
	}
	if count != 0 && fromFirstPage {
		return fmt.Errorf("%s returned TotalCount %d without the ReservedInstances member, so the records of a non-empty result were dropped; it is not treated as an empty result", reservedInstancesAction, count)
	}
	return nil
}

// reservedInstancesDecodeRecordList decodes a list of records nested under
// `container.member`. The service omits the container or the member when nothing
// matched, so an absent container or member is a valid empty result. An explicit
// null container or member is a malformed response and is reported instead of
// being read as an empty result. A member of the wrong shape, or an entry that is
// not an object, is reported as well, so a partial or empty result is never
// returned for a response that carried records.
func reservedInstancesDecodeRecordList(response map[string]interface{}) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	if response == nil {
		return nil, fmt.Errorf("%s returned an empty response body", reservedInstancesAction)
	}
	rawContainer, ok := response["ReservedInstances"]
	if !ok {
		return items, nil
	}
	if rawContainer == nil {
		return nil, fmt.Errorf("%s returned ReservedInstances as null instead of an object; the response is not treated as an empty result", reservedInstancesAction)
	}
	containerMap, ok := rawContainer.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%s returned ReservedInstances as %T instead of an object", reservedInstancesAction, rawContainer)
	}
	rawList, ok := containerMap["ReservedInstance"]
	if !ok {
		return items, nil
	}
	if rawList == nil {
		return nil, fmt.Errorf("%s returned ReservedInstances.ReservedInstance as null instead of a list; the response is not treated as an empty result", reservedInstancesAction)
	}
	list, ok := rawList.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s returned ReservedInstances.ReservedInstance as %T instead of a list", reservedInstancesAction, rawList)
	}
	for index, entry := range list {
		item, ok := entry.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s returned ReservedInstances.ReservedInstance[%d] as %T instead of an object", reservedInstancesAction, index, entry)
		}
		items = append(items, item)
	}
	return items, nil
}

// reservedInstancesParseCount parses a response count. Only a non-negative whole
// number is a valid count: a fractional value would be silently truncated, and a
// negative value is never a count.
func reservedInstancesParseCount(value interface{}) (int64, bool) {
	count, ok := reservedInstancesParseInt64(value)
	if !ok || count < 0 {
		return 0, false
	}
	return count, true
}

// reservedInstancesDecodeCount reads an integer response member. A missing
// member is reported as unknown so the caller can fall back; a member that is
// present but is not a non-negative whole number is a malformed response, not a
// count of 0.
func reservedInstancesDecodeCount(response map[string]interface{}) (int64, bool, error) {
	raw, ok := response["TotalCount"]
	if !ok || raw == nil {
		return 0, false, nil
	}
	value, ok := reservedInstancesParseCount(raw)
	if !ok {
		return 0, false, fmt.Errorf("%s returned TotalCount as %T instead of a non-negative integer", reservedInstancesAction, raw)
	}
	return value, true, nil
}

// reservedInstancesDecodeOperationLocks preserves the nested-list response
// contract while decoding the only nested object used by this data source.
func reservedInstancesDecodeOperationLocks(object map[string]interface{}) (types.List, error) {
	value, ok := object["OperationLocks"]
	if !ok {
		return types.ListValueMust(reservedInstanceOperationLockObjectType, []attr.Value{}), nil
	}
	if value == nil {
		return types.ListNull(reservedInstanceOperationLockObjectType), fmt.Errorf("OperationLocks is null instead of an object")
	}
	container, ok := value.(map[string]interface{})
	if !ok {
		return types.ListNull(reservedInstanceOperationLockObjectType), fmt.Errorf("OperationLocks is %T instead of an object", value)
	}
	raw, ok := container["OperationLock"]
	if !ok {
		return types.ListValueMust(reservedInstanceOperationLockObjectType, []attr.Value{}), nil
	}
	if raw == nil {
		return types.ListNull(reservedInstanceOperationLockObjectType), fmt.Errorf("OperationLocks.OperationLock is null instead of a list")
	}
	list, ok := raw.([]interface{})
	if !ok {
		return types.ListNull(reservedInstanceOperationLockObjectType), fmt.Errorf("OperationLocks.OperationLock is %T instead of a list", raw)
	}
	elements := make([]attr.Value, 0, len(list))
	for index, entry := range list {
		item, ok := entry.(map[string]interface{})
		if !ok {
			return types.ListNull(reservedInstanceOperationLockObjectType), fmt.Errorf("OperationLocks.OperationLock[%d] is %T instead of an object", index, entry)
		}
		lockReason, err := reservedInstancesDecodeStringMember(item, "LockReason")
		if err != nil {
			return types.ListNull(reservedInstanceOperationLockObjectType), fmt.Errorf("OperationLocks.OperationLock[%d]: %w", index, err)
		}
		elements = append(elements, types.ObjectValueMust(reservedInstanceOperationLockAttrTypes, map[string]attr.Value{
			"lock_reason": lockReason,
		}))
	}
	return types.ListValueMust(reservedInstanceOperationLockObjectType, elements), nil
}

// reservedInstancesRecordIdentity returns the stable identity of a record. The
// identity must be a non-empty string: a value of another type cannot be
// deduplicated or addressed, so it is reported as a malformed result instead of
// being turned into an identity by formatting whatever the server sent.
func reservedInstancesRecordIdentity(object map[string]interface{}) (string, error) {
	raw, ok := object["ReservedInstanceId"]
	if !ok || raw == nil {
		return "", fmt.Errorf("%s returned a record without ReservedInstanceId; the result was not returned because the record cannot be identified", reservedInstancesAction)
	}
	identity, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s returned ReservedInstanceId as %T instead of a string; the result was not returned because the record cannot be identified", reservedInstancesAction, raw)
	}
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return "", fmt.Errorf("%s returned an empty ReservedInstanceId; the result was not returned because the record cannot be identified", reservedInstancesAction)
	}
	return identity, nil
}

// reservedInstancesUniqueStrings drops blank and repeated values while
// preserving order, so a repeated filter value is sent once.
func reservedInstancesUniqueStrings(values []string) []string {
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

// reservedInstancesNumberedPage is one page of a numbered query. totalKnown
// distinguishes an explicit TotalCount of 0 from a response that does not report
// the total.
type reservedInstancesNumberedPage struct {
	items      []map[string]interface{}
	total      int64
	totalKnown bool
}

// reservedInstancesCollectNumberedPages reads numbered pages until a short or
// empty page, the server total is reached, or the page budget is exhausted. The
// context carries the read deadline and cancellation signal, so the loop stops
// at the same point as the requests it issues. Integrity is enforced before
// anything is returned: every record must carry a unique, non-empty string
// identity, the server total must stay consistent, a short or empty page is only
// accepted when it accounts for every row that is readable from the start page,
// and a full page for the whole budget means the endpoint ignored the page
// parameters.
//
// A read that starts after the first page cannot see the rows of the earlier
// pages, so completeness is measured against the readable suffix of the total
// (TotalCount minus the rows the start page skips), not against the whole total.
func reservedInstancesCollectNumberedPages(ctx context.Context, maxPages int, startPage, pageSize int64, fetch func(pageNumber int64) (reservedInstancesNumberedPage, error)) ([]map[string]interface{}, int64, error) {
	// PageNumber and PageSize are int32 on the wire. Bounding the two inputs to
	// that range keeps `skipped` inside int64 and keeps every page number this
	// loop sends representable, instead of wrapping around after it is cast:
	// (MaxInt32-1)*MaxInt32 is still far below MaxInt64.
	if startPage < 1 || startPage > math.MaxInt32 {
		return nil, 0, fmt.Errorf("%s cannot start at page %d: the start page must be between 1 and %d", reservedInstancesAction, startPage, math.MaxInt32)
	}
	if pageSize < 1 || pageSize > math.MaxInt32 {
		return nil, 0, fmt.Errorf("%s cannot use PageSize %d: the page size must be between 1 and %d", reservedInstancesAction, pageSize, math.MaxInt32)
	}
	objects := make([]map[string]interface{}, 0)
	totalCount := int64(0)
	totalKnown := false
	skipped := (startPage - 1) * pageSize
	pageNumber := startPage
	seen := make(map[string]struct{})
	for page := 0; page < maxPages; page++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, 0, fmt.Errorf("%s read stopped before page %d was requested: %w; the result was not returned because it would be incomplete", reservedInstancesAction, pageNumber, ctxErr)
		}
		result, err := fetch(pageNumber)
		if err != nil {
			return nil, 0, err
		}
		if result.totalKnown {
			if totalKnown && result.total != totalCount {
				return nil, 0, fmt.Errorf("%s reported TotalCount %d on page %d after %d on an earlier page; the result was not returned because the server total is inconsistent", reservedInstancesAction, result.total, pageNumber, totalCount)
			}
			totalCount = result.total
			totalKnown = true
		}
		readableTotal := int64(0)
		if totalKnown && totalCount > skipped {
			readableTotal = totalCount - skipped
		}
		for _, item := range result.items {
			identity, err := reservedInstancesRecordIdentity(item)
			if err != nil {
				return nil, 0, err
			}
			if _, duplicate := seen[identity]; duplicate {
				return nil, 0, fmt.Errorf("%s returned ReservedInstanceId %q more than once; the result was not returned because it would contain duplicate rows", reservedInstancesAction, identity)
			}
			seen[identity] = struct{}{}
		}
		objects = append(objects, result.items...)
		if totalKnown && int64(len(objects)) > readableTotal {
			return nil, 0, fmt.Errorf("%s returned %d rows from page %d while TotalCount %d leaves %d readable rows; the result was not returned because the response is inconsistent", reservedInstancesAction, len(objects), startPage, totalCount, readableTotal)
		}
		if len(result.items) == 0 {
			if totalKnown && int64(len(objects)) < readableTotal {
				return nil, 0, fmt.Errorf("%s returned an empty page at page %d while TotalCount %d still expects %d more readable rows; the result was not returned because it would be incomplete", reservedInstancesAction, pageNumber, totalCount, readableTotal-int64(len(objects)))
			}
			break
		}
		if int64(len(result.items)) < pageSize {
			if totalKnown && int64(len(objects)) < readableTotal {
				return nil, 0, fmt.Errorf("%s returned a short page at page %d while TotalCount %d still expects %d more readable rows; the result was not returned because it would be incomplete", reservedInstancesAction, pageNumber, totalCount, readableTotal-int64(len(objects)))
			}
			break
		}
		if totalKnown && int64(len(objects)) >= readableTotal {
			break
		}
		if pageNumber >= math.MaxInt32 {
			return nil, 0, fmt.Errorf("%s cannot continue past page %d, which is the largest page number the API accepts; the result was not returned because it would be incomplete", reservedInstancesAction, pageNumber)
		}
		pageNumber++
		if page == maxPages-1 {
			return nil, 0, fmt.Errorf("%s returned a full page for %d consecutive requests (PageSize=%d), so pagination did not terminate; the result was not returned because it would be incomplete", reservedInstancesAction, maxPages, pageSize)
		}
	}
	if !totalKnown && len(objects) > 0 {
		totalCount = int64(len(objects))
	}
	return objects, totalCount, nil
}
