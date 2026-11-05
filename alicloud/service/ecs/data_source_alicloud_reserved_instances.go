package ecs

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
				Description: "The zone ID. The value is only effective when `scope` is `Zone`; when it is set without `scope`, the data source requests `Scope = Zone`.",
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
				Description: "The scope of the reserved instance. Valid values: `Region` and `Zone`.",
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
				Description: "A list of statuses to filter by. Valid values: `Creating`, `Active`, `Retired`, `Updating` and `Expired`.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(100),
					listvalidator.ValueStringsAre(stringvalidator.OneOf(reservedInstanceStatuses...)),
				},
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A mapping of tags to filter by. Each entry is sent as a `Tag` filter. A maximum of 20 entries is allowed.",
				Validators: []validator.Map{
					mapvalidator.SizeAtMost(20),
				},
			},
			"page_number": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of the first page to read. Valid values: 1 or greater. Defaults to 1.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page. Valid values: 1 to 100. Defaults to 100.",
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of reserved instances that match all configured filters.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of reserved instance names that match the configured filters.",
			},
			"instances": schema.ListAttribute{
				Computed:    true,
				ElementType: reservedInstanceObjectType,
				Description: "A list of reserved instances that match the configured filters. The `operation_locks` member of the API response is not exported.",
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

	request := map[string]interface{}{
		"RegionId":   d.Meta().RegionId,
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
		elements = append(elements, reservedInstanceObject(object))
	}

	model.Id = types.StringValue(helper.DataResourceIdHash(ids))
	model.Ids = ecsFwStringListValue(ids)
	model.Names = ecsFwStringListValue(names)
	model.TotalCount = types.Int64Value(totalCount)
	model.Instances = types.ListValueMust(reservedInstanceObjectType, elements)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (d *reservedInstancesDataSource) reservedInstancesDescribeAll(ctx context.Context, request map[string]interface{}, startPage, pageSize int64) ([]map[string]interface{}, int64, error) {
	deadline := time.Now().Add(reservedInstancesReadBudget)
	pagesRead := 0
	return ecsFwCollectNumberedPages(reservedInstancesMaxPages, startPage, pageSize, reservedInstancesAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		if time.Now().After(deadline) {
			return nil, 0, fmt.Errorf("%s read exceeded the %s budget after %d pages; the result was not returned because it would be incomplete", reservedInstancesAction, reservedInstancesReadBudget, pagesRead)
		}
		pagesRead++
		request["PageNumber"] = int(pageNumber)
		request["PageSize"] = int(pageSize)
		response, err := ecsFwRpcPost(ctx, d.Meta(), reservedInstancesTypeName, reservedInstancesAction, request)
		if err != nil {
			return nil, 0, err
		}
		totalCount := int64(0)
		if value, ok := response["TotalCount"]; ok {
			totalCount = ecsFwInt64(value).ValueInt64()
		}
		return reservedInstancesItems(response), totalCount, nil
	})
}

func reservedInstancesItems(response map[string]interface{}) []map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	if response == nil {
		return items
	}
	raw, err := jsonpath.Get("$.ReservedInstances.ReservedInstance", response)
	if err != nil {
		return items
	}
	list, ok := raw.([]interface{})
	if !ok {
		return items
	}
	for _, entry := range list {
		if item, ok := entry.(map[string]interface{}); ok {
			items = append(items, item)
		}
	}
	return items
}

func reservedInstanceObject(object map[string]interface{}) types.Object {
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
		"tags":                   ecsFwTagsValue(object["Tags"]),
	})
}
