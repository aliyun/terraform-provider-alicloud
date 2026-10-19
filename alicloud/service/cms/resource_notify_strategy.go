package cms

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const notifyStrategyTypeName = "alicloud_cms_notify_strategy"

var (
	_ resource.Resource                = &notifyStrategyResource{}
	_ resource.ResourceWithImportState = &notifyStrategyResource{}
)

func NewNotifyStrategyResource() resource.Resource {
	return &notifyStrategyResource{}
}

type notifyStrategyResource struct {
	fwadapt.ResourceBase
}

type notifyStrategyResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	NotifyStrategyName         types.String `tfsdk:"notify_strategy_name"`
	Description                types.String `tfsdk:"description"`
	IgnoreRestoredNotification types.Bool   `tfsdk:"ignore_restored_notification"`
	Workspace                  types.String `tfsdk:"workspace"`
	NotifyStrategyId           types.String `tfsdk:"notify_strategy_id"`
	Enable                     types.Bool   `tfsdk:"enable"`
	CreateTime                 types.String `tfsdk:"create_time"`
	UpdateTime                 types.String `tfsdk:"update_time"`
	UserId                     types.String `tfsdk:"user_id"`
	GroupingSetting            types.Object `tfsdk:"grouping_setting"`
	Routes                     types.List   `tfsdk:"routes"`
	CustomTemplateEntries      types.List   `tfsdk:"custom_template_entries"`
}

// Attribute type maps for the nested-block objects. These mirror the schema
// below and are used to build types.Object / types.List values from API data.
var (
	groupingSettingAttrTypes = map[string]attr.Type{
		"grouping_keys": types.ListType{ElemType: types.StringType},
		"period_min":    types.Int64Type,
		"silence_sec":   types.Int64Type,
		"times":         types.Int64Type,
	}
	channelAttrTypes = map[string]attr.Type{
		"channel_type":         types.StringType,
		"enabled_sub_channels": types.ListType{ElemType: types.StringType},
		"receivers":            types.ListType{ElemType: types.StringType},
	}
	effectTimeRangeAttrTypes = map[string]attr.Type{
		"day_in_week":          types.ListType{ElemType: types.Int64Type},
		"end_time_in_minute":   types.Int64Type,
		"start_time_in_minute": types.Int64Type,
		"time_zone":            types.StringType,
	}
	conditionAttrTypes = map[string]attr.Type{
		"field": types.StringType,
		"op":    types.StringType,
		"value": types.StringType,
	}
	filterSettingAttrTypes = map[string]attr.Type{
		"conditions": types.ListType{ElemType: types.ObjectType{AttrTypes: conditionAttrTypes}},
		"expression": types.StringType,
		"relation":   types.StringType,
	}
	routeAttrTypes = map[string]attr.Type{
		"channels":          types.ListType{ElemType: types.ObjectType{AttrTypes: channelAttrTypes}},
		"effect_time_range": types.ObjectType{AttrTypes: effectTimeRangeAttrTypes},
		"filter_setting":    types.ObjectType{AttrTypes: filterSettingAttrTypes},
		"severities":        types.ListType{ElemType: types.StringType},
	}
	customTemplateEntryAttrTypes = map[string]attr.Type{
		"target_type":   types.StringType,
		"template_uuid": types.StringType,
	}
)

func (r *notifyStrategyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = notifyStrategyTypeName
}

