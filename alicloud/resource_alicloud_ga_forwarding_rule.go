package alicloud

import (
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGaForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaForwardingRuleCreate,
		Read:   resourceAliCloudGaForwardingRuleRead,
		Update: resourceAliCloudGaForwardingRuleUpdate,
		Delete: resourceAliCloudGaForwardingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
							ValidateFunc: StringInSlice([]string{"Host", "Path", "RequestHeader", "Query", "Method", "Cookie", "SourceIP"}, false),
						},
						"rule_condition_value": {
							Type:     schema.TypeString,
							Optional: true,
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"ForwardGroup", "Redirect", "FixResponse", "Rewrite", "AddHeader", "RemoveHeader", "Drop"}, false),
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

func resourceAliCloudGaForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateForwardingRules"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateForwardingRules")
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["ListenerId"] = d.Get("listener_id")

	forwardingRulesMaps := make([]map[string]interface{}, 0)
	forwardingRulesMap := map[string]interface{}{}

	if v, ok := d.GetOk("priority"); ok {
		forwardingRulesMap["Priority"] = v
	}

	if v, ok := d.GetOk("forwarding_rule_name"); ok {
		forwardingRulesMap["ForwardingRuleName"] = v
	}

	ruleConditions := d.Get("rule_conditions").(*schema.Set).List()
	ruleConditionsMaps := make([]map[string]interface{}, 0)
	for _, ruleConditionsList := range ruleConditions {
		ruleConditionsMap := map[string]interface{}{}
		ruleConditionsArg := ruleConditionsList.(map[string]interface{})

		ruleConditionsMap["RuleConditionType"] = ruleConditionsArg["rule_condition_type"]

		if ruleConditionValue, ok := ruleConditionsArg["rule_condition_value"]; ok {
			ruleConditionsMap["RuleConditionValue"] = ruleConditionValue
		}

		if pathConfig, ok := ruleConditionsArg["path_config"]; ok {
			pathConfigMap := map[string]interface{}{}
			for _, pathConfigList := range pathConfig.(*schema.Set).List() {
				pathConfigArg := pathConfigList.(map[string]interface{})

				if values, ok := pathConfigArg["values"]; ok {
					pathConfigMap["Values"] = values
				}
			}

			ruleConditionsMap["PathConfig"] = pathConfigMap
		}

		if hostConfig, ok := ruleConditionsArg["host_config"]; ok {
			hostConfigMap := map[string]interface{}{}
			for _, hostConfigList := range hostConfig.(*schema.Set).List() {
				hostConfigArg := hostConfigList.(map[string]interface{})

				if values, ok := hostConfigArg["values"]; ok {
					hostConfigMap["Values"] = values
				}
			}

			ruleConditionsMap["HostConfig"] = hostConfigMap
		}

		ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
	}

	forwardingRulesMap["RuleConditions"] = ruleConditionsMaps

	ruleActions := d.Get("rule_actions").(*schema.Set).List()
	ruleActionsMaps := make([]map[string]interface{}, 0)
	for _, ruleActionsList := range ruleActions {
		ruleActionsMap := map[string]interface{}{}
		ruleActionsArg := ruleActionsList.(map[string]interface{})

		ruleActionsMap["Order"] = ruleActionsArg["order"]
		ruleActionsMap["RuleActionType"] = ruleActionsArg["rule_action_type"]

		if ruleActionValue, ok := ruleActionsArg["rule_action_value"]; ok && fmt.Sprint(ruleActionValue) != "" && fmt.Sprint(ruleActionsMap["RuleActionType"]) != "Drop" {
			ruleActionsMap["RuleActionValue"] = ruleActionValue
		}

		if forwardGroupConfig, ok := ruleActionsArg["forward_group_config"]; ok {
			forwardGroupConfigMap := map[string]interface{}{}
			for _, forwardGroupConfigList := range forwardGroupConfig.(*schema.Set).List() {
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

			ruleActionsMap["ForwardGroupConfig"] = forwardGroupConfigMap
		}

		ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
	}

	forwardingRulesMap["RuleActions"] = ruleActionsMaps

	forwardingRulesMaps = append(forwardingRulesMaps, forwardingRulesMap)
	request["ForwardingRules"] = forwardingRulesMaps

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator"}) || NeedRetry(err) {
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

	return resourceAliCloudGaForwardingRuleRead(d, meta)
}

func resourceAliCloudGaForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("forwarding_rule_name", object["ForwardingRuleName"])
	d.Set("forwarding_rule_id", object["ForwardingRuleId"])
	d.Set("forwarding_rule_status", object["ForwardingRuleStatus"])

	ruleConditionsMap := make([]map[string]interface{}, 0)
	isPathConfigExist := false
	isHostConfigExist := false

	for _, ruleCondition := range object["RuleConditions"].([]interface{}) {
		ruleConditionArg := ruleCondition.(map[string]interface{})
		ruleConditionMap := map[string]interface{}{}
		ruleConditionMap["rule_condition_type"] = ruleConditionArg["RuleConditionType"]

		if ruleConditionArg["PathConfig"].(map[string]interface{})["Values"] != nil {
			isPathConfigExist = true

			ruleConditionMap["path_config"] = []map[string]interface{}{
				{
					"values": ruleConditionArg["PathConfig"].(map[string]interface{})["Values"],
				},
			}
		}

		if ruleConditionArg["HostConfig"].(map[string]interface{})["Values"] != nil {
			isHostConfigExist = true

			ruleConditionMap["host_config"] = []map[string]interface{}{
				{
					"values": ruleConditionArg["HostConfig"].(map[string]interface{})["Values"],
				},
			}
		}

		if !isPathConfigExist && !isHostConfigExist {
			if ruleConditionValue, ok := ruleConditionArg["RuleConditionValue"]; ok {
				ruleConditionMap["rule_condition_value"] = ruleConditionValue
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

func resourceAliCloudGaForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"ClientToken":   buildClientToken("UpdateForwardingRules"),
		"AcceleratorId": parts[0],
		"ListenerId":    parts[1],
	}

	forwardingRulesMaps := make([]map[string]interface{}, 0)
	forwardingRulesMap := map[string]interface{}{}

	forwardingRulesMap["ForwardingRuleId"] = parts[2]

	if d.HasChange("priority") {
		update = true
	}
	forwardingRulesMap["Priority"] = d.Get("priority")

	if d.HasChange("forwarding_rule_name") {
		update = true
	}
	if v, ok := d.GetOk("forwarding_rule_name"); ok {
		forwardingRulesMap["ForwardingRuleName"] = v
	} else {
		forwardingRulesMap["ForwardingRuleName"] = tea.String("")
	}

	if d.HasChange("rule_conditions") {
		update = true
	}
	ruleConditions := d.Get("rule_conditions").(*schema.Set).List()
	ruleConditionsMaps := make([]map[string]interface{}, 0)
	for _, ruleConditionsList := range ruleConditions {
		ruleConditionsMap := map[string]interface{}{}
		ruleConditionsArg := ruleConditionsList.(map[string]interface{})

		ruleConditionsMap["RuleConditionType"] = ruleConditionsArg["rule_condition_type"]

		if ruleConditionValue, ok := ruleConditionsArg["rule_condition_value"]; ok {
			ruleConditionsMap["RuleConditionValue"] = ruleConditionValue
		}

		if pathConfig, ok := ruleConditionsArg["path_config"]; ok {
			pathConfigMap := map[string]interface{}{}
			for _, pathConfigList := range pathConfig.(*schema.Set).List() {
				pathConfigArg := pathConfigList.(map[string]interface{})

				if values, ok := pathConfigArg["values"]; ok {
					pathConfigMap["Values"] = values
				}
			}

			ruleConditionsMap["PathConfig"] = pathConfigMap
		}

		if hostConfig, ok := ruleConditionsArg["host_config"]; ok {
			hostConfigMap := map[string]interface{}{}
			for _, hostConfigList := range hostConfig.(*schema.Set).List() {
				hostConfigArg := hostConfigList.(map[string]interface{})

				if values, ok := hostConfigArg["values"]; ok {
					hostConfigMap["Values"] = values
				}
			}

			ruleConditionsMap["HostConfig"] = hostConfigMap
		}

		ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
	}

	forwardingRulesMap["RuleConditions"] = ruleConditionsMaps

	if d.HasChange("rule_actions") {
		update = true
	}
	ruleActions := d.Get("rule_actions").(*schema.Set).List()
	ruleActionsMaps := make([]map[string]interface{}, 0)
	for _, ruleActionsList := range ruleActions {
		ruleActionsMap := map[string]interface{}{}
		ruleActionsArg := ruleActionsList.(map[string]interface{})

		ruleActionsMap["Order"] = ruleActionsArg["order"]
		ruleActionsMap["RuleActionType"] = ruleActionsArg["rule_action_type"]

		if ruleActionValue, ok := ruleActionsArg["rule_action_value"]; ok && fmt.Sprint(ruleActionValue) != "" && fmt.Sprint(ruleActionsMap["RuleActionType"]) != "Drop" {
			ruleActionsMap["RuleActionValue"] = ruleActionValue
		}

		if forwardGroupConfig, ok := ruleActionsArg["forward_group_config"]; ok {
			forwardGroupConfigMap := map[string]interface{}{}
			for _, forwardGroupConfigList := range forwardGroupConfig.(*schema.Set).List() {
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

			ruleActionsMap["ForwardGroupConfig"] = forwardGroupConfigMap
		}

		ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
	}

	forwardingRulesMap["RuleActions"] = ruleActionsMaps

	forwardingRulesMaps = append(forwardingRulesMaps, forwardingRulesMap)
	request["ForwardingRules"] = forwardingRulesMaps

	if update {
		action := "UpdateForwardingRules"
		var err error

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.ForwardingRule"}) || NeedRetry(err) {
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

	return resourceAliCloudGaForwardingRuleRead(d, meta)
}

func resourceAliCloudGaForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteForwardingRules"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":          client.RegionId,
		"ClientToken":       buildClientToken("DeleteForwardingRules"),
		"AcceleratorId":     parts[0],
		"ListenerId":        parts[1],
		"ForwardingRuleIds": []string{parts[2]},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.ForwardingRule"}) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaForwardingRuleStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
