package apigateway

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const apiGatewayAppKeyTypeName = "alicloud_api_gateway_app_key"

// apiGatewayAppKeyTimeout replaces the SDKv2 resource-level default timeouts;
// terraform-plugin-framework does not support user-configurable resource
// timeouts yet.
const apiGatewayAppKeyTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &apiGatewayAppKeyResource{}
	_ resource.ResourceWithImportState = &apiGatewayAppKeyResource{}
)

// NewApiGatewayAppKeyResource is the Framework factory registered in
// ServicePackage.Resources for alicloud_api_gateway_app_key.
func NewApiGatewayAppKeyResource() resource.Resource {
	return &apiGatewayAppKeyResource{}
}

type apiGatewayAppKeyResource struct {
	fwadapt.ResourceBase
}

type apiGatewayAppKeyResourceModel struct {
	Id               types.String `tfsdk:"id"`
	AppId            types.String `tfsdk:"app_id"`
	SecretResetToken types.String `tfsdk:"secret_reset_token"`
	AppKey           types.String `tfsdk:"app_key"`
	AppSecret        types.String `tfsdk:"app_secret"`
	AppCode          types.String `tfsdk:"app_code"`
	CreatedTime      types.String `tfsdk:"created_time"`
	ModifiedTime     types.String `tfsdk:"modified_time"`
}

func (r *apiGatewayAppKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = apiGatewayAppKeyTypeName
}

