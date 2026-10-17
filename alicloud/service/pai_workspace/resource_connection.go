package pai_workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const connectionTypeName = "alicloud_pai_workspace_connection"

var (
	_ resource.Resource                = &connectionResource{}
	_ resource.ResourceWithImportState = &connectionResource{}
)

func NewConnectionResource() resource.Resource {
	return &connectionResource{}
}

type connectionResource struct {
	fwadapt.ResourceBase
}

// connectionResourceModel is the state shape of the connection resource. The
// id and connection_id fields both carry ConnectionId, mirroring the ims user
// resource where id mirrors the primary identifier.
type connectionResourceModel struct {
	Id              types.String `tfsdk:"id"`
	ConnectionId    types.String `tfsdk:"connection_id"`
	ConnectionName  types.String `tfsdk:"connection_name"`
	ConnectionType  types.String `tfsdk:"connection_type"`
	WorkspaceId     types.String `tfsdk:"workspace_id"`
	Description     types.String `tfsdk:"description"`
	Accessibility   types.String `tfsdk:"accessibility"`
	Configs         types.Map    `tfsdk:"configs"`
	Secrets         types.Map    `tfsdk:"secrets"`
	Models          types.List   `tfsdk:"models"`
	ResourceMeta    types.Object `tfsdk:"resource_meta"`
	Creator         types.String `tfsdk:"creator"`
	GmtCreateTime   types.String `tfsdk:"gmt_create_time"`
	GmtModifiedTime types.String `tfsdk:"gmt_modified_time"`
}

// connectionModelModel maps a single Models element. It is used to unmarshal
// plan list elements via ElementsAs. The tfsdk tags match the nested object
// attribute names declared in the Schema.
type connectionModelModel struct {
	Model                 types.String `tfsdk:"model"`
	DisplayName           types.String `tfsdk:"display_name"`
	ModelType             types.String `tfsdk:"model_type"`
	ToolCall              types.Bool   `tfsdk:"tool_call"`
	SupportReasoning      types.Bool   `tfsdk:"support_reasoning"`
	SupportResponseSchema types.Bool   `tfsdk:"support_response_schema"`
	SupportVision         types.Bool   `tfsdk:"support_vision"`
	MaxModelLength        types.Int64  `tfsdk:"max_model_length"`
}

func (r *connectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = connectionTypeName
}

