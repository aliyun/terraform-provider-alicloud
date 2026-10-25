// Package tag. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package tag

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

const createdByTypeName = "alicloud_tag_created_by"

// The CreatedBy switch is account-scoped: OpenCreatedBy, CheckCreatedByEnabled and
// CloseCreatedBy all require RegionId cn-shanghai and are served only from that
// endpoint. Calling through the provider's regional endpoint fails with
// WrongEndpoint, so the resource pins the control plane instead of the provider
// region and persists the control region as the resource ID.
const createdByControlRegion = "cn-shanghai"

func createdByEndpoint() string {
	return connectivity.LoadRegionalEndpoint(createdByControlRegion, "tag")
}

var (
	_ resource.Resource                = &createdByResource{}
	_ resource.ResourceWithImportState = &createdByResource{}
)

func NewCreatedByResource() resource.Resource {
	return &createdByResource{}
}

type createdByResource struct {
	fwadapt.ResourceBase
}

type createdByResourceModel struct {
	Id         types.String `tfsdk:"id"`
	FwTimeouts types.Object `tfsdk:"timeouts"`
	OpenStatus types.Bool   `tfsdk:"open_status"`
}

func (r *createdByResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = createdByTypeName
}

func (r *createdByResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"open_status": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the createdby tag is enabled. Valid values:\n\n*   true\n*   false\n",
			},
		},
	}
}

func (r *createdByResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan createdByResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "create", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"ClientToken": helper.BuildClientToken("OpenCreatedBy"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	request["RegionId"] = createdByControlRegion
	createValues := map[string]interface{}{"RegionId": createdByControlRegion}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(ctx, r.Client(), "Tag", "2018-08-28", "OpenCreatedBy", request, operationTimeout, createdByEndpoint())
	// The switch is a singleton: an already-enabled account is the desired state, and
	// OpenCreatedBy reports that as CreatedByAlreadyOpen instead of succeeding.
	if err != nil && !isExpectedErrors(err, []string{"CreatedByAlreadyOpen"}) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling OpenCreatedBy", createdByTypeName),
			errs.WrapErrorf(err, "Creating %s: calling OpenCreatedBy", createdByTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(request, fwMappingPath("RegionId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_tag_created_by: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling OpenCreatedBy", createdByTypeName),
			fmt.Sprintf("The response carried no RegionId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.OpenStatus = fwKnownValue(trackedPlan.OpenStatus).(types.Bool)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := createdByWaitForState(ctx, r.Client(), id, "OpenStatus", []string{"true"}, []string{}, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for OpenCreatedBy", createdByTypeName),
			errs.WrapErrorf(err, "Creating %s: waiting for OpenCreatedBy (%s)", createdByTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", createdByTypeName),
			fmt.Sprintf("The resource %s is not visible via CheckCreatedByEnabled right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, createdByState(object, &plan))...)
}

func (r *createdByResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state createdByResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := createdByGet(ctx, r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling CheckCreatedByEnabled", createdByTypeName),
			errs.WrapErrorf(err, "Reading %s: calling CheckCreatedByEnabled", createdByTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, createdByState(object, &state))...)
}

func (r *createdByResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state createdByResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only the local timeouts argument can change: the switch is account-scoped and exposes no
	// mutable attribute. open_status is Computed, so a timeouts-only change leaves it unknown in
	// the plan; read the live switch back instead of persisting that unknown value.
	object, err := createdByGet(ctx, r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling CheckCreatedByEnabled", createdByTypeName),
			errs.WrapErrorf(err, "Updating %s: calling CheckCreatedByEnabled", createdByTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, createdByState(object, &plan))...)
}

func (r *createdByResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state createdByResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "delete", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The switch is account-scoped, so deletion always targets the control region;
	// a state ID left over from an import with another value must not decide it.
	deleteValues := map[string]interface{}{"RegionId": createdByControlRegion}
	selected := false
	deleteEligible0, deleteConditionErr0 := fwDeleteEligible(deleteValues, []fwDeleteCondition{})
	if deleteConditionErr0 != nil {
		resp.Diagnostics.AddError("Selecting deletion operation", deleteConditionErr0.Error())
		return
	}
	selected = selected || deleteEligible0
	if !selected {
		resp.Diagnostics.AddError("Selecting deletion operation", "No deletion operation matches the resource state; the resource remains tracked")
		return
	}
	if deleteEligible0 {
		request0 := map[string]interface{}{"ClientToken": helper.BuildClientToken("CloseCreatedBy")}
		request0["RegionId"] = createdByControlRegion
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building CloseCreatedBy request", err.Error())
			return
		}
		if _, err := retryRpcPost(ctx, r.Client(), "Tag", "2018-08-28", "CloseCreatedBy", request0, operationTimeout, createdByEndpoint()); err != nil && !isExpectedErrors(err, []string{"ResourceNotFound.CreatedBy", "ResourceNotFound"}) {
			resp.Diagnostics.AddError("Deleting alicloud_tag_created_by: CloseCreatedBy", err.Error())
			return
		}
		// A nil target waits for absence: StateChangeConf only treats a nil refresh
		// result as success when Target is empty, and createdByGet returns nil once the
		// switch reports OpenStatus=false.
		if _, err := createdByWaitForState(ctx, r.Client(), state.Id.ValueString(), "OpenStatus", nil, []string{}, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Deleting alicloud_tag_created_by: waiting for CloseCreatedBy", err.Error())
			return
		}
	}
}

func (r *createdByResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// createdByState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func createdByState(object map[string]interface{}, prior *createdByResourceModel) createdByResourceModel {
	// Get has already merged and normalized every response mapping.
	boolVal := func(key string) types.Bool {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.BoolNull()
			}
			cur, ok = m[seg]
			if !ok {
				return types.BoolNull()
			}
		}
		if v, ok := cur.(bool); ok {
			return types.BoolValue(v)
		}
		return types.BoolNull()
	}
	id := types.StringNull()
	if v, ok := object["RegionId"].(string); ok && v != "" {
		id = types.StringValue(v)
	} else if prior != nil && !prior.Id.IsNull() && !prior.Id.IsUnknown() {
		id = prior.Id
	}
	m := createdByResourceModel{
		FwTimeouts: types.ObjectNull(fwTimeoutTypes()),
		Id:         id,
		OpenStatus: boolVal("OpenStatus"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
	}
	return m
}
