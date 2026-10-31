package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const rcClusterNodePoolTypeName = "alicloud_rds_rc_cluster_node_pool"

var (
	_ resource.Resource                = &rcClusterNodePoolResource{}
	_ resource.ResourceWithImportState = &rcClusterNodePoolResource{}
)

func NewRCClusterNodePoolResource() resource.Resource {
	return &rcClusterNodePoolResource{}
}

type rcClusterNodePoolResource struct {
	fwadapt.ResourceBase
}

type rcClusterNodePoolResourceModel struct {
	Id               types.String `tfsdk:"id"`
	ClusterId        types.String `tfsdk:"cluster_id"`
	NodePoolId       types.String `tfsdk:"node_pool_id"`
	Name             types.String `tfsdk:"name"`
	RegionId         types.String `tfsdk:"region_id"`
	Created          types.String `tfsdk:"created"`
	Updated          types.String `tfsdk:"updated"`
	ScalingGroup     types.List   `tfsdk:"scaling_group"`
	KubernetesConfig types.List   `tfsdk:"kubernetes_config"`
	Status           types.List   `tfsdk:"status"`
}

var (
	rcClusterNodePoolLabelAttrTypes = map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}

	rcClusterNodePoolTaintAttrTypes = map[string]attr.Type{
		"key":    types.StringType,
		"value":  types.StringType,
		"effect": types.StringType,
	}

	rcClusterNodePoolLabelType = types.ObjectType{AttrTypes: rcClusterNodePoolLabelAttrTypes}
	rcClusterNodePoolTaintType = types.ObjectType{AttrTypes: rcClusterNodePoolTaintAttrTypes}

	rcClusterNodePoolLabelsType = types.ListType{ElemType: rcClusterNodePoolLabelType}
	rcClusterNodePoolTaintsType = types.ListType{ElemType: rcClusterNodePoolTaintType}

	rcClusterNodePoolKubernetesConfigAttrTypes = map[string]attr.Type{
		"labels": rcClusterNodePoolLabelsType,
		"taints": rcClusterNodePoolTaintsType,
	}

	rcClusterNodePoolScalingGroupAttrTypes = map[string]attr.Type{
		"v_switch_ids":         types.ListType{ElemType: types.StringType},
		"instance_types":       types.ListType{ElemType: types.StringType},
		"instance_charge_type": types.StringType,
		"period":               types.Int64Type,
		"period_unit":          types.StringType,
		"auto_renew":           types.BoolType,
		"security_group_ids":   types.ListType{ElemType: types.StringType},
		"desired_size":         types.Int64Type,
	}

	rcClusterNodePoolStatusAttrTypes = map[string]attr.Type{
		"state":                     types.StringType,
		"total_nodes":               types.Int64Type,
		"retry_times":               types.Int64Type,
		"scaling_error_description": types.StringType,
	}

	rcClusterNodePoolKubernetesConfigType = types.ObjectType{AttrTypes: rcClusterNodePoolKubernetesConfigAttrTypes}
	rcClusterNodePoolScalingGroupType     = types.ObjectType{AttrTypes: rcClusterNodePoolScalingGroupAttrTypes}
	rcClusterNodePoolStatusType           = types.ObjectType{AttrTypes: rcClusterNodePoolStatusAttrTypes}
)

func (r *rcClusterNodePoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = rcClusterNodePoolTypeName
}

