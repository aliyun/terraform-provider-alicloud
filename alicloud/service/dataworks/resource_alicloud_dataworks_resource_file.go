package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const dataWorksResourceFileTypeName = "alicloud_dataworks_resource_file"

// dataWorksResourceFileTimeout replaces the SDKv2 resource-level default
// timeouts (Create/Update/Delete each 5 minutes); terraform-plugin-framework
// does not support user-configurable resource timeouts yet.
const dataWorksResourceFileTimeout = 5 * time.Minute

var (
	_ resource.Resource                = &dataWorksResourceFileResource{}
	_ resource.ResourceWithImportState = &dataWorksResourceFileResource{}
)

// NewDataWorksResourceFileResource is the Framework resource factory for the
// DataWorks resource file. It is registered alongside the existing SDKv2
// DataWorks resources via conns.ServicePackage.Resources.
func NewDataWorksResourceFileResource() resource.Resource {
	return &dataWorksResourceFileResource{}
}

type dataWorksResourceFileResource struct {
	fwadapt.ResourceBase
}

// dataWorksResourceFileResourceModel maps the Terraform schema of
// alicloud_dataworks_resource_file. The composite id is "<FileId>:<ProjectId>",
// mirroring alicloud_data_works_folder.
//
// ProjectId and FileId are exposed as strings, matching
// alicloud_data_works_folder (whose id-derived parts are strings) and avoiding
// cross-resource type coercion when project_id references
// alicloud_data_works_project.id (an SDKv2 string id). The dataworks-public
// API accepts string values for these int64 formData fields.
//
// Field lifecycle follows the verified OpenAPI contract
// (dataworks-public/2020-05-18):
//   - CreateResourceFile accepts: ProjectId, FileFolderPath, FileName,
//     FileType, FileDescription, Owner, Content, OriginResourceName,
//     RegisterToCalcEngine, StorageURL.
//   - GetFile returns (under $.Data.File): Content, FileDescription, FileId,
//     FileName, FileType, Owner. FileFolderPath/OriginResourceName/
//     RegisterToCalcEngine/StorageURL are never read back, so they are kept
//     from the prior state on refresh.
//   - UpdateFile accepts: FileFolderPath, FileName, FileDescription, Content,
//     Owner, FileId, ProjectId (and many scheduler-only fields). It does NOT
//     accept FileType/OriginResourceName/RegisterToCalcEngine/StorageURL, so
//     those are ForceNew.
type dataWorksResourceFileResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	FileId               types.String `tfsdk:"file_id"`
	ProjectId            types.String `tfsdk:"project_id"`
	RegionId             types.String `tfsdk:"region_id"`
	FileName             types.String `tfsdk:"file_name"`
	FileType             types.Int64  `tfsdk:"file_type"`
	Content              types.String `tfsdk:"content"`
	FileDescription      types.String `tfsdk:"file_description"`
	Owner                types.String `tfsdk:"owner"`
	FileFolderPath       types.String `tfsdk:"file_folder_path"`
	OriginResourceName   types.String `tfsdk:"origin_resource_name"`
	RegisterToCalcEngine types.Bool   `tfsdk:"register_to_calc_engine"`
	StorageURL           types.String `tfsdk:"storage_url"`
}

func (r *dataWorksResourceFileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = dataWorksResourceFileTypeName
}

