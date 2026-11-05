package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	dataWorksFileTypeName = "alicloud_data_works_file"
	// dataWorksFileTimeout replaces the SDKv2 resource-level create/update/delete
	// default timeouts; terraform-plugin-framework does not support
	// user-configurable resource timeouts yet.
	dataWorksFileTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &dataWorksFileResource{}
	_ resource.ResourceWithImportState = &dataWorksFileResource{}
)

func NewDataWorksFileResource() resource.Resource {
	return &dataWorksFileResource{}
}

type dataWorksFileResource struct {
	fwadapt.ResourceBase
}

// dataWorksFileResourceModel holds every configurable attribute of
// alicloud_data_works_file. Numeric API parameters (FileType,
// AutoRerunTimes, ...) are modelled as strings to match the DataWorks RPC
// convention used by the existing dataworks resources; the server coerces
// numeric strings on the wire.
type dataWorksFileResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	FileId                          types.String `tfsdk:"file_id"`
	NodeId                          types.String `tfsdk:"node_id"`
	FileName                        types.String `tfsdk:"file_name"`
	FileType                        types.String `tfsdk:"file_type"`
	ProjectId                       types.String `tfsdk:"project_id"`
	ProjectIdentifier               types.String `tfsdk:"project_identifier"`
	Content                         types.String `tfsdk:"content"`
	Owner                           types.String `tfsdk:"owner"`
	FileDescription                 types.String `tfsdk:"file_description"`
	FileFolderPath                  types.String `tfsdk:"file_folder_path"`
	ConnectionName                  types.String `tfsdk:"connection_name"`
	CycleType                       types.String `tfsdk:"cycle_type"`
	DependentType                   types.String `tfsdk:"dependent_type"`
	SchedulerType                   types.String `tfsdk:"scheduler_type"`
	ResourceGroupIdentifier         types.String `tfsdk:"resource_group_identifier"`
	AutoRerunTimes                  types.String `tfsdk:"auto_rerun_times"`
	AutoRerunIntervalMillis         types.String `tfsdk:"auto_rerun_interval_millis"`
	RerunMode                       types.String `tfsdk:"rerun_mode"`
	Stop                            types.Bool   `tfsdk:"stop"`
	ParaValue                       types.String `tfsdk:"para_value"`
	StartEffectDate                 types.String `tfsdk:"start_effect_date"`
	EndEffectDate                   types.String `tfsdk:"end_effect_date"`
	CronExpress                     types.String `tfsdk:"cron_express"`
	DependentNodeIdList             types.String `tfsdk:"dependent_node_id_list"`
	InputList                       types.String `tfsdk:"input_list"`
	OutputList                      types.String `tfsdk:"output_list"`
	AutoParsing                     types.Bool   `tfsdk:"auto_parsing"`
	AdvancedSettings                types.String `tfsdk:"advanced_settings"`
	StartImmediately                types.Bool   `tfsdk:"start_immediately"`
	InputParameters                 types.String `tfsdk:"input_parameters"`
	OutputParameters                types.String `tfsdk:"output_parameters"`
	IgnoreParentSkipRunningProperty types.Bool   `tfsdk:"ignore_parent_skip_running_property"`
	ApplyScheduleImmediately        types.Bool   `tfsdk:"apply_schedule_immediately"`
	Timeout                         types.String `tfsdk:"timeout"`
	ImageId                         types.String `tfsdk:"image_id"`
	Deploy                          types.Bool   `tfsdk:"deploy"`
	Submit                          types.Bool   `tfsdk:"submit"`
}

func (r *dataWorksFileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = dataWorksFileTypeName
}

