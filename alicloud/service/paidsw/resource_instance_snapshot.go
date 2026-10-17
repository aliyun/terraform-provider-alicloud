package paidsw

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const instanceSnapshotTypeName = "alicloud_pai_dsw_api_instance_snapshot"

// PAI DSW (Data Science Workshop) exposes a RESTful/ROA API. "pai-dsw" is the
// gateway product code; the connectivity layer kebab→snake converts and
// lower-cases it to "pai_dsw" and resolves the endpoint via the Alibaba
// location service (no static endpoint table entry), the same path every
// other GA POP product without a pinned endpoint takes.
const (
	paiDswProduct = "pai-dsw"
	paiDswVersion = "2022-01-01"
)

var (
	_ resource.Resource                = &instanceSnapshotResource{}
	_ resource.ResourceWithImportState = &instanceSnapshotResource{}
)

func NewInstanceSnapshotResource() resource.Resource {
	return &instanceSnapshotResource{}
}

type instanceSnapshotResource struct {
	fwadapt.ResourceBase
}

// labelObjectType is the Framework object type for a single snapshot label
// {key, value}, used as the ElementType of the labels ListAttribute below.
// It is needed to rebuild a types.List from the Get response and to keep
// plan/state element types aligned.
//
// labels is modeled as schema.ListAttribute (a plain typed attribute) rather
// than schema.ListNestedAttribute deliberately: this provider serves Protocol
// V5 (see ProtoV5ProviderFactories in the AccTests), and the framework's V5
// schema converter rejects every fwschema.NestedAttribute with
// `protocol version 5 cannot have Attributes set`
// (terraform-plugin-framework internal/toproto5/schema_attribute.go). A
// ListAttribute with an object ElementType serializes its Type via the normal
// V5 path and round-trips {key, value} objects without per-field Required
// enforcement; the API still rejects a label missing its key.
var labelObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

type instanceSnapshotLabelModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type instanceSnapshotResourceModel struct {
	// Id is the composite "<InstanceId>:<SnapshotId>". Get and Delete need
	// both segments, so unlike most single-key resources the ID carries two.
	Id types.String `tfsdk:"id"`

	// Inputs. The PAI DSW snapshot API has no Update operation, so every
	// input attribute forces replacement.
	InstanceId          types.String `tfsdk:"instance_id"`
	SnapshotName        types.String `tfsdk:"snapshot_name"`
	SnapshotDescription types.String `tfsdk:"snapshot_description"`
	ImageUrl            types.String `tfsdk:"image_url"`
	Overwrite           types.Bool   `tfsdk:"overwrite"`
	Labels              types.List   `tfsdk:"labels"`
	ExcludePaths        types.List   `tfsdk:"exclude_paths"`

	// Computed outputs from GetInstanceSnapshot.
	SnapshotId      types.String `tfsdk:"snapshot_id"`
	Status          types.String `tfsdk:"status"`
	ImageId         types.String `tfsdk:"image_id"`
	ReasonCode      types.String `tfsdk:"reason_code"`
	ReasonMessage   types.String `tfsdk:"reason_message"`
	GmtCreateTime   types.String `tfsdk:"gmt_create_time"`
	GmtModifiedTime types.String `tfsdk:"gmt_modified_time"`
}

func (r *instanceSnapshotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = instanceSnapshotTypeName
}

