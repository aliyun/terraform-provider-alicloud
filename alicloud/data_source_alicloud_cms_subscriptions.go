package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCmsSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsSubscriptionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"subscription_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notify_strategy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_setting": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"op": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"expression": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"relation": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"pushing_setting": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alert_action_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"restore_action_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"template_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"response_plan_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCmsSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/subscriptions"
	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]*string)
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("subscription_name"); ok {
		query["subscriptionName"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("enable"); ok {
		query["enable"] = StringPointer(fmt.Sprint(v.(bool)))
	}

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

	nextToken := ""
	for {
		if nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		}
		query["maxResults"] = StringPointer(fmt.Sprint(PageSizeLarge))
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, query)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_subscriptions", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.subscriptionList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.subscriptionList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["subscriptionName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["subscriptionId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		tokenRaw, _ := jsonpath.Get("$.nextToken", response)
		nextToken = fmt.Sprint(tokenRaw)
		if nextToken == "" || nextToken == "<nil>" || len(result) == 0 {
			break
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"subscription_id":    fmt.Sprint(object["subscriptionId"]),
			"subscription_name":  object["subscriptionName"],
			"description":        object["description"],
			"notify_strategy_id": object["notifyStrategyId"],
			"workspace":          object["workspace"],
			"user_id":            object["userId"],
			"create_time":        object["createTime"],
			"update_time":        object["updateTime"],
			"enable":             object["enable"],
			"region_id":          object["regionId"],
			"subscription_type":  object["subscriptionType"],
		}
		mapping["filter_setting"] = buildCmsSubscriptionFilterSettingOutput(object["filterSetting"])
		mapping["pushing_setting"] = buildCmsSubscriptionPushingSettingOutput(object["pushingSetting"])
		ids = append(ids, fmt.Sprint(mapping["subscription_id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
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

func buildCmsSubscriptionFilterSettingOutput(raw interface{}) []map[string]interface{} {
	filterResult := make([]map[string]interface{}, 0)
	if raw == nil {
		return filterResult
	}
	filterMap, ok := raw.(map[string]interface{})
	if !ok {
		return filterResult
	}
	fs := map[string]interface{}{}
	if conditions, ok := filterMap["conditions"].([]interface{}); ok {
		condList := make([]map[string]interface{}, 0)
		for _, c := range conditions {
			if cm, ok := c.(map[string]interface{}); ok {
				condList = append(condList, map[string]interface{}{
					"field": cm["field"],
					"value": cm["value"],
					"op":    cm["op"],
				})
			}
		}
		fs["conditions"] = condList
	}
	fs["expression"] = filterMap["expression"]
	fs["relation"] = filterMap["relation"]
	return append(filterResult, fs)
}

func buildCmsSubscriptionPushingSettingOutput(raw interface{}) []map[string]interface{} {
	pushingResult := make([]map[string]interface{}, 0)
	if raw == nil {
		return pushingResult
	}
	pushingMap, ok := raw.(map[string]interface{})
	if !ok {
		return pushingResult
	}
	ps := map[string]interface{}{}
	if alertActionIds, ok := pushingMap["alertActionIds"].([]interface{}); ok {
		ids := make([]string, 0)
		for _, v := range alertActionIds {
			ids = append(ids, fmt.Sprint(v))
		}
		ps["alert_action_ids"] = ids
	}
	if restoreActionIds, ok := pushingMap["restoreActionIds"].([]interface{}); ok {
		ids := make([]string, 0)
		for _, v := range restoreActionIds {
			ids = append(ids, fmt.Sprint(v))
		}
		ps["restore_action_ids"] = ids
	}
	ps["template_uuid"] = pushingMap["templateUuid"]
	ps["response_plan_id"] = pushingMap["responsePlanId"]
	return append(pushingResult, ps)
}
