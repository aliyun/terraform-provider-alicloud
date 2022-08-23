package alicloud

import (
	"fmt"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCmsEventRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsEventRuleCreate,
		Read:   resourceAlicloudCmsEventRuleRead,
		Update: resourceAlicloudCmsEventRuleUpdate,
		Delete: resourceAlicloudCmsEventRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLED", "ENABLED"}, false),
			},
			"event_pattern": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product": {
							Type:     schema.TypeString,
							Required: true,
						},
						"event_type_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"level_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"sql_filter": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"silence_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCmsEventRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "PutEventRule"
	request := make(map[string]interface{})
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}

	request["RuleName"] = d.Get("rule_name")
	eventPatternMaps := make([]map[string]interface{}, 0)
	for _, eventPattern := range d.Get("event_pattern").(*schema.Set).List() {
		eventPatternMap := make(map[string]interface{})
		eventPatternArg := eventPattern.(map[string]interface{})

		if v, ok := eventPatternArg["product"].(string); ok {
			eventPatternMap["Product"] = v
		}
		if v, ok := eventPatternArg["event_type_list"]; ok {
			eventPatternMap["EventTypeList"] = v
		}
		if v, ok := eventPatternArg["level_list"]; ok {
			eventPatternMap["LevelList"] = v
		}
		if v, ok := eventPatternArg["name_list"]; ok {
			eventPatternMap["NameList"] = v
		}
		if v, ok := eventPatternArg["sql_filter"].(string); ok && v != "" {
			eventPatternMap["SQLFilter"] = v
		}

		eventPatternMaps = append(eventPatternMaps, eventPatternMap)
	}
	request["EventPattern"] = eventPatternMaps

	if v, ok := d.GetOk("group_id"); ok {
		request["GroupId"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["State"] = v
	}
	if v, ok := d.GetOkExists("silence_time"); ok {
		request["SilenceTime"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_event_rule", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	d.SetId(fmt.Sprint(request["RuleName"]))

	return resourceAlicloudCmsEventRuleRead(d, meta)
}

func resourceAlicloudCmsEventRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsEventRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("rule_name", object["Name"])
	d.Set("group_id", object["GroupId"])
	d.Set("description", object["Description"])
	d.Set("status", object["State"])

	if v, ok := object["EventPattern"]; ok {
		eventPatternMaps := make([]map[string]interface{}, 0)
		eventPattern := v.(map[string]interface{})
		eventPatternMap := make(map[string]interface{})
		eventPatternMap["product"] = eventPattern["Product"]
		eventPatternMap["event_type_list"] = eventPattern["EventTypeList"].(map[string]interface{})["EventTypeList"]
		eventPatternMap["level_list"] = eventPattern["LevelList"].(map[string]interface{})["LevelList"]
		eventPatternMap["name_list"] = eventPattern["NameList"].(map[string]interface{})["NameList"]
		eventPatternMap["sql_filter"] = eventPattern["SQLFilter"]
		eventPatternMaps = append(eventPatternMaps, eventPatternMap)
		d.Set("event_pattern", eventPatternMaps)
	}
	silenceTime, _ := strconv.Atoi(fmt.Sprint(object["SilenceTime"]))
	d.Set("silence_time", silenceTime)

	return nil
}

func resourceAlicloudCmsEventRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"RuleName": d.Id(),
	}

	eventPatternMaps := make([]map[string]interface{}, 0)
	for _, eventPattern := range d.Get("event_pattern").(*schema.Set).List() {
		eventPatternMap := make(map[string]interface{})
		eventPatternArg := eventPattern.(map[string]interface{})

		if v, ok := eventPatternArg["product"].(string); ok {
			eventPatternMap["Product"] = v
		}
		if v, ok := eventPatternArg["event_type_list"]; ok {
			eventPatternMap["EventTypeList"] = v
		}
		if v, ok := eventPatternArg["level_list"]; ok {
			eventPatternMap["LevelList"] = v
		}
		if v, ok := eventPatternArg["name_list"]; ok {
			eventPatternMap["NameList"] = v
		}
		if v, ok := eventPatternArg["sql_filter"].(string); ok && v != "" {
			eventPatternMap["SQLFilter"] = v
		}

		eventPatternMaps = append(eventPatternMaps, eventPatternMap)
	}
	request["EventPattern"] = eventPatternMaps

	if !d.IsNewResource() && d.HasChange("group_id") {
		update = true
	}
	if v, ok := d.GetOk("group_id"); ok {
		request["GroupId"] = v
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if !d.IsNewResource() && d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok {
		request["State"] = v
	}

	if !d.IsNewResource() && d.HasChange("event_pattern") {
		update = true
		eventPatternMaps := make([]map[string]interface{}, 0)
		for _, eventPattern := range d.Get("event_pattern").(*schema.Set).List() {
			eventPatternMap := make(map[string]interface{}, 0)
			eventPatternArg := eventPattern.(map[string]interface{})

			if v, ok := eventPatternArg["product"].(string); ok {
				eventPatternMap["Product"] = v
			}
			if v, ok := eventPatternArg["event_type_list"]; ok {
				eventPatternMap["EventTypeList"] = v
			}
			if v, ok := eventPatternArg["level_list"]; ok {
				eventPatternMap["LevelList"] = v
			}
			if v, ok := eventPatternArg["name_list"]; ok {
				eventPatternMap["NameList"] = v
			}
			if v, ok := eventPatternArg["sql_filter"].(string); ok && v != "" {
				eventPatternMap["SQLFilter"] = v
			}

			eventPatternMaps = append(eventPatternMaps, eventPatternMap)
		}
		request["EventPattern"] = eventPatternMaps

	}

	if !d.IsNewResource() && d.HasChange("silence_time") {
		update = true
	}
	if v, ok := d.GetOkExists("silence_time"); ok {
		request["SilenceTime"] = v
	}

	if update {
		action := "PutEventRule"
		conn, err := client.NewCmsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	return resourceAlicloudCmsEventRuleRead(d, meta)
}

func resourceAlicloudCmsEventRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEventRules"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RuleNames": []string{d.Id()},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
