package sas

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	groupedInstancesTypeName   = "alicloud_security_center_grouped_instances"
	groupedInstancesAction     = "DescribeGroupedInstances"
	groupedInstancesAPIVersion = "2018-12-03"
	groupedInstancesReadWindow = 3 * time.Minute
)

// groupedInstanceAttrTypes must not use nested attributes: the provider is
// served over protocol version 5, which rejects them.
var groupedInstanceAttrTypes = map[string]attr.Type{
	"asap_vul_instance_count":    types.Int64Type,
	"unprotected_instance_count": types.StringType,
	"field_alias_name":           types.StringType,
	"instance_count":             types.StringType,
	"risk_instance_count":        types.StringType,
	"group_flag":                 types.Int64Type,
	"auth_version_check_count":   types.Int64Type,
	"os":                         types.StringType,
	"instance_core_count":        types.Int64Type,
}

var (
	_ datasource.DataSource              = &groupedInstancesDataSource{}
	_ datasource.DataSourceWithConfigure = &groupedInstancesDataSource{}
)

func NewGroupedInstancesDataSource() datasource.DataSource {
	return &groupedInstancesDataSource{}
}

type groupedInstancesDataSource struct {
	fwadapt.DataSourceBase
}

type groupedInstancesModel struct {
	Id                         types.String `tfsdk:"id"`
	GroupField                 types.String `tfsdk:"group_field"`
	MachineTypes               types.String `tfsdk:"machine_types"`
	FieldValue                 types.String `tfsdk:"field_value"`
	Vendor                     types.Int64  `tfsdk:"vendor"`
	Vendors                    types.String `tfsdk:"vendors"`
	SaleVersionCheckCode       types.String `tfsdk:"sale_version_check_code"`
	Lang                       types.String `tfsdk:"lang"`
	ResourceDirectoryAccountId types.Int64  `tfsdk:"resource_directory_account_id"`
	NoPage                     types.Bool   `tfsdk:"no_page"`
	PageSize                   types.Int64  `tfsdk:"page_size"`
	CurrentPage                types.Int64  `tfsdk:"current_page"`
	TotalCount                 types.Int64  `tfsdk:"total_count"`
	PageItemCount              types.Int64  `tfsdk:"page_item_count"`
	Instances                  types.List   `tfsdk:"instances"`
}

func (d *groupedInstancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = groupedInstancesTypeName
}

func (d *groupedInstancesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the statistics of the assets in Security Center, grouped by the specified dimension.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, which is derived from the query condition.",
			},
			"group_field": schema.StringAttribute{
				Required:    true,
				Description: "The dimension by which the assets are grouped. Valid values: `groupId`, `regionId` and `vpcInstanceId`.",
				Validators: []validator.String{
					stringvalidator.OneOf("groupId", "regionId", "vpcInstanceId"),
				},
			},
			"machine_types": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the assets to query. The value is fixed to `ecs`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ecs"),
				},
			},
			"field_value": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the asset group to query. Fuzzy match is supported.",
			},
			"vendor": schema.Int64Attribute{
				Optional:    true,
				Description: "The vendor of the server. Valid values: `0` (Alibaba Cloud), `1` (outside Alibaba Cloud), `2` (IDC), `3`, `4`, `5`, `7` (other clouds) and `8` (lightweight server).",
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2, 3, 4, 5, 7, 8),
				},
			},
			"vendors": schema.StringAttribute{
				Optional:    true,
				Description: "The vendors of the servers. Separate multiple values with commas (,). Valid values are the same as those of `vendor`.",
			},
			"sale_version_check_code": schema.StringAttribute{
				Optional:    true,
				Description: "The sale condition used to filter the assets. Valid values: `sas_gte_advanced`, `sas_gte_enterprise`, `sas_gt_basic`, `sas_eq_advanced` and `sas_gt_anti_virus`.",
				Validators: []validator.String{
					stringvalidator.OneOf("sas_gte_advanced", "sas_gte_enterprise", "sas_gt_basic", "sas_eq_advanced", "sas_gt_anti_virus"),
				},
			},
			"lang": schema.StringAttribute{
				Optional:    true,
				Description: "The language of the request and the response. Default value: `zh`. Valid values: `zh` and `en`.",
				Validators: []validator.String{
					stringvalidator.OneOf("zh", "en"),
				},
			},
			"resource_directory_account_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The ID of the Alibaba Cloud account that is a member of the resource directory.",
			},
			"no_page": schema.BoolAttribute{
				Optional:    true,
				Description: "Specifies whether to disable the paged query. Set it to `true` to return all the results in one response, in which case `page_size` and `current_page` do not take effect. The query is paged when it is not set or set to `false`.",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of entries per page. Default value: `20`.",
			},
			"current_page": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of the page to return. Default value: `1`.",
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of the assets that match the query condition.",
			},
			"page_item_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of the assets on the current page.",
			},
			"instances": schema.ListAttribute{
				Computed:    true,
				Description: "A list of the asset statistics that match the query condition.",
				ElementType: types.ObjectType{AttrTypes: groupedInstanceAttrTypes},
			},
		},
	}
}