func (r *dataWorksResourceFileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// FileId is the server-assigned primary key returned by
			// CreateResourceFile ($.Data) and GetFile ($.Data.File.FileId);
			// it is read-only.
			"file_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// ProjectId is part of the composite identity and required by every
			// CRUD call; changing it means a different resource.
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_name": schema.StringAttribute{
				Required: true,
			},
			// FileType is required at create, accepted by CreateResourceFile
			// only (UpdateFile has no FileType), and read back by GetFile, so it
			// is ForceNew while still refreshable.
			"file_type": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"content":          schema.StringAttribute{Optional: true},
			"file_description": schema.StringAttribute{Optional: true},
			"owner":            schema.StringAttribute{Optional: true},
			// FileFolderPath is writable on Create and Update but never returned
			// by GetFile; refresh keeps the configured value.
			"file_folder_path": schema.StringAttribute{Optional: true},
			// The next three are create-only (CreateResourceFile accepts them,
			// UpdateFile does not) and never read back, so they are ForceNew and
			// preserved from state on refresh.
			"origin_resource_name": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"register_to_calc_engine": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"storage_url": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *dataWorksResourceFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dataWorksResourceFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action := "CreateResourceFile"
	request := map[string]interface{}{
		"ProjectId": plan.ProjectId.ValueString(),
		"FileName":  plan.FileName.ValueString(),
		"FileType":  plan.FileType.ValueInt64(),
	}
	if !plan.FileFolderPath.IsNull() {
		request["FileFolderPath"] = plan.FileFolderPath.ValueString()
	}
	if !plan.FileDescription.IsNull() {
		request["FileDescription"] = plan.FileDescription.ValueString()
	}
	if !plan.Owner.IsNull() {
		request["Owner"] = plan.Owner.ValueString()
	}
	if !plan.Content.IsNull() {
		request["Content"] = plan.Content.ValueString()
	}
	if !plan.OriginResourceName.IsNull() {
		request["OriginResourceName"] = plan.OriginResourceName.ValueString()
	}
	if !plan.RegisterToCalcEngine.IsNull() {
		request["RegisterToCalcEngine"] = plan.RegisterToCalcEngine.ValueBool()
	}
	if !plan.StorageURL.IsNull() {
		request["StorageURL"] = plan.StorageURL.ValueString()
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(dataWorksResourceFileTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateResourceFile", dataWorksResourceFileTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateResourceFile", dataWorksResourceFileTypeName).Error(),
		)
		return
	}

	fileId := scalarStr(response["Data"])
	if fileId.IsNull() {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: parsing FileId", dataWorksResourceFileTypeName),
			fmt.Sprintf("CreateResourceFile did not return a FileId: %v", response["Data"]),
		)
		return
	}
	plan.FileId = fileId
	plan.Id = types.StringValue(fmt.Sprintf("%s:%s", fileId.ValueString(), plan.ProjectId.ValueString()))
	plan.RegionId = types.StringValue(r.Meta().RegionId)

	object, err := NewService(r.Meta()).DescribeDataWorksResourceFile(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetFile", dataWorksResourceFileTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetFile", dataWorksResourceFileTypeName).Error(),
		)
		return
	}
	state := dataWorksResourceFileState(object)
	// Overlay the write-only / create-only values GetFile never returns.
	state.Id = plan.Id
	state.ProjectId = plan.ProjectId
	state.RegionId = plan.RegionId
	state.FileFolderPath = plan.FileFolderPath
	state.OriginResourceName = plan.OriginResourceName
	state.RegisterToCalcEngine = plan.RegisterToCalcEngine
	state.StorageURL = plan.StorageURL
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *dataWorksResourceFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dataWorksResourceFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := NewService(r.Meta()).DescribeDataWorksResourceFile(state.Id.ValueString())
	if err != nil {
		if errs.NotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetFile", dataWorksResourceFileTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetFile", dataWorksResourceFileTypeName).Error(),
		)
		return
	}
	model := dataWorksResourceFileState(object)
	// Preserve the composite id and the write-only / create-only values that
	// GetFile does not return, so drift detection keeps the configured value.
	model.Id = state.Id
	model.ProjectId = state.ProjectId
	model.RegionId = state.RegionId
	model.FileFolderPath = state.FileFolderPath
	model.OriginResourceName = state.OriginResourceName
	model.RegisterToCalcEngine = state.RegisterToCalcEngine
	model.StorageURL = state.StorageURL
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dataWorksResourceFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state dataWorksResourceFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts, err := helper.ParseResourceId(state.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: parsing resource id", dataWorksResourceFileTypeName),
			errs.WrapErrorf(err, "Updating %s: parsing resource id", dataWorksResourceFileTypeName).Error(),
		)
		return
	}

	action := "UpdateFile"
	request := map[string]interface{}{
		"FileId":    parts[0],
		"ProjectId": state.ProjectId.ValueString(),
		"FileName":  plan.FileName.ValueString(),
	}
	if !plan.FileFolderPath.IsNull() {
		request["FileFolderPath"] = plan.FileFolderPath.ValueString()
	}
	if !plan.FileDescription.IsNull() {
		request["FileDescription"] = plan.FileDescription.ValueString()
	}
	if !plan.Content.IsNull() {
		request["Content"] = plan.Content.ValueString()
	}
	if !plan.Owner.IsNull() {
		request["Owner"] = plan.Owner.ValueString()
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(dataWorksResourceFileTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling UpdateFile", dataWorksResourceFileTypeName),
			errs.WrapErrorf(err, "Updating %s: calling UpdateFile", dataWorksResourceFileTypeName).Error(),
		)
		return
	}

	object, err := NewService(r.Meta()).DescribeDataWorksResourceFile(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetFile", dataWorksResourceFileTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetFile", dataWorksResourceFileTypeName).Error(),
		)
		return
	}
	model := dataWorksResourceFileState(object)
	model.Id = state.Id
	model.ProjectId = state.ProjectId
	model.RegionId = state.RegionId
	model.FileFolderPath = plan.FileFolderPath
	model.OriginResourceName = state.OriginResourceName
	model.RegisterToCalcEngine = state.RegisterToCalcEngine
	model.StorageURL = state.StorageURL
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *dataWorksResourceFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dataWorksResourceFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts, err := helper.ParseResourceId(state.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: parsing resource id", dataWorksResourceFileTypeName),
			errs.WrapErrorf(err, "Deleting %s: parsing resource id", dataWorksResourceFileTypeName).Error(),
		)
		return
	}

	action := "DeleteFile"
	request := map[string]interface{}{
		"FileId":    parts[0],
		"ProjectId": state.ProjectId.ValueString(),
	}

	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(dataWorksResourceFileTimeout, func() *retry.RetryError {
		response, err = r.Meta().RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		// Idempotent delete: a missing file is already gone.
		if errs.IsExpectedErrors(err, []string{"Invalid.File.NotFound", "File.NotFound", "Invalid.FileId.NotFound"}) || errs.NotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteFile", dataWorksResourceFileTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteFile", dataWorksResourceFileTypeName).Error(),
		)
		return
	}

	// DeleteFile returns a top-level DeploymentId when the file had been
	// committed and an async scheduler-side deletion was triggered. An empty or
	// zero value means the deletion already completed (no scheduler work); a
	// non-zero value is polled via GetDeployment until the deployment reaches a
	// terminal status (1 success, 2 failed).
	deploymentId := scalarStr(response["DeploymentId"])
	if depID := deploymentId.ValueString(); !deploymentId.IsNull() && depID != "" && depID != "0" {
		if pollErr := NewService(r.Meta()).WaitDataWorksResourceFileDeployment(state.ProjectId.ValueString(), depID); pollErr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Deleting %s: waiting for deployment", dataWorksResourceFileTypeName),
				errs.WrapErrorf(pollErr, "Deleting %s: waiting for deployment %s", dataWorksResourceFileTypeName, deploymentId.ValueString()).Error(),
			)
			return
		}
	}
}

