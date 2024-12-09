package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudArmsDispatchRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudArmsDispatchRulesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"dispatch_rule_name": {
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
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dispatch_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_wait_time": {
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
									"repeat_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label_match_expression_grid": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"label_match_expression_groups": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"label_match_expressions": {
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
								},
							},
						},
						"notify_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notify_objects": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"notify_object_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"notify_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"notify_channels": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"notify_start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"notify_end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dispatch_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
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
									"tts_content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tts_recover_content": {
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

func dataSourceAlicloudArmsDispatchRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListNotificationPolicies"
	request := map[string]interface{}{
		"Page":     1,
		"Size":     PageSizeXLarge,
		"RegionId": client.RegionId,
		"IsDetail": true,
	}
	if v, ok := d.GetOk("dispatch_rule_name"); ok {
		request["Name"] = v
	}
	var objects []map[string]interface{}
	var dispatchRuleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		dispatchRuleNameRegex = r
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
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_arms_dispatch_rules", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.PageBean.NotificationPolicies", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PageBean.NotificationPolicies", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if dispatchRuleNameRegex != nil && !dispatchRuleNameRegex.MatchString(fmt.Sprint(item["Name"])) {
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
			"id":                 fmt.Sprint(object["Id"]),
			"dispatch_rule_id":   fmt.Sprint(object["Id"]),
			"dispatch_rule_name": object["Name"],
			"status":             object["State"],
		}
		ids = append(ids, fmt.Sprint(object["Id"]))
		names = append(names, object["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		if groupRule, ok := object["GroupRule"]; ok && groupRule != nil {
			groupRuleMaps := make([]map[string]interface{}, 0)
			if groupRuleItemMap, ok := groupRule.(map[string]interface{}); ok {
				groupRuleMap := make(map[string]interface{}, 0)
				groupRuleMap["group_interval"] = groupRuleItemMap["GroupInterval"]
				groupRuleMap["group_wait_time"] = groupRuleItemMap["GroupWait"]
				groupRuleMap["grouping_fields"] = groupRuleItemMap["GroupingFields"]
				groupRuleMap["repeat_interval"] = object["RepeatInterval"]
				groupRuleMaps = append(groupRuleMaps, groupRuleMap)
			}
			mapping["group_rules"] = groupRuleMaps
		}
		if matchingRulesList, ok := object["MatchingRules"]; ok && matchingRulesList != nil {
			matchingRulesMap := make(map[string]interface{}, 0)
			matchingRulesMaps := make([]map[string]interface{}, 0)
			for _, matchingRulesListItem := range matchingRulesList.([]interface{}) {
				matchingRulesListItemArg := matchingRulesListItem.(map[string]interface{})
				if matchingConditionsList, ok := matchingRulesListItemArg["MatchingConditions"]; ok && matchingConditionsList != nil {
					matchingConditionsMaps := make([]map[string]interface{}, 0)
					for _, matchingConditionsListItem := range matchingConditionsList.([]interface{}) {
						matchingConditionsMap := make(map[string]interface{}, 0)
						matchingConditionsArg := matchingConditionsListItem.(map[string]interface{})
						matchingConditionsMap["operator"] = matchingConditionsArg["Operator"]
						matchingConditionsMap["key"] = matchingConditionsArg["Key"]
						matchingConditionsMap["value"] = matchingConditionsArg["Value"]
						matchingConditionsMaps = append(matchingConditionsMaps, matchingConditionsMap)
					}
					matchingRulesMap["label_match_expressions"] = matchingConditionsMaps
					matchingRulesMaps = append(matchingRulesMaps, matchingRulesMap)
				}
			}
			mapping["label_match_expression_grid"] = []map[string]interface{}{
				{"label_match_expression_groups": matchingRulesMaps},
			}
		}
		if notifyTemplateItem, ok := object["NotifyTemplate"]; ok && notifyTemplateItem != nil {
			notifyTemplateMaps := make([]map[string]interface{}, 0)
			if notifyTemplateItemMap, ok := notifyTemplateItem.(map[string]interface{}); ok {
				notifyTemplateMap := make(map[string]interface{}, 0)
				notifyTemplateMap["email_title"] = notifyTemplateItemMap["EmailTitle"]
				notifyTemplateMap["email_content"] = notifyTemplateItemMap["EmailContent"]
				notifyTemplateMap["email_recover_title"] = notifyTemplateItemMap["EmailRecoverContent"]
				notifyTemplateMap["email_recover_content"] = notifyTemplateItemMap["EmailRecoverTitle"]
				notifyTemplateMap["sms_content"] = notifyTemplateItemMap["SmsContent"]
				notifyTemplateMap["sms_recover_content"] = notifyTemplateItemMap["SmsRecoverContent"]
				notifyTemplateMap["tts_content"] = notifyTemplateItemMap["TtsContent"]
				notifyTemplateMap["tts_recover_content"] = notifyTemplateItemMap["TtsRecoverContent"]
				notifyTemplateMap["robot_content"] = notifyTemplateItemMap["RobotContent"]
				notifyTemplateMaps = append(notifyTemplateMaps, notifyTemplateMap)
			}
			mapping["notify_template"] = notifyTemplateMaps
		}
		if notifyRuleItem, ok := object["NotifyRule"]; ok && notifyRuleItem != nil {
			notifyRulesMaps := make([]map[string]interface{}, 0)
			if notifyRuleItemMap, ok := notifyRuleItem.(map[string]interface{}); ok {
				notifyRulesMap := make(map[string]interface{}, 0)
				notifyObjectsMaps := make([]map[string]interface{}, 0)
				for _, notifyObjects := range notifyRuleItemMap["NotifyObjects"].([]interface{}) {
					notifyObjectsArg := notifyObjects.(map[string]interface{})
					notifyObjectsMap := make(map[string]interface{}, 0)
					notifyObjectsMap["notify_type"] = convertArmsDispatchRuleNotifyTypeResponse(notifyObjectsArg["NotifyObjectType"])
					notifyObjectsMap["notify_object_id"] = notifyObjectsArg["NotifyObjectId"]
					notifyObjectsMap["name"] = notifyObjectsArg["NotifyObjectName"]
					notifyObjectsMaps = append(notifyObjectsMaps, notifyObjectsMap)
				}
				notifyRulesMap["notify_objects"] = notifyObjectsMaps
				notifyRulesMap["notify_channels"] = notifyRuleItemMap["NotifyChannels"]
				notifyRulesMap["notify_start_time"] = notifyRuleItemMap["NotifyStartTime"]
				notifyRulesMap["notify_end_time"] = notifyRuleItemMap["NotifyEndTime"]
				notifyRulesMaps = append(notifyRulesMaps, notifyRulesMap)
			}
			mapping["notify_rules"] = notifyRulesMaps
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

	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
