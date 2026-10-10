package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBGatewayModelServices() *schema.Resource {
	return &schema.Resource{Read: dataSourceAlicloudPolarDBGatewayModelServicesRead, Schema: map[string]*schema.Schema{
		"gateway_id":     {Type: schema.TypeString, Required: true},
		"ids":            {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"name":           {Type: schema.TypeString, Optional: true},
		"model_category": {Type: schema.TypeString, Optional: true, ValidateFunc: StringInSlice([]string{"text", "embedding", "rerank"}, false)},
		"protocol":       {Type: schema.TypeString, Optional: true, ValidateFunc: StringInSlice([]string{"openai", "anthropic", "bailian", "vllm"}, false)},
		"status":         {Type: schema.TypeString, Optional: true},
		"services": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id": {Type: schema.TypeString, Computed: true}, "name": {Type: schema.TypeString, Computed: true},
			"model_category": {Type: schema.TypeString, Computed: true}, "protocol": {Type: schema.TypeString, Computed: true},
			"status": {Type: schema.TypeString, Computed: true}, "base_url": {Type: schema.TypeString, Computed: true},
			"vendor": {Type: schema.TypeString, Computed: true}, "create_time": {Type: schema.TypeString, Computed: true},
			"input_cost_points_per_million":  {Type: schema.TypeString, Computed: true},
			"output_cost_points_per_million": {Type: schema.TypeString, Computed: true},
			"request_cost_points":            {Type: schema.TypeString, Computed: true},
		}}},
	}}
}

func dataSourceAlicloudPolarDBGatewayModelServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{"GwClusterId": d.Get("gateway_id")}
	copyPolarDBGatewayAIFilters(d, request, map[string]string{"name": "Name", "model_category": "ModelCategory", "protocol": "Protocol", "status": "Status"})
	if ids := stringListFromSchema(d.Get("ids")); len(ids) > 0 {
		request["ModelServiceIds"] = strings.Join(ids, ",")
	}
	objects, err := listPolarDBGatewayAI(client, "DescribeModelServices", request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_gateway_model_services", "DescribeModelServices", AlibabaCloudSdkGoERROR)
	}
	items := make([]map[string]interface{}, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, object := range objects {
		id := fmt.Sprint(object["ModelServiceId"])
		ids = append(ids, id)
		items = append(items, map[string]interface{}{
			"id": id, "name": object["Name"], "model_category": object["ModelCategory"], "protocol": object["Protocol"],
			"status": object["Status"], "base_url": object["BaseUrl"], "vendor": object["Vendor"], "create_time": object["GmtCreated"],
			"input_cost_points_per_million": object["InputCostPointsPerMillion"], "output_cost_points_per_million": object["OutputCostPointsPerMillion"],
			"request_cost_points": object["RequestCostPoints"],
		})
	}
	if err = d.Set("services", items); err != nil {
		return WrapError(err)
	}
	d.SetId(dataResourceIdHash(ids))
	return nil
}
