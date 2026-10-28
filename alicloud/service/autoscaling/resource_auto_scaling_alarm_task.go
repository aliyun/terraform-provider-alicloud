// Package autoscaling. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package autoscaling

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const alarmTaskTypeName = "alicloud_auto_scaling_alarm_task"

var (
	_ resource.Resource                = &alarmTaskResource{}
	_ resource.ResourceWithImportState = &alarmTaskResource{}
)

func NewAlarmTaskResource() resource.Resource {
	return &alarmTaskResource{}
}

type alarmTaskResource struct {
	fwadapt.ResourceBase
}

type alarmTaskResourceModel struct {
	Id                       types.String  `tfsdk:"id"`
	FwTimeouts               types.Object  `tfsdk:"timeouts"`
	AlarmActions             types.List    `tfsdk:"alarm_actions"`
	AlarmTaskName            types.String  `tfsdk:"alarm_task_name"`
	AlarmTriggerState        types.String  `tfsdk:"alarm_trigger_state"`
	ComparisonOperator       types.String  `tfsdk:"comparison_operator"`
	Description              types.String  `tfsdk:"description"`
	Dimension                types.List    `tfsdk:"dimension"`
	Effective                types.String  `tfsdk:"effective"`
	EvaluationCount          types.Int64   `tfsdk:"evaluation_count"`
	Expression               types.List    `tfsdk:"expression"`
	ExpressionsLogicOperator types.String  `tfsdk:"expressions_logic_operator"`
	MetricName               types.String  `tfsdk:"metric_name"`
	MetricType               types.String  `tfsdk:"metric_type"`
	Period                   types.Int64   `tfsdk:"period"`
	ScalingGroupId           types.String  `tfsdk:"scaling_group_id"`
	Statistics               types.String  `tfsdk:"statistics"`
	Status                   types.String  `tfsdk:"status"`
	Threshold                types.Float64 `tfsdk:"threshold"`
}

func (r *alarmTaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = alarmTaskTypeName
}

func (r *alarmTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": fwTimeoutBlock(),
			// dimension/expression are ListNestedBlock (not ListNestedAttribute) so the
			// SDKv2+framework mux can serve them under protocol v5: NestedAttributeObject
			// requires protocol 6, which this provider does not advertise.
			"dimension": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"dimension_value": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The dimension information value associated with the monitoring item.",
						},
						"dimension_key": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The Dimension Information key associated with the monitoring item.",
						},
					},
				},
				Description: "The list of dimension information associated with the monitoring item.",
			},
			"expression": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"metric_name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Metric name of the multi-metric alarm rule",
						},
						"comparison_operator": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The comparison between the metric statistics value and the threshold value of the multi-indicator alarm rule.",
						},
						"period": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "The period of statistical monitoring item data in multi-indicator alarm rules.",
						},
						"statistics": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The method of statistical monitoring item data in multi-indicator alarm rules.",
						},
						"threshold": schema.Float64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "The threshold value of the multi-indicator alarm rule.",
						},
					},
				},
				Description: "List of monitoring indicators for multi-indicator alarm rules",
			},
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"alarm_actions": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Unique identifier of the alarm task Association scaling rule",
			},
			"alarm_task_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the resource",
			},
			"alarm_trigger_state": schema.StringAttribute{
				Computed:    true,
				Description: "The trigger status of the alarm task. Possible values:\n\nALARM: ALARM, the ALARM condition has been met.\nOK: normal, the alarm condition has not been met.\nINSUFFICIENT_DATA: insufficient data to determine whether the alarm condition is met.",
			},
			"comparison_operator": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Comparison between the statistical value of the monitoring item and the threshold.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the alarm task",
			},
			"effective": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Specify the effective time period of the alarm task",
			},
			"evaluation_count": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of times the threshold expression must be met to trigger the execution of a scaling rule.",
			},
			"expressions_logic_operator": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Relationship between multi-indicator alarm conditions",
			},
			"metric_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Monitoring Item Name",
			},
			"metric_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Monitoring Item Type",
			},
			"period": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Cycle of statistical monitoring item data",
			},
			"scaling_group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The ID of the scaling group associated with the alarm task.",
			},
			"statistics": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Methods for statistical monitoring item data",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The status of the resource",
			},
			"threshold": schema.Float64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The threshold value at which the alarm is triggered.",
			},
		},
	}
}

