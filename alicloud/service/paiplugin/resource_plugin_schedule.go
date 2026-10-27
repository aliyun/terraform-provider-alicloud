package paiplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const pluginScheduleTypeName = "alicloud_pai_plugin_schedule"

// PAIPlugin returns schedule errors in the response body as ErrorCode/ErrorMessage
// (int32), which connectivity.formatError does not inspect — it only looks at
// Code/code/Success/Message. So every call must check response["ErrorCode"]
// manually. The not-found signal for GetSchedule is ErrorCode == 400, per the
// CloudSpec resourceNotExistCondition (checkErrorCode, "400").
const (
	pluginScheduleProduct  = "PaiPlugin"
	pluginScheduleVersion  = "2022-01-12"
	pluginScheduleBasePath = "/api/v2/schedules"
	pluginScheduleNotFound = 400 // ErrorCode value meaning the schedule does not exist
)

var (
	_ resource.Resource                = &pluginScheduleResource{}
	_ resource.ResourceWithImportState = &pluginScheduleResource{}
)

func NewPluginScheduleResource() resource.Resource {
	return &pluginScheduleResource{}
}

type pluginScheduleResource struct {
	fwadapt.ResourceBase
}

type pluginScheduleModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	GroupId         types.String `tfsdk:"group_id"`
	SignatureId     types.String `tfsdk:"signature_id"`
	TemplateId      types.String `tfsdk:"template_id"`
	TemplateCode    types.String `tfsdk:"template_code"`
	SignName        types.String `tfsdk:"sign_name"`
	PaymentType     types.String `tfsdk:"payment_type"`
	ExecuteTime     types.String `tfsdk:"execute_time"`
	EndTime         types.String `tfsdk:"end_time"`
	AISendStartDate types.String `tfsdk:"ai_send_start_date"`
	AISendEndDate   types.String `tfsdk:"ai_send_end_date"`
	RepeatCycle     types.Int64  `tfsdk:"repeat_cycle"`
	RepeatCycleUnit types.Int64  `tfsdk:"repeat_cycle_unit"`
	RepeatTimes     types.Int64  `tfsdk:"repeat_times"`
	Status          types.Int64  `tfsdk:"status"`
	CreateTime      types.String `tfsdk:"create_time"`
	UpdatedTime     types.String `tfsdk:"updated_time"`
}

func (r *pluginScheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = pluginScheduleTypeName
}

