// Package alicloud. This file is hand-written from the Cms 2024-03-30 CloudSpec Mission definition.
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

func resourceAlicloudCmsMission() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsMissionCreate,
		Read:   resourceAlicloudCmsMissionRead,
		Update: resourceAlicloudCmsMissionUpdate,
		Delete: resourceAlicloudCmsMissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"digital_employee_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"variables": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"notification": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dingtalk": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"feishu": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"slack": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"wechat": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"call": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"sms": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"email": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"webhook": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"notification_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"blueprint": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: cmsMissionBlueprintCronSchema(),
							},
						},
						"calendar": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: cmsMissionBlueprintCalendarSchema(),
							},
						},
						"event": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: cmsMissionBlueprintEventSchema(),
							},
						},
					},
				},
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credits": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func cmsMissionBlueprintCronSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"prompt": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"cron_expression": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cron expression in 5-field Unix format (min hour day-of-month month day-of-week), e.g. `0 0 * * *`. Six-field (with seconds) and Quartz `?` syntax are not supported by the CMS API.",
		},
		"variables": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"time_zone": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "+0800",
		},
		"delay": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"run_immediately": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"priority": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"concurrency_policy": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "skip",
		},
		"timeout_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func cmsMissionBlueprintCalendarSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"prompt": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"rrule": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"variables": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"time_zone": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "+0800",
		},
		"delay": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"run_immediately": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"priority": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"concurrency_policy": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "skip",
		},
		"timeout_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func cmsMissionBlueprintEventSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"prompt": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"workspace": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"max_concurrency": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"debounce_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"variables": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"priority": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"concurrency_policy": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "skip",
		},
		"timeout_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

// DescribeCmsMission <<< Encapsulated get interface for Cms Mission.

