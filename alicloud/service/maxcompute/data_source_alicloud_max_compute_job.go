package maxcompute

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const maxComputeJobTypeName = "alicloud_max_compute_job"

// maxComputeJobNotFoundCodes are the error codes that GetJobInfo returns (with
// HTTP 200) when the job instance does not exist; they are not transient.
var maxComputeJobNotFoundCodes = []string{
	"OBJECT_NOT_EXIST",
}

var (
	_ datasource.DataSource              = &maxComputeJobDataSource{}
	_ datasource.DataSourceWithConfigure = &maxComputeJobDataSource{}
)

func NewMaxComputeJobDataSource() datasource.DataSource {
	return &maxComputeJobDataSource{}
}

type maxComputeJobDataSource struct {
	fwadapt.DataSourceBase
}

// maxComputeJobModel is the flat state of the data source. GetJobInfo returns a
// single JobInfo object, so the data source has one required input (instance_id)
// and every API field as a computed output.
type maxComputeJobModel struct {
	Id               types.String `tfsdk:"id"`
	InstanceId       types.String `tfsdk:"instance_id"`
	TenantId         types.String `tfsdk:"tenant_id"`
	Region           types.String `tfsdk:"region"`
	JobOwner         types.String `tfsdk:"job_owner"`
	Project          types.String `tfsdk:"project"`
	JobType          types.String `tfsdk:"job_type"`
	Priority         types.Int64  `tfsdk:"priority"`
	Status           types.String `tfsdk:"status"`
	QuotaNickname    types.String `tfsdk:"quota_nickname"`
	QuotaType        types.String `tfsdk:"quota_type"`
	CuUsage          types.Int64  `tfsdk:"cu_usage"`
	MemoryUsage      types.Int64  `tfsdk:"memory_usage"`
	SubmittedAtTime  types.Int64  `tfsdk:"submitted_at_time"`
	WaitingTime      types.Int64  `tfsdk:"waiting_time"`
	RunningAtTime    types.Int64  `tfsdk:"running_at_time"`
	RunningTime      types.Int64  `tfsdk:"running_time"`
	EndAtTime        types.Int64  `tfsdk:"end_at_time"`
	TotalTime        types.Int64  `tfsdk:"total_time"`
	Signature        types.String `tfsdk:"signature"`
	ExtPlantFrom     types.String `tfsdk:"ext_plant_from"`
	ExtNodeId        types.String `tfsdk:"ext_node_id"`
	ExtNodeOnDuty    types.String `tfsdk:"ext_node_on_duty"`
	InputBytes       types.Int64  `tfsdk:"input_bytes"`
	SceneResults     types.List   `tfsdk:"scene_results"`
	JobSubStatusList types.List   `tfsdk:"job_sub_status_list"`
}

type maxComputeJobSceneResultModel struct {
	SceneTag    types.String `tfsdk:"scene_tag"`
	Scene       types.String `tfsdk:"scene"`
	Type        types.String `tfsdk:"type"`
	Summary     types.String `tfsdk:"summary"`
	Description types.String `tfsdk:"description"`
	Params      types.Map    `tfsdk:"params"`
}

type maxComputeJobSubStatusModel struct {
	Code        types.Int64  `tfsdk:"code"`
	StartTime   types.String `tfsdk:"start_time"`
	Description types.String `tfsdk:"description"`
}

var maxComputeJobSceneResultObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"scene_tag":   types.StringType,
	"scene":       types.StringType,
	"type":        types.StringType,
	"summary":     types.StringType,
	"description": types.StringType,
	"params":      types.MapType{ElemType: types.StringType},
}}

var maxComputeJobSubStatusObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"code":        types.Int64Type,
	"start_time":  types.StringType,
	"description": types.StringType,
}}

func (d *maxComputeJobDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = maxComputeJobTypeName
}

