package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBGatewayModelAPIs() *schema.Resource {
	return &schema.Resource{Read: dataSourceAlicloudPolarDBGatewayModelAPIsRead, Schema: map[string]*schema.Schema{
		"gateway_id":     {Type: schema.TypeString, Required: true},
		"ids":            {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"name":           {Type: schema.TypeString, Optional: true},
		"model_category": {Type: schema.TypeString, Optional: true, ValidateFunc: StringInSlice([]string{"text", "embedding", "rerank"}, false)},
		"path_prefix":    {Type: schema.TypeString, Optional: true},
		"protocol":       {Type: schema.TypeString, Optional: true, ValidateFunc: StringInSlice([]string{"openai", "anthropic", "bailian", "vllm"}, false)},
		"status":         {Type: schema.TypeString, Optional: true},
		"apis": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id": {Type: schema.TypeString, Computed: true}, "name": {Type: schema.TypeString, Computed: true},
			"model_category": {Type: schema.TypeString, Computed: true}, "path_prefix": {Type: schema.TypeString, Computed: true},
			"protocol": {Type: schema.TypeString, Computed: true}, "status": {Type: schema.TypeString, Computed: true},
			"record_input": {Type: schema.TypeString, Computed: true}, "record_output": {Type: schema.TypeString, Computed: true},
			"route_rules": {Type: schema.TypeString, Computed: true}, "force_model": {Type: schema.TypeString, Computed: true},
			"invoke_endpoint": {Type: schema.TypeString, Computed: true}, "create_time": {Type: schema.TypeString, Computed: true},
		}}},
	}}
}

func dataSourceAlicloudPolarDBGatewayModelAPIsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{"GwClusterId": d.Get("gateway_id")}
	copyPolarDBGatewayAIFilters(d, request, map[string]string{"name": "Name", "model_category": "ModelCategory", "path_prefix": "PathPrefix", "protocol": "Protocol", "status": "Status"})
	if ids := stringListFromSchema(d.Get("ids")); len(ids) > 0 {
		request["ModelApiIds"] = strings.Join(ids, ",")
	}
	objects, err := listPolarDBGatewayAI(client, "DescribeModelApis", request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_gateway_model_apis", "DescribeModelApis", AlibabaCloudSdkGoERROR)
	}
	items := make([]map[string]interface{}, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, object := range objects {
		id := fmt.Sprint(object["ModelApiId"])
		category := object["Category"]
		if category == nil {
			category = object["ModelCategory"]
		}
		ids = append(ids, id)
		items = append(items, map[string]interface{}{
			"id": id, "name": object["Name"], "model_category": category, "path_prefix": object["PathPrefix"],
			"protocol": object["Protocol"], "status": object["Status"], "record_input": object["RecordInput"],
			"record_output": object["RecordOutput"], "route_rules": normalizePolarDBGatewayAIJSON(object["RouteRules"]),
			"force_model": object["ForceModel"], "invoke_endpoint": object["InvokeEndpoint"], "create_time": object["GmtCreated"],
		})
	}
	if err = d.Set("apis", items); err != nil {
		return WrapError(err)
	}
	d.SetId(dataResourceIdHash(ids))
	return nil
}
