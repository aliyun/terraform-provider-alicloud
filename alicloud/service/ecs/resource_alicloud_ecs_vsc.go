// Package ecs implements the Framework-native ECS resources.
//
// alicloud_ecs_vsc is a Framework-native managed resource (terraform-plugin-
// framework), registered through provider/registry and served by the mux
// provider alongside the SDKv2 ECS resources. It rewrites the earlier SDKv2
// implementation (PR #10815) on the Framework: the CreateVsc/DescribeVscs/
// DeleteVsc call shapes and the Tag.N.Key/Tag.N.Value + system-tag filtering
// semantics are preserved, and tags gain an in-place Update path the SDKv2
// ForceNew implementation lacked. The tags Update uses ResourceType="vsc" for
// TagResources/UntagResources, which the live ACC confirms accept "vsc"
// (the global tagging APIs' public enum is non-exhaustive). resource_group_id
// stays RequiresReplace, as on the SDKv2 line: CreateVsc accepts ResourceGroupId
// at create, but JoinResourceGroup's ResourceType enum does not include "vsc"
// — confirmed live (404 InvalidResourceType.NotFound; the enum is
// instance/disk/snapshot/image/securitygroup/ddh/ddhcluster/eni/keypair/
// launchtemplate/command/activation/managedinstance), so an in-place
// resource-group move is unsupported and a change forces replace.
package ecs

import (
	"context"
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ecsVscTypeName = "alicloud_ecs_vsc"

var (
	_ resource.Resource                = &ecsVscResource{}
	_ resource.ResourceWithImportState = &ecsVscResource{}
)

func NewEcsVscResource() resource.Resource {
	return &ecsVscResource{}
}

type ecsVscResource struct {
	fwadapt.ResourceBase
}

type ecsVscResourceModel struct {
	Id              types.String `tfsdk:"id"`
	InstanceId      types.String `tfsdk:"instance_id"`
	Description     types.String `tfsdk:"description"`
	ResourceGroupId types.String `tfsdk:"resource_group_id"`
	Status          types.String `tfsdk:"status"`
	Tags            types.Map    `tfsdk:"tags"`
	VscName         types.String `tfsdk:"vsc_name"`
	VscType         types.String `tfsdk:"vsc_type"`
}

func (r *ecsVscResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ecsVscTypeName
}

func (r *ecsVscResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an ECS Vsc resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Vsc.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the ECS instance to which the Vsc belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the Vsc.",
				// No ModifyVsc/UpdateVsc exists, so description changes replace.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the resource group to which the Vsc belongs.",
				// Settable on Create (CreateVsc takes ResourceGroupId) but not
				// movable in-place: JoinResourceGroup's ResourceType enum does not
				// include "vsc" (404 InvalidResourceType.NotFound), so a change
				// forces replace. UseStateForUnknown keeps the server-supplied value
				// when the field is omitted from config.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the Vsc.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vsc_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the Vsc.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vsc_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the Vsc.",
				// Settable on Create; the server defaults it (e.g. "Primary") when
				// omitted, so it is Computed — an unset value is Unknown in the
				// plan and accepts the server value after Read. Otherwise Framework
				// flags planned null -> applied "Primary" as an inconsistent result
				// after apply. No ModifyVsc exists, so a known change replaces.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			// tags is a MapAttribute (map[string]string) mirroring the alicloud
			// SDKv2 tags convention (schema.TypeMap, HCL tags = { key = "value" }).
			// A map attribute is protocol-v5-compatible: the resource is served by
			// the ProtoV5 mux (Framework v6 -> tf6to5server -> v5), and nested
			// *Attributes* (SetNestedAttribute's NestedObject.Attributes) cannot
			// round-trip through v5 ("protocol version 5 cannot have Attributes
			// set"). A flat MapAttribute has no nested Attributes, so it converts
			// cleanly. The Tag.N.Key/Tag.N.Value CreateVsc shape and the
			// TagResources/UntagResources Update diff are built from the decoded
			// map in the unit helpers, so the API mapping is unchanged.
			"tags": schema.MapAttribute{
				Optional:    true,
				Description: "A mapping of tags to assign to the Vsc.",
				ElementType: types.StringType,
			},
		},
	}
}

