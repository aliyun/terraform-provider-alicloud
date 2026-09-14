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
						"config_rule_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"risk_level": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntInSlice([]int{1, 2, 3}),
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
			"template_content": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tag_key_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_value_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_resource_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_resource_group_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_region_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags_scope": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"exclude_tags_scope": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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

	if v, ok := d.GetOk("template_content"); ok {
		request["TemplateContent"] = v
	}

	if v, ok := d.GetOk("tag_key_scope"); ok {
		request["TagKeyScope"] = v
	}
	if v, ok := d.GetOk("tag_value_scope"); ok {
		request["TagValueScope"] = v
	}
	if v, ok := d.GetOk("resource_ids_scope"); ok {
		request["ResourceIdsScope"] = v
	}
	if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
		request["ExcludeResourceIdsScope"] = v
	}
	if v, ok := d.GetOk("resource_group_ids_scope"); ok {
		request["ResourceGroupIdsScope"] = v
	}
	if v, ok := d.GetOk("exclude_resource_group_ids_scope"); ok {
		request["ExcludeResourceGroupIdsScope"] = v
	}
	if v, ok := d.GetOk("region_ids_scope"); ok {
		request["RegionIdsScope"] = v
	}
	if v, ok := d.GetOk("exclude_region_ids_scope"); ok {
		request["ExcludeRegionIdsScope"] = v
	}
	if v, ok := d.GetOk("tags_scope"); ok {
		for i, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			request[fmt.Sprintf("TagsScope.%d.TagKey", i+1)] = dataLoopTmp["tag_key"]
			request[fmt.Sprintf("TagsScope.%d.TagValue", i+1)] = dataLoopTmp["tag_value"]
		}
	}
	if v, ok := d.GetOk("exclude_tags_scope"); ok {
		for i, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			request[fmt.Sprintf("ExcludeTagsScope.%d.TagKey", i+1)] = dataLoopTmp["tag_key"]
			request[fmt.Sprintf("ExcludeTagsScope.%d.TagValue", i+1)] = dataLoopTmp["tag_value"]
		}
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

				if configRuleName, ok := configRulesArg["config_rule_name"]; ok {
					configRulesMap["ConfigRuleName"] = configRuleName
				}

				if description, ok := configRulesArg["description"]; ok {
					configRulesMap["Description"] = description
				}

				if riskLevel, ok := configRulesArg["risk_level"]; ok {
					configRulesMap["RiskLevel"] = riskLevel
				}

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
	d.Set("template_content", object["TemplateContent"])

	if scope, ok := object["Scope"].(map[string]interface{}); ok {
		d.Set("tag_key_scope", scope["TagKeyScope"])
		d.Set("tag_value_scope", scope["TagValueScope"])
		d.Set("resource_ids_scope", scope["ResourceIdsScope"])
		d.Set("exclude_resource_ids_scope", scope["ExcludeResourceIdsScope"])
		d.Set("resource_group_ids_scope", scope["ResourceGroupIdsScope"])
		d.Set("exclude_resource_group_ids_scope", scope["ExcludeResourceGroupIdsScope"])
		d.Set("region_ids_scope", scope["RegionIdsScope"])
		d.Set("exclude_region_ids_scope", scope["ExcludeRegionIdsScope"])

		if tagsScopeList, ok := scope["TagsScope"].([]interface{}); ok {
			tagsScopeMaps := make([]map[string]interface{}, 0)
			for _, tagsScope := range tagsScopeList {
				tagsScopeArg := tagsScope.(map[string]interface{})
				tagsScopeMap := map[string]interface{}{
					"tag_key":   tagsScopeArg["TagKey"],
					"tag_value": tagsScopeArg["TagValue"],
				}
				tagsScopeMaps = append(tagsScopeMaps, tagsScopeMap)
			}
			if err := d.Set("tags_scope", tagsScopeMaps); err != nil {
				return WrapError(err)
			}
		}

		if excludeTagsScopeList, ok := scope["ExcludeTagsScope"].([]interface{}); ok {
			excludeTagsScopeMaps := make([]map[string]interface{}, 0)
			for _, excludeTagsScope := range excludeTagsScopeList {
				excludeTagsScopeArg := excludeTagsScope.(map[string]interface{})
				excludeTagsScopeMap := map[string]interface{}{
					"tag_key":   excludeTagsScopeArg["TagKey"],
					"tag_value": excludeTagsScopeArg["TagValue"],
				}
				excludeTagsScopeMaps = append(excludeTagsScopeMaps, excludeTagsScopeMap)
			}
			if err := d.Set("exclude_tags_scope", excludeTagsScopeMaps); err != nil {
				return WrapError(err)
			}
		}
	}

	if _, ok := d.GetOk("config_rules"); ok {
		if configRulesList, ok := object["ConfigRules"]; ok {
			configRulesMaps := make([]map[string]interface{}, 0)
			for _, configRules := range configRulesList.([]interface{}) {
				configRulesArg := configRules.(map[string]interface{})
				configRulesMap := map[string]interface{}{}

				if managedRuleIdentifier, ok := configRulesArg["ManagedRuleIdentifier"]; ok {
					configRulesMap["managed_rule_identifier"] = managedRuleIdentifier
				}

				if configRuleName, ok := configRulesArg["ConfigRuleName"]; ok {
					configRulesMap["config_rule_name"] = configRuleName
				}

				if description, ok := configRulesArg["Description"]; ok {
					configRulesMap["description"] = description
				}

				if riskLevel, ok := configRulesArg["RiskLevel"]; ok {
					configRulesMap["risk_level"] = riskLevel
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

			if configRuleName, ok := configRulesArg["config_rule_name"]; ok {
				configRulesMap["ConfigRuleName"] = configRuleName
			}

			if description, ok := configRulesArg["description"]; ok {
				configRulesMap["Description"] = description
			}

			if riskLevel, ok := configRulesArg["risk_level"]; ok {
				configRulesMap["RiskLevel"] = riskLevel
			}

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

	if d.HasChange("tag_key_scope") {
		update = true
		request["TagKeyScope"] = d.Get("tag_key_scope")
	}
	if d.HasChange("tag_value_scope") {
		update = true
		request["TagValueScope"] = d.Get("tag_value_scope")
	}
	if d.HasChange("resource_ids_scope") {
		update = true
		request["ResourceIdsScope"] = d.Get("resource_ids_scope")
	}
	if d.HasChange("exclude_resource_ids_scope") {
		update = true
		request["ExcludeResourceIdsScope"] = d.Get("exclude_resource_ids_scope")
	}
	if d.HasChange("resource_group_ids_scope") {
		update = true
		request["ResourceGroupIdsScope"] = d.Get("resource_group_ids_scope")
	}
	if d.HasChange("exclude_resource_group_ids_scope") {
		update = true
		request["ExcludeResourceGroupIdsScope"] = d.Get("exclude_resource_group_ids_scope")
	}
	if d.HasChange("region_ids_scope") {
		update = true
		request["RegionIdsScope"] = d.Get("region_ids_scope")
	}
	if d.HasChange("exclude_region_ids_scope") {
		update = true
		request["ExcludeRegionIdsScope"] = d.Get("exclude_region_ids_scope")
	}
	if d.HasChange("tags_scope") {
		update = true
		for i, dataLoop := range d.Get("tags_scope").([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			request[fmt.Sprintf("TagsScope.%d.TagKey", i+1)] = dataLoopTmp["tag_key"]
			request[fmt.Sprintf("TagsScope.%d.TagValue", i+1)] = dataLoopTmp["tag_value"]
		}
	}
	if d.HasChange("exclude_tags_scope") {
		update = true
		for i, dataLoop := range d.Get("exclude_tags_scope").([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			request[fmt.Sprintf("ExcludeTagsScope.%d.TagKey", i+1)] = dataLoopTmp["tag_key"]
			request[fmt.Sprintf("ExcludeTagsScope.%d.TagValue", i+1)] = dataLoopTmp["tag_value"]
		}
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
		d.SetPartial("tag_key_scope")
		d.SetPartial("tag_value_scope")
		d.SetPartial("resource_ids_scope")
		d.SetPartial("exclude_resource_ids_scope")
		d.SetPartial("resource_group_ids_scope")
		d.SetPartial("exclude_resource_group_ids_scope")
		d.SetPartial("region_ids_scope")
		d.SetPartial("exclude_region_ids_scope")
		d.SetPartial("tags_scope")
		d.SetPartial("exclude_tags_scope")
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
