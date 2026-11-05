package threatdetection

import (
	"context"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const clientFileProtectEventsTypeName = "alicloud_threat_detection_client_file_protect_events"

const (
	clientFileProtectEventsAction   = "ListFileProtectClientEvent"
	clientFileProtectEventsPageSize = helper.PageSizeLarge
	// clientFileProtectEventsMaxPages bounds a pagination loop that never
	// terminates. The file protect client events are bounded by the time window
	// and the product quotas, so reaching this limit means the endpoint ignored
	// the page parameters; returning the accumulated rows would be partial.
	clientFileProtectEventsMaxPages = 1000
)

var (
	_ datasource.DataSource              = &clientFileProtectEventsDataSource{}
	_ datasource.DataSourceWithConfigure = &clientFileProtectEventsDataSource{}
)

func NewClientFileProtectEventsDataSource() datasource.DataSource {
	return &clientFileProtectEventsDataSource{}
}

type clientFileProtectEventsDataSource struct {
	fwadapt.DataSourceBase
}

// clientFileProtectEventsModel is the data source model. The tfsdk tags map
// fields to the schema attributes; Events is a slice of nested object models
// so the framework can serialize the list directly.
type clientFileProtectEventsModel struct {
	Id           types.String                  `tfsdk:"id"`
	RuleName     types.String                  `tfsdk:"rule_name"`
	Status       types.Int64                   `tfsdk:"status"`
	AlertLevels  types.List                    `tfsdk:"alert_levels"`
	InstanceId   types.String                  `tfsdk:"instance_id"`
	Uuid         types.String                  `tfsdk:"uuid"`
	InstanceName types.String                  `tfsdk:"instance_name"`
	InternetIp   types.String                  `tfsdk:"internet_ip"`
	IntranetIp   types.String                  `tfsdk:"intranet_ip"`
	StartTime    types.Int64                   `tfsdk:"start_time"`
	EndTime      types.Int64                   `tfsdk:"end_time"`
	Operation    types.String                  `tfsdk:"operation"`
	ProcPath     types.String                  `tfsdk:"proc_path"`
	FilePath     types.String                  `tfsdk:"file_path"`
	Ids          types.List                    `tfsdk:"ids"`
	PageSize     types.Int64                   `tfsdk:"page_size"`
	Events       []clientFileProtectEventModel `tfsdk:"events"`
}

// clientFileProtectEventModel is one file protect client event returned by
// ListFileProtectClientEvent. The Id is stringified so the ids filter and the
// nested id stay consistent regardless of the JSON number encoding.
type clientFileProtectEventModel struct {
	Id             types.String `tfsdk:"id"`
	Uuid           types.String `tfsdk:"uuid"`
	RuleName       types.String `tfsdk:"rule_name"`
	Status         types.Int64  `tfsdk:"status"`
	AlertLevel     types.Int64  `tfsdk:"alert_level"`
	InstanceName   types.String `tfsdk:"instance_name"`
	InternetIp     types.String `tfsdk:"internet_ip"`
	IntranetIp     types.String `tfsdk:"intranet_ip"`
	FilePath       types.String `tfsdk:"file_path"`
	Operation      types.String `tfsdk:"operation"`
	Remark         types.String `tfsdk:"remark"`
	ProcPermission types.String `tfsdk:"proc_permission"`
	CmdLine        types.String `tfsdk:"cmd_line"`
	ProcPath       types.String `tfsdk:"proc_path"`
	HandleTime     types.Int64  `tfsdk:"handle_time"`
	FirstTime      types.Int64  `tfsdk:"first_time"`
	LatestTime     types.Int64  `tfsdk:"latest_time"`
	ProcessId      types.String `tfsdk:"process_id"`
	Platform       types.String `tfsdk:"platform"`
	RuleAction     types.String `tfsdk:"rule_action"`
	Count          types.Int64  `tfsdk:"count"`
}

var clientFileProtectEventObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":              types.StringType,
	"uuid":            types.StringType,
	"rule_name":       types.StringType,
	"status":          types.Int64Type,
	"alert_level":     types.Int64Type,
	"instance_name":   types.StringType,
	"internet_ip":     types.StringType,
	"intranet_ip":     types.StringType,
	"file_path":       types.StringType,
	"operation":       types.StringType,
	"remark":          types.StringType,
	"proc_permission": types.StringType,
	"cmd_line":        types.StringType,
	"proc_path":       types.StringType,
	"handle_time":     types.Int64Type,
	"first_time":      types.Int64Type,
	"latest_time":     types.Int64Type,
	"process_id":      types.StringType,
	"platform":        types.StringType,
	"rule_action":     types.StringType,
	"count":           types.Int64Type,
}}

