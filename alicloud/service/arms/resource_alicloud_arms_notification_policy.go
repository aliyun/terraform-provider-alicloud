// Hand-written ARMS notification_policy resource on the ARMS 2019-08-08 API
// line (raw RpcPost for CreateOrUpdateNotificationPolicy/DeleteNotificationPolicy,
// RpcGet for ListNotificationPolicies).
package arms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceNotificationPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.FromErr(ResourceNotificationPolicyCreateContext(ctx, d, meta))
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.FromErr(ResourceNotificationPolicyReadContext(ctx, d, meta))
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.FromErr(ResourceNotificationPolicyDeleteContext(ctx, d, meta))
		},
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
				Computed: true,
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
				Computed: true,
				ForceNew: true,
			},
			"directed_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "enable",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, false),
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
				Computed: true,
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
				Computed: true,
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
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"eq", "neq", "in", "nin", "re", "nre"}, false),
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
				Default:  5,
				ForceNew: true,
			},
			"group_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
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
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"dingTalk", "email", "sms", "tts", "webhook"}, false),
				},
			},
			"notify_objects": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notify_object_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"CONTACT", "CONTACT_GROUP", "ARMS_CONTACT", "ARMS_CONTACT_GROUP", "DING_ROBOT_GROUP", "CONTACT_SCHEDULE"}, false),
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
							Computed: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"email", "sms", "tts"}, false),
							},
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

func ResourceNotificationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	return ResourceNotificationPolicyCreateContext(context.Background(), d, meta)
}

func ResourceNotificationPolicyCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request, err := armsNpCreateRequest(d, d.GetRawConfig(), client.RegionId)
	if err != nil {
		return err
	}
	return armsNpCreateWithRequest(ctx, d, meta, request)
}

// ArmsNpConfiguredBool reports whether key is present in the raw configuration
// with any value, including an explicit false. It returns false when the raw
// configuration is unavailable, which is the case for ResourceData built
// directly in unit tests.
func ArmsNpConfiguredBool(rawConfig cty.Value, key string) bool {
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return false
	}
	ty := rawConfig.Type()
	if !ty.IsObjectType() || !ty.HasAttribute(key) {
		return false
	}
	value := rawConfig.GetAttr(key)
	return value.IsKnown() && !value.IsNull()
}

// armsNpCreateRequest builds the CreateOrUpdateNotificationPolicy form body.
func armsNpCreateRequest(d *schema.ResourceData, rawConfig cty.Value, regionId string) (map[string]interface{}, error) {
	request := make(map[string]interface{})

	if d.Get("repeat").(bool) && d.Get("repeat_interval").(int) == 0 {
		return nil, fmt.Errorf("repeat_interval is required when repeat is true")
	}
	if !d.Get("repeat").(bool) && d.Get("escalation_policy_id").(int) == 0 {
		return nil, fmt.Errorf("escalation_policy_id is required when repeat is false")
	}
	request["RegionId"] = regionId
	request["Name"] = d.Get("name")

	if v, ok := d.GetOk("matching_rules"); ok {
		rulesJson, jErr := armsNpMatchingRulesToJson(v.([]interface{}))
		if jErr != nil {
			return nil, errs.WrapError(jErr)
		}
		request["MatchingRules"] = rulesJson
	}
	if v, ok := d.GetOk("group_rule"); ok {
		groupJson, jErr := ArmsNpGroupRuleToJson(v.([]interface{}))
		if jErr != nil {
			return nil, errs.WrapError(jErr)
		}
		request["GroupRule"] = groupJson
	}
	if v, ok := d.GetOk("notify_rule"); ok {
		ruleJson, jErr := armsNpNotifyRuleToJson(v.([]interface{}))
		if jErr != nil {
			return nil, errs.WrapError(jErr)
		}
		request["NotifyRule"] = ruleJson
	}
	if v, ok := d.GetOk("notify_template"); ok {
		tmplJson, jErr := armsNpNotifyTemplateToJson(v.([]interface{}))
		if jErr != nil {
			return nil, errs.WrapError(jErr)
		}
		request["NotifyTemplate"] = tmplJson
	}

	// GetOk reports ok=false for zero values, so a configured `false` on
	// these Default:true bools would never reach the API; always send the
	// materialized value.
	request["SendRecoverMessage"] = d.Get("send_recover_message")
	request["Repeat"] = d.Get("repeat")
	if v, ok := d.GetOk("repeat_interval"); ok {
		request["RepeatInterval"] = v
	}
	if v, ok := d.GetOk("escalation_policy_id"); ok {
		request["EscalationPolicyId"] = v
	}
	if v, ok := d.GetOk("integration_id"); ok {
		request["IntegrationId"] = v
	}
	// directed_mode is Optional+Computed, so GetOk alone cannot tell an omitted
	// argument (keep the service default) from an explicit false (must be sent).
	// GetRawConfig is the supported way to distinguish those two cases.
	if v, ok := d.GetOk("directed_mode"); ok {
		request["DirectedMode"] = v
	} else if ArmsNpConfiguredBool(rawConfig, "directed_mode") {
		request["DirectedMode"] = d.Get("directed_mode")
	}
	if v, ok := d.GetOk("state"); ok {
		request["State"] = v
	}

	return request, nil
}

