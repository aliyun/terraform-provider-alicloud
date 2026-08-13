// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCmsEventNotifyPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsEventNotifyPolicyCreate,
		Read:   resourceAliCloudCmsEventNotifyPolicyRead,
		Update: resourceAliCloudCmsEventNotifyPolicyUpdate,
		Delete: resourceAliCloudCmsEventNotifyPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ignore_restored_notification": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"grouping_setting": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"grouping_keys": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"period_min": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"times": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"silence_sec": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"routes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"digital_employee_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enable_rca": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"filter_setting": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"relation": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"expression": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"conditions": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"field": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"op": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
									"effect_time_range": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_zone": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"start_time_in_minute": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"end_time_in_minute": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"day_in_week": {
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeInt},
												},
											},
										},
									},
									"channels": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"receivers": {
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"channel_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"enabled_sub_channels": {
													Type:     schema.TypeSet,
													Optional: true,
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
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"template_uuid": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"response_plan": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repeat_notify_setting": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_incident_state": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"repeat_interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"pushing_setting": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"restore_action_ids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"alert_action_ids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"escalation_id": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"auto_recover_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"subscription": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workspace_filter_setting": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workspace_uuids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"tag_selector": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"relation": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"expression": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"conditions": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"field": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"op": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"value": {
																Type:     schema.TypeString,
																Optional: true,
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
							Optional: true,
						},
						"filter_setting": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"relation": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"expression": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"conditions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"op": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
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
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsEventNotifyPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/eventbase/notify-policy/create")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["workspace"] = StringPointer(d.Get("workspace").(string))
	query["regionId"] = StringPointer(client.RegionId)

	notifyStrategy := make(map[string]interface{})

	if v := d.Get("notify_strategy"); !IsNil(v) {
		groupingSetting := make(map[string]interface{})
		times1, _ := jsonpath.Get("$[0].grouping_setting[0].times", d.Get("notify_strategy"))
		if times1 != nil && times1 != "" {
			groupingSetting["times"] = times1
		}
		groupingKeys1, _ := jsonpath.Get("$[0].grouping_setting[0].grouping_keys", d.Get("notify_strategy"))
		if groupingKeys1 != nil && groupingKeys1 != "" {
			groupingSetting["groupingKeys"] = groupingKeys1
		}
		periodMin1, _ := jsonpath.Get("$[0].grouping_setting[0].period_min", d.Get("notify_strategy"))
		if periodMin1 != nil && periodMin1 != "" {
			groupingSetting["periodMin"] = periodMin1
		}
		silenceSec1, _ := jsonpath.Get("$[0].grouping_setting[0].silence_sec", d.Get("notify_strategy"))
		if silenceSec1 != nil && silenceSec1 != "" {
			groupingSetting["silenceSec"] = silenceSec1
		}

		if len(groupingSetting) > 0 {
			notifyStrategy["groupingSetting"] = groupingSetting
		}
		if v, ok := d.GetOk("notify_strategy"); ok {
			localData, err := jsonpath.Get("$[0].routes", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				localData1 := make(map[string]interface{})
				startTimeInMinute1, _ := jsonpath.Get("$[0].start_time_in_minute", dataLoopTmp["effect_time_range"])
				if startTimeInMinute1 != nil && startTimeInMinute1 != "" {
					localData1["startTimeInMinute"] = startTimeInMinute1
				}
				endTimeInMinute1, _ := jsonpath.Get("$[0].end_time_in_minute", dataLoopTmp["effect_time_range"])
				if endTimeInMinute1 != nil && endTimeInMinute1 != "" {
					localData1["endTimeInMinute"] = endTimeInMinute1
				}
				dayInWeek1, _ := jsonpath.Get("$[0].day_in_week", dataLoopTmp["effect_time_range"])
				if dayInWeek1 != nil && dayInWeek1 != "" {
					localData1["dayInWeek"] = dayInWeek1
				}
				timeZone1, _ := jsonpath.Get("$[0].time_zone", dataLoopTmp["effect_time_range"])
				if timeZone1 != nil && timeZone1 != "" {
					localData1["timeZone"] = timeZone1
				}
				if len(localData1) > 0 {
					dataLoopMap["effectTimeRange"] = localData1
				}
				localMaps2 := make([]interface{}, 0)
				localData2 := dataLoopTmp["channels"]
				for _, dataLoop2 := range convertToInterfaceArray(localData2) {
					dataLoop2Tmp := dataLoop2.(map[string]interface{})
					dataLoop2Map := make(map[string]interface{})
					dataLoop2Map["receivers"] = dataLoop2Tmp["receivers"]
					dataLoop2Map["channelType"] = dataLoop2Tmp["channel_type"]
					if enabledSubChannelsSet, ok := dataLoop2Tmp["enabled_sub_channels"].(*schema.Set); ok {
						dataLoop2Map["enabledSubChannels"] = enabledSubChannelsSet.List()
					} else {
						dataLoop2Map["enabledSubChannels"] = dataLoop2Tmp["enabled_sub_channels"]
					}
					localMaps2 = append(localMaps2, dataLoop2Map)
				}
				dataLoopMap["channels"] = localMaps2
				localData3 := make(map[string]interface{})
				if v, ok := dataLoopTmp["filter_setting"]; ok {
					localData4, err := jsonpath.Get("$[0].conditions", v)
					if err != nil {
						localData4 = make([]interface{}, 0)
					}
					localMaps4 := make([]interface{}, 0)
					for _, dataLoop4 := range convertToInterfaceArray(localData4) {
						dataLoop4Tmp := make(map[string]interface{})
						if dataLoop4 != nil {
							dataLoop4Tmp = dataLoop4.(map[string]interface{})
						}
						dataLoop4Map := make(map[string]interface{})
						dataLoop4Map["op"] = dataLoop4Tmp["op"]
						dataLoop4Map["field"] = dataLoop4Tmp["field"]
						dataLoop4Map["value"] = dataLoop4Tmp["value"]
						localMaps4 = append(localMaps4, dataLoop4Map)
					}
					localData3["conditions"] = localMaps4
				}

				relation1, _ := jsonpath.Get("$[0].relation", dataLoopTmp["filter_setting"])
				if relation1 != nil && relation1 != "" {
					localData3["relation"] = relation1
				}
				expression1, _ := jsonpath.Get("$[0].expression", dataLoopTmp["filter_setting"])
				if expression1 != nil && expression1 != "" {
					localData3["expression"] = expression1
				}
				if len(localData3) > 0 {
					dataLoopMap["filterSetting"] = localData3
				}
				dataLoopMap["enableRca"] = dataLoopTmp["enable_rca"]
				dataLoopMap["digitalEmployeeName"] = dataLoopTmp["digital_employee_name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			notifyStrategy["routes"] = localMaps
		}

		if v, ok := d.GetOk("notify_strategy"); ok {
			localData5, err := jsonpath.Get("$[0].custom_template_entries", v)
			if err != nil {
				localData5 = make([]interface{}, 0)
			}
			localMaps5 := make([]interface{}, 0)
			for _, dataLoop5 := range convertToInterfaceArray(localData5) {
				dataLoop5Tmp := make(map[string]interface{})
				if dataLoop5 != nil {
					dataLoop5Tmp = dataLoop5.(map[string]interface{})
				}
				dataLoop5Map := make(map[string]interface{})
				dataLoop5Map["templateUuid"] = dataLoop5Tmp["template_uuid"]
				localMaps5 = append(localMaps5, dataLoop5Map)
			}
			notifyStrategy["customTemplateEntries"] = localMaps5
		}

		ignoreRestoredNotification1, _ := jsonpath.Get("$[0].ignore_restored_notification", v)
		if ignoreRestoredNotification1 != nil && ignoreRestoredNotification1 != "" {
			notifyStrategy["ignoreRestoredNotification"] = ignoreRestoredNotification1
		}
		description1, _ := jsonpath.Get("$[0].description", v)
		if description1 != nil && description1 != "" {
			notifyStrategy["description"] = description1
		}

		request["notifyStrategy"] = notifyStrategy
	}

	responsePlan := make(map[string]interface{})

	if v := d.Get("response_plan"); !IsNil(v) {
		autoRecoverSeconds1, _ := jsonpath.Get("$[0].auto_recover_seconds", v)
		if autoRecoverSeconds1 != nil && autoRecoverSeconds1 != "" {
			responsePlan["autoRecoverSeconds"] = autoRecoverSeconds1
		}
		pushingSetting := make(map[string]interface{})
		restoreActionIds1, _ := jsonpath.Get("$[0].pushing_setting[0].restore_action_ids", d.Get("response_plan"))
		if restoreActionIds1 != nil && restoreActionIds1 != "" {
			pushingSetting["restoreActionIds"] = restoreActionIds1
		}
		alertActionIds1, _ := jsonpath.Get("$[0].pushing_setting[0].alert_action_ids", d.Get("response_plan"))
		if alertActionIds1 != nil && alertActionIds1 != "" {
			pushingSetting["alertActionIds"] = alertActionIds1
		}

		if len(pushingSetting) > 0 {
			responsePlan["pushingSetting"] = pushingSetting
		}
		repeatNotifySetting := make(map[string]interface{})
		repeatInterval1, _ := jsonpath.Get("$[0].repeat_notify_setting[0].repeat_interval", d.Get("response_plan"))
		if repeatInterval1 != nil && repeatInterval1 != "" {
			repeatNotifySetting["repeatInterval"] = repeatInterval1
		}
		endIncidentState1, _ := jsonpath.Get("$[0].repeat_notify_setting[0].end_incident_state", d.Get("response_plan"))
		if endIncidentState1 != nil && endIncidentState1 != "" {
			repeatNotifySetting["endIncidentState"] = endIncidentState1
		}

		if len(repeatNotifySetting) > 0 {
			responsePlan["repeatNotifySetting"] = repeatNotifySetting
		}
		escalationId1, _ := jsonpath.Get("$[0].escalation_id", v)
		if escalationId1 != nil && escalationId1 != "" {
			responsePlan["escalationId"] = escalationId1
		}

		request["responsePlan"] = responsePlan
	}

	subscription := make(map[string]interface{})

	if v := d.Get("subscription"); !IsNil(v) {
		// workspace_filter_setting is read with plain type assertions instead of
		// jsonpath: the nested tag_selector list makes the jsonpath lookups return
		// nothing, which silently dropped both tagSelector and workspaceUuids.
		workspaceFilterSetting := buildEventNotifyPolicyWorkspaceFilterSetting(d.Get("subscription"))

		if len(workspaceFilterSetting) > 0 {
			subscription["workspaceFilterSetting"] = workspaceFilterSetting
		}
		filterSetting1 := make(map[string]interface{})
		if v, ok := d.GetOk("subscription"); ok {
			localData7, err := jsonpath.Get("$[0].filter_setting[0].conditions", v)
			if err != nil {
				localData7 = make([]interface{}, 0)
			}
			localMaps7 := make([]interface{}, 0)
			for _, dataLoop7 := range convertToInterfaceArray(localData7) {
				dataLoop7Tmp := make(map[string]interface{})
				if dataLoop7 != nil {
					dataLoop7Tmp = dataLoop7.(map[string]interface{})
				}
				dataLoop7Map := make(map[string]interface{})
				dataLoop7Map["field"] = dataLoop7Tmp["field"]
				dataLoop7Map["value"] = dataLoop7Tmp["value"]
				dataLoop7Map["op"] = dataLoop7Tmp["op"]
				localMaps7 = append(localMaps7, dataLoop7Map)
			}
			filterSetting1["conditions"] = localMaps7
		}

		relation5, _ := jsonpath.Get("$[0].filter_setting[0].relation", d.Get("subscription"))
		if relation5 != nil && relation5 != "" {
			filterSetting1["relation"] = relation5
		}
		expression5, _ := jsonpath.Get("$[0].filter_setting[0].expression", d.Get("subscription"))
		if expression5 != nil && expression5 != "" {
			filterSetting1["expression"] = expression5
		}

		if len(filterSetting1) > 0 {
			subscription["filterSetting"] = filterSetting1
		}
		subscribeLegacyEvent1, _ := jsonpath.Get("$[0].subscribe_legacy_event", v)
		if subscribeLegacyEvent1 != nil && subscribeLegacyEvent1 != "" {
			subscription["subscribeLegacyEvent"] = subscribeLegacyEvent1
		}

		request["subscription"] = subscription
	}

	if v, ok := d.GetOk("name"); ok {
		request["name"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_event_notify_policy", action, AlibabaCloudSdkGoERROR)
	}

	notifyPolicyuuidVar, _ := jsonpath.Get("$.notifyPolicy.uuid", response)
	notifyPolicyworkspaceVar, _ := jsonpath.Get("$.notifyPolicy.workspace", response)
	d.SetId(fmt.Sprintf("%v:%v", notifyPolicyuuidVar, notifyPolicyworkspaceVar))

	return resourceAliCloudCmsEventNotifyPolicyUpdate(d, meta)
}

func resourceAliCloudCmsEventNotifyPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsEventNotifyPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_event_notify_policy DescribeCmsEventNotifyPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("enabled", objectRaw["enabled"])
	d.Set("name", objectRaw["name"])
	d.Set("update_time", objectRaw["updateTime"])
	d.Set("user_id", objectRaw["userId"])
	d.Set("version", objectRaw["version"])
	d.Set("uuid", objectRaw["uuid"])
	d.Set("workspace", objectRaw["workspace"])

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
	if err := d.Set("notify_strategy", notifyStrategyMaps); err != nil {
		return err
	}
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
	if err := d.Set("response_plan", responsePlanMaps); err != nil {
		return err
	}
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
	if err := d.Set("subscription", subscriptionMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudCmsEventNotifyPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	cmsServiceV2 := CmsServiceV2{client}
	objectRaw, _ := cmsServiceV2.DescribeCmsEventNotifyPolicy(d.Id())

	initedEnabled := false
	if _, ok := d.GetOkExists("enabled"); ok && d.IsNewResource() {
		initedEnabled = true
	}
	if initedEnabled || d.HasChange("enabled") {
		var err error
		target := d.Get("enabled").(bool)

		currentStatus, err := jsonpath.Get("enabled", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "enabled", objectRaw)
		}
		if formatBool(currentStatus) != target {
			if target == true {
				parts := strings.Split(d.Id(), ":")
				uuid := parts[0]
				action := fmt.Sprintf("/api/eventbase/notify-policy/%s/enable", uuid)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				query["workspace"] = StringPointer(parts[1])
				query["regionId"] = StringPointer(client.RegionId)
				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
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

			}
			if target == false {
				parts := strings.Split(d.Id(), ":")
				uuid := parts[0]
				action := fmt.Sprintf("/api/eventbase/notify-policy/%s/disable", uuid)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				query["workspace"] = StringPointer(parts[1])
				query["regionId"] = StringPointer(client.RegionId)
				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
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

			}
		}
	}

	var err error
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/api/eventbase/notify-policy/update")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["workspace"] = StringPointer(parts[1])
	request["uuid"] = parts[0]
	query["regionId"] = StringPointer(client.RegionId)
	// The update API uses version as an optimistic lock. It must carry the value
	// currently held by the server, otherwise the request fails with OptimisticLockFailed.
	// The enable/disable calls above bump the version, so the policy is re-read here
	// instead of reusing the snapshot taken at the beginning of this function.
	versionSource := objectRaw
	if latestRaw, describeErr := cmsServiceV2.DescribeCmsEventNotifyPolicy(d.Id()); describeErr == nil {
		versionSource = latestRaw
	}
	if currentVersion, ok := versionSource["version"]; ok && currentVersion != nil {
		request["version"] = currentVersion
	}
	if d.HasChange("notify_strategy") {
		update = true
	}
	notifyStrategy := make(map[string]interface{})

	if v := d.Get("notify_strategy"); !IsNil(v) || d.HasChange("notify_strategy") {
		groupingSetting := make(map[string]interface{})
		times1, _ := jsonpath.Get("$[0].grouping_setting[0].times", d.Get("notify_strategy"))
		if times1 != nil && times1 != "" {
			groupingSetting["times"] = times1
		}
		groupingKeys1, _ := jsonpath.Get("$[0].grouping_setting[0].grouping_keys", d.Get("notify_strategy"))
		if groupingKeys1 != nil && groupingKeys1 != "" {
			groupingSetting["groupingKeys"] = groupingKeys1
		}
		periodMin1, _ := jsonpath.Get("$[0].grouping_setting[0].period_min", d.Get("notify_strategy"))
		if periodMin1 != nil && periodMin1 != "" {
			groupingSetting["periodMin"] = periodMin1
		}
		silenceSec1, _ := jsonpath.Get("$[0].grouping_setting[0].silence_sec", d.Get("notify_strategy"))
		if silenceSec1 != nil && silenceSec1 != "" {
			groupingSetting["silenceSec"] = silenceSec1
		}

		if len(groupingSetting) > 0 {
			notifyStrategy["groupingSetting"] = groupingSetting
		}
		if v, ok := d.GetOk("notify_strategy"); ok {
			localData, err := jsonpath.Get("$[0].routes", v)
			if err != nil {
				localData = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(localData) {
				dataLoopTmp := make(map[string]interface{})
				if dataLoop != nil {
					dataLoopTmp = dataLoop.(map[string]interface{})
				}
				dataLoopMap := make(map[string]interface{})
				if !IsNil(dataLoopTmp["effect_time_range"]) {
					localData1 := make(map[string]interface{})
					startTimeInMinute1, _ := jsonpath.Get("$[0].start_time_in_minute", dataLoopTmp["effect_time_range"])
					if startTimeInMinute1 != nil && startTimeInMinute1 != "" {
						localData1["startTimeInMinute"] = startTimeInMinute1
					}
					endTimeInMinute1, _ := jsonpath.Get("$[0].end_time_in_minute", dataLoopTmp["effect_time_range"])
					if endTimeInMinute1 != nil && endTimeInMinute1 != "" {
						localData1["endTimeInMinute"] = endTimeInMinute1
					}
					dayInWeek1, _ := jsonpath.Get("$[0].day_in_week", dataLoopTmp["effect_time_range"])
					if dayInWeek1 != nil && dayInWeek1 != "" {
						localData1["dayInWeek"] = dayInWeek1
					}
					timeZone1, _ := jsonpath.Get("$[0].time_zone", dataLoopTmp["effect_time_range"])
					if timeZone1 != nil && timeZone1 != "" {
						localData1["timeZone"] = timeZone1
					}
					if len(localData1) > 0 {
						dataLoopMap["effectTimeRange"] = localData1
					}
				}
				localMaps2 := make([]interface{}, 0)
				localData2 := dataLoopTmp["channels"]
				for _, dataLoop2 := range convertToInterfaceArray(localData2) {
					dataLoop2Tmp := dataLoop2.(map[string]interface{})
					dataLoop2Map := make(map[string]interface{})
					dataLoop2Map["receivers"] = dataLoop2Tmp["receivers"]
					dataLoop2Map["channelType"] = dataLoop2Tmp["channel_type"]
					if enabledSubChannelsSet, ok := dataLoop2Tmp["enabled_sub_channels"].(*schema.Set); ok {
						dataLoop2Map["enabledSubChannels"] = enabledSubChannelsSet.List()
					} else {
						dataLoop2Map["enabledSubChannels"] = dataLoop2Tmp["enabled_sub_channels"]
					}
					localMaps2 = append(localMaps2, dataLoop2Map)
				}
				dataLoopMap["channels"] = localMaps2
				if !IsNil(dataLoopTmp["filter_setting"]) {
					localData3 := make(map[string]interface{})
					if v, ok := dataLoopTmp["filter_setting"]; ok {
						localData4, err := jsonpath.Get("$[0].conditions", v)
						if err != nil {
							localData4 = make([]interface{}, 0)
						}
						localMaps4 := make([]interface{}, 0)
						for _, dataLoop4 := range convertToInterfaceArray(localData4) {
							dataLoop4Tmp := make(map[string]interface{})
							if dataLoop4 != nil {
								dataLoop4Tmp = dataLoop4.(map[string]interface{})
							}
							dataLoop4Map := make(map[string]interface{})
							dataLoop4Map["op"] = dataLoop4Tmp["op"]
							dataLoop4Map["field"] = dataLoop4Tmp["field"]
							dataLoop4Map["value"] = dataLoop4Tmp["value"]
							localMaps4 = append(localMaps4, dataLoop4Map)
						}
						localData3["conditions"] = localMaps4
					}

					relation1, _ := jsonpath.Get("$[0].relation", dataLoopTmp["filter_setting"])
					if relation1 != nil && relation1 != "" {
						localData3["relation"] = relation1
					}
					expression1, _ := jsonpath.Get("$[0].expression", dataLoopTmp["filter_setting"])
					if expression1 != nil && expression1 != "" {
						localData3["expression"] = expression1
					}
					if len(localData3) > 0 {
						dataLoopMap["filterSetting"] = localData3
					}
				}
				dataLoopMap["enableRca"] = dataLoopTmp["enable_rca"]
				dataLoopMap["digitalEmployeeName"] = dataLoopTmp["digital_employee_name"]
				localMaps = append(localMaps, dataLoopMap)
			}
			notifyStrategy["routes"] = localMaps
		}

		if v, ok := d.GetOk("notify_strategy"); ok {
			localData5, err := jsonpath.Get("$[0].custom_template_entries", v)
			if err != nil {
				localData5 = make([]interface{}, 0)
			}
			localMaps5 := make([]interface{}, 0)
			for _, dataLoop5 := range convertToInterfaceArray(localData5) {
				dataLoop5Tmp := make(map[string]interface{})
				if dataLoop5 != nil {
					dataLoop5Tmp = dataLoop5.(map[string]interface{})
				}
				dataLoop5Map := make(map[string]interface{})
				dataLoop5Map["templateUuid"] = dataLoop5Tmp["template_uuid"]
				localMaps5 = append(localMaps5, dataLoop5Map)
			}
			notifyStrategy["customTemplateEntries"] = localMaps5
		}

		ignoreRestoredNotification1, _ := jsonpath.Get("$[0].ignore_restored_notification", v)
		if ignoreRestoredNotification1 != nil && ignoreRestoredNotification1 != "" {
			notifyStrategy["ignoreRestoredNotification"] = ignoreRestoredNotification1
		}
		description1, _ := jsonpath.Get("$[0].description", v)
		if description1 != nil && description1 != "" {
			notifyStrategy["description"] = description1
		}

		request["notifyStrategy"] = notifyStrategy
	}

	if !d.IsNewResource() && d.HasChange("response_plan") {
		update = true
	}
	responsePlan := make(map[string]interface{})

	if v := d.Get("response_plan"); !IsNil(v) || d.HasChange("response_plan") {
		autoRecoverSeconds1, _ := jsonpath.Get("$[0].auto_recover_seconds", v)
		if autoRecoverSeconds1 != nil && autoRecoverSeconds1 != "" {
			responsePlan["autoRecoverSeconds"] = autoRecoverSeconds1
		}
		pushingSetting := make(map[string]interface{})
		restoreActionIds1, _ := jsonpath.Get("$[0].pushing_setting[0].restore_action_ids", d.Get("response_plan"))
		if restoreActionIds1 != nil && restoreActionIds1 != "" {
			pushingSetting["restoreActionIds"] = restoreActionIds1
		}
		alertActionIds1, _ := jsonpath.Get("$[0].pushing_setting[0].alert_action_ids", d.Get("response_plan"))
		if alertActionIds1 != nil && alertActionIds1 != "" {
			pushingSetting["alertActionIds"] = alertActionIds1
		}

		if len(pushingSetting) > 0 {
			responsePlan["pushingSetting"] = pushingSetting
		}
		repeatNotifySetting := make(map[string]interface{})
		repeatInterval1, _ := jsonpath.Get("$[0].repeat_notify_setting[0].repeat_interval", d.Get("response_plan"))
		if repeatInterval1 != nil && repeatInterval1 != "" {
			repeatNotifySetting["repeatInterval"] = repeatInterval1
		}
		endIncidentState1, _ := jsonpath.Get("$[0].repeat_notify_setting[0].end_incident_state", d.Get("response_plan"))
		if endIncidentState1 != nil && endIncidentState1 != "" {
			repeatNotifySetting["endIncidentState"] = endIncidentState1
		}

		if len(repeatNotifySetting) > 0 {
			responsePlan["repeatNotifySetting"] = repeatNotifySetting
		}
		escalationId1, _ := jsonpath.Get("$[0].escalation_id", v)
		if escalationId1 != nil && escalationId1 != "" {
			responsePlan["escalationId"] = escalationId1
		}

		request["responsePlan"] = responsePlan
	}

	if d.HasChange("subscription") {
		update = true
	}
	subscription := make(map[string]interface{})

	if v := d.Get("subscription"); !IsNil(v) || d.HasChange("subscription") {
		// workspace_filter_setting is read with plain type assertions instead of
		// jsonpath: the nested tag_selector list makes the jsonpath lookups return
		// nothing, which silently dropped both tagSelector and workspaceUuids.
		workspaceFilterSetting := buildEventNotifyPolicyWorkspaceFilterSetting(d.Get("subscription"))

		if len(workspaceFilterSetting) > 0 {
			subscription["workspaceFilterSetting"] = workspaceFilterSetting
		}
		filterSetting1 := make(map[string]interface{})
		if v, ok := d.GetOk("subscription"); ok {
			localData7, err := jsonpath.Get("$[0].filter_setting[0].conditions", v)
			if err != nil {
				localData7 = make([]interface{}, 0)
			}
			localMaps7 := make([]interface{}, 0)
			for _, dataLoop7 := range convertToInterfaceArray(localData7) {
				dataLoop7Tmp := make(map[string]interface{})
				if dataLoop7 != nil {
					dataLoop7Tmp = dataLoop7.(map[string]interface{})
				}
				dataLoop7Map := make(map[string]interface{})
				dataLoop7Map["field"] = dataLoop7Tmp["field"]
				dataLoop7Map["value"] = dataLoop7Tmp["value"]
				dataLoop7Map["op"] = dataLoop7Tmp["op"]
				localMaps7 = append(localMaps7, dataLoop7Map)
			}
			filterSetting1["conditions"] = localMaps7
		}

		relation5, _ := jsonpath.Get("$[0].filter_setting[0].relation", d.Get("subscription"))
		if relation5 != nil && relation5 != "" {
			filterSetting1["relation"] = relation5
		}
		expression5, _ := jsonpath.Get("$[0].filter_setting[0].expression", d.Get("subscription"))
		if expression5 != nil && expression5 != "" {
			filterSetting1["expression"] = expression5
		}

		if len(filterSetting1) > 0 {
			subscription["filterSetting"] = filterSetting1
		}
		subscribeLegacyEvent1, _ := jsonpath.Get("$[0].subscribe_legacy_event", v)
		if subscribeLegacyEvent1 != nil && subscribeLegacyEvent1 != "" {
			subscription["subscribeLegacyEvent"] = subscribeLegacyEvent1
		}

		request["subscription"] = subscription
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
	}
	if v, ok := d.GetOk("name"); ok || d.HasChange("name") {
		request["name"] = v
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
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
	}

	return resourceAliCloudCmsEventNotifyPolicyRead(d, meta)
}

func resourceAliCloudCmsEventNotifyPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/api/eventbase/notify-policy")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["uuid"] = StringPointer(parts[0])
	query["workspace"] = StringPointer(parts[1])
	query["regionId"] = StringPointer(client.RegionId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// buildEventNotifyPolicyWorkspaceFilterSetting converts the subscription
// workspace_filter_setting block into the request payload shape.
func buildEventNotifyPolicyWorkspaceFilterSetting(raw interface{}) map[string]interface{} {
	workspaceFilterSetting := make(map[string]interface{})

	subscriptionList, ok := raw.([]interface{})
	if !ok || len(subscriptionList) == 0 {
		return workspaceFilterSetting
	}
	subscriptionMap, ok := subscriptionList[0].(map[string]interface{})
	if !ok {
		return workspaceFilterSetting
	}
	settingList, ok := subscriptionMap["workspace_filter_setting"].([]interface{})
	if !ok || len(settingList) == 0 {
		return workspaceFilterSetting
	}
	settingMap, ok := settingList[0].(map[string]interface{})
	if !ok {
		return workspaceFilterSetting
	}

	if workspaceUuids, ok := settingMap["workspace_uuids"].([]interface{}); ok && len(workspaceUuids) > 0 {
		workspaceFilterSetting["workspaceUuids"] = workspaceUuids
	}

	selectorList, ok := settingMap["tag_selector"].([]interface{})
	if !ok || len(selectorList) == 0 {
		return workspaceFilterSetting
	}
	selectorMap, ok := selectorList[0].(map[string]interface{})
	if !ok {
		return workspaceFilterSetting
	}

	tagSelector := make(map[string]interface{})
	if relation, ok := selectorMap["relation"].(string); ok && relation != "" {
		tagSelector["relation"] = relation
	}
	if expression, ok := selectorMap["expression"].(string); ok && expression != "" {
		tagSelector["expression"] = expression
	}
	conditions := make([]interface{}, 0)
	if conditionList, ok := selectorMap["conditions"].([]interface{}); ok {
		for _, condition := range conditionList {
			conditionMap, ok := condition.(map[string]interface{})
			if !ok {
				continue
			}
			conditions = append(conditions, map[string]interface{}{
				"field": conditionMap["field"],
				"op":    conditionMap["op"],
				"value": conditionMap["value"],
			})
		}
	}
	tagSelector["conditions"] = conditions

	if len(tagSelector) > 0 {
		workspaceFilterSetting["tagSelector"] = tagSelector
	}
	return workspaceFilterSetting
}
