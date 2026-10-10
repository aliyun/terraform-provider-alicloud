package alicloud

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBGatewayCostRules() *schema.Resource {
	return &schema.Resource{Read: dataSourceAlicloudPolarDBGatewayCostRulesRead, Schema: map[string]*schema.Schema{
		"gateway_id":       {Type: schema.TypeString, Required: true},
		"model_name":       {Type: schema.TypeString, Optional: true},
		"model_service_id": {Type: schema.TypeString, Optional: true},
		"rules": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id": {Type: schema.TypeString, Computed: true}, "model_name": {Type: schema.TypeString, Computed: true},
			"model_service_id": {Type: schema.TypeString, Computed: true}, "input_cost_points_per_million": {Type: schema.TypeString, Computed: true},
			"output_cost_points_per_million": {Type: schema.TypeString, Computed: true}, "cache_cost_points_per_million": {Type: schema.TypeString, Computed: true},
			"create_time": {Type: schema.TypeString, Computed: true}, "modify_time": {Type: schema.TypeString, Computed: true},
		}}},
	}}
}

func dataSourceAlicloudPolarDBGatewayCostRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{"GwClusterId": d.Get("gateway_id")}
	copyPolarDBGatewayAIFilters(d, request, map[string]string{"model_name": "ModelName", "model_service_id": "ModelServiceId"})
	objects, err := listPolarDBGatewayAI(client, "DescribeCostRules", request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_gateway_cost_rules", "DescribeCostRules", AlibabaCloudSdkGoERROR)
	}
	items := make([]map[string]interface{}, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, object := range objects {
		id := fmt.Sprint(object["CostRuleId"])
		modelName := object["Model"]
		if modelName == nil {
			modelName = object["ModelName"]
		}
		ids = append(ids, id)
		items = append(items, map[string]interface{}{
			"id": id, "model_name": modelName, "model_service_id": object["ModelServiceId"],
			"input_cost_points_per_million": object["InputCostPointsPerMillion"], "output_cost_points_per_million": object["OutputCostPointsPerMillion"],
			"cache_cost_points_per_million": object["CacheCostPointsPerMillion"], "create_time": object["GmtCreated"], "modify_time": object["GmtModified"],
		})
	}
	if err = d.Set("rules", items); err != nil {
		return WrapError(err)
	}
	d.SetId(dataResourceIdHash(ids))
	return nil
}
