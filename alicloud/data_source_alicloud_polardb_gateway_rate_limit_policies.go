package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBGatewayRateLimitPolicies() *schema.Resource {
	return &schema.Resource{Read: dataSourceAlicloudPolarDBGatewayRateLimitPoliciesRead, Schema: map[string]*schema.Schema{
		"gateway_id":   {Type: schema.TypeString, Required: true},
		"ids":          {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"scope_type":   {Type: schema.TypeString, Optional: true, ValidateFunc: StringInSlice([]string{"ConsumerGroup", "Consumer"}, false)},
		"scope_ref_id": {Type: schema.TypeString, Optional: true},
		"policies": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id": {Type: schema.TypeString, Computed: true}, "policy_type": {Type: schema.TypeString, Computed: true},
			"scope_type": {Type: schema.TypeString, Computed: true}, "scope_ref_id": {Type: schema.TypeString, Computed: true},
			"rate_limit_rpm": {Type: schema.TypeString, Computed: true}, "rate_limit_tpm": {Type: schema.TypeString, Computed: true},
			"status": {Type: schema.TypeString, Computed: true}, "create_time": {Type: schema.TypeString, Computed: true},
			"modify_time": {Type: schema.TypeString, Computed: true},
		}}},
	}}
}

func dataSourceAlicloudPolarDBGatewayRateLimitPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{"GwClusterId": d.Get("gateway_id")}
	copyPolarDBGatewayAIFilters(d, request, map[string]string{"scope_type": "ScopeType", "scope_ref_id": "ScopeRefId"})
	if ids := stringListFromSchema(d.Get("ids")); len(ids) == 1 {
		request["PolicyId"] = ids[0]
	}
	objects, err := listPolarDBGatewayAI(client, "DescribeRateLimitPolicy", request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_gateway_rate_limit_policies", "DescribeRateLimitPolicy", AlibabaCloudSdkGoERROR)
	}
	requestedIDs := make(map[string]struct{})
	for _, id := range stringListFromSchema(d.Get("ids")) {
		requestedIDs[id] = struct{}{}
	}
	items := make([]map[string]interface{}, 0, len(objects))
	ids := make([]string, 0, len(objects))
	for _, object := range objects {
		id := strings.TrimSpace(fmt.Sprint(object["PolicyId"]))
		if len(requestedIDs) > 0 {
			if _, ok := requestedIDs[id]; !ok {
				continue
			}
		}
		ids = append(ids, id)
		items = append(items, map[string]interface{}{
			"id": id, "policy_type": object["PolicyType"], "scope_type": object["ScopeType"],
			"scope_ref_id": object["ScopeRefId"], "rate_limit_rpm": object["RateLimitRpm"],
			"rate_limit_tpm": object["RateLimitTpm"], "status": object["Status"],
			"create_time": object["GmtCreated"], "modify_time": object["GmtModified"],
		})
	}
	if err = d.Set("policies", items); err != nil {
		return WrapError(err)
	}
	d.SetId(dataResourceIdHash(ids))
	return nil
}
