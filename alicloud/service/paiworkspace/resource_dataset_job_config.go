package paiworkspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	datasetJobConfigTypeName = "alicloud_pai_workspace_dataset_job_config"
	aiWorkSpaceProduct       = "AIWorkSpace"
	aiWorkSpaceVersion       = "2021-02-04"
	idSeparator              = ":"
)

var (
	_ resource.Resource                = &datasetJobConfigResource{}
	_ resource.ResourceWithImportState = &datasetJobConfigResource{}
)

func NewDatasetJobConfigResource() resource.Resource {
	return &datasetJobConfigResource{}
}

type datasetJobConfigResource struct {
	fwadapt.ResourceBase
}

// datasetJobConfigResourceModel is the Terraform state model for the
// DatasetJobConfig resource. The `id` is a composite of
// <workspace_id>:<dataset_id>:<dataset_job_config_id> so that import and
// Read can reconstruct the three identifiers Get/Delete/Update require.
type datasetJobConfigResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	WorkspaceId        types.String `tfsdk:"workspace_id"`
	DatasetId          types.String `tfsdk:"dataset_id"`
	DatasetJobConfigId types.String `tfsdk:"dataset_job_config_id"`
	ConfigType         types.String `tfsdk:"config_type"`
	Config             types.String `tfsdk:"config"`
	DatasetVersion     types.String `tfsdk:"dataset_version"`
	CreateTime         types.String `tfsdk:"create_time"`
	ModifyTime         types.String `tfsdk:"modify_time"`
}

func (r *datasetJobConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = datasetJobConfigTypeName
}

