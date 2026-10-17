package paifs

// Package paifs registers the PaiFeatureStore service's Framework-native
// resources and data sources with the provider assembly loop.
//
// FeatureStoreTask is a read-only object: the upstream PaiFeatureStore
// 2023-06-21 API exposes GetTask/ListTasks (query) and StopTask/ListTaskLogs
// (auxiliary) but no Create/Delete/Update lifecycle, so it is surfaced as a
// data source only. StopTask/ListLogs are write/auxiliary actions and are not
// exposed on a read-only data source.

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	paiFeatureStoreTasksTypeName = "alicloud_pai_feature_store_tasks"
	paiFeatureStoreProduct       = "PaiFeatureStore"
	paiFeatureStoreAPIVersion    = "2023-06-21"
	// pageTasksSize bounds a single ListTasks page. The API reports TotalCount
	// but the loop terminates on a short page, so this is a fetch width only.
	pageTasksSize = 50
	// taskNotFoundCode is returned by GetTask/ListTasks when the targeted
	// FeatureStore instance or task does not exist; the data source treats it
	// as "absent" rather than a hard read failure.
	taskNotFoundCode = "DoConfigCenterRequestErrorProblem"
)

var (
	_ datasource.DataSource              = &tasksDataSource{}
	_ datasource.DataSourceWithConfigure = &tasksDataSource{}
)

// NewPaiFeatureStoreTasksDataSource is the registration factory.
func NewPaiFeatureStoreTasksDataSource() datasource.DataSource {
	return &tasksDataSource{}
}

type tasksDataSource struct {
	fwadapt.DataSourceBase
}

// tasksDataSourceModel serves both config and state: input filters are
// Optional, the computed outputs are id and tasks.
type tasksDataSourceModel struct {
	InstanceID types.String `tfsdk:"instance_id"`
	TaskID     types.String `tfsdk:"task_id"`
	ProjectID  types.String `tfsdk:"project_id"`
	ObjectType types.String `tfsdk:"object_type"`
	ObjectID   types.String `tfsdk:"object_id"`
	Status     types.String `tfsdk:"status"`
	Type       types.String `tfsdk:"type"`
	ID         types.String `tfsdk:"id"`
	Tasks      []taskAttr   `tfsdk:"tasks"`
}

// taskAttr is one FeatureStoreTask row. Gmt* response fields are mapped to
// *_time attributes; instance_id is passed through from the input because the
// Task object itself does not carry it.
type taskAttr struct {
	TaskID        types.String `tfsdk:"task_id"`
	InstanceID    types.String `tfsdk:"instance_id"`
	ProjectID     types.String `tfsdk:"project_id"`
	ProjectName   types.String `tfsdk:"project_name"`
	ObjectType    types.String `tfsdk:"object_type"`
	ObjectID      types.String `tfsdk:"object_id"`
	Status        types.String `tfsdk:"status"`
	Type          types.String `tfsdk:"type"`
	CreateTime    types.String `tfsdk:"create_time"`
	ExecutedTime  types.String `tfsdk:"executed_time"`
	FinishedTime  types.String `tfsdk:"finished_time"`
	ModifiedTime  types.String `tfsdk:"modified_time"`
	Config        types.String `tfsdk:"config"`
	RunningConfig types.String `tfsdk:"running_config"`
}

func (d *tasksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = paiFeatureStoreTasksTypeName
}