func (r *dataWorksResourceFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// dataWorksResourceFileState maps a GetFile ($.Data.File) object onto the
// resource model. Write-only and create-only fields are deliberately left null;
// every caller overlays the plan/state value before writing state.
func dataWorksResourceFileState(object map[string]interface{}) dataWorksResourceFileResourceModel {
	str := func(key string) types.String {
		switch v := object[key].(type) {
		case string:
			if v != "" {
				return types.StringValue(v)
			}
			return types.StringNull()
		// GetFile may return FileId/FileType as json.Number, int, int64 or
		// float64; normalize them to their string/int64 representations so the
		// planned values do not collapse to null ("inconsistent result after
		// apply").
		case json.Number:
			return types.StringValue(v.String())
		case float64:
			return types.StringValue(strconv.FormatInt(int64(v), 10))
		case int:
			return types.StringValue(strconv.Itoa(v))
		case int64:
			return types.StringValue(strconv.FormatInt(v, 10))
		}
		return types.StringNull()
	}
	int64Val := func(key string) types.Int64 {
		switch v := object[key].(type) {
		case float64:
			return types.Int64Value(int64(v))
		case json.Number:
			if n, err := v.Int64(); err == nil {
				return types.Int64Value(n)
			}
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				return types.Int64Value(n)
			}
		}
		return types.Int64Null()
	}
	return dataWorksResourceFileResourceModel{
		FileId:          str("FileId"),
		FileName:        str("FileName"),
		FileType:        int64Val("FileType"),
		Content:         str("Content"),
		FileDescription: str("FileDescription"),
		Owner:           str("Owner"),
	}
}

