package sas

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const scaDetailsTypeName = "alicloud_threat_detection_sca_details"

var (
	_ datasource.DataSource              = &scaDetailsDataSource{}
	_ datasource.DataSourceWithConfigure = &scaDetailsDataSource{}
)

func NewScaDetailsDataSource() datasource.DataSource {
	return &scaDetailsDataSource{}
}

type scaDetailsDataSource struct {
	fwadapt.DataSourceBase
}

type scaDetailsModel struct {
	SearchItem    types.String `tfsdk:"search_item"`
	SearchInfo    types.String `tfsdk:"search_info"`
	SearchItemSub types.String `tfsdk:"search_item_sub"`
	SearchInfoSub types.String `tfsdk:"search_info_sub"`
	Biz           types.String `tfsdk:"biz"`
	Lang          types.String `tfsdk:"lang"`
	OutputFile    types.String `tfsdk:"output_file"`
	Id            types.String `tfsdk:"id"`
	Propertys     types.List   `tfsdk:"details"`
}

var scaDetailsPropertyAttrTypes = map[string]attr.Type{
	"type":                types.StringType,
	"biz_type_dispaly":    types.StringType,
	"process_started":     types.Int64Type,
	"internet_ip":         types.StringType,
	"config_path":         types.StringType,
	"pid":                 types.StringType,
	"port":                types.StringType,
	"cmdline":             types.StringType,
	"biz_type":            types.StringType,
	"listen_ip":           types.StringType,
	"version":             types.StringType,
	"instance_name":       types.StringType,
	"listen_status":       types.StringType,
	"name":                types.StringType,
	"ip":                  types.StringType,
	"process_user":        types.StringType,
	"instance_id":         types.StringType,
	"web_path":            types.StringType,
	"ppid":                types.StringType,
	"intranet_ip":         types.StringType,
	"uuid":                types.StringType,
	"listen_protocol":     types.StringType,
	"image_name":          types.StringType,
	"path":                types.StringType,
	"container_name":      types.StringType,
	"proof":               types.StringType,
	"create_timestamp":    types.Int64Type,
	"runtime_env_version": types.StringType,
	"pod_name":            types.StringType,
}

func (d *scaDetailsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = scaDetailsTypeName
}

func (d *scaDetailsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the asset fingerprint details of the middleware, databases and web services collected by Security Center.",
		Attributes: map[string]schema.Attribute{
			"search_item": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the search condition. Valid values: `name` and `type`. It is used together with `search_info`.",
			},
			"search_info": schema.StringAttribute{
				Optional:    true,
				Description: "The content of the search condition. It is used together with `search_item`.",
			},
			"search_item_sub": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the sub-condition. Valid values: `port`, `pid`, `version` and `user`. It is used together with `search_info_sub`.",
			},
			"search_info_sub": schema.StringAttribute{
				Optional:    true,
				Description: "The content of the sub-condition. It is used together with `search_item_sub`.",
			},
			"biz": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the asset fingerprint. Valid values: `sca` (middleware), `sca_database` (database) and `sca_web` (web service).",
			},
			"lang": schema.StringAttribute{
				Optional:    true,
				Description: "The language of the request and the response. Defaults to `zh`. Valid values: `zh` and `en`.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save the data source results.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, which is derived from the query condition.",
			},
			"details": schema.ListAttribute{
				Computed:    true,
				Description: "The details of the Asset Fingerprints list.",
				ElementType: types.ObjectType{AttrTypes: scaDetailsPropertyAttrTypes},
			},
		},
	}
}

