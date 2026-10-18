package threatdetection

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const customCheckItemTypeName = "alicloud_threat_detection_custom_check_item"

var (
	_ resource.Resource                = &customCheckItemResource{}
	_ resource.ResourceWithImportState = &customCheckItemResource{}

	// descObjectType is the framework object type of the description / assist_info
	// / solution blocks: a single object with Type and Value strings. The API
	// round-trips each block as {"Type":..,"Value":..}; the block is a list of
	// at most one such object, matching the SDK v2 MaxItems:1 shape.
	descObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":  types.StringType,
			"value": types.StringType,
		},
	}
)

func NewCustomCheckItemResource() resource.Resource {
	return &customCheckItemResource{}
}

type customCheckItemResource struct {
	fwadapt.ResourceBase
}

// checkItemStructModel is the element type of the description / assist_info /
// / solution blocks.
type checkItemStructModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

type customCheckItemModel struct {
	Id              types.String `tfsdk:"id"`
	CheckId         types.Int64  `tfsdk:"check_id"`
	CheckShowName   types.String `tfsdk:"check_show_name"`
	SectionIds      types.List   `tfsdk:"section_ids"`
	Vendor          types.String `tfsdk:"vendor"`
	InstanceType    types.String `tfsdk:"instance_type"`
	InstanceSubType types.String `tfsdk:"instance_sub_type"`
	RiskLevel       types.String `tfsdk:"risk_level"`
	Status          types.String `tfsdk:"status"`
	CheckRule       types.String `tfsdk:"check_rule"`
	Remark          types.String `tfsdk:"remark"`
	Description     types.List   `tfsdk:"description"`
	AssistInfo      types.List   `tfsdk:"assist_info"`
	Solution        types.List   `tfsdk:"solution"`
	CheckType       types.String `tfsdk:"check_type"`
}

func (r *customCheckItemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = customCheckItemTypeName
}

