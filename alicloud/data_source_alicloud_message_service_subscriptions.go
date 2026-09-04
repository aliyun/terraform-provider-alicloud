// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudMessageServiceSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudMessageServiceSubscriptionRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"subscription_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dlq_policy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dead_letter_target_queue": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modify_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"notify_content_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notify_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_rate_limit_policy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"max_receives_per_second": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"topic_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudMessageServiceSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListSubscriptionByTopic"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["TopicName"] = d.Get("topic_name")
	if v, ok := d.GetOk("subscription_name"); ok {
		request["SubscriptionName"] = v
	}

	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNum"] = v.(int)
	} else {
		request["PageNum"] = 1
	}

	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)

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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.Data.PageData[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["SubscriptionName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TopicName"], ":", item["SubscriptionName"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(objectRaw["TopicName"], ":", objectRaw["SubscriptionName"])

		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["endpoint"] = objectRaw["Endpoint"]
		mapping["filter_tag"] = objectRaw["FilterTag"]
		mapping["last_modify_time"] = objectRaw["LastModifyTime"]
		mapping["notify_content_format"] = objectRaw["NotifyContentFormat"]
		mapping["notify_strategy"] = objectRaw["NotifyStrategy"]
		mapping["topic_owner"] = objectRaw["TopicOwner"]
		mapping["subscription_name"] = objectRaw["SubscriptionName"]
		mapping["subscription_url"] = objectRaw["SubscriptionURL"]
		mapping["topic_name"] = objectRaw["TopicName"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["SubscriptionName"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["TopicName"], ":", objectRaw["SubscriptionName"])
		mapping, err = dataSourceAliCloudMessageServiceSubscriptionReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["SubscriptionName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("subscriptions", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudMessageServiceSubscriptionReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	messageServiceServiceV2 := MessageServiceServiceV2{client}
	getResp, err := messageServiceServiceV2.DescribeMessageServiceSubscription(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["create_time"] = objectRaw["CreateTime"]
	mapping["endpoint"] = objectRaw["Endpoint"]
	mapping["filter_tag"] = objectRaw["FilterTag"]
	mapping["last_modify_time"] = objectRaw["LastModifyTime"]
	mapping["notify_content_format"] = objectRaw["NotifyContentFormat"]
	mapping["notify_strategy"] = objectRaw["NotifyStrategy"]
	mapping["topic_owner"] = objectRaw["TopicOwner"]
	mapping["subscription_name"] = objectRaw["SubscriptionName"]
	mapping["topic_name"] = objectRaw["TopicName"]

	dlqPolicyMaps := make([]map[string]interface{}, 0)
	dlqPolicyMap := make(map[string]interface{})
	dlqPolicyRaw := make(map[string]interface{})
	if objectRaw["DlqPolicy"] != nil {
		dlqPolicyRaw = objectRaw["DlqPolicy"].(map[string]interface{})
	}
	if len(dlqPolicyRaw) > 0 {
		dlqPolicyMap["dead_letter_target_queue"] = dlqPolicyRaw["DeadLetterTargetQueue"]
		dlqPolicyMap["enabled"] = dlqPolicyRaw["Enabled"]

		dlqPolicyMaps = append(dlqPolicyMaps, dlqPolicyMap)
	}
	mapping["dlq_policy"] = dlqPolicyMaps
	tenantRateLimitPolicyMaps := make([]map[string]interface{}, 0)
	tenantRateLimitPolicyMap := make(map[string]interface{})
	tenantRateLimitPolicyRaw := make(map[string]interface{})
	if objectRaw["TenantRateLimitPolicy"] != nil {
		tenantRateLimitPolicyRaw = objectRaw["TenantRateLimitPolicy"].(map[string]interface{})
	}
	if len(tenantRateLimitPolicyRaw) > 0 {
		tenantRateLimitPolicyMap["enabled"] = tenantRateLimitPolicyRaw["Enabled"]
		tenantRateLimitPolicyMap["max_receives_per_second"] = tenantRateLimitPolicyRaw["MaxReceivesPerSecond"]

		tenantRateLimitPolicyMaps = append(tenantRateLimitPolicyMaps, tenantRateLimitPolicyMap)
	}
	mapping["tenant_rate_limit_policy"] = tenantRateLimitPolicyMaps

	return mapping, nil
}
