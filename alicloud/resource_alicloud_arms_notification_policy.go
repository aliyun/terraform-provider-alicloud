// Package alicloud. Hand-written for the ARMS notification_policy resource,
// following the alicloud_arms_webhook_contact precedent on the ARMS 2019-08-08 line.
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudArmsNotificationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudArmsNotificationPolicyCreate,
		Read:   resourceAliCloudArmsNotificationPolicyRead,
		Delete: resourceAliCloudArmsNotificationPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// Top-level scalar fields (form params on the wire).
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"send_recover_message": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"escalation_policy_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"repeat": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"repeat_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"integration_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"directed_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "enable",
				ForceNew: true,
			},
			"notification_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Four JSON-typed form params serialized from nested structs.
			// The published ARMS API models these as formData string params
			// carrying a JSON payload with camelCase nested keys; the
			// CloudSpec model exposes them as typed structs, so the resource
			// serializes the typed struct to the JSON wire format on Create
			// and normalizes the parsed-object or string response on Read.
			"matching_rules": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     armsNpMatchingRuleSchema(),
			},
			"group_rule": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     armsNpGroupRuleSchema(),
			},
			"notify_rule": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     armsNpNotifyRuleSchema(),
			},
			"notify_template": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     armsNpNotifyTemplateSchema(),
			},
		},
	}
}

// matching_rules[*]: { matching_conditions: [{ key, value, operator }] }
func armsNpMatchingRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"matching_conditions": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

// group_rule: { grouping_fields: [string], group_wait, group_interval }
func armsNpGroupRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"grouping_fields": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"group_wait": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"group_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

// notify_rule: { notify_start_time, notify_end_time, notify_channels: [string], notify_objects: [{ notify_object_type, notify_object_id, notify_object_name, notify_channels: [string] }] }
func armsNpNotifyRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"notify_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"notify_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"notify_channels": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_objects": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notify_object_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"notify_object_id": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"notify_object_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"notify_channels": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

// notify_template: email_title, email_content, email_recover_title, email_recover_content,
// sms_content, sms_recover_content, tts_content, tts_recover_content, robot_content
func armsNpNotifyTemplateSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"email_title": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"email_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"email_recover_title": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"email_recover_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sms_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sms_recover_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tts_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tts_recover_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"robot_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudArmsNotificationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateOrUpdateNotificationPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Name"] = d.Get("name")

	if v, ok := d.GetOk("matching_rules"); ok {
		if rulesJson, jErr := armsNpMatchingRulesToJson(v.([]interface{})); jErr == nil {
			request["MatchingRules"] = rulesJson
		}
	}
	if v, ok := d.GetOk("group_rule"); ok {
		if groupJson, jErr := armsNpGroupRuleToJson(v.([]interface{})); jErr == nil {
			request["GroupRule"] = groupJson
		}
	}
	if v, ok := d.GetOk("notify_rule"); ok {
		if ruleJson, jErr := armsNpNotifyRuleToJson(v.([]interface{})); jErr == nil {
			request["NotifyRule"] = ruleJson
		}
	}
	if v, ok := d.GetOk("notify_template"); ok {
		if tmplJson, jErr := armsNpNotifyTemplateToJson(v.([]interface{})); jErr == nil {
			request["NotifyTemplate"] = tmplJson
		}
	}

	if v, ok := d.GetOk("send_recover_message"); ok {
		request["SendRecoverMessage"] = v
	}
	if v, ok := d.GetOk("escalation_policy_id"); ok {
		request["EscalationPolicyId"] = v
	}
	if v, ok := d.GetOk("repeat"); ok {
		request["Repeat"] = v
	}
	if v, ok := d.GetOk("repeat_interval"); ok {
		request["RepeatInterval"] = v
	}
	if v, ok := d.GetOk("integration_id"); ok {
		request["IntegrationId"] = v
	}
	if v, ok := d.GetOk("directed_mode"); ok {
		request["DirectedMode"] = v
	}
	if v, ok := d.GetOk("state"); ok {
		request["State"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)
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

	id, _ := jsonpath.Get("$.NotificationPolicy.Id", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudArmsNotificationPolicyRead(d, meta)
}

func resourceAliCloudArmsNotificationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ArmsServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsNotificationPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_notification_policy DescribeArmsNotificationPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", objectRaw["Name"])
	d.Set("send_recover_message", objectRaw["SendRecoverMessage"])
	d.Set("escalation_policy_id", objectRaw["EscalationPolicyId"])
	d.Set("repeat", objectRaw["Repeat"])
	d.Set("repeat_interval", objectRaw["RepeatInterval"])
	d.Set("integration_id", objectRaw["IntegrationId"])
	d.Set("directed_mode", objectRaw["DirectedMode"])
	d.Set("state", objectRaw["State"])
	d.Set("notification_policy_id", d.Id())

	if v, ok := objectRaw["MatchingRules"]; ok {
		d.Set("matching_rules", armsNpMatchingRulesFromResponse(v))
	}
	if v, ok := objectRaw["GroupRule"]; ok {
		d.Set("group_rule", armsNpGroupRuleFromResponse(v))
	}
	if v, ok := objectRaw["NotifyRule"]; ok {
		d.Set("notify_rule", armsNpNotifyRuleFromResponse(v))
	}
	if v, ok := objectRaw["NotifyTemplate"]; ok {
		d.Set("notify_template", armsNpNotifyTemplateFromResponse(v))
	}

	return nil
}

func resourceAliCloudArmsNotificationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DeleteNotificationPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// === Serialization (TF struct -> JSON wire format with camelCase nested keys) ===

// armsNpMatchingRulesToJson turns a TypeList of {matching_conditions:[{key,value,operator}]}
// into the JSON array wire format: [{"matchingConditions":[{"key":..,"value":..,"operator":..}]},...]
func armsNpMatchingRulesToJson(raw []interface{}) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	out := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		condsRaw, _ := m["matching_conditions"].([]interface{})
		conds := make([]map[string]interface{}, 0, len(condsRaw))
		for _, c := range condsRaw {
			cm, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			conds = append(conds, map[string]interface{}{
				"key":      cm["key"],
				"value":    cm["value"],
				"operator": cm["operator"],
			})
		}
		out = append(out, map[string]interface{}{
			"matchingConditions": conds,
		})
	}
	b, err := json.Marshal(out)
	return string(b), err
}

// armsNpGroupRuleToJson turns a MaxItems:1 TypeList into a JSON object:
// {"groupingFields":[...],"groupWait":N,"groupInterval":N}
func armsNpGroupRuleToJson(raw []interface{}) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected group_rule element type")
	}
	out := map[string]interface{}{}
	if v, ok := m["grouping_fields"].([]interface{}); ok && len(v) > 0 {
		fields := make([]string, 0, len(v))
		for _, f := range v {
			fields = append(fields, fmt.Sprint(f))
		}
		out["groupingFields"] = fields
	}
	if v, ok := m["group_wait"]; ok {
		out["groupWait"] = v
	}
	if v, ok := m["group_interval"]; ok {
		out["groupInterval"] = v
	}
	b, err := json.Marshal(out)
	return string(b), err
}

// armsNpNotifyRuleToJson turns a MaxItems:1 TypeList into a JSON object:
// {"notifyStartTime":"..","notifyEndTime":"..","notifyChannels":[...],"notifyObjects":[{...}]}
func armsNpNotifyRuleToJson(raw []interface{}) (string, error) {
	if len(raw) == 0 {
		return "{}", nil
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected notify_rule element type")
	}
	out := map[string]interface{}{}
	if v, ok := m["notify_start_time"]; ok {
		out["notifyStartTime"] = v
	}
	if v, ok := m["notify_end_time"]; ok {
		out["notifyEndTime"] = v
	}
	if v, ok := m["notify_channels"].([]interface{}); ok && len(v) > 0 {
		chans := make([]string, 0, len(v))
		for _, c := range v {
			chans = append(chans, fmt.Sprint(c))
		}
		out["notifyChannels"] = chans
	}
	if v, ok := m["notify_objects"].([]interface{}); ok && len(v) > 0 {
		objs := make([]map[string]interface{}, 0, len(v))
		for _, o := range v {
			om, ok := o.(map[string]interface{})
			if !ok {
				continue
			}
			obj := map[string]interface{}{}
			if v, ok := om["notify_object_type"]; ok {
				obj["notifyObjectType"] = v
			}
			if v, ok := om["notify_object_id"]; ok {
				obj["notifyObjectId"] = v
			}
			if v, ok := om["notify_object_name"]; ok {
				obj["notifyObjectName"] = v
			}
			if v, ok := om["notify_channels"].([]interface{}); ok && len(v) > 0 {
				oc := make([]string, 0, len(v))
				for _, c := range v {
					oc = append(oc, fmt.Sprint(c))
				}
				obj["notifyChannels"] = oc
			}
			objs = append(objs, obj)
		}
		out["notifyObjects"] = objs
	}
	b, err := json.Marshal(out)
	return string(b), err
}

// armsNpNotifyTemplateToJson turns a MaxItems:1 TypeList into a JSON object
// with camelCase keys (emailTitle, emailContent, ...).
func armsNpNotifyTemplateToJson(raw []interface{}) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected notify_template element type")
	}
	out := map[string]interface{}{}
	for tfKey, wireKey := range map[string]string{
		"email_title":           "emailTitle",
		"email_content":         "emailContent",
		"email_recover_title":   "emailRecoverTitle",
		"email_recover_content": "emailRecoverContent",
		"sms_content":           "smsContent",
		"sms_recover_content":   "smsRecoverContent",
		"tts_content":           "ttsContent",
		"tts_recover_content":   "ttsRecoverContent",
		"robot_content":         "robotContent",
	} {
		if v, ok := m[tfKey]; ok {
			out[wireKey] = v
		}
	}
	b, err := json.Marshal(out)
	return string(b), err
}

