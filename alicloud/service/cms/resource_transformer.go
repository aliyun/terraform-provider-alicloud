package cms

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &transformerResource{}
	_ resource.ResourceWithImportState = &transformerResource{}
)

// NewTransformerResource is the Framework factory registered for
// alicloud_cms_transformer.
func NewTransformerResource() resource.Resource {
	return &transformerResource{}
}

// transformerResource implements the alicloud_cms_transformer resource against
// the CMS RESTful API (2024-03-30). It is the Framework-native replacement of
// the SDKv2 prototype that lived under the provider root.
type transformerResource struct {
	fwadapt.ResourceBase
}

// ---- tfsdk models ---------------------------------------------------------

// transformerResourceModel mirrors the resource schema. The nested Actions and
// FilterSetting fields are types.List of opaque object values; the framework
// reflection in ElementsAs/ObjectValueFrom handles (de)serialization against
// the NestedBlockObject definitions declared in Schema.
type transformerResourceModel struct {
	Id              types.String `tfsdk:"id"`
	TransformerId   types.String `tfsdk:"transformer_id"`
	TransformerName types.String `tfsdk:"transformer_name"`
	Workspace       types.String `tfsdk:"workspace"`
	Description     types.String `tfsdk:"description"`
	QuitAfterMatch  types.Bool   `tfsdk:"quit_after_match"`
	SortId          types.Int64  `tfsdk:"sort_id"`
	Enable          types.Bool   `tfsdk:"enable"`
	Actions         types.List   `tfsdk:"actions"`
	FilterSetting   types.List   `tfsdk:"filter_setting"`
	CreateTime      types.String `tfsdk:"create_time"`
	UpdateTime      types.String `tfsdk:"update_time"`
	UserId          types.String `tfsdk:"user_id"`
	RegionId        types.String `tfsdk:"region_id"`
}

// actionModel mirrors the actions NestedBlockObject. FilterSetting is a
// list capped at one element by the schema.
type actionModel struct {
	FilterSetting types.List   `tfsdk:"filter_setting"`
	LabelKey      types.String `tfsdk:"label_key"`
	Mapping       types.Map    `tfsdk:"mapping"`
	RegExp        types.String `tfsdk:"reg_exp"`
	Source        types.String `tfsdk:"source"`
	Target        types.String `tfsdk:"target"`
	Type          types.String `tfsdk:"type"`
	Value         types.String `tfsdk:"value"`
	Variable      types.String `tfsdk:"variable"`
}

// filterSettingModel mirrors the filter_setting NestedBlockObject, reused
// at the top level and inside each action.
type filterSettingModel struct {
	Conditions types.List   `tfsdk:"conditions"`
	Expression types.String `tfsdk:"expression"`
	Relation   types.String `tfsdk:"relation"`
}

// conditionModel mirrors the conditions NestedBlockObject.
type conditionModel struct {
	Field types.String `tfsdk:"field"`
	Op    types.String `tfsdk:"op"`
	Value types.String `tfsdk:"value"`
}

// ---- object element types --------------------------------------------------
//
// The framework represents a nested block element as an ObjectType whose
// AttrTypes match the NestedBlockObject attributes. These helpers return
// the canonical element types so the state lists can be built (and nulled)
// without a live object to derive the type from.

func conditionObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"field": types.StringType,
		"op":    types.StringType,
		"value": types.StringType,
	}}
}

func filterSettingObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"conditions": types.ListType{ElemType: conditionObjectType()},
		"expression": types.StringType,
		"relation":   types.StringType,
	}}
}

func actionObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"filter_setting": types.ListType{ElemType: filterSettingObjectType()},
		"label_key":      types.StringType,
		"mapping":        types.MapType{ElemType: types.StringType},
		"reg_exp":        types.StringType,
		"source":         types.StringType,
		"target":         types.StringType,
		"type":           types.StringType,
		"value":          types.StringType,
		"variable":       types.StringType,
	}}
}

// ---- Metadata / Schema ----------------------------------------------------

func (r *transformerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = cmsTypeName
}

