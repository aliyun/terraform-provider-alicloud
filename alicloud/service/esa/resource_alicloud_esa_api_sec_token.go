package esa

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const apiSecTokenTypeName = "alicloud_esa_api_sec_token"

// apiSecTokenTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Update/Delete each 5 minutes); terraform-plugin-framework does not
// support user-configurable resource timeouts yet.
const apiSecTokenTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &apiSecTokenResource{}
	_ resource.ResourceWithImportState = &apiSecTokenResource{}
)

func NewApiSecTokenResource() resource.Resource {
	return &apiSecTokenResource{}
}

type apiSecTokenResource struct {
	fwadapt.ResourceBase
}

// apiSecTokenResourceModel holds the resource state.
//
// TokenSources, Credentials, ApiRuleIds, EndpointIds and IdentifierIds are
// modeled as opaque JSON strings: the frozen CloudSpec evidence
// (cloudspecCommit 5ff43bdc) types them as array<struct>/array<JWK> but does
// not enumerate their inner members, and the AMP withMeta supplementary read
// never materialized. They pass through verbatim (no fabricated inner
// attributes), round-trip on read for drift detection, and remain subject to
// ESA product-owner confirmation before the Private APIs publish.
//
// SiteId and the composite Id are not part of the GetApiSecToken TokenInfos
// output, so apiSecTokenState leaves them null and each caller overlays the
// plan/state value (cf. snapshot_lock's write-only DryRun overlay).
type apiSecTokenResourceModel struct {
	Id               types.String `tfsdk:"id"`
	SiteId           types.Int64  `tfsdk:"site_id"`
	TokenId          types.Int64  `tfsdk:"token_id"`
	TokenName        types.String `tfsdk:"token_name"`
	TokenType        types.String `tfsdk:"token_type"`
	TokenSources     types.String `tfsdk:"token_sources"`
	Credentials      types.String `tfsdk:"credentials"`
	CredentialsCount types.Int64  `tfsdk:"credentials_count"`
	ApiRuleIds       types.String `tfsdk:"api_rule_ids"`
	LastUpdateAt     types.String `tfsdk:"last_update_at"`
	EndpointIds      types.String `tfsdk:"endpoint_ids"`
	IdentifierIds    types.String `tfsdk:"identifier_ids"`
}

func (r *apiSecTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = apiSecTokenTypeName
}

func (r *apiSecTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// SiteId is an identity component (routePlan identity=[SiteId,TokenId]);
			// changing it means a different token, hence RequiresReplace.
			"site_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"token_id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"token_name": schema.StringAttribute{
				Required: true,
			},
			// TokenType is an enum per CloudSpec, but the enum members are not
			// enumerated in the frozen evidence; a OneOf validator is therefore
			// intentionally omitted to avoid fabricating allowed values.
			"token_type": schema.StringAttribute{
				Required: true,
			},
			"token_sources": schema.StringAttribute{
				Required: true,
			},
			// Credentials is create-only: AddApiSecToken requires it but
			// SetApiSecToken omits it, so a change forces replacement.
			"credentials": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"credentials_count": schema.Int64Attribute{Computed: true},
			"api_rule_ids":      schema.StringAttribute{Computed: true},
			"last_update_at":    schema.StringAttribute{Computed: true},
			"endpoint_ids":      schema.StringAttribute{Computed: true},
			"identifier_ids":    schema.StringAttribute{Computed: true},
		},
	}
}

