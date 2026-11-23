package ecs

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
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
// diagnostics and the state write.
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
// records and the server total.
func reservedInstancesResultModel(objects []map[string]interface{}, totalCount int64, model reservedInstancesModel) (reservedInstancesModel, error) {
	type matchedInstance struct {
		id     string
		name   string
		object map[string]interface{}
	}
	sorted := make([]matchedInstance, 0, len(objects))
	for _, object := range objects {
		id, err := ecsRecordIdentity(object, reservedInstancesAction, "ReservedInstanceId")
		if err != nil {
			return model, err
		}
		name, err := ecsDecodeStringMember(object, "ReservedInstanceName")
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
	model.Ids = ecsStringListValue(ids)
	model.Names = ecsStringListValue(names)
	model.TotalCount = types.Int64Value(totalCount)
	model.Instances = types.ListValueMust(reservedInstanceObjectType, elements)
	return model, nil
}

// reservedInstancesRequest builds the full DescribeReservedInstances request
// from the configuration, separate from the read loop that sends it.
func reservedInstancesRequest(model reservedInstancesModel, regionId, scope string, pageNumber, pageSize int64) map[string]interface{} {
	request := map[string]interface{}{
		"RegionId":   regionId,
		"PageNumber": int(pageNumber),
		"PageSize":   int(pageSize),
	}
	ecsSetRequestString(request, "ReservedInstanceName", model.ReservedInstanceName)
	ecsSetRequestString(request, "ZoneId", model.ZoneId)
	ecsSetRequestString(request, "LockReason", model.LockReason)
	ecsSetRequestString(request, "InstanceType", model.InstanceType)
	ecsSetRequestString(request, "InstanceTypeFamily", model.InstanceTypeFamily)
	ecsSetRequestString(request, "Scope", types.StringValue(scope))
	ecsSetRequestString(request, "OfferingType", model.OfferingType)
	ecsSetRequestString(request, "AllocationType", model.AllocationType)
	ecsApplyRepeatList(request, "ReservedInstanceId", ecsElementsAsStrings(model.Ids))
	ecsApplyRepeatList(request, "Status", ecsElementsAsStrings(model.Status))
	ecsApplyTags(request, model.Tags)
	return request
}

func (d *reservedInstancesDataSource) reservedInstancesDescribeAll(ctx context.Context, client *connectivity.AliyunClient, request map[string]interface{}, startPage, pageSize int64) ([]map[string]interface{}, int64, error) {
	return ecsCollectNumberedPages(ctx, reservedInstancesAction, "ReservedInstanceId", reservedInstancesMaxPages, startPage, pageSize, func(pageNumber int64) (ecsNumberedPage, error) {
		request["PageNumber"] = int(pageNumber)
		request["PageSize"] = int(pageSize)
		response, err := reservedInstancesRpcPost(ctx, client, request)
		if err != nil {
			return ecsNumberedPage{}, err
		}
		if err := ecsValidateQueryEnvelope(response, reservedInstancesAction, "ReservedInstances", "TotalCount", pageNumber == 1); err != nil {
			return ecsNumberedPage{}, err
		}
		items, err := ecsDecodeRecordList(response, reservedInstancesAction, "ReservedInstances", "ReservedInstance")
		if err != nil {
			return ecsNumberedPage{}, err
		}
		totalCount, totalKnown, err := ecsDecodeCount(response, reservedInstancesAction, "TotalCount")
		if err != nil {
			return ecsNumberedPage{}, err
		}
		return ecsNumberedPage{items: items, total: totalCount, totalKnown: totalKnown}, nil
	})
}

func reservedInstanceObject(object map[string]interface{}) (types.Object, error) {
	identity, err := ecsRecordIdentity(object, reservedInstancesAction, "ReservedInstanceId")
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
	if err := ecsDecodeStringFields(object, []ecsStringField{
		{"reserved_instance_id", "ReservedInstanceId", false},
		{"reserved_instance_name", "ReservedInstanceName", false},
		{"status", "Status", false},
		{"creation_time", "CreationTime", false},
		{"start_time", "StartTime", false},
		{"expired_time", "ExpiredTime", false},
		{"instance_type", "InstanceType", false},
		{"region_id", "RegionId", false},
		{"zone_id", "ZoneId", false},
		{"scope", "Scope", false},
		{"platform", "Platform", false},
		{"offering_type", "OfferingType", false},
		{"allocation_status", "AllocationStatus", false},
		{"description", "Description", false},
		{"resource_group_id", "ResourceGroupId", false},
	}, attributes); err != nil {
		return nil, err
	}

	tags, err := ecsDecodeTags(object, ecsTagDecodeOptions{reservedKey: ecsReservedInstancesTagKey})
	if err != nil {
		return nil, err
	}
	attributes["tags"] = tags

	operationLocks, err := ecsOptionalObjectListMember(object, "OperationLocks", "OperationLock", reservedInstanceOperationLockObjectType, func(item map[string]interface{}) (attr.Value, error) {
		lockReason, err := ecsDecodeStringMember(item, "LockReason")
		if err != nil {
			return nil, err
		}
		return types.ObjectValueMust(reservedInstanceOperationLockAttrTypes, map[string]attr.Value{"lock_reason": lockReason}), nil
	})
	if err != nil {
		return nil, err
	}
	attributes["operation_locks"] = operationLocks

	instanceAmount, err := ecsDecodeInt64Member(object, "InstanceAmount", "InstanceAmount")
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
		if message := ecsValidateTagKey(key); message != "" {
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

func validateReservedInstancesTagValue(value string) string {
	if value == "" {
		return "A tag value must not be empty."
	}
	if utf8.RuneCountInString(value) > ecsTagMaxLength {
		return fmt.Sprintf("A tag value must be at most %d characters.", ecsTagMaxLength)
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

// reservedInstancesRpcClient is the context-bound RPC capability used by the
// read loop. The interface names exactly the one call the loop needs, so the
// loop owns every retry and each attempt issues one transport call.
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
		// The transport timeouts are derived from this read context's deadline, so
		// the transport can expire a moment before ctx.Err() is set. When the
		// deadline is that close, let the context report it first so the failure
		// always names the read deadline instead of whichever timer won.
		if errors.Is(rpcErr, context.DeadlineExceeded) && ctx.Err() == nil {
			if deadline, ok := ctx.Deadline(); ok {
				if remaining := time.Until(deadline); remaining < 250*time.Millisecond {
					timer := time.NewTimer(remaining + 100*time.Millisecond)
					select {
					case <-ctx.Done():
					case <-timer.C:
					}
					timer.Stop()
				}
			}
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("%s was cancelled while the request was in flight: %w", reservedInstancesAction, ctxErr)
		}
		if errors.Is(rpcErr, context.Canceled) || errors.Is(rpcErr, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%s reached its own cancellation or request deadline before the read context ended: %w", reservedInstancesAction, rpcErr)
		}
		if !ecsNeedRetry(rpcErr) || !time.Now().Add(backoff).Before(retryDeadline) {
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
