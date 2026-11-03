// Package polardb. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package polardb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const networkChannelTypeName = "alicloud_polardb_network_channel"

var (
	_ resource.Resource                = &networkChannelResource{}
	_ resource.ResourceWithImportState = &networkChannelResource{}
)

func NewNetworkChannelResource() resource.Resource {
	return &networkChannelResource{}
}

type networkChannelResource struct {
	fwadapt.ResourceBase
}

type networkChannelResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	FwTimeouts           types.Object `tfsdk:"timeouts"`
	ChannelInfos         types.List   `tfsdk:"channel_infos"`
	ChannelName          types.String `tfsdk:"channel_name"`
	DBClusterId          types.String `tfsdk:"db_cluster_id"`
	Notes                types.String `tfsdk:"notes"`
	OwnerAccount         types.String `tfsdk:"owner_account"`
	OwnerId              types.Int64  `tfsdk:"owner_id"`
	ResourceGroupId      types.String `tfsdk:"resource_group_id"`
	ResourceOwnerAccount types.String `tfsdk:"resource_owner_account"`
	ResourceOwnerId      types.Int64  `tfsdk:"resource_owner_id"`
	TargetDBClusterId    types.String `tfsdk:"target_db_cluster_id"`
	TargetIp             types.String `tfsdk:"target_ip"`
	TargetPort           types.String `tfsdk:"target_port"`
	VpcId                types.String `tfsdk:"vpc_id"`
}

func (r *networkChannelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = networkChannelTypeName
}

func (r *networkChannelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel_infos": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"channel_name": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Description: "Channel Name",
						},
						"target_port": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Description: "Target Port",
						},
						"vpc_id": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Description: "VPC ID",
						},
						"target_db_cluster_id": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Description: "Target Database Cluster ID",
						},
						"db_cluster_id": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Description: "Database Cluster ID",
						},
						"target_type": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Description: "Target Type",
						},
						"target_ip": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Description: "Target IP",
						},
						"region_id": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Description: "Region ID",
						},
						"notes": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
							Description: "Notes",
						},
					},
				},
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIf(func(_ context.Context, req planmodifier.ListRequest, resp *listplanmodifier.RequiresReplaceIfFuncResponse) {
						resp.RequiresReplace = fwConfiguredChanged(req.ConfigValue, req.PlanValue, req.StateValue)
					}, "Replace when configured nested values change", "Replace when configured nested values change"),
				},
				Description: "Channel Network List",
			},
			"channel_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The name of the network channel.",
			},
			"db_cluster_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The ID of the cluster.",
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "Note information.",
			},
			"owner_account": schema.StringAttribute{
				Optional:    true,
				Description: "Owner Account",
			},
			"owner_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Owner ID",
			},
			"resource_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the resource group.",
			},
			"resource_owner_account": schema.StringAttribute{
				Optional:    true,
				Description: "Resource Owner Account",
			},
			"resource_owner_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Resource Owner ID",
			},
			"target_db_cluster_id": schema.StringAttribute{
				Optional:    true,
				Description: "The target instance name.",
			},
			"target_ip": schema.StringAttribute{
				Required:    true,
				Description: "The IP address of the target instance.",
			},
			"target_port": schema.StringAttribute{
				Optional:    true,
				Description: "The port of the target instance.",
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the VPC to which the access endpoint resides.",
			},
		},
	}
}

