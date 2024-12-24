package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDcdnWafRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDcdnWafRulesRead,
		Schema: map[string]*schema.Schema{
			"query_args": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  20,
			},
			"waf_rules": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"action": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"cc_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"cn_region_list": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"conditions": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"op_value": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"sub_key": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"values": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"defense_scene": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"effect": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"gmt_modified": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"other_region_list": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"policy_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"rate_limit": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interval": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"status": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"code": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"count": {
													Computed: true,
													Type:     schema.TypeInt,
												},
												"ratio": {
													Computed: true,
													Type:     schema.TypeInt,
												},
											},
										},
									},
									"sub_key": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"target": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"threshold": {
										Computed: true,
										Type:     schema.TypeInt,
									},
									"ttl": {
										Computed: true,
										Type:     schema.TypeInt,
									},
								},
							},
						},
						"regular_rules": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"regular_types": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"remote_addr": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"rule_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"scenes": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"waf_group_ids": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"waf_rule_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDcdnWafRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("query_args"); ok {
		request["QueryArgs"] = v
	}
	setPagingRequest(d, request, PageSizeLarge)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeDcdnWafRules"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("dcdn", "2018-01-15", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dcdn_waf_rules", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Rules", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Rules", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RuleId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})

		mapping := map[string]interface{}{
			"id": fmt.Sprint(object["RuleId"]),
		}

		mapping["defense_scene"] = object["DefenseScene"]

		mapping["gmt_modified"] = object["GmtModified"]

		mapping["policy_id"] = fmt.Sprint(object["PolicyId"])

		mapping["rule_name"] = object["RuleName"]

		mapping["status"] = object["RuleStatus"]

		mapping["waf_rule_id"] = object["RuleId"]

		ruleConfig, err := jsonpath.Get("$.RuleConfig", object)
		if err != nil {
			return WrapError(err)
		}
		ruleConfigObj, err := convertJsonStringToMap(ruleConfig.(string))
		if err != nil {
			return WrapError(err)
		}

		mapping["action"] = ruleConfigObj["action"]
		mapping["cc_status"] = ruleConfigObj["ccStatus"]
		mapping["cn_region_list"] = ruleConfigObj["cnRegionList"]
		if conditions52Raw, ok := ruleConfigObj["conditions"]; ok {
			conditions52Maps := make([]map[string]interface{}, 0)
			for _, value0 := range conditions52Raw.([]interface{}) {
				conditions52 := value0.(map[string]interface{})
				conditions52Map := make(map[string]interface{})
				conditions52Map["key"] = conditions52["key"]
				conditions52Map["op_value"] = conditions52["opValue"]
				conditions52Map["sub_key"] = conditions52["subKey"]
				conditions52Map["values"] = conditions52["values"]
				conditions52Maps = append(conditions52Maps, conditions52Map)
			}
			mapping["conditions"] = conditions52Maps
		}

		if v, ok := ruleConfigObj["rateLimit"]; ok {
			rateLimitMap := make(map[string]interface{}, 0)
			rateLimitObj := v.(map[string]interface{})
			rateLimitMap["ttl"] = rateLimitObj["ttl"]
			rateLimitMap["sub_key"] = rateLimitObj["subKey"]
			rateLimitMap["target"] = rateLimitObj["target"]
			rateLimitMap["interval"] = rateLimitObj["interval"]
			rateLimitMap["threshold"] = rateLimitObj["threshold"]
			if v, ok := rateLimitObj["status"]; ok {
				statusMap := make(map[string]interface{}, 0)
				statusObj := v.(map[string]interface{})
				statusMap["code"] = statusObj["code"]
				statusMap["ratio"] = formatInt(statusObj["ratio"])
				statusMap["count"] = formatInt(statusObj["count"])
				rateLimitMap["status"] = []interface{}{statusMap}
			}
			mapping["rate_limit"] = []interface{}{rateLimitMap}
		}

		mapping["effect"] = ruleConfigObj["effect"]
		mapping["other_region_list"] = ruleConfigObj["otherRegionList"]
		mapping["regular_rules"] = ruleConfigObj["regularRules"]
		mapping["regular_types"] = ruleConfigObj["regularTypes"]
		mapping["remote_addr"] = ruleConfigObj["remoteAddr"]
		mapping["scenes"] = ruleConfigObj["tags"]
		mapping["waf_group_ids"] = ruleConfigObj["wafGroupIds"]

		ids = append(ids, fmt.Sprint(object["RuleId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("waf_rules", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
