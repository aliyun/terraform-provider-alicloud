package rds

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const dbInstancePerformanceTypeName = "alicloud_rds_db_instance_performance"

var (
	_ datasource.DataSource              = &dbInstancePerformanceDataSource{}
	_ datasource.DataSourceWithConfigure = &dbInstancePerformanceDataSource{}
)

func NewDbInstancePerformanceDataSource() datasource.DataSource {
	return &dbInstancePerformanceDataSource{}
}

type dbInstancePerformanceDataSource struct {
	fwadapt.DataSourceBase
}

// dbInstancePerformanceMaxKeys mirrors the API limit: Key accepts at most 30
// performance metrics per request.
const dbInstancePerformanceMaxKeys = 30

const dbInstancePerformanceQueryTimeoutCode = "InvalidOthers.Timeout"

// dbInstancePerformanceKeyValidator rejects a `key` that the API would reject,
// so the mismatch surfaces at plan time and no request is sent.
type dbInstancePerformanceKeyValidator struct{}

func (v dbInstancePerformanceKeyValidator) Description(context.Context) string {
	return "The value must contain between 1 and 30 distinct comma-separated performance metrics."
}

func (v dbInstancePerformanceKeyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v dbInstancePerformanceKeyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if err := validatePerformanceKey(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Performance Key", err.Error())
	}
}

// validatePerformanceKey applies the documented `key` rules: split on commas,
// trim each entry, ignore empty entries, and require 1 to 30 distinct metrics.
func validatePerformanceKey(key string) error {
	tokens := performanceKeyTokens(key)
	if len(tokens) == 0 {
		return fmt.Errorf("`key` must contain at least one performance metric, separated by commas")
	}
	if len(tokens) > dbInstancePerformanceMaxKeys {
		return fmt.Errorf("`key` accepts at most %d performance metrics, but %d distinct metrics were provided", dbInstancePerformanceMaxKeys, len(tokens))
	}
	return nil
}

// performanceKeyTokens returns the distinct, trimmed metrics of `key` in the
// request comma order. It is used for the request itself and for rebuilding the
// output list, so both share exactly the same token semantics.
func performanceKeyTokens(key string) []string {
	tokens := make([]string, 0, dbInstancePerformanceMaxKeys)
	seen := make(map[string]struct{}, dbInstancePerformanceMaxKeys)
	for _, raw := range strings.Split(key, ",") {
		token := strings.TrimSpace(raw)
		if token == "" {
			continue
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		tokens = append(tokens, token)
	}
	return tokens
}

type dbInstancePerformanceValueModel struct {
	Date  types.String `tfsdk:"date"`
	Value types.String `tfsdk:"value"`
}

type dbInstancePerformanceKeyModel struct {
	Key         types.String                      `tfsdk:"key"`
	Unit        types.String                      `tfsdk:"unit"`
	ValueFormat types.String                      `tfsdk:"value_format"`
	Values      []dbInstancePerformanceValueModel `tfsdk:"values"`
}

type dbInstancePerformanceModel struct {
	Id              types.String                    `tfsdk:"id"`
	DbInstanceId    types.String                    `tfsdk:"db_instance_id"`
	Key             types.String                    `tfsdk:"key"`
	StartTime       types.String                    `tfsdk:"start_time"`
	EndTime         types.String                    `tfsdk:"end_time"`
	NodeId          types.String                    `tfsdk:"node_id"`
	OutputFile      types.String                    `tfsdk:"output_file"`
	Engine          types.String                    `tfsdk:"engine"`
	PerformanceKeys []dbInstancePerformanceKeyModel `tfsdk:"performance_keys"`
}

func (d *dbInstancePerformanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = dbInstancePerformanceTypeName
}

func (d *dbInstancePerformanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of performance metrics of an ApsaraDB RDS instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source.",
			},
			"db_instance_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the RDS instance.",
			},
			"key": schema.StringAttribute{
				Required:    true,
				Description: "The performance metrics that you want to query. Separate multiple metrics with commas (,). Whitespace around each metric is trimmed and empty entries are ignored. You must provide at least one and at most 30 distinct metrics; otherwise the plan fails before any API call. The returned `performance_keys` list is ordered by this comma order.",
				Validators: []validator.String{
					dbInstancePerformanceKeyValidator{},
				},
			},
			"start_time": schema.StringAttribute{
				Required:    true,
				Description: "The beginning of the time range to query. Specify the time in the `yyyy-MM-ddTHH:mmZ` format. The time must be in UTC.",
			},
			"end_time": schema.StringAttribute{
				Required:    true,
				Description: "The end of the time range to query. Specify the time in the `yyyy-MM-ddTHH:mmZ` format. The time must be in UTC.",
			},
			"node_id": schema.StringAttribute{
				Optional:    true,
				Description: "The unique ID of the instance node.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"engine": schema.StringAttribute{
				Computed:    true,
				Description: "The database engine of the instance.",
			},
			"performance_keys": schema.ListAttribute{
				Computed:    true,
				Description: "A list of performance metrics rebuilt from the request `key`: each distinct request token appears in the request comma order, and a requested metric the API did not echo is kept in place with an empty `unit`, an empty `value_format` and an empty `values` list. Duplicate request tokens are collapsed and empty tokens are skipped. Metrics returned outside the request are appended after the requested ones, sorted by `key` in ascending lexicographical (case-sensitive) order. Because `key` must contain at least one metric, the list is never empty for a valid request.",
				// A plain list attribute with an object element type keeps the
				// same wire shape as the SDK v2 form (list of resources). The
				// nested attribute variants are not usable here: this provider
				// serves every half behind one protocol 5 server, which cannot
				// express nested attributes.
				ElementType: dbInstancePerformanceEntryObjectType,
			},
		},
	}
}

