package cms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	metricListTypeName      = "alicloud_cms_metric_list"
	metricListAction        = "DescribeMetricList"
	metricListProduct       = "Cms"
	metricListAPIVersion    = "2019-01-01"
	metricListDefaultLength = 1000
	metricListRetryTimeout  = 5 * time.Minute
)

var (
	_ datasource.DataSource              = (*metricListDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*metricListDataSource)(nil)
)

func NewMetricListDataSource() datasource.DataSource {
	return &metricListDataSource{}
}

type metricListDataSource struct {
	fwadapt.DataSourceBase
}

type metricListModel struct {
	Id              types.String `tfsdk:"id"`
	MetricName      types.String `tfsdk:"metric_name"`
	MetricNamespace types.String `tfsdk:"metric_namespace"`
	Dimensions      types.String `tfsdk:"dimensions"`
	StartTime       types.String `tfsdk:"start_time"`
	EndTime         types.String `tfsdk:"end_time"`
	Period          types.String `tfsdk:"period"`
	Express         types.String `tfsdk:"express"`
	Length          types.Int64  `tfsdk:"length"`
	NextToken       types.String `tfsdk:"next_token"`
	OutputFile      types.String `tfsdk:"output_file"`
	ActualPeriod    types.String `tfsdk:"actual_period"`
	Datapoints      types.List   `tfsdk:"datapoints"`
}

var metricListDatapointAttrTypes = map[string]attr.Type{
	"uuid":        types.StringType,
	"metric_name": types.StringType,
	"timestamp":   types.StringType,
	"average":     types.Float64Type,
	"sum":         types.Float64Type,
	"maximum":     types.Float64Type,
	"minimum":     types.Float64Type,
	"count":       types.Float64Type,
	"user_id":     types.StringType,
	"instance_id": types.StringType,
}

var metricListDatapointType = types.ObjectType{AttrTypes: metricListDatapointAttrTypes}

func (d *metricListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = metricListTypeName
}

func (d *metricListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// The datapoint list is a ListAttribute of ObjectType rather than a
	// NestedAttribute: protocol version 5 (what the SDKv2-shared server serves)
	// rejects nested attributes.
	resp.Schema = schema.Schema{
		Description: "Provides the datapoints of a specified Cloud Monitor metric within a time range.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, computed from the metric name and namespace.",
			},
			"metric_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the monitoring metric.",
			},
			"metric_namespace": schema.StringAttribute{
				Required:    true,
				Description: "The namespace of the monitoring metric.",
			},
			"dimensions": schema.StringAttribute{
				Optional:    true,
				Description: "The dimensions that classify the monitoring metric, as a JSON string.",
			},
			"start_time": schema.StringAttribute{
				Optional:    true,
				Description: "The beginning of the time range to query.",
			},
			"end_time": schema.StringAttribute{
				Optional:    true,
				Description: "The end of the time range to query.",
			},
			"period": schema.StringAttribute{
				Optional:    true,
				Description: "The interval at which monitoring data is collected. Unit: seconds.",
			},
			"express": schema.StringAttribute{
				Optional:    true,
				Description: "The aggregation method that is used to display the monitoring data, as a JSON string.",
			},
			"length": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum number of entries returned in a single query. Defaults to 1000.",
			},
			"next_token": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The pagination token used to retrieve the next page of datapoints.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save the data source results.",
			},
			"actual_period": schema.StringAttribute{
				Computed:    true,
				Description: "The actual period at which monitoring data is collected.",
			},
			"datapoints": schema.ListAttribute{
				Computed:    true,
				ElementType: metricListDatapointType,
				Description: "A list of metric datapoints.",
			},
		},
	}
}