func (r *apiSecTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiSecTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tokenSources, err := decodeOpaqueJSON(plan.TokenSources, "token_sources")
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s: invalid token_sources", apiSecTokenTypeName), err.Error())
		return
	}
	credentials, err := decodeOpaqueJSON(plan.Credentials, "credentials")
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s: invalid credentials", apiSecTokenTypeName), err.Error())
		return
	}

	action := "AddApiSecToken"
	request := map[string]interface{}{
		"SiteId":       plan.SiteId.ValueInt64(),
		"TokenName":    plan.TokenName.ValueString(),
		"TokenType":    plan.TokenType.ValueString(),
		"TokenSources": tokenSources,
		"Credentials":  credentials,
		"RegionId":     r.Meta().RegionId,
		"ClientToken":  helper.BuildClientToken(action),
	}

	var response map[string]interface{}
	var err2 error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err2 = retry.Retry(apiSecTokenTimeout, func() *retry.RetryError {
		response, err2 = r.Meta().RpcPost("ESA", "2024-09-10", action, nil, request, true)
		if err2 != nil {
			if helper.NeedRetry(err2) {
				wait()
				return retry.RetryableError(err2)
			}
			return retry.NonRetryableError(err2)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err2 != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling AddApiSecToken", apiSecTokenTypeName),
			errs.WrapErrorf(err2, "Creating %s: calling AddApiSecToken", apiSecTokenTypeName).Error(),
		)
		return
	}

	tokenId, err := extractTokenId(response)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: parsing TokenId from AddApiSecToken", apiSecTokenTypeName),
			err.Error(),
		)
		return
	}
	id := fmt.Sprintf("%d:%s", plan.SiteId.ValueInt64(), tokenId)

	object, err := NewService(r.Meta()).DescribeApiSecToken(id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetApiSecToken", apiSecTokenTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetApiSecToken", apiSecTokenTypeName).Error(),
		)
		return
	}
	state := apiSecTokenState(object)
	state.Id = types.StringValue(id)
	state.SiteId = plan.SiteId
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *apiSecTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiSecTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Meta()).DescribeApiSecToken(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetApiSecToken", apiSecTokenTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetApiSecToken", apiSecTokenTypeName).Error(),
		)
		return
	}
	model := apiSecTokenState(object)
	model.Id = state.Id
	model.SiteId = state.SiteId
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *apiSecTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state apiSecTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// SetApiSecToken accepts SiteId, TokenId, TokenName, TokenType and
	// TokenSources; it does not accept Credentials (create-only).
	update := false
	if !plan.TokenName.Equal(state.TokenName) {
		update = true
	}
	if !plan.TokenType.Equal(state.TokenType) {
		update = true
	}
	if !plan.TokenSources.Equal(state.TokenSources) {
		update = true
	}

	if update {
		tokenSources, err := decodeOpaqueJSON(plan.TokenSources, "token_sources")
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Updating %s: invalid token_sources", apiSecTokenTypeName), err.Error())
			return
		}
		action := "SetApiSecToken"
		request := map[string]interface{}{
			"SiteId":       state.SiteId.ValueInt64(),
			"TokenId":      state.TokenId.ValueInt64(),
			"TokenName":    plan.TokenName.ValueString(),
			"TokenType":    plan.TokenType.ValueString(),
			"TokenSources": tokenSources,
			"RegionId":     r.Meta().RegionId,
			"ClientToken":  helper.BuildClientToken(action),
		}

		var response map[string]interface{}
		var err2 error
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err2 = retry.Retry(apiSecTokenTimeout, func() *retry.RetryError {
			response, err2 = r.Meta().RpcPost("ESA", "2024-09-10", action, nil, request, true)
			if err2 != nil {
				if helper.NeedRetry(err2) {
					wait()
					return retry.RetryableError(err2)
				}
				return retry.NonRetryableError(err2)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err2 != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling SetApiSecToken", apiSecTokenTypeName),
				errs.WrapErrorf(err2, "Updating %s: calling SetApiSecToken", apiSecTokenTypeName).Error(),
			)
			return
		}
	}

	object, err := NewService(r.Meta()).DescribeApiSecToken(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetApiSecToken", apiSecTokenTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetApiSecToken", apiSecTokenTypeName).Error(),
		)
		return
	}
	model := apiSecTokenState(object)
	model.Id = state.Id
	model.SiteId = state.SiteId
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *apiSecTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiSecTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts := strings.SplitN(state.Id.ValueString(), ":", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: invalid id", apiSecTokenTypeName),
			fmt.Sprintf("invalid ApiSecToken resource id %q: expected SiteId:TokenId", state.Id.ValueString()),
		)
		return
	}
	action := "DeleteApiSecToken"
	query := map[string]interface{}{
		"SiteId":   parts[0],
		"TokenId":  parts[1],
		"RegionId": r.Meta().RegionId,
	}
	request := make(map[string]interface{})

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(apiSecTokenTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcGet("ESA", "2024-09-10", action, query, request)
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
		if errs.IsExpectedErrors(err, []string{"SiteNotFound", "TokenNotFound", "101"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteApiSecToken", apiSecTokenTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteApiSecToken", apiSecTokenTypeName).Error(),
		)
		return
	}
}