// === Deserialization (response -> TF nested list) ===
// The published API may return each JSON field as a parsed object/array
// or as a JSON string; both are normalized to the TF list form.

func armsNpParseJSON(raw interface{}) (interface{}, bool) {
	switch v := raw.(type) {
	case string:
		if v == "" {
			return nil, false
		}
		var parsed interface{}
		if err := json.Unmarshal([]byte(v), &parsed); err != nil {
			return nil, false
		}
		return parsed, true
	case nil:
		return nil, false
	default:
		return raw, true
	}
}

func armsNpMatchingRulesFromResponse(raw interface{}) []map[string]interface{} {
	parsed, ok := armsNpParseJSON(raw)
	if !ok {
		return []map[string]interface{}{}
	}
	list, ok := parsed.([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		conds := []map[string]interface{}{}
		if condsRaw, ok := m["matchingConditions"].([]interface{}); ok {
			for _, c := range condsRaw {
				cm, ok := c.(map[string]interface{})
				if !ok {
					continue
				}
				conds = append(conds, map[string]interface{}{
					"key":      cm["Key"],
					"value":    cm["Value"],
					"operator": cm["Operator"],
				})
			}
		}
		out = append(out, map[string]interface{}{
			"matching_conditions": conds,
		})
	}
	return out
}

func armsNpGroupRuleFromResponse(raw interface{}) []map[string]interface{} {
	parsed, ok := armsNpParseJSON(raw)
	if !ok {
		return []map[string]interface{}{}
	}
	m, ok := parsed.(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	out := map[string]interface{}{}
	if v, ok := m["GroupingFields"].([]interface{}); ok {
		fields := make([]string, 0, len(v))
		for _, f := range v {
			fields = append(fields, fmt.Sprint(f))
		}
		out["grouping_fields"] = fields
	}
	if v, ok := m["GroupWait"]; ok {
		out["group_wait"] = v
	}
	if v, ok := m["GroupInterval"]; ok {
		out["group_interval"] = v
	}
	return []map[string]interface{}{out}
}

func armsNpNotifyRuleFromResponse(raw interface{}) []map[string]interface{} {
	parsed, ok := armsNpParseJSON(raw)
	if !ok {
		return []map[string]interface{}{}
	}
	m, ok := parsed.(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	out := map[string]interface{}{}
	if v, ok := m["NotifyStartTime"]; ok {
		out["notify_start_time"] = v
	}
	if v, ok := m["NotifyEndTime"]; ok {
		out["notify_end_time"] = v
	}
	if v, ok := m["NotifyChannels"].([]interface{}); ok {
		chans := make([]string, 0, len(v))
		for _, c := range v {
			chans = append(chans, fmt.Sprint(c))
		}
		out["notify_channels"] = chans
	}
	if v, ok := m["NotifyObjects"].([]interface{}); ok {
		objs := make([]map[string]interface{}, 0, len(v))
		for _, o := range v {
			om, ok := o.(map[string]interface{})
			if !ok {
				continue
			}
			obj := map[string]interface{}{}
			if v, ok := om["NotifyObjectType"]; ok {
				obj["notify_object_type"] = v
			}
			if v, ok := om["NotifyObjectId"]; ok {
				obj["notify_object_id"] = v
			}
			if v, ok := om["NotifyObjectName"]; ok {
				obj["notify_object_name"] = v
			}
			if v, ok := om["NotifyChannels"].([]interface{}); ok {
				oc := make([]string, 0, len(v))
				for _, c := range v {
					oc = append(oc, fmt.Sprint(c))
				}
				obj["notify_channels"] = oc
			}
			objs = append(objs, obj)
		}
		out["notify_objects"] = objs
	}
	return []map[string]interface{}{out}
}

func armsNpNotifyTemplateFromResponse(raw interface{}) []map[string]interface{} {
	parsed, ok := armsNpParseJSON(raw)
	if !ok {
		return []map[string]interface{}{}
	}
	m, ok := parsed.(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	out := map[string]interface{}{}
	for wireKey, tfKey := range map[string]string{
		"EmailTitle":          "email_title",
		"EmailContent":        "email_content",
		"EmailRecoverTitle":   "email_recover_title",
		"EmailRecoverContent": "email_recover_content",
		"SmsContent":          "sms_content",
		"SmsRecoverContent":   "sms_recover_content",
		"TtsContent":          "tts_content",
		"TtsRecoverContent":   "tts_recover_content",
		"RobotContent":        "robot_content",
	} {
		if v, ok := m[wireKey]; ok {
			out[tfKey] = v
		}
	}
	return []map[string]interface{}{out}
}
