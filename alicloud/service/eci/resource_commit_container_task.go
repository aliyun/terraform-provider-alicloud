// Package eci. This file is generated automatically. Please do not modify it manually, thank you!
// Nested blocks are used for list-of-objects fields so the schema works under
// protocol version 5 (the muxed provider server); nested attributes require
// protocol 6 and are not available through the v5 mux.
package eci

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const commitContainerTaskTypeName = "alicloud_eci_commit_container_task"

var (
	_ resource.Resource                = &commitContainerTaskResource{}
	_ resource.ResourceWithImportState = &commitContainerTaskResource{}
)

func NewCommitContainerTaskResource() resource.Resource {
	return &commitContainerTaskResource{}
}

type commitContainerTaskResource struct {
	fwadapt.ResourceBase
}

type commitContainerTaskResourceModel struct {
	Id               types.String `tfsdk:"id"`
	FwTimeouts       types.Object `tfsdk:"timeouts"`
	AcrRegistryInfo  types.List   `tfsdk:"acr_registry_info"`
	Arn              types.List   `tfsdk:"arn"`
	ContainerGroupId types.String `tfsdk:"container_group_id"`
	ContainerName    types.String `tfsdk:"container_name"`
	CreateTime       types.String `tfsdk:"create_time"`
	Image            types.List   `tfsdk:"image"`
}

func (r *commitContainerTaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = commitContainerTaskTypeName
}

func (r *commitContainerTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": fwTimeoutBlock(),
			"acr_registry_info": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"arn_service": schema.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Description: "The RAM role ARN of the user (authorized party) in cross-account authorization scenarios.",
						},
						"instance_id": schema.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Description: "ACR Enterprise Edition instance ID.",
						},
						"arn_user": schema.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Description: "The RAM role ARN of the authorizer in cross-account authorization scenarios.",
						},
						"region_id": schema.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Description: "The region of the ACR Enterprise Edition instance.",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIf(func(_ context.Context, req planmodifier.ListRequest, resp *listplanmodifier.RequiresReplaceIfFuncResponse) {
						resp.RequiresReplace = fwConfiguredChanged(req.ConfigValue, req.PlanValue, req.StateValue)
					}, "Replace when configured nested values change", "Replace when configured nested values change"),
				},
				Description: "Access credential configuration information for ACR Enterprise Edition instance.",
			},
			"arn": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"role_type": schema.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Description: "Authorization type. Can be configured as service, indicating using RAM role for authorization.",
						},
						"role_arn": schema.StringAttribute{
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Description: "The ARN of the authorization role.",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIf(func(_ context.Context, req planmodifier.ListRequest, resp *listplanmodifier.RequiresReplaceIfFuncResponse) {
						resp.RequiresReplace = fwConfiguredChanged(req.ConfigValue, req.PlanValue, req.StateValue)
					}, "Replace when configured nested values change", "Replace when configured nested values change"),
				},
				Description: "ARN information required for authorization.",
			},
			"image": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"repository": schema.StringAttribute{
							Sensitive: true,
							Required:  true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Description: "ACR image repository address.",
						},
						"message": schema.StringAttribute{
							Sensitive: true,
							Optional:  true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Description: "Image description message.",
						},
						"author": schema.StringAttribute{
							Sensitive: true,
							Optional:  true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Description: "Image author information.",
						},
						"tag": schema.StringAttribute{
							Sensitive: true,
							Required:  true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Description: "Image tag. Required; an empty string is accepted by the API and means no tag modification.",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(1),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIf(func(_ context.Context, req planmodifier.ListRequest, resp *listplanmodifier.RequiresReplaceIfFuncResponse) {
						resp.RequiresReplace = fwConfiguredChanged(req.ConfigValue, req.PlanValue, req.StateValue)
					}, "Replace when configured nested values change", "Replace when configured nested values change"),
				},
				Description: "Container image information, including repository, tag, message and author.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"container_group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "ContainerGroupId",
			},
			"container_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "ContainerName",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the resource",
			},
		},
	}
}

