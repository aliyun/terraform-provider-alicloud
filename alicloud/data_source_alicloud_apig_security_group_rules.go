// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudApigSecurityGroupRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudApigSecurityGroupRulesRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"security_group_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_cidrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudApigSecurityGroupRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayId := d.Get("gateway_id").(string)
	action := fmt.Sprintf("/v1/gateways/%s/authorized-security-groups-rules", gatewayId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var filterRuleId string
	if v, ok := d.GetOk("security_group_rule_id"); ok {
		filterRuleId = v.(string)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	itemsRaw, err := jsonpath.Get("$.data.items", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.data.items", response)
	}
	items, ok := itemsRaw.([]interface{})
	if !ok {
		items = make([]interface{}, 0)
	}

	objects := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		ruleId := fmt.Sprint(item["securityGroupRuleId"])
		if filterRuleId != "" && ruleId != filterRuleId {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[ruleId]; !ok {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":                       fmt.Sprintf("%v:%v", gatewayId, ruleId),
			"gateway_id":               gatewayId,
			"security_group_rule_id":   ruleId,
			"security_group_id":        item["securityGroupId"],
			"port_range":               item["portRange"],
			"ip_protocol":              item["ipProtocol"],
			"description":              item["description"],
			"source_security_group_id": item["sourceSecurityGroupId"],
			"security_group_name":      item["securityGroupName"],
			"vpc_id":                   item["vpcId"],
			"auth_cidrs":               item["authCidrs"],
		}
		objects = append(objects, mapping)
		ids = append(ids, ruleId)
	}

	d.SetId(fmt.Sprintf("%v", gatewayId))
	d.Set("rules", objects)
	d.Set("ids", ids)

	return nil
}
