package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudConfigAggregateCompliancePack() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudConfigAggregateCompliancePackCreate,
		Read:   resourceAliCloudConfigAggregateCompliancePackRead,
		Update: resourceAliCloudConfigAggregateCompliancePackUpdate,
		Delete: resourceAliCloudConfigAggregateCompliancePackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aggregator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aggregate_compliance_pack_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"risk_level": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3}),
			},
			"compliance_pack_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"config_rule_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"config_rules"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_rule_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"config_rules": {
				Type:          schema.TypeSet,
				Optional:      true,
				Deprecated:    "Field `config_rules` has been deprecated from provider version 1.141.0. New field `config_rule_ids` instead.",
				ConflictsWith: []string{"config_rule_ids"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"managed_rule_identifier": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config_rule_parameters": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameter_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"parameter_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"aggregator_compliance_pack_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudConfigAggregateCompliancePackCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	var err error
	action := "CreateAggregateCompliancePack"
	request := make(map[string]interface{})

	request["ClientToken"] = buildClientToken("CreateAggregateCompliancePack")
	request["AggregatorId"] = d.Get("aggregator_id")
	request["CompliancePackName"] = d.Get("aggregate_compliance_pack_name")
	request["Description"] = d.Get("description")
	request["RiskLevel"] = d.Get("risk_level")

	if v, ok := d.GetOk("compliance_pack_template_id"); ok {
		request["CompliancePackTemplateId"] = v
	}

	configRulesMaps := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("config_rule_ids"); ok {
		for _, configRuleIds := range v.(*schema.Set).List() {
			configRuleIdsMap := map[string]interface{}{}
			configRuleIdsArg := configRuleIds.(map[string]interface{})

			if configRuleId, ok := configRuleIdsArg["config_rule_id"]; ok {
				configRuleIdsMap["ConfigRuleId"] = configRuleId
			}

			configRulesMaps = append(configRulesMaps, configRuleIdsMap)
		}
	} else {
		if v, ok := d.GetOk("config_rules"); ok {
			for _, configRules := range v.(*schema.Set).List() {
				configRulesMap := map[string]interface{}{}
				configRulesArg := configRules.(map[string]interface{})

				configRulesMap["ManagedRuleIdentifier"] = configRulesArg["managed_rule_identifier"]

				if configRuleParameters, ok := configRulesArg["config_rule_parameters"]; ok {
					configRuleParametersMaps := make([]map[string]interface{}, 0)
					configRuleParametersMap := map[string]interface{}{}

					for _, configRuleParametersList := range configRuleParameters.(*schema.Set).List() {
						configRuleParametersArg := configRuleParametersList.(map[string]interface{})

						if parameterName, ok := configRuleParametersArg["parameter_name"]; ok {
							configRuleParametersMap["ParameterName"] = parameterName
						}

						if parameterValue, ok := configRuleParametersArg["parameter_value"]; ok {
							configRuleParametersMap["ParameterValue"] = parameterValue
						}

						configRuleParametersMaps = append(configRuleParametersMaps, configRuleParametersMap)
					}

					configRulesMap["ConfigRuleParameters"] = configRuleParametersMaps
				}

				configRulesMaps = append(configRulesMaps, configRulesMap)
			}
		}
	}

	if len(configRulesMaps) > 0 {
		configRulesJson, err := convertListMapToJsonString(configRulesMaps)
		if err != nil {
			return WrapError(err)
		}

		request["ConfigRules"] = configRulesJson
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_aggregate_compliance_pack", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["AggregatorId"], response["CompliancePackId"]))

	stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, configService.ConfigAggregateCompliancePackStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudConfigAggregateCompliancePackRead(d, meta)
}

func resourceAliCloudConfigAggregateCompliancePackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}

	object, err := configService.DescribeConfigAggregateCompliancePack(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_aggregate_compliance_pack configService.DescribeConfigAggregateCompliancePack Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("aggregator_id", object["AggregatorId"])
	d.Set("aggregator_compliance_pack_id", object["CompliancePackId"])
	d.Set("aggregate_compliance_pack_name", object["CompliancePackName"])
	d.Set("description", object["Description"])
	d.Set("risk_level", formatInt(object["RiskLevel"]))
	d.Set("compliance_pack_template_id", object["CompliancePackTemplateId"])
	d.Set("status", object["Status"])

	if _, ok := d.GetOk("config_rules"); ok {
		if configRulesList, ok := object["ConfigRules"]; ok {
			configRulesMaps := make([]map[string]interface{}, 0)
			for _, configRules := range configRulesList.([]interface{}) {
				configRulesArg := configRules.(map[string]interface{})
				configRulesMap := map[string]interface{}{}

				if managedRuleIdentifier, ok := configRulesArg["ManagedRuleIdentifier"]; ok {
					configRulesMap["managed_rule_identifier"] = managedRuleIdentifier
				}

				if configRuleParameters, ok := configRulesArg["ConfigRuleParameters"]; ok {
					configRuleParametersMaps := make([]map[string]interface{}, 0)
					for _, configRuleParametersList := range configRuleParameters.([]interface{}) {
						configRuleParametersArg := configRuleParametersList.(map[string]interface{})
						configRuleParametersMap := map[string]interface{}{}

						if logType, ok := configRuleParametersArg["ParameterName"]; ok {
							configRuleParametersMap["parameter_name"] = logType
						}

						if logType, ok := configRuleParametersArg["ParameterValue"]; ok {
							configRuleParametersMap["parameter_value"] = logType
						}

						configRuleParametersMaps = append(configRuleParametersMaps, configRuleParametersMap)
					}

					configRulesMap["config_rule_parameters"] = configRuleParametersMaps
				}

				configRulesMaps = append(configRulesMaps, configRulesMap)
			}

			d.Set("config_rules", configRulesMaps)
		}
	} else {
		if configRulesList, ok := object["ConfigRules"]; ok {
			configRulesMaps := make([]map[string]interface{}, 0)
			for _, configRules := range configRulesList.([]interface{}) {
				configRulesArg := configRules.(map[string]interface{})
				configRulesMap := map[string]interface{}{}

				if configRuleId, ok := configRulesArg["ConfigRuleId"]; ok {
					configRulesMap["config_rule_id"] = configRuleId
				}

				configRulesMaps = append(configRulesMaps, configRulesMap)
			}

			d.Set("config_rule_ids", configRulesMaps)
		}
	}

	return nil
}

func resourceAliCloudConfigAggregateCompliancePackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	d.Partial(true)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	update := false
	request := map[string]interface{}{
		"ClientToken":      buildClientToken("UpdateAggregateCompliancePack"),
		"AggregatorId":     parts[0],
		"CompliancePackId": parts[1],
	}

	if d.HasChange("aggregate_compliance_pack_name") {
		update = true
	}
	request["CompliancePackName"] = d.Get("aggregate_compliance_pack_name")

	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")

	if d.HasChange("risk_level") {
		update = true
	}
	request["RiskLevel"] = d.Get("risk_level")

	if d.HasChange("config_rules") {
		update = true
	}
	if v, ok := d.GetOk("config_rules"); ok {
		configRulesMaps := make([]map[string]interface{}, 0)
		for _, configRules := range v.(*schema.Set).List() {
			configRulesMap := map[string]interface{}{}
			configRulesArg := configRules.(map[string]interface{})

			configRulesMap["ManagedRuleIdentifier"] = configRulesArg["managed_rule_identifier"]

			if configRuleParameters, ok := configRulesArg["config_rule_parameters"]; ok {
				configRuleParametersMaps := make([]map[string]interface{}, 0)
				configRuleParametersMap := map[string]interface{}{}

				for _, configRuleParametersList := range configRuleParameters.(*schema.Set).List() {
					configRuleParametersArg := configRuleParametersList.(map[string]interface{})

					if parameterName, ok := configRuleParametersArg["parameter_name"]; ok {
						configRuleParametersMap["ParameterName"] = parameterName
					}

					if parameterValue, ok := configRuleParametersArg["parameter_value"]; ok {
						configRuleParametersMap["ParameterValue"] = parameterValue
					}

					configRuleParametersMaps = append(configRuleParametersMaps, configRuleParametersMap)
				}

				configRulesMap["ConfigRuleParameters"] = configRuleParametersMaps
			}

			configRulesMaps = append(configRulesMaps, configRulesMap)
		}

		configRulesJson, err := convertListMapToJsonString(configRulesMaps)
		if err != nil {
			return WrapError(err)
		}

		request["ConfigRules"] = configRulesJson
	}

	if update {
		action := "UpdateAggregateCompliancePack"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, true)
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

		stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, configService.ConfigAggregateCompliancePackStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("aggregate_compliance_pack_name")
		d.SetPartial("description")
		d.SetPartial("risk_level")
		d.SetPartial("config_rules")
	}

	if d.HasChange("config_rule_ids") {
		oldConfigRuleIds, newConfigRuleIds := d.GetChange("config_rule_ids")
		remove := oldConfigRuleIds.(*schema.Set).Difference(newConfigRuleIds.(*schema.Set)).List()
		create := newConfigRuleIds.(*schema.Set).Difference(oldConfigRuleIds.(*schema.Set)).List()
		if len(remove) > 0 {
			action := "DetachAggregateConfigRuleToCompliancePack"

			detachAggregateConfigRuleReq := map[string]interface{}{
				"AggregatorId":     parts[0],
				"CompliancePackId": parts[1],
			}

			configRuleIds := make([]interface{}, 0)
			for _, configRuleIdsList := range remove {
				configRuleIdsArg := configRuleIdsList.(map[string]interface{})

				if configRuleId, ok := configRuleIdsArg["config_rule_id"]; ok {
					configRuleIds = append(configRuleIds, configRuleId)
				}
			}

			detachAggregateConfigRuleReq["ConfigRuleIds"] = convertListToCommaSeparate(configRuleIds)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Config", "2020-09-07", action, nil, detachAggregateConfigRuleReq, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			addDebug(action, response, detachAggregateConfigRuleReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		if len(create) > 0 {
			action := "AttachAggregateConfigRuleToCompliancePack"

			attachAggregateConfigRuleReq := map[string]interface{}{
				"AggregatorId":     parts[0],
				"CompliancePackId": parts[1],
			}

			configRuleIds := make([]interface{}, 0)
			for _, configRuleIdsList := range create {
				configRuleIdsArg := configRuleIdsList.(map[string]interface{})

				if configRuleId, ok := configRuleIdsArg["config_rule_id"]; ok {
					configRuleIds = append(configRuleIds, configRuleId)
				}
			}

			attachAggregateConfigRuleReq["ConfigRuleIds"] = convertListToCommaSeparate(configRuleIds)

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Config", "2020-09-07", action, nil, attachAggregateConfigRuleReq, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, attachAggregateConfigRuleReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		d.SetPartial("config_rule_ids")
	}

	d.Partial(false)

	return resourceAliCloudConfigAggregateCompliancePackRead(d, meta)
}

func resourceAliCloudConfigAggregateCompliancePackDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAggregateCompliancePacks"
	var response map[string]interface{}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":       buildClientToken("DeleteAggregateCompliancePacks"),
		"AggregatorId":      parts[0],
		"CompliancePackIds": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Invalid.AggregatorId.Value", "Invalid.CompliancePackId.Value"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
