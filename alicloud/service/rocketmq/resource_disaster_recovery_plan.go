package rocketmq

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const disasterRecoveryPlanTypeName = "alicloud_rocketmq_disaster_recovery_plan"

var (
	_ resource.Resource                = (*disasterRecoveryPlanResource)(nil)
	_ resource.ResourceWithImportState = (*disasterRecoveryPlanResource)(nil)
)

func NewDisasterRecoveryPlanResource() resource.Resource {
	return &disasterRecoveryPlanResource{}
}

type disasterRecoveryPlanResource struct {
	fwadapt.ResourceBase
}

// instanceModel is the element type of the instances nested block.
type instanceModel struct {
	AuthType        types.String `tfsdk:"auth_type"`
	ConsumerGroupId types.String `tfsdk:"consumer_group_id"`
	EndpointUrl     types.String `tfsdk:"endpoint_url"`
	InstanceId      types.String `tfsdk:"instance_id"`
	InstanceRole    types.String `tfsdk:"instance_role"`
	InstanceType    types.String `tfsdk:"instance_type"`
	MessageProperty types.List   `tfsdk:"message_property"`
	NetworkType     types.String `tfsdk:"network_type"`
	Password        types.String `tfsdk:"password"`
	RegionId        types.String `tfsdk:"region_id"`
	SecurityGroupId types.String `tfsdk:"security_group_id"`
	Username        types.String `tfsdk:"username"`
	VpcId           types.String `tfsdk:"vpc_id"`
	VswitchId       types.String `tfsdk:"vswitch_id"`
}

type messagePropertyModel struct {
	PropertyKey   types.String `tfsdk:"property_key"`
	PropertyValue types.String `tfsdk:"property_value"`
}

type disasterRecoveryPlanModel struct {
	Id                    types.String `tfsdk:"id"`
	AutoSyncCheckpoint    types.Bool   `tfsdk:"auto_sync_checkpoint"`
	CreateTime            types.String `tfsdk:"create_time"`
	Instances             types.List   `tfsdk:"instances"`
	PlanDesc              types.String `tfsdk:"plan_desc"`
	PlanId                types.String `tfsdk:"plan_id"`
	PlanName              types.String `tfsdk:"plan_name"`
	PlanType              types.String `tfsdk:"plan_type"`
	RegionId              types.String `tfsdk:"region_id"`
	Status                types.String `tfsdk:"status"`
	SyncCheckpointEnabled types.Bool   `tfsdk:"sync_checkpoint_enabled"`
	UpdateTime            types.String `tfsdk:"update_time"`
}

func (r *disasterRecoveryPlanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = disasterRecoveryPlanTypeName
}

