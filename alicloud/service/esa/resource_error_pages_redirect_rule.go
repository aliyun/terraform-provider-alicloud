package esa

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// errorPagesRedirectRuleProduct and errorPagesRedirectRuleVersion pin the ESA
	// RPC product/version used by every Create/Read/Update/Delete call. They are
	// kept local so a rename here is the only edit point.
	errorPagesRedirectRuleProduct = "ESA"
	errorPagesRedirectRuleVersion = "2024-09-10"
)

var (
	_ resource.Resource                = &errorPagesRedirectRuleResource{}
	_ resource.ResourceWithImportState = &errorPagesRedirectRuleResource{}
)

// NewErrorPagesRedirectRuleResource is the factory the framework provider
// registry calls to instantiate the resource.
func NewErrorPagesRedirectRuleResource() resource.Resource {
	return &errorPagesRedirectRuleResource{}
}

type errorPagesRedirectRuleResource struct {
	fwadapt.ResourceBase
}

// errorPagesRedirectItemType is the framework object type of a single
// error_pages_redirect block, used to (re)build a types.List for state. The
// attribute keys mirror the tfsdk tags on errorPagesRedirectItem below.
var errorPagesRedirectItemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"target_url":  types.StringType,
		"status_code": types.StringType,
	},
}

type errorPagesRedirectItem struct {
	TargetUrl  types.String `tfsdk:"target_url"`
	StatusCode types.String `tfsdk:"status_code"`
}

type errorPagesRedirectRuleModel struct {
	Id                 types.String `tfsdk:"id"`
	SiteId             types.String `tfsdk:"site_id"`
	ConfigId           types.String `tfsdk:"config_id"`
	ConfigType         types.String `tfsdk:"config_type"`
	ErrorPagesRedirect types.List   `tfsdk:"error_pages_redirect"`
	Rule               types.String `tfsdk:"rule"`
	RuleEnable         types.String `tfsdk:"rule_enable"`
	RuleName           types.String `tfsdk:"rule_name"`
	Sequence           types.Int64  `tfsdk:"sequence"`
	SiteVersion        types.Int64  `tfsdk:"site_version"`
}

func (r *errorPagesRedirectRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = errorPagesRedirectRuleTypeName
}

