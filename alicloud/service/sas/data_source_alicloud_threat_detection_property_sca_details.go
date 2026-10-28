package sas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const propertyScaDetailsTypeName = "alicloud_threat_detection_property_sca_details"

const (
	propertyScaDetailsAction           = "DescribePropertyScaDetail"
	propertyScaDetailsDefaultPageSize  = 10
	propertyScaDetailsMaxCriteriaItems = 20
	propertyScaDetailsQueryTimeout     = 5 * time.Minute
)

var propertyScaDetailsBizTypes = []string{"system_service", "software_library", "docker_component", "database", "web_container", "jar", "web_framework"}

var (
	_ datasource.DataSource                   = &propertyScaDetailsDataSource{}
	_ datasource.DataSourceWithConfigure      = &propertyScaDetailsDataSource{}
	_ datasource.DataSourceWithValidateConfig = &propertyScaDetailsDataSource{}
)

func NewThreatDetectionPropertyScaDetailsDataSource() datasource.DataSource {
	return &propertyScaDetailsDataSource{}
}

type propertyScaDetailsDataSource struct {
	fwadapt.DataSourceBase
}

type propertyScaDetailsCriteriaModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type propertyScaDetailsPropertyModel struct {
	Type              types.String `tfsdk:"type"`
	BizTypeDisplay    types.String `tfsdk:"biz_type_display"`
	ProcessStarted    types.Int64  `tfsdk:"process_started"`
	InternetIp        types.String `tfsdk:"internet_ip"`
	ConfigPath        types.String `tfsdk:"config_path"`
	Pid               types.String `tfsdk:"pid"`
	Port              types.String `tfsdk:"port"`
	Cmdline           types.String `tfsdk:"cmdline"`
	BizType           types.String `tfsdk:"biz_type"`
	ListenIp          types.String `tfsdk:"listen_ip"`
	Version           types.String `tfsdk:"version"`
	InstanceName      types.String `tfsdk:"instance_name"`
	ListenStatus      types.String `tfsdk:"listen_status"`
	Name              types.String `tfsdk:"name"`
	Ip                types.String `tfsdk:"ip"`
	ProcessUser       types.String `tfsdk:"process_user"`
	InstanceId        types.String `tfsdk:"instance_id"`
	WebPath           types.String `tfsdk:"web_path"`
	Ppid              types.String `tfsdk:"ppid"`
	IntranetIp        types.String `tfsdk:"intranet_ip"`
	Uuid              types.String `tfsdk:"uuid"`
	ListenProtocol    types.String `tfsdk:"listen_protocol"`
	ImageName         types.String `tfsdk:"image_name"`
	Path              types.String `tfsdk:"path"`
	ContainerName     types.String `tfsdk:"container_name"`
	Proof             types.String `tfsdk:"proof"`
	CreateTimestamp   types.Int64  `tfsdk:"create_timestamp"`
	RuntimeEnvVersion types.String `tfsdk:"runtime_env_version"`
	PodName           types.String `tfsdk:"pod_name"`
}

type propertyScaDetailsPageInfoValue struct {
	CurrentPage types.Int64  `tfsdk:"current_page"`
	PageSize    types.Int64  `tfsdk:"page_size"`
	TotalCount  types.Int64  `tfsdk:"total_count"`
	Count       types.Int64  `tfsdk:"count"`
	NextToken   types.String `tfsdk:"next_token"`
}

type propertyScaDetailsModel struct {
	Id                         types.String                      `tfsdk:"id"`
	Lang                       types.String                      `tfsdk:"lang"`
	Biz                        types.String                      `tfsdk:"biz"`
	BizType                    types.String                      `tfsdk:"biz_type"`
	ScaName                    types.String                      `tfsdk:"sca_name"`
	ScaNamePattern             types.String                      `tfsdk:"sca_name_pattern"`
	ScaVersion                 types.String                      `tfsdk:"sca_version"`
	Uuid                       types.String                      `tfsdk:"uuid"`
	SearchItem                 types.String                      `tfsdk:"search_item"`
	SearchInfo                 types.String                      `tfsdk:"search_info"`
	SearchItemSub              types.String                      `tfsdk:"search_item_sub"`
	SearchInfoSub              types.String                      `tfsdk:"search_info_sub"`
	SearchCriteriaList         []propertyScaDetailsCriteriaModel `tfsdk:"search_criteria_list"`
	Remark                     types.String                      `tfsdk:"remark"`
	Port                       types.String                      `tfsdk:"port"`
	User                       types.String                      `tfsdk:"user"`
	Pid                        types.String                      `tfsdk:"pid"`
	ProcessStartedStart        types.Int64                       `tfsdk:"process_started_start"`
	ProcessStartedEnd          types.Int64                       `tfsdk:"process_started_end"`
	ResourceDirectoryAccountId types.Int64                       `tfsdk:"resource_directory_account_id"`
	UseNextToken               types.Bool                        `tfsdk:"use_next_token"`
	NextToken                  types.String                      `tfsdk:"next_token"`
	CurrentPage                types.Int64                       `tfsdk:"current_page"`
	PageSize                   types.Int64                       `tfsdk:"page_size"`
	Properties                 []propertyScaDetailsPropertyModel `tfsdk:"properties"`
	PageInfo                   types.Object                      `tfsdk:"page_info"`
}

var propertyScaDetailsPropertyObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"type":                types.StringType,
	"biz_type_display":    types.StringType,
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
}}

var propertyScaDetailsCriteriaObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}}

var propertyScaDetailsPageInfoObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"current_page": types.Int64Type,
	"page_size":    types.Int64Type,
	"total_count":  types.Int64Type,
	"count":        types.Int64Type,
	"next_token":   types.StringType,
}}

// propertyScaDetailsCriteriaValidator rejects incomplete search conditions at
// plan time, because the API expects every condition to carry both fields.
type propertyScaDetailsCriteriaValidator struct{}

func (v propertyScaDetailsCriteriaValidator) Description(context.Context) string {
	return "Each search condition must set `name` and `value`."
}

func (v propertyScaDetailsCriteriaValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v propertyScaDetailsCriteriaValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	for i, element := range req.ConfigValue.Elements() {
		entryPath := req.Path.AtListIndex(i)
		entry, ok := element.(types.Object)
		// An unknown condition is checked once its value is known.
		if !ok || entry.IsUnknown() {
			continue
		}
		if entry.IsNull() {
			resp.Diagnostics.AddAttributeError(entryPath.AtName("name"), "Invalid Search Condition",
				"`name` must be set for every search condition.")
			continue
		}

		attributes := entry.Attributes()
		if name, ok := attributes["name"].(types.String); ok && !name.IsUnknown() {
			if name.IsNull() || name.ValueString() == "" {
				resp.Diagnostics.AddAttributeError(entryPath.AtName("name"), "Invalid Search Condition",
					"`name` must be set for every search condition, for example \"Name\".")
			}
		}
		if value, ok := attributes["value"].(types.String); ok && !value.IsUnknown() {
			if value.IsNull() || value.ValueString() == "" {
				resp.Diagnostics.AddAttributeError(entryPath.AtName("value"), "Invalid Search Condition",
					"`value` must be set for every search condition.")
			}
		}
	}
}

func (d *propertyScaDetailsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = propertyScaDetailsTypeName
}

