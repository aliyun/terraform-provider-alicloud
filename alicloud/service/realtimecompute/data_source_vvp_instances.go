package realtimecompute

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	vvpInstancesTypeName = "alicloud_realtime_compute_vvp_instances"
	//vvpInstancesAction   = "DescribeInstances"
	//vvpInstancesProduct  = "foasconsole"
	//// vvpInstancesVersion matches the foasconsole version already used by the
	//// alicloud_realtime_compute_vvp_instance resource and its service helpers.
	//vvpInstancesVersion = "2019-06-01"
	//// vvpInstancesPageSize is the maximum page size allowed by DescribeInstances.
	//vvpInstancesPageSize = 100
	//// vvpInstancesMaxPages bounds the pagination loop so an endpoint that
	//// ignores the page parameters cannot loop forever.
	//vvpInstancesMaxPages = 100
	//// vvpInstancesRequestPrefix is the foasconsole RPC parameter prefix; every
	//// DescribeInstances query parameter is flattened under it (see service.go).
	//vvpInstancesRequestPrefix = "DescribeInstancesRequest."
)

var (
	_ datasource.DataSource              = &vvpInstancesDataSource{}
	_ datasource.DataSourceWithConfigure = &vvpInstancesDataSource{}
)

func NewVvpInstancesDataSource() datasource.DataSource {
	return &vvpInstancesDataSource{}
}

type vvpInstancesDataSource struct {
	fwadapt.DataSourceBase
}

type vvpInstancesModel struct {
	Ids             types.List   `tfsdk:"ids"`
	NameRegex       types.String `tfsdk:"name_regex"`
	InstanceId      types.String `tfsdk:"instance_id"`
	InstanceName    types.String `tfsdk:"instance_name"`
	PaymentType     types.String `tfsdk:"payment_type"`
	NamespaceName   types.String `tfsdk:"namespace_name"`
	ResourceGroupId types.String `tfsdk:"resource_group_id"`
	OutputFile      types.String `tfsdk:"output_file"`
	Id              types.String `tfsdk:"id"`
	Names           types.List   `tfsdk:"names"`
	Instances       types.List   `tfsdk:"instances"`
}

var vvpInstanceAttrTypes = map[string]attr.Type{
	"id":                  types.StringType,
	"instance_id":         types.StringType,
	"instance_name":       types.StringType,
	"payment_type":        types.StringType,
	"status":              types.StringType,
	"region":              types.StringType,
	"zone_id":             types.StringType,
	"ha_zone_id":          types.StringType,
	"resource_group_id":   types.StringType,
	"resource_id":         types.StringType,
	"vpc_id":              types.StringType,
	"create_time":         types.StringType,
	"expire_time":         types.StringType,
	"architecture_type":   types.StringType,
	"order_state":         types.StringType,
	"monitor_type":        types.StringType,
	"elastic_instance_id": types.StringType,
	"elastic_order_state": types.StringType,
	"uid":                 types.StringType,
	"deletion_protection": types.BoolType,
	"elastic":             types.BoolType,
	"ha":                  types.BoolType,
	"ansm":                types.BoolType,
	"vswitch_ids":         types.ListType{ElemType: types.StringType},
	"tags":                types.MapType{ElemType: types.StringType},
	"resource_spec":       types.ObjectType{AttrTypes: resourceSpecAttrTypes},
	"storage":             types.ObjectType{AttrTypes: storageAttrTypes},
	"vpc_info":            types.ObjectType{AttrTypes: vpcInfoAttrTypes},
}

var resourceSpecAttrTypes = map[string]attr.Type{
	"cpu":       types.Int64Type,
	"memory_gb": types.Int64Type,
}

var storageOssAttrTypes = map[string]attr.Type{
	"bucket": types.StringType,
}

var storageAttrTypes = map[string]attr.Type{
	"fully_managed": types.BoolType,
	"order_state":   types.StringType,
	"oss":           types.ObjectType{AttrTypes: storageOssAttrTypes},
}

var vpcInfoAttrTypes = map[string]attr.Type{
	"vpc_id":      types.StringType,
	"cidr_block":  types.StringType,
	"vpc_name":    types.StringType,
	"description": types.StringType,
	"region_id":   types.StringType,
	"status":      types.StringType,
}

func (d *vvpInstancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = vvpInstancesTypeName
}

