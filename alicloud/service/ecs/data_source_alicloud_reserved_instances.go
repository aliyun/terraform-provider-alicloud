package ecs

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
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
				Description: "A list of reserved instance IDs to filter by. The filter is sent to the server.",
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
				Description: "The zone ID. It is required when `scope` is `Zone` and the request is rejected without it. When `zone_id` is set without `scope`, the data source requests `Scope = Zone`; when neither is set, the query is not restricted to a scope.",
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
				Description: "The scope of the reserved instance. Valid values: `Region` and `Zone`. When set to `Zone`, `zone_id` is required.",
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
				Description: "A mapping of tags to filter by. Each entry is sent as a `Tag` filter. A maximum of 20 entries is allowed. A key and a value must be 1 to 128 characters; a key must not start with `aliyun` or `acs:`, a value must not start with `acs:`, and neither may contain a URL. A tag filter returns at most 1,000 reserved instances; when more may match, query `ListTagResources` instead.",
				Validators: []validator.Map{
					mapvalidator.SizeAtMost(20),
					ecsFwTagMapValidator{valueForbiddenPrefixes: []string{"acs:"}},
				},
			},
			"page_number": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of the first page to read. Valid values: 1 or greater. Defaults to 1. The data source continues with the following pages until a short or empty page, the server total or the read budget is reached.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page. Valid values: 1 to 100. Defaults to 100. A page shorter than this value ends the read.",
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of reserved instances that match all configured filters, as reported by the `TotalCount` of the API response. It counts every matching reserved instance in the region, not only the rows returned from `page_number`; when the page parameter is ignored and the read is aborted, no value is written.",
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
	if scope == "Zone" && zoneId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: invalid configuration", reservedInstancesTypeName),
			"`zone_id` must be set when `scope` is `Zone`.",
		)
		return
	}
	if scope == "" && zoneId != "" {
		scope = "Zone"
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

	request := reservedInstancesRequest(model, d.Meta().RegionId, scope, pageNumber, pageSize)

	objects, totalCount, err := d.reservedInstancesDescribeAll(ctx, request, pageNumber, pageSize)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", reservedInstancesTypeName, reservedInstancesAction),
			err.Error(),
		)
		return
	}

	sorted := append([]map[string]interface{}(nil), objects...)
	sort.Slice(sorted, func(i, j int) bool {
		return ecsFwString(sorted[i]["ReservedInstanceId"]).ValueString() < ecsFwString(sorted[j]["ReservedInstanceId"]).ValueString()
	})

	ids := make([]string, 0, len(sorted))
	names := make([]string, 0, len(sorted))
	elements := make([]attr.Value, 0, len(sorted))
	for _, object := range sorted {
		ids = append(ids, ecsFwString(object["ReservedInstanceId"]).ValueString())
		names = append(names, ecsFwString(object["ReservedInstanceName"]).ValueString())
		element, err := reservedInstanceObject(object)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: decoding a record", reservedInstancesTypeName),
				err.Error(),
			)
			return
		}
		elements = append(elements, element)
	}

	// A cancelled or expired read must not write a partial result to state.
	if ctxErr := ctx.Err(); ctxErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: the read did not finish", reservedInstancesTypeName),
			fmt.Sprintf("The read was cancelled or ran out of time before the result was written: %v", ctxErr),
		)
		return
	}

	model.Id = types.StringValue(helper.DataResourceIdHash(ids))
	model.Ids = ecsFwStringListValue(ids)
	model.Names = ecsFwStringListValue(names)
	model.TotalCount = types.Int64Value(totalCount)
	model.Instances = types.ListValueMust(reservedInstanceObjectType, elements)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
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
	ecsFwSetRequestString(request, "ReservedInstanceName", model.ReservedInstanceName)
	ecsFwSetRequestString(request, "ZoneId", model.ZoneId)
	ecsFwSetRequestString(request, "LockReason", model.LockReason)
	ecsFwSetRequestString(request, "InstanceType", model.InstanceType)
	ecsFwSetRequestString(request, "InstanceTypeFamily", model.InstanceTypeFamily)
	ecsFwSetRequestString(request, "Scope", types.StringValue(scope))
	ecsFwSetRequestString(request, "OfferingType", model.OfferingType)
	ecsFwSetRequestString(request, "AllocationType", model.AllocationType)
	ecsFwApplyRepeatList(request, "ReservedInstanceId", ecsFwElementsAsStrings(model.Ids))
	ecsFwApplyRepeatList(request, "Status", ecsFwElementsAsStrings(model.Status))
	ecsFwApplyTags(request, model.Tags)
	return request
}