func (d *propertyScaDetailsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of asset fingerprint details, including middleware, database and web service fingerprints, collected by Cloud Security Center (SAS).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source.",
			},
			"lang": schema.StringAttribute{
				Optional:    true,
				Description: "The language of the request and the response. Valid values: `zh` and `en`. Default value: `zh`.",
				Validators:  []validator.String{stringvalidator.OneOf("zh", "en")},
			},
			"biz": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the asset fingerprint. Valid values: `sca`, `sca_database` and `sca_web`. Default value: `sca`.",
				Validators:  []validator.String{stringvalidator.OneOf("sca", "sca_database", "sca_web")},
			},
			"biz_type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the middleware, database or web service. Valid values: `system_service`, `software_library`, `docker_component`, `database`, `web_container`, `jar` and `web_framework`.",
				Validators:  []validator.String{stringvalidator.OneOf(propertyScaDetailsBizTypes...)},
			},
			"sca_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the asset fingerprint to query.",
			},
			"sca_name_pattern": schema.StringAttribute{
				Optional:    true,
				Description: "The process name of the asset fingerprint to query. It is matched by prefix.",
			},
			"sca_version": schema.StringAttribute{
				Optional:    true,
				Description: "The version of the middleware, database or web service.",
			},
			"uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the server that hosts the asset fingerprint.",
			},
			"search_item": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the search condition. Valid values: `name` and `type`. It must be set together with `search_info`.",
				Validators: []validator.String{
					stringvalidator.OneOf("name", "type"),
					stringvalidator.AlsoRequires(path.MatchRoot("search_info")),
				},
			},
			"search_info": schema.StringAttribute{
				Optional:    true,
				Description: "The search content. When `search_item` is `name`, it is the asset fingerprint name. When `search_item` is `type`, it is one of the `biz_type` values. It must be set together with `search_item` and must not be empty.",
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("search_item")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"search_item_sub": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the sub-search condition. Valid values: `port`, `pid`, `version` and `user`. It must be set together with `search_info_sub`.",
				Validators: []validator.String{
					stringvalidator.OneOf("port", "pid", "version", "user"),
					stringvalidator.AlsoRequires(path.MatchRoot("search_info_sub")),
				},
			},
			"search_info_sub": schema.StringAttribute{
				Optional:    true,
				Description: "The content of the sub-search condition, matching the type set in `search_item_sub`. It must be set together with `search_item_sub` and must not be empty.",
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("search_item_sub")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"search_criteria_list": schema.ListAttribute{
				Optional:    true,
				ElementType: propertyScaDetailsCriteriaObjectType,
				Description: fmt.Sprintf("The search conditions. Up to %d conditions are supported. Each condition sets `name` and `value`.", propertyScaDetailsMaxCriteriaItems),
				Validators: []validator.List{
					listvalidator.SizeAtMost(propertyScaDetailsMaxCriteriaItems),
					propertyScaDetailsCriteriaValidator{},
				},
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Description: "The server name or IP address of the asset fingerprint. Fuzzy match is supported.",
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
				Description: "The beginning of the process start time range to query, as a second-level timestamp.",
			},
			"process_started_end": schema.Int64Attribute{
				Optional:    true,
				Description: "The end of the process start time range to query, as a second-level timestamp.",
			},
			"resource_directory_account_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The ID of the Resource Directory member account.",
			},
			"use_next_token": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to page through the results with `NextToken` instead of `CurrentPage` and `PageSize`. When it is not set or is `false`, the offset pagination applies and `page_info.total_count` is returned. When it is `true`, `current_page` must not be set and `page_info.total_count` is not returned.",
			},
			"next_token": schema.StringAttribute{
				Optional:    true,
				Description: "The `NextToken` to start reading from. Only used when `use_next_token` is `true`; leave it empty to start from the beginning. All remaining pages are then collected in one read.",
			},
			"current_page": schema.Int64Attribute{
				Optional:    true,
				Description: "The page to start reading from when `use_next_token` is `false`. Default value: 1. All pages from this page to the last one are collected in one read.",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: fmt.Sprintf("The number of entries to return on each page. Default value: %d.", propertyScaDetailsDefaultPageSize),
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"properties": schema.ListAttribute{
				Computed:    true,
				ElementType: propertyScaDetailsPropertyObjectType,
				Description: "The asset fingerprint entries returned by `DescribePropertyScaDetail`. Entries collected over multiple pages are concatenated in page order.",
			},
			"page_info": schema.ObjectAttribute{
				Computed:       true,
				AttributeTypes: propertyScaDetailsPageInfoObjectType.AttrTypes,
				Description:    "The pagination information of the query. `current_page` is the first page read (`null` when `use_next_token` is `true`), `page_size` is the requested page size, `total_count` is the total number of matches reported by the API (`null` when `use_next_token` is `true`), `count` is the number of returned `properties` after all pages were collected, and `next_token` is the trailing `NextToken` of the NextToken pagination (`null` when `use_next_token` is `false`).",
			},
		},
	}
}

func (d *propertyScaDetailsDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var model propertyScaDetailsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if model.UseNextToken.IsUnknown() {
		return
	}

	if model.UseNextToken.ValueBool() {
		if !model.CurrentPage.IsNull() && !model.CurrentPage.IsUnknown() {
			resp.Diagnostics.AddAttributeError(path.Root("current_page"), "Invalid Pagination Configuration",
				"`current_page` cannot be set when `use_next_token` is `true`; NextToken pagination does not use page numbers.")
		}
		return
	}
	if !model.NextToken.IsNull() && !model.NextToken.IsUnknown() && model.NextToken.ValueString() != "" {
		resp.Diagnostics.AddAttributeError(path.Root("next_token"), "Invalid Pagination Configuration",
			"`next_token` can only be set when `use_next_token` is `true`.")
		return
	}
}