func (d *vvpInstancesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Realtime Compute for Apache Flink (VVP) Instances of the current Alibaba Cloud user.",
		Attributes: map[string]schema.Attribute{
			"ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of VVP Instance IDs.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by VVP Instance name.",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the Instance.",
			},
			"instance_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the Instance.",
			},
			"payment_type": schema.StringAttribute{
				Optional:    true,
				Description: "The billing method. Valid values: `PayAsYouGo` and `Subscription`.",
			},
			"namespace_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the namespace.",
			},
			"resource_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the resource group.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to save the data source results.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, which is derived from the query condition.",
			},
			"names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of VVP instance names.",
			},
			"instances": schema.ListAttribute{
				Computed:    true,
				ElementType: types.ObjectType{AttrTypes: vvpInstanceAttrTypes},
				Description: "The list of VVP instances that match the configured filters.",
			},
		},
	}
}

func (d *vvpInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	const action = "DescribeInstances"

	var config vvpInstancesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{}

	setString := func(key string, value types.String) {
		if !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
			request[key] = value.ValueString()
		}
	}

	setString("InstanceId", config.InstanceId)
	setString("InstanceName", config.InstanceName)
	setString("ChargeType", convertRealtimeComputeCreateInstanceRequestChargeTypeRequest(config.PaymentType))
	setString("NamespaceName", config.NamespaceName)
	setString("ResourceGroupId", config.ResourceGroupId)

	idsMap := make(map[string]string)
	if !config.Ids.IsNull() && !config.Ids.IsUnknown() {
		var idList []string
		resp.Diagnostics.Append(config.Ids.ElementsAs(ctx, &idList, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range idList {
			idsMap[v] = v
		}
	}

	var nameRegex *regexp.Regexp
	if v := config.NameRegex.ValueString(); !config.NameRegex.IsNull() && v != "" {
		compiled, err := regexp.Compile(v)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: invalid name_regex", vvpInstancesTypeName),
				err.Error(),
			)
			return
		}
		nameRegex = compiled
	}

	var objects []map[string]interface{}
	for pageIndex := 1; pageIndex <= vvpInstancesMaxPages; pageIndex++ {
		query := map[string]interface{}{
			vvpInstancesRequestPrefix + "Region":    client.RegionId,
			vvpInstancesRequestPrefix + "PageIndex": pageIndex,
			vvpInstancesRequestPrefix + "PageSize":  vvpInstancesPageSize,
		}
		setQuery := func(field string, value types.String) {
			if !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
				query[vvpInstancesRequestPrefix+field] = value.ValueString()
			}
		}
		setQuery("InstanceId", config.InstanceId)
		setQuery("InstanceName", config.InstanceName)
		setQuery("NamespaceName", config.NamespaceName)
		setQuery("ResourceGroupId", config.ResourceGroupId)
		if v := config.PaymentType.ValueString(); !config.PaymentType.IsNull() && v != "" {
			query[vvpInstancesRequestPrefix+"ChargeType"] = convertRealtimeComputeCreateInstanceRequestChargeTypeRequest(v)
		}

		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
			var rpcErr error
			response, rpcErr = client.RpcGet(vvpInstancesProduct, vvpInstancesVersion, vvpInstancesAction, query, nil)
			if rpcErr != nil {
				if helper.NeedRetry(rpcErr) {
					wait()
					return retry.RetryableError(rpcErr)
				}
				return retry.NonRetryableError(rpcErr)
			}
			return nil
		})
		logs.AddDebug(vvpInstancesAction, response, query)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: calling %s", vvpInstancesTypeName, vvpInstancesAction),
				err.Error(),
			)
			return
		}

		pageItems := vvpInstancesItems(response)
		objects = append(objects, pageItems...)

		totalCount := int(vvpInstanceInt64(response["TotalCount"]))
		if len(pageItems) < vvpInstancesPageSize {
			break
		}
		if totalCount > 0 && len(objects) >= totalCount {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	items := make([]attr.Value, 0, len(objects))
	records := make([]map[string]interface{}, 0, len(objects))
	for _, object := range objects {
		instanceId := vvpInstanceString(object["InstanceId"])
		if len(idsMap) > 0 {
			if _, ok := idsMap[instanceId]; !ok {
				continue
			}
		}

		instanceName := vvpInstanceString(object["InstanceName"])
		if nameRegex != nil && !nameRegex.MatchString(instanceName) {
			continue
		}

		ids = append(ids, instanceId)
		names = append(names, instanceName)
		items = append(items, vvpInstanceObject(object))
		records = append(records, vvpInstanceRecord(object))
	}

	idsList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	namesList, diags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instancesList, diags := types.ListValue(vvpInstanceObjectType, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := config
	state.Id = types.StringValue(helper.DataResourceIdHash(ids))
	state.Ids = idsList
	state.Names = namesList
	state.Instances = instancesList

	if output := config.OutputFile.ValueString(); output != "" {
		if err := helper.WriteToFile(output, records); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: writing output_file", vvpInstancesTypeName),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// vvpInstancesItems parses the Instances array of one DescribeInstances page.
// A missing key or an unexpected shape yields an empty slice, never an error.
func vvpInstancesItems(response map[string]interface{}) []map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	if response == nil {
		return items
	}
	raw, err := jsonpath.Get("$.Instances", response)
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

// vvpInstanceObject maps one raw instance to its framework object value. The
// keys must stay in sync with vvpInstanceAttrTypes.
func vvpInstanceObject(object map[string]interface{}) attr.Value {
	instanceId := vvpInstanceString(object["InstanceId"])
	return types.ObjectValueMust(vvpInstanceAttrTypes, map[string]attr.Value{
		"id":                    types.StringValue(instanceId),
		"instance_id":           types.StringValue(instanceId),
		"instance_name":         vvpInstanceStringValue(object["InstanceName"]),
		"payment_type":          types.StringValue(vvpInstanceString(convertRealtimeComputeInstancesChargeTypeResponse(object["ChargeType"]))),
		"status":                vvpInstanceStringValue(object["ClusterStatus"]),
		"region":                vvpInstanceStringValue(object["Region"]),
		"zone_id":               vvpInstanceStringValue(object["ZoneId"]),
		"ha_zone_id":            vvpInstanceStringValue(object["HaZoneId"]),
		"resource_group_id":     vvpInstanceStringValue(object["ResourceGroupId"]),
		"resource_id":           vvpInstanceStringValue(object["ResourceId"]),
		"vpc_id":                vvpInstanceStringValue(object["VpcId"]),
		"create_time":           types.StringValue(vvpInstanceInt64String(object["ResourceCreateTime"])),
		"expire_time":           types.StringValue(vvpInstanceInt64String(object["ResourceExpiredTime"])),
		"architecture_type":     vvpInstanceStringValue(object["ArchitectureType"]),
		"order_state":           vvpInstanceStringValue(object["OrderState"]),
		"monitor_type":          vvpInstanceStringValue(object["MonitorType"]),
		"elastic_instance_id":   vvpInstanceStringValue(object["ElasticInstanceId"]),
		"elastic_order_state":   vvpInstanceStringValue(object["ElasticOrderState"]),
		"uid":                   vvpInstanceStringValue(object["Uid"]),
		"deletion_protection":   types.BoolValue(vvpInstanceBool(object["DeletionProtection"])),
		"elastic":               types.BoolValue(vvpInstanceBool(object["Elastic"])),
		"ha":                    types.BoolValue(vvpInstanceBool(object["Ha"])),
		"ansm":                  types.BoolValue(vvpInstanceBool(object["Ansm"])),
		"vswitch_ids":           vvpInstanceStringList(object["VSwitchIds"]),
		"tags":                  vvpInstanceTagsValue(object["Tags"]),
		"resource_spec":         vvpInstanceResourceSpecObject(object["ResourceSpec"]),
		"ha_resource_spec":      vvpInstanceResourceSpecObject(object["HaResourceSpec"]),
		"elastic_resource_spec": vvpInstanceResourceSpecObject(object["ElasticResourceSpec"]),
		"storage":               vvpInstanceStorageObject(object["Storage"]),
		"vpc_info":              vvpInstanceVpcInfoObject(object["VpcInfo"]),
		"oss_info":              vvpInstanceOssInfoObject(object["OssInfo"]),
	})
}

// vvpInstanceRecord builds the output_file payload for one instance. It keeps
// the most useful fields; the exhaustive set lives in the instances attribute.
func vvpInstanceRecord(object map[string]interface{}) map[string]interface{} {
	instanceId := vvpInstanceString(object["InstanceId"])
	return map[string]interface{}{
		"id":                instanceId,
		"instance_id":       instanceId,
		"instance_name":     vvpInstanceString(object["InstanceName"]),
		"payment_type":      vvpInstanceString(convertRealtimeComputeInstancesChargeTypeResponse(object["ChargeType"])),
		"status":            vvpInstanceString(object["ClusterStatus"]),
		"region":            vvpInstanceString(object["Region"]),
		"zone_id":           vvpInstanceString(object["ZoneId"]),
		"resource_group_id": vvpInstanceString(object["ResourceGroupId"]),
		"resource_id":       vvpInstanceString(object["ResourceId"]),
		"vpc_id":            vvpInstanceString(object["VpcId"]),
		"create_time":       vvpInstanceInt64String(object["ResourceCreateTime"]),
		"expire_time":       vvpInstanceInt64String(object["ResourceExpiredTime"]),
		"tags":              tags.TagsToMap(object["Tags"]),
	}
}

func vvpInstanceResourceSpecObject(raw interface{}) attr.Value {
	m, ok := raw.(map[string]interface{})
	if !ok || len(m) == 0 {
		return types.ObjectNull(vvpResourceSpecAttrTypes)
	}
	return types.ObjectValueMust(vvpResourceSpecAttrTypes, map[string]attr.Value{
		"cpu":       vvpInstanceInt64Value(m["Cpu"]),
		"memory_gb": vvpInstanceInt64Value(m["MemoryGB"]),
	})
}

func vvpInstanceStorageObject(raw interface{}) attr.Value {
	m, ok := raw.(map[string]interface{})
	if !ok || len(m) == 0 {
		return types.ObjectNull(vvpStorageAttrTypes)
	}
	return types.ObjectValueMust(vvpStorageAttrTypes, map[string]attr.Value{
		"fully_managed": types.BoolValue(vvpInstanceBool(m["FullyManaged"])),
		"order_state":   vvpInstanceStringValue(m["OrderState"]),
		"oss":           vvpInstanceStorageOssObject(m["Oss"]),
	})
}

func vvpInstanceStorageOssObject(raw interface{}) attr.Value {
	m, ok := raw.(map[string]interface{})
	if !ok || len(m) == 0 {
		return types.ObjectNull(vvpStorageOssAttrTypes)
	}
	return types.ObjectValueMust(vvpStorageOssAttrTypes, map[string]attr.Value{
		"bucket": vvpInstanceStringValue(m["Bucket"]),
	})
}

func vvpInstanceVpcInfoObject(raw interface{}) attr.Value {
	m, ok := raw.(map[string]interface{})
	if !ok || len(m) == 0 {
		return types.ObjectNull(vvpVpcInfoAttrTypes)
	}
	return types.ObjectValueMust(vvpVpcInfoAttrTypes, map[string]attr.Value{
		"cidr_block":  vvpInstanceStringValue(m["CidrBlock"]),
		"description": vvpInstanceStringValue(m["Description"]),
		"region_id":   vvpInstanceStringValue(m["RegionId"]),
		"status":      vvpInstanceStringValue(m["Status"]),
		"vpc_id":      vvpInstanceStringValue(m["VpcId"]),
		"vpc_name":    vvpInstanceStringValue(m["VpcName"]),
	})
}

func vvpInstanceOssInfoObject(raw interface{}) attr.Value {
	m, ok := raw.(map[string]interface{})
	if !ok || len(m) == 0 {
		return types.ObjectNull(vvpOssInfoAttrTypes)
	}
	return types.ObjectValueMust(vvpOssInfoAttrTypes, map[string]attr.Value{
		"bucket":                   vvpInstanceStringValue(m["Bucket"]),
		"endpoint":                 vvpInstanceStringValue(m["Endpoint"]),
		"bucket_versioning_status": vvpInstanceStringValue(m["BucketVersioningStatus"]),
	})
}

func vvpInstanceStringList(raw interface{}) attr.Value {
	arr, ok := raw.([]interface{})
	if !ok {
		return types.ListNull(types.StringType)
	}
	elements := make([]attr.Value, 0, len(arr))
	for _, v := range arr {
		elements = append(elements, vvpInstanceStringValue(v))
	}
	return types.ListValueMust(types.StringType, elements)
}

func vvpInstanceTagsValue(raw interface{}) attr.Value {
	tagsMap := tags.TagsToMap(raw)
	elements := make(map[string]attr.Value, len(tagsMap))
	for k, v := range tagsMap {
		elements[k] = types.StringValue(fmt.Sprint(v))
	}
	return types.MapValueMust(types.StringType, elements)
}

func vvpInstanceString(v interface{}) string {
	switch value := v.(type) {
	case nil:
		return ""
	case string:
		return value
	default:
		return fmt.Sprint(value)
	}
}

func vvpInstanceStringValue(v interface{}) attr.Value {
	return types.StringValue(vvpInstanceString(v))
}

func vvpInstanceBool(v interface{}) bool {
	switch value := v.(type) {
	case bool:
		return value
	case string:
		return value == "true"
	default:
		return false
	}
}

func vvpInstanceInt64(v interface{}) int64 {
	switch value := v.(type) {
	case float64:
		return int64(value)
	case float32:
		return int64(value)
	case int:
		return int64(value)
	case int64:
		return value
	case json.Number:
		i, _ := value.Int64()
		return i
	case string:
		i, _ := strconv.ParseInt(value, 10, 64)
		return i
	default:
		return 0
	}
}

func vvpInstanceInt64Value(v interface{}) attr.Value {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(vvpInstanceInt64(v))
}

// vvpInstanceInt64String renders a millisecond epoch as a decimal string so a
// float64 decoded from JSON is never emitted in scientific notation.
func vvpInstanceInt64String(v interface{}) string {
	if v == nil {
		return ""
	}
	return strconv.FormatInt(vvpInstanceInt64(v), 10)
}
