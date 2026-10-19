// Package cms hosts the CloudMonitor (Cms) service's native Plugin Framework
// resources. This resource is hand-written; it is not produced by the resource
// generator.
package cms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	contextStoreAPIKeyTypeName = "alicloud_cms_context_store_api_key"

	cmsProductCode = "Cms"
	cmsAPIVersion  = "2024-03-30"

	// Bounded retry windows, matching the legacy resource's timeouts. The Cms
	// API has no Get or Update; Create and Delete occasionally need to wait out
	// transient ServiceUnavailable responses.
	createTimeout = 5 * time.Minute
	deleteTimeout = 5 * time.Minute
	readTimeout   = 1 * time.Minute
)

// notFoundErrorCodes are returned by the Cms ROA API when the workspace, the
// context store or the API key no longer exists. Delete treats them as success.
var notFoundErrorCodes = []string{
	"WorkspaceNotExist",
	"ContextStoreNotExist",
	"NotFound",
	"EntityNotExist",
}

var (
	_ resource.Resource                = (*contextStoreAPIKeyResource)(nil)
	_ resource.ResourceWithImportState = (*contextStoreAPIKeyResource)(nil)
)

func NewContextStoreAPIKeyResource() resource.Resource {
	return &contextStoreAPIKeyResource{}
}

type contextStoreAPIKeyResource struct {
	fwadapt.ResourceBase
}

type contextStoreAPIKeyModel struct {
	Id               types.String `tfsdk:"id"`
	Workspace        types.String `tfsdk:"workspace"`
	ContextStoreName types.String `tfsdk:"context_store_name"`
	Name             types.String `tfsdk:"name"`
	APIKey           types.String `tfsdk:"api_key"`
	CreateTime       types.String `tfsdk:"create_time"`
	RegionId         types.String `tfsdk:"region_id"`
}

func (r *contextStoreAPIKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = contextStoreAPIKeyTypeName
}

