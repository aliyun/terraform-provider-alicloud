package paidsw

import (
	"context"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	tempFileTaskTypeName = "alicloud_pai_dsw_temp_file_task"

	// PAI DSW product code and API version.
	paiDswProduct = "pai-dsw"
	paiDswVersion = "2022-01-01"

	// Resource not-found arrives as a *tea.SDKError with Code "ValidationError",
	// per the CloudSpec GetTempFileTask resourceNotExistCondition.
	tempFileTaskNotFoundCode = "ValidationError"
)

var (
	_ resource.Resource                = &tempFileTaskResource{}
	_ resource.ResourceWithImportState = &tempFileTaskResource{}
)

func NewTempFileTaskResource() resource.Resource {
	return &tempFileTaskResource{}
}

type tempFileTaskResource struct {
	fwadapt.ResourceBase
}

type tempFileTaskResourceModel struct {
	Id              types.String `tfsdk:"id"`
	InstanceId      types.String `tfsdk:"instance_id"`
	GmtExpiredTime  types.String `tfsdk:"gmt_expired_time"`
	OwnerId         types.String `tfsdk:"owner_id"`
	UserId          types.String `tfsdk:"user_id"`
	CreateTime      types.String `tfsdk:"create_time"`
	GmtModifiedTime types.String `tfsdk:"gmt_modified_time"`
	RegionId        types.String `tfsdk:"region_id"`
}

func (r *tempFileTaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = tempFileTaskTypeName
}

func (r *tempFileTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI DSW Temp File Task resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Temp File Task. Same as `temp_file_task_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the DSW instance.",
				PlanModifiers: []planmodifier.String{
					// Create-only on the API side: UpdateTempFileTask does not accept
					// InstanceId, so a change forces replacement.
					stringplanmodifier.RequiresReplace(),
				},
			},
			"gmt_expired_time": schema.StringAttribute{
				Optional:    true,
				Description: "The expiry time of the temp file task.",
			},
			"owner_id": schema.StringAttribute{
				Computed:    true,
				Description: "The owner ID of the temp file task.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The user ID of the temp file task.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the temp file task.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gmt_modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The modification time of the temp file task.",
				// No UseStateForUnknown: GmtModifiedTime is advanced by the
				// service on every Update. Pinning the prior state value in the
				// plan (UseStateForUnknown) would disagree with the refreshed
				// value Read returns after Update, producing a deterministic
				// inconsistent-result-after-apply. Leaving it unknown in plan so
				// the post-Read value applies cleanly (R3.8 server-advanced).
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the temp file task.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *tempFileTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tempFileTaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"InstanceId": plan.InstanceId.ValueString(),
	}
	if !plan.GmtExpiredTime.IsNull() {
		body["GmtExpiredTime"] = plan.GmtExpiredTime.ValueString()
	}

	response, err := r.Client().RoaPost(paiDswProduct, paiDswVersion, "/api/v2/tempfiletasks", nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateTempFileTask", tempFileTaskTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateTempFileTask", tempFileTaskTypeName).Error(),
		)
		return
	}

	id := strings.TrimSpace(fmt.Sprint(response["TempFileTaskId"]))
	if id == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateTempFileTask", tempFileTaskTypeName),
			fmt.Sprintf("The response carried no TempFileTaskId: %v", response),
		)
		return
	}

	object, err := getTempFileTask(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetTempFileTask", tempFileTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetTempFileTask", tempFileTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateTempFileTask", tempFileTaskTypeName),
			fmt.Sprintf("The temp file task %s is not visible via GetTempFileTask right after creation", id),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tempFileTaskState(object, r.Client().RegionId))...)
}

func (r *tempFileTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tempFileTaskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getTempFileTask(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetTempFileTask", tempFileTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetTempFileTask", tempFileTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tempFileTaskState(object, r.Client().RegionId))...)
}

func (r *tempFileTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state tempFileTaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// UpdateTempFileTask accepts only GmtExpiredTime; instance_id is ForceNew.
	if !plan.GmtExpiredTime.Equal(state.GmtExpiredTime) {
		body := map[string]interface{}{
			"GmtExpiredTime": plan.GmtExpiredTime.ValueString(),
		}
		if _, err := r.Client().RoaPut(paiDswProduct, paiDswVersion, "/api/v2/tempfiletasks/"+state.Id.ValueString(), nil, nil, body, true); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Updating %s: calling UpdateTempFileTask", tempFileTaskTypeName),
				errs.WrapErrorf(err, "Updating %s: calling UpdateTempFileTask", tempFileTaskTypeName).Error(),
			)
			return
		}
	}

	object, err := getTempFileTask(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetTempFileTask", tempFileTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetTempFileTask", tempFileTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, tempFileTaskState(object, r.Client().RegionId))...)
}

func (r *tempFileTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tempFileTaskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.Client().RoaDelete(paiDswProduct, paiDswVersion, "/api/v2/tempfiletasks/"+state.Id.ValueString(), nil, nil, nil, true)
	if err != nil && !isTempFileTaskNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteTempFileTask", tempFileTaskTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteTempFileTask", tempFileTaskTypeName).Error(),
		)
		return
	}
}

func (r *tempFileTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getTempFileTask queries GetTempFileTask by task ID and returns the response
// object, or nil when the task no longer exists.
func getTempFileTask(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := client.RoaGet(paiDswProduct, paiDswVersion, "/api/v2/tempfiletasks/"+id, nil, nil, nil)
	if err != nil {
		if isTempFileTaskNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return response, nil
}

// tempFileTaskState maps a GetTempFileTask response onto the resource model.
// Empty strings become null: dates and optional fields the API omits read
// back as absent rather than "". The task ID is returned as "Uuid" by
// GetTempFileTask, while CreateTempFileTask returns "TempFileTaskId".
func tempFileTaskState(object map[string]interface{}, regionId string) tempFileTaskResourceModel {
	str := func(key string) types.String {
		if v, ok := object[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	return tempFileTaskResourceModel{
		Id:              str("Uuid"),
		InstanceId:      str("InstanceId"),
		GmtExpiredTime:  str("GmtExpiredTime"),
		OwnerId:         str("OwnerId"),
		UserId:          str("UserId"),
		CreateTime:      str("GmtCreateTime"),
		GmtModifiedTime: str("GmtModifiedTime"),
		RegionId:        types.StringValue(regionId),
	}
}

// isTempFileNotExist reports whether err is the PAI DSW task-not-found error.
// The not-found code arrives on a *tea.SDKError from the ROA client as
// "ValidationError", per the CloudSpec resourceNotExistCondition; the substring
// fallback covers wrapped forms, matching IsExpectedErrors' behavior.
func isTempFileTaskNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == tempFileTaskNotFoundCode {
		return true
	}
	return strings.Contains(err.Error(), tempFileTaskNotFoundCode)
}
