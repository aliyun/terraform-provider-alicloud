package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBGatewayConsumerGroups() *schema.Resource {
	return &schema.Resource{Read: dataSourceAlicloudPolarDBGatewayConsumerGroupsRead, Schema: map[string]*schema.Schema{
		"gateway_id": {Type: schema.TypeString, Required: true},
		"ids":        {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"groups": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id": {Type: schema.TypeString, Computed: true}, "name": {Type: schema.TypeString, Computed: true},
			"nickname": {Type: schema.TypeString, Computed: true}, "is_default": {Type: schema.TypeString, Computed: true},
			"allowed_models": {Type: schema.TypeString, Computed: true}, "create_time": {Type: schema.TypeString, Computed: true},
			"modify_time": {Type: schema.TypeString, Computed: true},
		}}},
	}}
}

func dataSourceAlicloudPolarDBGatewayConsumerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{"GwClusterId": d.Get("gateway_id")}
	if ids := stringListFromSchema(d.Get("ids")); len(ids) == 1 {
		request["ConsumerGroupId"] = ids[0]
	}
	objects, err := listPolarDBGatewayAI(client, "DescribeConsumerGroups", request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_gateway_consumer_groups", "DescribeConsumerGroups", AlibabaCloudSdkGoERROR)
	}
	requestedIDs := make(map[string]struct{})
	for _, id := range stringListFromSchema(d.Get("ids")) {
		requestedIDs[id] = struct{}{}
	}
	items := make([]map[string]interface{}, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, object := range objects {
		id := strings.TrimSpace(fmt.Sprint(object["ConsumerGroupId"]))
		if len(requestedIDs) > 0 {
			if _, ok := requestedIDs[id]; !ok {
				continue
			}
		}
		ids = append(ids, id)
		items = append(items, map[string]interface{}{
			"id": id, "name": object["ConsumerGroupName"], "nickname": object["NickName"],
			"is_default": object["IsDefault"], "allowed_models": object["AllowedModels"],
			"create_time": object["GmtCreated"], "modify_time": object["GmtModified"],
		})
	}
	if err = d.Set("groups", items); err != nil {
		return WrapError(err)
	}
	d.SetId(dataResourceIdHash(ids))
	return nil
}
