package ecs

import (
	"context"
	"fmt"

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
				PlanModifiers: []planmodifier.Int64{
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
				PlanModifiers: []planmodifier.String{
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
				PlanModifiers: []planmodifier.String{
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

	request := map[string]interface{}{
		"RegionId":        r.Client().RegionId,
		"CommittedAmount": plan.CommittedAmount.ValueString(),
		"ClientToken":     helper.BuildClientToken("PurchaseSavingPlanOffering"),
	}
	if !plan.PlanType.IsNull() {
		request["PlanType"] = plan.PlanType.ValueString()
	}
	if !plan.OfferingType.IsNull() {
		request["OfferingType"] = plan.OfferingType.ValueString()
	}
	if !plan.PurchaseMethod.IsNull() {
		request["PurchaseMethod"] = plan.PurchaseMethod.ValueString()
	}
	if !plan.ChargeType.IsNull() {
		request["ChargeType"] = plan.ChargeType.ValueString()
	}
	if !plan.Period.IsNull() {
		request["Period"] = plan.Period.ValueInt64()
	}
	if !plan.PeriodUnit.IsNull() {
		request["PeriodUnit"] = plan.PeriodUnit.ValueString()
	}
	if !plan.InstanceFamily.IsNull() {
		request["InstanceFamily"] = plan.InstanceFamily.ValueString()
	}
	if !plan.InstanceFamilySet.IsNull() {
		request["InstanceFamilySet"] = plan.InstanceFamilySet.ValueString()
	}
	if !plan.SavingPlanName.IsNull() {
		request["Name"] = plan.SavingPlanName.ValueString()
	}
	if !plan.Description.IsNull() {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.StartTime.IsNull() {
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

	// Financial-commitment safeguard: persist the purchased ID to state BEFORE
	// refreshing. If the subsequent Query fails or the plan is not yet visible,
	// the ID stays in state so the next apply refreshes instead of re-purchasing
	// — never lose a paid Saving Plan ID.
	persisted := plan
	persisted.Id = types.StringValue(id)
	persisted.SavingPlanId = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &persisted)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := querySavingPlan(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Reading %s: calling QuerySavingsPlansInstance", savingPlanTypeName),
			fmt.Sprintf("Purchase succeeded (saving_plan_id=%s) but the initial refresh failed: %s. The ID has been persisted to state; the next apply will refresh without re-purchasing.", id, err.Error()),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Creating %s: calling PurchaseSavingPlanOffering", savingPlanTypeName),
			fmt.Sprintf("Purchase succeeded (saving_plan_id=%s) but the saving plan is not yet visible via QuerySavingsPlansInstance. The ID has been persisted to state; the next apply will refresh without re-purchasing.", id),
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
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling QuerySavingsPlansInstance", savingPlanTypeName),
			errs.WrapErrorf(err, "Reading %s: calling QuerySavingsPlansInstance", savingPlanTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
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
// preserved from the prior model so the state does not regress.
func savingPlanState(object map[string]interface{}, prior savingPlanResourceModel, id string) savingPlanResourceModel {
	apiStrOrPrior := func(key string, priorVal types.String) types.String {
		if v, ok := object[key]; ok && v != nil {
			if s, ok := v.(string); ok && s != "" {
				return types.StringValue(s)
			}
			return types.StringValue(fmt.Sprint(v))
		}
		return priorVal
	}

	period := types.Int64Null()
	if v, ok := object["Cycle"]; ok {
		period = periodFromCycle(v)
	}
	if period.IsNull() {
		period = prior.Period
	}

	return savingPlanResourceModel{
		Id:                types.StringValue(id),
		CommittedAmount:   apiStrOrPrior("PoolValue", prior.CommittedAmount),
		PlanType:          apiStrOrPrior("SavingsType", prior.PlanType),
		OfferingType:      apiStrOrPrior("PayMode", prior.OfferingType),
		PurchaseMethod:    prior.PurchaseMethod,
		ChargeType:        prior.ChargeType,
		Period:            period,
		PeriodUnit:        prior.PeriodUnit,
		InstanceFamily:    apiStrOrPrior("InstanceFamily", prior.InstanceFamily),
		InstanceFamilySet: prior.InstanceFamilySet,
		SavingPlanName:    prior.SavingPlanName,
		Description:       prior.Description,
		StartTime:         apiStrOrPrior("StartTime", prior.StartTime),
		RegionId:          apiStrOrPrior("Region", prior.RegionId),
		SavingPlanId:      types.StringValue(id),
		CreateTime:        apiStrOrPrior("StartTimestamp", prior.CreateTime),
		PaymentType:       apiStrOrPrior("PaymentType", prior.PaymentType),
	}
}