func (r *transformerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a CMS Transformer resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The composite ID of the resource, formatted `<transformer_id>:<workspace>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"transformer_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier generated by the server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"transformer_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the transformer.",
			},
			"workspace": schema.StringAttribute{
				Required:    true,
				Description: "The workspace ID for resource isolation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the transformer.",
			},
			"quit_after_match": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to stop after the first match.",
			},
			"sort_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The sort order of the transformer.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"enable": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the transformer is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the transformer.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "The last update time of the transformer.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The user ID of the creator.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"actions": schema.ListNestedBlock{
				Description: "Transform actions with field mapping, regex extraction and filter settings.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"label_key": schema.StringAttribute{Optional: true},
						"mapping": schema.MapAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						"reg_exp":  schema.StringAttribute{Optional: true},
						"source":   schema.StringAttribute{Optional: true},
						"target":   schema.StringAttribute{Optional: true},
						"type":     schema.StringAttribute{Optional: true},
						"value":    schema.StringAttribute{Optional: true},
						"variable": schema.StringAttribute{Optional: true},
					},
					Blocks: map[string]schema.Block{
						"filter_setting": schema.ListNestedBlock{
							NestedObject: filterSettingBlockObject(),
						},
					},
				},
			},
			"filter_setting": schema.ListNestedBlock{
				Description:  "The top-level event filter setting.",
				NestedObject: filterSettingBlockObject(),
			},
		},
	}
}

// filterSettingBlockObject returns the shared filter_setting nested block
// object, used both at the top level and inside each action. conditions is a
// nested block (proto5-compatible); expression and relation are scalar
// attributes.
func filterSettingBlockObject() schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"expression": schema.StringAttribute{Optional: true},
			"relation":   schema.StringAttribute{Optional: true},
		},
		Blocks: map[string]schema.Block{
			"conditions": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"field": schema.StringAttribute{Optional: true},
						"op":    schema.StringAttribute{Optional: true},
						"value": schema.StringAttribute{Optional: true},
					},
				},
			},
		},
	}
}

// ---- CRUD -----------------------------------------------------------------

func (r *transformerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan transformerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()

	query := map[string]*string{
		"workspace": stringPointer(plan.Workspace.ValueString()),
		"regionId":  stringPointer(client.RegionId),
	}
	body := buildTransformerRequestBody(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := client.RoaPost(cmsProduct, cmsVersion, "/transformers", query, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateTransformer", cmsTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateTransformer", cmsTypeName).Error(),
		)
		return
	}

	transformerId := strVal(response, "transformerId")
	if transformerId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateTransformer", cmsTypeName),
			fmt.Sprintf("the response carried no transformerId: %v", response),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s:%s", transformerId, plan.Workspace.ValueString()))
	plan.TransformerId = types.StringValue(transformerId)

	// The create response only carries transformerId; refresh the full object
	// from GetTransformer so the state reflects the stored scalar and nested
	// values rather than the partial create payload.
	object, err := describeCmsTransformer(client, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetTransformer", cmsTypeName),
			errs.WrapErrorf(err, "Creating %s: calling GetTransformer", cmsTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling GetTransformer", cmsTypeName),
			fmt.Sprintf("transformer %s vanished right after create", plan.Id.ValueString()),
		)
		return
	}

	// The create endpoint does not accept enable; the server applies a default.
	// If the user requested a different enable state, reconcile it via the
	// dedicated sub-operation so the first apply converges instead of leaving a
	// plan diff for the next run to resolve.
	if !plan.Enable.IsNull() && !plan.Enable.IsUnknown() {
		if !plan.Enable.Equal(boolType(object["enable"])) {
			op := "disable"
			if plan.Enable.ValueBool() {
				op = "enable"
			}
			if err := toggleCmsTransformer(client, transformerId, plan.Workspace.ValueString(), op); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: calling %sTransformer", cmsTypeName, op),
					err.Error(),
				)
				return
			}
			object, err = describeCmsTransformer(client, plan.Id.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: calling GetTransformer", cmsTypeName),
					errs.WrapErrorf(err, "Creating %s: calling GetTransformer", cmsTypeName).Error(),
				)
				return
			}
			if object == nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Creating %s: calling GetTransformer", cmsTypeName),
					fmt.Sprintf("transformer %s vanished after %s", plan.Id.ValueString(), op),
				)
				return
			}
		}
	}

	applyObjectToState(ctx, &plan, object, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *transformerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state transformerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()

	object, err := describeCmsTransformer(client, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetTransformer", cmsTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetTransformer", cmsTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	applyObjectToState(ctx, &state, object, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *transformerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state transformerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()

	transformerId, workspace, _ := splitTransformerId(state.Id.ValueString())

	// Enable/disable is a dedicated sub-operation; it must not ride along on
	// the general Update payload. Toggle only when the user changed it.
	if !plan.Enable.Equal(state.Enable) {
		op := "disable"
		if plan.Enable.ValueBool() {
			op = "enable"
		}
		if err := toggleCmsTransformer(client, transformerId, workspace, op); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling %sTransformer", cmsTypeName, op),
				err.Error(),
			)
			return
		}
	}

	query := map[string]*string{
		"workspace": stringPointer(workspace),
		"regionId":  stringPointer(client.RegionId),
	}
	body := buildTransformerRequestBody(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	action := fmt.Sprintf("/transformers/%s", transformerId)
	if _, err := client.RoaPut(cmsProduct, cmsVersion, action, query, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateTransformer", cmsTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateTransformer", cmsTypeName).Error(),
		)
		return
	}

	// Preserve the composite id across the update; the server-generated half is
	// immutable. Refresh from the server so the state reflects what
	// GetTransformer returns rather than the request payload we just sent.
	plan.Id = state.Id
	plan.TransformerId = state.TransformerId
	object, err := describeCmsTransformer(client, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetTransformer", cmsTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetTransformer", cmsTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	applyObjectToState(ctx, &plan, object, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *transformerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state transformerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()

	transformerId, workspace, _ := splitTransformerId(state.Id.ValueString())
	action := fmt.Sprintf("/transformers/%s", transformerId)
	query := map[string]*string{
		"workspace": stringPointer(workspace),
		"regionId":  stringPointer(client.RegionId),
	}

	_, err := client.RoaDelete(cmsProduct, cmsVersion, action, query, nil, nil, true)
	if err != nil && !errs.NotFoundError(err) && !errs.IsExpectedErrors(err, []string{"NotFound", "ResourceNotFound"}) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteTransformer", cmsTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteTransformer", cmsTypeName).Error(),
		)
		return
	}
}

func (r *transformerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// ---- request body ---------------------------------------------------------

// buildTransformerRequestBody mirrors the SDKv2 buildCmsTransformerRequestBody:
// only populated fields are sent, and the actions/filter_setting subtrees are
// shaped into the camelCase API contract.
func buildTransformerRequestBody(ctx context.Context, m *transformerResourceModel, diags *diag.Diagnostics) map[string]interface{} {
	request := map[string]interface{}{}
	if v := m.TransformerName.ValueString(); v != "" {
		request["transformerName"] = v
	}
	if v := m.Description.ValueString(); v != "" {
		request["description"] = v
	}
	if !m.QuitAfterMatch.IsNull() && !m.QuitAfterMatch.IsUnknown() {
		request["quitAfterMatch"] = m.QuitAfterMatch.ValueBool()
	}
	if !m.SortId.IsNull() && !m.SortId.IsUnknown() {
		request["sortId"] = m.SortId.ValueInt64()
	}

	if !m.FilterSetting.IsNull() && !m.FilterSetting.IsUnknown() {
		var fsModels []filterSettingModel
		diags.Append(m.FilterSetting.ElementsAs(ctx, &fsModels, false)...)
		if len(fsModels) > 0 {
			if fs := buildFilterSettingRequest(ctx, &fsModels[0], diags); len(fs) > 0 {
				request["filterSetting"] = fs
			}
		}
	}

	if !m.Actions.IsNull() && !m.Actions.IsUnknown() {
		var actionModels []actionModel
		diags.Append(m.Actions.ElementsAs(ctx, &actionModels, false)...)
		actions := make([]map[string]interface{}, 0, len(actionModels))
		for i := range actionModels {
			actions = append(actions, buildActionRequest(ctx, &actionModels[i], diags))
		}
		if len(actions) > 0 {
			request["actions"] = actions
		}
	}
	return request
}

func buildActionRequest(ctx context.Context, a *actionModel, diags *diag.Diagnostics) map[string]interface{} {
	out := map[string]interface{}{}
	if v := a.LabelKey.ValueString(); v != "" {
		out["labelKey"] = v
	}
	if v := a.RegExp.ValueString(); v != "" {
		out["regExp"] = v
	}
	if v := a.Source.ValueString(); v != "" {
		out["source"] = v
	}
	if v := a.Target.ValueString(); v != "" {
		out["target"] = v
	}
	if v := a.Type.ValueString(); v != "" {
		out["type"] = v
	}
	if v := a.Value.ValueString(); v != "" {
		out["value"] = v
	}
	if v := a.Variable.ValueString(); v != "" {
		out["variable"] = v
	}
	if !a.Mapping.IsNull() && !a.Mapping.IsUnknown() {
		mapping := map[string]string{}
		diags.Append(a.Mapping.ElementsAs(ctx, &mapping, false)...)
		if len(mapping) > 0 {
			out["mapping"] = mapping
		}
	}
	if !a.FilterSetting.IsNull() && !a.FilterSetting.IsUnknown() {
		var fsModels []filterSettingModel
		diags.Append(a.FilterSetting.ElementsAs(ctx, &fsModels, false)...)
		if len(fsModels) > 0 {
			if fs := buildFilterSettingRequest(ctx, &fsModels[0], diags); len(fs) > 0 {
				out["filterSetting"] = fs
			}
		}
	}
	return out
}

func buildFilterSettingRequest(ctx context.Context, fs *filterSettingModel, diags *diag.Diagnostics) map[string]interface{} {
	out := map[string]interface{}{}
	if !fs.Conditions.IsNull() && !fs.Conditions.IsUnknown() {
		var condModels []conditionModel
		diags.Append(fs.Conditions.ElementsAs(ctx, &condModels, false)...)
		conditions := make([]map[string]string, 0, len(condModels))
		for i := range condModels {
			c := map[string]string{}
			if v := condModels[i].Field.ValueString(); v != "" {
				c["field"] = v
			}
			if v := condModels[i].Op.ValueString(); v != "" {
				c["op"] = v
			}
			if v := condModels[i].Value.ValueString(); v != "" {
				c["value"] = v
			}
			conditions = append(conditions, c)
		}
		if len(conditions) > 0 {
			out["conditions"] = conditions
		}
	}
	if v := fs.Expression.ValueString(); v != "" {
		out["expression"] = v
	}
	if v := fs.Relation.ValueString(); v != "" {
		out["relation"] = v
	}
	return out
}

// ---- state mapping --------------------------------------------------------

// applyObjectToState refreshes the model from the GetTransformer "transformer"
// object (or the create response carrying transformerId), including the
// nested actions/filter_setting subtrees.
func applyObjectToState(ctx context.Context, m *transformerResourceModel, object map[string]interface{}, diags *diag.Diagnostics) {
	m.TransformerName = strType(strVal(object, "transformerName"))
	m.Workspace = strType(strVal(object, "workspace"))
	m.Description = strType(strVal(object, "description"))
	m.QuitAfterMatch = boolType(object["quitAfterMatch"])
	m.SortId = intType(object["sortId"])
	m.Enable = boolType(object["enable"])
	m.CreateTime = strType(strVal(object, "createTime"))
	m.UpdateTime = strType(strVal(object, "updateTime"))
	m.UserId = strType(strVal(object, "userId"))
	m.RegionId = strType(strVal(object, "regionId"))
	if m.TransformerId.IsNull() || m.TransformerId.ValueString() == "" {
		m.TransformerId = strType(strVal(object, "transformerId"))
	}
	if m.Id.IsNull() || m.Id.ValueString() == "" {
		if tid := strVal(object, "transformerId"); tid != "" {
			m.Id = types.StringValue(fmt.Sprintf("%s:%s", tid, strVal(object, "workspace")))
		}
	}

	m.FilterSetting = filterSettingListFromAPI(ctx, object["filterSetting"], diags)
	m.Actions = actionsListFromAPI(ctx, object["actions"], diags)
}

// actionsListFromAPI builds the actions types.List from the API array.
func actionsListFromAPI(ctx context.Context, raw interface{}, diags *diag.Diagnostics) types.List {
	items := toInterfaceSlice(raw)
	if len(items) == 0 {
		return types.ListNull(actionObjectType())
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		am := actionModel{
			LabelKey:      strType(strVal(itemMap(item), "labelKey")),
			RegExp:        strType(strVal(itemMap(item), "regExp")),
			Source:        strType(strVal(itemMap(item), "source")),
			Target:        strType(strVal(itemMap(item), "target")),
			Type:          strType(strVal(itemMap(item), "type")),
			Value:         strType(strVal(itemMap(item), "value")),
			Variable:      strType(strVal(itemMap(item), "variable")),
			Mapping:       mappingMapFromAPI(ctx, itemMap(item)["mapping"], diags),
			FilterSetting: filterSettingListFromAPI(ctx, itemMap(item)["filterSetting"], diags),
		}
		obj, d := types.ObjectValueFrom(ctx, actionObjectType().AttrTypes, &am)
		diags.Append(d...)
		values = append(values, obj)
	}
	list, d := types.ListValue(actionObjectType(), values)
	diags.Append(d...)
	return list
}

func filterSettingListFromAPI(ctx context.Context, raw interface{}, diags *diag.Diagnostics) types.List {
	items := toInterfaceSlice(raw)
	if len(items) == 0 {
		return types.ListNull(filterSettingObjectType())
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		fm := filterSettingModel{
			Conditions: conditionsListFromAPI(ctx, itemMap(item)["conditions"], diags),
			Expression: strType(strVal(itemMap(item), "expression")),
			Relation:   strType(strVal(itemMap(item), "relation")),
		}
		obj, d := types.ObjectValueFrom(ctx, filterSettingObjectType().AttrTypes, &fm)
		diags.Append(d...)
		values = append(values, obj)
	}
	list, d := types.ListValue(filterSettingObjectType(), values)
	diags.Append(d...)
	return list
}

func conditionsListFromAPI(ctx context.Context, raw interface{}, diags *diag.Diagnostics) types.List {
	items := toInterfaceSlice(raw)
	if len(items) == 0 {
		return types.ListNull(conditionObjectType())
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		cm := conditionModel{
			Field: strType(strVal(itemMap(item), "field")),
			Op:    strType(strVal(itemMap(item), "op")),
			Value: strType(strVal(itemMap(item), "value")),
		}
		obj, d := types.ObjectValueFrom(ctx, conditionObjectType().AttrTypes, &cm)
		diags.Append(d...)
		values = append(values, obj)
	}
	list, d := types.ListValue(conditionObjectType(), values)
	diags.Append(d...)
	return list
}

func mappingMapFromAPI(ctx context.Context, raw interface{}, diags *diag.Diagnostics) types.Map {
	src := itemMap(raw)
	m := make(map[string]string, len(src))
	for k, v := range src {
		m[k] = toStr(v)
	}
	if len(m) == 0 {
		return types.MapNull(types.StringType)
	}
	mv, d := types.MapValueFrom(ctx, types.StringType, m)
	diags.Append(d...)
	return mv
}

// ---- value coercion helpers ----------------------------------------------

func strVal(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	return toStr(m[key])
}

func toStr(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case bool:
		return strconv.FormatBool(t)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	default:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%v", v)
	}
}

func itemMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return nil
}

func toInterfaceSlice(v interface{}) []interface{} {
	switch t := v.(type) {
	case []interface{}:
		return t
	case []map[string]interface{}:
		out := make([]interface{}, len(t))
		for i := range t {
			out[i] = t[i]
		}
		return out
	default:
		return nil
	}
}

func strType(v string) types.String {
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(v)
}

func boolType(v interface{}) types.Bool {
	switch t := v.(type) {
	case bool:
		return types.BoolValue(t)
	case string:
		if b, err := strconv.ParseBool(t); err == nil {
			return types.BoolValue(b)
		}
	}
	return types.BoolNull()
}

func intType(v interface{}) types.Int64 {
	switch t := v.(type) {
	case float64:
		return types.Int64Value(int64(t))
	case int:
		return types.Int64Value(int64(t))
	case string:
		if n, err := strconv.ParseInt(t, 10, 64); err == nil {
			return types.Int64Value(n)
		}
	}
	return types.Int64Null()
}