func armsNpCreateWithRequest(ctx context.Context, d *schema.ResourceData, meta interface{}, request map[string]interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateOrUpdateNotificationPolicy"
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_notification_policy", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, err := jsonpath.Get("$.NotificationPolicy.Id", response)
	if err != nil || id == nil {
		return errs.WrapErrorf(fmt.Errorf("missing NotificationPolicy.Id in create response"), errs.DefaultErrorMsg, "alicloud_arms_notification_policy", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(id))

	return ResourceNotificationPolicyReadContext(ctx, d, meta)
}

func ResourceNotificationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	return ResourceNotificationPolicyReadContext(context.Background(), d, meta)
}

func ResourceNotificationPolicyReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsNotificationPolicyContext(ctx, d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.IsProviderNotFound(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_notification_policy DescribeArmsNotificationPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	// Validate and convert every nested JSON field before any d.Set call: a
	// malformed or contract-drifted response must fail the Read instead of
	// writing a rule-less state for these ForceNew fields.
	nested := make(map[string]interface{}, 4)
	for _, key := range []string{"MatchingRules", "GroupRule", "NotifyRule", "NotifyTemplate"} {
		value := objectRaw[key]
		if text, ok := value.(string); ok && text != "" {
			if err := json.Unmarshal([]byte(text), &value); err != nil {
				return errs.WrapError(fmt.Errorf("invalid %s JSON response: %w", key, err))
			}
			objectRaw[key] = value
		}
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "MatchingRules":
			rules, err := armsNpMatchingRulesFromResponse(value)
			if err != nil {
				return errs.WrapError(err)
			}
			nested["matching_rules"] = rules
		case "GroupRule":
			rule, err := ArmsNpGroupRuleFromResponse(value)
			if err != nil {
				return errs.WrapError(err)
			}
			nested["group_rule"] = rule
		case "NotifyRule":
			rule, err := armsNpNotifyRuleFromResponse(value)
			if err != nil {
				return errs.WrapError(err)
			}
			nested["notify_rule"] = rule
		case "NotifyTemplate":
			if _, ok := value.(map[string]interface{}); !ok {
				return errs.WrapError(fmt.Errorf("unexpected %s response type %T", key, value))
			}
			nested["notify_template"] = armsNpNotifyTemplateFromResponse(value)
		}
	}

	if err := d.Set("name", objectRaw["Name"]); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("send_recover_message", objectRaw["SendRecoverMessage"]); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("escalation_policy_id", objectRaw["EscalationPolicyId"]); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("repeat", objectRaw["Repeat"]); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("repeat_interval", objectRaw["RepeatInterval"]); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("integration_id", objectRaw["IntegrationId"]); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("directed_mode", objectRaw["DirectedMode"]); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("state", objectRaw["State"]); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("notification_policy_id", d.Id()); err != nil {
		return errs.WrapError(err)
	}

	if v, ok := nested["matching_rules"]; ok {
		if err := d.Set("matching_rules", v); err != nil {
			return errs.WrapError(err)
		}
	}
	if v, ok := nested["group_rule"]; ok {
		if err := d.Set("group_rule", v); err != nil {
			return errs.WrapError(err)
		}
	}
	if v, ok := nested["notify_rule"]; ok {
		if err := d.Set("notify_rule", v); err != nil {
			return errs.WrapError(err)
		}
	}
	if v, ok := nested["notify_template"]; ok {
		if err := d.Set("notify_template", v); err != nil {
			return errs.WrapError(err)
		}
	}

	return nil
}

func ResourceNotificationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	return ResourceNotificationPolicyDeleteContext(context.Background(), d, meta)
}

func ResourceNotificationPolicyDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DeleteNotificationPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)
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
		if errs.IsTypedNotFound(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	if success, ok := response["IsSuccess"].(bool); !ok || !success {
		return errs.WrapErrorf(fmt.Errorf("delete response did not confirm IsSuccess=true"), errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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

// ArmsNpGroupRuleToJson turns a MaxItems:1 TypeList into a JSON object:
// {"groupingFields":[...],"groupWait":N,"groupInterval":N}
func ArmsNpGroupRuleToJson(raw []interface{}) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected group_rule element type")
	}
	out := map[string]interface{}{}
	if v, ok := m["grouping_fields"].([]interface{}); ok {
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
		// Omit empty values: the SDK represents unset optional strings as ""
		// and CreateOrUpdateNotificationPolicy is an upsert — sending ""
		// would clear server-side template fields the config does not manage.
		if v, ok := m[tfKey]; ok && v != "" {
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

func armsNpMatchingRulesFromResponse(raw interface{}) ([]map[string]interface{}, error) {
	parsed, ok := armsNpParseJSON(raw)
	if !ok {
		return nil, fmt.Errorf("invalid MatchingRules response value %T", raw)
	}
	list, ok := parsed.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected MatchingRules response type %T", parsed)
	}
	out := make([]map[string]interface{}, 0, len(list))
	for i, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected MatchingRules[%d] element type %T", i, item)
		}
		conds := []map[string]interface{}{}
		if condsRaw, exists := m["MatchingConditions"]; exists && condsRaw != nil {
			condsList, ok := condsRaw.([]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected MatchingRules[%d].MatchingConditions type %T", i, condsRaw)
			}
			for j, c := range condsList {
				cm, ok := c.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("unexpected MatchingRules[%d].MatchingConditions[%d] element type %T", i, j, c)
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
	return out, nil
}

func ArmsNpGroupRuleFromResponse(raw interface{}) ([]map[string]interface{}, error) {
	parsed, ok := armsNpParseJSON(raw)
	if !ok {
		return nil, fmt.Errorf("invalid GroupRule response value %T", raw)
	}
	m, ok := parsed.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected GroupRule response type %T", parsed)
	}
	out := map[string]interface{}{}
	if fieldsRaw, exists := m["GroupingFields"]; exists && fieldsRaw != nil {
		fieldsList, ok := fieldsRaw.([]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected GroupRule.GroupingFields type %T", fieldsRaw)
		}
		fields := make([]string, 0, len(fieldsList))
		for _, f := range fieldsList {
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
	return []map[string]interface{}{out}, nil
}

func armsNpNotifyRuleFromResponse(raw interface{}) ([]map[string]interface{}, error) {
	parsed, ok := armsNpParseJSON(raw)
	if !ok {
		return nil, fmt.Errorf("invalid NotifyRule response value %T", raw)
	}
	m, ok := parsed.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected NotifyRule response type %T", parsed)
	}
	out := map[string]interface{}{}
	if v, ok := m["NotifyStartTime"]; ok {
		out["notify_start_time"] = v
	}
	if v, ok := m["NotifyEndTime"]; ok {
		out["notify_end_time"] = v
	}
	if chansRaw, exists := m["NotifyChannels"]; exists && chansRaw != nil {
		chansList, ok := chansRaw.([]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected NotifyRule.NotifyChannels type %T", chansRaw)
		}
		chans := make([]string, 0, len(chansList))
		for _, c := range chansList {
			chans = append(chans, fmt.Sprint(c))
		}
		out["notify_channels"] = chans
	}
	if objsRaw, exists := m["NotifyObjects"]; exists && objsRaw != nil {
		objsList, ok := objsRaw.([]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected NotifyRule.NotifyObjects type %T", objsRaw)
		}
		objs := make([]map[string]interface{}, 0, len(objsList))
		for i, o := range objsList {
			om, ok := o.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected NotifyRule.NotifyObjects[%d] element type %T", i, o)
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
			if ocRaw, exists := om["NotifyChannels"]; exists && ocRaw != nil {
				ocList, ok := ocRaw.([]interface{})
				if !ok {
					return nil, fmt.Errorf("unexpected NotifyRule.NotifyObjects[%d].NotifyChannels type %T", i, ocRaw)
				}
				oc := make([]string, 0, len(ocList))
				for _, c := range ocList {
					oc = append(oc, fmt.Sprint(c))
				}
				obj["notify_channels"] = oc
			}
			objs = append(objs, obj)
		}
		out["notify_objects"] = objs
	}
	return []map[string]interface{}{out}, nil
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