func (r *apiGatewayAppKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an API Gateway App Key resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the App Key, formatted as `<app_id>:<app_key>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the App to which the App Key belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// AppKey, AppSecret and AppCode are server-generated credentials.
			// They are read-only Computed so that exposing them does not
			// re-introduce a perpetual plan diff. AppSecret rotation no longer
			// requires recreating the resource: changing secret_reset_token
			// drives an in-place ResetAppSecret (CloudAPI/2016-07-14) in Update.
			// Changing the AppKey identity itself would change the resource's
			// primary identifier and is therefore a ForceNew, handled via the
			// app_id RequiresReplace above (the id is `<app_id>:<app_key>`).
			"app_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The AppKey generated for the App.",
			},
			"app_secret": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The AppSecret generated for the App. Rotate it by changing secret_reset_token.",
			},
			"secret_reset_token": schema.StringAttribute{
				Optional:    true,
				Description: "An arbitrary value that, when changed, rotates the AppSecret by calling ResetAppSecret. The value itself is not sent to the server; only a change between configurations triggers rotation.",
			},
			"app_code": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The AppCode generated for the App.",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the App Key was created.",
			},
			"modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the App Key was last modified.",
			},
		},
	}
}

func (r *apiGatewayAppKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiGatewayAppKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appIdStr := plan.AppId.ValueString()
	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: parsing app_id", apiGatewayAppKeyTypeName),
			errs.WrapErrorf(err, "Creating %s: app_id %q is not an integer", apiGatewayAppKeyTypeName, appIdStr).Error(),
		)
		return
	}

	cloudApiService := NewCloudApiService(r.Meta())

	// CreateAppKey's response only carries a RequestId, so the server-generated
	// AppKey must be read back from DescribeAppSecurities. Snapshot the existing
	// AppKeys first to identify the new one by diff. This is best-effort: a
	// concurrent CreateAppKey on the same App could surface another key.
	existing := make(map[string]struct{})
	if before, err := cloudApiService.DescribeAppSecurities(appId); err != nil {
		// A NotFound here means the App has no AppKey yet (or is gone); the
		// subsequent CreateAppKey surfaces a missing-App error if applicable.
		if !errs.NotFoundError(err) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating %s: calling DescribeAppSecurities", apiGatewayAppKeyTypeName),
				err.Error(),
			)
			return
		}
	} else {
		for _, m := range before {
			existing[fmt.Sprint(m["AppKey"])] = struct{}{}
		}
	}

	// CreateAppKey is not idempotent: a single call avoids duplicate keys. The
	// server generates AppKey/AppSecret/AppCode since none are supplied.
	if err := cloudApiService.CreateAppKey(appIdStr, "", ""); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateAppKey", apiGatewayAppKeyTypeName),
			err.Error(),
		)
		return
	}

	// Locate the new AppKey by diffing against the snapshot; retry for the
	// eventual consistency of DescribeAppSecurities.
	var object map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(apiGatewayAppKeyTimeout, func() *retry.RetryError {
		items, e := cloudApiService.DescribeAppSecurities(appId)
		if e != nil {
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		for _, m := range items {
			if _, ok := existing[fmt.Sprint(m["AppKey"])]; !ok {
				object = m
				return nil
			}
		}
		wait()
		return retry.RetryableError(fmt.Errorf("the newly created AppKey for App %s is not yet visible", appIdStr))
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading back the new AppKey", apiGatewayAppKeyTypeName),
			errs.WrapErrorf(err, "Creating %s: reading back the new AppKey for App %s", apiGatewayAppKeyTypeName, appIdStr).Error(),
		)
		return
	}

	model := apiGatewayAppKeyState(object)
	model.Id = types.StringValue(fmt.Sprintf("%s:%s", appIdStr, fmt.Sprint(object["AppKey"])))
	model.AppId = plan.AppId
	// Persist the rotation trigger so a later change is detectable; the value
	// itself never reaches the server.
	model.SecretResetToken = plan.SecretResetToken
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *apiGatewayAppKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiGatewayAppKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts, err := helper.ParseResourceId(state.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing id", apiGatewayAppKeyTypeName),
			errs.WrapErrorf(err, "Reading %s: parsing id %q", apiGatewayAppKeyTypeName, state.Id.ValueString()).Error(),
		)
		return
	}
	appIdStr, appKey := parts[0], parts[1]
	appId, perr := strconv.ParseInt(appIdStr, 10, 64)
	if perr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing app_id", apiGatewayAppKeyTypeName),
			errs.WrapErrorf(perr, "Reading %s: app_id %q is not an integer", apiGatewayAppKeyTypeName, appIdStr).Error(),
		)
		return
	}

	cloudApiService := NewCloudApiService(r.Meta())
	object, err := cloudApiService.FindAppSecurity(appId, appKey)
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeAppSecurities", apiGatewayAppKeyTypeName),
			err.Error(),
		)
		return
	}
	model := apiGatewayAppKeyState(object)
	model.Id = state.Id
	model.AppId = types.StringValue(appIdStr)
	// DescribeAppSecurities does not return secret_reset_token; preserve the
	// stored value so refresh does not erase the trigger and cause a spurious
	// rotation on the next plan.
	model.SecretResetToken = state.SecretResetToken
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// Update rotates the AppSecret when secret_reset_token changes. The token
// value is never sent to the server: only a divergence between the configured
// value and the stored value signals an intent to rotate, mirroring the
// polardb gateway_consumer api_key_reset_token idiom. Rotation calls
// CloudAPI/2016-07-14 ResetAppSecret (only AppKey is supplied; NewAppKey is
// omitted because changing the AppKey identity is a ForceNew, and NewAppSecret
// is omitted so the server auto-generates a fresh secret). ResetAppSecret's
// response only carries a RequestId, so the new AppSecret is read back
// through DescribeAppSecurities (FindAppSecurity) before setting state.
// app_id is RequiresReplace, so no other attribute reaches Update.
func (r *apiGatewayAppKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state apiGatewayAppKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts, err := helper.ParseResourceId(state.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing id", apiGatewayAppKeyTypeName),
			errs.WrapErrorf(err, "Updating %s: parsing id %q", apiGatewayAppKeyTypeName, state.Id.ValueString()).Error(),
		)
		return
	}
	appIdStr, appKey := parts[0], parts[1]
	appId, perr := strconv.ParseInt(appIdStr, 10, 64)
	if perr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing app_id", apiGatewayAppKeyTypeName),
			errs.WrapErrorf(perr, "Updating %s: app_id %q is not an integer", apiGatewayAppKeyTypeName, appIdStr).Error(),
		)
		return
	}

	cloudApiService := NewCloudApiService(r.Meta())

	// Rotate only when the trigger actually changed and is set. An unset or
	// unchanged token leaves the AppSecret untouched.
	if !plan.SecretResetToken.Equal(state.SecretResetToken) && !plan.SecretResetToken.IsNull() && plan.SecretResetToken.ValueString() != "" {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(apiGatewayAppKeyTimeout, func() *retry.RetryError {
			if e := cloudApiService.ResetAppSecret(appKey); e != nil {
				if helper.NeedRetry(e) {
					wait()
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			return nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ResetAppSecret", apiGatewayAppKeyTypeName),
				errs.WrapErrorf(err, "Updating %s: calling ResetAppSecret for AppKey %s", apiGatewayAppKeyTypeName, appKey).Error(),
			)
			return
		}
	}

	// Read back the (possibly rotated) AppSecret. ResetAppSecret's response
	// carries only a RequestId, so the current values come from
	// DescribeAppSecurities. This also converges eventual consistency after
	// rotation.
	var object map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(apiGatewayAppKeyTimeout, func() *retry.RetryError {
		m, e := cloudApiService.FindAppSecurity(appId, appKey)
		if e != nil {
			if errs.NotFoundError(e) {
				resp.State.RemoveResource(ctx)
				return nil
			}
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		object = m
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: reading back the AppKey", apiGatewayAppKeyTypeName),
			errs.WrapErrorf(err, "Updating %s: reading back the AppKey for App %s", apiGatewayAppKeyTypeName, appIdStr).Error(),
		)
		return
	}
	if object == nil {
		// FindAppSecurity returned NotFound and removed the resource from state.
		return
	}

	model := apiGatewayAppKeyState(object)
	model.Id = state.Id
	model.AppId = plan.AppId
	model.SecretResetToken = plan.SecretResetToken
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *apiGatewayAppKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiGatewayAppKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts, err := helper.ParseResourceId(state.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing id", apiGatewayAppKeyTypeName),
			errs.WrapErrorf(err, "Deleting %s: parsing id %q", apiGatewayAppKeyTypeName, state.Id.ValueString()).Error(),
		)
		return
	}
	appIdStr, appKey := parts[0], parts[1]

	cloudApiService := NewCloudApiService(r.Meta())
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(apiGatewayAppKeyTimeout, func() *retry.RetryError {
		if e := cloudApiService.DeleteAppKey(appIdStr, appKey); e != nil {
			if helper.NeedRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, appKeyExpectedNotFoundCodes) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteAppKey", apiGatewayAppKeyTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteAppKey", apiGatewayAppKeyTypeName).Error(),
		)
		return
	}
}

func (r *apiGatewayAppKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// apiGatewayAppKeyState maps a DescribeAppSecurities AppSecurity item onto the
// resource model. AppId and Id are set by the callers, since the AppSecurity
// item carries no AppId. fmt.Sprint covers the string, json.Number and numeric
// decodings the tea-rpc client may produce for these otherwise-string fields.
func apiGatewayAppKeyState(object map[string]interface{}) apiGatewayAppKeyResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key]; ok && v != nil {
			if s := fmt.Sprint(v); s != "" {
				return types.StringValue(s)
			}
		}
		return types.StringNull()
	}
	return apiGatewayAppKeyResourceModel{
		AppKey:       str("AppKey"),
		AppSecret:    str("AppSecret"),
		AppCode:      str("AppCode"),
		CreatedTime:  str("CreatedTime"),
		ModifiedTime: str("ModifiedTime"),
	}
}