func (r *contextStoreAPIKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Cms Context Store API Key resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. The value is formulated as `<workspace>:<context_store_name>:<name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace": schema.StringAttribute{
				Required:    true,
				Description: "The name of the workspace to which the context store belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"context_store_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the context store.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the API key, which is used to identify the purpose of the key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// The API does not expose the key value on a separate Get endpoint, so
			// it is read back from the Create response and the List endpoint. It is
			// computed and sensitive: terraform hides it in plan/apply output.
			"api_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The value of the API key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the API key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the API key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *contextStoreAPIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan contextStoreAPIKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	workspace := plan.Workspace.ValueString()
	contextStoreName := plan.ContextStoreName.ValueString()
	name := plan.Name.ValueString()

	action := fmt.Sprintf("/workspace/%s/contextstore/%s/apikey", workspace, contextStoreName)
	body := map[string]interface{}{
		"name": name,
	}

	response, err := callWithRetry(ctx, createTimeout, func() (map[string]interface{}, error) {
		return client.RoaPost(cmsProductCode, cmsAPIVersion, action, query(), nil, body, true)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateContextStoreAPIKey", contextStoreAPIKeyTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, contextStoreAPIKeyTypeName, action, string(errs.AlibabaCloudSdkGoERROR)).Error(),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s:%s:%s", workspace, contextStoreName, name))
	if v, ok := response["apiKey"]; ok && fmt.Sprint(v) != "" {
		plan.APIKey = types.StringValue(fmt.Sprint(v))
	}

	// Refresh the computed fields (create_time, region_id) and reconcile the key
	// from the authoritative List endpoint. A still-warm API key list is the only
	// read path; the original resource behaved the same way.
	object, err := describeContextStoreAPIKey(ctx, client, workspace, contextStoreName, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading back after Create", contextStoreAPIKeyTypeName),
			errs.WrapErrorf(err, "Creating %s: reading back after Create", contextStoreAPIKeyTypeName).Error(),
		)
		return
	}
	if object != nil {
		fillFromObject(&plan, object, client.RegionId)
	}
	// region_id is provider-local and known at create time; set it even when
	// the List read-back has not yet caught up, so the plan converges without a
	// follow-up refresh.
	if plan.RegionId.IsNull() && client.RegionId != "" {
		plan.RegionId = types.StringValue(client.RegionId)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *contextStoreAPIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state contextStoreAPIKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspace, contextStoreName, name, ok := parseID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing resource id", contextStoreAPIKeyTypeName),
			fmt.Sprintf("invalid resource id %q: expected <workspace>:<context_store_name>:<name>", state.Id.ValueString()),
		)
		return
	}

	client := r.Client()
	object, err := describeContextStoreAPIKey(ctx, client, workspace, contextStoreName, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListContextStoreAPIKeys", contextStoreAPIKeyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling ListContextStoreAPIKeys", contextStoreAPIKeyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	fillFromObject(&state, object, client.RegionId)
	// Keep the configured identity arguments stable on import: the list response
	// omits them on some endpoints, and the parsed id is authoritative there.
	state.Workspace = types.StringValue(workspace)
	state.ContextStoreName = types.StringValue(contextStoreName)
	state.Name = types.StringValue(name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *contextStoreAPIKeyResource) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// No field is updatable: workspace, context_store_name and name are all
	// RequiresReplace, and the Cms API offers no UpdateContextStoreAPIKey.
}

func (r *contextStoreAPIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state contextStoreAPIKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspace, contextStoreName, name, ok := parseID(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing resource id", contextStoreAPIKeyTypeName),
			fmt.Sprintf("invalid resource id %q: expected <workspace>:<context_store_name>:<name>", state.Id.ValueString()),
		)
		return
	}

	client := r.Client()
	action := fmt.Sprintf("/workspace/%s/contextstore/%s/apikey/%s", workspace, contextStoreName, name)

	_, err := callWithRetry(ctx, deleteTimeout, func() (map[string]interface{}, error) {
		return client.RoaDelete(cmsProductCode, cmsAPIVersion, action, query(), nil, nil, true)
	})
	if err != nil && !isNotFound(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteContextStoreAPIKey", contextStoreAPIKeyTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), action, string(errs.AlibabaCloudSdkGoERROR)).Error(),
		)
		return
	}
}

func (r *contextStoreAPIKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeContextStoreAPIKey lists the API keys of a context store and returns the
// one whose name matches, or nil when it is gone. The Cms API exposes no Get.
func describeContextStoreAPIKey(ctx context.Context, client *connectivity.AliyunClient, workspace, contextStoreName, name string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/workspace/%s/contextstore/%s/apikey", workspace, contextStoreName)
	params := query()
	params["maxResults"] = sp("100")

	for {
		response, err := callWithRetry(ctx, readTimeout, func() (map[string]interface{}, error) {
			return client.RoaGet(cmsProductCode, cmsAPIVersion, action, params, nil, nil)
		})
		if err != nil {
			if isNotFound(err) {
				return nil, nil
			}
			return nil, err
		}

		if results, ok := response["results"].([]interface{}); ok {
			for _, item := range results {
				itemRaw, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				if fmt.Sprint(itemRaw["name"]) == name {
					return itemRaw, nil
				}
			}
		}

		nextToken, _ := response["nextToken"].(string)
		if strings.TrimSpace(nextToken) == "" {
			break
		}
		params["nextToken"] = sp(nextToken)
	}
	return nil, nil
}

// fillFromObject maps a list item onto the resource model. The configured
// identity arguments stay authoritative; only the computed attributes are taken
// from the API response here.
func fillFromObject(m *contextStoreAPIKeyModel, object map[string]interface{}, regionId string) {
	if v, ok := object["apiKey"]; ok && fmt.Sprint(v) != "" {
		m.APIKey = types.StringValue(fmt.Sprint(v))
	}
	if v, ok := object["createTime"].(string); ok && v != "" {
		m.CreateTime = types.StringValue(v)
	} else if v, ok := object["create_time"].(string); ok && v != "" {
		m.CreateTime = types.StringValue(v)
	}
	if regionId != "" {
		m.RegionId = types.StringValue(regionId)
	}
}

// parseID splits <workspace>:<context_store_name>:<name> into its three parts.
// SplitN bounds the split so a name that itself contains a colon is preserved.
func parseID(id string) (workspace, contextStoreName, name string, ok bool) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// isNotFound reports whether err means the API key (or its workspace/context
// store) is gone. It covers the Cms-specific not-found codes plus the generic
// 404 path handled by errs.NotFoundError.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	for _, code := range notFoundErrorCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

// callWithRetry retries a Cms ROA call while helper.NeedRetry reports the error
// as transient, using the same incremental backoff the legacy resource used. A
// non-retryable error or an expired deadline returns the last error unchanged.
func callWithRetry(ctx context.Context, timeout time.Duration, call func() (map[string]interface{}, error)) (map[string]interface{}, error) {
	deadline := time.Now().Add(timeout)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		response, err := call()
		if err == nil {
			return response, nil
		}
		if !helper.NeedRetry(err) || time.Now().After(deadline) {
			return response, err
		}
		wait()
	}
}

// query returns a fresh, empty Cms ROA query map. The Roa* helpers take a
// map[string]*string; a non-nil map lets callers add query parameters later.
func query() map[string]*string {
	return make(map[string]*string)
}

// sp is the local string-pointer helper; the root package's StringPointer is not
// importable from a service package without a cycle.
func sp(s string) *string {
	return &s
}
