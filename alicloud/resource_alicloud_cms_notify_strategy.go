package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsNotifyStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsNotifyStrategyCreate,
		Read:   resourceAliCloudCmsNotifyStrategyRead,
		Update: resourceAliCloudCmsNotifyStrategyUpdate,
		Delete: resourceAliCloudCmsNotifyStrategyDelete,
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
			"custom_template_entries": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"template_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"grouping_setting": {
				Type:     schema.TypeList,
				Required: true,
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
						"silence_sec": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"ignore_restored_notification": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"notify_strategy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notify_strategy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"enabled_sub_channels": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"receivers": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
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
									"day_in_week": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"end_time_in_minute": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"start_time_in_minute": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"time_zone": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"filter_setting": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"conditions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Required: true,
												},
												"op": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"expression": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"relation": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"severities": {
							Type:     schema.TypeList,
							Optional: true,
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
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsNotifyStrategyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/notifyStrategies")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}

	request["name"] = d.Get("notify_strategy_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOkExists("ignore_restored_notification"); ok {
		request["ignoreRestoredNotification"] = v
	}
	if v := d.Get("grouping_setting"); !IsNil(v) {
		request["groupingSetting"] = expandCmsNotifyStrategyGroupingSetting(v)
	}
	if v := d.Get("routes"); !IsNil(v) {
		request["routes"] = expandCmsNotifyStrategyRoutes(v)
	}
	if v := d.Get("custom_template_entries"); !IsNil(v) {
		request["customTemplateEntries"] = expandCmsNotifyStrategyCustomTemplateEntries(v)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_notify_strategy", action, AlibabaCloudSdkGoERROR)
	}

	id, err := jsonpath.Get("$.data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_cms_notify_strategy", "$.data", response)
	}
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCmsNotifyStrategyRead(d, meta)
}

func resourceAliCloudCmsNotifyStrategyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsNotifyStrategy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_notify_strategy DescribeCmsNotifyStrategy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	objectRawObj, _ := jsonpath.Get("$.data", objectRaw)
	objectRawMap := make(map[string]interface{})
	if objectRawObj != nil {
		objectRawMap, _ = objectRawObj.(map[string]interface{})
	}
	d.Set("create_time", objectRawMap["createTime"])
	d.Set("description", objectRawMap["description"])
	d.Set("enable", objectRawMap["enabled"])
	d.Set("ignore_restored_notification", objectRawMap["ignoreRestoredNotification"])
	d.Set("notify_strategy_id", objectRawMap["notifyStrategyId"])
	d.Set("notify_strategy_name", objectRawMap["notifyStrategyName"])
	d.Set("update_time", objectRawMap["updateTime"])
	d.Set("user_id", objectRawMap["userId"])
	d.Set("workspace", objectRawMap["workspace"])

	groupingSettingMaps := make([]map[string]interface{}, 0)
	groupingSettingMap := make(map[string]interface{})
	groupingSettingRaw := make(map[string]interface{})
	if objectRawMap["groupingSetting"] != nil {
		groupingSettingRaw, _ = objectRawMap["groupingSetting"].(map[string]interface{})
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
	if err := d.Set("grouping_setting", groupingSettingMaps); err != nil {
		return err
	}
	customTemplateEntriesRaw := objectRawMap["customTemplateEntries"]
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
	if err := d.Set("custom_template_entries", customTemplateEntriesMaps); err != nil {
		return err
	}
	routesRaw := objectRawMap["routes"]
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
	if err := d.Set("routes", routesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudCmsNotifyStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var body map[string]interface{}
	update := false

	var err error
	notifyStrategyId := d.Id()
	action := fmt.Sprintf("/notifyStrategies/%s", notifyStrategyId)
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}
	body = make(map[string]interface{})

	if d.HasChange("notify_strategy_name") {
		update = true
	}
	request["name"] = d.Get("notify_strategy_name")
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOkExists("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("ignore_restored_notification") {
		update = true
	}
	if v, ok := d.GetOkExists("ignore_restored_notification"); ok || d.HasChange("ignore_restored_notification") {
		request["ignoreRestoredNotification"] = v
	}
	if d.HasChange("grouping_setting") {
		update = true
	}
	if v := d.Get("grouping_setting"); !IsNil(v) || d.HasChange("grouping_setting") {
		request["groupingSetting"] = expandCmsNotifyStrategyGroupingSetting(v)
	}
	if d.HasChange("routes") {
		update = true
	}
	if v := d.Get("routes"); !IsNil(v) || d.HasChange("routes") {
		request["routes"] = expandCmsNotifyStrategyRoutes(v)
	}
	if d.HasChange("custom_template_entries") {
		update = true
	}
	if v := d.Get("custom_template_entries"); !IsNil(v) || d.HasChange("custom_template_entries") {
		request["customTemplateEntries"] = expandCmsNotifyStrategyCustomTemplateEntries(v)
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

	return resourceAliCloudCmsNotifyStrategyRead(d, meta)
}

func resourceAliCloudCmsNotifyStrategyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	notifyStrategyId := d.Id()
	action := fmt.Sprintf("/notifyStrategies/%s", notifyStrategyId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}

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
		if IsExpectedErrors(err, []string{"NotFound", "ResourceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func expandCmsNotifyStrategyGroupingSetting(v interface{}) map[string]interface{} {
	groupingSetting := make(map[string]interface{})
	groupingSettingMap, _ := jsonpath.Get("$[0]", v)
	if groupingSettingMap == nil {
		return groupingSetting
	}
	groupingSettingArg := groupingSettingMap.(map[string]interface{})

	if groupingKeys, ok := groupingSettingArg["grouping_keys"]; ok {
		groupingSetting["groupingKeys"] = convertToInterfaceArray(groupingKeys)
	}
	if periodMin, ok := groupingSettingArg["period_min"]; ok {
		if periodMinInt, isInt := periodMin.(int); isInt && periodMinInt != 0 {
			groupingSetting["periodMin"] = periodMin
		}
	}
	if silenceSec, ok := groupingSettingArg["silence_sec"]; ok {
		if silenceSecInt, isInt := silenceSec.(int); isInt && silenceSecInt != 0 {
			groupingSetting["silenceSec"] = silenceSec
		}
	}
	if times, ok := groupingSettingArg["times"]; ok {
		if timesInt, isInt := times.(int); isInt && timesInt != 0 {
			groupingSetting["times"] = times
		}
	}

	return groupingSetting
}

func expandCmsNotifyStrategyRoutes(v interface{}) []interface{} {
	routesList := make([]interface{}, 0)
	for _, routesChild := range convertToInterfaceArray(v) {
		routesChildArg := routesChild.(map[string]interface{})
		routesMap := make(map[string]interface{})

		if channels, ok := routesChildArg["channels"]; ok {
			channelsList := make([]interface{}, 0)
			for _, channelsChild := range convertToInterfaceArray(channels) {
				channelsChildArg := channelsChild.(map[string]interface{})
				channelsMap := make(map[string]interface{})
				if channelType, ok := channelsChildArg["channel_type"]; ok && channelType != "" {
					channelsMap["channelType"] = channelType
				}
				if enabledSubChannels, ok := channelsChildArg["enabled_sub_channels"]; ok {
					channelsMap["enabledSubChannels"] = convertToInterfaceArray(enabledSubChannels)
				}
				if receivers, ok := channelsChildArg["receivers"]; ok {
					channelsMap["receivers"] = convertToInterfaceArray(receivers)
				}
				channelsList = append(channelsList, channelsMap)
			}
			routesMap["channels"] = channelsList
		}

		if effectTimeRange, ok := routesChildArg["effect_time_range"]; ok {
			effectTimeRangeList := convertToInterfaceArray(effectTimeRange)
			if len(effectTimeRangeList) > 0 {
				effectTimeRangeArg := effectTimeRangeList[0].(map[string]interface{})
				effectTimeRangeMap := make(map[string]interface{})
				if timeZone, ok := effectTimeRangeArg["time_zone"]; ok && timeZone != "" {
					effectTimeRangeMap["timeZone"] = timeZone
				}
				if dayInWeek, ok := effectTimeRangeArg["day_in_week"]; ok {
					effectTimeRangeMap["dayInWeek"] = convertToInterfaceArray(dayInWeek)
				}
				if startTimeInMinute, ok := effectTimeRangeArg["start_time_in_minute"]; ok {
					effectTimeRangeMap["startTimeInMinute"] = startTimeInMinute
				}
				if endTimeInMinute, ok := effectTimeRangeArg["end_time_in_minute"]; ok {
					effectTimeRangeMap["endTimeInMinute"] = endTimeInMinute
				}
				routesMap["effectTimeRange"] = effectTimeRangeMap
			}
		}

		if filterSetting, ok := routesChildArg["filter_setting"]; ok {
			filterSettingList := convertToInterfaceArray(filterSetting)
			if len(filterSettingList) > 0 {
				filterSettingArg := filterSettingList[0].(map[string]interface{})
				filterSettingMap := make(map[string]interface{})
				if expression, ok := filterSettingArg["expression"]; ok && expression != "" {
					filterSettingMap["expression"] = expression
				}
				if relation, ok := filterSettingArg["relation"]; ok && relation != "" {
					filterSettingMap["relation"] = relation
				}
				if conditions, ok := filterSettingArg["conditions"]; ok {
					conditionsList := make([]interface{}, 0)
					for _, conditionsChild := range convertToInterfaceArray(conditions) {
						conditionsChildArg := conditionsChild.(map[string]interface{})
						conditionsMap := make(map[string]interface{})
						if field, ok := conditionsChildArg["field"]; ok && field != "" {
							conditionsMap["field"] = field
						}
						if op, ok := conditionsChildArg["op"]; ok && op != "" {
							conditionsMap["op"] = op
						}
						if value, ok := conditionsChildArg["value"]; ok && value != "" {
							conditionsMap["value"] = value
						}
						conditionsList = append(conditionsList, conditionsMap)
					}
					filterSettingMap["conditions"] = conditionsList
				}
				routesMap["filterSetting"] = filterSettingMap
			}
		}

		if severities, ok := routesChildArg["severities"]; ok {
			routesMap["severities"] = convertToInterfaceArray(severities)
		}

		routesList = append(routesList, routesMap)
	}

	return routesList
}

func expandCmsNotifyStrategyCustomTemplateEntries(v interface{}) []interface{} {
	customTemplateEntriesList := make([]interface{}, 0)
	for _, customTemplateEntriesChild := range convertToInterfaceArray(v) {
		customTemplateEntriesChildArg := customTemplateEntriesChild.(map[string]interface{})
		customTemplateEntriesMap := make(map[string]interface{})
		if targetType, ok := customTemplateEntriesChildArg["target_type"]; ok && targetType != "" {
			customTemplateEntriesMap["targetType"] = targetType
		}
		if templateUuid, ok := customTemplateEntriesChildArg["template_uuid"]; ok && templateUuid != "" {
			customTemplateEntriesMap["templateUuid"] = templateUuid
		}
		customTemplateEntriesList = append(customTemplateEntriesList, customTemplateEntriesMap)
	}

	return customTemplateEntriesList
}