func (r *pluginScheduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Plugin Schedule resource. " +
			"For information about PAI Plugin and how to use it, refer to the product documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the PAI Plugin Schedule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the schedule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "The crowd (audience) ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"signature_id": schema.StringAttribute{
				Optional:    true,
				Description: "The signature ID, or specify a signature.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"template_id": schema.StringAttribute{
				Optional:    true,
				Description: "The template ID, or specify a template code.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"template_code": schema.StringAttribute{
				Optional:    true,
				Description: "The template code.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sign_name": schema.StringAttribute{
				Optional:    true,
				Description: "The signature name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// PaymentType / ExecuteTime / EndTime are echoed by GetSchedule, and the
			// API supplies a default/computed value when the user omits them (e.g.
			// PaymentType defaults to PayAsYouGo). Marking them Computed +
			// UseStateForUnknown keeps an omitted field unknown during planning so a
			// server-supplied default does not diverge from a null plan and trigger a
			// ForceNew recreate loop. A user-supplied value still wins in the plan and
			// is matched against the echoed value on read.
			"payment_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The billing type. Valid values: `PayAsYouGo` (postpaid) and `Subscription` (prepaid). " +
					"Defaults to postpaid when empty.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"execute_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The execution time (UTC+8). Empty means execute immediately.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"end_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The end time (UTC+8).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// AISendStartDate / AISendEndDate are write-only: they are accepted by
			// CreateSchedule but are not returned by GetSchedule, so the read path
			// preserves the configured value from state instead of overwriting it.
			"ai_send_start_date": schema.StringAttribute{
				Optional:    true,
				Description: "The AI send start date, e.g. `2022-01-01 12:00:00`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ai_send_end_date": schema.StringAttribute{
				Optional:    true,
				Description: "The AI send end date, e.g. `2022-02-01 12:00:00`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"repeat_cycle": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The repeat cycle, executed together with the repeat cycle unit.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"repeat_cycle_unit": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The repeat cycle unit. Valid values: `0` (never, default), `1` (hour), `2` (day), `3` (week), `4` (month).",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"repeat_times": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The repeat times. `-1` means repeat with no end time; a positive `N` means terminate after N repeats.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.Int64Attribute{
				Computed:    true,
				Description: "The status of the schedule. `0`: checking, `1`: check success, `2`: check failed, `3`: sending, `4`: send success, `5`: send failed.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time (UTC+8).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "The update time (UTC+8).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *pluginScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pluginScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"Name":    plan.Name.ValueString(),
		"GroupId": plan.GroupId.ValueString(),
	}
	if !plan.SignatureId.IsNull() {
		body["SignatureId"] = plan.SignatureId.ValueString()
	}
	if !plan.TemplateId.IsNull() {
		body["TemplateId"] = plan.TemplateId.ValueString()
	}
	if !plan.TemplateCode.IsNull() {
		body["TemplateCode"] = plan.TemplateCode.ValueString()
	}
	if !plan.SignName.IsNull() {
		body["SignName"] = plan.SignName.ValueString()
	}
	if !plan.PaymentType.IsNull() {
		body["PaymentType"] = plan.PaymentType.ValueString()
	}
	if !plan.ExecuteTime.IsNull() {
		body["ExecuteTime"] = plan.ExecuteTime.ValueString()
	}
	if !plan.EndTime.IsNull() {
		body["EndTime"] = plan.EndTime.ValueString()
	}
	if !plan.AISendStartDate.IsNull() {
		body["AISendStartDate"] = plan.AISendStartDate.ValueString()
	}
	if !plan.AISendEndDate.IsNull() {
		body["AISendEndDate"] = plan.AISendEndDate.ValueString()
	}
	if !plan.RepeatCycle.IsNull() {
		body["RepeatCycle"] = plan.RepeatCycle.ValueInt64()
	}
	if !plan.RepeatCycleUnit.IsNull() {
		body["RepeatCycleUnit"] = plan.RepeatCycleUnit.ValueInt64()
	}
	if !plan.RepeatTimes.IsNull() {
		body["RepeatTimes"] = plan.RepeatTimes.ValueInt64()
	}

	response, err := r.Client().RoaPost(pluginScheduleProduct, pluginScheduleVersion, pluginScheduleBasePath, nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSchedule", pluginScheduleTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateSchedule", pluginScheduleTypeName).Error(),
		)
		return
	}
	if code := errorCode(response); code != 0 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSchedule", pluginScheduleTypeName),
			fmt.Sprintf("CreateSchedule failed: ErrorCode=%d, ErrorMessage=%s", code, errorMessage(response)),
		)
		return
	}

	data, _ := response["Data"].(map[string]interface{})
	scheduleId := stringField(data, "Id")
	if scheduleId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSchedule", pluginScheduleTypeName),
			fmt.Sprintf("The CreateSchedule response carried no Data.Id: %v", response),
		)
		return
	}

	object, err := getSchedule(r.Client(), scheduleId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSchedule", pluginScheduleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetSchedule", pluginScheduleTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSchedule", pluginScheduleTypeName),
			fmt.Sprintf("The schedule %s is not visible via GetSchedule right after creation", scheduleId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, modelFromObject(object, plan))...)
}

func (r *pluginScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pluginScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getSchedule(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetSchedule", pluginScheduleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetSchedule", pluginScheduleTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, modelFromObject(object, state))...)
}

func (r *pluginScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Every settable attribute is ForceNew (RequiresReplace); PAIPlugin exposes
	// no UpdateSchedule API. Framework routes changes to Delete+Create, so this
	// method is unreachable in practice but must exist to satisfy the interface.
	var plan pluginScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *pluginScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pluginScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := fmt.Sprintf("%s/%s", pluginScheduleBasePath, state.Id.ValueString())
	response, err := r.Client().RoaDelete(pluginScheduleProduct, pluginScheduleVersion, action, nil, nil, nil, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteSchedule", pluginScheduleTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteSchedule", pluginScheduleTypeName).Error(),
		)
		return
	}
	// ErrorCode 400 means the schedule is already gone; tolerate it as a
	// successful delete, matching the GetSchedule not-found condition.
	if code := errorCode(response); code != 0 && code != pluginScheduleNotFound {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteSchedule", pluginScheduleTypeName),
			fmt.Sprintf("DeleteSchedule failed: ErrorCode=%d, ErrorMessage=%s", code, errorMessage(response)),
		)
		return
	}
}

