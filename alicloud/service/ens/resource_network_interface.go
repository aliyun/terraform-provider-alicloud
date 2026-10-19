package ens

import (
	"context"
	"encoding/json"
	"fmt"

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

const networkInterfaceTypeName = "alicloud_ens_network_interface"

var (
	_ resource.Resource                = &networkInterfaceResource{}
	_ resource.ResourceWithImportState = &networkInterfaceResource{}
)

func NewNetworkInterfaceResource() resource.Resource {
	return &networkInterfaceResource{}
}

type networkInterfaceResource struct {
	fwadapt.ResourceBase
}

// networkInterfaceModel mirrors the schema one-to-one. security_group_ids is
// the only collection; the API treats it as create-only, so it is Required and
// RequiresReplace. The remaining fields are either Optional create/update
// arguments or Computed values populated from DescribeNetworkInterfaces.
type networkInterfaceModel struct {
	Id                   types.String `tfsdk:"id"`
	SecurityGroupIds     types.List   `tfsdk:"security_group_ids"`
	Description          types.String `tfsdk:"description"`
	EnsRegionId          types.String `tfsdk:"ens_region_id"`
	InstanceId           types.String `tfsdk:"instance_id"`
	MacAddress           types.String `tfsdk:"mac_address"`
	NetworkId            types.String `tfsdk:"network_id"`
	NetworkInterfaceName types.String `tfsdk:"network_interface_name"`
	PrimaryIp            types.String `tfsdk:"primary_ip"`
	PrimaryIpType        types.String `tfsdk:"primary_ip_type"`
	Status               types.String `tfsdk:"status"`
	VswitchId            types.String `tfsdk:"vswitch_id"`
	VmncLearn            types.Bool   `tfsdk:"vmnc_learn"`
	CreateTime           types.String `tfsdk:"create_time"`
}

func (r *networkInterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = networkInterfaceTypeName
}

func (r *networkInterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a ENS Network Interface resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as the NetworkInterfaceId returned by CreateNetworkInterface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_group_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The IDs of the security groups. Create-only: changing this value forces replacement.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the ENI.",
			},
			"network_interface_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the ENI.",
			},
			"vswitch_id": schema.StringAttribute{
				Optional:    true,
				Description: "The vSwitch ID. Create-only: changing this value forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ens_region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The node ID of ENS.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the instance bound to the ENI.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mac_address": schema.StringAttribute{
				Computed:    true,
				Description: "The MAC address of the ENI.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_id": schema.StringAttribute{
				Computed:    true,
				Description: "The network ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"primary_ip": schema.StringAttribute{
				Computed:    true,
				Description: "The primary private network IP address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"primary_ip_type": schema.StringAttribute{
				Computed:    true,
				Description: "The primary IP address type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the network card.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// vmnc_learn is read-only: ENS does not expose a public action to
			// modify VmncLearn after creation. ModifyNetworkInterfaceAttribute
			// only accepts Description and NetworkInterfaceName, and the
			// ModifyNetworkInterfaceVmncLearn action is not present on the ENS
			// endpoint. The value is populated from DescribeNetworkInterfaces.
			"vmnc_learn": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to enable NIC route learning. Read-only.",
			},
		},
	}
}

func (r *networkInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan networkInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"SecurityGroupIds": securityGroupIDsJSON(plan.SecurityGroupIds),
	}
	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.NetworkInterfaceName.IsNull() && plan.NetworkInterfaceName.ValueString() != "" {
		// CreateNetworkInterface names the parameter "Name"; the update
		// action ModifyNetworkInterfaceAttribute uses "NetworkInterfaceName".
		request["Name"] = plan.NetworkInterfaceName.ValueString()
	}
	if !plan.VswitchId.IsNull() && plan.VswitchId.ValueString() != "" {
		request["VSwitchId"] = plan.VswitchId.ValueString()
	}

	response, err := r.Client().RpcPost(ensProduct, ensVersion, "CreateNetworkInterface", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateNetworkInterface", networkInterfaceTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateNetworkInterface", networkInterfaceTypeName).Error(),
		)
		return
	}

	id := createNetworkInterfaceId(response)
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateNetworkInterface", networkInterfaceTypeName),
			fmt.Sprintf("The response carried no NetworkInterfaceIds: %v", response),
		)
		return
	}

	object, err := describeNetworkInterface(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNetworkInterfaces", networkInterfaceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeNetworkInterfaces", networkInterfaceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateNetworkInterface", networkInterfaceTypeName),
			fmt.Sprintf("The ENI %s is not visible via DescribeNetworkInterfaces right after creation", id),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, networkInterfaceState(id, object))...)
}

