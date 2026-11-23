package arms

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDispatchRule() *schema.Resource {
	return &schema.Resource{
		Create: ResourceDispatchRuleCreate,
		Read:   ResourceDispatchRuleRead,
		Update: ResourceDispatchRuleUpdate,
		Delete: ResourceDispatchRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"is_recover": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"group_rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"group_wait_time": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"group_interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"grouping_fields": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"repeat_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("dispatch_type").(string) == "DISCARD_ALERT"
				},
			},
			"dispatch_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CREATE_ALERT", "DISCARD_ALERT"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label_match_expression_grid": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label_match_expression_groups": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"label_match_expressions": {
										Type:     schema.TypeSet,
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
													Type:     schema.TypeString,
													Required: true,
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
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notify_objects": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notify_object_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"notify_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"ARMS_ROBOT", "ARMS_CONTACT", "ARMS_CONTACT_GROUP"}, false),
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"notify_channels": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"notify_start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notify_end_time": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("dispatch_type").(string) == "DISCARD_ALERT"
				},
			},
			"notify_template": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"email_content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"email_recover_title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"email_recover_content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sms_content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sms_recover_content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tts_content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tts_recover_content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"robot_content": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"dispatch_rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func ResourceDispatchRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDispatchRule"
	request := make(map[string]interface{})
	var err error
	dispatchRuleMap := make(map[string]interface{}, 0)
	if v, ok := d.GetOk("is_recover"); ok {
		dispatchRuleMap["isRecover"] = v
	}

	if v, ok := d.GetOk("group_rules"); ok {
		groupRulesMaps := make([]map[string]interface{}, 0)
		for _, groupRules := range v.(*schema.Set).List() {
			groupRulesArg := groupRules.(map[string]interface{})
			groupRulesMap := map[string]interface{}{
				"groupWait":      groupRulesArg["group_wait_time"],
				"groupInterval":  groupRulesArg["group_interval"],
				"groupingFields": groupRulesArg["grouping_fields"],
				"repeatInterval": groupRulesArg["repeat_interval"],
			}
			groupRulesMaps = append(groupRulesMaps, groupRulesMap)
		}
		dispatchRuleMap["groupRules"] = groupRulesMaps
	}
	if v, ok := d.GetOk("dispatch_type"); ok {
		dispatchRuleMap["dispatchType"] = v
	}
	if v, ok := d.GetOk("label_match_expression_grid"); ok {
		labelMatchExpressionGrid := v.(*schema.Set).List()[0]
		labelMatchExpressionGridArg := labelMatchExpressionGrid.(map[string]interface{})
		labelMatchExpressionGroupsMaps := make([]map[string]interface{}, 0)
		for _, labelMatchExpressionGroups := range labelMatchExpressionGridArg["label_match_expression_groups"].(*schema.Set).List() {
			labelMatchExpressionGroupsArg := labelMatchExpressionGroups.(map[string]interface{})
			labelMatchExpressionsMaps := make([]map[string]interface{}, 0)
			for _, labelMatchExpressions := range labelMatchExpressionGroupsArg["label_match_expressions"].(*schema.Set).List() {
				labelMatchExpressionsArg := labelMatchExpressions.(map[string]interface{})
				labelMatchExpressionsMap := map[string]interface{}{
					"key":      labelMatchExpressionsArg["key"],
					"value":    labelMatchExpressionsArg["value"],
					"operator": labelMatchExpressionsArg["operator"],
				}
				labelMatchExpressionsMaps = append(labelMatchExpressionsMaps, labelMatchExpressionsMap)
			}
			labelMatchExpressionGroupsMaps = append(labelMatchExpressionGroupsMaps, map[string]interface{}{
				"labelMatchExpressions": labelMatchExpressionsMaps,
			})
		}
		dispatchRuleMap["labelMatchExpressionGrid"] = map[string]interface{}{
			"labelMatchExpressionGroups": labelMatchExpressionGroupsMaps,
		}
	}

	if v, ok := d.GetOk("notify_rules"); ok {
		notifyRulesMaps := make([]map[string]interface{}, 0)
		for _, notifyRules := range v.(*schema.Set).List() {
			notifyRulesArg := notifyRules.(map[string]interface{})
			notifyObjectsMaps := make([]map[string]interface{}, 0)
			for _, notifyObjects := range notifyRulesArg["notify_objects"].(*schema.Set).List() {
				notifyObjectsArg := notifyObjects.(map[string]interface{})
				notifyObjectsMap := map[string]interface{}{
					"notifyType":     convertArmsDispatchRuleNotifyTypeRequest(notifyObjectsArg["notify_type"]),
					"name":           notifyObjectsArg["name"],
					"notifyObjectId": notifyObjectsArg["notify_object_id"],
				}
				notifyObjectsMaps = append(notifyObjectsMaps, notifyObjectsMap)
			}
			notifyRulesMap := map[string]interface{}{
				"notifyObjects":   notifyObjectsMaps,
				"notifyChannels":  notifyRulesArg["notify_channels"].([]interface{}),
				"notifyStartTime": notifyRulesArg["notify_start_time"],
				"notifyEndTime":   notifyRulesArg["notify_end_time"],
			}
			notifyRulesMaps = append(notifyRulesMaps, notifyRulesMap)
		}
		dispatchRuleMap["notifyRules"] = notifyRulesMaps
	}

	if v, ok := d.GetOk("notify_template"); ok {
		for _, groupRules := range v.(*schema.Set).List() {
			notifyTemplateArg := groupRules.(map[string]interface{})
			notifyTemplateMap := map[string]interface{}{
				"emailTitle":          notifyTemplateArg["email_title"],
				"emailContent":        notifyTemplateArg["email_content"],
				"emailRecoverTitle":   notifyTemplateArg["email_recover_title"],
				"emailRecoverContent": notifyTemplateArg["email_recover_content"],
				"ttsContent":          notifyTemplateArg["tts_content"],
				"ttsRecoverContent":   notifyTemplateArg["tts_recover_content"],
				"smsContent":          notifyTemplateArg["sms_content"],
				"smsRecoverContent":   notifyTemplateArg["sms_recover_content"],
				"dingContent":         notifyTemplateArg["robot_content"],
			}
			dispatchRuleMap["notifyTemplate"] = notifyTemplateMap
		}
	}

	if v, ok := d.GetOk("dispatch_rule_name"); ok {
		dispatchRuleMap["name"] = v
	}

	request["RegionId"] = client.RegionId
	if v, err := convert.ConvertMaptoJsonString(dispatchRuleMap); err != nil {
		return errs.WrapError(err)
	} else {
		request["DispatchRule"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_dispatch_rule", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DispatchRuleId"]))

	return ResourceDispatchRuleRead(d, meta)
}

func ResourceDispatchRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := Service{client}
	object, err := armsService.DescribeArmsDispatchRule(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_dispatch_rule armsService.DescribeArmsDispatchRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	// DescribeDispatchRule only carries the rule itself: it neither reports the
	// notify time window nor the notify template. Both are read from the
	// notification policy detail, looked up by id so that an account with more
	// rules than one page still resolves the right policy.
	detail, err := NewServiceV2(client).DescribeArmsNotificationPolicy(d.Id())
	if err != nil {
		if !errs.NotFoundError(err) {
			return errs.WrapError(err)
		}
		// The rule exists (DescribeDispatchRule answered), so a missing policy
		// detail is not evidence of deletion: leave the notify attributes alone
		// instead of writing fabricated defaults.
		log.Printf("[DEBUG] Resource alicloud_arms_dispatch_rule notification policy %s is not listed, keeping the known notify configuration: %s", d.Id(), err)
		detail = nil
	}
	if groupRulesList, ok := object["GroupRules"]; ok && groupRulesList != nil {
		groupRulesMaps := make([]map[string]interface{}, 0)
		for _, groupRulesListItem := range groupRulesList.([]interface{}) {
			if groupRulesItemMap, ok := groupRulesListItem.(map[string]interface{}); ok {
				groupRulesMap := make(map[string]interface{}, 0)
				groupRulesMap["group_interval"] = groupRulesItemMap["GroupInterval"]
				groupRulesMap["group_wait_time"] = groupRulesItemMap["GroupWaitTime"]
				groupRulesMap["group_id"] = groupRulesItemMap["GroupId"]
				groupRulesMap["grouping_fields"] = groupRulesItemMap["GroupingFields"]
				groupRulesMap["repeat_interval"] = groupRulesItemMap["RepeatInterval"]
				groupRulesMaps = append(groupRulesMaps, groupRulesMap)
			}
		}
		d.Set("group_rules", groupRulesMaps)
	}
	d.Set("status", object["State"])
	if labelMatchExpressionGrid, ok := object["LabelMatchExpressionGrid"]; ok && labelMatchExpressionGrid != nil {
		labelMatchExpressionGridMaps := make([]map[string]interface{}, 0)
		labelMatchExpressionGridMap := make(map[string]interface{})

		labelMatchExpressionGroupsMaps := make([]map[string]interface{}, 0)
		if v, ok := labelMatchExpressionGrid.(map[string]interface{})["LabelMatchExpressionGroups"]; ok && v != nil {
			for _, labelMatchExpressionGroups := range v.([]interface{}) {
				labelMatchExpressionGroupsMap := make(map[string]interface{})
				labelMatchExpressionsMaps := make([]map[string]interface{}, 0)
				if v, ok := labelMatchExpressionGroups.(map[string]interface{})["LabelMatchExpressions"]; ok && v != nil {
					for _, labelMatchExpressions := range v.([]interface{}) {
						labelMatchExpressionsArg := labelMatchExpressions.(map[string]interface{})
						labelMatchExpressionsMap := make(map[string]interface{}, 0)
						labelMatchExpressionsMap["operator"] = labelMatchExpressionsArg["Operator"]
						labelMatchExpressionsMap["key"] = labelMatchExpressionsArg["Key"]
						labelMatchExpressionsMap["value"] = labelMatchExpressionsArg["Value"]
						labelMatchExpressionsMaps = append(labelMatchExpressionsMaps, labelMatchExpressionsMap)
					}
				}
				labelMatchExpressionGroupsMap["label_match_expressions"] = labelMatchExpressionsMaps
				labelMatchExpressionGroupsMaps = append(labelMatchExpressionGroupsMaps, labelMatchExpressionGroupsMap)
			}
		}
		labelMatchExpressionGridMap["label_match_expression_groups"] = labelMatchExpressionGroupsMaps
		labelMatchExpressionGridMaps = append(labelMatchExpressionGridMaps, labelMatchExpressionGridMap)
		if err := d.Set("label_match_expression_grid", labelMatchExpressionGridMaps); err != nil {
			return errs.WrapError(err)
		}
	}

	if notifyRules := buildArmsDispatchRuleNotifyRules(object, detail); notifyRules != nil {
		if err := d.Set("notify_rules", notifyRules); err != nil {
			return errs.WrapError(err)
		}
	}
	if notifyTemplate := buildArmsDispatchRuleNotifyTemplate(detail); notifyTemplate != nil {
		if err := d.Set("notify_template", []map[string]interface{}{notifyTemplate}); err != nil {
			return errs.WrapError(err)
		}
	}

	d.Set("dispatch_rule_name", object["Name"])
	return nil
}

