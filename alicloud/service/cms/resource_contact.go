package cms

import (
	"context"
	"fmt"

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

const contactTypeName = "alicloud_cms_contact"

var (
	_ resource.Resource                = &contactResource{}
	_ resource.ResourceWithImportState = &contactResource{}
)

// NewContactResource is the Framework factory for the alicloud_cms_contact
// managed resource, registered by ServicePackage.
func NewContactResource() resource.Resource {
	return &contactResource{}
}

// contactResource implements alicloud_cms_contact on the
// terraform-plugin-framework. It embeds fwadapt.ResourceBase so the muxed
// provider injects the *connectivity.AliyunClient during Configure; the CRUD
// logic delegates to CmsService so the API path is shared with the data
// source and exercised by the no-cloud unit tests.
type contactResource struct {
	fwadapt.ResourceBase
}

// contactResourceModel is the tfsdk state shape of alicloud_cms_contact.
type contactResourceModel struct {
	ContactName types.String `tfsdk:"contact_name"`
	Email       types.String `tfsdk:"email"`
	Phone       types.String `tfsdk:"phone"`
	Lang        types.String `tfsdk:"lang"`
	ImUserIds   types.Map    `tfsdk:"im_user_ids"`
	Workspace   types.String `tfsdk:"workspace"`
	ContactId   types.String `tfsdk:"contact_id"`
	RegionId    types.String `tfsdk:"region_id"`
	Id          types.String `tfsdk:"id"`
}

func (r *contactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = contactTypeName
}

func (r *contactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Cloud Monitor Service (CMS) Contact resource.",
		Attributes: map[string]schema.Attribute{
			"contact_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the alarm contact.",
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Description: "The email address of the alarm contact.",
			},
			"phone": schema.StringAttribute{
				Optional:    true,
				Description: "The phone number of the alarm contact.",
			},
			"lang": schema.StringAttribute{
				Optional:    true,
				Description: "The language type of the alarm. Valid values: `zh_CN`, `en_US`.",
				Validators: []validator.String{
					stringvalidator.OneOf("zh_CN", "en_US"),
				},
			},
			"im_user_ids": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The DingTalk and other communication tools user IDs. It is a map of string.",
			},
			"workspace": schema.StringAttribute{
				Computed:    true,
				Description: "The workspace ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"contact_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the alarm contact. Same as `id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `contact_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *contactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan contactResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc := NewService(r.Client())
	createResp, err := svc.CreateContact(ctx, contactRequestBody(ctx, plan))
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s", contactTypeName), err.Error())
		return
	}
	id := fmt.Sprint(createResp["contactId"])

	contact, err := svc.GetContact(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s after create", contactTypeName), err.Error())
		return
	}
	if contact == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s", contactTypeName),
			fmt.Sprintf("The CMS Contact %s is not visible via ListContacts right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, contactState(ctx, contact, r.Client().RegionId, plan.ImUserIds))...)
}

func (r *contactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state contactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc := NewService(r.Client())
	contact, err := svc.GetContact(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s", contactTypeName), err.Error())
		return
	}
	if contact == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, contactState(ctx, contact, r.Client().RegionId, state.ImUserIds))...)
}

func (r *contactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state contactResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc := NewService(r.Client())
	if _, err := svc.UpdateContact(ctx, state.Id.ValueString(), contactRequestBody(ctx, plan)); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Updating %s", contactTypeName), err.Error())
		return
	}

	contact, err := svc.GetContact(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s after update", contactTypeName), err.Error())
		return
	}
	if contact == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, contactState(ctx, contact, r.Client().RegionId, plan.ImUserIds))...)
}

func (r *contactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state contactResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc := NewService(r.Client())
	if err := svc.DeleteContact(ctx, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Deleting %s", contactTypeName), err.Error())
		return
	}
}

func (r *contactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// contactRequestBody maps the resource model onto the CMS CreateContact /
// UpdateContact (PutContact) request body. The CMS API persists "" for email
// and phone when the caller clears them, so a null plan value is sent as an
// empty string to match the legacy SDKv2 behavior; im_user_ids is sent as a
// plain map (absent keys become an empty map).
func contactRequestBody(ctx context.Context, m contactResourceModel) map[string]interface{} {
	body := map[string]interface{}{
		"name":      m.ContactName.ValueString(),
		"email":     m.Email.ValueString(),
		"phone":     m.Phone.ValueString(),
		"lang":      m.Lang.ValueString(),
		"imUserIds": imUserIdsToMap(ctx, m.ImUserIds),
	}
	return body
}

// imUserIdsToMap converts a types.Map to the plain map[string]interface{}
// the CMS API expects. An unknown/empty map yields an empty map, matching the
// SDKv2 behavior of sending d.Get("im_user_ids") for an unset attribute.
func imUserIdsToMap(ctx context.Context, m types.Map) map[string]interface{} {
	out := make(map[string]interface{})
	if m.IsNull() || m.IsUnknown() {
		return out
	}
	_ = m.ElementsAs(ctx, &out, false)
	return out
}

// contactState maps a CMS contact object onto the resource model. Strings are
// preserved including "" (present-empty) so a cleared email/phone reads back
// as "" rather than null, matching the legacy SDKv2 state semantics and the
// ACC assertions; absent keys become null.
//
// im_user_ids is refreshed from the API only when the CMS ListContacts response
// actually echoes a populated imUserIds map. The CMS ListContacts endpoint does
// not reliably return imUserIds (the legacy SDKv2 Read only d.Set it when the
// key was present and non-nil, otherwise preserving the prior value). To match
// that behavior under the Framework — where unconditionally setting null would
// produce a "Provider produced inconsistent result after apply" error against a
// known plan value — the caller passes the prior (plan/state) im_user_ids as
// fallbackImUserIds; it is kept verbatim unless the API overrides it.
func contactState(ctx context.Context, obj map[string]interface{}, regionId string, fallbackImUserIds types.Map) contactResourceModel {
	str := func(key string) types.String {
		if v, ok := obj[key].(string); ok {
			return types.StringValue(v)
		}
		return types.StringNull()
	}

	imUserIds := fallbackImUserIds
	if raw, ok := obj["imUserIds"].(map[string]interface{}); ok && len(raw) > 0 {
		m := make(map[string]string, len(raw))
		for k, v := range raw {
			m[k] = fmt.Sprint(v)
		}
		if mv, diags := types.MapValueFrom(ctx, types.StringType, m); !diags.HasError() {
			imUserIds = mv
		}
	}

	contactId := str("contactId")
	id := contactId
	if id.IsNull() {
		id = types.StringNull()
	}
	if regionId == "" {
		regionId = "cn-hangzhou"
	}
	return contactResourceModel{
		ContactName: str("name"),
		Email:       str("email"),
		Phone:       str("phone"),
		Lang:        str("lang"),
		ImUserIds:   imUserIds,
		Workspace:   str("workspace"),
		ContactId:   contactId,
		RegionId:    types.StringValue(regionId),
		Id:          id,
	}
}
