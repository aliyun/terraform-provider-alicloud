package ecs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const portRangeListTypeName = "alicloud_ecs_port_range_list"

var (
	_ resource.Resource                = &portRangeListResource{}
	_ resource.ResourceWithImportState = &portRangeListResource{}
)

// entryObjectType is the object type of an entry element, used to build and
// null-out the entry set on read.
var entryObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"port_range":  types.StringType,
		"description": types.StringType,
	},
}

func NewPortRangeListResource() resource.Resource {
	return &portRangeListResource{}
}

type portRangeListResource struct {
	fwadapt.ResourceBase
}

type portRangeListResourceModel struct {
	Id                types.String `tfsdk:"id"`
	MaxEntries        types.Int64  `tfsdk:"max_entries"`
	PortRangeListName types.String `tfsdk:"port_range_list_name"`
	Description       types.String `tfsdk:"description"`
	Entry             types.Set    `tfsdk:"entry"`
	ResourceGroupId   types.String `tfsdk:"resource_group_id"`
	Tags              types.Map    `tfsdk:"tags"`
	AssociationCount  types.String `tfsdk:"association_count"`
	CreateTime        types.String `tfsdk:"create_time"`
}

// entryModel mirrors one Entry element of the CreatePortRangeList/ModifyPortRangeList
// arrays: a {PortRange, Description} pair.
type entryModel struct {
	PortRange   types.String `tfsdk:"port_range"`
	Description types.String `tfsdk:"description"`
}

func (r *portRangeListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = portRangeListTypeName
}

func (r *portRangeListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a ECS Port Range List resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Port Range List. Same as `port_range_list_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_entries": schema.Int64Attribute{
				Required:    true,
				Description: "The maximum number of entries that the port range list can contain. Valid values: `1` to `2000`. This can be updated by ModifyPortRangeList.",
			},
			"port_range_list_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the port range list. The name must be 2 to 128 characters in length, must start with an English letter or a Chinese character, and can contain Chinese characters, English letters, digits, colons (:), underscores (_), periods (.), and hyphens (-). It cannot start with `http://`, `https://`, `com.aliyun`, or `com.alibabacloud`.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the port range list. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.",
			},
			"resource_group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the resource group to which the port range list belongs. The resource group can only be set at Create: PortRangeList does not support JoinResourceGroup and ModifyPortRangeList does not accept ResourceGroupId.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The tags of the port range list.",
			},
			"association_count": schema.StringAttribute{
				Computed:    true,
				Description: "The number of resources associated with the port range list.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The time when the port range list was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"entry": schema.SetNestedBlock{
				Description: "The entries of the port range list.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"port_range": schema.StringAttribute{
							Optional:    true,
							Description: "The port range. Format: `<start>/<end>`, for example `80/80`. Port ranges in a list cannot be duplicated.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "The description of the port range. The description must be 2 to 32 characters in length and cannot start with `http://` or `https://`.",
						},
					},
				},
			},
		},
	}
}

