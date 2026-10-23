// Package polardb. This file is generated automatically. Please do not modify it manually, thank you!
package polardb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const dbNodeTypeName = "alicloud_polardb_db_node"

var (
	_ resource.Resource                     = &dbNodeResource{}
	_ resource.ResourceWithImportState      = &dbNodeResource{}
	_ resource.ResourceWithConfigValidators = &dbNodeResource{}
)

func NewDbNodeResource() resource.Resource {
	return &dbNodeResource{}
}

func (r *dbNodeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	// ModifyDBNodeClass requires both parameters, so they are only usable together.
	return []resource.ConfigValidator{
		resourcevalidator.RequiredTogether(path.MatchRoot("db_node_target_class"), path.MatchRoot("modify_type")),
	}
}

type dbNodeResource struct {
	fwadapt.ResourceBase
}

type dbNodeResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	FwTimeouts              types.Object `tfsdk:"timeouts"`
	AutoUseCoupon           types.Bool   `tfsdk:"auto_use_coupon"`
	CloudProvider           types.String `tfsdk:"cloud_provider"`
	DBClusterId             types.String `tfsdk:"db_cluster_id"`
	DBNode                  types.List   `tfsdk:"db_node"`
	DBNodeIds               types.List   `tfsdk:"db_node_ids"`
	DBNodeTargetClass       types.String `tfsdk:"db_node_target_class"`
	DBNodeType              types.String `tfsdk:"db_node_type"`
	EndpointBindList        types.String `tfsdk:"endpoint_bind_list"`
	FromTimeService         types.String `tfsdk:"from_time_service"`
	ImciSwitch              types.String `tfsdk:"imci_switch"`
	ModifyType              types.String `tfsdk:"modify_type"`
	PlannedEndTime          types.String `tfsdk:"planned_end_time"`
	PlannedFlashingOffTime  types.String `tfsdk:"planned_flashing_off_time"`
	PlannedStartTime        types.String `tfsdk:"planned_start_time"`
	PromotionCode           types.String `tfsdk:"promotion_code"`
	ResourceGroupId         types.String `tfsdk:"resource_group_id"`
	RestartPlannedEndTime   types.String `tfsdk:"restart_planned_end_time"`
	RestartPlannedStartTime types.String `tfsdk:"restart_planned_start_time"`
	SccMode                 types.String `tfsdk:"scc_mode"`
	SubCategory             types.String `tfsdk:"sub_category"`
}

func (r *dbNodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = dbNodeTypeName
}