func (r *networkChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan networkChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "create", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"ClientToken": buildClientToken("CreateNetworkChannel"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.ChannelInfos.IsNull() && !plan.ChannelInfos.IsUnknown() {
		createValues["ChannelInfos"] = fwValueToInterface(plan.ChannelInfos, map[string]fwWireField{"channel_name": {Name: "ChannelName", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_port": {Name: "TargetPort", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "vpc_id": {Name: "VpcId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_db_cluster_id": {Name: "TargetDBClusterId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "db_cluster_id": {Name: "DBClusterId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_type": {Name: "TargetType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_ip": {Name: "TargetIp", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "region_id": {Name: "RegionId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "notes": {Name: "Notes", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.ChannelName.IsNull() && !plan.ChannelName.IsUnknown() {
		createValues["ChannelName"] = fwValueToInterface(plan.ChannelName, map[string]fwWireField{}, false)
	}
	if !plan.DBClusterId.IsNull() && !plan.DBClusterId.IsUnknown() {
		createValues["DBClusterId"] = fwValueToInterface(plan.DBClusterId, map[string]fwWireField{}, false)
	}
	if !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		createValues["Notes"] = fwValueToInterface(plan.Notes, map[string]fwWireField{}, false)
	}
	if !plan.OwnerAccount.IsNull() && !plan.OwnerAccount.IsUnknown() {
		createValues["OwnerAccount"] = fwValueToInterface(plan.OwnerAccount, map[string]fwWireField{}, false)
	}
	if !plan.OwnerId.IsNull() && !plan.OwnerId.IsUnknown() {
		createValues["OwnerId"] = fwValueToInterface(plan.OwnerId, map[string]fwWireField{}, false)
	}
	if !plan.ResourceGroupId.IsNull() && !plan.ResourceGroupId.IsUnknown() {
		createValues["ResourceGroupId"] = fwValueToInterface(plan.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !plan.ResourceOwnerAccount.IsNull() && !plan.ResourceOwnerAccount.IsUnknown() {
		createValues["ResourceOwnerAccount"] = fwValueToInterface(plan.ResourceOwnerAccount, map[string]fwWireField{}, false)
	}
	if !plan.ResourceOwnerId.IsNull() && !plan.ResourceOwnerId.IsUnknown() {
		createValues["ResourceOwnerId"] = fwValueToInterface(plan.ResourceOwnerId, map[string]fwWireField{}, false)
	}
	if !plan.TargetDBClusterId.IsNull() && !plan.TargetDBClusterId.IsUnknown() {
		createValues["TargetDBClusterId"] = fwValueToInterface(plan.TargetDBClusterId, map[string]fwWireField{}, false)
	}
	if !plan.TargetIp.IsNull() && !plan.TargetIp.IsUnknown() {
		createValues["TargetIp"] = fwValueToInterface(plan.TargetIp, map[string]fwWireField{}, false)
	}
	if !plan.TargetPort.IsNull() && !plan.TargetPort.IsUnknown() {
		createValues["TargetPort"] = fwValueToInterface(plan.TargetPort, map[string]fwWireField{}, false)
	}
	if !plan.VpcId.IsNull() && !plan.VpcId.IsUnknown() {
		createValues["VpcId"] = fwValueToInterface(plan.VpcId, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.TargetIp", Target: "TargetIp", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ChannelInfos[*].TargetPort", Target: "ChannelInfos[*].TargetPort", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceGroupId", Target: "ResourceGroupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ChannelInfos[*].DBClusterId", Target: "ChannelInfos[*].DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ChannelInfos[*].TargetDBClusterId", Target: "ChannelInfos[*].TargetDBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Notes", Target: "Notes", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ChannelInfos[*].VpcId", Target: "ChannelInfos[*].VpcId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ChannelInfos[*].ChannelName", Target: "ChannelInfos[*].ChannelName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceOwnerId", Target: "ResourceOwnerId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ChannelName", Target: "ChannelName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ChannelInfos[*].RegionId", Target: "ChannelInfos[*].RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceOwnerAccount", Target: "ResourceOwnerAccount", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.TargetDBClusterId", Target: "TargetDBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ChannelInfos[*].TargetIp", Target: "ChannelInfos[*].TargetIp", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.OwnerId", Target: "OwnerId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.VpcId", Target: "VpcId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ChannelInfos[*].Notes", Target: "ChannelInfos[*].Notes", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ChannelInfos[*].TargetType", Target: "ChannelInfos[*].TargetType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.TargetPort", Target: "TargetPort", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.OwnerAccount", Target: "OwnerAccount", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "CreateNetworkChannel", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateNetworkChannel", networkChannelTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateNetworkChannel", networkChannelTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(request, fwMappingPath("DBClusterId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_polardb_network_channel: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	identityValues1 := fwWalkMappedPath(request, fwMappingPath("ChannelName"), nil)
	if len(identityValues1) == 0 || identityValues1[0].Value == nil || fmt.Sprint(identityValues1[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_polardb_network_channel: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues1[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateNetworkChannel", networkChannelTypeName),
			fmt.Sprintf("The response carried no DBClusterId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.ChannelInfos = fwKnownValue(trackedPlan.ChannelInfos).(types.List)
	trackedPlan.ChannelName = fwKnownValue(trackedPlan.ChannelName).(types.String)
	trackedPlan.DBClusterId = fwKnownValue(trackedPlan.DBClusterId).(types.String)
	trackedPlan.Notes = fwKnownValue(trackedPlan.Notes).(types.String)
	trackedPlan.OwnerAccount = fwKnownValue(trackedPlan.OwnerAccount).(types.String)
	trackedPlan.OwnerId = fwKnownValue(trackedPlan.OwnerId).(types.Int64)
	trackedPlan.ResourceGroupId = fwKnownValue(trackedPlan.ResourceGroupId).(types.String)
	trackedPlan.ResourceOwnerAccount = fwKnownValue(trackedPlan.ResourceOwnerAccount).(types.String)
	trackedPlan.ResourceOwnerId = fwKnownValue(trackedPlan.ResourceOwnerId).(types.Int64)
	trackedPlan.TargetDBClusterId = fwKnownValue(trackedPlan.TargetDBClusterId).(types.String)
	trackedPlan.TargetIp = fwKnownValue(trackedPlan.TargetIp).(types.String)
	trackedPlan.TargetPort = fwKnownValue(trackedPlan.TargetPort).(types.String)
	trackedPlan.VpcId = fwKnownValue(trackedPlan.VpcId).(types.String)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := networkChannelGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", networkChannelTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", networkChannelTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", networkChannelTypeName),
			fmt.Sprintf("The resource %s is not visible via DescribeNetworkChannel right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, networkChannelState(object, &plan))...)
}

func (r *networkChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state networkChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := networkChannelGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNetworkChannel", networkChannelTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeNetworkChannel", networkChannelTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, networkChannelState(object, &state))...)
}

func (r *networkChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state networkChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *networkChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state networkChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "delete", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.ChannelInfos.IsNull() && !state.ChannelInfos.IsUnknown() {
		deleteValues["ChannelInfos"] = fwValueToInterface(state.ChannelInfos, map[string]fwWireField{"channel_name": {Name: "ChannelName", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_port": {Name: "TargetPort", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "vpc_id": {Name: "VpcId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_db_cluster_id": {Name: "TargetDBClusterId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "db_cluster_id": {Name: "DBClusterId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_type": {Name: "TargetType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_ip": {Name: "TargetIp", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "region_id": {Name: "RegionId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "notes": {Name: "Notes", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.ChannelName.IsNull() && !state.ChannelName.IsUnknown() {
		deleteValues["ChannelName"] = fwValueToInterface(state.ChannelName, map[string]fwWireField{}, false)
	}
	if !state.DBClusterId.IsNull() && !state.DBClusterId.IsUnknown() {
		deleteValues["DBClusterId"] = fwValueToInterface(state.DBClusterId, map[string]fwWireField{}, false)
	}
	if !state.Notes.IsNull() && !state.Notes.IsUnknown() {
		deleteValues["Notes"] = fwValueToInterface(state.Notes, map[string]fwWireField{}, false)
	}
	if !state.OwnerAccount.IsNull() && !state.OwnerAccount.IsUnknown() {
		deleteValues["OwnerAccount"] = fwValueToInterface(state.OwnerAccount, map[string]fwWireField{}, false)
	}
	if !state.OwnerId.IsNull() && !state.OwnerId.IsUnknown() {
		deleteValues["OwnerId"] = fwValueToInterface(state.OwnerId, map[string]fwWireField{}, false)
	}
	if !state.ResourceGroupId.IsNull() && !state.ResourceGroupId.IsUnknown() {
		deleteValues["ResourceGroupId"] = fwValueToInterface(state.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !state.ResourceOwnerAccount.IsNull() && !state.ResourceOwnerAccount.IsUnknown() {
		deleteValues["ResourceOwnerAccount"] = fwValueToInterface(state.ResourceOwnerAccount, map[string]fwWireField{}, false)
	}
	if !state.ResourceOwnerId.IsNull() && !state.ResourceOwnerId.IsUnknown() {
		deleteValues["ResourceOwnerId"] = fwValueToInterface(state.ResourceOwnerId, map[string]fwWireField{}, false)
	}
	if !state.TargetDBClusterId.IsNull() && !state.TargetDBClusterId.IsUnknown() {
		deleteValues["TargetDBClusterId"] = fwValueToInterface(state.TargetDBClusterId, map[string]fwWireField{}, false)
	}
	if !state.TargetIp.IsNull() && !state.TargetIp.IsUnknown() {
		deleteValues["TargetIp"] = fwValueToInterface(state.TargetIp, map[string]fwWireField{}, false)
	}
	if !state.TargetPort.IsNull() && !state.TargetPort.IsUnknown() {
		deleteValues["TargetPort"] = fwValueToInterface(state.TargetPort, map[string]fwWireField{}, false)
	}
	if !state.VpcId.IsNull() && !state.VpcId.IsUnknown() {
		deleteValues["VpcId"] = fwValueToInterface(state.VpcId, map[string]fwWireField{}, false)
	}
	identityParts := strings.Split(state.Id.ValueString(), ":")
	if len(identityParts) > 0 {
		deleteValues["DBClusterId"] = identityParts[0]
	}
	if len(identityParts) > 1 {
		deleteValues["ChannelName"] = identityParts[1]
	}
	needsRefresh := false
	if needsRefresh {
		object, err := networkChannelGet(r.Client(), state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Reading deletion conditions", err.Error())
			return
		}
		if object == nil {
			return
		}
		for key, value := range object {
			if value != nil {
				deleteValues[key] = value
			}
		}
	}
	selected := false
	deleteEligible0, deleteConditionErr0 := fwDeleteEligible(deleteValues, []fwDeleteCondition{})
	if deleteConditionErr0 != nil {
		resp.Diagnostics.AddError("Selecting deletion operation", deleteConditionErr0.Error())
		return
	}
	selected = selected || deleteEligible0
	if !selected {
		resp.Diagnostics.AddError("Selecting deletion operation", "No deletion operation matches the resource state; the resource remains tracked")
		return
	}
	if deleteEligible0 {
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteNetworkChannel")}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.ResourceOwnerId", Target: "ResourceOwnerId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ChannelName", Target: "ChannelName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ResourceOwnerAccount", Target: "ResourceOwnerAccount", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ResourceGroupId", Target: "ResourceGroupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.OwnerId", Target: "OwnerId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.VpcId", Target: "VpcId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.OwnerAccount", Target: "OwnerAccount", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DeleteNetworkChannel request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "DeleteNetworkChannel", request0, operationTimeout); err != nil && !isExpectedErrors(err, []string{"ResourceNotFound.NetworkChannel", "ResourceNotFound"}) {
			resp.Diagnostics.AddError("Deleting alicloud_polardb_network_channel: DeleteNetworkChannel", err.Error())
			return
		}
	}
}

func (r *networkChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// networkChannelState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func networkChannelState(object map[string]interface{}, prior *networkChannelResourceModel) networkChannelResourceModel {
	// Get has already merged and normalized every response mapping.
	str := func(key string) types.String {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.StringNull()
			}
			cur, ok = m[seg]
			if !ok {
				return types.StringNull()
			}
		}
		// gap#7c: present-but-empty strings must return ValueString(""), not Null.
		// v1 d.Set(field, "") treats empty as set (not absence); TestCheckResourceAttrSet
		// expects a non-null value for API-returned "" (e.g. DhcpOptionsSetId when VPC has
		// no DHCP set). The earlier && v != "" collapsed present-empty to Null, failing Set
		// checks. Absent fields still return Null via the !ok branches above.
		if v, ok := cur.(string); ok {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	// lookup traverses a dotted responseKey (e.g. "RegionStatuses.RegionStatus") segment by segment
	// from the Get response object, returning the raw value (array/map/scalar) or nil. gap#5: array
	// fields' responseKey is the array's true path from the read-tree [*] child, which may be a
	// two-level unwrap (object["RegionStatuses"]["RegionStatus"]) that a literal object[key] lookup
	// cannot reach. Single-segment keys (e.g. "IpWhiteList") traverse one step == object[key].
	lookup := func(key string) interface{} {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return nil
			}
			cur, ok = m[seg]
			if !ok {
				return nil
			}
		}
		return cur
	}
	// gap#11: composite-ID resource — build the id from all composite fields in the response,
	// joined by ":" (e.g. AccessGroupName:FileSystemType). The Get response carries every
	// part, so the state id round-trips through Read/Delete correctly (otherwise Read would
	// lose FileSystemType and DescribeAccessGroups would query the wrong namespace → 404).
	_cidParts := []string{
		str("DBClusterId").ValueString(),
		str("ChannelName").ValueString(),
	}
	m := networkChannelResourceModel{
		FwTimeouts:   types.ObjectNull(fwTimeoutTypes()),
		Id:           types.StringValue(strings.Join(_cidParts, ":")),
		ChannelInfos: networkChannelChannelInfosToTypeList(lookup("ChannelInfos")),
		ChannelName:  str("ChannelName"),
		DBClusterId:  str("DBClusterId"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		m.ChannelInfos = fwPreserveRequestValues(m.ChannelInfos, prior.ChannelInfos, map[string]fwWireField{"channel_name": {Name: "ChannelName", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_port": {Name: "TargetPort", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "vpc_id": {Name: "VpcId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_db_cluster_id": {Name: "TargetDBClusterId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "db_cluster_id": {Name: "DBClusterId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_type": {Name: "TargetType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_ip": {Name: "TargetIp", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "region_id": {Name: "RegionId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "notes": {Name: "Notes", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}).(types.List)
		m.Notes = prior.Notes
		m.OwnerAccount = prior.OwnerAccount
		m.OwnerId = prior.OwnerId
		m.ResourceGroupId = prior.ResourceGroupId
		m.ResourceOwnerAccount = prior.ResourceOwnerAccount
		m.ResourceOwnerId = prior.ResourceOwnerId
		m.TargetDBClusterId = prior.TargetDBClusterId
		m.TargetIp = prior.TargetIp
		m.TargetPort = prior.TargetPort
		m.VpcId = prior.VpcId
	}
	return m
}
