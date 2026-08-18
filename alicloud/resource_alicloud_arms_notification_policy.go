package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudArmsNotificationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudArmsNotificationPolicyCreate,
		Read:   resourceAlicloudArmsNotificationPolicyRead,
		Update: resourceAlicloudArmsNotificationPolicyUpdate,
		Delete: resourceAlicloudArmsNotificationPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_rule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_wait": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
						},
						"group_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
						"grouping_fields": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"matching_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching_conditions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"operator": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"eq", "neq", "in", "nin", "re", "nre"}, false),
									},
								},
							},
						},
					},
				},
			},
			"notify_rule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notify_start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notify_end_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notify_channels": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"notify_objects": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notify_object_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"CONTACT", "CONTACT_GROUP", "ARMS_CONTACT", "ARMS_CONTACT_GROUP", "DING_ROBOT_GROUP", "CONTACT_SCHEDULE"}, false),
									},
									"notify_object_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"notify_object_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"notify_channels": {
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
			"notify_template": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_title": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"email_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"email_recover_title": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"email_recover_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sms_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sms_recover_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tts_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tts_recover_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"robot_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"escalation_policy_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"integration_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"repeat": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"repeat_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"send_recover_message": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, false),
			},
			"directed_mode": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudArmsNotificationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateOrUpdateNotificationPolicy"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["Name"] = d.Get("name")

	if v, ok := d.GetOk("group_rule"); ok {
		groupRuleList := v.([]interface{})
		if len(groupRuleList) > 0 {
			groupRuleArg := groupRuleList[0].(map[string]interface{})
			groupRuleMap := map[string]interface{}{
				"groupWait":     groupRuleArg["group_wait"],
				"groupInterval": groupRuleArg["group_interval"],
			}
			if v, ok := groupRuleArg["grouping_fields"]; ok {
				groupRuleMap["groupingFields"] = v
			}
			if jsonStr, err := convertMaptoJsonString(groupRuleMap); err != nil {
				return WrapError(err)
			} else {
				request["GroupRule"] = jsonStr
			}
		}
	}

	if v, ok := d.GetOk("matching_rules"); ok {
		matchingRulesList := v.([]interface{})
		matchingRulesMaps := make([]map[string]interface{}, 0)
		for _, matchingRules := range matchingRulesList {
			matchingRulesArg := matchingRules.(map[string]interface{})
			matchingConditionsMaps := make([]map[string]interface{}, 0)
			for _, condition := range matchingRulesArg["matching_conditions"].([]interface{}) {
				conditionArg := condition.(map[string]interface{})
				matchingConditionsMaps = append(matchingConditionsMaps, map[string]interface{}{
					"key":      conditionArg["key"],
					"value":    conditionArg["value"],
					"operator": conditionArg["operator"],
				})
			}
			matchingRulesMaps = append(matchingRulesMaps, map[string]interface{}{
				"matchingConditions": matchingConditionsMaps,
			})
		}
		if jsonStr, err := convertListMapToJsonString(matchingRulesMaps); err != nil {
			return WrapError(err)
		} else {
			request["MatchingRules"] = jsonStr
		}
	}

	if v, ok := d.GetOk("notify_rule"); ok {
		notifyRuleList := v.([]interface{})
		if len(notifyRuleList) > 0 {
			notifyRuleArg := notifyRuleList[0].(map[string]interface{})
			notifyObjectsMaps := make([]map[string]interface{}, 0)
			for _, notifyObject := range notifyRuleArg["notify_objects"].([]interface{}) {
				notifyObjectArg := notifyObject.(map[string]interface{})
				notifyObjectMap := map[string]interface{}{
					"notifyObjectType": notifyObjectArg["notify_object_type"],
					"notifyObjectId":   notifyObjectArg["notify_object_id"],
					"notifyObjectName": notifyObjectArg["notify_object_name"],
				}
				if v, ok := notifyObjectArg["notify_channels"]; ok {
					notifyObjectMap["notifyChannels"] = v
				}
				notifyObjectsMaps = append(notifyObjectsMaps, notifyObjectMap)
			}
			notifyRuleMap := map[string]interface{}{
				"notifyStartTime": notifyRuleArg["notify_start_time"],
				"notifyEndTime":   notifyRuleArg["notify_end_time"],
				"notifyChannels":  notifyRuleArg["notify_channels"],
				"notifyObjects":   notifyObjectsMaps,
			}
			if jsonStr, err := convertMaptoJsonString(notifyRuleMap); err != nil {
				return WrapError(err)
			} else {
				request["NotifyRule"] = jsonStr
			}
		}
	}

	if v, ok := d.GetOk("notify_template"); ok {
		notifyTemplateList := v.([]interface{})
		if len(notifyTemplateList) > 0 {
			notifyTemplateArg := notifyTemplateList[0].(map[string]interface{})
			notifyTemplateMap := map[string]interface{}{
				"emailTitle":          notifyTemplateArg["email_title"],
				"emailContent":        notifyTemplateArg["email_content"],
				"emailRecoverTitle":   notifyTemplateArg["email_recover_title"],
				"emailRecoverContent": notifyTemplateArg["email_recover_content"],
				"smsContent":          notifyTemplateArg["sms_content"],
				"smsRecoverContent":   notifyTemplateArg["sms_recover_content"],
				"ttsContent":          notifyTemplateArg["tts_content"],
				"ttsRecoverContent":   notifyTemplateArg["tts_recover_content"],
				"robotContent":        notifyTemplateArg["robot_content"],
			}
			if jsonStr, err := convertMaptoJsonString(notifyTemplateMap); err != nil {
				return WrapError(err)
			} else {
				request["NotifyTemplate"] = jsonStr
			}
		}
	}

	if v, ok := d.GetOk("escalation_policy_id"); ok {
		request["EscalationPolicyId"] = v
	}
	if v, ok := d.GetOk("integration_id"); ok {
		request["IntegrationId"] = v
	}
	if v, ok := d.GetOkExists("repeat"); ok {
		request["Repeat"] = v
	}
	if v, ok := d.GetOk("repeat_interval"); ok {
		request["RepeatInterval"] = v
	}
	if v, ok := d.GetOkExists("send_recover_message"); ok {
		request["SendRecoverMessage"] = v
	}
	if v, ok := d.GetOk("state"); ok {
		request["State"] = v
	}
	if v, ok := d.GetOkExists("directed_mode"); ok {
		request["DirectedMode"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_notification_policy", action, AlibabaCloudSdkGoERROR)
	}

	notificationPolicy, ok := response["NotificationPolicy"].(map[string]interface{})
	if !ok {
		return WrapErrorf(fmt.Errorf("NotificationPolicy is not valid"), DefaultErrorMsg, "alicloud_arms_notification_policy", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(notificationPolicy["Id"]))

	return resourceAlicloudArmsNotificationPolicyRead(d, meta)
}

func resourceAlicloudArmsNotificationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := ArmsService{client}
	object, err := armsService.ListArmsNotificationPolicies(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_notification_policy armsService.ListArmsNotificationPolicies Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object["Name"])
	d.Set("state", object["State"])
	d.Set("directed_mode", object["DirectedMode"])
	d.Set("escalation_policy_id", object["EscalationPolicyId"])
	d.Set("integration_id", object["IntegrationId"])
	d.Set("repeat", object["Repeat"])
	d.Set("repeat_interval", object["RepeatInterval"])
	d.Set("send_recover_message", object["SendRecoverMessage"])

	if groupRule, ok := object["GroupRule"]; ok && groupRule != nil {
		groupRuleObj := groupRule.(map[string]interface{})
		if len(groupRuleObj) > 0 {
			groupRuleMaps := []map[string]interface{}{
				{
					"group_wait":      groupRuleObj["GroupWait"],
					"group_interval":  groupRuleObj["GroupInterval"],
					"grouping_fields": groupRuleObj["GroupingFields"],
				},
			}
			d.Set("group_rule", groupRuleMaps)
		}
	}

	if matchingRules, ok := object["MatchingRules"]; ok && matchingRules != nil {
		matchingRulesList := matchingRules.([]interface{})
		matchingRulesMaps := make([]map[string]interface{}, 0)
		for _, matchingRule := range matchingRulesList {
			matchingRuleObj := matchingRule.(map[string]interface{})
			matchingConditionsMaps := make([]map[string]interface{}, 0)
			if conditions, ok := matchingRuleObj["MatchingConditions"]; ok && conditions != nil {
				for _, condition := range conditions.([]interface{}) {
					conditionObj := condition.(map[string]interface{})
					matchingConditionsMaps = append(matchingConditionsMaps, map[string]interface{}{
						"key":      conditionObj["Key"],
						"value":    conditionObj["Value"],
						"operator": conditionObj["Operator"],
					})
				}
			}
			matchingRulesMaps = append(matchingRulesMaps, map[string]interface{}{
				"matching_conditions": matchingConditionsMaps,
			})
		}
		d.Set("matching_rules", matchingRulesMaps)
	}

	if notifyRule, ok := object["NotifyRule"]; ok && notifyRule != nil {
		notifyRuleObj := notifyRule.(map[string]interface{})
		notifyObjectsMaps := make([]map[string]interface{}, 0)
		if notifyObjects, ok := notifyRuleObj["NotifyObjects"]; ok && notifyObjects != nil {
			for _, notifyObject := range notifyObjects.([]interface{}) {
				notifyObjectObj := notifyObject.(map[string]interface{})
				notifyObjectMap := map[string]interface{}{
					"notify_object_type": notifyObjectObj["NotifyObjectType"],
					"notify_object_id":   fmt.Sprint(notifyObjectObj["NotifyObjectId"]),
					"notify_object_name": notifyObjectObj["NotifyObjectName"],
				}
				if v, ok := notifyObjectObj["NotifyChannels"]; ok {
					notifyObjectMap["notify_channels"] = v
				}
				notifyObjectsMaps = append(notifyObjectsMaps, notifyObjectMap)
			}
		}
		notifyRuleMaps := []map[string]interface{}{
			{
				"notify_start_time": notifyRuleObj["NotifyStartTime"],
				"notify_end_time":   notifyRuleObj["NotifyEndTime"],
				"notify_channels":   notifyRuleObj["NotifyChannels"],
				"notify_objects":    notifyObjectsMaps,
			},
		}
		d.Set("notify_rule", notifyRuleMaps)
	}

	if notifyTemplate, ok := object["NotifyTemplate"]; ok && notifyTemplate != nil {
		notifyTemplateObj := notifyTemplate.(map[string]interface{})
		if len(notifyTemplateObj) > 0 {
			notifyTemplateMaps := []map[string]interface{}{
				{
					"email_title":           notifyTemplateObj["EmailTitle"],
					"email_content":         notifyTemplateObj["EmailContent"],
					"email_recover_title":   notifyTemplateObj["EmailRecoverTitle"],
					"email_recover_content": notifyTemplateObj["EmailRecoverContent"],
					"sms_content":           notifyTemplateObj["SmsContent"],
					"sms_recover_content":   notifyTemplateObj["SmsRecoverContent"],
					"tts_content":           notifyTemplateObj["TtsContent"],
					"tts_recover_content":   notifyTemplateObj["TtsRecoverContent"],
					"robot_content":         notifyTemplateObj["RobotContent"],
				},
			}
			d.Set("notify_template", notifyTemplateMaps)
		}
	}

	return nil
}

func resourceAlicloudArmsNotificationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateOrUpdateNotificationPolicy"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["Id"] = d.Id()
	request["Name"] = d.Get("name")

	if v, ok := d.GetOk("group_rule"); ok {
		groupRuleList := v.([]interface{})
		if len(groupRuleList) > 0 {
			groupRuleArg := groupRuleList[0].(map[string]interface{})
			groupRuleMap := map[string]interface{}{
				"groupWait":     groupRuleArg["group_wait"],
				"groupInterval": groupRuleArg["group_interval"],
			}
			if v, ok := groupRuleArg["grouping_fields"]; ok {
				groupRuleMap["groupingFields"] = v
			}
			if jsonStr, err := convertMaptoJsonString(groupRuleMap); err != nil {
				return WrapError(err)
			} else {
				request["GroupRule"] = jsonStr
			}
		}
	}

	if v, ok := d.GetOk("matching_rules"); ok {
		matchingRulesList := v.([]interface{})
		matchingRulesMaps := make([]map[string]interface{}, 0)
		for _, matchingRules := range matchingRulesList {
			matchingRulesArg := matchingRules.(map[string]interface{})
			matchingConditionsMaps := make([]map[string]interface{}, 0)
			for _, condition := range matchingRulesArg["matching_conditions"].([]interface{}) {
				conditionArg := condition.(map[string]interface{})
				matchingConditionsMaps = append(matchingConditionsMaps, map[string]interface{}{
					"key":      conditionArg["key"],
					"value":    conditionArg["value"],
					"operator": conditionArg["operator"],
				})
			}
			matchingRulesMaps = append(matchingRulesMaps, map[string]interface{}{
				"matchingConditions": matchingConditionsMaps,
			})
		}
		if jsonStr, err := convertListMapToJsonString(matchingRulesMaps); err != nil {
			return WrapError(err)
		} else {
			request["MatchingRules"] = jsonStr
		}
	}

	if v, ok := d.GetOk("notify_rule"); ok {
		notifyRuleList := v.([]interface{})
		if len(notifyRuleList) > 0 {
			notifyRuleArg := notifyRuleList[0].(map[string]interface{})
			notifyObjectsMaps := make([]map[string]interface{}, 0)
			for _, notifyObject := range notifyRuleArg["notify_objects"].([]interface{}) {
				notifyObjectArg := notifyObject.(map[string]interface{})
				notifyObjectMap := map[string]interface{}{
					"notifyObjectType": notifyObjectArg["notify_object_type"],
					"notifyObjectId":   notifyObjectArg["notify_object_id"],
					"notifyObjectName": notifyObjectArg["notify_object_name"],
				}
				if v, ok := notifyObjectArg["notify_channels"]; ok {
					notifyObjectMap["notifyChannels"] = v
				}
				notifyObjectsMaps = append(notifyObjectsMaps, notifyObjectMap)
			}
			notifyRuleMap := map[string]interface{}{
				"notifyStartTime": notifyRuleArg["notify_start_time"],
				"notifyEndTime":   notifyRuleArg["notify_end_time"],
				"notifyChannels":  notifyRuleArg["notify_channels"],
				"notifyObjects":   notifyObjectsMaps,
			}
			if jsonStr, err := convertMaptoJsonString(notifyRuleMap); err != nil {
				return WrapError(err)
			} else {
				request["NotifyRule"] = jsonStr
			}
		}
	}

	if v, ok := d.GetOk("notify_template"); ok {
		notifyTemplateList := v.([]interface{})
		if len(notifyTemplateList) > 0 {
			notifyTemplateArg := notifyTemplateList[0].(map[string]interface{})
			notifyTemplateMap := map[string]interface{}{
				"emailTitle":          notifyTemplateArg["email_title"],
				"emailContent":        notifyTemplateArg["email_content"],
				"emailRecoverTitle":   notifyTemplateArg["email_recover_title"],
				"emailRecoverContent": notifyTemplateArg["email_recover_content"],
				"smsContent":          notifyTemplateArg["sms_content"],
				"smsRecoverContent":   notifyTemplateArg["sms_recover_content"],
				"ttsContent":          notifyTemplateArg["tts_content"],
				"ttsRecoverContent":   notifyTemplateArg["tts_recover_content"],
				"robotContent":        notifyTemplateArg["robot_content"],
			}
			if jsonStr, err := convertMaptoJsonString(notifyTemplateMap); err != nil {
				return WrapError(err)
			} else {
				request["NotifyTemplate"] = jsonStr
			}
		}
	}

	if v, ok := d.GetOk("escalation_policy_id"); ok {
		request["EscalationPolicyId"] = v
	}
	if v, ok := d.GetOk("integration_id"); ok {
		request["IntegrationId"] = v
	}
	if v, ok := d.GetOkExists("repeat"); ok {
		request["Repeat"] = v
	}
	if v, ok := d.GetOk("repeat_interval"); ok {
		request["RepeatInterval"] = v
	}
	if v, ok := d.GetOkExists("send_recover_message"); ok {
		request["SendRecoverMessage"] = v
	}
	if v, ok := d.GetOk("state"); ok {
		request["State"] = v
	}
	if v, ok := d.GetOkExists("directed_mode"); ok {
		request["DirectedMode"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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

	return resourceAlicloudArmsNotificationPolicyRead(d, meta)
}

func resourceAlicloudArmsNotificationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteNotificationPolicy"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"Id":       d.Id(),
		"RegionId": client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
	return nil
}
