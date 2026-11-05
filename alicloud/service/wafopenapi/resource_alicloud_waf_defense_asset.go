// Package wafopenapi. Framework-native alicloud_waf_defense_asset resource.
//
// DefenseAsset is a WAF protective asset — a protected target (an AI prompt
// endpoint, an API-sec API, an account-protector login surface, ...) fronted by
// a WAF instance. It is managed by the waf-openapi/2021-10-01 private API set
// whose authoritative online contract (zhenyuan released-entity) is:
//
//   - CreateDefenseAsset(InstanceId, AssetType, Resource, Conditions[], AssetConfig,
//     [AssetName], [Description], [AssetOrigin], [AssetStatus], [ResourceManagerResourceGroupId])
//     -> AssetId
//   - DescribeDefenseAssets(InstanceId, [AssetIds], ...) -> {Assets[], TotalCount, ...}
//   - ModifyDefenseAsset(InstanceId, AssetId, AssetConfig, [Description], [AssetName],
//     [ResourceManagerResourceGroupId])   // Resource/Conditions/AssetType/AssetStatus are NOT accepted
//   - ModifyDefenseAssetStatus(InstanceId, AssetId, AssetStatus)  // int32, "0"|"1"
//   - DeleteDefenseAsset(InstanceId, AssetId)  // Defense.Control.DefenseResourceNotExist when gone
//
// AssetType is a fixed enum validated on Create: AI, Client_Injection, Apisec,
// Account_Protector, SSO. Resource and Conditions are required on Create and are
// not accepted by Modify, so they are immutable (RequiresReplace). AssetConfig is
// a required, type-dependent JSON object (style "json"): for AssetType "AI" it carries
// PromptPosition/ResponseContents/InjectionRatio/IsPayment/SensitiveLevel/OperationPurpose/
// SuccessConditions/FailureConditions/ExtractAccounts/CookieDomain; for "Apisec" it carries
// Swagger; for "Account_Protector" it carries ExtractAccounts/CookieDomain. Because the shape
// depends on AssetType, AssetConfig is exposed as an opaque canonical JSON string the caller
// supplies; Conditions likewise is exposed as an opaque canonical JSON string of the
// [{Domain, UrlPath, HttpMethod}] array. A JSON-semantic plan modifier suppresses formatting-
// only drift between the caller's JSON and the service echo so neither field churns.
package wafopenapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const defenseAssetTypeName = "alicloud_waf_defense_asset"

// defenseAssetValidAssetTypes is the CreateDefenseAsset AssetType enum, taken from the
// authoritative online (zhenyuan) published API definition. CreateDefenseAsset rejects
// unknown values with InvalidAssetType, so the schema validates against this set.
var defenseAssetValidAssetTypes = []string{"AI", "Client_Injection", "Apisec", "Account_Protector", "SSO"}

var (
	_ resource.Resource                = &defenseAssetResource{}
	_ resource.ResourceWithImportState = &defenseAssetResource{}
)

func NewWafDefenseAssetResource() resource.Resource {
	return &defenseAssetResource{}
}

type defenseAssetResource struct {
	fwadapt.ResourceBase
}

type defenseAssetResourceModel struct {
	Id          types.String `tfsdk:"id"`
	FwTimeouts  types.Object `tfsdk:"timeouts"`
	InstanceId  types.String `tfsdk:"instance_id"`
	AssetType   types.String `tfsdk:"asset_type"`
	Resource    types.String `tfsdk:"resource"`
	Conditions  types.String `tfsdk:"conditions"`
	AssetConfig types.String `tfsdk:"asset_config"`
	AssetName   types.String `tfsdk:"asset_name"`
	Description types.String `tfsdk:"description"`
	AssetStatus types.String `tfsdk:"asset_status"`
}

func (r *defenseAssetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = defenseAssetTypeName
}