func (r *errorPagesRedirectRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an ESA ErrorPagesRedirectRule resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID in the format `<site_id>:<config_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the ESA site. Changing this creates a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the error pages redirect rule, assigned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"config_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"error_pages_redirect": schema.ListNestedAttribute{
				Required:    true,
				Description: "The error pages redirect targets.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"target_url": schema.StringAttribute{
							Required:    true,
							Description: "The URL to which the error page is redirected.",
						},
						"status_code": schema.StringAttribute{
							Required:    true,
							Description: "The HTTP status code that triggers the redirect.",
						},
					},
				},
			},
			"rule": schema.StringAttribute{
				Optional:    true,
				Description: "The rule expression of the error pages redirect rule.",
			},
			"rule_enable": schema.StringAttribute{
				Optional:    true,
				Description: "Whether the rule is enabled. Valid values: `on`, `off`.",
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off"),
				},
			},
			"rule_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the error pages redirect rule.",
			},
			"sequence": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The sequence of the rule. Changing this creates a new resource.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"site_version": schema.Int64Attribute{
				Optional:    true,
				Description: "The version of the ESA site. Changing this creates a new resource.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *errorPagesRedirectRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan errorPagesRedirectRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"SiteId": plan.SiteId.ValueString(),
	}
	if !plan.Sequence.IsNull() {
		request["Sequence"] = plan.Sequence.ValueInt64()
	}
	if !plan.SiteVersion.IsNull() {
		request["SiteVersion"] = plan.SiteVersion.ValueInt64()
	}
	if epr, ok := errorPagesRedirectRequestJSON(ctx, plan.ErrorPagesRedirect, &resp.Diagnostics); ok {
		request["ErrorPagesRedirect"] = epr
	} else {
		return
	}
	if !plan.RuleEnable.IsNull() {
		request["RuleEnable"] = plan.RuleEnable.ValueString()
	}
	if !plan.RuleName.IsNull() {
		request["RuleName"] = plan.RuleName.ValueString()
	}
	if !plan.Rule.IsNull() {
		request["Rule"] = plan.Rule.ValueString()
	}

	response, err := r.Client().RpcPost(errorPagesRedirectRuleProduct, errorPagesRedirectRuleVersion, "CreateErrorPagesRedirectRule", nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateErrorPagesRedirectRule", errorPagesRedirectRuleTypeName).Error(),
		)
		return
	}

	configId := toString(response["ConfigId"])
	if configId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			fmt.Sprintf("The response carried no ConfigId: %v", response),
		)
		return
	}
	siteId := plan.SiteId.ValueString()

	object, err := getErrorPagesRedirectRule(r.Client(), siteId, configId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			fmt.Sprintf("The rule %s:%s is not visible via GetErrorPagesRedirectRule right after creation", siteId, configId),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, errorPagesRedirectRuleState(ctx, siteId, configId, object, &resp.Diagnostics))...)
}

func (r *errorPagesRedirectRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state errorPagesRedirectRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteId, configId, diags := parseErrorPagesRedirectRuleID(state.Id.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getErrorPagesRedirectRule(r.Client(), siteId, configId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, errorPagesRedirectRuleState(ctx, siteId, configId, object, &resp.Diagnostics))...)
}

func (r *errorPagesRedirectRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state errorPagesRedirectRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteId, configId, diags := parseErrorPagesRedirectRuleID(state.Id.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"SiteId":   siteId,
		"ConfigId": configId,
	}
	if !plan.Sequence.Equal(state.Sequence) {
		if !plan.Sequence.IsNull() {
			request["Sequence"] = plan.Sequence.ValueInt64()
		}
	}
	if !plan.SiteVersion.Equal(state.SiteVersion) {
		if !plan.SiteVersion.IsNull() {
			request["SiteVersion"] = plan.SiteVersion.ValueInt64()
		}
	}
	if epr, ok := errorPagesRedirectRequestJSON(ctx, plan.ErrorPagesRedirect, &resp.Diagnostics); ok {
		request["ErrorPagesRedirect"] = epr
	} else {
		return
	}
	if !plan.RuleEnable.Equal(state.RuleEnable) {
		request["RuleEnable"] = plan.RuleEnable.ValueString()
	}
	if !plan.RuleName.Equal(state.RuleName) {
		request["RuleName"] = plan.RuleName.ValueString()
	}
	if !plan.Rule.Equal(state.Rule) {
		request["Rule"] = plan.Rule.ValueString()
	}

	if _, err := r.Client().RpcPost(errorPagesRedirectRuleProduct, errorPagesRedirectRuleVersion, "UpdateErrorPagesRedirectRule", nil, request, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateErrorPagesRedirectRule", errorPagesRedirectRuleTypeName).Error(),
		)
		return
	}

	object, err := getErrorPagesRedirectRule(r.Client(), siteId, configId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetErrorPagesRedirectRule", errorPagesRedirectRuleTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, errorPagesRedirectRuleState(ctx, siteId, configId, object, &resp.Diagnostics))...)
}

func (r *errorPagesRedirectRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state errorPagesRedirectRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteId, configId, diags := parseErrorPagesRedirectRuleID(state.Id.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Client().RpcPost(errorPagesRedirectRuleProduct, errorPagesRedirectRuleVersion, "DeleteErrorPagesRedirectRule", nil, map[string]interface{}{
		"SiteId":   siteId,
		"ConfigId": configId,
	}, true)
	if err != nil && !errs.NotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteErrorPagesRedirectRule", errorPagesRedirectRuleTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteErrorPagesRedirectRule", errorPagesRedirectRuleTypeName).Error(),
		)
		return
	}
}

func (r *errorPagesRedirectRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import token must be `<site_id>:<config_id>`. ImportStatePassthroughID
	// stores it on `id`; Read then parses it with the same strict length check
	// used by Read/Update/Delete, so a malformed import token fails cleanly
	// instead of panicking on an out-of-range slice index.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// parseErrorPagesRedirectRuleID splits a composite `<site_id>:<config_id>` ID
// into its parts. Any token that is not exactly two colon-separated parts is
// rejected with an error; callers must surface it and return, never index.
func parseErrorPagesRedirectRuleID(id string) (siteId, configId string, diags diag.Diagnostics) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		diags.AddError(
			fmt.Sprintf("Parsing %s ID", errorPagesRedirectRuleTypeName),
			fmt.Sprintf("invalid resource ID %q: expected exactly two non-empty parts separated by a colon (<site_id>:<config_id>), got %d part(s)", id, len(parts)),
		)
		return "", "", diags
	}
	return parts[0], parts[1], diags
}

// getErrorPagesRedirectRule queries GetErrorPagesRedirectRule and returns the
// raw response object, or nil when the rule no longer exists.
func getErrorPagesRedirectRule(client *connectivity.AliyunClient, siteId, configId string) (map[string]interface{}, error) {
	response, err := client.RpcGet(errorPagesRedirectRuleProduct, errorPagesRedirectRuleVersion, "GetErrorPagesRedirectRule", map[string]interface{}{
		"SiteId":   siteId,
		"ConfigId": configId,
	}, nil)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// errorPagesRedirectRuleState maps a GetErrorPagesRedirectRule response onto the
// resource model. The composite id is rebuilt from the parsed parts so the state
// is authoritative even when the response omits SiteId.
func errorPagesRedirectRuleState(ctx context.Context, siteId, configId string, obj map[string]interface{}, diags *diag.Diagnostics) errorPagesRedirectRuleModel {
	return errorPagesRedirectRuleModel{
		Id:                 types.StringValue(fmt.Sprintf("%s:%s", siteId, configId)),
		SiteId:             types.StringValue(siteId),
		ConfigId:           types.StringValue(configId),
		ConfigType:         stringVal(obj["ConfigType"]),
		ErrorPagesRedirect: errorPagesRedirectList(ctx, obj["ErrorPagesRedirect"], diags),
		Rule:               stringVal(obj["Rule"]),
		RuleEnable:         stringVal(obj["RuleEnable"]),
		RuleName:           stringVal(obj["RuleName"]),
		Sequence:           int64Val(obj["Sequence"]),
		SiteVersion:        int64Val(obj["SiteVersion"]),
	}
}

// errorPagesRedirectRequestJSON marshals a plan's error_pages_redirect list
// into the JSON string the ESA Create/Update API expects
// ([{"StatusCode":..,"TargetURL":..},...]). ok=false means a decoding error
// was appended and the caller must abort.
func errorPagesRedirectRequestJSON(ctx context.Context, list types.List, diags *diag.Diagnostics) (string, bool) {
	if list.IsNull() || list.IsUnknown() {
		return "[]", true
	}
	var items []errorPagesRedirectItem
	diags.Append(list.ElementsAs(ctx, &items, false)...)
	if diags.HasError() {
		return "", false
	}
	arr := make([]map[string]interface{}, 0, len(items))
	for _, it := range items {
		arr = append(arr, map[string]interface{}{
			"StatusCode": it.StatusCode.ValueString(),
			"TargetURL":  it.TargetUrl.ValueString(),
		})
	}
	b, err := json.Marshal(arr)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Creating/updating %s", errorPagesRedirectRuleTypeName),
			fmt.Sprintf("marshalling error_pages_redirect: %v", err),
		)
		return "", false
	}
	return string(b), true
}

// errorPagesRedirectList builds the types.List for state from the raw API
// response value. nil/missing maps to an empty (non-null) list.
func errorPagesRedirectList(ctx context.Context, raw interface{}, diags *diag.Diagnostics) types.List {
	items := make([]errorPagesRedirectItem, 0)
	if raw != nil {
		if arr, ok := raw.([]interface{}); ok {
			for _, e := range arr {
				m, _ := e.(map[string]interface{})
				items = append(items, errorPagesRedirectItem{
					TargetUrl:  stringVal(m["TargetURL"]),
					StatusCode: stringVal(m["StatusCode"]),
				})
			}
		}
	}
	list, ldiags := types.ListValueFrom(ctx, errorPagesRedirectItemType, items)
	diags.Append(ldiags...)
	return list
}

// stringVal converts a response scalar to types.String, treating empty/missing
// as null so optional and computed attributes read back absent rather than "".
func stringVal(v interface{}) types.String {
	switch x := v.(type) {
	case string:
		if x == "" {
			return types.StringNull()
		}
		return types.StringValue(x)
	case nil:
		return types.StringNull()
	default:
		s := fmt.Sprint(x)
		if s == "" {
			return types.StringNull()
		}
		return types.StringValue(s)
	}
}

// int64Val converts a response scalar to types.Int64. Missing/invalid maps to
// null, matching the Optional+Computed behavior of sequence and site_version.
func int64Val(v interface{}) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	switch x := v.(type) {
	case float64:
		return types.Int64Value(int64(x))
	case int:
		return types.Int64Value(int64(x))
	case int64:
		return types.Int64Value(x)
	case string:
		if x == "" {
			return types.Int64Null()
		}
		n, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return types.Int64Null()
		}
		return types.Int64Value(n)
	default:
		return types.Int64Null()
	}
}

// toString flattens a response scalar to a plain string for use as an ID
// component. JSON numbers arrive as float64, so the integer form is preferred.
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatInt(int64(x), 10)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	default:
		return fmt.Sprint(v)
	}
}
