package rds

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const rcClusterTypeName = "alicloud_rds_rc_cluster"

var (
	_ resource.Resource                = &rcClusterResource{}
	_ resource.ResourceWithImportState = &rcClusterResource{}
)

func NewRCClusterResource() resource.Resource {
	return &rcClusterResource{}
}

type rcClusterResource struct {
	fwadapt.ResourceBase
}

type rcClusterResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	ClusterId                   types.String `tfsdk:"cluster_id"`
	RegionId                    types.String `tfsdk:"region_id"`
	ClusterName                 types.String `tfsdk:"cluster_name"`
	VClusterStatus              types.String `tfsdk:"vcluster_status"`
	VpcId                       types.String `tfsdk:"vpc_id"`
	VSwitchId                   types.String `tfsdk:"v_switch_id"`
	ClientKey                   types.String `tfsdk:"client_key"`
	ClientCertificate           types.String `tfsdk:"client_certificate"`
	SupportDiskPerformanceLevel types.List   `tfsdk:"support_disk_performance_level"`
	SupportDiskCategory         types.List   `tfsdk:"support_disk_category"`
}

func (r *rcClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = rcClusterTypeName
}

func (r *rcClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a RDS Custom container cluster resource. The resource imports an existing ACK managed cluster into RDS Custom; there is no update operation, so changing an argument replaces the resource and re-imports the cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, which is the ACK cluster ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the ACK managed cluster to import into RDS Custom.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region ID of the cluster. Defaults to the region of the provider.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vcluster_status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the cluster. UNINITIALIZED means the cluster is not imported into RDS Custom.",
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the VPC that the ACK cluster belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"v_switch_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the vSwitch that the ACK cluster belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The client-key-data of the ACK cluster KubeConfig.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_certificate": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The client-certificate-data of the ACK cluster KubeConfig.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"support_disk_performance_level": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The disk performance levels supported by the cluster.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"support_disk_category": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The disk categories supported by the cluster.",
			},
		},
	}
}

func (r *rcClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rcClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := plan.ClusterId.ValueString()
	regionId := plan.RegionId.ValueString()
	if regionId == "" {
		regionId = r.Client().RegionId
	}

	request := map[string]interface{}{
		"RegionId":          regionId,
		"ClusterId":         clusterId,
		"VpcId":             plan.VpcId.ValueString(),
		"VSwitchId":         plan.VSwitchId.ValueString(),
		"ClientKey":         plan.ClientKey.ValueString(),
		"ClientCertificate": plan.ClientCertificate.ValueString(),
	}
	if values, diags := expandRCClusterStringList(ctx, plan.SupportDiskPerformanceLevel); len(values) > 0 {
		resp.Diagnostics.Append(diags...)
		request["SupportDiskPerformanceLevel"] = values
	}
	if values, diags := expandRCClusterStringList(ctx, plan.SupportDiskCategory); len(values) > 0 {
		resp.Diagnostics.Append(diags...)
		request["SupportDiskCategory"] = values
	}
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.Client().RpcPost(rcClusterRdsProduct, rcClusterRdsVersion, "CreateRCVCluster", nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s: calling CreateRCVCluster", rcClusterTypeName),
			errs.WrapErrorf(err, "Importing %s: calling CreateRCVCluster", rcClusterTypeName).Error(),
		)
		return
	}

	object, err := waitRCVClusterAvailable(ctx, r.Client(), clusterId, rcClusterCreateWindow)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s: waiting for cluster %s", rcClusterTypeName, clusterId),
			errs.WrapErrorf(err, "Importing %s: waiting for cluster %s", rcClusterTypeName, clusterId).Error(),
		)
		return
	}

	state, diags := rcClusterState(object, clusterId, regionId, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *rcClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rcClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.ClusterId.ValueString()
	if clusterId == "" {
		clusterId = state.Id.ValueString()
	}
	if clusterId == "" || strings.Contains(clusterId, ":") {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: invalid resource ID", rcClusterTypeName),
			fmt.Sprintf("The resource ID %q is not an ACK cluster ID", state.Id.ValueString()),
		)
		return
	}

	object, err := describeRCVCluster(r.Client(), clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeRCVCluster", rcClusterTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeRCVCluster", rcClusterTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState, diags := rcClusterState(object, clusterId, r.Client().RegionId, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *rcClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Every argument is ForceNew, so an update can only be a no-op refresh.
	var plan rcClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *rcClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state rcClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.ClusterId.ValueString()
	if clusterId == "" {
		clusterId = state.Id.ValueString()
	}

	lastStatus, err := deleteRCVCluster(ctx, r.Client(), clusterId, rcClusterDeleteWindow)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteRCVCluster", rcClusterTypeName),
			fmt.Sprintf("%s (last VClusterStatus=%q)", errs.WrapErrorf(err, "Deleting %s: calling DeleteRCVCluster", rcClusterTypeName).Error(), lastStatus),
		)
		return
	}
}

func (r *rcClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// rcClusterState maps a DescribeRCVCluster payload onto the resource model.
// Values that the API does not return keep the value already known to the
// prior state or plan.
func rcClusterState(object map[string]interface{}, clusterId, defaultRegion string, prior rcClusterResourceModel) (rcClusterResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	regionId := rcClusterString(object["Region"])
	if regionId.IsNull() {
		regionId = rcClusterString(object["RegionId"])
	}
	if regionId.IsNull() && defaultRegion != "" {
		regionId = types.StringValue(defaultRegion)
	}
	if regionId.IsNull() {
		regionId = prior.RegionId
	}

	vpcId := rcClusterString(object["VpcId"])
	if vpcId.IsNull() {
		vpcId = prior.VpcId
	}
	performanceLevels := rcClusterStringList(object["SupportDiskPerformanceLevel"])
	if performanceLevels.IsNull() {
		if !prior.SupportDiskPerformanceLevel.IsUnknown() {
			performanceLevels = prior.SupportDiskPerformanceLevel
		} else {
			performanceLevels = types.ListNull(types.StringType)
		}
	}

	state := rcClusterResourceModel{
		Id:                          types.StringValue(clusterId),
		ClusterId:                   types.StringValue(clusterId),
		RegionId:                    regionId,
		ClusterName:                 rcClusterString(object["ClusterName"]),
		VClusterStatus:              rcClusterString(object["VClusterStatus"]),
		VpcId:                       vpcId,
		VSwitchId:                   prior.VSwitchId,
		ClientKey:                   prior.ClientKey,
		ClientCertificate:           prior.ClientCertificate,
		SupportDiskPerformanceLevel: performanceLevels,
		SupportDiskCategory:         prior.SupportDiskCategory,
	}
	return state, diags
}

// expandRCClusterStringList reads an optional list-of-string value.
func expandRCClusterStringList(ctx context.Context, value types.List) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value.IsNull() || value.IsUnknown() {
		return nil, diags
	}
	var values []string
	diags.Append(value.ElementsAs(ctx, &values, false)...)
	return values, diags
}

func rcClusterString(raw interface{}) types.String {
	value, ok := raw.(string)
	if !ok || value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

func rcClusterStringList(raw interface{}) types.List {
	items, ok := raw.([]interface{})
	if !ok || len(items) == 0 {
		return types.ListNull(types.StringType)
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		value, ok := item.(string)
		if !ok {
			continue
		}
		values = append(values, types.StringValue(value))
	}
	if len(values) == 0 {
		return types.ListNull(types.StringType)
	}
	return types.ListValueMust(types.StringType, values)
}
