// Package ddoscoo. This file is a Framework-native handwritten resource for the
// alicloud_ddoscoo_web_cache_custom_rule resource, implementing the static page cache custom
// rule lifecycle against the ddoscoo/2020-01-01 OpenAPI family.
//
// Lifecycle mapping (verified against the published OpenAPI on 2026-11-05):
//   - Create / Update: ModifyWebCacheCustomRule — sets one named custom cache rule (Domain + Rules JSON).
//   - Read / List:     DescribeWebCacheConfigs — returns DomainCacheConfigs[].CustomRules[].
//   - Delete:          DeleteWebCacheCustomRule — removes rules by RuleNames.
//
// ModifyWebCacheCustomRule takes the rule set as a JSON array of named rules and the separate
// Delete API removes rules by name, the same per-rule upsert shape the existing
// alicloud_ddoscoo_web_cc_rule resource uses against ConfigWebCCRuleV2 / DeleteWebCCRuleV2; the
// resource therefore owns a single named rule and writes a one-element rule set, matching the PD
// route plan assumption. Domain and Name form the primary key (Domain:Name) and are both
// ForceNew. ResourceGroupId is accepted by every API but is not returned by DescribeWebCacheConfigs,
// so it is not modelled: a write-only optional attribute would drift on every plan, and the
// existing ddoscoo resources do not model it either.
package ddoscoo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const webCacheCustomRuleTypeName = "alicloud_ddoscoo_web_cache_custom_rule"

// ddoscoo product / version / actions used by this resource.
const (
	webCacheProduct  = "ddoscoo"
	webCacheVersion  = "2020-01-01"
	actionModifyRule = "ModifyWebCacheCustomRule"
	actionDeleteRule = "DeleteWebCacheCustomRule"
	actionDescribe   = "DescribeWebCacheConfigs"

	// The DdosCoo describe APIs are eventually consistent and retry on the shared
	// transient classifier, matching the SDKv2 ddoscoo service helpers.
	webCacheRetryTimeout = 1 * time.Minute
)

var (
	_ resource.Resource                = &webCacheCustomRuleResource{}
	_ resource.ResourceWithImportState = &webCacheCustomRuleResource{}
)

func NewWebCacheCustomRuleResource() resource.Resource {
	return &webCacheCustomRuleResource{}
}

type webCacheCustomRuleResource struct {
	fwadapt.ResourceBase
}

type webCacheCustomRuleModel struct {
	Id        types.String `tfsdk:"id"`
	Domain    types.String `tfsdk:"domain"`
	Name      types.String `tfsdk:"name"`
	Uri       types.String `tfsdk:"uri"`
	Mode      types.String `tfsdk:"mode"`
	CacheTtl  types.Int64  `tfsdk:"cache_ttl"`
	RequestId types.String `tfsdk:"request_id"`
}

func (r *webCacheCustomRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = webCacheCustomRuleTypeName
}