func (r *portRangeListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan portRangeListResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"RegionId":          r.Client().RegionId,
		"ClientToken":       helper.BuildClientToken("CreatePortRangeList"),
		"MaxEntries":        plan.MaxEntries.ValueInt64(),
		"PortRangeListName": plan.PortRangeListName.ValueString(),
	}
	if !plan.Description.IsNull() {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.Entry.IsNull() {
		entryArray, d := entriesToRequest(ctx, plan.Entry)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		request["Entry"] = entryArray
	}
	if !plan.ResourceGroupId.IsNull() {
		request["ResourceGroupId"] = plan.ResourceGroupId.ValueString()
	}
	if !plan.Tags.IsNull() {
		i := 1
		var tags map[string]string
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for key, value := range tags {
			request[fmt.Sprintf("Tag.%d.Key", i)] = key
			request[fmt.Sprintf("Tag.%d.Value", i)] = value
			i++
		}
	}

	response, err := callRpcWithRetry(r.Client(), "CreatePortRangeList", request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePortRangeList", portRangeListTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreatePortRangeList", portRangeListTypeName).Error(),
		)
		return
	}

	id, _ := response["PortRangeListId"].(string)
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePortRangeList", portRangeListTypeName),
			fmt.Sprintf("The CreatePortRangeList response carried no PortRangeListId: %v", response),
		)
		return
	}

	// Persist the created id and the planned inputs before readback. If the
	// follow-up Describe fails or the list is not yet visible, the resource is
	// still tracked by id: a subsequent apply can Read, Update, or Delete it
	// instead of re-creating, and the CreatePortRangeList ClientToken dedups a
	// retry. Computed fields stay unknown here and are filled by the readback
	// below when it succeeds.
	plan.Id = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, d := buildPortRangeListState(ctx, r.Client(), id, plan.Tags)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreatePortRangeList", portRangeListTypeName),
			fmt.Sprintf("The port range list %s is not visible via DescribePortRangeLists right after creation; the id is persisted in state so a subsequent apply can reconcile it", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *portRangeListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var priorState portRangeListResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, d := buildPortRangeListState(ctx, r.Client(), priorState.Id.ValueString(), priorState.Tags)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *portRangeListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, priorState portRangeListResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ModifyPortRangeList updates name/description/max_entries and add/removes
	// entries in a single call. resource_group_id is ForceNew and handled below
	// by RequiresReplace; tags are reconciled separately via TagResources.
	if !plan.MaxEntries.Equal(priorState.MaxEntries) ||
		!plan.PortRangeListName.Equal(priorState.PortRangeListName) ||
		!plan.Description.Equal(priorState.Description) ||
		!plan.Entry.Equal(priorState.Entry) {
		request := map[string]interface{}{
			"PortRangeListId":   priorState.Id.ValueString(),
			"RegionId":          r.Client().RegionId,
			"ClientToken":       helper.BuildClientToken("ModifyPortRangeList"),
			"MaxEntries":        plan.MaxEntries.ValueInt64(),
			"PortRangeListName": plan.PortRangeListName.ValueString(),
		}
		if !plan.Description.Equal(priorState.Description) {
			request["Description"] = plan.Description.ValueString()
		}
		if !plan.Entry.Equal(priorState.Entry) {
			addEntry, removeEntry, d := diffEntries(ctx, priorState.Entry, plan.Entry)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			request["AddEntry"] = addEntry
			request["RemoveEntry"] = removeEntry
		}
		if _, err := callRpcWithRetry(r.Client(), "ModifyPortRangeList", request, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling ModifyPortRangeList", portRangeListTypeName),
				errs.WrapErrorf(err, "Updating %s: calling ModifyPortRangeList", portRangeListTypeName).Error(),
			)
			return
		}
	}

	if !plan.Tags.Equal(priorState.Tags) {
		var oldTags, newTags map[string]string
		resp.Diagnostics.Append(priorState.Tags.ElementsAs(ctx, &oldTags, false)...)
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &newTags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if err := setResourceTags(r.Client(), priorState.Id.ValueString(), oldTags, newTags); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: reconciling tags", portRangeListTypeName),
				errs.WrapErrorf(err, "Updating %s: reconciling tags", portRangeListTypeName).Error(),
			)
			return
		}
	}

	state, d := buildPortRangeListState(ctx, r.Client(), priorState.Id.ValueString(), plan.Tags)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *portRangeListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state portRangeListResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"PortRangeListId": state.Id.ValueString(),
		"RegionId":        r.Client().RegionId,
	}
	_, err := callRpcWithRetry(r.Client(), "DeletePortRangeList", request, true)
	if err != nil {
		if isPortRangeListNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeletePortRangeList", portRangeListTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeletePortRangeList", portRangeListTypeName).Error(),
		)
		return
	}
}

func (r *portRangeListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildPortRangeListState reads a PortRangeList and its entries and maps them
// onto the resource model. A nil model with no diagnostics means the list no
// longer exists (Read/Update remove the resource); a non-nil model is set as
// state; diagnostics with errors abort the caller.
//
// priorTags is the plan/prior-state Tags value so the readback preserves the
// null-vs-explicit-empty distinction the user configured. ECS legitimately
// accepts tags = {} (Tag is an optional array of length 0~20), and
// Create/Update both overwrite state with the value returned here: rewriting
// an explicit empty map to null would violate the HashiCorp Create/Update data
// consistency contract (every null/known plan value is preserved verbatim) and
// surface as an inconsistent-result-after-apply / refresh drift. When the API
// reports no tags the prior null-ness is preserved (omitted stays null, an
// explicit {} stays an explicit empty map); the prior values are never copied
// back, so a non-empty map cleared out-of-band is not masked as still-present.
func buildPortRangeListState(ctx context.Context, client *connectivity.AliyunClient, id string, priorTags types.Map) (*portRangeListResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	object, err := describePortRangeList(client, id)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Reading %s: calling DescribePortRangeLists", portRangeListTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribePortRangeLists", portRangeListTypeName).Error(),
		)
		return nil, diags
	}
	if object == nil {
		return nil, diags
	}

	entries, err := describePortRangeListEntries(client, id)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Reading %s: calling DescribePortRangeListEntries", portRangeListTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribePortRangeListEntries", portRangeListTypeName).Error(),
		)
		return nil, diags
	}

	model := &portRangeListResourceModel{
		Id:                types.StringValue(id),
		MaxEntries:        int64Value(object["MaxEntries"]),
		PortRangeListName: anyStringValue(object["PortRangeListName"]),
		Description:       anyStringValue(object["Description"]),
		ResourceGroupId:   anyStringValue(object["ResourceGroupId"]),
		AssociationCount:  anyStringValue(object["AssociationCount"]),
		CreateTime:        anyStringValue(object["CreationTime"]),
	}

	entryObjects := make([]entryModel, 0, len(entries))
	for _, raw := range entries {
		e, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		entryObjects = append(entryObjects, entryModel{
			PortRange:   anyStringValue(e["PortRange"]),
			Description: anyStringValue(e["Description"]),
		})
	}
	if len(entryObjects) > 0 {
		entrySet, d := types.SetValueFrom(ctx, entryObjectType, entryObjects)
		diags.Append(d...)
		model.Entry = entrySet
	} else {
		// No entries: keep the set null so it matches a config without entry
		// blocks (avoids a perpetual null-vs-empty-set diff and ImportStateVerify).
		model.Entry = types.SetNull(entryObjectType)
	}

	tags := tagsToMap(object["Tags"])
	switch {
	case len(tags) > 0:
		// The API reports real tags: refresh to the actual values instead of
		// copying the plan, so an out-of-band change is surfaced as a diff.
		tagMap, d := types.MapValueFrom(ctx, types.StringType, tags)
		diags.Append(d...)
		model.Tags = tagMap
	case priorTags.IsNull():
		// No API tags and the config omitted tags: keep the attribute null so
		// it matches a config without a tags block.
		model.Tags = types.MapNull(types.StringType)
	default:
		// No API tags but the config set tags explicitly (possibly tags = {}).
		// Preserve an explicit empty map instead of rewriting it to null: a
		// null would make an explicit tags = {} drift. The prior values are
		// not copied back — the API truth is "no tags", surfaced as an empty
		// map rather than null so the configured non-null shape is retained.
		emptyTags, d := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		diags.Append(d...)
		model.Tags = emptyTags
	}

	if diags.HasError() {
		return nil, diags
	}
	return model, diags
}