func (r *datasetJobConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI Workspace Dataset Job Config resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID in format `<workspace_id>:<dataset_id>:<dataset_job_config_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				Required:    true,
				Description: "The workspace ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dataset_id": schema.StringAttribute{
				Required:    true,
				Description: "The dataset ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dataset_job_config_id": schema.StringAttribute{
				Computed:    true,
				Description: "The dataset job config ID, assigned by the server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"config_type": schema.StringAttribute{
				Required:    true,
				Description: "The config type. Valid values: `MultimodalIntelligentTag`, `MultimodalSemanticIndex`.",
				Validators: []validator.String{
					stringvalidator.OneOf("MultimodalIntelligentTag", "MultimodalSemanticIndex"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config": schema.StringAttribute{
				Required:    true,
				Description: "The config content in JSON format. The structure depends on `config_type`.",
			},
			"dataset_version": schema.StringAttribute{
				Optional:    true,
				Description: "The dataset version name. Write-only: set at create time, not returned by Get.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time, in ISO 8601 (UTC) format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modify_time": schema.StringAttribute{
				Computed:    true,
				Description: "The modification time, in ISO 8601 (UTC) format.",
				// No UseStateForUnknown: UpdateDatasetJobConfig bumps the
				// timestamp, so a planned value can never survive an apply.
			},
		},
	}
}

func (r *datasetJobConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan datasetJobConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"WorkspaceId": plan.WorkspaceId.ValueString(),
		"ConfigType":  plan.ConfigType.ValueString(),
		"Config":      plan.Config.ValueString(),
	}
	if !plan.DatasetVersion.IsNull() {
		body["DatasetVersion"] = plan.DatasetVersion.ValueString()
	}

	action := fmt.Sprintf("/api/v1/datasets/%s/datasetjobconfigs", plan.DatasetId.ValueString())
	response, err := r.Client().RoaPost(aiWorkSpaceProduct, aiWorkSpaceVersion, action, nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDatasetJobConfig", datasetJobConfigTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateDatasetJobConfig", datasetJobConfigTypeName).Error(),
		)
		return
	}

	datasetJobConfigId, _ := response["DatasetJobConfigId"].(string)
	if datasetJobConfigId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDatasetJobConfig", datasetJobConfigTypeName),
			fmt.Sprintf("The response carried no DatasetJobConfigId: %v", response),
		)
		return
	}

	plan.DatasetJobConfigId = types.StringValue(datasetJobConfigId)
	plan.Id = types.StringValue(compositeID(plan.WorkspaceId.ValueString(), plan.DatasetId.ValueString(), datasetJobConfigId))

	object, err := getDatasetJobConfig(r.Client(), plan.WorkspaceId.ValueString(), plan.DatasetId.ValueString(), datasetJobConfigId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDatasetJobConfig", datasetJobConfigTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDatasetJobConfig", datasetJobConfigTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateDatasetJobConfig", datasetJobConfigTypeName),
			fmt.Sprintf("The dataset job config %s is not visible via GetDatasetJobConfig right after creation", datasetJobConfigId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, datasetJobConfigState(object, plan))...)
}

func (r *datasetJobConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state datasetJobConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getDatasetJobConfig(r.Client(), state.WorkspaceId.ValueString(), state.DatasetId.ValueString(), state.DatasetJobConfigId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDatasetJobConfig", datasetJobConfigTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDatasetJobConfig", datasetJobConfigTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, datasetJobConfigState(object, state))...)
}

func (r *datasetJobConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state datasetJobConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// config_type and the identity fields are ForceNew, so the only in-place
	// mutable field is config. WorkspaceId is part of the body per the API.
	body := map[string]interface{}{
		"WorkspaceId": plan.WorkspaceId.ValueString(),
		"ConfigType":  plan.ConfigType.ValueString(),
		"Config":      plan.Config.ValueString(),
	}

	action := fmt.Sprintf("/api/v1/datasets/%s/datasetjobconfigs/%s", plan.DatasetId.ValueString(), plan.DatasetJobConfigId.ValueString())
	if _, err := r.Client().RoaPut(aiWorkSpaceProduct, aiWorkSpaceVersion, action, queryWorkspace(plan.WorkspaceId.ValueString()), nil, body, true); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateDatasetJobConfig", datasetJobConfigTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateDatasetJobConfig", datasetJobConfigTypeName).Error(),
		)
		return
	}

	object, err := getDatasetJobConfig(r.Client(), plan.WorkspaceId.ValueString(), plan.DatasetId.ValueString(), plan.DatasetJobConfigId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetDatasetJobConfig", datasetJobConfigTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetDatasetJobConfig", datasetJobConfigTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, datasetJobConfigState(object, plan))...)
}

func (r *datasetJobConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state datasetJobConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := fmt.Sprintf("/api/v1/datasets/%s/datasetjobconfigs/%s", state.DatasetId.ValueString(), state.DatasetJobConfigId.ValueString())
	_, err := r.Client().RoaDelete(aiWorkSpaceProduct, aiWorkSpaceVersion, action, queryWorkspace(state.WorkspaceId.ValueString()), nil, nil, true)
	if err != nil && !isDatasetJobConfigNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteDatasetJobConfig", datasetJobConfigTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteDatasetJobConfig", datasetJobConfigTypeName).Error(),
		)
		return
	}
}

// ImportState accepts an ID of the form <workspace_id>:<dataset_id>:<dataset_job_config_id>.
func (r *datasetJobConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, idSeparator)
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s", datasetJobConfigTypeName),
			fmt.Sprintf("Unexpected import ID format %q, expected <workspace_id>:<dataset_id>:<dataset_job_config_id>", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_id"), types.StringValue(parts[0]))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("dataset_id"), types.StringValue(parts[1]))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("dataset_job_config_id"), types.StringValue(parts[2]))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}

// getDatasetJobConfig queries GetDatasetJobConfig and returns the object, or nil
// when the config no longer exists.
func getDatasetJobConfig(client *connectivity.AliyunClient, workspaceId, datasetId, datasetJobConfigId string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/api/v1/datasets/%s/datasetjobconfigs/%s", datasetId, datasetJobConfigId)
	response, err := client.RoaGet(aiWorkSpaceProduct, aiWorkSpaceVersion, action, queryWorkspace(workspaceId), nil, nil)
	if err != nil {
		if isDatasetJobConfigNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if response == nil {
		return nil, nil
	}
	return response, nil
}

// datasetJobConfigState maps a GetDatasetJobConfig response onto the state model,
// preserving the create-only DatasetVersion and the server-assigned identifiers
// from the prior state (Get returns neither).
func datasetJobConfigState(resp map[string]interface{}, prior datasetJobConfigResourceModel) datasetJobConfigResourceModel {
	str := func(key string) types.String {
		if v, ok := resp[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	// WorkspaceId and DatasetId are echoed by Get; fall back to the prior
	// state value if the API omits one (defensive). DatasetVersion and the
	// server-assigned identifiers are not in the Get response, so they are
	// carried forward from the prior state.
	workspaceId := str("WorkspaceId")
	if workspaceId.IsNull() {
		workspaceId = prior.WorkspaceId
	}
	datasetId := str("DatasetId")
	if datasetId.IsNull() {
		datasetId = prior.DatasetId
	}
	return datasetJobConfigResourceModel{
		Id:                 prior.Id,
		WorkspaceId:        workspaceId,
		DatasetId:          datasetId,
		DatasetJobConfigId: prior.DatasetJobConfigId,
		ConfigType:         str("ConfigType"),
		Config:             str("Config"),
		DatasetVersion:     prior.DatasetVersion,
		CreateTime:         str("CreateTime"),
		ModifyTime:         str("ModifyTime"),
	}
}

func compositeID(workspaceId, datasetId, datasetJobConfigId string) string {
	return fmt.Sprintf("%s%s%s%s%s", workspaceId, idSeparator, datasetId, idSeparator, datasetJobConfigId)
}

func queryWorkspace(workspaceId string) map[string]*string {
	if workspaceId == "" {
		return nil
	}
	v := workspaceId
	return map[string]*string{"WorkspaceId": &v}
}

// isDatasetJobConfigNotExist reports whether err is a not-found error. ROA
// errors arrive on a *tea.SDKError (404 / NotFound message); errs.NotFoundError
// covers that, and the substring fallback covers the PAI-specific error codes.
//
// PAI also surfaces a deleted DatasetJobConfig as NoPermissionError (HTTP 400,
// "No permission for this action: PaiDataset:GetDatasetJobConfig") instead of a
// 404. This was verified via ACC: GetDatasetJobConfig succeeds during the
// create/import steps and returns NoPermissionError only after the destroy
// step, so it is the API's gone signal. Treating it as not-found lets Read drop
// the resource from state and keeps Delete idempotent for an already-deleted
// config.
func isDatasetJobConfigNotExist(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	msg := err.Error()
	lowerMsg := strings.ToLower(msg)
	if strings.Contains(msg, "201300002") || strings.Contains(msg, "201300003") {
		return true
	}
	if strings.Contains(lowerMsg, "notfound") || strings.Contains(lowerMsg, "not exist") {
		return true
	}
	return strings.Contains(lowerMsg, "nopermissionerror") || strings.Contains(lowerMsg, "no permission for this action")
}