func (d *tasksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides a list of PaiFeatureStore FeatureStoreTask resources in the current account by default.\n\n" +
			"When `task_id` is set, the data source fetches a single task via GetTask; otherwise it lists tasks via ListTasks with the given filters.\n\n" +
			"FeatureStoreTask is read-only (no Create/Delete/Update lifecycle API).",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the PaiFeatureStore instance that owns the tasks.",
			},
			"task_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of a single task. When set, the data source returns one task fetched by GetTask instead of listing.",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the ID of the project.",
			},
			"object_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the object type of the task. Valid values: ModelFeature, FeatureView.",
			},
			"object_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the ID of the task object.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the status of the task. Valid values: Initializing, Running, Success, Failure.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the type of the task. Valid values: OfflineToOnline, ExportTrainingSet.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source. It is the instance_id when listing, or instance_id/task_id when fetching a single task.",
			},
			"tasks": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"task_id":        schema.StringAttribute{Computed: true, Description: "The ID of the task."},
						"instance_id":    schema.StringAttribute{Computed: true, Description: "The ID of the PaiFeatureStore instance that owns the task."},
						"project_id":     schema.StringAttribute{Computed: true, Description: "The ID of the project."},
						"project_name":   schema.StringAttribute{Computed: true, Description: "The name of the project."},
						"object_type":    schema.StringAttribute{Computed: true, Description: "The object type of the task."},
						"object_id":      schema.StringAttribute{Computed: true, Description: "The ID of the task object."},
						"status":         schema.StringAttribute{Computed: true, Description: "The status of the task."},
						"type":           schema.StringAttribute{Computed: true, Description: "The type of the task."},
						"create_time":    schema.StringAttribute{Computed: true, Description: "The creation time (GmtCreateTime)."},
						"executed_time":  schema.StringAttribute{Computed: true, Description: "The execution time (GmtExecutedTime)."},
						"finished_time":  schema.StringAttribute{Computed: true, Description: "The finish time (GmtFinishedTime)."},
						"modified_time":  schema.StringAttribute{Computed: true, Description: "The modification time (GmtModifiedTime), only from GetTask."},
						"config":         schema.StringAttribute{Computed: true, Description: "The task config, only from GetTask."},
						"running_config": schema.StringAttribute{Computed: true, Description: "The task running config, only from GetTask."},
					},
				},
			},
		},
	}
}

func (d *tasksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg tasksDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()
	if client == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", paiFeatureStoreTasksTypeName),
			"The provider client was not injected. Please report this issue to the provider developers.",
		)
		return
	}

	filters := listFilters{
		ProjectID:  cfg.ProjectID,
		ObjectType: cfg.ObjectType,
		ObjectID:   cfg.ObjectID,
		Status:     cfg.Status,
		Type:       cfg.Type,
	}
	id, tasks, diags := fetchTasks(ctx, client, cfg.InstanceID.ValueString(), cfg.TaskID.ValueString(), filters)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := tasksDataSourceModel{
		InstanceID: cfg.InstanceID,
		TaskID:     cfg.TaskID,
		ProjectID:  cfg.ProjectID,
		ObjectType: cfg.ObjectType,
		ObjectID:   cfg.ObjectID,
		Status:     cfg.Status,
		Type:       cfg.Type,
		ID:         types.StringValue(id),
		Tasks:      tasks,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// paiRpcClient is the seam the unit tests fake; *connectivity.AliyunClient
// satisfies it because RpcPost has this signature. Production reads never reach
// the SDK v2 testing helper, so the test binary links at most one sweep flag.
type paiRpcClient interface {
	RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error)
}

type listFilters struct {
	ProjectID  types.String
	ObjectType types.String
	ObjectID   types.String
	Status     types.String
	Type       types.String
}

// fetchTasks is the core read logic, factored out so unit tests can drive it
// with a fake client and no Terraform plumbing.
func fetchTasks(_ context.Context, client paiRpcClient, instanceID, taskID string, f listFilters) (string, []taskAttr, diag.Diagnostics) {
	var diags diag.Diagnostics

	if instanceID == "" {
		diags.AddError(
			fmt.Sprintf("Reading %s: missing instance_id", paiFeatureStoreTasksTypeName),
			"instance_id is required: FeatureStoreTask belongs to a PaiFeatureStore instance.",
		)
		return "", nil, diags
	}

	if taskID != "" {
		query := map[string]interface{}{
			"InstanceId": instanceID,
			"TaskId":     taskID,
		}
		response, err := client.RpcPost(paiFeatureStoreProduct, paiFeatureStoreAPIVersion, "GetTask", query, map[string]interface{}{}, true)
		if err != nil {
			if isTaskNotFound(err) {
				diags.AddWarning(
					fmt.Sprintf("Reading %s", paiFeatureStoreTasksTypeName),
					fmt.Sprintf("The task %s/%s was not found: %s", instanceID, taskID, err.Error()),
				)
				return singleID(instanceID, taskID), nil, diags
			}
			diags.AddError(
				fmt.Sprintf("Reading %s: calling GetTask", paiFeatureStoreTasksTypeName),
				err.Error(),
			)
			return "", nil, diags
		}
		return singleID(instanceID, taskID), []taskAttr{mapGetTaskToTaskAttr(response, instanceID, taskID)}, diags
	}

	tasks, err := listAllTasks(client, instanceID, f)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Reading %s: calling ListTasks", paiFeatureStoreTasksTypeName),
			err.Error(),
		)
		return "", nil, diags
	}
	return instanceID, tasks, diags
}