func (r *ecsVscResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ecsVscResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	action := "CreateVsc"
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken(action),
		"InstanceId":  plan.InstanceId.ValueString(),
	}
	// Optional/Computed values must not be sent when null or unknown:
	// ValueString on an unknown would yield "" and poison the request.
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request["Description"] = plan.Description.ValueString()
	}
	if !plan.ResourceGroupId.IsNull() && !plan.ResourceGroupId.IsUnknown() {
		request["ResourceGroupId"] = plan.ResourceGroupId.ValueString()
	}
	if !plan.VscName.IsNull() && !plan.VscName.IsUnknown() {
		request["VscName"] = plan.VscName.ValueString()
	}
	if !plan.VscType.IsNull() && !plan.VscType.IsUnknown() {
		request["VscType"] = plan.VscType.ValueString()
	}
	if !plan.Tags.IsNull() {
		if plan.Tags.IsUnknown() {
			// Pre-change validation: an unknown tags value cannot be mapped to
			// the Tag.N request shape, so the CreateVsc call must not proceed.
			resp.Diagnostics.AddError(
				fmt.Sprintf("Creating the %s resource", ecsVscTypeName),
				"The tags argument is unknown; the CreateVsc request cannot be built before the plan is fully known.",
			)
			return
		}
		// CreateVsc takes tags inline as Tag.<i>.Key/Tag.<i>.Value, 1-based,
		// the same shape CreateVsc used on the SDKv2 line.
		tags := tagsToMap(ctx, &resp.Diagnostics, plan.Tags)
		if resp.Diagnostics.HasError() {
			return
		}
		for k, v := range vscCreateTagParams(tags) {
			request[k] = v
		}
	}

	response, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVsc", ecsVscTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateVsc", ecsVscTypeName).Error(),
		)
		return
	}

	// Strictly assert the returned VscId is a non-empty string: fmt.Sprintf's
	// "%v" would accept a number/map and synthesise a fake id.
	vscId, ok := response["VscId"].(string)
	if !ok || vscId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVsc", ecsVscTypeName),
			fmt.Sprintf("The CreateVsc response carried no string VscId: %v", response),
		)
		return
	}

	// Persist the created identity before reading: if DescribeVscs fails or
	// the Vsc is not yet visible (eventual consistency), the id is still in
	// state so Terraform can refresh it on the next apply instead of treating
	// a created resource as absent and recreating it.
	plan.Id = types.StringValue(vscId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	object, err := getEcsVsc(client, vscId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeVscs", ecsVscTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeVscs", ecsVscTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateVsc", ecsVscTypeName),
			fmt.Sprintf("The Vsc %s is not visible via DescribeVscs right after creation", vscId),
		)
		return
	}
	state, diags := ecsVscState(ctx, vscId, object)
	resp.Diagnostics.Append(diags...)
	// Preserve plan-known Optional values the API round-trip normalizes away.
	// tags={} (a known empty map) and explicit empty-string description/vsc_name
	// are legitimate plan values; ecsVscState returns MapNull/StringNull when
	// Describe carries no user tag or an empty/absent string, which would
	// diverge from the plan and trip Framework's "inconsistent result after
	// apply". This mirrors the plan retention Update already applies (see the
	// tagsChanged branch below). Omitted (null) values stay null — only known
	// (non-null, non-unknown) values are carried over, so the null-vs-empty
	// plan/state semantics are preserved.
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		state.Tags = plan.Tags
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		state.Description = plan.Description
	}
	if !plan.VscName.IsNull() && !plan.VscName.IsUnknown() {
		state.VscName = plan.VscName
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ecsVscResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ecsVscResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getEcsVsc(r.Client(), state.Id.ValueString())
	newState, diags, remove := applyReadResult(ctx, state, object, err)
	resp.Diagnostics.Append(diags...)
	if remove {
		// Legitimate NotFound: an empty DescribeVscs array (the Ecs 2014-05-26
		// contract's NotFound shape), not a matched entry missing Status.
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *ecsVscResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ecsVscResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	id := state.Id.ValueString()

	// description / vsc_name / vsc_type / instance_id / resource_group_id are
	// RequiresReplace, so an Update only runs for tags — the one field the
	// model marks updatable (TagResources/UntagResources). resource_group_id
	// cannot be moved in-place: JoinResourceGroup rejects ResourceType="vsc"
	// (404 InvalidResourceType.NotFound; its enum is instance/disk/snapshot/
	// image/securitygroup/ddh/ddhcluster/eni/keypair/launchtemplate/command/
	// activation/managedinstance), so a change forces replace, not Update.

	// ResourceType="vsc" below is the resource type the CloudSpec model
	// (ACS::Ecs::Vsc) binds to TagResources/UntagResources; the global tagging
	// APIs' public enum is non-exhaustive (see the package doc), and the live
	// ACC confirms "vsc" is accepted.

	tagsChanged := !plan.Tags.IsUnknown() && !plan.Tags.Equal(state.Tags)
	if tagsChanged {
		oldTags := tagsToMap(ctx, &resp.Diagnostics, state.Tags)
		newTags := tagsToMap(ctx, &resp.Diagnostics, plan.Tags)
		if resp.Diagnostics.HasError() {
			return
		}

		var removed []string
		added := map[string]string{}
		for k, oldV := range oldTags {
			if newV, ok := newTags[k]; !ok {
				removed = append(removed, k)
			} else if newV != oldV {
				// Changed value: re-tag overwrites the key, no untag needed.
				added[k] = newV
			}
		}
		for k, newV := range newTags {
			if _, ok := oldTags[k]; !ok {
				added[k] = newV
			}
		}

		if len(removed) > 0 {
			request := map[string]interface{}{
				"RegionId":     client.RegionId,
				"ResourceType": "vsc",
				"ResourceId.1": id,
			}
			for i, k := range removed {
				request[fmt.Sprintf("TagKey.%d", i+1)] = k
			}
			if _, err := client.RpcPost("Ecs", "2014-05-26", "UntagResources", nil, request, true); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s: calling UntagResources", ecsVscTypeName),
					errs.WrapErrorf(err, "Updating %s: calling UntagResources", ecsVscTypeName).Error(),
				)
				return
			}
		}
		if len(added) > 0 {
			request := map[string]interface{}{
				"RegionId":     client.RegionId,
				"ResourceType": "vsc",
				"ResourceId.1": id,
			}
			count := 1
			for k, v := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = k
				request[fmt.Sprintf("Tag.%d.Value", count)] = v
				count++
			}
			if _, err := client.RpcPost("Ecs", "2014-05-26", "TagResources", nil, request, true); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Updating %s: calling TagResources", ecsVscTypeName),
					errs.WrapErrorf(err, "Updating %s: calling TagResources", ecsVscTypeName).Error(),
				)
				return
			}
		}
	}

	object, err := getEcsVsc(client, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeVscs", ecsVscTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeVscs", ecsVscTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	newState, diags := ecsVscState(ctx, id, object)
	resp.Diagnostics.Append(diags...)
	// tags is a user-controlled Optional attribute; after a successful
	// TagResources/UntagResources the desired value is the plan's. Re-reading
	// immediately can race the API's eventual consistency and surface a stale
	// value Framework's strict plan-vs-applied check flags ("new element … has
	// appeared"). Trust the plan for the field just updated; the next Refresh
	// reconciles any real drift. resource_group_id is RequiresReplace, so it
	// never changes in Update and is read verbatim from DescribeVscs.
	if tagsChanged {
		newState.Tags = plan.Tags
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *ecsVscResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ecsVscResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "DeleteVsc"
	_, err := r.Client().RpcPost("Ecs", "2014-05-26", action, nil, map[string]interface{}{
		"VscId":       state.Id.ValueString(),
		"RegionId":    r.Client().RegionId,
		"ClientToken": helper.BuildClientToken(action),
	}, true)
	if err != nil && !errs.NotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteVsc", ecsVscTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteVsc", ecsVscTypeName).Error(),
		)
		return
	}
}

