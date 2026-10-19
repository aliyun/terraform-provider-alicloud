package cloudapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
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

const apiGatewayDomainTypeName = "alicloud_api_gateway_domain"

// idSep separates the two parts of a domain resource id: <group_id>:<domain_name>.
// A literal ":" is used rather than a shared constant because the framework
// service packages do not import package alicloud, where COLON_SEPARATED lives.
const idSep = ":"

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

type apiGatewayDomainModel struct {
	Id                       types.String `tfsdk:"id"`
	GroupId                  types.String `tfsdk:"group_id"`
	DomainName               types.String `tfsdk:"domain_name"`
	SslOcspEnable            types.Bool   `tfsdk:"ssl_ocsp_enable"`
	SslOcspCacheEnable       types.Bool   `tfsdk:"ssl_ocsp_cache_enable"`
	ClientCertSDnPassThrough types.Bool   `tfsdk:"client_cert_s_dn_pass_through"`
	WssEnable                types.String `tfsdk:"wss_enable"`
	SubDomain                types.String `tfsdk:"sub_domain"`
	DomainBindingStatus      types.String `tfsdk:"domain_binding_status"`
	DomainRemark             types.String `tfsdk:"domain_remark"`
	DomainWebSocketStatus    types.String `tfsdk:"domain_web_socket_status"`
	DomainLegalStatus        types.String `tfsdk:"domain_legal_status"`
	DomainCNAMEStatus        types.String `tfsdk:"domain_cname_status"`
	CertificateId            types.String `tfsdk:"certificate_id"`
	CertificateName          types.String `tfsdk:"certificate_name"`
	CertificateValidStart    types.Int64  `tfsdk:"certificate_valid_start"`
	CertificateValidEnd      types.Int64  `tfsdk:"certificate_valid_end"`
}

func (r *apiGatewayDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = apiGatewayDomainTypeName
}