// buildArmsDispatchRuleNotifyRules merges the DescribeDispatchRule payload with the
// matching notification policy detail into the notify_rules state. Notify objects
// and channels belong to the dispatch rule, while the notify time window is only
// published on the notification policy detail. A nil result means the attribute must
// be left alone: either the payload carries no usable notify rule, or the window is
// not confirmed by the detail. In the latter case Read keeps the existing window
// (e.g. a non-default 08:30-20:15) instead of rewriting it to the all-day default.
func buildArmsDispatchRuleNotifyRules(describe, detail map[string]interface{}) []map[string]interface{} {
	notifyRules, ok := describe["NotifyRules"].([]interface{})
	if !ok || len(notifyRules) == 0 {
		return nil
	}
	notifyStartTime, notifyEndTime, windowConfirmed := armsDispatchRuleNotifyTimeRange(detail)
	if !windowConfirmed {
		log.Printf("[DEBUG] Resource alicloud_arms_dispatch_rule notify window is not confirmed by the notification policy detail, keeping the known notify_rules")
		return nil
	}
	notifyRulesMaps := make([]map[string]interface{}, 0, len(notifyRules))
	for _, item := range notifyRules {
		notifyRule, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		notifyRulesMaps = append(notifyRulesMaps, map[string]interface{}{
			"notify_objects":    buildArmsDispatchRuleNotifyObjects(notifyRule["NotifyObjects"]),
			"notify_channels":   armsDispatchRuleStringList(notifyRule["NotifyChannels"]),
			"notify_start_time": notifyStartTime,
			"notify_end_time":   notifyEndTime,
		})
	}
	if len(notifyRulesMaps) == 0 {
		return nil
	}
	return notifyRulesMaps
}