func (r *defenseAssetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The AssetId of the defense asset, assigned by CreateDefenseAsset.",
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The WAF instance id owning the defense asset. Required on every DefenseAsset API call; changing it recreates the resource.",
			},
			"asset_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(defenseAssetValidAssetTypes...),
				},
				Description: "The type of the protective asset. CreateDefenseAsset validates this against a fixed enum and rejects unknown values with InvalidAssetType; changing it recreates the resource. Valid values: AI, Client_Injection, Apisec, Account_Protector, SSO.",
			},
			"resource": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The protected resource identifier (e.g. an API or domain reference) bound to the asset at creation. Required by CreateDefenseAsset and not accepted by ModifyDefenseAsset, so changing it recreates the resource.",
			},
			"conditions": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					jsonStringSemanticEquals{},
					stringplanmodifier.RequiresReplace(),
				},
				Description: "A JSON array of match conditions, each {\"Domain\":string,\"UrlPath\":string,\"HttpMethod\":string}. Required by CreateDefenseAsset and not accepted by ModifyDefenseAsset, so changing it recreates the resource. Stored as canonical JSON; formatting-only differences are not treated as changes.",
			},
			"asset_config": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					jsonStringSemanticEquals{},
				},
				Description: "The asset configuration as a JSON object, required by CreateDefenseAsset and ModifyDefenseAsset. The shape depends on asset_type (AI: InjectionRatio/PromptPosition/ResponseContents/IsPayment/SensitiveLevel/OperationPurpose/SuccessConditions/FailureConditions/ExtractAccounts/CookieDomain; Apisec: Swagger; Account_Protector: ExtractAccounts/CookieDomain). Updated in place through ModifyDefenseAsset. Stored as canonical JSON; formatting-only differences are not treated as changes.",
			},
			"asset_name": schema.StringAttribute{
				Optional:    true,
				Description: "A human-readable name for the asset. Updatable through ModifyDefenseAsset.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for the asset. Updatable through ModifyDefenseAsset.",
			},
			"asset_status": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("0", "1"),
				},
				Description: "The protection status of the asset, toggled through ModifyDefenseAssetStatus. \"1\" means enabled, \"0\" means disabled. Defaults to 1 (enabled) on creation when unset.",
			},
		},
	}
}

// jsonStringSemanticEquals is a string plan modifier that suppresses the diff between a
// caller-supplied JSON string and the service echo when the two are semantically equal
// (same structure after key-sorted canonicalization). It prevents perpetual diffs where the
// service reorders or re-formats an opaque JSON field the caller owns. When the values differ
// semantically it leaves the plan value untouched so the change drives the real update.
type jsonStringSemanticEquals struct{}

func (jsonStringSemanticEquals) Description(ctx context.Context) string {
	return "suppresses the diff when the plan and state JSON are semantically equal"
}

func (jsonStringSemanticEquals) MarkdownDescription(ctx context.Context) string {
	return jsonStringSemanticEquals{}.Description(ctx)
}

func (jsonStringSemanticEquals) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() || req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if jsonStringsEqual(req.PlanValue.ValueString(), req.StateValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}

// canonicalJSON unmarshals a JSON document into an interface value and re-marshals it so that
// object keys are sorted; the resulting string is a stable canonical form suitable for equality.
func canonicalJSON(s string) (string, bool) {
	var i interface{}
	if err := json.Unmarshal([]byte(s), &i); err != nil {
		return "", false
	}
	b, err := json.Marshal(i)
	if err != nil {
		return "", false
	}
	return string(b), true
}

// jsonStringsEqual reports whether two JSON strings are semantically equal. If either fails to
// parse it falls back to a plain string comparison so a non-JSON value still compares honestly.
func jsonStringsEqual(a, b string) bool {
	ca, oka := canonicalJSON(a)
	cb, okb := canonicalJSON(b)
	if !oka || !okb {
		return a == b
	}
	return ca == cb
}