func (d *propertyScaDetailsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model propertyScaDetailsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pageSize := int64(propertyScaDetailsDefaultPageSize)
	if !model.PageSize.IsNull() && !model.PageSize.IsUnknown() {
		pageSize = model.PageSize.ValueInt64()
	}
	useNextToken := model.UseNextToken.ValueBool()
	firstPage := int64(1)
	if !model.CurrentPage.IsNull() && !model.CurrentPage.IsUnknown() {
		firstPage = model.CurrentPage.ValueInt64()
	}

	baseRequest := propertyScaDetailsRequest(&model)
	properties := make([]interface{}, 0)
	lastPageInfo := map[string]interface{}{}

	if useNextToken {
		nextToken := ""
		if !model.NextToken.IsNull() && !model.NextToken.IsUnknown() {
			nextToken = model.NextToken.ValueString()
		}
		for {
			request := propertyScaDetailsPagedRequest(baseRequest, pageSize)
			request["UseNextToken"] = true
			if nextToken != "" {
				request["NextToken"] = nextToken
			}
			response, err := d.describePropertyScaDetail(ctx, request)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Reading %s: calling %s", propertyScaDetailsTypeName, propertyScaDetailsAction),
					err.Error(),
				)
				return
			}

			properties = append(properties, propertyScaDetailsResponseProperties(response)...)
			lastPageInfo = propertyScaDetailsResponsePageInfo(response)
			nextToken = propertyScaDetailsPageInfoString(lastPageInfo, "NextToken")
			if nextToken == "" {
				break
			}
		}
	} else {
		currentPage := firstPage
		for {
			request := propertyScaDetailsPagedRequest(baseRequest, pageSize)
			request["CurrentPage"] = currentPage
			response, err := d.describePropertyScaDetail(ctx, request)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Reading %s: calling %s", propertyScaDetailsTypeName, propertyScaDetailsAction),
					err.Error(),
				)
				return
			}

			pageItems := propertyScaDetailsResponseProperties(response)
			properties = append(properties, pageItems...)
			lastPageInfo = propertyScaDetailsResponsePageInfo(response)

			// The API can cap PageSize below the requested value, so the loop
			// follows the returned page size and the collected count instead of
			// the requested page size.
			returnedPageSize := propertyScaDetailsPageInfoInt64(lastPageInfo, "PageSize")
			if returnedPageSize <= 0 {
				returnedPageSize = pageSize
			}
			if len(pageItems) == 0 || int64(len(properties)) >= propertyScaDetailsPageInfoInt64(lastPageInfo, "TotalCount") || int64(len(pageItems)) < returnedPageSize {
				break
			}
			currentPage++
		}
	}

	propertyModels := make([]propertyScaDetailsPropertyModel, 0, len(properties))
	for _, property := range properties {
		item, ok := property.(map[string]interface{})
		if !ok {
			continue
		}
		propertyModels = append(propertyModels, propertyScaDetailsProperty(item))
	}

	pageInfo := propertyScaDetailsPageInfoValue{
		CurrentPage: types.Int64Null(),
		PageSize:    types.Int64Value(pageSize),
		TotalCount:  types.Int64Null(),
		Count:       types.Int64Value(int64(len(propertyModels))),
		NextToken:   types.StringNull(),
	}
	if useNextToken {
		pageInfo.NextToken = types.StringValue(propertyScaDetailsPageInfoString(lastPageInfo, "NextToken"))
	} else {
		pageInfo.CurrentPage = types.Int64Value(firstPage)
		pageInfo.TotalCount = types.Int64Value(propertyScaDetailsPageInfoInt64(lastPageInfo, "TotalCount"))
	}
	pageInfoObject, diags := types.ObjectValueFrom(ctx, propertyScaDetailsPageInfoObjectType.AttrTypes, pageInfo)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.Id = types.StringValue(propertyScaDetailsId(d.Client().RegionId, baseRequest, useNextToken))
	model.Properties = propertyModels
	model.PageInfo = pageInfoObject

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// describePropertyScaDetail reads one page and retries the throttling and
// transient errors that a multi-page collection can hit. The API allows 10
// requests per second, so a full collection has to back off on contention.
func (d *propertyScaDetailsDataSource) describePropertyScaDetail(ctx context.Context, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := sdkretry.RetryContext(ctx, propertyScaDetailsQueryTimeout, func() *sdkretry.RetryError {
		result, rpcErr := d.Client().RpcPost("Sas", "2018-12-03", propertyScaDetailsAction, nil, request, true)
		if rpcErr != nil {
			if helper.NeedRetry(rpcErr) {
				wait()
				return sdkretry.RetryableError(rpcErr)
			}
			return sdkretry.NonRetryableError(rpcErr)
		}
		response = result
		return nil
	})
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, propertyScaDetailsTypeName, propertyScaDetailsAction, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(propertyScaDetailsAction, response, request)
	return response, nil
}

