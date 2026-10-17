package paiplugin

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const messageTypeName = "alicloud_pai_plugin_message"

var (
	_ resource.Resource                = &messageResource{}
	_ resource.ResourceWithImportState = &messageResource{}
)

func NewMessageResource() resource.Resource {
	return &messageResource{}
}

type messageResource struct {
	fwadapt.ResourceBase
}

type messageModel struct {
	Id               types.String `tfsdk:"id"`
	MessageId        types.String `tfsdk:"message_id"`
	PhoneNumber      types.String `tfsdk:"phone_number"`
	Signature        types.String `tfsdk:"signature"`
	SignatureId      types.String `tfsdk:"signature_id"`
	TemplateCode     types.String `tfsdk:"template_code"`
	TemplateId       types.String `tfsdk:"template_id"`
	TemplateParams   types.List   `tfsdk:"template_params"`
	PaymentType      types.String `tfsdk:"payment_type"`
	ScheduleIdCreate types.String `tfsdk:"schedule_id_create"`
	SmsUpExtendCodes types.List   `tfsdk:"sms_up_extend_codes"`
	OutIds           types.List   `tfsdk:"out_ids"`
	GroupId          types.String `tfsdk:"group_id"`
	OutId            types.String `tfsdk:"out_id"`
	ScheduleId       types.String `tfsdk:"schedule_id"`
	ErrorCode        types.String `tfsdk:"error_code"`
	Status           types.Int64  `tfsdk:"status"`
	TemplateType     types.Int64  `tfsdk:"template_type"`
	SendDate         types.String `tfsdk:"send_date"`
}

func (r *messageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = messageTypeName
}

