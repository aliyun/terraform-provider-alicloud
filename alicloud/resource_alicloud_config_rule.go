package alicloud

import (
	"fmt"
	"log"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudConfigRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigRuleCreate,
		Read:   resourceAlicloudConfigRuleRead,
		Update: resourceAlicloudConfigRuleUpdate,
		Delete: resourceAlicloudConfigRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"input_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"member_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"multi_account": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"risk_level": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope_compliance_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope_compliance_resource_types": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_detail_message_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ConfigurationItemChangeNotification", "OversizedConfigurationItemChangeNotification", "ScheduledNotification"}, false),
			},
			"source_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_maximum_execution_frequency": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"One_Hour", "Six_Hours", "Three_Hours", "Twelve_Hours", "TwentyFour_Hours"}, false),
			},
			"source_owner": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALIYUN", "CUSTOM_FC"}, false),
			},
		},
	}
}

func resourceAlicloudConfigRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "PutConfigRule"
	request := make(map[string]interface{})
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("input_parameters"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_rule", action, AlibabaCloudSdkGoERROR)
		}
		request["InputParameters"] = respJson
	}

	if v, ok := d.GetOk("member_id"); ok {
		request["MemberId"] = v
	}

	if v, ok := d.GetOkExists("multi_account"); ok {
		request["MultiAccount"] = v
	}

	request["RiskLevel"] = d.Get("risk_level")
	request["ConfigRuleName"] = d.Get("rule_name")
	if v, ok := d.GetOk("scope_compliance_resource_id"); ok {
		request["ScopeComplianceResourceId"] = v
	}

	request["ScopeComplianceResourceTypes"] = convertListToJsonString(d.Get("scope_compliance_resource_types").(*schema.Set).List())
	request["SourceDetailMessageType"] = d.Get("source_detail_message_type")
	request["SourceIdentifier"] = d.Get("source_identifier")
	if v, ok := d.GetOk("source_maximum_execution_frequency"); ok {
		request["SourceMaximumExecutionFrequency"] = v
	}

	request["SourceOwner"] = d.Get("source_owner")
	request["ClientToken"] = buildClientToken("PutConfigRule")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_rule", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["ConfigRuleId"]))

	return resourceAlicloudConfigRuleRead(d, meta)
}
func resourceAlicloudConfigRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_rule configService.DescribeConfigRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("input_parameters", object["InputParameters"])
	d.Set("risk_level", formatInt(object["RiskLevel"]))
	d.Set("rule_name", object["ConfigRuleName"])
	d.Set("scope_compliance_resource_id", object["Scope"].(map[string]interface{})["ComplianceResourceId"])
	d.Set("scope_compliance_resource_types", object["Scope"].(map[string]interface{})["ComplianceResourceTypes"])
	d.Set("source_identifier", object["Source"].(map[string]interface{})["Identifier"])
	d.Set("source_maximum_execution_frequency", object["MaximumExecutionFrequency"])
	d.Set("source_owner", object["Source"].(map[string]interface{})["Owner"])
	if v := object["Source"].(map[string]interface{})["SourceDetails"].([]interface{}); len(v) > 0 {
		d.Set("source_detail_message_type", v[0].(map[string]interface{})["MessageType"])
	}
	return nil
}
func resourceAlicloudConfigRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ConfigRuleId": d.Id(),
	}
	if d.HasChange("risk_level") {
		update = true
	}
	request["RiskLevel"] = d.Get("risk_level")
	request["ConfigRuleName"] = d.Get("rule_name")
	if d.HasChange("scope_compliance_resource_types") {
		update = true
	}
	request["ScopeComplianceResourceTypes"] = convertListToJsonString(d.Get("scope_compliance_resource_types").(*schema.Set).List())
	if d.HasChange("source_detail_message_type") {
		update = true
	}
	request["SourceDetailMessageType"] = d.Get("source_detail_message_type")
	request["SourceIdentifier"] = d.Get("source_identifier")
	request["SourceOwner"] = d.Get("source_owner")
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("input_parameters") {
		update = true
		respJson, err := convertMaptoJsonString(d.Get("input_parameters").(map[string]interface{}))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_rule", "PutConfigRule", AlibabaCloudSdkGoERROR)
		}
		request["InputParameters"] = respJson
	}
	if d.HasChange("scope_compliance_resource_id") {
		update = true
		request["ScopeComplianceResourceId"] = d.Get("scope_compliance_resource_id")
	}
	if d.HasChange("source_maximum_execution_frequency") {
		update = true
		request["SourceMaximumExecutionFrequency"] = d.Get("source_maximum_execution_frequency")
	}
	if update {
		if _, ok := d.GetOk("member_id"); ok {
			request["MemberId"] = d.Get("member_id")
		}
		if _, ok := d.GetOkExists("multi_account"); ok {
			request["MultiAccount"] = d.Get("multi_account")
		}
		action := "PutConfigRule"
		conn, err := client.NewConfigClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("PutConfigRule")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudConfigRuleRead(d, meta)
}
func resourceAlicloudConfigRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteConfigRules"
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ConfigRuleIds": d.Id(),
	}

	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"AccountNotExisted", "ConfigRuleNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
