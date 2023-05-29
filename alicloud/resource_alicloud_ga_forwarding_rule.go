package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudGaForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaForwardingRuleCreate,
		Read:   resourceAlicloudGaForwardingRuleRead,
		Update: resourceAlicloudGaForwardingRuleUpdate,
		Delete: resourceAlicloudGaForwardingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"forwarding_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_conditions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_condition_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"Host", "Path"}, false),
						},
						"path_config": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"host_config": {
							Type:     schema.TypeSet,
							MinItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"rule_actions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"rule_action_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rule_action_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"forward_group_config": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_group_tuples": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"endpoint_group_id": {
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
			"forwarding_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"forwarding_rule_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateForwardingRules"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["ListenerId"] = d.Get("listener_id")
	forwardingRule := make(map[string]interface{})

	if v, ok := d.GetOk("priority"); ok {
		forwardingRule["Priority"] = v
	}

	ruleConditions := d.Get("rule_conditions").(*schema.Set).List()
	ruleConditionsMap := make([]map[string]interface{}, 0)
	for _, ruleCondition := range ruleConditions {
		ruleCondition := ruleCondition.(map[string]interface{})
		ruleConditionMap := map[string]interface{}{}
		ruleConditionMap["RuleConditionType"] = ruleCondition["rule_condition_type"]
		if len(ruleCondition["path_config"].(*schema.Set).List()) > 0 {
			ruleConditionMap["PathConfig"] = map[string]interface{}{
				"Values": ruleCondition["path_config"].(*schema.Set).List()[0].(map[string]interface{})["values"],
			}
		}
		if len(ruleCondition["host_config"].(*schema.Set).List()) > 0 {
			ruleConditionMap["HostConfig"] = map[string]interface{}{
				"Values": ruleCondition["host_config"].(*schema.Set).List()[0].(map[string]interface{})["values"],
			}
		}
		ruleConditionsMap = append(ruleConditionsMap, ruleConditionMap)
	}
	forwardingRule["RuleConditions"] = ruleConditionsMap

	ruleActions := d.Get("rule_actions").(*schema.Set).List()
	ruleActionsMap := make([]map[string]interface{}, 0)
	for _, ruleAction := range ruleActions {
		ruleAction := ruleAction.(map[string]interface{})
		ruleActionMap := map[string]interface{}{}
		ruleActionMap["Order"] = ruleAction["order"]
		ruleActionMap["RuleActionType"] = ruleAction["rule_action_type"]

		if v, ok := ruleAction["rule_action_value"]; ok {
			ruleActionMap["RuleActionValue"] = v
		}

		if v, ok := ruleAction["forward_group_config"]; ok {
			forwardGroupConfigMap := map[string]interface{}{}
			for _, forwardGroupConfigList := range v.(*schema.Set).List() {
				forwardGroupConfigArg := forwardGroupConfigList.(map[string]interface{})
				if serverGroupTuples, ok := forwardGroupConfigArg["server_group_tuples"]; ok {
					serverGroupTuplesMaps := make([]map[string]interface{}, 0)
					for _, serverGroupTuplesList := range serverGroupTuples.(*schema.Set).List() {
						serverGroupTuplesArg := serverGroupTuplesList.(map[string]interface{})
						serverGroupTuplesMap := map[string]interface{}{}

						if endpointGroupId, ok := serverGroupTuplesArg["endpoint_group_id"]; ok {
							serverGroupTuplesMap["EndpointGroupId"] = endpointGroupId
						}

						serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
					}

					forwardGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMaps
				}
			}
			ruleActionMap["ForwardGroupConfig"] = forwardGroupConfigMap
		}

		ruleActionsMap = append(ruleActionsMap, ruleActionMap)
	}
	forwardingRule["RuleActions"] = ruleActionsMap

	if v, ok := d.GetOk("forwarding_rule_name"); ok {
		forwardingRule["ForwardingRuleName"] = v
	}

	request["ForwardingRules"] = []interface{}{forwardingRule}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken(action)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"StateError.Accelerator"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_forwarding_rule", action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ForwardingRules", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	response = v.([]interface{})[0].(map[string]interface{})

	d.SetId(fmt.Sprintf("%s:%s:%s", request["AcceleratorId"].(string), request["ListenerId"].(string), fmt.Sprint(response["ForwardingRuleId"])))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaForwardingRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaForwardingRuleRead(d, meta)
}

func resourceAlicloudGaForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	object, err := gaService.DescribeGaForwardingRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_ip_set gaService.DescribeGaForwardingRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	d.Set("accelerator_id", parts[0])
	d.Set("listener_id", object["ListenerId"])
	d.Set("priority", object["Priority"])
	d.Set("forwarding_rule_id", object["ForwardingRuleId"])
	d.Set("forwarding_rule_name", object["ForwardingRuleName"])
	d.Set("forwarding_rule_status", object["ForwardingRuleStatus"])

	ruleConditionsMap := make([]map[string]interface{}, 0)
	for _, ruleCondition := range object["RuleConditions"].([]interface{}) {
		ruleCondition := ruleCondition.(map[string]interface{})
		ruleConditionMap := map[string]interface{}{}
		ruleConditionMap["rule_condition_type"] = ruleCondition["RuleConditionType"]
		if ruleCondition["PathConfig"].(map[string]interface{})["Values"] != nil {
			ruleConditionMap["path_config"] = []map[string]interface{}{
				{
					"values": ruleCondition["PathConfig"].(map[string]interface{})["Values"],
				},
			}
		}
		if ruleCondition["HostConfig"].(map[string]interface{})["Values"] != nil {
			ruleConditionMap["host_config"] = []map[string]interface{}{
				{
					"values": ruleCondition["HostConfig"].(map[string]interface{})["Values"],
				},
			}
		}
		ruleConditionsMap = append(ruleConditionsMap, ruleConditionMap)
	}
	d.Set("rule_conditions", ruleConditionsMap)

	if ruleActionsList, ok := object["RuleActions"]; ok {
		ruleActionsMap := make([]map[string]interface{}, 0)
		for _, ruleActions := range ruleActionsList.([]interface{}) {
			ruleActionArg := ruleActions.(map[string]interface{})
			ruleActionMap := map[string]interface{}{}

			if ruleActionOrder, ok := ruleActionArg["Order"]; ok {
				ruleActionMap["order"] = ruleActionOrder
			}

			if ruleActionType, ok := ruleActionArg["RuleActionType"]; ok {
				ruleActionMap["rule_action_type"] = ruleActionType
			}

			if forwardGroupConfig, ok := ruleActionArg["ForwardGroupConfig"]; ok {
				forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
				if len(forwardGroupConfigArg) > 0 {
					serverGroupTuplesMaps := make([]map[string]interface{}, 0)
					if forwardGroupConfigArgs, ok := forwardGroupConfigArg["ServerGroupTuples"].([]interface{}); ok {
						for _, serverGroupTuples := range forwardGroupConfigArgs {
							serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
							serverGroupTuplesMap := map[string]interface{}{}
							serverGroupTuplesMap["endpoint_group_id"] = serverGroupTuplesArg["EndpointGroupId"]
							serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
						}
					}
					if len(serverGroupTuplesMaps) > 0 {
						forwardGroupConfigMaps := make([]map[string]interface{}, 0)
						forwardGroupConfigMap := map[string]interface{}{}
						forwardGroupConfigMap["server_group_tuples"] = serverGroupTuplesMaps
						forwardGroupConfigMaps = append(forwardGroupConfigMaps, forwardGroupConfigMap)
						ruleActionMap["forward_group_config"] = forwardGroupConfigMaps
					}
				} else {
					if ruleActionValue, ok := ruleActionArg["RuleActionValue"]; ok {
						ruleActionMap["rule_action_value"] = ruleActionValue
					}
				}
			}
			ruleActionsMap = append(ruleActionsMap, ruleActionMap)
		}

		d.Set("rule_actions", ruleActionsMap)
	}

	return nil
}

func resourceAlicloudGaForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"AcceleratorId": parts[0],
		"ListenerId":    parts[1],
	}

	forwardingRule := make(map[string]interface{})
	forwardingRule["ForwardingRuleId"] = parts[2]

	if d.HasChange("priority") {
		update = true
	}
	forwardingRule["Priority"] = d.Get("priority")

	if d.HasChange("forwarding_rule_name") {
		update = true
	}
	if v, ok := d.GetOk("forwarding_rule_name"); ok {
		forwardingRule["ForwardingRuleName"] = v
	}

	if d.HasChange("rule_conditions") {
		update = true
		ruleConditions := d.Get("rule_conditions").(*schema.Set).List()
		ruleConditionsMap := make([]map[string]interface{}, 0)
		for _, ruleCondition := range ruleConditions {
			ruleCondition := ruleCondition.(map[string]interface{})
			ruleConditionMap := map[string]interface{}{}
			ruleConditionMap["RuleConditionType"] = ruleCondition["rule_condition_type"]
			if len(ruleCondition["path_config"].(*schema.Set).List()) > 0 {
				ruleConditionMap["PathConfig"] = map[string]interface{}{
					"Values": ruleCondition["path_config"].(*schema.Set).List()[0].(map[string]interface{})["values"],
				}
			}
			if len(ruleCondition["host_config"].(*schema.Set).List()) > 0 {
				ruleConditionMap["HostConfig"] = map[string]interface{}{
					"Values": ruleCondition["host_config"].(*schema.Set).List()[0].(map[string]interface{})["values"],
				}
			}
			ruleConditionsMap = append(ruleConditionsMap, ruleConditionMap)
		}
		forwardingRule["RuleConditions"] = ruleConditionsMap
	}

	if d.HasChange("rule_actions") {
		update = true
		ruleActions := d.Get("rule_actions").(*schema.Set).List()
		ruleActionsMap := make([]map[string]interface{}, 0)
		for _, ruleAction := range ruleActions {
			ruleAction := ruleAction.(map[string]interface{})
			ruleActionMap := map[string]interface{}{}
			ruleActionMap["Order"] = ruleAction["order"]
			ruleActionMap["RuleActionType"] = ruleAction["rule_action_type"]

			if v, ok := ruleAction["rule_action_value"]; ok {
				ruleActionMap["RuleActionValue"] = v
			}

			if v, ok := ruleAction["forward_group_config"]; ok {
				forwardGroupConfigMap := map[string]interface{}{}
				for _, forwardGroupConfigList := range v.(*schema.Set).List() {
					forwardGroupConfigArg := forwardGroupConfigList.(map[string]interface{})
					if serverGroupTuples, ok := forwardGroupConfigArg["server_group_tuples"]; ok {
						serverGroupTuplesMaps := make([]map[string]interface{}, 0)
						for _, serverGroupTuplesList := range serverGroupTuples.(*schema.Set).List() {
							serverGroupTuplesArg := serverGroupTuplesList.(map[string]interface{})
							serverGroupTuplesMap := map[string]interface{}{}

							if endpointGroupId, ok := serverGroupTuplesArg["endpoint_group_id"]; ok {
								serverGroupTuplesMap["EndpointGroupId"] = endpointGroupId
							}

							serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
						}

						forwardGroupConfigMap["ServerGroupTuples"] = serverGroupTuplesMaps
					}
				}
				ruleActionMap["ForwardGroupConfig"] = forwardGroupConfigMap
			}

			ruleActionsMap = append(ruleActionsMap, ruleActionMap)
		}
		forwardingRule["RuleActions"] = ruleActionsMap
	}

	request["ForwardingRules"] = []interface{}{forwardingRule}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("UpdateForwardingRules")

	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	if update {
		action := "UpdateForwardingRules"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken(action)
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.ForwardingRule"}) {
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

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaForwardingRuleStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudGaForwardingRuleRead(d, meta)
}

func resourceAlicloudGaForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteForwardingRules"
	var response map[string]interface{}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AcceleratorId": parts[0],
		"ListenerId":    parts[1],
	}
	request["ForwardingRuleIds"] = []string{parts[2]}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken(action)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.ForwardingRule"}) {
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

	return nil
}
