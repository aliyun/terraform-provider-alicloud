package oos

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const opsItemConfigurationTypeName = "alicloud_oos_ops_item_configuration"

var (
	_ resource.Resource                = &opsItemConfigurationResource{}
	_ resource.ResourceWithImportState = &opsItemConfigurationResource{}
)

func NewOpsItemConfigurationResource() resource.Resource {
	return &opsItemConfigurationResource{}
}

type opsItemConfigurationResource struct {
	fwadapt.ResourceBase
}

type opsItemConfigurationResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	ConfigurationId          types.String `tfsdk:"configuration_id"`
	OpsItemConfigurationName types.String `tfsdk:"ops_item_configuration_name"`
	DataSource               types.String `tfsdk:"data_source"`
	ResourceGroupId          types.String `tfsdk:"resource_group_id"`
	CreateTime               types.String `tfsdk:"create_time"`
	RegionId                 types.String `tfsdk:"region_id"`
}

func (r *opsItemConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = opsItemConfigurationTypeName
}

func (r *opsItemConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an OOS OpsItem Configuration resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `configuration_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"configuration_id": schema.StringAttribute{
				Computed:    true,
				Description: "The first-level ID of the OpsItem configuration. It is the unique identifier assigned by the system on creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ops_item_configuration_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the OpsItem configuration.",
			},
			"resource_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the resource group.",
			},
			"data_source": schema.StringAttribute{
				Required:    true,
				Description: "The data source of the OpsItem configuration, JSON-encoded. `Type` is required and must be one of `ALIYUN/ROS`, `ALIYUN/ITSM`, `ALIYUN/ARMS`, `ALIYUN/CMS`, `ALIYUN/CIAM`, `ALIYUN/OOS`. Optional fields: `Arn`, `Path`, `ExcludePath`, `NodeId`, `Interval`, `FilterCondition`, `Arguments`.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *opsItemConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan opsItemConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	request := map[string]interface{}{
		"ConfigurationName": plan.OpsItemConfigurationName.ValueString(),
		"RegionId":          client.RegionId,
		"DataSource":        plan.DataSource.ValueString(),
	}
	if !plan.ResourceGroupId.IsNull() {
		request["ResourceGroupId"] = plan.ResourceGroupId.ValueString()
	}

	response, err := client.RpcPost("oos", "2019-06-01", "CreateOpsItemConfiguration", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateOpsItemConfiguration", opsItemConfigurationTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateOpsItemConfiguration", opsItemConfigurationTypeName).Error(),
		)
		return
	}

	object, _ := response["OpsItemConfiguration"].(map[string]interface{})
	configurationId, _ := object["ConfigurationId"].(string)
	if configurationId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateOpsItemConfiguration", opsItemConfigurationTypeName),
			fmt.Sprintf("The response carried no OpsItemConfiguration.ConfigurationId: %v", response),
		)
		return
	}

	created, err := getOpsItemConfiguration(client, configurationId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetOpsItemConfiguration", opsItemConfigurationTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetOpsItemConfiguration", opsItemConfigurationTypeName).Error(),
		)
		return
	}
	if created == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateOpsItemConfiguration", opsItemConfigurationTypeName),
			fmt.Sprintf("The OpsItem configuration %s is not visible via GetOpsItemConfiguration right after creation", configurationId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, opsItemConfigurationState(created, client.RegionId, plan.DataSource.ValueString()))...)
}

func (r *opsItemConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state opsItemConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	object, err := getOpsItemConfiguration(client, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetOpsItemConfiguration", opsItemConfigurationTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetOpsItemConfiguration", opsItemConfigurationTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, opsItemConfigurationState(object, client.RegionId, state.DataSource.ValueString()))...)
}

func (r *opsItemConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state opsItemConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	// UpdateOpsItemConfiguration requires both ConfigurationId and ConfigurationName.
	request := map[string]interface{}{
		"ConfigurationId":   state.Id.ValueString(),
		"ConfigurationName": plan.OpsItemConfigurationName.ValueString(),
		"RegionId":          client.RegionId,
		"DataSource":        plan.DataSource.ValueString(),
	}
	if !plan.ResourceGroupId.IsNull() {
		request["ResourceGroupId"] = plan.ResourceGroupId.ValueString()
	}

	if _, err := client.RpcPost("oos", "2019-06-01", "UpdateOpsItemConfiguration", nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateOpsItemConfiguration", opsItemConfigurationTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateOpsItemConfiguration", opsItemConfigurationTypeName).Error(),
		)
		return
	}

	object, err := getOpsItemConfiguration(client, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetOpsItemConfiguration", opsItemConfigurationTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetOpsItemConfiguration", opsItemConfigurationTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, opsItemConfigurationState(object, client.RegionId, plan.DataSource.ValueString()))...)
}

func (r *opsItemConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state opsItemConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	// DeleteOpsItemConfigurations accepts a comma-separated list of ConfigurationIds.
	_, err := client.RpcPost("oos", "2019-06-01", "DeleteOpsItemConfigurations", nil, map[string]interface{}{
		"ConfigurationIds": state.Id.ValueString(),
		"RegionId":         client.RegionId,
	}, true)
	if err != nil && !isOpsItemConfigurationNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteOpsItemConfigurations", opsItemConfigurationTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteOpsItemConfigurations", opsItemConfigurationTypeName).Error(),
		)
		return
	}
}

func (r *opsItemConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getOpsItemConfiguration queries GetOpsItemConfiguration by ConfigurationId and returns
// the OpsItemConfiguration object from the response, or nil when it no longer exists.
func getOpsItemConfiguration(client *connectivity.AliyunClient, configurationId string) (map[string]interface{}, error) {
	response, err := client.RpcPost("oos", "2019-06-01", "GetOpsItemConfiguration", nil, map[string]interface{}{
		"ConfigurationId": configurationId,
		"RegionId":        client.RegionId,
	}, true)
	if err != nil {
		if isOpsItemConfigurationNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	object, _ := response["OpsItemConfiguration"].(map[string]interface{})
	if object == nil {
		return nil, fmt.Errorf("the GetOpsItemConfiguration response carried no OpsItemConfiguration object: %v", response)
	}
	return object, nil
}

// opsItemConfigurationState maps the OpsItemConfiguration object of a GetOpsItemConfiguration
// (or CreateOpsItemConfiguration) response onto the resource model. Empty strings become null.
// DataSource is read from the response when present; otherwise the dataSourceFallback
// (the prior plan/state value) is preserved, because the API does not always echo it back.
func opsItemConfigurationState(object map[string]interface{}, regionId, dataSourceFallback string) opsItemConfigurationResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	configurationId := str("ConfigurationId")
	region := regionId
	if region == "" {
		region = str("RegionId").ValueString()
	}
	dataSource := str("DataSource")
	if dataSource.IsNull() && dataSourceFallback != "" {
		dataSource = types.StringValue(dataSourceFallback)
	}
	return opsItemConfigurationResourceModel{
		Id:                       configurationId,
		ConfigurationId:          configurationId,
		OpsItemConfigurationName: str("ConfigurationName"),
		DataSource:               dataSource,
		ResourceGroupId:          str("ResourceGroupId"),
		CreateTime:               str("CreateDate"),
		RegionId:                 types.StringValue(region),
	}
}

// isOpsItemConfigurationNotExist reports whether err is the OOS OpsItemConfiguration
// not-found error. The code arrives on a *tea.SDKError from RpcPost; the substring
// fallback covers wrapped forms, and errs.NotFoundError covers the generic 404 path.
func isOpsItemConfigurationNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "EntityNotExists.OpsItemConfiguration" {
		return true
	}
	if strings.Contains(err.Error(), "EntityNotExists.OpsItemConfiguration") {
		return true
	}
	return errs.NotFoundError(err)
}