// toCanonicalJSONString normalizes a response field that may arrive as a JSON string, a decoded
// object/array, or nil into a canonical JSON types.String. Empty/missing -> StringNull.
func toCanonicalJSONString(v interface{}) types.String {
	var s string
	switch x := v.(type) {
	case nil:
		return types.StringNull()
	case string:
		s = x
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return types.StringNull()
		}
		s = string(b)
	}
	if c, ok := canonicalJSON(s); ok {
		s = c
	}
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// toStatusString normalizes the AssetStatus response (int32 in JSON, possibly a string) into
// the "0"/"1" types.String the schema uses. Empty/missing -> StringNull.
func toStatusString(v interface{}) types.String {
	switch x := v.(type) {
	case nil:
		return types.StringNull()
	case string:
		if x == "" {
			return types.StringNull()
		}
		return types.StringValue(x)
	case float64:
		return types.StringValue(strconv.Itoa(int(x)))
	case int:
		return types.StringValue(strconv.Itoa(x))
	case int64:
		return types.StringValue(strconv.FormatInt(x, 10))
	}
	return types.StringNull()
}

// readString navigates a dotted key in a response object and returns a types.String
// (missing or non-string -> StringNull). It mirrors the package-local helpers in other
// framework service packages without dragging their full mapping machinery.
func readString(obj map[string]interface{}, key string) types.String {
	cur := interface{}(obj)
	for _, seg := range strings.Split(key, ".") {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return types.StringNull()
		}
		cur, ok = m[seg]
		if !ok {
			return types.StringNull()
		}
	}
	if v, ok := cur.(string); ok {
		if v == "" {
			return types.StringNull()
		}
		return types.StringValue(v)
	}
	return types.StringNull()
}

// defenseAssetRead describes a single asset by calling DescribeDefenseAssets with an AssetIds
// filter and matching by AssetId. Returns (nil, nil) when the asset is gone so the caller
// removes state. DescribeDefenseAssets does not echo InstanceId on each asset, so the caller
// preserves InstanceId from prior state.
func defenseAssetRead(ctx context.Context, client *connectivity.AliyunClient, instanceId, assetId string, timeout time.Duration) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"InstanceId": instanceId,
		"AssetIds":   assetId,
	}
	response, err := wafRetryRpcPost(ctx, client, "DescribeDefenseAssets", request, timeout)
	if err != nil {
		if defenseAssetNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	assets, _ := response["Assets"].([]interface{})
	for _, raw := range assets {
		obj, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := obj["AssetId"].(string); ok && id == assetId {
			return obj, nil
		}
	}
	return nil, nil
}

// defenseAssetState maps the Describe response object onto the resource model. InstanceId is
// preserved from prior state because DescribeDefenseAssets does not echo it per asset; the
// JSON fields (asset_config, conditions) are canonicalized so they round-trip without drift.
func defenseAssetState(obj map[string]interface{}, prior *defenseAssetResourceModel) defenseAssetResourceModel {
	m := defenseAssetResourceModel{
		FwTimeouts:  types.ObjectNull(fwTimeoutTypes()),
		Id:          readString(obj, "AssetId"),
		AssetType:   readString(obj, "AssetType"),
		Resource:    readString(obj, "Resource"),
		Conditions:  toCanonicalJSONString(obj["Conditions"]),
		AssetConfig: toCanonicalJSONString(obj["AssetConfig"]),
		AssetName:   readString(obj, "AssetName"),
		Description: readString(obj, "Description"),
		AssetStatus: toStatusString(obj["AssetStatus"]),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		// InstanceId is not echoed by DescribeDefenseAssets; keep the caller's value.
		m.InstanceId = prior.InstanceId
		if m.Id.IsNull() && !prior.Id.IsNull() {
			m.Id = prior.Id
		}
	}
	return m
}

