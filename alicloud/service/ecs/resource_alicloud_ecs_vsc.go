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