func (d *groupedInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data groupedInstancesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"GroupField": data.GroupField.ValueString(),
	}
	if !data.MachineTypes.IsNull() && !data.MachineTypes.IsUnknown() {
		request["MachineTypes"] = data.MachineTypes.ValueString()
	}
	if !data.FieldValue.IsNull() && !data.FieldValue.IsUnknown() {
		request["FieldValue"] = data.FieldValue.ValueString()
	}
	if !data.Vendor.IsNull() && !data.Vendor.IsUnknown() {
		request["Vendor"] = data.Vendor.ValueInt64()
	}
	if !data.Vendors.IsNull() && !data.Vendors.IsUnknown() {
		request["Vendors"] = data.Vendors.ValueString()
	}
	if !data.SaleVersionCheckCode.IsNull() && !data.SaleVersionCheckCode.IsUnknown() {
		request["SaleVersionCheckCode"] = data.SaleVersionCheckCode.ValueString()
	}
	if !data.Lang.IsNull() && !data.Lang.IsUnknown() {
		request["Lang"] = data.Lang.ValueString()
	}
	if !data.ResourceDirectoryAccountId.IsNull() && !data.ResourceDirectoryAccountId.IsUnknown() {
		request["ResourceDirectoryAccountId"] = data.ResourceDirectoryAccountId.ValueInt64()
	}
	if !data.NoPage.IsNull() && !data.NoPage.IsUnknown() {
		request["NoPage"] = data.NoPage.ValueBool()
	}
	if !data.PageSize.IsNull() && !data.PageSize.IsUnknown() {
		request["PageSize"] = data.PageSize.ValueInt64()
	}
	if !data.CurrentPage.IsNull() && !data.CurrentPage.IsUnknown() {
		request["CurrentPage"] = data.CurrentPage.ValueInt64()
	}

	response, err := helper.RetryRpcPost(ctx, groupedInstancesAction, func() (map[string]interface{}, error) {
		return d.Client().RpcPost("Sas", groupedInstancesAPIVersion, groupedInstancesAction, nil, request, true)
	}, groupedInstancesReadWindow, helper.NeedRetry)
	logs.AddDebug(groupedInstancesAction, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", groupedInstancesTypeName, groupedInstancesAction),
			err.Error(),
		)
		return
	}

	rawInstances, _ := response["Instances"].([]interface{})
	instances := make([]attr.Value, 0, len(rawInstances))
	for _, raw := range rawInstances {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		instances = append(instances, types.ObjectValueMust(groupedInstanceAttrTypes, map[string]attr.Value{
			"asap_vul_instance_count":    toInt64Value(item["AsapVulInstanceCount"]),
			"unprotected_instance_count": toStringValue(item["UnProtectedInstanceCount"]),
			"field_alias_name":           toStringValue(item["FieldAliasName"]),
			"instance_count":             toStringValue(item["InstanceCount"]),
			"risk_instance_count":        toStringValue(item["RiskInstanceCount"]),
			"group_flag":                 toInt64Value(item["GroupFlag"]),
			"auth_version_check_count":   toInt64Value(item["AuthVersionCheckCount"]),
			"os":                         toStringValue(item["Os"]),
			"instance_core_count":        toInt64Value(item["InstanceCoreCount"]),
		}))
	}
	instanceList, diags := types.ListValue(types.ObjectType{AttrTypes: groupedInstanceAttrTypes}, instances)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pageInfo, _ := response["PageInfo"].(map[string]interface{})

	state := groupedInstancesModel{
		Id:                         types.StringValue(groupedInstancesId(request)),
		GroupField:                 data.GroupField,
		MachineTypes:               data.MachineTypes,
		FieldValue:                 data.FieldValue,
		Vendor:                     data.Vendor,
		Vendors:                    data.Vendors,
		SaleVersionCheckCode:       data.SaleVersionCheckCode,
		Lang:                       data.Lang,
		ResourceDirectoryAccountId: data.ResourceDirectoryAccountId,
		NoPage:                     data.NoPage,
		PageSize:                   echoInt64(data.PageSize, pageInfo["PageSize"]),
		CurrentPage:                echoInt64(data.CurrentPage, pageInfo["CurrentPage"]),
		TotalCount:                 toInt64Value(pageInfo["TotalCount"]),
		PageItemCount:              toInt64Value(pageInfo["Count"]),
		Instances:                  instanceList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// groupedInstancesId derives a stable ID from the query condition so that a
// different condition is a different data source instance.
func groupedInstancesId(request map[string]interface{}) string {
	keys := make([]string, 0, len(request))
	for k := range request {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, request[k]))
	}
	return helper.DataResourceIdHash(parts)
}

// echoInt64 keeps a configured value because a paged query result cannot
// change it, and falls back to the value the API reports.
func echoInt64(configured types.Int64, reported interface{}) types.Int64 {
	if !configured.IsNull() && !configured.IsUnknown() {
		return configured
	}
	return toInt64Value(reported)
}

func toInt64Value(value interface{}) types.Int64 {
	switch v := value.(type) {
	case nil:
		return types.Int64Null()
	case float64:
		return types.Int64Value(int64(v))
	case float32:
		return types.Int64Value(int64(v))
	case int:
		return types.Int64Value(int64(v))
	case int32:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case string:
		if v == "" {
			return types.Int64Null()
		}
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return types.Int64Null()
		}
		return types.Int64Value(parsed)
	default:
		parsed, err := strconv.ParseInt(fmt.Sprint(v), 10, 64)
		if err != nil {
			return types.Int64Null()
		}
		return types.Int64Value(parsed)
	}
}

func toStringValue(value interface{}) types.String {
	switch v := value.(type) {
	case nil:
		return types.StringNull()
	case string:
		return types.StringValue(v)
	default:
		return types.StringValue(fmt.Sprint(v))
	}
}