func (r *webCacheCustomRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a DdosCoo (Anti-DDoS Pro) Web Cache Custom Rule resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, in the format `<domain>:<name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "The website domain that the static page cache custom rule belongs to. The domain must already be configured with a website forwarding rule and associated with an enhanced package DDoS Pro instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the custom cache rule. Together with `domain` it forms the resource ID; changing it forces a replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uri": schema.StringAttribute{
				Required:    true,
				Description: "The path of the cached page, e.g. `/blog/`.",
			},
			"mode": schema.StringAttribute{
				Required:    true,
				Description: "The cache mode. Valid values: `standard` (standard mode), `aggressive` (aggressive mode) and `bypass` (no caching).",
				Validators: []validator.String{
					stringvalidator.OneOf("standard", "aggressive", "bypass"),
				},
			},
			"cache_ttl": schema.Int64Attribute{
				Required:    true,
				Description: "The expiration time of the cached page, in seconds.",
			},
			"request_id": schema.StringAttribute{
				Computed:    true,
				Description: "The request ID returned by the most recent ModifyWebCacheCustomRule call, for traceability. It is refreshed on every apply and is not stable across reads.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *webCacheCustomRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan webCacheCustomRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := webCacheCustomRuleModify(ctx, r.Meta(), plan.Domain.ValueString(), plan.Name.ValueString(), plan.Uri.ValueString(), plan.Mode.ValueString(), plan.CacheTtl.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", webCacheCustomRuleTypeName, actionModifyRule),
			errs.WrapErrorf(err, "Creating %s: calling %s", webCacheCustomRuleTypeName, actionModifyRule).Error(),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s:%s", plan.Domain.ValueString(), plan.Name.ValueString()))
	plan.RequestId = webCacheRequestID(response)

	// Set the ID before the read so an eventual-consistency miss during refresh does not lose the
	// identity the create just persisted; the read below reconciles the attributes.
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	object, err := webCacheCustomRuleGet(ctx, r.Meta(), plan.Domain.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", webCacheCustomRuleTypeName, actionDescribe),
			errs.WrapErrorf(err, "Reading %s: calling %s", webCacheCustomRuleTypeName, actionDescribe).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling %s", webCacheCustomRuleTypeName, actionModifyRule),
			fmt.Sprintf("The rule %s:%s was not visible via %s right after creation", plan.Domain.ValueString(), plan.Name.ValueString(), actionDescribe),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, webCacheCustomRuleState(object, plan.Domain.ValueString(), plan.RequestId))...)
}

func (r *webCacheCustomRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webCacheCustomRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, name := webCacheCustomRuleParseID(state.Id.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := webCacheCustomRuleGet(ctx, r.Meta(), domain, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", webCacheCustomRuleTypeName, actionDescribe),
			errs.WrapErrorf(err, "Reading %s: calling %s", webCacheCustomRuleTypeName, actionDescribe).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, webCacheCustomRuleState(object, domain, state.RequestId))...)
}

// webCacheCustomRuleParseID splits the `<domain>:<name>` resource ID into its parts. The domain
// itself never contains a colon in this API family, so a single split is sufficient; a malformed ID
// is reported as a diagnostic rather than a silent panic, matching the SDKv2 service helpers.
func webCacheCustomRuleParseID(id string, diags *diag.Diagnostics) (string, string) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		diags.AddError(
			fmt.Sprintf("Parsing %s ID", webCacheCustomRuleTypeName),
			fmt.Sprintf("The resource ID %q could not be parsed into `<domain>:<name>`: expected exactly two non-empty parts separated by a single colon.", id),
		)
		return "", ""
	}
	return parts[0], parts[1]
}

func (r *webCacheCustomRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state webCacheCustomRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Domain and Name are RequiresReplace, so the upsert here only re-sends the writable attributes
	// (uri / mode / cache_ttl). ModifyWebCacheCustomRule is a per-rule set, so a single-element rule
	// set keyed by the existing name is enough; only re-issue the write when a mutable attribute
	// changed to avoid a needless call on every plan.
	if !plan.Uri.Equal(state.Uri) || !plan.Mode.Equal(state.Mode) || !plan.CacheTtl.Equal(state.CacheTtl) {
		response, err := webCacheCustomRuleModify(ctx, r.Meta(), state.Domain.ValueString(), state.Name.ValueString(), plan.Uri.ValueString(), plan.Mode.ValueString(), plan.CacheTtl.ValueInt64())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling %s", webCacheCustomRuleTypeName, actionModifyRule),
				errs.WrapErrorf(err, "Updating %s: calling %s", webCacheCustomRuleTypeName, actionModifyRule).Error(),
			)
			return
		}
		plan.RequestId = webCacheRequestID(response)
	} else {
		plan.RequestId = state.RequestId
	}

	plan.Id = state.Id
	plan.Domain = state.Domain
	plan.Name = state.Name

	object, err := webCacheCustomRuleGet(ctx, r.Meta(), state.Domain.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", webCacheCustomRuleTypeName, actionDescribe),
			errs.WrapErrorf(err, "Reading %s: calling %s", webCacheCustomRuleTypeName, actionDescribe).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, webCacheCustomRuleState(object, state.Domain.ValueString(), plan.RequestId))...)
}

func (r *webCacheCustomRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webCacheCustomRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, name := webCacheCustomRuleParseID(state.Id.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// RuleNames is a repeatList; the indexed RuleNames.1 form is the encoding the existing
	// DeleteWebPreciseAccessRule call in this package uses for the same shape.
	request := map[string]interface{}{
		"Domain":      domain,
		"RuleNames.1": name,
		"RegionId":    r.Meta().RegionId,
	}
	if _, err := helper.RetryRpcPost(ctx, actionDeleteRule, func() (map[string]interface{}, error) {
		return r.Meta().RpcPost(webCacheProduct, webCacheVersion, actionDeleteRule, nil, request, true)
	}, webCacheRetryTimeout, helper.NeedRetry); err != nil {
		// A missing rule or a missing domain already reached the desired state; both are swallowed.
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"DomainNotExist", "RuleNotExist", "NotFound"}) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling %s", webCacheCustomRuleTypeName, actionDeleteRule),
			errs.WrapErrorf(err, "Deleting %s: calling %s", webCacheCustomRuleTypeName, actionDeleteRule).Error(),
		)
		return
	}
}

func (r *webCacheCustomRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// webCacheCustomRuleModify issues ModifyWebCacheCustomRule for a single named rule. The Rules
// parameter is a JSON string array of rule objects, matching the published example and the SDKv2
// ConfigWebCCRuleV2 call shape used by alicloud_ddoscoo_web_cc_rule.
func webCacheCustomRuleModify(ctx context.Context, client *connectivity.AliyunClient, domain, name, uri, mode string, cacheTtl int64) (map[string]interface{}, error) {
	rule := map[string]interface{}{
		"Name":     name,
		"Uri":      uri,
		"Mode":     mode,
		"CacheTtl": cacheTtl,
	}
	request := map[string]interface{}{
		"Domain":   domain,
		"Rules":    convert.ConvertObjectToJsonString([]interface{}{rule}),
		"RegionId": client.RegionId,
	}
	return helper.RetryRpcPost(ctx, actionModifyRule, func() (map[string]interface{}, error) {
		return client.RpcPost(webCacheProduct, webCacheVersion, actionModifyRule, nil, request, true)
	}, webCacheRetryTimeout, helper.NeedRetry)
}

// webCacheCustomRuleGet calls DescribeWebCacheConfigs for one domain and returns the CustomRule
// whose Name matches, or nil when the rule (or the domain's cache config) no longer exists.
func webCacheCustomRuleGet(ctx context.Context, client *connectivity.AliyunClient, domain, name string) (map[string]interface{}, error) {
	// Domains is a repeatList; the indexed Domains.1 form is the encoding the existing
	// DescribeWebPreciseAccessRule helper in this package uses for DescribeWebCacheConfigs.
	request := map[string]interface{}{
		"Domains.1": domain,
		"RegionId":  client.RegionId,
	}
	response, err := helper.RetryRpcPost(ctx, actionDescribe, func() (map[string]interface{}, error) {
		return client.RpcPost(webCacheProduct, webCacheVersion, actionDescribe, nil, request, true)
	}, webCacheRetryTimeout, helper.NeedRetry)
	if err != nil {
		// DomainNotExist means the domain is no longer protected, so the rule is gone.
		if errs.IsExpectedErrors(err, []string{"DomainNotExist"}) {
			return nil, nil
		}
		return nil, err
	}

	// DomainCacheConfigs[].CustomRules[] is the rule list; walk both levels and return the rule
	// whose Name matches. A missing array or an empty response is treated as "rule not found".
	configs, _ := response["DomainCacheConfigs"].([]interface{})
	for _, c := range configs {
		cfg, _ := c.(map[string]interface{})
		rules, _ := cfg["CustomRules"].([]interface{})
		for _, r := range rules {
			item, _ := r.(map[string]interface{})
			if fmt.Sprint(item["Name"]) == name {
				return item, nil
			}
		}
	}
	return nil, nil
}

// webCacheCustomRuleState maps a DescribeWebCacheConfigs CustomRule onto the resource model.
// Domain is not a field of the CustomRule object (it lives on the enclosing DomainCacheConfig), so
// it is threaded in from the resource ID; Name/Uri/Mode/CacheTtl come from the rule object itself.
// request_id is preserved from the last write because DescribeWebCacheConfigs does not return it.
func webCacheCustomRuleState(object map[string]interface{}, domain string, requestId types.String) webCacheCustomRuleModel {
	str := func(key string) types.String {
		if v, ok := object[key]; ok {
			if s, ok := v.(string); ok && s != "" {
				return types.StringValue(s)
			}
		}
		return types.StringNull()
	}
	var ttl types.Int64
	if v, ok := object["CacheTtl"]; ok {
		switch n := v.(type) {
		case float64:
			ttl = types.Int64Value(int64(n))
		case int64:
			ttl = types.Int64Value(n)
		case int:
			ttl = types.Int64Value(int64(n))
		default:
			ttl = types.Int64Null()
		}
	} else {
		ttl = types.Int64Null()
	}

	name := str("Name")
	domainVal := types.StringNull()
	if domain != "" {
		domainVal = types.StringValue(domain)
	}
	m := webCacheCustomRuleModel{
		Id:        types.StringNull(),
		Domain:    domainVal,
		Name:      name,
		Uri:       str("Uri"),
		Mode:      str("Mode"),
		CacheTtl:  ttl,
		RequestId: requestId,
	}
	if !m.Domain.IsNull() && !m.Name.IsNull() {
		m.Id = types.StringValue(fmt.Sprintf("%s:%s", m.Domain.ValueString(), m.Name.ValueString()))
	}
	return m
}

// webCacheRequestID extracts the RequestId from a write response for the request_id attribute.
func webCacheRequestID(response map[string]interface{}) types.String {
	if v, ok := response["RequestId"].(string); ok && v != "" {
		return types.StringValue(v)
	}
	return types.StringNull()
}
