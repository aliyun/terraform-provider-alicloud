package ims

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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

const userTypeName = "alicloud_ims_user"

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	fwadapt.ResourceBase
}

type userResourceModel struct {
	Id                types.String `tfsdk:"id"`
	UserPrincipalName types.String `tfsdk:"user_principal_name"`
	DisplayName       types.String `tfsdk:"display_name"`
	Email             types.String `tfsdk:"email"`
	MobilePhone       types.String `tfsdk:"mobile_phone"`
	Comments          types.String `tfsdk:"comments"`
	UserId            types.String `tfsdk:"user_id"`
	UserName          types.String `tfsdk:"user_name"`
	ProvisionType     types.String `tfsdk:"provision_type"`
	CreateDate        types.String `tfsdk:"create_date"`
	UpdateDate        types.String `tfsdk:"update_date"`
	LastLoginDate     types.String `tfsdk:"last_login_date"`
}

func (r *userResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = userTypeName
}

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a IMS User resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `user_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_principal_name": schema.StringAttribute{
				Required: true,
				Description: "The logon name of the RAM user, in the format `<username>@<AccountAlias>.onaliyun.com`. " +
					"The total length is 1 to 128 characters, and `<username>` is 1 to 64 characters in length. " +
					"Only letters, digits, periods (.), hyphens (-) and underscores (_) are allowed. " +
					"The default domain can be read from the `alicloud_ims_default_domain` data source.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9._@-]+$`),
						"must only contain letters, digits, periods (.), hyphens (-), underscores (_) and the at sign (@)"),
				},
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the RAM user. The length is 1 to 24 characters.",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 24),
				},
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Description: "The email address of the RAM user. This parameter is available only on the China site.",
			},
			"mobile_phone": schema.StringAttribute{
				Optional: true,
				Description: "The mobile phone number of the RAM user, in the format `<international area code>-<number>`, " +
					"for example `86-1868888****`. This parameter is available only on the China site.",
			},
			"comments": schema.StringAttribute{
				Optional:    true,
				Description: "The remarks of the RAM user. The length is 1 to 128 characters.",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
				},
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the RAM user, for example `20732900249392****`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the RAM user, the part of the logon name before the at sign (@).",
				// No UseStateForUnknown: the value follows user_principal_name,
				// so it changes on a rename. Keeping the state value in the plan
				// would trip the inconsistent-result-after-apply check.
			},
			"provision_type": schema.StringAttribute{
				Computed:    true,
				Description: "How the RAM user was created. Valid values: `Manual`, `SCIM` and `CloudSSO`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_date": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the RAM user was created, in RFC 3339 format (UTC).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_date": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the RAM user was last updated, in RFC 3339 format (UTC).",
				// No UseStateForUnknown: UpdateUser bumps the timestamp, so the
				// planned state value can never survive an apply.
			},
			"last_login_date": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the RAM user last logged on to the console, in RFC 3339 format (UTC). Empty if the user has never logged on.",
				// No UseStateForUnknown: the value moves outside provider control
				// whenever the user logs on.
			},
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"UserPrincipalName": plan.UserPrincipalName.ValueString(),
		"DisplayName":       plan.DisplayName.ValueString(),
	}
	if !plan.Email.IsNull() {
		request["Email"] = plan.Email.ValueString()
	}
	if !plan.MobilePhone.IsNull() {
		request["MobilePhone"] = plan.MobilePhone.ValueString()
	}
	if !plan.Comments.IsNull() {
		request["Comments"] = plan.Comments.ValueString()
	}

	response, err := r.Client().RpcPost("Ims", "2019-08-15", "CreateUser", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateUser", userTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateUser", userTypeName).Error(),
		)
		return
	}

	user, _ := response["User"].(map[string]interface{})
	userId, _ := user["UserId"].(string)
	if userId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateUser", userTypeName),
			fmt.Sprintf("The response carried no User.UserId: %v", response),
		)
		return
	}

	object, err := getUser(r.Client(), userId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetUser", userTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetUser", userTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateUser", userTypeName),
			fmt.Sprintf("The user %s is not visible via GetUser right after creation", userId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, userState(object))...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getUser(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetUser", userTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetUser", userTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, userState(object))...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"UserId": state.Id.ValueString(),
	}
	if !plan.UserPrincipalName.Equal(state.UserPrincipalName) {
		request["NewUserPrincipalName"] = plan.UserPrincipalName.ValueString()
	}
	if !plan.DisplayName.Equal(state.DisplayName) {
		request["NewDisplayName"] = plan.DisplayName.ValueString()
	}
	// Optional attributes: a null plan value clears the attribute, matching
	// the legacy alicloud_ram_user update behavior of sending an empty string.
	if !plan.Email.Equal(state.Email) {
		request["NewEmail"] = plan.Email.ValueString()
	}
	if !plan.MobilePhone.Equal(state.MobilePhone) {
		request["NewMobilePhone"] = plan.MobilePhone.ValueString()
	}
	if !plan.Comments.Equal(state.Comments) {
		request["NewComments"] = plan.Comments.ValueString()
	}

	if len(request) > 1 {
		if _, err := r.Client().RpcPost("Ims", "2019-08-15", "UpdateUser", nil, request, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateUser", userTypeName),
				errs.WrapErrorf(err, "Updating %s: calling UpdateUser", userTypeName).Error(),
			)
			return
		}
	}

	object, err := getUser(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetUser", userTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetUser", userTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, userState(object))...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Client().RpcPost("Ims", "2019-08-15", "DeleteUser", nil, map[string]interface{}{
		"UserId": state.Id.ValueString(),
	}, true)
	if err != nil && !isUserNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteUser", userTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteUser", userTypeName).Error(),
		)
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getUser queries GetUser by user ID and returns the User object from the
// response, or nil when the user no longer exists.
func getUser(client *connectivity.AliyunClient, userId string) (map[string]interface{}, error) {
	response, err := client.RpcPost("Ims", "2019-08-15", "GetUser", nil, map[string]interface{}{
		"UserId": userId,
	}, true)
	if err != nil {
		if isUserNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	user, _ := response["User"].(map[string]interface{})
	if user == nil {
		return nil, fmt.Errorf("the GetUser response carried no User object: %v", response)
	}
	return user, nil
}

// userState maps the User object of a GetUser response onto the resource
// model. Empty strings become null: dates and optional fields the API omits
// read back as absent rather than "".
func userState(user map[string]interface{}) userResourceModel {
	str := func(key string) types.String {
		if v, ok := user[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return userResourceModel{
		Id:                str("UserId"),
		UserId:            str("UserId"),
		UserPrincipalName: str("UserPrincipalName"),
		DisplayName:       str("DisplayName"),
		Email:             str("Email"),
		MobilePhone:       str("MobilePhone"),
		Comments:          str("Comments"),
		UserName:          str("UserName"),
		ProvisionType:     str("ProvisionType"),
		CreateDate:        str("CreateDate"),
		UpdateDate:        str("UpdateDate"),
		LastLoginDate:     str("LastLoginDate"),
	}
}

// isUserNotExist reports whether err is the IMS user-not-found error. The
// code arrives on a *tea.SDKError from RpcPost; the substring fallback covers
// wrapped forms, matching IsExpectedErrors' behavior on the SDK line.
func isUserNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "EntityNotExist.User" {
		return true
	}
	return strings.Contains(err.Error(), "EntityNotExist.User")
}
