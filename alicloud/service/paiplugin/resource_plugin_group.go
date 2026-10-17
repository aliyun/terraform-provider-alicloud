package paiplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	tea "github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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

const pluginGroupTypeName = "alicloud_pai_plugin_group"

var (
	_ resource.Resource                = &pluginGroupResource{}
	_ resource.ResourceWithImportState = &pluginGroupResource{}
)

func NewPluginGroupResource() resource.Resource {
	return &pluginGroupResource{}
}

type pluginGroupResource struct {
	fwadapt.ResourceBase
}

// pluginGroupResourceModel maps the PAIPlugin Group object onto the resource
// state. The API carries numbers as JSON numbers, so the int helpers accept
// both float64 (encoding/json default) and json.Number (UseNumber clients).
type pluginGroupResourceModel struct {
	Id             types.String `tfsdk:"id"`
	PluginGroupId  types.String `tfsdk:"plugin_group_id"`
	Name           types.String `tfsdk:"name"`
	Source         types.Int64  `tfsdk:"source"`
	Project        types.String `tfsdk:"project"`
	Algorithm      types.String `tfsdk:"algorithm"`
	Filter         types.String `tfsdk:"filter"`
	Column         types.String `tfsdk:"column"`
	Table          types.String `tfsdk:"table"`
	Text           types.String `tfsdk:"text"`
	Uri            types.String `tfsdk:"uri"`
	Remark         types.String `tfsdk:"remark"`
	PhoneNumber    types.Bool   `tfsdk:"phone_number"`
	Amount         types.Int64  `tfsdk:"amount"`
	CampaignId     types.String `tfsdk:"campaign_id"`
	CreatedTime    types.String `tfsdk:"created_time"`
	UpdatedTime    types.String `tfsdk:"updated_time"`
	Status         types.Int64  `tfsdk:"status"`
	InferenceJobId types.String `tfsdk:"inference_job_id"`
	History        types.String `tfsdk:"history"`
}

func (r *pluginGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = pluginGroupTypeName
}

