// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsEventNotifyPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsEventNotifyPolicyRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notify_strategy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ignore_restored_notification": {
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
												"times": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"silence_sec": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"routes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"digital_employee_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"enable_rca": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"filter_setting": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"relation": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"expression": {
																Type:     schema.TypeString,
																Computed: true,
															},
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
														},
													},
												},
												"effect_time_range": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"time_zone": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"start_time_in_minute": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"end_time_in_minute": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"day_in_week": {
																Type:     schema.TypeList,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeInt},
															},
														},
													},
												},
												"channels": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"receivers": {
																Type:     schema.TypeList,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
															"channel_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"enabled_sub_channels": {
																Type:     schema.TypeSet,
																Computed: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
											},
										},
									},
									"custom_template_entries": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"template_uuid": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"response_plan": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repeat_notify_setting": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_incident_state": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"repeat_interval": {
													Type:     schema.TypeInt,
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
												"restore_action_ids": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"alert_action_ids": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"escalation_id": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"auto_recover_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"subscription": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workspace_filter_setting": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"workspace_uuids": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"tag_selector": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"relation": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"expression": {
																Type:     schema.TypeString,
																Computed: true,
															},
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
														},
													},
												},
											},
										},
									},
									"subscribe_legacy_event": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"filter_setting": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"relation": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"expression": {
													Type:     schema.TypeString,
													Computed: true,
												},
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
											},
										},
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
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeInt,
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudCmsEventNotifyPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

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
	// ListNotifyPolicies
	action := fmt.Sprintf("/api/eventbase/notify-policies")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["regionId"] = StringPointer(client.RegionId)
	query["workspace"] = StringPointer(d.Get("workspace").(string))
	if v, ok := d.GetOk("name"); ok {
		query["name"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		query["orderBy"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("order_desc"); ok {
		query["orderDesc"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["maxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))
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

		resp, _ := jsonpath.Get("$.notifyPolicyList[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["uuid"], ":", item["workspace"])]; !ok {
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
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(objectRaw["uuid"], ":", objectRaw["workspace"])

		mapping["create_time"] = objectRaw["createTime"]
		mapping["description"] = objectRaw["description"]
		mapping["enabled"] = objectRaw["enabled"]
		mapping["name"] = objectRaw["name"]
		mapping["update_time"] = objectRaw["updateTime"]
		mapping["user_id"] = objectRaw["userId"]
		mapping["uuid"] = objectRaw["uuid"]
		mapping["version"] = objectRaw["version"]
		mapping["workspace"] = objectRaw["workspace"]

		notifyStrategyMaps := make([]map[string]interface{}, 0)
		notifyStrategyMap := make(map[string]interface{})
		notifyStrategyRaw := make(map[string]interface{})
		if objectRaw["notifyStrategy"] != nil {
			notifyStrategyRaw = objectRaw["notifyStrategy"].(map[string]interface{})
		}
		if len(notifyStrategyRaw) > 0 {
			notifyStrategyMap["description"] = notifyStrategyRaw["description"]
			notifyStrategyMap["ignore_restored_notification"] = notifyStrategyRaw["ignoreRestoredNotification"]

			customTemplateEntriesRaw := notifyStrategyRaw["customTemplateEntries"]
			customTemplateEntriesMaps := make([]map[string]interface{}, 0)
			if customTemplateEntriesRaw != nil {
				for _, customTemplateEntriesChildRaw := range convertToInterfaceArray(customTemplateEntriesRaw) {
					customTemplateEntriesMap := make(map[string]interface{})
					customTemplateEntriesChildRaw := customTemplateEntriesChildRaw.(map[string]interface{})
					customTemplateEntriesMap["template_uuid"] = customTemplateEntriesChildRaw["templateUuid"]

					customTemplateEntriesMaps = append(customTemplateEntriesMaps, customTemplateEntriesMap)
				}
			}
			notifyStrategyMap["custom_template_entries"] = customTemplateEntriesMaps
			groupingSettingMaps := make([]map[string]interface{}, 0)
			groupingSettingMap := make(map[string]interface{})
			groupingSettingRaw := make(map[string]interface{})
			if notifyStrategyRaw["groupingSetting"] != nil {
				groupingSettingRaw = notifyStrategyRaw["groupingSetting"].(map[string]interface{})
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
			notifyStrategyMap["grouping_setting"] = groupingSettingMaps
			routesRaw := notifyStrategyRaw["routes"]
			routesMaps := make([]map[string]interface{}, 0)
			if routesRaw != nil {
				for _, routesChildRaw := range convertToInterfaceArray(routesRaw) {
					routesMap := make(map[string]interface{})
					routesChildRaw := routesChildRaw.(map[string]interface{})
					routesMap["digital_employee_name"] = routesChildRaw["digitalEmployeeName"]
					routesMap["enable_rca"] = routesChildRaw["enableRca"]

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
						effectTimeRangeRaw = routesChildRaw["effectTimeRange"].(map[string]interface{})
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
						filterSettingRaw = routesChildRaw["filterSetting"].(map[string]interface{})
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
					routesMaps = append(routesMaps, routesMap)
				}
			}
			notifyStrategyMap["routes"] = routesMaps
			notifyStrategyMaps = append(notifyStrategyMaps, notifyStrategyMap)
		}
		mapping["notify_strategy"] = notifyStrategyMaps

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["uuid"], ":", objectRaw["workspace"])
		mapping, err = dataSourceAliCloudCmsEventNotifyPolicyReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudCmsEventNotifyPolicyReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	cmsServiceV2 := CmsServiceV2{client}
	getResp, err := cmsServiceV2.DescribeCmsEventNotifyPolicy(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["create_time"] = objectRaw["createTime"]
	mapping["description"] = objectRaw["description"]
	mapping["enabled"] = objectRaw["enabled"]
	mapping["name"] = objectRaw["name"]
	mapping["update_time"] = objectRaw["updateTime"]
	mapping["user_id"] = objectRaw["userId"]
	mapping["uuid"] = objectRaw["uuid"]
	mapping["version"] = objectRaw["version"]
	mapping["workspace"] = objectRaw["workspace"]

	notifyStrategyMaps := make([]map[string]interface{}, 0)
	notifyStrategyMap := make(map[string]interface{})
	notifyStrategyRaw := make(map[string]interface{})
	if objectRaw["notifyStrategy"] != nil {
		notifyStrategyRaw = objectRaw["notifyStrategy"].(map[string]interface{})
	}
	if len(notifyStrategyRaw) > 0 {
		notifyStrategyMap["description"] = notifyStrategyRaw["description"]
		notifyStrategyMap["ignore_restored_notification"] = notifyStrategyRaw["ignoreRestoredNotification"]

		customTemplateEntriesRaw := notifyStrategyRaw["customTemplateEntries"]
		customTemplateEntriesMaps := make([]map[string]interface{}, 0)
		if customTemplateEntriesRaw != nil {
			for _, customTemplateEntriesChildRaw := range convertToInterfaceArray(customTemplateEntriesRaw) {
				customTemplateEntriesMap := make(map[string]interface{})
				customTemplateEntriesChildRaw := customTemplateEntriesChildRaw.(map[string]interface{})
				customTemplateEntriesMap["template_uuid"] = customTemplateEntriesChildRaw["templateUuid"]

				customTemplateEntriesMaps = append(customTemplateEntriesMaps, customTemplateEntriesMap)
			}
		}
		notifyStrategyMap["custom_template_entries"] = customTemplateEntriesMaps
		groupingSettingMaps := make([]map[string]interface{}, 0)
		groupingSettingMap := make(map[string]interface{})
		groupingSettingRaw := make(map[string]interface{})
		if notifyStrategyRaw["groupingSetting"] != nil {
			groupingSettingRaw = notifyStrategyRaw["groupingSetting"].(map[string]interface{})
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
		notifyStrategyMap["grouping_setting"] = groupingSettingMaps
		routesRaw := notifyStrategyRaw["routes"]
		routesMaps := make([]map[string]interface{}, 0)
		if routesRaw != nil {
			for _, routesChildRaw := range convertToInterfaceArray(routesRaw) {
				routesMap := make(map[string]interface{})
				routesChildRaw := routesChildRaw.(map[string]interface{})
				routesMap["digital_employee_name"] = routesChildRaw["digitalEmployeeName"]
				routesMap["enable_rca"] = routesChildRaw["enableRca"]

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
					effectTimeRangeRaw = routesChildRaw["effectTimeRange"].(map[string]interface{})
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
					filterSettingRaw = routesChildRaw["filterSetting"].(map[string]interface{})
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
				routesMaps = append(routesMaps, routesMap)
			}
		}
		notifyStrategyMap["routes"] = routesMaps
		notifyStrategyMaps = append(notifyStrategyMaps, notifyStrategyMap)
	}
	mapping["notify_strategy"] = notifyStrategyMaps
	responsePlanMaps := make([]map[string]interface{}, 0)
	responsePlanMap := make(map[string]interface{})
	responsePlanRaw := make(map[string]interface{})
	if objectRaw["responsePlan"] != nil {
		responsePlanRaw = objectRaw["responsePlan"].(map[string]interface{})
	}
	if len(responsePlanRaw) > 0 {
		responsePlanMap["auto_recover_seconds"] = responsePlanRaw["autoRecoverSeconds"]

		escalationIdRaw := make([]interface{}, 0)
		if responsePlanRaw["escalationId"] != nil {
			escalationIdRaw = convertToInterfaceArray(responsePlanRaw["escalationId"])
		}

		responsePlanMap["escalation_id"] = escalationIdRaw
		pushingSettingMaps := make([]map[string]interface{}, 0)
		pushingSettingMap := make(map[string]interface{})
		pushingSettingRaw := make(map[string]interface{})
		if responsePlanRaw["pushingSetting"] != nil {
			pushingSettingRaw = responsePlanRaw["pushingSetting"].(map[string]interface{})
		}
		if len(pushingSettingRaw) > 0 {

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
		responsePlanMap["pushing_setting"] = pushingSettingMaps
		repeatNotifySettingMaps := make([]map[string]interface{}, 0)
		repeatNotifySettingMap := make(map[string]interface{})
		repeatNotifySettingRaw := make(map[string]interface{})
		if responsePlanRaw["repeatNotifySetting"] != nil {
			repeatNotifySettingRaw = responsePlanRaw["repeatNotifySetting"].(map[string]interface{})
		}
		if len(repeatNotifySettingRaw) > 0 {
			repeatNotifySettingMap["end_incident_state"] = repeatNotifySettingRaw["endIncidentState"]
			repeatNotifySettingMap["repeat_interval"] = repeatNotifySettingRaw["repeatInterval"]

			repeatNotifySettingMaps = append(repeatNotifySettingMaps, repeatNotifySettingMap)
		}
		responsePlanMap["repeat_notify_setting"] = repeatNotifySettingMaps
		responsePlanMaps = append(responsePlanMaps, responsePlanMap)
	}
	mapping["response_plan"] = responsePlanMaps
	subscriptionMaps := make([]map[string]interface{}, 0)
	subscriptionMap := make(map[string]interface{})
	subscriptionRaw := make(map[string]interface{})
	if objectRaw["subscription"] != nil {
		subscriptionRaw = objectRaw["subscription"].(map[string]interface{})
	}
	if len(subscriptionRaw) > 0 {
		subscriptionMap["subscribe_legacy_event"] = subscriptionRaw["subscribeLegacyEvent"]

		filterSettingMaps := make([]map[string]interface{}, 0)
		filterSettingMap := make(map[string]interface{})
		filterSettingRaw := make(map[string]interface{})
		if subscriptionRaw["filterSetting"] != nil {
			filterSettingRaw = subscriptionRaw["filterSetting"].(map[string]interface{})
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
		subscriptionMap["filter_setting"] = filterSettingMaps
		workspaceFilterSettingMaps := make([]map[string]interface{}, 0)
		workspaceFilterSettingMap := make(map[string]interface{})
		workspaceFilterSettingRaw := make(map[string]interface{})
		if subscriptionRaw["workspaceFilterSetting"] != nil {
			workspaceFilterSettingRaw = subscriptionRaw["workspaceFilterSetting"].(map[string]interface{})
		}
		if len(workspaceFilterSettingRaw) > 0 {

			tagSelectorMaps := make([]map[string]interface{}, 0)
			tagSelectorMap := make(map[string]interface{})
			tagSelectorRaw := make(map[string]interface{})
			if workspaceFilterSettingRaw["tagSelector"] != nil {
				tagSelectorRaw = workspaceFilterSettingRaw["tagSelector"].(map[string]interface{})
			}
			if len(tagSelectorRaw) > 0 {
				tagSelectorMap["expression"] = tagSelectorRaw["expression"]
				tagSelectorMap["relation"] = tagSelectorRaw["relation"]

				conditionsRaw := tagSelectorRaw["conditions"]
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
				tagSelectorMap["conditions"] = conditionsMaps
				tagSelectorMaps = append(tagSelectorMaps, tagSelectorMap)
			}
			workspaceFilterSettingMap["tag_selector"] = tagSelectorMaps
			workspaceUuidsRaw := make([]interface{}, 0)
			if workspaceFilterSettingRaw["workspaceUuids"] != nil {
				workspaceUuidsRaw = convertToInterfaceArray(workspaceFilterSettingRaw["workspaceUuids"])
			}

			workspaceFilterSettingMap["workspace_uuids"] = workspaceUuidsRaw
			workspaceFilterSettingMaps = append(workspaceFilterSettingMaps, workspaceFilterSettingMap)
		}
		subscriptionMap["workspace_filter_setting"] = workspaceFilterSettingMaps
		subscriptionMaps = append(subscriptionMaps, subscriptionMap)
	}
	mapping["subscription"] = subscriptionMaps

	return mapping, nil
}