// dbInstancePerformanceEntryObjectType and
// dbInstancePerformanceValueEntryObjectType describe the element shapes of the
// performance_keys list attribute.
var dbInstancePerformanceEntryObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"key":          types.StringType,
	"unit":         types.StringType,
	"value_format": types.StringType,
	"values":       types.ListType{ElemType: dbInstancePerformanceValueEntryObjectType},
}}

var dbInstancePerformanceValueEntryObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"date":  types.StringType,
	"value": types.StringType,
}}

func (d *dbInstancePerformanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model dbInstancePerformanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := validatePerformanceKey(model.Key.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: invalid key", dbInstancePerformanceTypeName),
			err.Error(),
		)
		return
	}

	response, err := d.describeDBInstancePerformance(ctx, &model)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeDBInstancePerformance", dbInstancePerformanceTypeName),
			err.Error(),
		)
		return
	}

	engine := types.StringNull()
	if v, ok := response["Engine"].(string); ok {
		engine = types.StringValue(v)
	}

	// Rebuild the list from the request tokens, so every requested metric
	// keeps its index across refreshes even when the API stops echoing it,
	// and response keys outside the request keep a sorted relative order.
	ordered := rebuildPerformanceKeys(model.Key.ValueString(), performanceKeyItems(response))
	model.Id = types.StringValue(dbInstancePerformanceId(&model))
	model.Engine = engine
	model.PerformanceKeys = make([]dbInstancePerformanceKeyModel, 0, len(ordered))
	for _, item := range ordered {
		keyModel := dbInstancePerformanceKeyModel{
			Key:         types.StringValue(item.Key),
			Unit:        types.StringValue(item.Unit),
			ValueFormat: types.StringValue(item.ValueFormat),
			Values:      []dbInstancePerformanceValueModel{},
		}
		for _, value := range item.Values {
			keyModel.Values = append(keyModel.Values, dbInstancePerformanceValueModel{
				Date:  types.StringValue(value.Date),
				Value: types.StringValue(value.Value),
			})
		}
		model.PerformanceKeys = append(model.PerformanceKeys, keyModel)
	}

	if output := model.OutputFile.ValueString(); output != "" {
		if err := helper.WriteToFile(output, ordered); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Writing %s: output file", dbInstancePerformanceTypeName),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (d *dbInstancePerformanceDataSource) describeDBInstancePerformance(ctx context.Context, model *dbInstancePerformanceModel) (map[string]interface{}, error) {
	client := d.Client()
	action := "DescribeDBInstancePerformance"
	request := map[string]interface{}{
		"SourceIp":     client.SourceIp,
		"DBInstanceId": model.DbInstanceId.ValueString(),
		"Key":          strings.Join(performanceKeyTokens(model.Key.ValueString()), ","),
		"StartTime":    model.StartTime.ValueString(),
		"EndTime":      model.EndTime.ValueString(),
	}
	if v := model.NodeId.ValueString(); !model.NodeId.IsNull() && v != "" {
		request["NodeId"] = v
	}

	var response map[string]interface{}
	// Cancellation can return before the RPC callback finishes.
	var responseMu sync.Mutex
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := sdkretry.RetryContext(ctx, 5*time.Minute, func() *sdkretry.RetryError {
		result, rpcErr := client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
		responseMu.Lock()
		response = result
		responseMu.Unlock()
		if rpcErr != nil {
			if helper.NeedRetry(rpcErr) || isDbInstancePerformanceQueryTimeout(rpcErr) {
				wait()
				return sdkretry.RetryableError(rpcErr)
			}
			return sdkretry.NonRetryableError(rpcErr)
		}
		return nil
	})
	responseMu.Lock()
	defer responseMu.Unlock()
	logs.AddDebug(action, response, request)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, dbInstancePerformanceTypeName, action, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// Prefix matching would also retry an empty SDK error code.
func isDbInstancePerformanceQueryTimeout(err error) bool {
	e, ok := err.(*tea.SDKError)
	return ok && tea.StringValue(e.Code) == dbInstancePerformanceQueryTimeoutCode
}

// performanceKeyItem is one PerformanceKey entry parsed from the API response;
// the JSON tags also shape the output_file content.
type performanceKeyItem struct {
	Key         string                 `json:"key"`
	Unit        string                 `json:"unit"`
	ValueFormat string                 `json:"value_format"`
	Values      []performanceValueItem `json:"values"`
}

type performanceValueItem struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

// performanceKeyItems parses the PerformanceKeys.PerformanceKey array and keeps
// every echoed entry, even when the entry carries no Values object or an empty
// value list. Only a response without any PerformanceKey yields no items.
func performanceKeyItems(response map[string]interface{}) []performanceKeyItem {
	items := make([]performanceKeyItem, 0)
	raw, err := jsonpath.Get("$.PerformanceKeys.PerformanceKey", response)
	if err != nil {
		return items
	}
	list, ok := raw.([]interface{})
	if !ok {
		return items
	}
	for _, v := range list {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		entry := performanceKeyItem{
			Key:         fmt.Sprint(item["Key"]),
			Unit:        fmt.Sprint(item["Unit"]),
			ValueFormat: fmt.Sprint(item["ValueFormat"]),
			Values:      []performanceValueItem{},
		}
		if valueObj, ok := item["Values"].(map[string]interface{}); ok {
			if valueItems, ok := valueObj["PerformanceValue"].([]interface{}); ok {
				for _, valueItem := range valueItems {
					if valueMap, ok := valueItem.(map[string]interface{}); ok {
						entry.Values = append(entry.Values, performanceValueItem{
							Date:  fmt.Sprint(valueMap["Date"]),
							Value: fmt.Sprint(valueMap["Value"]),
						})
					}
				}
			}
		}
		items = append(items, entry)
	}
	return items
}

// rebuildPerformanceKeys rebuilds the performance key list from the request key
// tokens and the parsed response entries. Every distinct request token appears
// in the request comma order; a requested token the API did not echo is
// synthesized as a placeholder entry with empty unit, value_format and values,
// so its index never shifts when the API starts or stops echoing it. Matching
// is case-sensitive on the exact comma-separated tokens, duplicate request
// tokens collapse into one entry, and empty tokens are skipped. Response keys
// outside the request are appended after the requested ones, sorted by key in
// ascending lexicographical (case-sensitive) order, so their relative order
// does not depend on the order in which the API returned them.
func rebuildPerformanceKeys(requestKey string, items []performanceKeyItem) []performanceKeyItem {
	byKey := make(map[string]performanceKeyItem, len(items))
	for _, item := range items {
		byKey[item.Key] = item
	}

	inRequest := make(map[string]bool)
	ordered := make([]performanceKeyItem, 0, len(items))
	for _, token := range performanceKeyTokens(requestKey) {
		inRequest[token] = true
		if item, ok := byKey[token]; ok {
			ordered = append(ordered, item)
			continue
		}
		ordered = append(ordered, performanceKeyItem{
			Key:    token,
			Values: []performanceValueItem{},
		})
	}

	extra := make([]string, 0, len(byKey))
	for key := range byKey {
		if !inRequest[key] {
			extra = append(extra, key)
		}
	}
	sort.Strings(extra)
	for _, key := range extra {
		ordered = append(ordered, byKey[key])
	}
	return ordered
}

func dbInstancePerformanceId(model *dbInstancePerformanceModel) string {
	names := []string{
		model.DbInstanceId.ValueString(),
		model.Key.ValueString(),
		model.StartTime.ValueString(),
		model.EndTime.ValueString(),
	}
	if v := model.NodeId.ValueString(); !model.NodeId.IsNull() && v != "" {
		names = append(names, v)
	}
	var buf bytes.Buffer
	for _, id := range names {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}
	return fmt.Sprintf("%d", helper.Hashcode(buf.String()))
}
