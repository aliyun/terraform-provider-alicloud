package hologram

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const holoWebLoginSettingTypeName = "alicloud_hologram_holo_web_login_setting"

var (
	_ resource.Resource                = &holoWebLoginSettingResource{}
	_ resource.ResourceWithImportState = &holoWebLoginSettingResource{}
)

func NewHoloWebLoginSettingResource() resource.Resource {
	return &holoWebLoginSettingResource{}
}

// holoWebLoginSettingResource manages the HoloWeb login setting of a Hologres
// instance. The underlying API only exposes Get and Update, so the resource
// uses a setting-toggle lifecycle: Create and Update call UpdateHoloWebLogin
// Setting with the configured boolean, Read calls GetHoloWebLoginSetting, and
// Delete reverts the setting to false. The Hologres instance is the resource
// identity and is never created or destroyed by this resource.
type holoWebLoginSettingResource struct {
	fwadapt.ResourceBase
}

type holoWebLoginSettingResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	InstanceId                 types.String `tfsdk:"instance_id"`
	AllowExternalAccountsLogin types.Bool   `tfsdk:"allow_external_accounts_login"`
}

func (r *holoWebLoginSettingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = holoWebLoginSettingTypeName
}

func (r *holoWebLoginSettingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"allow_external_accounts_login": schema.BoolAttribute{
				Required: true,
			},
		},
	}
}

func (r *holoWebLoginSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan holoWebLoginSettingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := plan.InstanceId.ValueString()
	service := NewService(r.Meta())
	if err := service.UpdateHoloWebLoginSetting(instanceId, plan.AllowExternalAccountsLogin.ValueBool()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling UpdateHoloWebLoginSetting", holoWebLoginSettingTypeName),
			err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(instanceId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *holoWebLoginSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state holoWebLoginSettingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// id is the instance id; instance_id may be null right after import.
	instanceId := state.Id.ValueString()
	if instanceId == "" {
		instanceId = state.InstanceId.ValueString()
	}

	service := NewService(r.Meta())
	object, err := service.DescribeHoloWebLoginSetting(instanceId)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetHoloWebLoginSetting", holoWebLoginSettingTypeName),
			err.Error(),
		)
		return
	}

	state.Id = types.StringValue(instanceId)
	state.InstanceId = types.StringValue(instanceId)
	state.AllowExternalAccountsLogin = boolValue(object["AllowExternalAccountsLogin"])
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *holoWebLoginSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state holoWebLoginSettingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	if instanceId == "" {
		instanceId = plan.InstanceId.ValueString()
	}

	service := NewService(r.Meta())
	if !plan.AllowExternalAccountsLogin.Equal(state.AllowExternalAccountsLogin) {
		if err := service.UpdateHoloWebLoginSetting(instanceId, plan.AllowExternalAccountsLogin.ValueBool()); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateHoloWebLoginSetting", holoWebLoginSettingTypeName),
				err.Error(),
			)
			return
		}
	}

	plan.Id = types.StringValue(instanceId)
	plan.InstanceId = types.StringValue(instanceId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *holoWebLoginSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state holoWebLoginSettingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	service := NewService(r.Meta())
	// Revert the setting to the service default (false). The instance is not
	// owned by this resource, so a NotFound on read does not apply here; if the
	// instance is already gone, the toggle is implicitly gone too.
	if err := service.UpdateHoloWebLoginSetting(instanceId, false); err != nil {
		if errs.IsExpectedErrors(err, []string{"resource not exists failed"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling UpdateHoloWebLoginSetting", holoWebLoginSettingTypeName),
			err.Error(),
		)
		return
	}
}

func (r *holoWebLoginSettingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// boolValue maps the JSON-decoded response field onto a Framework Bool. The
// tea ROA client decodes JSON booleans as bool, but older paths may surface
// them as strings.
func boolValue(v interface{}) types.Bool {
	switch b := v.(type) {
	case bool:
		return types.BoolValue(b)
	case string:
		if b == "true" {
			return types.BoolValue(true)
		}
		if b == "false" {
			return types.BoolValue(false)
		}
	}
	return types.BoolNull()
}
