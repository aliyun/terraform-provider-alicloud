package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBGatewayConsumers() *schema.Resource {
	return &schema.Resource{Read: dataSourceAlicloudPolarDBGatewayConsumersRead, Schema: map[string]*schema.Schema{
		"gateway_id":        {Type: schema.TypeString, Required: true},
		"ids":               {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"consumer_group_id": {Type: schema.TypeString, Optional: true},
		"consumers": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id": {Type: schema.TypeString, Computed: true}, "name": {Type: schema.TypeString, Computed: true},
			"consumer_group_id": {Type: schema.TypeString, Computed: true}, "consumer_group_name": {Type: schema.TypeString, Computed: true},
			"nickname": {Type: schema.TypeString, Computed: true}, "allowed_models": {Type: schema.TypeString, Computed: true},
			"month_to_date_cost_count": {Type: schema.TypeInt, Computed: true}, "lifetime_cost_count": {Type: schema.TypeInt, Computed: true},
			"month_to_date_token_count": {Type: schema.TypeInt, Computed: true}, "lifetime_token_count": {Type: schema.TypeInt, Computed: true},
			"create_time": {Type: schema.TypeString, Computed: true}, "modify_time": {Type: schema.TypeString, Computed: true},
		}}},
	}}
}

func dataSourceAlicloudPolarDBGatewayConsumersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{"GwClusterId": d.Get("gateway_id")}
	if value, ok := d.GetOk("consumer_group_id"); ok {
		request["ConsumerGroupId"] = value
	}
	if ids := stringListFromSchema(d.Get("ids")); len(ids) == 1 {
		request["ConsumerId"] = ids[0]
	}
	objects, err := listPolarDBGatewayAI(client, "DescribeConsumers", request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_gateway_consumers", "DescribeConsumers", AlibabaCloudSdkGoERROR)
	}
	requestedIDs := make(map[string]struct{})
	for _, id := range stringListFromSchema(d.Get("ids")) {
		requestedIDs[id] = struct{}{}
	}
	items := make([]map[string]interface{}, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, object := range objects {
		id := strings.TrimSpace(fmt.Sprint(object["ConsumerId"]))
		if len(requestedIDs) > 0 {
			if _, ok := requestedIDs[id]; !ok {
				continue
			}
		}
		ids = append(ids, id)
		items = append(items, map[string]interface{}{
			"id": id, "name": object["Name"], "consumer_group_id": object["ConsumerGroupId"],
			"consumer_group_name": object["ConsumerGroupName"], "nickname": object["NickName"],
			"allowed_models": object["AllowedModels"], "month_to_date_cost_count": object["MtdCostCount"],
			"lifetime_cost_count": object["LifetimeCostCount"], "month_to_date_token_count": object["MtdTokenCount"],
			"lifetime_token_count": object["LifetimeTokenCount"], "create_time": object["GmtCreated"],
			"modify_time": object["GmtModified"],
		})
	}
	if err = d.Set("consumers", items); err != nil {
		return WrapError(err)
	}
	d.SetId(dataResourceIdHash(ids))
	return nil
}