func (r *defenseAssetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan defenseAssetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "create", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"ClientToken": helper.BuildClientToken("CreateDefenseAsset"),
		"RegionId":    r.Meta().RegionId,
		"InstanceId":  plan.InstanceId.ValueString(),
		"AssetType":   plan.AssetType.ValueString(),
		"Resource":    plan.Resource.ValueString(),
		"Conditions":  plan.Conditions.ValueString(),
		"AssetConfig": plan.AssetConfig.ValueString(),
	}
	if !plan.AssetName.IsNull() && !plan.AssetName.IsUnknown() {
		if v := plan.AssetName.ValueString(); v != "" {
			request["AssetName"] = v
		}
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		if v := plan.Description.ValueString(); v != "" {
			request["Description"] = v
		}
	}
	if !plan.AssetStatus.IsNull() && !plan.AssetStatus.IsUnknown() {
		if s := plan.AssetStatus.ValueString(); s != "" {
			request["AssetStatus"] = s
		}
	}

	response, err := wafRetryRpcPost(ctx, r.Meta(), "CreateDefenseAsset", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDefenseAsset", defenseAssetTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateDefenseAsset", defenseAssetTypeName).Error(),
		)
		return
	}

	assetId := ""
	if v, ok := response["AssetId"].(string); ok {
		assetId = v
	}
	if assetId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: missing identity", defenseAssetTypeName),
			fmt.Sprintf("CreateDefenseAsset response did not carry an AssetId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any follow-up call so a failure here is recoverable;
	// the read back then refreshes the computed fields.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(assetId)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout := operationTimeout
	object, err := defenseAssetRead(ctx, r.Meta(), plan.InstanceId.ValueString(), assetId, readTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling DescribeDefenseAssets", defenseAssetTypeName),
			errs.WrapErrorf(err, "Creating %s: DescribeDefenseAssets after create (%s)", defenseAssetTypeName, assetId).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", defenseAssetTypeName),
			fmt.Sprintf("The defense asset %s is not visible via DescribeDefenseAssets right after creation", assetId),
		)
		return
	}

	// CreateDefenseAsset accepts AssetStatus, but to be robust against a type whose create path
	// ignores it, reconcile the status through ModifyDefenseAssetStatus when the read-back differs.
	if !plan.AssetStatus.IsNull() && !plan.AssetStatus.IsUnknown() {
		desired := plan.AssetStatus.ValueString()
		current := toStatusString(object["AssetStatus"]).ValueString()
		if desired != "" && desired != current {
			statusReq := map[string]interface{}{
				"ClientToken": helper.BuildClientToken("ModifyDefenseAssetStatus"),
				"RegionId":    r.Meta().RegionId,
				"InstanceId":  plan.InstanceId.ValueString(),
				"AssetId":     assetId,
				"AssetStatus": desired,
			}
			if _, err := wafRetryRpcPost(ctx, r.Meta(), "ModifyDefenseAssetStatus", statusReq, operationTimeout); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: setting asset_status", defenseAssetTypeName),
					errs.WrapErrorf(err, "Creating %s: ModifyDefenseAssetStatus after create (%s)", defenseAssetTypeName, assetId).Error(),
				)
				return
			}
			object, err = defenseAssetRead(ctx, r.Meta(), plan.InstanceId.ValueString(), assetId, readTimeout)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: calling DescribeDefenseAssets", defenseAssetTypeName),
					errs.WrapErrorf(err, "Creating %s: DescribeDefenseAssets after status set (%s)", defenseAssetTypeName, assetId).Error(),
				)
				return
			}
			if object == nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: after create", defenseAssetTypeName),
					fmt.Sprintf("The defense asset %s disappeared after setting asset_status", assetId),
				)
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, defenseAssetState(object, &plan))...)
}

func (r *defenseAssetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state defenseAssetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "create", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := defenseAssetRead(ctx, r.Meta(), state.InstanceId.ValueString(), state.Id.ValueString(), readTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDefenseAssets", defenseAssetTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeDefenseAssets (%s)", defenseAssetTypeName, state.Id.ValueString()).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, defenseAssetState(object, &state))...)
}