func (r *instanceSnapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// SnapshotDescription and Overwrite are accepted by Create but are NOT
	// returned by GetInstanceSnapshot. They are create-only inputs, so they
	// require replacement and their values are preserved from the
	// plan/state during Read rather than refreshed from the API (refreshing
	// them to null would produce a permanent post-apply diff).
	resp.Schema = schema.Schema{
		Description: "Provides a PAI DSW (Data Science Workshop) API Instance Snapshot resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID in the form `<instance_id>:<snapshot_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the PAI DSW instance the snapshot belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_name": schema.StringAttribute{
				Required:    true,
				Description: "The snapshot name. Only lowercase letters, digits and hyphens are supported.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9-]+$`),
						"must only contain lowercase letters, digits and hyphens"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the snapshot. Only set at creation; not returned by Get.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_url": schema.StringAttribute{
				Required:    true,
				Description: "The image URL used to create the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"overwrite": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to overwrite an existing snapshot with the same name. Only set at creation; not returned by Get.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"labels": schema.ListAttribute{
				ElementType: labelObjectType,
				Optional:    true,
				Description: "The labels attached to the snapshot. Each element is an object with a key and an optional value.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"exclude_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The paths excluded from the snapshot.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"image_id": schema.StringAttribute{
				Computed:    true,
				Description: "The image ID of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reason_code": schema.StringAttribute{
				Computed:    true,
				Description: "The failure reason code, populated when the snapshot fails.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reason_message": schema.StringAttribute{
				Computed:    true,
				Description: "The failure reason message, populated when the snapshot fails.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The last modification time of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *instanceSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan instanceSnapshotResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	labels, err := listToLabels(ctx, plan.Labels)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s: reading labels", instanceSnapshotTypeName), err.Error())
		return
	}
	excludePaths, err := listToExcludePaths(ctx, plan.ExcludePaths)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s: reading exclude_paths", instanceSnapshotTypeName), err.Error())
		return
	}

	instanceId := plan.InstanceId.ValueString()
	body := map[string]interface{}{
		"ImageUrl":     plan.ImageUrl.ValueString(),
		"SnapshotName": plan.SnapshotName.ValueString(),
	}
	if !plan.SnapshotDescription.IsNull() {
		body["SnapshotDescription"] = plan.SnapshotDescription.ValueString()
	}
	if !plan.Overwrite.IsNull() {
		body["Overwrite"] = plan.Overwrite.ValueBool()
	}
	if len(labels) > 0 {
		body["Labels"] = labels
	}
	if len(excludePaths) > 0 {
		body["ExcludePaths"] = excludePaths
	}

	action := fmt.Sprintf("/api/v2/instances/%s/snapshots", instanceId)
	response, err := r.Client().RoaPost(paiDswProduct, paiDswVersion, action, nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateInstanceSnapshot", instanceSnapshotTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateInstanceSnapshot for instance %s", instanceSnapshotTypeName, instanceId).Error(),
		)
		return
	}

	snapshotId, _ := response["SnapshotId"].(string)
	if snapshotId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateInstanceSnapshot", instanceSnapshotTypeName),
			fmt.Sprintf("The CreateInstanceSnapshot response carried no SnapshotId: %v", response),
		)
		return
	}

	// The snapshot exists in the cloud from this point on. Persist its identity
	// before the read-back: Terraform keeps a non-null state returned together
	// with an error, so a failing or lagging GetInstanceSnapshot must not leave
	// the snapshot untracked (which would orphan it and make the next apply
	// create a duplicate). Read/Delete only need the composite ID, so the
	// resource stays manageable and destroyable from this state alone.
	identity := plan
	identity.Id = types.StringValue(fmt.Sprintf("%s:%s", instanceId, snapshotId))
	identity.SnapshotId = types.StringValue(snapshotId)
	identity.Status = types.StringNull()
	identity.ImageId = types.StringNull()
	identity.ReasonCode = types.StringNull()
	identity.ReasonMessage = types.StringNull()
	identity.GmtCreateTime = types.StringNull()
	identity.GmtModifiedTime = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getInstanceSnapshot(r.Client(), instanceId, snapshotId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetInstanceSnapshot", instanceSnapshotTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetInstanceSnapshot for instance %s snapshot %s", instanceSnapshotTypeName, instanceId, snapshotId).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateInstanceSnapshot", instanceSnapshotTypeName),
			fmt.Sprintf("The snapshot %s:%s is not visible via GetInstanceSnapshot right after creation", instanceId, snapshotId),
		)
		return
	}

	state, err := instanceSnapshotState(ctx, object, &plan)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s: mapping state", instanceSnapshotTypeName), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *instanceSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state instanceSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, snapshotId, err := parseCompositeId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s", instanceSnapshotTypeName), err.Error())
		return
	}

	object, err := getInstanceSnapshot(r.Client(), instanceId, snapshotId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetInstanceSnapshot", instanceSnapshotTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetInstanceSnapshot for instance %s snapshot %s", instanceSnapshotTypeName, instanceId, snapshotId).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	updated, err := instanceSnapshotState(ctx, object, &state)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Reading %s: mapping state", instanceSnapshotTypeName), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &updated)...)
}

