package sas

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	scaDetailsTypeName   = "alicloud_security_center_sca_details"
	scaDetailsAction     = "DescribePropertyScaDetail"
	scaDetailsAPIVersion = "2018-12-03"
	scaDetailsReadWindow = 3 * time.Minute
)

// scaDetailsPropertyAttrTypes must stay flat: the provider is served over
// protocol version 5, which rejects nested attributes.
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

var scaDetailsCriteriaAttrTypes = map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}

var (
	_ datasource.DataSource                     = &scaDetailsDataSource{}
	_ datasource.DataSourceWithConfigure        = &scaDetailsDataSource{}
	_ datasource.DataSourceWithConfigValidators = &scaDetailsDataSource{}
)

func NewScaDetailsDataSource() datasource.DataSource {
	return &scaDetailsDataSource{}
}

type scaDetailsDataSource struct {
	fwadapt.DataSourceBase
}

type scaDetailsModel struct {
	Id                         types.String `tfsdk:"id"`
	Lang                       types.String `tfsdk:"lang"`
	Biz                        types.String `tfsdk:"biz"`
	SearchItem                 types.String `tfsdk:"search_item"`
	SearchInfo                 types.String `tfsdk:"search_info"`
	SearchCriteriaList         types.List   `tfsdk:"search_criteria_list"`
	ScaName                    types.String `tfsdk:"sca_name"`
	BizType                    types.String `tfsdk:"biz_type"`
	SearchItemSub              types.String `tfsdk:"search_item_sub"`
	SearchInfoSub              types.String `tfsdk:"search_info_sub"`
	Remark                     types.String `tfsdk:"remark"`
	Uuid                       types.String `tfsdk:"uuid"`
	ScaVersion                 types.String `tfsdk:"sca_version"`
	Port                       types.String `tfsdk:"port"`
	User                       types.String `tfsdk:"user"`
	Pid                        types.String `tfsdk:"pid"`
	ProcessStartedStart        types.Int64  `tfsdk:"process_started_start"`
	ProcessStartedEnd          types.Int64  `tfsdk:"process_started_end"`
	ScaNamePattern             types.String `tfsdk:"sca_name_pattern"`
	ResourceDirectoryAccountId types.Int64  `tfsdk:"resource_directory_account_id"`
	CurrentPage                types.Int64  `tfsdk:"current_page"`
	PageSize                   types.Int64  `tfsdk:"page_size"`
	UseNextToken               types.Bool   `tfsdk:"use_next_token"`
	NextToken                  types.String `tfsdk:"next_token"`
	TotalCount                 types.Int64  `tfsdk:"total_count"`
	PageItemCount              types.Int64  `tfsdk:"page_item_count"`
	Propertys                  types.List   `tfsdk:"propertys"`
	NextTokenResult            types.String `tfsdk:"next_token_result"`
}

func (d *scaDetailsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = scaDetailsTypeName
}

// ConfigValidators keeps the search condition and the sub-condition each
// complete: a lone search_item or search_item_sub makes the service ignore the
// condition and return a wider result set than the configuration asks for.
func (d *scaDetailsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.RequiredTogether(
			path.MatchRoot("search_item"),
			path.MatchRoot("search_info"),
		),
		datasourcevalidator.RequiredTogether(
			path.MatchRoot("search_item_sub"),
			path.MatchRoot("search_info_sub"),
		),
	}
}

