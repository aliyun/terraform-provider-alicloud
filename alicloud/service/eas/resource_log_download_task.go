package eas

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

const logDownloadTaskTypeName = "alicloud_pai_log_download_task"

// EAS LogDownloadTask lives under product code "EAS", version "2021-07-01". The
// API is restful: CreateLogDownloadTask POSTs to /api/v2/log-download-tasks with a
// JSON body, and Describe/Delete route on /api/v2/log-download-tasks/{ClusterId}/{TaskId}.
// ClusterId is the regional EAS cluster, supplied by client.RegionId (matching the
// existing alicloud_pai_service resource), so the Terraform id is the single TaskId.

var (
	_ resource.Resource                = &logDownloadTaskResource{}
	_ resource.ResourceWithImportState = &logDownloadTaskResource{}
)

func NewLogDownloadTaskResource() resource.Resource {
	return &logDownloadTaskResource{}
}

type logDownloadTaskResource struct {
	fwadapt.ResourceBase
}

type logDownloadTaskModel struct {
	Id           types.String `tfsdk:"id"`
	TaskId       types.String `tfsdk:"task_id"`
	ServiceName  types.String `tfsdk:"service_name"`
	AliasName    types.String `tfsdk:"alias_name"`
	FromTime     types.String `tfsdk:"from_time"`
	ToTime       types.String `tfsdk:"to_time"`
	InstanceName types.String `tfsdk:"instance_name"`
	Container    types.String `tfsdk:"container"`
	CreateTime   types.String `tfsdk:"create_time"`
	FinishTime   types.String `tfsdk:"finish_time"`
	Status       types.String `tfsdk:"status"`
	DownloadUrl  types.String `tfsdk:"download_url"`
	RegionId     types.String `tfsdk:"region_id"`
}

func (r *logDownloadTaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = logDownloadTaskTypeName
}

func (r *logDownloadTaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a PAI EAS LogDownloadTask resource. " +
			"A log download task collects service logs for a specified EAS service " +
			"within a time range and produces a downloadable archive.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource. Same as `task_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"task_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the log download task assigned by EAS.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the EAS service whose logs are downloaded.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The alias name of the log download task.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"from_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The start time of the log collection range, formatted as `yyyy-MM-dd HH:mm:ss`. " +
					"Defaults to seven days before `to_time`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"to_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The end time of the log collection range, formatted as `yyyy-MM-dd HH:mm:ss`. " +
					"Defaults to the current time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The name of the EAS instance to filter logs by. " +
					"An empty value means all instances.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"container": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The name of the business container to filter logs by. " +
					"An empty value means all business containers.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the log download task.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"finish_time": schema.StringAttribute{
				Computed:    true,
				Description: "The finish time of the log download task.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the log download task, such as `Running`, `Finished` or `Failed`.",
			},
			"download_url": schema.StringAttribute{
				Computed:    true,
				Description: "The download URL of the produced log archive, populated when the task finishes.",
			},
			"region_id": schema.StringAttribute{
				Computed:    true,
				Description: "The region ID of the EAS cluster, same as the provider region.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *logDownloadTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan logDownloadTaskModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]interface{}{
		"ServiceName": plan.ServiceName.ValueString(),
	}
	if !plan.AliasName.IsNull() {
		body["AliasName"] = plan.AliasName.ValueString()
	}
	if !plan.FromTime.IsNull() {
		body["FromTime"] = plan.FromTime.ValueString()
	}
	if !plan.ToTime.IsNull() {
		body["ToTime"] = plan.ToTime.ValueString()
	}
	if !plan.InstanceName.IsNull() {
		body["InstanceName"] = plan.InstanceName.ValueString()
	}
	if !plan.Container.IsNull() {
		body["Container"] = plan.Container.ValueString()
	}

	response, err := r.Client().RoaPost("EAS", "2021-07-01", "/api/v2/log-download-tasks", nil, nil, body, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateLogDownloadTask", logDownloadTaskTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateLogDownloadTask", logDownloadTaskTypeName).Error(),
		)
		return
	}

	taskId, _ := response["TaskId"].(string)
	if taskId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateLogDownloadTask", logDownloadTaskTypeName),
			fmt.Sprintf("The response carried no TaskId: %v", response),
		)
		return
	}

	object, err := getLogDownloadTask(r.Client(), taskId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLogDownloadTask", logDownloadTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLogDownloadTask", logDownloadTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateLogDownloadTask", logDownloadTaskTypeName),
			fmt.Sprintf("The log download task %s is not visible via DescribeLogDownloadTask right after creation", taskId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, logDownloadTaskState(object, r.Client().RegionId))...)
}