func (r *defenseAssetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state defenseAssetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "update", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	assetId := state.Id.ValueString()

	// ModifyDefenseAsset accepts AssetConfig (required), Description and AssetName only; it does
	// not accept Resource/Conditions/AssetType/AssetStatus. Trigger it when any of its accepted
	// fields changed. AssetConfig is always sent because the API marks it required on Modify.
	if !plan.AssetConfig.Equal(state.AssetConfig) || !plan.AssetName.Equal(state.AssetName) || !plan.Description.Equal(state.Description) {
		request := map[string]interface{}{
			"ClientToken": helper.BuildClientToken("ModifyDefenseAsset"),
			"RegionId":    r.Meta().RegionId,
			"InstanceId":  instanceId,
			"AssetId":     assetId,
			"AssetConfig": plan.AssetConfig.ValueString(),
		}
		if !plan.AssetName.IsNull() && !plan.AssetName.IsUnknown() {
			if v := plan.AssetName.ValueString(); v != "" {
				request["AssetName"] = v
			}
		}
		if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
			if v := plan.Description.ValueString(); v != "" {
				request["Description"] = v
			}
		}
		if _, err := wafRetryRpcPost(ctx, r.Meta(), "ModifyDefenseAsset", request, operationTimeout); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyDefenseAsset", defenseAssetTypeName),
				errs.WrapErrorf(err, "Updating %s: ModifyDefenseAsset (%s)", defenseAssetTypeName, assetId).Error(),
			)
			return
		}
	}

	// A configured status change goes through the dedicated status API.
	if !plan.AssetStatus.IsNull() && !plan.AssetStatus.IsUnknown() && !plan.AssetStatus.Equal(state.AssetStatus) {
		request := map[string]interface{}{
			"ClientToken": helper.BuildClientToken("ModifyDefenseAssetStatus"),
			"RegionId":    r.Meta().RegionId,
			"InstanceId":  instanceId,
			"AssetId":     assetId,
			"AssetStatus": plan.AssetStatus.ValueString(),
		}
		if _, err := wafRetryRpcPost(ctx, r.Meta(), "ModifyDefenseAssetStatus", request, operationTimeout); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyDefenseAssetStatus", defenseAssetTypeName),
				errs.WrapErrorf(err, "Updating %s: ModifyDefenseAssetStatus (%s)", defenseAssetTypeName, assetId).Error(),
			)
			return
		}
	}

	readTimeout := operationTimeout
	object, err := defenseAssetRead(ctx, r.Meta(), instanceId, assetId, readTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling DescribeDefenseAssets", defenseAssetTypeName),
			errs.WrapErrorf(err, "Updating %s: DescribeDefenseAssets after update (%s)", defenseAssetTypeName, assetId).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, defenseAssetState(object, &plan))...)
}

func (r *defenseAssetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state defenseAssetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "delete", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"ClientToken": helper.BuildClientToken("DeleteDefenseAsset"),
		"RegionId":    r.Meta().RegionId,
		"InstanceId":  state.InstanceId.ValueString(),
		"AssetId":     state.Id.ValueString(),
	}
	if _, err := wafRetryRpcPost(ctx, r.Meta(), "DeleteDefenseAsset", request, operationTimeout); err != nil {
		// DeleteDefenseAsset reports Defense.Control.DefenseResourceNotExist once the asset is
		// gone; that is the desired terminal state for delete.
		if defenseAssetNotExist(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDefenseAsset", defenseAssetTypeName),
			errs.WrapErrorf(err, "Deleting %s: DeleteDefenseAsset (%s)", defenseAssetTypeName, state.Id.ValueString()).Error(),
		)
		return
	}
}

// ImportState accepts "<InstanceId>:<AssetId>" so the read can address the asset.
// DescribeDefenseAssets requires InstanceId on every call, so a bare AssetId cannot be
// resolved; the composite form is the only recoverable identity. asset_config and conditions
// are opaque caller-owned JSON and are not refreshed on import; supply them in config after
// import to avoid a diff, or accept the generated diff.
func (r *defenseAssetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Importing "+defenseAssetTypeName,
			fmt.Sprintf("expected the import id in <InstanceId>:<AssetId> form, got %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, defenseAssetResourceModel{
		Id:         types.StringValue(parts[1]),
		InstanceId: types.StringValue(parts[0]),
		FwTimeouts: types.ObjectNull(fwTimeoutTypes()),
	})...)
}