func buildArmsDispatchRuleNotifyObjects(source interface{}) []map[string]interface{} {
	notifyObjectsMaps := make([]map[string]interface{}, 0)
	notifyObjects, ok := source.([]interface{})
	if !ok {
		return notifyObjectsMaps
	}
	for _, item := range notifyObjects {
		notifyObject, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		notifyObjectsMaps = append(notifyObjectsMaps, map[string]interface{}{
			"notify_type":      convertArmsDispatchRuleNotifyTypeResponse(notifyObject["NotifyType"]),
			"notify_object_id": notifyObject["NotifyObjectId"],
			"name":             notifyObject["Name"],
		})
	}
	return notifyObjectsMaps
}

// armsDispatchRuleNotifyTimeRange reads the notify window from the notification
// policy detail and reports whether it is confirmed. Only a detail that actually
// publishes a NotifyRule object counts: when both window fields are absent or empty
// it is the verified all-day case and normalizes to 00:00/23:59; when both are
// non-empty strings they are used verbatim. A missing detail, an unparseable or
// non-object NotifyRule, a partially published window, or a non-string field is NOT
// confirmed, so the caller must keep the existing state rather than fill business
// defaults.
func armsDispatchRuleNotifyTimeRange(detail map[string]interface{}) (string, string, bool) {
	if detail == nil {
		return "", "", false
	}
	raw, ok := detail["NotifyRule"]
	if !ok || raw == nil {
		return "", "", false
	}
	parsed, ok := armsNpParseJSON(raw)
	if !ok {
		return "", "", false
	}
	notifyRule, ok := parsed.(map[string]interface{})
	if !ok {
		return "", "", false
	}
	start, startOK := armsDispatchRuleWindowField(notifyRule["NotifyStartTime"])
	end, endOK := armsDispatchRuleWindowField(notifyRule["NotifyEndTime"])
	if !startOK || !endOK {
		return "", "", false
	}
	if start == "" && end == "" {
		startTime, endTime := normalizeArmsNotifyTimeRange(start, end)
		return startTime, endTime, true
	}
	if start == "" || end == "" {
		return "", "", false
	}
	return start, end, true
}

// armsDispatchRuleWindowField treats an absent or null field as an empty string and
// rejects any non-string value, so a malformed window is reported as unconfirmed.
func armsDispatchRuleWindowField(value interface{}) (string, bool) {
	if value == nil {
		return "", true
	}
	text, ok := value.(string)
	return text, ok
}

// buildArmsDispatchRuleNotifyTemplate maps the notification policy template onto
// the notify_template schema. The published field may arrive as a parsed object, a
// JSON string or a single item array; an empty or unexpected value means the rule
// has no template to publish and refresh must keep the previous state.
func buildArmsDispatchRuleNotifyTemplate(detail map[string]interface{}) map[string]interface{} {
	raw, ok := detail["NotifyTemplate"]
	if !ok || raw == nil {
		return nil
	}
	parsed, ok := armsNpParseJSON(raw)
	if !ok {
		return nil
	}
	var template map[string]interface{}
	switch value := parsed.(type) {
	case map[string]interface{}:
		template = value
	case []interface{}:
		if len(value) == 0 {
			return nil
		}
		template, _ = value[0].(map[string]interface{})
	default:
		return nil
	}
	if len(template) == 0 {
		return nil
	}
	notifyTemplateMap := map[string]interface{}{
		"email_title":           armsDispatchRuleString(template["EmailTitle"]),
		"email_content":         armsDispatchRuleString(template["EmailContent"]),
		"email_recover_title":   armsDispatchRuleString(template["EmailRecoverTitle"]),
		"email_recover_content": armsDispatchRuleString(template["EmailRecoverContent"]),
		"tts_content":           armsDispatchRuleString(template["TtsContent"]),
		"tts_recover_content":   armsDispatchRuleString(template["TtsRecoverContent"]),
		"sms_content":           armsDispatchRuleString(template["SmsContent"]),
		"sms_recover_content":   armsDispatchRuleString(template["SmsRecoverContent"]),
		"robot_content":         armsDispatchRuleString(template["RobotContent"]),
	}
	if notifyTemplateMap["robot_content"] == "" {
		notifyTemplateMap["robot_content"] = armsDispatchRuleString(template["DingContent"])
	}
	return notifyTemplateMap
}