func (d *scaDetailsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the asset fingerprint details of the middleware, databases and web services collected by Security Center.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, which is derived from the query condition.",
			},
			"lang": schema.StringAttribute{
				Optional:    true,
				Description: "The language of the request and the response. Defaults to `zh`. Valid values: `zh` and `en`.",
				Validators: []validator.String{
					stringvalidator.OneOf("zh", "en"),
				},
			},
			"biz": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the asset fingerprint. Valid values: `sca` (middleware), `sca_database` (database) and `sca_web` (web service).",
				Validators: []validator.String{
					stringvalidator.OneOf("sca", "sca_database", "sca_web"),
				},
			},
			"search_item": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the search condition. Valid values: `name` and `type`. It is used together with `search_info`.",
				Validators: []validator.String{
					stringvalidator.OneOf("name", "type"),
				},
			},
			"search_info": schema.StringAttribute{
				Optional:    true,
				Description: "The content of the search condition. It is used together with `search_item`.",
			},
			"search_criteria_list": schema.ListAttribute{
				Optional:    true,
				Description: "The search conditions. Each element contains `name` and `value`. A maximum of 20 conditions are supported.",
				ElementType: types.ObjectType{AttrTypes: scaDetailsCriteriaAttrTypes},
				Validators: []validator.List{
					listvalidator.SizeAtMost(20),
				},
			},
			"sca_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the asset fingerprint.",
			},
			"biz_type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the middleware, database or web service. Valid values: `system_service`, `software_library`, `docker_component`, `database`, `web_container`, `jar` and `web_framework`.",
				Validators: []validator.String{
					stringvalidator.OneOf("system_service", "software_library", "docker_component", "database", "web_container", "jar", "web_framework"),
				},
			},
			"search_item_sub": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the sub-condition. Valid values: `port`, `pid`, `version` and `user`. It is used together with `search_info_sub`.",
				Validators: []validator.String{
					stringvalidator.OneOf("port", "pid", "version", "user"),
				},
			},
			"search_info_sub": schema.StringAttribute{
				Optional:    true,
				Description: "The content of the sub-condition. It is used together with `search_item_sub`.",
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Description: "The name or IP address of the server. Fuzzy match is supported.",
			},
			"uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the server on which the asset fingerprint is collected.",
			},
			"sca_version": schema.StringAttribute{
				Optional:    true,
				Description: "The version of the middleware, database or web service.",
			},
			"port": schema.StringAttribute{
				Optional:    true,
				Description: "The listening port of the process.",
			},
			"user": schema.StringAttribute{
				Optional:    true,
				Description: "The user that runs the process.",
			},
			"pid": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the process.",
			},
			"process_started_start": schema.Int64Attribute{
				Optional:    true,
				Description: "The start of the process start time range, in seconds.",
			},
			"process_started_end": schema.Int64Attribute{
				Optional:    true,
				Description: "The end of the process start time range, in seconds.",
			},
			"sca_name_pattern": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the process.",
			},
			"resource_directory_account_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The ID of the Alibaba Cloud account that is a member of the resource directory.",
			},
			"current_page": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of the page to return. Defaults to `1`.",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of entries per page. Defaults to `10`.",
			},
			"use_next_token": schema.BoolAttribute{
				Optional:    true,
				Description: "Specifies whether to use `next_token` instead of the paged query.",
			},
			"next_token": schema.StringAttribute{
				Optional:    true,
				Description: "The token that starts the query. It is used together with `use_next_token`.",
			},
			"next_token_result": schema.StringAttribute{
				Computed:    true,
				Description: "The token that is used to retrieve the next page of results.",
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of the asset fingerprints that match the query condition. It is left unset when the API does not return it.",
			},
			"page_item_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of the asset fingerprints on the current page.",
			},
			"propertys": schema.ListAttribute{
				Computed:    true,
				Description: "A list of the asset fingerprint details.",
				ElementType: types.ObjectType{AttrTypes: scaDetailsPropertyAttrTypes},
			},
		},
	}
}