func propertyScaDetailsRequest(model *propertyScaDetailsModel) map[string]interface{} {
	request := make(map[string]interface{}, 21)
	propertyScaDetailsStringParam(request, "Lang", model.Lang)
	propertyScaDetailsStringParam(request, "Biz", model.Biz)
	propertyScaDetailsStringParam(request, "BizType", model.BizType)
	propertyScaDetailsStringParam(request, "ScaName", model.ScaName)
	propertyScaDetailsStringParam(request, "ScaNamePattern", model.ScaNamePattern)
	propertyScaDetailsStringParam(request, "ScaVersion", model.ScaVersion)
	propertyScaDetailsStringParam(request, "Uuid", model.Uuid)
	propertyScaDetailsStringParam(request, "SearchItem", model.SearchItem)
	propertyScaDetailsStringParam(request, "SearchInfo", model.SearchInfo)
	propertyScaDetailsStringParam(request, "SearchItemSub", model.SearchItemSub)
	propertyScaDetailsStringParam(request, "SearchInfoSub", model.SearchInfoSub)
	propertyScaDetailsStringParam(request, "Remark", model.Remark)
	propertyScaDetailsStringParam(request, "Port", model.Port)
	propertyScaDetailsStringParam(request, "User", model.User)
	propertyScaDetailsStringParam(request, "Pid", model.Pid)
	propertyScaDetailsInt64Param(request, "ProcessStartedStart", model.ProcessStartedStart)
	propertyScaDetailsInt64Param(request, "ProcessStartedEnd", model.ProcessStartedEnd)
	propertyScaDetailsInt64Param(request, "ResourceDirectoryAccountId", model.ResourceDirectoryAccountId)
	if criteria := propertyScaDetailsExpandCriteria(model.SearchCriteriaList); len(criteria) > 0 {
		request["SearchCriteriaList"] = criteria
	}
	return request
}

func propertyScaDetailsPagedRequest(baseRequest map[string]interface{}, pageSize int64) map[string]interface{} {
	request := make(map[string]interface{}, len(baseRequest)+3)
	for key, value := range baseRequest {
		request[key] = value
	}
	request["PageSize"] = pageSize
	return request
}

func propertyScaDetailsStringParam(request map[string]interface{}, key string, value types.String) {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return
	}
	request[key] = value.ValueString()
}

func propertyScaDetailsInt64Param(request map[string]interface{}, key string, value types.Int64) {
	if value.IsNull() || value.IsUnknown() {
		return
	}
	request[key] = value.ValueInt64()
}

func propertyScaDetailsExpandCriteria(criteria []propertyScaDetailsCriteriaModel) []interface{} {
	expanded := make([]interface{}, 0, len(criteria))
	for _, criterion := range criteria {
		entry := map[string]interface{}{}
		if !criterion.Name.IsNull() && !criterion.Name.IsUnknown() {
			entry["Name"] = criterion.Name.ValueString()
		}
		if !criterion.Value.IsNull() && !criterion.Value.IsUnknown() {
			entry["Value"] = criterion.Value.ValueString()
		}
		expanded = append(expanded, entry)
	}
	return expanded
}

func propertyScaDetailsResponseProperties(response map[string]interface{}) []interface{} {
	raw, ok := response["Propertys"]
	if !ok || raw == nil {
		return nil
	}
	properties, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	return properties
}

func propertyScaDetailsResponsePageInfo(response map[string]interface{}) map[string]interface{} {
	raw, ok := response["PageInfo"]
	if !ok || raw == nil {
		return map[string]interface{}{}
	}
	pageInfo, ok := raw.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return pageInfo
}

