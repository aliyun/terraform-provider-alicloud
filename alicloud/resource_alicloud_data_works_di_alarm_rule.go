// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDataWorksDiAlarmRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataWorksDiAlarmRuleCreate,
		Read:   resourceAliCloudDataWorksDiAlarmRuleRead,
		Update: resourceAliCloudDataWorksDiAlarmRuleUpdate,
		Delete: resourceAliCloudDataWorksDiAlarmRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"di_alarm_rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"di_alarm_rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"di_job_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"metric_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notification_settings": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notification_channels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"channels": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"inhibition_interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"notification_receivers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"receiver_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"receiver_values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"trigger_conditions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ddl_report_tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"severity": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudDataWorksDiAlarmRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDIAlarmRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["DIJobId"] = d.Get("di_job_id")
	query["RegionId"] = client.RegionId
	query["ClientToken"] = StringPointer(buildClientToken(action))

	if v, ok := d.GetOk("description"); ok {
		query["Description"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("metric_type"); ok {
		query["MetricType"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("enabled"); ok {
		query["Enabled"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOk("trigger_conditions"); ok {
		triggerConditionsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Duration"] = dataLoopTmp["duration"]
			dataLoopMap["Threshold"] = dataLoopTmp["threshold"]
			dataLoopMap["Severity"] = dataLoopTmp["severity"]
			dataLoopMap["DdlReportTags"] = dataLoopTmp["ddl_report_tags"]
			triggerConditionsMapsArray = append(triggerConditionsMapsArray, dataLoopMap)
		}
		triggerConditionsMapsJson, err := json.Marshal(triggerConditionsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		query["TriggerConditions"] = string(triggerConditionsMapsJson)
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("notification_settings"); v != nil {
		if v, ok := d.GetOk("notification_settings"); ok {
			localData1, err := jsonpath.Get("$[0].notification_channels", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["Channels"] = dataLoop1Tmp["channels"]
				dataLoop1Map["Severity"] = dataLoop1Tmp["severity"]
				localMaps = append(localMaps, dataLoop1Map)
			}
			objectDataLocalMap["NotificationChannels"] = localMaps
		}

		if v, ok := d.GetOk("notification_settings"); ok {
			localData2, err := jsonpath.Get("$[0].notification_receivers", v)
			if err != nil {
				localData2 = make([]interface{}, 0)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop2 := range localData2.([]interface{}) {
				dataLoop2Tmp := make(map[string]interface{})
				if dataLoop2 != nil {
					dataLoop2Tmp = dataLoop2.(map[string]interface{})
				}
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["ReceiverType"] = dataLoop2Tmp["receiver_type"]
				dataLoop2Map["ReceiverValues"] = dataLoop2Tmp["receiver_values"]
				localMaps1 = append(localMaps1, dataLoop2Map)
			}
			objectDataLocalMap["NotificationReceivers"] = localMaps1
		}

		inhibitionInterval1, _ := jsonpath.Get("$[0].inhibition_interval", v)
		if inhibitionInterval1 != nil && inhibitionInterval1 != "" {
			objectDataLocalMap["InhibitionInterval"] = inhibitionInterval1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		query["NotificationSettings"] = string(objectDataLocalMapJson)
	}

	if v, ok := d.GetOk("di_alarm_rule_name"); ok {
		query["Name"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_di_alarm_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DIJobId"], response["DIAlarmRuleId"]))

	return resourceAliCloudDataWorksDiAlarmRuleRead(d, meta)
}

func resourceAliCloudDataWorksDiAlarmRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksDiAlarmRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_di_alarm_rule DescribeDataWorksDiAlarmRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["Name"] != nil {
		d.Set("di_alarm_rule_name", objectRaw["Name"])
	}
	if objectRaw["Enabled"] != nil {
		d.Set("enabled", objectRaw["Enabled"])
	}
	if objectRaw["MetricType"] != nil {
		d.Set("metric_type", objectRaw["MetricType"])
	}
	if objectRaw["DIAlarmRuleId"] != nil {
		d.Set("di_alarm_rule_id", objectRaw["DIAlarmRuleId"])
	}
	if objectRaw["DIJobId"] != nil {
		d.Set("di_job_id", objectRaw["DIJobId"])
	}

	notificationSettingsMaps := make([]map[string]interface{}, 0)
	notificationSettingsMap := make(map[string]interface{})
	notificationSettings1Raw := make(map[string]interface{})
	if objectRaw["NotificationSettings"] != nil {
		notificationSettings1Raw = objectRaw["NotificationSettings"].(map[string]interface{})
	}
	if len(notificationSettings1Raw) > 0 {
		notificationSettingsMap["inhibition_interval"] = notificationSettings1Raw["InhibitionInterval"]

		notificationChannels1Raw := notificationSettings1Raw["NotificationChannels"]
		notificationChannelsMaps := make([]map[string]interface{}, 0)
		if notificationChannels1Raw != nil {
			for _, notificationChannelsChild1Raw := range notificationChannels1Raw.([]interface{}) {
				notificationChannelsMap := make(map[string]interface{})
				notificationChannelsChild1Raw := notificationChannelsChild1Raw.(map[string]interface{})
				notificationChannelsMap["severity"] = notificationChannelsChild1Raw["Severity"]

				channels1Raw := make([]interface{}, 0)
				if notificationChannelsChild1Raw["Channels"] != nil {
					channels1Raw = notificationChannelsChild1Raw["Channels"].([]interface{})
				}

				notificationChannelsMap["channels"] = channels1Raw
				notificationChannelsMaps = append(notificationChannelsMaps, notificationChannelsMap)
			}
		}
		notificationSettingsMap["notification_channels"] = notificationChannelsMaps
		notificationReceivers1Raw := notificationSettings1Raw["NotificationReceivers"]
		notificationReceiversMaps := make([]map[string]interface{}, 0)
		if notificationReceivers1Raw != nil {
			for _, notificationReceiversChild1Raw := range notificationReceivers1Raw.([]interface{}) {
				notificationReceiversMap := make(map[string]interface{})
				notificationReceiversChild1Raw := notificationReceiversChild1Raw.(map[string]interface{})
				notificationReceiversMap["receiver_type"] = notificationReceiversChild1Raw["ReceiverType"]

				receiverValues1Raw := make([]interface{}, 0)
				if notificationReceiversChild1Raw["ReceiverValues"] != nil {
					receiverValues1Raw = notificationReceiversChild1Raw["ReceiverValues"].([]interface{})
				}

				notificationReceiversMap["receiver_values"] = receiverValues1Raw
				notificationReceiversMaps = append(notificationReceiversMaps, notificationReceiversMap)
			}
		}
		notificationSettingsMap["notification_receivers"] = notificationReceiversMaps
		notificationSettingsMaps = append(notificationSettingsMaps, notificationSettingsMap)
	}
	if objectRaw["NotificationSettings"] != nil {
		if err := d.Set("notification_settings", notificationSettingsMaps); err != nil {
			return err
		}
	}
	triggerConditions1Raw := objectRaw["TriggerConditions"]
	triggerConditionsMaps := make([]map[string]interface{}, 0)
	if triggerConditions1Raw != nil {
		for _, triggerConditionsChild1Raw := range triggerConditions1Raw.([]interface{}) {
			triggerConditionsMap := make(map[string]interface{})
			triggerConditionsChild1Raw := triggerConditionsChild1Raw.(map[string]interface{})
			triggerConditionsMap["duration"] = triggerConditionsChild1Raw["Duration"]
			triggerConditionsMap["severity"] = triggerConditionsChild1Raw["Severity"]
			triggerConditionsMap["threshold"] = triggerConditionsChild1Raw["Threshold"]

			ddlReportTags1Raw := make([]interface{}, 0)
			if triggerConditionsChild1Raw["DdlReportTags"] != nil {
				ddlReportTags1Raw = triggerConditionsChild1Raw["DdlReportTags"].([]interface{})
			}

			triggerConditionsMap["ddl_report_tags"] = ddlReportTags1Raw
			triggerConditionsMaps = append(triggerConditionsMaps, triggerConditionsMap)
		}
	}
	if objectRaw["TriggerConditions"] != nil {
		if err := d.Set("trigger_conditions", triggerConditionsMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceAliCloudDataWorksDiAlarmRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	parts := strings.Split(d.Id(), ":")
	action := "UpdateDIAlarmRule"
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DIAlarmRuleId"] = parts[1]
	query["DIJobId"] = parts[0]
	query["RegionId"] = client.RegionId
	query["ClientToken"] = buildClientToken(action)
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			query["Description"] = StringPointer(v.(string))
		}
	}

	if d.HasChange("metric_type") {
		update = true
	}
	if v, ok := d.GetOk("metric_type"); ok {
		query["MetricType"] = StringPointer(v.(string))
	}

	if d.HasChange("enabled") {
		update = true
		if v, ok := d.GetOk("enabled"); ok {
			query["Enabled"] = StringPointer(strconv.FormatBool(v.(bool)))
		}
	}

	if d.HasChange("trigger_conditions") {
		update = true
	}
	if v, ok := d.GetOk("trigger_conditions"); ok || d.HasChange("trigger_conditions") {
		triggerConditionsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Duration"] = dataLoopTmp["duration"]
			dataLoopMap["Threshold"] = dataLoopTmp["threshold"]
			dataLoopMap["DdlReportTags"] = dataLoopTmp["ddl_report_tags"]
			dataLoopMap["Severity"] = dataLoopTmp["severity"]
			triggerConditionsMapsArray = append(triggerConditionsMapsArray, dataLoopMap)
		}
		triggerConditionsMapsJson, err := json.Marshal(triggerConditionsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		query["TriggerConditions"] = string(triggerConditionsMapsJson)
	}

	if d.HasChange("notification_settings") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("notification_settings"); v != nil {
		if v, ok := d.GetOk("notification_settings"); ok {
			localData1, err := jsonpath.Get("$[0].notification_channels", v)
			if err != nil {
				localData1 = make([]interface{}, 0)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := make(map[string]interface{})
				if dataLoop1 != nil {
					dataLoop1Tmp = dataLoop1.(map[string]interface{})
				}
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["Channels"] = dataLoop1Tmp["channels"]
				dataLoop1Map["Severity"] = dataLoop1Tmp["severity"]
				localMaps = append(localMaps, dataLoop1Map)
			}
			objectDataLocalMap["NotificationChannels"] = localMaps
		}

		if v, ok := d.GetOk("notification_settings"); ok {
			localData2, err := jsonpath.Get("$[0].notification_receivers", v)
			if err != nil {
				localData2 = make([]interface{}, 0)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop2 := range localData2.([]interface{}) {
				dataLoop2Tmp := make(map[string]interface{})
				if dataLoop2 != nil {
					dataLoop2Tmp = dataLoop2.(map[string]interface{})
				}
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["ReceiverType"] = dataLoop2Tmp["receiver_type"]
				dataLoop2Map["ReceiverValues"] = dataLoop2Tmp["receiver_values"]
				localMaps1 = append(localMaps1, dataLoop2Map)
			}
			objectDataLocalMap["NotificationReceivers"] = localMaps1
		}

		inhibitionInterval1, _ := jsonpath.Get("$[0].inhibition_interval", v)
		if inhibitionInterval1 != nil && (d.HasChange("notification_settings.0.inhibition_interval") || inhibitionInterval1 != "") {
			objectDataLocalMap["InhibitionInterval"] = inhibitionInterval1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		query["NotificationSettings"] = string(objectDataLocalMapJson)
	}

	if d.HasChange("di_alarm_rule_name") {
		update = true
	}
	if v, ok := d.GetOk("di_alarm_rule_name"); ok {
		query["Name"] = StringPointer(v.(string))
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)
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

	return resourceAliCloudDataWorksDiAlarmRuleRead(d, meta)
}

func resourceAliCloudDataWorksDiAlarmRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDIAlarmRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["DIAlarmRuleId"] = parts[1]
	query["DIJobId"] = parts[0]
	query["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)

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
		if IsExpectedErrors(err, []string{"500130"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
