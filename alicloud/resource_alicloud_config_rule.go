package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"config_rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"input_parameters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"risk_level": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
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
				ValidateFunc: validation.StringInSlice([]string{"ConfigurationItemChangeNotification", "ConfigurationSnapshotDeliveryCompleted", "OversizedConfigurationItemChangeNotification", "Schedule"}, false),
			},
			"source_identifier": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_maximum_execution_frequency": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"One_Hour", "Six_Hours", "Three_Hours", "Twelve_Hours", "TwentyFour_Hours"}, false),
			},
			"source_owner": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALIYUN", "CUSTOM_FC"}, false),
			},
		},
	}
}

func resourceAlicloudConfigRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := config.CreatePutConfigRuleRequest()
	request.ConfigRuleName = d.Get("config_rule_name").(string)
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("input_parameters"); ok {
		request.InputParameters = v.(string)
	}
	request.RiskLevel = requests.NewInteger(d.Get("risk_level").(int))
	if v, ok := d.GetOk("scope_compliance_resource_id"); ok {
		request.ScopeComplianceResourceId = v.(string)
	}
	request.ScopeComplianceResourceTypes = convertListToJsonString(d.Get("scope_compliance_resource_types").(*schema.Set).List())
	request.SourceDetailMessageType = d.Get("source_detail_message_type").(string)
	request.SourceIdentifier = d.Get("source_identifier").(string)
	if v, ok := d.GetOk("source_maximum_execution_frequency"); ok {
		request.SourceMaximumExecutionFrequency = v.(string)
	}
	request.SourceOwner = d.Get("source_owner").(string)

	request.ClientToken = buildClientToken(request.GetActionName())
	raw, err := client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
		return configClient.PutConfigRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_rule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*config.PutConfigRuleResponse)
	d.SetId(fmt.Sprintf("%v", response.ConfigRuleId))

	return resourceAlicloudConfigRuleRead(d, meta)
}
func resourceAlicloudConfigRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config_rule_name", object.ConfigRuleName)
	d.Set("description", object.Description)
	d.Set("risk_level", object.RiskLevel)
	d.Set("source_identifier", object.Source.Identifier)
	d.Set("source_owner", object.Source.Owner)
	return nil
}
func resourceAlicloudConfigRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := config.CreatePutConfigRuleRequest()
	request.ConfigRuleId = d.Id()
	if d.HasChange("config_rule_name") {
		update = true
	}
	request.ConfigRuleName = d.Get("config_rule_name").(string)
	if d.HasChange("risk_level") {
		update = true
	}
	request.RiskLevel = requests.NewInteger(d.Get("risk_level").(int))
	if d.HasChange("scope_compliance_resource_types") {
		update = true
	}
	request.ScopeComplianceResourceTypes = convertListToJsonString(d.Get("scope_compliance_resource_types").(*schema.Set).List())
	if d.HasChange("source_detail_message_type") {
		update = true
	}
	request.SourceDetailMessageType = d.Get("source_detail_message_type").(string)
	if d.HasChange("source_identifier") {
		update = true
	}
	request.SourceIdentifier = d.Get("source_identifier").(string)
	if d.HasChange("source_owner") {
		update = true
	}
	request.SourceOwner = d.Get("source_owner").(string)
	if d.HasChange("description") {
		update = true
		request.Description = d.Get("description").(string)
	}
	if d.HasChange("input_parameters") {
		update = true
		request.InputParameters = d.Get("input_parameters").(string)
	}
	if d.HasChange("scope_compliance_resource_id") {
		update = true
		request.ScopeComplianceResourceId = d.Get("scope_compliance_resource_id").(string)
	}
	if d.HasChange("source_maximum_execution_frequency") {
		update = true
		request.SourceMaximumExecutionFrequency = d.Get("source_maximum_execution_frequency").(string)
	}
	if update {
		request.ClientToken = buildClientToken(request.GetActionName())
		raw, err := client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
			return configClient.PutConfigRule(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudConfigRuleRead(d, meta)
}
func resourceAlicloudConfigRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := config.CreateDeleteConfigRulesRequest()
	request.ConfigRuleIds = d.Id()
	raw, err := client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
		return configClient.DeleteConfigRules(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"AccountNotExisted", "ConfigRuleNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
