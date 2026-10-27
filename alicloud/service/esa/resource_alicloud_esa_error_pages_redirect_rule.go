// Package esa. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package esa

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const errorPagesRedirectRuleTypeName = "alicloud_esa_error_pages_redirect_rule"

var (
	_ resource.Resource                = &errorPagesRedirectRuleResource{}
	_ resource.ResourceWithImportState = &errorPagesRedirectRuleResource{}
)

func NewErrorPagesRedirectRuleResource() resource.Resource {
	return &errorPagesRedirectRuleResource{}
}

type errorPagesRedirectRuleResource struct {
	fwadapt.ResourceBase
}

type errorPagesRedirectRuleResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	FwTimeouts         types.Object `tfsdk:"timeouts"`
	ConfigId           types.Int64  `tfsdk:"config_id"`
	ConfigType         types.String `tfsdk:"config_type"`
	ErrorPagesRedirect types.List   `tfsdk:"error_pages_redirect"`
	Rule               types.String `tfsdk:"rule"`
	RuleEnable         types.String `tfsdk:"rule_enable"`
	RuleName           types.String `tfsdk:"rule_name"`
	Sequence           types.Int64  `tfsdk:"sequence"`
	SiteId             types.Int64  `tfsdk:"site_id"`
	SiteVersion        types.Int64  `tfsdk:"site_version"`
}

func (r *errorPagesRedirectRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = errorPagesRedirectRuleTypeName
}

func (r *errorPagesRedirectRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"config_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Configuration ID.",
			},
			"config_type": schema.StringAttribute{
				Computed:    true,
				Description: "Configuration type. Valid values:\n- `global`: Query global configuration.\n- `rule`: Query rule-based configuration.",
			},
			// A plain list attribute with an object element type (not a
			// ListNestedAttribute): this provider serves every half behind one
			// protocol 5 server (tf5muxserver + NewProtocol5), which cannot
			// express nested attributes. The element object type is shared with
			// the runtime value builder in service.go.
			"error_pages_redirect": schema.ListAttribute{
				ElementType: errorPagesRedirectRuleErrorPagesRedirectObjType(),
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				Required:    true,
				Description: "Error code redirection. Each element is an object with a required `target_url` (the redirect destination) and a required `status_code` (one of 400, 403, 404, 405, 414, 416, 500, 501, 502, 503, 504).",
			},
			"rule": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Rule content, which uses conditional expressions to match user requests. This parameter is not required when adding a global configuration. There are two usage scenarios:\n- Match all incoming requests: Set the value to `true`.\n- Match specific requests: Set the value to a custom expression, such as `(http.host eq \"video.example.com\")`.",
			},
			"rule_enable": schema.StringAttribute{
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off"),
				},
				Optional:    true,
				Computed:    true,
				Description: "Rule switch. This parameter is not required when adding a global configuration. Valid values:\n- `on`: Enable.\n- `off`: Disable.",
			},
			"rule_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Rule name. This parameter is not required when adding a global configuration.",
			},
			"sequence": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "Rule execution order. Rules with smaller values take higher priority.",
			},
			"site_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Description: "The site ID, which can be obtained by calling the [ListSites](https://help.aliyun.com/document_detail/2850189.html) operation.",
			},
			"site_version": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "The version number of the site configuration. For sites with configuration version management enabled, you can use this parameter to specify the site version for which the configuration takes effect. The default value is version 0.",
			},
		},
	}
}