// entriesToRequest flattens a planned entry set into the []interface{} array of
// {PortRange, Description} maps that CreatePortRangeList expects; RpcPost
// expands the array into Entry.N.PortRange / Entry.N.Description.
func entriesToRequest(ctx context.Context, set types.Set) ([]interface{}, diag.Diagnostics) {
	var entries []entryModel
	diags := set.ElementsAs(ctx, &entries, false)
	array := make([]interface{}, 0, len(entries))
	for _, e := range entries {
		array = append(array, map[string]interface{}{
			"PortRange":   e.PortRange.ValueString(),
			"Description": e.Description.ValueString(),
		})
	}
	return array, diags
}

// diffEntries computes the AddEntry and RemoveEntry arrays for ModifyPortRangeList
// from the old and new entry sets. A newly added entry and an entry whose
// description changed are both sent via AddEntry; the changed/removed entry is
// also sent via RemoveEntry (re-added above with its new description), matching
// the SDKv2 Update behavior.
func diffEntries(ctx context.Context, oldSet, newSet types.Set) (addEntry, removeEntry []interface{}, d diag.Diagnostics) {
	var oldEntries, newEntries []entryModel
	d.Append(oldSet.ElementsAs(ctx, &oldEntries, false)...)
	d.Append(newSet.ElementsAs(ctx, &newEntries, false)...)
	if d.HasError() {
		return nil, nil, d
	}
	oldMap := make(map[string]string, len(oldEntries))
	for _, e := range oldEntries {
		oldMap[e.PortRange.ValueString()] = e.Description.ValueString()
	}
	newMap := make(map[string]string, len(newEntries))
	for _, e := range newEntries {
		newMap[e.PortRange.ValueString()] = e.Description.ValueString()
	}
	addEntry = make([]interface{}, 0)
	for pr, desc := range newMap {
		if oldDesc, ok := oldMap[pr]; !ok || oldDesc != desc {
			addEntry = append(addEntry, map[string]interface{}{
				"PortRange":   pr,
				"Description": desc,
			})
		}
	}
	removeEntry = make([]interface{}, 0)
	for pr, desc := range oldMap {
		if newDesc, ok := newMap[pr]; !ok || newDesc != desc {
			removeEntry = append(removeEntry, map[string]interface{}{
				"PortRange": pr,
			})
		}
	}
	return addEntry, removeEntry, d
}

// anyStringValue maps a response field to a String. Empty strings and nil
// become null (absent) rather than "", matching the Framework idiom and the
// SDKv2 Computed-string behavior; numbers and other scalars are stringified so
// fields like AssociationCount survive as String.
func anyStringValue(v interface{}) types.String {
	if v == nil {
		return types.StringNull()
	}
	if s, ok := v.(string); ok {
		if s == "" {
			return types.StringNull()
		}
		return types.StringValue(s)
	}
	s := fmt.Sprint(v)
	if s == "" || s == "<nil>" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// int64Value maps a response number field to an Int64. POP responses decode
// numbers as float64, int, json.Number, or string; all are handled.
func int64Value(v interface{}) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	switch n := v.(type) {
	case int:
		return types.Int64Value(int64(n))
	case int64:
		return types.Int64Value(n)
	case float64:
		return types.Int64Value(int64(n))
	case string:
		if i, err := strconv.ParseInt(n, 10, 64); err == nil {
			return types.Int64Value(i)
		}
		return types.Int64Null()
	default:
		if i, err := strconv.ParseInt(fmt.Sprint(v), 10, 64); err == nil {
			return types.Int64Value(i)
		}
		return types.Int64Null()
	}
}
