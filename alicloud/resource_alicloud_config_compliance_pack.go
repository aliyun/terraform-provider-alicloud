package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudConfigCompliancePack() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudConfigCompliancePackCreate,
		Read:   resourceAliCloudConfigCompliancePackRead,
		Update: resourceAliCloudConfigCompliancePackUpdate,
		Delete: resourceAliCloudConfigCompliancePackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"compliance_pack_name": {
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
				Computed:      true,
				ConflictsWith: []string{"config_rules"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_rule_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
						"config_rule_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"risk_level": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: IntInSlice([]int{1, 2, 3}),
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_content": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func resourceAliCloudConfigCompliancePackCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	var err error
	action := "CreateCompliancePack"
	request := make(map[string]interface{})

	request["ClientToken"] = buildClientToken("CreateCompliancePack")
	request["CompliancePackName"] = d.Get("compliance_pack_name")
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

				if configRuleName := configRulesArg["config_rule_name"]; configRuleName.(string) != "" {
					configRulesMap["ConfigRuleName"] = configRuleName
				}

				if description := configRulesArg["description"]; description.(string) != "" {
					configRulesMap["Description"] = description
				}

				if riskLevel := configRulesArg["risk_level"]; riskLevel.(int) != 0 {
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
		setTagsScopeParams(request, "TagsScope", v.([]interface{}))
	}
	if v, ok := d.GetOk("exclude_tags_scope"); ok {
		setTagsScopeParams(request, "ExcludeTagsScope", v.([]interface{}))
	}

	wait := incrementalWait(3*time.Second, 30*time.Second)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_compliance_pack", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["CompliancePackId"]))

	stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, configService.ConfigCompliancePackStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudConfigCompliancePackRead(d, meta)
}

func resourceAliCloudConfigCompliancePackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}

	object, err := configService.DescribeConfigCompliancePack(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_compliance_pack configService.DescribeConfigCompliancePack Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("compliance_pack_name", object["CompliancePackName"])
	d.Set("description", object["Description"])
	d.Set("risk_level", formatInt(object["RiskLevel"]))
	d.Set("compliance_pack_template_id", object["CompliancePackTemplateId"])
	d.Set("status", object["Status"])
	d.Set("template_content", object["TemplateContent"])

	if scope, ok := object["Scope"]; ok {
		scopeMap := scope.(map[string]interface{})
		d.Set("tag_key_scope", scopeMap["TagKeyScope"])
		d.Set("tag_value_scope", scopeMap["TagValueScope"])
		d.Set("resource_ids_scope", scopeMap["ResourceIdsScope"])
		d.Set("exclude_resource_ids_scope", scopeMap["ExcludeResourceIdsScope"])
		d.Set("resource_group_ids_scope", scopeMap["ResourceGroupIdsScope"])
		d.Set("exclude_resource_group_ids_scope", scopeMap["ExcludeResourceGroupIdsScope"])
		d.Set("region_ids_scope", scopeMap["RegionIdsScope"])
		d.Set("exclude_region_ids_scope", scopeMap["ExcludeRegionIdsScope"])

		if tagsScopeList, ok := scopeMap["TagsScope"]; ok {
			tagsScopeMaps := make([]map[string]interface{}, 0)
			for _, ts := range tagsScopeList.([]interface{}) {
				tsArg := ts.(map[string]interface{})
				tsMap := map[string]interface{}{}
				if tagKey, ok := tsArg["TagKey"]; ok {
					tsMap["tag_key"] = tagKey
				}
				if tagValue, ok := tsArg["TagValue"]; ok {
					tsMap["tag_value"] = tagValue
				}
				tagsScopeMaps = append(tagsScopeMaps, tsMap)
			}
			d.Set("tags_scope", tagsScopeMaps)
		}

		if excludeTagsScopeList, ok := scopeMap["ExcludeTagsScope"]; ok {
			excludeTagsScopeMaps := make([]map[string]interface{}, 0)
			for _, ets := range excludeTagsScopeList.([]interface{}) {
				etsArg := ets.(map[string]interface{})
				etsMap := map[string]interface{}{}
				if tagKey, ok := etsArg["TagKey"]; ok {
					etsMap["tag_key"] = tagKey
				}
				if tagValue, ok := etsArg["TagValue"]; ok {
					etsMap["tag_value"] = tagValue
				}
				excludeTagsScopeMaps = append(excludeTagsScopeMaps, etsMap)
			}
			d.Set("exclude_tags_scope", excludeTagsScopeMaps)
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
					configRulesMap["risk_level"] = formatInt(riskLevel)
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

func resourceAliCloudConfigCompliancePackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"ClientToken":      buildClientToken("UpdateCompliancePack"),
		"CompliancePackId": d.Id(),
	}

	if d.HasChange("compliance_pack_name") {
		update = true
	}
	request["CompliancePackName"] = d.Get("compliance_pack_name")

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

			if configRuleName := configRulesArg["config_rule_name"]; configRuleName.(string) != "" {
				configRulesMap["ConfigRuleName"] = configRuleName
			}

			if description := configRulesArg["description"]; description.(string) != "" {
				configRulesMap["Description"] = description
			}

			if riskLevel := configRulesArg["risk_level"]; riskLevel.(int) != 0 {
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
	}
	if v, ok := d.GetOk("tag_key_scope"); ok {
		request["TagKeyScope"] = v
	}
	if d.HasChange("tag_value_scope") {
		update = true
	}
	if v, ok := d.GetOk("tag_value_scope"); ok {
		request["TagValueScope"] = v
	}
	if d.HasChange("resource_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("resource_ids_scope"); ok {
		request["ResourceIdsScope"] = v
	}
	if d.HasChange("exclude_resource_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
		request["ExcludeResourceIdsScope"] = v
	}
	if d.HasChange("resource_group_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_ids_scope"); ok {
		request["ResourceGroupIdsScope"] = v
	}
	if d.HasChange("exclude_resource_group_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("exclude_resource_group_ids_scope"); ok {
		request["ExcludeResourceGroupIdsScope"] = v
	}
	if d.HasChange("region_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("region_ids_scope"); ok {
		request["RegionIdsScope"] = v
	}
	if d.HasChange("exclude_region_ids_scope") {
		update = true
	}
	if v, ok := d.GetOk("exclude_region_ids_scope"); ok {
		request["ExcludeRegionIdsScope"] = v
	}
	if d.HasChange("tags_scope") {
		update = true
	}
	if v, ok := d.GetOk("tags_scope"); ok {
		setTagsScopeParams(request, "TagsScope", v.([]interface{}))
	}
	if d.HasChange("exclude_tags_scope") {
		update = true
	}
	if v, ok := d.GetOk("exclude_tags_scope"); ok {
		setTagsScopeParams(request, "ExcludeTagsScope", v.([]interface{}))
	}

	if update {
		action := "UpdateCompliancePack"
		wait := incrementalWait(3*time.Second, 30*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"CompliancePackAlreadyPending"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, configService.ConfigCompliancePackStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("compliance_pack_name")
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
			action := "DetachConfigRuleToCompliancePack"

			detachConfigRuleReq := map[string]interface{}{
				"CompliancePackId": d.Id(),
			}

			configRuleIds := make([]interface{}, 0)
			for _, configRuleIdsList := range remove {
				configRuleIdsArg := configRuleIdsList.(map[string]interface{})

				if configRuleId, ok := configRuleIdsArg["config_rule_id"]; ok {
					configRuleIds = append(configRuleIds, configRuleId)
				}
			}

			detachConfigRuleReq["ConfigRuleIds"] = convertListToCommaSeparate(configRuleIds)

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Config", "2020-09-07", action, nil, detachConfigRuleReq, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, detachConfigRuleReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		if len(create) > 0 {
			action := "AttachConfigRuleToCompliancePack"

			attachConfigRuleReq := map[string]interface{}{
				"CompliancePackId": d.Id(),
			}

			configRuleIds := make([]interface{}, 0)
			for _, configRuleIdsList := range create {
				configRuleIdsArg := configRuleIdsList.(map[string]interface{})

				if configRuleId, ok := configRuleIdsArg["config_rule_id"]; ok {
					configRuleIds = append(configRuleIds, configRuleId)
				}
			}

			attachConfigRuleReq["ConfigRuleIds"] = convertListToCommaSeparate(configRuleIds)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Config", "2020-09-07", action, nil, attachConfigRuleReq, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, attachConfigRuleReq)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		d.SetPartial("config_rule_ids")
	}

	d.Partial(false)

	return resourceAliCloudConfigCompliancePackRead(d, meta)
}

func resourceAliCloudConfigCompliancePackDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCompliancePacks"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ClientToken":       buildClientToken("DeleteCompliancePacks"),
		"CompliancePackIds": d.Id(),
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// setTagsScopeParams sets tags_scope/exclude_tags_scope parameters in POP flat
// format (TagsScope.N.TagKey / TagsScope.N.TagValue) as required by the
// CreateCompliancePack/UpdateCompliancePack API. The API rejects JSON array
// serialization with "InvalidTagsScope: Flat format is required."
func setTagsScopeParams(request map[string]interface{}, key string, tagsScopeList []interface{}) {
	for i, ts := range tagsScopeList {
		tsArg, ok := ts.(map[string]interface{})
		if !ok {
			continue
		}
		if tagKey, ok := tsArg["tag_key"]; ok {
			request[fmt.Sprintf("%s.%d.TagKey", key, i+1)] = tagKey
		}
		if tagValue, ok := tsArg["tag_value"]; ok {
			request[fmt.Sprintf("%s.%d.TagValue", key, i+1)] = tagValue
		}
	}
}
