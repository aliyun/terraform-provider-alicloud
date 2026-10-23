// Package eflo. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package eflo

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const vccGrantRuleTypeName = "alicloud_eflo_vcc_grant_rule"

var (
	_ resource.Resource                = &vccGrantRuleResource{}
	_ resource.ResourceWithImportState = &vccGrantRuleResource{}
)

func NewVccGrantRuleResource() resource.Resource {
	return &vccGrantRuleResource{}
}

type vccGrantRuleResource struct {
	fwadapt.ResourceBase
}

type vccGrantRuleResourceModel struct {
	Id            types.String `tfsdk:"id"`
	FwTimeouts    types.Object `tfsdk:"timeouts"`
	CreateTime    types.String `tfsdk:"create_time"`
	ErId          types.String `tfsdk:"er_id"`
	GrantRuleId   types.String `tfsdk:"grant_rule_id"`
	GrantTenantId types.String `tfsdk:"grant_tenant_id"`
	InstanceId    types.String `tfsdk:"instance_id"`
	InstanceName  types.String `tfsdk:"instance_name"`
	Product       types.String `tfsdk:"product"`
	TenantId      types.String `tfsdk:"tenant_id"`
	Used          types.Bool   `tfsdk:"used"`
}

func (r *vccGrantRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = vccGrantRuleTypeName
}

func (r *vccGrantRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time",
			},
			"er_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The ID of the ER instance",
			},
			"grant_rule_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the authorized resource instance",
			},
			"grant_tenant_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Authorized tenant ID",
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Network instance ID",
			},
			"instance_name": schema.StringAttribute{
				Computed:    true,
				Description: "Network instance name",
			},
			"product": schema.StringAttribute{
				Computed:    true,
				Description: "Network product CODE: VPD, VCC",
			},
			"tenant_id": schema.StringAttribute{
				Computed:    true,
				Description: "Tenant ID",
			},
			"used": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the current authorization information has been used; The default is false",
			},
		},
	}
}