// listAllTasks pages through ListTasks until a short page signals the end.
// Order is preserved as returned by the API.
func listAllTasks(client paiRpcClient, instanceID string, f listFilters) ([]taskAttr, error) {
	pageNumber := 1
	var out []taskAttr
	for {
		request := buildListRequest(instanceID, f)
		request["PageSize"] = pageTasksSize
		request["PageNumber"] = pageNumber

		response, err := client.RpcPost(paiFeatureStoreProduct, paiFeatureStoreAPIVersion, "ListTasks", request, map[string]interface{}{}, true)
		if err != nil {
			return nil, err
		}

		rawTasks, _ := response["Tasks"].([]interface{})
		for _, v := range rawTasks {
			m, _ := v.(map[string]interface{})
			out = append(out, mapListTaskToTaskAttr(m, instanceID))
		}

		if len(rawTasks) < pageTasksSize {
			break
		}
		pageNumber++
	}
	return out, nil
}

// buildListRequest assembles the ListTasks query from the instance id and the
// optional filters. Absent filters are omitted so the API applies its default
// scope rather than an empty-string filter.
func buildListRequest(instanceID string, f listFilters) map[string]interface{} {
	request := map[string]interface{}{"InstanceId": instanceID}
	if v := f.ProjectID.ValueString(); v != "" {
		request["ProjectId"] = v
	}
	if v := f.ObjectType.ValueString(); v != "" {
		request["ObjectType"] = v
	}
	if v := f.ObjectID.ValueString(); v != "" {
		request["ObjectId"] = v
	}
	if v := f.Status.ValueString(); v != "" {
		request["Status"] = v
	}
	if v := f.Type.ValueString(); v != "" {
		request["Type"] = v
	}
	return request
}

// isTaskNotFound reports whether an RPC error means the targeted instance/task
// is absent, so the data source can return an empty result rather than fail.
func isTaskNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), taskNotFoundCode)
}

func mapGetTaskToTaskAttr(response map[string]interface{}, instanceID, taskID string) taskAttr {
	return taskAttr{
		TaskID:        types.StringValue(taskID),
		InstanceID:    types.StringValue(instanceID),
		ProjectID:     types.StringValue(strVal(response["ProjectId"])),
		ProjectName:   types.StringValue(strVal(response["ProjectName"])),
		ObjectType:    types.StringValue(strVal(response["ObjectType"])),
		ObjectID:      types.StringValue(strVal(response["ObjectId"])),
		Status:        types.StringValue(strVal(response["Status"])),
		Type:          types.StringValue(strVal(response["Type"])),
		CreateTime:    types.StringValue(strVal(response["GmtCreateTime"])),
		ExecutedTime:  types.StringValue(strVal(response["GmtExecutedTime"])),
		FinishedTime:  types.StringValue(strVal(response["GmtFinishedTime"])),
		ModifiedTime:  types.StringValue(strVal(response["GmtModifiedTime"])),
		Config:        types.StringValue(strVal(response["Config"])),
		RunningConfig: types.StringValue(strVal(response["RunningConfig"])),
	}
}

func mapListTaskToTaskAttr(m map[string]interface{}, instanceID string) taskAttr {
	return taskAttr{
		TaskID:       types.StringValue(strVal(m["TaskId"])),
		InstanceID:   types.StringValue(instanceID),
		ProjectID:    types.StringValue(strVal(m["ProjectId"])),
		ProjectName:  types.StringValue(strVal(m["ProjectName"])),
		ObjectType:   types.StringValue(strVal(m["ObjectType"])),
		ObjectID:     types.StringValue(strVal(m["ObjectId"])),
		Status:       types.StringValue(strVal(m["Status"])),
		Type:         types.StringValue(strVal(m["Type"])),
		CreateTime:   types.StringValue(strVal(m["GmtCreateTime"])),
		ExecutedTime: types.StringValue(strVal(m["GmtExecutedTime"])),
		FinishedTime: types.StringValue(strVal(m["GmtFinishedTime"])),
		// ListTasks does not return these fields.
		ModifiedTime:  types.StringNull(),
		Config:        types.StringNull(),
		RunningConfig: types.StringNull(),
	}
}

func singleID(instanceID, taskID string) string {
	return fmt.Sprintf("%s/%s", instanceID, taskID)
}

// strVal stringifies a response value, returning "" for nil so absent fields
// map to an empty (not unknown) attribute.
func strVal(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}
