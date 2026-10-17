package paifs

// Package paifs registers the PaiFeatureStore service's Framework-native
// resources and data sources with the provider assembly loop.
//
// FeatureStoreTask is a read-only object: the upstream PaiFeatureStore
// 2023-06-21 API exposes GetTask/ListTasks (read-only) and StopTask/ListTaskLogs
// (auxiliary) but no Create/Delete/Update lifecycle, so it is surfaced as a
// data source only. StopTask/ListLogs are write/auxiliary actions and are not
// exposed on a read-only data source.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	paiFeatureStoreTasksTypeName = "alicloud_pai_feature_store_tasks"
	paiFeatureStoreProduct       = "PaiFeatureStore"
	paiFeatureStoreAPIVersion    = "2023-06-21"
	// pageTasksSize bounds a single ListTasks page. TotalCount, an empty page
	// and a short page each terminate the paging loop.
	pageTasksSize = 50
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
	Tasks      types.List   `tfsdk:"tasks"`
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

// taskAttrTypes is the element type of the tasks list. Protocol 5 cannot carry
// nested attributes, so the list is declared as a list of plain objects.
var taskAttrTypes = types.ObjectType{AttrTypes: map[string]attr.Type{
	"task_id":        types.StringType,
	"instance_id":    types.StringType,
	"project_id":     types.StringType,
	"project_name":   types.StringType,
	"object_type":    types.StringType,
	"object_id":      types.StringType,
	"status":         types.StringType,
	"type":           types.StringType,
	"create_time":    types.StringType,
	"executed_time":  types.StringType,
	"finished_time":  types.StringType,
	"modified_time":  types.StringType,
	"config":         types.StringType,
	"running_config": types.StringType,
}}

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
				Description: "The ID of a single task. When set, the data source returns one task fetched by GetTask instead of listing, and the filter arguments are ignored.",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the ID of the project. Ignored when task_id is set.",
			},
			"object_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the object type of the task. Valid values: ModelFeature, FeatureView. Ignored when task_id is set.",
				Validators: []validator.String{
					stringvalidator.OneOf("ModelFeature", "FeatureView"),
				},
			},
			"object_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the ID of the task object. Ignored when task_id is set.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the status of the task. Valid values: Initializing, Running, Success, Failure. Ignored when task_id is set.",
				Validators: []validator.String{
					stringvalidator.OneOf("Initializing", "Running", "Success", "Failure"),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter: the type of the task. Valid values: OfflineToOnline, ExportTrainingSet. Ignored when task_id is set.",
				Validators: []validator.String{
					stringvalidator.OneOf("OfflineToOnline", "ExportTrainingSet"),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source. It is the instance_id when listing, or instance_id/task_id when fetching a single task.",
			},
			"tasks": schema.ListAttribute{
				Computed:    true,
				ElementType: taskAttrTypes,
				Description: "A list of FeatureStoreTask resources.",
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
		Tasks:      taskAttrList(tasks),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// paiRoaClient is the seam the unit tests fake; *connectivity.AliyunClient
// satisfies it because RoaGet has this signature.
type paiRoaClient interface {
	RoaGet(apiProductCode string, apiVersion string, pathName string, query map[string]*string, headers map[string]*string, body interface{}) (map[string]interface{}, error)
}

type listFilters struct {
	ProjectID  types.String
	ObjectType types.String
	ObjectID   types.String
	Status     types.String
	Type       types.String
}

// fetchTasks is the core read logic, factored out so unit tests can drive it
// with a fake client and no Terraform plumbing. PaiFeatureStore is a ROA
// service: the instance id travels in the path and filters in the query.
func fetchTasks(_ context.Context, client paiRoaClient, instanceID, taskID string, f listFilters) (string, []taskAttr, diag.Diagnostics) {
	var diags diag.Diagnostics

	if instanceID == "" {
		diags.AddError(
			fmt.Sprintf("Reading %s: missing instance_id", paiFeatureStoreTasksTypeName),
			"instance_id is required: FeatureStoreTask belongs to a PaiFeatureStore instance.",
		)
		return "", nil, diags
	}

	if taskID != "" {
		path := fmt.Sprintf("%s/%s", tasksPath(instanceID), taskID)
		response, err := client.RoaGet(paiFeatureStoreProduct, paiFeatureStoreAPIVersion, path, nil, nil, nil)
		if err != nil {
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

// tasksPath is the ROA resource path shared by ListTasks and GetTask.
func tasksPath(instanceID string) string {
	return fmt.Sprintf("/api/v1/instances/%s/tasks", instanceID)
}

// listAllTasks pages through ListTasks until the accumulated rows reach the
// reported TotalCount, or a short or empty page signals the end. Order is
// preserved as returned by the API.
func listAllTasks(client paiRoaClient, instanceID string, f listFilters) ([]taskAttr, error) {
	var out []taskAttr
	path := tasksPath(instanceID)
	for pageNumber := 1; ; pageNumber++ {
		query := buildListQuery(f)
		query["PageSize"] = stringPtr(strconv.Itoa(pageTasksSize))
		query["PageNumber"] = stringPtr(strconv.Itoa(pageNumber))

		response, err := client.RoaGet(paiFeatureStoreProduct, paiFeatureStoreAPIVersion, path, query, nil, nil)
		if err != nil {
			return nil, err
		}

		rows, err := taskRows(response)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			out = append(out, mapListTaskToTaskAttr(row, instanceID))
		}

		if len(rows) == 0 || len(rows) < pageTasksSize {
			break
		}
		if total, ok := intValue(response["TotalCount"]); ok && len(out) >= total {
			break
		}
	}
	return out, nil
}

// taskRows extracts the Tasks array of a ListTasks response. A response shape
// that does not match the contract is an error: silently returning no rows
// would make a false empty result indistinguishable from a real one.
func taskRows(response map[string]interface{}) ([]map[string]interface{}, error) {
	raw, ok := response["Tasks"]
	if !ok || raw == nil {
		return nil, nil
	}
	items, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected ListTasks response: Tasks is %T, want array", raw)
	}
	rows := make([]map[string]interface{}, 0, len(items))
	for i, item := range items {
		row, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected ListTasks response: Tasks[%d] is %T, want object", i, item)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// intValue reads an integer response field regardless of the concrete numeric
// type the JSON decoder produced.
func intValue(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	case string:
		i, err := strconv.Atoi(n)
		if err != nil {
			return 0, false
		}
		return i, true
	}
	return 0, false
}

// buildListQuery assembles the ListTasks query from the optional filters.
// Absent filters are omitted so the API applies its default scope rather than
// an empty-string filter.
func buildListQuery(f listFilters) map[string]*string {
	query := make(map[string]*string, 5)
	addQueryParam(query, "ProjectId", f.ProjectID)
	addQueryParam(query, "ObjectType", f.ObjectType)
	addQueryParam(query, "ObjectId", f.ObjectID)
	addQueryParam(query, "Status", f.Status)
	addQueryParam(query, "Type", f.Type)
	return query
}

func addQueryParam(query map[string]*string, key string, value types.String) {
	if v := value.ValueString(); v != "" {
		query[key] = &v
	}
}

func stringPtr(v string) *string {
	return &v
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

func taskAttrList(rows []taskAttr) types.List {
	values := make([]attr.Value, 0, len(rows))
	for _, r := range rows {
		values = append(values, types.ObjectValueMust(taskAttrTypes.AttrTypes, map[string]attr.Value{
			"task_id":        r.TaskID,
			"instance_id":    r.InstanceID,
			"project_id":     r.ProjectID,
			"project_name":   r.ProjectName,
			"object_type":    r.ObjectType,
			"object_id":      r.ObjectID,
			"status":         r.Status,
			"type":           r.Type,
			"create_time":    r.CreateTime,
			"executed_time":  r.ExecutedTime,
			"finished_time":  r.FinishedTime,
			"modified_time":  r.ModifiedTime,
			"config":         r.Config,
			"running_config": r.RunningConfig,
		}))
	}
	return types.ListValueMust(taskAttrTypes, values)
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
