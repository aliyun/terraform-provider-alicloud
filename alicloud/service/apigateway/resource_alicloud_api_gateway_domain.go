package apigateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const apiGatewayDomainTypeName = "alicloud_api_gateway_domain"

// apiGatewayDomainTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Update/Delete each 5 minutes); terraform-plugin-framework does not
// support user-configurable resource timeouts yet.
const apiGatewayDomainTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &apiGatewayDomainResource{}
	_ resource.ResourceWithImportState = &apiGatewayDomainResource{}
)

func NewApiGatewayDomainResource() resource.Resource {
	return &apiGatewayDomainResource{}
}

type apiGatewayDomainResource struct {
	fwadapt.ResourceBase
}

type apiGatewayDomainResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	GroupId               types.String `tfsdk:"group_id"`
	DomainName            types.String `tfsdk:"domain_name"`
	BindStageName         types.String `tfsdk:"bind_stage_name"`
	CustomDomainType      types.String `tfsdk:"custom_domain_type"`
	IsHttpRedirectToHttps types.Bool   `tfsdk:"is_http_redirect_to_https"`
	IsForce               types.Bool   `tfsdk:"is_force"`
	SubDomain             types.String `tfsdk:"sub_domain"`
	DomainBindingStatus   types.String `tfsdk:"domain_binding_status"`
	DomainLegalStatus     types.String `tfsdk:"domain_legal_status"`
	DomainWebSocketStatus types.String `tfsdk:"domain_websocket_status"`
	DomainRemark          types.String `tfsdk:"domain_remark"`
	DomainCNAMEStatus     types.String `tfsdk:"domain_cname_status"`
	CertificateId         types.String `tfsdk:"certificate_id"`
	CertificateName       types.String `tfsdk:"certificate_name"`
	CertificateValidStart types.Int64  `tfsdk:"certificate_valid_start"`
	CertificateValidEnd   types.Int64  `tfsdk:"certificate_valid_end"`
}

func (r *apiGatewayDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = apiGatewayDomainTypeName
}

func (r *apiGatewayDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// BindStageName is write-only: DescribeDomain never returns it, so the
			// configured value is preserved across refresh. An empty string means the
			// domain can access all environments (TEST/PRE/RELEASE).
			"bind_stage_name": schema.StringAttribute{
				Optional: true,
			},
			"custom_domain_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("INTERNET", "INTRANET"),
				},
			},
			"is_http_redirect_to_https": schema.BoolAttribute{
				Optional: true,
			},
			"is_force": schema.BoolAttribute{
				Optional: true,
			},
			// The API-managed computed fields carry no UseStateForUnknown: binding
			// status, CNAME resolution and certificate timestamps follow server-side
			// state transitions, so the planned value could not survive an apply.
			"sub_domain":              schema.StringAttribute{Computed: true},
			"domain_binding_status":   schema.StringAttribute{Computed: true},
			"domain_legal_status":     schema.StringAttribute{Computed: true},
			"domain_websocket_status": schema.StringAttribute{Computed: true},
			"domain_remark":           schema.StringAttribute{Computed: true},
			"domain_cname_status":     schema.StringAttribute{Computed: true},
			"certificate_id":          schema.StringAttribute{Computed: true},
			"certificate_name":        schema.StringAttribute{Computed: true},
			"certificate_valid_start": schema.Int64Attribute{Computed: true},
			"certificate_valid_end":   schema.Int64Attribute{Computed: true},
		},
	}
}

func (r *apiGatewayDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiGatewayDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.bindDomain(plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling SetDomain", apiGatewayDomainTypeName),
			err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(domainId(plan.GroupId.ValueString(), plan.DomainName.ValueString()))
	object, err := NewApiGatewayServiceV2(r.Meta()).DescribeApiGatewayDomain(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDomain", apiGatewayDomainTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeDomain", apiGatewayDomainTypeName).Error(),
		)
		return
	}
	state := apiGatewayDomainState(object)
	state.preserveWriteOnly(plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *apiGatewayDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiGatewayDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewApiGatewayServiceV2(r.Meta()).DescribeApiGatewayDomain(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDomain", apiGatewayDomainTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeDomain", apiGatewayDomainTypeName).Error(),
		)
		return
	}
	model := apiGatewayDomainState(object)
	model.preserveWriteOnly(state)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *apiGatewayDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan apiGatewayDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.bindDomain(plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling SetDomain", apiGatewayDomainTypeName),
			err.Error(),
		)
		return
	}

	object, err := NewApiGatewayServiceV2(r.Meta()).DescribeApiGatewayDomain(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDomain", apiGatewayDomainTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeDomain", apiGatewayDomainTypeName).Error(),
		)
		return
	}
	state := apiGatewayDomainState(object)
	state.preserveWriteOnly(plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *apiGatewayDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiGatewayDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, domainName, err := splitApiGatewayDomainId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing id", apiGatewayDomainTypeName),
			err.Error(),
		)
		return
	}
	action := "DeleteDomain"
	request := map[string]interface{}{
		"GroupId":    groupId,
		"DomainName": domainName,
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(apiGatewayDomainTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		// DeleteDomain returns success for a domain that no longer exists, so a
		// NotFound-style error is treated as a completed delete.
		if errs.IsExpectedErrors(err, []string{"NotFoundDomain"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDomain", apiGatewayDomainTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteDomain", apiGatewayDomainTypeName).Error(),
		)
	}
}