func propertyScaDetailsProperty(item map[string]interface{}) propertyScaDetailsPropertyModel {
	return propertyScaDetailsPropertyModel{
		Type:              propertyScaDetailsResponseString(item, "Type"),
		BizTypeDisplay:    propertyScaDetailsResponseString(item, "BizTypeDispaly"),
		ProcessStarted:    propertyScaDetailsResponseInt64(item, "ProcessStarted"),
		InternetIp:        propertyScaDetailsResponseString(item, "InternetIp"),
		ConfigPath:        propertyScaDetailsResponseString(item, "ConfigPath"),
		Pid:               propertyScaDetailsResponseString(item, "Pid"),
		Port:              propertyScaDetailsResponseString(item, "Port"),
		Cmdline:           propertyScaDetailsResponseString(item, "Cmdline"),
		BizType:           propertyScaDetailsResponseString(item, "BizType"),
		ListenIp:          propertyScaDetailsResponseString(item, "ListenIp"),
		Version:           propertyScaDetailsResponseString(item, "Version"),
		InstanceName:      propertyScaDetailsResponseString(item, "InstanceName"),
		ListenStatus:      propertyScaDetailsResponseString(item, "ListenStatus"),
		Name:              propertyScaDetailsResponseString(item, "Name"),
		Ip:                propertyScaDetailsResponseString(item, "Ip"),
		ProcessUser:       propertyScaDetailsResponseString(item, "ProcessUser"),
		InstanceId:        propertyScaDetailsResponseString(item, "InstanceId"),
		WebPath:           propertyScaDetailsResponseString(item, "WebPath"),
		Ppid:              propertyScaDetailsResponseString(item, "Ppid"),
		IntranetIp:        propertyScaDetailsResponseString(item, "IntranetIp"),
		Uuid:              propertyScaDetailsResponseString(item, "Uuid"),
		ListenProtocol:    propertyScaDetailsResponseString(item, "ListenProtocol"),
		ImageName:         propertyScaDetailsResponseString(item, "ImageName"),
		Path:              propertyScaDetailsResponseString(item, "Path"),
		ContainerName:     propertyScaDetailsResponseString(item, "ContainerName"),
		Proof:             propertyScaDetailsResponseString(item, "Proof"),
		CreateTimestamp:   propertyScaDetailsResponseInt64(item, "CreateTimestamp"),
		RuntimeEnvVersion: propertyScaDetailsResponseString(item, "RuntimeEnvVersion"),
		PodName:           propertyScaDetailsResponseString(item, "PodName"),
	}
}

func propertyScaDetailsResponseString(item map[string]interface{}, key string) types.String {
	value, ok := item[key]
	if !ok || value == nil {
		return types.StringNull()
	}
	if text, ok := value.(string); ok {
		return types.StringValue(text)
	}
	return types.StringValue(fmt.Sprint(value))
}

func propertyScaDetailsResponseInt64(item map[string]interface{}, key string) types.Int64 {
	value, ok := item[key]
	if !ok || value == nil {
		return types.Int64Null()
	}
	return propertyScaDetailsInt64Value(value)
}

func propertyScaDetailsPageInfoString(pageInfo map[string]interface{}, key string) string {
	value, ok := pageInfo[key]
	if !ok || value == nil {
		return ""
	}
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return text
}

func propertyScaDetailsPageInfoInt64(pageInfo map[string]interface{}, key string) int64 {
	value, ok := pageInfo[key]
	if !ok || value == nil {
		return 0
	}
	return propertyScaDetailsInt64Value(value).ValueInt64()
}

func propertyScaDetailsInt64Value(value interface{}) types.Int64 {
	switch number := value.(type) {
	case json.Number:
		parsed, err := number.Int64()
		if err != nil {
			return types.Int64Null()
		}
		return types.Int64Value(parsed)
	case float64:
		return types.Int64Value(int64(number))
	case int64:
		return types.Int64Value(number)
	case int:
		return types.Int64Value(int64(number))
	case string:
		var parsed int64
		if _, err := fmt.Sscan(number, &parsed); err != nil {
			return types.Int64Null()
		}
		return types.Int64Value(parsed)
	}
	return types.Int64Null()
}

func propertyScaDetailsId(regionId string, request map[string]interface{}, useNextToken bool) string {
	var buf bytes.Buffer
	buf.WriteString(regionId)
	buf.WriteString(fmt.Sprintf("-t%t-", useNextToken))
	raw, err := json.Marshal(request)
	if err == nil {
		buf.Write(raw)
	}
	return fmt.Sprintf("%d", helper.Hashcode(buf.String()))
}