func (r *dbNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_use_coupon": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to use coupons automatically. Valid values: `true` and `false`. Default value: `true`.",
			},
			"cloud_provider": schema.StringAttribute{
				Optional:    true,
				Description: "The cloud service provider of the node.",
			},
			"db_cluster_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The ID of the cluster.",
			},
			// Protocol 5 cannot serve nested attributes; the element type carries the
			// same sub-attributes and state paths (db_node.0.target_class, ...).
			"db_node": schema.ListAttribute{
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"target_class": types.StringType,
						"zone_id":      types.StringType,
					},
				},
				Validators: []validator.List{
					// One node per resource: Read resolves exactly the tracked node, so a
					// longer list could never converge and would leak the extra nodes.
					listvalidator.SizeBetween(1, 1),
				},
				Required: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIf(func(_ context.Context, req planmodifier.ListRequest, resp *listplanmodifier.RequiresReplaceIfFuncResponse) {
						resp.RequiresReplace = fwConfiguredChanged(req.ConfigValue, req.PlanValue, req.StateValue)
					}, "Replace when configured nested values change", "Replace when configured nested values change"),
				},
				Description: "Add node information.",
			},
			"db_node_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				Description: "The DBNodeId list created.",
			},
			"db_node_target_class": schema.StringAttribute{
				Optional:    true,
				Description: "For the target specifications of all nodes, see [Compute Node Specifications](~~ 102542 ~~).",
			},
			"db_node_type": schema.StringAttribute{
				Optional:    true,
				Description: "Node type. Value range:\n- RO\n- STANDBY\n- DLNode",
			},
			"endpoint_bind_list": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The address ID of the cluster connection address to which the new node needs to join. If you need to add multiple connection addresses at the same time, separate multiple connection address IDs with commas (,).\n> * You can call the [DescribeDBClusterEndpoints](~~ 98205 ~~) interface to view the details of the cluster connection address, including the connection address ID.\n* You can enter the default cluster connection address and the address ID of the custom cluster connection address.\n* If this parameter is left blank, the new node is added to all cluster connection addresses that enable the **New node auto-add** function (that is, the value of 'AutoAddNewNodes' is 'Enable').",
			},
			"from_time_service": schema.StringAttribute{
				Optional:    true,
				Description: "Whether to execute the restart node task immediately or on schedule. Value range:\n\n- **false** (default): scheduled execution\n\n- **true**: immediate execution",
			},
			"imci_switch": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Whether to enable column storage index. Value range:\n\n-**ON**: ON\n\n-**OFF**: OFF",
			},
			"modify_type": schema.StringAttribute{
				Optional:    true,
				Description: "Change type, the range of values is as follows:\n* **Upgrade**: Upgrade specification\n* **Downgrade**: Downgrade specification",
			},
			"planned_end_time": schema.StringAttribute{
				Optional:    true,
				Description: "The latest time to start the target scheduled task. The format is' YYYY-MM-DDThh:mm:ssZ'(UTC).\n> * The latest time should be 30 minutes or more later than the start time.\n* If plannedstarttime' is set but this parameter is not specified, the latest time to execute the target task is 'start time +30 minutes' by default'. For example, when the' plannedstarttime' is set to' 2021-01-14 T09:00:00Z' and this parameter is left blank, the target task will start execution at the latest '2021-01-14 T09:30:00Z.",
			},
			"planned_flashing_off_time": schema.StringAttribute{
				Optional:    true,
				Description: "Planned flashing-off time. The time uses ISO 8601 format with UTC timezone (format: YYYY-MM-DDTHH:mm:ssZ).",
			},
			"planned_start_time": schema.StringAttribute{
				Optional:    true,
				Description: "Start execution timing (that is, the earliest time when the new node task is added within the target time period). The format is' YYYY-MM-DDThh:mm:ssZ'(UTC).\n> * The value range of the start time is any time point in the next 24 hours. For example, the current time is' 2021-01-14 T09:00:00Z', and the starting time range allowed here is' 2021-01-14 T09:00:00Z' ~ '2021-01-15 T09:00:00Z'.\n* If this parameter is left blank, the new node task will be executed immediately by default.",
			},
			"promotion_code": schema.StringAttribute{
				Optional:    true,
				Description: "The coupon code. If it is left empty, the default coupon is used.",
			},
			"resource_group_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The ID of the resource group",
			},
			"restart_planned_end_time": schema.StringAttribute{
				Optional:    true,
				Description: "The planned end time for the RestartDBNode operation. Format: YYYY-MM-DDThh:mm:ssZ (UTC). The latest time should be 30 minutes or more later than the start time.",
			},
			"restart_planned_start_time": schema.StringAttribute{
				Optional:    true,
				Description: "The planned start time for the RestartDBNode operation. Format: YYYY-MM-DDThh:mm:ssZ (UTC).",
			},
			"scc_mode": schema.StringAttribute{
				Validators: []validator.String{
					stringvalidator.OneOf("ON", "OFF"),
				},
				Optional:    true,
				Description: "Cluster strong consistent SCC",
			},
			"sub_category": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					// sub_category is only carried by ModifyDBNodeClass, which requires
					// the resize pair, so it cannot be set on its own.
					stringvalidator.AlsoRequires(path.MatchRoot("db_node_target_class"), path.MatchRoot("modify_type")),
					stringvalidator.OneOf("normal_exclusive", "normal_general"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The sub-series of the cluster. Valid values: `normal_exclusive` (exclusive) and `normal_general` (general).",
			},
		},
	}
}