func (r *rcClusterNodePoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a RDS Custom (RC) cluster node pool resource. A node pool scales the ACK worker nodes that carry RDS Custom instances of a RDS Custom container cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, in the format `<cluster_id>:<node_pool_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the RDS Custom container cluster (the imported ACK cluster) that owns the node pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"node_pool_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the node pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the node pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the node pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the node pool was created.",
			},
			"updated": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the node pool was last updated.",
			},
			"scaling_group": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: rcClusterNodePoolScalingGroupType,
				Description: "The scaling group configuration of the node pool. Exactly one element is expected.",
			},
			"kubernetes_config": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: rcClusterNodePoolKubernetesConfigType,
				Description: "The Kubernetes configuration (node labels and taints) of the node pool. Exactly one element is expected.",
			},
			"status": schema.ListAttribute{
				Computed:    true,
				ElementType: rcClusterNodePoolStatusType,
				Description: "The status of the node pool. Exactly one element is returned.",
			},
		},
	}
}

func (r *rcClusterNodePoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rcClusterNodePoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := plan.ClusterId.ValueString()
	request := map[string]interface{}{
		"RegionId":  r.Client().RegionId,
		"ClusterId": clusterId,
	}

	nodePoolInfo, err := expandRCClusterNodePoolNodePoolInfo(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s: building NodePoolInfo", rcClusterNodePoolTypeName), err.Error())
		return
	}
	if nodePoolInfo != "" {
		request["NodePoolInfo"] = nodePoolInfo
	}

	scalingGroup, diags := expandRCClusterNodePoolScalingGroup(ctx, plan.ScalingGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if scalingGroup != "" {
		request["ScalingGroup"] = scalingGroup
	}

	kubernetesConfig, diags := expandRCClusterNodePoolKubernetesConfig(ctx, plan.KubernetesConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if kubernetesConfig != "" {
		request["KubernetesConfig"] = kubernetesConfig
	}

	response, err := r.Client().RpcPost(rcClusterNodePoolRdsProduct, rcClusterNodePoolRdsVersion, "CreateRCClusterNodePool", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateRCClusterNodePool", rcClusterNodePoolTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateRCClusterNodePool", rcClusterNodePoolTypeName).Error(),
		)
		return
	}

	nodePoolId, _ := response["NodePoolId"].(string)
	if nodePoolId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateRCClusterNodePool", rcClusterNodePoolTypeName),
			fmt.Sprintf("The response carried no NodePoolId: %v", response),
		)
		return
	}

	object, err := waitRCClusterNodePoolVisible(ctx, r.Client(), clusterId, nodePoolId, rcClusterNodePoolCreateWindow)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for node pool %s", rcClusterNodePoolTypeName, nodePoolId),
			errs.WrapErrorf(err, "Creating %s: waiting for node pool %s", rcClusterNodePoolTypeName, nodePoolId).Error(),
		)
		return
	}

	state, diags := rcClusterNodePoolState(ctx, object, clusterId, nodePoolId, r.Client().RegionId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.KubernetesConfig = rcClusterNodePoolApplyConfigShape(plan.KubernetesConfig, state.KubernetesConfig)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *rcClusterNodePoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rcClusterNodePoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, nodePoolId, diags := rcClusterNodePoolIdentity(state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeRCClusterNodePool(r.Client(), clusterId, nodePoolId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeRCClusterNodePoolDetail", rcClusterNodePoolTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeRCClusterNodePoolDetail", rcClusterNodePoolTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState, diags := rcClusterNodePoolState(ctx, object, clusterId, nodePoolId, r.Client().RegionId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState.KubernetesConfig = rcClusterNodePoolApplyConfigShape(state.KubernetesConfig, newState.KubernetesConfig)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *rcClusterNodePoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state rcClusterNodePoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, nodePoolId, diags := rcClusterNodePoolIdentity(state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"RegionId":   r.Client().RegionId,
		"ClusterId":  clusterId,
		"NodePoolId": nodePoolId,
	}

	nodePoolInfo, err := expandRCClusterNodePoolNodePoolInfo(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Updating %s: building NodePoolInfo", rcClusterNodePoolTypeName), err.Error())
		return
	}
	if nodePoolInfo != "" {
		request["NodePoolInfo"] = nodePoolInfo
	}

	scalingGroup, diags := expandRCClusterNodePoolScalingGroup(ctx, plan.ScalingGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if scalingGroup != "" {
		request["ScalingGroup"] = scalingGroup
	}

	kubernetesConfig, diags := expandRCClusterNodePoolKubernetesConfig(ctx, plan.KubernetesConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if kubernetesConfig != "" {
		request["KubernetesConfig"] = kubernetesConfig
	}

	if _, err := r.Client().RpcPost(rcClusterNodePoolRdsProduct, rcClusterNodePoolRdsVersion, "ModifyRCClusterNodePool", nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling ModifyRCClusterNodePool", rcClusterNodePoolTypeName),
			errs.WrapErrorf(err, "Updating %s: calling ModifyRCClusterNodePool", rcClusterNodePoolTypeName).Error(),
		)
		return
	}

	object, err := waitRCClusterNodePoolVisible(ctx, r.Client(), clusterId, nodePoolId, rcClusterNodePoolUpdateWindow)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: waiting for node pool %s", rcClusterNodePoolTypeName, nodePoolId),
			errs.WrapErrorf(err, "Updating %s: waiting for node pool %s", rcClusterNodePoolTypeName, nodePoolId).Error(),
		)
		return
	}

	newState, diags := rcClusterNodePoolState(ctx, object, clusterId, nodePoolId, r.Client().RegionId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState.KubernetesConfig = rcClusterNodePoolApplyConfigShape(plan.KubernetesConfig, newState.KubernetesConfig)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *rcClusterNodePoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state rcClusterNodePoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, nodePoolId, diags := rcClusterNodePoolIdentity(state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Client().RpcPost(rcClusterNodePoolRdsProduct, rcClusterNodePoolRdsVersion, "DeleteRCClusterNodePool", nil, map[string]interface{}{
		"RegionId":   r.Client().RegionId,
		"ClusterId":  clusterId,
		"NodePoolId": nodePoolId,
	}, true)
	if err != nil && !isRCClusterNodePoolNotFound(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteRCClusterNodePool", rcClusterNodePoolTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteRCClusterNodePool", rcClusterNodePoolTypeName).Error(),
		)
		return
	}

	if err := waitRCClusterNodePoolDeleted(ctx, r.Client(), clusterId, nodePoolId, rcClusterNodePoolDeleteWindow); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for node pool %s", rcClusterNodePoolTypeName, nodePoolId),
			errs.WrapErrorf(err, "Deleting %s: waiting for node pool %s", rcClusterNodePoolTypeName, nodePoolId).Error(),
		)
		return
	}
}

func (r *rcClusterNodePoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// rcClusterNodePoolIdentity resolves the cluster and node pool IDs of a
// resource, falling back to parsing the resource ID for imported state.
func rcClusterNodePoolIdentity(state rcClusterNodePoolResourceModel) (string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	clusterId := state.ClusterId.ValueString()
	nodePoolId := state.NodePoolId.ValueString()
	if clusterId != "" && nodePoolId != "" {
		return clusterId, nodePoolId, diags
	}

	parts := strings.Split(state.Id.ValueString(), ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		diags.AddError(
			fmt.Sprintf("Reading %s: invalid resource ID", rcClusterNodePoolTypeName),
			fmt.Sprintf("The resource ID %q is not in the format <cluster_id>:<node_pool_id>", state.Id.ValueString()),
		)
		return "", "", diags
	}
	return parts[0], parts[1], diags
}

// rcClusterNodePoolState maps a DescribeRCClusterNodePoolDetail response onto
// the resource model.
func rcClusterNodePoolState(ctx context.Context, object map[string]interface{}, clusterId, nodePoolId, defaultRegionId string) (rcClusterNodePoolResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	info, _ := object["NodePoolInfo"].(map[string]interface{})

	regionId := rcClusterNodePoolString(info["RegionId"])
	if regionId.IsNull() && defaultRegionId != "" {
		regionId = types.StringValue(defaultRegionId)
	}

	state := rcClusterNodePoolResourceModel{
		Id:               types.StringValue(clusterId + ":" + nodePoolId),
		ClusterId:        types.StringValue(clusterId),
		NodePoolId:       types.StringValue(nodePoolId),
		Name:             rcClusterNodePoolString(info["Name"]),
		RegionId:         regionId,
		Created:          rcClusterNodePoolString(info["Created"]),
		Updated:          rcClusterNodePoolString(info["Updated"]),
		ScalingGroup:     rcClusterNodePoolScalingGroupValue(object["ScalingGroup"]),
		KubernetesConfig: rcClusterNodePoolKubernetesConfigValue(object["KubernetesConfig"]),
		Status:           rcClusterNodePoolStatusValue(object["Status"]),
	}
	return state, diags
}

func expandRCClusterNodePoolScalingGroup(ctx context.Context, list types.List) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return "", diags
	}
	elements := list.Elements()
	if len(elements) == 0 {
		return "", diags
	}
	object, ok := elements[0].(types.Object)
	if !ok {
		diags.AddError("Invalid scaling_group", "The scaling_group value is not an object")
		return "", diags
	}
	attributes := object.Attributes()

	request := map[string]interface{}{}
	if values, diags2 := expandRCClusterNodePoolStringSlice(ctx, attributes["v_switch_ids"]); len(values) > 0 {
		diags.Append(diags2...)
		request["VSwitchIds"] = values
	}
	if values, diags2 := expandRCClusterNodePoolStringSlice(ctx, attributes["instance_types"]); len(values) > 0 {
		diags.Append(diags2...)
		request["InstanceTypes"] = values
	}
	if v, ok := attributes["instance_charge_type"]; ok {
		if s := v.(types.String); !s.IsNull() && !s.IsUnknown() {
			request["InstanceChargeType"] = s.ValueString()
		}
	}
	if v, ok := attributes["period"]; ok {
		if p := v.(types.Int64); !p.IsNull() && !p.IsUnknown() {
			request["Period"] = p.ValueInt64()
		}
	}
	if v, ok := attributes["period_unit"]; ok {
		if s := v.(types.String); !s.IsNull() && !s.IsUnknown() {
			request["PeriodUnit"] = s.ValueString()
		}
	}
	if v, ok := attributes["auto_renew"]; ok {
		if b := v.(types.Bool); !b.IsNull() && !b.IsUnknown() {
			request["AutoRenew"] = b.ValueBool()
		}
	}
	if values, diags2 := expandRCClusterNodePoolStringSlice(ctx, attributes["security_group_ids"]); len(values) > 0 {
		diags.Append(diags2...)
		request["SecurityGroupIds"] = values
	}
	if v, ok := attributes["desired_size"]; ok {
		if s := v.(types.Int64); !s.IsNull() && !s.IsUnknown() {
			request["DesiredSize"] = s.ValueInt64()
		}
	}
	if diags.HasError() {
		return "", diags
	}
	if len(request) == 0 {
		return "", diags
	}
	payload, err := json.Marshal(request)
	if err != nil {
		diags.AddError("Invalid scaling_group", err.Error())
		return "", diags
	}
	return string(payload), diags
}

func expandRCClusterNodePoolKubernetesConfig(ctx context.Context, list types.List) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return "", diags
	}
	elements := list.Elements()
	if len(elements) == 0 {
		return "", diags
	}
	object, ok := elements[0].(types.Object)
	if !ok {
		diags.AddError("Invalid kubernetes_config", "The kubernetes_config value is not an object")
		return "", diags
	}
	attributes := object.Attributes()

	request := map[string]interface{}{}
	if list, ok := attributes["labels"].(types.List); ok && !list.IsNull() && !list.IsUnknown() {
		values := make([]map[string]interface{}, 0, len(list.Elements()))
		for _, element := range list.Elements() {
			label, ok := element.(types.Object)
			if !ok {
				continue
			}
			labelAttributes := label.Attributes()
			values = append(values, map[string]interface{}{
				"Key":   rcClusterNodePoolPlanString(labelAttributes["key"]),
				"Value": rcClusterNodePoolPlanString(labelAttributes["value"]),
			})
		}
		if len(values) > 0 {
			request["Labels"] = values
		}
	}
	if list, ok := attributes["taints"].(types.List); ok && !list.IsNull() && !list.IsUnknown() {
		values := make([]map[string]interface{}, 0, len(list.Elements()))
		for _, element := range list.Elements() {
			taint, ok := element.(types.Object)
			if !ok {
				continue
			}
			taintAttributes := taint.Attributes()
			values = append(values, map[string]interface{}{
				"Key":    rcClusterNodePoolPlanString(taintAttributes["key"]),
				"Value":  rcClusterNodePoolPlanString(taintAttributes["value"]),
				"Effect": rcClusterNodePoolPlanString(taintAttributes["effect"]),
			})
		}
		if len(values) > 0 {
			request["Taints"] = values
		}
	}
	if diags.HasError() {
		return "", diags
	}
	if len(request) == 0 {
		return "", diags
	}
	payload, err := json.Marshal(request)
	if err != nil {
		diags.AddError("Invalid kubernetes_config", err.Error())
		return "", diags
	}
	return string(payload), diags
}

// expandRCClusterNodePoolStringSlice reads an optional list-of-string value.
func expandRCClusterNodePoolStringSlice(ctx context.Context, value attr.Value) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	list, ok := value.(types.List)
	if !ok || list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var values []string
	diags.Append(list.ElementsAs(ctx, &values, false)...)
	return values, diags
}

// rcClusterNodePoolPlanString returns the string of a planned attribute, or ""
// when it is null or unknown.
func rcClusterNodePoolPlanString(value attr.Value) string {
	stringValue, ok := value.(types.String)
	if !ok || stringValue.IsNull() || stringValue.IsUnknown() {
		return ""
	}
	return stringValue.ValueString()
}

func rcClusterNodePoolScalingGroupValue(raw interface{}) types.List {
	object, ok := raw.(map[string]interface{})
	if !ok || len(object) == 0 {
		return types.ListNull(rcClusterNodePoolScalingGroupType)
	}
	value := types.ObjectValueMust(rcClusterNodePoolScalingGroupAttrTypes, map[string]attr.Value{
		"v_switch_ids":         rcClusterNodePoolStringList(object["VSwitchIds"]),
		"instance_types":       rcClusterNodePoolStringList(object["InstanceTypes"]),
		"instance_charge_type": rcClusterNodePoolString(object["InstanceChargeType"]),
		"period":               rcClusterNodePoolInt64(object["Period"]),
		"period_unit":          rcClusterNodePoolString(object["PeriodUnit"]),
		"auto_renew":           rcClusterNodePoolBool(object["AutoRenew"]),
		"security_group_ids":   rcClusterNodePoolStringList(object["SecurityGroupIds"]),
		"desired_size":         rcClusterNodePoolInt64(object["DesiredSize"]),
	})
	return types.ListValueMust(rcClusterNodePoolScalingGroupType, []attr.Value{value})
}

func rcClusterNodePoolKubernetesConfigValue(raw interface{}) types.List {
	object, ok := raw.(map[string]interface{})
	if !ok || len(object) == 0 {
		return types.ListNull(rcClusterNodePoolKubernetesConfigType)
	}
	value := types.ObjectValueMust(rcClusterNodePoolKubernetesConfigAttrTypes, map[string]attr.Value{
		"labels": rcClusterNodePoolLabelList(object["Labels"]),
		"taints": rcClusterNodePoolTaintList(object["Taints"]),
	})
	return types.ListValueMust(rcClusterNodePoolKubernetesConfigType, []attr.Value{value})
}

// rcClusterNodePoolApplyConfigShape restores the null-versus-empty-list shape of
// the labels and taints that the configuration declared. DescribeRCClusterNodePoolDetail
// omits collections that are empty, so the response alone cannot distinguish
// `taints = []` from an absent taints argument; reporting null for a planned
// empty list makes Terraform reject the apply with "produced an unexpected new
// value".
func rcClusterNodePoolApplyConfigShape(prior, observed types.List) types.List {
	if observed.IsNull() || observed.IsUnknown() || len(observed.Elements()) == 0 {
		return observed
	}
	observedObject, ok := observed.Elements()[0].(types.Object)
	if !ok {
		return observed
	}
	attributes := make(map[string]attr.Value, len(rcClusterNodePoolKubernetesConfigAttrTypes))
	for name, value := range observedObject.Attributes() {
		attributes[name] = value
	}
	var priorObject types.Object
	hasPrior := !prior.IsNull() && !prior.IsUnknown() && len(prior.Elements()) > 0
	if hasPrior {
		priorObject, hasPrior = prior.Elements()[0].(types.Object)
	}
	for name, elementType := range map[string]attr.Type{
		"labels": rcClusterNodePoolLabelType,
		"taints": rcClusterNodePoolTaintType,
	} {
		observedList, ok := attributes[name].(types.List)
		if !ok || !observedList.IsNull() {
			continue
		}
		if hasPrior {
			if priorList, ok := priorObject.Attributes()[name].(types.List); ok && !priorList.IsUnknown() {
				attributes[name] = priorList
				continue
			}
		}
		attributes[name] = types.ListValueMust(elementType, []attr.Value{})
	}
	return types.ListValueMust(rcClusterNodePoolKubernetesConfigType, []attr.Value{
		types.ObjectValueMust(rcClusterNodePoolKubernetesConfigAttrTypes, attributes),
	})
}

func rcClusterNodePoolStatusValue(raw interface{}) types.List {
	object, ok := raw.(map[string]interface{})
	if !ok || len(object) == 0 {
		return types.ListNull(rcClusterNodePoolStatusType)
	}
	value := types.ObjectValueMust(rcClusterNodePoolStatusAttrTypes, map[string]attr.Value{
		"state":                     rcClusterNodePoolString(object["State"]),
		"total_nodes":               rcClusterNodePoolInt64(object["TotalNodes"]),
		"retry_times":               rcClusterNodePoolInt64(object["RetryTimes"]),
		"scaling_error_description": rcClusterNodePoolString(object["ScalingErrorDescription"]),
	})
	return types.ListValueMust(rcClusterNodePoolStatusType, []attr.Value{value})
}

func rcClusterNodePoolLabelList(raw interface{}) types.List {
	items, ok := raw.([]interface{})
	if !ok || len(items) == 0 {
		return types.ListNull(rcClusterNodePoolLabelType)
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		object, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		values = append(values, types.ObjectValueMust(rcClusterNodePoolLabelAttrTypes, map[string]attr.Value{
			"key":   rcClusterNodePoolString(object["Key"]),
			"value": rcClusterNodePoolString(object["Value"]),
		}))
	}
	if len(values) == 0 {
		return types.ListNull(rcClusterNodePoolLabelType)
	}
	return types.ListValueMust(rcClusterNodePoolLabelType, values)
}

func rcClusterNodePoolTaintList(raw interface{}) types.List {
	items, ok := raw.([]interface{})
	if !ok || len(items) == 0 {
		return types.ListNull(rcClusterNodePoolTaintType)
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		object, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		values = append(values, types.ObjectValueMust(rcClusterNodePoolTaintAttrTypes, map[string]attr.Value{
			"key":    rcClusterNodePoolString(object["Key"]),
			"value":  rcClusterNodePoolString(object["Value"]),
			"effect": rcClusterNodePoolString(object["Effect"]),
		}))
	}
	if len(values) == 0 {
		return types.ListNull(rcClusterNodePoolTaintType)
	}
	return types.ListValueMust(rcClusterNodePoolTaintType, values)
}

func rcClusterNodePoolStringList(raw interface{}) types.List {
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

func rcClusterNodePoolString(raw interface{}) types.String {
	value, ok := raw.(string)
	if !ok || value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

func rcClusterNodePoolInt64(raw interface{}) types.Int64 {
	switch value := raw.(type) {
	case json.Number:
		if parsed, err := value.Int64(); err == nil {
			return types.Int64Value(parsed)
		}
	case float64:
		return types.Int64Value(int64(value))
	case float32:
		return types.Int64Value(int64(value))
	case int:
		return types.Int64Value(int64(value))
	case int32:
		return types.Int64Value(int64(value))
	case int64:
		return types.Int64Value(value)
	case string:
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return types.Int64Value(parsed)
		}
	}
	return types.Int64Null()
}

func rcClusterNodePoolBool(raw interface{}) types.Bool {
	switch value := raw.(type) {
	case bool:
		return types.BoolValue(value)
	case string:
		if parsed, err := strconv.ParseBool(value); err == nil {
			return types.BoolValue(parsed)
		}
	}
	return types.BoolNull()
}