func (r *errorPagesRedirectRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan errorPagesRedirectRuleResourceModel
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
		"ClientToken": buildClientToken("CreateErrorPagesRedirectRule"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.ConfigId.IsNull() && !plan.ConfigId.IsUnknown() {
		createValues["ConfigId"] = fwValueToInterface(plan.ConfigId, map[string]fwWireField{}, false)
	}
	if !plan.ConfigType.IsNull() && !plan.ConfigType.IsUnknown() {
		createValues["ConfigType"] = fwValueToInterface(plan.ConfigType, map[string]fwWireField{}, false)
	}
	if !plan.ErrorPagesRedirect.IsNull() && !plan.ErrorPagesRedirect.IsUnknown() {
		createValues["ErrorPagesRedirect"] = fwValueToInterface(plan.ErrorPagesRedirect, map[string]fwWireField{"target_url": {Name: "TargetURL", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "status_code": {Name: "StatusCode", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.Rule.IsNull() && !plan.Rule.IsUnknown() {
		createValues["Rule"] = fwValueToInterface(plan.Rule, map[string]fwWireField{}, false)
	}
	if !plan.RuleEnable.IsNull() && !plan.RuleEnable.IsUnknown() {
		createValues["RuleEnable"] = fwValueToInterface(plan.RuleEnable, map[string]fwWireField{}, false)
	}
	if !plan.RuleName.IsNull() && !plan.RuleName.IsUnknown() {
		createValues["RuleName"] = fwValueToInterface(plan.RuleName, map[string]fwWireField{}, false)
	}
	if !plan.Sequence.IsNull() && !plan.Sequence.IsUnknown() {
		createValues["Sequence"] = fwValueToInterface(plan.Sequence, map[string]fwWireField{}, false)
	}
	if !plan.SiteId.IsNull() && !plan.SiteId.IsUnknown() {
		createValues["SiteId"] = fwValueToInterface(plan.SiteId, map[string]fwWireField{}, false)
	}
	if !plan.SiteVersion.IsNull() && !plan.SiteVersion.IsUnknown() {
		createValues["SiteVersion"] = fwValueToInterface(plan.SiteVersion, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.Sequence", Target: "Sequence", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.SiteVersion", Target: "SiteVersion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.SiteId", Target: "SiteId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ErrorPagesRedirect[*].StatusCode", Target: "ErrorPagesRedirect[*].StatusCode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ErrorPagesRedirect[*].TargetURL", Target: "ErrorPagesRedirect[*].TargetURL", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RuleEnable", Target: "RuleEnable", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RuleName", Target: "RuleName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Rule", Target: "Rule", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(ctx, r.Client(), "ESA", "2024-09-10", "CreateErrorPagesRedirectRule", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateErrorPagesRedirectRule", errorPagesRedirectRuleTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(request, fwMappingPath("SiteId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_esa_error_pages_redirect_rule: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	identityValues1 := fwWalkMappedPath(response, fwMappingPath("$.ConfigId"), nil)
	if len(identityValues1) == 0 || identityValues1[0].Value == nil || fmt.Sprint(identityValues1[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_esa_error_pages_redirect_rule: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues1[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			fmt.Sprintf("The response carried no SiteId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.ConfigId = fwKnownValue(trackedPlan.ConfigId).(types.Int64)
	trackedPlan.ConfigType = fwKnownValue(trackedPlan.ConfigType).(types.String)
	trackedPlan.ErrorPagesRedirect = fwKnownValue(trackedPlan.ErrorPagesRedirect).(types.List)
	trackedPlan.Rule = fwKnownValue(trackedPlan.Rule).(types.String)
	trackedPlan.RuleEnable = fwKnownValue(trackedPlan.RuleEnable).(types.String)
	trackedPlan.RuleName = fwKnownValue(trackedPlan.RuleName).(types.String)
	trackedPlan.Sequence = fwKnownValue(trackedPlan.Sequence).(types.Int64)
	trackedPlan.SiteId = fwKnownValue(trackedPlan.SiteId).(types.Int64)
	trackedPlan.SiteVersion = fwKnownValue(trackedPlan.SiteVersion).(types.Int64)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := DescribeErrorPagesRedirectRule(ctx, r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", errorPagesRedirectRuleTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", errorPagesRedirectRuleTypeName),
			fmt.Sprintf("The resource %s is not visible via GetErrorPagesRedirectRule right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, errorPagesRedirectRuleState(object, &plan))...)
}

func (r *errorPagesRedirectRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state errorPagesRedirectRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := DescribeErrorPagesRedirectRule(ctx, r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, errorPagesRedirectRuleState(object, &state))...)
}

func (r *errorPagesRedirectRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state errorPagesRedirectRuleResourceModel
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
	if !plan.ConfigId.IsNull() && !plan.ConfigId.IsUnknown() {
		planValues["ConfigId"] = fwValueToInterface(plan.ConfigId, map[string]fwWireField{}, false)
	}
	if !plan.ConfigType.IsNull() && !plan.ConfigType.IsUnknown() {
		planValues["ConfigType"] = fwValueToInterface(plan.ConfigType, map[string]fwWireField{}, false)
	}
	if !plan.ErrorPagesRedirect.IsNull() && !plan.ErrorPagesRedirect.IsUnknown() {
		planValues["ErrorPagesRedirect"] = fwValueToInterface(plan.ErrorPagesRedirect, map[string]fwWireField{"target_url": {Name: "TargetURL", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "status_code": {Name: "StatusCode", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.Rule.IsNull() && !plan.Rule.IsUnknown() {
		planValues["Rule"] = fwValueToInterface(plan.Rule, map[string]fwWireField{}, false)
	}
	if !plan.RuleEnable.IsNull() && !plan.RuleEnable.IsUnknown() {
		planValues["RuleEnable"] = fwValueToInterface(plan.RuleEnable, map[string]fwWireField{}, false)
	}
	if !plan.RuleName.IsNull() && !plan.RuleName.IsUnknown() {
		planValues["RuleName"] = fwValueToInterface(plan.RuleName, map[string]fwWireField{}, false)
	}
	if !plan.Sequence.IsNull() && !plan.Sequence.IsUnknown() {
		planValues["Sequence"] = fwValueToInterface(plan.Sequence, map[string]fwWireField{}, false)
	}
	if !plan.SiteId.IsNull() && !plan.SiteId.IsUnknown() {
		planValues["SiteId"] = fwValueToInterface(plan.SiteId, map[string]fwWireField{}, false)
	}
	if !plan.SiteVersion.IsNull() && !plan.SiteVersion.IsUnknown() {
		planValues["SiteVersion"] = fwValueToInterface(plan.SiteVersion, map[string]fwWireField{}, false)
	}
	stateValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.ConfigId.IsNull() && !state.ConfigId.IsUnknown() {
		stateValues["ConfigId"] = fwValueToInterface(state.ConfigId, map[string]fwWireField{}, false)
	}
	if !state.ConfigType.IsNull() && !state.ConfigType.IsUnknown() {
		stateValues["ConfigType"] = fwValueToInterface(state.ConfigType, map[string]fwWireField{}, false)
	}
	if !state.ErrorPagesRedirect.IsNull() && !state.ErrorPagesRedirect.IsUnknown() {
		stateValues["ErrorPagesRedirect"] = fwValueToInterface(state.ErrorPagesRedirect, map[string]fwWireField{"target_url": {Name: "TargetURL", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "status_code": {Name: "StatusCode", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.Rule.IsNull() && !state.Rule.IsUnknown() {
		stateValues["Rule"] = fwValueToInterface(state.Rule, map[string]fwWireField{}, false)
	}
	if !state.RuleEnable.IsNull() && !state.RuleEnable.IsUnknown() {
		stateValues["RuleEnable"] = fwValueToInterface(state.RuleEnable, map[string]fwWireField{}, false)
	}
	if !state.RuleName.IsNull() && !state.RuleName.IsUnknown() {
		stateValues["RuleName"] = fwValueToInterface(state.RuleName, map[string]fwWireField{}, false)
	}
	if !state.Sequence.IsNull() && !state.Sequence.IsUnknown() {
		stateValues["Sequence"] = fwValueToInterface(state.Sequence, map[string]fwWireField{}, false)
	}
	if !state.SiteId.IsNull() && !state.SiteId.IsUnknown() {
		stateValues["SiteId"] = fwValueToInterface(state.SiteId, map[string]fwWireField{}, false)
	}
	if !state.SiteVersion.IsNull() && !state.SiteVersion.IsUnknown() {
		stateValues["SiteVersion"] = fwValueToInterface(state.SiteVersion, map[string]fwWireField{}, false)
	}
	updated0 := false
	if fwPlannedChanged(plan.ErrorPagesRedirect, state.ErrorPagesRedirect) {
		updated0 = true
	}
	if fwPlannedChanged(plan.RuleEnable, state.RuleEnable) {
		updated0 = true
	}
	if fwPlannedChanged(plan.RuleName, state.RuleName) {
		updated0 = true
	}
	if fwPlannedChanged(plan.Rule, state.Rule) {
		updated0 = true
	}
	if updated0 {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"SiteId":      state.Id.ValueString(),
			"ClientToken": buildClientToken("UpdateErrorPagesRedirectRule"),
		}
		parts0 := strings.Split(state.Id.ValueString(), ":")
		if len(parts0) > 0 {
			request0["SiteId"] = parts0[0]
		}
		if len(parts0) > 1 {
			request0["ConfigId"] = parts0[1]
		}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, operationValues0, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.ConfigId", Target: "ConfigId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Sequence", Target: "Sequence", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.SiteId", Target: "SiteId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ErrorPagesRedirect[*].StatusCode", Target: "ErrorPagesRedirect[*].StatusCode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ErrorPagesRedirect[*].TargetURL", Target: "ErrorPagesRedirect[*].TargetURL", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RuleEnable", Target: "RuleEnable", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RuleName", Target: "RuleName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Rule", Target: "Rule", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues0, stateValues, []string{"$.SiteId", "$.ConfigId", "$.ErrorPagesRedirect[*].StatusCode", "$.ErrorPagesRedirect[*].TargetURL"}), false); err != nil {
			resp.Diagnostics.AddError("Building UpdateErrorPagesRedirectRule request", err.Error())
			return
		}
		if _, err := retryRpcPost(ctx, r.Client(), "ESA", "2024-09-10", "UpdateErrorPagesRedirectRule", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_esa_error_pages_redirect_rule: UpdateErrorPagesRedirectRule", err.Error())
			return
		}
	}
	object, err := DescribeErrorPagesRedirectRule(ctx, r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, errorPagesRedirectRuleState(object, &plan))...)
}

func (r *errorPagesRedirectRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state errorPagesRedirectRuleResourceModel
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
	if !state.ConfigId.IsNull() && !state.ConfigId.IsUnknown() {
		deleteValues["ConfigId"] = fwValueToInterface(state.ConfigId, map[string]fwWireField{}, false)
	}
	if !state.ConfigType.IsNull() && !state.ConfigType.IsUnknown() {
		deleteValues["ConfigType"] = fwValueToInterface(state.ConfigType, map[string]fwWireField{}, false)
	}
	if !state.ErrorPagesRedirect.IsNull() && !state.ErrorPagesRedirect.IsUnknown() {
		deleteValues["ErrorPagesRedirect"] = fwValueToInterface(state.ErrorPagesRedirect, map[string]fwWireField{"target_url": {Name: "TargetURL", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "status_code": {Name: "StatusCode", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.Rule.IsNull() && !state.Rule.IsUnknown() {
		deleteValues["Rule"] = fwValueToInterface(state.Rule, map[string]fwWireField{}, false)
	}
	if !state.RuleEnable.IsNull() && !state.RuleEnable.IsUnknown() {
		deleteValues["RuleEnable"] = fwValueToInterface(state.RuleEnable, map[string]fwWireField{}, false)
	}
	if !state.RuleName.IsNull() && !state.RuleName.IsUnknown() {
		deleteValues["RuleName"] = fwValueToInterface(state.RuleName, map[string]fwWireField{}, false)
	}
	if !state.Sequence.IsNull() && !state.Sequence.IsUnknown() {
		deleteValues["Sequence"] = fwValueToInterface(state.Sequence, map[string]fwWireField{}, false)
	}
	if !state.SiteId.IsNull() && !state.SiteId.IsUnknown() {
		deleteValues["SiteId"] = fwValueToInterface(state.SiteId, map[string]fwWireField{}, false)
	}
	if !state.SiteVersion.IsNull() && !state.SiteVersion.IsUnknown() {
		deleteValues["SiteVersion"] = fwValueToInterface(state.SiteVersion, map[string]fwWireField{}, false)
	}
	identityParts := strings.Split(state.Id.ValueString(), ":")
	if len(identityParts) > 0 {
		deleteValues["SiteId"] = identityParts[0]
	}
	if len(identityParts) > 1 {
		deleteValues["ConfigId"] = identityParts[1]
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
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteErrorPagesRedirectRule")}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.ConfigId", Target: "ConfigId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.SiteId", Target: "SiteId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DeleteErrorPagesRedirectRule request", err.Error())
			return
		}
		if _, err := retryRpcPost(ctx, r.Client(), "ESA", "2024-09-10", "DeleteErrorPagesRedirectRule", request0, operationTimeout); err != nil && !isExpectedErrors(err, []string{"ResourceNotFound.ErrorPagesRedirectRule", "ResourceNotFound"}) {
			resp.Diagnostics.AddError("Deleting alicloud_esa_error_pages_redirect_rule: DeleteErrorPagesRedirectRule", err.Error())
			return
		}
	}
}

func (r *errorPagesRedirectRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// errorPagesRedirectRuleState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func errorPagesRedirectRuleState(object map[string]interface{}, prior *errorPagesRedirectRuleResourceModel) errorPagesRedirectRuleResourceModel {
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
		str("SiteId").ValueString(),
		str("ConfigId").ValueString(),
	}
	m := errorPagesRedirectRuleResourceModel{
		FwTimeouts:         types.ObjectNull(fwTimeoutTypes()),
		Id:                 types.StringValue(strings.Join(_cidParts, ":")),
		ConfigId:           intVal("ConfigId"),
		ConfigType:         str("ConfigType"),
		ErrorPagesRedirect: errorPagesRedirectRuleErrorPagesRedirectToTypeList(lookup("ErrorPagesRedirect")),
		Rule:               str("Rule"),
		RuleEnable:         str("RuleEnable"),
		RuleName:           str("RuleName"),
		Sequence:           intVal("Sequence"),
		SiteId:             intVal("SiteId"),
		SiteVersion:        intVal("SiteVersion"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		m.ErrorPagesRedirect = fwPreserveRequestValues(m.ErrorPagesRedirect, prior.ErrorPagesRedirect, map[string]fwWireField{"target_url": {Name: "TargetURL", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "status_code": {Name: "StatusCode", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}).(types.List)
	}
	return m
}