func (d *reservedInstancesDataSource) reservedInstancesDescribeAll(ctx context.Context, request map[string]interface{}, startPage, pageSize int64) ([]map[string]interface{}, int64, error) {
	return ecsFwCollectNumberedPages(ctx, reservedInstancesMaxPages, startPage, pageSize, reservedInstancesAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		request["PageNumber"] = int(pageNumber)
		request["PageSize"] = int(pageSize)
		response, err := ecsFwRpcPost(ctx, d.Meta(), reservedInstancesTypeName, reservedInstancesAction, request)
		if err != nil {
			return nil, 0, err
		}
		items, err := ecsFwDecodeRecordList(response, "ReservedInstances", "ReservedInstance", reservedInstancesAction)
		if err != nil {
			return nil, 0, err
		}
		totalCount, _, err := ecsFwDecodeCount(response, "TotalCount", reservedInstancesAction)
		if err != nil {
			return nil, 0, err
		}
		return items, totalCount, nil
	})
}

func reservedInstanceObject(object map[string]interface{}) (types.Object, error) {
	tags, err := ecsFwDecodeTags(object["Tags"])
	if err != nil {
		return types.ObjectNull(reservedInstanceAttrTypes), fmt.Errorf("%s: %w", ecsFwString(object["ReservedInstanceId"]).ValueString(), err)
	}
	return types.ObjectValueMust(reservedInstanceAttrTypes, map[string]attr.Value{
		"reserved_instance_id":   ecsFwString(object["ReservedInstanceId"]),
		"reserved_instance_name": ecsFwString(object["ReservedInstanceName"]),
		"status":                 ecsFwString(object["Status"]),
		"creation_time":          ecsFwString(object["CreationTime"]),
		"start_time":             ecsFwString(object["StartTime"]),
		"expired_time":           ecsFwString(object["ExpiredTime"]),
		"instance_type":          ecsFwString(object["InstanceType"]),
		"instance_amount":        ecsFwInt64(object["InstanceAmount"]),
		"region_id":              ecsFwString(object["RegionId"]),
		"zone_id":                ecsFwString(object["ZoneId"]),
		"scope":                  ecsFwString(object["Scope"]),
		"platform":               ecsFwString(object["Platform"]),
		"offering_type":          ecsFwString(object["OfferingType"]),
		"allocation_status":      ecsFwString(object["AllocationStatus"]),
		"description":            ecsFwString(object["Description"]),
		"resource_group_id":      ecsFwString(object["ResourceGroupId"]),
		"operation_locks":        reservedInstanceOperationLocks(object["OperationLocks"]),
		"tags":                   tags,
	}), nil
}

func reservedInstanceOperationLocks(value interface{}) types.List {
	elements := make([]attr.Value, 0)
	container, ok := value.(map[string]interface{})
	if !ok {
		return types.ListValueMust(reservedInstanceOperationLockObjectType, elements)
	}
	raw, ok := container["OperationLock"]
	if !ok {
		return types.ListValueMust(reservedInstanceOperationLockObjectType, elements)
	}
	list, ok := raw.([]interface{})
	if !ok {
		return types.ListValueMust(reservedInstanceOperationLockObjectType, elements)
	}
	for _, entry := range list {
		lock, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		elements = append(elements, types.ObjectValueMust(reservedInstanceOperationLockAttrTypes, map[string]attr.Value{
			"lock_reason": ecsFwString(lock["LockReason"]),
		}))
	}
	return types.ListValueMust(reservedInstanceOperationLockObjectType, elements)
}