// Update satisfies the resource.Resource interface but is never invoked at
// runtime: the PAI DSW snapshot API exposes no Update operation, and every
// input attribute carries RequiresReplace, so any planned change forces
// Delete + Create rather than an in-place update. The body is intentionally
// empty — if the framework ever routed here, leaving state untouched would be
// the only safe option for an un-updateable resource, and RequiresReplace
// guarantees that path is unreachable.
func (r *instanceSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete removes the snapshot. A not-found on delete is treated as success so
// a re-run after an out-of-band deletion converges instead of erroring.
func (r *instanceSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state instanceSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, snapshotId, err := parseCompositeId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Deleting %s", instanceSnapshotTypeName), err.Error())
		return
	}

	action := fmt.Sprintf("/api/v2/instances/%s/snapshots/%s", instanceId, snapshotId)
	if _, err := r.Client().RoaDelete(paiDswProduct, paiDswVersion, action, nil, nil, nil, true); err != nil {
		if !isSnapshotNotExist(err) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Deleting %s: calling DeleteInstanceSnapshot", instanceSnapshotTypeName),
				errs.WrapErrorf(err, "Deleting %s: calling DeleteInstanceSnapshot for instance %s snapshot %s", instanceSnapshotTypeName, instanceId, snapshotId).Error(),
			)
		}
	}
}

func (r *instanceSnapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	instanceId, snapshotId, err := parseCompositeId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Importing %s", instanceSnapshotTypeName), err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), types.StringValue(instanceId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("snapshot_id"), types.StringValue(snapshotId))...)
}

// getInstanceSnapshot fetches a snapshot and returns the response body, or nil
// when the snapshot no longer exists.
func getInstanceSnapshot(client *connectivity.AliyunClient, instanceId, snapshotId string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/api/v2/instances/%s/snapshots/%s", instanceId, snapshotId)
	response, err := client.RoaGet(paiDswProduct, paiDswVersion, action, nil, nil, nil)
	if err != nil {
		if isSnapshotNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// parseCompositeId splits "<InstanceId>:<SnapshotId>".
func parseCompositeId(id string) (string, string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected the resource ID in the form <instance_id>:<snapshot_id>, got %q", id)
	}
	return parts[0], parts[1], nil
}

// instanceSnapshotState maps a GetInstanceSnapshot response onto the resource
// model. preserved carries the create-only inputs (SnapshotDescription,
// Overwrite) which GetInstanceSnapshot does not return, so their values are
// kept from the prior plan/state instead of being nulled out.
func instanceSnapshotState(ctx context.Context, object map[string]interface{}, preserved *instanceSnapshotResourceModel) (instanceSnapshotResourceModel, error) {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}

	var snapshotDescription types.String
	var overwrite types.Bool
	var preservedLabels, preservedExcludePaths types.List
	if preserved != nil {
		snapshotDescription = preserved.SnapshotDescription
		overwrite = preserved.Overwrite
		preservedLabels = preserved.Labels
		preservedExcludePaths = preserved.ExcludePaths
	}

	labels, err := labelsToList(ctx, object["Labels"], preservedLabels)
	if err != nil {
		return instanceSnapshotResourceModel{}, err
	}
	excludePaths, err := excludePathsToList(ctx, object["ExcludePaths"], preservedExcludePaths)
	if err != nil {
		return instanceSnapshotResourceModel{}, err
	}

	instanceId := str("InstanceId")
	snapshotId := str("SnapshotId")
	var id types.String
	if !instanceId.IsNull() && !snapshotId.IsNull() {
		id = types.StringValue(fmt.Sprintf("%s:%s", instanceId.ValueString(), snapshotId.ValueString()))
	} else {
		id = types.StringNull()
	}

	return instanceSnapshotResourceModel{
		Id:                  id,
		InstanceId:          instanceId,
		SnapshotName:        str("SnapshotName"),
		SnapshotDescription: snapshotDescription,
		ImageUrl:            str("ImageUrl"),
		Overwrite:           overwrite,
		Labels:              labels,
		ExcludePaths:        excludePaths,
		SnapshotId:          snapshotId,
		Status:              str("Status"),
		ImageId:             str("ImageId"),
		ReasonCode:          str("ReasonCode"),
		ReasonMessage:       str("ReasonMessage"),
		GmtCreateTime:       str("GmtCreateTime"),
		GmtModifiedTime:     str("GmtModifiedTime"),
	}, nil
}