func (r *customCheckItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Threat Detection Custom Check Item resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `check_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"check_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the custom check item.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"check_show_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the check item.",
			},
			"section_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Required:    true,
				Description: "The list of section IDs associated with the custom check item. " +
					"The API does not return this attribute, so the configured value is preserved in state.",
			},
			"vendor": schema.StringAttribute{
				Required:    true,
				Description: "The vendor to which the check item belongs. Changing this creates a new check item.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_type": schema.StringAttribute{
				Required:    true,
				Description: "The asset type to which the check item belongs. Changing this creates a new check item.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_sub_type": schema.StringAttribute{
				Required:    true,
				Description: "The asset subtype to which the check item belongs. Changing this creates a new check item.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"risk_level": schema.StringAttribute{
				Required:    true,
				Description: "The risk level of the check item. Valid values: `HIGH`, `MEDIUM`, `LOW`.",
				Validators: []validator.String{
					stringvalidator.OneOf("HIGH", "MEDIUM", "LOW"),
				},
			},
			"status": schema.StringAttribute{
				Required:    true,
				Description: "The status of the check item. Valid values: `RELEASE`, `EDIT`.",
				Validators: []validator.String{
					stringvalidator.OneOf("RELEASE", "EDIT"),
				},
			},
			"check_rule": schema.StringAttribute{
				Required: true,
				Description: "The check rule of the check item. It must be a JSON definition string that contains " +
					"AssociatedData and MatchProperty; arbitrary strings are rejected by the API.",
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Description: "The remarks of the check item.",
			},
			"check_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the check item.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"description": schema.ListNestedBlock{
				Description: "The description information of the check item.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "The type of the description information.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "The value of the description information.",
						},
					},
				},
			},
			"assist_info": schema.ListNestedBlock{
				Description: "The description help information of the check item.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "The type of the description help information.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "The value of the description help information.",
						},
					},
				},
			},
			"solution": schema.ListNestedBlock{
				Description: "The solution of the check item.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "The type of the solution.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "The value of the solution.",
						},
					},
				},
			},
		},
	}
}

func (r *customCheckItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customCheckItemModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"CheckShowName":   plan.CheckShowName.ValueString(),
		"Vendor":          plan.Vendor.ValueString(),
		"InstanceType":    plan.InstanceType.ValueString(),
		"InstanceSubType": plan.InstanceSubType.ValueString(),
		"RiskLevel":       plan.RiskLevel.ValueString(),
		"Status":          plan.Status.ValueString(),
		"CheckRule":       plan.CheckRule.ValueString(),
		"SectionIds":      sectionIDsToSlice(ctx, plan.SectionIds),
	}
	if !plan.Remark.IsNull() {
		request["Remark"] = plan.Remark.ValueString()
	}
	if m := buildStructMap(ctx, plan.Description); len(m) > 0 {
		if b, err := json.Marshal(m); err == nil {
			request["Description"] = string(b)
		}
	}
	if m := buildStructMap(ctx, plan.AssistInfo); len(m) > 0 {
		if b, err := json.Marshal(m); err == nil {
			request["AssistInfo"] = string(b)
		}
	}
	if m := buildStructMap(ctx, plan.Solution); len(m) > 0 {
		if b, err := json.Marshal(m); err == nil {
			request["Solution"] = string(b)
		}
	}

	// CreateCheckItem is an RPC POST: business parameters travel in the request
	// body while the query map stays empty, the same shape the SDK v2 half and
	// the existing attack_path_whitelist resource use for Sas/2018-12-03.
	response, err := r.Client().RpcPost("Sas", "2018-12-03", "CreateCheckItem", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateCheckItem", customCheckItemTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateCheckItem", customCheckItemTypeName).Error(),
		)
		return
	}

	checkId, ok := extractCheckId(response)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateCheckItem", customCheckItemTypeName),
			fmt.Sprintf("The CreateCheckItem response carried no Data.CheckId: %v", response),
		)
		return
	}

	object, err := listCustomCheckItem(r.Client(), strconv.FormatInt(checkId, 10))
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListCheckItems", customCheckItemTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListCheckItems", customCheckItemTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateCheckItem", customCheckItemTypeName),
			fmt.Sprintf("The custom check item %d is not visible via ListCheckItems right after creation", checkId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, buildCustomCheckItemState(strconv.FormatInt(checkId, 10), object, plan.SectionIds))...)
}

func (r *customCheckItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state customCheckItemModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := listCustomCheckItem(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListCheckItems", customCheckItemTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListCheckItems", customCheckItemTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	// section_ids is not returned by ListCheckItems, so the prior state value is
	// carried forward unchanged.
	resp.Diagnostics.Append(resp.State.Set(ctx, buildCustomCheckItemState(state.Id.ValueString(), object, state.SectionIds))...)
}

func (r *customCheckItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state customCheckItemModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkId, err := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing check_id", customCheckItemTypeName),
			errs.WrapErrorf(err, "Updating %s: parsing check_id %q", customCheckItemTypeName, state.Id.ValueString()).Error(),
		)
		return
	}

	request := map[string]interface{}{
		"CheckId": checkId,
	}
	// Only fields that changed are sent, matching the SDK v2 setCustomCheckItemRequest
	// for-update behaviour. vendor / instance_type / instance_sub_type are
	// RequiresReplace and never reach Update.
	if !plan.CheckShowName.Equal(state.CheckShowName) {
		request["CheckShowName"] = plan.CheckShowName.ValueString()
	}
	if !plan.RiskLevel.Equal(state.RiskLevel) {
		request["RiskLevel"] = plan.RiskLevel.ValueString()
	}
	if !plan.Status.Equal(state.Status) {
		request["Status"] = plan.Status.ValueString()
	}
	if !plan.CheckRule.Equal(state.CheckRule) {
		request["CheckRule"] = plan.CheckRule.ValueString()
	}
	if !plan.SectionIds.Equal(state.SectionIds) {
		request["SectionIds"] = sectionIDsToSlice(ctx, plan.SectionIds)
	}
	if !plan.Remark.Equal(state.Remark) && !plan.Remark.IsNull() {
		request["Remark"] = plan.Remark.ValueString()
	}
	// description / assist_info / solution are sent whenever the planned block
	// carries content; an absent block is omitted so the API clears the field.
	if m := buildStructMap(ctx, plan.Description); len(m) > 0 {
		if b, err := json.Marshal(m); err == nil {
			request["Description"] = string(b)
		}
	}
	if m := buildStructMap(ctx, plan.AssistInfo); len(m) > 0 {
		if b, err := json.Marshal(m); err == nil {
			request["AssistInfo"] = string(b)
		}
	}
	if m := buildStructMap(ctx, plan.Solution); len(m) > 0 {
		if b, err := json.Marshal(m); err == nil {
			request["Solution"] = string(b)
		}
	}

	if _, err := r.Client().RpcPost("Sas", "2018-12-03", "UpdateCheckItem", nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateCheckItem", customCheckItemTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateCheckItem", customCheckItemTypeName).Error(),
		)
		return
	}

	object, err := listCustomCheckItem(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListCheckItems", customCheckItemTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListCheckItems", customCheckItemTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, buildCustomCheckItemState(state.Id.ValueString(), object, plan.SectionIds))...)
}

func (r *customCheckItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customCheckItemModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkId, err := strconv.ParseInt(state.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing check_id", customCheckItemTypeName),
			errs.WrapErrorf(err, "Deleting %s: parsing check_id %q", customCheckItemTypeName, state.Id.ValueString()).Error(),
		)
		return
	}

	_, err = r.Client().RpcPost("Sas", "2018-12-03", "DeleteCheckItem", nil, map[string]interface{}{
		"CheckIds.1": checkId,
	}, true)
	if err != nil && !errs.NotFoundError(err) && !isExpectedErrorCode(err, "CspmDeleteCheckCustomItemError") {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteCheckItem", customCheckItemTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteCheckItem", customCheckItemTypeName).Error(),
		)
	}
}

func (r *customCheckItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// listCustomCheckItem queries ListCheckItems by check ID and returns the
// matching item, or nil when the check item no longer exists. ListCheckItems is
// an RPC GET, so the filter parameters travel in the signed query map.
func listCustomCheckItem(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := client.RpcGet("Sas", "2018-12-03", "ListCheckItems", map[string]interface{}{
		"CheckId":     id,
		"PageSize":    1,
		"CurrentPage": 1,
	}, nil)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	items, _ := response["CheckItems"].([]interface{})
	if len(items) == 0 {
		return nil, nil
	}
	object, _ := items[0].(map[string]interface{})
	return object, nil
}

// extractCheckId reads $.Data.CheckId from a CreateCheckItem response. The SDK
// may surface the JSON number as float64, int64, json.Number or a string-backed
// type depending on the decoder, so every plausible form is handled and a
// format-and-parse fallback covers pointer / decimal wrappers.
func extractCheckId(response map[string]interface{}) (int64, bool) {
	data, _ := response["Data"].(map[string]interface{})
	if data == nil {
		return 0, false
	}
	raw, ok := data["CheckId"]
	if !ok || raw == nil {
		return 0, false
	}
	if id, ok := toInt64Loose(raw); ok {
		return id, true
	}
	return 0, false
}

// toInt64Loose coerces a response scalar to int64 across the decoder-dependent
// representations (float64, int64, int, json.Number, string) and falls back to a
// format-and-parse so pointer / decimal-like wrappers still resolve.
func toInt64Loose(raw interface{}) (int64, bool) {
	switch v := raw.(type) {
	case float64:
		return int64(v), true
	case int64:
		return v, true
	case int:
		return int64(v), true
	case json.Number:
		id, err := strconv.ParseInt(string(v), 10, 64)
		if err == nil {
			return id, true
		}
		if f, err := strconv.ParseFloat(string(v), 64); err == nil {
			return int64(f), true
		}
	case string:
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			return id, true
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(f), true
		}
	}
	s := fmt.Sprintf("%v", raw)
	if id, err := strconv.ParseInt(s, 10, 64); err == nil {
		return id, true
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return int64(f), true
	}
	return 0, false
}

// buildCustomCheckItemState maps a ListCheckItems item onto the resource model.
// sectionIDs is carried from the plan or prior state because the API never
// returns section_ids.
func buildCustomCheckItemState(id string, object map[string]interface{}, sectionIDs types.List) customCheckItemModel {
	return customCheckItemModel{
		Id:              types.StringValue(id),
		CheckId:         int64Value(object["CheckId"]),
		CheckShowName:   stringValue(object["CheckShowName"]),
		Vendor:          stringValue(object["Vendor"]),
		InstanceType:    stringValue(object["InstanceType"]),
		InstanceSubType: stringValue(object["InstanceSubType"]),
		RiskLevel:       stringValue(object["RiskLevel"]),
		Status:          stringValue(object["Status"]),
		CheckRule:       stringValue(object["CheckRule"]),
		Remark:          stringValue(object["Remark"]),
		CheckType:       stringValue(object["CheckType"]),
		Description:     structList(object["Description"]),
		AssistInfo:      structList(object["AssistInfo"]),
		Solution:        structList(object["Solution"]),
		SectionIds:      sectionIDs,
	}
}

// sectionIDsToSlice extracts the int64 elements of a section_ids list as the
// []interface{} shape the SectionIds RPC parameter expects.
func sectionIDsToSlice(ctx context.Context, list types.List) []interface{} {
	out := make([]interface{}, 0, len(list.Elements()))
	for _, e := range list.Elements() {
		if v, ok := e.(types.Int64); ok {
			out = append(out, v.ValueInt64())
		}
	}
	return out
}

// buildStructMap extracts the first element of a description / assist_info /
// / solution block into the {"Type","Value"} map the API expects, omitting
// empty entries. An absent or empty block yields an empty map.
func buildStructMap(ctx context.Context, list types.List) map[string]interface{} {
	m := make(map[string]interface{})
	if list.IsNull() || list.IsUnknown() {
		return m
	}
	var vals []checkItemStructModel
	resp := list.ElementsAs(ctx, &vals, false)
	if resp.HasError() || len(vals) == 0 {
		return m
	}
	item := vals[0]
	if !item.Type.IsNull() && item.Type.ValueString() != "" {
		m["Type"] = item.Type.ValueString()
	}
	if !item.Value.IsNull() && item.Value.ValueString() != "" {
		m["Value"] = item.Value.ValueString()
	}
	return m
}

// structList flattens a {"Type","Value"} response object into the block list
// form. A nil or empty object yields a null list.
func structList(raw interface{}) types.List {
	m, ok := raw.(map[string]interface{})
	if !ok || len(m) == 0 {
		return types.ListNull(descObjectType)
	}
	obj, _ := types.ObjectValue(descObjectType.AttrTypes, map[string]attr.Value{
		"type":  stringValue(m["Type"]),
		"value": stringValue(m["Value"]),
	})
	l, _ := types.ListValue(descObjectType, []attr.Value{obj})
	return l
}

// stringValue converts a response scalar to a String value, treating empty
// strings as null so absent optional fields read back as absent.
func stringValue(v interface{}) types.String {
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

// int64Value converts a response numeric scalar to an Int64 value. It routes
// through toInt64Loose so json.Number and string-backed numerics resolve the
// same way as extractCheckId.
func int64Value(v interface{}) types.Int64 {
	if id, ok := toInt64Loose(v); ok {
		return types.Int64Value(id)
	}
	return types.Int64Null()
}

// isExpectedErrorCode reports whether err is a *tea.SDKError whose code matches
// one of the expected codes, with a substring fallback for wrapped forms. It
// mirrors IsExpectedErrors' behaviour without importing package alicloud.
func isExpectedErrorCode(err error, codes ...string) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		for _, c := range codes {
			if *e.Code == c {
				return true
			}
		}
	}
	msg := err.Error()
	for _, c := range codes {
		if strings.Contains(msg, c) {
			return true
		}
	}
	return false
}