func (r *apiGatewayDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an API Gateway Domain resource that binds a custom domain name to an API Gateway group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. The value formats as `<group_id>:<domain_name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the API Gateway group that the custom domain binds to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "The custom domain name to bind to the API Gateway group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// ssl_ocsp_enable, ssl_ocsp_cache_enable, client_cert_s_dn_pass_through
			// and wss_enable are write-only: they are sent to the backing API only on
			// Update, and the read API (DescribeDomain) does not return them. Their
			// values are therefore preserved from the prior state on refresh.
			"ssl_ocsp_enable": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to enable SSL OCSP stapling for the bound certificate.",
			},
			"ssl_ocsp_cache_enable": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to enable SSL OCSP cache for the domain.",
			},
			"client_cert_s_dn_pass_through": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to pass through the client certificate Subject DN.",
			},
			"wss_enable": schema.StringAttribute{
				Optional:    true,
				Description: "Whether to enable WebSocket for the domain. Valid values: `ON`, `OFF`.",
			},
			"sub_domain": schema.StringAttribute{
				Computed:    true,
				Description: "The system allocated sub domain of the API Gateway group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_binding_status": schema.StringAttribute{
				Computed:    true,
				Description: "The binding status of the domain.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_remark": schema.StringAttribute{
				Computed:    true,
				Description: "The remark of the domain.",
			},
			"domain_web_socket_status": schema.StringAttribute{
				Computed:    true,
				Description: "The WebSocket status of the domain.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_legal_status": schema.StringAttribute{
				Computed:    true,
				Description: "The legal status of the domain.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_cname_status": schema.StringAttribute{
				Computed:    true,
				Description: "The CNAME status of the domain.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"certificate_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the bound SSL certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"certificate_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the bound SSL certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"certificate_valid_start": schema.Int64Attribute{
				Computed:    true,
				Description: "The valid start timestamp (milliseconds) of the bound SSL certificate.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"certificate_valid_end": schema.Int64Attribute{
				Computed:    true,
				Description: "The valid end timestamp (milliseconds) of the bound SSL certificate.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *apiGatewayDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiGatewayDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := plan.GroupId.ValueString()
	domainName := plan.DomainName.ValueString()

	if _, err := r.Client().RpcPost("CloudAPI", "2016-07-14", "SetDomain", map[string]interface{}{
		"GroupId":    groupId,
		"DomainName": domainName,
	}, nil, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling SetDomain", apiGatewayDomainTypeName),
			errs.WrapErrorf(err, "Creating %s: calling SetDomain", apiGatewayDomainTypeName).Error(),
		)
		return
	}

	plan.Id = types.StringValue(domainId(groupId, domainName))
	// The write-only attributes (ssl_ocsp_*, client_cert_s_dn_pass_through,
	// wss_enable) ride on the plan; refresh copies them back from this model,
	// since DescribeDomain never returns them.
	if err := r.applyUpdate(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(err.summary, err.detail)
		return
	}
	if removed, err := r.refresh(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(err.summary, err.detail)
		return
	} else if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiGatewayDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiGatewayDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed, err := r.refresh(ctx, &state); err != nil {
		resp.Diagnostics.AddError(err.summary, err.detail)
		return
	} else if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *apiGatewayDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan apiGatewayDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Carry the resource id and write-only attributes from the prior state so
	// applyUpdate/refresh operate on a model that already reflects the saved id.
	var state apiGatewayDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Id = state.Id

	if err := r.applyUpdate(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(err.summary, err.detail)
		return
	}
	if removed, err := r.refresh(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(err.summary, err.detail)
		return
	} else if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiGatewayDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiGatewayDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, domainName, ok := parseDomainId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing resource id", apiGatewayDomainTypeName),
			fmt.Sprintf("invalid resource id format: %s, expected <group_id>:<domain_name>", state.Id.ValueString()),
		)
		return
	}

	if _, err := r.Client().RpcPost("CloudAPI", "2016-07-14", "DeleteDomain", map[string]interface{}{
		"GroupId":    groupId,
		"DomainName": domainName,
	}, nil, true); err != nil && !isDomainNotFound(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDomain", apiGatewayDomainTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteDomain", apiGatewayDomainTypeName).Error(),
		)
		return
	}
}

func (r *apiGatewayDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// opError carries a diagnostic summary/detail pair out of applyUpdate/refresh so
// the Create/Read/Update handlers can forward it through resp.Diagnostics.
type opError struct {
	summary string
	detail  string
}

// applyUpdate sends the update-only attributes to the backing APIs. The three
// boolean SSL/OCSP/client-cert parameters ride on SetDomainCertificate; wss_enable
// rides on SetDomainWebSocketStatus. SetDomainCertificate's Go SDK request struct
// does not carry them as fields, so they are sent as query parameters, matching the
// SDK v2 resource behavior. RpcPost forwards the query map verbatim.
func (r *apiGatewayDomainResource) applyUpdate(ctx context.Context, plan *apiGatewayDomainModel) *opError {
	groupId, domainName, ok := parseDomainId(plan.Id.ValueString())
	if !ok {
		return &opError{
			summary: fmt.Sprintf("Updating %s: parsing resource id", apiGatewayDomainTypeName),
			detail:  fmt.Sprintf("invalid resource id format: %s, expected <group_id>:<domain_name>", plan.Id.ValueString()),
		}
	}

	// SetDomainCertificate carries SslOcspEnable, SslOcspCacheEnable and
	// ClientCertSDnPassThrough. Send the request whenever the plan carries any of
	// them: a false value is meaningful and must reach the API rather than be
	// dropped as "unset".
	if !plan.SslOcspEnable.IsNull() || !plan.SslOcspCacheEnable.IsNull() || !plan.ClientCertSDnPassThrough.IsNull() {
		if _, err := r.Client().RpcPost("CloudAPI", "2016-07-14", "SetDomainCertificate", map[string]interface{}{
			"GroupId":                  groupId,
			"DomainName":               domainName,
			"SslOcspEnable":            plan.SslOcspEnable.ValueBool(),
			"SslOcspCacheEnable":       plan.SslOcspCacheEnable.ValueBool(),
			"ClientCertSDnPassThrough": plan.ClientCertSDnPassThrough.ValueBool(),
		}, nil, true); err != nil {
			return &opError{summary: fmt.Sprintf("Updating %s: calling SetDomainCertificate", apiGatewayDomainTypeName),
				detail: errs.WrapErrorf(err, "Updating %s: calling SetDomainCertificate", apiGatewayDomainTypeName).Error()}
		}
	}

	if !plan.WssEnable.IsNull() {
		if _, err := r.Client().RpcPost("CloudAPI", "2016-07-14", "SetDomainWebSocketStatus", map[string]interface{}{
			"GroupId":    groupId,
			"DomainName": domainName,
			"WSSEnable":  plan.WssEnable.ValueString(),
		}, nil, true); err != nil {
			return &opError{summary: fmt.Sprintf("Updating %s: calling SetDomainWebSocketStatus", apiGatewayDomainTypeName),
				detail: errs.WrapErrorf(err, "Updating %s: calling SetDomainWebSocketStatus", apiGatewayDomainTypeName).Error()}
		}
	}
	return nil
}

// refresh calls DescribeDomain, removes the resource when it is gone, otherwise
// fills the computed fields from the API. The four write-only attributes are
// preserved from the prior state: DescribeDomain never returns them.
func (r *apiGatewayDomainResource) refresh(ctx context.Context, state *apiGatewayDomainModel) (bool, *opError) {
	groupId, domainName, ok := parseDomainId(state.Id.ValueString())
	if !ok {
		return false, &opError{
			summary: fmt.Sprintf("Reading %s: parsing resource id", apiGatewayDomainTypeName),
			detail:  fmt.Sprintf("invalid resource id format: %s, expected <group_id>:<domain_name>", state.Id.ValueString()),
		}
	}

	response, err := describeDomain(r.Client(), groupId, domainName)
	if err != nil {
		return false, &opError{
			summary: fmt.Sprintf("Reading %s: calling DescribeDomain", apiGatewayDomainTypeName),
			detail:  errs.WrapErrorf(err, "Reading %s: calling DescribeDomain", apiGatewayDomainTypeName).Error(),
		}
	}
	if response == nil {
		return true, nil
	}

	state.GroupId = strValue(response, "GroupId")
	state.DomainName = strValue(response, "DomainName")
	state.Id = types.StringValue(domainId(groupId, domainName))
	state.SubDomain = strValue(response, "SubDomain")
	state.DomainBindingStatus = strValue(response, "DomainBindingStatus")
	state.DomainRemark = strValue(response, "DomainRemark")
	state.DomainWebSocketStatus = strValue(response, "DomainWebSocketStatus")
	state.DomainLegalStatus = strValue(response, "DomainLegalStatus")
	state.DomainCNAMEStatus = strValue(response, "DomainCNAMEStatus")
	state.CertificateId = strValue(response, "CertificateId")
	state.CertificateName = strValue(response, "CertificateName")
	state.CertificateValidStart = int64Value(response, "CertificateValidStart")
	state.CertificateValidEnd = int64Value(response, "CertificateValidEnd")
	return false, nil
}

func describeDomain(client domainClient, groupId, domainName string) (map[string]interface{}, error) {
	response, err := client.RpcPost("CloudAPI", "2016-07-14", "DescribeDomain", map[string]interface{}{
		"GroupId":    groupId,
		"DomainName": domainName,
	}, nil, true)
	if err != nil {
		if isDomainNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// domainClient is the subset of *connectivity.AliyunClient that refresh/applyUpdate
// use, extracted so the DescribeDomain helper is testable without a live client.
type domainClient interface {
	RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error)
}

func domainId(groupId, domainName string) string {
	return groupId + idSep + domainName
}

func parseDomainId(id string) (string, string, bool) {
	parts := strings.SplitN(id, idSep, 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func strValue(response map[string]interface{}, key string) types.String {
	if v, ok := response[key].(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func int64Value(response map[string]interface{}, key string) types.Int64 {
	switch v := response[key].(type) {
	case float64:
		if v != 0 {
			return types.Int64Value(int64(v))
		}
	case int64:
		if v != 0 {
			return types.Int64Value(v)
		}
	case int:
		if v != 0 {
			return types.Int64Value(int64(v))
		}
	}
	return types.Int64Null()
}

// isDomainNotFound reports whether err is the CloudAPI domain-not-found error.
// DescribeDomain surfaces NotFoundDomain/NotFound as a *tea.SDKError code from
// RpcPost; the substring fallback covers wrapped forms.
func isDomainNotFound(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		switch *e.Code {
		case "NotFoundDomain", "NotFound", "ResourceNotFound":
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "NotFoundDomain")
}