func (d *clientFileProtectEventsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = clientFileProtectEventsTypeName
}

func (d *clientFileProtectEventsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the file protect client events of Cloud Security Center (Threat Detection). The events are system-generated security detections and are read-only; this data source queries `ListFileProtectClientEvent` and returns all matching events.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, derived from the configured filters so it stays stable across refreshes.",
			},
			"rule_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the file protect rule. The value is sent to the server as the `RuleName` filter.",
			},
			"status": schema.Int64Attribute{
				Optional:    true,
				Description: "The status of the event. The value is sent to the server as the `Status` filter.",
			},
			"alert_levels": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "The alert levels to filter. Valid values: `0` (no alert), `1` (reminder), `2` (suspicious) and `3` (high risk). The value is sent to the server as the `AlertLevels` repeatList filter.",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the instance. The value is sent to the server as the `InstanceId` filter.",
			},
			"uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the instance. The value is sent to the server as the `Uuid` filter.",
			},
			"instance_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the instance. The value is sent to the server as the `InstanceName` filter.",
			},
			"internet_ip": schema.StringAttribute{
				Optional:    true,
				Description: "The internet IP of the instance. The value is sent to the server as the `InternetIp` filter.",
			},
			"intranet_ip": schema.StringAttribute{
				Optional:    true,
				Description: "The intranet IP of the instance. The value is sent to the server as the `IntranetIp` filter.",
			},
			"start_time": schema.Int64Attribute{
				Optional:    true,
				Description: "The start time of the query range, as a Unix timestamp in milliseconds. The value is sent to the server as the `StartTime` filter.",
			},
			"end_time": schema.Int64Attribute{
				Optional:    true,
				Description: "The end time of the query range, as a Unix timestamp in milliseconds. The value is sent to the server as the `EndTime` filter.",
			},
			"operation": schema.StringAttribute{
				Optional:    true,
				Description: "The file operation. The value is sent to the server as the `Operation` filter.",
			},
			"proc_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path of the process. The value is sent to the server as the `ProcPath` filter.",
			},
			"file_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path of the file. The value is sent to the server as the `FilePath` filter.",
			},
			"ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A list of event IDs to filter. The filter is applied on the provider side after all pages are read.",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of entries per page when querying the server. If unset, the provider uses the default page size of the product. The provider reads all pages regardless of this value.",
			},
			"events": schema.ListAttribute{
				Computed:    true,
				ElementType: clientFileProtectEventObjectType,
				Description: "A list of file protect client events that match the configured filters. Each element contains the fields returned by ListFileProtectClientEvent.",
			},
		},
	}
}

