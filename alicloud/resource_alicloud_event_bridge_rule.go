package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEventBridgeRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEventBridgeRuleCreate,
		Read:   resourceAliCloudEventBridgeRuleRead,
		Update: resourceAliCloudEventBridgeRuleUpdate,
		Delete: resourceAliCloudEventBridgeRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"event_bus_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"filter_pattern": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ENABLE", "DISABLE"}, false),
			},
			"targets": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: StringInSlice([]string{
								"acs.alikafka",
								"acs.api.destination",
								"acs.arms.loki",
								"acs.datahub",
								"acs.dingtalk",
								"acs.eventbridge",
								"acs.eventbridge.olap",
								"acs.eventbus.SLSCloudLens",
								"acs.fc.function",
								"acs.fnf",
								"acs.k8s",
								"acs.mail",
								"acs.mns.queue",
								"acs.mns.topic",
								"acs.openapi",
								"acs.rabbitmq",
								"acs.rds.mysql",
								"acs.rocketmq",
								"acs.sae",
								"acs.sls",
								"acs.sms",
								"http",
								"https",
								"mysql"}, false),
						},
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"push_retry_strategy": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"BACKOFF_RETRY", "EXPONENTIAL_DECAY_RETRY"}, false),
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
						"param_list": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"form": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: StringInSlice([]string{"ORIGINAL", "TEMPLATE", "JSONPATH", "CONSTANT"}, false),
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
					},
				},
			},
		},
	}
}

func resourceAliCloudEventBridgeRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateRule"
	request := make(map[string]interface{})
	var err error

	request["EventBusName"] = d.Get("event_bus_name")
	request["RuleName"] = d.Get("rule_name")
	request["FilterPattern"] = d.Get("filter_pattern")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	targets := d.Get("targets")
	targetsMaps := make([]map[string]interface{}, 0)
	for _, targetsList := range targets.([]interface{}) {
		targetsMap := map[string]interface{}{}
		targetsArg := targetsList.(map[string]interface{})

		targetsMap["Id"] = targetsArg["target_id"]
		targetsMap["Type"] = targetsArg["type"]
		targetsMap["Endpoint"] = targetsArg["endpoint"]

		if pushRetryStrategy, ok := targetsArg["push_retry_strategy"]; ok && fmt.Sprint(pushRetryStrategy) != "" {
			targetsMap["PushRetryStrategy"] = pushRetryStrategy
		}

		if deadLetterQueue, ok := targetsArg["dead_letter_queue"]; ok {
			deadLetterQueueMap := map[string]interface{}{}
			for _, deadLetterQueueList := range deadLetterQueue.(*schema.Set).List() {
				deadLetterQueueArg := deadLetterQueueList.(map[string]interface{})

				if arn, ok := deadLetterQueueArg["arn"]; ok {
					deadLetterQueueMap["Arn"] = arn
				}
			}

			targetsMap["DeadLetterQueue"] = deadLetterQueueMap
		}

		paramList := targetsArg["param_list"]
		paramListMaps := make([]map[string]interface{}, 0)
		for _, paramListArgList := range paramList.(*schema.Set).List() {
			paramListMap := map[string]interface{}{}
			paramListArg := paramListArgList.(map[string]interface{})

			paramListMap["ResourceKey"] = paramListArg["resource_key"]
			paramListMap["Form"] = paramListArg["form"]

			if paramListMap["Form"] == "TEMPLATE" {
				if template, ok := paramListArg["template"]; ok {
					paramListMap["Template"] = template
				}

				if value, ok := paramListArg["value"]; ok {
					paramListMap["Value"] = value
				}
			}

			if paramListMap["Form"] == "JSONPATH" || paramListMap["Form"] == "CONSTANT" {
				if value, ok := paramListArg["value"]; ok {
					paramListMap["Value"] = value
				}
			}

			paramListMaps = append(paramListMaps, paramListMap)
		}

		targetsMap["ParamList"] = paramListMaps

		targetsMaps = append(targetsMaps, targetsMap)
	}

	targetsMapsJson, err := convertArrayObjectToJsonString(targetsMaps)
	if err != nil {
		return WrapError(err)
	}

	request["EventTargets"] = targetsMapsJson
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprintf("%v:%v", request["EventBusName"], request["RuleName"]))

	return resourceAliCloudEventBridgeRuleUpdate(d, meta)
}

func resourceAliCloudEventBridgeRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeService := EventbridgeService{client}

	object, err := eventBridgeService.DescribeEventBridgeRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_rule eventBridgeService.DescribeEventBridgeRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("event_bus_name", object["EventBusName"])
	d.Set("rule_name", object["RuleName"])
	d.Set("filter_pattern", object["FilterPattern"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])

	if targetsList, ok := object["Targets"]; ok {
		targetsMaps := make([]map[string]interface{}, 0)
		for _, targets := range targetsList.([]interface{}) {
			targetsArg := targets.(map[string]interface{})
			targetsMap := make(map[string]interface{})

			targetsMap["target_id"] = targetsArg["Id"]
			targetsMap["type"] = targetsArg["Type"]
			targetsMap["endpoint"] = targetsArg["Endpoint"]

			if pushRetryStrategy, ok := targetsArg["PushRetryStrategy"]; ok {
				targetsMap["push_retry_strategy"] = pushRetryStrategy
			}

			if deadLetterQueue, ok := targetsArg["DeadLetterQueue"]; ok && deadLetterQueue != nil {
				deadLetterQueueMaps := make([]map[string]interface{}, 0)
				deadLetterQueueArg := deadLetterQueue.(map[string]interface{})
				deadLetterQueueMap := map[string]interface{}{}

				if deadLetterQueueArgArn, ok := deadLetterQueueArg["Arn"]; ok {
					deadLetterQueueMap["arn"] = deadLetterQueueArgArn
				}

				if len(deadLetterQueueMap) > 0 {
					deadLetterQueueMaps = append(deadLetterQueueMaps, deadLetterQueueMap)
				}

				targetsMap["dead_letter_queue"] = deadLetterQueueMaps
			}

			if paramList, ok := targetsArg["ParamList"]; ok {
				paramListMaps := make([]map[string]interface{}, 0)
				for _, paramLists := range paramList.([]interface{}) {
					paramListMap := make(map[string]interface{})
					paramListArg := paramLists.(map[string]interface{})

					// There is an api bug that the event bridge service will return a default param which ResourceKey is IsBase64Encode and the value is false
					if fmt.Sprint(paramListArg["ResourceKey"]) == "IsBase64Encode" && fmt.Sprint(paramListArg["Value"]) == "false" {
						continue
					}

					if resourceKey, ok := paramListArg["ResourceKey"]; ok {
						paramListMap["resource_key"] = resourceKey
					}

					if form, ok := paramListArg["Form"]; ok {
						paramListMap["form"] = form
					}

					if template, ok := paramListArg["Template"]; ok {
						paramListMap["template"] = template
					}

					if value, ok := paramListArg["Value"]; ok {
						paramListMap["value"] = value
					}

					paramListMaps = append(paramListMaps, paramListMap)
				}

				targetsMap["param_list"] = paramListMaps
			}

			targetsMaps = append(targetsMaps, targetsMap)
		}

		d.Set("targets", targetsMaps)
	}

	return nil
}

func resourceAliCloudEventBridgeRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeService := EventbridgeService{client}
	var response map[string]interface{}
	d.Partial(true)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	update := false
	request := map[string]interface{}{
		"EventBusName": parts[0],
		"RuleName":     parts[1],
	}

	if !d.IsNewResource() && d.HasChange("filter_pattern") {
		update = true
	}
	request["FilterPattern"] = d.Get("filter_pattern")

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if update {
		action := "UpdateRule"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		d.SetPartial("filter_pattern")
		d.SetPartial("description")
	}

	update = false
	putTargetsReq := map[string]interface{}{
		"EventBusName": parts[0],
		"RuleName":     parts[1],
	}

	if !d.IsNewResource() && d.HasChange("targets") {
		update = true
	}
	targets := d.Get("targets")
	targetsMaps := make([]map[string]interface{}, 0)
	for _, targetsList := range targets.([]interface{}) {
		targetsMap := map[string]interface{}{}
		targetsArg := targetsList.(map[string]interface{})

		targetsMap["Id"] = targetsArg["target_id"]
		targetsMap["Type"] = targetsArg["type"]
		targetsMap["Endpoint"] = targetsArg["endpoint"]

		if pushRetryStrategy, ok := targetsArg["push_retry_strategy"]; ok && fmt.Sprint(pushRetryStrategy) != "" {
			targetsMap["PushRetryStrategy"] = pushRetryStrategy
		}

		if deadLetterQueue, ok := targetsArg["dead_letter_queue"]; ok {
			deadLetterQueueMap := map[string]interface{}{}
			for _, deadLetterQueueList := range deadLetterQueue.(*schema.Set).List() {
				deadLetterQueueArg := deadLetterQueueList.(map[string]interface{})

				if arn, ok := deadLetterQueueArg["arn"]; ok {
					deadLetterQueueMap["Arn"] = arn
				}
			}

			targetsMap["DeadLetterQueue"] = deadLetterQueueMap
		}

		paramList := targetsArg["param_list"]
		paramListMaps := make([]map[string]interface{}, 0)
		for _, paramListArgList := range paramList.(*schema.Set).List() {
			paramListMap := map[string]interface{}{}
			paramListArg := paramListArgList.(map[string]interface{})

			paramListMap["ResourceKey"] = paramListArg["resource_key"]
			paramListMap["Form"] = paramListArg["form"]

			if paramListMap["Form"] == "TEMPLATE" {
				if template, ok := paramListArg["template"]; ok {
					paramListMap["Template"] = template
				}

				if value, ok := paramListArg["value"]; ok {
					paramListMap["Value"] = value
				}
			}

			if paramListMap["Form"] == "JSONPATH" || paramListMap["Form"] == "CONSTANT" {
				if value, ok := paramListArg["value"]; ok {
					paramListMap["Value"] = value
				}
			}

			paramListMaps = append(paramListMaps, paramListMap)
		}

		targetsMap["ParamList"] = paramListMaps

		targetsMaps = append(targetsMaps, targetsMap)
	}

	targetsMapsJson, err := convertArrayObjectToJsonString(targetsMaps)
	if err != nil {
		return WrapError(err)
	}

	putTargetsReq["Targets"] = targetsMapsJson

	if update {
		action := "PutTargets"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, putTargetsReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, putTargetsReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		d.SetPartial("targets")
	}

	if d.HasChange("status") {
		object, err := eventBridgeService.DescribeEventBridgeRule(d.Id())
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
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
					return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
				}

				stateConf := BuildStateConf([]string{}, []string{"DISABLE"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, eventBridgeService.EventBridgeRuleStateRefreshFunc(d.Id(), []string{}))
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
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
					return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
				}

				stateConf := BuildStateConf([]string{}, []string{"ENABLE"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, eventBridgeService.EventBridgeRuleStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			d.SetPartial("status")
		}
	}

	d.Partial(false)

	return resourceAliCloudEventBridgeRuleRead(d, meta)
}

func resourceAliCloudEventBridgeRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventBridgeService := EventbridgeService{client}
	action := "DeleteRule"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"EventBusName": parts[0],
		"RuleName":     parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"EventRuleNotExisted"}) {
		return nil
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, eventBridgeService.EventBridgeRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
