package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEventBridgeRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEventBridgeRuleCreate,
		Read:   resourceAlicloudEventBridgeRuleRead,
		Update: resourceAlicloudEventBridgeRuleUpdate,
		Delete: resourceAlicloudEventBridgeRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_bus_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"filter_pattern": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLE", "ENABLE"}, false),
			},
			"targets": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"param_list": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"form": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"ORIGINAL", "TEMPLATE", "JSONPATH", "CONSTANT"}, false),
									},
									"resource_key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"template": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"target_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"acs.fc.function", "acs.mns.topic", "acs.mns.queue", "http", "acs.sms", "acs.mail", "acs.dingtalk", "https", "acs.eventbridge", "acs.rabbitmq", "acs.rocketmq"}, false),
						},
						"push_retry_strategy": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"BACKOFF_RETRY", "EXPONENTIAL_DECAY_RETRY"}, false),
						},
						"dead_letter_queue": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"arn": {
										Type:     schema.TypeString,
										Optional: true,
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

func resourceAlicloudEventBridgeRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateRule"
	request := make(map[string]interface{})
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	request["EventBusName"] = d.Get("event_bus_name")
	if v, ok := d.GetOk("filter_pattern"); ok {
		request["FilterPattern"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["RuleName"] = d.Get("rule_name")
	if v, ok := d.GetOk("targets"); ok {
		targetsMaps := make([]map[string]interface{}, 0)
		for _, targets := range v.([]interface{}) {
			targetsArg := targets.(map[string]interface{})
			targetsMap := map[string]interface{}{}
			targetsMap["endpoint"] = targetsArg["endpoint"]
			targetsMap["id"] = targetsArg["target_id"]
			targetsMap["type"] = targetsArg["type"]
			if pushRetryStrategy, ok := targetsMap["pushRetryStrategy"]; ok && pushRetryStrategy != "" {
				targetsMap["pushRetryStrategy"] = targetsArg["push_retry_strategy"]
			}
			if v, ok := targetsArg["dead_letter_queue"]; ok {
				deadLetterQueueMap := map[string]interface{}{}
				for _, deadLetterQueue := range v.(*schema.Set).List() {
					deadLetterQueueArg := deadLetterQueue.(map[string]interface{})
					deadLetterQueueMap["arn"] = deadLetterQueueArg["arn"]
				}
				targetsMap["deadLetterQueue"] = deadLetterQueueMap
			}
			paramListMaps := make([]map[string]interface{}, 0)
			for _, paramList := range targetsArg["param_list"].(*schema.Set).List() {
				paramListArg := paramList.(map[string]interface{})
				paramListMap := map[string]interface{}{}
				paramListMap["form"] = paramListArg["form"]
				paramListMap["resourceKey"] = paramListArg["resource_key"]
				if paramListMap["form"] == "TEMPLATE" {
					paramListMap["template"] = paramListArg["template"]
					paramListMap["value"] = paramListArg["value"]
				}
				if paramListMap["form"] == "JSONPATH" || paramListMap["form"] == "CONSTANT" {
					paramListMap["value"] = paramListArg["value"]
				}
				paramListMaps = append(paramListMaps, paramListMap)
			}
			targetsMap["paramList"] = paramListMaps
			targetsMaps = append(targetsMaps, targetsMap)
		}
		if v, err := convertArrayObjectToJsonString(targetsMaps); err == nil {
			request["Targets"] = v
		} else {
			return WrapError(err)
		}
	}
	request["ClientToken"] = buildClientToken("CreateRule")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_event_bridge_rule", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("CreateRule failed, response: %v", response))
	}

	d.SetId(fmt.Sprint(request["EventBusName"], ":", request["RuleName"]))

	return resourceAlicloudEventBridgeRuleUpdate(d, meta)
}

func resourceAlicloudEventBridgeRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventbridgeService := EventbridgeService{client}
	object, err := eventbridgeService.DescribeEventBridgeRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_rule eventbridgeService.DescribeEventBridgeRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("event_bus_name", parts[0])
	d.Set("rule_name", parts[1])
	d.Set("description", object["Description"])
	d.Set("filter_pattern", object["FilterPattern"])
	d.Set("status", object["Status"])

	if targetsList, ok := object["Targets"]; ok && targetsList != nil {
		targetsMaps := make([]map[string]interface{}, 0)
		for _, targetsListItem := range targetsList.([]interface{}) {
			if targetsListItemMap, ok := targetsListItem.(map[string]interface{}); ok {
				targetsListMap := make(map[string]interface{})
				targetsListMap["endpoint"] = targetsListItemMap["Endpoint"]
				targetsListMap["target_id"] = targetsListItemMap["Id"]
				targetsListMap["type"] = targetsListItemMap["Type"]
				targetsListMap["push_retry_strategy"] = targetsListItemMap["PushRetryStrategy"]

				if deadLetterQueue, ok := targetsListItemMap["DeadLetterQueue"]; ok && deadLetterQueue != nil {
					deadLetterQueueMaps := make([]map[string]interface{}, 0)
					deadLetterQueueArg := deadLetterQueue.(map[string]interface{})
					deadLetterQueueMap := map[string]interface{}{}
					if deadLetterQueueArgArn, ok := deadLetterQueueArg["Arn"]; ok && deadLetterQueueArgArn != nil {
						deadLetterQueueMap["arn"] = deadLetterQueueArgArn
						deadLetterQueueMaps = append(deadLetterQueueMaps, deadLetterQueueMap)
						targetsListMap["dead_letter_queue"] = deadLetterQueueMaps
					}
				}

				if paramListMap, ok := targetsListItemMap["ParamList"]; ok && paramListMap != nil {
					paramListMaps := make([]map[string]interface{}, 0)
					for _, paramListMapItem := range paramListMap.([]interface{}) {
						paramListMap := make(map[string]interface{})
						if paramListMapItemMap, ok := paramListMapItem.(map[string]interface{}); ok {
							// There is an api bug that the event bridge service will return a default param which ResourceKey is IsBase64Encode and the value is false
							if fmt.Sprint(paramListMapItemMap["ResourceKey"]) == "IsBase64Encode" && fmt.Sprint(paramListMapItemMap["Value"]) == "false" {
								continue
							}
							paramListMap["form"] = paramListMapItemMap["Form"]
							paramListMap["resource_key"] = paramListMapItemMap["ResourceKey"]
							paramListMap["template"] = paramListMapItemMap["Template"]
							paramListMap["value"] = paramListMapItemMap["Value"]
							paramListMaps = append(paramListMaps, paramListMap)
						}
					}
					targetsListMap["param_list"] = paramListMaps
				}

				targetsMaps = append(targetsMaps, targetsListMap)
			}
		}

		err = d.Set("targets", targetsMaps)
		if err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func resourceAlicloudEventBridgeRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventbridgeService := EventbridgeService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"EventBusName": parts[0],
		"RuleName":     parts[1],
	}

	if !d.IsNewResource() && d.HasChange("targets") {
		update = true
		targetsMaps := make([]map[string]interface{}, 0)
		for _, targets := range d.Get("targets").([]interface{}) {
			targetsArg := targets.(map[string]interface{})
			targetsMap := map[string]interface{}{}
			targetsMap["endpoint"] = targetsArg["endpoint"]
			targetsMap["id"] = targetsArg["target_id"]
			targetsMap["type"] = targetsArg["type"]
			targetsMap["pushRetryStrategy"] = targetsArg["push_retry_strategy"]
			if v, ok := targetsArg["dead_letter_queue"]; ok {
				deadLetterQueueMap := map[string]interface{}{}
				for _, deadLetterQueue := range v.(*schema.Set).List() {
					deadLetterQueueArg := deadLetterQueue.(map[string]interface{})
					deadLetterQueueMap["arn"] = deadLetterQueueArg["arn"]
				}
				targetsMap["deadLetterQueue"] = deadLetterQueueMap
			}
			targetsMaps = append(targetsMaps, targetsMap)
			paramListMaps := make([]map[string]interface{}, 0)
			for _, paramList := range targetsArg["param_list"].(*schema.Set).List() {
				paramListArg := paramList.(map[string]interface{})
				paramListMap := map[string]interface{}{}
				paramListMap["form"] = paramListArg["form"]
				paramListMap["resourceKey"] = paramListArg["resource_key"]
				if paramListMap["form"] == "TEMPLATE" {
					paramListMap["template"] = paramListArg["template"]
					paramListMap["value"] = paramListArg["value"]
				}
				if paramListMap["form"] == "JSONPATH" || paramListMap["form"] == "CONSTANT" {
					paramListMap["value"] = paramListArg["value"]
				}
				paramListMaps = append(paramListMaps, paramListMap)
			}
			targetsMap["paramList"] = paramListMaps
		}
		if v, err := convertArrayObjectToJsonString(targetsMaps); err == nil {
			request["Targets"] = v
		} else {
			return WrapError(err)
		}
	}
	if update {
		action := "UpdateTargets"
		conn, err := client.NewEventbridgeClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateTargets")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &runtime)
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
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("UpdateTargets failed, response: %v", response))
		}
		d.SetPartial("targets")
	}

	update = false
	updateRuleReq := map[string]interface{}{
		"EventBusName": parts[0],
		"RuleName":     parts[1],
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		updateRuleReq["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("filter_pattern") {
		update = true
	}
	updateRuleReq["FilterPattern"] = d.Get("filter_pattern")
	if update {
		action := "UpdateRule"
		conn, err := client.NewEventbridgeClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateRule")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, updateRuleReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateRuleReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("UpdateRule failed, response: %v", response))
		}
		d.SetPartial("description")
		d.SetPartial("filter_pattern")
	}
	if d.HasChange("status") {
		object, err := eventbridgeService.DescribeEventBridgeRule(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "DISABLE" {
				request := map[string]interface{}{
					"EventBusName": parts[0],
					"RuleName":     parts[1],
				}
				action := "DisableRule"
				conn, err := client.NewEventbridgeClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				if fmt.Sprint(response["Code"]) != "Success" {
					return WrapError(fmt.Errorf("DisableRule failed, response: %v", response))
				}
				stateConf := BuildStateConf([]string{}, []string{"DISABLE"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, eventbridgeService.EventBridgeRuleStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "ENABLE" {
				request := map[string]interface{}{
					"EventBusName": parts[0],
					"RuleName":     parts[1],
				}
				action := "EnableRule"
				conn, err := client.NewEventbridgeClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				if fmt.Sprint(response["Code"]) != "Success" {
					return WrapError(fmt.Errorf("EnableRule failed, response: %v", response))
				}
				stateConf := BuildStateConf([]string{}, []string{"ENABLE"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, eventbridgeService.EventBridgeRuleStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudEventBridgeRuleRead(d, meta)
}

func resourceAlicloudEventBridgeRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	eventbridgeService := EventbridgeService{client}
	action := "DeleteRule"
	var response map[string]interface{}
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"EventBusName": parts[0],
		"RuleName":     parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"EventRuleNotExisted"}) {
		return nil
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("DeleteRule failed, response: %v", response))
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, eventbridgeService.EventBridgeRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