func (r *pluginGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Plugin Group resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `plugin_group_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_group_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the PAI Plugin Group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source": schema.Int64Attribute{
				Optional: true,
				Description: "The data source of the group. Valid values: `0` to `4`. " +
					"`0`: text (one phone number per line, up to 100); `4`: algorithm output.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"project": schema.StringAttribute{
				Optional:    true,
				Description: "The project of the group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"algorithm": schema.StringAttribute{
				Optional:    true,
				Description: "The algorithm of the group, required when `source` is `4`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"filter": schema.StringAttribute{
				Optional:    true,
				Description: "The filter of the group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"column": schema.StringAttribute{
				Optional:    true,
				Description: "The column of the group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"table": schema.StringAttribute{
				Optional:    true,
				Description: "The table of the group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"text": schema.StringAttribute{
				Optional:    true,
				Description: "The text content of the group, one phone number per line and up to 100 when `source` is `0`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uri": schema.StringAttribute{
				Optional:    true,
				Description: "The URI of the group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Description: "The remark of the group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"phone_number": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the group contains phone numbers.",
			},
			"amount": schema.Int64Attribute{
				Computed:    true,
				Description: "The amount of the group.",
			},
			"campaign_id": schema.StringAttribute{
				Computed:    true,
				Description: "The campaign ID of the group.",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the group was created.",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the group was last updated.",
			},
			"status": schema.Int64Attribute{
				Computed:    true,
				Description: "The status of the group.",
			},
			"inference_job_id": schema.StringAttribute{
				Computed:    true,
				Description: "The inference job ID of the group.",
			},
			"history": schema.StringAttribute{
				Computed:    true,
				Description: "The history of the group.",
			},
		},
	}
}

func (r *pluginGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pluginGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"Name": plan.Name.ValueString(),
	}
	if !plan.Source.IsNull() {
		request["Source"] = plan.Source.ValueInt64()
	}
	if !plan.Project.IsNull() {
		request["Project"] = plan.Project.ValueString()
	}
	if !plan.Algorithm.IsNull() {
		request["Algorithm"] = plan.Algorithm.ValueString()
	}
	if !plan.Filter.IsNull() {
		request["Filter"] = plan.Filter.ValueString()
	}
	if !plan.Column.IsNull() {
		request["Column"] = plan.Column.ValueString()
	}
	if !plan.Table.IsNull() {
		request["Table"] = plan.Table.ValueString()
	}
	if !plan.Text.IsNull() {
		request["Text"] = plan.Text.ValueString()
	}
	if !plan.Uri.IsNull() {
		request["Uri"] = plan.Uri.ValueString()
	}
	if !plan.Remark.IsNull() {
		request["Remark"] = plan.Remark.ValueString()
	}

	// PAIPlugin 2022-01-12 is ROA/RESTful. The official Go SDK
	// (alibabacloud-go/paiplugin-20220112) sends Action=CreateGroup,
	// Version=2022-01-12 over POST /api/v2/groups with a JSON body. The
	// generic RoaPost passes apiName="" to DoRequest, which omits the
	// x-acs-action header the backend routes on; RoaPostWithApiName drives
	// DoRequestWithAction so the operation Action is sent explicitly, matching
	// the SDK. The same applies to GetGroup/DeleteGroup below.
	response, err := r.Client().RoaPostWithApiName("PAIPlugin", "2022-01-12", "CreateGroup", "/api/v2/groups", nil, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateGroup", pluginGroupTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateGroup", pluginGroupTypeName).Error(),
		)
		return
	}

	data, _ := response["Data"].(map[string]interface{})
	groupId := toStr(data["Id"])
	if groupId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateGroup", pluginGroupTypeName),
			fmt.Sprintf("The response carried no Data.Id: %v", response),
		)
		return
	}

	// F1: persist a Framework state carrying the freshly assigned id/plugin_group_id
	// and the original plan configuration BEFORE the follow-up GetGroup read. A
	// transient GetGroup failure (404-not-yet-visible / 500 / malformed body) then
	// leaves the id in state instead of dropping it, so the resource cannot become
	// an orphan or be silently rebuilt. The read below refreshes the computed
	// attributes when it succeeds; on failure the persisted state is retained.
	resp.Diagnostics.Append(resp.State.Set(ctx, pluginGroupStateFromCreate(plan, groupId))...)

	object, err := FindPluginGroup(r.Client(), groupId)
	if err != nil {
		// Routing/auth/5xx errors are surfaced but do not wipe the id already in state.
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetGroup", pluginGroupTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetGroup", pluginGroupTypeName).Error(),
		)
		return
	}
	if object == nil {
		// Genuine not-found right after creation: keep the persisted id state and
		// surface the inconsistency rather than removing the resource.
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetGroup", pluginGroupTypeName),
			fmt.Sprintf("The group %s is not visible via GetGroup right after creation", groupId),
		)
		return
	}
	// F3 (AC-R3): a 200 GetGroup body whose Data.Id is missing/invalid is
	// malformed — surface it instead of silently leaving the refreshed state
	// unset. The create-state persisted above is retained (not nulled) because
	// no Set/RemoveResource runs on this branch.
	if state, ok := pluginGroupState(object); ok {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: parsing GetGroup", pluginGroupTypeName),
			fmt.Sprintf("The GetGroup response Data is malformed (missing or invalid Id): %v", object),
		)
	}
}

func (r *pluginGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pluginGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := FindPluginGroup(r.Client(), state.Id.ValueString())
	if err != nil {
		// F2: only a genuine product not-found removes the resource; routing/auth
		// (InvalidApi.NotFound/InvalidAction), 403/500 and malformed responses keep
		// the existing state (returning without RemoveResource/Set preserves it).
		if isPluginGroupNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetGroup", pluginGroupTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetGroup", pluginGroupTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	// F3 (AC-R3): a 200 GetGroup body whose Data.Id is missing/invalid is
	// malformed — surface a diagnostic instead of silently doing nothing. The
	// prior state is retained (no Set/RemoveResource runs on this branch).
	if state, ok := pluginGroupState(object); ok {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	} else {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing GetGroup", pluginGroupTypeName),
			fmt.Sprintf("The GetGroup response Data is malformed (missing or invalid Id): %v", object),
		)
	}
}

func (r *pluginGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// PAIPlugin Group exposes no in-place Update API; every configurable attribute
	// is ForceNew, so the engine routes any config change to replacement and this
	// method is unreachable for plan diffs. It exists to satisfy the
	// resource.Resource interface (which mandates Update). Adopt the planned
	// state verbatim so state stays coherent if it is ever invoked.
	var plan pluginGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *pluginGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pluginGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Client().RoaDeleteWithApiName("PAIPlugin", "2022-01-12", "DeleteGroup", fmt.Sprintf("/api/v2/groups/%s", state.Id.ValueString()), nil, nil, nil, true)
	if err != nil && !isPluginGroupNotExist(err) {
		// F2: a genuine not-found is tolerated (idempotent delete); routing/auth/5xx
		// errors are surfaced and the resource is left in place.
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteGroup", pluginGroupTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteGroup", pluginGroupTypeName).Error(),
		)
		return
	}
}

func (r *pluginGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// FindPluginGroup queries GetGroup by group ID and returns the Data object from
// the response, or nil when the group no longer exists. Exported so the
// acceptance test's CheckDestroy reuses the same ROA lookup + not-found contract
// as the resource Read/Delete instead of replicating the call.
func FindPluginGroup(client *connectivity.AliyunClient, groupId string) (map[string]interface{}, error) {
	response, err := client.RoaGetWithApiName("PAIPlugin", "2022-01-12", "GetGroup", fmt.Sprintf("/api/v2/groups/%s", groupId), nil, nil, nil)
	if err != nil {
		if isPluginGroupNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	data, _ := response["Data"].(map[string]interface{})
	if data == nil {
		return nil, fmt.Errorf("the GetGroup response carried no Data object: %v", response)
	}
	// F1 (AC-R1): GetGroup must answer for the requested group. A response
	// whose Data.Id does not match the requested groupId is malformed — a
	// different resource's body reaching this call path (cross-resource
	// takeover) — so it is a hard error (not not-found): the caller surfaces a
	// diagnostic and preserves state rather than adopting a foreign group or
	// removing the resource. toStr yields "" for a missing/non-scalar Id, which
	// never equals a non-empty groupId, so a missing/invalid Id is covered too.
	if toStr(data["Id"]) != groupId {
		return nil, fmt.Errorf("the GetGroup response Data.Id %q does not match the requested groupId %q: %v", toStr(data["Id"]), groupId, data)
	}
	return data, nil
}

// pluginGroupState maps the Data object of a GetGroup response onto the
// resource model. Empty strings become null: dates and optional fields the API
// omits read back as absent rather than "". It reports ok=false when the Data
// object is nil or carries no valid group Id, so the caller can preserve the
// prior state on a malformed response instead of nulling it.
func pluginGroupState(data map[string]interface{}) (pluginGroupResourceModel, bool) {
	if data == nil {
		return pluginGroupResourceModel{}, false
	}
	id := toStr(data["Id"])
	if id == "" {
		return pluginGroupResourceModel{}, false
	}
	str := func(key string) types.String {
		if v, ok := data[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	boolVal := func(key string) types.Bool {
		if v, ok := data[key].(bool); ok {
			return types.BoolValue(v)
		}
		return types.BoolNull()
	}
	intVal := func(key string) types.Int64 {
		switch v := data[key].(type) {
		case float64:
			return types.Int64Value(int64(v))
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.Int64Value(n)
			}
			return types.Int64Null()
		default:
			return types.Int64Null()
		}
	}
	idv := types.StringValue(id)
	return pluginGroupResourceModel{
		Id:             idv,
		PluginGroupId:  idv,
		Name:           str("Name"),
		Source:         intVal("Source"),
		Project:        str("Project"),
		Algorithm:      str("Algorithm"),
		Filter:         str("Filter"),
		Column:         str("Column"),
		Table:          str("Table"),
		Text:           str("Text"),
		Uri:            str("Uri"),
		Remark:         str("Remark"),
		PhoneNumber:    boolVal("PhoneNumber"),
		Amount:         intVal("Amount"),
		CampaignId:     str("CampaignId"),
		CreatedTime:    str("CreatedTime"),
		UpdatedTime:    str("UpdatedTime"),
		Status:         intVal("Status"),
		InferenceJobId: str("InferenceJobId"),
		History:        str("History"),
	}, true
}

// pluginGroupStateFromCreate persists the freshly assigned id and the original
// plan configuration before the follow-up GetGroup read (see F1). It takes the
// plan by value so the caller's plan is never mutated.
func pluginGroupStateFromCreate(plan pluginGroupResourceModel, groupId string) pluginGroupResourceModel {
	m := plan
	m.Id = types.StringValue(groupId)
	m.PluginGroupId = types.StringValue(groupId)
	return m
}

// isPluginGroupNotExist reports whether err is the genuine PAIPlugin
// group-not-found contract: an HTTP 404 whose code is an EXPLICIT not-found
// marker (allow-list). The ROA SDK surfaces HTTP 404 as a *tea.SDKError; a
// blanket "any 404 is not-found" rule would also catch InvalidApi.NotFound /
// InvalidAction (operation not routed to the backend) and auth-layer 404s
// (Forbidden / InvalidAccessKeyId / SignatureDoesNotMatch), silently
// RemoveResource-ing or ignoring a failed Delete and orphaning/rebuilding live
// resources. So an empty code, an unknown code and a routing/auth code are all
// hard errors that preserve state; only an explicit not-found code removes.
// errs wraps ROA errors in *ComplexError, so unwrap .Cause to reach the SDK error.
func isPluginGroupNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*errs.ComplexError); ok {
		return isPluginGroupNotExist(e.Cause)
	}
	if e, ok := err.(*tea.SDKError); ok {
		statusCode := 0
		if e.StatusCode != nil {
			statusCode = *e.StatusCode
		}
		if statusCode != 404 {
			return false
		}
		code := ""
		if e.Code != nil {
			code = strings.ToLower(*e.Code)
		}
		// F2 (AC-R2): only an explicit product not-found code removes the
		// resource or tolerates a failed Delete. An empty code, an unknown code
		// and routing/auth-layer 404s (InvalidApi.NotFound / InvalidAction /
		// auth) are hard errors that preserve state — never silently
		// RemoveResource, which would orphan or rebuild a live resource.
		if code == "" || isPluginGroupRoutingOrAuthCode(code) {
			return false
		}
		return isPluginGroupNotFoundCode(code)
	}
	// Non-SDK errors are not the PAIPlugin not-found contract; preserve state.
	return false
}

// isPluginGroupNotFoundCode reports whether a lower-cased, non-empty SDK error
// code that survived the routing/auth deny-list is an explicit PAIPlugin
// group-not-found contract. Unknown codes are NOT treated as not-found: a
// RemoveResource on an unrecognized 404 would orphan a live resource, so the
// caller keeps the state and surfaces the error instead. The routing prefixes
// (invalidapi/invalidaction/auth) are excluded before this is reached, so
// InvalidApi.NotFound — which carries a "notfound" suffix — is never matched.
func isPluginGroupNotFoundCode(code string) bool {
	switch code {
	case "notfound", "resourcenotfound", "groupnotfound", "nosuchgroup",
		"notfoundexception", "entitynotexist", "nosuchentity":
		return true
	}
	// Accept "<scope>.notfound" / "<scope>notfound" forms emitted for a genuine
	// missing entity; routing/auth forms were already excluded upstream.
	if strings.HasSuffix(code, ".notfound") || strings.HasSuffix(code, "notfound") {
		return true
	}
	return false
}

// isPluginGroupRoutingOrAuthCode reports whether a lower-cased SDK error code
// indicates a routing/auth-layer 404 rather than a genuine missing group.
func isPluginGroupRoutingOrAuthCode(code string) bool {
	prefixes := []string{
		"invalidapi",            // InvalidApi.NotFound / InvalidApi.* — service/routing
		"invalidaction",         // InvalidAction — operation not routed
		"invalidaccesskey",      // InvalidAccessKeyId.* — auth
		"missingaccesskey",      // MissingAccessKeyId.* — auth
		"signaturedoesnotmatch", // auth
		"forbidden",             // Forbidden / Forbidden.* — auth
		"accessdenied",          // AccessDenied — auth
		"unauthorized",          // Unauthorized — auth
	}
	for _, p := range prefixes {
		if strings.HasPrefix(code, p) {
			return true
		}
	}
	return false
}

// toStr renders a scalar group Id from the response as a string, covering the
// JSON number forms the SDK may deliver. Only string and integer scalars are
// accepted; bool/map/slice/nil and non-integer floats yield "" so a malformed
// body cannot synthesize a pseudo-id.
func toStr(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		if val == float64(int64(val)) {
			return strconv.FormatInt(int64(val), 10)
		}
		return "" // non-integer float is not a valid group id
	case json.Number:
		if _, err := val.Int64(); err == nil {
			return val.String()
		}
		return ""
	default:
		return "" // nil/bool/map/slice/other cannot be a group id
	}
}
