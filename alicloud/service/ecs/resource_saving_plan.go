package ecs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const savingPlanTypeName = "alicloud_ecs_saving_plan"

var (
	_ resource.Resource                = &savingPlanResource{}
	_ resource.ResourceWithImportState = &savingPlanResource{}
)

func NewSavingPlanResource() resource.Resource {
	return &savingPlanResource{}
}

type savingPlanResource struct {
	fwadapt.ResourceBase
}

type savingPlanResourceModel struct {
	Id                types.String `tfsdk:"id"`
	CommittedAmount   types.String `tfsdk:"committed_amount"`
	PlanType          types.String `tfsdk:"plan_type"`
	OfferingType      types.String `tfsdk:"offering_type"`
	PurchaseMethod    types.String `tfsdk:"purchase_method"`
	ChargeType        types.String `tfsdk:"charge_type"`
	Period            types.Int64  `tfsdk:"period"`
	PeriodUnit        types.String `tfsdk:"period_unit"`
	InstanceFamily    types.String `tfsdk:"instance_family"`
	InstanceFamilySet types.String `tfsdk:"instance_family_set"`
	SavingPlanName    types.String `tfsdk:"saving_plan_name"`
	Description       types.String `tfsdk:"description"`
	StartTime         types.String `tfsdk:"start_time"`
	RegionId          types.String `tfsdk:"region_id"`
	SavingPlanId      types.String `tfsdk:"saving_plan_id"`
	CreateTime        types.String `tfsdk:"create_time"`
	PaymentType       types.String `tfsdk:"payment_type"`
}

func (r *savingPlanResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = savingPlanTypeName
}