func (r *disasterRecoveryPlanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a RocketMQ Disaster Recovery Plan resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, same as `plan_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_sync_checkpoint": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"create_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plan_desc": schema.StringAttribute{Optional: true},
			"plan_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plan_name": schema.StringAttribute{Optional: true},
			"plan_type": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sync_checkpoint_enabled": schema.BoolAttribute{Optional: true, Computed: true},
			"update_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instances": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"auth_type":         schema.StringAttribute{Optional: true},
						"consumer_group_id": schema.StringAttribute{Optional: true, Computed: true},
						"endpoint_url": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"instance_id": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"instance_role": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"instance_type": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"message_property": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"property_key":   schema.StringAttribute{Optional: true},
									"property_value": schema.StringAttribute{Optional: true},
								},
							},
						},
						"network_type": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"password": schema.StringAttribute{Optional: true, Sensitive: true},
						"region_id": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"security_group_id": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"username": schema.StringAttribute{Optional: true},
						"vpc_id": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"vswitch_id": schema.StringAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
					},
				},
			},
		},
	}
}

// expandInstances builds the API request body slice from the planned instances.
// Empty strings are omitted so the server applies its own defaults, matching the
// SDKv2 expand semantics; consumerGroupId is included whenever set.
func expandInstances(ctx context.Context, instances types.List) ([]map[string]interface{}, error) {
	var models []instanceModel
	if diags := instances.ElementsAs(ctx, &models, false); diags.HasError() {
		return nil, fmt.Errorf("expanding instances: %s", diags.Errors())
	}
	out := make([]map[string]interface{}, 0, len(models))
	for _, m := range models {
		item := make(map[string]interface{})
		setIf := func(key string, v types.String) {
			if !v.IsNull() && v.ValueString() != "" {
				item[key] = v.ValueString()
			}
		}
		setIf("instanceType", m.InstanceType)
		setIf("instanceRole", m.InstanceRole)
		setIf("securityGroupId", m.SecurityGroupId)
		setIf("password", m.Password)
		setIf("instanceId", m.InstanceId)
		setIf("consumerGroupId", m.ConsumerGroupId)
		setIf("regionId", m.RegionId)
		setIf("vSwitchId", m.VswitchId)
		setIf("endpointUrl", m.EndpointUrl)
		setIf("vpcId", m.VpcId)
		setIf("authType", m.AuthType)
		setIf("networkType", m.NetworkType)
		setIf("username", m.Username)
		var mp []messagePropertyModel
		if mpDiags := m.MessageProperty.ElementsAs(ctx, &mp, false); !mpDiags.HasError() && len(mp) > 0 {
			messageProperty := make(map[string]interface{})
			if pk := mp[0].PropertyKey; !pk.IsNull() && pk.ValueString() != "" {
				messageProperty["propertyKey"] = pk.ValueString()
			}
			if pv := mp[0].PropertyValue; !pv.IsNull() && pv.ValueString() != "" {
				messageProperty["propertyValue"] = pv.ValueString()
			}
			if len(messageProperty) > 0 {
				item["messageProperty"] = messageProperty
			}
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *disasterRecoveryPlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan disasterRecoveryPlanModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{}
	if v := plan.PlanName; !v.IsNull() && v.ValueString() != "" {
		body["planName"] = v.ValueString()
	}
	if v := plan.PlanDesc; !v.IsNull() && v.ValueString() != "" {
		body["planDesc"] = v.ValueString()
	}
	if v := plan.PlanType; !v.IsNull() && v.ValueString() != "" {
		body["planType"] = v.ValueString()
	}
	if !plan.AutoSyncCheckpoint.IsNull() {
		body["autoSyncCheckpoint"] = plan.AutoSyncCheckpoint.ValueBool()
	}
	if !plan.SyncCheckpointEnabled.IsNull() {
		body["syncCheckpointEnabled"] = plan.SyncCheckpointEnabled.ValueBool()
	}
	if !plan.Instances.IsNull() && len(plan.Instances.Elements()) > 0 {
		expanded, err := expandInstances(ctx, plan.Instances)
		if err != nil {
			resp.Diagnostics.AddError("Expanding instances for CreateDisasterRecoveryPlan", err.Error())
			return
		}
		body["instances"] = expanded
	}

	// Create retries on Instance.NotFound: the DisasterRecovery backend
	// occasionally 404s a freshly-referenced instance before its ledger
	// settles. Bounded so a unit-test mock returns on the first call.
	client := r.Client()
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	for attempt := 0; attempt < createMaxAttempts; attempt++ {
		response, err = client.RoaPost(rocketmqProduct, rocketmqVersion, "/disaster_recovery", map[string]*string{}, nil, body, true)
		if err == nil {
			break
		}
		if isPlanNotFound(err) || strings.Contains(err.Error(), "Instance.NotFound") {
			wait()
			continue
		}
		break
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDisasterRecoveryPlan", disasterRecoveryPlanTypeName),
			err.Error(),
		)
		return
	}

	planId := fmt.Sprint(response["data"])
	plan.Id = types.StringValue(planId)
	plan.PlanId = types.StringValue(planId)

	if err := waitForPlanStatus(client, planId, []string{planStatusCreated, planStatusRunning}, createSettleTimeout, 4*time.Minute); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s: waiting for plan status", disasterRecoveryPlanTypeName), err.Error())
		return
	}

	model, diags := readPlanIntoModel(ctx, client, planId, plan.Instances)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *disasterRecoveryPlanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state disasterRecoveryPlanModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	model, diags := readPlanIntoModel(ctx, r.Client(), state.Id.ValueString(), state.Instances)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if model == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *disasterRecoveryPlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan disasterRecoveryPlanModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state disasterRecoveryPlanModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{}
	if v := plan.PlanName; !v.IsNull() && v.ValueString() != "" {
		body["planName"] = v.ValueString()
	}
	if v := plan.PlanDesc; !v.IsNull() && v.ValueString() != "" {
		body["planDesc"] = v.ValueString()
	}
	if !plan.AutoSyncCheckpoint.IsNull() {
		body["autoSyncCheckpoint"] = plan.AutoSyncCheckpoint.ValueBool()
	}
	if !plan.SyncCheckpointEnabled.IsNull() {
		body["syncCheckpointEnabled"] = plan.SyncCheckpointEnabled.ValueBool()
	}
	if !plan.Instances.IsNull() && len(plan.Instances.Elements()) > 0 {
		expanded, err := expandInstances(ctx, plan.Instances)
		if err != nil {
			resp.Diagnostics.AddError("Expanding instances for UpdateDisasterRecoveryPlan", err.Error())
			return
		}
		body["instances"] = expanded
	}

	action := fmt.Sprintf("/disaster_recovery/%s", state.Id.ValueString())
	if _, err := r.Client().RoaPatch(rocketmqProduct, rocketmqVersion, action, map[string]*string{}, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateDisasterRecoveryPlan", disasterRecoveryPlanTypeName),
			err.Error(),
		)
		return
	}

	model, diags := readPlanIntoModel(ctx, r.Client(), state.Id.ValueString(), plan.Instances)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *disasterRecoveryPlanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state disasterRecoveryPlanModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	action := fmt.Sprintf("/disaster_recovery/%s", state.Id.ValueString())
	if _, err := r.Client().RoaDelete(rocketmqProduct, rocketmqVersion, action, map[string]*string{}, nil, nil, true); err != nil && !isPlanNotFound(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDisasterRecoveryPlan", disasterRecoveryPlanTypeName),
			err.Error(),
		)
		return
	}
	if err := waitForPlanGone(r.Client(), state.Id.ValueString(), deleteSettleTimeout, 5*time.Minute); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Deleting %s: waiting for plan removal", disasterRecoveryPlanTypeName), err.Error())
		return
	}
}

func (r *disasterRecoveryPlanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// readPlanIntoModel fetches the plan and maps it onto the resource model,
// carrying over write-only fields the API never returns (password,
// endpointUrl, networkType, vpcId, vswitchId) from the prior state. A nil
// model signals the plan is gone.
func readPlanIntoModel(ctx context.Context, client *connectivity.AliyunClient, planId string, priorInstances types.List) (*disasterRecoveryPlanModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	object, err := describeDisasterRecoveryPlan(client, planId)
	if err != nil {
		diags.AddError(fmt.Sprintf("Reading %s: calling GetDisasterRecoveryPlan", disasterRecoveryPlanTypeName), err.Error())
		return nil, diags
	}
	if object == nil {
		return nil, diags
	}
	var prior []instanceModel
	_ = priorInstances.ElementsAs(ctx, &prior, false)
	instances := flattenInstances(object["instances"], prior)
	instancesList, d := types.ListValueFrom(ctx, instancesElementType(), instances)
	diags.Append(d...)
	return &disasterRecoveryPlanModel{
		Id:                    types.StringValue(planId),
		AutoSyncCheckpoint:    boolFrom(object["autoSyncCheckpoint"]),
		CreateTime:            strFrom(object["createTime"]),
		Instances:             instancesList,
		PlanDesc:              strFrom(object["planDesc"]),
		PlanId:                types.StringValue(planId),
		PlanName:              strFrom(object["planName"]),
		PlanType:              strFrom(object["planType"]),
		RegionId:              strFrom(client.RegionId),
		Status:                strFrom(object["planStatus"]),
		SyncCheckpointEnabled: boolFrom(object["syncCheckpointEnabled"]),
		UpdateTime:            strFrom(object["updateTime"]),
	}, diags
}
