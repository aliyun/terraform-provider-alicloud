package ims

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const imsAppSecretTypeName = "alicloud_ims_app_secret"

var (
	_ resource.Resource                = &imsAppSecretResource{}
	_ resource.ResourceWithImportState = &imsAppSecretResource{}
)

// NewImsAppSecretResource returns a Framework resource for the IMS AppSecret.
//
// An application secret is an immutable credential bound to an IMS
// application: it is created by CreateAppSecret, read by GetAppSecret, listed
// by ListAppSecretIds and destroyed by DeleteAppSecret. There is no in-place
// update, so any change to the user-configurable app_id forces replacement.
// The AppSecretId is server-generated and is, together with AppId, the
// composite identity imported as "app_id:app_secret_id".
func NewImsAppSecretResource() resource.Resource {
	return &imsAppSecretResource{}
}

type imsAppSecretResource struct {
	fwadapt.ResourceBase
}

type imsAppSecretResourceModel struct {
	Id             types.String `tfsdk:"id"`
	AppId          types.String `tfsdk:"app_id"`
	AppSecretId    types.String `tfsdk:"app_secret_id"`
	AppSecretValue types.String `tfsdk:"app_secret_value"`
	CreateDate     types.String `tfsdk:"create_date"`
}

func (r *imsAppSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = imsAppSecretTypeName
}

func (r *imsAppSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an IMS App Secret resource. " +
			"An application secret is a credential of an IMS application. " +
			"The application (referenced by `app_id`) is managed out-of-band; " +
			"this resource only creates and destroys a single secret for it. " +
			"The secret is immutable: changing `app_id` forces replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, in the format `<app_id>:<app_secret_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Required: true,
				Description: "The ID of the IMS application that the secret belongs to. " +
					"The application is managed out-of-band (there is no `alicloud_ims_application` resource).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"app_secret_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the application secret, server-generated.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_secret_value": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The value of the application secret, usable as an OAuth client secret. Only returned by CreateAppSecret and GetAppSecret.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_date": schema.StringAttribute{
				Computed:    true,
				Description: "The time the secret was created, in RFC 3339 format (UTC).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *imsAppSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan imsAppSecretResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId := plan.AppId.ValueString()
	response, err := r.Meta().RpcPost("Ims", "2019-08-15", "CreateAppSecret", nil, map[string]interface{}{
		"AppId": appId,
	}, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateAppSecret", imsAppSecretTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateAppSecret", imsAppSecretTypeName).Error(),
		)
		return
	}

	appSecret, _ := response["AppSecret"].(map[string]interface{})
	appSecretId, _ := appSecret["AppSecretId"].(string)
	if appSecretId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateAppSecret", imsAppSecretTypeName),
			fmt.Sprintf("The response carried no AppSecret.AppSecretId: %v", response),
		)
		return
	}

	// Read back via GetAppSecret so the state reflects the authoritative view
	// of the server rather than the create-only response projection.
	object, err := getImsAppSecret(r.Meta(), appId, appSecretId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAppSecret", imsAppSecretTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAppSecret", imsAppSecretTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateAppSecret", imsAppSecretTypeName),
			fmt.Sprintf("The app secret %s:%s is not visible via GetAppSecret right after creation", appId, appSecretId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, imsAppSecretState(object))...)
}

func (r *imsAppSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state imsAppSecretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId, appSecretId, err := parseImsAppSecretId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", imsAppSecretTypeName),
			err.Error(),
		)
		return
	}

	object, err := getImsAppSecret(r.Meta(), appId, appSecretId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAppSecret", imsAppSecretTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAppSecret", imsAppSecretTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, imsAppSecretState(object))...)
}

// Update satisfies the resource.Resource interface. The secret is immutable:
// app_id carries RequiresReplace and every other attribute is Computed, so the
// framework never reaches Update with a real diff. Should it be reached (for
// example to reconcile a drifted computed attribute), the authoritative state
// is read back from GetAppSecret rather than applied from the plan.
func (r *imsAppSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state imsAppSecretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId, appSecretId, err := parseImsAppSecretId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing resource id", imsAppSecretTypeName),
			err.Error(),
		)
		return
	}

	object, err := getImsAppSecret(r.Meta(), appId, appSecretId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetAppSecret", imsAppSecretTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetAppSecret", imsAppSecretTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, imsAppSecretState(object))...)
}

func (r *imsAppSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state imsAppSecretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId, appSecretId, err := parseImsAppSecretId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing resource id", imsAppSecretTypeName),
			err.Error(),
		)
		return
	}

	_, err = r.Meta().RpcPost("Ims", "2019-08-15", "DeleteAppSecret", nil, map[string]interface{}{
		"AppId":       appId,
		"AppSecretId": appSecretId,
	}, true)
	if err != nil && !isImsAppSecretNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteAppSecret", imsAppSecretTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteAppSecret", imsAppSecretTypeName).Error(),
		)
		return
	}
}

func (r *imsAppSecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getImsAppSecret queries GetAppSecret by app id and secret id and returns the
// AppSecret object from the response, or nil when the secret no longer exists.
func getImsAppSecret(client *connectivity.AliyunClient, appId, appSecretId string) (map[string]interface{}, error) {
	response, err := client.RpcPost("Ims", "2019-08-15", "GetAppSecret", nil, map[string]interface{}{
		"AppId":       appId,
		"AppSecretId": appSecretId,
	}, true)
	if err != nil {
		if isImsAppSecretNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	appSecret, _ := response["AppSecret"].(map[string]interface{})
	if appSecret == nil {
		return nil, fmt.Errorf("the GetAppSecret response carried no AppSecret object: %v", response)
	}
	return appSecret, nil
}

// imsAppSecretState maps the AppSecret object of a GetAppSecret response onto
// the resource model. Empty strings become null: the create date and secret
// value read back absent rather than "".
func imsAppSecretState(appSecret map[string]interface{}) imsAppSecretResourceModel {
	str := func(key string) types.String {
		if v, ok := appSecret[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	appId := str("AppId")
	appSecretId := str("AppSecretId")
	var id types.String
	if appId.ValueString() != "" && appSecretId.ValueString() != "" {
		id = types.StringValue(fmt.Sprintf("%s:%s", appId.ValueString(), appSecretId.ValueString()))
	} else {
		id = types.StringNull()
	}
	return imsAppSecretResourceModel{
		Id:             id,
		AppId:          appId,
		AppSecretId:    appSecretId,
		AppSecretValue: str("AppSecretValue"),
		CreateDate:     str("CreateDate"),
	}
}

// parseImsAppSecretId splits the composite resource id "app_id:app_secret_id".
func parseImsAppSecretId(id string) (appId, appSecretId string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid resource id %q: expected <app_id>:<app_secret_id>", id)
	}
	return parts[0], parts[1], nil
}

// isImsAppSecretNotExist reports whether err means the secret (or its parent
// application) is gone. The code arrives on a *tea.SDKError from RpcPost; the
// substring fallback covers wrapped forms, matching isUserNotExist's behavior.
func isImsAppSecretNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil &&
		(*e.Code == "EntityNotExist.AppSecret" || *e.Code == "EntityNotExist.Application") {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "EntityNotExist.AppSecret") ||
		strings.Contains(msg, "EntityNotExist.Application")
}