func (s *CmsServiceV2) DescribeCmsMission(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/mission/%s", id)

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
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
		if IsExpectedErrors(err, []string{"MissionNotFound", "ResourceNotFound"}) || NotFoundError(err) {
			return object, WrapErrorf(NotFoundErr("Mission", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// DescribeCmsMission >>> Encapsulated.

func (s *CmsServiceV2) CmsMissionStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.CmsMissionStateRefreshFuncWithApi(id, field, failStates, s.DescribeCmsMission)
}

func (s *CmsServiceV2) CmsMissionStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, _ := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func resourceAlicloudCmsMissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/mission"
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error

	body := buildCmsMissionRequestBody(d, true)

	wait := incrementalWait(3*time.Second, 3*time.Second)
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
	addDebug(action, response, body)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_mission", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["name"]))

	return resourceAlicloudCmsMissionRead(d, meta)
}

func resourceAlicloudCmsMissionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsMission(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_mission DescribeCmsMission Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", objectRaw["name"])
	d.Set("display_name", objectRaw["displayName"])
	d.Set("description", objectRaw["description"])
	d.Set("digital_employee_name", objectRaw["digitalEmployeeName"])
	d.Set("enabled", objectRaw["enabled"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("update_time", objectRaw["updateTime"])
	if v, ok := objectRaw["variables"]; ok && v != nil {
		d.Set("variables", flattenCmsMissionMap(v))
	}
	if v, ok := objectRaw["notificationPolicy"]; ok && v != nil {
		d.Set("notification_policy", flattenCmsMissionNotificationPolicy(v))
	}
	if v, ok := objectRaw["notification"]; ok && v != nil {
		d.Set("notification", flattenCmsMissionNotification(v))
	}
	if v, ok := objectRaw["configuration"]; ok && v != nil {
		d.Set("configuration", flattenCmsMissionConfiguration(v))
	}
	if v, ok := objectRaw["blueprint"]; ok && v != nil {
		d.Set("blueprint", flattenCmsMissionBlueprint(v))
	}
	return nil
}

func resourceAlicloudCmsMissionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/mission/%s", d.Id())
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error

	body := buildCmsMissionRequestBody(d, false)

	wait := incrementalWait(3*time.Second, 3*time.Second)
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
	addDebug(action, response, body)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudCmsMissionRead(d, meta)
}

func resourceAlicloudCmsMissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/mission/%s", d.Id())
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error

	wait := incrementalWait(3*time.Second, 3*time.Second)
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
	addDebug(action, response, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"MissionNotFound", "ResourceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

// buildCmsMissionRequestBody builds the JSON request body for Create (includeName=true)
// and Update (includeName=false, name is in the path).
func buildCmsMissionRequestBody(d *schema.ResourceData, includeName bool) map[string]interface{} {
	request := make(map[string]interface{})
	if includeName {
		request["name"] = d.Get("name")
	}
	if v, ok := d.GetOk("display_name"); ok {
		request["displayName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["digitalEmployeeName"] = d.Get("digital_employee_name")
	if v, ok := d.GetOk("variables"); ok && len(v.(map[string]interface{})) > 0 {
		request["variables"] = expandCmsMissionMap(v)
	}
	request["enabled"] = d.Get("enabled")
	if v, ok := d.GetOk("notification"); ok && len(v.([]interface{})) > 0 {
		request["notification"] = expandCmsMissionNotification(v.([]interface{}))
	}
	if v, ok := d.GetOk("notification_policy"); ok && len(v.([]interface{})) > 0 {
		request["notificationPolicy"] = expandCmsMissionNotificationPolicy(v.([]interface{}))
	}
	if v, ok := d.GetOk("configuration"); ok && len(v.([]interface{})) > 0 {
		request["configuration"] = expandCmsMissionConfiguration(v.([]interface{}))
	}
	if v, ok := d.GetOk("blueprint"); ok && len(v.([]interface{})) > 0 {
		request["blueprint"] = expandCmsMissionBlueprint(v.([]interface{}))
	}
	return request
}

func expandCmsMissionMap(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, val := range v.(map[string]interface{}) {
		result[key] = fmt.Sprint(val)
	}
	return result
}

func flattenCmsMissionMap(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if m, ok := v.(map[string]interface{}); ok {
		for key, val := range m {
			result[key] = fmt.Sprint(val)
		}
	}
	return result
}

func expandCmsMissionNotification(list []interface{}) map[string]interface{} {
	if len(list) == 0 || list[0] == nil {
		return map[string]interface{}{}
	}
	m := list[0].(map[string]interface{})
	result := make(map[string]interface{})
	channels := []string{"dingtalk", "feishu", "slack", "wechat", "call", "sms", "email", "webhook"}
	for _, ch := range channels {
		if v, ok := m[ch]; ok && v != nil {
			result[ch] = expandStringList(v.([]interface{}))
		}
	}
	return result
}

func flattenCmsMissionNotification(v interface{}) []map[string]interface{} {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]interface{})
	channels := []string{"dingtalk", "feishu", "slack", "wechat", "call", "sms", "email", "webhook"}
	for _, ch := range channels {
		if raw, ok := m[ch]; ok && raw != nil {
			if list, ok := raw.([]interface{}); ok {
				result[ch] = list
			}
		}
	}
	return []map[string]interface{}{result}
}

func expandCmsMissionNotificationPolicy(list []interface{}) map[string]interface{} {
	if len(list) == 0 || list[0] == nil {
		return map[string]interface{}{}
	}
	m := list[0].(map[string]interface{})
	result := make(map[string]interface{})
	if v, ok := m["region"]; ok {
		result["region"] = v
	}
	if v, ok := m["workspace"]; ok {
		result["workspace"] = v
	}
	return result
}

func flattenCmsMissionNotificationPolicy(v interface{}) []map[string]interface{} {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]interface{})
	if v, ok := m["region"]; ok {
		result["region"] = v
	}
	if v, ok := m["workspace"]; ok {
		result["workspace"] = v
	}
	return []map[string]interface{}{result}
}

func expandCmsMissionConfiguration(list []interface{}) map[string]interface{} {
	if len(list) == 0 || list[0] == nil {
		return map[string]interface{}{}
	}
	m := list[0].(map[string]interface{})
	result := make(map[string]interface{})
	if v, ok := m["credits"]; ok {
		result["credits"] = v
	}
	return result
}

func flattenCmsMissionConfiguration(v interface{}) []map[string]interface{} {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]interface{})
	if v, ok := m["credits"]; ok {
		result["credits"] = v
	}
	return []map[string]interface{}{result}
}