func (d *clientFileProtectEventsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model clientFileProtectEventsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	idsFilter := make(map[string]struct{})
	if !model.Ids.IsNull() && !model.Ids.IsUnknown() {
		var ids []string
		resp.Diagnostics.Append(model.Ids.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, id := range ids {
			if id == "" {
				continue
			}
			idsFilter[id] = struct{}{}
		}
	}

	var alertLevels []int64
	if !model.AlertLevels.IsNull() && !model.AlertLevels.IsUnknown() {
		resp.Diagnostics.Append(model.AlertLevels.ElementsAs(ctx, &alertLevels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	pageSize := int(clientFileProtectEventsPageSize)
	if !model.PageSize.IsNull() && !model.PageSize.IsUnknown() && model.PageSize.ValueInt64() > 0 {
		pageSize = int(model.PageSize.ValueInt64())
	}

	items, err := d.collectAllEvents(ctx, &model, alertLevels, pageSize)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", clientFileProtectEventsTypeName, clientFileProtectEventsAction),
			err.Error(),
		)
		return
	}

	events := make([]clientFileProtectEventModel, 0, len(items))
	ids := make([]string, 0, len(items))
	for _, item := range items {
		id := clientFileProtectEventString(item, "Id")
		if len(idsFilter) > 0 {
			if _, ok := idsFilter[id]; !ok {
				continue
			}
		}
		ids = append(ids, id)
		events = append(events, clientFileProtectEventModel{
			Id:             types.StringValue(id),
			Uuid:           types.StringValue(clientFileProtectEventString(item, "Uuid")),
			RuleName:       types.StringValue(clientFileProtectEventString(item, "RuleName")),
			Status:         types.Int64Value(clientFileProtectEventInt64(item, "Status")),
			AlertLevel:     types.Int64Value(clientFileProtectEventInt64(item, "AlertLevel")),
			InstanceName:   types.StringValue(clientFileProtectEventString(item, "InstanceName")),
			InternetIp:     types.StringValue(clientFileProtectEventString(item, "InternetIp")),
			IntranetIp:     types.StringValue(clientFileProtectEventString(item, "IntranetIp")),
			FilePath:       types.StringValue(clientFileProtectEventString(item, "FilePath")),
			Operation:      types.StringValue(clientFileProtectEventString(item, "Operation")),
			Remark:         types.StringValue(clientFileProtectEventString(item, "Remark")),
			ProcPermission: types.StringValue(clientFileProtectEventString(item, "ProcPermission")),
			CmdLine:        types.StringValue(clientFileProtectEventString(item, "CmdLine")),
			ProcPath:       types.StringValue(clientFileProtectEventString(item, "ProcPath")),
			HandleTime:     types.Int64Value(clientFileProtectEventInt64(item, "HandleTime")),
			FirstTime:      types.Int64Value(clientFileProtectEventInt64(item, "FirstTime")),
			LatestTime:     types.Int64Value(clientFileProtectEventInt64(item, "LatestTime")),
			ProcessId:      types.StringValue(clientFileProtectEventString(item, "ProcessId")),
			Platform:       types.StringValue(clientFileProtectEventString(item, "Platform")),
			RuleAction:     types.StringValue(clientFileProtectEventString(item, "RuleAction")),
			Count:          types.Int64Value(clientFileProtectEventInt64(item, "Count")),
		})
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.Id = types.StringValue(clientFileProtectEventsId(&model))
	model.Ids = idsList
	model.Events = events
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// collectAllEvents reads every page of ListFileProtectClientEvent. The API has
// no cursor, so the loop stops at the first page shorter than PageSize; an
// exact multiple of PageSize therefore costs one extra, empty page.
func (d *clientFileProtectEventsDataSource) collectAllEvents(ctx context.Context, model *clientFileProtectEventsModel, alertLevels []int64, pageSize int) ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	for pageNumber := 1; pageNumber <= clientFileProtectEventsMaxPages; pageNumber++ {
		pageItems, err := d.readEventsPage(ctx, model, alertLevels, pageNumber, pageSize)
		if err != nil {
			return nil, err
		}
		items = append(items, pageItems...)
		if len(pageItems) < pageSize {
			return items, nil
		}
	}
	return nil, fmt.Errorf("%s returned a full page for %d consecutive requests (PageSize=%d), so pagination did not terminate; the result was not returned because it would be incomplete", clientFileProtectEventsAction, clientFileProtectEventsMaxPages, pageSize)
}

func (d *clientFileProtectEventsDataSource) readEventsPage(ctx context.Context, model *clientFileProtectEventsModel, alertLevels []int64, pageNumber, pageSize int) ([]map[string]interface{}, error) {
	request := map[string]interface{}{
		"CurrentPage": pageNumber,
		"PageSize":    pageSize,
	}
	if v := model.RuleName.ValueString(); !model.RuleName.IsNull() && v != "" {
		request["RuleName"] = v
	}
	if !model.Status.IsNull() && !model.Status.IsUnknown() {
		request["Status"] = model.Status.ValueInt64()
	}
	if len(alertLevels) > 0 {
		levels := make([]interface{}, 0, len(alertLevels))
		for _, lv := range alertLevels {
			levels = append(levels, lv)
		}
		request["AlertLevels"] = levels
	}
	if v := model.InstanceId.ValueString(); !model.InstanceId.IsNull() && v != "" {
		request["InstanceId"] = v
	}
	if v := model.Uuid.ValueString(); !model.Uuid.IsNull() && v != "" {
		request["Uuid"] = v
	}
	if v := model.InstanceName.ValueString(); !model.InstanceName.IsNull() && v != "" {
		request["InstanceName"] = v
	}
	if v := model.InternetIp.ValueString(); !model.InternetIp.IsNull() && v != "" {
		request["InternetIp"] = v
	}
	if v := model.IntranetIp.ValueString(); !model.IntranetIp.IsNull() && v != "" {
		request["IntranetIp"] = v
	}
	if !model.StartTime.IsNull() && !model.StartTime.IsUnknown() {
		request["StartTime"] = model.StartTime.ValueInt64()
	}
	if !model.EndTime.IsNull() && !model.EndTime.IsUnknown() {
		request["EndTime"] = model.EndTime.ValueInt64()
	}
	if v := model.Operation.ValueString(); !model.Operation.IsNull() && v != "" {
		request["Operation"] = v
	}
	if v := model.ProcPath.ValueString(); !model.ProcPath.IsNull() && v != "" {
		request["ProcPath"] = v
	}
	if v := model.FilePath.ValueString(); !model.FilePath.IsNull() && v != "" {
		request["FilePath"] = v
	}

	// ListFileProtectClientEvent is a GET RPC (methods: ["get"], operationType:
	// "read"); all filters travel as query parameters. RetryRpcPost's retry loop
	// is method-agnostic, so it wraps the GET call the same way it wraps POST.
	response, err := helper.RetryRpcPost(ctx, clientFileProtectEventsAction, func() (map[string]interface{}, error) {
		return d.Meta().RpcGet("Sas", "2018-12-03", clientFileProtectEventsAction, request, nil)
	}, 5*time.Minute, helper.NeedRetry)
	logs.AddDebug(clientFileProtectEventsAction, response, request)
	if err != nil {
		return nil, err
	}
	return clientFileProtectEventsItems(response), nil
}

// clientFileProtectEventsItems parses the EventList array. A response without
// the key or with an empty array yields an empty slice; it is not an error.
func clientFileProtectEventsItems(response map[string]interface{}) []map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	if response == nil {
		return items
	}
	raw, err := jsonpath.Get("$.EventList", response)
	if err != nil {
		return items
	}
	list, ok := raw.([]interface{})
	if !ok {
		return items
	}
	for _, v := range list {
		if item, ok := v.(map[string]interface{}); ok {
			items = append(items, item)
		}
	}
	return items
}

// clientFileProtectEventString returns the string value of a field. JSON
// numbers are stringified so id-like fields stay stable regardless of the
// number encoding the API chooses.
func clientFileProtectEventString(item map[string]interface{}, key string) string {
	v, ok := item[key]
	if !ok || v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int:
		return fmt.Sprint(val)
	case int64:
		return fmt.Sprint(val)
	default:
		return fmt.Sprint(val)
	}
}

// clientFileProtectEventInt64 returns the int64 value of a numeric field. JSON
// numbers arrive as float64, so the conversion goes through float64.
func clientFileProtectEventInt64(item map[string]interface{}, key string) int64 {
	v, ok := item[key]
	if !ok || v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	case int64:
		return val
	case string:
		var n int64
		_, _ = fmt.Sscanf(val, "%d", &n)
		return n
	default:
		return 0
	}
}

func clientFileProtectEventsId(model *clientFileProtectEventsModel) string {
	var alertLevels []string
	_ = model.AlertLevels.ElementsAs(context.Background(), &alertLevels, false)
	return helper.DataResourceIdHash([]string{
		model.RuleName.ValueString(),
		fmt.Sprint(model.Status.ValueInt64()),
		fmt.Sprint(alertLevels),
		model.InstanceId.ValueString(),
		model.Uuid.ValueString(),
		model.InstanceName.ValueString(),
		model.InternetIp.ValueString(),
		model.IntranetIp.ValueString(),
		fmt.Sprint(model.StartTime.ValueInt64()),
		fmt.Sprint(model.EndTime.ValueInt64()),
		model.Operation.ValueString(),
		model.ProcPath.ValueString(),
		model.FilePath.ValueString(),
	})
}