func (r *pluginScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getSchedule queries GetSchedule by schedule ID and returns the Data object,
// or nil when the schedule no longer exists (ErrorCode 400).
func getSchedule(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	action := fmt.Sprintf("%s/%s", pluginScheduleBasePath, id)
	response, err := client.RoaGet(pluginScheduleProduct, pluginScheduleVersion, action, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if code := errorCode(response); code == pluginScheduleNotFound {
		return nil, nil
	} else if code != 0 {
		return nil, fmt.Errorf("GetSchedule failed: ErrorCode=%d, ErrorMessage=%s", code, errorMessage(response))
	}
	data, ok := response["Data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("the GetSchedule response carried no Data object: %v", response)
	}
	return data, nil
}

// modelFromObject maps a GetSchedule Data object onto the resource model. The
// base (plan or prior state) supplies the write-only AISend* fields, which the
// API never echoes back. Empty strings become null, matching the ims pattern.
func modelFromObject(object map[string]interface{}, base pluginScheduleModel) pluginScheduleModel {
	return pluginScheduleModel{
		Id:              stringAttr(object, "Id"),
		Name:            stringAttr(object, "Name"),
		GroupId:         stringAttr(object, "GroupId"),
		SignatureId:     stringAttr(object, "SignatureId"),
		TemplateId:      stringAttr(object, "TemplateId"),
		TemplateCode:    stringAttr(object, "TemplateCode"),
		SignName:        stringAttr(object, "SignName"),
		PaymentType:     stringAttr(object, "PaymentType"),
		ExecuteTime:     stringAttr(object, "ExecuteTime"),
		EndTime:         stringAttr(object, "EndTime"),
		AISendStartDate: base.AISendStartDate, // write-only: preserved from base
		AISendEndDate:   base.AISendEndDate,   // write-only: preserved from base
		RepeatCycle:     intAttr(object, "RepeatCycle"),
		RepeatCycleUnit: intAttr(object, "RepeatCycleUnit"),
		RepeatTimes:     intAttr(object, "RepeatTimes"),
		Status:          intAttr(object, "Status"),
		CreateTime:      stringAttr(object, "CreatedTime"),
		UpdatedTime:     stringAttr(object, "UpdatedTime"),
	}
}

// stringField returns object[key] as a plain string, or "" when absent.
func stringField(object map[string]interface{}, key string) string {
	if object == nil {
		return ""
	}
	if v, ok := object[key].(string); ok {
		return v
	}
	return ""
}

// stringAttr returns object[key] as a types.String, mapping "" or absence to
// null so optional fields the API omits read back as absent rather than "".
func stringAttr(object map[string]interface{}, key string) types.String {
	if v := stringField(object, key); v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

// intAttr returns object[key] as a types.Int64, accepting json.Number, float64
// and integer forms. Absence or null maps to a null value.
func intAttr(object map[string]interface{}, key string) types.Int64 {
	if object == nil {
		return types.Int64Null()
	}
	v, ok := object[key]
	if !ok || v == nil {
		return types.Int64Null()
	}
	switch n := v.(type) {
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return types.Int64Null()
		}
		return types.Int64Value(i)
	case float64:
		return types.Int64Value(int64(n))
	case int:
		return types.Int64Value(int64(n))
	case int64:
		return types.Int64Value(n)
	case string:
		i, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return types.Int64Null()
		}
		return types.Int64Value(i)
	}
	return types.Int64Null()
}

// errorCode returns response["ErrorCode"] as an int. PAIPlugin reports business
// errors (including not-found) here; 0 means success.
func errorCode(response map[string]interface{}) int {
	if response == nil {
		return 0
	}
	v, ok := response["ErrorCode"]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case string:
		i, _ := strconv.Atoi(n)
		return i
	}
	return 0
}

// errorMessage returns response["ErrorMessage"] as a string, or "" when absent.
func errorMessage(response map[string]interface{}) string {
	if v, ok := response["ErrorMessage"].(string); ok {
		return v
	}
	return ""
}
