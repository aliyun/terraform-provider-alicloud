// Package cr. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package cr

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const repositoryTypeName = "alicloud_cr_repository"

var (
	_ resource.Resource                = &repositoryResource{}
	_ resource.ResourceWithImportState = &repositoryResource{}
)

func NewRepositoryResource() resource.Resource {
	return &repositoryResource{}
}

type repositoryResource struct {
	fwadapt.ResourceBase
}

type repositoryResourceModel struct {
	Id                types.String `tfsdk:"id"`
	FwTimeouts        types.Object `tfsdk:"timeouts"`
	CreateTime        types.String `tfsdk:"create_time"`
	Detail            types.String `tfsdk:"detail"`
	InstanceId        types.String `tfsdk:"instance_id"`
	ModifiedTime      types.String `tfsdk:"modified_time"`
	RepoBuildType     types.String `tfsdk:"repo_build_type"`
	RepoNamespaceName types.String `tfsdk:"repo_namespace_name"`
	RepoType          types.String `tfsdk:"repo_type"`
	RepositoryId      types.String `tfsdk:"repository_id"`
	RepositoryName    types.String `tfsdk:"repository_name"`
	ResourceGroupId   types.String `tfsdk:"resource_group_id"`
	Summary           types.String `tfsdk:"summary"`
	TagImmutability   types.Bool   `tfsdk:"tag_immutability"`
}

func (r *repositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = repositoryTypeName
}

func (r *repositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{"timeouts": fwTimeoutBlock()},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the resource",
			},
			"detail": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Detail",
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "InstanceId",
			},
			"modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last modification time",
			},
			"repo_build_type": schema.StringAttribute{
				Computed:    true,
				Description: "repository build type",
			},
			"repo_namespace_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "RepoNamespaceName",
			},
			"repo_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "RepoType",
			},
			"repository_id": schema.StringAttribute{
				Computed:    true,
				Description: "The first ID of the resource",
			},
			"repository_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "RepoName",
			},
			"resource_group_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource group",
			},
			"summary": schema.StringAttribute{
				Required:    true,
				Description: "Summary",
			},
			"tag_immutability": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "TagImmutability",
			},
		},
	}
}