func (r *ecsVscResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getEcsVsc queries DescribeVscs for a single Vsc by ID and returns the Vsc
// object, or nil for a legitimate NotFound (an empty Vscs array, the Ecs
// 2014-05-26 contract's NotFound shape). A matched entry that omits Status is
// surfaced as an error by parseEcsVscResponse so Read emits a diagnostic and
// PRESERVES state instead of silently RemoveResource-ing a live, ID'd resource.
// RpcPost's autoRetry already retries throttling, so the SDKv2 retry.Retry/
// incrementalWait wrapper is not duplicated here. The response is parsed by
// parseEcsVscResponse, which is a pure function so the robustness rules
// (malformed/empty/no-match handling) can be exercised without the cloud.
func getEcsVsc(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := client.RpcPost("Ecs", "2014-05-26", "DescribeVscs", nil, map[string]interface{}{
		"VscIds.1": id,
		"RegionId": client.RegionId,
	}, true)
	if err != nil {
		return nil, err
	}
	return parseEcsVscResponse(response, id)
}

// parseEcsVscResponse decodes a DescribeVscs response for a single Vsc id.
// It returns the matched Vsc object, or nil for a legitimate NotFound (an
// empty Vscs array, the Ecs 2014-05-26 contract's NotFound shape). A matching
// entry that omits Status is NOT NotFound: it is surfaced as an error (see the
// Status-absent branch below) so Read preserves state instead of silently
// RemoveResource-ing a live, ID'd resource.
//
// Robustness: a malformed response (absent/null/non-array $.Vscs root) and a
// non-empty result whose entries do not contain the queried VscId are surfaced
// as errors rather than silently treated as NotFound, so a garbage response
// can never remove a live resource from state. The returned object is matched
// by VscId rather than blindly taking arr[0]: a sibling entry must never be
// mapped onto this resource.
func parseEcsVscResponse(response map[string]interface{}, id string) (map[string]interface{}, error) {
	// Query $.Vscs (not $.Vscs[*]): the wildcard form collapses an absent,
	// null, or non-array root into an empty slice with no error, which would
	// read as NotFound and silently RemoveResource a live entry. The bare key
	// distinguishes a missing root (error), a null root (nil), and a non-array
	// root (its own type) so each malformed shape is surfaced as an error.
	v, err := jsonpath.Get("$.Vscs", response)
	if err != nil {
		return nil, fmt.Errorf("DescribeVscs response missing the $.Vscs root: %w", err)
	}
	if v == nil {
		return nil, fmt.Errorf("DescribeVscs response has a null $.Vscs root")
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("DescribeVscs $.Vscs is not an array: %T", v)
	}
	if len(arr) == 0 {
		// Legitimate empty array: the Vsc no longer exists.
		return nil, nil
	}
	for _, raw := range arr {
		obj, ok := raw.(map[string]interface{})
		if !ok {
			// A non-object entry is malformed; keep scanning for a real match
			// before deciding the response is unusable.
			continue
		}
		if vid, _ := obj["VscId"].(string); vid != id {
			// A sibling entry; skip rather than misattribute it.
			continue
		}
		if obj["Status"] == nil {
			// A matched entry must carry a Status: DescribeVscs returns one per
			// the Ecs 2014-05-26 contract (enum In_use/Attaching/Detaching/
			// AttachFailed/DetachFailed), and a legitimate NotFound is an empty
			// Vscs array, not a present-but-Status-less entry. An entry matching
			// the queried VscId yet omitting Status is an inconsistent/malformed
			// response, so surface it as an error: Read then emits a diagnostic and
			// PRESERVES state instead of silently RemoveResource-ing a live, ID'd
			// resource. (No authoritative contract treats matched-ID+absent-Status
			// as NotFound; the prior nil,nil was an unverified legacy assumption
			// flagged in static review.)
			return nil, fmt.Errorf("DescribeVscs entry for VscId %q is missing Status: %v", id, obj)
		}
		return obj, nil
	}
	// Non-empty result but no entry whose VscId matches the queried id: the
	// response is inconsistent (DescribeVscs was queried with VscIds.1=id),
	// so surface it as an error rather than a silent NotFound.
	return nil, fmt.Errorf("DescribeVscs result does not contain the queried VscId %q", id)
}

// vscCreateTagParams renders the decoded tags map as the CreateVsc request shape
// Tag.<1-based i>.Key / Tag.<i>.Value. It is a pure function so the request
// shape can be verified without the cloud; the unknown-tags guard that
// precedes it in Create is plan-time validation and stays there. Map iteration
// order is unspecified, so the Tag.<i> index numbering is too — the API keys
// the pair by its Tag.<i>.Key/Tag.<i>.Value, not by position, so order is
// irrelevant; tests assert the pair invariant, not exact indices.
func vscCreateTagParams(tags map[string]string) map[string]string {
	params := make(map[string]string, len(tags)*2)
	i := 0
	for k, v := range tags {
		i++
		params[fmt.Sprintf("Tag.%d.Key", i)] = k
		params[fmt.Sprintf("Tag.%d.Value", i)] = v
	}
	return params
}

// ecsVscState maps a DescribeVscs Vsc object onto the resource model. Empty
// strings become null: optional fields the API omits read back as absent
// rather than "". The id is the caller-known VscId (d.Id on the SDKv2 line).
func ecsVscState(ctx context.Context, id string, object map[string]interface{}) (ecsVscResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}

	var tagsMap types.Map
	tagRaw, _ := jsonpath.Get("$.Tags[*]", object)
	if arr, ok := tagRaw.([]interface{}); ok && len(arr) > 0 {
		m := map[string]string{}
		for _, raw := range arr {
			t, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			key, _ := t["TagKey"].(string)
			// Skip Alibaba Cloud system tags (acs:*, aliyun, http(s)://*) so they
			// don't pollute state and cause perpetual drift; mirrors
			// EcsService.ecsTagIgnored.
			if strings.HasPrefix(key, "acs:") || strings.HasPrefix(key, "aliyun") || strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
				continue
			}
			// A map value is a string, never null, so the prior "value is
			// Required, never null in state" state-fidelity invariant holds
			// structurally — a map either carries a key with a string value or
			// omits the key. Preserve the API value verbatim, including the
			// legitimate empty string, and default an absent TagValue to "".
			if v, ok := t["TagValue"].(string); ok {
				m[key] = v
			} else {
				m[key] = ""
			}
		}
		if len(m) > 0 {
			var d diag.Diagnostics
			tagsMap, d = types.MapValueFrom(ctx, types.StringType, m)
			diags.Append(d...)
		} else {
			tagsMap = types.MapNull(types.StringType)
		}
	} else {
		// No tags on the API side: null (not an empty map) so a null config
		// does not drift into a no-op tag update.
		tagsMap = types.MapNull(types.StringType)
	}

	return ecsVscResourceModel{
		Id:              types.StringValue(id),
		InstanceId:      str("InstanceId"),
		Description:     str("Description"),
		ResourceGroupId: str("ResourceGroupId"),
		Status:          str("Status"),
		VscName:         str("VscName"),
		VscType:         str("VscType"),
		Tags:            tagsMap,
	}, diags
}