func expandCmsMissionBlueprint(list []interface{}) map[string]interface{} {
	if len(list) == 0 || list[0] == nil {
		return map[string]interface{}{}
	}
	m := list[0].(map[string]interface{})
	result := make(map[string]interface{})
	if v, ok := m["cron"]; ok && v != nil {
		result["cron"] = expandCmsMissionBlueprintCronList(v.([]interface{}))
	}
	if v, ok := m["calendar"]; ok && v != nil {
		result["calendar"] = expandCmsMissionBlueprintCalendarList(v.([]interface{}))
	}
	if v, ok := m["event"]; ok && v != nil {
		result["event"] = expandCmsMissionBlueprintEventList(v.([]interface{}))
	}
	return result
}

func flattenCmsMissionBlueprint(v interface{}) []map[string]interface{} {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]interface{})
	if v, ok := m["cron"]; ok && v != nil {
		result["cron"] = flattenCmsMissionBlueprintCronList(v)
	}
	if v, ok := m["calendar"]; ok && v != nil {
		result["calendar"] = flattenCmsMissionBlueprintCalendarList(v)
	}
	if v, ok := m["event"]; ok && v != nil {
		result["event"] = flattenCmsMissionBlueprintEventList(v)
	}
	return []map[string]interface{}{result}
}

func expandCmsMissionBlueprintCronList(list []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		entry := make(map[string]interface{})
		if v, ok := m["id"]; ok {
			entry["id"] = v
		}
		if v, ok := m["display_name"]; ok {
			entry["displayName"] = v
		}
		if v, ok := m["description"]; ok {
			entry["description"] = v
		}
		if v, ok := m["prompt"]; ok {
			entry["prompt"] = v
		}
		if v, ok := m["cron_expression"]; ok {
			entry["cronExpression"] = v
		}
		if v, ok := m["variables"]; ok && len(v.(map[string]interface{})) > 0 {
			entry["variables"] = expandCmsMissionMap(v)
		}
		if v, ok := m["time_zone"]; ok {
			entry["timeZone"] = v
		}
		if v, ok := m["delay"]; ok {
			entry["delay"] = v
		}
		if v, ok := m["run_immediately"]; ok {
			entry["runImmediately"] = v
		}
		if v, ok := m["priority"]; ok {
			entry["priority"] = v
		}
		if v, ok := m["concurrency_policy"]; ok {
			entry["concurrencyPolicy"] = v
		}
		if v, ok := m["timeout_seconds"]; ok {
			entry["timeoutSeconds"] = v
		}
		result = append(result, entry)
	}
	return result
}

func flattenCmsMissionBlueprintCronList(v interface{}) []map[string]interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		entry := make(map[string]interface{})
		entry["id"] = m["id"]
		entry["display_name"] = m["displayName"]
		entry["description"] = m["description"]
		entry["prompt"] = m["prompt"]
		entry["cron_expression"] = m["cronExpression"]
		if vv, ok := m["variables"]; ok && vv != nil {
			entry["variables"] = flattenCmsMissionMap(vv)
		}
		entry["time_zone"] = m["timeZone"]
		entry["delay"] = m["delay"]
		entry["run_immediately"] = m["runImmediately"]
		entry["priority"] = m["priority"]
		entry["concurrency_policy"] = m["concurrencyPolicy"]
		entry["timeout_seconds"] = m["timeoutSeconds"]
		entry["enabled"] = m["enabled"]
		result = append(result, entry)
	}
	return result
}