func (r *connectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Workspace Connection resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `connection_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_id": schema.StringAttribute{
				Computed:    true,
				Description: "The first ID of the resource, the connection ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connection_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the connection, such as `MilvusConnection`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the PAI workspace to which the connection belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the connection.",
			},
			"accessibility": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The accessibility of the connection, such as `PRIVATE` or `PUBLIC`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"configs": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The configuration info of the connection, a map of string key-value pairs.",
			},
			"secrets": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Sensitive:   true,
				Description: "The secrets of the connection, a map of string key-value pairs. The API masks secret values on read-back, so the provider preserves the configured values in state.",
			},
			"models": schema.ListAttribute{
				Optional:    true,
				Description: "The model configuration list of the connection. Each element is an object with model, display_name, model_type, tool_call, support_reasoning, support_response_schema, support_vision and max_model_length fields.",
				ElementType: types.ObjectType{AttrTypes: modelAttrTypes()},
			},
			"resource_meta": schema.ObjectAttribute{
				Optional:       true,
				Computed:       true,
				Description:    "The resource meta info of the connection. Set only at create time.",
				AttributeTypes: resourceMetaAttrTypes(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
					objectplanmodifier.RequiresReplace(),
				},
			},
			"creator": schema.StringAttribute{
				Computed:    true,
				Description: "The creator of the connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The create time of the connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The modified time of the connection. The API bumps this value on every UpdateConnection; the update path keeps the prior value in state so apply does not flag an inconsistent result, and the fresh value is picked up on the next refresh.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *connectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan connectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"ConnectionName": plan.ConnectionName.ValueString(),
		"ConnectionType": plan.ConnectionType.ValueString(),
		"WorkspaceId":    plan.WorkspaceId.ValueString(),
		"Configs":        mapStringToBody(plan.Configs),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.Accessibility.IsNull() && !plan.Accessibility.IsUnknown() {
		request["Accessibility"] = plan.Accessibility.ValueString()
	}
	if !plan.Secrets.IsNull() && !plan.Secrets.IsUnknown() {
		request["Secrets"] = mapStringToBody(plan.Secrets)
	}
	if !plan.ResourceMeta.IsNull() && !plan.ResourceMeta.IsUnknown() {
		request["ResourceMeta"] = resourceMetaToBody(plan.ResourceMeta)
	}
	if !plan.Models.IsNull() && !plan.Models.IsUnknown() {
		modelsBody, d := modelsToCreateBody(ctx, plan.Models)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(modelsBody) > 0 {
			request["Models"] = modelsBody
		}
	}

	action := "/api/v1/connections"
	response, err := r.Client().RoaPost("AIWorkSpace", "2021-02-04", action, map[string]*string{}, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateConnection", connectionTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateConnection", connectionTypeName).Error(),
		)
		return
	}

	id, _ := response["ConnectionId"].(string)
	if id == "" {
		// The Roa* helper unwraps response["body"]; ConnectionId sits at the top level.
		if v := response["ConnectionId"]; v != nil {
			id = fmt.Sprint(v)
		}
	}
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateConnection", connectionTypeName),
			fmt.Sprintf("The CreateConnection response carried no ConnectionId: %v", response),
		)
		return
	}

	object, err := getConnection(ctx, r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetConnection", connectionTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetConnection after create", connectionTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateConnection", connectionTypeName),
			fmt.Sprintf("The connection %s is not visible via GetConnection right after creation", id),
		)
		return
	}

	state, d := connectionState(ctx, object)
	resp.Diagnostics.Append(d...)
	// GetConnection masks secret values and omits the create-only capability
	// fields, so the configured values are kept instead of being overwritten by
	// masked/null read-back values.
	state = preserveWriteOnly(state, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *connectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state connectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getConnection(ctx, r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetConnection", connectionTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetConnection", connectionTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState, d := connectionState(ctx, object)
	resp.Diagnostics.Append(d...)
	newState = preserveWriteOnly(newState, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *connectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state connectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	request := map[string]interface{}{}
	// Optional and updatable fields are sent only when changed, so a null plan
	// value clears the attribute (sending "") instead of leaving the old value
	// behind. This mirrors the ims user resource update behavior.
	if !plan.Description.Equal(state.Description) {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.Secrets.Equal(state.Secrets) {
		request["Secrets"] = mapStringToBody(plan.Secrets)
	}
	if !plan.Configs.Equal(state.Configs) {
		request["Configs"] = mapStringToBody(plan.Configs)
	}
	if !plan.Models.Equal(state.Models) {
		modelsBody, d := modelsToUpdateBody(ctx, plan.Models)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		request["Models"] = modelsBody
	}

	action := fmt.Sprintf("/api/v1/connections/%s", id)
	if len(request) > 0 {
		if _, err := r.Client().RoaPut("AIWorkSpace", "2021-02-04", action, map[string]*string{}, nil, request, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateConnection", connectionTypeName),
				errs.WrapErrorf(err, "Updating %s: calling UpdateConnection", connectionTypeName).Error(),
			)
			return
		}
	}

	object, err := getConnection(ctx, r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetConnection", connectionTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetConnection after update", connectionTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState, d := connectionState(ctx, object)
	resp.Diagnostics.Append(d...)
	newState = preserveWriteOnly(newState, plan)
	// GmtModifiedTime is refreshed by the API on every UpdateConnection, but
	// the plan (via UseStateForUnknown) keeps the prior value. Align the new
	// state with the plan so apply does not flag an inconsistent result after
	// apply; the fresh value is picked up on the next refresh/Read.
	newState.GmtModifiedTime = plan.GmtModifiedTime
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *connectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state connectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	action := fmt.Sprintf("/api/v1/connections/%s", id)
	_, err := r.Client().RoaDelete("AIWorkSpace", "2021-02-04", action, map[string]*string{}, nil, nil, true)
	if err != nil && !isConnectionNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteConnection", connectionTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteConnection", connectionTypeName).Error(),
		)
		return
	}
}

func (r *connectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getConnection queries GetConnection by connection ID and returns the
// connection object from the response, or nil when the connection no longer
// exists. The Roa* helper unwraps response["body"], so the returned map is the
// flat connection object.
func getConnection(ctx context.Context, client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/api/v1/connections/%s", id)
	response, err := client.RoaGet("AIWorkSpace", "2021-02-04", action, map[string]*string{}, nil, nil)
	if err != nil {
		if isConnectionNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// connectionState maps a GetConnection response object onto the resource model.
// Empty strings become null: dates and optional fields the API omits read back
// as absent rather than "".
func connectionState(ctx context.Context, obj map[string]interface{}) (connectionResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	str := func(key string) types.String {
		if v, ok := obj[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}

	id := str("ConnectionId")
	m := connectionResourceModel{
		Id:              id,
		ConnectionId:    id,
		ConnectionName:  str("ConnectionName"),
		ConnectionType:  str("ConnectionType"),
		WorkspaceId:     str("WorkspaceId"),
		Description:     str("Description"),
		Accessibility:   str("Accessibility"),
		Creator:         str("Creator"),
		GmtCreateTime:   str("GmtCreateTime"),
		GmtModifiedTime: str("GmtModifiedTime"),
	}

	var d diag.Diagnostics
	m.Configs, d = mapStringFromResponse(asMap(obj["Configs"]))
	diags = append(diags, d...)
	m.Secrets, d = mapStringFromResponse(asMap(obj["Secrets"]))
	diags = append(diags, d...)
	m.ResourceMeta, d = buildResourceMetaObject(asMap(obj["ResourceMeta"]))
	diags = append(diags, d...)

	m.Models = types.ListNull(types.ObjectType{AttrTypes: modelAttrTypes()})
	if raw, ok := obj["Models"].([]interface{}); ok && len(raw) > 0 {
		objs := make([]attr.Value, 0, len(raw))
		for _, e := range raw {
			mm, ok := e.(map[string]interface{})
			if !ok {
				continue
			}
			o, dd := buildModelObject(mm)
			diags = append(diags, dd...)
			objs = append(objs, o)
		}
		if len(objs) > 0 {
			list, dd := types.ListValue(types.ObjectType{AttrTypes: modelAttrTypes()}, objs)
			diags = append(diags, dd...)
			m.Models = list
		}
	}
	return m, diags
}

// createOnlyModelFields are Models element fields that GetConnection does not
// reliably echo on read-back: the four CreateConnection-only capability
// fields. The API omits MaxModelLength and returns stale false defaults for
// the boolean capability flags after UpdateConnection (the API ignores them on
// update), so the prior configured value is kept instead of being overwritten.
var createOnlyModelFields = map[string]bool{
	"support_reasoning":       true,
	"support_response_schema": true,
	"support_vision":          true,
	"max_model_length":        true,
}

// preserveWriteOnly retains values the API does not echo on read-back. Secret
// values are masked by the API and the create-only Models capability fields
// are omitted from GetConnection, so refresh would otherwise null out the
// configured values. When the API returns a real value it is kept; only null
// API values fall back to the prior configured value.
func preserveWriteOnly(dst, prior connectionResourceModel) connectionResourceModel {
	if !prior.Secrets.IsNull() {
		dst.Secrets = prior.Secrets
	}
	dst.Models = preserveCreateOnlyModels(dst.Models, prior.Models)
	return dst
}

// preserveCreateOnlyModels merges prior create-only capability field values
// into the API-derived model list for elements the API left null. Elements
// are matched by index; mismatched lengths leave the API list unchanged.
func preserveCreateOnlyModels(api, prior types.List) types.List {
	if api.IsNull() || prior.IsNull() {
		return api
	}
	apiElems := api.Elements()
	priorElems := prior.Elements()
	if len(apiElems) != len(priorElems) {
		return api
	}
	out := make([]attr.Value, len(apiElems))
	for i, ae := range apiElems {
		aObj, ok := ae.(types.Object)
		if !ok || aObj.IsNull() {
			out[i] = ae
			continue
		}
		pObj, ok := priorElems[i].(types.Object)
		if !ok || pObj.IsNull() {
			out[i] = ae
			continue
		}
		out[i] = mergeCreateOnlyModelObject(aObj, pObj)
	}
	if l, diags := types.ListValue(types.ObjectType{AttrTypes: modelAttrTypes()}, out); !diags.HasError() {
		return l
	}
	return api
}

// mergeCreateOnlyModelObject keeps the API value for settable fields and, for
// the create-only capability fields, retains the prior configured value
// whenever it was set. GetConnection does not reliably echo these fields: a
// non-null false/zero read-back (common after UpdateConnection) is a stale
// default rather than a deliberate change, since the API ignores them on
// update. Matching only null API values left the false read-back disagreeing
// with the configured true, triggering an "inconsistent result after apply".
func mergeCreateOnlyModelObject(api, prior types.Object) types.Object {
	apiAttrs := api.Attributes()
	priorAttrs := prior.Attributes()
	vals := make(map[string]attr.Value, len(modelAttrTypes()))
	for k := range modelAttrTypes() {
		v := apiAttrs[k]
		if createOnlyModelFields[k] {
			if pv, ok := priorAttrs[k]; ok && !pv.IsNull() {
				v = pv
			}
		}
		vals[k] = v
	}
	obj, _ := types.ObjectValue(modelAttrTypes(), vals)
	return obj
}

// modelAttrTypes returns the attribute type map for a Models element object.
func modelAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"model":                   types.StringType,
		"display_name":            types.StringType,
		"model_type":              types.StringType,
		"tool_call":               types.BoolType,
		"support_reasoning":       types.BoolType,
		"support_response_schema": types.BoolType,
		"support_vision":          types.BoolType,
		"max_model_length":        types.Int64Type,
	}
}

func resourceMetaAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"instance_name": types.StringType,
		"instance_id":   types.StringType,
		"extra":         types.StringType,
	}
}

// buildModelObject builds a framework Object for one Models element from the
// raw API response map. JSON booleans arrive as bool; integer fields arrive as
// float64 because the response is JSON-decoded into map[string]interface{}.
func buildModelObject(m map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	str := func(k string) types.String {
		if v, ok := m[k].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	b := func(k string) types.Bool {
		if v, ok := m[k].(bool); ok {
			return types.BoolValue(v)
		}
		return types.BoolNull()
	}
	i := func(k string) types.Int64 {
		switch v := m[k].(type) {
		case float64:
			return types.Int64Value(int64(v))
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		}
		return types.Int64Null()
	}
	vals := map[string]attr.Value{
		"model":                   str("Model"),
		"display_name":            str("DisplayName"),
		"model_type":              str("ModelType"),
		"tool_call":               b("ToolCall"),
		"support_reasoning":       b("SupportReasoning"),
		"support_response_schema": b("SupportResponseSchema"),
		"support_vision":          b("SupportVision"),
		"max_model_length":        i("MaxModelLength"),
	}
	obj, d := types.ObjectValue(modelAttrTypes(), vals)
	return obj, append(diags, d...)
}

func buildResourceMetaObject(m map[string]interface{}) (types.Object, diag.Diagnostics) {
	if m == nil {
		return types.ObjectNull(resourceMetaAttrTypes()), nil
	}
	str := func(k string) types.String {
		if v, ok := m[k].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	vals := map[string]attr.Value{
		"instance_name": str("InstanceName"),
		"instance_id":   str("InstanceId"),
		"extra":         str("Extra"),
	}
	obj, d := types.ObjectValue(resourceMetaAttrTypes(), vals)
	return obj, d
}

// mapStringToBody converts a framework string map attribute into the
// map[string]interface{} shape sent in request bodies. Null/unknown maps
// yield an empty map so the request field is always present for Configs.
func mapStringToBody(m types.Map) map[string]interface{} {
	out := map[string]interface{}{}
	if m.IsNull() || m.IsUnknown() {
		return out
	}
	for k, v := range m.Elements() {
		if s, ok := v.(types.String); ok && !s.IsNull() {
			out[k] = s.ValueString()
		}
	}
	return out
}

// mapStringFromResponse builds a framework string map from the raw API response
// object. Non-string values are stringified to keep the map[string]string
// contract the SDKv2 line relied on.
func mapStringFromResponse(m map[string]interface{}) (types.Map, diag.Diagnostics) {
	if m == nil {
		return types.MapNull(types.StringType), nil
	}
	vals := make(map[string]attr.Value, len(m))
	for k, v := range m {
		if s, ok := v.(string); ok {
			vals[k] = types.StringValue(s)
		} else {
			vals[k] = types.StringValue(fmt.Sprint(v))
		}
	}
	return types.MapValue(types.StringType, vals)
}

// resourceMetaToBody converts a framework resource_meta Object into the request
// body shape, omitting null fields.
func resourceMetaToBody(o types.Object) map[string]interface{} {
	out := map[string]interface{}{}
	if o.IsNull() || o.IsUnknown() {
		return out
	}
	attrs := o.Attributes()
	strVal := func(k string) (string, bool) {
		if v, ok := attrs[k].(types.String); ok && !v.IsNull() && !v.IsUnknown() {
			return v.ValueString(), true
		}
		return "", false
	}
	if v, ok := strVal("instance_name"); ok {
		out["InstanceName"] = v
	}
	if v, ok := strVal("instance_id"); ok {
		out["InstanceId"] = v
	}
	if v, ok := strVal("extra"); ok {
		out["Extra"] = v
	}
	return out
}

// modelsToCreateBody builds the Models array sent to CreateConnection. All
// eight element fields are eligible: the four settable fields plus the four
// create-only capability fields governed by the CloudSpec mapping.
func modelsToCreateBody(ctx context.Context, models types.List) ([]map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	if models.IsNull() || models.IsUnknown() {
		return nil, diags
	}
	var elems []connectionModelModel
	diags.Append(models.ElementsAs(ctx, &elems, false)...)
	if diags.HasError() {
		return nil, diags
	}
	out := make([]map[string]interface{}, 0, len(elems))
	for _, m := range elems {
		item := map[string]interface{}{}
		if !m.Model.IsNull() && !m.Model.IsUnknown() {
			item["Model"] = m.Model.ValueString()
		}
		if !m.DisplayName.IsNull() && !m.DisplayName.IsUnknown() {
			item["DisplayName"] = m.DisplayName.ValueString()
		}
		if !m.ModelType.IsNull() && !m.ModelType.IsUnknown() {
			item["ModelType"] = m.ModelType.ValueString()
		}
		if !m.ToolCall.IsNull() && !m.ToolCall.IsUnknown() {
			item["ToolCall"] = m.ToolCall.ValueBool()
		}
		if !m.SupportReasoning.IsNull() && !m.SupportReasoning.IsUnknown() {
			item["SupportReasoning"] = m.SupportReasoning.ValueBool()
		}
		if !m.SupportResponseSchema.IsNull() && !m.SupportResponseSchema.IsUnknown() {
			item["SupportResponseSchema"] = m.SupportResponseSchema.ValueBool()
		}
		if !m.SupportVision.IsNull() && !m.SupportVision.IsUnknown() {
			item["SupportVision"] = m.SupportVision.ValueBool()
		}
		if !m.MaxModelLength.IsNull() && !m.MaxModelLength.IsUnknown() {
			item["MaxModelLength"] = m.MaxModelLength.ValueInt64()
		}
		out = append(out, item)
	}
	return out, diags
}

// modelsToUpdateBody builds the Models array sent to UpdateConnection. The
// update mapping carries only the four settable fields; the create-only
// capability fields are omitted because the API ignores them on update.
func modelsToUpdateBody(ctx context.Context, models types.List) ([]map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	if models.IsNull() || models.IsUnknown() {
		return []map[string]interface{}{}, diags
	}
	var elems []connectionModelModel
	diags.Append(models.ElementsAs(ctx, &elems, false)...)
	if diags.HasError() {
		return nil, diags
	}
	out := make([]map[string]interface{}, 0, len(elems))
	for _, m := range elems {
		item := map[string]interface{}{}
		if !m.Model.IsNull() && !m.Model.IsUnknown() {
			item["Model"] = m.Model.ValueString()
		}
		if !m.DisplayName.IsNull() && !m.DisplayName.IsUnknown() {
			item["DisplayName"] = m.DisplayName.ValueString()
		}
		if !m.ModelType.IsNull() && !m.ModelType.IsUnknown() {
			item["ModelType"] = m.ModelType.ValueString()
		}
		if !m.ToolCall.IsNull() && !m.ToolCall.IsUnknown() {
			item["ToolCall"] = m.ToolCall.ValueBool()
		}
		out = append(out, item)
	}
	return out, diags
}

// asMap coerces an interface{} to a map[string]interface{}, returning nil for
// non-map values.
func asMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return nil
}

// isConnectionNotExist reports whether err is the connection not-found error.
// The CloudSpec model declares RequestFailedErrorProblem as the not-exist error
// code (notExistCheckType: checkErrorCode) for both GetConnection and
// DeleteConnection. This mirrors errs.IsExpectedErrors' tea.SDKError branch:
// code equality, the historical prefix check, and the *e.Data fallback.
func isConnectionNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok {
		if e.Code != nil {
			if *e.Code == "RequestFailedErrorProblem" ||
				strings.HasPrefix("RequestFailedErrorProblem", *e.Code) {
				return true
			}
		}
		if e.Data != nil && strings.Contains(*e.Data, "RequestFailedErrorProblem") {
			return true
		}
		return false
	}
	return strings.Contains(err.Error(), "RequestFailedErrorProblem")
}
