package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudCmsSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsSubscriptionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"agent_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"routes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"channels": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"channel_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"enabled_sub_channels": {
																Type:     schema.TypeList,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
															"receivers": {
																Type:     schema.TypeList,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
												"effect_time_range": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day_in_week": {
																Type:     schema.TypeList,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeInt},
															},
															"end_time_in_minute": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"start_time_in_minute": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"time_zone": {
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
							},
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
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
												"op": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
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
						"notify_strategy_id": {
							Type:     schema.TypeString,
							Computed: true,
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
									"response_plan_id": {
										Type:     schema.TypeString,
										Computed: true,
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
								},
							},
						},
						"subscription_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
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
		},
	}
}

func dataSourceAliCloudCmsSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
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
	var query map[string]*string
	// ListSubscriptions
	action := fmt.Sprintf("/subscriptions")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["RegionId"] = StringPointer(client.RegionId)
	query["maxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)

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

		resp, _ := jsonpath.Get("$.dataList[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil {
				if !nameRegex.MatchString(fmt.Sprint(item["subscriptionName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["subscriptionId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["subscriptionId"]

		mapping["create_time"] = objectRaw["createTime"]
		mapping["description"] = objectRaw["description"]
		mapping["enable"] = objectRaw["enabled"]
		mapping["notify_strategy_id"] = objectRaw["notifyStrategyId"]
		mapping["subscription_id"] = objectRaw["subscriptionId"]
		mapping["subscription_name"] = objectRaw["subscriptionName"]
		mapping["subscription_type"] = objectRaw["subscriptionType"]
		mapping["update_time"] = objectRaw["updateTime"]
		mapping["user_id"] = objectRaw["userId"]
		mapping["workspace"] = objectRaw["workspace"]

		filterSettingMaps := make([]map[string]interface{}, 0)
		filterSettingMap := make(map[string]interface{})
		filterSettingRaw := make(map[string]interface{})
		if objectRaw["filterSetting"] != nil {
			filterSettingRaw, _ = objectRaw["filterSetting"].(map[string]interface{})
		}
		if len(filterSettingRaw) > 0 {
			filterSettingMap["expression"] = filterSettingRaw["expression"]
			filterSettingMap["relation"] = filterSettingRaw["relation"]

			conditionsRaw := filterSettingRaw["conditions"]
			conditionsMaps := make([]map[string]interface{}, 0)
			if conditionsRaw != nil {
				for _, conditionsChildRaw := range convertToInterfaceArray(conditionsRaw) {
					conditionsMap := make(map[string]interface{})
					conditionsChildRaw := conditionsChildRaw.(map[string]interface{})
					conditionsMap["field"] = conditionsChildRaw["field"]
					conditionsMap["op"] = conditionsChildRaw["op"]
					conditionsMap["value"] = conditionsChildRaw["value"]

					conditionsMaps = append(conditionsMaps, conditionsMap)
				}
			}
			filterSettingMap["conditions"] = conditionsMaps
			filterSettingMaps = append(filterSettingMaps, filterSettingMap)
		}
		mapping["filter_setting"] = filterSettingMaps
		pushingSettingMaps := make([]map[string]interface{}, 0)
		pushingSettingMap := make(map[string]interface{})
		pushingSettingRaw := make(map[string]interface{})
		if objectRaw["pushingSetting"] != nil {
			pushingSettingRaw, _ = objectRaw["pushingSetting"].(map[string]interface{})
		}
		if len(pushingSettingRaw) > 0 {
			pushingSettingMap["response_plan_id"] = pushingSettingRaw["responsePlanId"]
			pushingSettingMap["template_uuid"] = pushingSettingRaw["templateUuid"]

			alertActionIdsRaw := make([]interface{}, 0)
			if pushingSettingRaw["alertActionIds"] != nil {
				alertActionIdsRaw = convertToInterfaceArray(pushingSettingRaw["alertActionIds"])
			}
			pushingSettingMap["alert_action_ids"] = alertActionIdsRaw

			restoreActionIdsRaw := make([]interface{}, 0)
			if pushingSettingRaw["restoreActionIds"] != nil {
				restoreActionIdsRaw = convertToInterfaceArray(pushingSettingRaw["restoreActionIds"])
			}
			pushingSettingMap["restore_action_ids"] = restoreActionIdsRaw

			pushingSettingMaps = append(pushingSettingMaps, pushingSettingMap)
		}
		mapping["pushing_setting"] = pushingSettingMaps
		agentConfigMaps := make([]map[string]interface{}, 0)
		agentConfigMap := make(map[string]interface{})
		agentConfigRaw := make(map[string]interface{})
		if objectRaw["agentConfig"] != nil {
			agentConfigRaw, _ = objectRaw["agentConfig"].(map[string]interface{})
		}
		if len(agentConfigRaw) > 0 {
			agentConfigMap["agent_uuid"] = agentConfigRaw["agentUuid"]

			routesRaw := agentConfigRaw["routes"]
			routesMaps := make([]map[string]interface{}, 0)
			if routesRaw != nil {
				for _, routesChildRaw := range convertToInterfaceArray(routesRaw) {
					routesMap := make(map[string]interface{})
					routesChildRaw := routesChildRaw.(map[string]interface{})

					channelsRaw := routesChildRaw["channels"]
					channelsMaps := make([]map[string]interface{}, 0)
					if channelsRaw != nil {
						for _, channelsChildRaw := range convertToInterfaceArray(channelsRaw) {
							channelsMap := make(map[string]interface{})
							channelsChildRaw := channelsChildRaw.(map[string]interface{})
							channelsMap["channel_type"] = channelsChildRaw["channelType"]

							enabledSubChannelsRaw := make([]interface{}, 0)
							if channelsChildRaw["enabledSubChannels"] != nil {
								enabledSubChannelsRaw = convertToInterfaceArray(channelsChildRaw["enabledSubChannels"])
							}
							channelsMap["enabled_sub_channels"] = enabledSubChannelsRaw

							receiversRaw := make([]interface{}, 0)
							if channelsChildRaw["receivers"] != nil {
								receiversRaw = convertToInterfaceArray(channelsChildRaw["receivers"])
							}
							channelsMap["receivers"] = receiversRaw

							channelsMaps = append(channelsMaps, channelsMap)
						}
					}
					routesMap["channels"] = channelsMaps

					effectTimeRangeMaps := make([]map[string]interface{}, 0)
					effectTimeRangeMap := make(map[string]interface{})
					effectTimeRangeRaw := make(map[string]interface{})
					if routesChildRaw["effectTimeRange"] != nil {
						effectTimeRangeRaw, _ = routesChildRaw["effectTimeRange"].(map[string]interface{})
					}
					if len(effectTimeRangeRaw) > 0 {
						effectTimeRangeMap["end_time_in_minute"] = effectTimeRangeRaw["endTimeInMinute"]
						effectTimeRangeMap["start_time_in_minute"] = effectTimeRangeRaw["startTimeInMinute"]
						effectTimeRangeMap["time_zone"] = effectTimeRangeRaw["timeZone"]

						dayInWeekRaw := make([]interface{}, 0)
						if effectTimeRangeRaw["dayInWeek"] != nil {
							dayInWeekRaw = convertToInterfaceArray(effectTimeRangeRaw["dayInWeek"])
						}
						effectTimeRangeMap["day_in_week"] = dayInWeekRaw

						effectTimeRangeMaps = append(effectTimeRangeMaps, effectTimeRangeMap)
					}
					routesMap["effect_time_range"] = effectTimeRangeMaps

					routesMaps = append(routesMaps, routesMap)
				}
			}
			agentConfigMap["routes"] = routesMaps

			agentConfigMaps = append(agentConfigMaps, agentConfigMap)
		}
		mapping["agent_config"] = agentConfigMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["subscriptionName"])
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