func (d *scaDetailsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	const action = "DescribePropertyScaDetail"

	var config scaDetailsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{}
	request["PageSize"] = helper.PageSizeLarge
	request["CurrentPage"] = 1

	setString := func(key string, value types.String) {
		if !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
			request[key] = value.ValueString()
		}
	}

	setString("SearchItem", config.SearchItem)
	setString("SearchInfo", config.SearchInfo)
	setString("SearchItemSub", config.SearchItemSub)
	setString("SearchInfoSub", config.SearchInfoSub)
	setString("Biz", config.Biz)
	setString("Lang", config.Lang)

	var objects []map[string]interface{}
	var callResponse map[string]interface{}
	var callErr error

	for {
		callResponse, callErr = helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
			return d.Meta().RpcPost("Sas", "2018-12-03", action, nil, request, true)
		}, 5*time.Minute, helper.NeedRetry)
		logs.AddDebug(action, callResponse, request)

		if callErr != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling %s", scaDetailsTypeName, action),
				callErr.Error(),
			)
			return
		}

		rawPropertys, _ := callResponse["Propertys"].([]interface{})
		for _, v := range rawPropertys {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			objects = append(objects, item)
		}

		if len(rawPropertys) < helper.PageSizeLarge {
			break
		}

		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	properties := make([]map[string]interface{}, 0, len(objects))
	items := make([]attr.Value, 0, len(objects))
	for _, raw := range objects {
		property := map[string]interface{}{
			"type":                raw["Type"],
			"biz_type_dispaly":    raw["BizTypeDispaly"],
			"process_started":     raw["ProcessStarted"],
			"internet_ip":         raw["InternetIp"],
			"config_path":         raw["ConfigPath"],
			"pid":                 raw["Pid"],
			"port":                raw["Port"],
			"cmdline":             raw["Cmdline"],
			"biz_type":            raw["BizType"],
			"listen_ip":           raw["ListenIp"],
			"version":             raw["Version"],
			"instance_name":       raw["InstanceName"],
			"listen_status":       raw["ListenStatus"],
			"name":                raw["Name"],
			"ip":                  raw["Ip"],
			"process_user":        raw["ProcessUser"],
			"instance_id":         raw["InstanceId"],
			"web_path":            raw["WebPath"],
			"ppid":                raw["Ppid"],
			"intranet_ip":         raw["IntranetIp"],
			"uuid":                raw["Uuid"],
			"listen_protocol":     raw["ListenProtocol"],
			"image_name":          raw["ImageName"],
			"path":                raw["Path"],
			"container_name":      raw["ContainerName"],
			"proof":               raw["Proof"],
			"create_timestamp":    raw["CreateTimestamp"],
			"runtime_env_version": raw["RuntimeEnvVersion"],
			"pod_name":            raw["PodName"],
		}

		properties = append(properties, property)

		items = append(items, types.ObjectValueMust(scaDetailsPropertyAttrTypes, map[string]attr.Value{
			"type":                scaToStringValue(raw["Type"]),
			"biz_type_dispaly":    scaToStringValue(raw["BizTypeDispaly"]),
			"process_started":     scaToInt64Value(raw["ProcessStarted"]),
			"internet_ip":         scaToStringValue(raw["InternetIp"]),
			"config_path":         scaToStringValue(raw["ConfigPath"]),
			"pid":                 scaToStringValue(raw["Pid"]),
			"port":                scaToStringValue(raw["Port"]),
			"cmdline":             scaToStringValue(raw["Cmdline"]),
			"biz_type":            scaToStringValue(raw["BizType"]),
			"listen_ip":           scaToStringValue(raw["ListenIp"]),
			"version":             scaToStringValue(raw["Version"]),
			"instance_name":       scaToStringValue(raw["InstanceName"]),
			"listen_status":       scaToStringValue(raw["ListenStatus"]),
			"name":                scaToStringValue(raw["Name"]),
			"ip":                  scaToStringValue(raw["Ip"]),
			"process_user":        scaToStringValue(raw["ProcessUser"]),
			"instance_id":         scaToStringValue(raw["InstanceId"]),
			"web_path":            scaToStringValue(raw["WebPath"]),
			"ppid":                scaToStringValue(raw["Ppid"]),
			"intranet_ip":         scaToStringValue(raw["IntranetIp"]),
			"uuid":                scaToStringValue(raw["Uuid"]),
			"listen_protocol":     scaToStringValue(raw["ListenProtocol"]),
			"image_name":          scaToStringValue(raw["ImageName"]),
			"path":                scaToStringValue(raw["Path"]),
			"container_name":      scaToStringValue(raw["ContainerName"]),
			"proof":               scaToStringValue(raw["Proof"]),
			"create_timestamp":    scaToInt64Value(raw["CreateTimestamp"]),
			"runtime_env_version": scaToStringValue(raw["RuntimeEnvVersion"]),
			"pod_name":            scaToStringValue(raw["PodName"]),
		}))
	}

	propertys, diags := types.ListValue(types.ObjectType{AttrTypes: scaDetailsPropertyAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := config
	state.Id = types.StringValue(helper.DataResourceIdHash([]string{fmt.Sprint(objects)}))
	state.Propertys = propertys

	if err := scaDetailsWriteOutputFile(state.OutputFile, properties); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: writing output_file", scaDetailsTypeName),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// scaEchoInt64 keeps a configured value because the query cannot change it,
// and falls back to the value the API reports.
func scaEchoInt64(configured types.Int64, reported interface{}) types.Int64 {
	if !configured.IsNull() && !configured.IsUnknown() {
		return configured
	}
	return scaToInt64Value(reported)
}

// scaNextTokenResult exposes the token of the next page without overwriting the
// token the practitioner configured for this query.
func scaNextTokenResult(reported interface{}) types.String {
	if reportedString, ok := reported.(string); ok && reportedString != "" {
		return types.StringValue(reportedString)
	}
	return types.StringNull()
}

func scaToInt64Value(value interface{}) types.Int64 {
	switch v := value.(type) {
	case nil:
		return types.Int64Null()
	case float64:
		return types.Int64Value(int64(v))
	case float32:
		return types.Int64Value(int64(v))
	case int:
		return types.Int64Value(int64(v))
	case int32:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case string:
		if v == "" {
			return types.Int64Null()
		}
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return types.Int64Null()
		}
		return types.Int64Value(parsed)
	default:
		parsed, err := strconv.ParseInt(fmt.Sprint(v), 10, 64)
		if err != nil {
			return types.Int64Null()
		}
		return types.Int64Value(parsed)
	}
}

func scaToStringValue(value interface{}) types.String {
	switch v := value.(type) {
	case nil:
		return types.StringNull()
	case string:
		return types.StringValue(v)
	default:
		return types.StringValue(fmt.Sprint(v))
	}
}

func scaDetailsWriteOutputFile(output types.String, objects []map[string]interface{}) error {
	if output.IsNull() || output.IsUnknown() || output.ValueString() == "" {
		return nil
	}

	if err := helper.WriteToFile(output.ValueString(), objects); err != nil {
		return err
	}

	return nil
}
