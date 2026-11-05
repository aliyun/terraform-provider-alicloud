package threatdetection

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var (
	_ resource.Resource                = &unknownThreatDetectStrategyResource{}
	_ resource.ResourceWithImportState = &unknownThreatDetectStrategyResource{}
)

func NewUnknownThreatDetectStrategyResource() resource.Resource {
	return &unknownThreatDetectStrategyResource{}
}

type unknownThreatDetectStrategyResource struct {
	fwadapt.ResourceBase
}

type unknownThreatDetectStrategyResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	AssetSelectionType    types.String `tfsdk:"asset_selection_type"`
	StudyMode             types.String `tfsdk:"study_mode"`
	DurationDaysAfterInit types.Int64  `tfsdk:"duration_days_after_init"`
	DurationDaysAfterStop types.Int64  `tfsdk:"duration_days_after_stop"`
	MachineCount          types.Int64  `tfsdk:"machine_count"`
}

func (r *unknownThreatDetectStrategyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = unknownThreatDetectStrategyTypeName
}

func (r *unknownThreatDetectStrategyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"asset_selection_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"study_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("hash", "path"),
				},
			},
			"duration_days_after_init": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"duration_days_after_stop": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"machine_count": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *unknownThreatDetectStrategyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan unknownThreatDetectStrategyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, err := unknownThreatDetectStrategyCreateRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: building CreateUnknownThreatDetectStrategy request", unknownThreatDetectStrategyTypeName),
			err.Error(),
		)
		return
	}

	action := "CreateUnknownThreatDetectStrategy"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, unknownThreatDetectStrategyWriteTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateUnknownThreatDetectStrategy", unknownThreatDetectStrategyTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, unknownThreatDetectStrategyTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	rawId, err := jsonpath.Get("$.Id", response)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading $.Id", unknownThreatDetectStrategyTypeName),
			errs.WrapErrorf(err, errs.FailedGetAttributeMsg, unknownThreatDetectStrategyTypeName, "$.Id", response).Error(),
		)
		return
	}
	id, err := UnknownThreatDetectStrategyId(rawId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading $.Id", unknownThreatDetectStrategyTypeName),
			err.Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeUnknownThreatDetectStrategyContext(ctx, id)
	if err != nil {
		// Keep the created ID in state so the strategy stays trackable even when
		// the read-back fails; Terraform marks the resource as tainted.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling ListUnknownThreatDetectStrategy", unknownThreatDetectStrategyTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "ListUnknownThreatDetectStrategy", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, err := unknownThreatDetectStrategyState(object)
	if err != nil {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: mapping ListUnknownThreatDetectStrategy response", unknownThreatDetectStrategyTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// Update applies the updatable fields (name, study_mode, duration_days_after_init,
// duration_days_after_stop). asset_selection_type is replacement-only, so
// changing it triggers Create/Delete through RequiresReplace instead of this method.
func (r *unknownThreatDetectStrategyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan unknownThreatDetectStrategyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueString()
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: missing id", unknownThreatDetectStrategyTypeName),
			"id must be set before updating",
		)
		return
	}

	request := map[string]interface{}{
		"Id":   id,
		"Name": plan.Name.ValueString(),
	}
	if !plan.StudyMode.IsNull() && !plan.StudyMode.IsUnknown() && plan.StudyMode.ValueString() != "" {
		request["StudyMode"] = plan.StudyMode.ValueString()
	}
	if !plan.DurationDaysAfterInit.IsNull() && !plan.DurationDaysAfterInit.IsUnknown() {
		request["DurationDaysAfterInit"] = plan.DurationDaysAfterInit.ValueInt64()
	}
	if !plan.DurationDaysAfterStop.IsNull() && !plan.DurationDaysAfterStop.IsUnknown() {
		request["DurationDaysAfterStop"] = plan.DurationDaysAfterStop.ValueInt64()
	}

	action := "UpdateUnknownThreatDetectStrategy"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.RetryContext(ctx, unknownThreatDetectStrategyWriteTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateUnknownThreatDetectStrategy", unknownThreatDetectStrategyTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeUnknownThreatDetectStrategyContext(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: calling ListUnknownThreatDetectStrategy", unknownThreatDetectStrategyTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "ListUnknownThreatDetectStrategy", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	model, err := unknownThreatDetectStrategyState(object)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: mapping ListUnknownThreatDetectStrategy response", unknownThreatDetectStrategyTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *unknownThreatDetectStrategyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state unknownThreatDetectStrategyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}

	object, err := NewService(r.Meta()).DescribeUnknownThreatDetectStrategyContext(ctx, id)
	if err != nil {
		if errs.IsProviderNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListUnknownThreatDetectStrategy", unknownThreatDetectStrategyTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "ListUnknownThreatDetectStrategy", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, err := unknownThreatDetectStrategyState(object)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: mapping ListUnknownThreatDetectStrategy response", unknownThreatDetectStrategyTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *unknownThreatDetectStrategyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state unknownThreatDetectStrategyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if err := NewService(r.Meta()).DeleteUnknownThreatDetectStrategyContext(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteUnknownThreatDetectStrategy", unknownThreatDetectStrategyTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DeleteUnknownThreatDetectStrategy", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *unknownThreatDetectStrategyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func unknownThreatDetectStrategyCreateRequest(plan unknownThreatDetectStrategyResourceModel) (map[string]interface{}, error) {
	if plan.Name.IsNull() || plan.Name.IsUnknown() || plan.Name.ValueString() == "" {
		return nil, fmt.Errorf("name is required")
	}
	if plan.AssetSelectionType.IsNull() || plan.AssetSelectionType.IsUnknown() || plan.AssetSelectionType.ValueString() == "" {
		return nil, fmt.Errorf("asset_selection_type is required")
	}
	request := map[string]interface{}{
		"Name":               plan.Name.ValueString(),
		"AssetSelectionType": plan.AssetSelectionType.ValueString(),
	}
	if !plan.StudyMode.IsNull() && !plan.StudyMode.IsUnknown() && plan.StudyMode.ValueString() != "" {
		request["StudyMode"] = plan.StudyMode.ValueString()
	}
	if !plan.DurationDaysAfterInit.IsNull() && !plan.DurationDaysAfterInit.IsUnknown() {
		request["DurationDaysAfterInit"] = plan.DurationDaysAfterInit.ValueInt64()
	}
	if !plan.DurationDaysAfterStop.IsNull() && !plan.DurationDaysAfterStop.IsUnknown() {
		request["DurationDaysAfterStop"] = plan.DurationDaysAfterStop.ValueInt64()
	}
	return request, nil
}

// unknownThreatDetectStrategyState maps a ListUnknownThreatDetectStrategy element
// into the resource model. Required attributes (name, asset_selection_type) are
// validated before writing state so a contract-drifted response fails the refresh
// instead of replacing configured values with empty ones.
func unknownThreatDetectStrategyState(object map[string]interface{}) (unknownThreatDetectStrategyResourceModel, error) {
	id, err := UnknownThreatDetectStrategyId(object["Id"])
	if err != nil {
		return unknownThreatDetectStrategyResourceModel{}, err
	}
	name, err := unknownThreatDetectStrategyRequiredString(object["Name"], "Name")
	if err != nil {
		return unknownThreatDetectStrategyResourceModel{}, err
	}
	assetSelectionType, err := unknownThreatDetectStrategyRequiredString(object["AssetSelectionType"], "AssetSelectionType")
	if err != nil {
		return unknownThreatDetectStrategyResourceModel{}, err
	}
	return unknownThreatDetectStrategyResourceModel{
		Id:                    types.StringValue(id),
		Name:                  name,
		AssetSelectionType:    assetSelectionType,
		StudyMode:             unknownThreatDetectStrategyString(object["StudyMode"]),
		DurationDaysAfterInit: unknownThreatDetectStrategyInt(object["DurationDaysAfterInit"]),
		DurationDaysAfterStop: unknownThreatDetectStrategyInt(object["DurationDaysAfterStop"]),
		MachineCount:          unknownThreatDetectStrategyInt(object["MachineCount"]),
	}, nil
}

func unknownThreatDetectStrategyRequiredString(raw interface{}, field string) (types.String, error) {
	v, ok := raw.(string)
	if !ok || v == "" {
		return types.StringNull(), fmt.Errorf("%s is required in the ListUnknownThreatDetectStrategy response, got %T %v", field, raw, raw)
	}
	return types.StringValue(v), nil
}

func unknownThreatDetectStrategyString(raw interface{}) types.String {
	if v, ok := raw.(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

// unknownThreatDetectStrategyInt converts a JSON-decoded integer (float64,
// json.Number, int, int64) into a Framework Int64 value. A missing or non-integer
// value yields null instead of corrupting state with 0.
func unknownThreatDetectStrategyInt(raw interface{}) types.Int64 {
	switch v := raw.(type) {
	case float64:
		return types.Int64Value(int64(v))
	case int:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case int32:
		return types.Int64Value(int64(v))
	case json.Number:
		n, err := v.Int64()
		if err == nil {
			return types.Int64Value(n)
		}
	case string:
		// fall through to null
		_ = v
	case nil:
		return types.Int64Null()
	}
	return types.Int64Null()
}