func (r *savingPlanResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an ECS Saving Plan resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `saving_plan_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"committed_amount": schema.StringAttribute{
				Required:    true,
				Description: "The hourly commitment amount of the saving plan, in USD per hour.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plan_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("EcsCompute"),
				Description: "The type of the saving plan. Valid values: `EcsCompute`, `General`. Default to `EcsCompute`.",
				Validators: []validator.String{
					stringvalidator.OneOf("EcsCompute", "General"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"offering_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("AllUpfront"),
				Description: "The payment method of the saving plan. Valid values: `AllUpfront`, `HalfUpfront`, `NoUpfront`. Default to `AllUpfront`.",
				Validators: []validator.String{
					stringvalidator.OneOf("AllUpfront", "HalfUpfront", "NoUpfront"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"purchase_method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ByInstanceFamily"),
				Description: "The purchase method of the saving plan. Valid values: `ByInstanceFamily`, `ByInstanceFamilySet`. Default to `ByInstanceFamily`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ByInstanceFamily", "ByInstanceFamilySet"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"charge_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("PrePaid"),
				Description: "The charge type. Default to `PrePaid`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"period": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The period of the saving plan.",
				// UseStateForUnknown must run before RequiresReplace so a refresh
				// that fills this Optional+Computed (no Default) attribute from the
				// API carries the prior state value into the plan instead of leaving
				// it unknown. Without it the next plan sees unknown != state and
				// schedules a replacement (re-purchase) on every refresh — the
				// defect the tfprotov6 next-plan test guards against. A user config
				// change yields a known plan value, so UseStateForUnknown is a no-op
				// and RequiresReplace still fires on a real config-vs-state diff.
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"period_unit": schema.StringAttribute{
				Optional:    true,
				Description: "The unit of the period. Valid values: `Year`, `Month`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_family": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The instance family. Required when `purchase_method` is `ByInstanceFamily`.",
				// Same Optional+Computed+RequiresReplace+no-Default shape as
				// period/start_time: UseStateForUnknown carries the refresh-filled
				// value into the plan so an API-returned instance_family does not
				// trigger a replacement on the next plan when the user omitted it.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_family_set": schema.StringAttribute{
				Optional:    true,
				Description: "The instance family set. Required when `purchase_method` is `ByInstanceFamilySet`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"saving_plan_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the saving plan.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the saving plan.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"start_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The effective start time of the saving plan.",
				// Same Optional+Computed+RequiresReplace+no-Default shape as
				// period: UseStateForUnknown carries the refresh-filled start_time
				// into the plan so an API-returned value does not trigger a
				// replacement (re-purchase) on the next plan when the user omitted
				// start_time from config.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the saving plan.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"saving_plan_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Saving Plan.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the saving plan.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"payment_type": schema.StringAttribute{
				Computed:    true,
				Description: "The payment type of the saving plan.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *savingPlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan savingPlanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Conditional-required validation before Purchase. The schema documents that
	// instance_family is required when purchase_method is ByInstanceFamily and
	// instance_family_set when ByInstanceFamilySet. instance_family is
	// Optional+Computed and instance_family_set is Optional, so a user who
	// selects a purchase method but omits the corresponding family leaves it
	// unknown or null at Create plan time. Framework's ValueString returns ""
	// for an unknown, so without this check the previous code would have sent an
	// empty InstanceFamily to an irreversible Purchase. For a financial
	// commitment, diagnose the misconfiguration BEFORE calling Purchase rather
	// than letting the API reject it (or worse, silently accepting an empty
	// family). purchase_method always carries its StaticString default or the
	// user-set value, so it is known here.
	switch plan.PurchaseMethod.ValueString() {
	case "ByInstanceFamily":
		if plan.InstanceFamily.IsUnknown() || plan.InstanceFamily.IsNull() {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: missing instance_family", savingPlanTypeName),
				"instance_family is required when purchase_method is ByInstanceFamily; set it to a valid instance family before applying.",
			)
			return
		}
	case "ByInstanceFamilySet":
		if plan.InstanceFamilySet.IsUnknown() || plan.InstanceFamilySet.IsNull() {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: missing instance_family_set", savingPlanTypeName),
				"instance_family_set is required when purchase_method is ByInstanceFamilySet; set it to a valid instance family set before applying.",
			)
			return
		}
	}

	request := map[string]interface{}{
		"RegionId":        r.Client().RegionId,
		"CommittedAmount": plan.CommittedAmount.ValueString(),
		"ClientToken":     helper.BuildClientToken("PurchaseSavingPlanOffering"),
	}
	// Only a concrete known non-null value is sent to Purchase. An unknown
	// Computed attribute (period, start_time — Optional+Computed with no
	// StaticString default, so an omitted-by-user value arrives as unknown at
	// Create plan time) must be omitted, not sent: Framework's ValueString/
	// ValueInt64 return "" / 0 for unknown, and IsNull is false for unknown, so
	// the previous `!IsNull()` guard would have silently coerced an unknown
	// period to 0 and an unknown start_time to "" in the Purchase request.
	// Default-bearing attributes (plan_type, offering_type, purchase_method,
	// charge_type) are known at plan time, so they are sent normally.
	if isKnownString(plan.PlanType) {
		request["PlanType"] = plan.PlanType.ValueString()
	}
	if isKnownString(plan.OfferingType) {
		request["OfferingType"] = plan.OfferingType.ValueString()
	}
	if isKnownString(plan.PurchaseMethod) {
		request["PurchaseMethod"] = plan.PurchaseMethod.ValueString()
	}
	if isKnownString(plan.ChargeType) {
		request["ChargeType"] = plan.ChargeType.ValueString()
	}
	if isKnownInt64(plan.Period) {
		request["Period"] = plan.Period.ValueInt64()
	}
	if isKnownString(plan.PeriodUnit) {
		request["PeriodUnit"] = plan.PeriodUnit.ValueString()
	}
	if isKnownString(plan.InstanceFamily) {
		request["InstanceFamily"] = plan.InstanceFamily.ValueString()
	}
	if isKnownString(plan.InstanceFamilySet) {
		request["InstanceFamilySet"] = plan.InstanceFamilySet.ValueString()
	}
	if isKnownString(plan.SavingPlanName) {
		request["Name"] = plan.SavingPlanName.ValueString()
	}
	if isKnownString(plan.Description) {
		request["Description"] = plan.Description.ValueString()
	}
	if isKnownString(plan.StartTime) {
		request["StartTime"] = plan.StartTime.ValueString()
	}

	response, err := r.Client().RpcPost("Ecs", "2016-03-14", "PurchaseSavingPlanOffering", nil, request, true)
	logs.AddDebug("PurchaseSavingPlanOffering", response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PurchaseSavingPlanOffering", savingPlanTypeName),
			errs.WrapErrorf(err, "Creating %s: calling PurchaseSavingPlanOffering", savingPlanTypeName).Error(),
		)
		return
	}

	idSetsRaw, ok := response["SavingPlanIdSets"]
	if !ok || idSetsRaw == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PurchaseSavingPlanOffering", savingPlanTypeName),
			fmt.Sprintf("The response carried no SavingPlanIdSets: %v", response),
		)
		return
	}
	idSets, ok := idSetsRaw.([]interface{})
	if !ok || len(idSets) == 0 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PurchaseSavingPlanOffering", savingPlanTypeName),
			fmt.Sprintf("SavingPlanIdSets is empty in response: %v", response),
		)
		return
	}
	// Strict type check: SavingPlanIdSets[0] must be a non-empty string.
	// fmt.Sprint on a nil or non-string entry would produce "<nil>" and mask
	// a malformed response — the ID must be real before it is trusted.
	firstID, ok := idSets[0].(string)
	if !ok || firstID == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PurchaseSavingPlanOffering", savingPlanTypeName),
			fmt.Sprintf("SavingPlanIdSets[0] is not a non-empty string: %v", idSets[0]),
		)
		return
	}
	id := firstID

	// Financial-commitment safeguard: persist the purchased ID to a fully-known
	// state BEFORE refreshing. If the subsequent Query fails or the plan is not
	// yet visible, the ID stays in state so the next apply refreshes instead of
	// re-purchasing — never lose a paid Saving Plan ID. The persisted state is
	// fully known (no unknown Computed values): a real Create plan carries
	// unknown Computed attributes (region_id, create_time, payment_type and any
	// omitted Optional+Computed field such as period, instance_family or
	// start_time) that the refresh would normally fill; when the refresh misses,
	// returning those unknowns in NewState would violate the Framework
	// ApplyResourceChange protocol. savingPlanPersistedModel coerces every
	// non-concrete attribute to null so the persisted state is legal while the
	// next refresh fills the concrete values once the plan becomes visible.
	persisted := savingPlanPersistedModel(plan, id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &persisted)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := querySavingPlan(r.Client(), id)
	if err != nil {
		// Purchase just succeeded, so every non-found query outcome — a
		// not-yet-visible empty Items (eventual consistency), an authoritative
		// NotFound (the Query API has not propagated the purchase yet), a
		// transient server error, or a malformed response — must keep the paid
		// ID rather than re-purchase. A warning (not an error) signals the
		// refresh miss so the next apply retries the refresh.
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Creating %s: calling PurchaseSavingPlanOffering", savingPlanTypeName),
			fmt.Sprintf("Purchase succeeded (saving_plan_id=%s) but the initial refresh did not return the plan: %s. The ID has been persisted to a fully-known state; the next apply will refresh without re-purchasing.", id, err.Error()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, savingPlanState(object, plan, id))...)
}

func (r *savingPlanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state savingPlanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := querySavingPlan(r.Client(), state.Id.ValueString())
	if errors.Is(err, errSavingPlanQueryNotFound) {
		// The query returned a not-found-class signal (a precise HTTP 404 or
		// Code=="NotFound"). The BssOpenApi/2017-12-14 QuerySavingsPlansInstance
		// contract does NOT guarantee that such a signal means THIS saving plan
		// instance is gone — a 404 may be routing/endpoint/scoped to a different
		// object, and the official error table documents no instance-NotFound
		// code — so it is not an authorisation to drop a paid Saving Plan. A
		// purchase that has not yet propagated can return the SAME not-found
		// signal immediately after Create (Create keeps the ID there for exactly
		// this reason), and a previously-visible plan can return a generic 404 for
		// reasons unrelated to the instance existing; prior visibility only proves
		// past existence, not that this 404 is the instance being gone. The paid
		// ID is therefore ALWAYS retained with a warning across consecutive
		// same-signal Reads — regardless of whether it was previously visible —
		// until the query succeeds and refreshes real data, or the resource is
		// removed from configuration. No query not-found signal ever drops a paid
		// financial commitment.
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Reading %s: calling QuerySavingsPlansInstance", savingPlanTypeName),
			fmt.Sprintf("Saving plan %s was reported not found by QuerySavingsPlansInstance (404 / NotFound). The BssOpenApi/2017-12-14 contract does not guarantee this signal means the instance is gone, so the paid ID is preserved in state rather than dropped on a not-found signal that is not pinned to this instance. Remove the resource from configuration if it was deleted out of band.", state.Id.ValueString()),
		)
		return
	}
	if errors.Is(err, errSavingPlanNotVisible) {
		// A legitimate success whose Items list does not contain the plan. This
		// is ambiguous — not yet visible (eventual consistency) or genuinely
		// absent — and the contract cannot tell. Preserve the paid ID and warn so
		// the next apply refreshes without re-purchasing; the resource is removed
		// only by removing it from configuration.
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Reading %s: calling QuerySavingsPlansInstance", savingPlanTypeName),
			fmt.Sprintf("Saving plan %s is not visible via QuerySavingsPlansInstance (legitimate success with no matching plan); the ID is preserved in state for the next apply. Remove the resource from configuration if it was deleted out of band.", state.Id.ValueString()),
		)
		return
	}
	if err != nil {
		// A business-failure envelope, a malformed response, or a transient
		// server/transport error: surface it as an error and preserve the paid ID
		// in state for the next refresh rather than RemoveResource'ing it.
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling QuerySavingsPlansInstance", savingPlanTypeName),
			errs.WrapErrorf(err, "Reading %s: calling QuerySavingsPlansInstance", savingPlanTypeName).Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, savingPlanState(object, state, state.Id.ValueString()))...)
}

func (r *savingPlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes are ForceNew (RequiresReplace); Terraform replaces the
	// resource on any change rather than calling Update. The method exists to
	// satisfy the interface and refreshes state from the plan.
	var plan savingPlanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *savingPlanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Saving Plans are financial commitments and cannot be cancelled via API.
	// The resource is removed from Terraform state only; the cloud resource
	// remains and must be cancelled via the console if needed.
	var state savingPlanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *savingPlanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// savingPlanState maps the Saving Plan object from a QuerySavingsPlansInstance
// response onto the resource model. Fields the API does not return are
// preserved from the prior model so the state does not regress, EXCEPT that an
// unknown prior value is coerced to null: prior may be a Create plan whose
// Computed attributes (region_id, create_time, payment_type, period,
// instance_family, start_time) are still unknown, and returning an unknown in
// NewState would violate the Framework ApplyResourceChange protocol. The next
// refresh fills the concrete value once the plan is visible.
func savingPlanState(object map[string]interface{}, prior savingPlanResourceModel, id string) savingPlanResourceModel {
	apiStrOrPrior := func(key string, priorVal types.String) types.String {
		if v, ok := object[key]; ok && v != nil {
			if s, ok := v.(string); ok && s != "" {
				return types.StringValue(s)
			}
			return types.StringValue(fmt.Sprint(v))
		}
		return knownStringOrNull(priorVal)
	}

	period := types.Int64Null()
	if v, ok := object["Cycle"]; ok {
		period = periodFromCycle(v)
	}
	if period.IsNull() {
		period = knownInt64OrNull(prior.Period)
	}

	return savingPlanResourceModel{
		Id:                types.StringValue(id),
		CommittedAmount:   apiStrOrPrior("PoolValue", prior.CommittedAmount),
		PlanType:          savingPlanPlanTypeFromQuery(object["SavingsType"], prior.PlanType),
		OfferingType:      savingPlanOfferingTypeFromQuery(object["PayMode"], prior.OfferingType),
		PurchaseMethod:    knownStringOrNull(prior.PurchaseMethod),
		ChargeType:        knownStringOrNull(prior.ChargeType),
		Period:            period,
		PeriodUnit:        knownStringOrNull(prior.PeriodUnit),
		InstanceFamily:    apiStrOrPrior("InstanceFamily", prior.InstanceFamily),
		InstanceFamilySet: knownStringOrNull(prior.InstanceFamilySet),
		SavingPlanName:    knownStringOrNull(prior.SavingPlanName),
		Description:       knownStringOrNull(prior.Description),
		StartTime:         apiStrOrPrior("StartTime", prior.StartTime),
		RegionId:          apiStrOrPrior("Region", prior.RegionId),
		SavingPlanId:      types.StringValue(id),
		CreateTime:        apiStrOrPrior("StartTimestamp", prior.CreateTime),
		PaymentType:       apiStrOrPrior("PaymentType", prior.PaymentType),
	}
}

// savingPlanPersistedModel builds a fully-known state model from a Create plan
// and the purchased ID, for the path where the post-purchase refresh did not
// return a full object (not yet visible, query error, or an authoritative
// NotFound immediately after a just-succeeded purchase). A real Framework
// Create plan carries unknown Computed attributes (region_id, create_time,
// payment_type and any omitted Optional+Computed field such as period,
// instance_family or start_time) that the refresh would normally fill. When the
// refresh misses, those unknowns cannot be returned in NewState — doing so
// would violate the ApplyResourceChange protocol (NewState for a newly-created
// resource must be fully known). Each attribute that is not a concrete known
// value in the plan is coerced to null so the persisted state is fully known
// while the purchased ID is retained; the next refresh fills the concrete
// values once the plan becomes visible.
func savingPlanPersistedModel(plan savingPlanResourceModel, id string) savingPlanResourceModel {
	return savingPlanResourceModel{
		Id:                types.StringValue(id),
		CommittedAmount:   knownStringOrNull(plan.CommittedAmount),
		PlanType:          knownStringOrNull(plan.PlanType),
		OfferingType:      knownStringOrNull(plan.OfferingType),
		PurchaseMethod:    knownStringOrNull(plan.PurchaseMethod),
		ChargeType:        knownStringOrNull(plan.ChargeType),
		Period:            knownInt64OrNull(plan.Period),
		PeriodUnit:        knownStringOrNull(plan.PeriodUnit),
		InstanceFamily:    knownStringOrNull(plan.InstanceFamily),
		InstanceFamilySet: knownStringOrNull(plan.InstanceFamilySet),
		SavingPlanName:    knownStringOrNull(plan.SavingPlanName),
		Description:       knownStringOrNull(plan.Description),
		StartTime:         knownStringOrNull(plan.StartTime),
		RegionId:          knownStringOrNull(plan.RegionId),
		SavingPlanId:      types.StringValue(id),
		CreateTime:        knownStringOrNull(plan.CreateTime),
		PaymentType:       knownStringOrNull(plan.PaymentType),
	}
}

// isKnownString reports whether v is a concrete known non-null string — the
// only state in which a value is sent to the cloud API. Null and unknown are
// both omitted so an unknown Computed attribute is never silently coerced to ""
// by Framework's ValueString (IsNull returns false for an unknown, so the bare
// !IsNull guard would have sent an empty value).
func isKnownString(v types.String) bool {
	return !v.IsNull() && !v.IsUnknown()
}

// isKnownInt64 is the Int64 counterpart of isKnownString.
func isKnownInt64(v types.Int64) bool {
	return !v.IsNull() && !v.IsUnknown()
}

// knownStringOrNull returns the value when it is a concrete known string (null
// is concrete and preserved) and types.StringNull() when it is unknown, so a
// model built for NewState never carries an unknown Computed value.
func knownStringOrNull(v types.String) types.String {
	if v.IsUnknown() {
		return types.StringNull()
	}
	return v
}

// knownInt64OrNull is the Int64 counterpart of knownStringOrNull.
func knownInt64OrNull(v types.Int64) types.Int64 {
	if v.IsUnknown() {
		return types.Int64Null()
	}
	return v
}

// savingPlanQueryStr returns the trimmed string form of a
// QuerySavingsPlansInstance response field value, or "" when the field is
// absent or nil. Non-string values (e.g. a float64 PoolValue) are stringified
// via fmt.Sprint, matching apiStr.
func savingPlanQueryStr(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

// savingPlanPlanTypeFromQuery translates the BSS QuerySavingsPlansInstance
// SavingsType field onto the resource plan_type schema enum. The Query API
// (BssOpenApi 2017-12-14) returns SavingsType in its own value space — "ECS"
// for an ECS compute plan and "universal"/"通用型" for a general plan — while
// the Purchase API (Ecs 2016-03-14 PurchaseSavingPlanOffering) and the
// resource schema use plan_type values "EcsCompute"/"General" (OneOf,
// RequiresReplace). Writing the raw BSS value into plan_type would violate the
// schema's OneOf and, because plan_type is RequiresReplace, schedule a
// replacement on the next plan whenever the API value differs from the
// configured value. The translation maps the known BSS forms to the schema
// forms (and tolerates the Purchase-form strings case-insensitively); an
// unrecognized value is NOT written into plan_type — the configured legal
// plan_type is retained (or null with no prior) so the BSS field's
// not-guaranteed semantics for unmapped values never produce an inconsistent
// result or a spurious replacement.
func savingPlanPlanTypeFromQuery(raw interface{}, prior types.String) types.String {
	switch strings.ToLower(savingPlanQueryStr(raw)) {
	case "ecs", "ecscompute":
		return types.StringValue("EcsCompute")
	case "universal", "general", "通用型", "通用":
		return types.StringValue("General")
	}
	return knownStringOrNull(prior)
}

// savingPlanOfferingTypeFromQuery translates the BSS PayMode field onto the
// resource offering_type schema enum. The Query API returns PayMode in its own
// value space — "total"/"half"/"zero" — while the resource schema uses
// offering_type values "AllUpfront"/"HalfUpfront"/"NoUpfront" (OneOf,
// RequiresReplace). Writing the raw BSS PayMode into offering_type would
// violate the schema's OneOf and schedule a RequiresReplace replacement.
// "total" (full prepayment) maps to AllUpfront, "half" to HalfUpfront, "zero"
// (no prepayment) to NoUpfront; an unrecognized value retains the configured
// legal offering_type rather than mis-changing it.
func savingPlanOfferingTypeFromQuery(raw interface{}, prior types.String) types.String {
	switch strings.ToLower(savingPlanQueryStr(raw)) {
	case "total", "allupfront":
		return types.StringValue("AllUpfront")
	case "half", "halfupfront":
		return types.StringValue("HalfUpfront")
	case "zero", "noupfront":
		return types.StringValue("NoUpfront")
	}
	return knownStringOrNull(prior)
}