func (d *maxComputeJobDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides information about an Alibaba Cloud MaxCompute job.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, set to the job instance ID.",
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the MaxCompute job instance. This is the `instanceId` path parameter of the GetJobInfo API.",
			},
			"tenant_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the tenant to which the job belongs.",
			},
			"region": schema.StringAttribute{
				Computed:    true,
				Description: "The region where the job runs.",
			},
			"job_owner": schema.StringAttribute{
				Computed:    true,
				Description: "The owner of the job.",
			},
			"project": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the MaxCompute project to which the job belongs.",
			},
			"job_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the job, such as `SQL`, `MAPREDUCE`, `GRAPH`, `STREAM` and `CUPID`.",
			},
			"priority": schema.Int64Attribute{
				Computed:    true,
				Description: "The priority of the job.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the job, such as `running`, `terminated` and `failed`.",
			},
			"quota_nickname": schema.StringAttribute{
				Computed:    true,
				Description: "The nickname of the quota used by the job.",
			},
			"quota_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the quota used by the job, such as `subscription` and `payAsYouGo`.",
			},
			"cu_usage": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of compute units (CUs) used by the job.",
			},
			"memory_usage": schema.Int64Attribute{
				Computed:    true,
				Description: "The memory usage of the job.",
			},
			"submitted_at_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The time at which the job was submitted, as a Unix timestamp in seconds.",
			},
			"waiting_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The time that the job spent waiting, in seconds.",
			},
			"running_at_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The time at which the job started running, as a Unix timestamp in seconds.",
			},
			"running_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The time that the job spent running, in seconds.",
			},
			"end_at_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The time at which the job ended, as a Unix timestamp in seconds.",
			},
			"total_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The total elapsed time of the job, in seconds.",
			},
			"signature": schema.StringAttribute{
				Computed:    true,
				Description: "The signature of the job.",
			},
			"ext_plant_from": schema.StringAttribute{
				Computed:    true,
				Description: "The external platform from which the job was submitted.",
			},
			"ext_node_id": schema.StringAttribute{
				Computed:    true,
				Description: "The external node ID associated with the job.",
			},
			"ext_node_on_duty": schema.StringAttribute{
				Computed:    true,
				Description: "The external on-duty node associated with the job.",
			},
			"input_bytes": schema.Int64Attribute{
				Computed:    true,
				Description: "The input data size of the job, in bytes.",
			},
			"scene_results": schema.ListAttribute{
				Computed:    true,
				ElementType: maxComputeJobSceneResultObjectType,
				Description: "The scene analysis results of the job. Each element contains `scene_tag`, `scene`, `type`, `summary`, `description` and `params`.",
			},
			"job_sub_status_list": schema.ListAttribute{
				Computed:    true,
				ElementType: maxComputeJobSubStatusObjectType,
				Description: "The sub-status list of the job. Each element contains `code`, `start_time` and `description`.",
			},
		},
	}
}

