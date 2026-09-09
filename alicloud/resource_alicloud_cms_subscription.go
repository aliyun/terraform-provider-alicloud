package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCmsSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsSubscriptionCreate,
		Read:   resourceAlicloudCmsSubscriptionRead,
		Update: resourceAlicloudCmsSubscriptionUpdate,
		Delete: resourceAlicloudCmsSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"subscription_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_strategy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"op": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"IN", "EQ"}, false),
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
						"restore_action_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"template_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"response_plan_id": {
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
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCmsSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/subscriptions"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["subscriptionName"] = d.Get("subscription_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("notify_strategy_id"); ok {
		request["notifyStrategyId"] = v
	}
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("filter_setting"); ok {
		if fs := buildCmsSubscriptionFilterSettingInput(v.([]interface{})); fs != nil {
			request["filterSetting"] = fs
		}
	}
	if v, ok := d.GetOk("pushing_setting"); ok {
		if ps := buildCmsSubscriptionPushingSettingInput(v.([]interface{})); ps != nil {
			request["pushingSetting"] = ps
		}
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

	id, _ := jsonpath.Get("$.data", response)
	d.SetId(fmt.Sprint(id))

	if !d.Get("enable").(bool) {
		if err := cmsSubscriptionToggleEnable(d, meta, false); err != nil {
			return err
		}
	}

	return resourceAlicloudCmsSubscriptionRead(d, meta)
}

func resourceAlicloudCmsSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/subscriptions/%s", d.Id())
	query := make(map[string]*string)
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if NotFoundError(err) || IsExpectedErrors(err, []string{"NotFound", "ResourceNotFound"}) {
				log.Printf("[DEBUG] Resource alicloud_cms_subscription GetSubscription Failed!!! %s", err)
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	objectRaw, err := jsonpath.Get("$.subscription", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "$.subscription", response)
	}
	object := make(map[string]interface{})
	if objectRaw != nil {
		object, _ = objectRaw.(map[string]interface{})
	}

	d.Set("subscription_id", object["subscriptionId"])
	d.Set("subscription_name", object["subscriptionName"])
	d.Set("description", object["description"])
	d.Set("notify_strategy_id", object["notifyStrategyId"])
	d.Set("workspace", object["workspace"])
	d.Set("user_id", object["userId"])
	d.Set("create_time", object["createTime"])
	d.Set("update_time", object["updateTime"])
	d.Set("enable", object["enable"])
	d.Set("region_id", object["regionId"])
	d.Set("subscription_type", object["subscriptionType"])
	if err := setCmsSubscriptionFilterSettingOutput(d, object["filterSetting"]); err != nil {
		return err
	}
	if err := setCmsSubscriptionPushingSettingOutput(d, object["pushingSetting"]); err != nil {
		return err
	}

	return nil
}

func resourceAlicloudCmsSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	if d.HasChange("enable") {
		if err := cmsSubscriptionToggleEnable(d, meta, d.Get("enable").(bool)); err != nil {
			return err
		}
		d.SetPartial("enable")
	}

	if d.HasChange("subscription_name") || d.HasChange("description") || d.HasChange("notify_strategy_id") ||
		d.HasChange("workspace") || d.HasChange("filter_setting") || d.HasChange("pushing_setting") {
		action := fmt.Sprintf("/subscriptions/%s", d.Id())
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		body := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})

		request["subscriptionName"] = d.Get("subscription_name")
		if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
			request["description"] = v
		}
		if v, ok := d.GetOk("notify_strategy_id"); ok || d.HasChange("notify_strategy_id") {
			request["notifyStrategyId"] = v
		}
		if v, ok := d.GetOk("workspace"); ok || d.HasChange("workspace") {
			query["workspace"] = StringPointer(v.(string))
		}
		if v, ok := d.GetOk("filter_setting"); ok || d.HasChange("filter_setting") {
			if fs := buildCmsSubscriptionFilterSettingInput(v.([]interface{})); fs != nil {
				request["filterSetting"] = fs
			}
		}
		if v, ok := d.GetOk("pushing_setting"); ok || d.HasChange("pushing_setting") {
			if ps := buildCmsSubscriptionPushingSettingInput(v.([]interface{})); ps != nil {
				request["pushingSetting"] = ps
			}
		}
		body = request
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
		d.SetPartial("subscription_name")
		d.SetPartial("description")
		d.SetPartial("notify_strategy_id")
		d.SetPartial("workspace")
		d.SetPartial("filter_setting")
		d.SetPartial("pushing_setting")
	}

	d.Partial(false)
	return resourceAlicloudCmsSubscriptionRead(d, meta)
}

func resourceAlicloudCmsSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/subscriptions/%s", d.Id())
	query := make(map[string]*string)
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}
	var response map[string]interface{}
	var err error
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
	addDebug(action, response, query)
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"NotFound", "ResourceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func cmsSubscriptionToggleEnable(d *schema.ResourceData, meta interface{}, enable bool) error {
	client := meta.(*connectivity.AliyunClient)

	suffix := "enable"
	if !enable {
		suffix = "disable"
	}
	action := fmt.Sprintf("/subscriptions/%s/%s", d.Id(), suffix)
	query := make(map[string]*string)
	if v, ok := d.GetOk("workspace"); ok {
		query["workspace"] = StringPointer(v.(string))
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func buildCmsSubscriptionFilterSettingInput(list []interface{}) map[string]interface{} {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	filterSetting := make(map[string]interface{})
	if conditions, ok := m["conditions"].([]interface{}); ok && len(conditions) > 0 {
		conditionList := make([]map[string]interface{}, 0)
		for _, c := range conditions {
			if cm, ok := c.(map[string]interface{}); ok {
				conditionList = append(conditionList, map[string]interface{}{
					"field": cm["field"],
					"value": cm["value"],
					"op":    cm["op"],
				})
			}
		}
		filterSetting["conditions"] = conditionList
	}
	if v, ok := m["expression"].(string); ok && v != "" {
		filterSetting["expression"] = v
	}
	if v, ok := m["relation"].(string); ok && v != "" {
		filterSetting["relation"] = v
	}
	return filterSetting
}

func buildCmsSubscriptionPushingSettingInput(list []interface{}) map[string]interface{} {
	if len(list) == 0 || list[0] == nil {
		return nil
	}
	m, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	pushingSetting := make(map[string]interface{})
	if alertActionIds, ok := m["alert_action_ids"].([]interface{}); ok && len(alertActionIds) > 0 {
		ids := make([]string, 0)
		for _, v := range alertActionIds {
			ids = append(ids, fmt.Sprint(v))
		}
		pushingSetting["alertActionIds"] = ids
	}
	if restoreActionIds, ok := m["restore_action_ids"].([]interface{}); ok && len(restoreActionIds) > 0 {
		ids := make([]string, 0)
		for _, v := range restoreActionIds {
			ids = append(ids, fmt.Sprint(v))
		}
		pushingSetting["restoreActionIds"] = ids
	}
	if v, ok := m["template_uuid"].(string); ok && v != "" {
		pushingSetting["templateUuid"] = v
	}
	if v, ok := m["response_plan_id"].(string); ok && v != "" {
		pushingSetting["responsePlanId"] = v
	}
	return pushingSetting
}

func setCmsSubscriptionFilterSettingOutput(d *schema.ResourceData, raw interface{}) error {
	filterResult := make([]map[string]interface{}, 0)
	if raw != nil {
		if filterMap, ok := raw.(map[string]interface{}); ok {
			fs := map[string]interface{}{}
			if conditions, ok := filterMap["conditions"].([]interface{}); ok {
				condList := make([]map[string]interface{}, 0)
				for _, c := range conditions {
					if cm, ok := c.(map[string]interface{}); ok {
						condList = append(condList, map[string]interface{}{
							"field": cm["field"],
							"value": cm["value"],
							"op":    cm["op"],
						})
					}
				}
				fs["conditions"] = condList
			}
			fs["expression"] = filterMap["expression"]
			fs["relation"] = filterMap["relation"]
			filterResult = append(filterResult, fs)
		}
	}
	return d.Set("filter_setting", filterResult)
}

func setCmsSubscriptionPushingSettingOutput(d *schema.ResourceData, raw interface{}) error {
	pushingResult := make([]map[string]interface{}, 0)
	if raw != nil {
		if pushingMap, ok := raw.(map[string]interface{}); ok {
			ps := map[string]interface{}{}
			if alertActionIds, ok := pushingMap["alertActionIds"].([]interface{}); ok {
				ids := make([]string, 0)
				for _, v := range alertActionIds {
					ids = append(ids, fmt.Sprint(v))
				}
				ps["alert_action_ids"] = ids
			}
			if restoreActionIds, ok := pushingMap["restoreActionIds"].([]interface{}); ok {
				ids := make([]string, 0)
				for _, v := range restoreActionIds {
					ids = append(ids, fmt.Sprint(v))
				}
				ps["restore_action_ids"] = ids
			}
			ps["template_uuid"] = pushingMap["templateUuid"]
			ps["response_plan_id"] = pushingMap["responsePlanId"]
			pushingResult = append(pushingResult, ps)
		}
	}
	return d.Set("pushing_setting", pushingResult)
}