func expandCmsMissionBlueprintCalendarList(list []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		entry := make(map[string]interface{})
		if v, ok := m["id"]; ok {
			entry["id"] = v
		}
		if v, ok := m["display_name"]; ok {
			entry["displayName"] = v
		}
		if v, ok := m["description"]; ok {
			entry["description"] = v
		}
		if v, ok := m["prompt"]; ok {
			entry["prompt"] = v
		}
		if v, ok := m["rrule"]; ok {
			entry["rrule"] = v
		}
		if v, ok := m["variables"]; ok && len(v.(map[string]interface{})) > 0 {
			entry["variables"] = expandCmsMissionMap(v)
		}
		if v, ok := m["time_zone"]; ok {
			entry["timeZone"] = v
		}
		if v, ok := m["delay"]; ok {
			entry["delay"] = v
		}
		if v, ok := m["run_immediately"]; ok {
			entry["runImmediately"] = v
		}
		if v, ok := m["priority"]; ok {
			entry["priority"] = v
		}
		if v, ok := m["concurrency_policy"]; ok {
			entry["concurrencyPolicy"] = v
		}
		if v, ok := m["timeout_seconds"]; ok {
			entry["timeoutSeconds"] = v
		}
		result = append(result, entry)
	}
	return result
}

func flattenCmsMissionBlueprintCalendarList(v interface{}) []map[string]interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		entry := make(map[string]interface{})
		entry["id"] = m["id"]
		entry["display_name"] = m["displayName"]
		entry["description"] = m["description"]
		entry["prompt"] = m["prompt"]
		entry["rrule"] = m["rrule"]
		if vv, ok := m["variables"]; ok && vv != nil {
			entry["variables"] = flattenCmsMissionMap(vv)
		}
		entry["time_zone"] = m["timeZone"]
		entry["delay"] = m["delay"]
		entry["run_immediately"] = m["runImmediately"]
		entry["priority"] = m["priority"]
		entry["concurrency_policy"] = m["concurrencyPolicy"]
		entry["timeout_seconds"] = m["timeoutSeconds"]
		entry["enabled"] = m["enabled"]
		result = append(result, entry)
	}
	return result
}

func expandCmsMissionBlueprintEventList(list []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		entry := make(map[string]interface{})
		if v, ok := m["id"]; ok {
			entry["id"] = v
		}
		if v, ok := m["display_name"]; ok {
			entry["displayName"] = v
		}
		if v, ok := m["description"]; ok {
			entry["description"] = v
		}
		if v, ok := m["prompt"]; ok {
			entry["prompt"] = v
		}
		if v, ok := m["workspace"]; ok {
			entry["workspace"] = v
		}
		if v, ok := m["max_concurrency"]; ok {
			entry["maxConcurrency"] = v
		}
		if v, ok := m["debounce_seconds"]; ok {
			entry["debounceSeconds"] = v
		}
		if v, ok := m["variables"]; ok && len(v.(map[string]interface{})) > 0 {
			entry["variables"] = expandCmsMissionMap(v)
		}
		if v, ok := m["priority"]; ok {
			entry["priority"] = v
		}
		if v, ok := m["concurrency_policy"]; ok {
			entry["concurrencyPolicy"] = v
		}
		if v, ok := m["timeout_seconds"]; ok {
			entry["timeoutSeconds"] = v
		}
		result = append(result, entry)
	}
	return result
}

func flattenCmsMissionBlueprintEventList(v interface{}) []map[string]interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		entry := make(map[string]interface{})
		entry["id"] = m["id"]
		entry["display_name"] = m["displayName"]
		entry["description"] = m["description"]
		entry["prompt"] = m["prompt"]
		entry["workspace"] = m["workspace"]
		entry["max_concurrency"] = m["maxConcurrency"]
		entry["debounce_seconds"] = m["debounceSeconds"]
		if vv, ok := m["variables"]; ok && vv != nil {
			entry["variables"] = flattenCmsMissionMap(vv)
		}
		entry["priority"] = m["priority"]
		entry["concurrency_policy"] = m["concurrencyPolicy"]
		entry["timeout_seconds"] = m["timeoutSeconds"]
		entry["enabled"] = m["enabled"]
		result = append(result, entry)
	}
	return result
}

// jsonpathGet wraps jsonpath.Get to avoid importing the package in every hand-written file.
func jsonpathGet(path string, object map[string]interface{}) (interface{}, error) {
	return jsonpath.Get(path, object)
}