func (r *logDownloadTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state logDownloadTaskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := getLogDownloadTask(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLogDownloadTask", logDownloadTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLogDownloadTask", logDownloadTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, logDownloadTaskState(object, r.Client().RegionId))...)
}

func (r *logDownloadTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Every settable attribute carries RequiresReplace, so an in-place update
	// never reaches this method. Re-reading keeps state honest if it ever does.
	var state logDownloadTaskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	object, err := getLogDownloadTask(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeLogDownloadTask", logDownloadTaskTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeLogDownloadTask", logDownloadTaskTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, logDownloadTaskState(object, r.Client().RegionId))...)
}

func (r *logDownloadTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state logDownloadTaskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := deleteLogDownloadTask(r.Client(), state.Id.ValueString()); err != nil && !isLogDownloadTaskNotExist(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteLogDownloadTask", logDownloadTaskTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling DeleteLogDownloadTask", logDownloadTaskTypeName).Error(),
		)
		return
	}
}

func (r *logDownloadTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getLogDownloadTask queries DescribeLogDownloadTask by task id and returns the
// Task object from the response, or nil when the task no longer exists. The
// regional EAS cluster id is supplied by client.RegionId.
func getLogDownloadTask(client *connectivity.AliyunClient, taskId string) (map[string]interface{}, error) {
	clusterId := client.RegionId
	action := fmt.Sprintf("/api/v2/log-download-tasks/%s/%s", clusterId, taskId)
	query := map[string]*string{"ClusterId": &clusterId}
	response, err := client.RoaGet("EAS", "2021-07-01", action, query, nil, nil)
	if err != nil {
		if isLogDownloadTaskNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	task, _ := response["Task"].(map[string]interface{})
	if task == nil {
		return nil, fmt.Errorf("the DescribeLogDownloadTask response carried no Task object: %v", response)
	}
	return task, nil
}

// deleteLogDownloadTask calls DeleteLogDownloadTask by task id.
func deleteLogDownloadTask(client *connectivity.AliyunClient, taskId string) (map[string]interface{}, error) {
	clusterId := client.RegionId
	action := fmt.Sprintf("/api/v2/log-download-tasks/%s/%s", clusterId, taskId)
	query := map[string]*string{"ClusterId": &clusterId}
	return client.RoaDelete("EAS", "2021-07-01", action, query, nil, nil, true)
}

// logDownloadTaskState maps the Task object of a DescribeLogDownloadTask response
// onto the resource model. Empty strings become null: dates and optional fields
// the API omits read back as absent rather than "". regionId is the regional EAS
// cluster, which the response does not carry.
func logDownloadTaskState(task map[string]interface{}, regionId string) logDownloadTaskModel {
	str := func(key string) types.String {
		if v, ok := task[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	config, _ := task["Config"].(map[string]interface{})
	cfgStr := func(key string) types.String {
		if config != nil {
			if v, ok := config[key].(string); ok && v != "" {
				return types.StringValue(v)
			}
		}
		return types.StringNull()
	}
	var region types.String
	if regionId != "" {
		region = types.StringValue(regionId)
	} else {
		region = types.StringNull()
	}
	taskId := str("TaskId")
	return logDownloadTaskModel{
		Id:           taskId,
		TaskId:       taskId,
		ServiceName:  str("ServiceName"),
		AliasName:    str("AliasName"),
		FromTime:     cfgStr("FromTime"),
		ToTime:       cfgStr("ToTime"),
		InstanceName: cfgStr("InstanceName"),
		Container:    cfgStr("Container"),
		CreateTime:   str("CreateTime"),
		FinishTime:   str("FinishTime"),
		Status:       str("Status"),
		DownloadUrl:  str("DownloadUrl"),
		RegionId:     region,
	}
}

// isLogDownloadTaskNotExist reports whether err signals the task is gone. The EAS
// restful API surfaces a *tea.SDKError whose code carries a NotFound/NotExist
// marker; the substring fallback covers wrapped forms, matching the SDK line's
// IsExpectedErrors behavior.
func isLogDownloadTaskNotExist(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil {
		code := *e.Code
		if strings.Contains(code, "NotFound") || strings.Contains(code, "NotExist") {
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "NotFound") || strings.Contains(msg, "NotExist")
}
