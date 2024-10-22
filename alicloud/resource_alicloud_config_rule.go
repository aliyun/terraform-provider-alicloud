// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudConfigRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigRuleCreate,
		Read:   resourceAlicloudConfigRuleRead,
		Update: resourceAlicloudConfigRuleUpdate,
		Delete: resourceAlicloudConfigRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"compliance": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compliance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"compliance_pack_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_rule_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_rule_trigger_types": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"config_rule_trigger_types", "source_detail_message_type"},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exclude_resource_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"input_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"maximum_execution_frequency": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"source_maximum_execution_frequency"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if fmt.Sprint(d.Get("config_rule_trigger_types")) == "ConfigurationItemChangeNotification" || fmt.Sprint(d.Get("source_detail_message_type")) == "ConfigurationItemChangeNotification" {
						return true
					}
					return false
				},
			},
			"modified_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_ids_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_types_scope": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"resource_types_scope", "scope_compliance_resource_types"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"risk_level": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_owner": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tag_key_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_value_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_detail_message_type": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'source_detail_message_type' has been deprecated from provider version 1.124.1. New field 'config_rule_trigger_types' instead.",
			},
			"source_maximum_execution_frequency": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if fmt.Sprint(d.Get("config_rule_trigger_types")) == "ConfigurationItemChangeNotification" || fmt.Sprint(d.Get("source_detail_message_type")) == "ConfigurationItemChangeNotification" {
						return true
					}
					return false
				},
				Deprecated: "Field 'source_maximum_execution_frequency' has been deprecated from provider version 1.124.1. New field 'maximum_execution_frequency' instead.",
			},
			"scope_compliance_resource_types": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'scope_compliance_resource_types' has been deprecated from provider version 1.124.1. New field 'resource_types_scope' instead.",
			},
		},
	}
}

func resourceAlicloudConfigRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateConfigRule"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("source_owner"); ok {
		request["SourceOwner"] = v
	}
	if v, ok := d.GetOk("source_identifier"); ok {
		request["SourceIdentifier"] = v
	}
	if v, ok := d.GetOk("config_rule_trigger_types"); ok {
		request["ConfigRuleTriggerTypes"] = v
	} else if v, ok := d.GetOk("source_detail_message_type"); ok {
		request["ConfigRuleTriggerTypes"] = v
	}

	if v, ok := d.GetOk("risk_level"); ok {
		request["RiskLevel"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["ConfigRuleName"] = v
	}
	if v, ok := d.GetOk("tag_key_scope"); ok {
		request["TagKeyScope"] = v
	}
	if v, ok := d.GetOk("tag_value_scope"); ok {
		request["TagValueScope"] = v
	}
	if v, ok := d.GetOk("region_ids_scope"); ok {
		request["RegionIdsScope"] = v
	}
	if v, ok := d.GetOk("resource_group_ids_scope"); ok {
		request["ResourceGroupIdsScope"] = v
	}
	if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
		request["ExcludeResourceIdsScope"] = v
	}
	if v, ok := d.GetOk("input_parameters"); ok {
		request["InputParameters"] = convertMapToJsonStringIgnoreError(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("resource_types_scope"); ok {
		jsonPathResult12, err := jsonpath.Get("$", v)
		if err != nil {
			return WrapError(err)
		}
		request["ResourceTypesScope"] = convertListToCommaSeparate(jsonPathResult12.([]interface{}))
	} else if v, ok := d.GetOk("scope_compliance_resource_types"); ok {
		request["ResourceTypesScope"] = convertListToCommaSeparate(v.([]interface{}))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ConfigRuleId"]))

	configServiceV2 := ConfigServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ACTIVE", "EVALUATING"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, configServiceV2.ConfigRuleStateRefreshFunc(d.Id(), "ConfigRuleState", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudConfigRuleUpdate(d, meta)
}

func resourceAlicloudConfigRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configServiceV2 := ConfigServiceV2{client}

	objectRaw, err := configServiceV2.DescribeConfigRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_rule .DescribeConfigRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("account_id", objectRaw["AccountId"])
	d.Set("config_rule_arn", objectRaw["ConfigRuleArn"])
	d.Set("create_time", objectRaw["CreateTimestamp"])
	d.Set("description", objectRaw["Description"])
	d.Set("exclude_resource_ids_scope", objectRaw["ExcludeResourceIdsScope"])
	d.Set("input_parameters", objectRaw["InputParameters"])
	d.Set("maximum_execution_frequency", objectRaw["MaximumExecutionFrequency"])
	d.Set("modified_timestamp", objectRaw["ModifiedTimestamp"])
	d.Set("region_ids_scope", objectRaw["RegionIdsScope"])
	d.Set("resource_group_ids_scope", objectRaw["ResourceGroupIdsScope"])
	d.Set("risk_level", objectRaw["RiskLevel"])
	d.Set("rule_name", objectRaw["ConfigRuleName"])
	d.Set("status", objectRaw["ConfigRuleState"])
	d.Set("tag_key_scope", objectRaw["TagKeyScope"])
	d.Set("tag_value_scope", objectRaw["TagValueScope"])
	d.Set("config_rule_id", objectRaw["ConfigRuleId"])
	createBy2RawObj, _ := jsonpath.Get("$.CreateBy", objectRaw)
	createBy2Raw := make(map[string]interface{})
	if createBy2RawObj != nil {
		createBy2Raw = createBy2RawObj.(map[string]interface{})
	}
	d.Set("compliance_pack_id", createBy2Raw["CompliancePackId"])
	source2RawObj, _ := jsonpath.Get("$.Source", objectRaw)
	source2Raw := make(map[string]interface{})
	if source2RawObj != nil {
		source2Raw = source2RawObj.(map[string]interface{})
	}
	d.Set("source_identifier", source2Raw["Identifier"])
	d.Set("source_owner", source2Raw["Owner"])
	sourceDetails2RawObj, _ := jsonpath.Get("$.Source.SourceDetails[*]", objectRaw)
	sourceDetails2Raw := make([]interface{}, 0)
	if sourceDetails2RawObj != nil {
		sourceDetails2Raw = sourceDetails2RawObj.([]interface{})
	}

	sourceDetailsChild2Raw := make(map[string]interface{})
	if len(sourceDetails2Raw) > 0 {
		sourceDetailsChild2Raw = sourceDetails2Raw[0].(map[string]interface{})
	}
	d.Set("config_rule_trigger_types", sourceDetailsChild2Raw["MessageType"])
	d.Set("event_source", sourceDetailsChild2Raw["EventSource"])
	complianceMaps := make([]map[string]interface{}, 0)
	complianceMap := make(map[string]interface{})
	compliance2Raw := make(map[string]interface{})
	if objectRaw["Compliance"] != nil {
		compliance2Raw = objectRaw["Compliance"].(map[string]interface{})
	}

	complianceMap["compliance_type"] = compliance2Raw["ComplianceType"]
	complianceMap["count"] = compliance2Raw["Count"]
	complianceMaps = append(complianceMaps, complianceMap)
	d.Set("compliance", complianceMaps)
	complianceResourceTypes2Raw, _ := jsonpath.Get("$.Scope.ComplianceResourceTypes", objectRaw)
	d.Set("resource_types_scope", complianceResourceTypes2Raw)

	d.Set("source_detail_message_type", d.Get("config_rule_trigger_types"))
	d.Set("source_maximum_execution_frequency", d.Get("maximum_execution_frequency"))
	d.Set("scope_compliance_resource_types", d.Get("resource_types_scope"))
	return nil
}

func resourceAlicloudConfigRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	update := false
	d.Partial(true)
	update = false
	action := "UpdateConfigRule"

	request = make(map[string]interface{})

	request["ConfigRuleId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChanges("maximum_execution_frequency", "source_maximum_execution_frequency") {
		update = true
		if v, ok := d.GetOk("maximum_execution_frequency"); ok {
			request["MaximumExecutionFrequency"] = v
		} else if v, ok := d.GetOk("source_maximum_execution_frequency"); ok {
			request["MaximumExecutionFrequency"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("tag_key_scope") {
		update = true
		if v, ok := d.GetOk("tag_key_scope"); ok {
			request["TagKeyScope"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("tag_value_scope") {
		update = true
		if v, ok := d.GetOk("tag_value_scope"); ok {
			request["TagValueScope"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("region_ids_scope") {
		update = true
		if v, ok := d.GetOk("region_ids_scope"); ok {
			request["RegionIdsScope"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("resource_group_ids_scope") {
		update = true
		if v, ok := d.GetOk("resource_group_ids_scope"); ok {
			request["ResourceGroupIdsScope"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("exclude_resource_ids_scope") {
		update = true
		if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
			request["ExcludeResourceIdsScope"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("risk_level") {
		update = true
		if v, ok := d.GetOk("risk_level"); ok {
			request["RiskLevel"] = v
		}
	}
	if !d.IsNewResource() && d.HasChanges("config_rule_trigger_types", "source_detail_message_type") {
		update = true
		if v, ok := d.GetOk("config_rule_trigger_types"); ok {
			request["ConfigRuleTriggerTypes"] = v
		} else if v, ok := d.GetOk("source_detail_message_type"); ok {
			request["ConfigRuleTriggerTypes"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("input_parameters") {
		update = true
		if v, ok := d.GetOk("input_parameters"); ok {
			request["InputParameters"] = convertMapToJsonStringIgnoreError(v.(map[string]interface{}))
		}
	}
	if !d.IsNewResource() && d.HasChanges("resource_types_scope", "scope_compliance_resource_types") {
		update = true
		if v, ok := d.GetOk("resource_types_scope"); ok {
			jsonPathResult10, err := jsonpath.Get("$", v)
			if err != nil {
				return WrapError(err)
			}
			request["ResourceTypesScope"] = convertListToCommaSeparate(jsonPathResult10.([]interface{}))
		} else if v, ok := d.GetOk("scope_compliance_resource_types"); ok {
			request["ResourceTypesScope"] = convertListToCommaSeparate(v.([]interface{}))
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		{
			configServiceV2 := ConfigServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"ACTIVE", "EVALUATING"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, configServiceV2.ConfigRuleStateRefreshFunc(d.Id(), "ConfigRuleState", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("description")
		d.SetPartial("maximum_execution_frequency")
		d.SetPartial("tag_key_scope")
		d.SetPartial("tag_value_scope")
		d.SetPartial("region_ids_scope")
		d.SetPartial("resource_group_ids_scope")
		d.SetPartial("exclude_resource_ids_scope")
		d.SetPartial("risk_level")
		d.SetPartial("config_rule_trigger_types")
		d.SetPartial("input_parameters")
		d.SetPartial("resource_types_scope")
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		configServiceV2 := ConfigServiceV2{client}
		object, err := configServiceV2.DescribeConfigRule(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["ConfigRuleState"].(string) != target {
			if target == "INACTIVE" {
				action = "StopConfigRules"
				request = make(map[string]interface{})

				request["ConfigRuleIds"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Config", "2019-01-08", action, nil, request, false)

					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == "ACTIVE" {
				action = "ActiveConfigRules"
				request = make(map[string]interface{})

				request["ConfigRuleIds"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Config", "2019-01-08", action, nil, request, false)

					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	d.Partial(false)
	return resourceAlicloudConfigRuleRead(d, meta)
}

func resourceAlicloudConfigRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteConfigRules"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["ConfigRuleIds"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2019-01-08", action, nil, request, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ConfigRuleNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	configServiceV2 := ConfigServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, configServiceV2.ConfigRuleStateRefreshFunc(d.Id(), "ConfigRuleState", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
