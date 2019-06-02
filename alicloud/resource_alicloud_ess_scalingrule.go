package alicloud

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEssScalingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingRuleCreate,
		Read:   resourceAliyunEssScalingRuleRead,
		Update: resourceAliyunEssScalingRuleUpdate,
		Delete: resourceAliyunEssScalingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"adjustment_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{string(QuantityChangeInCapacity),
					string(PercentChangeInCapacity), string(TotalCapacity)}),
			},
			"adjustment_value": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"scaling_rule_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"ari": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cooldown": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, 86400),
			},
		},
	}
}

func resourceAliyunEssScalingRuleCreate(d *schema.ResourceData, meta interface{}) error {

	request, err := buildAlicloudEssScalingRuleArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateScalingRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scalingrule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ess.CreateScalingRuleResponse)
	d.SetId(response.ScalingRuleId)

	return resourceAliyunEssScalingRuleRead(d, meta)
}

func resourceAliyunEssScalingRuleRead(d *schema.ResourceData, meta interface{}) error {

	//Compatible with older versions id
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		parts, _ := ParseResourceId(d.Id(), 2)
		d.SetId(parts[1])
	}

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	object, err := essService.DescribeEssScalingRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("ari", object.ScalingRuleAri)
	d.Set("adjustment_type", object.AdjustmentType)
	d.Set("adjustment_value", object.AdjustmentValue)
	d.Set("scaling_rule_name", object.ScalingRuleName)
	d.Set("cooldown", object.Cooldown)

	return nil
}

func resourceAliyunEssScalingRuleDelete(d *schema.ResourceData, meta interface{}) error {

	//Compatible with older versions id
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		parts, _ := ParseResourceId(d.Id(), 2)
		d.SetId(parts[1])
	}

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	request := ess.CreateDeleteScalingRuleRequest()
	request.ScalingRuleId = d.Id()

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingRule(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidScalingRuleIdNotFound}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return WrapError(essService.WaitForEssScalingRule(d.Id(), Deleted, DefaultTimeout))
}

func resourceAliyunEssScalingRuleUpdate(d *schema.ResourceData, meta interface{}) error {

	//Compatible with older versions id
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		parts, _ := ParseResourceId(d.Id(), 2)
		d.SetId(parts[1])
	}

	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateModifyScalingRuleRequest()
	request.ScalingRuleId = d.Id()

	if d.HasChange("adjustment_type") {
		request.AdjustmentType = d.Get("adjustment_type").(string)
	}

	if d.HasChange("adjustment_value") {
		request.AdjustmentValue = requests.NewInteger(d.Get("adjustment_value").(int))
	}

	if d.HasChange("scaling_rule_name") {
		request.ScalingRuleName = d.Get("scaling_rule_name").(string)
	}

	if d.HasChange("cooldown") {
		request.Cooldown = requests.NewInteger(d.Get("cooldown").(int))
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return resourceAliyunEssScalingRuleRead(d, meta)
}

func buildAlicloudEssScalingRuleArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingRuleRequest, error) {
	request := ess.CreateCreateScalingRuleRequest()
	request.ScalingGroupId = d.Get("scaling_group_id").(string)
	request.AdjustmentType = d.Get("adjustment_type").(string)
	request.AdjustmentValue = requests.NewInteger(d.Get("adjustment_value").(int))

	if v, ok := d.GetOk("scaling_rule_name"); ok && v.(string) != "" {
		request.ScalingRuleName = v.(string)
	}

	if v, ok := d.GetOk("cooldown"); ok && v.(int) != 0 {
		request.Cooldown = requests.NewInteger(v.(int))
	}

	return request, nil
}