func (r *dataWorksFileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_name": schema.StringAttribute{
				Required: true,
			},
			"file_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_identifier": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content":                    schema.StringAttribute{Optional: true, Computed: true},
			"owner":                      schema.StringAttribute{Optional: true, Computed: true},
			"file_description":           schema.StringAttribute{Optional: true, Computed: true},
			"file_folder_path":           schema.StringAttribute{Optional: true, Computed: true},
			"connection_name":            schema.StringAttribute{Optional: true, Computed: true},
			"cycle_type":                 schema.StringAttribute{Optional: true, Computed: true},
			"dependent_type":             schema.StringAttribute{Optional: true, Computed: true},
			"scheduler_type":             schema.StringAttribute{Optional: true, Computed: true},
			"resource_group_identifier":  schema.StringAttribute{Optional: true, Computed: true},
			"auto_rerun_times":           schema.StringAttribute{Optional: true, Computed: true},
			"auto_rerun_interval_millis": schema.StringAttribute{Optional: true, Computed: true},
			"rerun_mode":                 schema.StringAttribute{Optional: true, Computed: true},
			"stop": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"para_value":             schema.StringAttribute{Optional: true, Computed: true},
			"start_effect_date":      schema.StringAttribute{Optional: true, Computed: true},
			"end_effect_date":        schema.StringAttribute{Optional: true, Computed: true},
			"cron_express":           schema.StringAttribute{Optional: true, Computed: true},
			"dependent_node_id_list": schema.StringAttribute{Optional: true, Computed: true},
			"input_list":             schema.StringAttribute{Optional: true, Computed: true},
			"output_list":            schema.StringAttribute{Optional: true, Computed: true},
			"auto_parsing": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"advanced_settings": schema.StringAttribute{Optional: true, Computed: true},
			"start_immediately": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"input_parameters":  schema.StringAttribute{Optional: true, Computed: true},
			"output_parameters": schema.StringAttribute{Optional: true, Computed: true},
			"ignore_parent_skip_running_property": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"apply_schedule_immediately": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"timeout":  schema.StringAttribute{Optional: true, Computed: true},
			"image_id": schema.StringAttribute{Optional: true, Computed: true},
			// deploy and submit expose the DeployFile / SubmitFile lifecycle
			// actions as opt-in triggers; when true the resource calls the
			// corresponding API after create/update.
			"deploy": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"submit": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *dataWorksFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dataWorksFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := dataWorksFileValidateProject(plan); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Creating %s: invalid project identity", dataWorksFileTypeName), err.Error())
		return
	}

	request := dataWorksFileCreateRequest(plan)
	action := "CreateFile"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.RetryContext(ctx, dataWorksFileTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateFile", dataWorksFileTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, dataWorksFileTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	fileId, err := dataWorksFileExtractFileId(response)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading FileId from CreateFile response", dataWorksFileTypeName),
			err.Error(),
		)
		return
	}
	plan.FileId = types.StringValue(fileId)
	plan.Id = types.StringValue(dataWorksFileCompositeId(fileId, plan))

	if err := r.applyDeploySubmit(ctx, plan); err != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: publishing file", dataWorksFileTypeName),
			err.Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeDataWorksFileContext(ctx, plan.Id.ValueString(), plan.ProjectId.ValueString(), plan.ProjectIdentifier.ValueString())
	if err != nil {
		// Keep the created id in state so the file stays trackable even when
		// the read-back fails; Terraform marks the resource as tainted.
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: calling GetFile", dataWorksFileTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.Id.ValueString(), "GetFile", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, err := dataWorksFileState(object, plan.Id.ValueString(), plan)
	if err != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after create: mapping GetFile response", dataWorksFileTypeName),
			err.Error(),
		)
		return
	}
	dataWorksFilePreserveEmptyConfigured(&model, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dataWorksFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dataWorksFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		return
	}

	object, err := NewService(r.Meta()).DescribeDataWorksFileContext(ctx, id, state.ProjectId.ValueString(), state.ProjectIdentifier.ValueString())
	if err != nil {
		if errs.IsProviderNotFound(err) || errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetFile", dataWorksFileTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetFile", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	model, err := dataWorksFileState(object, state.Id.ValueString(), state)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: mapping GetFile response", dataWorksFileTypeName),
			err.Error(),
		)
		return
	}
	dataWorksFilePreserveEmptyConfigured(&model, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dataWorksFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dataWorksFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts, err := helper.ParseResourceId(plan.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing resource id", dataWorksFileTypeName),
			errs.WrapError(err).Error(),
		)
		return
	}

	request := dataWorksFileUpdateRequest(plan, parts[0])
	action := "UpdateFile"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.RetryContext(ctx, dataWorksFileTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateFile", dataWorksFileTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	if err := r.applyDeploySubmit(ctx, plan); err != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: publishing file", dataWorksFileTypeName),
			err.Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeDataWorksFileContext(ctx, plan.Id.ValueString(), plan.ProjectId.ValueString(), plan.ProjectIdentifier.ValueString())
	if err != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: calling GetFile", dataWorksFileTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.Id.ValueString(), "GetFile", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	model, err := dataWorksFileState(object, plan.Id.ValueString(), plan)
	if err != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s after update: mapping GetFile response", dataWorksFileTypeName),
			err.Error(),
		)
		return
	}
	dataWorksFilePreserveEmptyConfigured(&model, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dataWorksFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dataWorksFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if err := NewService(r.Meta()).DeleteDataWorksFileContext(ctx, id, state.ProjectId.ValueString(), state.ProjectIdentifier.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteFile", dataWorksFileTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DeleteFile", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *dataWorksFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// applyDeploySubmit calls DeployFile and/or SubmitFile when the corresponding
// opt-in flag is true. Both are best-effort publish actions; a failure is
// surfaced so the user knows the file was not published even though it was
// created/updated.
func (r *dataWorksFileResource) applyDeploySubmit(ctx context.Context, plan dataWorksFileResourceModel) error {
	if plan.Deploy.IsNull() || !plan.Deploy.ValueBool() {
		return nil
	}
	if err := r.callDeployFile(ctx, plan); err != nil {
		return fmt.Errorf("calling DeployFile: %w", err)
	}
	if plan.Submit.IsNull() || !plan.Submit.ValueBool() {
		return nil
	}
	if err := r.callSubmitFile(ctx, plan); err != nil {
		return fmt.Errorf("calling SubmitFile: %w", err)
	}
	return nil
}

func (r *dataWorksFileResource) callDeployFile(ctx context.Context, plan dataWorksFileResourceModel) error {
	parts, err := helper.ParseResourceId(plan.Id.ValueString(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{"FileId": parts[0]}
	dataWorksFileSetWorkspace(request, plan.ProjectId.ValueString(), plan.ProjectIdentifier.ValueString(), parts[1])
	if v := plan.NodeId.ValueString(); v != "" {
		request["NodeId"] = v
	}
	action := "DeployFile"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.RetryContext(ctx, dataWorksFileTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (r *dataWorksFileResource) callSubmitFile(ctx context.Context, plan dataWorksFileResourceModel) error {
	parts, err := helper.ParseResourceId(plan.Id.ValueString(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{"FileId": parts[0]}
	dataWorksFileSetWorkspace(request, plan.ProjectId.ValueString(), plan.ProjectIdentifier.ValueString(), parts[1])
	action := "SubmitFile"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.RetryContext(ctx, dataWorksFileTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = r.Meta().RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, plan.Id.ValueString(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// dataWorksFileValidateProject enforces the ProjectId/ProjectIdentifier
// "one-of" identity contract locally so a missing workspace identity is
// reported before reaching the API.
func dataWorksFileValidateProject(plan dataWorksFileResourceModel) error {
	pid := plan.ProjectId.ValueString()
	ident := plan.ProjectIdentifier.ValueString()
	if pid == "" && ident == "" {
		return fmt.Errorf("one of project_id or project_identifier must be set")
	}
	return nil
}

// dataWorksFileCompositeId builds "fileId:<workspace>" from the configured
// workspace identity: project_id when set, otherwise project_identifier. The
// same identity is later sent under the matching key on GetFile/DeleteFile so
// the resource resolves the workspace regardless of which form was configured.
func dataWorksFileCompositeId(fileId string, plan dataWorksFileResourceModel) string {
	workspace := plan.ProjectId.ValueString()
	if workspace == "" {
		workspace = plan.ProjectIdentifier.ValueString()
	}
	return fmt.Sprintf("%s:%s", fileId, workspace)
}

func dataWorksFileCreateRequest(plan dataWorksFileResourceModel) map[string]interface{} {
	request := map[string]interface{}{
		"FileName": plan.FileName.ValueString(),
		"FileType": plan.FileType.ValueString(),
	}
	dataWorksFileSetStr(request, "ProjectId", plan.ProjectId)
	dataWorksFileSetStr(request, "ProjectIdentifier", plan.ProjectIdentifier)
	dataWorksFileSetStr(request, "Content", plan.Content)
	dataWorksFileSetStr(request, "Owner", plan.Owner)
	dataWorksFileSetStr(request, "FileDescription", plan.FileDescription)
	dataWorksFileSetStr(request, "FileFolderPath", plan.FileFolderPath)
	dataWorksFileSetStr(request, "ConnectionName", plan.ConnectionName)
	dataWorksFileSetStr(request, "CycleType", plan.CycleType)
	dataWorksFileSetStr(request, "DependentType", plan.DependentType)
	dataWorksFileSetStr(request, "SchedulerType", plan.SchedulerType)
	dataWorksFileSetStr(request, "ResourceGroupIdentifier", plan.ResourceGroupIdentifier)
	dataWorksFileSetStr(request, "AutoRerunTimes", plan.AutoRerunTimes)
	dataWorksFileSetStr(request, "AutoRerunIntervalMillis", plan.AutoRerunIntervalMillis)
	dataWorksFileSetStr(request, "RerunMode", plan.RerunMode)
	dataWorksFileSetBool(request, "Stop", plan.Stop)
	dataWorksFileSetStr(request, "ParaValue", plan.ParaValue)
	dataWorksFileSetStr(request, "StartEffectDate", plan.StartEffectDate)
	dataWorksFileSetStr(request, "EndEffectDate", plan.EndEffectDate)
	dataWorksFileSetStr(request, "CronExpress", plan.CronExpress)
	dataWorksFileSetStr(request, "DependentNodeIdList", plan.DependentNodeIdList)
	dataWorksFileSetStr(request, "InputList", plan.InputList)
	dataWorksFileSetStr(request, "OutputList", plan.OutputList)
	dataWorksFileSetBool(request, "AutoParsing", plan.AutoParsing)
	dataWorksFileSetStr(request, "AdvancedSettings", plan.AdvancedSettings)
	dataWorksFileSetBool(request, "StartImmediately", plan.StartImmediately)
	dataWorksFileSetStr(request, "InputParameters", plan.InputParameters)
	dataWorksFileSetStr(request, "OutputParameters", plan.OutputParameters)
	dataWorksFileSetBool(request, "IgnoreParentSkipRunningProperty", plan.IgnoreParentSkipRunningProperty)
	dataWorksFileSetBool(request, "ApplyScheduleImmediately", plan.ApplyScheduleImmediately)
	dataWorksFileSetStr(request, "Timeout", plan.Timeout)
	dataWorksFileSetStr(request, "ImageId", plan.ImageId)
	return request
}

func dataWorksFileUpdateRequest(plan dataWorksFileResourceModel, fileId string) map[string]interface{} {
	request := map[string]interface{}{
		"FileId":   fileId,
		"FileName": plan.FileName.ValueString(),
	}
	dataWorksFileSetStr(request, "ProjectId", plan.ProjectId)
	dataWorksFileSetStr(request, "ProjectIdentifier", plan.ProjectIdentifier)
	dataWorksFileSetStr(request, "Content", plan.Content)
	dataWorksFileSetStr(request, "Owner", plan.Owner)
	dataWorksFileSetStr(request, "FileDescription", plan.FileDescription)
	dataWorksFileSetStr(request, "FileFolderPath", plan.FileFolderPath)
	dataWorksFileSetStr(request, "ConnectionName", plan.ConnectionName)
	dataWorksFileSetStr(request, "CycleType", plan.CycleType)
	dataWorksFileSetStr(request, "DependentType", plan.DependentType)
	dataWorksFileSetStr(request, "SchedulerType", plan.SchedulerType)
	dataWorksFileSetStr(request, "ResourceGroupIdentifier", plan.ResourceGroupIdentifier)
	dataWorksFileSetStr(request, "AutoRerunTimes", plan.AutoRerunTimes)
	dataWorksFileSetStr(request, "AutoRerunIntervalMillis", plan.AutoRerunIntervalMillis)
	dataWorksFileSetStr(request, "RerunMode", plan.RerunMode)
	dataWorksFileSetBool(request, "Stop", plan.Stop)
	dataWorksFileSetStr(request, "ParaValue", plan.ParaValue)
	dataWorksFileSetStr(request, "StartEffectDate", plan.StartEffectDate)
	dataWorksFileSetStr(request, "EndEffectDate", plan.EndEffectDate)
	dataWorksFileSetStr(request, "CronExpress", plan.CronExpress)
	dataWorksFileSetStr(request, "DependentNodeIdList", plan.DependentNodeIdList)
	dataWorksFileSetStr(request, "InputList", plan.InputList)
	dataWorksFileSetStr(request, "OutputList", plan.OutputList)
	dataWorksFileSetBool(request, "AutoParsing", plan.AutoParsing)
	dataWorksFileSetStr(request, "AdvancedSettings", plan.AdvancedSettings)
	dataWorksFileSetBool(request, "StartImmediately", plan.StartImmediately)
	dataWorksFileSetStr(request, "InputParameters", plan.InputParameters)
	dataWorksFileSetStr(request, "OutputParameters", plan.OutputParameters)
	dataWorksFileSetBool(request, "IgnoreParentSkipRunningProperty", plan.IgnoreParentSkipRunningProperty)
	dataWorksFileSetBool(request, "ApplyScheduleImmediately", plan.ApplyScheduleImmediately)
	dataWorksFileSetStr(request, "Timeout", plan.Timeout)
	dataWorksFileSetStr(request, "ImageId", plan.ImageId)
	return request
}

// dataWorksFileExtractFileId reads FileId from a CreateFile response. The
// published contract returns the new FileId directly as the top-level Data
// value (an int64, e.g. {"Data": 1000001}); DataWorksFileId accepts numeric
// and string forms. Legacy object shapes ($.FileId, $.Data.FileId) are kept
// as fallbacks for contract drift.
func dataWorksFileExtractFileId(response map[string]interface{}) (string, error) {
	if raw, ok := response["Data"]; ok && raw != nil {
		if id, ferr := DataWorksFileId(raw); ferr == nil {
			return id, nil
		}
	}
	if raw, ok := response["FileId"]; ok && raw != nil {
		return DataWorksFileId(raw)
	}
	if v, err := jsonpath.Get("$.Data.FileId", response); err == nil && v != nil {
		return DataWorksFileId(v)
	}
	return "", fmt.Errorf("FileId not found in CreateFile response: %v", response)
}

// dataWorksFileState maps the GetFile response object to the resource model.
// A contract-drifted response fails the refresh instead of writing null over
// configured values. The workspace identity is recovered from the configured
// state because GetFile does not echo ProjectId; the composite id segment is
// the final fallback for imports that supply "fileId:projectId".
func dataWorksFileState(object map[string]interface{}, id string, source dataWorksFileResourceModel) (dataWorksFileResourceModel, error) {
	fileId, err := DataWorksFileId(object["FileId"])
	if err != nil {
		return dataWorksFileResourceModel{}, err
	}
	name := dataWorksFileStr(object["FileName"])
	if name.IsNull() {
		return dataWorksFileResourceModel{}, fmt.Errorf("FileName is required in the GetFile response, got %T %v", object["FileName"], object["FileName"])
	}
	fileType := dataWorksFileStr(object["FileType"])
	if fileType.IsNull() {
		return dataWorksFileResourceModel{}, fmt.Errorf("FileType is required in the GetFile response, got %T %v", object["FileType"], object["FileType"])
	}
	// GetFile does not return ProjectId, so the workspace identity comes from
	// the configured state. When only project_identifier was set the id carries
	// the identifier; otherwise the id segment is a project_id for imports.
	projectId := types.StringNull()
	projectIdentifier := dataWorksFileStr(object["ProjectIdentifier"])
	switch {
	case !source.ProjectId.IsNull() && !source.ProjectId.IsUnknown():
		projectId = source.ProjectId
	case !source.ProjectIdentifier.IsNull() && !source.ProjectIdentifier.IsUnknown():
		projectIdentifier = source.ProjectIdentifier
	default:
		if parts, perr := helper.ParseResourceId(id, 2); perr == nil && len(parts) == 2 {
			projectId = types.StringValue(parts[1])
		}
	}
	workspace := projectId
	if workspace.IsNull() && !projectIdentifier.IsNull() {
		workspace = projectIdentifier
	}

	model := dataWorksFileResourceModel{
		Id:                              types.StringValue(fmt.Sprintf("%s:%s", fileId, workspace.ValueString())),
		FileId:                          types.StringValue(fileId),
		NodeId:                          dataWorksFileStr(object["NodeId"]),
		FileName:                        name,
		FileType:                        fileType,
		ProjectId:                       projectId,
		ProjectIdentifier:               projectIdentifier,
		Content:                         dataWorksFileStr(object["Content"]),
		Owner:                           dataWorksFileStr(object["Owner"]),
		FileDescription:                 dataWorksFileStr(object["FileDescription"]),
		FileFolderPath:                  dataWorksFileStr(object["FileFolderPath"]),
		ConnectionName:                  dataWorksFileStr(object["ConnectionName"]),
		CycleType:                       dataWorksFileStr(object["CycleType"]),
		DependentType:                   dataWorksFileStr(object["DependentType"]),
		SchedulerType:                   dataWorksFileStr(object["SchedulerType"]),
		ResourceGroupIdentifier:         dataWorksFileStr(object["ResourceGroupIdentifier"]),
		AutoRerunTimes:                  dataWorksFileStr(object["AutoRerunTimes"]),
		AutoRerunIntervalMillis:         dataWorksFileStr(object["AutoRerunIntervalMillis"]),
		RerunMode:                       dataWorksFileStr(object["RerunMode"]),
		Stop:                            dataWorksFileBool(object["Stop"]),
		ParaValue:                       dataWorksFileStr(object["ParaValue"]),
		StartEffectDate:                 dataWorksFileStr(object["StartEffectDate"]),
		EndEffectDate:                   dataWorksFileStr(object["EndEffectDate"]),
		CronExpress:                     dataWorksFileStr(object["CronExpress"]),
		DependentNodeIdList:             dataWorksFileStr(object["DependentNodeIdList"]),
		InputList:                       dataWorksFileStr(object["InputList"]),
		OutputList:                      dataWorksFileStr(object["OutputList"]),
		AutoParsing:                     dataWorksFileBool(object["AutoParsing"]),
		AdvancedSettings:                dataWorksFileStr(object["AdvancedSettings"]),
		StartImmediately:                dataWorksFileBool(object["StartImmediately"]),
		InputParameters:                 dataWorksFileStr(object["InputParameters"]),
		OutputParameters:                dataWorksFileStr(object["OutputParameters"]),
		IgnoreParentSkipRunningProperty: dataWorksFileBool(object["IgnoreParentSkipRunningProperty"]),
		ApplyScheduleImmediately:        dataWorksFileBool(object["ApplyScheduleImmediately"]),
		Timeout:                         dataWorksFileStr(object["Timeout"]),
		ImageId:                         dataWorksFileStr(object["ImageId"]),
	}
	return model, nil
}

// dataWorksFilePreserveEmptyConfigured keeps explicitly-empty configured values
// in state. CreateFile/UpdateFile omit empty optional parameters, so
// normalizing an unset response field to null would otherwise plan a
// perpetual diff for attributes the user intentionally left empty.
func dataWorksFilePreserveEmptyConfigured(model *dataWorksFileResourceModel, source dataWorksFileResourceModel) {
	dataWorksFilePreserveStr(&model.Content, source.Content)
	dataWorksFilePreserveStr(&model.Owner, source.Owner)
	dataWorksFilePreserveStr(&model.FileDescription, source.FileDescription)
	dataWorksFilePreserveStr(&model.FileFolderPath, source.FileFolderPath)
	dataWorksFilePreserveStr(&model.ConnectionName, source.ConnectionName)
	dataWorksFilePreserveStr(&model.CycleType, source.CycleType)
	dataWorksFilePreserveStr(&model.DependentType, source.DependentType)
	dataWorksFilePreserveStr(&model.SchedulerType, source.SchedulerType)
	dataWorksFilePreserveStr(&model.ResourceGroupIdentifier, source.ResourceGroupIdentifier)
	dataWorksFilePreserveStr(&model.AutoRerunTimes, source.AutoRerunTimes)
	dataWorksFilePreserveStr(&model.AutoRerunIntervalMillis, source.AutoRerunIntervalMillis)
	dataWorksFilePreserveStr(&model.RerunMode, source.RerunMode)
	dataWorksFilePreserveStr(&model.ParaValue, source.ParaValue)
	dataWorksFilePreserveStr(&model.StartEffectDate, source.StartEffectDate)
	dataWorksFilePreserveStr(&model.EndEffectDate, source.EndEffectDate)
	dataWorksFilePreserveStr(&model.CronExpress, source.CronExpress)
	dataWorksFilePreserveStr(&model.DependentNodeIdList, source.DependentNodeIdList)
	dataWorksFilePreserveStr(&model.InputList, source.InputList)
	dataWorksFilePreserveStr(&model.OutputList, source.OutputList)
	dataWorksFilePreserveStr(&model.AdvancedSettings, source.AdvancedSettings)
	dataWorksFilePreserveStr(&model.InputParameters, source.InputParameters)
	dataWorksFilePreserveStr(&model.OutputParameters, source.OutputParameters)
	dataWorksFilePreserveStr(&model.Timeout, source.Timeout)
	dataWorksFilePreserveStr(&model.ImageId, source.ImageId)
	dataWorksFilePreserveBool(&model.Stop, source.Stop)
	dataWorksFilePreserveBool(&model.AutoParsing, source.AutoParsing)
	dataWorksFilePreserveBool(&model.StartImmediately, source.StartImmediately)
	dataWorksFilePreserveBool(&model.IgnoreParentSkipRunningProperty, source.IgnoreParentSkipRunningProperty)
	dataWorksFilePreserveBool(&model.ApplyScheduleImmediately, source.ApplyScheduleImmediately)
	// project_identifier is a ForceNew workspace identity that GetFile does not
	// echo back; preserve any non-null configured value to avoid a perpetual diff.
	if model.ProjectIdentifier.IsNull() && !source.ProjectIdentifier.IsNull() && !source.ProjectIdentifier.IsUnknown() {
		model.ProjectIdentifier = source.ProjectIdentifier
	}
	if model.Deploy.IsNull() && !source.Deploy.IsNull() && !source.Deploy.IsUnknown() {
		model.Deploy = source.Deploy
	}
	if model.Submit.IsNull() && !source.Submit.IsNull() && !source.Submit.IsUnknown() {
		model.Submit = source.Submit
	}
}

// dataWorksFileSetStr adds a string parameter only when the plan value is
// non-null and non-empty, so the API request omits unset optional fields.
func dataWorksFileSetStr(request map[string]interface{}, key string, v types.String) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	if s := v.ValueString(); s != "" {
		request[key] = s
	}
}

func dataWorksFileSetBool(request map[string]interface{}, key string, v types.Bool) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	request[key] = v.ValueBool()
}

// dataWorksFileStr normalizes a response value into a string attribute. nil,
// empty strings and absent numbers yield null so unset fields are not written
// over configured state; numbers are formatted without trailing zeros.
func dataWorksFileStr(raw interface{}) types.String {
	switch v := raw.(type) {
	case nil:
		return types.StringNull()
	case string:
		if v == "" {
			return types.StringNull()
		}
		return types.StringValue(v)
	case float64:
		// FormatFloat with -1 precision renders integers without a trailing
		// ".0" and preserves fractional values verbatim.
		return types.StringValue(strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		return types.StringValue(strconv.FormatBool(v))
	case json.Number:
		return types.StringValue(v.String())
	default:
		return types.StringValue(fmt.Sprint(v))
	}
}

func dataWorksFileBool(raw interface{}) types.Bool {
	switch v := raw.(type) {
	case nil:
		return types.BoolNull()
	case bool:
		return types.BoolValue(v)
	default:
		return types.BoolNull()
	}
}

func dataWorksFilePreserveStr(model *types.String, source types.String) {
	if model.IsNull() && !source.IsNull() && !source.IsUnknown() && source.ValueString() == "" {
		*model = source
	}
}

func dataWorksFilePreserveBool(model *types.Bool, source types.Bool) {
	if model.IsNull() && !source.IsNull() && !source.IsUnknown() {
		*model = source
	}
}