// tagsToMap decodes a tags map (plan/state) into a key->value map, the shape
// the CreateVsc Tag.<i>.Key/Tag.<i>.Value and the TagResources/UntagResources
// diff need. An unknown or null map yields an empty map (no tags); callers
// must check diagnostics before acting on the result.
func tagsToMap(ctx context.Context, diags *diag.Diagnostics, m types.Map) map[string]string {
	out := map[string]string{}
	if m.IsNull() || m.IsUnknown() {
		return out
	}
	diags.Append(m.ElementsAs(ctx, &out, false)...)
	return out
}

// applyReadResult decides, from a DescribeVscs fetch (object + err), how Read
// updates state. A fetch error — including a malformed response surfaced by
// parseEcsVscResponse, such as a matched entry missing Status — returns the
// prior state with a visible diagnostic and remove=false, so an ID'd resource is
// never silently forgotten. A nil object with no error is the legitimate NotFound
// (an empty Vscs array, the Ecs 2014-05-26 contract's NotFound shape) and removes
// the resource. A present object refreshes state. Separated from Read so the
// state-preservation rule is exercisable without a live client or a Framework
// ReadRequest.
func applyReadResult(ctx context.Context, prior ecsVscResourceModel, object map[string]interface{}, fetchErr error) (ecsVscResourceModel, diag.Diagnostics, bool) {
	var diags diag.Diagnostics
	if fetchErr != nil {
		diags.AddError(
			fmt.Sprintf("Reading %s: calling DescribeVscs", ecsVscTypeName),
			errs.WrapErrorf(fetchErr, "Reading %s: calling DescribeVscs", ecsVscTypeName).Error(),
		)
		return prior, diags, false
	}
	if object == nil {
		return prior, diags, true
	}
	newState, d := ecsVscState(ctx, prior.Id.ValueString(), object)
	diags.Append(d...)
	// Preserve the null-vs-empty distinction the API round-trip erases. ecsVscState
	// renders "no value on the API side" as null: a known empty tags map (tags={})
	// and an explicit empty description/vsc_name all decode to null when Describe
	// carries no user tag or an empty/absent string. That would diverge from a
	// prior state holding the known empty value and trip a perpetual refresh drift
	// — a tags={} config refreshes state tags to null and plans +tags={}. Carry
	// the prior empty value forward only when the freshly-read value is null AND
	// the prior value is a known empty value (empty map / empty string): a
	// non-empty prior that the API no longer reports is a real drift to surface,
	// not to mask; omitted (null) prior values stay null. This mirrors the plan
	// retention Create/Update apply for the same fields; import's prior state
	// lacks the empty map, so ImportStateVerify on tags={} still differs (see the
	// EmptyTags acceptance test comment).
	if newState.Tags.IsNull() && !prior.Tags.IsNull() && !prior.Tags.IsUnknown() && len(prior.Tags.Elements()) == 0 {
		newState.Tags = prior.Tags
	}
	if newState.Description.IsNull() && !prior.Description.IsNull() && !prior.Description.IsUnknown() && prior.Description.ValueString() == "" {
		newState.Description = prior.Description
	}
	if newState.VscName.IsNull() && !prior.VscName.IsNull() && !prior.VscName.IsUnknown() && prior.VscName.ValueString() == "" {
		newState.VscName = prior.VscName
	}
	return newState, diags, false
}
