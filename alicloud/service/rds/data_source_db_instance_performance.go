package rds

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
				Description: "The performance metrics that you want to query. Separate multiple metrics with commas (,). You can query up to 30 metrics at a time. The returned `performance_keys` list is ordered by this comma order.",
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
				Description: "A list of performance metrics for the requested `key` list, rebuilt in the comma order of the request `key`. Every metric echoed by the API is kept, even when the specified time range holds no monitoring data: such a metric has an empty `values` list, and only a response without any performance key yields an empty `performance_keys` list.",
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

	response, err := d.describeDBInstancePerformance(&model)
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

	// Keep every PerformanceKey the API returns and rebuild the list in the
	// request key comma order, so indexes stay stable across refreshes no
	// matter the order the API returns the keys in or whether monitoring
	// samples exist.
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

func (d *dbInstancePerformanceDataSource) describeDBInstancePerformance(model *dbInstancePerformanceModel) (map[string]interface{}, error) {
	client := d.Client()
	action := "DescribeDBInstancePerformance"
	request := map[string]interface{}{
		"SourceIp":     client.SourceIp,
		"DBInstanceId": model.DbInstanceId.ValueString(),
		"Key":          model.Key.ValueString(),
		"StartTime":    model.StartTime.ValueString(),
		"EndTime":      model.EndTime.ValueString(),
	}
	if v := model.NodeId.ValueString(); !model.NodeId.IsNull() && v != "" {
		request["NodeId"] = v
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = sdkretry.Retry(5*time.Minute, func() *sdkretry.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return sdkretry.RetryableError(err)
			}
			return sdkretry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, dbInstancePerformanceTypeName, action, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
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

// rebuildPerformanceKeys orders the parsed performance keys by the request key
// comma order. Matching is case-sensitive on the exact comma-separated tokens:
// a requested token the API did not echo is not synthesized, and response keys
// outside the request are appended afterwards in response order instead of
// being dropped. The result therefore does not depend on the order in which
// the API returned the keys, nor on whether monitoring samples exist.
func rebuildPerformanceKeys(requestKey string, items []performanceKeyItem) []performanceKeyItem {
	byKey := make(map[string]performanceKeyItem, len(items))
	for _, item := range items {
		byKey[item.Key] = item
	}
	inRequest := make(map[string]bool)
	ordered := make([]performanceKeyItem, 0, len(items))
	for _, token := range strings.Split(requestKey, ",") {
		if inRequest[token] {
			continue
		}
		inRequest[token] = true
		if item, ok := byKey[token]; ok {
			ordered = append(ordered, item)
		}
	}
	for _, item := range items {
		if !inRequest[item.Key] {
			ordered = append(ordered, item)
		}
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