// scalarStr decodes a scalar value (json.Number/int/int64/float64/string) the
// tea-rpc client may produce into a types.String. It is used for
// CreateResourceFile's $.Data (FileId) and DeleteFile's DeploymentId responses.
func scalarStr(v interface{}) types.String {
	switch n := v.(type) {
	case string:
		if n == "" {
			return types.StringNull()
		}
		return types.StringValue(n)
	case json.Number:
		return types.StringValue(n.String())
	case float64:
		return types.StringValue(strconv.FormatInt(int64(n), 10))
	case int:
		return types.StringValue(strconv.Itoa(n))
	case int64:
		return types.StringValue(strconv.FormatInt(n, 10))
	}
	return types.StringNull()
}

// DescribeDataWorksResourceFile reads a single resource file via GetFile. The
// id is "<FileId>:<ProjectId>". It returns the $.Data.File object; a missing
// file is surfaced as a NotFound error so Framework Read can remove it.
func (s *Service) DescribeDataWorksResourceFile(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}
	action := "GetFile"
	request := map[string]interface{}{
		"FileId":    parts[0],
		"ProjectId": parts[1],
	}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(dataWorksResourceFileTimeout, func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"Throttling.Api", "Throttling.System", "Throttling.User"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	data, ok := response["Data"].(map[string]interface{})
	if !ok {
		return object, errs.WrapErrorf(errs.NotFoundErr("ResourceFile", id), errs.NotFoundMsg, response)
	}
	file, ok := data["File"].(map[string]interface{})
	if !ok || file == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("ResourceFile", id), errs.NotFoundMsg, response)
	}
	// A File object without a FileId means the file does not exist.
	if _, exists := file["FileId"]; !exists {
		return object, errs.WrapErrorf(errs.NotFoundErr("ResourceFile", id), errs.NotFoundMsg, response)
	}
	return file, nil
}

// WaitDataWorksResourceFileDeployment polls GetDeployment until the async
// delete triggered by DeleteFile reaches a terminal status. Status is 0
// (ready/running), 1 (success) or 2 (failed); 2 is surfaced as an error.
func (s *Service) WaitDataWorksResourceFileDeployment(projectId, deploymentId string) error {
	client := s.client
	action := "GetDeployment"
	request := map[string]interface{}{
		"DeploymentId": deploymentId,
		"ProjectId":    projectId,
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	return retry.Retry(dataWorksResourceFileTimeout, func() *retry.RetryError {
		response, err := client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"Throttling.Api", "Throttling.System", "Throttling.User"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		data, _ := response["Data"].(map[string]interface{})
		deployment, _ := data["Deployment"].(map[string]interface{})
		status := int64Scalar(deployment["Status"])
		switch status.ValueInt64() {
		case 1:
			return nil
		case 2:
			return retry.NonRetryableError(fmt.Errorf("deployment %s for DataWorks resource file failed (status=2)", deploymentId))
		default:
			// 0 = still running; keep waiting.
			wait()
			return retry.RetryableError(fmt.Errorf("deployment %s is still running (status=0)", deploymentId))
		}
	})
}

// int64Scalar decodes a scalar value (json.Number/int/int64/float64/string)
// the tea-rpc client may produce into a types.Int64. It is used for the
// GetDeployment Deployment.Status response.
func int64Scalar(v interface{}) types.Int64 {
	switch n := v.(type) {
	case float64:
		return types.Int64Value(int64(n))
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return types.Int64Value(i)
		}
	case int:
		return types.Int64Value(int64(n))
	case int64:
		return types.Int64Value(n)
	case string:
		if i, err := strconv.ParseInt(n, 10, 64); err == nil {
			return types.Int64Value(i)
		}
	}
	return types.Int64Null()
}