func (r *alarmTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alarmTaskResourceModel
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
		"ClientToken": buildClientToken("CreateAlarm"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.AlarmActions.IsNull() && !plan.AlarmActions.IsUnknown() {
		createValues["AlarmActions"] = fwValueToInterface(plan.AlarmActions, map[string]fwWireField{}, false)
	}
	if !plan.AlarmTaskName.IsNull() && !plan.AlarmTaskName.IsUnknown() {
		createValues["AlarmTaskName"] = fwValueToInterface(plan.AlarmTaskName, map[string]fwWireField{}, false)
	}
	if !plan.AlarmTriggerState.IsNull() && !plan.AlarmTriggerState.IsUnknown() {
		createValues["AlarmTriggerState"] = fwValueToInterface(plan.AlarmTriggerState, map[string]fwWireField{}, false)
	}
	if !plan.ComparisonOperator.IsNull() && !plan.ComparisonOperator.IsUnknown() {
		createValues["ComparisonOperator"] = fwValueToInterface(plan.ComparisonOperator, map[string]fwWireField{}, false)
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createValues["Description"] = fwValueToInterface(plan.Description, map[string]fwWireField{}, false)
	}
	if !plan.Dimension.IsNull() && !plan.Dimension.IsUnknown() {
		createValues["Dimension"] = fwValueToInterface(plan.Dimension, map[string]fwWireField{"dimension_value": {Name: "DimensionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dimension_key": {Name: "DimensionKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.Effective.IsNull() && !plan.Effective.IsUnknown() {
		createValues["Effective"] = fwValueToInterface(plan.Effective, map[string]fwWireField{}, false)
	}
	if !plan.EvaluationCount.IsNull() && !plan.EvaluationCount.IsUnknown() {
		createValues["EvaluationCount"] = fwValueToInterface(plan.EvaluationCount, map[string]fwWireField{}, false)
	}
	if !plan.Expression.IsNull() && !plan.Expression.IsUnknown() {
		createValues["Expression"] = fwValueToInterface(plan.Expression, map[string]fwWireField{"metric_name": {Name: "MetricName", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "comparison_operator": {Name: "ComparisonOperator", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "period": {Name: "Period", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "statistics": {Name: "Statistics", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "threshold": {Name: "Threshold", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.ExpressionsLogicOperator.IsNull() && !plan.ExpressionsLogicOperator.IsUnknown() {
		createValues["ExpressionsLogicOperator"] = fwValueToInterface(plan.ExpressionsLogicOperator, map[string]fwWireField{}, false)
	}
	if !plan.MetricName.IsNull() && !plan.MetricName.IsUnknown() {
		createValues["MetricName"] = fwValueToInterface(plan.MetricName, map[string]fwWireField{}, false)
	}
	if !plan.MetricType.IsNull() && !plan.MetricType.IsUnknown() {
		createValues["MetricType"] = fwValueToInterface(plan.MetricType, map[string]fwWireField{}, false)
	}
	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		createValues["Period"] = fwValueToInterface(plan.Period, map[string]fwWireField{}, false)
	}
	if !plan.ScalingGroupId.IsNull() && !plan.ScalingGroupId.IsUnknown() {
		createValues["ScalingGroupId"] = fwValueToInterface(plan.ScalingGroupId, map[string]fwWireField{}, false)
	}
	if !plan.Statistics.IsNull() && !plan.Statistics.IsUnknown() {
		createValues["Statistics"] = fwValueToInterface(plan.Statistics, map[string]fwWireField{}, false)
	}
	if !plan.Threshold.IsNull() && !plan.Threshold.IsUnknown() {
		createValues["Threshold"] = fwValueToInterface(plan.Threshold, map[string]fwWireField{}, false)
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		createValues["Status"] = fwValueToInterface(plan.Status, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Description", Target: "Description", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ScalingGroupId", Target: "ScalingGroupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		// Multi-metric alarms send Expressions[] via the whole-object remap below;
		// single-metric sends top-level MetricName/Period/Statistics/Threshold/
		// ComparisonOperator via the $.X -> X mappings. The $.Expression[*].X -> X
		// extracts were a generator dual-cross-mapping defect: they scatter
		// multi-metric expression fields onto single-metric top-level params
		// (uniform-scalar conflicts on differing thresholds) and are no-ops for
		// single-metric (no Expression block). v1 omits top-level metric fields for
		// multi-metric too, so Expressions-only is the correct request.
		{Source: "$.MetricType", Target: "MetricType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.EvaluationCount", Target: "EvaluationCount", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Effective", Target: "Effective", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AlarmActions[*]", Target: "AlarmActions[*]", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.Dimension[*].DimensionKey", Target: "Dimensions[*].DimensionKey", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.Dimension[*].DimensionValue", Target: "Dimensions[*].DimensionValue", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		// Single-metric alarms send top-level MetricName/Period/Statistics/Threshold/
		// ComparisonOperator (CreateAlarm 2022-02-22 accepts both top-level and
		// Expressions[]). Broadcasting $.X to Expressions[*].X produced a malformed
		// one-element Expressions the API rejects ("Expressions is not valid") for
		// single-metric configs, so these map top-level scalars to top-level params;
		// multi-metric configs have no top-level metric fields, making these no-ops,
		// and Expressions is built by the whole-object $.Expression[*] remap below.
		{Source: "$.MetricName", Target: "MetricName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ComparisonOperator", Target: "ComparisonOperator", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Period", Target: "Period", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Statistics", Target: "Statistics", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Threshold", Target: "Threshold", Encoding: "uniform-scalar|number", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ExpressionsLogicOperator", Target: "ExpressionsLogicOperator", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AlarmTaskName", Target: "Name", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Expression[*]", Target: "Expressions[*]", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(r.Client(), "Ess", "2022-02-22", "CreateAlarm", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateAlarm", alarmTaskTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateAlarm", alarmTaskTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(response, fwMappingPath("$.AlarmTaskId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_auto_scaling_alarm_task: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateAlarm", alarmTaskTypeName),
			fmt.Sprintf("The response carried no AlarmTaskId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.AlarmActions = fwKnownValue(trackedPlan.AlarmActions).(types.List)
	trackedPlan.AlarmTaskName = fwKnownValue(trackedPlan.AlarmTaskName).(types.String)
	trackedPlan.AlarmTriggerState = fwKnownValue(trackedPlan.AlarmTriggerState).(types.String)
	trackedPlan.ComparisonOperator = fwKnownValue(trackedPlan.ComparisonOperator).(types.String)
	trackedPlan.Description = fwKnownValue(trackedPlan.Description).(types.String)
	trackedPlan.Dimension = fwKnownValue(trackedPlan.Dimension).(types.List)
	trackedPlan.Effective = fwKnownValue(trackedPlan.Effective).(types.String)
	trackedPlan.EvaluationCount = fwKnownValue(trackedPlan.EvaluationCount).(types.Int64)
	trackedPlan.Expression = fwKnownValue(trackedPlan.Expression).(types.List)
	trackedPlan.ExpressionsLogicOperator = fwKnownValue(trackedPlan.ExpressionsLogicOperator).(types.String)
	trackedPlan.MetricName = fwKnownValue(trackedPlan.MetricName).(types.String)
	trackedPlan.MetricType = fwKnownValue(trackedPlan.MetricType).(types.String)
	trackedPlan.Period = fwKnownValue(trackedPlan.Period).(types.Int64)
	trackedPlan.ScalingGroupId = fwKnownValue(trackedPlan.ScalingGroupId).(types.String)
	trackedPlan.Statistics = fwKnownValue(trackedPlan.Statistics).(types.String)
	trackedPlan.Threshold = fwKnownValue(trackedPlan.Threshold).(types.Float64)
	trackedPlan.Status = fwKnownValue(trackedPlan.Status).(types.String)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := AlarmTaskGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", alarmTaskTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", alarmTaskTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", alarmTaskTypeName),
			fmt.Sprintf("The resource %s is not visible via DescribeAlarms right after creation", id),
		)
		return
	}
	// Apply configurable properties represented only by secondary update APIs.
	state := alarmTaskState(object, &plan)
	state.Id = types.StringValue(id)
	plan.Id = types.StringValue(id)
	planValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.AlarmActions.IsNull() && !plan.AlarmActions.IsUnknown() {
		planValues["AlarmActions"] = fwValueToInterface(plan.AlarmActions, map[string]fwWireField{}, false)
	}
	if !plan.AlarmTaskName.IsNull() && !plan.AlarmTaskName.IsUnknown() {
		planValues["AlarmTaskName"] = fwValueToInterface(plan.AlarmTaskName, map[string]fwWireField{}, false)
	}
	if !plan.AlarmTriggerState.IsNull() && !plan.AlarmTriggerState.IsUnknown() {
		planValues["AlarmTriggerState"] = fwValueToInterface(plan.AlarmTriggerState, map[string]fwWireField{}, false)
	}
	if !plan.ComparisonOperator.IsNull() && !plan.ComparisonOperator.IsUnknown() {
		planValues["ComparisonOperator"] = fwValueToInterface(plan.ComparisonOperator, map[string]fwWireField{}, false)
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		planValues["Description"] = fwValueToInterface(plan.Description, map[string]fwWireField{}, false)
	}
	if !plan.Dimension.IsNull() && !plan.Dimension.IsUnknown() {
		planValues["Dimension"] = fwValueToInterface(plan.Dimension, map[string]fwWireField{"dimension_value": {Name: "DimensionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dimension_key": {Name: "DimensionKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.Effective.IsNull() && !plan.Effective.IsUnknown() {
		planValues["Effective"] = fwValueToInterface(plan.Effective, map[string]fwWireField{}, false)
	}
	if !plan.EvaluationCount.IsNull() && !plan.EvaluationCount.IsUnknown() {
		planValues["EvaluationCount"] = fwValueToInterface(plan.EvaluationCount, map[string]fwWireField{}, false)
	}
	if !plan.Expression.IsNull() && !plan.Expression.IsUnknown() {
		planValues["Expression"] = fwValueToInterface(plan.Expression, map[string]fwWireField{"metric_name": {Name: "MetricName", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "comparison_operator": {Name: "ComparisonOperator", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "period": {Name: "Period", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "statistics": {Name: "Statistics", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "threshold": {Name: "Threshold", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.ExpressionsLogicOperator.IsNull() && !plan.ExpressionsLogicOperator.IsUnknown() {
		planValues["ExpressionsLogicOperator"] = fwValueToInterface(plan.ExpressionsLogicOperator, map[string]fwWireField{}, false)
	}
	if !plan.MetricName.IsNull() && !plan.MetricName.IsUnknown() {
		planValues["MetricName"] = fwValueToInterface(plan.MetricName, map[string]fwWireField{}, false)
	}
	if !plan.MetricType.IsNull() && !plan.MetricType.IsUnknown() {
		planValues["MetricType"] = fwValueToInterface(plan.MetricType, map[string]fwWireField{}, false)
	}
	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		planValues["Period"] = fwValueToInterface(plan.Period, map[string]fwWireField{}, false)
	}
	if !plan.ScalingGroupId.IsNull() && !plan.ScalingGroupId.IsUnknown() {
		planValues["ScalingGroupId"] = fwValueToInterface(plan.ScalingGroupId, map[string]fwWireField{}, false)
	}
	if !plan.Statistics.IsNull() && !plan.Statistics.IsUnknown() {
		planValues["Statistics"] = fwValueToInterface(plan.Statistics, map[string]fwWireField{}, false)
	}
	if !plan.Threshold.IsNull() && !plan.Threshold.IsUnknown() {
		planValues["Threshold"] = fwValueToInterface(plan.Threshold, map[string]fwWireField{}, false)
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		planValues["Status"] = fwValueToInterface(plan.Status, map[string]fwWireField{}, false)
	}
	stateValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.AlarmActions.IsNull() && !state.AlarmActions.IsUnknown() {
		stateValues["AlarmActions"] = fwValueToInterface(state.AlarmActions, map[string]fwWireField{}, false)
	}
	if !state.AlarmTaskName.IsNull() && !state.AlarmTaskName.IsUnknown() {
		stateValues["AlarmTaskName"] = fwValueToInterface(state.AlarmTaskName, map[string]fwWireField{}, false)
	}
	if !state.AlarmTriggerState.IsNull() && !state.AlarmTriggerState.IsUnknown() {
		stateValues["AlarmTriggerState"] = fwValueToInterface(state.AlarmTriggerState, map[string]fwWireField{}, false)
	}
	if !state.ComparisonOperator.IsNull() && !state.ComparisonOperator.IsUnknown() {
		stateValues["ComparisonOperator"] = fwValueToInterface(state.ComparisonOperator, map[string]fwWireField{}, false)
	}
	if !state.Description.IsNull() && !state.Description.IsUnknown() {
		stateValues["Description"] = fwValueToInterface(state.Description, map[string]fwWireField{}, false)
	}
	if !state.Dimension.IsNull() && !state.Dimension.IsUnknown() {
		stateValues["Dimension"] = fwValueToInterface(state.Dimension, map[string]fwWireField{"dimension_value": {Name: "DimensionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dimension_key": {Name: "DimensionKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.Effective.IsNull() && !state.Effective.IsUnknown() {
		stateValues["Effective"] = fwValueToInterface(state.Effective, map[string]fwWireField{}, false)
	}
	if !state.EvaluationCount.IsNull() && !state.EvaluationCount.IsUnknown() {
		stateValues["EvaluationCount"] = fwValueToInterface(state.EvaluationCount, map[string]fwWireField{}, false)
	}
	if !state.Expression.IsNull() && !state.Expression.IsUnknown() {
		stateValues["Expression"] = fwValueToInterface(state.Expression, map[string]fwWireField{"metric_name": {Name: "MetricName", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "comparison_operator": {Name: "ComparisonOperator", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "period": {Name: "Period", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "statistics": {Name: "Statistics", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "threshold": {Name: "Threshold", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.ExpressionsLogicOperator.IsNull() && !state.ExpressionsLogicOperator.IsUnknown() {
		stateValues["ExpressionsLogicOperator"] = fwValueToInterface(state.ExpressionsLogicOperator, map[string]fwWireField{}, false)
	}
	if !state.MetricName.IsNull() && !state.MetricName.IsUnknown() {
		stateValues["MetricName"] = fwValueToInterface(state.MetricName, map[string]fwWireField{}, false)
	}
	if !state.MetricType.IsNull() && !state.MetricType.IsUnknown() {
		stateValues["MetricType"] = fwValueToInterface(state.MetricType, map[string]fwWireField{}, false)
	}
	if !state.Period.IsNull() && !state.Period.IsUnknown() {
		stateValues["Period"] = fwValueToInterface(state.Period, map[string]fwWireField{}, false)
	}
	if !state.ScalingGroupId.IsNull() && !state.ScalingGroupId.IsUnknown() {
		stateValues["ScalingGroupId"] = fwValueToInterface(state.ScalingGroupId, map[string]fwWireField{}, false)
	}
	if !state.Statistics.IsNull() && !state.Statistics.IsUnknown() {
		stateValues["Statistics"] = fwValueToInterface(state.Statistics, map[string]fwWireField{}, false)
	}
	if !state.Threshold.IsNull() && !state.Threshold.IsUnknown() {
		stateValues["Threshold"] = fwValueToInterface(state.Threshold, map[string]fwWireField{}, false)
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		stateValues["Status"] = fwValueToInterface(state.Status, map[string]fwWireField{}, false)
	}
	updated0 := false
	if fwPlannedChanged(plan.Status, state.Status) {
		updated0 = true
	}
	updated0 = false
	statusValues0_0 := fwWalkMappedPath(planValues, fwMappingPath("$.Status"), nil)
	previousStatus0_0 := fwWalkMappedPath(stateValues, fwMappingPath("$.Status"), nil)
	if len(statusValues0_0) > 0 && !fwValuesEqual(statusValues0_0[0].Value, func() interface{} {
		if len(previousStatus0_0) > 0 {
			return previousStatus0_0[0].Value
		}
		return nil
	}()) {
		updated0 = true
	}
	if len(statusValues0_0) == 0 || fmt.Sprint(statusValues0_0[0].Value) != "true" {
		updated0 = false
	}
	if updated0 {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"AlarmTaskId": state.Id.ValueString(),
			"ClientToken": buildClientToken("EnableAlarm"),
		}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, operationValues0, []fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building EnableAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2014-08-28", "EnableAlarm", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_auto_scaling_alarm_task: EnableAlarm", err.Error())
			return
		}
	}
	updated1 := false
	if fwPlannedChanged(plan.Status, state.Status) {
		updated1 = true
	}
	updated1 = false
	statusValues1_0 := fwWalkMappedPath(planValues, fwMappingPath("$.Status"), nil)
	previousStatus1_0 := fwWalkMappedPath(stateValues, fwMappingPath("$.Status"), nil)
	if len(statusValues1_0) > 0 && !fwValuesEqual(statusValues1_0[0].Value, func() interface{} {
		if len(previousStatus1_0) > 0 {
			return previousStatus1_0[0].Value
		}
		return nil
	}()) {
		updated1 = true
	}
	if len(statusValues1_0) == 0 || fmt.Sprint(statusValues1_0[0].Value) != "false" {
		updated1 = false
	}
	if updated1 {
		operationValues1 := planValues
		request1 := map[string]interface{}{
			"AlarmTaskId": state.Id.ValueString(),
			"ClientToken": buildClientToken("DisableAlarm"),
		}
		request1["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request1, operationValues1, []fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DisableAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2014-08-28", "DisableAlarm", request1, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_auto_scaling_alarm_task: DisableAlarm", err.Error())
			return
		}
	}
	updated2 := false
	if fwPlannedChanged(plan.Status, state.Status) {
		updated2 = true
	}
	updated2 = false
	statusValues2_0 := fwWalkMappedPath(planValues, fwMappingPath("$.Status"), nil)
	previousStatus2_0 := fwWalkMappedPath(stateValues, fwMappingPath("$.Status"), nil)
	if len(statusValues2_0) > 0 && !fwValuesEqual(statusValues2_0[0].Value, func() interface{} {
		if len(previousStatus2_0) > 0 {
			return previousStatus2_0[0].Value
		}
		return nil
	}()) {
		updated2 = true
	}
	if len(statusValues2_0) == 0 || fmt.Sprint(statusValues2_0[0].Value) != "true" {
		updated2 = false
	}
	if updated2 {
		operationValues2 := planValues
		request2 := map[string]interface{}{
			"AlarmTaskId": state.Id.ValueString(),
			"ClientToken": buildClientToken("EnableAlarm"),
		}
		request2["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request2, operationValues2, []fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building EnableAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2022-02-22", "EnableAlarm", request2, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_auto_scaling_alarm_task: EnableAlarm", err.Error())
			return
		}
	}
	updated3 := false
	if fwPlannedChanged(plan.Status, state.Status) {
		updated3 = true
	}
	updated3 = false
	statusValues3_0 := fwWalkMappedPath(planValues, fwMappingPath("$.Status"), nil)
	previousStatus3_0 := fwWalkMappedPath(stateValues, fwMappingPath("$.Status"), nil)
	if len(statusValues3_0) > 0 && !fwValuesEqual(statusValues3_0[0].Value, func() interface{} {
		if len(previousStatus3_0) > 0 {
			return previousStatus3_0[0].Value
		}
		return nil
	}()) {
		updated3 = true
	}
	if len(statusValues3_0) == 0 || fmt.Sprint(statusValues3_0[0].Value) != "false" {
		updated3 = false
	}
	if updated3 {
		operationValues3 := planValues
		request3 := map[string]interface{}{
			"AlarmTaskId": state.Id.ValueString(),
			"ClientToken": buildClientToken("DisableAlarm"),
		}
		request3["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request3, operationValues3, []fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DisableAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2022-02-22", "DisableAlarm", request3, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_auto_scaling_alarm_task: DisableAlarm", err.Error())
			return
		}
	}
	object, err = AlarmTaskGet(r.Client(), id)
	if err != nil || object == nil {
		resp.Diagnostics.AddError("Reading after post-create updates", fmt.Sprint(err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alarmTaskState(object, &plan))...)
}

func (r *alarmTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alarmTaskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := AlarmTaskGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeAlarms", alarmTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeAlarms", alarmTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alarmTaskState(object, &state))...)
}

func (r *alarmTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state alarmTaskResourceModel
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
	if !plan.AlarmActions.IsNull() && !plan.AlarmActions.IsUnknown() {
		planValues["AlarmActions"] = fwValueToInterface(plan.AlarmActions, map[string]fwWireField{}, false)
	}
	if !plan.AlarmTaskName.IsNull() && !plan.AlarmTaskName.IsUnknown() {
		planValues["AlarmTaskName"] = fwValueToInterface(plan.AlarmTaskName, map[string]fwWireField{}, false)
	}
	if !plan.AlarmTriggerState.IsNull() && !plan.AlarmTriggerState.IsUnknown() {
		planValues["AlarmTriggerState"] = fwValueToInterface(plan.AlarmTriggerState, map[string]fwWireField{}, false)
	}
	if !plan.ComparisonOperator.IsNull() && !plan.ComparisonOperator.IsUnknown() {
		planValues["ComparisonOperator"] = fwValueToInterface(plan.ComparisonOperator, map[string]fwWireField{}, false)
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		planValues["Description"] = fwValueToInterface(plan.Description, map[string]fwWireField{}, false)
	}
	if !plan.Dimension.IsNull() && !plan.Dimension.IsUnknown() {
		planValues["Dimension"] = fwValueToInterface(plan.Dimension, map[string]fwWireField{"dimension_value": {Name: "DimensionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dimension_key": {Name: "DimensionKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.Effective.IsNull() && !plan.Effective.IsUnknown() {
		planValues["Effective"] = fwValueToInterface(plan.Effective, map[string]fwWireField{}, false)
	}
	if !plan.EvaluationCount.IsNull() && !plan.EvaluationCount.IsUnknown() {
		planValues["EvaluationCount"] = fwValueToInterface(plan.EvaluationCount, map[string]fwWireField{}, false)
	}
	if !plan.Expression.IsNull() && !plan.Expression.IsUnknown() {
		planValues["Expression"] = fwValueToInterface(plan.Expression, map[string]fwWireField{"metric_name": {Name: "MetricName", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "comparison_operator": {Name: "ComparisonOperator", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "period": {Name: "Period", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "statistics": {Name: "Statistics", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "threshold": {Name: "Threshold", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.ExpressionsLogicOperator.IsNull() && !plan.ExpressionsLogicOperator.IsUnknown() {
		planValues["ExpressionsLogicOperator"] = fwValueToInterface(plan.ExpressionsLogicOperator, map[string]fwWireField{}, false)
	}
	if !plan.MetricName.IsNull() && !plan.MetricName.IsUnknown() {
		planValues["MetricName"] = fwValueToInterface(plan.MetricName, map[string]fwWireField{}, false)
	}
	if !plan.MetricType.IsNull() && !plan.MetricType.IsUnknown() {
		planValues["MetricType"] = fwValueToInterface(plan.MetricType, map[string]fwWireField{}, false)
	}
	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		planValues["Period"] = fwValueToInterface(plan.Period, map[string]fwWireField{}, false)
	}
	if !plan.ScalingGroupId.IsNull() && !plan.ScalingGroupId.IsUnknown() {
		planValues["ScalingGroupId"] = fwValueToInterface(plan.ScalingGroupId, map[string]fwWireField{}, false)
	}
	if !plan.Statistics.IsNull() && !plan.Statistics.IsUnknown() {
		planValues["Statistics"] = fwValueToInterface(plan.Statistics, map[string]fwWireField{}, false)
	}
	if !plan.Threshold.IsNull() && !plan.Threshold.IsUnknown() {
		planValues["Threshold"] = fwValueToInterface(plan.Threshold, map[string]fwWireField{}, false)
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		planValues["Status"] = fwValueToInterface(plan.Status, map[string]fwWireField{}, false)
	}
	stateValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.AlarmActions.IsNull() && !state.AlarmActions.IsUnknown() {
		stateValues["AlarmActions"] = fwValueToInterface(state.AlarmActions, map[string]fwWireField{}, false)
	}
	if !state.AlarmTaskName.IsNull() && !state.AlarmTaskName.IsUnknown() {
		stateValues["AlarmTaskName"] = fwValueToInterface(state.AlarmTaskName, map[string]fwWireField{}, false)
	}
	if !state.AlarmTriggerState.IsNull() && !state.AlarmTriggerState.IsUnknown() {
		stateValues["AlarmTriggerState"] = fwValueToInterface(state.AlarmTriggerState, map[string]fwWireField{}, false)
	}
	if !state.ComparisonOperator.IsNull() && !state.ComparisonOperator.IsUnknown() {
		stateValues["ComparisonOperator"] = fwValueToInterface(state.ComparisonOperator, map[string]fwWireField{}, false)
	}
	if !state.Description.IsNull() && !state.Description.IsUnknown() {
		stateValues["Description"] = fwValueToInterface(state.Description, map[string]fwWireField{}, false)
	}
	if !state.Dimension.IsNull() && !state.Dimension.IsUnknown() {
		stateValues["Dimension"] = fwValueToInterface(state.Dimension, map[string]fwWireField{"dimension_value": {Name: "DimensionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dimension_key": {Name: "DimensionKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.Effective.IsNull() && !state.Effective.IsUnknown() {
		stateValues["Effective"] = fwValueToInterface(state.Effective, map[string]fwWireField{}, false)
	}
	if !state.EvaluationCount.IsNull() && !state.EvaluationCount.IsUnknown() {
		stateValues["EvaluationCount"] = fwValueToInterface(state.EvaluationCount, map[string]fwWireField{}, false)
	}
	if !state.Expression.IsNull() && !state.Expression.IsUnknown() {
		stateValues["Expression"] = fwValueToInterface(state.Expression, map[string]fwWireField{"metric_name": {Name: "MetricName", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "comparison_operator": {Name: "ComparisonOperator", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "period": {Name: "Period", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "statistics": {Name: "Statistics", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "threshold": {Name: "Threshold", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.ExpressionsLogicOperator.IsNull() && !state.ExpressionsLogicOperator.IsUnknown() {
		stateValues["ExpressionsLogicOperator"] = fwValueToInterface(state.ExpressionsLogicOperator, map[string]fwWireField{}, false)
	}
	if !state.MetricName.IsNull() && !state.MetricName.IsUnknown() {
		stateValues["MetricName"] = fwValueToInterface(state.MetricName, map[string]fwWireField{}, false)
	}
	if !state.MetricType.IsNull() && !state.MetricType.IsUnknown() {
		stateValues["MetricType"] = fwValueToInterface(state.MetricType, map[string]fwWireField{}, false)
	}
	if !state.Period.IsNull() && !state.Period.IsUnknown() {
		stateValues["Period"] = fwValueToInterface(state.Period, map[string]fwWireField{}, false)
	}
	if !state.ScalingGroupId.IsNull() && !state.ScalingGroupId.IsUnknown() {
		stateValues["ScalingGroupId"] = fwValueToInterface(state.ScalingGroupId, map[string]fwWireField{}, false)
	}
	if !state.Statistics.IsNull() && !state.Statistics.IsUnknown() {
		stateValues["Statistics"] = fwValueToInterface(state.Statistics, map[string]fwWireField{}, false)
	}
	if !state.Threshold.IsNull() && !state.Threshold.IsUnknown() {
		stateValues["Threshold"] = fwValueToInterface(state.Threshold, map[string]fwWireField{}, false)
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		stateValues["Status"] = fwValueToInterface(state.Status, map[string]fwWireField{}, false)
	}
	updated0 := false
	if fwPlannedChanged(plan.Description, state.Description) {
		updated0 = true
	}
	if fwPlannedChanged(plan.Expression, state.Expression) {
		updated0 = true
	}
	if fwPlannedChanged(plan.MetricType, state.MetricType) {
		updated0 = true
	}
	if fwPlannedChanged(plan.EvaluationCount, state.EvaluationCount) {
		updated0 = true
	}
	if fwPlannedChanged(plan.Effective, state.Effective) {
		updated0 = true
	}
	if fwPlannedChanged(plan.AlarmActions, state.AlarmActions) {
		updated0 = true
	}
	if fwPlannedChanged(plan.Dimension, state.Dimension) {
		updated0 = true
	}
	if fwPlannedChanged(plan.MetricName, state.MetricName) {
		updated0 = true
	}
	if fwPlannedChanged(plan.ComparisonOperator, state.ComparisonOperator) {
		updated0 = true
	}
	if fwPlannedChanged(plan.Period, state.Period) {
		updated0 = true
	}
	if fwPlannedChanged(plan.Statistics, state.Statistics) {
		updated0 = true
	}
	if fwPlannedChanged(plan.Threshold, state.Threshold) {
		updated0 = true
	}
	if fwPlannedChanged(plan.ExpressionsLogicOperator, state.ExpressionsLogicOperator) {
		updated0 = true
	}
	if fwPlannedChanged(plan.AlarmTaskName, state.AlarmTaskName) {
		updated0 = true
	}
	if updated0 {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"AlarmTaskId": state.Id.ValueString(),
			"ClientToken": buildClientToken("ModifyAlarm"),
		}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, operationValues0, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Description", Target: "Description", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.MetricType", Target: "MetricType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.EvaluationCount", Target: "EvaluationCount", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Effective", Target: "Effective", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmActions[*]", Target: "AlarmAction[*]", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.Dimension[*].DimensionKey", Target: "Dimension[*].DimensionKey", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.Dimension[*].DimensionValue", Target: "Dimension[*].DimensionValue", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.MetricName", Target: "MetricName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ComparisonOperator", Target: "ComparisonOperator", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Period", Target: "Period", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Statistics", Target: "Statistics", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Threshold", Target: "Threshold", Encoding: "uniform-scalar|number", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ExpressionsLogicOperator", Target: "ExpressionsLogicOperator", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskName", Target: "Name", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Expression[*]", Target: "Expression[*]", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			// Dimension is a full-replace list on ModifyAlarm: omitting it when
			// unchanged clears the existing dimensions (the read then drifts
			// dimension.# 1 -> 0). v1 always sends request.Dimension in
			// ModifyAlarm, so force the binding through fwChangedUpdateMappings by
			// marking $.Dimension as required for both the 2014 and 2022 calls.
		}, operationValues0, stateValues, []string{"$.AlarmTaskId", "$.RegionId", "$.Dimension"}), false); err != nil {
			resp.Diagnostics.AddError("Building ModifyAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2014-08-28", "ModifyAlarm", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_auto_scaling_alarm_task: ModifyAlarm", err.Error())
			return
		}
	}
	updated1 := false
	if fwPlannedChanged(plan.Status, state.Status) {
		updated1 = true
	}
	updated1 = false
	statusValues1_0 := fwWalkMappedPath(planValues, fwMappingPath("$.Status"), nil)
	previousStatus1_0 := fwWalkMappedPath(stateValues, fwMappingPath("$.Status"), nil)
	if len(statusValues1_0) > 0 && !fwValuesEqual(statusValues1_0[0].Value, func() interface{} {
		if len(previousStatus1_0) > 0 {
			return previousStatus1_0[0].Value
		}
		return nil
	}()) {
		updated1 = true
	}
	if len(statusValues1_0) == 0 || fmt.Sprint(statusValues1_0[0].Value) != "true" {
		updated1 = false
	}
	if updated1 {
		operationValues1 := planValues
		request1 := map[string]interface{}{
			"AlarmTaskId": state.Id.ValueString(),
			"ClientToken": buildClientToken("EnableAlarm"),
		}
		request1["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request1, operationValues1, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues1, stateValues, []string{"$.AlarmTaskId", "$.RegionId"}), false); err != nil {
			resp.Diagnostics.AddError("Building EnableAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2014-08-28", "EnableAlarm", request1, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_auto_scaling_alarm_task: EnableAlarm", err.Error())
			return
		}
	}
	updated2 := false
	if fwPlannedChanged(plan.Status, state.Status) {
		updated2 = true
	}
	updated2 = false
	statusValues2_0 := fwWalkMappedPath(planValues, fwMappingPath("$.Status"), nil)
	previousStatus2_0 := fwWalkMappedPath(stateValues, fwMappingPath("$.Status"), nil)
	if len(statusValues2_0) > 0 && !fwValuesEqual(statusValues2_0[0].Value, func() interface{} {
		if len(previousStatus2_0) > 0 {
			return previousStatus2_0[0].Value
		}
		return nil
	}()) {
		updated2 = true
	}
	if len(statusValues2_0) == 0 || fmt.Sprint(statusValues2_0[0].Value) != "false" {
		updated2 = false
	}
	if updated2 {
		operationValues2 := planValues
		request2 := map[string]interface{}{
			"AlarmTaskId": state.Id.ValueString(),
			"ClientToken": buildClientToken("DisableAlarm"),
		}
		request2["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request2, operationValues2, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues2, stateValues, []string{"$.AlarmTaskId", "$.RegionId"}), false); err != nil {
			resp.Diagnostics.AddError("Building DisableAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2014-08-28", "DisableAlarm", request2, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_auto_scaling_alarm_task: DisableAlarm", err.Error())
			return
		}
	}
	updated3 := false
	if fwPlannedChanged(plan.Description, state.Description) {
		updated3 = true
	}
	if fwPlannedChanged(plan.Expression, state.Expression) {
		updated3 = true
	}
	if fwPlannedChanged(plan.MetricType, state.MetricType) {
		updated3 = true
	}
	if fwPlannedChanged(plan.EvaluationCount, state.EvaluationCount) {
		updated3 = true
	}
	if fwPlannedChanged(plan.Effective, state.Effective) {
		updated3 = true
	}
	if fwPlannedChanged(plan.AlarmActions, state.AlarmActions) {
		updated3 = true
	}
	if fwPlannedChanged(plan.Dimension, state.Dimension) {
		updated3 = true
	}
	if fwPlannedChanged(plan.MetricName, state.MetricName) {
		updated3 = true
	}
	if fwPlannedChanged(plan.ComparisonOperator, state.ComparisonOperator) {
		updated3 = true
	}
	if fwPlannedChanged(plan.Period, state.Period) {
		updated3 = true
	}
	if fwPlannedChanged(plan.Statistics, state.Statistics) {
		updated3 = true
	}
	if fwPlannedChanged(plan.Threshold, state.Threshold) {
		updated3 = true
	}
	if fwPlannedChanged(plan.ExpressionsLogicOperator, state.ExpressionsLogicOperator) {
		updated3 = true
	}
	if fwPlannedChanged(plan.AlarmTaskName, state.AlarmTaskName) {
		updated3 = true
	}
	if updated3 {
		operationValues3 := planValues
		request3 := map[string]interface{}{
			"AlarmTaskId": state.Id.ValueString(),
			"ClientToken": buildClientToken("ModifyAlarm"),
		}
		request3["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request3, operationValues3, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Description", Target: "Description", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.MetricType", Target: "MetricType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.EvaluationCount", Target: "EvaluationCount", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Effective", Target: "Effective", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmActions[*]", Target: "AlarmActions[*]", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.Dimension[*].DimensionKey", Target: "Dimensions[*].DimensionKey", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.Dimension[*].DimensionValue", Target: "Dimensions[*].DimensionValue", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
			{Source: "$.MetricName", Target: "MetricName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ComparisonOperator", Target: "ComparisonOperator", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Period", Target: "Period", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Statistics", Target: "Statistics", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Threshold", Target: "Threshold", Encoding: "uniform-scalar|number", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ExpressionsLogicOperator", Target: "ExpressionsLogicOperator", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskName", Target: "Name", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Expression[*]", Target: "Expressions[*]", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		}, operationValues3, stateValues, []string{"$.AlarmTaskId", "$.RegionId", "$.Dimension"}), false); err != nil {
			resp.Diagnostics.AddError("Building ModifyAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2022-02-22", "ModifyAlarm", request3, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_auto_scaling_alarm_task: ModifyAlarm", err.Error())
			return
		}
	}
	updated4 := false
	if fwPlannedChanged(plan.Status, state.Status) {
		updated4 = true
	}
	updated4 = false
	statusValues4_0 := fwWalkMappedPath(planValues, fwMappingPath("$.Status"), nil)
	previousStatus4_0 := fwWalkMappedPath(stateValues, fwMappingPath("$.Status"), nil)
	if len(statusValues4_0) > 0 && !fwValuesEqual(statusValues4_0[0].Value, func() interface{} {
		if len(previousStatus4_0) > 0 {
			return previousStatus4_0[0].Value
		}
		return nil
	}()) {
		updated4 = true
	}
	if len(statusValues4_0) == 0 || fmt.Sprint(statusValues4_0[0].Value) != "true" {
		updated4 = false
	}
	if updated4 {
		operationValues4 := planValues
		request4 := map[string]interface{}{
			"AlarmTaskId": state.Id.ValueString(),
			"ClientToken": buildClientToken("EnableAlarm"),
		}
		request4["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request4, operationValues4, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues4, stateValues, []string{"$.AlarmTaskId", "$.RegionId"}), false); err != nil {
			resp.Diagnostics.AddError("Building EnableAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2022-02-22", "EnableAlarm", request4, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_auto_scaling_alarm_task: EnableAlarm", err.Error())
			return
		}
	}
	updated5 := false
	if fwPlannedChanged(plan.Status, state.Status) {
		updated5 = true
	}
	updated5 = false
	statusValues5_0 := fwWalkMappedPath(planValues, fwMappingPath("$.Status"), nil)
	previousStatus5_0 := fwWalkMappedPath(stateValues, fwMappingPath("$.Status"), nil)
	if len(statusValues5_0) > 0 && !fwValuesEqual(statusValues5_0[0].Value, func() interface{} {
		if len(previousStatus5_0) > 0 {
			return previousStatus5_0[0].Value
		}
		return nil
	}()) {
		updated5 = true
	}
	if len(statusValues5_0) == 0 || fmt.Sprint(statusValues5_0[0].Value) != "false" {
		updated5 = false
	}
	if updated5 {
		operationValues5 := planValues
		request5 := map[string]interface{}{
			"AlarmTaskId": state.Id.ValueString(),
			"ClientToken": buildClientToken("DisableAlarm"),
		}
		request5["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request5, operationValues5, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues5, stateValues, []string{"$.AlarmTaskId", "$.RegionId"}), false); err != nil {
			resp.Diagnostics.AddError("Building DisableAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2022-02-22", "DisableAlarm", request5, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_auto_scaling_alarm_task: DisableAlarm", err.Error())
			return
		}
	}
	object, err := AlarmTaskGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeAlarms", alarmTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeAlarms", alarmTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alarmTaskState(object, &plan))...)
}

func (r *alarmTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alarmTaskResourceModel
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
	if !state.AlarmActions.IsNull() && !state.AlarmActions.IsUnknown() {
		deleteValues["AlarmActions"] = fwValueToInterface(state.AlarmActions, map[string]fwWireField{}, false)
	}
	if !state.AlarmTaskName.IsNull() && !state.AlarmTaskName.IsUnknown() {
		deleteValues["AlarmTaskName"] = fwValueToInterface(state.AlarmTaskName, map[string]fwWireField{}, false)
	}
	if !state.AlarmTriggerState.IsNull() && !state.AlarmTriggerState.IsUnknown() {
		deleteValues["AlarmTriggerState"] = fwValueToInterface(state.AlarmTriggerState, map[string]fwWireField{}, false)
	}
	if !state.ComparisonOperator.IsNull() && !state.ComparisonOperator.IsUnknown() {
		deleteValues["ComparisonOperator"] = fwValueToInterface(state.ComparisonOperator, map[string]fwWireField{}, false)
	}
	if !state.Description.IsNull() && !state.Description.IsUnknown() {
		deleteValues["Description"] = fwValueToInterface(state.Description, map[string]fwWireField{}, false)
	}
	if !state.Dimension.IsNull() && !state.Dimension.IsUnknown() {
		deleteValues["Dimension"] = fwValueToInterface(state.Dimension, map[string]fwWireField{"dimension_value": {Name: "DimensionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dimension_key": {Name: "DimensionKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.Effective.IsNull() && !state.Effective.IsUnknown() {
		deleteValues["Effective"] = fwValueToInterface(state.Effective, map[string]fwWireField{}, false)
	}
	if !state.EvaluationCount.IsNull() && !state.EvaluationCount.IsUnknown() {
		deleteValues["EvaluationCount"] = fwValueToInterface(state.EvaluationCount, map[string]fwWireField{}, false)
	}
	if !state.Expression.IsNull() && !state.Expression.IsUnknown() {
		deleteValues["Expression"] = fwValueToInterface(state.Expression, map[string]fwWireField{"metric_name": {Name: "MetricName", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "comparison_operator": {Name: "ComparisonOperator", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "period": {Name: "Period", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "statistics": {Name: "Statistics", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "threshold": {Name: "Threshold", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.ExpressionsLogicOperator.IsNull() && !state.ExpressionsLogicOperator.IsUnknown() {
		deleteValues["ExpressionsLogicOperator"] = fwValueToInterface(state.ExpressionsLogicOperator, map[string]fwWireField{}, false)
	}
	if !state.MetricName.IsNull() && !state.MetricName.IsUnknown() {
		deleteValues["MetricName"] = fwValueToInterface(state.MetricName, map[string]fwWireField{}, false)
	}
	if !state.MetricType.IsNull() && !state.MetricType.IsUnknown() {
		deleteValues["MetricType"] = fwValueToInterface(state.MetricType, map[string]fwWireField{}, false)
	}
	if !state.Period.IsNull() && !state.Period.IsUnknown() {
		deleteValues["Period"] = fwValueToInterface(state.Period, map[string]fwWireField{}, false)
	}
	if !state.ScalingGroupId.IsNull() && !state.ScalingGroupId.IsUnknown() {
		deleteValues["ScalingGroupId"] = fwValueToInterface(state.ScalingGroupId, map[string]fwWireField{}, false)
	}
	if !state.Statistics.IsNull() && !state.Statistics.IsUnknown() {
		deleteValues["Statistics"] = fwValueToInterface(state.Statistics, map[string]fwWireField{}, false)
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		deleteValues["Status"] = fwValueToInterface(state.Status, map[string]fwWireField{}, false)
	}
	deleteValues["AlarmTaskId"] = state.Id.ValueString()
	needsRefresh := false
	if needsRefresh {
		object, err := AlarmTaskGet(r.Client(), state.Id.ValueString())
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
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteAlarm")}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AlarmTaskId", Target: "AlarmTaskId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DeleteAlarm request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "Ess", "2022-02-22", "DeleteAlarm", request0, operationTimeout); err != nil && !isExpectedErrors(err, []string{"ResourceNotFound.AlarmTask", "ResourceNotFound"}) {
			resp.Diagnostics.AddError("Deleting alicloud_auto_scaling_alarm_task: DeleteAlarm", err.Error())
			return
		}
	}
}

func (r *alarmTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// alarmTaskState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func alarmTaskState(object map[string]interface{}, prior *alarmTaskResourceModel) alarmTaskResourceModel {
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
		// The API returns Enable as a bool; the schema models status as the
		// string "true"/"false" (v1 stores enable as a bool). Without this, str()
		// returns Null on a bool, so status drifts to null after apply.
		if b, ok := cur.(bool); ok {
			if b {
				return types.StringValue("true")
			}
			return types.StringValue("false")
		}
		return types.StringNull()
	}
	intVal := func(key string) types.Int64 {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.Int64Null()
			}
			cur, ok = m[seg]
			if !ok {
				return types.Int64Null()
			}
		}
		switch v := cur.(type) {
		case float64:
			return types.Int64Value(int64(v))
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.Int64Value(n)
			}
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				return types.Int64Value(n)
			}
		}
		return types.Int64Null()
	}
	floatVal := func(key string) types.Float64 {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.Float64Null()
			}
			cur, ok = m[seg]
			if !ok {
				return types.Float64Null()
			}
		}
		switch v := cur.(type) {
		case float64:
			return types.Float64Value(v)
		case int:
			return types.Float64Value(float64(v))
		case int64:
			return types.Float64Value(float64(v))
		case json.Number:
			if n, err := v.Float64(); err == nil {
				return types.Float64Value(n)
			}
		case string:
			if n, err := strconv.ParseFloat(v, 64); err == nil {
				return types.Float64Value(n)
			}
		}
		return types.Float64Null()
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
	m := alarmTaskResourceModel{
		FwTimeouts:               types.ObjectNull(fwTimeoutTypes()),
		Id:                       str("AlarmTaskId"),
		AlarmActions:             alarmTaskAlarmActionsToTypeList(lookup("AlarmActions")),
		AlarmTaskName:            str("AlarmTaskName"),
		AlarmTriggerState:        str("AlarmTriggerState"),
		ComparisonOperator:       str("ComparisonOperator"),
		Description:              str("Description"),
		Dimension:                alarmTaskDimensionToTypeList(lookup("Dimension")),
		Effective:                str("Effective"),
		EvaluationCount:          intVal("EvaluationCount"),
		Expression:               alarmTaskExpressionToTypeList(lookup("Expression")),
		ExpressionsLogicOperator: str("ExpressionsLogicOperator"),
		MetricName:               str("MetricName"),
		MetricType:               str("MetricType"),
		Period:                   intVal("Period"),
		ScalingGroupId:           str("ScalingGroupId"),
		Statistics:               str("Statistics"),
		Threshold:                floatVal("Threshold"),
		Status:                   str("Status"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		m.Dimension = fwPreserveRequestValues(m.Dimension, prior.Dimension, map[string]fwWireField{"dimension_value": {Name: "DimensionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dimension_key": {Name: "DimensionKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}).(types.List)
		// The API returns a one-element Expressions echo for single-metric alarms
		// (mirroring the top-level metric) and the real Expressions[] for
		// multi-metric alarms; a one-element list is ambiguous, so the metric_name
		// field alone cannot tell them apart (the API mirrors MetricName for
		// multi-metric too). Decide by the prior plan when it carries the
		// configuration, and fall back to the API Expressions count for import,
		// where the prior state only has the id and neither plan field is set.
		expressionFields := map[string]fwWireField{"metric_name": {Name: "MetricName", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "comparison_operator": {Name: "ComparisonOperator", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "period": {Name: "Period", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "statistics": {Name: "Statistics", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "threshold": {Name: "Threshold", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}
		switch {
		case !prior.Expression.IsNull() && !prior.Expression.IsUnknown() && len(prior.Expression.Elements()) > 0:
			// multi-metric: keep the API expressions, preserving request-only fields
			m.Expression = fwPreserveRequestValues(m.Expression, prior.Expression, expressionFields).(types.List)
		case !prior.MetricName.IsNull() && prior.MetricName.ValueString() != "":
			// single-metric (top-level metric_name, no expression blocks): suppress
			// the one-element Expressions echo so expression.# stays 0
			m.Expression = alarmTaskExpressionToTypeList(nil)
		default:
			// import: prior carries only the id, so neither plan field is set.
			// >=2 API expressions => multi-metric (keep); <=1 => single-metric echo.
			if !m.Expression.IsNull() && !m.Expression.IsUnknown() && len(m.Expression.Elements()) >= 2 {
				m.Expression = fwPreserveRequestValues(m.Expression, prior.Expression, expressionFields).(types.List)
			} else {
				m.Expression = alarmTaskExpressionToTypeList(nil)
			}
		}
	}
	return m
}
