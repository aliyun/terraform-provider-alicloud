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

func dataSourceAliCloudCmsNotifyStrategies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsNotifyStrategiesRead,
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
			"notify_strategies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_template_entries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"template_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"grouping_setting": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"grouping_keys": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"period_min": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"silence_sec": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"times": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"ignore_restored_notification": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"notify_strategy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notify_strategy_name": {
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
									"severities": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
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

func dataSourceAliCloudCmsNotifyStrategiesRead(d *schema.ResourceData, meta interface{}) error {
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
	// ListNotifyStrategies
	action := fmt.Sprintf("/notifyStrategies")
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
				if !nameRegex.MatchString(fmt.Sprint(item["notifyStrategyName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["notifyStrategyId"])]; !ok {
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

		mapping["id"] = objectRaw["notifyStrategyId"]

		mapping["create_time"] = objectRaw["createTime"]
		mapping["description"] = objectRaw["description"]
		mapping["enable"] = objectRaw["enabled"]
		mapping["ignore_restored_notification"] = objectRaw["ignoreRestoredNotification"]
		mapping["notify_strategy_id"] = objectRaw["notifyStrategyId"]
		mapping["notify_strategy_name"] = objectRaw["notifyStrategyName"]
		mapping["update_time"] = objectRaw["updateTime"]
		mapping["user_id"] = objectRaw["userId"]
		mapping["workspace"] = objectRaw["workspace"]

		groupingSettingMaps := make([]map[string]interface{}, 0)
		groupingSettingMap := make(map[string]interface{})
		groupingSettingRaw := make(map[string]interface{})
		if objectRaw["groupingSetting"] != nil {
			groupingSettingRaw, _ = objectRaw["groupingSetting"].(map[string]interface{})
		}
		if len(groupingSettingRaw) > 0 {
			groupingSettingMap["period_min"] = groupingSettingRaw["periodMin"]
			groupingSettingMap["silence_sec"] = groupingSettingRaw["silenceSec"]
			groupingSettingMap["times"] = groupingSettingRaw["times"]

			groupingKeysRaw := make([]interface{}, 0)
			if groupingSettingRaw["groupingKeys"] != nil {
				groupingKeysRaw = convertToInterfaceArray(groupingSettingRaw["groupingKeys"])
			}
			groupingSettingMap["grouping_keys"] = groupingKeysRaw

			groupingSettingMaps = append(groupingSettingMaps, groupingSettingMap)
		}
		mapping["grouping_setting"] = groupingSettingMaps
		customTemplateEntriesRaw := objectRaw["customTemplateEntries"]
		customTemplateEntriesMaps := make([]map[string]interface{}, 0)
		if customTemplateEntriesRaw != nil {
			for _, customTemplateEntriesChildRaw := range convertToInterfaceArray(customTemplateEntriesRaw) {
				customTemplateEntriesMap := make(map[string]interface{})
				customTemplateEntriesChildRaw := customTemplateEntriesChildRaw.(map[string]interface{})
				customTemplateEntriesMap["target_type"] = customTemplateEntriesChildRaw["targetType"]
				customTemplateEntriesMap["template_uuid"] = customTemplateEntriesChildRaw["templateUuid"]

				customTemplateEntriesMaps = append(customTemplateEntriesMaps, customTemplateEntriesMap)
			}
		}
		mapping["custom_template_entries"] = customTemplateEntriesMaps
		routesRaw := objectRaw["routes"]
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

				filterSettingMaps := make([]map[string]interface{}, 0)
				filterSettingMap := make(map[string]interface{})
				filterSettingRaw := make(map[string]interface{})
				if routesChildRaw["filterSetting"] != nil {
					filterSettingRaw, _ = routesChildRaw["filterSetting"].(map[string]interface{})
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
				routesMap["filter_setting"] = filterSettingMaps

				severitiesRaw := make([]interface{}, 0)
				if routesChildRaw["severities"] != nil {
					severitiesRaw = convertToInterfaceArray(routesChildRaw["severities"])
				}
				routesMap["severities"] = severitiesRaw

				routesMaps = append(routesMaps, routesMap)
			}
		}
		mapping["routes"] = routesMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["notifyStrategyName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("notify_strategies", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