func (d *maxComputeJobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model maxComputeJobModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := model.InstanceId.ValueString()
	if instanceId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", maxComputeJobTypeName),
			"instance_id is required to look up a MaxCompute job.",
		)
		return
	}

	client := d.Meta()
	action := fmt.Sprintf("/api/v1/jobs/%s/info", instanceId)
	query := make(map[string]*string)

	var response map[string]interface{}
	// Cancellation can return before the RPC callback finishes writing the
	// response, so the write is guarded; the same pattern is used by the other
	// Framework data sources that retry with RetryContext.
	var responseMu sync.Mutex
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := sdkretry.RetryContext(ctx, 1*time.Minute, func() *sdkretry.RetryError {
		result, rpcErr := client.RoaGet("MaxCompute", "2022-01-04", action, query, nil, nil)
		responseMu.Lock()
		response = result
		responseMu.Unlock()
		if rpcErr != nil {
			if errs.IsExpectedErrors(rpcErr, []string{"500"}) || helper.NeedRetry(rpcErr) {
				wait()
				return sdkretry.RetryableError(rpcErr)
			}
			return sdkretry.NonRetryableError(rpcErr)
		}
		return nil
	})
	responseMu.Lock()
	defer responseMu.Unlock()
	logs.AddDebug(action, response, map[string]interface{}{"instanceId": instanceId})
	if err != nil {
		if errs.IsExpectedErrors(err, maxComputeJobNotFoundCodes) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: job not found", maxComputeJobTypeName),
				fmt.Sprintf("The MaxCompute job instance %q does not exist: %s", instanceId, err.Error()),
			)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetJobInfo", maxComputeJobTypeName),
			err.Error(),
		)
		return
	}

	data, ok := maxComputeJobData(response)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: parsing response", maxComputeJobTypeName),
			fmt.Sprintf("The GetJobInfo response did not contain a data object for instance %q.", instanceId),
		)
		return
	}

	sceneResults, diags := maxComputeJobSceneResults(data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	jobSubStatusList, diags := maxComputeJobSubStatusList(data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.Id = types.StringValue(instanceId)
	model.InstanceId = types.StringValue(maxComputeJobString(data, "instanceId"))
	model.TenantId = types.StringValue(maxComputeJobString(data, "tenantId"))
	model.Region = types.StringValue(maxComputeJobString(data, "region"))
	model.JobOwner = types.StringValue(maxComputeJobString(data, "jobOwner"))
	model.Project = types.StringValue(maxComputeJobString(data, "project"))
	model.JobType = types.StringValue(maxComputeJobString(data, "jobType"))
	model.Priority = types.Int64Value(maxComputeJobInt(data, "priority"))
	model.Status = types.StringValue(maxComputeJobString(data, "status"))
	model.QuotaNickname = types.StringValue(maxComputeJobString(data, "quotaNickname"))
	model.QuotaType = types.StringValue(maxComputeJobString(data, "quotaType"))
	model.CuUsage = types.Int64Value(maxComputeJobInt(data, "cuUsage"))
	model.MemoryUsage = types.Int64Value(maxComputeJobInt(data, "memoryUsage"))
	model.SubmittedAtTime = types.Int64Value(maxComputeJobInt(data, "submittedAtTime"))
	model.WaitingTime = types.Int64Value(maxComputeJobInt(data, "waitingTime"))
	model.RunningAtTime = types.Int64Value(maxComputeJobInt(data, "runningAtTime"))
	model.RunningTime = types.Int64Value(maxComputeJobInt(data, "runningTime"))
	model.EndAtTime = types.Int64Value(maxComputeJobInt(data, "endAtTime"))
	model.TotalTime = types.Int64Value(maxComputeJobInt(data, "totalTime"))
	model.Signature = types.StringValue(maxComputeJobString(data, "signature"))
	model.ExtPlantFrom = types.StringValue(maxComputeJobString(data, "extPlantFrom"))
	model.ExtNodeId = types.StringValue(maxComputeJobString(data, "extNodeId"))
	model.ExtNodeOnDuty = types.StringValue(maxComputeJobString(data, "extNodeOnDuty"))
	model.InputBytes = types.Int64Value(maxComputeJobInt(data, "inputBytes"))
	model.SceneResults = sceneResults
	model.JobSubStatusList = jobSubStatusList

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// maxComputeJobData extracts the JobInfo object from the PopResult envelope.
// GetJobInfo wraps the payload under $.data, matching the other MaxCompute v2
// ROA APIs.
func maxComputeJobData(response map[string]interface{}) (map[string]interface{}, bool) {
	if response == nil {
		return nil, false
	}
	raw, err := jsonpath.Get("$.data", response)
	if err != nil {
		return nil, false
	}
	data, ok := raw.(map[string]interface{})
	return data, ok
}

func maxComputeJobSceneResults(data map[string]interface{}) (types.List, diag.Diagnostics) {
	items := maxComputeJobSlice(data, "sceneResults")
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		paramsMap := make(map[string]attr.Value)
		if raw, ok := item["params"].(map[string]interface{}); ok {
			for k, v := range raw {
				paramsMap[k] = types.StringValue(fmt.Sprint(v))
			}
		}
		params, diags := types.MapValue(types.StringType, paramsMap)
		if diags.HasError() {
			return types.ListNull(maxComputeJobSceneResultObjectType), diags
		}
		obj, diags := types.ObjectValue(maxComputeJobSceneResultObjectType.AttrTypes, map[string]attr.Value{
			"scene_tag":   types.StringValue(maxComputeJobString(item, "sceneTag")),
			"scene":       types.StringValue(maxComputeJobString(item, "scene")),
			"type":        types.StringValue(maxComputeJobString(item, "type")),
			"summary":     types.StringValue(maxComputeJobString(item, "summary")),
			"description": types.StringValue(maxComputeJobString(item, "description")),
			"params":      params,
		})
		if diags.HasError() {
			return types.ListNull(maxComputeJobSceneResultObjectType), diags
		}
		values = append(values, obj)
	}
	list, diags := types.ListValue(maxComputeJobSceneResultObjectType, values)
	return list, diags
}

func maxComputeJobSubStatusList(data map[string]interface{}) (types.List, diag.Diagnostics) {
	items := maxComputeJobSlice(data, "jobSubStatusList")
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		obj, diags := types.ObjectValue(maxComputeJobSubStatusObjectType.AttrTypes, map[string]attr.Value{
			"code":        types.Int64Value(maxComputeJobInt(item, "code")),
			"start_time":  types.StringValue(maxComputeJobString(item, "startTime")),
			"description": types.StringValue(maxComputeJobString(item, "description")),
		})
		if diags.HasError() {
			return types.ListNull(maxComputeJobSubStatusObjectType), diags
		}
		values = append(values, obj)
	}
	list, diags := types.ListValue(maxComputeJobSubStatusObjectType, values)
	return list, diags
}

// maxComputeJobSlice returns the list at key as a slice of object maps; a missing
// key or non-list value yields an empty slice, never an error.
func maxComputeJobSlice(data map[string]interface{}, key string) []map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	if data == nil {
		return items
	}
	raw, ok := data[key].([]interface{})
	if !ok {
		return items
	}
	for _, v := range raw {
		if item, ok := v.(map[string]interface{}); ok {
			items = append(items, item)
		}
	}
	return items
}

func maxComputeJobString(item map[string]interface{}, key string) string {
	if v, ok := item[key].(string); ok {
		return v
	}
	return ""
}

// maxComputeJobInt coerces a JSON number to int64. The connectivity layer
// decodes response bodies into map[string]interface{}, so numbers arrive as
// float64; json.Number is handled as a fallback.
func maxComputeJobInt(item map[string]interface{}, key string) int64 {
	switch v := item[key].(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return n
		}
	}
	return 0
}