func (r *apiSecTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// extractTokenId reads the API-assigned TokenId from an AddApiSecToken
// response. The frozen CloudSpec evidence gives backendName=data.content.Id;
// the top-level "Id"/"TokenId" fallbacks cover envelope variants without
// fabricating a path, and surface a clear error if none match.
func extractTokenId(response map[string]interface{}) (string, error) {
	if v, err := jsonpath.Get("$.data.content.Id", response); err == nil && v != nil {
		return fmt.Sprint(v), nil
	}
	for _, key := range []string{"Id", "TokenId"} {
		if v, ok := response[key]; ok && v != nil {
			return fmt.Sprint(v), nil
		}
	}
	return "", fmt.Errorf("TokenId not found in AddApiSecToken response: %v", response)
}

// apiSecTokenState maps a GetApiSecToken TokenInfos item onto the resource
// model. Numeric fields are decoded across json.Number/float64/int/int64/
// string encodings the tea-rpc client may produce, mirroring the snapshot_lock
// Framework mapper. Opaque array fields are re-serialized to stable JSON.
// Id and SiteId are intentionally left null: GetApiSecToken does not return
// them, so callers overlay the plan/state value.
func apiSecTokenState(object map[string]interface{}) apiSecTokenResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	flexStr := func(key string) types.String {
		switch v := object[key].(type) {
		case string:
			if v != "" {
				return types.StringValue(v)
			}
		case json.Number:
			return types.StringValue(v.String())
		case float64:
			return types.StringValue(strconv.FormatFloat(v, 'f', -1, 64))
		case int:
			return types.StringValue(strconv.Itoa(v))
		case int64:
			return types.StringValue(strconv.FormatInt(v, 10))
		case bool:
			return types.StringValue(strconv.FormatBool(v))
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
	return apiSecTokenResourceModel{
		TokenId:          int64Val("TokenId"),
		TokenName:        str("TokenName"),
		TokenType:        str("TokenType"),
		TokenSources:     encodeOpaqueJSON(object["TokenSources"]),
		Credentials:      encodeOpaqueJSON(object["Credentials"]),
		CredentialsCount: int64Val("CredentialsCount"),
		ApiRuleIds:       encodeOpaqueJSON(object["ApiRuleIds"]),
		LastUpdateAt:     flexStr("LastUpdateAt"),
		EndpointIds:      encodeOpaqueJSON(object["EndpointIds"]),
		IdentifierIds:    encodeOpaqueJSON(object["IdentifierIds"]),
	}
}

// decodeOpaqueJSON parses an opaque JSON string attribute into a generic value
// for the ESA request body so the API receives a real array/object rather than
// a string-encoded blob. Null/empty yields nil (field omitted).
func decodeOpaqueJSON(s types.String, field string) (interface{}, error) {
	if s.IsNull() || s.ValueString() == "" {
		return nil, nil
	}
	var v interface{}
	if err := json.Unmarshal([]byte(s.ValueString()), &v); err != nil {
		return nil, fmt.Errorf("decoding %s: %s", field, err.Error())
	}
	return v, nil
}

// encodeOpaqueJSON re-serializes a raw response field to a stable JSON string
// for state. json.Marshal sorts map keys alphabetically, so drift detection is
// deterministic for equal structures.
func encodeOpaqueJSON(raw interface{}) types.String {
	if raw == nil {
		return types.StringNull()
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return types.StringNull()
	}
	return types.StringValue(string(b))
}
