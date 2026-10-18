package sas

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const customCheckStandardPolicyTypeName = "alicloud_threat_detection_custom_check_standard_policy"

var (
	_ resource.Resource                = (*customCheckStandardPolicyResource)(nil)
	_ resource.ResourceWithImportState = (*customCheckStandardPolicyResource)(nil)
)

func NewCustomCheckStandardPolicyResource() resource.Resource {
	return &customCheckStandardPolicyResource{}
}

type customCheckStandardPolicyResource struct {
	fwadapt.ResourceBase
}

type customCheckStandardPolicyModel struct {
	Id                types.String `tfsdk:"id"`
	PolicyId          types.String `tfsdk:"policy_id"`
	PolicyShowName    types.String `tfsdk:"policy_show_name"`
	PolicyType        types.String `tfsdk:"policy_type"`
	DependentPolicyId types.String `tfsdk:"dependent_policy_id"`
	Type              types.String `tfsdk:"type"`
	CheckType         types.String `tfsdk:"check_type"`
}

func (r *customCheckStandardPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = customCheckStandardPolicyTypeName
}

func (r *customCheckStandardPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Threat Detection Custom Check Standard Policy resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID in Terraform, formatted as `<policy_id>:<policy_type>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the custom policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_show_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the custom policy.",
			},
			"policy_type": schema.StringAttribute{
				Required:    true,
				Description: "The policy category type for custom check rules. Valid values: `STANDARD`, `REQUIREMENT`, `SECTION`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("STANDARD", "REQUIREMENT", "SECTION"),
				},
			},
			"dependent_policy_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the parent policy. The dependency order from low to high is: Section -> Requirement -> Standard.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the major policy category, required when `policy_type` is `STANDARD`. Valid values: `AISPM`, `IDENTITY_PERMISSION`, `RISK`, `COMPLIANCE`.",
				Validators: []validator.String{
					stringvalidator.OneOf("AISPM", "IDENTITY_PERMISSION", "RISK", "COMPLIANCE"),
				},
			},
			"check_type": schema.StringAttribute{
				Computed:    true,
				Description: "The check type of the policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *customCheckStandardPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customCheckStandardPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"PolicyShowName": plan.PolicyShowName.ValueString(),
		"PolicyType":     plan.PolicyType.ValueString(),
	}
	if !plan.DependentPolicyId.IsNull() && plan.DependentPolicyId.ValueString() != "" {
		request["DependentPolicyId"] = plan.DependentPolicyId.ValueString()
	}
	if !plan.Type.IsNull() && plan.Type.ValueString() != "" {
		request["Type"] = plan.Type.ValueString()
	}

	response, err := r.Client().RpcPost(sasProduct, sasVersion, "CreateCheckPolicy", nil, request, false)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateCheckPolicy", customCheckStandardPolicyTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateCheckPolicy", customCheckStandardPolicyTypeName).Error(),
		)
		return
	}

	policyId := policyIdFromCreateResponse(response)
	if policyId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateCheckPolicy", customCheckStandardPolicyTypeName),
			fmt.Sprintf("The CreateCheckPolicy response carried no CheckCustomPolicy.PolicyId: %v", response),
		)
		return
	}

	object, err := describeCustomCheckStandardPolicy(r.Client(), fmt.Sprintf("%s:%s", policyId, plan.PolicyType.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListCheckPolicies", customCheckStandardPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListCheckPolicies after create", customCheckStandardPolicyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateCheckPolicy", customCheckStandardPolicyTypeName),
			fmt.Sprintf("The policy %s:%s is not visible via ListCheckPolicies right after creation", policyId, plan.PolicyType.ValueString()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, policyState(object))...)
}

func (r *customCheckStandardPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state customCheckStandardPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeCustomCheckStandardPolicy(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListCheckPolicies", customCheckStandardPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListCheckPolicies", customCheckStandardPolicyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, policyState(object))...)
}

func (r *customCheckStandardPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state customCheckStandardPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId, policyType, err := parsePolicyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s", customCheckStandardPolicyTypeName),
			err.Error(),
		)
		return
	}

	// UpdateCheckPolicy only when an updatable field changed; dependent_policy_id
	// is ForceNew and policy_type is ForceNew, so only policy_show_name and type
	// reach Update.
	changed := !plan.PolicyShowName.Equal(state.PolicyShowName) || !plan.Type.Equal(state.Type)
	if changed {
		request := map[string]interface{}{
			"PolicyId":       policyId,
			"PolicyType":     policyType,
			"PolicyShowName": plan.PolicyShowName.ValueString(),
		}
		if !plan.Type.IsNull() && plan.Type.ValueString() != "" {
			request["Type"] = plan.Type.ValueString()
		}
		if _, err := r.Client().RpcPost(sasProduct, sasVersion, "UpdateCheckPolicy", nil, request, false); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateCheckPolicy", customCheckStandardPolicyTypeName),
				errs.WrapErrorf(err, "Updating %s: calling UpdateCheckPolicy", customCheckStandardPolicyTypeName).Error(),
			)
			return
		}
	}

	object, err := describeCustomCheckStandardPolicy(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListCheckPolicies", customCheckStandardPolicyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListCheckPolicies after update", customCheckStandardPolicyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, policyState(object))...)
}

func (r *customCheckStandardPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customCheckStandardPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId, policyType, err := parsePolicyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s", customCheckStandardPolicyTypeName),
			err.Error(),
		)
		return
	}

	_, err = r.Client().RpcPost(sasProduct, sasVersion, "DeleteCheckPolicy", nil, map[string]interface{}{
		"PolicyIds.1": policyId,
		"PolicyType":  policyType,
	}, false)
	if err != nil && !isPolicyNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteCheckPolicy", customCheckStandardPolicyTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteCheckPolicy", customCheckStandardPolicyTypeName).Error(),
		)
		return
	}
}

func (r *customCheckStandardPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// policyIdFromCreateResponse extracts PolicyId from the CheckCustomPolicy
// object of a CreateCheckPolicy response, returning "" when absent. The id is
// numeric, so it flows through policyFieldString to avoid float64 scientific
// notation.
func policyIdFromCreateResponse(response map[string]interface{}) string {
	if response == nil {
		return ""
	}
	if cp, ok := response["CheckCustomPolicy"].(map[string]interface{}); ok {
		return policyFieldString(cp, "PolicyId")
	}
	return ""
}

// policyState maps a ListCheckPolicies policy object onto the resource model.
// Empty values become null so optional/omitted API fields read back as absent
// rather than "".
func policyState(object map[string]interface{}) customCheckStandardPolicyModel {
	str := func(key string) types.String {
		if s := policyFieldString(object, key); s != "" {
			return types.StringValue(s)
		}
		return types.StringNull()
	}

	m := customCheckStandardPolicyModel{
		PolicyId:       str("PolicyId"),
		PolicyShowName: str("PolicyShowName"),
		PolicyType:     str("PolicyType"),
		CheckType:      str("CheckType"),
		Type:           str("Type"),
	}

	policyId := str("PolicyId")
	policyType := str("PolicyType")
	if !policyId.IsNull() && !policyType.IsNull() {
		m.Id = types.StringValue(fmt.Sprintf("%s:%s", policyId.ValueString(), policyType.ValueString()))
	}

	// DependentPolicyId comes back as 0 for SECTION-level policies; treat 0 and
	// empty as absent, matching the legacy SasService Describe behavior.
	if dependent := policyFieldString(object, "DependentPolicyId"); dependent != "" && dependent != "0" {
		m.DependentPolicyId = types.StringValue(dependent)
	}

	return m
}

// isPolicyNotExist reports whether err is a policy-not-found error from
// DeleteCheckPolicy, so Delete is idempotent when the policy is already gone.
func isPolicyNotExist(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "notfound") || strings.Contains(msg, "not exist") || strings.Contains(msg, "does not exist")
}