func (r *dbNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dbNodeResourceModel
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
		"ClientToken": buildClientToken("CreateDBNodes"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.AutoUseCoupon.IsNull() && !plan.AutoUseCoupon.IsUnknown() {
		createValues["AutoUseCoupon"] = fwValueToInterface(plan.AutoUseCoupon, map[string]fwWireField{}, false)
	}
	if !plan.CloudProvider.IsNull() && !plan.CloudProvider.IsUnknown() {
		createValues["CloudProvider"] = fwValueToInterface(plan.CloudProvider, map[string]fwWireField{}, false)
	}
	if !plan.DBClusterId.IsNull() && !plan.DBClusterId.IsUnknown() {
		createValues["DBClusterId"] = fwValueToInterface(plan.DBClusterId, map[string]fwWireField{}, false)
	}
	if !plan.DBNode.IsNull() && !plan.DBNode.IsUnknown() {
		createValues["DBNode"] = fwValueToInterface(plan.DBNode, map[string]fwWireField{"zone_id": {Name: "ZoneId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_class": {Name: "TargetClass", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.DBNodeIds.IsNull() && !plan.DBNodeIds.IsUnknown() {
		createValues["DBNodeIds"] = fwValueToInterface(plan.DBNodeIds, map[string]fwWireField{}, false)
	}
	if !plan.DBNodeTargetClass.IsNull() && !plan.DBNodeTargetClass.IsUnknown() {
		createValues["DBNodeTargetClass"] = fwValueToInterface(plan.DBNodeTargetClass, map[string]fwWireField{}, false)
	}
	if !plan.DBNodeType.IsNull() && !plan.DBNodeType.IsUnknown() {
		createValues["DBNodeType"] = fwValueToInterface(plan.DBNodeType, map[string]fwWireField{}, false)
	}
	if !plan.EndpointBindList.IsNull() && !plan.EndpointBindList.IsUnknown() {
		createValues["EndpointBindList"] = fwValueToInterface(plan.EndpointBindList, map[string]fwWireField{}, false)
	}
	if !plan.FromTimeService.IsNull() && !plan.FromTimeService.IsUnknown() {
		createValues["FromTimeService"] = fwValueToInterface(plan.FromTimeService, map[string]fwWireField{}, false)
	}
	if !plan.ImciSwitch.IsNull() && !plan.ImciSwitch.IsUnknown() {
		createValues["ImciSwitch"] = fwValueToInterface(plan.ImciSwitch, map[string]fwWireField{}, false)
	}
	if !plan.ModifyType.IsNull() && !plan.ModifyType.IsUnknown() {
		createValues["ModifyType"] = fwValueToInterface(plan.ModifyType, map[string]fwWireField{}, false)
	}
	if !plan.PlannedEndTime.IsNull() && !plan.PlannedEndTime.IsUnknown() {
		createValues["PlannedEndTime"] = fwValueToInterface(plan.PlannedEndTime, map[string]fwWireField{}, false)
	}
	if !plan.PlannedFlashingOffTime.IsNull() && !plan.PlannedFlashingOffTime.IsUnknown() {
		createValues["PlannedFlashingOffTime"] = fwValueToInterface(plan.PlannedFlashingOffTime, map[string]fwWireField{}, false)
	}
	if !plan.PlannedStartTime.IsNull() && !plan.PlannedStartTime.IsUnknown() {
		createValues["PlannedStartTime"] = fwValueToInterface(plan.PlannedStartTime, map[string]fwWireField{}, false)
	}
	if !plan.PromotionCode.IsNull() && !plan.PromotionCode.IsUnknown() {
		createValues["PromotionCode"] = fwValueToInterface(plan.PromotionCode, map[string]fwWireField{}, false)
	}
	if !plan.ResourceGroupId.IsNull() && !plan.ResourceGroupId.IsUnknown() {
		createValues["ResourceGroupId"] = fwValueToInterface(plan.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !plan.RestartPlannedEndTime.IsNull() && !plan.RestartPlannedEndTime.IsUnknown() {
		createValues["RestartPlannedEndTime"] = fwValueToInterface(plan.RestartPlannedEndTime, map[string]fwWireField{}, false)
	}
	if !plan.RestartPlannedStartTime.IsNull() && !plan.RestartPlannedStartTime.IsUnknown() {
		createValues["RestartPlannedStartTime"] = fwValueToInterface(plan.RestartPlannedStartTime, map[string]fwWireField{}, false)
	}
	if !plan.SccMode.IsNull() && !plan.SccMode.IsUnknown() {
		createValues["SccMode"] = fwValueToInterface(plan.SccMode, map[string]fwWireField{}, false)
	}
	if !plan.SubCategory.IsNull() && !plan.SubCategory.IsUnknown() {
		createValues["SubCategory"] = fwValueToInterface(plan.SubCategory, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.AutoUseCoupon", Target: "AutoUseCoupon", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.CloudProvider", Target: "CloudProvider", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBNodeType", Target: "DBNodeType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBNode[*].TargetClass", Target: "DBNode[*].TargetClass", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.DBNode[*].ZoneId", Target: "DBNode[*].ZoneId", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.EndpointBindList", Target: "EndpointBindList", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ImciSwitch", Target: "ImciSwitch", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.PlannedEndTime", Target: "PlannedEndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.PlannedStartTime", Target: "PlannedStartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.PromotionCode", Target: "PromotionCode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ResourceGroupId", Target: "ResourceGroupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "CreateDBNodes", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDBNodes", dbNodeTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateDBNodes", dbNodeTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	// CreateDBNodes returns the allocated nodes as DBNodeIds.DBNodeId[*]; the resource
	// tracks the first node, whose id is the unique key declared by the model.
	identityValues0 := fwWalkMappedPath(response, fwMappingPath("$.DBNodeIds.DBNodeId[*]"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_polardb_db_node: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDBNodes", dbNodeTypeName),
			fmt.Sprintf("The response carried no DBNodeId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.AutoUseCoupon = fwKnownValue(trackedPlan.AutoUseCoupon).(types.Bool)
	trackedPlan.CloudProvider = fwKnownValue(trackedPlan.CloudProvider).(types.String)
	trackedPlan.DBClusterId = fwKnownValue(trackedPlan.DBClusterId).(types.String)
	trackedPlan.DBNode = fwKnownValue(trackedPlan.DBNode).(types.List)
	trackedPlan.DBNodeIds = fwKnownValue(trackedPlan.DBNodeIds).(types.List)
	trackedPlan.DBNodeTargetClass = fwKnownValue(trackedPlan.DBNodeTargetClass).(types.String)
	trackedPlan.DBNodeType = fwKnownValue(trackedPlan.DBNodeType).(types.String)
	trackedPlan.EndpointBindList = fwKnownValue(trackedPlan.EndpointBindList).(types.String)
	trackedPlan.FromTimeService = fwKnownValue(trackedPlan.FromTimeService).(types.String)
	trackedPlan.ImciSwitch = fwKnownValue(trackedPlan.ImciSwitch).(types.String)
	trackedPlan.ModifyType = fwKnownValue(trackedPlan.ModifyType).(types.String)
	trackedPlan.PlannedEndTime = fwKnownValue(trackedPlan.PlannedEndTime).(types.String)
	trackedPlan.PlannedFlashingOffTime = fwKnownValue(trackedPlan.PlannedFlashingOffTime).(types.String)
	trackedPlan.PlannedStartTime = fwKnownValue(trackedPlan.PlannedStartTime).(types.String)
	trackedPlan.PromotionCode = fwKnownValue(trackedPlan.PromotionCode).(types.String)
	trackedPlan.ResourceGroupId = fwKnownValue(trackedPlan.ResourceGroupId).(types.String)
	trackedPlan.RestartPlannedEndTime = fwKnownValue(trackedPlan.RestartPlannedEndTime).(types.String)
	trackedPlan.RestartPlannedStartTime = fwKnownValue(trackedPlan.RestartPlannedStartTime).(types.String)
	trackedPlan.SccMode = fwKnownValue(trackedPlan.SccMode).(types.String)
	trackedPlan.SubCategory = fwKnownValue(trackedPlan.SubCategory).(types.String)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := dbNodeGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", dbNodeTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", dbNodeTypeName, id).Error(),
		)
		return
	}
	// DescribeDBClusters reflects a freshly added node after a short delay, and
	// lists it before its creation task finishes: DeleteDBNodes and
	// ModifyDBNodeClass are rejected with OperationDenied.DBClusterStatus while
	// the cluster is still working. Wait for the cluster to settle so apply
	// reflects a node that later operations can act on.
	deadline := time.Now().Add(10 * time.Minute)
	for {
		if object != nil {
			if status, ok := object["DBClusterStatus"].(string); !ok || status == "Running" {
				break
			}
		}
		if time.Now().After(deadline) {
			if object == nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: after create", dbNodeTypeName),
					fmt.Sprintf("The resource %s is not visible via DescribeDBClusters right after creation", id),
				)
			} else {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: waiting for node readiness", dbNodeTypeName),
					fmt.Sprintf("The node %s remained in cluster status %v for 10 minutes", id, object["DBClusterStatus"]),
				)
			}
			return
		}
		time.Sleep(5 * time.Second)
		object, err = dbNodeGet(r.Client(), id)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: reading after create", dbNodeTypeName),
				errs.WrapErrorf(err, "Creating %s: reading after create (%s)", dbNodeTypeName, id).Error(),
			)
			return
		}
	}
	// CreateDBNodes carries every create-time parameter; RestartDBNode, ModifyDBNodeClass
	// and ModifyDBNodeSccMode are update-only and are driven from Update.
	resp.Diagnostics.Append(resp.State.Set(ctx, dbNodeState(object, &plan))...)
}

func (r *dbNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dbNodeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := dbNodeGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDBClusters", dbNodeTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeDBClusters", dbNodeTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, dbNodeState(object, &state))...)
}

func (r *dbNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state dbNodeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "update", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	planValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.AutoUseCoupon.IsNull() && !plan.AutoUseCoupon.IsUnknown() {
		planValues["AutoUseCoupon"] = fwValueToInterface(plan.AutoUseCoupon, map[string]fwWireField{}, false)
	}
	if !plan.CloudProvider.IsNull() && !plan.CloudProvider.IsUnknown() {
		planValues["CloudProvider"] = fwValueToInterface(plan.CloudProvider, map[string]fwWireField{}, false)
	}
	if !plan.DBClusterId.IsNull() && !plan.DBClusterId.IsUnknown() {
		planValues["DBClusterId"] = fwValueToInterface(plan.DBClusterId, map[string]fwWireField{}, false)
	}
	if !plan.DBNode.IsNull() && !plan.DBNode.IsUnknown() {
		planValues["DBNode"] = fwValueToInterface(plan.DBNode, map[string]fwWireField{"zone_id": {Name: "ZoneId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_class": {Name: "TargetClass", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.DBNodeIds.IsNull() && !plan.DBNodeIds.IsUnknown() {
		planValues["DBNodeIds"] = fwValueToInterface(plan.DBNodeIds, map[string]fwWireField{}, false)
	}
	if !plan.DBNodeTargetClass.IsNull() && !plan.DBNodeTargetClass.IsUnknown() {
		planValues["DBNodeTargetClass"] = fwValueToInterface(plan.DBNodeTargetClass, map[string]fwWireField{}, false)
	}
	if !plan.DBNodeType.IsNull() && !plan.DBNodeType.IsUnknown() {
		planValues["DBNodeType"] = fwValueToInterface(plan.DBNodeType, map[string]fwWireField{}, false)
	}
	if !plan.EndpointBindList.IsNull() && !plan.EndpointBindList.IsUnknown() {
		planValues["EndpointBindList"] = fwValueToInterface(plan.EndpointBindList, map[string]fwWireField{}, false)
	}
	if !plan.FromTimeService.IsNull() && !plan.FromTimeService.IsUnknown() {
		planValues["FromTimeService"] = fwValueToInterface(plan.FromTimeService, map[string]fwWireField{}, false)
	}
	if !plan.ImciSwitch.IsNull() && !plan.ImciSwitch.IsUnknown() {
		planValues["ImciSwitch"] = fwValueToInterface(plan.ImciSwitch, map[string]fwWireField{}, false)
	}
	if !plan.ModifyType.IsNull() && !plan.ModifyType.IsUnknown() {
		planValues["ModifyType"] = fwValueToInterface(plan.ModifyType, map[string]fwWireField{}, false)
	}
	if !plan.PlannedEndTime.IsNull() && !plan.PlannedEndTime.IsUnknown() {
		planValues["PlannedEndTime"] = fwValueToInterface(plan.PlannedEndTime, map[string]fwWireField{}, false)
	}
	if !plan.PlannedFlashingOffTime.IsNull() && !plan.PlannedFlashingOffTime.IsUnknown() {
		planValues["PlannedFlashingOffTime"] = fwValueToInterface(plan.PlannedFlashingOffTime, map[string]fwWireField{}, false)
	}
	if !plan.PlannedStartTime.IsNull() && !plan.PlannedStartTime.IsUnknown() {
		planValues["PlannedStartTime"] = fwValueToInterface(plan.PlannedStartTime, map[string]fwWireField{}, false)
	}
	if !plan.PromotionCode.IsNull() && !plan.PromotionCode.IsUnknown() {
		planValues["PromotionCode"] = fwValueToInterface(plan.PromotionCode, map[string]fwWireField{}, false)
	}
	if !plan.ResourceGroupId.IsNull() && !plan.ResourceGroupId.IsUnknown() {
		planValues["ResourceGroupId"] = fwValueToInterface(plan.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !plan.RestartPlannedEndTime.IsNull() && !plan.RestartPlannedEndTime.IsUnknown() {
		planValues["RestartPlannedEndTime"] = fwValueToInterface(plan.RestartPlannedEndTime, map[string]fwWireField{}, false)
	}
	if !plan.RestartPlannedStartTime.IsNull() && !plan.RestartPlannedStartTime.IsUnknown() {
		planValues["RestartPlannedStartTime"] = fwValueToInterface(plan.RestartPlannedStartTime, map[string]fwWireField{}, false)
	}
	if !plan.SccMode.IsNull() && !plan.SccMode.IsUnknown() {
		planValues["SccMode"] = fwValueToInterface(plan.SccMode, map[string]fwWireField{}, false)
	}
	if !plan.SubCategory.IsNull() && !plan.SubCategory.IsUnknown() {
		planValues["SubCategory"] = fwValueToInterface(plan.SubCategory, map[string]fwWireField{}, false)
	}
	stateValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.AutoUseCoupon.IsNull() && !state.AutoUseCoupon.IsUnknown() {
		stateValues["AutoUseCoupon"] = fwValueToInterface(state.AutoUseCoupon, map[string]fwWireField{}, false)
	}
	if !state.CloudProvider.IsNull() && !state.CloudProvider.IsUnknown() {
		stateValues["CloudProvider"] = fwValueToInterface(state.CloudProvider, map[string]fwWireField{}, false)
	}
	if !state.DBClusterId.IsNull() && !state.DBClusterId.IsUnknown() {
		stateValues["DBClusterId"] = fwValueToInterface(state.DBClusterId, map[string]fwWireField{}, false)
	}
	if !state.DBNode.IsNull() && !state.DBNode.IsUnknown() {
		stateValues["DBNode"] = fwValueToInterface(state.DBNode, map[string]fwWireField{"zone_id": {Name: "ZoneId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_class": {Name: "TargetClass", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.DBNodeIds.IsNull() && !state.DBNodeIds.IsUnknown() {
		stateValues["DBNodeIds"] = fwValueToInterface(state.DBNodeIds, map[string]fwWireField{}, false)
	}
	if !state.DBNodeTargetClass.IsNull() && !state.DBNodeTargetClass.IsUnknown() {
		stateValues["DBNodeTargetClass"] = fwValueToInterface(state.DBNodeTargetClass, map[string]fwWireField{}, false)
	}
	if !state.DBNodeType.IsNull() && !state.DBNodeType.IsUnknown() {
		stateValues["DBNodeType"] = fwValueToInterface(state.DBNodeType, map[string]fwWireField{}, false)
	}
	if !state.EndpointBindList.IsNull() && !state.EndpointBindList.IsUnknown() {
		stateValues["EndpointBindList"] = fwValueToInterface(state.EndpointBindList, map[string]fwWireField{}, false)
	}
	if !state.FromTimeService.IsNull() && !state.FromTimeService.IsUnknown() {
		stateValues["FromTimeService"] = fwValueToInterface(state.FromTimeService, map[string]fwWireField{}, false)
	}
	if !state.ImciSwitch.IsNull() && !state.ImciSwitch.IsUnknown() {
		stateValues["ImciSwitch"] = fwValueToInterface(state.ImciSwitch, map[string]fwWireField{}, false)
	}
	if !state.ModifyType.IsNull() && !state.ModifyType.IsUnknown() {
		stateValues["ModifyType"] = fwValueToInterface(state.ModifyType, map[string]fwWireField{}, false)
	}
	if !state.PlannedEndTime.IsNull() && !state.PlannedEndTime.IsUnknown() {
		stateValues["PlannedEndTime"] = fwValueToInterface(state.PlannedEndTime, map[string]fwWireField{}, false)
	}
	if !state.PlannedFlashingOffTime.IsNull() && !state.PlannedFlashingOffTime.IsUnknown() {
		stateValues["PlannedFlashingOffTime"] = fwValueToInterface(state.PlannedFlashingOffTime, map[string]fwWireField{}, false)
	}
	if !state.PlannedStartTime.IsNull() && !state.PlannedStartTime.IsUnknown() {
		stateValues["PlannedStartTime"] = fwValueToInterface(state.PlannedStartTime, map[string]fwWireField{}, false)
	}
	if !state.PromotionCode.IsNull() && !state.PromotionCode.IsUnknown() {
		stateValues["PromotionCode"] = fwValueToInterface(state.PromotionCode, map[string]fwWireField{}, false)
	}
	if !state.ResourceGroupId.IsNull() && !state.ResourceGroupId.IsUnknown() {
		stateValues["ResourceGroupId"] = fwValueToInterface(state.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !state.RestartPlannedEndTime.IsNull() && !state.RestartPlannedEndTime.IsUnknown() {
		stateValues["RestartPlannedEndTime"] = fwValueToInterface(state.RestartPlannedEndTime, map[string]fwWireField{}, false)
	}
	if !state.RestartPlannedStartTime.IsNull() && !state.RestartPlannedStartTime.IsUnknown() {
		stateValues["RestartPlannedStartTime"] = fwValueToInterface(state.RestartPlannedStartTime, map[string]fwWireField{}, false)
	}
	if !state.SccMode.IsNull() && !state.SccMode.IsUnknown() {
		stateValues["SccMode"] = fwValueToInterface(state.SccMode, map[string]fwWireField{}, false)
	}
	if !state.SubCategory.IsNull() && !state.SubCategory.IsUnknown() {
		stateValues["SubCategory"] = fwValueToInterface(state.SubCategory, map[string]fwWireField{}, false)
	}
	updated0 := false
	if fwPlannedChanged(plan.FromTimeService, state.FromTimeService) {
		updated0 = true
	}
	if fwPlannedChanged(plan.RestartPlannedEndTime, state.RestartPlannedEndTime) {
		updated0 = true
	}
	if fwPlannedChanged(plan.RestartPlannedStartTime, state.RestartPlannedStartTime) {
		updated0 = true
	}
	if updated0 {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"DBNodeId":    state.Id.ValueString(),
			"ClientToken": buildClientToken("RestartDBNode"),
		}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, operationValues0, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.DBNodeId", Target: "DBNodeId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.FromTimeService", Target: "FromTimeService", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RestartPlannedEndTime", Target: "PlannedEndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RestartPlannedStartTime", Target: "PlannedStartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues0, stateValues, []string{"$.DBNodeId"}), false); err != nil {
			resp.Diagnostics.AddError("Building RestartDBNode request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "RestartDBNode", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_polardb_db_node: RestartDBNode", err.Error())
			return
		}
	}
	// ModifyDBNodeClass is the resize operation: its other optional parameters
	// (coupons, provider, node type, schedule) only ride along when a resize runs.
	resizeRequested := false
	if fwPlannedChanged(plan.DBNodeTargetClass, state.DBNodeTargetClass) {
		resizeRequested = true
	}
	if fwPlannedChanged(plan.ModifyType, state.ModifyType) {
		resizeRequested = true
	}
	if fwPlannedChanged(plan.SubCategory, state.SubCategory) {
		resizeRequested = true
	}
	// Removing the pair is a configuration cleanup, not an API action, and both
	// parameters are mandatory whenever the operation is called.
	updated1 := resizeRequested &&
		!plan.DBNodeTargetClass.IsNull() && !plan.DBNodeTargetClass.IsUnknown() &&
		!plan.ModifyType.IsNull() && !plan.ModifyType.IsUnknown()
	if updated1 {
		operationValues1 := planValues
		request1 := map[string]interface{}{
			"DBNodeId":    state.Id.ValueString(),
			"ClientToken": buildClientToken("ModifyDBNodeClass"),
		}
		request1["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request1, operationValues1, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.AutoUseCoupon", Target: "AutoUseCoupon", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.CloudProvider", Target: "CloudProvider", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBNodeTargetClass", Target: "DBNodeTargetClass", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBNodeType", Target: "DBNodeType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ModifyType", Target: "ModifyType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.PlannedEndTime", Target: "PlannedEndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.PlannedFlashingOffTime", Target: "PlannedFlashingOffTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.PlannedStartTime", Target: "PlannedStartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.PromotionCode", Target: "PromotionCode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.SubCategory", Target: "SubCategory", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues1, stateValues, []string{"$.DBNodeId", "$.DBClusterId", "$.DBNodeTargetClass", "$.ModifyType"}), false); err != nil {
			resp.Diagnostics.AddError("Building ModifyDBNodeClass request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "ModifyDBNodeClass", request1, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_polardb_db_node: ModifyDBNodeClass", err.Error())
			return
		}
	}
	updated2 := false
	if fwPlannedChanged(plan.SccMode, state.SccMode) {
		updated2 = true
	}
	if updated2 {
		operationValues2 := planValues
		request2 := map[string]interface{}{
			"DBNodeId":    state.Id.ValueString(),
			"ClientToken": buildClientToken("ModifyDBNodeSccMode"),
		}
		request2["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request2, operationValues2, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBNodeId", Target: "DBNodeId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.SccMode", Target: "SccMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues2, stateValues, []string{"$.DBNodeId", "$.DBClusterId", "$.SccMode"}), false); err != nil {
			resp.Diagnostics.AddError("Building ModifyDBNodeSccMode request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "ModifyDBNodeSccMode", request2, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_polardb_db_node: ModifyDBNodeSccMode", err.Error())
			return
		}
	}
	// ModifyDBNodesClass is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	object, err := dbNodeGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDBClusters", dbNodeTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeDBClusters", dbNodeTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, dbNodeState(object, &plan))...)
}

func (r *dbNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dbNodeResourceModel
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
	if !state.AutoUseCoupon.IsNull() && !state.AutoUseCoupon.IsUnknown() {
		deleteValues["AutoUseCoupon"] = fwValueToInterface(state.AutoUseCoupon, map[string]fwWireField{}, false)
	}
	if !state.CloudProvider.IsNull() && !state.CloudProvider.IsUnknown() {
		deleteValues["CloudProvider"] = fwValueToInterface(state.CloudProvider, map[string]fwWireField{}, false)
	}
	if !state.DBClusterId.IsNull() && !state.DBClusterId.IsUnknown() {
		deleteValues["DBClusterId"] = fwValueToInterface(state.DBClusterId, map[string]fwWireField{}, false)
	}
	if !state.DBNode.IsNull() && !state.DBNode.IsUnknown() {
		deleteValues["DBNode"] = fwValueToInterface(state.DBNode, map[string]fwWireField{"zone_id": {Name: "ZoneId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "target_class": {Name: "TargetClass", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.DBNodeIds.IsNull() && !state.DBNodeIds.IsUnknown() {
		deleteValues["DBNodeIds"] = fwValueToInterface(state.DBNodeIds, map[string]fwWireField{}, false)
	}
	if !state.DBNodeTargetClass.IsNull() && !state.DBNodeTargetClass.IsUnknown() {
		deleteValues["DBNodeTargetClass"] = fwValueToInterface(state.DBNodeTargetClass, map[string]fwWireField{}, false)
	}
	if !state.DBNodeType.IsNull() && !state.DBNodeType.IsUnknown() {
		deleteValues["DBNodeType"] = fwValueToInterface(state.DBNodeType, map[string]fwWireField{}, false)
	}
	if !state.EndpointBindList.IsNull() && !state.EndpointBindList.IsUnknown() {
		deleteValues["EndpointBindList"] = fwValueToInterface(state.EndpointBindList, map[string]fwWireField{}, false)
	}
	if !state.FromTimeService.IsNull() && !state.FromTimeService.IsUnknown() {
		deleteValues["FromTimeService"] = fwValueToInterface(state.FromTimeService, map[string]fwWireField{}, false)
	}
	if !state.ImciSwitch.IsNull() && !state.ImciSwitch.IsUnknown() {
		deleteValues["ImciSwitch"] = fwValueToInterface(state.ImciSwitch, map[string]fwWireField{}, false)
	}
	if !state.ModifyType.IsNull() && !state.ModifyType.IsUnknown() {
		deleteValues["ModifyType"] = fwValueToInterface(state.ModifyType, map[string]fwWireField{}, false)
	}
	if !state.PlannedEndTime.IsNull() && !state.PlannedEndTime.IsUnknown() {
		deleteValues["PlannedEndTime"] = fwValueToInterface(state.PlannedEndTime, map[string]fwWireField{}, false)
	}
	if !state.PlannedFlashingOffTime.IsNull() && !state.PlannedFlashingOffTime.IsUnknown() {
		deleteValues["PlannedFlashingOffTime"] = fwValueToInterface(state.PlannedFlashingOffTime, map[string]fwWireField{}, false)
	}
	if !state.PlannedStartTime.IsNull() && !state.PlannedStartTime.IsUnknown() {
		deleteValues["PlannedStartTime"] = fwValueToInterface(state.PlannedStartTime, map[string]fwWireField{}, false)
	}
	if !state.PromotionCode.IsNull() && !state.PromotionCode.IsUnknown() {
		deleteValues["PromotionCode"] = fwValueToInterface(state.PromotionCode, map[string]fwWireField{}, false)
	}
	if !state.ResourceGroupId.IsNull() && !state.ResourceGroupId.IsUnknown() {
		deleteValues["ResourceGroupId"] = fwValueToInterface(state.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !state.RestartPlannedEndTime.IsNull() && !state.RestartPlannedEndTime.IsUnknown() {
		deleteValues["RestartPlannedEndTime"] = fwValueToInterface(state.RestartPlannedEndTime, map[string]fwWireField{}, false)
	}
	if !state.RestartPlannedStartTime.IsNull() && !state.RestartPlannedStartTime.IsUnknown() {
		deleteValues["RestartPlannedStartTime"] = fwValueToInterface(state.RestartPlannedStartTime, map[string]fwWireField{}, false)
	}
	if !state.SccMode.IsNull() && !state.SccMode.IsUnknown() {
		deleteValues["SccMode"] = fwValueToInterface(state.SccMode, map[string]fwWireField{}, false)
	}
	if !state.SubCategory.IsNull() && !state.SubCategory.IsUnknown() {
		deleteValues["SubCategory"] = fwValueToInterface(state.SubCategory, map[string]fwWireField{}, false)
	}
	// The resource identity is authoritative for deletion; db_node_ids is a
	// create-time output and may be absent from configuration.
	deleteValues["DBNodeIds"] = []interface{}{state.Id.ValueString()}
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
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteDBNodes")}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBNodeIds[*]", Target: "DBNodeId[*]", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DeleteDBNodes request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "DeleteDBNodes", request0, operationTimeout); err != nil && !isExpectedErrors(err, []string{"ResourceNotFound.DbNode", "ResourceNotFound"}) {
			resp.Diagnostics.AddError("Deleting alicloud_polardb_db_node: DeleteDBNodes", err.Error())
			return
		}
	}
}

func (r *dbNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// dbNodeState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func dbNodeState(object map[string]interface{}, prior *dbNodeResourceModel) dbNodeResourceModel {
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
	m := dbNodeResourceModel{
		FwTimeouts:      types.ObjectNull(fwTimeoutTypes()),
		Id:              str("DBNodeId"),
		DBClusterId:     str("DBClusterId"),
		DBNode:          dbNodeDBNodeToTypeList(lookup("DBNode")),
		DBNodeIds:       dbNodeTrackedIds(str("DBNodeId")),
		ImciSwitch:      str("ImciSwitch"),
		ResourceGroupId: str("ResourceGroupId"),
		SubCategory:     dbNodeSubCategoryState(str("SubCategory")),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		// db_node is the CreateDBNodes request: a resize through db_node_target_class
		// changes the observed class, so state keeps the configured request values to
		// avoid turning the next plan into a destructive replacement.
		m.DBNode = fwPreserveRequestValues(m.DBNode, prior.DBNode, map[string]fwWireField{"zone_id": {Name: "ZoneId", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "target_class": {Name: "TargetClass", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}}).(types.List)
		// Configured values win for parameters the read API only exposes cluster-wide.
		if !prior.ResourceGroupId.IsNull() && !prior.ResourceGroupId.IsUnknown() {
			m.ResourceGroupId = prior.ResourceGroupId
		}
		if !prior.SubCategory.IsNull() && !prior.SubCategory.IsUnknown() {
			m.SubCategory = prior.SubCategory
		}
		m.AutoUseCoupon = prior.AutoUseCoupon
		m.CloudProvider = prior.CloudProvider
		m.DBNodeTargetClass = prior.DBNodeTargetClass
		m.DBNodeType = prior.DBNodeType
		m.EndpointBindList = prior.EndpointBindList
		m.FromTimeService = prior.FromTimeService
		m.ModifyType = prior.ModifyType
		m.PlannedEndTime = prior.PlannedEndTime
		m.PlannedFlashingOffTime = prior.PlannedFlashingOffTime
		m.PlannedStartTime = prior.PlannedStartTime
		m.PromotionCode = prior.PromotionCode
		m.RestartPlannedEndTime = prior.RestartPlannedEndTime
		m.RestartPlannedStartTime = prior.RestartPlannedStartTime
		m.SccMode = prior.SccMode
	}
	return m
}

// dbNodeTrackedIds exposes the single node managed by this resource as the list
// shape CreateDBNodes returns.
func dbNodeTrackedIds(id types.String) types.List {
	if id.IsNull() || id.ValueString() == "" {
		return types.ListValueMust(types.StringType, nil)
	}
	return types.ListValueMust(types.StringType, []attr.Value{id})
}

// DescribeDBClusters reports the sub-series as Exclusive/General while the write APIs
// accept normal_exclusive/normal_general; state keeps the writable spelling so a
// configured value converges on the next plan.
func dbNodeSubCategoryState(value types.String) types.String {
	if value.IsNull() || value.IsUnknown() {
		return value
	}
	switch strings.ToLower(value.ValueString()) {
	case "exclusive":
		return types.StringValue("normal_exclusive")
	case "general":
		return types.StringValue("normal_general")
	}
	return value
}