func (r *notifyStrategyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Cms Notify Strategy resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `notify_strategy_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"notify_strategy_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the notify strategy.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the notify strategy.",
			},
			"ignore_restored_notification": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to ignore restored notifications.",
			},
			"workspace": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The workspace id of the notify strategy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"notify_strategy_id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the notify strategy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enable": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the notify strategy is enabled.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the notify strategy was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the notify strategy was last updated.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The user id of the notify strategy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"grouping_setting": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"grouping_keys": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "The grouping keys of the grouping setting.",
					},
					"period_min": schema.Int64Attribute{
						Optional:    true,
						Description: "The period (in minutes) of the grouping setting.",
					},
					"silence_sec": schema.Int64Attribute{
						Optional:    true,
						Description: "The silence duration (in seconds) of the grouping setting.",
					},
					"times": schema.Int64Attribute{
						Optional:    true,
						Description: "The times of the grouping setting.",
					},
				},
			},
			"routes": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"severities": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "The severities of the route.",
						},
					},
					Blocks: map[string]schema.Block{
						"channels": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"channel_type": schema.StringAttribute{
										Required:    true,
										Description: "The channel type of the channel.",
									},
									"enabled_sub_channels": schema.ListAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Description: "The enabled sub channels of the channel.",
									},
									"receivers": schema.ListAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Description: "The receivers of the channel.",
									},
								},
							},
						},
						"effect_time_range": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"day_in_week": schema.ListAttribute{
									ElementType: types.Int64Type,
									Optional:    true,
									Description: "The days in a week of the effect time range.",
								},
								"end_time_in_minute": schema.Int64Attribute{
									Optional:    true,
									Description: "The end time (in minute) of the effect time range.",
								},
								"start_time_in_minute": schema.Int64Attribute{
									Optional:    true,
									Description: "The start time (in minute) of the effect time range.",
								},
								"time_zone": schema.StringAttribute{
									Optional:    true,
									Description: "The time zone of the effect time range.",
								},
							},
						},
						"filter_setting": schema.SingleNestedBlock{
							Attributes: map[string]schema.Attribute{
								"expression": schema.StringAttribute{
									Optional:    true,
									Description: "The expression of the filter setting.",
								},
								"relation": schema.StringAttribute{
									Optional:    true,
									Description: "The relation of the filter setting.",
								},
							},
							Blocks: map[string]schema.Block{
								"conditions": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"field": schema.StringAttribute{
												Required:    true,
												Description: "The field of the condition.",
											},
											"op": schema.StringAttribute{
												Required:    true,
												Description: "The operator of the condition.",
											},
											"value": schema.StringAttribute{
												Required:    true,
												Description: "The value of the condition.",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"custom_template_entries": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"target_type": schema.StringAttribute{
							Required:    true,
							Description: "The target type of the custom template entry.",
						},
						"template_uuid": schema.StringAttribute{
							Required:    true,
							Description: "The template uuid of the custom template entry.",
						},
					},
				},
			},
		},
	}
}

func (r *notifyStrategyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notifyStrategyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := map[string]*string{"RegionId": stringPtr(r.Client().RegionId)}
	if v := plan.Workspace.ValueString(); v != "" {
		query["workspace"] = stringPtr(v)
	}

	body := map[string]interface{}{
		"name": plan.NotifyStrategyName.ValueString(),
	}
	if v := plan.Description.ValueString(); v != "" {
		body["description"] = v
	}
	if !plan.IgnoreRestoredNotification.IsNull() {
		body["ignoreRestoredNotification"] = plan.IgnoreRestoredNotification.ValueBool()
	}
	if !plan.GroupingSetting.IsNull() {
		body["groupingSetting"] = expandGroupingSetting(plan.GroupingSetting)
	}
	if !plan.Routes.IsNull() {
		body["routes"] = expandRoutes(plan.Routes)
	}
	if !plan.CustomTemplateEntries.IsNull() {
		body["customTemplateEntries"] = expandCustomTemplateEntries(plan.CustomTemplateEntries)
	}

	response, err := r.Client().RoaPost("Cms", "2024-03-30", "/notifyStrategies", query, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateNotifyStrategy", notifyStrategyTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateNotifyStrategy", notifyStrategyTypeName).Error(),
		)
		return
	}

	id := fmt.Sprint(response["data"])
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateNotifyStrategy", notifyStrategyTypeName),
			fmt.Sprintf("The response carried no id at $.data: %v", response),
		)
		return
	}

	object, err := readNotifyStrategy(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetNotifyStrategy", notifyStrategyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetNotifyStrategy after create", notifyStrategyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateNotifyStrategy", notifyStrategyTypeName),
			fmt.Sprintf("The notify strategy %s is not visible via GetNotifyStrategy right after creation", id),
		)
		return
	}

	state := notifyStrategyResourceModel{}
	setNotifyStrategyState(&state, object)
	state.Id = types.StringValue(id)
	state.NotifyStrategyId = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *notifyStrategyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notifyStrategyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := readNotifyStrategy(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetNotifyStrategy", notifyStrategyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetNotifyStrategy", notifyStrategyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState := notifyStrategyResourceModel{}
	setNotifyStrategyState(&newState, object)
	newState.Id = state.Id
	newState.NotifyStrategyId = state.Id
	if newState.Workspace.IsNull() && !state.Workspace.IsNull() {
		newState.Workspace = state.Workspace
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *notifyStrategyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state notifyStrategyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := map[string]*string{"RegionId": stringPtr(r.Client().RegionId)}
	if v := plan.Workspace.ValueString(); v != "" {
		query["workspace"] = stringPtr(v)
	}

	body := map[string]interface{}{
		"name": plan.NotifyStrategyName.ValueString(),
	}
	if v := plan.Description.ValueString(); v != "" {
		body["description"] = v
	}
	if !plan.IgnoreRestoredNotification.IsNull() {
		body["ignoreRestoredNotification"] = plan.IgnoreRestoredNotification.ValueBool()
	}
	if !plan.GroupingSetting.IsNull() {
		body["groupingSetting"] = expandGroupingSetting(plan.GroupingSetting)
	}
	if !plan.Routes.IsNull() {
		body["routes"] = expandRoutes(plan.Routes)
	}
	if !plan.CustomTemplateEntries.IsNull() {
		body["customTemplateEntries"] = expandCustomTemplateEntries(plan.CustomTemplateEntries)
	}

	action := fmt.Sprintf("/notifyStrategies/%s", state.Id.ValueString())
	if _, err := r.Client().RoaPut("Cms", "2024-03-30", action, query, nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateNotifyStrategy", notifyStrategyTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateNotifyStrategy", notifyStrategyTypeName).Error(),
		)
		return
	}

	object, err := readNotifyStrategy(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetNotifyStrategy", notifyStrategyTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetNotifyStrategy after update", notifyStrategyTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState := notifyStrategyResourceModel{}
	setNotifyStrategyState(&newState, object)
	newState.Id = state.Id
	newState.NotifyStrategyId = state.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *notifyStrategyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notifyStrategyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := map[string]*string{"RegionId": stringPtr(r.Client().RegionId)}
	if v := state.Workspace.ValueString(); v != "" {
		query["workspace"] = stringPtr(v)
	}
	action := fmt.Sprintf("/notifyStrategies/%s", state.Id.ValueString())
	if _, err := r.Client().RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true); err != nil && !isNotifyStrategyNotFound(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteNotifyStrategy", notifyStrategyTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteNotifyStrategy", notifyStrategyTypeName).Error(),
		)
		return
	}
}

func (r *notifyStrategyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// readNotifyStrategy calls GetNotifyStrategy and returns the data object, or
// nil when the strategy no longer exists.
func readNotifyStrategy(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	query := map[string]*string{"RegionId": stringPtr(client.RegionId)}
	action := fmt.Sprintf("/notifyStrategies/%s", id)
	response, err := client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
	if err != nil {
		if isNotifyStrategyNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	data, _ := response["data"].(map[string]interface{})
	return data, nil
}

// setNotifyStrategyState maps the API data object onto the resource model.
func setNotifyStrategyState(m *notifyStrategyResourceModel, object map[string]interface{}) {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	m.NotifyStrategyName = str("notifyStrategyName")
	m.Description = str("description")
	m.NotifyStrategyId = str("notifyStrategyId")
	m.CreateTime = str("createTime")
	m.UpdateTime = str("updateTime")
	m.UserId = str("userId")
	m.Workspace = str("workspace")
	if v, ok := object["enabled"].(bool); ok {
		m.Enable = types.BoolValue(v)
	} else {
		m.Enable = types.BoolNull()
	}
	if v, ok := object["ignoreRestoredNotification"].(bool); ok {
		m.IgnoreRestoredNotification = types.BoolValue(v)
	} else {
		m.IgnoreRestoredNotification = types.BoolNull()
	}
	m.GroupingSetting = flattenGroupingSetting(object["groupingSetting"])
	m.Routes = flattenRoutes(object["routes"])
	m.CustomTemplateEntries = flattenCustomTemplateEntries(object["customTemplateEntries"])
}

// isNotifyStrategyNotFound reports whether err is the not-found error for a
// notify strategy. The code arrives as a wrapped tea SDK error; the substring
// fallback covers wrapped forms, matching IsExpectedErrors' behavior.
func isNotifyStrategyNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "NotFound") || strings.Contains(msg, "ResourceNotFound")
}

func stringPtr(s string) *string { return &s }

// ---- expand (model -> API request) ----

func expandGroupingSetting(obj types.Object) map[string]interface{} {
	out := map[string]interface{}{}
	if obj.IsNull() {
		return out
	}
	attrs := obj.Attributes()
	if v, ok := attrs["grouping_keys"].(types.List); ok && !v.IsNull() {
		out["groupingKeys"] = stringSliceFromList(v)
	}
	if v, ok := attrs["period_min"].(types.Int64); ok && !v.IsNull() && v.ValueInt64() != 0 {
		out["periodMin"] = v.ValueInt64()
	}
	if v, ok := attrs["silence_sec"].(types.Int64); ok && !v.IsNull() && v.ValueInt64() != 0 {
		out["silenceSec"] = v.ValueInt64()
	}
	if v, ok := attrs["times"].(types.Int64); ok && !v.IsNull() && v.ValueInt64() != 0 {
		out["times"] = v.ValueInt64()
	}
	return out
}

func expandRoutes(list types.List) []interface{} {
	if list.IsNull() {
		return nil
	}
	out := make([]interface{}, 0, len(list.Elements()))
	for _, el := range list.Elements() {
		obj, ok := el.(types.Object)
		if !ok {
			continue
		}
		attrs := obj.Attributes()
		rm := map[string]interface{}{}

		if channels, ok := attrs["channels"].(types.List); ok && !channels.IsNull() {
			chList := make([]interface{}, 0, len(channels.Elements()))
			for _, ce := range channels.Elements() {
				co, ok := ce.(types.Object)
				if !ok {
					continue
				}
				ca := co.Attributes()
				cm := map[string]interface{}{}
				if ct, ok := ca["channel_type"].(types.String); ok && ct.ValueString() != "" {
					cm["channelType"] = ct.ValueString()
				}
				if esc, ok := ca["enabled_sub_channels"].(types.List); ok && !esc.IsNull() {
					cm["enabledSubChannels"] = stringSliceFromList(esc)
				}
				if rec, ok := ca["receivers"].(types.List); ok && !rec.IsNull() {
					cm["receivers"] = stringSliceFromList(rec)
				}
				chList = append(chList, cm)
			}
			rm["channels"] = chList
		}

		if etr, ok := attrs["effect_time_range"].(types.Object); ok && !etr.IsNull() {
			ea := etr.Attributes()
			em := map[string]interface{}{}
			if tz, ok := ea["time_zone"].(types.String); ok && tz.ValueString() != "" {
				em["timeZone"] = tz.ValueString()
			}
			if diw, ok := ea["day_in_week"].(types.List); ok && !diw.IsNull() {
				em["dayInWeek"] = intSliceFromList(diw)
			}
			if stm, ok := ea["start_time_in_minute"].(types.Int64); ok && !stm.IsNull() {
				em["startTimeInMinute"] = stm.ValueInt64()
			}
			if etm, ok := ea["end_time_in_minute"].(types.Int64); ok && !etm.IsNull() {
				em["endTimeInMinute"] = etm.ValueInt64()
			}
			rm["effectTimeRange"] = em
		}

		if fs, ok := attrs["filter_setting"].(types.Object); ok && !fs.IsNull() {
			fa := fs.Attributes()
			fm := map[string]interface{}{}
			if expr, ok := fa["expression"].(types.String); ok && expr.ValueString() != "" {
				fm["expression"] = expr.ValueString()
			}
			if rel, ok := fa["relation"].(types.String); ok && rel.ValueString() != "" {
				fm["relation"] = rel.ValueString()
			}
			if conds, ok := fa["conditions"].(types.List); ok && !conds.IsNull() {
				cl := make([]interface{}, 0, len(conds.Elements()))
				for _, ce := range conds.Elements() {
					co, ok := ce.(types.Object)
					if !ok {
						continue
					}
					ca := co.Attributes()
					cm := map[string]interface{}{}
					if f, ok := ca["field"].(types.String); ok && f.ValueString() != "" {
						cm["field"] = f.ValueString()
					}
					if o, ok := ca["op"].(types.String); ok && o.ValueString() != "" {
						cm["op"] = o.ValueString()
					}
					if v, ok := ca["value"].(types.String); ok && v.ValueString() != "" {
						cm["value"] = v.ValueString()
					}
					cl = append(cl, cm)
				}
				fm["conditions"] = cl
			}
			rm["filterSetting"] = fm
		}

		if sev, ok := attrs["severities"].(types.List); ok && !sev.IsNull() {
			rm["severities"] = stringSliceFromList(sev)
		}

		out = append(out, rm)
	}
	return out
}

func expandCustomTemplateEntries(list types.List) []interface{} {
	if list.IsNull() {
		return nil
	}
	out := make([]interface{}, 0, len(list.Elements()))
	for _, el := range list.Elements() {
		obj, ok := el.(types.Object)
		if !ok {
			continue
		}
		ca := obj.Attributes()
		cm := map[string]interface{}{}
		if tt, ok := ca["target_type"].(types.String); ok && tt.ValueString() != "" {
			cm["targetType"] = tt.ValueString()
		}
		if tu, ok := ca["template_uuid"].(types.String); ok && tu.ValueString() != "" {
			cm["templateUuid"] = tu.ValueString()
		}
		out = append(out, cm)
	}
	return out
}

func stringSliceFromList(list types.List) []interface{} {
	if list.IsNull() {
		return nil
	}
	out := make([]interface{}, 0, len(list.Elements()))
	for _, el := range list.Elements() {
		if s, ok := el.(types.String); ok {
			out = append(out, s.ValueString())
		}
	}
	return out
}

func intSliceFromList(list types.List) []interface{} {
	if list.IsNull() {
		return nil
	}
	out := make([]interface{}, 0, len(list.Elements()))
	for _, el := range list.Elements() {
		if v, ok := el.(types.Int64); ok {
			out = append(out, v.ValueInt64())
		}
	}
	return out
}

// ---- flatten (API response -> model attr.Value) ----

func flattenGroupingSetting(v interface{}) types.Object {
	m, _ := v.(map[string]interface{})
	if len(m) == 0 {
		return types.ObjectNull(groupingSettingAttrTypes)
	}
	values := map[string]attr.Value{
		"grouping_keys": flattenStringSlice(m["groupingKeys"]),
		"period_min":    flattenInt64(m["periodMin"]),
		"silence_sec":   flattenInt64(m["silenceSec"]),
		"times":         flattenInt64(m["times"]),
	}
	obj, _ := types.ObjectValue(groupingSettingAttrTypes, values)
	return obj
}

func flattenRoutes(v interface{}) types.List {
	arr := toInterfaceSlice(v)
	if len(arr) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: routeAttrTypes})
	}
	vals := make([]attr.Value, 0, len(arr))
	for _, x := range arr {
		rm, _ := x.(map[string]interface{})
		vals = append(vals, flattenRouteObject(rm))
	}
	list, _ := types.ListValue(types.ObjectType{AttrTypes: routeAttrTypes}, vals)
	return list
}

func flattenRouteObject(m map[string]interface{}) types.Object {
	values := map[string]attr.Value{
		"channels":          flattenChannels(m["channels"]),
		"effect_time_range": flattenEffectTimeRange(m["effectTimeRange"]),
		"filter_setting":    flattenFilterSetting(m["filterSetting"]),
		"severities":        flattenStringSlice(m["severities"]),
	}
	obj, _ := types.ObjectValue(routeAttrTypes, values)
	return obj
}

func flattenChannels(v interface{}) types.List {
	arr := toInterfaceSlice(v)
	if len(arr) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: channelAttrTypes})
	}
	vals := make([]attr.Value, 0, len(arr))
	for _, x := range arr {
		cm, _ := x.(map[string]interface{})
		values := map[string]attr.Value{
			"channel_type":         flattenString(cm["channelType"]),
			"enabled_sub_channels": flattenStringSlice(cm["enabledSubChannels"]),
			"receivers":            flattenStringSlice(cm["receivers"]),
		}
		obj, _ := types.ObjectValue(channelAttrTypes, values)
		vals = append(vals, obj)
	}
	list, _ := types.ListValue(types.ObjectType{AttrTypes: channelAttrTypes}, vals)
	return list
}

func flattenEffectTimeRange(v interface{}) types.Object {
	m, _ := v.(map[string]interface{})
	if len(m) == 0 {
		return types.ObjectNull(effectTimeRangeAttrTypes)
	}
	values := map[string]attr.Value{
		"day_in_week":          flattenInt64Slice(m["dayInWeek"]),
		"end_time_in_minute":   flattenInt64(m["endTimeInMinute"]),
		"start_time_in_minute": flattenInt64(m["startTimeInMinute"]),
		"time_zone":            flattenString(m["timeZone"]),
	}
	obj, _ := types.ObjectValue(effectTimeRangeAttrTypes, values)
	return obj
}

func flattenFilterSetting(v interface{}) types.Object {
	m, _ := v.(map[string]interface{})
	if len(m) == 0 {
		return types.ObjectNull(filterSettingAttrTypes)
	}
	values := map[string]attr.Value{
		"conditions": flattenConditions(m["conditions"]),
		"expression": flattenString(m["expression"]),
		"relation":   flattenString(m["relation"]),
	}
	obj, _ := types.ObjectValue(filterSettingAttrTypes, values)
	return obj
}

func flattenConditions(v interface{}) types.List {
	arr := toInterfaceSlice(v)
	if len(arr) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: conditionAttrTypes})
	}
	vals := make([]attr.Value, 0, len(arr))
	for _, x := range arr {
		cm, _ := x.(map[string]interface{})
		values := map[string]attr.Value{
			"field": flattenString(cm["field"]),
			"op":    flattenString(cm["op"]),
			"value": flattenString(cm["value"]),
		}
		obj, _ := types.ObjectValue(conditionAttrTypes, values)
		vals = append(vals, obj)
	}
	list, _ := types.ListValue(types.ObjectType{AttrTypes: conditionAttrTypes}, vals)
	return list
}

func flattenCustomTemplateEntries(v interface{}) types.List {
	arr := toInterfaceSlice(v)
	if len(arr) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: customTemplateEntryAttrTypes})
	}
	vals := make([]attr.Value, 0, len(arr))
	for _, x := range arr {
		cm, _ := x.(map[string]interface{})
		values := map[string]attr.Value{
			"target_type":   flattenString(cm["targetType"]),
			"template_uuid": flattenString(cm["templateUuid"]),
		}
		obj, _ := types.ObjectValue(customTemplateEntryAttrTypes, values)
		vals = append(vals, obj)
	}
	list, _ := types.ListValue(types.ObjectType{AttrTypes: customTemplateEntryAttrTypes}, vals)
	return list
}

func flattenString(v interface{}) types.String {
	if v == nil {
		return types.StringNull()
	}
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

func flattenInt64(v interface{}) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(toInt64(v))
}

func flattenStringSlice(v interface{}) types.List {
	arr := toInterfaceSlice(v)
	vals := make([]attr.Value, 0, len(arr))
	for _, x := range arr {
		if s, ok := x.(string); ok && s != "" {
			vals = append(vals, types.StringValue(s))
		}
	}
	list, _ := types.ListValue(types.StringType, vals)
	return list
}

func flattenInt64Slice(v interface{}) types.List {
	arr := toInterfaceSlice(v)
	vals := make([]attr.Value, 0, len(arr))
	for _, x := range arr {
		if x != nil {
			vals = append(vals, types.Int64Value(toInt64(x)))
		}
	}
	list, _ := types.ListValue(types.Int64Type, vals)
	return list
}

func toInterfaceSlice(v interface{}) []interface{} {
	if a, ok := v.([]interface{}); ok {
		return a
	}
	return nil
}

func toInt64(v interface{}) int64 {
	switch x := v.(type) {
	case float64:
		return int64(x)
	case int:
		return int64(x)
	case int64:
		return x
	case string:
		if n, err := strconv.ParseInt(x, 10, 64); err == nil {
			return n
		}
	}
	return 0
}
