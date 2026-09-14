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

func resourceAliCloudCmsSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsSubscriptionCreate,
		Read:   resourceAliCloudCmsSubscriptionRead,
		Update: resourceAliCloudCmsSubscriptionUpdate,
		Delete: resourceAliCloudCmsSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agent_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"routes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channels": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"channel_type": {
													Type:     schema.TypeString,
													Optional: true,
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
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
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
			"notify_strategy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pushing_setting": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_action_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"response_plan_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"restore_action_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"template_uuid": {
							Type:     schema.TypeString,
							Optional: true,
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
				Required: true,
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
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/subscriptions")
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

	request["subscriptionName"] = d.Get("subscription_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("notify_strategy_id"); ok {
		request["notifyStrategyId"] = v
	}
	if v := d.Get("filter_setting"); !IsNil(v) {
		request["filterSetting"] = expandCmsSubscriptionFilterSetting(v)
	}
	if v := d.Get("pushing_setting"); !IsNil(v) {
		request["pushingSetting"] = expandCmsSubscriptionPushingSetting(v)
	}
	if v := d.Get("agent_config"); !IsNil(v) {
		request["agentConfig"] = expandCmsSubscriptionAgentConfig(v)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_subscription", action, AlibabaCloudSdkGoERROR)
	}

	id, err := jsonpath.Get("$.data", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_cms_subscription", "$.data", response)
	}
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCmsSubscriptionRead(d, meta)
}

func resourceAliCloudCmsSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsSubscription(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_subscription DescribeCmsSubscription Failed!!! %s", err)
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
	d.Set("notify_strategy_id", objectRawMap["notifyStrategyId"])
	d.Set("subscription_id", objectRawMap["subscriptionId"])
	d.Set("subscription_name", objectRawMap["subscriptionName"])
	d.Set("subscription_type", objectRawMap["subscriptionType"])
	d.Set("update_time", objectRawMap["updateTime"])
	d.Set("user_id", objectRawMap["userId"])
	d.Set("workspace", objectRawMap["workspace"])

	filterSettingMaps := make([]map[string]interface{}, 0)
	filterSettingMap := make(map[string]interface{})
	filterSettingRaw := make(map[string]interface{})
	if objectRawMap["filterSetting"] != nil {
		filterSettingRaw, _ = objectRawMap["filterSetting"].(map[string]interface{})
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
	if err := d.Set("filter_setting", filterSettingMaps); err != nil {
		return err
	}
	pushingSettingMaps := make([]map[string]interface{}, 0)
	pushingSettingMap := make(map[string]interface{})
	pushingSettingRaw := make(map[string]interface{})
	if objectRawMap["pushingSetting"] != nil {
		pushingSettingRaw, _ = objectRawMap["pushingSetting"].(map[string]interface{})
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
	if err := d.Set("pushing_setting", pushingSettingMaps); err != nil {
		return err
	}
	agentConfigMaps := make([]map[string]interface{}, 0)
	agentConfigMap := make(map[string]interface{})
	agentConfigRaw := make(map[string]interface{})
	if objectRawMap["agentConfig"] != nil {
		agentConfigRaw, _ = objectRawMap["agentConfig"].(map[string]interface{})
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
	if err := d.Set("agent_config", agentConfigMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudCmsSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var body map[string]interface{}
	update := false

	var err error
	subscriptionId := d.Id()
	action := fmt.Sprintf("/subscriptions/%s", subscriptionId)
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}
	body = make(map[string]interface{})

	if d.HasChange("subscription_name") {
		update = true
	}
	request["subscriptionName"] = d.Get("subscription_name")
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOkExists("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("notify_strategy_id") {
		update = true
	}
	if v, ok := d.GetOkExists("notify_strategy_id"); ok || d.HasChange("notify_strategy_id") {
		request["notifyStrategyId"] = v
	}
	if d.HasChange("filter_setting") {
		update = true
	}
	if v := d.Get("filter_setting"); !IsNil(v) || d.HasChange("filter_setting") {
		request["filterSetting"] = expandCmsSubscriptionFilterSetting(v)
	}
	if d.HasChange("pushing_setting") {
		update = true
	}
	if v := d.Get("pushing_setting"); !IsNil(v) || d.HasChange("pushing_setting") {
		request["pushingSetting"] = expandCmsSubscriptionPushingSetting(v)
	}
	if d.HasChange("agent_config") {
		update = true
	}
	if v := d.Get("agent_config"); !IsNil(v) || d.HasChange("agent_config") {
		request["agentConfig"] = expandCmsSubscriptionAgentConfig(v)
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

	return resourceAliCloudCmsSubscriptionRead(d, meta)
}

func resourceAliCloudCmsSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	subscriptionId := d.Id()
	action := fmt.Sprintf("/subscriptions/%s", subscriptionId)
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

func expandCmsSubscriptionFilterSetting(v interface{}) map[string]interface{} {
	filterSetting := make(map[string]interface{})
	filterSettingMap, _ := jsonpath.Get("$[0]", v)
	if filterSettingMap == nil {
		return filterSetting
	}
	filterSettingArg := filterSettingMap.(map[string]interface{})

	if expression, ok := filterSettingArg["expression"]; ok && expression != "" {
		filterSetting["expression"] = expression
	}
	if relation, ok := filterSettingArg["relation"]; ok && relation != "" {
		filterSetting["relation"] = relation
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
		filterSetting["conditions"] = conditionsList
	}

	return filterSetting
}

func expandCmsSubscriptionPushingSetting(v interface{}) map[string]interface{} {
	pushingSetting := make(map[string]interface{})
	pushingSettingMap, _ := jsonpath.Get("$[0]", v)
	if pushingSettingMap == nil {
		return pushingSetting
	}
	pushingSettingArg := pushingSettingMap.(map[string]interface{})

	if responsePlanId, ok := pushingSettingArg["response_plan_id"]; ok && responsePlanId != "" {
		pushingSetting["responsePlanId"] = responsePlanId
	}
	if templateUuid, ok := pushingSettingArg["template_uuid"]; ok && templateUuid != "" {
		pushingSetting["templateUuid"] = templateUuid
	}
	if alertActionIds, ok := pushingSettingArg["alert_action_ids"]; ok {
		pushingSetting["alertActionIds"] = convertToInterfaceArray(alertActionIds)
	}
	if restoreActionIds, ok := pushingSettingArg["restore_action_ids"]; ok {
		pushingSetting["restoreActionIds"] = convertToInterfaceArray(restoreActionIds)
	}

	return pushingSetting
}

func expandCmsSubscriptionAgentConfig(v interface{}) map[string]interface{} {
	agentConfig := make(map[string]interface{})
	agentConfigMap, _ := jsonpath.Get("$[0]", v)
	if agentConfigMap == nil {
		return agentConfig
	}
	agentConfigArg := agentConfigMap.(map[string]interface{})

	if agentUuid, ok := agentConfigArg["agent_uuid"]; ok && agentUuid != "" {
		agentConfig["agentUuid"] = agentUuid
	}
	if routes, ok := agentConfigArg["routes"]; ok {
		routesList := make([]interface{}, 0)
		for _, routesChild := range convertToInterfaceArray(routes) {
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

			routesList = append(routesList, routesMap)
		}
		agentConfig["routes"] = routesList
	}

	return agentConfig
}