func (r *repositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan repositoryResourceModel
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
		"ClientToken": buildClientToken("CreateRepository"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	createValues["Status"] = "ALL"
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		createValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.Detail.IsNull() && !plan.Detail.IsUnknown() {
		createValues["Detail"] = fwValueToInterface(plan.Detail, map[string]fwWireField{}, false)
	}
	if !plan.InstanceId.IsNull() && !plan.InstanceId.IsUnknown() {
		createValues["InstanceId"] = fwValueToInterface(plan.InstanceId, map[string]fwWireField{}, false)
	}
	if !plan.ModifiedTime.IsNull() && !plan.ModifiedTime.IsUnknown() {
		createValues["ModifiedTime"] = fwValueToInterface(plan.ModifiedTime, map[string]fwWireField{}, false)
	}
	if !plan.RepoBuildType.IsNull() && !plan.RepoBuildType.IsUnknown() {
		createValues["RepoBuildType"] = fwValueToInterface(plan.RepoBuildType, map[string]fwWireField{}, false)
	}
	if !plan.RepoNamespaceName.IsNull() && !plan.RepoNamespaceName.IsUnknown() {
		createValues["RepoNamespaceName"] = fwValueToInterface(plan.RepoNamespaceName, map[string]fwWireField{}, false)
	}
	if !plan.RepoType.IsNull() && !plan.RepoType.IsUnknown() {
		createValues["RepoType"] = fwValueToInterface(plan.RepoType, map[string]fwWireField{}, false)
	}
	if !plan.RepositoryId.IsNull() && !plan.RepositoryId.IsUnknown() {
		createValues["RepositoryId"] = fwValueToInterface(plan.RepositoryId, map[string]fwWireField{}, false)
	}
	if !plan.RepositoryName.IsNull() && !plan.RepositoryName.IsUnknown() {
		createValues["RepositoryName"] = fwValueToInterface(plan.RepositoryName, map[string]fwWireField{}, false)
	}
	if !plan.ResourceGroupId.IsNull() && !plan.ResourceGroupId.IsUnknown() {
		createValues["ResourceGroupId"] = fwValueToInterface(plan.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !plan.Summary.IsNull() && !plan.Summary.IsUnknown() {
		createValues["Summary"] = fwValueToInterface(plan.Summary, map[string]fwWireField{}, false)
	}
	if !plan.TagImmutability.IsNull() && !plan.TagImmutability.IsUnknown() {
		createValues["TagImmutability"] = fwValueToInterface(plan.TagImmutability, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.Detail", Target: "Detail", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.InstanceId", Target: "InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepoNamespaceName", Target: "RepoNamespaceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepoType", Target: "RepoType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.RepositoryName", Target: "RepoName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Summary", Target: "Summary", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.TagImmutability", Target: "TagImmutability", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(r.Client(), "cr", "2018-12-01", "CreateRepository", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateRepository", repositoryTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateRepository", repositoryTypeName).Error(),
		)
		return
	}

	identityParts := []string{}
	identityValues0 := fwWalkMappedPath(request, fwMappingPath("InstanceId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_cr_repository: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues0[0].Value))
	identityValues1 := fwWalkMappedPath(response, fwMappingPath("$.RepoId"), nil)
	if len(identityValues1) == 0 || identityValues1[0].Value == nil || fmt.Sprint(identityValues1[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_cr_repository: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	identityParts = append(identityParts, fmt.Sprint(identityValues1[0].Value))
	id := strings.Join(identityParts, ":")
	if id == "" || id == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateRepository", repositoryTypeName),
			fmt.Sprintf("The response carried no InstanceId: %v", response),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.CreateTime = fwKnownValue(trackedPlan.CreateTime).(types.String)
	trackedPlan.Detail = fwKnownValue(trackedPlan.Detail).(types.String)
	trackedPlan.InstanceId = fwKnownValue(trackedPlan.InstanceId).(types.String)
	trackedPlan.ModifiedTime = fwKnownValue(trackedPlan.ModifiedTime).(types.String)
	trackedPlan.RepoBuildType = fwKnownValue(trackedPlan.RepoBuildType).(types.String)
	trackedPlan.RepoNamespaceName = fwKnownValue(trackedPlan.RepoNamespaceName).(types.String)
	trackedPlan.RepoType = fwKnownValue(trackedPlan.RepoType).(types.String)
	trackedPlan.RepositoryId = fwKnownValue(trackedPlan.RepositoryId).(types.String)
	trackedPlan.RepositoryName = fwKnownValue(trackedPlan.RepositoryName).(types.String)
	trackedPlan.ResourceGroupId = fwKnownValue(trackedPlan.ResourceGroupId).(types.String)
	trackedPlan.Summary = fwKnownValue(trackedPlan.Summary).(types.String)
	trackedPlan.TagImmutability = fwKnownValue(trackedPlan.TagImmutability).(types.Bool)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := repositoryGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", repositoryTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", repositoryTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", repositoryTypeName),
			fmt.Sprintf("The resource %s is not visible via GetRepository right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, repositoryState(object, &plan))...)
}

func (r *repositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state repositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := repositoryGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetRepository", repositoryTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetRepository", repositoryTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, repositoryState(object, &state))...)
}

func (r *repositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state repositoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "update", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	planValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	planValues["Status"] = "ALL"
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		planValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.Detail.IsNull() && !plan.Detail.IsUnknown() {
		planValues["Detail"] = fwValueToInterface(plan.Detail, map[string]fwWireField{}, false)
	}
	if !plan.InstanceId.IsNull() && !plan.InstanceId.IsUnknown() {
		planValues["InstanceId"] = fwValueToInterface(plan.InstanceId, map[string]fwWireField{}, false)
	}
	if !plan.ModifiedTime.IsNull() && !plan.ModifiedTime.IsUnknown() {
		planValues["ModifiedTime"] = fwValueToInterface(plan.ModifiedTime, map[string]fwWireField{}, false)
	}
	if !plan.RepoBuildType.IsNull() && !plan.RepoBuildType.IsUnknown() {
		planValues["RepoBuildType"] = fwValueToInterface(plan.RepoBuildType, map[string]fwWireField{}, false)
	}
	if !plan.RepoNamespaceName.IsNull() && !plan.RepoNamespaceName.IsUnknown() {
		planValues["RepoNamespaceName"] = fwValueToInterface(plan.RepoNamespaceName, map[string]fwWireField{}, false)
	}
	if !plan.RepoType.IsNull() && !plan.RepoType.IsUnknown() {
		planValues["RepoType"] = fwValueToInterface(plan.RepoType, map[string]fwWireField{}, false)
	}
	if !plan.RepositoryId.IsNull() && !plan.RepositoryId.IsUnknown() {
		planValues["RepositoryId"] = fwValueToInterface(plan.RepositoryId, map[string]fwWireField{}, false)
	}
	if !plan.RepositoryName.IsNull() && !plan.RepositoryName.IsUnknown() {
		planValues["RepositoryName"] = fwValueToInterface(plan.RepositoryName, map[string]fwWireField{}, false)
	}
	if !plan.ResourceGroupId.IsNull() && !plan.ResourceGroupId.IsUnknown() {
		planValues["ResourceGroupId"] = fwValueToInterface(plan.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !plan.Summary.IsNull() && !plan.Summary.IsUnknown() {
		planValues["Summary"] = fwValueToInterface(plan.Summary, map[string]fwWireField{}, false)
	}
	if !plan.TagImmutability.IsNull() && !plan.TagImmutability.IsUnknown() {
		planValues["TagImmutability"] = fwValueToInterface(plan.TagImmutability, map[string]fwWireField{}, false)
	}
	stateValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	stateValues["Status"] = "ALL"
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		stateValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.Detail.IsNull() && !state.Detail.IsUnknown() {
		stateValues["Detail"] = fwValueToInterface(state.Detail, map[string]fwWireField{}, false)
	}
	if !state.InstanceId.IsNull() && !state.InstanceId.IsUnknown() {
		stateValues["InstanceId"] = fwValueToInterface(state.InstanceId, map[string]fwWireField{}, false)
	}
	if !state.ModifiedTime.IsNull() && !state.ModifiedTime.IsUnknown() {
		stateValues["ModifiedTime"] = fwValueToInterface(state.ModifiedTime, map[string]fwWireField{}, false)
	}
	if !state.RepoBuildType.IsNull() && !state.RepoBuildType.IsUnknown() {
		stateValues["RepoBuildType"] = fwValueToInterface(state.RepoBuildType, map[string]fwWireField{}, false)
	}
	if !state.RepoNamespaceName.IsNull() && !state.RepoNamespaceName.IsUnknown() {
		stateValues["RepoNamespaceName"] = fwValueToInterface(state.RepoNamespaceName, map[string]fwWireField{}, false)
	}
	if !state.RepoType.IsNull() && !state.RepoType.IsUnknown() {
		stateValues["RepoType"] = fwValueToInterface(state.RepoType, map[string]fwWireField{}, false)
	}
	if !state.RepositoryId.IsNull() && !state.RepositoryId.IsUnknown() {
		stateValues["RepositoryId"] = fwValueToInterface(state.RepositoryId, map[string]fwWireField{}, false)
	}
	if !state.RepositoryName.IsNull() && !state.RepositoryName.IsUnknown() {
		stateValues["RepositoryName"] = fwValueToInterface(state.RepositoryName, map[string]fwWireField{}, false)
	}
	if !state.ResourceGroupId.IsNull() && !state.ResourceGroupId.IsUnknown() {
		stateValues["ResourceGroupId"] = fwValueToInterface(state.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !state.Summary.IsNull() && !state.Summary.IsUnknown() {
		stateValues["Summary"] = fwValueToInterface(state.Summary, map[string]fwWireField{}, false)
	}
	if !state.TagImmutability.IsNull() && !state.TagImmutability.IsUnknown() {
		stateValues["TagImmutability"] = fwValueToInterface(state.TagImmutability, map[string]fwWireField{}, false)
	}
	updated0 := false
	if fwPlannedChanged(plan.Detail, state.Detail) {
		updated0 = true
	}
	if fwPlannedChanged(plan.Summary, state.Summary) {
		updated0 = true
	}
	if fwPlannedChanged(plan.TagImmutability, state.TagImmutability) {
		updated0 = true
	}
	if updated0 {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"InstanceId":  state.Id.ValueString(),
			"ClientToken": buildClientToken("UpdateRepository"),
		}
		parts0 := strings.Split(state.Id.ValueString(), ":")
		if len(parts0) > 0 {
			request0["InstanceId"] = parts0[0]
		}
		if len(parts0) > 1 {
			request0["RepositoryId"] = parts0[1]
			// gap: the API's UpdateRepository identifies the repo by the RepoId query param.
			// Set it authoritatively from the persisted state identity (mirroring repositoryGet
			// and v1 d.Get("repo_id")), rather than relying on plan.RepositoryId being carried
			// through the update plan. If the Computed repository_id decays to Unknown during
			// update planning, the $.RepositoryId->RepoId binding below would emit no RepoId and
			// the API answers InvalidNamespaceName (it cannot locate the repo). Pin RepoId here.
			request0["RepoId"] = parts0[1]
		}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, operationValues0, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.Detail", Target: "Detail", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.InstanceId", Target: "InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RepoNamespaceName", Target: "RepoNamespaceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RepoType", Target: "RepoType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RepositoryId", Target: "RepoId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RepositoryName", Target: "RepoName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Summary", Target: "Summary", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.TagImmutability", Target: "TagImmutability", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			// gap: UpdateRepository surfaces InvalidNamespaceName when RepoNamespaceName is absent even though
			// the spec marks it optional and RepoId (docRequired) is sent. The runtime needs the namespace
			// context alongside RepoId, mirroring CreateRepository (which always sends RepoNamespaceName+RepoName).
			// Send the unchanged identity fields so the API can locate the repo within its namespace hierarchy.
		}, operationValues0, stateValues, []string{"$.InstanceId", "$.RepositoryId", "$.RegionId", "$.RepoNamespaceName", "$.RepositoryName", "$.RepoType", "$.Summary"}), false); err != nil {
			resp.Diagnostics.AddError("Building UpdateRepository request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "cr", "2018-12-01", "UpdateRepository", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_cr_repository: UpdateRepository", err.Error())
			return
		}
	}
	// CancelRepoBuildRecord is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// CreateBuildRecordByRecord is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// CreateBuildRecordByRule is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// CreateRepoBuildRule is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// CreateRepoSourceCodeRepo is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// CreateRepoSyncTaskByRule is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// CreateRepoTagScanTask is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// CreateRepoTrigger is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// DeleteRepoBuildRule is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// DeleteRepoTag is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// DeleteRepoTrigger is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// UpdateRepoBuildRule is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// UpdateRepoSourceCodeRepo is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// UpdateRepoTrigger is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// CreateRepoSyncTask is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	object, err := repositoryGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetRepository", repositoryTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetRepository", repositoryTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, repositoryState(object, &plan))...)
}

func (r *repositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state repositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "delete", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	deleteValues["Status"] = "ALL"
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		deleteValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.Detail.IsNull() && !state.Detail.IsUnknown() {
		deleteValues["Detail"] = fwValueToInterface(state.Detail, map[string]fwWireField{}, false)
	}
	if !state.InstanceId.IsNull() && !state.InstanceId.IsUnknown() {
		deleteValues["InstanceId"] = fwValueToInterface(state.InstanceId, map[string]fwWireField{}, false)
	}
	if !state.ModifiedTime.IsNull() && !state.ModifiedTime.IsUnknown() {
		deleteValues["ModifiedTime"] = fwValueToInterface(state.ModifiedTime, map[string]fwWireField{}, false)
	}
	if !state.RepoBuildType.IsNull() && !state.RepoBuildType.IsUnknown() {
		deleteValues["RepoBuildType"] = fwValueToInterface(state.RepoBuildType, map[string]fwWireField{}, false)
	}
	if !state.RepoNamespaceName.IsNull() && !state.RepoNamespaceName.IsUnknown() {
		deleteValues["RepoNamespaceName"] = fwValueToInterface(state.RepoNamespaceName, map[string]fwWireField{}, false)
	}
	if !state.RepoType.IsNull() && !state.RepoType.IsUnknown() {
		deleteValues["RepoType"] = fwValueToInterface(state.RepoType, map[string]fwWireField{}, false)
	}
	if !state.RepositoryId.IsNull() && !state.RepositoryId.IsUnknown() {
		deleteValues["RepositoryId"] = fwValueToInterface(state.RepositoryId, map[string]fwWireField{}, false)
	}
	if !state.RepositoryName.IsNull() && !state.RepositoryName.IsUnknown() {
		deleteValues["RepositoryName"] = fwValueToInterface(state.RepositoryName, map[string]fwWireField{}, false)
	}
	if !state.ResourceGroupId.IsNull() && !state.ResourceGroupId.IsUnknown() {
		deleteValues["ResourceGroupId"] = fwValueToInterface(state.ResourceGroupId, map[string]fwWireField{}, false)
	}
	if !state.Summary.IsNull() && !state.Summary.IsUnknown() {
		deleteValues["Summary"] = fwValueToInterface(state.Summary, map[string]fwWireField{}, false)
	}
	if !state.TagImmutability.IsNull() && !state.TagImmutability.IsUnknown() {
		deleteValues["TagImmutability"] = fwValueToInterface(state.TagImmutability, map[string]fwWireField{}, false)
	}
	identityParts := strings.Split(state.Id.ValueString(), ":")
	if len(identityParts) > 0 {
		deleteValues["InstanceId"] = identityParts[0]
	}
	if len(identityParts) > 1 {
		deleteValues["RepositoryId"] = identityParts[1]
	}
	needsRefresh := false
	if needsRefresh {
		object, err := repositoryGet(r.Client(), state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Reading deletion conditions", err.Error())
			return
		}
		if object == nil {
			return
		}
		for key, value := range object {
			if value != nil {
				deleteValues[key] = value
			}
		}
	}
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
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteRepository")}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.InstanceId", Target: "InstanceId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RepoNamespaceName", Target: "RepoNamespaceName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RepositoryId", Target: "RepoId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RepositoryName", Target: "RepoName", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DeleteRepository request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "cr", "2018-12-01", "DeleteRepository", request0, operationTimeout); err != nil && !isExpectedErrors(err, []string{"ResourceNotFound.Repository", "ResourceNotFound"}) {
			resp.Diagnostics.AddError("Deleting alicloud_cr_repository: DeleteRepository", err.Error())
			return
		}
		// gap: CR repo deletion is synchronous — the repository is gone by the first poll.
		// A non-empty Target (["DELETING"]) with a gone resource raises NotFoundError instead of
		// success (see comment below). Use a nil (empty) Target so gone-wait succeeds, matching the
		// SDKv2 WaitForCrEERepo Deleted semantics.
		if _, err := repositoryWaitForState(r.Client(), state.Id.ValueString(), "Status", nil, nil, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Waiting after DeleteRepository", err.Error())
			return
		}
	}
	// gap: v1 Delete 用 BuildStateConf(target=[]) 轮询 StateRefreshFunc 直到 Get 返回 nil（gone），
	// 解决 gateway-route-table 异步删除后 DeleteVpc 竞态 DependencyViolation.GatewayRouteTable。
	// v2 WaitForState：空 deleteTargetStatus（gone-wait）发 Target=[]string{}（空 slice, len 0），
	// 镜像 v1 BuildStateConf(target=[])。hashicorp retry 当 Refresh 返回 (nil,"",nil)（gone）且
	// Target 空 → success（done）；若误用 Target=[""]（len 1, 非空），gone 时 helper 报
	// NotFoundError("couldn't find resource (N retries)") 而非 success——此前 route_table destroy
	// 即栽此。非空 deleteTargetStatus 则发 Target=["<status>"] 等待该终态。Get 在 not-found 返回
	// (nil,nil)（非 error），故 wait 不挂；同步删除资源首检即 gone（nil,"",nil→Target 空→done），无 hang。
	if _, err := repositoryWaitForState(r.Client(), state.Id.ValueString(), "Status", nil, nil, operationTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: waiting for gone", repositoryTypeName),
			errs.WrapErrorf(err, "Deleting %s: waiting for gone (%s)", repositoryTypeName, state.Id.ValueString()).Error(),
		)
		return
	}
}

func (r *repositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// repositoryState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func repositoryState(object map[string]interface{}, prior *repositoryResourceModel) repositoryResourceModel {
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
		// gap#7c: present-but-empty strings must return ValueString(""), not Null.
		// v1 d.Set(field, "") treats empty as set (not absence); TestCheckResourceAttrSet
		// expects a non-null value for API-returned "" (e.g. DhcpOptionsSetId when VPC has
		// no DHCP set). The earlier && v != "" collapsed present-empty to Null, failing Set
		// checks. Absent fields still return Null via the !ok branches above.
		if v, ok := cur.(string); ok {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
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
	// gap#11: composite-ID resource — build the id from all composite fields in the response,
	// joined by ":" (e.g. AccessGroupName:FileSystemType). The Get response carries every
	// part, so the state id round-trips through Read/Delete correctly (otherwise Read would
	// lose FileSystemType and DescribeAccessGroups would query the wrong namespace → 404).
	_cidParts := []string{
		str("InstanceId").ValueString(),
		str("RepositoryId").ValueString(),
	}
	m := repositoryResourceModel{
		FwTimeouts:        types.ObjectNull(fwTimeoutTypes()),
		Id:                types.StringValue(strings.Join(_cidParts, ":")),
		CreateTime:        str("CreateTime"),
		Detail:            str("Detail"),
		InstanceId:        str("InstanceId"),
		ModifiedTime:      str("ModifiedTime"),
		RepoBuildType:     str("RepoBuildType"),
		RepoNamespaceName: str("RepoNamespaceName"),
		RepoType:          str("RepoType"),
		RepositoryId:      str("RepositoryId"),
		RepositoryName:    str("RepositoryName"),
		ResourceGroupId:   str("ResourceGroupId"),
		Summary:           str("Summary"),
		TagImmutability:   boolVal("TagImmutability"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
	}
	return m
}