func (d *scaDetailsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data scaDetailsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{}
	setString := func(key string, value types.String) {
		if !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
			request[key] = value.ValueString()
		}
	}
	setInt64 := func(key string, value types.Int64) {
		if !value.IsNull() && !value.IsUnknown() {
			request[key] = value.ValueInt64()
		}
	}

	setString("Lang", data.Lang)
	setString("Biz", data.Biz)
	setString("SearchItem", data.SearchItem)
	setString("SearchInfo", data.SearchInfo)
	setString("ScaName", data.ScaName)
	setString("BizType", data.BizType)
	setString("SearchItemSub", data.SearchItemSub)
	setString("SearchInfoSub", data.SearchInfoSub)
	setString("Remark", data.Remark)
	setString("Uuid", data.Uuid)
	setString("ScaVersion", data.ScaVersion)
	setString("Port", data.Port)
	setString("User", data.User)
	setString("Pid", data.Pid)
	setInt64("ProcessStartedStart", data.ProcessStartedStart)
	setInt64("ProcessStartedEnd", data.ProcessStartedEnd)
	setString("ScaNamePattern", data.ScaNamePattern)
	setInt64("ResourceDirectoryAccountId", data.ResourceDirectoryAccountId)
	setInt64("CurrentPage", data.CurrentPage)
	setInt64("PageSize", data.PageSize)
	setString("NextToken", data.NextToken)
	if !data.UseNextToken.IsNull() && !data.UseNextToken.IsUnknown() {
		request["UseNextToken"] = data.UseNextToken.ValueBool()
	}

	if !data.SearchCriteriaList.IsNull() && !data.SearchCriteriaList.IsUnknown() {
		elements := data.SearchCriteriaList.Elements()
		criterias := make([]map[string]interface{}, 0, len(elements))
		for _, elementValue := range elements {
			element, ok := elementValue.(types.Object)
			if !ok {
				continue
			}
			attributes := element.Attributes()
			name, _ := attributes["name"].(types.String)
			value, _ := attributes["value"].(types.String)
			item := map[string]interface{}{}
			if !name.IsNull() && !name.IsUnknown() && name.ValueString() != "" {
				item["Name"] = name.ValueString()
			}
			if !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
				item["Value"] = value.ValueString()
			}
			criterias = append(criterias, item)
		}
		if len(criterias) > 0 {
			request["SearchCriteriaList"] = criterias
		}
	}

	response, err := helper.RetryRpcPost(ctx, scaDetailsAction, func() (map[string]interface{}, error) {
		return d.Client().RpcPost("Sas", scaDetailsAPIVersion, scaDetailsAction, nil, request, true)
	}, scaDetailsReadWindow, helper.NeedRetry)
	logs.AddDebug(scaDetailsAction, response, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling %s", scaDetailsTypeName, scaDetailsAction),
			err.Error(),
		)
		return
	}

	rawPropertys, _ := response["Propertys"].([]interface{})
	propertys := make([]attr.Value, 0, len(rawPropertys))
	for _, raw := range rawPropertys {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		propertys = append(propertys, types.ObjectValueMust(scaDetailsPropertyAttrTypes, map[string]attr.Value{
			"type":                scaToStringValue(item["Type"]),
			"biz_type_dispaly":    scaToStringValue(item["BizTypeDispaly"]),
			"process_started":     scaToInt64Value(item["ProcessStarted"]),
			"internet_ip":         scaToStringValue(item["InternetIp"]),
			"config_path":         scaToStringValue(item["ConfigPath"]),
			"pid":                 scaToStringValue(item["Pid"]),
			"port":                scaToStringValue(item["Port"]),
			"cmdline":             scaToStringValue(item["Cmdline"]),
			"biz_type":            scaToStringValue(item["BizType"]),
			"listen_ip":           scaToStringValue(item["ListenIp"]),
			"version":             scaToStringValue(item["Version"]),
			"instance_name":       scaToStringValue(item["InstanceName"]),
			"listen_status":       scaToStringValue(item["ListenStatus"]),
			"name":                scaToStringValue(item["Name"]),
			"ip":                  scaToStringValue(item["Ip"]),
			"process_user":        scaToStringValue(item["ProcessUser"]),
			"instance_id":         scaToStringValue(item["InstanceId"]),
			"web_path":            scaToStringValue(item["WebPath"]),
			"ppid":                scaToStringValue(item["Ppid"]),
			"intranet_ip":         scaToStringValue(item["IntranetIp"]),
			"uuid":                scaToStringValue(item["Uuid"]),
			"listen_protocol":     scaToStringValue(item["ListenProtocol"]),
			"image_name":          scaToStringValue(item["ImageName"]),
			"path":                scaToStringValue(item["Path"]),
			"container_name":      scaToStringValue(item["ContainerName"]),
			"proof":               scaToStringValue(item["Proof"]),
			"create_timestamp":    scaToInt64Value(item["CreateTimestamp"]),
			"runtime_env_version": scaToStringValue(item["RuntimeEnvVersion"]),
			"pod_name":            scaToStringValue(item["PodName"]),
		}))
	}
	propertyList, diags := types.ListValue(types.ObjectType{AttrTypes: scaDetailsPropertyAttrTypes}, propertys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pageInfo, _ := response["PageInfo"].(map[string]interface{})

	state := scaDetailsModel{
		Id:                         types.StringValue(helper.DataResourceIdHashFromMap(request)),
		Lang:                       data.Lang,
		Biz:                        data.Biz,
		SearchItem:                 data.SearchItem,
		SearchInfo:                 data.SearchInfo,
		SearchCriteriaList:         data.SearchCriteriaList,
		ScaName:                    data.ScaName,
		BizType:                    data.BizType,
		SearchItemSub:              data.SearchItemSub,
		SearchInfoSub:              data.SearchInfoSub,
		Remark:                     data.Remark,
		Uuid:                       data.Uuid,
		ScaVersion:                 data.ScaVersion,
		Port:                       data.Port,
		User:                       data.User,
		Pid:                        data.Pid,
		ProcessStartedStart:        data.ProcessStartedStart,
		ProcessStartedEnd:          data.ProcessStartedEnd,
		ScaNamePattern:             data.ScaNamePattern,
		ResourceDirectoryAccountId: data.ResourceDirectoryAccountId,
		CurrentPage:                scaEchoInt64(data.CurrentPage, pageInfo["CurrentPage"]),
		PageSize:                   scaEchoInt64(data.PageSize, pageInfo["PageSize"]),
		UseNextToken:               data.UseNextToken,
		NextToken:                  data.NextToken,
		NextTokenResult:            scaNextTokenResult(pageInfo["NextToken"]),
		TotalCount:                 scaToInt64Value(pageInfo["TotalCount"]),
		PageItemCount:              scaToInt64Value(pageInfo["Count"]),
		Propertys:                  propertyList,
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