func (r *messageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Plugin Message resource. SendMessage creates the message; ListMessages reads its status.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `message_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"message_id": schema.StringAttribute{
				Computed:    true,
				Description: "The message id, returned by SendMessage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"phone_number": schema.StringAttribute{
				Required:    true,
				Description: "The phone number. Sent as the single-element PhoneNumbers array.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"signature": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The signature name. Either signature or signature_id may be set. The server echoes it back on Read, so it is also Computed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"signature_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The signature id. Either signature or signature_id may be set. The server resolves and echoes it on Read, so it is also Computed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"template_code": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The template code. Either template_code or template_id may be set. The server echoes it back on Read, so it is also Computed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"template_id": schema.StringAttribute{
				Optional:    true,
				Description: "The template id. Either template_code or template_id may be set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"template_params": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The template params, JSON strings keyed by variable name. The server echoes them back on Read, so it is also Computed.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"payment_type": schema.StringAttribute{
				Optional:    true,
				Description: "The payment type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schedule_id_create": schema.StringAttribute{
				Optional:    true,
				Description: "The schedule id used to associate a reach plan at send time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sms_up_extend_codes": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The sms up extend codes.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"out_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The out ids.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.StringAttribute{
				Computed:    true,
				Description: "The associated group id, empty when unassociated.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"out_id": schema.StringAttribute{
				Computed:    true,
				Description: "The out id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"schedule_id": schema.StringAttribute{
				Computed:    true,
				Description: "The associated schedule id, empty when unassociated.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"error_code": schema.StringAttribute{
				Computed:    true,
				Description: "The error code of the message.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.Int64Attribute{
				Computed:    true,
				Description: "The send status. 0: sending, 1: success, 2: failure.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"template_type": schema.Int64Attribute{
				Computed:    true,
				Description: "The template type. 0: verification code, 1: notification, 2: marketing.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"send_date": schema.StringAttribute{
				Computed:    true,
				Description: "The send date (yyyymmdd) captured at creation, used to read the message status via ListMessages which requires a Datetime filter.",
			},
		},
	}
}

func (r *messageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan messageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"PhoneNumbers": []string{plan.PhoneNumber.ValueString()},
	}
	if !plan.Signature.IsNull() {
		body["SignName"] = plan.Signature.ValueString()
	}
	if !plan.SignatureId.IsNull() {
		body["SignatureId"] = plan.SignatureId.ValueString()
	}
	if !plan.TemplateCode.IsNull() {
		body["TemplateCode"] = plan.TemplateCode.ValueString()
	}
	if !plan.TemplateId.IsNull() {
		body["TemplateId"] = plan.TemplateId.ValueString()
	}
	if params := strSliceFromList(plan.TemplateParams); len(params) > 0 {
		body["TemplateParams"] = params
	}
	if !plan.PaymentType.IsNull() {
		body["PaymentType"] = plan.PaymentType.ValueString()
	}
	if !plan.ScheduleIdCreate.IsNull() {
		body["ScheduleId"] = plan.ScheduleIdCreate.ValueString()
	}
	if codes := strSliceFromList(plan.SmsUpExtendCodes); len(codes) > 0 {
		body["SmsUpExtendCodes"] = codes
	}
	if outIds := strSliceFromList(plan.OutIds); len(outIds) > 0 {
		body["OutIds"] = outIds
	}

	response, err := r.Client().RoaPost(paiPluginProduct, paiPluginVersion, pathMessages, nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling SendMessage", messageTypeName),
			errs.WrapErrorf(err, "Creating %s: calling SendMessage", messageTypeName).Error(),
		)
		return
	}

	data := dataObject(response)
	messageId := firstMessageId(data)
	if messageId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling SendMessage", messageTypeName),
			fmt.Sprintf("The response carried no Data.Messages[].Id: %v", response),
		)
		return
	}

	state := messageStateFromModel(&plan)
	state.Id = types.StringValue(messageId)
	state.MessageId = types.StringValue(messageId)
	state.SendDate = types.StringValue(time.Now().Format("20060102"))

	// Refresh from ListMessages so the apply result reflects the server-side
	// status fields rather than the planned inputs alone.
	if object := findMessage(r.Client(), state.SendDate.ValueString(), messageId); object != nil {
		state = mergeMessageState(state, object)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *messageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state messageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Without a send date the message cannot be re-fetched: ListMessages
	// requires a Datetime filter and SendMessage returns none. Preserve state
	// rather than guessing, so imports without send_date stay stable.
	if state.SendDate.ValueString() == "" {
		return
	}

	object := findMessage(r.Client(), state.SendDate.ValueString(), state.MessageId.ValueString())
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, mergeMessageState(state, object))...)
}

func (r *messageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// SendMessage is fire-and-forget: every input change RequiresReplace, so
	// an in-place update never reaches the provider.
	var state messageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *messageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// A sent message cannot be unsent: PAIPlugin exposes no delete, so delete
	// only removes the resource from state.
}

func (r *messageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// firstMessageId reads Data.Messages[0].Id from a SendMessage response.
func firstMessageId(data map[string]interface{}) string {
	if data == nil {
		return ""
	}
	messages, ok := data["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		return ""
	}
	first, ok := messages[0].(map[string]interface{})
	if !ok {
		return ""
	}
	return str(first, "Id")
}

// findMessage queries ListMessages by Datetime + MessageId and returns the
// matching message object, or nil when it is absent from that day's batch.
func findMessage(client *connectivity.AliyunClient, datetime, messageId string) map[string]interface{} {
	query := queryFrom("Datetime", datetime, "MessageId", messageId)
	response, err := roaGet(client, pathMessages, query)
	if err != nil {
		return nil
	}
	data := dataObject(response)
	if data == nil {
		return nil
	}
	messages, ok := data["Messages"].([]interface{})
	if !ok {
		return nil
	}
	for _, m := range messages {
		if obj, ok := m.(map[string]interface{}); ok && str(obj, "Id") == messageId {
			return obj
		}
	}
	return nil
}

// mergeMessageState overlays the ListMessages object onto the planned state,
// keeping the planned inputs that ListMessages does not echo back (signature,
// template params, etc.) while refreshing the server-side status fields.
func mergeMessageState(state messageModel, obj map[string]interface{}) messageModel {
	if v := str(obj, "PhoneNumber"); v != "" {
		state.PhoneNumber = types.StringValue(v)
	}
	state.GroupId = strValue(obj, "GroupId")
	state.OutId = strValue(obj, "OutId")
	state.ScheduleId = strValue(obj, "ScheduleId")
	state.ErrorCode = strValue(obj, "ErrorCode")
	state.Signature = strValue(obj, "Signature")
	state.SignatureId = strValue(obj, "SignatureId")
	state.TemplateCode = strValue(obj, "TemplateCode")
	state.TemplateParams = strSliceValue(obj, "TemplateParams")
	if s := intValue(obj, "Status"); s >= 0 {
		state.Status = types.Int64Value(s)
	}
	if t := intValue(obj, "TemplateType"); t >= 0 {
		state.TemplateType = types.Int64Value(t)
	}
	return state
}

func messageStateFromModel(m *messageModel) messageModel {
	return *m
}