func (r *apiGatewayDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// bindDomain calls SetDomain with the configured attributes. SetDomain binds a
// custom domain to an API group and is idempotent, so it serves both Create and
// Update. group_id and domain_name are RequiresReplace, so the update path only
// reaches here for changes to the bind configuration.
func (r *apiGatewayDomainResource) bindDomain(plan apiGatewayDomainResourceModel) error {
	request := map[string]interface{}{
		"GroupId":    plan.GroupId.ValueString(),
		"DomainName": plan.DomainName.ValueString(),
	}
	if !plan.BindStageName.IsNull() {
		request["BindStageName"] = plan.BindStageName.ValueString()
	}
	if !plan.CustomDomainType.IsNull() {
		request["CustomDomainType"] = plan.CustomDomainType.ValueString()
	}
	if !plan.IsHttpRedirectToHttps.IsNull() {
		request["IsHttpRedirectToHttps"] = plan.IsHttpRedirectToHttps.ValueBool()
	}
	if !plan.IsForce.IsNull() {
		request["IsForce"] = plan.IsForce.ValueBool()
	}

	action := "SetDomain"
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(apiGatewayDomainTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, domainId(plan.GroupId.ValueString(), plan.DomainName.ValueString()), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// preserveWriteOnly overlays the write-only attributes from a prior model onto the
// refreshed state. DescribeDomain never returns BindStageName, CustomDomainType,
// IsHttpRedirectToHttps or IsForce, so they must be carried from plan/state.
func (m *apiGatewayDomainResourceModel) preserveWriteOnly(from apiGatewayDomainResourceModel) {
	m.BindStageName = from.BindStageName
	m.CustomDomainType = from.CustomDomainType
	m.IsHttpRedirectToHttps = from.IsHttpRedirectToHttps
	m.IsForce = from.IsForce
}

func domainId(groupId, domainName string) string {
	return fmt.Sprintf("%s%s%s", groupId, COLON_SEPARATED, domainName)
}

// apiGatewayDomainState maps a DescribeDomain response onto the resource model.
// The tea-rpc client decodes numeric fields as json.Number/int rather than
// float64, so every numeric encoding the transport produces must map back to
// Int64 or planned values would collapse to null ("inconsistent result after apply").
func apiGatewayDomainState(object map[string]interface{}) apiGatewayDomainResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	int64Val := func(key string) types.Int64 {
		switch v := object[key].(type) {
		case float64:
			return types.Int64Value(int64(v))
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.Int64Value(n)
			}
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				return types.Int64Value(n)
			}
		}
		return types.Int64Null()
	}
	groupId := str("GroupId")
	domainName := str("DomainName")
	var id types.String
	if !groupId.IsNull() && !domainName.IsNull() {
		id = types.StringValue(domainId(groupId.ValueString(), domainName.ValueString()))
	}
	return apiGatewayDomainResourceModel{
		Id:                    id,
		GroupId:               groupId,
		DomainName:            domainName,
		SubDomain:             str("SubDomain"),
		DomainBindingStatus:   str("DomainBindingStatus"),
		DomainLegalStatus:     str("DomainLegalStatus"),
		DomainWebSocketStatus: str("DomainWebSocketStatus"),
		DomainRemark:          str("DomainRemark"),
		DomainCNAMEStatus:     str("DomainCNAMEStatus"),
		CertificateId:         str("CertificateId"),
		CertificateName:       str("CertificateName"),
		CertificateValidStart: int64Val("CertificateValidStart"),
		CertificateValidEnd:   int64Val("CertificateValidEnd"),
	}
}