func (d *metricListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state metricListModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"MetricName": state.MetricName.ValueString(),
		"Namespace":  state.MetricNamespace.ValueString(),
		"Length":     metricListDefaultLength,
	}
	setMetricListRequestString(request, "Dimensions", state.Dimensions)
	setMetricListRequestString(request, "StartTime", state.StartTime)
	setMetricListRequestString(request, "EndTime", state.EndTime)
	setMetricListRequestString(request, "Period", state.Period)
	setMetricListRequestString(request, "Express", state.Express)
	if !state.Length.IsNull() && !state.Length.IsUnknown() {
		request["Length"] = int(state.Length.ValueInt64())
	}
	if !state.NextToken.IsNull() && !state.NextToken.IsUnknown() && state.NextToken.ValueString() != "" {
		request["NextToken"] = state.NextToken.ValueString()
	}

	var objects []map[string]interface{}
	var callResponse map[string]interface{}
	var callErr error
	currentToken, _ := request["NextToken"].(string)
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(metricListRetryTimeout, func() *retry.RetryError {
			callResponse, callErr = d.Meta().RpcPost(metricListProduct, metricListAPIVersion, metricListAction, nil, request, false)
			if callErr != nil {
				// Only transient failures are retried: timeouts, throttling,
				// 5xx and InternalError. A 400/403/404 (including the 10002
				// "metric is not exist" family) is terminal and must fail or
				// resolve on the first request instead of burning the window.
				if metricListRetryable(callErr) {
					wait()
					return retry.RetryableError(callErr)
				}
				return retry.NonRetryableError(callErr)
			}
			return nil
		})
		logs.AddDebug(metricListAction, callResponse, request)
		if err != nil {
			if metricListNotExist(err) {
				d.setMetricListEmptyResult(ctx, state, resp)
				return
			}
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling %s", metricListTypeName, metricListAction),
				err.Error(),
			)
			return
		}
		if metricListResponseNotExist(callResponse) {
			d.setMetricListEmptyResult(ctx, state, resp)
			return
		}
		if err := metricListResponseError(callResponse); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling %s", metricListTypeName, metricListAction),
				err.Error(),
			)
			return
		}

		if datapointsStr, ok := callResponse["Datapoints"].(string); ok && datapointsStr != "" {
			var datapoints []map[string]interface{}
			if jsonErr := json.Unmarshal([]byte(datapointsStr), &datapoints); jsonErr != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Reading %s: decoding Datapoints of %s", metricListTypeName, metricListAction),
					jsonErr.Error(),
				)
				return
			}
			objects = append(objects, datapoints...)
		}

		nextToken := metricListNextToken(callResponse)
		if nextToken == "" || nextToken == currentToken {
			break
		}
		currentToken = nextToken
		request["NextToken"] = nextToken
	}

	records := make([]map[string]interface{}, 0, len(objects))
	items := make([]attr.Value, 0, len(objects))
	metricName := state.MetricName.ValueString()
	for _, dp := range objects {
		// The API does not echo MetricName in datapoint objects; datapoints
		// only exist for the queried metric, so backfill the request value.
		name := metricName
		if v := dp["MetricName"]; v != nil {
			name = metricListString(v)
		}
		record := map[string]interface{}{
			"uuid":        metricListString(dp["uuid"]),
			"metric_name": name,
			"timestamp":   metricListTimestampString(dp["timestamp"]),
			"average":     metricListNumber(dp["Average"]),
			"sum":         metricListNumber(dp["Sum"]),
			"maximum":     metricListNumber(dp["Maximum"]),
			"minimum":     metricListNumber(dp["Minimum"]),
			"count":       metricListNumber(dp["Count"]),
			"user_id":     metricListString(dp["userId"]),
			"instance_id": metricListString(dp["instanceId"]),
		}
		records = append(records, record)
		items = append(items, types.ObjectValueMust(metricListDatapointAttrTypes, map[string]attr.Value{
			"uuid":        types.StringValue(record["uuid"].(string)),
			"metric_name": types.StringValue(name),
			"timestamp":   types.StringValue(record["timestamp"].(string)),
			"average":     metricListFloatValue(dp["Average"]),
			"sum":         metricListFloatValue(dp["Sum"]),
			"maximum":     metricListFloatValue(dp["Maximum"]),
			"minimum":     metricListFloatValue(dp["Minimum"]),
			"count":       metricListFloatValue(dp["Count"]),
			"user_id":     types.StringValue(record["user_id"].(string)),
			"instance_id": types.StringValue(record["instance_id"].(string)),
		}))
	}

	datapoints, diags := types.ListValue(metricListDatapointType, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Id = types.StringValue(helper.DataResourceIdHash([]string{state.MetricName.ValueString(), state.MetricNamespace.ValueString()}))
	state.Datapoints = datapoints
	state.NextToken = types.StringValue(metricListNextToken(callResponse))
	if v, ok := callResponse["Period"]; ok && v != nil {
		state.ActualPeriod = types.StringValue(fmt.Sprint(v))
	} else {
		state.ActualPeriod = types.StringValue("")
	}

	if err := metricListWriteOutputFile(state.OutputFile, records); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: writing output_file", metricListTypeName),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// setMetricListEmptyResult records a legal empty query result. The 10002
// response means the metric does not exist for the namespace, which for a
// read-only data source is an empty result rather than a failure. Every
// computed value is cleared because a re-read must not keep the previous
// query's datapoints, and output_file must be overwritten with "[]" so a
// stale export is never presented as the current result.
func (d *metricListDataSource) setMetricListEmptyResult(ctx context.Context, state metricListModel, resp *datasource.ReadResponse) {
	empty, diags := types.ListValue(metricListDatapointType, []attr.Value{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Id = types.StringValue(helper.DataResourceIdHash([]string{state.MetricName.ValueString(), state.MetricNamespace.ValueString()}))
	state.Datapoints = empty
	state.NextToken = types.StringValue("")
	state.ActualPeriod = types.StringValue("")
	if err := metricListWriteOutputFile(state.OutputFile, make([]map[string]interface{}, 0)); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: writing output_file", metricListTypeName),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func setMetricListRequestString(request map[string]interface{}, key string, value types.String) {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return
	}
	request[key] = value.ValueString()
}

func metricListWriteOutputFile(output types.String, records []map[string]interface{}) error {
	if output.IsNull() || output.IsUnknown() || output.ValueString() == "" {
		return nil
	}
	if err := helper.WriteToFile(output.ValueString(), records); err != nil {
		return err
	}
	return nil
}

// metricListRetryable limits retries to transient failures: helper.NeedRetry
// covers timeouts, throttling and 5xx, and InternalError is the CMS-specific
// transient code. Everything else fails after the first request.
func metricListRetryable(err error) bool {
	if err == nil {
		return false
	}
	if helper.NeedRetry(err) {
		return true
	}
	var sdkErr *tea.SDKError
	if errors.As(err, &sdkErr) && sdkErr.Code != nil {
		return strings.TrimSpace(*sdkErr.Code) == "InternalError"
	}
	return false
}

// metricListNotExist reports only the CMS 10002 response. A broad message
// match could hide a different terminal error and overwrite output_file.
func metricListNotExist(err error) bool {
	if err == nil {
		return false
	}
	var sdkErr *tea.SDKError
	if errors.As(err, &sdkErr) {
		if sdkErr.Code != nil && strings.TrimSpace(*sdkErr.Code) == "10002" {
			return true
		}
		return sdkErr.Message != nil && metricListWrappedNotExistMessage(*sdkErr.Message)
	}
	return false
}

// CMS RPC can wrap business code 10002 in transport code 400. Recognize only
// its complete metric-not-found message, never an arbitrary "is not exist".
func metricListWrappedNotExistMessage(message string) bool {
	return strings.Contains(message, "[400-10002]the metric(") &&
		strings.Contains(message, ") of project(") &&
		strings.Contains(message, " is not exist.")
}

// metricListResponseNotExist covers the same 10002 answer when it arrives as
// a response body instead of an error.
func metricListResponseNotExist(response map[string]interface{}) bool {
	if response == nil {
		return false
	}
	if code, ok := response["Code"]; ok && code != nil && strings.TrimSpace(fmt.Sprint(code)) == "10002" {
		return true
	}
	message, ok := response["Message"].(string)
	return ok && metricListWrappedNotExistMessage(message)
}

// metricListResponseError reports a non-success response Code that arrives in
// the body of an otherwise successful call. formatError only raises an error
// for such a body when it also carries a Message, so a code-only failure would
// otherwise be presented as an empty result. The 10002 family is resolved as
// an empty result before this check runs.
func metricListResponseError(response map[string]interface{}) error {
	if response == nil {
		return nil
	}
	code, ok := response["Code"]
	if !ok || code == nil {
		return nil
	}
	switch strings.ToLower(strings.TrimSpace(fmt.Sprint(code))) {
	case "", "0", "200", "success", "ok":
		return nil
	}
	return fmt.Errorf("%s failed, code: %s, message: %s", metricListAction, fmt.Sprint(code), metricListString(response["Message"]))
}

func metricListString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

// metricListNextToken returns the page token of the current response. The API
// omits NextToken (or returns null) on the last page; any non-string value is
// treated as "last page" so the state never holds a bogus token.
func metricListNextToken(response map[string]interface{}) string {
	if v, ok := response["NextToken"]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

// metricListTimestampString renders the millisecond epoch as a decimal string:
// json.Unmarshal decodes the JSON number into float64 and fmt.Sprint would
// emit scientific notation (e.g. "1.54877766e+12").
func metricListTimestampString(v interface{}) string {
	switch t := v.(type) {
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case json.Number:
		return t.String()
	case string:
		return t
	default:
		return ""
	}
}

func metricListFloatValue(v interface{}) attr.Value {
	if f, ok := metricListFloat64(v); ok {
		return types.Float64Value(f)
	}
	return types.Float64Null()
}

// metricListNumber mirrors metricListFloatValue for the output_file export,
// where a missing or non-numeric value stays null instead of becoming 0.
func metricListNumber(v interface{}) interface{} {
	if f, ok := metricListFloat64(v); ok {
		return f
	}
	return nil
}

func metricListFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case string:
		if n == "" {
			return 0, false
		}
		f, err := strconv.ParseFloat(n, 64)
		return f, err == nil
	default:
		return 0, false
	}
}
