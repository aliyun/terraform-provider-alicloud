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

func dataSourceAlicloudArmsNotificationPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudArmsNotificationPoliciesRead,
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
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repeat": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"repeat_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"send_recover_message": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"escalation_policy_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"integration_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"directed_mode": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"group_rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_wait": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"group_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"grouping_fields": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"matching_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"matching_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"notify_rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notify_start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"notify_end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"notify_channels": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"notify_objects": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"notify_object_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"notify_object_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"notify_object_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"notify_channels": {
													Type:     schema.TypeList,
													Computed: true,
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
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email_title": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"email_content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"email_recover_title": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"email_recover_content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sms_content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sms_recover_content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tts_content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tts_recover_content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"robot_content": {
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

func dataSourceAlicloudArmsNotificationPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListNotificationPolicies"
	request := map[string]interface{}{
		"Page":     1,
		"Size":     PageSizeXLarge,
		"RegionId": client.RegionId,
		"IsDetail": true,
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
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
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_arms_notification_policies", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.PageBean.NotificationPolicies", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PageBean.NotificationPolicies", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["Name"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":    fmt.Sprint(object["Id"]),
			"name":  object["Name"],
			"state": object["State"],
		}
		ids = append(ids, fmt.Sprint(object["Id"]))
		names = append(names, object["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		mapping["repeat"] = object["Repeat"]
		mapping["repeat_interval"] = object["RepeatInterval"]
		mapping["send_recover_message"] = object["SendRecoverMessage"]
		mapping["escalation_policy_id"] = object["EscalationPolicyId"]
		mapping["integration_id"] = object["IntegrationId"]
		mapping["directed_mode"] = object["DirectedMode"]

		if groupRule, ok := object["GroupRule"]; ok && groupRule != nil {
			if groupRuleObj, ok := groupRule.(map[string]interface{}); ok {
				mapping["group_rule"] = []map[string]interface{}{
					{
						"group_wait":      groupRuleObj["GroupWait"],
						"group_interval":  groupRuleObj["GroupInterval"],
						"grouping_fields": groupRuleObj["GroupingFields"],
					},
				}
			}
		}

		if matchingRules, ok := object["MatchingRules"]; ok && matchingRules != nil {
			matchingRulesMaps := make([]map[string]interface{}, 0)
			for _, matchingRule := range matchingRules.([]interface{}) {
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
			mapping["matching_rules"] = matchingRulesMaps
		}

		if notifyRule, ok := object["NotifyRule"]; ok && notifyRule != nil {
			if notifyRuleObj, ok := notifyRule.(map[string]interface{}); ok {
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
				mapping["notify_rule"] = []map[string]interface{}{
					{
						"notify_start_time": notifyRuleObj["NotifyStartTime"],
						"notify_end_time":   notifyRuleObj["NotifyEndTime"],
						"notify_channels":   notifyRuleObj["NotifyChannels"],
						"notify_objects":    notifyObjectsMaps,
					},
				}
			}
		}

		if notifyTemplate, ok := object["NotifyTemplate"]; ok && notifyTemplate != nil {
			if notifyTemplateObj, ok := notifyTemplate.(map[string]interface{}); ok {
				mapping["notify_template"] = []map[string]interface{}{
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
			}
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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