func (r *commitContainerTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan commitContainerTaskResourceModel
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
		"ClientToken": buildClientToken("CommitContainer"),
	}
	// CommitContainer is an RPC API: RegionId travels in the request body/query in
	// addition to resolving the endpoint, and the provider region is authoritative.
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.AcrRegistryInfo.IsNull() && !plan.AcrRegistryInfo.IsUnknown() {
		createValues["AcrRegistryInfo"] = fwValueToInterface(plan.AcrRegistryInfo, map[string]fwWireField{"arn_service": {Name: "ArnService", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "instance_id": {Name: "InstanceId", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "arn_user": {Name: "ArnUser", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "region_id": {Name: "RegionId", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}}, true)
	}
	if !plan.Arn.IsNull() && !plan.Arn.IsUnknown() {
		createValues["Arn"] = fwValueToInterface(plan.Arn, map[string]fwWireField{"role_type": {Name: "RoleType", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "role_arn": {Name: "RoleArn", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}}, true)
	}
	if !plan.ContainerGroupId.IsNull() && !plan.ContainerGroupId.IsUnknown() {
		createValues["ContainerGroupId"] = fwValueToInterface(plan.ContainerGroupId, map[string]fwWireField{}, false)
	}
	if !plan.ContainerName.IsNull() && !plan.ContainerName.IsUnknown() {
		createValues["ContainerName"] = fwValueToInterface(plan.ContainerName, map[string]fwWireField{}, false)
	}
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		createValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.Image.IsNull() && !plan.Image.IsUnknown() {
		createValues["Image"] = fwValueToInterface(plan.Image, map[string]fwWireField{"repository": {Name: "Repository", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "message": {Name: "Message", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "author": {Name: "Author", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "tag": {Name: "Tag", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}}, true)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.AcrRegistryInfo.ArnService", Target: "AcrRegistryInfo.ArnService", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.AcrRegistryInfo.ArnUser", Target: "AcrRegistryInfo.ArnUser", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.AcrRegistryInfo.InstanceId", Target: "AcrRegistryInfo.InstanceId", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.AcrRegistryInfo.RegionId", Target: "AcrRegistryInfo.RegionId", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.Arn.RoleArn", Target: "Arn.RoleArn", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.Arn.RoleType", Target: "Arn.RoleType", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.ContainerGroupId", Target: "ContainerGroupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ContainerName", Target: "ContainerName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Image.Author", Target: "Image.Author", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.Image.Message", Target: "Image.Message", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.Image.Repository", Target: "Image.Repository", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.Image.Tag", Target: "Image.Tag", Encoding: "", Constant: "", Repeat: true, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(r.Client(), "Eci", "2018-08-08", "CommitContainer", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CommitContainer", commitContainerTaskTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CommitContainer", commitContainerTaskTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(response, fwMappingPath("$.TaskId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_eci_commit_container_task: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CommitContainer", commitContainerTaskTypeName),
			fmt.Sprintf("The response carried no CommitContainerTaskId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.AcrRegistryInfo = fwKnownValue(trackedPlan.AcrRegistryInfo).(types.List)
	trackedPlan.Arn = fwKnownValue(trackedPlan.Arn).(types.List)
	trackedPlan.ContainerGroupId = fwKnownValue(trackedPlan.ContainerGroupId).(types.String)
	trackedPlan.ContainerName = fwKnownValue(trackedPlan.ContainerName).(types.String)
	trackedPlan.CreateTime = fwKnownValue(trackedPlan.CreateTime).(types.String)
	trackedPlan.Image = fwKnownValue(trackedPlan.Image).(types.List)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := commitContainerTaskGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", commitContainerTaskTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", commitContainerTaskTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", commitContainerTaskTypeName),
			fmt.Sprintf("The resource %s is not visible via DescribeCommitContainerTask right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, commitContainerTaskState(object, &plan))...)
}

func (r *commitContainerTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state commitContainerTaskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := commitContainerTaskGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeCommitContainerTask", commitContainerTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeCommitContainerTask", commitContainerTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, commitContainerTaskState(object, &state))...)
}

func (r *commitContainerTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state commitContainerTaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *commitContainerTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// CommitContainer (Eci/2018-08-08) is an asynchronous task API with no
	// delete operation: it submits an image build/push task and returns a
	// CommitContainerTaskId, then describes the task via
	// DescribeCommitContainerTask. The task record is managed by the service
	// (retained for query after the task reaches a terminal state) and has no
	// delete API. Destroy removes the resource from Terraform state without
	// invoking any delete API or deleting the task record; Read already drops
	// the resource from state via RemoveResource when the task is no longer
	// visible. Destroy is a no-op rather than a hard error so that
	// `terraform destroy` can complete; no false destroy-success is reported
	// because no delete API is invoked and the task record is left intact for
	// the service to manage.
}

func (r *commitContainerTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// commitContainerTaskState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func commitContainerTaskState(object map[string]interface{}, prior *commitContainerTaskResourceModel) commitContainerTaskResourceModel {
	// Get has already merged and normalized every response mapping.
	str := func(key string) types.String {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.StringNull()
			}
			cur, ok = m[seg]
			if !ok {
				return types.StringNull()
			}
		}
		// Present-but-empty strings must read back as "", not null: the API returns
		// empty strings for some fields and Terraform treats "" as set, so collapsing
		// them to null would produce a permanent diff. Absent keys still return null.
		if v, ok := cur.(string); ok {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	m := commitContainerTaskResourceModel{
		FwTimeouts: types.ObjectNull(fwTimeoutTypes()),
		Id:         str("CommitContainerTaskId"),
		CreateTime: str("CreateTime"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		m.AcrRegistryInfo = prior.AcrRegistryInfo
		m.Arn = prior.Arn
		m.ContainerGroupId = prior.ContainerGroupId
		m.ContainerName = prior.ContainerName
		m.Image = prior.Image
	}
	// ContainerName is returned by DescribeCommitContainerTask, so read it back instead of
	// keeping the prior value; keep the configured value when the API omits it.
	if v := str("ContainerName"); !v.IsNull() {
		m.ContainerName = v
	}
	return m
}