func (r *vccGrantRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vccGrantRuleResourceModel
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
		"ClientToken": buildClientToken("CreateVccGrantRule"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		createValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.ErId.IsNull() && !plan.ErId.IsUnknown() {
		createValues["ErId"] = fwValueToInterface(plan.ErId, map[string]fwWireField{}, false)
	}
	if !plan.GrantRuleId.IsNull() && !plan.GrantRuleId.IsUnknown() {
		createValues["GrantRuleId"] = fwValueToInterface(plan.GrantRuleId, map[string]fwWireField{}, false)
	}
	if !plan.GrantTenantId.IsNull() && !plan.GrantTenantId.IsUnknown() {
		createValues["GrantTenantId"] = fwValueToInterface(plan.GrantTenantId, map[string]fwWireField{}, false)
	}
	if !plan.InstanceId.IsNull() && !plan.InstanceId.IsUnknown() {
		createValues["InstanceId"] = fwValueToInterface(plan.InstanceId, map[string]fwWireField{}, false)
	}
	if !plan.InstanceName.IsNull() && !plan.InstanceName.IsUnknown() {
		createValues["InstanceName"] = fwValueToInterface(plan.InstanceName, map[string]fwWireField{}, false)
	}
	if !plan.Product.IsNull() && !plan.Product.IsUnknown() {
		createValues["Product"] = fwValueToInterface(plan.Product, map[string]fwWireField{}, false)
	}
	if !plan.TenantId.IsNull() && !plan.TenantId.IsUnknown() {
		createValues["TenantId"] = fwValueToInterface(plan.TenantId, map[string]fwWireField{}, false)
	}
	if !plan.Used.IsNull() && !plan.Used.IsUnknown() {
		createValues["Used"] = fwValueToInterface(plan.Used, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.InstanceId", Target: "InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ErId", Target: "ErId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.GrantTenantId", Target: "GrantTenantId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(r.Client(), "eflo", "2022-05-30", "CreateVccGrantRule", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVccGrantRule", vccGrantRuleTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateVccGrantRule", vccGrantRuleTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(request, fwMappingPath("InstanceId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_eflo_vcc_grant_rule: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	identityValues1 := fwWalkMappedPath(response, fwMappingPath("$.Content.GrantRuleId"), nil)
	if len(identityValues1) == 0 || identityValues1[0].Value == nil || fmt.Sprint(identityValues1[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_eflo_vcc_grant_rule: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues1[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVccGrantRule", vccGrantRuleTypeName),
			fmt.Sprintf("The response carried no InstanceId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.CreateTime = fwKnownValue(trackedPlan.CreateTime).(types.String)
	trackedPlan.ErId = fwKnownValue(trackedPlan.ErId).(types.String)
	trackedPlan.GrantRuleId = fwKnownValue(trackedPlan.GrantRuleId).(types.String)
	trackedPlan.GrantTenantId = fwKnownValue(trackedPlan.GrantTenantId).(types.String)
	trackedPlan.InstanceId = fwKnownValue(trackedPlan.InstanceId).(types.String)
	trackedPlan.InstanceName = fwKnownValue(trackedPlan.InstanceName).(types.String)
	trackedPlan.Product = fwKnownValue(trackedPlan.Product).(types.String)
	trackedPlan.TenantId = fwKnownValue(trackedPlan.TenantId).(types.String)
	trackedPlan.Used = fwKnownValue(trackedPlan.Used).(types.Bool)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := vccGrantRuleGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", vccGrantRuleTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", vccGrantRuleTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", vccGrantRuleTypeName),
			fmt.Sprintf("The resource %s is not visible via GetVccGrantRule right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, vccGrantRuleState(object, &plan))...)
}

func (r *vccGrantRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vccGrantRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := vccGrantRuleGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetVccGrantRule", vccGrantRuleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetVccGrantRule", vccGrantRuleTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, vccGrantRuleState(object, &state))...)
}

func (r *vccGrantRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state vccGrantRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *vccGrantRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vccGrantRuleResourceModel
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
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		deleteValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.ErId.IsNull() && !state.ErId.IsUnknown() {
		deleteValues["ErId"] = fwValueToInterface(state.ErId, map[string]fwWireField{}, false)
	}
	if !state.GrantRuleId.IsNull() && !state.GrantRuleId.IsUnknown() {
		deleteValues["GrantRuleId"] = fwValueToInterface(state.GrantRuleId, map[string]fwWireField{}, false)
	}
	if !state.GrantTenantId.IsNull() && !state.GrantTenantId.IsUnknown() {
		deleteValues["GrantTenantId"] = fwValueToInterface(state.GrantTenantId, map[string]fwWireField{}, false)
	}
	if !state.InstanceId.IsNull() && !state.InstanceId.IsUnknown() {
		deleteValues["InstanceId"] = fwValueToInterface(state.InstanceId, map[string]fwWireField{}, false)
	}
	if !state.InstanceName.IsNull() && !state.InstanceName.IsUnknown() {
		deleteValues["InstanceName"] = fwValueToInterface(state.InstanceName, map[string]fwWireField{}, false)
	}
	if !state.Product.IsNull() && !state.Product.IsUnknown() {
		deleteValues["Product"] = fwValueToInterface(state.Product, map[string]fwWireField{}, false)
	}
	if !state.TenantId.IsNull() && !state.TenantId.IsUnknown() {
		deleteValues["TenantId"] = fwValueToInterface(state.TenantId, map[string]fwWireField{}, false)
	}
	if !state.Used.IsNull() && !state.Used.IsUnknown() {
		deleteValues["Used"] = fwValueToInterface(state.Used, map[string]fwWireField{}, false)
	}
	identityParts := strings.Split(state.Id.ValueString(), ":")
	if len(identityParts) > 0 {
		deleteValues["InstanceId"] = identityParts[0]
	}
	if len(identityParts) > 1 {
		deleteValues["GrantRuleId"] = identityParts[1]
	}
	needsRefresh := false
	if needsRefresh {
		object, err := vccGrantRuleGet(r.Client(), state.Id.ValueString())
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
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteVccGrantRule")}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.InstanceId", Target: "InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.GrantRuleId", Target: "GrantRuleId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.ErId", Target: "ErId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DeleteVccGrantRule request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "eflo", "2022-05-30", "DeleteVccGrantRule", request0, operationTimeout); err != nil && !isExpectedErrors(err, []string{"1003"}) {
			resp.Diagnostics.AddError("Deleting alicloud_eflo_vcc_grant_rule: DeleteVccGrantRule", err.Error())
			return
		}
	}
}

func (r *vccGrantRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// vccGrantRuleState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func vccGrantRuleState(object map[string]interface{}, prior *vccGrantRuleResourceModel) vccGrantRuleResourceModel {
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
	boolVal := func(key string) types.Bool {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.BoolNull()
			}
			cur, ok = m[seg]
			if !ok {
				return types.BoolNull()
			}
		}
		if v, ok := cur.(bool); ok {
			return types.BoolValue(v)
		}
		return types.BoolNull()
	}
	// gap#11: composite-ID resource — build the id from all composite fields in the response,
	// joined by ":" (e.g. AccessGroupName:FileSystemType). The Get response carries every
	// part, so the state id round-trips through Read/Delete correctly (otherwise Read would
	// lose FileSystemType and DescribeAccessGroups would query the wrong namespace → 404).
	_cidParts := []string{
		str("InstanceId").ValueString(),
		str("GrantRuleId").ValueString(),
	}
	m := vccGrantRuleResourceModel{
		FwTimeouts:    types.ObjectNull(fwTimeoutTypes()),
		Id:            types.StringValue(strings.Join(_cidParts, ":")),
		CreateTime:    str("CreateTime"),
		ErId:          str("ErId"),
		GrantRuleId:   str("GrantRuleId"),
		GrantTenantId: str("GrantTenantId"),
		InstanceId:    str("InstanceId"),
		InstanceName:  str("InstanceName"),
		Product:       str("Product"),
		TenantId:      str("TenantId"),
		Used:          boolVal("Used"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
	}
	return m
}