func armsDispatchRuleString(source interface{}) string {
	if source == nil {
		return ""
	}
	if value, ok := source.(string); ok {
		return value
	}
	return fmt.Sprint(source)
}

func armsDispatchRuleStringList(source interface{}) []string {
	values, ok := source.([]interface{})
	if !ok {
		return []string{}
	}
	list := make([]string, 0, len(values))
	for _, value := range values {
		list = append(list, armsDispatchRuleString(value))
	}
	return list
}

// The ARMS API omits NotifyStartTime and NotifyEndTime for all-day rules.
// Normalize the empty response to the values used when creating the rule so
// that refresh does not produce a perpetual no-op diff.
func normalizeArmsNotifyTimeRange(start, end interface{}) (string, string) {
	startTime, _ := start.(string)
	if startTime == "" {
		startTime = "00:00"
	}
	endTime, _ := end.(string)
	if endTime == "" {
		endTime = "23:59"
	}
	return startTime, endTime
}

func ResourceDispatchRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "UpdateDispatchRule"
	request := make(map[string]interface{})
	var err error
	dispatchRuleMap := make(map[string]interface{}, 0)
	dispatchRuleMap["isRecover"] = d.Get("is_recover")

	dispatchRuleMap["ruleid"] = d.Id()
	dispatchRuleMap["id"] = d.Id()
	if _, ok := d.GetOk("group_rules"); ok {
		oraw, nraw := d.GetChange("group_rules")
		groupRulesMaps := make([]map[string]interface{}, 0)
		groupId := 0
		for _, groupRules := range oraw.(*schema.Set).List() {
			groupRulesArg := groupRules.(map[string]interface{})
			groupId = groupRulesArg["group_id"].(int)
		}
		for _, groupRules := range nraw.(*schema.Set).List() {
			groupRulesArg := groupRules.(map[string]interface{})
			groupRulesMap := map[string]interface{}{
				"groupId":        groupId,
				"groupWait":      groupRulesArg["group_wait_time"],
				"groupInterval":  groupRulesArg["group_interval"],
				"groupingFields": groupRulesArg["grouping_fields"],
				"repeatInterval": groupRulesArg["repeat_interval"],
			}
			groupRulesMaps = append(groupRulesMaps, groupRulesMap)
		}
		dispatchRuleMap["groupRules"] = groupRulesMaps
	}
	if v, ok := d.GetOk("dispatch_type"); ok {
		dispatchRuleMap["dispatchType"] = v
	}
	if v, ok := d.GetOk("label_match_expression_grid"); ok {
		labelMatchExpressionGrid := v.(*schema.Set).List()[0]
		labelMatchExpressionGridArg := labelMatchExpressionGrid.(map[string]interface{})
		labelMatchExpressionGroupsMaps := make([]map[string]interface{}, 0)
		for _, labelMatchExpressionGroups := range labelMatchExpressionGridArg["label_match_expression_groups"].(*schema.Set).List() {
			labelMatchExpressionGroupsArg := labelMatchExpressionGroups.(map[string]interface{})
			labelMatchExpressionsMaps := make([]map[string]interface{}, 0)
			for _, labelMatchExpressions := range labelMatchExpressionGroupsArg["label_match_expressions"].(*schema.Set).List() {
				labelMatchExpressionsArg := labelMatchExpressions.(map[string]interface{})
				labelMatchExpressionsMap := map[string]interface{}{
					"key":      labelMatchExpressionsArg["key"],
					"value":    labelMatchExpressionsArg["value"],
					"operator": labelMatchExpressionsArg["operator"],
				}
				labelMatchExpressionsMaps = append(labelMatchExpressionsMaps, labelMatchExpressionsMap)
			}
			labelMatchExpressionGroupsMaps = append(labelMatchExpressionGroupsMaps, map[string]interface{}{
				"labelMatchExpressions": labelMatchExpressionsMaps,
			})
		}
		dispatchRuleMap["labelMatchExpressionGrid"] = map[string]interface{}{
			"labelMatchExpressionGroups": labelMatchExpressionGroupsMaps,
		}
	}

	if v, ok := d.GetOk("notify_rules"); ok {
		notifyRulesMaps := make([]map[string]interface{}, 0)
		for _, notifyRules := range v.(*schema.Set).List() {
			notifyRulesArg := notifyRules.(map[string]interface{})
			notifyObjectsMaps := make([]map[string]interface{}, 0)
			for _, notifyObjects := range notifyRulesArg["notify_objects"].(*schema.Set).List() {
				notifyObjectsArg := notifyObjects.(map[string]interface{})
				notifyObjectsMap := map[string]interface{}{
					"notifyType":     convertArmsDispatchRuleNotifyTypeRequest(notifyObjectsArg["notify_type"]),
					"name":           notifyObjectsArg["name"],
					"notifyObjectId": notifyObjectsArg["notify_object_id"],
				}
				notifyObjectsMaps = append(notifyObjectsMaps, notifyObjectsMap)
			}
			notifyRulesMap := map[string]interface{}{
				"notifyObjects":   notifyObjectsMaps,
				"notifyChannels":  notifyRulesArg["notify_channels"].([]interface{}),
				"notifyStartTime": notifyRulesArg["notify_start_time"],
				"notifyEndTime":   notifyRulesArg["notify_end_time"],
			}
			notifyRulesMaps = append(notifyRulesMaps, notifyRulesMap)
		}
		dispatchRuleMap["notifyRules"] = notifyRulesMaps
	}

	if v, ok := d.GetOk("notify_template"); ok {
		for _, groupRules := range v.(*schema.Set).List() {
			if notifyTemplateArg, ok := groupRules.(map[string]interface{}); ok {
				notifyTemplateMap := map[string]interface{}{
					"emailTitle":          notifyTemplateArg["email_title"],
					"emailContent":        notifyTemplateArg["email_content"],
					"emailRecoverTitle":   notifyTemplateArg["email_recover_title"],
					"emailRecoverContent": notifyTemplateArg["email_recover_content"],
					"ttsContent":          notifyTemplateArg["tts_content"],
					"ttsRecoverContent":   notifyTemplateArg["tts_recover_content"],
					"smsContent":          notifyTemplateArg["sms_content"],
					"smsRecoverContent":   notifyTemplateArg["sms_recover_content"],
					"dingContent":         notifyTemplateArg["robot_content"],
				}
				dispatchRuleMap["notifyTemplate"] = notifyTemplateMap
			}
		}
	}

	if v, ok := d.GetOk("dispatch_rule_name"); ok {
		dispatchRuleMap["name"] = v
	}

	request["RegionId"] = client.RegionId
	if v, err := convert.ConvertMaptoJsonString(dispatchRuleMap); err != nil {
		return errs.WrapError(err)
	} else {
		request["DispatchRule"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_dispatch_rule", action, errs.AlibabaCloudSdkGoERROR)
	}

	return ResourceDispatchRuleRead(d, meta)
}
func ResourceDispatchRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDispatchRule"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"Id":       d.Id(),
		"RegionId": client.RegionId,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertArmsDispatchRuleNotifyTypeRequest(source interface{}) interface{} {
	switch source {
	case "ARMS_ROBOT":
		return "DING_ROBOT_GROUP"
	}
	return source
}
func convertArmsDispatchRuleNotifyTypeResponse(source interface{}) interface{} {
	switch source {
	case "CONTACT":
		return "ARMS_CONTACT"
	case "DING_ROBOT_GROUP":
		return "ARMS_ROBOT"
	case "CONTACT_GROUP":
		return "ARMS_CONTACT_GROUP"
	}
	return source
}