// listToLabels converts a plan/state Labels list into the request body form
// []map[string]interface{}{{"Key": .., "Value": ..}}.
func listToLabels(ctx context.Context, list types.List) ([]map[string]interface{}, error) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}
	var in []instanceSnapshotLabelModel
	if diags := list.ElementsAs(ctx, &in, false); diags.HasError() {
		return nil, fmt.Errorf("reading labels from plan: %v", diags)
	}
	out := make([]map[string]interface{}, 0, len(in))
	for _, l := range in {
		entry := map[string]interface{}{"Key": l.Key.ValueString()}
		if !l.Value.IsNull() {
			entry["Value"] = l.Value.ValueString()
		}
		out = append(out, entry)
	}
	return out, nil
}

// listToExcludePaths converts a plan/state ExcludePaths list into []string for
// the request body.
func listToExcludePaths(ctx context.Context, list types.List) ([]string, error) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}
	var out []string
	if diags := list.ElementsAs(ctx, &out, false); diags.HasError() {
		return nil, fmt.Errorf("reading exclude_paths from plan: %v", diags)
	}
	return out, nil
}

// labelsToList rebuilds a types.List from a Get response Labels array. When the
// response carries no elements it keeps the configured or prior-state
// collection instead of forcing null: labels is Optional and not Computed, so
// turning an explicit [] into null would surface a diff (and, with
// RequiresReplace, a replacement) on the next plan.
func labelsToList(ctx context.Context, raw interface{}, preserved types.List) (types.List, error) {
	arr, _ := raw.([]interface{})
	if len(arr) == 0 {
		return preservedOrNull(preserved, labelObjectType), nil
	}
	out := make([]instanceSnapshotLabelModel, 0, len(arr))
	for _, item := range arr {
		m, _ := item.(map[string]interface{})
		out = append(out, instanceSnapshotLabelModel{
			Key:   objectString(m["Key"]),
			Value: objectString(m["Value"]),
		})
	}
	list, diags := types.ListValueFrom(ctx, labelObjectType, out)
	if diags.HasError() {
		return types.ListNull(labelObjectType), fmt.Errorf("converting labels to list: %v", diags)
	}
	return list, nil
}

// excludePathsToList rebuilds a types.List from a Get response ExcludePaths
// array of strings, with the same empty-vs-null preservation as labelsToList.
func excludePathsToList(ctx context.Context, raw interface{}, preserved types.List) (types.List, error) {
	arr, _ := raw.([]interface{})
	if len(arr) == 0 {
		return preservedOrNull(preserved, types.StringType), nil
	}
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return preservedOrNull(preserved, types.StringType), nil
	}
	list, diags := types.ListValueFrom(ctx, types.StringType, out)
	if diags.HasError() {
		return types.ListNull(types.StringType), fmt.Errorf("converting exclude_paths to list: %v", diags)
	}
	return list, nil
}

// preservedOrNull returns the configured/prior-state collection when the API
// returned no elements, and a typed null otherwise (unset input, or a state
// written before this attribute was configured).
func preservedOrNull(preserved types.List, elemType attr.Type) types.List {
	if !preserved.IsNull() && !preserved.IsUnknown() {
		return preserved
	}
	return types.ListNull(elemType)
}

// objectString returns a non-null types.String for a non-empty string value
// from a map[string]interface{} API object, else types.StringNull().
func objectString(v interface{}) types.String {
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}

// isSnapshotNotExist reports whether err means the snapshot is gone. The ROA
// client surfaces 4xx as a *tea.SDKError carrying the HTTP status, which
// errs.NotFoundError catches via StatusCode == 404; the substring fallback
// covers any wrapped message variant the generic detector misses.
func isSnapshotNotExist(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not found") || strings.Contains(msg, "not exist")
}