func (r *networkInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state networkInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeNetworkInterface(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNetworkInterfaces", networkInterfaceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeNetworkInterfaces", networkInterfaceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, networkInterfaceState(state.Id.ValueString(), object))...)
}

func (r *networkInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state networkInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"NetworkInterfaceId": state.Id.ValueString(),
	}
	// ModifyNetworkInterfaceAttribute accepts only Description and
	// NetworkInterfaceName. A null plan value clears the attribute by sending
	// an empty string, matching the legacy SDKv2 update behavior.
	if !plan.Description.Equal(state.Description) {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.NetworkInterfaceName.Equal(state.NetworkInterfaceName) {
		request["NetworkInterfaceName"] = plan.NetworkInterfaceName.ValueString()
	}

	if len(request) > 1 {
		if _, err := r.Client().RpcPost(ensProduct, ensVersion, "ModifyNetworkInterfaceAttribute", nil, request, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyNetworkInterfaceAttribute", networkInterfaceTypeName),
				errs.WrapErrorf(err, "Updating %s: calling ModifyNetworkInterfaceAttribute", networkInterfaceTypeName).Error(),
			)
			return
		}
	}

	object, err := describeNetworkInterface(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeNetworkInterfaces", networkInterfaceTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeNetworkInterfaces", networkInterfaceTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, networkInterfaceState(state.Id.ValueString(), object))...)
}

func (r *networkInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state networkInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ENS RPC POST reads array parameters from the query string. Placing
	// NetworkInterfaceIds in the body made ENS report MissingNetworkInterfaceIds,
	// so the single id is sent as NetworkInterfaceIds.0 on the query.
	query := map[string]interface{}{
		"NetworkInterfaceIds.0": state.Id.ValueString(),
	}
	_, err := r.Client().RpcPost(ensProduct, ensVersion, "DeleteNetworkInterfaces", query, map[string]interface{}{}, true)
	if err != nil {
		if errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteNetworkInterfaces", networkInterfaceTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteNetworkInterfaces", networkInterfaceTypeName).Error(),
		)
		return
	}
}

func (r *networkInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// createNetworkInterfaceId extracts the first id from a CreateNetworkInterface
// response. ENS returns NetworkInterfaceIds as an array; the string form is also
// accepted for robustness.
func createNetworkInterfaceId(response map[string]interface{}) string {
	switch v := response["NetworkInterfaceIds"].(type) {
	case []interface{}:
		if len(v) > 0 {
			return fmt.Sprint(v[0])
		}
	case string:
		return v
	}
	return ""
}

// securityGroupIDsJSON serializes a planned security_group_ids list to the JSON
// array string ENS expects on CreateNetworkInterface. An empty list yields "[]".
func securityGroupIDsJSON(list types.List) string {
	elements := make([]string, 0, len(list.Elements()))
	for _, el := range list.Elements() {
		if s, ok := el.(types.String); ok {
			elements = append(elements, s.ValueString())
		}
	}
	raw, _ := json.Marshal(elements)
	return string(raw)
}

// networkInterfaceState maps a DescribeNetworkInterfaces element onto the
// resource model. Empty API strings become null so optional and computed
// fields read back absent rather than "".
func networkInterfaceState(id string, object map[string]interface{}) networkInterfaceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	ids := securityGroupIDs(object)
	securityGroupIds, diags := types.ListValueFrom(context.Background(), types.StringType, ids)
	if diags.HasError() {
		securityGroupIds = types.ListNull(types.StringType)
	}
	var vmnc types.Bool
	if v, ok := object["VmncLearn"].(bool); ok {
		vmnc = types.BoolValue(v)
	} else {
		vmnc = types.BoolNull()
	}
	return networkInterfaceModel{
		Id:                   types.StringValue(id),
		SecurityGroupIds:     securityGroupIds,
		Description:          str("Description"),
		EnsRegionId:          str("EnsRegionId"),
		InstanceId:           str("InstanceId"),
		MacAddress:           str("MacAddress"),
		NetworkId:            str("NetworkId"),
		NetworkInterfaceName: str("NetworkInterfaceName"),
		PrimaryIp:            str("PrimaryIp"),
		PrimaryIpType:        str("PrimaryIpType"),
		Status:               str("Status"),
		VswitchId:            str("VSwitchId"),
		VmncLearn:            vmnc,
		CreateTime:           str("CreationTime"),
	}
}
